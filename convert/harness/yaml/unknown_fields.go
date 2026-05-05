// Copyright 2022 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package yaml

import (
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// collectUnknownFields compares the raw parsed JSON tree against the typed
// *Config value and returns a sorted list of JSON paths for map keys that
// have no matching json tag in the schema. Parsing never fails from unknown
// fields — they are silently ignored by encoding/json — so this walker is
// purely observational.
//
// Walk strategy: descend the raw tree and the typed reflect.Value in parallel.
// For each struct field with a json tag, look up the key in the raw map and
// recurse into the matched child. Raw keys that have no match are reported.
// For interface{} fields (Step.Spec, Stage.Spec, Runtime.Spec, Infrastructure.Spec,
// Volume.Spec), the concrete type has been assigned by the custom
// UnmarshalJSON based on a "type" discriminator, so we recurse using the
// runtime type of the parsed value.
//
// Types intentionally NOT descended into:
//   - Variable.Value / Variable.Default (user-defined arbitrary values)
//   - any map[string]... with non-struct element types (all keys user-defined)
func collectUnknownFields(cfg *Config, raw interface{}) []string {
	if cfg == nil || raw == nil {
		return nil
	}
	w := &unknownWalker{seen: make(map[string]struct{})}
	w.walk("$", reflect.ValueOf(cfg), raw)
	if len(w.seen) == 0 {
		return nil
	}
	out := make([]string, 0, len(w.seen))
	for p := range w.seen {
		out = append(out, p)
	}
	sort.Strings(out)
	return out
}

type unknownWalker struct {
	seen map[string]struct{}
}

func (w *unknownWalker) add(path string) {
	w.seen[path] = struct{}{}
}

// walk descends the raw value alongside the typed reflect.Value.
// path is the JSON path accumulated so far (dot+bracket form, e.g. "$.pipeline.stages[0].stage").
func (w *unknownWalker) walk(path string, typed reflect.Value, raw interface{}) {
	if raw == nil {
		return
	}

	// Unwrap typed pointer/interface layers.
	for typed.IsValid() && (typed.Kind() == reflect.Ptr || typed.Kind() == reflect.Interface) {
		if typed.IsNil() {
			// No typed value to compare against. This happens for interface{}
			// Spec fields with a nil/unknown type, or fields omitted in the
			// struct. Without a schema we cannot classify children as
			// known/unknown, so stop.
			return
		}
		typed = typed.Elem()
	}
	if !typed.IsValid() {
		return
	}

	// flexible.Field[T]'s internal Value field carries either a string
	// (expression — treat as leaf) or a concrete T (descend).
	if isFlexibleField(typed.Type()) {
		inner := typed.FieldByName("Value")
		w.walk(path, inner, raw)
		return
	}

	switch typed.Kind() {
	case reflect.Struct:
		rawMap, ok := raw.(map[string]interface{})
		if !ok {
			return
		}
		tags := structTags(typed.Type())
		for key, rawChild := range rawMap {
			fieldPath, known := tags[key]
			if !known {
				w.add(path + "." + key)
				continue
			}
			childTyped := typed.FieldByIndex(fieldPath)
			w.walk(path+"."+key, childTyped, rawChild)
		}
	case reflect.Slice, reflect.Array:
		rawSlice, ok := raw.([]interface{})
		if !ok {
			return
		}
		for i, rawItem := range rawSlice {
			childPath := path + "[" + strconv.Itoa(i) + "]"
			var childTyped reflect.Value
			if i < typed.Len() {
				childTyped = typed.Index(i)
			}
			w.walk(childPath, childTyped, rawItem)
		}
	case reflect.Map:
		// Typed maps with struct-shaped values (e.g. map[string]*Variable)
		// carry schemas on the values but not the keys. Descend into each
		// value. Maps with primitive/interface elements (map[string]string,
		// map[string]interface{}) have no schema to compare — skip.
		elemKind := typed.Type().Elem().Kind()
		if elemKind != reflect.Struct && elemKind != reflect.Ptr && elemKind != reflect.Interface {
			return
		}
		rawMap, ok := raw.(map[string]interface{})
		if !ok {
			return
		}
		for key, rawChild := range rawMap {
			childTyped := typed.MapIndex(reflect.ValueOf(key))
			if !childTyped.IsValid() {
				continue
			}
			// Map values from reflect are not addressable; copy so walk()
			// can unwrap interface/pointer layers.
			copyVal := reflect.New(childTyped.Type()).Elem()
			copyVal.Set(childTyped)
			w.walk(path+"."+key, copyVal, rawChild)
		}
	}
}

// isFlexibleField returns true when t is github.com/drone/go-convert/internal/flexible.Field[T].
// The type is parameterized, so match on package + name prefix.
func isFlexibleField(t reflect.Type) bool {
	if t.Kind() != reflect.Struct {
		return false
	}
	if t.PkgPath() != "github.com/drone/go-convert/internal/flexible" {
		return false
	}
	return strings.HasPrefix(t.Name(), "Field[")
}

// structTagCache memoizes the tag-to-field-index-path map for each struct type.
var structTagCache sync.Map // reflect.Type -> map[string][]int

// structTags returns a map from json tag name to the reflect field index path
// (suitable for Value.FieldByIndex). Embedded/anonymous fields are traversed
// so promoted tags resolve to their leaf field.
func structTags(t reflect.Type) map[string][]int {
	if cached, ok := structTagCache.Load(t); ok {
		return cached.(map[string][]int)
	}
	out := make(map[string][]int)
	collectStructTags(t, nil, out)
	structTagCache.Store(t, out)
	return out
}

func collectStructTags(t reflect.Type, prefix []int, out map[string][]int) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		idxPath := append(append([]int{}, prefix...), i)
		if f.Anonymous {
			// Anonymous struct fields promote their tags to the outer level
			// per Go's visibility rules — the JSON package matches them the
			// same way. Recurse using the extended index path.
			ft := f.Type
			if ft.Kind() == reflect.Ptr {
				ft = ft.Elem()
			}
			if ft.Kind() == reflect.Struct {
				collectStructTags(ft, idxPath, out)
			}
			continue
		}
		tag := f.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		name := tag
		if idx := strings.Index(tag, ","); idx >= 0 {
			name = tag[:idx]
		}
		if name == "" {
			continue
		}
		if _, exists := out[name]; !exists {
			out[name] = idxPath
		}
	}
}
