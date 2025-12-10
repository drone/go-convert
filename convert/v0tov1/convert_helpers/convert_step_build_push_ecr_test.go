package converthelpers

import (
	"testing"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	"github.com/google/go-cmp/cmp"
)

func TestConvertStepBuildAndPushECR(t *testing.T) {
	tests := []struct {
		name     string
		step     *v0.Step
		expected map[string]interface{}
	}{
		{
			name: "basic ECR build and push",
			step: &v0.Step{
				Spec: &v0.StepBuildAndPushECR{
					ConnectorRef: "aws-connector",
					Region:       "us-east-1",
					Account:      "123456789012",
					ImageName:    "my-app",
					Tags:         []string{"latest", "v1.0.0"},
				},
			},
			expected: map[string]interface{}{
				"connector": "aws-connector",
				"region":    "us-east-1",
				"registry":  "123456789012.dkr.ecr.us-east-1.amazonaws.com",
				"repo":      "my-app",
				"tags":      []string{"latest", "v1.0.0"},
			},
		},
		{
			name: "BaseImageConnectorRefs with single string value",
			step: &v0.Step{
				Spec: &v0.StepBuildAndPushECR{
					ConnectorRef:           "aws-connector",
					Region:                 "eu-west-1",
					Account:                "111222333444",
					ImageName:              "frontend",
					BaseImageConnectorRefs: "docker-hub-connector",
				},
			},
			expected: map[string]interface{}{
				"connector":          "aws-connector",
				"region":             "eu-west-1",
				"registry":           "111222333444.dkr.ecr.eu-west-1.amazonaws.com",
				"repo":               "frontend",
				"baseimageconnector": "docker-hub-connector",
			},
		},
		{
			name: "BaseImageConnectorRefs with multiple values array",
			step: &v0.Step{
				Spec: &v0.StepBuildAndPushECR{
					ConnectorRef: "aws-connector",
					Region:       "ap-south-1",
					Account:      "555666777888",
					ImageName:    "multi-base-app",
					BaseImageConnectorRefs: []interface{}{
						"docker-hub-connector",
						"gcr-connector",
						"private-registry-connector",
					},
				},
			},
			expected: map[string]interface{}{
				"connector": "aws-connector",
				"region":    "ap-south-1",
				"registry":  "555666777888.dkr.ecr.ap-south-1.amazonaws.com",
				"repo":      "multi-base-app",
				"baseimageconnector": []interface{}{
					"docker-hub-connector",
					"gcr-connector",
					"private-registry-connector",
				},
			},
		},
		{
			name: "registry URL construction with missing account",
			step: &v0.Step{
				Spec: &v0.StepBuildAndPushECR{
					ConnectorRef: "aws-connector",
					Region:       "us-east-1",
					Account:      "",
					ImageName:    "my-app",
				},
			},
			expected: map[string]interface{}{
				"connector": "aws-connector",
				"region":    "us-east-1",
				"repo":      "my-app",
			},
		},
		{
			name: "registry URL construction with missing region",
			step: &v0.Step{
				Spec: &v0.StepBuildAndPushECR{
					ConnectorRef: "aws-connector",
					Region:       "",
					Account:      "123456789012",
					ImageName:    "my-app",
				},
			},
			expected: map[string]interface{}{
				"connector": "aws-connector",
				"repo":      "my-app",
			},
		},
		{
			name: "registry URL construction with both missing",
			step: &v0.Step{
				Spec: &v0.StepBuildAndPushECR{
					ConnectorRef: "aws-connector",
					Region:       "",
					Account:      "",
					ImageName:    "my-app",
				},
			},
			expected: map[string]interface{}{
				"connector": "aws-connector",
				"repo":      "my-app",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertStepBuildAndPushECR(tt.step)
			if result == nil {
				t.Fatal("expected non-nil result")
			}

			if result.Uses != "buildAndPushToECR" {
				t.Errorf("expected Uses to be buildAndPushToECR, got %s", result.Uses)
			}

			if diff := cmp.Diff(tt.expected, result.With); diff != "" {
				t.Errorf("With mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertStepBuildAndPushECR_NilCases(t *testing.T) {
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
			result := ConvertStepBuildAndPushECR(tt.step)
			if result != nil {
				t.Errorf("expected nil result, got %v", result)
			}
		})
	}
}
