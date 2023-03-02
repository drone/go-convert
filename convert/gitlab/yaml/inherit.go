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
	Default   *InheritKeys `yaml:"default"`
	Variables *InheritKeys `yaml:"variables"`
}

type InheritKeys struct {
	None bool
	Keys []string
}

// UnmarshalYAML implements the unmarshal interface.
func (v *InheritKeys) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 bool
	var out2 []string

	if err := unmarshal(&out1); err == nil {
		v.None = !out1
		return nil
	}

	if err := unmarshal(&out2); err == nil {
		v.Keys = out2
		return nil
	}

	return errors.New("failed to unmarshal inherit")
}
