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

// Job defines a gitlab job.
// https://docs.gitlab.com/ee/ci/yaml/#job-keywords
type Job struct {
	After         Stringorslice        `yaml:"after_script,omitempty"`
	Artifacts     *Artifacts           `yaml:"artifacts,omitempty"`
	AllowFailure  *AllowFailure        `yaml:"allow_failure,omitempty"`
	Before        Stringorslice        `yaml:"before_script,omitempty"`
	Cache         *Cache               `yaml:"cache,omitempty"`
	Environment   *Environment         `yaml:"environment,omitempty"`
	Extends       Stringorslice        `yaml:"extends,omitempty"`
	Image         *Image               `yaml:"image,omitempty"`
	Inherit       *Inherit             `yaml:"inherit,omitempty"`
	Interruptible bool                 `yaml:"interruptible,omitempty"`
	Needs         []*Needs             `yaml:"needs,omitempty"`
	Pages         *Job                 `yaml:"pages,omitempty"`
	Parallel      *Parallel            `yaml:"parallel,omitempty"`
	Release       *Release             `yaml:"release,omitempty"`
	ResourceGroup string               `yaml:"resource_group,omitempty"`
	Retry         *Retry               `yaml:"retry,omitempty"`
	Rules         []*Rule              `yaml:"rules,omitempty"`
	Script        Stringorslice        `yaml:"script,omitempty"`
	Secrets       map[string]*Secret   `yaml:"secrets,omitempty"`
	Services      []*Image             `yaml:"services,omitempty"`
	Stage         string               `yaml:"stage,omitempty"`
	Tags          Stringorslice        `yaml:"tags,omitempty"`
	Timeout       string               `yaml:"timeout,omitempty"`
	Trigger       *Trigger             `yaml:"trigger,omitempty"`
	Variables     map[string]*Variable `yaml:"variables,omitempty"`
	When          string               `yaml:"when,omitempty"` // on_success, manual, always, on_failure, delayed, never
}
