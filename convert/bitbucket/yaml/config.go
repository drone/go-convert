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
		Clone       *Clone       `yaml:"clone,omitempty"`
		Definitions *Definitions `yaml:"definitions,omitempty"`
		Image       *Image       `yaml:"image,omitempty"`
		Options     *Options     `yaml:"options,omitempty"`
		Pipelines   Pipelines    `yaml:"pipelines,omitempty"`
	}

	Clone struct {
		Depth      *Depth `yaml:"depth,omitempty"`
		Enabled    *bool  `yaml:"enabled,omitempty"`
		LFS        bool   `yaml:"lfs,omitempty"`
		SkipVerify bool   `yaml:"skip-ssl-verify,omitempty"`
	}

	Condition struct {
		Changesets *Changesets `yaml:"changesets,omitempty"`
	}

	Changesets struct {
		IncludePaths []string `yaml:"includePaths,omitempty"`
	}

	Definitions struct {
		Caches   map[string]*Cache   `yaml:"caches,omitempty"`
		Services map[string]*Service `yaml:"services,omitempty"`
	}

	Options struct {
		Docker  bool `yaml:"docker,omitempty"`
		MaxTime int  `yaml:"max-time,omitempty"`
		Size    Size `yaml:"size,omitempty"`
	}

	Parallel struct {
		FailFast bool     `yaml:"fail-fast,omitempty"`
		Steps    []*Steps `yaml:"steps,omitempty"`
	}

	Pipelines struct {
		Default      []*Steps            `yaml:"default,omitempty"`
		Branches     map[string][]*Steps `yaml:"branches,omitempty"`
		PullRequests map[string][]*Steps `yaml:"pull-requests,omitempty"`
		Tags         map[string][]*Steps `yaml:"tags,omitempty"`
		Custom       []*Steps            `yaml:"custom,omitempty"`
	}

	Service struct {
		Image     *Image            `yaml:"image,omitempty"`
		Memory    int               `yaml:"memory,omitempty"` // default 1024
		Type      string            `yaml:"type,omitempty"`
		Variables map[string]string `yaml:"variables,omitempty"`
	}

	Stage struct {
		Condition  *Condition `yaml:"condition,omitempty"`
		Deployment string     `yaml:"deployment,omitempty"` // test, staging, production
		Steps      []*Steps   `yaml:"steps,omitempty"`
		Name       string     `yaml:"name,omitempty"`
		Trigger    string     `yaml:"trigger,omitempty"`
	}

	Steps struct {
		Step     *Step     `yaml:"step,omitempty"`
		Stage    *Stage    `yaml:"stage,omitempty"`
		Parallel *Parallel `yaml:"parallel,omitempty"`
	}

	Step struct {
		Artifacts   *Artifacts    `yaml:"artifacts,omitempty"`
		Caches      []string      `yaml:"caches,omitempty"`
		Clone       *Clone        `yaml:"clone,omitempty"`
		Condition   *Condition    `yaml:"condition,omitempty"`
		Deployment  string        `yaml:"deployment,omitempty"` // test, staging, production
		FailFast    bool          `yaml:"fail-fast,omitempty"`
		Image       *Image        `yaml:"image,omitempty"`
		MaxTime     int           `yaml:"max-time,omitempty"`
		Name        string        `yaml:"name,omitempty"`
		Oidc        bool          `yaml:"oidc,omitempty"`
		RunsOn      Stringorslice `yaml:"runs-on,omitempty"`
		Script      []*Script     `yaml:"script"`
		ScriptAfter []*Script     `yaml:"after-script,omitempty"`
		Services    []string      `yaml:"services,omitempty"`
		Size        Size          `yaml:"size,omitempty"`
		Trigger     string        `yaml:"trigger,omitempty"` // automatic, manual
	}

	Variable struct {
		AllowedValues []string `yaml:"allowed-values,omitempty"`
		Default       string   `yaml:"default,omitempty"`
		Name          string   `yaml:"name,omitempty"`
	}

	AWS struct {
		AccessKey string `yaml:"access-key,omitempty"`
		SecretKey string `yaml:"secret-key,omitempty"`
		OIDCRole  string `yaml:"oidc-role,omitempty"`
	}
)
