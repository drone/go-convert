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
	Config struct {
		Commands   map[string]*Command   `yaml:"commands,omitempty"`
		Executors  map[string]*Executor  `yaml:"executors,omitempty"`
		Jobs       map[string]*Job       `yaml:"jobs,omitempty"`
		Orbs       map[string]*Orb       `yaml:"orbs,omitempty"`
		Parameters map[string]*Parameter `yaml:"parameters,omitempty"`
		Setup      bool                  `yaml:"setup,omitempty"`
		Workflows  *Workflows            `yaml:"workflows,omitempty"`
		Version    string                `yaml:"version,omitempty"`
	}

	Docker struct {
		Auth        *DockerAuth       `yaml:"auth,omitempty"`
		AuthAWS     *DockerAuthAWS    `yaml:"aws_auth,omitempty"`
		Command     Stringorslice     `yaml:"command,omitempty"`
		Entrypoint  Stringorslice     `yaml:"entrypoint,omitempty"`
		Environment map[string]string `yaml:"environment,omitempty"`
		Image       string            `yaml:"image,omitempty"`
		Name        string            `yaml:"name,omitempty"`
		User        string            `yaml:"user,omitempty"`
	}

	Command struct {
		Description string                `yaml:"description,omitempty"`
		Parameters  map[string]*Parameter `yaml:"parameters,omitempty"`
		Steps       []*Step               `yaml:"steps,omitempty"`
	}

	DockerAuth struct {
		Username string `yaml:"username,omitempty"`
		Password string `yaml:"password,omitempty"`
	}

	DockerAuthAWS struct {
		AccessKey string `yaml:"aws_access_key_id,omitempty"`
		SecretKey string `yaml:"aws_secret_access_key,omitempty"`
	}

	Executor struct {
		Docker        []*Docker         `yaml:"docker,omitempty"`
		ResourceClass string            `yaml:"resource_class,omitempty"`
		Machine       *Machine          `yaml:"machine,omitempty"`
		Macos         *Macos            `yaml:"macos,omitempty"`
		Windows       interface{}       `yaml:"widows,omitempty"`
		Shell         string            `yaml:"shell,omitempty"`
		WorkingDir    string            `yaml:"working_directory,omitempty"`
		Environment   map[string]string `yaml:"environment,omitempty"`
	}

	Filters struct {
		Branches *Filter `yaml:"branches,omitempty"`
		Tags     *Filter `yaml:"tags,omitempty"`
	}

	Filter struct {
		Only   Stringorslice `yaml:"only,omitempty"`
		Ignore Stringorslice `yaml:"ignore,omitempty"`
	}

	Job struct {
		Branches      *Filter               `yaml:"branches,omitempty"`
		Docker        []*Docker             `yaml:"docker,omitempty"`
		Environment   map[string]string     `yaml:"environment,omitempty"`
		Executor      *JobExecutor          `yaml:"executor,omitempty"`
		IPRanges      bool                  `yaml:"circleci_ip_ranges,omitempty"`
		Machine       *Machine              `yaml:"machine,omitempty"`
		Macos         *Macos                `yaml:"macos,omitempty"`
		Parallelism   int                   `yaml:"parallelism,omitempty"`
		Parameters    map[string]*Parameter `yaml:"parameters,omitempty"`
		ResourceClass string                `yaml:"resource_class,omitempty"`
		Shell         string                `yaml:"shell,omitempty"`
		Steps         []*Step               `yaml:"steps,omitempty"`
		WorkingDir    string                `yaml:"working_directory,omitempty"`
	}

	Macos struct {
		Xcode string `yaml:"xcode,omitempty"`
	}

	Matrix struct {
		Alias      string                   `yaml:"alias,omitempty"`
		Exclude    []map[string]interface{} `yaml:"exclude,omitempty"`
		Parameters map[string][]interface{} `yaml:"parameters,omitempty"`
	}

	Parameter struct {
		Description string      `yaml:"description,omitempty"`
		Type        string      `yaml:"type,omitempty"` // string, boolean, integer, enum, executor, steps
		Default     interface{} `yaml:"default,omitempty"`
	}

	Schedule struct {
		Cron    string   `yaml:"cron,omitempty"`
		Filters *Filters `yaml:"filters,omitempty"`
	}

	Trigger struct {
		Schedule *Schedule `yaml:"schedule,omitempty"`
	}
)
