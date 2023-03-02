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

type Imports struct {
	Items []*Import
}

// UnmarshalYAML implements the unmarshal interface.
func (v *Imports) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 *Import
	var out2 []*Import
	if err := unmarshal(&out1); err == nil {
		v.Items = append(v.Items, out1)
		return nil
	}
	if err := unmarshal(&out2); err == nil {
		v.Items = out2
		return nil
	}
	return errors.New("failed to unmarshal imports")
}

// MarshalYAML implements the marshal interface.
func (v *Imports) MarshalYAML() (interface{}, error) {
	return v.Items, nil
}

type Import struct {
	Source string `yaml:"source,omitempty"`
	Mode   string `yaml:"mode,omitempty"` // merge, deep_merge, deep_merge_append, deep_merge_prepend
	If     string `yaml:"if,omitempty"`
}

// UnmarshalYAML implements the unmarshal interface.
func (v *Import) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 string
	var out2 = struct {
		Source string `yaml:"source"`
		Mode   string `yaml:"mode"`
		If     string `yaml:"if"`
	}{}
	if err := unmarshal(&out1); err == nil {
		v.Source = out1
		return nil
	}
	if err := unmarshal(&out2); err == nil {
		v.Source = out2.Source
		v.Mode = out2.Mode
		v.If = out2.If
		return nil
	}
	return errors.New("failed to unmarshal import")
}
