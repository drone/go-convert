// Copyright 2022 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package converthelpers

import (
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// ConvertDeploymentService converts v0 DeploymentService to v1 ServiceRef
// In v1, service inputs are not required - only the serviceRef is kept
func ConvertDeploymentService(src *v0.DeploymentService) *v1.ServiceRef {
	if src == nil {
		return nil
	}

	// For single service, return simple string reference
	if src.ServiceRef != "" {
		return &v1.ServiceRef{
			Items: []string{src.ServiceRef},
		}
	}

	return nil
}

// ConvertDeploymentServices converts multiple v0 services to v1 ServiceRef
func ConvertDeploymentServices(src *v0.DeploymentServices) *v1.ServiceRef {
	if src == nil || len(src.Values) == 0 {
		return nil
	}

	serviceRefs := make([]string, 0, len(src.Values))
	for _, service := range src.Values {
		if service != nil && service.ServiceRef != "" {
			serviceRefs = append(serviceRefs, service.ServiceRef)
		}
	}

	if len(serviceRefs) > 0 {
		return &v1.ServiceRef{
			Items: serviceRefs,
		}
	}

	return nil
}

func ConvertDeploymentServiceConfig(src *v0.DeploymentServiceConfig) *v1.ServiceRef {
	if src == nil {
		return nil
	}

	return &v1.ServiceRef{
		Items: []string{src.ServiceRef},
	}
}
	

// ConvertEnvironment converts v0 Environment to v1 EnvironmentRef
func ConvertEnvironment(src *v0.Environment) *v1.EnvironmentRef {
	if src == nil {
		return nil
	}

	// Single environment deploying to all infrastructures
	if src.EnvironmentRef != "" {
		deployTo := "all"
		if src.DeployToAll {
			deployTo = "all"
		} else if len(src.InfrastructureDefinitions) == 1 {
			// Single infrastructure
			deployTo = src.InfrastructureDefinitions[0].Identifier
		} else if len(src.InfrastructureDefinitions) > 1 {
			// Multiple infrastructures
			infraIds := make([]string, len(src.InfrastructureDefinitions))
			for i, infra := range src.InfrastructureDefinitions {
				infraIds[i] = infra.Identifier
			}
			return &v1.EnvironmentRef{
				Items: []*v1.EnvironmentItem{
					{
						Name:     src.EnvironmentRef,
						DeployTo: infraIds,
					},
				},
			}
		}

		return &v1.EnvironmentRef{
			Items: []*v1.EnvironmentItem{
				{
					Name:     src.EnvironmentRef,
					DeployTo: deployTo,
				},
			},
		}
	}

	return nil
}

// ConvertEnvironments converts v0 Environments to v1 EnvironmentRef
func ConvertEnvironments(src *v0.Environments) *v1.EnvironmentRef {
	if src == nil || len(src.Values) == 0 {
		return nil
	}

	items := make([]*v1.EnvironmentItem, 0, len(src.Values))
	for _, env := range src.Values {
		if env.EnvironmentRef != "" {
			deployTo := "all"
			if env.DeployToAll {
				deployTo = "all"
			} else if len(env.InfrastructureDefinitions) == 1 {
				deployTo = env.InfrastructureDefinitions[0].Identifier
			} else if len(env.InfrastructureDefinitions) > 1 {
				infraIds := make([]string, len(env.InfrastructureDefinitions))
				for i, infra := range env.InfrastructureDefinitions {
					infraIds[i] = infra.Identifier
				}
				items = append(items, &v1.EnvironmentItem{
					Name:     env.EnvironmentRef,
					DeployTo: infraIds,
				})
				continue
			}

			items = append(items, &v1.EnvironmentItem{
				Name:     env.EnvironmentRef,
				DeployTo: deployTo,
			})
		}
	}

	if len(items) > 0 {
		result := &v1.EnvironmentRef{
			Items: items,
		}
		if src.Metadata != nil {
			result.Parallel = src.Metadata.Parallel
		}
		return result
	}

	return nil
}

// ConvertEnvironmentGroup converts v0 EnvironmentGroup to v1 EnvironmentRef
func ConvertEnvironmentGroup(src *v0.EnvironmentGroup) *v1.EnvironmentRef {
	if src == nil {
		return nil
	}

	// Shorthand reference to environment group
	if src.EnvGroupRef != "" && src.Environments == nil {
		// Simple group reference format: environment: { group: "my-group" }
		return &v1.EnvironmentRef{
			Group: src.EnvGroupRef,
		}
	}

	// Handle explicit group with environments
	if src.Environments != nil {
		switch envs := src.Environments.(type) {
		case []interface{}:
			items := make([]*v1.EnvironmentItem, 0, len(envs))
			for _, env := range envs {
				if envMap, ok := env.(map[string]interface{}); ok {
					if envRef, ok := envMap["environmentRef"].(string); ok {
						deployTo := "all"
						if deployToAll, ok := envMap["deployToAll"].(bool); ok && deployToAll {
							deployTo = "all"
						}
						items = append(items, &v1.EnvironmentItem{
							Name:     envRef,
							DeployTo: deployTo,
						})
					}
				}
			}
			if len(items) > 0 {
				// Explicit group with environments
				groupConfig := map[string]interface{}{
					"name":     src.EnvGroupRef,
					"parallel": false,
					"items":    items,
				}
				if src.Metadata != nil {
					groupConfig["parallel"] = src.Metadata.Parallel
				}

				return &v1.EnvironmentRef{
					Group: groupConfig,
				}
			}
		}
	}

	return nil
}
