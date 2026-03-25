package pipelineconverter

import (
	"reflect"
	"strings"

	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	convertexpressions "github.com/peroxidemonke7/v0tov1_expressions/convertexpressions"
)

// expressionProcessor walks v1 pipeline structs and converts Harness expressions
// using context-aware trie-based matching.
type expressionProcessor struct {
	stepTypeMap       map[string]*StepInfo // step ID → {Type, V0Path, V1Path}
	currentStepID     string               // set when inside a v1.Step
	currentStepType   string               // set when inside a v1.Step
	currentStepV1Path string               // set when inside a v1.Step

	// Cached derived maps, built lazily on first processString call.
	flatTypeMap   map[string]string // step ID → type
	stepV1PathMap map[string]string // step ID → v1 path
	mapsBuilt     bool
}

// PostProcessExpressions walks the converted v1 pipeline and converts all Harness
// expressions in string and flexible.Field values using ConvertExpressionWithTrie.
func PostProcessExpressions(pipeline *v1.Pipeline, stepTypeMap map[string]*StepInfo) {
	if pipeline == nil {
		return
	}

	p := &expressionProcessor{
		stepTypeMap: stepTypeMap,
	}

	p.processValue(reflect.ValueOf(pipeline))
}

// tryConvertString checks if s contains a Harness expression and, if so, converts it.
// Returns the converted string and true if a conversion was made, or the original string and false otherwise.
func (p *expressionProcessor) tryConvertString(s string) (string, bool) {
	if !strings.Contains(s, "<+") {
		return s, false
	}
	converted := p.processString(s)
	return converted, converted != s
}

// processValue recursively walks a reflect.Value and converts expressions.
func (p *expressionProcessor) processValue(val reflect.Value) {
	switch val.Kind() {
	case reflect.Ptr:
		if val.IsNil() {
			return
		}
		p.processValue(val.Elem())

	case reflect.Struct:
		p.processStruct(val)

	case reflect.Slice:
		p.processSlice(val)

	case reflect.Map:
		p.processMap(val)

	case reflect.Interface:
		if val.IsNil() {
			return
		}
		elem := val.Elem()
		if elem.Kind() == reflect.String {
			if converted, changed := p.tryConvertString(elem.String()); changed && val.CanSet() {
				val.Set(reflect.ValueOf(converted))
			}
		} else {
			p.processValue(elem)
		}

	case reflect.String:
		// Strings inside structs are handled in processStruct via field setting.
		// Strings in slices are handled in processSlice.
		// This case is reached for non-settable strings (e.g., map keys) - nothing to do.
	}
}

// withStepContext saves the current step context, sets it from the given step ID,
// executes fn, then restores the previous context.
func (p *expressionProcessor) withStepContext(stepID string, fn func()) {
	savedID, savedType, savedPath := p.currentStepID, p.currentStepType, p.currentStepV1Path

	p.currentStepID = stepID
	if info, ok := p.stepTypeMap[stepID]; ok {
		p.currentStepType = info.Type
		p.currentStepV1Path = info.V1Path
	} else {
		p.currentStepType = ""
		p.currentStepV1Path = ""
	}

	fn()

	p.currentStepID = savedID
	p.currentStepType = savedType
	p.currentStepV1Path = savedPath
}

// processStruct handles struct values. It detects v1.Step to set step context,
// and flexible.Field[T] to process expression values.
func (p *expressionProcessor) processStruct(val reflect.Value) {
	t := val.Type()

	// Check if this is a flexible.Field[T] (single field "Value" of type interface{})
	if isFlexibleField(t) {
		p.processFlexibleField(val)
		return
	}

	// Check if this is a v1.Step — set step context
	if t == reflect.TypeOf(v1.Step{}) {
		idField := val.FieldByName("Id")
		if idField.IsValid() && idField.Kind() == reflect.String && idField.String() != "" {
			p.withStepContext(idField.String(), func() {
				p.processStructFields(val)
			})
			return
		}
	}

	p.processStructFields(val)
}

// processStructFields iterates over exported struct fields and converts expressions.
func (p *expressionProcessor) processStructFields(val reflect.Value) {
	t := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := t.Field(i)

		// Skip unexported fields
		if !fieldType.IsExported() {
			continue
		}

		// Skip Context field (metadata, not pipeline content)
		if fieldType.Name == "Context" {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			if field.CanSet() {
				if converted, changed := p.tryConvertString(field.String()); changed {
					field.SetString(converted)
				}
			}
		default:
			p.processValue(field)
		}
	}
}

// processFlexibleField handles flexible.Field[T] values.
// If the inner Value is a string containing expressions, convert it.
// If it's a struct, recurse into it.
func (p *expressionProcessor) processFlexibleField(val reflect.Value) {
	valueField := val.FieldByName("Value")
	if !valueField.IsValid() || valueField.IsNil() {
		return
	}

	elem := valueField.Elem()
	if elem.Kind() == reflect.String {
		if converted, changed := p.tryConvertString(elem.String()); changed && valueField.CanSet() {
			valueField.Set(reflect.ValueOf(converted))
		}
	} else {
		// Recurse into the struct/map/slice value
		p.processValue(elem)
	}
}

