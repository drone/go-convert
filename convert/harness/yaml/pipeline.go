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

import "github.com/drone/go-convert/internal/flexible"

type (
	// Config defines resource configuration.
	Config struct {
		Pipeline Pipeline `json:"pipeline" yaml:"pipeline"`
	}

	// Pipeline defines a pipeline.
	Pipeline struct {
		ID                string              `json:"identifier,omitempty"        yaml:"identifier,omitempty"`
		Name              string              `json:"name,omitempty"              yaml:"name,omitempty"`
		Desc              string              `json:"description,omitempty"       yaml:"description,omitempty"`
		Account           string              `json:"accountIdentifier,omitempty" yaml:"accountIdentifier,omitempty"`
		Project           string              `json:"projectIdentifier,omitempty" yaml:"projectIdentifier,omitempty"`
		Org               string              `json:"orgIdentifier,omitempty"     yaml:"orgIdentifier,omitempty"`
		Props             Properties          `json:"properties,omitempty"        yaml:"properties,omitempty"`
		Stages            []*Stages           `json:"stages,omitempty"            yaml:"stages"`
		Variables         []*Variable         `json:"variables,omitempty"         yaml:"variables,omitempty"`
		Tags              map[string]string   `json:"tags,omitempty"              yaml:"tags,omitempty"`
		FlowControl       *FlowControl        `json:"flowControl,omitempty"       yaml:"flowControl,omitempty"`
		NotificationRules []*NotificationRule `json:"notificationRules,omitempty" yaml:"notificationRules,omitempty"`
		DelegateSelectors *flexible.Field[[]string] `json:"delegateSelectors,omitempty" yaml:"delegateSelectors,omitempty"`
	}

	FlowControl struct {
		Barriers []*Barrier `json:"barriers,omitempty" yaml:"barriers,omitempty"`
	}

	Barrier struct {
		Name       string `json:"name,omitempty" yaml:"name,omitempty"`
		Identifier string `json:"identifier,omitempty" yaml:"identifier,omitempty"`
	}

	// Properties defines pipeline properties.
	Properties struct {
		CI CI `json:"ci,omitempty" yaml:"ci"`
	}

	// CI defines CI pipeline properties.
	CI struct {
		Codebase *Codebase `json:"codebase,omitempty" yaml:"codebase,omitempty"`
	}

	// BuildIntelligence defines the build intelligence settings.
	BuildIntelligence struct {
		Enabled *flexible.Field[bool] `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	}

	// Cache defines the cache settings.
	Cache struct {
		Enabled bool     `json:"enabled,omitempty" yaml:"enabled,omitempty"`
		Key     string   `json:"key,omitempty"     yaml:"key,omitempty"`
		Paths   []string `json:"paths,omitempty"   yaml:"paths,omitempty"`
		Policy  string   `json:"policy,omitempty"  yaml:"policy,omitempty"`
		Override bool    `json:"override,omitempty"  yaml:"override,omitempty"`
	}

	// Codebase defines a codebase.
	Codebase struct {
		Name  string `json:"repoName,omitempty"     yaml:"repoName,omitempty"`
		Conn  string `json:"connectorRef,omitempty" yaml:"connectorRef,omitempty"`
		Build flexible.Field[Build] `json:"build,omitempty"        yaml:"build,omitempty"` // branch|tag
		Depth *flexible.Field[int64] `json:"depth,omitempty"        yaml:"depth,omitempty"`
		SslVerify *flexible.Field[bool] `json:"sslVerify,omitempty"  yaml:"sslVerify,omitempty"`
		PrCloneStrategy string `json:"prCloneStrategy,omitempty" yaml:"prCloneStrategy,omitempty"`
		Resources            *Resources                  `json:"resources,omitempty"          yaml:"resources,omitempty"`
		Lfs                  *flexible.Field[bool]        `json:"lfs,omitempty"                yaml:"lfs,omitempty"`
		Debug                *flexible.Field[bool]        `json:"debug,omitempty"              yaml:"debug,omitempty"`
		FetchTags            *flexible.Field[bool]        `json:"fetchTags,omitempty"          yaml:"fetchTags,omitempty"`
		PersistCredentials   *flexible.Field[bool]        `json:"persistCredentials,omitempty" yaml:"persistCredentials,omitempty"`
		SparseCheckout       []string                    `json:"sparseCheckout,omitempty"     yaml:"sparseCheckout,omitempty"`
		SubmoduleStrategy    string 					 `json:"submoduleStrategy,omitempty"  yaml:"submoduleStrategy,omitempty"`
		CloneDirectory       string                      `json:"cloneDirectory,omitempty"     yaml:"cloneDirectory,omitempty"`
		PreFetchCommand      string                      `json:"preFetchCommand,omitempty"    yaml:"preFetchCommand,omitempty"`
	}

	Build struct {
		Type string `json:"type,omitempty" yaml:"type,omitempty"`
		Spec BuildSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
	}

	BuildSpec struct {
		Branch string `json:"branch,omitempty" yaml:"branch,omitempty"`
		Tag string `json:"tag,omitempty" yaml:"tag,omitempty"`
		Number *flexible.Field[int] `json:"number,omitempty" yaml:"number,omitempty"`
		CommitSha string `json:"commitSha,omitempty" yaml:"commitSha,omitempty"`
	}

	Stages struct {
		Stage    *Stage    `json:"stage,omitempty"    yaml:"stage,omitempty"`
		Parallel []*Stages `json:"parallel,omitempty" yaml:"parallel,omitempty"`
	}

	Platform struct {
		OS   string `json:"os,omitempty"   yaml:"os,omitempty"`
		Arch string `json:"arch,omitempty" yaml:"arch,omitempty"`
	}

	Variable struct {
		Name  string `json:"name,omitempty"  yaml:"name,omitempty"`
		Type  string `json:"type,omitempty"  yaml:"type,omitempty"` // Secret|String|Number
		Value interface{} `json:"value,omitempty" yaml:"value,omitempty"`
		Required bool `json:"required,omitempty" yaml:"required,omitempty"`
		Default interface{} `json:"default,omitempty" yaml:"default,omitempty"`
	}

	Execution struct {
		Steps []*Steps `json:"steps,omitempty" yaml:"steps,omitempty"` // Un-necessary
	}

	Steps struct {
		Step      *Step      `json:"step,omitempty" yaml:"step,omitempty"` // Un-necessary
		Parallel  []*Steps   `json:"parallel,omitempty" yaml:"parallel,omitempty"`
		StepGroup *StepGroup `json:"stepGroup,omitempty" yaml:"stepGroup,omitempty"`
	}

	Report struct {
		Type string       `json:"type" yaml:"type,omitempty"` // JUnit|JUnit
		Spec *ReportJunit `json:"spec" yaml:"spec,omitempty"` // TODO
	}

	STOTarget struct {
		Type      string `json:"type" yaml:"type,omitempty"`
		Detection string `json:"detection" yaml:"detection,omitempty"`
	}

	STOAdvanced struct {
		Log *STOAdvancedLog `json:"log" yaml:"log,omitempty"`
	}

	STOAdvancedLog struct {
		Level string `json:"level" yaml:"level,omitempty"`
	}

	STOImage struct {
		Tag  string `json:"tag,omitempty"   yaml:"tag,omitempty"`
		Name string `json:"name,omitempty"         yaml:"name,omitempty"`
		Type string `json:"type,omitempty"         yaml:"type,omitempty"`
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

	Output struct {
		Name  string `json:"name,omitempty"      yaml:"name,omitempty"`
		Type  string `json:"type,omitempty" yaml:"type,omitempty"`
		Value string `json:"value,omitempty" yaml:"value,omitempty"`
	}

	ServiceSpec struct {
		Env        map[string]string `json:"envVariables,omitempty"   yaml:"envVariables,omitempty"`
		Entrypoint []string          `json:"entrypoint,omitempty"     yaml:"entrypoint,omitempty"`
		Args       []string          `json:"args,omitempty"           yaml:"args,omitempty"`
		Conn       string            `json:"connectorRef,omitempty"   yaml:"connectorRef,omitempty"`
		Image      string            `json:"image,omitempty"          yaml:"image,omitempty"`
		Resources  *Resources        `json:"resources,omitempty"      yaml:"resources,omitempty"`
		Privileged *flexible.Field[bool]              `json:"privileged,omitempty"     yaml:"privileged,omitempty"`
	}

	Resources struct {
		Limits *Limits `json:"limits,omitempty" yaml:"limits,omitempty"`
	}

	Limits struct {
		Memory *flexible.Field[*BytesSize] `json:"memory,omitempty" yaml:"memory,omitempty"`
		CPU    *flexible.Field[*MilliSize] `json:"cpu,omitempty"    yaml:"cpu,omitempty"`
	}
)

// GetCPUString returns the CPU value as a string, handling both expressions and parsed values
func (l *Limits) GetCPUString() string {
    if l == nil || l.CPU == nil {
        return ""
    }
    
	if expr, ok := l.CPU.AsString(); ok {
		return expr
	}
    
    
    // If it's a struct value, convert to string with "m" suffix
    if cpu, ok := l.CPU.AsStruct(); ok {
        return cpu.String() + "m"
    }
    
    return ""
}

// GetMemoryString returns the memory value as a string, handling both expressions and parsed values
func (l *Limits) GetMemoryString() string {
    if l == nil || l.Memory == nil {
        return ""
    }
    
	if expr, ok := l.Memory.AsString(); ok {
		return expr
	}
    
    // If it's a struct value, convert to string
    if memory, ok := l.Memory.AsStruct(); ok {
        return memory.String()
    }
    
    return ""
}