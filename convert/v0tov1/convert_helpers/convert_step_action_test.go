package converthelpers

import (
	"encoding/json"
	"strings"
	"testing"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/google/go-cmp/cmp"
)

func TestConvertStepAction(t *testing.T) {
	tests := []struct {
		name           string
		step           *v0.Step
		expectedScript v1.Stringorslice
		expectedWith   map[string]interface{} // decoded from PLUGIN_WITH JSON; nil if not expected
		expectedExtra  map[string]interface{} // non-PLUGIN_WITH env entries
	}{
		{
			name: "basic action step with uses only",
			step: &v0.Step{
				Spec: &v0.StepAction{
					Uses: "actions/checkout@v3",
				},
			},
			expectedScript: v1.Stringorslice{"plugin -kind action -name actions/checkout@v3"},
		},
		{
			name: "action step with with params encoded as PLUGIN_WITH JSON",
			step: &v0.Step{
				Spec: &v0.StepAction{
					Uses: "actions/setup-go@v4",
					With: map[string]interface{}{
						"go-version": "1.21",
						"cache":      true,
					},
				},
			},
			expectedScript: v1.Stringorslice{"plugin -kind action -name actions/setup-go@v4"},
			expectedWith: map[string]interface{}{
				"go-version": "1.21",
				"cache":      true,
			},
		},
		{
			name: "action step with version constraint value preserves leading =",
			step: &v0.Step{
				Spec: &v0.StepAction{
					Uses: "actions/setup-go@v3.5.0",
					With: map[string]interface{}{
						"go-version": "=1.20.1",
					},
				},
			},
			expectedScript: v1.Stringorslice{"plugin -kind action -name actions/setup-go@v3.5.0"},
			expectedWith: map[string]interface{}{
				"go-version": "=1.20.1",
			},
		},
		{
			name: "action step with env vars",
			step: &v0.Step{
				Spec: &v0.StepAction{
					Uses: "actions/upload-artifact@v3",
					Envs: map[string]string{
						"GITHUB_TOKEN": "my-token",
					},
				},
			},
			expectedScript: v1.Stringorslice{"plugin -kind action -name actions/upload-artifact@v3"},
			expectedExtra: map[string]interface{}{
				"GITHUB_TOKEN": "my-token",
			},
		},
		{
			name: "action step with both with and env",
			step: &v0.Step{
				Spec: &v0.StepAction{
					Uses: "docker/build-push-action@v5",
					With: map[string]interface{}{
						"context": ".",
						"push":    true,
						"tags":    "myapp:latest",
					},
					Envs: map[string]string{
						"DOCKER_BUILDKIT": "1",
					},
				},
			},
			expectedScript: v1.Stringorslice{"plugin -kind action -name docker/build-push-action@v5"},
			expectedWith: map[string]interface{}{
				"context": ".",
				"push":    true,
				"tags":    "myapp:latest",
			},
			expectedExtra: map[string]interface{}{
				"DOCKER_BUILDKIT": "1",
			},
		},
		{
			name: "action step with empty with and env",
			step: &v0.Step{
				Spec: &v0.StepAction{
					Uses: "actions/cache@v3",
					With: map[string]interface{}{},
					Envs: map[string]string{},
				},
			},
			expectedScript: v1.Stringorslice{"plugin -kind action -name actions/cache@v3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertStepAction(tt.step)
			if result == nil {
				t.Fatal("expected non-nil result")
			}

			if diff := cmp.Diff(tt.expectedScript, result.Script); diff != "" {
				t.Errorf("Script mismatch (-want +got):\n%s", diff)
			}

			if tt.expectedWith == nil && tt.expectedExtra == nil {
				if result.Env != nil {
					t.Errorf("expected nil Env, got %+v", result.Env.Value)
				}
				return
			}

			if result.Env == nil {
				t.Fatalf("expected non-nil Env")
			}

			got, ok := result.Env.Value.(map[string]interface{})
			if !ok {
				t.Fatalf("Env.Value is %T, expected map[string]interface{}", result.Env.Value)
			}

			if tt.expectedWith != nil {
				rawJSON, ok := got["PLUGIN_WITH"].(string)
				if !ok {
					t.Fatalf("expected PLUGIN_WITH string entry, got %v", got["PLUGIN_WITH"])
				}
				var decoded map[string]interface{}
				if err := json.Unmarshal([]byte(rawJSON), &decoded); err != nil {
					t.Fatalf("PLUGIN_WITH is not valid JSON: %v", err)
				}
				if diff := cmp.Diff(tt.expectedWith, decoded); diff != "" {
					t.Errorf("PLUGIN_WITH decoded mismatch (-want +got):\n%s", diff)
				}
			} else if _, present := got["PLUGIN_WITH"]; present {
				t.Errorf("unexpected PLUGIN_WITH entry: %v", got["PLUGIN_WITH"])
			}

			extra := map[string]interface{}{}
			for k, v := range got {
				if k == "PLUGIN_WITH" {
					continue
				}
				extra[k] = v
			}
			expectedExtra := tt.expectedExtra
			if expectedExtra == nil {
				expectedExtra = map[string]interface{}{}
			}
			if diff := cmp.Diff(expectedExtra, extra); diff != "" {
				t.Errorf("Extra env mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertStepAction_PluginWithNotHTMLEscaped(t *testing.T) {
	step := &v0.Step{
		Spec: &v0.StepAction{
			Uses: "actions/setup-go@v4",
			With: map[string]interface{}{
				"check-latest": "<+pipeline.variables.checkLatest>",
				"go-version":   "1.20.1",
			},
		},
	}

	result := ConvertStepAction(step)
	if result == nil || result.Env == nil {
		t.Fatal("expected non-nil result with Env")
	}

	got, ok := result.Env.Value.(map[string]interface{})
	if !ok {
		t.Fatalf("Env.Value is %T, expected map[string]interface{}", result.Env.Value)
	}

	rawJSON, ok := got["PLUGIN_WITH"].(string)
	if !ok {
		t.Fatalf("expected PLUGIN_WITH string entry, got %v", got["PLUGIN_WITH"])
	}

	if strings.Contains(rawJSON, `\u003c`) || strings.Contains(rawJSON, `\u003e`) {
		t.Errorf("PLUGIN_WITH should not HTML-escape expressions, got: %s", rawJSON)
	}
	if !strings.Contains(rawJSON, "<+pipeline.variables.checkLatest>") {
		t.Errorf("PLUGIN_WITH should contain the raw expression, got: %s", rawJSON)
	}
}

func TestConvertStepAction_NilCases(t *testing.T) {
	tests := []struct {
		name string
		step *v0.Step
	}{
		{
			name: "nil step",
			step: nil,
		},
		{
			name: "nil spec",
			step: &v0.Step{
				Spec: nil,
			},
		},
		{
			name: "wrong spec type",
			step: &v0.Step{
				Spec: &v0.StepRun{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertStepAction(tt.step)
			if result != nil {
				t.Errorf("expected nil result, got %v", result)
			}
		})
	}
}
