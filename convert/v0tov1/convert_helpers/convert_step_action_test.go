package converthelpers

import (
	"testing"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/drone/go-convert/internal/flexible"
	"github.com/google/go-cmp/cmp"
)

func TestConvertStepAction(t *testing.T) {
	tests := []struct {
		name           string
		step           *v0.Step
		expectedScript v1.Stringorslice
		expectedEnv    *flexible.Field[map[string]interface{}]
	}{
		{
			name: "basic action step with uses only",
			step: &v0.Step{
				Spec: &v0.StepAction{
					Uses: "actions/checkout@v3",
				},
			},
			expectedScript: v1.Stringorslice{"plugin -kind action -name actions/checkout@v3"},
			expectedEnv:    nil,
		},
		{
			name: "action step with with params prefixed PLUGIN_WITH_",
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
			expectedEnv: &flexible.Field[map[string]interface{}]{Value: map[string]interface{}{
				"PLUGIN_WITH_go-version": "1.21",
				"PLUGIN_WITH_cache":      true,
			}},
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
			expectedEnv: &flexible.Field[map[string]interface{}]{Value: map[string]interface{}{
				"GITHUB_TOKEN": "my-token",
			}},
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
			expectedEnv: &flexible.Field[map[string]interface{}]{Value: map[string]interface{}{
				"PLUGIN_WITH_context": ".",
				"PLUGIN_WITH_push":    true,
				"PLUGIN_WITH_tags":    "myapp:latest",
				"DOCKER_BUILDKIT":     "1",
			}},
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
			expectedEnv:    nil,
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

			if diff := cmp.Diff(tt.expectedEnv, result.Env); diff != "" {
				t.Errorf("Env mismatch (-want +got):\n%s", diff)
			}
		})
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
