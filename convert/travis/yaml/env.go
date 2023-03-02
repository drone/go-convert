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
	"strings"
)

type Env struct {
	Global []map[string]string `yaml:"global,omitempty"`
	Jobs   []map[string]string `yaml:"jobs,omitempty"`
}

// UnmarshalYAML implements the unmarshal interface.
func (v *Env) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 *Envmap
	var out2 []*Envmap
	var out3 = struct {
		Global []*Envmap `yaml:"global"`
		Jobs   []*Envmap `yaml:"jobs"`
	}{}

	if err := unmarshal(&out1); err == nil {
		v.Global = append(v.Global, out1.Items)
		return nil
	}

	if err := unmarshal(&out2); err == nil {
		for _, vv := range out2 {
			v.Global = append(v.Global, vv.Items)
		}
		return nil
	}

	if err := unmarshal(&out3); err == nil {
		for _, vv := range out3.Global {
			v.Global = append(v.Global, vv.Items)
		}
		for _, vv := range out3.Jobs {
			v.Jobs = append(v.Jobs, vv.Items)
		}
		return nil
	}

	return errors.New("failed to unmarshal env")
}

type Envmap struct {
	Items map[string]string
}

// UnmarshalYAML implements the unmarshal interface.
func (v *Envmap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 string
	var out2 []string
	var out3 map[string]string

	// fist we attempt to unmarshal a string in key=value format
	if err := unmarshal(&out1); err == nil {
		v.Items = map[string]string{}
		parts := strings.SplitN(out1, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			val := parts[1]
			v.Items[key] = val
		}
		return nil
	}
	// then we attempt to unmarshal a string slice in key=value format
	if err := unmarshal(&out2); err == nil {
		v.Items = map[string]string{}
		for _, vv := range out2 {
			parts := strings.SplitN(vv, "=", 2)
			if len(parts) == 2 {
				key := parts[0]
				val := parts[1]
				v.Items[key] = val
			}
		}
		return nil
	}
	// then we attempt to unmarshal a map
	if err := unmarshal(&out3); err == nil {
		v.Items = map[string]string{}
		for key, val := range out3 {
			v.Items[key] = val
		}
		return nil
	}

	return errors.New("failed to unmarshal env")
}
