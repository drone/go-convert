package converthelpers

import (
	"testing"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/drone/go-convert/internal/flexible"
	"github.com/google/go-cmp/cmp"
)

func TestConvertStepBitrise(t *testing.T) {
	tests := []struct {
		name           string
		step           *v0.Step
		expectedScript v1.Stringorslice
		expectedEnv    *flexible.Field[map[string]interface{}]
	}{
		{
			name: "basic bitrise step with uses only",
			step: &v0.Step{
				Spec: &v0.StepBitrise{
					Uses: "script@1.2.0",
				},
			},
			expectedScript: v1.Stringorslice{"plugin -kind bitrise -name script@1.2.0"},
			expectedEnv:    nil,
		},
		{
			name: "bitrise step with with params",
			step: &v0.Step{
				Spec: &v0.StepBitrise{
					Uses: "xcode-build@2.0.0",
					With: map[string]interface{}{
						"project_path": "MyApp.xcodeproj",
						"scheme":       "Release",
					},
				},
			},
			expectedScript: v1.Stringorslice{"plugin -kind bitrise -name xcode-build@2.0.0"},
			expectedEnv: &flexible.Field[map[string]interface{}]{Value: map[string]interface{}{
				"project_path": "MyApp.xcodeproj",
				"scheme":       "Release",
			}},
		},
		{
			name: "bitrise step with env vars",
			step: &v0.Step{
				Spec: &v0.StepBitrise{
					Uses: "deploy-to-bitrise-io@2.1.0",
					Envs: map[string]string{
						"BITRISE_TOKEN": "my-token",
					},
				},
			},
			expectedScript: v1.Stringorslice{"plugin -kind bitrise -name deploy-to-bitrise-io@2.1.0"},
			expectedEnv: &flexible.Field[map[string]interface{}]{Value: map[string]interface{}{
				"BITRISE_TOKEN": "my-token",
			}},
		},
		{
			name: "bitrise step with both with and env",
			step: &v0.Step{
				Spec: &v0.StepBitrise{
					Uses: "gradle-runner@2.0.0",
					With: map[string]interface{}{
						"gradle_task": "assembleRelease",
					},
					Envs: map[string]string{
						"JAVA_HOME": "/usr/lib/jvm/java-11",
					},
				},
			},
			expectedScript: v1.Stringorslice{"plugin -kind bitrise -name gradle-runner@2.0.0"},
			expectedEnv: &flexible.Field[map[string]interface{}]{Value: map[string]interface{}{
				"gradle_task": "assembleRelease",
				"JAVA_HOME":   "/usr/lib/jvm/java-11",
			}},
		},
		{
			name: "bitrise step with empty with and env",
			step: &v0.Step{
				Spec: &v0.StepBitrise{
					Uses: "cache-pull@2.7.0",
					With: map[string]interface{}{},
					Envs: map[string]string{},
				},
			},
			expectedScript: v1.Stringorslice{"plugin -kind bitrise -name cache-pull@2.7.0"},
			expectedEnv:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertStepBitrise(tt.step)
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

func TestConvertStepBitrise_NilCases(t *testing.T) {
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
			result := ConvertStepBitrise(tt.step)
			if result != nil {
				t.Errorf("expected nil result, got %v", result)
			}
		})
	}
}
