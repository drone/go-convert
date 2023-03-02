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
	Script struct {
		Text string
		Pipe *Pipe
	}

	Pipe struct {
		Name      string            `yaml:"name,omitempty"`
		Image     string            `yaml:"pipe,omitempty"`
		Variables map[string]string `yaml:"variables,omitempty"`
	}
)

// UnmarshalYAML implements the unmarshal interface.
func (v *Script) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 string
	var out2 *Pipe
	if err := unmarshal(&out1); err == nil {
		v.Text = out1
		return nil
	}
	if err := unmarshal(&out2); err == nil {
		v.Pipe = out2
		return nil
	}
	return errors.New("failed to unmarshal script")
}

// MarshalYAML implements the marshal interface.
func (v *Script) MarshalYAML() (interface{}, error) {
	// marshal the pipe
	if v.Pipe != nil {
		return v.Pipe, nil
	}
	// else marshal the script text.
	return v.Text, nil
}
