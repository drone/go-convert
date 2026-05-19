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
)

// EnvironmentRef is the unified v1 environment configuration.
type EnvironmentRef struct {
	Items      []*EnvironmentItem    `json:"items,omitempty" yaml:"items,omitempty"`
	Sequential *flexible.Field[bool] `json:"sequential,omitempty" yaml:"sequential,omitempty"`
	Group      interface{}           `json:"group,omitempty" yaml:"group,omitempty"`
	Filters    []*Filter             `json:"filters,omitempty" yaml:"filters,omitempty"`
	MultiEnv   bool                  `json:"-" yaml:"-"`
}

// MarshalJSON implements json.Marshaler following the ServiceRef pattern.
func (v EnvironmentRef) MarshalJSON() ([]byte, error) {
	// Variant: single env – marshal as flat EnvironmentItem (like ServiceRef string)
	if len(v.Items) == 1 && !v.MultiEnv {
		return json.Marshal(v.Items[0])
	}

	// Variant: multi env – marshal as struct with items
	type Alias EnvironmentRef
	return json.Marshal((*Alias)(&v))
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (v *EnvironmentRef) UnmarshalJSON(data []byte) error {
	var out1 string
	var out2 = struct {
		Items      []*EnvironmentItem    `json:"items,omitempty" yaml:"items,omitempty"`
		Sequential *flexible.Field[bool] `json:"sequential,omitempty" yaml:"sequential,omitempty"`
		Group      interface{}           `json:"group,omitempty" yaml:"group,omitempty"`
	}{}

	if err := json.Unmarshal(data, &out1); err == nil {
		v.Items = []*EnvironmentItem{
			{Id: out1},
		}
		v.MultiEnv = false
		return nil
	}

	if err := json.Unmarshal(data, &out2); err == nil {
		v.Sequential = out2.Sequential
		v.Items = out2.Items
		v.Group = out2.Group
		// Set MultiEnv flag based on whether this is a multi-environment config
		// MultiEnv is true if: multiple items, or has sequential/group fields
		v.MultiEnv = len(out2.Items) > 0 || out2.Sequential != nil || out2.Group != nil
		return nil
	} else {
		return err
	}
}
