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

// ConvertEnvironmentFilters converts a slice of v0 EnvironmentFilter to v1 []*Filter.
func ConvertEnvironmentFilters(filters []*v0.EnvironmentFilter) []*v1.Filter {
	if len(filters) == 0 {
		return nil
	}

	var result []*v1.Filter

	for _, f := range filters {
		if f == nil {
			continue
		}

		// Handle Entities as flexible.Field - can be expression or []string
		if f.Entities == nil || f.Entities.IsNil() {
			continue
		}

		// If Entities is an expression, skip (can't resolve at conversion time)
		if f.Entities.IsExpression() {
			continue
		}

		// Normal case: Entities is []string
		entities, ok := f.Entities.AsStruct()
		if !ok {
			continue
		}

		for _, entity := range entities {
			filter := convertSingleFilter(f, entity)
			if filter != nil {
				result = append(result, filter)
			}
		}
	}

	if len(result) == 0 {
		return nil
	}
	return result
}

// convertSingleFilter converts a single v0 filter for one entity type to a concrete v1.Filter.
func convertSingleFilter(f *v0.EnvironmentFilter, entity string) *v1.Filter {
	filterEntity := buildFilterEntity(f)
	if filterEntity == nil {
		return nil
	}

	filter := &v1.Filter{}
	switch entity {
	case "environments":
		filter.Environments = filterEntity
	case "infrastructures":
		filter.Infrastructures = filterEntity
	default:
		return nil
	}
	return filter
}

// buildFilterEntity builds a v1.FilterEntity from a v0 filter's type and spec.
func buildFilterEntity(f *v0.EnvironmentFilter) *v1.FilterEntity {
	switch f.Type {
	case "all":
		return v1.NewFilterAll()

	case "tags":
		if f.Spec == nil || f.Spec.Tags == nil || f.Spec.Tags.IsNil() {
			return nil
		}

		matchKey := "all"
		if f.Spec.MatchType == "any" {
			matchKey = "in"
		}

		// Tags is flexible.Field[map[string]string] - can be expression or map
		if f.Spec.Tags.IsExpression() {
			if exprStr, ok := f.Spec.Tags.AsString(); ok {
				return v1.NewFilterTagsExpression(matchKey, exprStr)
			}
			return nil
		}

		// Normal case: Tags is map[string]string
		tags, ok := f.Spec.Tags.AsStruct()
		if !ok || len(tags) == 0 {
			return nil
		}

		tagsMap := make(map[string]interface{}, len(tags))
		for k, v := range tags {
			tagsMap[k] = v
		}
		return v1.NewFilterTags(matchKey, tagsMap)

	default:
		return nil
	}
}
