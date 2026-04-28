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
	"fmt"
	"log"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/drone/go-convert/internal/flexible"
)

// ConvertDeploymentService converts v0 DeploymentService to v1 ServiceRef
func ConvertDeploymentService(src *v0.DeploymentService, ctx *StageConversionContext) *v1.ServiceRef {
	if src == nil {
		return nil
	}

	// Handle useFromStage
	if src.UseFromStage != nil && src.UseFromStage.Stage != "" {
		if ref := ctx.GetService(src.UseFromStage.Stage); ref != nil {
			return ref
		}
		log.Printf("Warning!!! service useFromStage '%s' not found in context\n", src.UseFromStage.Stage)
		return nil
	}

	// For single service, return simple string reference
	if src.ServiceRef != "" {
		return &v1.ServiceRef{
			Items:        []string{src.ServiceRef},
			MultiService: false,
		}
	}

	return nil
}

// ConvertDeploymentServices converts multiple v0 services to v1 ServiceRef
func ConvertDeploymentServices(src *v0.DeploymentServices, ctx *StageConversionContext) *v1.ServiceRef {
	if src == nil {
		return nil
	}
	// Handle useFromStage
	if src.UseFromStage != nil && src.UseFromStage.Stage != "" {
		if ref := ctx.GetService(src.UseFromStage.Stage); ref != nil {
			return ref
		}
		log.Printf("Warning!!! services useFromStage '%s' not found in context\n", src.UseFromStage.Stage)
		return nil
	}

	// Handle Values as flexible.Field
	if src.Values == nil {
		return nil
	}

	// Check if Values is an expression
	if expr, ok := src.Values.AsString(); ok {
		// If expression is not <+input> or empty string, return nil and log
		if expr != "<+input>" && expr != "" {
			log.Printf("Warning: services.values contains unsupported expression '%s', skipping conversion\n", expr)
			return nil
		}
		// For <+input> or empty string, continue with nil values (no services)
		return nil
	}

	// Values is a struct (array of services)
	values, ok := src.Values.AsStruct()
	if !ok || len(values) == 0 {
		return nil
	}

	serviceRefs := make([]string, 0, len(values))
	for _, service := range values {
		if service != nil && service.ServiceRef != "" {
			serviceRefs = append(serviceRefs, service.ServiceRef)
		}
	}
	// by default set to true
	var sequential *flexible.Field[bool] = &flexible.Field[bool]{Value: true}
	if src.Metadata != nil && src.Metadata.Parallel != nil {
		sequential = flexible.NegateBool(src.Metadata.Parallel)
	}
	if len(serviceRefs) > 0 {
		return &v1.ServiceRef{
			Items:        serviceRefs,
			MultiService: true,
			Sequential:   sequential,
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
func ConvertEnvironment(src *v0.Environment, ctx *StageConversionContext) *v1.EnvironmentRef {
	if src == nil {
		return nil
	}
	// Handle useFromStage
	if src.UseFromStage != nil && src.UseFromStage.Stage != "" {
		if ref := ctx.GetEnvironment(src.UseFromStage.Stage); ref != nil {
			return ref
		}
		log.Printf("Warning!!! environment useFromStage '%s' not found in context\n", src.UseFromStage.Stage)
		return nil
	}

	// Single environment deploying to all infrastructures
	if src.EnvironmentRef != "" {
		deployTo := resolveDeployTo(src.DeployToAll, &src.InfrastructureDefinitions)
		return &v1.EnvironmentRef{
			Name:     src.EnvironmentRef,
			Id:       src.EnvironmentRef,
			DeployTo: deployTo,
		}
	}

	return nil
}

// ConvertEnvironments converts v0 Environments to v1 EnvironmentRef
func ConvertEnvironments(src *v0.Environments, ctx *StageConversionContext) *v1.EnvironmentRef {
	if src == nil {
		return nil
	}
	// Handle useFromStage
	if src.UseFromStage != nil && src.UseFromStage.Stage != "" {
		if ref := ctx.GetEnvironment(src.UseFromStage.Stage); ref != nil {
			return ref
		}
		log.Printf("Warning!!! environments useFromStage '%s' not found in context\n", src.UseFromStage.Stage)
		return nil
	}

	if len(src.Values) == 0 {
		return nil
	}

	// by default set to true
	var sequential *flexible.Field[bool] = &flexible.Field[bool]{Value: true}
	if src.Metadata != nil && src.Metadata.Parallel != nil {
		sequential = flexible.NegateBool(src.Metadata.Parallel)
	}

	items := make([]*v1.EnvironmentItem, 0, len(src.Values))
	for _, env := range src.Values {
		if env.EnvironmentRef != "" {
			deployTo := resolveDeployTo(env.DeployToAll, &env.InfrastructureDefinitions)
			items = append(items, &v1.EnvironmentItem{
				Name:     env.EnvironmentRef,
				Id:       env.EnvironmentRef,
				DeployTo: deployTo,
			})
		}
	}

	if len(items) > 0 {
		result := &v1.EnvironmentRef{
			Items:      items,
			Sequential: sequential,
		}
		return result
	}

	return nil
}

// ConvertEnvironmentGroup converts v0 EnvironmentGroup to v1 EnvironmentRef
func ConvertEnvironmentGroup(src *v0.EnvironmentGroup, ctx *StageConversionContext) *v1.EnvironmentRef {
	if src == nil {
		return nil
	}

	// Handle useFromStage
	if src.UseFromStage != nil && src.UseFromStage.Stage != "" {
		if ref := ctx.GetEnvironment(src.UseFromStage.Stage); ref != nil {
			return ref
		}
		log.Printf("Warning!!! environmentGroup useFromStage '%s' not found in context\n", src.UseFromStage.Stage)
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
		items, sequential := parseEnvironmentGroupEnvironments(src.Environments, src.Metadata)
		if len(items) > 0 {
			groupConfig := map[string]interface{}{
				"name":       src.EnvGroupRef,
				"sequential": sequential,
				"items":      items,
			}
			return &v1.EnvironmentRef{
				Group: groupConfig,
			}
		}
	}

	return nil
}

// parseEnvironmentGroupEnvironments parses the environments field which can be:
// 1. []interface{} - direct array of environment objects
// 2. map[string]interface{} with "metadata" and "values" keys
func parseEnvironmentGroupEnvironments(environments interface{}, groupMetadata *v0.EnvironmentMetadata) ([]*v1.EnvironmentItem, bool) {
	// Default sequential to true (parallel: false)
	sequential := true

	// Check group-level metadata first
	if groupMetadata != nil && groupMetadata.Parallel != nil {
		if val, ok := groupMetadata.Parallel.AsStruct(); ok {
			sequential = !val
		}
	}

	var envList []interface{}

	switch envs := environments.(type) {
	case []interface{}:
		// Direct array of environments
		envList = envs

	case map[string]interface{}:
		// Object with metadata and values
		if metadata, ok := envs["metadata"].(map[string]interface{}); ok {
			if parallel, ok := metadata["parallel"].(bool); ok {
				sequential = !parallel
			}
		}
		if values, ok := envs["values"].([]interface{}); ok {
			envList = values
		}

	default:
		return nil, sequential
	}

	items := make([]*v1.EnvironmentItem, 0, len(envList))
	for _, env := range envList {
		if envMap, ok := env.(map[string]interface{}); ok {
			envRef, _ := envMap["environmentRef"].(string)
			if envRef == "" {
				continue
			}

			deployTo := extractDeployTo(envMap)
			items = append(items, &v1.EnvironmentItem{
				Name:     envRef,
				Id:       envRef,
				DeployTo: deployTo,
			})
		}
	}

	return items, sequential
}

// extractDeployTo extracts the deploy-to value from an environment map
func extractDeployTo(envMap map[string]interface{}) interface{} {
	// Check deployToAll first
	if deployToAll, ok := envMap["deployToAll"].(bool); ok && deployToAll {
		return "all"
	}

	// Check infrastructureDefinitions
	if infraDefs, ok := envMap["infrastructureDefinitions"].([]interface{}); ok {
		if len(infraDefs) == 1 {
			if infraMap, ok := infraDefs[0].(map[string]interface{}); ok {
				if id, ok := infraMap["identifier"].(string); ok {
					return id
				}
			}
		} else if len(infraDefs) > 1 {
			infraList := make([]string, 0, len(infraDefs))
			for _, infra := range infraDefs {
				if infraMap, ok := infra.(map[string]interface{}); ok {
					if id, ok := infraMap["identifier"].(string); ok {
						infraList = append(infraList, id)
					}
				}
			}
			if len(infraList) > 0 {
				return infraList
			}
		}
	}

	return "all"
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

// resolveDeployTo determines the deploy-to value based on DeployToAll and infrastructure definitions.
// - If DeployToAll is a boolean true → "all"
// - If DeployToAll is <+input> → returns the infra value only (single string or list)
// - If DeployToAll is any other expression → builds ternary: <+ <+expr> ? "all" : "infraValue" >
// - If DeployToAll is nil/false → returns the infra value
func resolveDeployTo(deployToAll *flexible.Field[bool], infraDefs *flexible.Field[[]*v0.InfrastructureDefinition]) interface{} {
	// Compute infra value from infrastructure definitions
	var infraValue interface{}
	if infraDefs != nil {
		if infra, ok := infraDefs.AsString(); ok {
			infraValue = infra
		} else if infra, ok := infraDefs.AsStruct(); ok {
			if len(infra) == 1 {
				infraValue = infra[0].Identifier
			} else if len(infra) > 1 {
				infraList := make([]string, 0, len(infra))
				for _, i := range infra {
					infraList = append(infraList, i.Identifier)
				}
				infraValue = infraList
			}
		}
	}

	if deployToAll == nil {
		return infraValue
	}

	// Boolean value
	if val, ok := deployToAll.AsStruct(); ok {
		if val {
			return "all"
		}
		return infraValue
	}

	// Expression value
	if expr, ok := deployToAll.AsString(); ok {
		if expr == "<+input>" {
			// <+input> means deploy-to is the infra only
			return infraValue
		}
		// Other expression: build ternary
		// <+ <+originalExpr> ? "all" : "infraFromYaml">
		infraStr := formatInfraForExpression(infraValue)
		return fmt.Sprintf("<+ %s ? \"all\" : %s>", expr, infraStr)
	}

	return infraValue
}

// formatInfraForExpression formats infrastructure value for use in ternary expressions.
// Single string → "infraName"
// List → ["infra1","infra2",...]
func formatInfraForExpression(infraValue interface{}) string {
	switch v := infraValue.(type) {
	case string:
		return fmt.Sprintf("%q", v)
	case []string:
		result := "["
		for i, s := range v {
			if i > 0 {
				result += ","
			}
			result += fmt.Sprintf("%q", s)
		}
		result += "]"
		return result
	default:
		return "\"\""
	}
}
