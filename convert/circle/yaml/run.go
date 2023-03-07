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
	Run struct {
		Background       bool
		Command          string
		Environment      map[string]string
		Name             string
		NoOutputTimeout  string
		Shell            string
		When             string
		WorkingDirectory string
	}

	run struct {
		Background       bool              `yaml:"background,omitempty"`
		Command          string            `yaml:"command"`
		Environment      map[string]string `yaml:"environment,omitempty"`
		Name             string            `yaml:"name,omitempty"`
		NoOutputTimeout  string            `yaml:"no_output_timeout,omitempty"`
		Shell            string            `yaml:"shell,omitempty"`
		When             string            `yaml:"when,omitempty"`
		WorkingDirectory string            `yaml:"working_directory,omitempty"`
	}
)

// UnmarshalYAML implements the unmarshal interface.
func (v *Run) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 string
	var out2 *run

	if err := unmarshal(&out1); err == nil {
		v.Command = out1
		return nil
	}

	if err := unmarshal(&out2); err == nil {
		v.Background = out2.Background
		v.Command = out2.Command
		v.Environment = out2.Environment
		v.Name = out2.Name
		v.NoOutputTimeout = out2.NoOutputTimeout
		v.Shell = out2.Shell
		v.When = out2.When
		v.WorkingDirectory = out2.WorkingDirectory
		return nil
	}

	return errors.New("failed to unmarshal machine")
}

// MarshalYAML implements the marshal interface.
func (v *Run) MarshalYAML() (interface{}, error) {
	if v.IsEmpty() {
		return v.Command, nil
	}
	return &run{
		Background:       v.Background,
		Command:          v.Command,
		Environment:      v.Environment,
		Name:             v.Name,
		NoOutputTimeout:  v.NoOutputTimeout,
		Shell:            v.Shell,
		When:             v.When,
		WorkingDirectory: v.WorkingDirectory,
	}, nil
}

// IsEmpty returns true if the struct is empty.
func (v *Run) IsEmpty() bool {
	return !v.Background &&
		len(v.Environment) == 0 &&
		len(v.Name) == 0 &&
		len(v.NoOutputTimeout) == 0 &&
		len(v.Shell) == 0 &&
		len(v.When) == 0 &&
		len(v.WorkingDirectory) == 0
}
