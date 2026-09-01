package converthelpers

import (
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
		expectedEnv    map[string]interface{} // nil if no env expected
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
			name: "action step with with params emitted as per-key PLUGIN_WITH_ env vars",
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
			expectedEnv: map[string]interface{}{
				"PLUGIN_WITH_go-version": "1.21",
				"PLUGIN_WITH_cache":      true,
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
			expectedEnv: map[string]interface{}{
				"PLUGIN_WITH_go-version": "=1.20.1",
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
			expectedEnv: map[string]interface{}{
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
			expectedEnv: map[string]interface{}{
				"PLUGIN_WITH_context": ".",
				"PLUGIN_WITH_push":    true,
				"PLUGIN_WITH_tags":    "myapp:latest",
				"DOCKER_BUILDKIT":     "1",
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
		{
			// CI-24451: a secret expression inside `with` must be emitted as a
			// per-key env var. As a JSON-packed PLUGIN_WITH blob, runtime
			// resolution of a multi-line secret injects raw newlines into the
			// pre-encoded JSON and the plugin's json.Unmarshal fails.
			name: "action step with secret expression in with stays per-key",
			step: &v0.Step{
				Spec: &v0.StepAction{
					Uses: "google-github-actions/auth@v2",
					With: map[string]interface{}{
						"credentials_json": `<+secrets.getValue("gcp-secret")>`,
						"project_id":       "secops-central-sa",
					},
				},
			},
			expectedScript: v1.Stringorslice{"plugin -kind action -name google-github-actions/auth@v2"},
			expectedEnv: map[string]interface{}{
				"PLUGIN_WITH_credentials_json": `<+secrets.getValue("gcp-secret")>`,
				"PLUGIN_WITH_project_id":       "secops-central-sa",
			},
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

			if tt.expectedEnv == nil {
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

			if blob, present := got["PLUGIN_WITH"]; present {
				t.Errorf("unexpected PLUGIN_WITH JSON blob entry: %v", blob)
			}

			if diff := cmp.Diff(tt.expectedEnv, got); diff != "" {
				t.Errorf("Env mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertStepAction_ExpressionsPreservedVerbatim(t *testing.T) {
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

	if got["PLUGIN_WITH_check-latest"] != "<+pipeline.variables.checkLatest>" {
		t.Errorf("expression should be preserved verbatim, got: %v", got["PLUGIN_WITH_check-latest"])
	}
	if got["PLUGIN_WITH_go-version"] != "1.20.1" {
		t.Errorf("expected PLUGIN_WITH_go-version to be 1.20.1, got: %v", got["PLUGIN_WITH_go-version"])
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
