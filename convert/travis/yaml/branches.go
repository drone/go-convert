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

type Branches struct {
	Only   []string `yaml:"only,omitempty"`
	Except []string `yaml:"except,omitempty"`
}

// UnmarshalYAML implements the unmarshal interface.
func (v *Branches) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 string
	var out2 []string
	var out3 = struct {
		Only   Stringorslice `yaml:"only"`
		Except Stringorslice `yaml:"except"`
	}{}

	if err := unmarshal(&out1); err == nil {
		v.Only = []string{out1}
		return nil
	}

	if err := unmarshal(&out2); err == nil {
		v.Only = out2
		return nil
	}

	if err := unmarshal(&out3); err == nil {
		v.Only = out3.Only
		v.Except = out3.Except
		return nil
	}

	return errors.New("failed to unmarshal branches")
}
