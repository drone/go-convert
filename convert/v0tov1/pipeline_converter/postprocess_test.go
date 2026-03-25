package pipelineconverter

import (
	"reflect"
	"testing"

	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/drone/go-convert/internal/flexible"
)

func TestIsFlexibleField(t *testing.T) {
	// flexible.Field[bool] should be detected
	f := flexible.Field[bool]{Value: true}
	if !isFlexibleField(reflect.TypeOf(f)) {
		t.Error("expected flexible.Field[bool] to be detected")
	}

	// A regular struct should not be detected
	type notFlexible struct {
		Name  string
		Value interface{}
	}
	nf := notFlexible{}
	if isFlexibleField(reflect.TypeOf(nf)) {
		t.Error("expected notFlexible to NOT be detected as flexible field")
	}
}

// Note: Step type resolution tests have been moved to convert_expressions package
// since resolution now happens lazily inside the trie during path matching.
// See trie.go resolveStepType() and matchRecursive() for the implementation.

func TestProcessString_SingleExpression(t *testing.T) {
	p := &expressionProcessor{
		stepTypeMap: map[string]*StepInfo{
			"runStep1": {Type: "Run", V0Path: "pipeline.stages.build.spec.execution.steps.runStep1"},
		},
		flatTypeMap:     map[string]string{"runStep1": "Run"},
		currentStepID:   "runStep1",
		currentStepType: "Run",
	}

	// Trie converts spec.execution.steps → steps
	result := p.processString("<+pipeline.stages.build.spec.execution.steps.runStep1.output>")
	expected := "<+pipeline.stages.build.steps.runStep1.output>"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestProcessString_NoExpression(t *testing.T) {
	p := &expressionProcessor{
		stepTypeMap: map[string]*StepInfo{},
		flatTypeMap: map[string]string{},
	}

	// Plain string should be returned as-is
	result := p.processString("hello world")
	if result != "hello world" {
		t.Errorf("expected 'hello world', got %q", result)
	}
}

func TestProcessString_MixedContent(t *testing.T) {
	p := &expressionProcessor{
		stepTypeMap: map[string]*StepInfo{
			"runStep1": {Type: "Run", V0Path: "pipeline.stages.build.spec.execution.steps.runStep1"},
		},
		flatTypeMap:     map[string]string{"runStep1": "Run"},
		currentStepID:   "runStep1",
		currentStepType: "Run",
	}

	// Mixed content with text and expression (trie converts spec.execution.steps → steps)
	result := p.processString("prefix <+pipeline.stages.build.spec.execution.steps.runStep1.output> suffix")
	expected := "prefix <+pipeline.stages.build.steps.runStep1.output> suffix"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestPostProcessExpressions_Pipeline(t *testing.T) {
	stepTypeMap := map[string]*StepInfo{
		"runStep1": {Type: "Run", V0Path: "pipeline.stages.build.spec.execution.steps.runStep1"},
	}

	pipeline := &v1.Pipeline{
		Id:   "myPipeline",
		Name: "Build Pipeline",
		Stages: []*v1.Stage{
			{
				Id:   "build",
				Name: "Build Stage",
				If:   "<+pipeline.stages.build.spec.execution.steps.runStep1.output>",
				Steps: []*v1.Step{
					{
						Id:      "runStep1",
						Name:    "Run Step",
						Timeout: "5m",
					},
				},
			},
		},
	}

	PostProcessExpressions(pipeline, stepTypeMap)

	// Stage If condition should have spec.execution.steps → steps converted
	expectedIf := "<+pipeline.stages.build.steps.runStep1.output>"
	if pipeline.Stages[0].If != expectedIf {
		t.Errorf("stage.If: expected %q, got %q", expectedIf, pipeline.Stages[0].If)
	}
}

func TestPostProcessExpressions_FlexibleFieldExpression(t *testing.T) {
	stepTypeMap := map[string]*StepInfo{
		"runStep1": {Type: "Run", V0Path: "pipeline.stages.build.spec.execution.steps.runStep1"},
		"step2":    {Type: "Run", V0Path: "pipeline.stages.build.spec.execution.steps.step2"},
	}

	envField := &flexible.Field[map[string]string]{}
	envField.SetExpression("<+pipeline.stages.build.spec.execution.steps.step2.output>")

	pipeline := &v1.Pipeline{
		Stages: []*v1.Stage{
			{
				Id: "build",
				Steps: []*v1.Step{
					{
						Id:  "runStep1",
						Env: envField,
					},
				},
			},
		},
	}

	PostProcessExpressions(pipeline, stepTypeMap)

	// FlexibleField with expression should have spec.execution.steps → steps converted
	str, ok := pipeline.Stages[0].Steps[0].Env.AsString()
	if !ok {
		t.Fatal("expected Env to still be a string expression")
	}
	expected := "<+pipeline.stages.build.steps.step2.output>"
	if str != expected {
		t.Errorf("step.Env expression: expected %q, got %q", expected, str)
	}
}

func TestPostProcessExpressions_MapStringValues(t *testing.T) {
	stepTypeMap := map[string]*StepInfo{
		"runStep1": {Type: "Run", V0Path: "pipeline.stages.build.spec.execution.steps.runStep1"},
	}

	pipeline := &v1.Pipeline{
		Stages: []*v1.Stage{
			{
				Id: "build",
				Env: map[string]interface{}{
					"STEP_OUTPUT": "<+pipeline.stages.build.spec.execution.steps.runStep1.output>",
					"PLAIN":       "no-expression",
				},
			},
		},
	}

	PostProcessExpressions(pipeline, stepTypeMap)

	if val, ok := pipeline.Stages[0].Env["STEP_OUTPUT"].(string); ok {
		expected := "<+pipeline.stages.build.steps.runStep1.output>"
		if val != expected {
			t.Errorf("stage.Env[STEP_OUTPUT]: expected %q, got %q", expected, val)
		}
	} else {
		t.Error("stage.Env[STEP_OUTPUT] should be a string")
	}

	if val, ok := pipeline.Stages[0].Env["PLAIN"].(string); ok {
		if val != "no-expression" {
			t.Errorf("stage.Env[PLAIN]: expected 'no-expression', got %q", val)
		}
	}
}

func TestPostProcessExpressions_NilPipeline(t *testing.T) {
	// Should not panic
	PostProcessExpressions(nil, nil)
}

func TestPostProcessExpressions_StageIdentifier(t *testing.T) {
	stepTypeMap := map[string]*StepInfo{
		"step1": {Type: "Run", V0Path: "pipeline.stages.build.spec.execution.steps.step1"},
	}

	pipeline := &v1.Pipeline{
		Stages: []*v1.Stage{
			{
				Id: "build",
				If: "<+stage.spec.execution.steps.step1.output>",
			},
		},
	}

	PostProcessExpressions(pipeline, stepTypeMap)

	expected := "<+stage.steps.step1.output>"
	if pipeline.Stages[0].If != expected {
		t.Errorf("stage.If: expected %q, got %q", expected, pipeline.Stages[0].If)
	}
}

func TestPostProcessExpressions_StringorsliceScript(t *testing.T) {
	// Test that expressions in Stringorslice ([]string) are converted
	stepTypeMap := map[string]*StepInfo{
		"runStep1": {Type: "Run", V0Path: "pipeline.stages.build.spec.execution.steps.runStep1"},
		"runStep2": {Type: "Run", V0Path: "pipeline.stages.build.spec.execution.steps.runStep2"},
	}

	pipeline := &v1.Pipeline{
		Stages: []*v1.Stage{
			{
				Id: "build",
				Steps: []*v1.Step{
					{
						Id: "runStep1",
						Run: &v1.StepRun{
							Script: v1.Stringorslice{
								"echo <+pipeline.stages.build.spec.execution.steps.runStep1.output>",
								"echo hello",
								"echo <+pipeline.stages.build.spec.execution.steps.runStep2.output>",
							},
						},
					},
				},
			},
		},
	}

	PostProcessExpressions(pipeline, stepTypeMap)

	// First script line should have spec.execution.steps → steps converted
	expected0 := "echo <+pipeline.stages.build.steps.runStep1.output>"
	if pipeline.Stages[0].Steps[0].Run.Script[0] != expected0 {
		t.Errorf("script[0]: expected %q, got %q", expected0, pipeline.Stages[0].Steps[0].Run.Script[0])
	}

	// Second line has no expression, should be unchanged
	if pipeline.Stages[0].Steps[0].Run.Script[1] != "echo hello" {
		t.Errorf("script[1]: expected 'echo hello', got %q", pipeline.Stages[0].Steps[0].Run.Script[1])
	}

	// Third line should also have spec.execution.steps → steps converted
	expected2 := "echo <+pipeline.stages.build.steps.runStep2.output>"
	if pipeline.Stages[0].Steps[0].Run.Script[2] != expected2 {
		t.Errorf("script[2]: expected %q, got %q", expected2, pipeline.Stages[0].Steps[0].Run.Script[2])
	}
}

func TestPostProcessExpressions_StepNeeds(t *testing.T) {
	// Test that expressions in Step.Needs (Stringorslice) are converted
	stepTypeMap := map[string]*StepInfo{
		"step1": {Type: "Run", V0Path: "pipeline.stages.build.spec.execution.steps.step1"},
		"step2": {Type: "Run", V0Path: "pipeline.stages.build.spec.execution.steps.step2"},
	}

	pipeline := &v1.Pipeline{
		Stages: []*v1.Stage{
			{
				Id: "build",
				Steps: []*v1.Step{
					{
						Id: "step2",
						Needs: v1.Stringorslice{
							"<+stage.spec.execution.steps.step1.status>",
						},
					},
				},
			},
		},
	}

	PostProcessExpressions(pipeline, stepTypeMap)

	// Needs should have spec.execution.steps → steps converted
	expected := "<+stage.steps.step1.status>"
	if pipeline.Stages[0].Steps[0].Needs[0] != expected {
		t.Errorf("needs[0]: expected %q, got %q", expected, pipeline.Stages[0].Steps[0].Needs[0])
	}
}

func TestPostProcessExpressions_InterfaceSlice(t *testing.T) {
	// Test expressions in []interface{} slices
	stepTypeMap := map[string]*StepInfo{
		"step1": {Type: "Run", V0Path: "pipeline.stages.build.spec.execution.steps.step1"},
	}

	pipeline := &v1.Pipeline{
		Stages: []*v1.Stage{
			{
				Id: "build",
				Outputs: map[string]interface{}{
					"items": []interface{}{
						"<+pipeline.stages.build.spec.execution.steps.step1.output>",
						"plain-value",
						"<+stage.spec.execution.steps.step1.output>",
					},
				},
			},
		},
	}

	PostProcessExpressions(pipeline, stepTypeMap)

	items, ok := pipeline.Stages[0].Outputs["items"].([]interface{})
	if !ok {
		t.Fatal("expected items to be []interface{}")
	}

	expected0 := "<+pipeline.stages.build.steps.step1.output>"
	if items[0] != expected0 {
		t.Errorf("items[0]: expected %q, got %v", expected0, items[0])
	}

	if items[1] != "plain-value" {
		t.Errorf("items[1]: expected 'plain-value', got %v", items[1])
	}

	expected2 := "<+stage.steps.step1.output>"
	if items[2] != expected2 {
		t.Errorf("items[2]: expected %q, got %v", expected2, items[2])
	}
}

func TestPostProcessExpressions_NestedMapInterface(t *testing.T) {
	// Test deeply nested map[string]interface{} with expressions
	stepTypeMap := map[string]*StepInfo{
		"step1": {Type: "Run", V0Path: "pipeline.stages.build.spec.execution.steps.step1"},
		"step2": {Type: "Run", V0Path: "pipeline.stages.build.spec.execution.steps.step2"},
	}

	pipeline := &v1.Pipeline{
		Stages: []*v1.Stage{
			{
				Id: "build",
				Steps: []*v1.Step{
					{
						Id: "step1",
						With: map[string]interface{}{
							"level1": map[string]interface{}{
								"level2": "<+pipeline.stages.build.spec.execution.steps.step2.output>",
							},
							"direct": "<+stage.spec.execution.steps.step2.output>",
						},
					},
				},
			},
		},
	}

	PostProcessExpressions(pipeline, stepTypeMap)

	// Check direct value
	if direct, ok := pipeline.Stages[0].Steps[0].With["direct"].(string); ok {
		expected := "<+stage.steps.step2.output>"
		if direct != expected {
			t.Errorf("with.direct: expected %q, got %q", expected, direct)
		}
	} else {
		t.Error("with.direct should be a string")
	}

	// Check nested value
	if level1, ok := pipeline.Stages[0].Steps[0].With["level1"].(map[string]interface{}); ok {
		if level2, ok := level1["level2"].(string); ok {
			expected := "<+pipeline.stages.build.steps.step2.output>"
			if level2 != expected {
				t.Errorf("with.level1.level2: expected %q, got %q", expected, level2)
			}
		} else {
			t.Error("with.level1.level2 should be a string")
		}
	} else {
		t.Error("with.level1 should be a map")
	}
}
