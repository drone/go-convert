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

// Stages defines a set of build stages. Build stages are run
// sequentially. Stages run their Jobs in parallel.
//
// https://config.travis-ci.com/ref/stages
type Stages struct {
	Items []*Stage
}

// UnmarshalYAML implements the unmarshal interface.
func (v *Stages) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 *Stage
	var out2 []*Stage
	if err := unmarshal(&out1); err == nil {
		v.Items = append(v.Items, out1)
		return nil
	}
	if err := unmarshal(&out2); err == nil {
		v.Items = out2
		return nil
	}
	return errors.New("failed to unmarshal stages")
}

// MarshalYAML implements the marshal interface.
func (v *Stages) MarshalYAML() (interface{}, error) {
	return v.Items, nil
}

// Stage defines a buld stage.
type Stage struct {
	Name string `yaml:"name,omitempty"`
	If   string `yaml:"if,omitempty"`
}

// UnmarshalYAML implements the unmarshal interface.
func (v *Stage) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 string
	var out2 = struct {
		Name string `yaml:"name"`
		If   string `yaml:"if"`
	}{}
	if err := unmarshal(&out1); err == nil {
		v.Name = out1
		return nil
	}
	if err := unmarshal(&out2); err == nil {
		v.Name = out2.Name
		v.If = out2.If
		return nil
	}
	return errors.New("failed to unmarshal stage")
}
