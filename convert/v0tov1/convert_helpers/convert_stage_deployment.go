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
func ConvertDeploymentService(src *v0.DeploymentService) *v1.ServiceRef {
	if src == nil {
		return nil
	}

	// For single service, return simple string reference
	if src.ServiceRef != "" {
		return &v1.ServiceRef{
			Items:      []string{src.ServiceRef},
			MultiService: false,
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
			Items:      serviceRefs,
			MultiService: true,
		}
	}

	return nil
}

func ConvertDeploymentServiceConfig(src *v0.DeploymentServiceConfig) *v1.ServiceRef {
	if src == nil {
		return nil
	}
	if src.ServiceRef != "" {
		return &v1.ServiceRef{
			Items: []string{src.ServiceRef},
		}
	}

	if src.ServiceItem != nil {
		return &v1.ServiceRef{
			Items: []string{src.ServiceItem.Identifier},
		}
	}
	return nil
}

// ConvertEnvironment converts v0 Environment to v1 EnvironmentRef
func ConvertEnvironment(src *v0.Environment) *v1.EnvironmentRef {
	if src == nil {
		return nil
	}

	// Single environment deploying to all infrastructures
	if src.EnvironmentRef != "" {
		var deployTo interface{}
		if infra, ok := src.InfrastructureDefinitions.AsString(); ok {
			deployTo = infra
		} else if infra, ok := src.InfrastructureDefinitions.AsStruct(); ok && len(infra) > 0 {
			deployTo = infra[0].Identifier
		}
		if src.DeployToAll {
			deployTo = "all"
		}
		return &v1.EnvironmentRef{
			Name:     src.EnvironmentRef,
			Id:       src.EnvironmentRef,
			DeployTo: deployTo,
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
			var deployTo interface{}
			if infra, ok := env.InfrastructureDefinitions.AsString(); ok {
				deployTo = infra
			} else if infra, ok := env.InfrastructureDefinitions.AsStruct(); ok {
				infraList := make([]string, 0, len(infra))
				for _, i := range infra {
					infraList = append(infraList, i.Identifier)
				}
				deployTo = infraList
			}
			if env.DeployToAll {
				deployTo = "all"
			}
			items = append(items, &v1.EnvironmentItem{
				Name:     env.EnvironmentRef,
				Id:       env.EnvironmentRef,
				DeployTo: deployTo,
			})
		}
	}

	if len(items) > 0 {
		result := &v1.EnvironmentRef{
			Items: items,
		}
		if src.Metadata != nil {
			result.Sequential = !src.Metadata.Parallel
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
						var deployTo interface{}
						deployTo = "all"
						if deployToAll, ok := envMap["deployToAll"].(bool); ok && deployToAll {
							deployTo = "all"
						}
						if infra, ok := envMap["infrastructureDefinitions"].(map[string]interface{}); ok {
							infraList := make([]string, 0, len(infra))
							for _, i := range infra {
								infraList = append(infraList, i.(string))
							}
							deployTo = infraList
						}

						items = append(items, &v1.EnvironmentItem{
							Name:     envRef,
							Id:       envRef,
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

// ConvertDeploymentInfrastructure converts v0 DeploymentInfrastructure to v1 EnvironmentRef
// It adds the environmentRef from infrastructure to the environments of the stage
func ConvertDeploymentInfrastructure(src *v0.DeploymentInfrastructure) *v1.EnvironmentRef {
	if src == nil || src.EnvironmentRef == "" {
		return nil
	}

	// Create environment item with the environmentRef from infrastructure
	envItem := &v1.EnvironmentItem{
		Name:     src.EnvironmentRef,
		Id:       src.EnvironmentRef,
		DeployTo: "all", // Default to all infrastructures
	}

	// If infrastructure definition is specified, use its identifier
	if src.InfrastructureDefinition.Identifier != "" {
		envItem.DeployTo = src.InfrastructureDefinition.Identifier
	}

	return &v1.EnvironmentRef{
		Items: []*v1.EnvironmentItem{envItem},
	}
}
