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
)

type (
	Step struct { // TODO missing failure strategies
		ID          string            `json:"identifier,omitempty"        yaml:"identifier,omitempty"`
		Description string            `json:"description,omitempty"       yaml:"description,omitempty"`
		Name        string            `json:"name,omitempty"              yaml:"name,omitempty"`
		Skip        string            `json:"skipCondition,omitempty"     yaml:"skipCondition,omitempty"`
		Spec        interface{}       `json:"spec,omitempty"              yaml:"spec,omitempty"`
		Timeout     Duration          `json:"timeout,omitempty"           yaml:"timeout,omitempty"`
		Type        string            `json:"type,omitempty"              yaml:"type,omitempty"`
		When        *StepWhen         `json:"when,omitempty"              yaml:"when,omitempty"`
		Env         map[string]string `json:"envVariables,omitempty"      yaml:"envVariables,omitempty"`
		Strategy    *Strategy         `json:"strategy,omitempty"     yaml:"strategy,omitempty"`
	}

	//
	// Step specifications
	//

	StepArtifactoryUpload struct {
		// TODO
	}

	StepBarrier struct {
		// TODO
	}

	StepBuildAndPushECR struct {
		// TODO
	}

	StepBuildAndPushGCR struct {
		// TODO
	}

	StepFlagConfiguration struct {
		// TODO
	}

	StepGCSUpload struct {
		// TODO
	}

	StepHarnessApproval struct {
		ApprovalMessage                 string           `json:"approvalMessage,omitempty"                 yaml:"approvalMessage,omitempty"`
		Approvers                       *Approvers       `json:"approvers,omitempty"                       yaml:"approvers,omitempty"`
		ApproverInputs                  []*ApproverInput `json:"approverInputs,omitempty"                  yaml:"approverInputs,omitempty"`
		IncludePipelineExecutionHistory string           `json:"includePipelineExecutionHistory,omitempty" yaml:"includePipelineExecutionHistory,omitempty"`
	}

	StepRestoreCacheGCS struct {
		// TODO
	}

	StepRestoreCacheS3 struct {
		// TODO
	}

	StepRunTests struct {
		// TODO
	}

	StepSaveCacheGCS struct {
		// TODO
	}

	StepSaveCacheS3 struct {
		// TODO
	}

	StepDocker struct {
		BuildsArgs      map[string]string `json:"buildArgs,omitempty"       yaml:"buildArgs,omitempty"`
		ConnectorRef    string            `json:"connectorRef,omitempty"    yaml:"connectorRef,omitempty"`
		Context         string            `json:"context,omitempty"         yaml:"context,omitempty"`
		Dockerfile      string            `json:"dockerfile,omitempty"      yaml:"dockerfile,omitempty"`
		Labels          map[string]string `json:"labels,omitempty"          yaml:"labels,omitempty"`
		Optimize        bool              `json:"optimize,omitempty"        yaml:"optimize,omitempty"`
		Privileged      bool              `json:"privileged,omitempty"      yaml:"privileged,omitempty"`
		RemoteCacheRepo string            `json:"remoteCacheRepo,omitempty" yaml:"remoteCacheRepo,omitempty"`
		Repo            string            `json:"repo,omitempty"            yaml:"repo,omitempty"`
		Reports         []*Report         `json:"reports,omitempty"         yaml:"reports,omitempty"`
		Resources       *Resources        `json:"resources,omitempty"       yaml:"resources,omitempty"`
		RunAsUser       string            `json:"runAsUser,omitempty"       yaml:"runAsUser,omitempty"`
		Tags            map[string]string `json:"tags,omitempty"            yaml:"tags,omitempty"`
		Target          string            `json:"target,omitempty"          yaml:"target,omitempty"`
	}

	StepHTTP struct {
		URL             string      `json:"url,omitempty"             yaml:"url,omitempty"`
		Method          string      `json:"method,omitempty"          yaml:"method,omitempty"`
		Headers         []*Variable `json:"headers,omitempty"         yaml:"headers,omitempty"`
		OutputVariables []*Variable `json:"outputVariables,omitempty" yaml:"outputVariables,omitempty"`
		RequestBody     string      `json:"requestBody,omitempty"     yaml:"requestBody,omitempty"`
		Assertion       string      `json:"assertion,omitempty"       yaml:"assertion,omitempty"`

		// NOTE the below fields are not part of the
		// official schema, however, they are useful for
		// executing docker pipelines.

		ConnRef         string `json:"connectorRef,omitempty"    yaml:"connectorRef,omitempty"`
		Image           string `json:"image,omitempty"           yaml:"image,omitempty"`
		ImagePullPolicy string `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty"`
		RunAsUser       string `json:"runAsUser,omitempty"       yaml:"runAsUser,omitempty"`
	}

	StepPlugin struct {
		Env             map[string]string      `json:"envVariables,omitempty"    yaml:"envVariables,omitempty"`
		ConnRef         string                 `json:"connectorRef,omitempty"    yaml:"connectorRef,omitempty"`
		Image           string                 `json:"image,omitempty"           yaml:"image,omitempty"`
		ImagePullPolicy string                 `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty"`
		Privileged      bool                   `json:"privileged,omitempty"      yaml:"privileged,omitempty"`
		Reports         []*Report              `json:"reports,omitempty"         yaml:"reports,omitempty"`
		Resources       *Resources             `json:"resources,omitempty"       yaml:"resources,omitempty"`
		RunAsUser       string                 `json:"runAsUser,omitempty"       yaml:"runAsUser,omitempty"`
		Settings        map[string]interface{} `json:"settings,omitempty"        yaml:"settings,omitempty"`
	}

	StepAction struct {
		Uses string                 `json:"uses,omitempty"            yaml:"uses,omitempty"`
		With map[string]interface{} `json:"with,omitempty"            yaml:"with,omitempty"`
		Envs map[string]string      `json:"env,omitempty"             yaml:"env,omitempty"`
	}

	StepBitrise struct {
		Uses string                 `json:"uses,omitempty"            yaml:"uses,omitempty"`
		With map[string]interface{} `json:"with,omitempty"            yaml:"with,omitempty"`
		Envs map[string]string      `json:"env,omitempty"             yaml:"env,omitempty"`
	}

	StepRun struct {
		Env             map[string]string `json:"envVariables,omitempty"    yaml:"envVariables,omitempty"`
		Command         string            `json:"command,omitempty"         yaml:"command,omitempty"`
		ConnRef         string            `json:"connectorRef,omitempty"    yaml:"connectorRef,omitempty"`
		Image           string            `json:"image,omitempty"           yaml:"image,omitempty"`
		ImagePullPolicy string            `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty"`
		Outputs         []string          `json:"outputVariables,omitempty" yaml:"outputVariables,omitempty"`
		Privileged      bool              `json:"privileged,omitempty"      yaml:"privileged,omitempty"`
		Resources       *Resources        `json:"resources,omitempty"       yaml:"resources,omitempty"`
		RunAsUser       string            `json:"runAsUser,omitempty"       yaml:"runAsUser,omitempty"`
		Reports         *Report           `json:"reports,omitempty"         yaml:"reports,omitempty"`
		Shell           string            `json:"shell,omitempty"           yaml:"shell,omitempty"`
	}

	StepBackground struct {
		Command         string            `json:"command,omitempty"         yaml:"command,omitempty"`
		ConnRef         string            `json:"connectorRef,omitempty"    yaml:"connectorRef,omitempty"`
		Entrypoint      []string          `json:"entrypoint,omitempty"      yaml:"entrypoint,omitempty"`
		Env             map[string]string `json:"envVariables,omitempty"    yaml:"envVariables,omitempty"`
		Image           string            `json:"image,omitempty"           yaml:"image,omitempty"`
		ImagePullPolicy string            `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty"`
		PortBindings    map[string]string `json:"portBindings,omitempty"    yaml:"portBindings,omitempty"`
		Privileged      bool              `json:"privileged,omitempty"      yaml:"privileged,omitempty"`
		Resources       *Resources        `json:"resources,omitempty"       yaml:"resources,omitempty"`
		RunAsUser       string            `json:"runAsUser,omitempty"       yaml:"runAsUser,omitempty"`
	}

	StepScript struct {
		DelegateSelectors    string           `json:"delegateSelectors,omitempty"    yaml:"delegateSelectors,omitempty"`
		EnvironmentVariables []*Variable      `json:"environmentVariables,omitempty" yaml:"environmentVariables,omitempty"`
		ExecutionTarget      *ExecutionTarget `json:"executionTarget,omitempty"      yaml:"executionTarget,omitempty"`
		Metadata             string           `json:"metadata,omitempty"             yaml:"metadata,omitempty"`
		OnDelegate           bool             `json:"onDelegate,omitempty"           yaml:"onDelegate,omitempty"`
		OutputVariables      []*Variable      `json:"outputVariables,omitempty"      yaml:"outputVariables,omitempty"`
		Shell                string           `json:"shell,omitempty"                yaml:"shell,omitempty"` // Bash|Powershell
		Source               *Source          `json:"source,omitempty"               yaml:"source,omitempty"`

		// NOTE the below fields are not part of the
		// official schema, however, they are useful for
		// executing docker pipelines.

		ConnRef         string `json:"connectorRef,omitempty"    yaml:"connectorRef,omitempty"`
		Image           string `json:"image,omitempty"           yaml:"image,omitempty"`
		ImagePullPolicy string `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty"`
		RunAsUser       string `json:"runAsUser,omitempty"       yaml:"runAsUser,omitempty"`
	}

	StepS3Upload struct {
		Bucket       string     `json:"bucket,omitempty"       yaml:"bucket,omitempty"`
		ConnectorRef string     `json:"connectorRef,omitempty" yaml:"connectorRef,omitempty"`
		Endpoint     string     `json:"endpoint,omitempty"     yaml:"endpoint,omitempty"`
		Region       string     `json:"region,omitempty"       yaml:"region,omitempty"`
		Resources    *Resources `json:"resources,omitempty"    yaml:"resources,omitempty"`
		RunAsUser    string     `json:"runAsUser,omitempty"    yaml:"runAsUser,omitempty"`
		SourcePath   string     `json:"sourcePath,omitempty"   yaml:"sourcePath,omitempty"`
		Target       string     `json:"target,omitempty"       yaml:"target,omitempty"`
	}

	//
	// Supporting structures
	//

	Approvers struct {
		MinimumCount             int      `json:"minimumCount,omitempty" yaml:"minimumCount,omitempty"`
		DisallowPipelineExecutor bool     `json:"disallowPipelineExecutor,omitempty" yaml:"disallowPipelineExecutor,omitempty"`
		UserGroups               []string `json:"userGroups,omitempty" yaml:"userGroups,omitempty"`
	}

	ApproverInput struct {
		Name         string `json:"name,omitempty" yaml:"name,omitempty"`
		DefaultValue string `json:"defaultValue,omitempty" yaml:"defaultValue,omitempty"`
	}

	ExecutionTarget struct {
		ConnectorRef     string `json:"connectorRef,omitempty"     yaml:"connectorRef,omitempty"`
		Script           string `json:"script,omitempty"           yaml:"script,omitempty"`
		WorkingDirectory string `json:"workingDirectory,omitempty" yaml:"workingDirectory,omitempty"`
	}

	Source struct {
		Type string     `json:"type,omitempty" yaml:"type,omitempty"`
		Spec SourceSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
	}

	SourceSpec struct {
		Script string `json:"script,omitempty" yaml:"script,omitempty"`
	}

	StepWhen struct {
		StageStatus string `json:"stageStatus,omitempty" yaml:"stageStatus,omitempty"`
		Condition   string `json:"condition,omitempty" yaml:"condition,omitempty"`
	}
)

// UnmarshalJSON implement the json.Unmarshaler interface.
func (s *Step) UnmarshalJSON(data []byte) error {
	type S Step
	type T struct {
		*S
		Spec json.RawMessage `json:"spec"`
	}

	obj := &T{S: (*S)(s)}
	if err := json.Unmarshal(data, obj); err != nil {
		return err
	}

	switch s.Type {
	case StepTypeRun:
		s.Spec = new(StepRun)
	case StepTypePlugin:
		s.Spec = new(StepPlugin)
	case StepTypeBuildAndPushDockerRegistry:
		s.Spec = new(StepDocker)
	case StepTypeS3Upload:
		s.Spec = new(StepS3Upload)
	case StepTypeShellScript:
		s.Spec = new(StepScript)
	case StepTypeHTTP:
		s.Spec = new(StepHTTP)
	default:
		return fmt.Errorf("unknown step type %s", s.Type)
	}

	return json.Unmarshal(obj.Spec, s.Spec)
}
