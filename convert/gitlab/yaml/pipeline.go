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
	// Pipeline defines a gitlab pipeline.
	Pipeline struct {
		Default   *Default             `yaml:"default,omitempty"`
		Include   []*Include           `yaml:"include,omitempty"`
		Image     *Image               `yaml:"image,omitempty"`
		Jobs      map[string]*Job      `yaml:",inline"`
		Stages    []string             `yaml:"stages,omitempty"`
		Variables map[string]*Variable `yaml:"variables,omitempty"`
		Workflow  *Workflow            `yaml:"workflow,omitempty"`
	}

	// Default defines global pipeline defaults.
	Default struct {
		After         Stringorslice `yaml:"after_script,omitempty"`
		Before        Stringorslice `yaml:"before_script,omitempty"`
		Artifacts     *Artifacts    `yaml:"artifacts,omitempty"`
		Cache         *Cache        `yaml:"cache,omitempty"`
		Image         *Image        `yaml:"image,omitempty"`
		Interruptible bool          `yaml:"interruptible,omitempty"`
		Retry         *Retry        `yaml:"retry,omitempty"`
		Services      []*Image      `yaml:"services,omitempty"`
		Tags          Stringorslice `yaml:"tags,omitempty"`
		Timeout       string        `yaml:"duration,omitempty"`
	}
)
