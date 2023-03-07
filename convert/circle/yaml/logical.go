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

type (
	Matches struct {
		Pattern string `yaml:"pattern,omitempty"`
		Value   string `yaml:"value,omitempty"`
	}

	Logical struct {
		Literal interface{}

		And     []*Logical
		Equal   []interface{}
		Matches *Matches
		Not     *Logical
		Or      []*Logical
	}

	logical struct {
		And     []*Logical    `yaml:"and,omitempty"`
		Equal   []interface{} `yaml:"equal,omitempty"`
		Matches *Matches      `yaml:"matches,omitempty"`
		Not     *Logical      `yaml:"not,omitempty"`
		Or      []*Logical    `yaml:"or,omitempty"`
	}
)

// UnmarshalYAML implements the unmarshal interface.
func (v *Logical) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 *logical
	var out2 interface{}

	if err := unmarshal(&out1); err == nil {
		v.And = out1.And
		v.Equal = out1.Equal
		v.Matches = out1.Matches
		v.Not = out1.Not
		v.Or = out1.Or

		// if the struct is empty, the user probably
		// assigned a yaml anchor
		//
		// wtf. yes this is really possible.
		if !v.IsEmpty() {
			return nil
		}
	}

	if err := unmarshal(&out2); err == nil {
		v.Literal = out2
		return nil
	}

	return errors.New("failed to unmarshal logical condition")
}

// MarshalYAML implements the marshal interface.
func (v *Logical) MarshalYAML() (interface{}, error) {
	if v.Literal != nil {
		return v.Literal, nil
	}
	return &logical{
		And:     v.And,
		Equal:   v.Equal,
		Matches: v.Matches,
		Not:     v.Not,
		Or:      v.Or,
	}, nil
}

// IsEmpty returns true if the struct is empty.
func (v *Logical) IsEmpty() bool {
	return v.Not == nil &&
		v.Matches == nil &&
		len(v.And) == 0 &&
		len(v.Equal) == 0 &&
		len(v.Or) == 0
}
