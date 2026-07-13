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
	"strings"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	"github.com/drone/go-convert/convert/v0tov1/messagelog"
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
		messagelog.GetMessageLogger().LogWarning(
			"USE_FROM_STAGE_NOT_FOUND",
			fmt.Sprintf("service useFromStage %q not found in context", src.UseFromStage.Stage),
			messagelog.WithContext(map[string]string{"from_stage": src.UseFromStage.Stage, "kind": "service"}),
		)
		return nil
	}

	// For single service, return simple string reference
	if src.ServiceRef != "" {
		var serviceWith map[string]interface{}
		if hasValidServiceInputs(src.ServiceInputs) {
			serviceWith = map[string]interface{}{
				"overlay": src.ServiceInputs,
			}
		}
		serviceItem := &v1.ServiceItem{
			Id:   src.ServiceRef,
			With: serviceWith,
			Ref:  src.GitBranch,
		}
		return &v1.ServiceRef{
			Items:        []*v1.ServiceItem{serviceItem},
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
		messagelog.GetMessageLogger().LogWarning(
			"USE_FROM_STAGE_NOT_FOUND",
			fmt.Sprintf("services useFromStage %q not found in context", src.UseFromStage.Stage),
			messagelog.WithContext(map[string]string{"from_stage": src.UseFromStage.Stage, "kind": "services"}),
		)
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
			messagelog.GetMessageLogger().LogWarning(
				"UNSUPPORTED_EXPRESSION",
				fmt.Sprintf("services.values contains unsupported expression %q; skipping conversion", expr),
				messagelog.WithContext(map[string]string{"expression": expr, "field": "services.values"}),
			)
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

	serviceItems := make([]*v1.ServiceItem, 0, len(values))

	for _, service := range values {
		if service != nil && service.ServiceRef != "" {
			var serviceWith map[string]interface{}
			if hasValidServiceInputs(service.ServiceInputs) {
				serviceWith = map[string]interface{}{
					"overlay": service.ServiceInputs,
				}
			}
			serviceItems = append(serviceItems, &v1.ServiceItem{
				Id:   service.ServiceRef,
				With: serviceWith,
				Ref:  service.GitBranch,
			})
		}
	}
	// v0 defaulted to parallel execution; preserve that in v1 (which now defaults to serial)
	// by emitting parallel: true unless v0 explicitly set parallel.
	var parallel *flexible.Field[bool] = &flexible.Field[bool]{Value: true}
	if src.Metadata != nil && src.Metadata.Parallel != nil {
		parallel = src.Metadata.Parallel
	}
	if len(serviceItems) > 0 {
		return &v1.ServiceRef{
			Items:        serviceItems,
			MultiService: true,
			Parallel:     parallel,
		}
	}

	return nil
}

// func ConvertDeploymentServiceConfig(src *v0.DeploymentServiceConfig) *v1.ServiceRef {
// 	if src == nil {
// 		return nil
// 	}
// 	if src.ServiceRef != "" {
// 		return &v1.ServiceRef{
// 			Items: []string{src.ServiceRef},
// 		}
// 	}

// 	if src.ServiceItem != nil {
// 		return &v1.ServiceRef{
// 			Items: []string{src.ServiceItem.Identifier},
// 		}
// 	}
// 	return nil
// }

// hasValidServiceInputs checks if serviceInputs should be included in v1 output.
// Returns false for nil, empty string, or expression values (strings containing <+).
func hasValidServiceInputs(serviceInputs interface{}) bool {
	if serviceInputs == nil {
		return false
	}

	// Check if it's a string (expression or empty)
	if str, ok := serviceInputs.(string); ok {
		// Empty string or expression - skip
		if str == "" || strings.Contains(str, "<+") {
			return false
		}
		// Non-empty, non-expression string - this shouldn't happen but skip anyway
		return false
	}

	// It's a struct/map - include it
	return true
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
		messagelog.GetMessageLogger().LogWarning(
			"USE_FROM_STAGE_NOT_FOUND",
			fmt.Sprintf("environment useFromStage %q not found in context", src.UseFromStage.Stage),
			messagelog.WithContext(map[string]string{"from_stage": src.UseFromStage.Stage, "kind": "environment"}),
		)
		return nil
	}

	if src.EnvironmentRef == "" {
		return nil
	}

	// Resolve deploy-to from both singular and plural infra definitions
	deployTo := resolveEnvironmentDeployTo(src)

	item := &v1.EnvironmentItem{
		Id:        src.EnvironmentRef,
		DeployTo:  deployTo,
		Overrides: buildEnvironmentOverrides(src.EnvironmentInputs, src.ServiceOverrideInputs),
	}

	return &v1.EnvironmentRef{
		Items:    []*v1.EnvironmentItem{item},
		MultiEnv: false,
	}
}

