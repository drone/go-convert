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
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

type (
	// Pipeline defines a gitlab pipeline.
	Pipeline struct {
		Default      *Default             `yaml:"default,omitempty"`
		Include      []*Include           `yaml:"include,omitempty"`
		Image        *Image               `yaml:"image,omitempty"`
		Jobs         map[string]*Job      `yaml:"jobs,omitempty"`
		TemplateJobs map[string]*Job      `yaml:"-"`
		Stages       []string             `yaml:"stages,omitempty"`
		Variables    map[string]*Variable `yaml:"variables,omitempty"`
		Workflow     *Workflow            `yaml:"workflow,omitempty"`
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

var globalKeys = map[string]struct{}{
	"after_script":  {},
	"artifacts":     {},
	"before_script": {},
	"cache":         {},
	"image":         {},
	"interruptible": {},
	"retry":         {},
	"services":      {},
	"tags":          {},
	"timeout":       {},
	"variables":     {},
	"workflow":      {},
	"stages":        {},
	"include":       {},
}

func (p *Pipeline) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var rawData map[string]interface{}
	if err := unmarshal(&rawData); err != nil {
		return err
	}

	if p.Default == nil {
		p.Default = &Default{}
	}

	for k, v := range rawData {
		// we check if the key is a global one
		if k == "default" {
			defaultYaml, err := yaml.Marshal(v)
			if err != nil {
				return err
			}
			if err := yaml.Unmarshal(defaultYaml, p.Default); err != nil {
				return err
			}
			// Remove the key to avoid processing it again
			delete(rawData, k)
		} else if _, isGlobal := globalKeys[k]; isGlobal {
			switch k {
			case "artifacts":
				artifactsYaml, err := yaml.Marshal(v)
				if err != nil {
					return err
				}
				if err := yaml.Unmarshal(artifactsYaml, &p.Default.Artifacts); err != nil {
					return err
				}
			case "image":
				imageYaml, err := yaml.Marshal(v)
				if err != nil {
					return err
				}
				if err := yaml.Unmarshal(imageYaml, &p.Default.Image); err != nil {
					return err
				}
			case "timeout":
				timeoutYaml, err := yaml.Marshal(v)
				if err != nil {
					return err
				}
				if err := yaml.Unmarshal(timeoutYaml, &p.Default.Timeout); err != nil {
					return err
				}
			case "before_script":
				beforeYaml, err := yaml.Marshal(v)
				if err != nil {
					return err
				}
				if err := yaml.Unmarshal(beforeYaml, &p.Default.Before); err != nil {
					return err
				}
			case "after_script":
				afterYaml, err := yaml.Marshal(v)
				if err != nil {
					return err
				}
				if err := yaml.Unmarshal(afterYaml, &p.Default.After); err != nil {
					return err
				}
			case "cache":
				cacheYaml, err := yaml.Marshal(v)
				if err != nil {
					return err
				}
				if err := yaml.Unmarshal(cacheYaml, &p.Default.Cache); err != nil {
					return err
				}
			case "services":
				servicesYaml, err := yaml.Marshal(v)
				if err != nil {
					return err
				}
				if err := yaml.Unmarshal(servicesYaml, &p.Default.Services); err != nil {
					return err
				}
			case "tags":
				tagsYaml, err := yaml.Marshal(v)
				if err != nil {
					return err
				}
				if err := yaml.Unmarshal(tagsYaml, &p.Default.Tags); err != nil {
					return err
				}
			case "interruptible":
				interruptibleYaml, err := yaml.Marshal(v)
				if err != nil {
					return err
				}
				if err := yaml.Unmarshal(interruptibleYaml, &p.Default.Interruptible); err != nil {
					return err
				}
			case "variables":
				variablesYaml, err := yaml.Marshal(v)
				if err != nil {
					return err
				}
				if p.Variables == nil {
					p.Variables = make(map[string]*Variable)
				}
				if err := yaml.Unmarshal(variablesYaml, &p.Variables); err != nil {
					return err
				}
			case "workflow":
				workflowYaml, err := yaml.Marshal(v)
				if err != nil {
					return err
				}
				if err := yaml.Unmarshal(workflowYaml, &p.Workflow); err != nil {
					return err
				}
			case "stages":
				stagesYaml, err := yaml.Marshal(v)
				if err != nil {
					return err
				}
				if err := yaml.Unmarshal(stagesYaml, &p.Stages); err != nil {
					return err
				}
			case "include":
				includeYaml, err := yaml.Marshal(v)
				if err != nil {
					return err
				}
				if err := yaml.Unmarshal(includeYaml, &p.Include); err != nil {
					return err
				}
			}

			// Remove the key to avoid processing it again as a job
			delete(rawData, k)
		}
	}

	// If Default is still empty, set it to nil
	if p.Default != nil && reflect.DeepEqual(*p.Default, Default{}) {
		p.Default = nil
	}

	p.Jobs = make(map[string]*Job)

	for k, v := range rawData {
		jobYaml, err := yaml.Marshal(v)
		if err != nil {
			return err
		}
		job := &Job{}
		if err := yaml.Unmarshal(jobYaml, job); err != nil {
			return err
		}
		// If the job name starts with a dot, it's a template job
		if strings.HasPrefix(k, ".") {
			if p.TemplateJobs == nil {
				p.TemplateJobs = make(map[string]*Job)
			}
			p.TemplateJobs[k] = job
		} else {
			p.Jobs[k] = job
		}
	}

	return nil
}
