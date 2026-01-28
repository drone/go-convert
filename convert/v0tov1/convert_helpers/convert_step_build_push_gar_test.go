package converthelpers

import (
	"testing"
	"github.com/drone/go-convert/internal/flexible"
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	"github.com/google/go-cmp/cmp"
)

func TestConvertStepBuildAndPushGAR(t *testing.T) {
	tests := []struct {
		name     string
		step     *v0.Step
		expected map[string]interface{}
	}{
		{
			name: "basic GAR build and push",
			step: &v0.Step{
				Spec: &v0.StepBuildAndPushGAR{
					ConnectorRef: "gcp-connector",
					Host:         "us-central1-docker.pkg.dev",
					ProjectID:    "my-gcp-project",
					ImageName:    "my-app",
					Tags:         []string{"latest", "v1.0.0"},
				},
			},
			expected: map[string]interface{}{
				"connector": "gcp-connector",
				"registry":  "us-central1-docker.pkg.dev/my-gcp-project",
				"repo":      "my-app",
				"tags":      []string{"latest", "v1.0.0"},
			},
		},
		{
			name: "with caching and build args",
			step: &v0.Step{
				Spec: &v0.StepBuildAndPushGAR{
					ConnectorRef: "gcp-connector",
					Host:         "europe-west1-docker.pkg.dev",
					ProjectID:    "prod-project",
					ImageName:    "backend-service",
					Caching:      &flexible.Field[bool]{Value: true},
					BuildArgs: map[string]string{
						"GO_VERSION": "1.21",
						"APP_ENV":    "production",
					},
				},
			},
			expected: map[string]interface{}{
				"connector": "gcp-connector",
				"registry":  "europe-west1-docker.pkg.dev/prod-project",
				"repo":      "backend-service",
				"caching":   &flexible.Field[bool]{Value: true},
				"buildargs": map[string]string{
					"GO_VERSION": "1.21",
					"APP_ENV":    "production",
				},
			},
		},
		{
			name: "BaseImageConnectorRefs with single string value",
			step: &v0.Step{
				Spec: &v0.StepBuildAndPushGAR{
					ConnectorRef:           "gcp-connector",
					Host:                   "asia-south1-docker.pkg.dev",
					ProjectID:              "dev-project",
					ImageName:              "frontend",
					BaseImageConnectorRefs: "docker-hub-connector",
				},
			},
			expected: map[string]interface{}{
				"connector":          "gcp-connector",
				"registry":           "asia-south1-docker.pkg.dev/dev-project",
				"repo":               "frontend",
				"baseimageconnector": "docker-hub-connector",
			},
		},
		{
			name: "BaseImageConnectorRefs with multiple values array",
			step: &v0.Step{
				Spec: &v0.StepBuildAndPushGAR{
					ConnectorRef: "gcp-connector",
					Host:         "us-west1-docker.pkg.dev",
					ProjectID:    "test-project",
					ImageName:    "multi-base-app",
					BaseImageConnectorRefs: []interface{}{
						"primary-connector",
						"secondary-connector",
						"tertiary-connector",
					},
				},
			},
			expected: map[string]interface{}{
				"connector": "gcp-connector",
				"registry":  "us-west1-docker.pkg.dev/test-project",
				"repo":      "multi-base-app",
				"baseimageconnector": []interface{}{
					"primary-connector",
					"secondary-connector",
					"tertiary-connector",
				},
			},
		},
		{
			name: "BaseImageConnectorRefs with map value",
			step: &v0.Step{
				Spec: &v0.StepBuildAndPushGAR{
					ConnectorRef: "gcp-connector",
					Host:         "us-east1-docker.pkg.dev",
					ProjectID:    "complex-project",
					ImageName:    "complex-app",
					BaseImageConnectorRefs: map[string]interface{}{
						"base":    "docker-connector",
						"builder": "gcr-connector",
					},
				},
			},
			expected: map[string]interface{}{
				"connector": "gcp-connector",
				"registry":  "us-east1-docker.pkg.dev/complex-project",
				"repo":      "complex-app",
				"baseimageconnector": map[string]interface{}{
					"base":    "docker-connector",
					"builder": "gcr-connector",
				},
			},
		},
		{
			name: "registry URL construction with missing host",
			step: &v0.Step{
				Spec: &v0.StepBuildAndPushGAR{
					ConnectorRef: "gcp-connector",
					Host:         "",
					ProjectID:    "my-project",
					ImageName:    "my-app",
				},
			},
			expected: map[string]interface{}{
				"connector": "gcp-connector",
				"repo":      "my-app",
			},
		},
		{
			name: "registry URL construction with missing projectID",
			step: &v0.Step{
				Spec: &v0.StepBuildAndPushGAR{
					ConnectorRef: "gcp-connector",
					Host:         "us-central1-docker.pkg.dev",
					ProjectID:    "",
					ImageName:    "my-app",
				},
			},
			expected: map[string]interface{}{
				"connector": "gcp-connector",
				"repo":      "my-app",
			},
		},
		{
			name: "registry URL construction with both missing",
			step: &v0.Step{
				Spec: &v0.StepBuildAndPushGAR{
					ConnectorRef: "gcp-connector",
					Host:         "",
					ProjectID:    "",
					ImageName:    "my-app",
				},
			},
			expected: map[string]interface{}{
				"connector": "gcp-connector",
				"repo":      "my-app",
			},
		},
		{
			name: "with all optional fields",
			step: &v0.Step{
				Spec: &v0.StepBuildAndPushGAR{
					ConnectorRef: "gcp-connector",
					Host:         "us-central1-docker.pkg.dev",
					ProjectID:    "complete-project",
					ImageName:    "complete-app",
					Tags:         []string{"v2.0.0"},
					Caching:      &flexible.Field[bool]{Value: true},
					Dockerfile:   "Dockerfile.prod",
					Context:      "./app",
					Target:       "production",
					Labels: map[string]string{
						"version": "2.0.0",
						"team":    "platform",
					},
					BuildArgs: map[string]string{
						"PYTHON_VERSION": "3.11",
					},
					Env: map[string]string{
						"BUILD_ENV": "prod",
						"LOG_LEVEL": "info",
					},
				},
			},
			expected: map[string]interface{}{
				"connector":  "gcp-connector",
				"registry":   "us-central1-docker.pkg.dev/complete-project",
				"repo":       "complete-app",
				"tags":       []string{"v2.0.0"},
				"caching":    &flexible.Field[bool]{Value: true},
				"dockerfile": "Dockerfile.prod",
				"context":    "./app",
				"target":     "production",
				"labels": map[string]string{
					"version": "2.0.0",
					"team":    "platform",
				},
				"buildargs": map[string]string{
					"PYTHON_VERSION": "3.11",
				},
				"envvars": map[string]string{
					"BUILD_ENV": "prod",
					"LOG_LEVEL": "info",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertStepBuildAndPushGAR(tt.step)
			if result == nil {
				t.Fatal("expected non-nil result")
			}

			if result.Uses != "buildAndPushToGAR" {
				t.Errorf("expected Uses to be buildAndPushToGAR, got %s", result.Uses)
			}

			if diff := cmp.Diff(tt.expected, result.With); diff != "" {
				t.Errorf("With mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertStepBuildAndPushGAR_NilCases(t *testing.T) {
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
			result := ConvertStepBuildAndPushGAR(tt.step)
			if result != nil {
				t.Errorf("expected nil result, got %v", result)
			}
		})
	}
}
