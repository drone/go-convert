package converthelpers

import (
	"testing"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	"github.com/google/go-cmp/cmp"
)

func TestConvertStepGCSUpload(t *testing.T) {
	tests := []struct {
		name     string
		step     *v0.Step
		expected map[string]interface{}
	}{
		{
			name: "all required fields with bucket and target",
			step: &v0.Step{
				Spec: &v0.StepGCSUpload{
					ConnectorRef: "gcp-connector",
					SourcePath:   "dist/",
					Bucket:       "my-bucket",
					Target:       "path/in/bucket",
				},
			},
			expected: map[string]interface{}{
				"connector": "gcp-connector",
				"source":    "dist/",
				"target":    "my-bucket/path/in/bucket",
			},
		},
		{
			name: "bucket only without target",
			step: &v0.Step{
				Spec: &v0.StepGCSUpload{
					ConnectorRef: "gcp-connector",
					SourcePath:   "artifacts/",
					Bucket:       "my-bucket",
				},
			},
			expected: map[string]interface{}{
				"connector": "gcp-connector",
				"source":    "artifacts/",
				"target":    "my-bucket",
			},
		},
		{
			name: "target only without bucket",
			step: &v0.Step{
				Spec: &v0.StepGCSUpload{
					ConnectorRef: "gcp-connector",
					SourcePath:   "artifacts/",
					Target:       "my-bucket/path",
				},
			},
			expected: map[string]interface{}{
				"connector": "gcp-connector",
				"source":    "artifacts/",
				"target":    "my-bucket/path",
			},
		},
		{
			name: "minimal empty spec",
			step: &v0.Step{
				Spec: &v0.StepGCSUpload{},
			},
			expected: map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertStepGCSUpload(tt.step)
			if result == nil {
				t.Fatal("expected non-nil result")
			}

			if result.Uses != "uploadArtifactsToGCS" {
				t.Errorf("expected Uses to be uploadArtifactsToGCS, got %s", result.Uses)
			}

			if diff := cmp.Diff(tt.expected, result.With); diff != "" {
				t.Errorf("With mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertStepGCSUpload_NilCases(t *testing.T) {
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
			result := ConvertStepGCSUpload(tt.step)
			if result != nil {
				t.Errorf("expected nil result, got %v", result)
			}
		})
	}
}