// buildEnvironmentOverrides converts v0 environmentInputs and serviceOverrideInputs
// into the v1 environment "overrides" block.
//
// environmentInputs    -> overrides.env-global
// serviceOverrideInputs -> overrides.env-service
//
// When the inputs are a concrete struct, the value is wrapped as:
//
//	env-global:
//	  with:
//	    overlay:
//	      <inputs>
//
// When the inputs are a runtime expression (e.g. <+input>), the key is set
// directly to the expression string:
//
//	env-global: <+input>
func buildEnvironmentOverrides(envInputs, serviceOverrideInputs interface{}) map[string]interface{} {
	overrides := make(map[string]interface{})

	if entry := buildOverrideEntry(envInputs); entry != nil {
		overrides["env-global"] = entry
	}
	if entry := buildOverrideEntry(serviceOverrideInputs); entry != nil {
		overrides["env-service"] = entry
	}

	if len(overrides) == 0 {
		return nil
	}
	return overrides
}

// buildOverrideEntry builds a single override entry from v0 inputs.
func buildOverrideEntry(inputs interface{}) interface{} {
	if inputs == nil {
		return nil
	}

	// Expression (e.g. <+input>) or empty string.
	if str, ok := inputs.(string); ok {
		if str == "" {
			return nil
		}
		// Pass the expression through directly.
		return str
	}

	return map[string]interface{}{
		"with": map[string]interface{}{
			"overlay": inputs,
		},
	}
}

