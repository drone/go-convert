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
	Machine struct {
		Default            bool   `yaml:"default,omitempty"`
		Image              string `yaml:"image,omitempty"`
		Shell              string `yaml:"shell,omitempty"`
		DockerLayerCaching bool   `yaml:"docker_layer_caching,omitempty"`
	}

	machine struct {
		Image              string `yaml:"image,omitempty"`
		Shell              string `yaml:"shell,omitempty"`
		DockerLayerCaching bool   `yaml:"docker_layer_caching,omitempty"`
	}
)

// UnmarshalYAML implements the unmarshal interface.
func (v *Machine) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 bool
	var out2 *machine

	if err := unmarshal(&out1); err == nil {
		v.Default = true
		return nil
	}

	if err := unmarshal(&out2); err == nil {
		v.Image = out2.Image
		v.Shell = out2.Shell
		v.DockerLayerCaching = out2.DockerLayerCaching
		return nil
	}

	return errors.New("failed to unmarshal machine")
}

// MarshalYAML implements the marshal interface.
func (v *Machine) MarshalYAML() (interface{}, error) {
	if v.Default {
		return true, nil
	}
	return &machine{
		Image:              v.Image,
		Shell:              v.Shell,
		DockerLayerCaching: v.DockerLayerCaching,
	}, nil
}
