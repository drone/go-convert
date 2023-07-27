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
	"errors"
)

type Variable struct {
	Value   string        `yaml:"value,omitempty"`
	Desc    string        `yaml:"description,omitempty"`
	Options Stringorslice `yaml:"options,omitempty"`
	Expand  *bool         `yaml:"expand,omitempty"`
}

// UnmarshalYAML implements the unmarshal interface.
func (v *Variable) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 string
	var out2 = struct {
		Value   string        `yaml:"value,omitempty"`
		Desc    string        `yaml:"description,omitempty"`
		Options Stringorslice `yaml:"options,omitempty"`
		Expand  interface{}   `yaml:"expand"`
	}{}

	if err := unmarshal(&out1); err == nil {
		expand := true
		v.Value = out1
		v.Expand = &expand
		return nil
	}

	if err := unmarshal(&out2); err == nil {
		v.Value = out2.Value
		v.Desc = out2.Desc
		v.Options = out2.Options

		switch out2.Expand.(type) {
		case bool:
			expand := out2.Expand.(bool)
			v.Expand = &expand
		default:
			expand := true
			v.Expand = &expand
		}
		return nil
	}

	return errors.New("failed to unmarshal variable")
}

func (v *Variable) MarshalYAML() (interface{}, error) {
	// If there's no description and expand is true or not set, just output the value as a string
	if v.Desc == "" && (v.Expand == nil || *v.Expand) {
		return v.Value, nil
	}

	return struct {
		Value   string        `yaml:"value,omitempty"`
		Desc    string        `yaml:"description,omitempty"`
		Options Stringorslice `yaml:"options,omitempty"`
		Expand  *bool         `yaml:"expand,omitempty"`
	}{
		Value:   v.Value,
		Desc:    v.Desc,
		Options: v.Options,
		Expand:  v.Expand,
	}, nil
}