// ConvertEnvironments converts v0 Environments to v1 EnvironmentRef (multi-env format).
func ConvertEnvironments(src *v0.Environments, ctx *StageConversionContext) *v1.EnvironmentRef {
	if src == nil {
		return nil
	}
	// Handle useFromStage
	if src.UseFromStage != nil && src.UseFromStage.Stage != "" {
		if ref := ctx.GetEnvironment(src.UseFromStage.Stage); ref != nil {
			return ref
		}
		messagelog.GetMessageLogger().LogWarning(
			"USE_FROM_STAGE_NOT_FOUND",
			fmt.Sprintf("environments useFromStage %q not found in context", src.UseFromStage.Stage),
			messagelog.WithContext(map[string]string{"from_stage": src.UseFromStage.Stage, "kind": "environments"}),
		)
		return nil
	}

	// v0 defaulted to parallel execution; preserve that in v1 (which now defaults to serial)
	// by emitting parallel: true unless v0 explicitly set parallel.
	var parallel *flexible.Field[bool] = &flexible.Field[bool]{Value: true}
	if src.Metadata != nil && src.Metadata.Parallel != nil {
		parallel = src.Metadata.Parallel
	}

	// Check if Values is nil or empty
	hasValues := false
	var values []*v0.Environment
	if src.Values != nil {
		if v, ok := src.Values.AsStruct(); ok && len(v) > 0 {
			hasValues = true
			values = v
		}
	}

	// Check if Filters exist
	hasFilters := false
	var v0Filters []*v0.EnvironmentFilter
	if src.Filters != nil {
		if f, ok := src.Filters.AsStruct(); ok && len(f) > 0 {
			hasFilters = true
			v0Filters = f
		}
	}

	// Case: top-level filters without values (filter-based multi-environment)
	if !hasValues && hasFilters {
		v1Filters := ConvertEnvironmentFilters(v0Filters)
		if len(v1Filters) > 0 {
			return &v1.EnvironmentRef{
				Parallel: parallel,
				Filters:  v1Filters,
			}
		}
		return nil
	}

	if !hasValues {
		return nil
	}

	items := make([]*v1.EnvironmentItem, 0, len(values))
	for _, env := range values {
		if env == nil || env.EnvironmentRef == "" {
			continue
		}
		item := &v1.EnvironmentItem{
			Id:        env.EnvironmentRef,
			Overrides: buildEnvironmentOverrides(env.EnvironmentInputs, env.ServiceOverrideInputs),
		}
		// Per-environment filters (infra-level filters on each env)
		if env.Filters != nil {
			if ef, ok := env.Filters.AsStruct(); ok && len(ef) > 0 {
				v1Filters := ConvertEnvironmentFilters(ef)
				if len(v1Filters) > 0 {
					item.Filters = v1Filters
				}
			}
		}
		if len(item.Filters) == 0 {
			item.DeployTo = resolveEnvironmentDeployTo(env)
		}
		items = append(items, item)
	}

	if len(items) > 0 {
		return &v1.EnvironmentRef{
			Items:    items,
			Parallel: parallel,
			MultiEnv: true,
		}
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
		messagelog.GetMessageLogger().LogWarning(
			"USE_FROM_STAGE_NOT_FOUND",
			fmt.Sprintf("environmentGroup useFromStage %q not found in context", src.UseFromStage.Stage),
			messagelog.WithContext(map[string]string{"from_stage": src.UseFromStage.Stage, "kind": "environmentGroup"}),
		)
		return nil
	}

	// v0 defaulted to parallel execution; preserve that in v1 (which now defaults to serial)
	// by emitting parallel: true unless v0 explicitly set parallel.
	var parallel *flexible.Field[bool] = &flexible.Field[bool]{Value: true}
	if src.Metadata != nil && src.Metadata.Parallel != nil {
		parallel = src.Metadata.Parallel
	}

	// Case: environments is an expression (e.g., <+input>)
	if src.Environments != nil {
		if expr, ok := src.Environments.AsString(); ok && expr != "" {
			// Expression for environments - pass through
			groupConfig := map[string]interface{}{
				"id": src.EnvGroupRef,
			}
			return &v1.EnvironmentRef{
				Parallel: parallel,
				Group:      groupConfig,
			}
		}
	}

	// Case: filters is an expression (e.g., <+input>)
	if src.Filters != nil {
		if expr, ok := src.Filters.AsString(); ok && expr != "" {
			groupConfig := map[string]interface{}{
				"id": src.EnvGroupRef,
			}
			return &v1.EnvironmentRef{
				Parallel: parallel,
				Group:      groupConfig,
			}
		}
	}

	// Case: group ref with top-level filters (no environments)
	if src.EnvGroupRef != "" && (src.Environments == nil || src.Environments.IsNil()) {
		// Check for filters
		if src.Filters != nil {
			if filters, ok := src.Filters.AsStruct(); ok && len(filters) > 0 {
				v1Filters := ConvertEnvironmentFilters(filters)
				groupConfig := map[string]interface{}{
					"id":      src.EnvGroupRef,
					"filters": v1Filters,
				}
				return &v1.EnvironmentRef{
					Parallel: parallel,
					Group:      groupConfig,
				}
			}
		}

		// Simple group reference format
		return &v1.EnvironmentRef{
			Parallel: parallel,
			Group:    map[string]interface{}{"id": src.EnvGroupRef},
		}
	}

	// Handle explicit group with environments array
	if src.Environments != nil {
		envItems, ok := src.Environments.AsStruct()
		if ok && len(envItems) > 0 {
			items := convertEnvironmentGroupEnvItems(envItems)
			if len(items) > 0 {
				groupConfig := map[string]interface{}{
					"id":    src.EnvGroupRef,
					"items": items,
				}
				return &v1.EnvironmentRef{
					Parallel: parallel,
					Group:      groupConfig,
				}
			}
		}
	}

	return nil
}

// convertEnvironmentGroupEnvItems converts v0 Environment array to v1 EnvironmentItem array
func convertEnvironmentGroupEnvItems(envItems []*v0.Environment) []*v1.EnvironmentItem {
	items := make([]*v1.EnvironmentItem, 0, len(envItems))
	for _, env := range envItems {
		if env == nil || env.EnvironmentRef == "" {
			continue
		}

		deployTo := resolveEnvironmentDeployTo(env)
		items = append(items, &v1.EnvironmentItem{
			Id:        env.EnvironmentRef,
			DeployTo:  deployTo,
			Overrides: buildEnvironmentOverrides(env.EnvironmentInputs, env.ServiceOverrideInputs),
		})
	}
	return items
}

