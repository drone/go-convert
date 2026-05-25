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

import "encoding/json"

// DeployToItem represents a single infrastructure deploy-to entry.
// Can be either a simple string "infraId" or an object with id and overlay.
type DeployToItem struct {
	Id   string                 `json:"id,omitempty"`
	With map[string]interface{} `json:"with,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler for DeployToItem.
// Handles: string "infraId" or object {id: "infraId", with: {...}}
func (d *DeployToItem) UnmarshalJSON(data []byte) error {
	// Try as simple string first
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		d.Id = str
		d.With = nil
		return nil
	}

	// Try as full object
	type Alias DeployToItem
	var alias Alias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*d = DeployToItem(alias)
	return nil
}

// MarshalJSON implements json.Marshaler for DeployToItem.
// Outputs: string "infraId" if no With, or object {id: "infraId", with: {...}} if With exists
func (d DeployToItem) MarshalJSON() ([]byte, error) {
	// If no With, marshal as simple string
	if len(d.With) == 0 {
		return json.Marshal(d.Id)
	}
	// Otherwise marshal as full object
	type Alias DeployToItem
	return json.Marshal((Alias)(d))
}

// EnvironmentItem represents a single environment entry.
type EnvironmentItem struct {
	Id       string      `json:"id,omitempty" yaml:"id,omitempty"`
	DeployTo interface{} `json:"deploy-to,omitempty" yaml:"deploy-to,omitempty"`
	Filters  []*Filter   `json:"filters,omitempty" yaml:"filters,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler for EnvironmentItem.
// Handles: string "envId" or object {id: "envId", deploy-to: ...}
func (e *EnvironmentItem) UnmarshalJSON(data []byte) error {
	// Try as simple string first
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		e.Id = str
		e.DeployTo = nil
		e.Filters = nil
		return nil
	}

	// Try as full object
	type Alias EnvironmentItem
	var alias Alias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*e = EnvironmentItem(alias)
	return nil
}

// MarshalJSON implements json.Marshaler for EnvironmentItem.
// Always outputs object format {id: "envId", ...} for consistency
func (e EnvironmentItem) MarshalJSON() ([]byte, error) {
	type Alias EnvironmentItem
	return json.Marshal((Alias)(e))
}
