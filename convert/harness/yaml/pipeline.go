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
	// Config defines resource configuration.
	Config struct {
		Pipeline Pipeline `json:"pipeline" yaml:"pipeline"`
	}

	// Pipeline defines a pipeline.
	Pipeline struct {
		ID        string      `json:"identifier,omitempty"        yaml:"identifier,omitempty"`
		Name      string      `json:"name,omitempty"              yaml:"name,omitempty"`
		Desc      string      `json:"description,omitempty"       yaml:"description,omitempty"`
		Account   string      `json:"accountIdentifier,omitempty" yaml:"accountIdentifier,omitempty"`
		Project   string      `json:"projectIdentifier,omitempty" yaml:"projectIdentifier,omitempty"`
		Org       string      `json:"orgIdentifier,omitempty"     yaml:"orgIdentifier,omitempty"`
		Props     Properties  `json:"properties,omitempty"        yaml:"properties,omitempty"`
		Stages    []*Stages   `json:"stages,omitempty"            yaml:"stages"`
		Variables []*Variable `json:"variables,omitempty"         yaml:"variables,omitempty"`
	}

	// Properties defines pipeline properties.
	Properties struct {
		CI CI `json:"ci,omitempty" yaml:"ci"`
	}

	// CI defines CI pipeline properties.
	CI struct {
		Codebase Codebase `json:"codebase,omitempty" yaml:"codebase,omitempty"`
	}

	// Cache defines the cache settings.
	Cache struct {
		Enabled bool     `json:"enabled,omitempty" yaml:"enabled,omitempty"`
		Key     string   `json:"key,omitempty"     yaml:"key,omitempty"`
		Paths   []string `json:"paths,omitempty"   yaml:"paths,omitempty"`
	}

	// Codebase defines a codebase.
	Codebase struct {
		Name  string `json:"repoName,omitempty"     yaml:"repoName,omitempty"`
		Conn  string `json:"connectorRef,omitempty" yaml:"connectorRef,omitempty"`
		Build string `json:"build,omitempty"        yaml:"build,omitempty"` // branch|tag
	}

	Stages struct {
		Stage    *Stage    `json:"stage,omitempty"    yaml:"stage,omitempty"`
		Parallel []*Stages `json:"parallel,omitempty" yaml:"parallel,omitempty"`
	}

	// Infrastructure provides pipeline infrastructure.
	Infrastructure struct {
		Type string     `json:"type,omitempty"          yaml:"type,omitempty"`
		From string     `json:"useFromStage,omitempty"  yaml:"useFromStage,omitempty"` // this is also weird
		Spec *InfraSpec `json:"spec,omitempty"          yaml:"spec,omitempty"`
	}

	// InfraSpec describes pipeline infastructure.
	InfraSpec struct {
		Namespace string `json:"namespace,omitempty"    yaml:"namespace,omitempty"`
		Conn      string `json:"connectorRef,omitempty" yaml:"connectorRef,omitempty"`
	}

	Platform struct {
		OS   string `json:"os,omitempty"   yaml:"os,omitempty"`
		Arch string `json:"arch,omitempty" yaml:"arch,omitempty"`
	}

	Runtime struct {
		Type string      `json:"type,omitempty"   yaml:"type,omitempty"`
		Spec interface{} `json:"spec,omitempty"   yaml:"spec,omitempty"`
	}

	Variable struct {
		Name  string `json:"name,omitempty"  yaml:"name,omitempty"`
		Type  string `json:"type,omitempty"  yaml:"type,omitempty"` // Secret|Text
		Value string `json:"value,omitempty" yaml:"value,omitempty"`
	}

	Execution struct {
		Steps []*Steps `json:"steps,omitempty" yaml:"steps,omitempty"` // Un-necessary
	}

	Steps struct {
		Step     *Step    `json:"step,omitempty" yaml:"step,omitempty"` // Un-necessary
		Parallel []*Steps `json:"parallel,omitempty" yaml:"parallel,omitempty"`
	}

	Report struct {
		Type string       `json:"type" yaml:"type,omitempty"` // JUnit|JUnit
		Spec *ReportJunit `json:"spec" yaml:"spec,omitempty"` // TODO
	}

	ReportJunit struct {
		Paths []string `json:"paths" yaml:"paths,omitempty"`
	}

	Service struct {
		ID   string       `json:"identifier,omitempty"   yaml:"identifier,omitempty"`
		Name string       `json:"name,omitempty"         yaml:"name,omitempty"`
		Type string       `json:"type,omitempty"         yaml:"type,omitempty"` // Service
		Desc string       `json:"description,omitempty"  yaml:"description,omitempty"`
		Spec *ServiceSpec `json:"spec,omitempty"         yaml:"spec,omitempty"`
	}

	ServiceSpec struct {
		Env        map[string]string `json:"envVariables,omitempty"   yaml:"envVariables,omitempty"`
		Entrypoint []string          `json:"entrypoint,omitempty"     yaml:"entrypoint,omitempty"`
		Args       []string          `json:"args,omitempty"           yaml:"args,omitempty"`
		Conn       string            `json:"connectorRef,omitempty"   yaml:"connectorRef,omitempty"`
		Image      string            `json:"image,omitempty"          yaml:"image,omitempty"`
		Resources  *Resources        `json:"resources,omitempty"      yaml:"resources,omitempty"`
	}

	Resources struct {
		Limits Limits `json:"limits,omitempty" yaml:"limits,omitempty"`
	}

	Limits struct {
		Memory BytesSize `json:"memory,omitempty" yaml:"memory,omitempty"`
		CPU    MilliSize `json:"cpu,omitempty"    yaml:"cpu,omitempty"` // TODO
	}
)
