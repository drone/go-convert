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

// Image defines a container image.
// https://docs.gitlab.com/ee/ci/yaml/#image
type Image struct {
	Name       string        `yaml:"name,omitempty"`
	Alias      string        `yaml:"alias,omitempty"`
	Entrypoint Stringorslice `yaml:"entrypoint,omitempty"`
	Command    Stringorslice `yaml:"command,omitempty"`
	PullPolicy Stringorslice `yaml:"pull_policy,omitempty"` // single string or array of strings
}

// UnmarshalYAML implements the unmarshal interface.
func (v *Image) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 string
	var out2 = struct {
		Name       string        `yaml:"name,omitempty"`
		Alias      string        `yaml:"alias,omitempty"`
		Entrypoint Stringorslice `yaml:"entrypoint,omitempty"`
		Command    Stringorslice `yaml:"command,omitempty"`
		PullPolicy Stringorslice `yaml:"pull_policy,omitempty"` // single string or array of strings
	}{}

	if err := unmarshal(&out1); err == nil {
		v.Name = out1
		return nil
	}

	if err := unmarshal(&out2); err == nil {
		v.Name = out2.Name
		v.Alias = out2.Alias
		v.Entrypoint = out2.Entrypoint
		v.Command = out2.Command
		v.PullPolicy = out2.PullPolicy
		return nil
	}

	return errors.New("failed to unmarshal image")
}
