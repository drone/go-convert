package converthelpers

import (
	"testing"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	"github.com/drone/go-convert/internal/flexible"
	"github.com/google/go-cmp/cmp"
)

func TestConvertDeploymentService_ServiceInputsWithoutSpec(t *testing.T) {
	tests := []struct {
		name          string
		serviceRef    string
		serviceInputs *flexible.Field[v0.ServiceInputs]
		wantWith      map[string]interface{}
	}{
		{
			name:       "serviceDefinition with only type and no spec yields empty with",
			serviceRef: "svc1",
			serviceInputs: func() *flexible.Field[v0.ServiceInputs] {
				f := &flexible.Field[v0.ServiceInputs]{}
				f.Set(v0.ServiceInputs{
					ServiceDefinition: map[string]interface{}{
						"type": "Kubernetes",
					},
				})
				return f
			}(),
			wantWith: nil,
		},
		{
			name:       "serviceDefinition with type and spec yields overlay",
			serviceRef: "svc1",
			serviceInputs: func() *flexible.Field[v0.ServiceInputs] {
				f := &flexible.Field[v0.ServiceInputs]{}
				f.Set(v0.ServiceInputs{
					ServiceDefinition: map[string]interface{}{
						"type": "Kubernetes",
						"spec": map[string]interface{}{
							"manifests": []interface{}{},
						},
					},
				})
				return f
			}(),
			wantWith: map[string]interface{}{
				"overlay": map[string]interface{}{
					"type": "Kubernetes",
					"spec": map[string]interface{}{
						"manifests": []interface{}{},
					},
				},
			},
		},
		{
			name:       "serviceDefinition as expression passes through",
			serviceRef: "svc1",
			serviceInputs: func() *flexible.Field[v0.ServiceInputs] {
				f := &flexible.Field[v0.ServiceInputs]{}
				f.SetExpression("<+input>")
				return f
			}(),
			wantWith: map[string]interface{}{
				"overlay": "<+input>",
			},
		},
		{
			name:          "nil serviceInputs yields empty with",
			serviceRef:    "svc1",
			serviceInputs: nil,
			wantWith:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := &v0.DeploymentService{
				ServiceRef:    tt.serviceRef,
				ServiceInputs: tt.serviceInputs,
			}
			got := ConvertDeploymentService(src, NewStageConversionContext())
			if got == nil {
				t.Fatal("expected non-nil ServiceRef")
			}
			items := got.ItemList()
			if len(items) != 1 {
				t.Fatalf("expected 1 service item, got %d", len(items))
			}
			if diff := cmp.Diff(tt.wantWith, items[0].With); diff != "" {
				t.Errorf("With mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertEnvironmentOverrides(t *testing.T) {
	tests := []struct {
		name     string
		env      *v0.Environment
		expected map[string]interface{}
	}{
		{
			name: "concrete environmentInputs and serviceOverrideInputs",
			env: &v0.Environment{
				EnvironmentRef: "qa",
				DeployToAll:    &flexible.Field[bool]{Value: false},
				InfrastructureDefinitions: &flexible.Field[[]*v0.InfrastructureDefinition]{
					Value: []*v0.InfrastructureDefinition{{Identifier: "k8s_gcp"}},
				},
				EnvironmentInputs: map[string]interface{}{
					"variables": []interface{}{
						map[string]interface{}{
							"name":  "env_global_variable",
							"type":  "String",
							"value": "<+input>",
						},
					},
				},
				ServiceOverrideInputs: map[string]interface{}{
					"variables": []interface{}{
						map[string]interface{}{
							"name":  "env_service_variable",
							"type":  "String",
							"value": "<+input>",
						},
					},
				},
			},
			expected: map[string]interface{}{
				"global": map[string]interface{}{
					"with": map[string]interface{}{
						"overlay": map[string]interface{}{
							"variables": []interface{}{
								map[string]interface{}{
									"name":  "env_global_variable",
									"type":  "String",
									"value": "<+input>",
								},
							},
						},
					},
				},
				"service": map[string]interface{}{
					"with": map[string]interface{}{
						"overlay": map[string]interface{}{
							"variables": []interface{}{
								map[string]interface{}{
									"name":  "env_service_variable",
									"type":  "String",
									"value": "<+input>",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "expression inputs are set directly",
			env: &v0.Environment{
				EnvironmentRef:        "qa",
				EnvironmentInputs:     "<+input>",
				ServiceOverrideInputs: "<+input>",
			},
			expected: map[string]interface{}{
				"global":  "<+input>",
				"service": "<+input>",
			},
		},
		{
			name: "only environmentInputs present",
			env: &v0.Environment{
				EnvironmentRef:    "qa",
				EnvironmentInputs: "<+input>",
			},
			expected: map[string]interface{}{
				"global": "<+input>",
			},
		},
		{
			name: "no inputs yields no overrides",
			env: &v0.Environment{
				EnvironmentRef: "qa",
			},
			expected: nil,
		},
		{
			name: "empty string inputs yield no overrides",
			env: &v0.Environment{
				EnvironmentRef:        "qa",
				EnvironmentInputs:     "",
				ServiceOverrideInputs: "",
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertEnvironment(tt.env, NewStageConversionContext())
			items := got.ItemList()
			if got == nil || len(items) != 1 {
				t.Fatalf("expected a single environment item, got %#v", got)
			}
			if diff := cmp.Diff(tt.expected, items[0].Overrides); diff != "" {
				t.Errorf("overrides mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// TestConvertEnvironment_FullyRuntimeInput verifies that when every non-null
// field on a single v0 Environment is the literal "<+input>" runtime input
// expression, the resulting EnvironmentRef collapses its items list into a
// single runtime input expression instead of a one-item list.
func TestConvertEnvironment_FullyRuntimeInput(t *testing.T) {
	tests := []struct {
		name          string
		env           *v0.Environment
		wantItemsExpr bool
	}{
		{
			name: "environmentRef, deployToAll and environmentInputs all <+input>",
			env: &v0.Environment{
				EnvironmentRef:    "<+input>",
				DeployToAll:       &flexible.Field[bool]{Value: "<+input>"},
				EnvironmentInputs: "<+input>",
			},
			wantItemsExpr: true,
		},
		{
			name: "environmentRef concrete - must not collapse",
			env: &v0.Environment{
				EnvironmentRef:    "qa",
				DeployToAll:       &flexible.Field[bool]{Value: "<+input>"},
				EnvironmentInputs: "<+input>",
			},
			wantItemsExpr: false,
		},
		{
			name: "infrastructureDefinitions concrete - must not collapse",
			env: &v0.Environment{
				EnvironmentRef: "<+input>",
				InfrastructureDefinitions: &flexible.Field[[]*v0.InfrastructureDefinition]{
					Value: []*v0.InfrastructureDefinition{{Identifier: "dev_infra"}},
				},
			},
			wantItemsExpr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertEnvironment(tt.env, NewStageConversionContext())
			if got == nil {
				t.Fatal("expected non-nil EnvironmentRef")
			}
			expr, isExpr := got.Items.AsString()
			if tt.wantItemsExpr {
				if !isExpr || expr != runtimeInputExpr {
					t.Errorf("expected items to be the runtime input expression, got %#v", got.Items.Value)
				}
			} else if isExpr {
				t.Errorf("did not expect items to collapse to an expression, got %q", expr)
			}
		})
	}
}
