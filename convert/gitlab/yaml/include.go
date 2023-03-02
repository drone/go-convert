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

// Include includes external yaml files.
// https://docs.gitlab.com/ee/ci/yaml/#include
type Include struct {
	Local    string        `yaml:"local,omitempty"`
	Project  string        `yaml:"project,omitempty"`
	Ref      string        `yaml:"ref,omitempty"`
	Remote   string        `yaml:"remote,omitempty"`
	Template string        `yaml:"template,omitempty"`
	File     Stringorslice `yaml:"file,omitempty"`
}

// UnmarshalYAML implements the unmarshal interface.
func (v *Include) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 string
	var out2 = struct {
		Local    string        `yaml:"local,omitempty"`
		Project  string        `yaml:"project,omitempty"`
		Ref      string        `yaml:"ref,omitempty"`
		Remote   string        `yaml:"remote,omitempty"`
		Template string        `yaml:"template,omitempty"`
		File     Stringorslice `yaml:"file,omitempty"`
	}{}

	if err := unmarshal(&out1); err == nil {
		v.Local = out1
		return nil
	}

	if err := unmarshal(&out2); err == nil {
		v.Local = out2.Local
		v.Project = out2.Project
		v.Ref = out2.Ref
		v.Remote = out2.Remote
		v.Template = out2.Template
		v.File = out2.File
		return nil
	}

	return errors.New("failed to unmarshal include")
}
