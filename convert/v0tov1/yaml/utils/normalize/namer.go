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

package normalize

import (
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// generator generates unique identifiers
type generator struct {
	store map[string]struct{}
}

func newGenerator() *generator {
	return &generator{
		store: map[string]struct{}{},
	}
}

// register registers a name with the store.
func (s *generator) register(name string) bool {
	if _, ok := s.store[name]; ok {
		return false
	}
	s.store[name] = struct{}{}
	return true
}

// generate generates and registeres a unique name with the
// store. If the base name is already registered, a unique
// suffix is appended to the name.
func (s *generator) generate(name ...string) string {
	var base string
	// choose the first non-empty name.
	for _, s := range name {
		if s != "" {
			base = s
			break
		}
	}

	// convert the name to a slug
	base = slugify(base)

	// attempt to register the name
	if s.register(base) {
		return base
	}
	// append a suffix to the name and register
	// the first unique combination.
	for i := 1; ; i++ {
		next := base + strconv.Itoa(i)
		if s.register(next) {
			return next
		}
	}
}

//
// slugs
//

var safeRanges = []*unicode.RangeTable{
	unicode.Letter,
	unicode.Number,
}

func safe(r rune) rune {
	switch {
	case unicode.IsOneOf(safeRanges, r):
		return unicode.ToLower(r)
	}
	return -1
}

func slugify(s string) string {
	s = norm.NFKD.String(s)
	s = strings.Map(safe, s)
	return s
}