// processSlice processes slice values, converting string elements that contain expressions.
func (p *expressionProcessor) processSlice(val reflect.Value) {
	for i := 0; i < val.Len(); i++ {
		elem := val.Index(i)

		// Handle string elements directly (e.g., []string, Stringorslice)
		if elem.Kind() == reflect.String && elem.CanSet() {
			if converted, changed := p.tryConvertString(elem.String()); changed {
				elem.SetString(converted)
			}
			continue
		}

		// Handle interface elements that contain strings
		if elem.Kind() == reflect.Interface && !elem.IsNil() {
			inner := elem.Elem()
			if inner.Kind() == reflect.String {
				if converted, changed := p.tryConvertString(inner.String()); changed && elem.CanSet() {
					elem.Set(reflect.ValueOf(converted))
				}
				continue
			}
		}

		// Recurse into other types (structs, pointers, nested slices, maps)
		p.processValue(elem)
	}
}

// processMap processes map values, converting string values that contain expressions.
func (p *expressionProcessor) processMap(val reflect.Value) {
	if val.IsNil() {
		return
	}
	for _, key := range val.MapKeys() {
		elem := val.MapIndex(key)

		// Unwrap interface
		if elem.Kind() == reflect.Interface && !elem.IsNil() {
			inner := elem.Elem()
			if inner.Kind() == reflect.String {
				if converted, changed := p.tryConvertString(inner.String()); changed {
					val.SetMapIndex(key, reflect.ValueOf(converted))
				}
				continue
			}
			// For non-string interface values, recurse
			p.processValue(inner)
			continue
		}

		if elem.Kind() == reflect.String {
			if converted, changed := p.tryConvertString(elem.String()); changed {
				val.SetMapIndex(key, reflect.ValueOf(converted))
			}
			continue
		}

		// Recurse into struct/ptr/slice map values
		p.processValue(elem)
	}
}

// processString processes a string that may contain one or more <+...> expressions.
// Each expression is converted using trie-based context-aware rules.
// Step type resolution happens lazily inside the trie when it encounters step.spec expressions.
// FQN mode is enabled to convert relative step expressions to fully qualified paths.
func (p *expressionProcessor) processString(s string) string {
	spans := convertexpressions.FindHarnessExprs(s)
	if len(spans) == 0 {
		return s
	}

	// Build ConversionContext directly from stepTypeMap — avoids redundant derived maps.
	// The trie will resolve step types only when needed (at step.spec nodes).
	// FQN mode replaces v1Path at step_node with the step's v1 FQN base path.
	ctx := p.buildConversionContext()

	var b strings.Builder
	prev := 0
	for _, span := range spans {
		b.WriteString(s[prev:span[0]])

		expr := s[span[0]:span[1]] // full <+...> expression

		// Apply trie-based context-aware conversion rules
		// Step type resolution happens lazily inside the trie
		converted := convertexpressions.ConvertExpressionWithTrie(expr, ctx, false)

		b.WriteString(converted)
		prev = span[1]
	}
	b.WriteString(s[prev:])
	return b.String()
}

// ensureDerivedMaps builds flatTypeMap and stepV1PathMap once from stepTypeMap.
func (p *expressionProcessor) ensureDerivedMaps() {
	if p.mapsBuilt {
		return
	}
	p.flatTypeMap = make(map[string]string, len(p.stepTypeMap))
	p.stepV1PathMap = make(map[string]string, len(p.stepTypeMap))
	for id, info := range p.stepTypeMap {
		p.flatTypeMap[id] = info.Type
		p.stepV1PathMap[id] = info.V1Path
	}
	p.mapsBuilt = true
}

// buildConversionContext creates a ConversionContext for the current step scope.
func (p *expressionProcessor) buildConversionContext() *convertexpressions.ConversionContext {
	p.ensureDerivedMaps()
	return &convertexpressions.ConversionContext{
		StepID:            p.currentStepID,
		CurrentStepType:   p.currentStepType,
		CurrentStepV1Path: p.currentStepV1Path,
		StepTypeMap:       p.flatTypeMap,
		StepV1PathMap:     p.stepV1PathMap,
		UseFQN:            true,
	}
}

// isFlexibleField checks if a reflect.Type is a flexible.Field[T].
// flexible.Field[T] is a struct with a single exported field "Value" of type interface{}.
func isFlexibleField(t reflect.Type) bool {
	if t.Kind() != reflect.Struct {
		return false
	}
	if t.NumField() != 1 {
		return false
	}
	f := t.Field(0)
	return f.Name == "Value" && f.Type.Kind() == reflect.Interface
}
