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

type (
	Pipeline struct {
		Concurrency *Concurrency      `yaml:"concurrency,omitempty"`
		Environment map[string]string `yaml:"env,omitempty"`
		Name        string            `yaml:"name,omitempty"`
		On          *On               `yaml:"on,omitempty"`
		Jobs        map[string]Job    `yaml:"jobs,omitempty"`
	}

	Event struct {
		Types []string `yaml:"types,omitempty"`
	}

	Input struct {
		Description string      `yaml:"description,omitempty"`
		Required    bool        `yaml:"required,omitempty"`
		Default     interface{} `yaml:"default,omitempty"`
		Type        string      `yaml:"type,omitempty"`
		Options     interface{} `yaml:"options,omitempty"`
	}

	Job struct {
		RunsOn      string              `yaml:"runs-on,omitempty"`
		Container   string              `yaml:"container,omitempty"`
		Services    map[string]*Service `yaml:"services,omitempty"`
		Steps       []*Step             `yaml:"steps,omitempty"`
		Environment map[string]string   `yaml:"env,omitempty"`
		If          string              `yaml:"if,omitempty"`
		Strategy    *Strategy           `yaml:"strategy,omitempty"`
	}

	Matrix struct {
		Matrix  map[string][]string      `yaml:",inline"`
		Include []map[string]interface{} `yaml:"include,omitempty"`
		Exclude []map[string]interface{} `yaml:"exclude,omitempty"`
	}

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
	}

	PullRequest struct {
		Branches        []string `yaml:"branches,omitempty"`
		BranchesIgnore  []string `yaml:"branches-ignore,omitempty"`
		Paths           []string `yaml:"paths,omitempty"`
		PathsIgnore     []string `yaml:"paths-ignore,omitempty"`
		Tags            []string `yaml:"tags,omitempty"`
		TagsIgnore      []string `yaml:"tags-ignore,omitempty"`
		Types           []string `yaml:"types,omitempty"`
		ReviewApproved  bool     `yaml:"review-approved,omitempty"`
		ReviewDismissed bool     `yaml:"review-dismissed,omitempty"`
	}

	PullRequestTarget struct {
		Branches       []string `yaml:"branches,omitempty"`
		BranchesIgnore []string `yaml:"branches-ignore,omitempty"`
		Types          []string `yaml:"types,omitempty"`
	}

	Push struct {
		Branches       []string `yaml:"branches,omitempty"`
		BranchesIgnore []string `yaml:"branches-ignore,omitempty"`
		Paths          []string `yaml:"paths,omitempty"`
		PathsIgnore    []string `yaml:"paths-ignore,omitempty"`
		Tags           []string `yaml:"tags,omitempty"`
		TagsIgnore     []string `yaml:"tags-ignore,omitempty"`
	}

	Service struct {
		Image    string            `yaml:"image,omitempty"`
		Env      map[string]string `yaml:"env,omitempty"`
		Ports    []string          `yaml:"ports,omitempty"`
		Options  []string          `yaml:"options,omitempty"`
		Volumes  []string          `yaml:"volumes,omitempty"`
		Networks []string          `yaml:"networks,omitempty"`
	}

	Step struct {
		Name        string                 `yaml:"name,omitempty"`
		Uses        string                 `yaml:"uses,omitempty"`
		With        map[string]interface{} `yaml:"with,omitempty"`
		Run         string                 `yaml:"run,omitempty"`
		If          string                 `yaml:"if,omitempty"`
		Environment map[string]string      `yaml:"env,omitempty"`
	}

	Strategy struct {
		Matrix *Matrix `yaml:"matrix,omitempty"`
	}

	WorkflowCall struct {
		Workflows []string                       `yaml:"workflows,omitempty"`
		Inputs    map[string]interface{}         `yaml:"inputs,omitempty"`
		Outputs   map[string]interface{}         `yaml:"outputs,omitempty"`
		Secrets   map[string]*WorkflowCallSecret `yaml:"secrets,omitempty"`
	}

	WorkflowCallSecret struct {
		Description string `yaml:"description,omitempty"`
		Required    bool   `yaml:"required,omitempty"`
	}

	WorkflowDispatch struct {
		Inputs map[string]*Input `yaml:"inputs,omitempty"`
	}

	WorkflowRun struct {
		Workflows      []string `yaml:"workflows,omitempty"`
		Types          []string `yaml:"types,omitempty"`
		Branches       []string `yaml:"branches,omitempty"`
		BranchesIgnore []string `yaml:"branches-ignore,omitempty"`
	}
)
