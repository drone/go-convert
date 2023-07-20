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
	"strings"
)

type (
	// Pipeline defines a gitlab pipeline.
	Pipeline struct {
		Artifacts    *Artifacts `yaml:"artifacts,omitempty"`
		Schema       string
		Default      *Default
		Include      []*Include
		Image        *Image
		Services     []*Image
		BeforeScript []string
		AfterScript  []string
		Variables    map[string]*Variable
		Cache        *Cache
		Stages       []string
		Pages        *Job
		Workflow     *Workflow
		Jobs         map[string]*Job
		TemplateJobs map[string]*Job `yaml:"-"`
	}

	// pipeline is a temporary structure to parse the pipeline.
	pipeline struct {
		Artifacts    *Artifacts           `yaml:"artifacts,omitempty"`
		Schema       string               `yaml:"$schema,omitempty"`
		Default      *Default             `yaml:"default,omitempty"`
		Include      []*Include           `yaml:"include,omitempty"`
		Image        *Image               `yaml:"image,omitempty"`
		Services     []*Image             `yaml:"services,omitempty"`
		BeforeScript []string             `yaml:"before_script,omitempty"`
		AfterScript  []string             `yaml:"after_script,omitempty"`
		Variables    map[string]*Variable `yaml:"variables,omitempty"`
		Cache        *Cache               `yaml:"cache,omitempty"`
		Stages       []string             `yaml:"stages,omitempty"`
		Pages        *Job                 `yaml:"pages,omitempty"`
		Workflow     *Workflow            `yaml:"workflow,omitempty"`
		Jobs         map[string]*Job      `yaml:",inline"`
	}

	// Default defines global pipeline defaults.
	Default struct {
		After         []string   `yaml:"after_script,omitempty"`
		Before        []string   `yaml:"before_script,omitempty"`
		Artifacts     *Artifacts `yaml:"artifacts,omitempty"`
		Cache         *Cache     `yaml:"cache,omitempty"`
		Image         *Image     `yaml:"image,omitempty"`
		Interruptible bool       `yaml:"interruptible,omitempty"`
		Retry         *Retry     `yaml:"retry,omitempty"`
		Services      []*Image   `yaml:"services,omitempty"`
		Tags          []string   `yaml:"tags,omitempty"`
		Timeout       string     `yaml:"duration,omitempty"`
	}
)

func (p *Pipeline) UnmarshalYAML(unmarshal func(interface{}) error) error {
	out := new(pipeline)
	if err := unmarshal(&out); err != nil {
		return err
	}
	p.Artifacts = out.Artifacts
	p.Schema = out.Schema
	p.Default = out.Default
	p.Include = out.Include
	p.Image = out.Image
	p.Services = out.Services
	p.BeforeScript = out.BeforeScript
	p.AfterScript = out.AfterScript
	p.Variables = out.Variables
	p.Cache = out.Cache
	p.Stages = out.Stages
	p.Pages = out.Pages
	p.Workflow = out.Workflow
	p.Jobs = out.Jobs
	hasTemplateJobs := false

	for k := range out.Jobs {
		if strings.HasPrefix(k, ".") {
			hasTemplateJobs = true
			break
		}
	}

	if hasTemplateJobs {
		p.TemplateJobs = make(map[string]*Job) // Initialize as an empty map only if there are template jobs

		for k, v := range out.Jobs {
			if strings.HasPrefix(k, ".") {
				delete(p.Jobs, k)
				p.TemplateJobs[k] = v
			}
		}
	}
	return nil
}

// MarshalYAML implements yaml marshalling.
func (p *Pipeline) MarshalYAML() (interface{}, error) {
	out := new(pipeline)
	out.Artifacts = p.Artifacts
	out.Schema = p.Schema
	out.Default = p.Default
	out.Include = p.Include
	out.Image = p.Image
	out.Services = p.Services
	out.BeforeScript = p.BeforeScript
	out.AfterScript = p.AfterScript
	out.Variables = p.Variables
	out.Cache = p.Cache
	out.Stages = p.Stages
	out.Pages = p.Pages
	out.Workflow = p.Workflow
	out.Jobs = p.Jobs

	for k, v := range p.TemplateJobs {
		out.Jobs[k] = v
	}
	return out, nil
}
