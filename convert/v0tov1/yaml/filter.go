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

package yaml

import (
	"encoding/json"

	"github.com/drone/go-convert/internal/flexible"
	"gopkg.in/yaml.v3"
)

// Filter represents a single v1 environment/infrastructure filter entry.
//
// v1 schema (Filters is an array of Filter):
//
//	filters:
//	  - environments: all
//	  - infrastructures:
//	      tags:
//	        all: { key: value }
//	  - environments:
//	      tags:
//	        in: { key: value }
type Filter struct {
	Environments    *FilterEntity `json:"environments,omitempty" yaml:"environments,omitempty"`
	Infrastructures *FilterEntity `json:"infrastructures,omitempty" yaml:"infrastructures,omitempty"`
}

// FilterEntity represents a filter entity value. It can be:
//   - "all" (string)
//   - A single FilterCondition (object with tags)
//   - An array of FilterCondition
//   - An object with or/and arrays of FilterCondition
type FilterEntity struct {
	All        bool               `json:"-" yaml:"-"`
	Conditions []*FilterCondition `json:"-" yaml:"-"`
}

// MarshalJSON implements the json.Marshaler interface.
func (e *FilterEntity) MarshalJSON() ([]byte, error) {
	if e.All {
		return json.Marshal("all")
	}
	if len(e.Conditions) >= 1 {
		return json.Marshal(e.Conditions)
	}
	return json.Marshal(nil)
}

// MarshalYAML implements the yaml.Marshaler interface.
func (e FilterEntity) MarshalYAML() (interface{}, error) {
	if e.All {
		return "all", nil
	}
	if len(e.Conditions) >= 1 {
		return e.Conditions, nil
	}
	return nil, nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (e *FilterEntity) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.ScalarNode && node.Value == "all" {
		e.All = true
		return nil
	}

	var single FilterCondition
	if err := node.Decode(&single); err == nil && single.Tags != nil {
		e.Conditions = []*FilterCondition{&single}
		return nil
	}

	var arr []*FilterCondition
	if err := node.Decode(&arr); err == nil {
		e.Conditions = arr
		return nil
	}

	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (e *FilterEntity) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		if str == "all" {
			e.All = true
		}
		return nil
	}

	var single FilterCondition
	if err := json.Unmarshal(data, &single); err == nil {
		e.Conditions = []*FilterCondition{&single}
		return nil
	}

	var arr []*FilterCondition
	if err := json.Unmarshal(data, &arr); err == nil {
		e.Conditions = arr
		return nil
	}

	return nil
}

// FilterCondition represents a tag-based filter condition.
//
//	tags:
//	  all: { key: value }   # matchType=all → AND
//	  in:  { key: value }   # matchType=any → OR
type FilterCondition struct {
	Tags *FilterTags `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// FilterTags holds tag filter criteria. Exactly one of All or In should be set.
// Both fields support expressions (e.g. <+service.tags>) via flexible.Field.
type FilterTags struct {
	All *flexible.Field[map[string]interface{}] `json:"all,omitempty" yaml:"all,omitempty"`
	In  *flexible.Field[map[string]interface{}] `json:"in,omitempty" yaml:"in,omitempty"`
}

// NewFilterAll creates a FilterEntity that matches all items.
func NewFilterAll() *FilterEntity {
	return &FilterEntity{All: true}
}

// NewFilterTags creates a FilterEntity with a single tag condition.
func NewFilterTags(matchType string, tags map[string]interface{}) *FilterEntity {
	ft := &FilterTags{}
	if matchType == "in" || matchType == "any" {
		ft.In = &flexible.Field[map[string]interface{}]{Value: tags}
	} else {
		ft.All = &flexible.Field[map[string]interface{}]{Value: tags}
	}
	return &FilterEntity{
		Conditions: []*FilterCondition{{Tags: ft}},
	}
}

// NewFilterTagsExpression creates a FilterEntity with an expression-valued tag condition.
func NewFilterTagsExpression(matchType string, expr string) *FilterEntity {
	ft := &FilterTags{}
	if matchType == "in" || matchType == "any" {
		ft.In = &flexible.Field[map[string]interface{}]{Value: expr}
	} else {
		ft.All = &flexible.Field[map[string]interface{}]{Value: expr}
	}
	return &FilterEntity{
		Conditions: []*FilterCondition{{Tags: ft}},
	}
}