// ConvertDeploymentInfrastructure converts v0 DeploymentInfrastructure to v1 EnvironmentRef
// It adds the environmentRef from infrastructure to the environments of the stage
func ConvertDeploymentInfrastructure(src *v0.DeploymentInfrastructure) *v1.EnvironmentRef {
	if src == nil || src.EnvironmentRef == "" {
		return nil
	}

	// Create environment item with the environmentRef from infrastructure
	envItem := &v1.EnvironmentItem{
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

// resolveEnvironmentDeployTo resolves the deploy-to value for a v0 Environment.
// Checks both singular infrastructureDefinition and plural infrastructureDefinitions.
func resolveEnvironmentDeployTo(env *v0.Environment) interface{} {
	infraDefs := collectInfraDefinitions(env)
	return resolveDeployTo(env.DeployToAll, infraDefs)
}

// collectInfraDefinitions merges singular and plural infrastructure definitions into one list.
func collectInfraDefinitions(env *v0.Environment) *flexible.Field[[]*v0.InfrastructureDefinition] {
	// Prefer plural if present
	if env.InfrastructureDefinitions != nil && !env.InfrastructureDefinitions.IsNil() {
		return env.InfrastructureDefinitions
	}
	// Fall back to singular
	if env.InfrastructureDefinition != nil && !env.InfrastructureDefinition.IsNil() {
		if env.InfrastructureDefinition.IsExpression() {
			if expr, ok := env.InfrastructureDefinition.AsString(); ok {
				f := &flexible.Field[[]*v0.InfrastructureDefinition]{}
				f.SetExpression(expr)
				return f
			}
		}
		if single, ok := env.InfrastructureDefinition.AsStruct(); ok {
			return &flexible.Field[[]*v0.InfrastructureDefinition]{Value: []*v0.InfrastructureDefinition{&single}}
		}
	}
	return nil
}

// hasValidInfraInputs checks if infrastructure inputs should be included in v1 output.
// Returns false for nil, empty string, or expression values (strings containing <+).
func hasValidInfraInputs(inputs interface{}) bool {
	if inputs == nil {
		return false
	}

	// Check if it's a string (expression or empty)
	if str, ok := inputs.(string); ok {
		// Empty string or expression - skip
		if str == "" || strings.Contains(str, "<+") {
			return false
		}
		// Non-empty, non-expression string - this shouldn't happen but skip anyway
		return false
	}

	// It's a struct/map - include it
	return true
}

// resolveDeployTo determines the deploy-to value based on DeployToAll and infrastructure definitions.
// - If DeployToAll is a boolean true → "all"
// - If DeployToAll is <+input> → returns the infra value only (single string or list)
// - If DeployToAll is any other expression → builds ternary: <+ <+expr> ? "all" : "infraValue" >
// - If DeployToAll is nil/false → returns the infra value
// - If infrastructure has inputs, returns DeployToItem with overlay
func resolveDeployTo(deployToAll *flexible.Field[bool], infraDefs *flexible.Field[[]*v0.InfrastructureDefinition]) interface{} {
	// Compute infra value from infrastructure definitions
	var infraValue interface{}
	if infraDefs != nil {
		if infra, ok := infraDefs.AsString(); ok {
			infraValue = infra
		} else if infra, ok := infraDefs.AsStruct(); ok {
			if len(infra) == 1 {
				// Single infrastructure - check for inputs
				if hasValidInfraInputs(infra[0].Inputs) {
					infraValue = &v1.DeployToItem{
						Id: infra[0].Identifier,
						With: map[string]interface{}{
							"overlay": infra[0].Inputs,
						},
					}
				} else {
					infraValue = infra[0].Identifier
				}
			} else if len(infra) > 1 {
				// Multiple infrastructures - check each for inputs
				hasAnyInputs := false
				for _, i := range infra {
					if hasValidInfraInputs(i.Inputs) {
						hasAnyInputs = true
						break
					}
				}

				if hasAnyInputs {
					// At least one has inputs - use DeployToItem array
					infraItems := make([]*v1.DeployToItem, 0, len(infra))
					for _, i := range infra {
						item := &v1.DeployToItem{Id: i.Identifier}
						if hasValidInfraInputs(i.Inputs) {
							item.With = map[string]interface{}{
								"overlay": i.Inputs,
							}
						}
						infraItems = append(infraItems, item)
					}
					infraValue = infraItems
				} else {
					// No inputs - use simple string array
					infraList := make([]string, 0, len(infra))
					for _, i := range infra {
						infraList = append(infraList, i.Identifier)
					}
					infraValue = infraList
				}
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
