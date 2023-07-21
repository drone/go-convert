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

import "errors"

type Inherit struct {
	Default   *InheritKeys `yaml:"default,omitempty"`
	Variables *InheritKeys `yaml:"variables,omitempty"`
}

type InheritKeys struct {
	All  bool     `yaml:"all,omitempty"`
	Keys []string `yaml:"keys,omitempty"`
}

// UnmarshalYAML implements the unmarshal interface for InheritKeys.
func (v *InheritKeys) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 bool
	var out2 []string

	if err := unmarshal(&out1); err == nil {
		v.All = !out1
		return nil
	}

	if err := unmarshal(&out2); err == nil {
		v.Keys = out2
		return nil
	}

	return errors.New("failed to unmarshal inherit keys")
}

// UnmarshalYAML implements the unmarshal interface for Inherit.
func (v *Inherit) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 bool
	var out2 struct {
		Default   *InheritKeys `yaml:"default,omitempty"`
		Variables *InheritKeys `yaml:"variables,omitempty"`
	}

	if err := unmarshal(&out1); err == nil {
		v.Default = &InheritKeys{All: !out1}
		v.Variables = &InheritKeys{All: !out1}
		return nil
	}

	if err := unmarshal(&out2); err == nil {
		v.Default = out2.Default
		v.Variables = out2.Variables
		return nil
	}

	return errors.New("failed to unmarshal inherit")
}

func (v *Inherit) MarshalYAML() (interface{}, error) {
	m := make(map[string]interface{})

	if v.Default != nil {
		if v.Default.All {
			m["default"] = false
		} else {
			m["default"] = v.Default.Keys
		}
	}

	if v.Variables != nil {
		if v.Variables.All {
			m["variables"] = false
		} else {
			m["variables"] = v.Variables.Keys
		}
	}

	return m, nil
}
