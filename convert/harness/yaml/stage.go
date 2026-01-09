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
	"encoding/json"
	"fmt"

	"github.com/drone/go-convert/internal/flexible"
)

type (
	// Stage defines a pipeline stage.
	Stage struct {
		ID          string      `json:"identifier,omitempty"   yaml:"identifier,omitempty"`
		Description string      `json:"description,omitempty"  yaml:"description,omitempty"`
		Name        string      `json:"name,omitempty"         yaml:"name,omitempty"`
		DelegateSelectors *flexible.Field[[]string] `json:"delegateSelectors,omitempty" yaml:"delegateSelectors,omitempty"`
		Spec        interface{} `json:"spec,omitempty"         yaml:"spec,omitempty"`
		Type        string      `json:"type,omitempty"         yaml:"type,omitempty"`
		Vars        []*Variable `json:"variables,omitempty"    yaml:"variables,omitempty"`
		When        *flexible.Field[StageWhen]  `json:"when,omitempty"         yaml:"when,omitempty"`
		Strategy    *Strategy   `json:"strategy,omitempty"     yaml:"strategy,omitempty"`
		FailureStrategies *flexible.Field[[]*FailureStrategy]   `json:"failureStrategies,omitempty" yaml:"failureStrategies,omitempty"`
	}

	StageCustom struct {
		Execution *Execution `json:"execution,omitempty" yaml:"execution,omitempty"`
		Environment *Environment `json:"environment,omitempty" yaml:"environment,omitempty"`
		Timeout string `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	}

	// StageApproval defines an approval stage.
	StageApproval struct {
		Execution *Execution `json:"execution,omitempty" yaml:"execution,omitempty"`
		Timeout string `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	}

	// StageCI defines a continuous integration stage.
	StageCI struct {
		BuildIntelligence *BuildIntelligence `json:"buildIntelligence,omitempty"   yaml:"buildIntelligence,omitempty"`
		Cache             *Cache             `json:"caching,omitempty"             yaml:"caching,omitempty"`
		Clone             bool               `json:"cloneCodebase,omitempty"       yaml:"cloneCodebase,omitempty"`
		Execution         Execution          `json:"execution,omitempty"           yaml:"execution,omitempty"`
		Infrastructure    *Infrastructure    `json:"infrastructure,omitempty"      yaml:"infrastructure,omitempty"`
		Platform          *Platform          `json:"platform,omitempty"            yaml:"platform,omitempty"`
		Runtime           *Runtime           `json:"runtime,omitempty"            yaml:"runtime,omitempty"`
		Services          []*Service         `json:"serviceDependencies,omitempty" yaml:"serviceDependencies,omitempty"`
		SharedPaths       []string           `json:"sharedPaths,omitempty"         yaml:"sharedPaths,omitempty"`
		Timeout string `json:"timeout,omitempty" yaml:"timeout,omitempty"` 
	}

	// StageDeployment defines a deployment stage.
	StageDeployment struct {
		DeploymentType    string               `json:"deploymentType,omitempty"    yaml:"deploymentType,omitempty"`
		Service           *DeploymentService   `json:"service,omitempty"           yaml:"service,omitempty"`
		ServiceConfig     *DeploymentServiceConfig `json:"serviceConfig,omitempty"     yaml:"serviceConfig,omitempty"`
		Services          *DeploymentServices  `json:"services,omitempty"          yaml:"services,omitempty"`
		Execution         *DeploymentExecution `json:"execution,omitempty"       yaml:"execution,omitempty"`
		EnvironmentGroup  *EnvironmentGroup    `json:"environmentGroup,omitempty"  yaml:"environmentGroup,omitempty"`
		Environment       *Environment         `json:"environment,omitempty"       yaml:"environment,omitempty"`
		Environments      *Environments        `json:"environments,omitempty" yaml:"environments,omitempty"`
		Infrastructure    *DeploymentInfrastructure      `json:"infrastructure,omitempty" yaml:"infrastructure,omitempty"`
		Tags              map[string]string    `json:"tags,omitempty"              yaml:"tags,omitempty"`
		Timeout string `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	}	

	DeploymentInfrastructure struct {
		EnvironmentRef string `json:"environmentRef,omitempty" yaml:"environmentRef,omitempty"`
		InfrastructureDefinition InfrastructureDefinition `json:"infrastructureDefinition,omitempty" yaml:"infrastructureDefinition,omitempty"`	
		AllowSimultaneousDeployments bool `json:"allowSimultaneousDeployments,omitempty" yaml:"allowSimultaneousDeployments,omitempty"`
	}

	// DeploymentExecution defines the deployment execution
	DeploymentExecution struct {
		Steps         []*Steps `json:"steps,omitempty"         yaml:"steps,omitempty"`
		RollbackSteps []*Steps `json:"rollbackSteps,omitempty" yaml:"rollbackSteps,omitempty"`
	}

	// StageFeatureFlag defines a feature flag stage.
	StageFeatureFlag struct {
		Execution *Execution `json:"execution,omitempty" yaml:"execution,omitempty"`
	}

	StageWhen struct {
		PipelineStatus string `json:"pipelineStatus,omitempty" yaml:"pipelineStatus,omitempty"`
		Condition      string `json:"condition,omitempty" yaml:"condition,omitempty"`
	}

	Strategy struct {
		Matrix      map[string]interface{} `json:"matrix,omitempty" yaml:"matrix,omitempty"`
		Parallelism *flexible.Field[int64]          `json:"parallelism,omitempty" yaml:"parallelism,omitempty"`
		Repeat      *Repeat                `json:"repeat,omitempty" yaml:"repeat,omitempty"`
	}

	Exclusion map[string]string

	// Parallelism struct {
	// 	Number         int `yaml:"parallelism"`
	// 	MaxConcurrency flexible.Field[int] `yaml:"maxConcurrency"`
	// }

	Repeat struct {
		Times          *flexible.Field[int64]      `json:"times,omitempty" yaml:"times,omitempty"`
		Items          *flexible.Field[[]string] `json:"items,omitempty" yaml:"items,omitempty"`
		MaxConcurrency *flexible.Field[int64]      `json:"maxConcurrency,omitempty" yaml:"maxConcurrency,omitempty"`
		Start          *flexible.Field[int64]      `json:"start,omitempty" yaml:"start,omitempty"`
		End            *flexible.Field[int64]      `json:"end,omitempty" yaml:"end,omitempty"`
		Unit           string     `json:"unit,omitempty" yaml:"unit,omitempty"`
		NodeName       string     `json:"nodeName,omitempty" yaml:"nodeName,omitempty"`
		PartitionSize  *flexible.Field[int64]      `json:"partitionSize,omitempty" yaml:"partitionSize,omitempty"`
	}
)

// UnmarshalJSON implement the json.Unmarshaler interface.
func (s *Stage) UnmarshalJSON(data []byte) error {
	type S Stage
	type T struct {
		*S
		Spec json.RawMessage `json:"spec"`
	}

	obj := &T{S: (*S)(s)}
	if err := json.Unmarshal(data, obj); err != nil {
		return err
	}

	switch s.Type {
	case StageTypeCI:
		s.Spec = new(StageCI)
	case StageTypeFeatureFlag:
		s.Spec = new(StageFeatureFlag)
	case StageTypeDeployment:
		s.Spec = new(StageDeployment)
	case StageTypeCustom:
		s.Spec = new(StageCustom)
	case StageTypeApproval:
		s.Spec = new(StageApproval)
	default:
		return fmt.Errorf("unknown stage type %s", s.Type)
	}
	return json.Unmarshal(obj.Spec, s.Spec)
}
