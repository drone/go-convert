package converthelpers

import (
	"testing"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	"github.com/google/go-cmp/cmp"
)

func TestConvertStepS3Upload(t *testing.T) {
	tests := []struct {
		name     string
		step     *v0.Step
		expected map[string]interface{}
	}{
		{
			name: "all required fields",
			step: &v0.Step{
				Spec: &v0.StepS3Upload{
					ConnectorRef: "aws-connector",
					Region:       "us-east-1",
					Bucket:       "my-bucket",
					SourcePath:   "dist/",
				},
			},
			expected: map[string]interface{}{
				"connector": "aws-connector",
				"region":    "us-east-1",
				"bucket":    "my-bucket",
				"source":    "dist/",
			},
		},
		{
			name: "with optional endpoint and target",
			step: &v0.Step{
				Spec: &v0.StepS3Upload{
					ConnectorRef: "aws-connector",
					Region:       "us-west-2",
					Bucket:       "my-bucket",
					SourcePath:   "artifacts/",
					Endpoint:     "http://minio.company.com",
					Target:       "path/in/bucket",
				},
			},
			expected: map[string]interface{}{
				"connector": "aws-connector",
				"region":    "us-west-2",
				"bucket":    "my-bucket",
				"source":    "artifacts/",
				"endpoint":  "http://minio.company.com",
				"target":    "path/in/bucket",
			},
		},
		{
			name: "minimal empty spec",
			step: &v0.Step{
				Spec: &v0.StepS3Upload{},
			},
			expected: map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertStepS3Upload(tt.step)
			if result == nil {
				t.Fatal("expected non-nil result")
			}

			if result.Uses != "uploadArtifactsToS3" {
				t.Errorf("expected Uses to be uploadArtifactsToS3, got %s", result.Uses)
			}

			if diff := cmp.Diff(tt.expected, result.With); diff != "" {
				t.Errorf("With mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertStepS3Upload_NilCases(t *testing.T) {
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
			result := ConvertStepS3Upload(tt.step)
			if result != nil {
				t.Errorf("expected nil result, got %v", result)
			}
		})
	}
}
