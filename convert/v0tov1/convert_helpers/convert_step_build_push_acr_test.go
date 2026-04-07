package converthelpers

import (
	"testing"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	"github.com/drone/go-convert/internal/flexible"
	"github.com/google/go-cmp/cmp"
)

func TestConvertStepBuildAndPushACR(t *testing.T) {
	tests := []struct {
		name     string
		step     *v0.Step
		expected map[string]interface{}
	}{
		{
			name: "basic ACR build and push with registry/image split",
			step: &v0.Step{
				Spec: &v0.StepBuildAndPushACR{
					ConnectorRef:   "azure-connector",
					Repository:     "myregistry.azurecr.io/myapp",
					SubscriptionId: "sub-123",
					Tags:           &flexible.Field[[]string]{Value: []string{"latest", "v1.0.0"}},
				},
			},
			expected: map[string]interface{}{
				"connector":      "azure-connector",
				"registry":       "myregistry.azurecr.io",
				"image_name":     "myapp",
				"subscriptionid": "sub-123",
				"tags":           &flexible.Field[[]string]{Value: []string{"latest", "v1.0.0"}},
				"caching":        true,
				"build_mode":     "build_and_push",
			},
		},
		{
			name: "repository with nested path",
			step: &v0.Step{
				Spec: &v0.StepBuildAndPushACR{
					ConnectorRef: "azure-connector",
					Repository:   "prodregistry.azurecr.io/team/backend-service",
				},
			},
			expected: map[string]interface{}{
				"connector":  "azure-connector",
				"registry":   "prodregistry.azurecr.io",
				"image_name": "team/backend-service",
				"caching":    true,
				"build_mode": "build_and_push",
			},
		},
		{
			name: "repository without slash (no registry)",
			step: &v0.Step{
				Spec: &v0.StepBuildAndPushACR{
					ConnectorRef: "azure-connector",
					Repository:   "myapp",
				},
			},
			expected: map[string]interface{}{
				"connector":  "azure-connector",
				"image_name": "myapp",
				"caching":    true,
				"build_mode": "build_and_push",
			},
		},
		{
			name: "with all optional fields",
			step: &v0.Step{
				Spec: &v0.StepBuildAndPushACR{
					ConnectorRef:   "azure-connector",
					Repository:     "complete.azurecr.io/myapp",
					SubscriptionId: "sub-456",
					Tags:           &flexible.Field[[]string]{Value: []string{"v2.0.0"}},
					Caching:        &flexible.Field[bool]{Value: true},
					Dockerfile:     "Dockerfile.prod",
					Context:        "./app",
					Target:         "production",
					Labels: &flexible.Field[map[string]string]{Value: map[string]string{
						"version": "2.0.0",
					}},
					BuildArgs: &flexible.Field[map[string]string]{Value: map[string]string{
						"NODE_VERSION": "20",
					}},
					Env: &flexible.Field[map[string]string]{Value: map[string]string{
						"BUILD_ENV": "prod",
					}},
					BaseImageConnectorRefs: "docker-connector",
				},
			},
			expected: map[string]interface{}{
				"connector":          "azure-connector",
				"registry":           "complete.azurecr.io",
				"image_name":         "myapp",
				"subscriptionid":     "sub-456",
				"tags":               &flexible.Field[[]string]{Value: []string{"v2.0.0"}},
				"caching":            &flexible.Field[bool]{Value: true},
				"build_mode":         "build_and_push",
				"dockerfile":         "Dockerfile.prod",
				"context":            "./app",
				"target":             "production",
				"baseimageconnector": "docker-connector",
				"labels": &flexible.Field[map[string]string]{Value: map[string]string{
					"version": "2.0.0",
				}},
				"buildargs": &flexible.Field[map[string]string]{Value: map[string]string{
					"NODE_VERSION": "20",
				}},
				"envvars": &flexible.Field[map[string]string]{Value: map[string]string{
					"BUILD_ENV": "prod",
				}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertStepBuildAndPushACR(tt.step)
			if result == nil {
				t.Fatal("expected non-nil result")
			}

			if result.Uses != "buildAndPushToACR" {
				t.Errorf("expected Uses to be buildAndPushToACR, got %s", result.Uses)
			}

			if diff := cmp.Diff(tt.expected, result.With); diff != "" {
				t.Errorf("With mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertStepBuildAndPushACR_NilCases(t *testing.T) {
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
			result := ConvertStepBuildAndPushACR(tt.step)
			if result != nil {
				t.Errorf("expected nil result, got %v", result)
			}
		})
	}
}
