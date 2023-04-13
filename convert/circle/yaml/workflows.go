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

type (
	Workflows struct {
		Version string               `yaml:"version,omitempty"`
		Items   map[string]*Workflow `yaml:",inline"`
	}

	Workflow struct {
		Jobs     []*WorkflowJob `yaml:"jobs,omitempty"`
		Triggers []*Trigger     `yaml:"triggers,omitempty"`
		Unless   *Logical       `yaml:"when,omitempty"`
		When     *Logical       `yaml:"unless,omitempty"`
	}

	WorkflowJob struct {
		Name     string
		Context  []string
		Filters  *Filters
		Matrix   *Matrix
		Type     string
		Requires []string
		Params   map[string]interface{} // custom params
	}

	workflowJob struct {
		Context  Stringorslice          `yaml:"context,omitempty"`
		Filters  *Filters               `yaml:"filters,omitempty"`
		Matrix   *Matrix                `yaml:"matrix,omitempty"`
		Type     string                 `yaml:"type,omitempty"` // approval
		Requires []string               `yaml:"requires,omitempty"`
		Params   map[string]interface{} `yaml:",inline"` // custom params
	}
)

// UnmarshalYAML implements the unmarshal interface.
func (v *WorkflowJob) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 string
	var out2 map[string]*workflowJob

	if err := unmarshal(&out1); err == nil {
		v.Name = out1
		return nil
	}

	if err := unmarshal(&out2); err == nil {
		if len(out2) == 0 {
			return errors.New("failed to unmarshal job")
		}
		for key, val := range out2 {
			v.Name = key
			v.Context = val.Context
			v.Filters = val.Filters
			v.Matrix = val.Matrix
			v.Type = val.Type
			v.Requires = val.Requires
			v.Params = val.Params
		}
		return nil
	}

	return errors.New("failed to unmarshal job")
}

// MarshalYAML implements the marshal interface.
func (v *WorkflowJob) MarshalYAML() (interface{}, error) {
	// if the structure is empty, output the string only
	if len(v.Context) == 0 && len(v.Params) == 0 && v.Filters == nil && v.Matrix == nil && v.Type == "" && len(v.Requires) == 0 {
		return v.Name, nil
	}
	return map[string]*workflowJob{
		v.Name: {
			Context:  v.Context,
			Filters:  v.Filters,
			Matrix:   v.Matrix,
			Type:     v.Type,
			Requires: v.Requires,
			Params:   v.Params,
		},
	}, nil
}
