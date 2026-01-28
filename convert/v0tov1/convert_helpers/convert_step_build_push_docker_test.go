package converthelpers

import (
	"testing"
	"github.com/drone/go-convert/internal/flexible"
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	"github.com/google/go-cmp/cmp"
)

func TestConvertStepBuildAndPushDockerRegistry(t *testing.T) {
	tests := []struct {
		name     string
		step     *v0.Step
		expected map[string]interface{}
	}{
		{
			name: "basic Docker build and push",
			step: &v0.Step{
				Spec: &v0.StepBuildAndPushDockerRegistry{
					ConnectorRef: "docker-connector",
					Repo:         "myorg/myapp",
					Tags:         []string{"latest", "v1.0.0"},
				},
			},
			expected: map[string]interface{}{
				"connector": "docker-connector",
				"repo":      "myorg/myapp",
				"tags":      []string{"latest", "v1.0.0"},
			},
		},
		{
			name: "BaseImageConnectorRefs with single string value",
			step: &v0.Step{
				Spec: &v0.StepBuildAndPushDockerRegistry{
					ConnectorRef:           "docker-connector",
					Repo:                   "myorg/frontend",
					BaseImageConnectorRefs: "gcr-connector",
				},
			},
			expected: map[string]interface{}{
				"connector":          "docker-connector",
				"repo":               "myorg/frontend",
				"baseimageconnector": "gcr-connector",
			},
		},
		{
			name: "BaseImageConnectorRefs with multiple values array takes first",
			step: &v0.Step{
				Spec: &v0.StepBuildAndPushDockerRegistry{
					ConnectorRef: "docker-connector",
					Repo:         "myorg/multi-base",
					BaseImageConnectorRefs: []interface{}{
						"primary-connector",
						"secondary-connector",
						"tertiary-connector",
					},
				},
			},
			expected: map[string]interface{}{
				"connector":          "docker-connector",
				"repo":               "myorg/multi-base",
				"baseimageconnector": "primary-connector",
			},
		},
		{
			name: "BaseImageConnectorRefs with string array takes first",
			step: &v0.Step{
				Spec: &v0.StepBuildAndPushDockerRegistry{
					ConnectorRef: "docker-connector",
					Repo:         "myorg/app",
					BaseImageConnectorRefs: []string{
						"first-connector",
						"second-connector",
					},
				},
			},
			expected: map[string]interface{}{
				"connector":          "docker-connector",
				"repo":               "myorg/app",
				"baseimageconnector": "first-connector",
			},
		},
		{
			name: "BaseImageConnectorRefs with empty array",
			step: &v0.Step{
				Spec: &v0.StepBuildAndPushDockerRegistry{
					ConnectorRef:           "docker-connector",
					Repo:                   "myorg/app",
					BaseImageConnectorRefs: []interface{}{},
				},
			},
			expected: map[string]interface{}{
				"connector": "docker-connector",
				"repo":      "myorg/app",
			},
		},
		{
			name: "with all optional fields",
			step: &v0.Step{
				Spec: &v0.StepBuildAndPushDockerRegistry{
					ConnectorRef: "docker-connector",
					Repo:         "myorg/complete",
					Tags:         []string{"v2.0.0"},
					Caching:      &flexible.Field[bool]{Value: true},
					Dockerfile:   "Dockerfile.prod",
					Context:      "./backend",
					Target:       "production",
					Labels: map[string]string{
						"version": "2.0.0",
					},
					BuildArgs: map[string]string{
						"NODE_VERSION": "20",
					},
					Env: map[string]string{
						"BUILD_ENV": "prod",
					},
				},
			},
			expected: map[string]interface{}{
				"connector":  "docker-connector",
				"repo":       "myorg/complete",
				"tags":       []string{"v2.0.0"},
				"caching":    &flexible.Field[bool]{Value: true},
				"dockerfile": "Dockerfile.prod",
				"context":    "./backend",
				"target":     "production",
				"labels": map[string]string{
					"version": "2.0.0",
				},
				"buildargs": map[string]string{
					"NODE_VERSION": "20",
				},
				"envvars": map[string]string{
					"BUILD_ENV": "prod",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertStepBuildAndPushDockerRegistry(tt.step)
			if result == nil {
				t.Fatal("expected non-nil result")
			}

			if result.Uses != "buildAndPushToDocker" {
				t.Errorf("expected Uses to be buildAndPushToDocker, got %s", result.Uses)
			}

			if diff := cmp.Diff(tt.expected, result.With); diff != "" {
				t.Errorf("With mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertStepBuildAndPushDockerRegistry_NilCases(t *testing.T) {
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
			result := ConvertStepBuildAndPushDockerRegistry(tt.step)
			if result != nil {
				t.Errorf("expected nil result, got %v", result)
			}
		})
	}
}
