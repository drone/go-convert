package converthelpers

import (
	"testing"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	"github.com/google/go-cmp/cmp"
)

func TestConvertStepArtifactoryUpload(t *testing.T) {
	tests := []struct {
		name     string
		step     *v0.Step
		expected map[string]interface{}
	}{
		{
			name: "all required fields",
			step: &v0.Step{
				Spec: &v0.StepArtifactoryUpload{
					ConnRef:    "artifactory-connector",
					Target:     "libs-release-local/com/example/1.0.0",
					SourcePath: "dist/",
				},
			},
			expected: map[string]interface{}{
				"connector": "artifactory-connector",
				"target":    "libs-release-local/com/example/1.0.0",
				"source":    "dist/",
			},
		},
		{
			name: "connector only",
			step: &v0.Step{
				Spec: &v0.StepArtifactoryUpload{
					ConnRef: "artifactory-connector",
				},
			},
			expected: map[string]interface{}{
				"connector": "artifactory-connector",
			},
		},
		{
			name: "minimal empty spec",
			step: &v0.Step{
				Spec: &v0.StepArtifactoryUpload{},
			},
			expected: map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertStepArtifactoryUpload(tt.step)
			if result == nil {
				t.Fatal("expected non-nil result")
			}

			if result.Uses != "uploadArtifactsToJfrogArtifactory" {
				t.Errorf("expected Uses to be uploadArtifactsToJfrogArtifactory, got %s", result.Uses)
			}

			if diff := cmp.Diff(tt.expected, result.With); diff != "" {
				t.Errorf("With mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertStepArtifactoryUpload_NilCases(t *testing.T) {
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
			step: &v0.Step{Spec: nil},
		},
		{
			name: "wrong spec type",
			step: &v0.Step{Spec: &v0.StepRun{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertStepArtifactoryUpload(tt.step)
			if result != nil {
				t.Errorf("expected nil result, got %v", result)
			}
		})
	}
}
