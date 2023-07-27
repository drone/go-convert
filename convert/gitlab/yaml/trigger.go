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

type Trigger struct {
	Project  string   `yaml:"-"`
	Branch   string   `yaml:"branch,omitempty"`
	Include  string   `yaml:"include,omitempty"`
	Strategy string   `yaml:"strategy,omitempty"`
	Forward  *Forward `yaml:"forward,omitempty"`
}

type Forward struct {
	YamlVariables     *bool `yaml:"yaml_variables,omitempty"`
	PipelineVariables bool  `yaml:"pipeline_variables,omitempty"`
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (v *Trigger) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 string
	var out2 = struct {
		Project  string   `yaml:"project,omitempty"`
		Branch   string   `yaml:"branch,omitempty"`
		Include  string   `yaml:"include,omitempty"`
		Strategy string   `yaml:"strategy,omitempty"`
		Forward  *Forward `yaml:"forward,omitempty"`
	}{}

	if err := unmarshal(&out1); err == nil {
		v.Project = out1
		return nil
	}

	if err := unmarshal(&out2); err == nil {
		v.Project = out2.Project
		v.Branch = out2.Branch
		v.Include = out2.Include
		v.Strategy = out2.Strategy
		v.Forward = out2.Forward
		return nil
	}

	return errors.New("failed to unmarshal trigger")
}

func (v *Trigger) MarshalYAML() (interface{}, error) {
	// Always marshal as a struct
	return struct {
		Project  string   `yaml:"project,omitempty"`
		Branch   string   `yaml:"branch,omitempty"`
		Include  string   `yaml:"include,omitempty"`
		Strategy string   `yaml:"strategy,omitempty"`
		Forward  *Forward `yaml:"forward,omitempty"`
	}{
		Project:  v.Project,
		Branch:   v.Branch,
		Include:  v.Include,
		Strategy: v.Strategy,
		Forward:  v.Forward,
	}, nil
}
