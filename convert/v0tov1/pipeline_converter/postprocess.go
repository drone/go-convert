package pipelineconverter

import (
	"reflect"
	"strings"

	convertexpressions "github.com/drone/go-convert/convert/convertexpressions"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// expressionProcessor walks v1 pipeline structs and converts Harness expressions
// using context-aware trie-based matching.
type expressionProcessor struct {
	currentStepID      string // set when inside a v1.Step
	currentStepType    string // set when inside a v1.Step
	currentStepV1Path  string // set when inside a v1.Step (the step's full v1 FQN)
	currentStageID     string // set when inside a v1.Stage
	currentStepGroupID string // ">"-joined ancestor group chain of the current scope
	inRollback         bool   // true while walking a stage's rollback steps
	useFQN             bool   // whether to use fully qualified names for step expressions

	// stepInfoByFQN is the sole step (and step group) registry, keyed by each
	// step's full v1 FQN. nil when no pipeline-level context is available.
	stepInfoByFQN map[string]*convertexpressions.StepInfoFQN
}

// PostProcessExpressions walks the converted v1 entity (Pipeline, Template,
// InputSet, Trigger, etc.) and converts all Harness expressions in string and
// flexible.Field values using ConvertExpressionWithTrie.
//
// stepInfoByFQN may be nil when no pipeline-level step context is available
// (e.g. for Template/InputSet/Trigger wrappers). useFQN controls whether to
// use fully qualified names for step expressions; pass false for context-free
// conversion of non-pipeline entities.
func PostProcessExpressions(target any, stepInfoByFQN map[string]*convertexpressions.StepInfoFQN, useFQN bool) {
	if target == nil {
		return
	}
	val := reflect.ValueOf(target)
	if !val.IsValid() {
		return
	}
	if val.Kind() == reflect.Ptr && val.IsNil() {
		return
	}

	p := &expressionProcessor{
		stepInfoByFQN: stepInfoByFQN,
		useFQN:        useFQN,
	}

	// Logger context flag is set once for the duration of the walk rather
	// than per-expression to avoid mutex thrashing on the global logger.
	GetExpressionLogger().SetIncludeContext(true)

	p.processValue(val)
}

// tryConvertString checks if s contains a Harness expression and, if so, converts it.
// Returns the converted string and true if a conversion was made, or the original string and false otherwise.
func (p *expressionProcessor) tryConvertString(s string) (string, bool) {
	if !strings.Contains(s, "<+") && !strings.Contains(s, "${{") {
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

// withStageContext saves/restores current stage context around fn.
func (p *expressionProcessor) withStageContext(stageID string, fn func()) {
	saved := p.currentStageID
	p.currentStageID = stageID
	fn()
	p.currentStageID = saved
}

// withStepGroupContext saves/restores the current step group context around fn.
// It accumulates a chain (e.g. "A>B>C") matching the chain-based key used
// during step registration in ConvertSteps.
func (p *expressionProcessor) withStepGroupContext(groupID string, fn func()) {
	saved := p.currentStepGroupID
	if p.currentStepGroupID != "" {
		p.currentStepGroupID = p.currentStepGroupID + ">" + groupID
	} else {
		p.currentStepGroupID = groupID
	}
	fn()
	p.currentStepGroupID = saved
}

// withRollbackContext marks the current scope as a stage's rollback steps so
// step FQNs are reconstructed with a "rollback" first segment.
func (p *expressionProcessor) withRollbackContext(fn func()) {
	saved := p.inRollback
	p.inRollback = true
	fn()
	p.inRollback = saved
}

// stepFQN reconstructs a step's full v1 FQN from its stage, ">"-joined ancestor
// group chain, and leaf ID. Rollback steps use a "rollback" first segment
// (matching convertV0PathToV1Path); deeper levels use "steps". "" if stage unknown.
func stepFQN(stageID, groupChain, stepID string, rollback bool) string {
	if stageID == "" {
		return ""
	}
	root := "steps"
	if rollback {
		root = "rollback"
	}
	var b strings.Builder
	b.WriteString("pipeline.stages.")
	b.WriteString(stageID)
	groups := splitGroupChain(groupChain)
	if len(groups) > 0 {
		b.WriteString(".")
		b.WriteString(root)
		b.WriteString(".")
		b.WriteString(groups[0])
		for _, g := range groups[1:] {
			b.WriteString(".steps.")
			b.WriteString(g)
		}
		b.WriteString(".steps.")
		b.WriteString(stepID)
	} else {
		b.WriteString(".")
		b.WriteString(root)
		b.WriteString(".")
		b.WriteString(stepID)
	}
	return b.String()
}

// withStepContext saves the current step context, sets it from the given step ID,
// executes fn, then restores the previous context.
func (p *expressionProcessor) withStepContext(stepID string, fn func()) {
	savedID, savedType, savedPath := p.currentStepID, p.currentStepType, p.currentStepV1Path

	p.currentStepID = stepID
	// Reconstruct the FQN and look up the step type; leave empty for untyped or
	// template steps (matching prior behavior).
	fqn := stepFQN(p.currentStageID, p.currentStepGroupID, stepID, p.inRollback)
	if info, ok := p.stepInfoByFQN[fqn]; ok {
		p.currentStepType = info.Type
		p.currentStepV1Path = fqn
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

	// Check if this is a v1.Stage — set stage context
	if t == reflect.TypeOf(v1.Stage{}) {
		idField := val.FieldByName("Id")
		if idField.IsValid() && idField.Kind() == reflect.String && idField.String() != "" {
			p.withStageContext(idField.String(), func() {
				p.processStructFields(val)
			})
			return
		}
	}

	// Check if this is a v1.Step — set step context (and group context if it has a Group)
	if t == reflect.TypeOf(v1.Step{}) {
		idField := val.FieldByName("Id")
		groupField := val.FieldByName("Group")
		hasGroup := groupField.IsValid() && groupField.Kind() == reflect.Ptr && !groupField.IsNil()
		stepID := ""
		if idField.IsValid() && idField.Kind() == reflect.String {
			stepID = idField.String()
		}
		if stepID != "" {
			if hasGroup {
				// Step group: set both step context and group context
				p.withStepContext(stepID, func() {
					p.withStepGroupContext(stepID, func() {
						p.processStructFields(val)
					})
				})
			} else {
				p.withStepContext(stepID, func() {
					p.processStructFields(val)
				})
			}
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

		// Skip Trigger.InputYaml: it holds an already-converted v1 YAML
		// document; re-walking its raw text could mangle structure.
		if fieldType.Name == "InputYaml" {
			continue
		}

		// Stage.Rollback steps live under a "rollback" FQN segment; mark the
		// scope so step FQNs are reconstructed correctly while walking them.
		if fieldType.Name == "Rollback" {
			p.withRollbackContext(func() {
				p.processValue(field)
			})
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
func (p *expressionProcessor) processString(s string) string {
	spans := convertexpressions.FindHarnessExprs(s)
	if len(spans) == 0 {
		return s
	}

	// Build the ConversionContext; the trie resolves step types lazily at
	// step.spec nodes and, in FQN mode, expands v1Path to the step's FQN base.
	ctx := p.buildConversionContext()

	// Get expression logger for logging conversions. The IncludeContext flag
	// is set once at the top of PostProcessExpressions, not per call.
	logger := GetExpressionLogger()

	var b strings.Builder
	prev := 0
	for _, span := range spans {
		b.WriteString(s[prev:span[0]])

		expr := s[span[0]:span[1]] // full <+...> expression

		// Apply trie-based context-aware conversion rules
		// Step type resolution happens lazily inside the trie
		converted := convertexpressions.ConvertExpressionWithTrie(expr, ctx, false)

		// Log the conversion with full context
		logger.LogConversion(expr, converted, ctx)

		b.WriteString(converted)
		prev = span[1]
	}
	b.WriteString(s[prev:])
	return b.String()
}

// buildConversionContext creates a ConversionContext for the current step scope.
func (p *expressionProcessor) buildConversionContext() *convertexpressions.ConversionContext {
	var chain []string
	if p.currentStepGroupID != "" {
		chain = strings.Split(p.currentStepGroupID, ">")
	}
	return &convertexpressions.ConversionContext{
		CurrentStepID:         p.currentStepID,
		CurrentStepType:       p.currentStepType,
		UseFQN:                p.useFQN,
		CurrentStageID:        p.currentStageID,
		StepInfoByFQN:         p.stepInfoByFQN,
		CurrentStepGroupChain: chain,
		CurrentFQN:            p.currentStepV1Path,
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
