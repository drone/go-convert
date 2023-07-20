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

type Needs struct {
	Items []*Need
}

// UnmarshalYAML implements the unmarshal interface.
func (v *Needs) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 *Need
	var out2 []*Need
	if err := unmarshal(&out1); err == nil {
		v.Items = append(v.Items, out1)
		return nil
	}
	if err := unmarshal(&out2); err == nil {
		v.Items = append(v.Items, out2...)
		return nil
	}
	return errors.New("failed to unmarshal needs list")
}

func (v *Needs) MarshalYAML() (interface{}, error) {
	if v.Items == nil || len(v.Items) == 0 {
		return []*Need{}, nil
	}
	return v.Items, nil
}

type Need struct {
	Job       string `yaml:"job,omitempty"`
	Ref       string `yaml:"ref,omitempty"`
	Project   string `yaml:"project,omitempty"`
	Pipeline  string `yaml:"pipeline,omitempty"`
	Artifacts bool   `yaml:"artifacts,omitempty"`
	Optional  bool   `yaml:"optional,omitempty"`
}

// UnmarshalYAML implements the unmarshal interface.
func (v *Need) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 string
	var out2 = struct {
		Job       string `yaml:"job,omitempty"`
		Ref       string `yaml:"ref,omitempty"`
		Project   string `yaml:"project,omitempty"`
		Pipeline  string `yaml:"pipeline,omitempty"`
		Artifacts bool   `yaml:"artifacts,omitempty"`
		Optional  bool   `yaml:"optional,omitempty"`
	}{}

	if err := unmarshal(&out1); err == nil {
		v.Job = out1
		return nil
	}

	if err := unmarshal(&out2); err == nil {
		v.Job = out2.Job
		v.Ref = out2.Ref
		v.Project = out2.Project
		v.Pipeline = out2.Pipeline
		v.Artifacts = out2.Artifacts
		v.Optional = out2.Optional
		return nil
	}

	return errors.New("failed to unmarshal needs")
}
