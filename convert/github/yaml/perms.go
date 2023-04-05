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
	Permissions struct {
		Actions            string `yaml:"actions,omitempty"`
		Checks             string `yaml:"checks,omitempty"`
		Contents           string `yaml:"contents,omitempty"`
		Deployments        string `yaml:"deployments,omitempty"`
		IDToken            string `yaml:"id-token,omitempty"`
		Issues             string `yaml:"issues,omitempty"`
		Discussions        string `yaml:"discussions,omitempty"`
		Packages           string `yaml:"packages,omitempty"`
		Pages              string `yaml:"pages,omitempty"`
		PullRequests       string `yaml:"pull-requests,omitempty"`
		RepositoryProjects string `yaml:"repository-projects,omitempty"`
		SecurityEvents     string `yaml:"security-events,omitempty"`
		Statuses           string `yaml:"statuses,omitempty"`
		ReadAll            bool   `yaml:"-"`
		WriteAll           bool   `yaml:"-"`
	}

	// temporary structure used to marshal and
	// unmarshal permissions.
	permissions struct {
		Actions            string `yaml:"actions,omitempty"`
		Checks             string `yaml:"checks,omitempty"`
		Contents           string `yaml:"contents,omitempty"`
		Deployments        string `yaml:"deployments,omitempty"`
		IDToken            string `yaml:"id-token,omitempty"`
		Issues             string `yaml:"issues,omitempty"`
		Discussions        string `yaml:"discussions,omitempty"`
		Packages           string `yaml:"packages,omitempty"`
		Pages              string `yaml:"pages,omitempty"`
		PullRequests       string `yaml:"pull-requests,omitempty"`
		RepositoryProjects string `yaml:"repository-projects,omitempty"`
		SecurityEvents     string `yaml:"security-events,omitempty"`
		Statuses           string `yaml:"statuses,omitempty"`
		ReadAll            bool   `yaml:"-"`
		WriteAll           bool   `yaml:"-"`
	}
)

// UnmarshalYAML implements the unmarshal interface for WorkflowTriggers.
func (v *Permissions) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 string
	var out2 permissions

	if err := unmarshal(&out1); err == nil {
		switch out1 {
		case "read-all":
			v.ReadAll = true
		case "write-all":
			v.WriteAll = true
		}
		return nil
	}
	if err := unmarshal(&out2); err == nil {
		*v = (Permissions)(out2)
		return nil
	}
	return errors.New("failed to unmarshal permissions")
}

// MarshalYAML implements the marshal interface.
func (v *Permissions) MarshalYAML() (interface{}, error) {
	switch {
	case v.ReadAll:
		return "read-all", nil
	case v.WriteAll:
		return "write-all", nil
	default:
		return permissions(*v), nil
	}
}
