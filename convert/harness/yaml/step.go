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
	Step struct { // TODO missing failure strategies
		ID                string                              `json:"identifier,omitempty"        yaml:"identifier,omitempty"`
		Description       string                              `json:"description,omitempty"       yaml:"description,omitempty"`
		Name              string                              `json:"name,omitempty"              yaml:"name,omitempty"`
		Skip              string                              `json:"skipCondition,omitempty"     yaml:"skipCondition,omitempty"`
		Spec              interface{}                         `json:"spec,omitempty"              yaml:"spec,omitempty"`
		Timeout           string                              `json:"timeout"                     yaml:"timeout"`
		Type              string                              `json:"type,omitempty"              yaml:"type,omitempty"`
		When              *flexible.Field[StepWhen]           `json:"when,omitempty"              yaml:"when,omitempty"`
		Env               map[string]string                   `json:"envVariables,omitempty"      yaml:"envVariables,omitempty"`
		Strategy          *Strategy                           `json:"strategy,omitempty"     yaml:"strategy,omitempty"`
		FailureStrategies *flexible.Field[[]*FailureStrategy] `json:"failureStrategies,omitempty" yaml:"failureStrategies,omitempty"`
	}

	StepGroup struct { // TODO missing failure strategies
		ID                string                              `json:"identifier,omitempty"        yaml:"identifier,omitempty"`
		Description       string                              `json:"description,omitempty"       yaml:"description,omitempty"`
		DelegateSelectors *flexible.Field[[]string]           `json:"delegateSelectors,omitempty" yaml:"delegateSelectors,omitempty"`
		Name              string                              `json:"name,omitempty"              yaml:"name,omitempty"`
		Skip              string                              `json:"skipCondition,omitempty"     yaml:"skipCondition,omitempty"`
		Steps             []*Steps                            `json:"steps,omitempty"              yaml:"steps,omitempty"`
		Timeout           string                              `json:"timeout"                     yaml:"timeout"`
		When              *flexible.Field[StepWhen]           `json:"when,omitempty"              yaml:"when,omitempty"`
		Env               map[string]string                   `json:"envVariables,omitempty"      yaml:"envVariables,omitempty"`
		Strategy          *Strategy                           `json:"strategy,omitempty"     yaml:"strategy,omitempty"`
		Variables         []*Variable                         `json:"variables,omitempty"       yaml:"variables,omitempty"`
		FailureStrategies *flexible.Field[[]*FailureStrategy] `json:"failureStrategies,omitempty" yaml:"failureStrategies,omitempty"`
	}

	//
	// Step specifications
	//

	CommonStepSpec struct {
		IncludeInfraSelectors bool                      `json:"includeInfraSelectors,omitempty" yaml:"includeInfraSelectors,omitempty"`
		DelegateSelectors     *flexible.Field[[]string] `json:"delegateSelectors,omitempty" yaml:"delegateSelectors,omitempty"`
	}

	StepArtifactoryUpload struct {
		CommonStepSpec
		ConnRef    string `json:"connectorRef,omitempty"    yaml:"connectorRef,omitempty"`
		Target     string `json:"target,omitempty"       yaml:"target,omitempty"`
		SourcePath string `json:"sourcePath,omitempty"   yaml:"sourcePath,omitempty"`
		RunAsUser  string `json:"runAsUser,omitempty"    yaml:"runAsUser,omitempty"`
	}

	StepBarrier struct {
		CommonStepSpec
		BarrierRef string `json:"barrierRef,omitempty" yaml:"barrierRef,omitempty"`
	}

	StepQueue struct {
		CommonStepSpec
		Key   string `json:"key,omitempty" yaml:"key,omitempty"`
		Scope string `json:"scope,omitempty" yaml:"scope,omitempty"`
	}

	StepBuildAndPushECR struct {
		CommonStepSpec
		ConnectorRef           string            `json:"connectorRef,omitempty"           yaml:"connectorRef,omitempty"`
		Region                 string            `json:"region,omitempty"                 yaml:"region,omitempty"`
		Account                string            `json:"account,omitempty"                yaml:"account,omitempty"`
		ImageName              string            `json:"imageName,omitempty"              yaml:"imageName,omitempty"`
		Tags                   *flexible.Field[[]string]          `json:"tags,omitempty"                   yaml:"tags,omitempty"`
		Caching                *flexible.Field[bool]              `json:"caching,omitempty"                yaml:"caching,omitempty"`
		BaseImageConnectorRefs interface{}       `json:"baseImageConnectorRefs,omitempty" yaml:"baseImageConnectorRefs,omitempty"`
		Dockerfile             string            `json:"dockerfile,omitempty"             yaml:"dockerfile,omitempty"`
		Context                string            `json:"context,omitempty"                yaml:"context,omitempty"`
		Labels                 *flexible.Field[map[string]string] `json:"labels,omitempty"                 yaml:"labels,omitempty"`
		BuildArgs              *flexible.Field[map[string]string]`json:"buildArgs,omitempty"              yaml:"buildArgs,omitempty"`
		Target                 string            `json:"target,omitempty"                 yaml:"target,omitempty"`
		Env                    *flexible.Field[map[string]string] `json:"envVariables,omitempty"           yaml:"envVariables,omitempty"`
		RunAsUser              string            `json:"runAsUser,omitempty"              yaml:"runAsUser,omitempty"`
	}

	// Deprecated
	StepBuildAndPushGCR struct {
		// TODO
	}

	StepBuildAndPushGAR struct {
		CommonStepSpec
		ConnectorRef           string            `json:"connectorRef,omitempty"           yaml:"connectorRef,omitempty"`
		Host                   string            `json:"host,omitempty"                   yaml:"host,omitempty"`
		ProjectID              string            `json:"projectID,omitempty"              yaml:"projectID,omitempty"`
		ImageName              string            `json:"imageName,omitempty"              yaml:"imageName,omitempty"`
		Tags                   *flexible.Field[[]string]          `json:"tags,omitempty"                   yaml:"tags,omitempty"`
		Caching                *flexible.Field[bool]              `json:"caching,omitempty"                yaml:"caching,omitempty"`
		BaseImageConnectorRefs interface{}       `json:"baseImageConnectorRefs,omitempty" yaml:"baseImageConnectorRefs,omitempty"`
		Dockerfile             string            `json:"dockerfile,omitempty"             yaml:"dockerfile,omitempty"`
		Context                string            `json:"context,omitempty"                yaml:"context,omitempty"`
		Labels                 *flexible.Field[map[string]string] `json:"labels,omitempty"                 yaml:"labels,omitempty"`
		BuildArgs              *flexible.Field[map[string]string] `json:"buildArgs,omitempty"              yaml:"buildArgs,omitempty"`
		Target                 string            `json:"target,omitempty"                 yaml:"target,omitempty"`
		Env                    *flexible.Field[map[string]string] `json:"envVariables,omitempty"           yaml:"envVariables,omitempty"`
		RunAsUser              string            `json:"runAsUser,omitempty"              yaml:"runAsUser,omitempty"`
	}
   
	StepBuildAndPushDockerRegistry struct {
		CommonStepSpec
		BuildArgs             *flexible.Field[map[string]string] `json:"buildArgs,omitempty"               yaml:"buildArgs,omitempty"`
		ConnectorRef           string            `json:"connectorRef,omitempty"            yaml:"connectorRef,omitempty"`
		Context                string            `json:"context,omitempty"                 yaml:"context,omitempty"`
		Dockerfile             string            `json:"dockerfile,omitempty"              yaml:"dockerfile,omitempty"`
		Labels                 *flexible.Field[map[string]string] `json:"labels,omitempty"                  yaml:"labels,omitempty"`
		Optimize               *flexible.Field[bool]              `json:"optimize,omitempty"                yaml:"optimize,omitempty"`
		Privileged             *flexible.Field[bool]              `json:"privileged,omitempty"              yaml:"privileged,omitempty"`
		RemoteCacheRepo        string            `json:"remoteCacheRepo,omitempty"         yaml:"remoteCacheRepo,omitempty"`
		Repo                   string            `json:"repo,omitempty"                    yaml:"repo,omitempty"`
		Reports                []*Report         `json:"reports,omitempty"                 yaml:"reports,omitempty"`
		Resources              *Resources        `json:"resources,omitempty"               yaml:"resources,omitempty"`
		RunAsUser              string            `json:"runAsUser,omitempty"               yaml:"runAsUser,omitempty"`
		Tags                   *flexible.Field[[]string]          `json:"tags,omitempty"                    yaml:"tags,omitempty"`
		Target                 string            `json:"target,omitempty"                  yaml:"target,omitempty"`
		Caching                *flexible.Field[bool]              `json:"caching,omitempty"                 yaml:"caching,omitempty"`
		Env                    *flexible.Field[map[string]string] `json:"envVariables,omitempty"            yaml:"envVariables,omitempty"`
		BaseImageConnectorRefs interface{}       `json:"baseImageConnectorRefs,omitempty"  yaml:"baseImageConnectorRefs,omitempty"`
	}

	StepBuildAndPushACR struct {
	}

	StepFlagConfiguration struct {
		// TODO
	}

	StepGCSUpload struct {
		CommonStepSpec
		ConnectorRef string `json:"connectorRef,omitempty" yaml:"connectorRef,omitempty"`
		Bucket       string `json:"bucket,omitempty"       yaml:"bucket,omitempty"`
		SourcePath   string `json:"sourcePath,omitempty"   yaml:"sourcePath,omitempty"`
		Target       string `json:"target,omitempty"       yaml:"target,omitempty"`
		RunAsUser    string `json:"runAsUser,omitempty"    yaml:"runAsUser,omitempty"`
	}

	StepHarnessApproval struct {
		CommonStepSpec
		ApprovalMessage                 string           `json:"approvalMessage,omitempty"                 yaml:"approvalMessage,omitempty"`
		Approvers                       *Approvers       `json:"approvers,omitempty"                       yaml:"approvers,omitempty"`
		ApproverInputs                  []*ApproverInput `json:"approverInputs,omitempty"                  yaml:"approverInputs,omitempty"`
		IncludePipelineExecutionHistory *flexible.Field[bool]             `json:"includePipelineExecutionHistory,omitempty" yaml:"includePipelineExecutionHistory,omitempty"`
		IsAutoRejectEnabled             *flexible.Field[bool]             `json:"isAutoRejectEnabled,omitempty"             yaml:"isAutoRejectEnabled,omitempty"`
		AutoApproval                    *AutoApproval    `json:"autoApproval,omitempty"                    yaml:"autoApproval,omitempty"`
	}

	StepCustomApproval struct {
		CommonStepSpec
		Shell                string      `json:"shell,omitempty"                yaml:"shell,omitempty"`
		RetryInterval        string      `json:"retryInterval,omitempty"        yaml:"retryInterval,omitempty"`
		ScriptTimeout        string      `json:"scriptTimeout,omitempty"        yaml:"scriptTimeout,omitempty"`
		Source               *Source     `json:"source,omitempty"               yaml:"source,omitempty"`
		EnvironmentVariables []*Variable `json:"environmentVariables,omitempty" yaml:"environmentVariables,omitempty"`
		ApprovalCriteria     *Criteria   `json:"approvalCriteria,omitempty"     yaml:"approvalCriteria,omitempty"`
		RejectionCriteria    *Criteria   `json:"rejectionCriteria,omitempty"    yaml:"rejectionCriteria,omitempty"`
		OutputVariables      []*Output   `json:"outputVariables,omitempty"      yaml:"outputVariables,omitempty"`
	}

	Criteria struct {
		Type string        `json:"type,omitempty" yaml:"type,omitempty"` // KeyValues or Jexl
		Spec *CriteriaSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
	}

	CriteriaSpec struct {
		MatchAnyCondition *flexible.Field[bool]         `json:"matchAnyCondition,omitempty" yaml:"matchAnyCondition,omitempty"`
		Conditions        []*Condition `json:"conditions,omitempty"        yaml:"conditions,omitempty"`
		Expression        string       `json:"expression,omitempty"        yaml:"expression,omitempty"` // For Jexl type
	}

	Condition struct {
		Key      string      `json:"key,omitempty"      yaml:"key,omitempty"`
		Operator string      `json:"operator,omitempty" yaml:"operator,omitempty"` // equals, not equals, in, not in
		Value    interface{} `json:"value,omitempty"    yaml:"value,omitempty"`
	}

	StepRestoreCacheGCS struct {
		CommonStepSpec
		ConnectorRef      string `json:"connectorRef,omitempty"      yaml:"connectorRef,omitempty"`
		Bucket            string `json:"bucket,omitempty"            yaml:"bucket,omitempty"`
		Key               string `json:"key,omitempty"               yaml:"key,omitempty"`
		ArchiveFormat     string `json:"archiveFormat,omitempty"     yaml:"archiveFormat,omitempty"`
		FailIfKeyNotFound *flexible.Field[bool]   `json:"failIfKeyNotFound,omitempty" yaml:"failIfKeyNotFound,omitempty"`
		RunAsUser         string `json:"runAsUser,omitempty"         yaml:"runAsUser,omitempty"`
	}

	StepRestoreCacheS3 struct {
		CommonStepSpec
		ConnectorRef      string `json:"connectorRef,omitempty"      yaml:"connectorRef,omitempty"`
		Region            string `json:"region,omitempty"            yaml:"region,omitempty"`
		Bucket            string `json:"bucket,omitempty"            yaml:"bucket,omitempty"`
		Key               string `json:"key,omitempty"               yaml:"key,omitempty"`
		Endpoint          string `json:"endpoint,omitempty"          yaml:"endpoint,omitempty"`
		ArchiveFormat     string `json:"archiveFormat,omitempty"     yaml:"archiveFormat,omitempty"`
		PathStyle         *flexible.Field[bool]   `json:"pathStyle,omitempty"         yaml:"pathStyle,omitempty"`
		FailIfKeyNotFound *flexible.Field[bool]   `json:"failIfKeyNotFound,omitempty" yaml:"failIfKeyNotFound,omitempty"`
		RunAsUser         string `json:"runAsUser,omitempty"         yaml:"runAsUser,omitempty"`
	}

	StepRunTests struct {
		// TODO
	}

	StepTestIntelligence struct {
		CommonStepSpec
		ConnRef          string            `json:"connectorRef,omitempty"     yaml:"connectorRef,omitempty"`
		Image            string            `json:"image,omitempty"            yaml:"image,omitempty"`
		Shell            string            `json:"shell,omitempty"            yaml:"shell,omitempty"`
		Command          string            `json:"command,omitempty"          yaml:"command,omitempty"`
		Privileged       *flexible.Field[bool]              `json:"privileged,omitempty"       yaml:"privileged,omitempty"`
		Reports          *Report           `json:"reports,omitempty"          yaml:"reports,omitempty"`
		Env              *flexible.Field[map[string]string] `json:"envVariables,omitempty"     yaml:"envVariables,omitempty"`
		Outputs          []*Output         `json:"outputVariables,omitempty"  yaml:"outputVariables,omitempty"`
		ImagePullPolicy  string            `json:"imagePullPolicy,omitempty"  yaml:"imagePullPolicy,omitempty"`
		IntelligenceMode *flexible.Field[bool]              `json:"intelligenceMode,omitempty" yaml:"intelligenceMode,omitempty"`
		Globs            []string          `json:"globs,omitempty"            yaml:"globs,omitempty"`
		RunAsUser        string            `json:"runAsUser,omitempty"        yaml:"runAsUser,omitempty"`
		Resources        *Resources        `json:"resources,omitempty"        yaml:"resources,omitempty"`
	}

	StepSaveCacheGCS struct {
		CommonStepSpec
		ConnectorRef  string   `json:"connectorRef,omitempty"  yaml:"connectorRef,omitempty"`
		Bucket        string   `json:"bucket,omitempty"        yaml:"bucket,omitempty"`
		Key           string   `json:"key,omitempty"           yaml:"key,omitempty"`
		SourcePaths   []string `json:"sourcePaths,omitempty"   yaml:"sourcePaths,omitempty"`
		ArchiveFormat string   `json:"archiveFormat,omitempty" yaml:"archiveFormat,omitempty"`
		Override      *flexible.Field[bool]     `json:"override,omitempty"      yaml:"override,omitempty"`
		RunAsUser     string   `json:"runAsUser,omitempty"     yaml:"runAsUser,omitempty"`
	}

	StepSaveCacheS3 struct {
		CommonStepSpec
		ConnectorRef  string   `json:"connectorRef,omitempty"  yaml:"connectorRef,omitempty"`
		Region        string   `json:"region,omitempty"        yaml:"region,omitempty"`
		Bucket        string   `json:"bucket,omitempty"        yaml:"bucket,omitempty"`
		Key           string   `json:"key,omitempty"           yaml:"key,omitempty"`
		SourcePaths   []string `json:"sourcePaths,omitempty"   yaml:"sourcePaths,omitempty"`
		Endpoint      string   `json:"endpoint,omitempty"      yaml:"endpoint,omitempty"`
		ArchiveFormat string   `json:"archiveFormat,omitempty" yaml:"archiveFormat,omitempty"`
		Override      *flexible.Field[bool]     `json:"override,omitempty"      yaml:"override,omitempty"`
		PathStyle     *flexible.Field[bool]     `json:"pathStyle,omitempty"     yaml:"pathStyle,omitempty"`
		RunAsUser     string   `json:"runAsUser,omitempty"     yaml:"runAsUser,omitempty"`
	}

	StepTrivy struct {
		CommonStepSpec
		Mode       string       `json:"mode,omitempty"    yaml:"mode,omitempty"`
		Config     string       `json:"config,omitempty"         yaml:"config,omitempty"`
		Target     *STOTarget   `json:"target,omitempty"         yaml:"target,omitempty"`
		Advanced   *STOAdvanced `json:"advanced,omitempty"         yaml:"advanced,omitempty"`
		Privileged bool         `json:"privileged,omitempty"      yaml:"privileged,omitempty"`
		Image      *STOImage    `json:"image,omitempty"      yaml:"image,omitempty"`
	}

	StepDocker struct {
		CommonStepSpec
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
		Tags            []string          `json:"tags,omitempty"            yaml:"tags,omitempty"`
		Target          string            `json:"target,omitempty"          yaml:"target,omitempty"`
		Caching         bool              `json:"caching,omitempty"         yaml:"caching,omitempty"`
	}

	StepHTTP struct {
		CommonStepSpec
		URL             string          `json:"url,omitempty"             yaml:"url,omitempty"`
		Method          string          `json:"method,omitempty"          yaml:"method,omitempty"`
		Headers         []*KeyValuePair `json:"headers,omitempty"         yaml:"headers,omitempty"`
		InputVariables  []*Variable     `json:"inputVariables,omitempty"  yaml:"inputVariables,omitempty"`
		OutputVariables []*Variable     `json:"outputVariables,omitempty" yaml:"outputVariables,omitempty"`
		RequestBody     string          `json:"requestBody,omitempty"     yaml:"requestBody,omitempty"`
		Assertion       string          `json:"assertion,omitempty"       yaml:"assertion,omitempty"`
		Certificate     string          `json:"certificate,omitempty"     yaml:"certificate,omitempty"`
		CertificateKey  string          `json:"certificateKey,omitempty"  yaml:"certificateKey,omitempty"`

		// NOTE the below fields are not part of the
		// official schema, however, they are useful for
		// executing docker pipelines.

		ConnRef         string `json:"connectorRef,omitempty"    yaml:"connectorRef,omitempty"`
		Image           string `json:"image,omitempty"           yaml:"image,omitempty"`
		ImagePullPolicy string `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty"`
		RunAsUser       string `json:"runAsUser,omitempty"       yaml:"runAsUser,omitempty"`
	}

	KeyValuePair struct {
		Key   string `json:"key,omitempty"   yaml:"key,omitempty"`
		Value string `json:"value,omitempty" yaml:"value,omitempty"`
	}

	StepEmail struct {
		CommonStepSpec
		To             string      `json:"to,omitempty"             yaml:"to,omitempty"`
		ToUserGroups   []string    `json:"toUserGroups,omitempty"   yaml:"toUserGroups,omitempty"`
		Cc             string      `json:"cc,omitempty"             yaml:"cc,omitempty"`
		CcUserGroups   []string    `json:"ccUserGroups,omitempty"   yaml:"ccUserGroups,omitempty"`
		Subject        string      `json:"subject,omitempty"        yaml:"subject,omitempty"`
		Body           string      `json:"body,omitempty"           yaml:"body,omitempty"`
		InputVariables []*Variable `json:"inputVariables,omitempty" yaml:"inputVariables,omitempty"`
	}

	StepPlugin struct {
		CommonStepSpec
		Env             *flexible.Field[map[string]interface{}]      `json:"envVariables,omitempty"    yaml:"envVariables,omitempty"`
		ConnRef         string                 `json:"connectorRef,omitempty"    yaml:"connectorRef,omitempty"`
		Image           string                 `json:"image,omitempty"           yaml:"image,omitempty"`
		ImagePullPolicy string                 `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty"`
		Privileged      *flexible.Field[bool]                    `json:"privileged,omitempty"      yaml:"privileged,omitempty"`
		Reports         *Report                `json:"reports,omitempty"         yaml:"reports,omitempty"`
		Resources       *Resources             `json:"resources,omitempty"       yaml:"resources,omitempty"`
		RunAsUser       string                 `json:"runAsUser,omitempty"       yaml:"runAsUser,omitempty"`
		Settings        map[string]interface{} `json:"settings,omitempty"        yaml:"settings,omitempty"`
		Entrypoint      *flexible.Field[[]string]               `json:"entrypoint,omitempty"      yaml:"entrypoint,omitempty"`
	}

	StepGitClone struct {
		CommonStepSpec
		Repository     string     `json:"repoName,omitempty"    yaml:"repoNama,omitempty"`
		ConnRef        string     `json:"connectorRef,omitempty"    yaml:"connectorRef,omitempty"`
		BuildType      *flexible.Field[Build]     `json:"build,omitempty"    yaml:"build,omitempty"`
		CloneDirectory string     `json:"cloneDirectory,omitempty"    yaml:"cloneDirectory,omitempty"`
		Privileged     *flexible.Field[bool]       `json:"privileged,omitempty"      yaml:"privileged,omitempty"`
		Depth          *flexible.Field[int]     `json:"depth,omitempty"    yaml:"cloneDirectory,omitempty"`
		Resources      *Resources `json:"resources,omitempty"       yaml:"resources,omitempty"`
		SSLVerify      string     `json:"sslVerify,omitempty"    yaml:"cloneDirectory,omitempty"`
		RunAsUser      string     `json:"runAsUser,omitempty"       yaml:"runAsUser,omitempty"`
		Timeout        string     `json:"timeout,omitempty"    yaml:"cloneDirectory,omitempty"`
	}

	StepAction struct {
		CommonStepSpec
		Uses string                 `json:"uses,omitempty"            yaml:"uses,omitempty"`
		With map[string]interface{} `json:"with,omitempty"            yaml:"with,omitempty"`
		Envs map[string]string      `json:"env,omitempty"             yaml:"env,omitempty"`
	}

	StepBitrise struct {
		CommonStepSpec
		Uses string                 `json:"uses,omitempty"            yaml:"uses,omitempty"`
		With map[string]interface{} `json:"with,omitempty"            yaml:"with,omitempty"`
		Envs map[string]string      `json:"env,omitempty"             yaml:"env,omitempty"`
	}

	StepRun struct {
		CommonStepSpec
		Env             *flexible.Field[map[string]interface{}] `json:"envVariables,omitempty"    yaml:"envVariables,omitempty"    v1path:".env"`
		Command         string            `json:"command,omitempty"         yaml:"command,omitempty"         v1path:".script"`
		ConnRef         string            `json:"connectorRef,omitempty"    yaml:"connectorRef,omitempty"    v1path:".container.connector"`
		Image           string            `json:"image,omitempty"           yaml:"image,omitempty"           v1path:".container.image"`
		ImagePullPolicy string            `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty" v1path:".container.pull"`
		Outputs         []*Output         `json:"outputVariables,omitempty" yaml:"outputVariables,omitempty" v1path:".outputs"`
		Privileged      *flexible.Field[bool]               `json:"privileged,omitempty"      yaml:"privileged,omitempty"      v1path:".container.privileged"`
		Resources       *Resources        `json:"resources,omitempty"       yaml:"resources,omitempty"       v1path:"-"`
		RunAsUser       *flexible.Field[int]            `json:"runAsUser,omitempty"       yaml:"runAsUser,omitempty"       v1path:"-"`
		Reports         *Report           `json:"reports,omitempty"         yaml:"reports,omitempty"         v1path:".report"`
		Shell           string            `json:"shell,omitempty"           yaml:"shell,omitempty"           v1path:".shell"`
		Alias           *OutputAlias      `json:"outputAlias,omitempty"     yaml:"outputAlias,omitempty"     v1path:".alias"`
	}

	OutputAlias struct {
		Key string `json:"key,omitempty" yaml:"key,omitempty"`
		Scope string `json:"scope,omitempty" yaml:"scope,omitempty"`
	}

	StepBackground struct {
		CommonStepSpec
		Command         string            `json:"command,omitempty"         yaml:"command,omitempty"`
		ConnRef         string            `json:"connectorRef,omitempty"    yaml:"connectorRef,omitempty"`
		Entrypoint      *flexible.Field[[]string]          `json:"entrypoint,omitempty"      yaml:"entrypoint,omitempty"`
		Env             *flexible.Field[map[string]interface{}] `json:"envVariables,omitempty"    yaml:"envVariables,omitempty"`
		Image           string            `json:"image,omitempty"           yaml:"image,omitempty"`
		ImagePullPolicy string            `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty"`
		PortBindings    map[string]string `json:"portBindings,omitempty"    yaml:"portBindings,omitempty"`
		Privileged      *flexible.Field[bool]               `json:"privileged,omitempty"      yaml:"privileged,omitempty"`
		Resources       *Resources        `json:"resources,omitempty"       yaml:"resources,omitempty"`
		RunAsUser       string            `json:"runAsUser,omitempty"       yaml:"runAsUser,omitempty"`
		Reports         *Report           `json:"reports,omitempty"         yaml:"reports,omitempty"`
		Shell           string            `json:"shell,omitempty"           yaml:"shell,omitempty"`
	}

	StepShellScript struct {
		CommonStepSpec
		EnvironmentVariables []*Variable      `json:"environmentVariables,omitempty" yaml:"environmentVariables,omitempty"`
		ExecutionTarget      *ExecutionTarget `json:"executionTarget,omitempty"      yaml:"executionTarget,omitempty"`
		Metadata             string           `json:"metadata,omitempty"             yaml:"metadata,omitempty"`
		OnDelegate           *flexible.Field[bool]             `json:"onDelegate,omitempty"           yaml:"onDelegate,omitempty"`
		OutputVariables      []*Output        `json:"outputVariables,omitempty"      yaml:"outputVariables,omitempty"`
		Shell                string           `json:"shell,omitempty"                yaml:"shell,omitempty"` // Bash|Powershell
		Source               *Source          `json:"source,omitempty"               yaml:"source,omitempty"`
		Alias                *OutputAlias     `json:"outputAlias,omitempty"          yaml:"outputAlias,omitempty"`
		// NOTE the below fields are not part of the
		// official schema, however, they are useful for
		// executing docker pipelines.

		ConnRef         string `json:"connectorRef,omitempty"    yaml:"connectorRef,omitempty"`
		Image           string `json:"image,omitempty"           yaml:"image,omitempty"`
		ImagePullPolicy string `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty"`
		RunAsUser       string `json:"runAsUser,omitempty"       yaml:"runAsUser,omitempty"`
	}

	StepS3Upload struct {
		CommonStepSpec
		Bucket       string     `json:"bucket,omitempty"       yaml:"bucket,omitempty"`
		ConnectorRef string     `json:"connectorRef,omitempty" yaml:"connectorRef,omitempty"`
		Endpoint     string     `json:"endpoint,omitempty"     yaml:"endpoint,omitempty"`
		Region       string     `json:"region,omitempty"       yaml:"region,omitempty"`
		Resources    *Resources `json:"resources,omitempty"    yaml:"resources,omitempty"`
		RunAsUser    string     `json:"runAsUser,omitempty"    yaml:"runAsUser,omitempty"`
		SourcePath   string     `json:"sourcePath,omitempty"   yaml:"sourcePath,omitempty"`
		Target       string     `json:"target,omitempty"       yaml:"target,omitempty"`
	}

	// Feature: Jira Create
	StepJiraCreate struct {
		CommonStepSpec
		ConnectorRef string      `json:"connectorRef,omitempty" yaml:"connectorRef,omitempty"`
		ProjectKey   string      `json:"projectKey,omitempty" yaml:"projectKey,omitempty"`
		IssueType    string      `json:"issueType,omitempty" yaml:"issueType,omitempty"`
		Fields       []*JiraField `json:"fields,omitempty" yaml:"fields,omitempty"`
	}

	// Feature: Jira Update
	StepJiraUpdate struct {
		CommonStepSpec
		ConnectorRef string          `json:"connectorRef,omitempty" yaml:"connectorRef,omitempty"`
		ProjectKey   string          `json:"projectKey,omitempty" yaml:"projectKey,omitempty"`
		IssueType    string          `json:"issueType,omitempty" yaml:"issueType,omitempty"`
		IssueKey     string          `json:"issueKey,omitempty" yaml:"issueKey,omitempty"`
		Fields       []*JiraField     `json:"fields,omitempty" yaml:"fields,omitempty"`
		TransitionTo *JiraTransition `json:"transitionTo,omitempty" yaml:"transitionTo,omitempty"`
	}

	JiraField struct {
		Name string `json:"name,omitempty" yaml:"name,omitempty"`
		Value string `json:"value,omitempty" yaml:"value,omitempty"`
	}

	StepJiraApproval struct {
		CommonStepSpec
		ConnectorRef      string    `json:"connectorRef,omitempty" yaml:"connectorRef,omitempty"`
		ApprovalCriteria  *Criteria `json:"approvalCriteria,omitempty"     yaml:"approvalCriteria,omitempty"`
		RejectionCriteria *Criteria `json:"rejectionCriteria,omitempty"    yaml:"rejectionCriteria,omitempty"`
		RetryInterval     string    `json:"retryInterval,omitempty" yaml:"retryInterval,omitempty"`
		IssueKey          string    `json:"issueKey,omitempty" yaml:"issueKey,omitempty"`
		ProjectKey        string    `json:"projectKey,omitempty" yaml:"projectKey,omitempty"`
		IssueType         string    `json:"issueType,omitempty" yaml:"issueType,omitempty"`
	}

	StepServiceNowApproval struct {
		CommonStepSpec
		ConnectorRef      string    `json:"connectorRef,omitempty" yaml:"connectorRef,omitempty"`
		ApprovalCriteria  *Criteria `json:"approvalCriteria,omitempty"     yaml:"approvalCriteria,omitempty"`
		RejectionCriteria *Criteria `json:"rejectionCriteria,omitempty"    yaml:"rejectionCriteria,omitempty"`
		RetryInterval     string    `json:"retryInterval,omitempty" yaml:"retryInterval,omitempty"`
		TicketNumber      string    `json:"ticketNumber,omitempty" yaml:"ticketNumber,omitempty"`
		TicketType        string    `json:"ticketType,omitempty" yaml:"ticketType,omitempty"`
		ChangeWindow      *ChangeWindow `json:"changeWindow,omitempty" yaml:"changeWindow,omitempty"`
	}

	ChangeWindow struct {
		StartField string `json:"startField,omitempty" yaml:"startField,omitempty"`
		EndField string `json:"endField,omitempty" yaml:"endField,omitempty"`
	}


	StepServiceNowCreate struct {
		CommonStepSpec
		ConnectorRef string      `json:"connectorRef,omitempty" yaml:"connectorRef,omitempty"`
		TicketType   string      `json:"ticketType,omitempty" yaml:"ticketType,omitempty"`
		Fields       []*ServiceNowField `json:"fields,omitempty" yaml:"fields,omitempty"`
		CreateType   string      `json:"createType,omitempty" yaml:"createType,omitempty"`
	}

	StepServiceNowUpdate struct {
		CommonStepSpec
		ConnectorRef          string      `json:"connectorRef,omitempty" yaml:"connectorRef,omitempty"`
		TicketType            string      `json:"ticketType,omitempty" yaml:"ticketType,omitempty"`
		TicketNumber          string      `json:"ticketNumber,omitempty" yaml:"ticketNumber,omitempty"`
		Fields                []*ServiceNowField `json:"fields,omitempty" yaml:"fields,omitempty"`
		UseServiceNowTemplate bool        `json:"useServiceNowTemplate,omitempty" yaml:"useServiceNowTemplate,omitempty"`
	}

	ServiceNowField struct {
		Name string `json:"name,omitempty" yaml:"name,omitempty"`
		Value string `json:"value,omitempty" yaml:"value,omitempty"`
	}

	JiraTransition struct {
		Status         string `json:"status,omitempty" yaml:"status,omitempty"`
		TransitionName string `json:"transitionName,omitempty" yaml:"transitionName,omitempty"`
	}

	// Wait step
	StepWait struct {
		CommonStepSpec
		Duration string `json:"duration,omitempty" yaml:"duration,omitempty"`
	}

	//
	// Supporting structures
	//

	Approvers struct {
		MinimumCount             *flexible.Field[int]      `json:"minimumCount,omitempty" yaml:"minimumCount,omitempty"`
		DisallowPipelineExecutor *flexible.Field[bool]     `json:"disallowPipelineExecutor,omitempty" yaml:"disallowPipelineExecutor,omitempty"`
		UserGroups               *flexible.Field[[]string] `json:"userGroups,omitempty" yaml:"userGroups,omitempty"`
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
		Condition   *flexible.Field[bool] `json:"condition,omitempty" yaml:"condition,omitempty"`
	}

	AutoApproval struct {
		Action            string             `json:"action,omitempty"            yaml:"action,omitempty"`
		ScheduledDeadline *ScheduledDeadline `json:"scheduledDeadline,omitempty" yaml:"scheduledDeadline,omitempty"`
		Comments          string             `json:"comments,omitempty"          yaml:"comments,omitempty"`
	}

	ScheduledDeadline struct {
		TimeZone string `json:"timeZone,omitempty" yaml:"timeZone,omitempty"`
		Time     string `json:"time,omitempty"     yaml:"time,omitempty"`
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
	case StepTypeCustomApproval:
		s.Spec = new(StepCustomApproval)
	case StepTypeAction:
		s.Spec = new(StepAction)
	case StepTypeBackground:
		s.Spec = new(StepBackground)
	case StepTypeGitClone:
		s.Spec = new(StepGitClone)
	case StepTypeRun:
		s.Spec = new(StepRun)
	case StepTypeBarrier:
		s.Spec = new(StepBarrier)
	case StepTypeQueue:
		s.Spec = new(StepQueue)
	case StepTypePlugin:
		s.Spec = new(StepPlugin)
	case StepTypeBuildAndPushDockerRegistry:
		s.Spec = new(StepBuildAndPushDockerRegistry)
	case StepTypeS3Upload:
		s.Spec = new(StepS3Upload)
	case StepTypeShellScript:
		s.Spec = new(StepShellScript)
	case StepTypeHTTP:
		s.Spec = new(StepHTTP)
	case StepTypeHarnessApproval:
		s.Spec = new(StepHarnessApproval)
	case StepTypeK8sRollingDeploy:
		s.Spec = new(StepK8sRollingDeploy)
	case StepTypeK8sRollingRollback:
		s.Spec = new(StepK8sRollingRollback)
	case StepTypeK8sApply:
		s.Spec = new(StepK8sApply)
	case StepTypeK8sBGSwapServices:
		s.Spec = new(StepK8sBGSwapServices)
	case StepTypeK8sBlueGreenStageScaleDown:
		s.Spec = new(StepK8sBlueGreenStageScaleDown)
	case StepTypeK8sCanaryDelete:
		s.Spec = new(StepK8sCanaryDelete)
	case StepTypeK8sDelete:
		s.Spec = new(StepK8sDelete)
	case StepTypeK8sDiff:
		s.Spec = new(StepK8sDiff)
	case StepTypeK8sRollout:
		s.Spec = new(StepK8sRollout)
	case StepTypeK8sScale:
		s.Spec = new(StepK8sScale)
	case StepTypeK8sDryRun:
		s.Spec = new(StepK8sDryRun)
	case StepTypeK8sTrafficRouting:
		s.Spec = new(StepK8sTrafficRouting)
	case StepTypeK8sCanaryDeploy:
		s.Spec = new(StepK8sCanaryDeploy)
	case StepTypeK8sBlueGreenDeploy:
		s.Spec = new(StepK8sBlueGreenDeploy)
	case StepTypeK8sPatch:
		s.Spec = new(StepK8sPatch)
	case StepTypeJiraCreate:
		s.Spec = new(StepJiraCreate)
	case StepTypeJiraUpdate:
		s.Spec = new(StepJiraUpdate)
	case StepTypeJiraApproval:
		s.Spec = new(StepJiraApproval)
	case StepTypeServiceNowApproval:
		s.Spec = new(StepServiceNowApproval)
	case StepTypeHelmBGDeploy:
		s.Spec = new(StepHelmBGDeploy)
	case StepTypeHelmBlueGreenSwapStep:
		s.Spec = new(StepHelmBlueGreenSwapStep)
	case StepTypeHelmCanaryDeploy:
		s.Spec = new(StepHelmCanaryDeploy)
	case StepTypeHelmCanaryDelete:
		s.Spec = new(StepHelmCanaryDelete)
	case StepTypeHelmDelete:
		s.Spec = new(StepHelmDelete)
	case StepTypeHelmDeploy:
		s.Spec = new(StepHelmDeploy)
	case StepTypeHelmRollback:
		s.Spec = new(StepHelmRollback)
	case StepTypeWait:
		s.Spec = new(StepWait)
	case StepTypeEmail:
		s.Spec = new(StepEmail)
	case StepTypeSaveCacheS3:
		s.Spec = new(StepSaveCacheS3)
	case StepTypeArtifactoryUpload:
		s.Spec = new(StepArtifactoryUpload)
	case StepTypeSaveCacheGCS:
		s.Spec = new(StepSaveCacheGCS)
	case StepTypeRestoreCacheGCS:
		s.Spec = new(StepRestoreCacheGCS)
	case StepTypeRestoreCacheS3:
		s.Spec = new(StepRestoreCacheS3)
	case StepTypeBuildAndPushECR:
		s.Spec = new(StepBuildAndPushECR)
	case StepTypeGCSUpload:
		s.Spec = new(StepGCSUpload)
	case StepTypeBuildAndPushGAR:
		s.Spec = new(StepBuildAndPushGAR)
	case StepTypeTest:
		s.Spec = new(StepTestIntelligence)
	case StepTypeServiceNowCreate:
		s.Spec = new(StepServiceNowCreate)
	case StepTypeServiceNowUpdate:
		s.Spec = new(StepServiceNowUpdate)
	case StepTypeIACMTerraformPlugin:
		s.Spec = new(StepIACMTerraformPlugin)
	case StepTypeIACMOpenTofuPlugin:
		s.Spec = new(StepIACMOpenTofuPlugin)
	default:
		// log.Printf("unknown step type while unmarshalling %s", s.Type)
		return fmt.Errorf("unknown step type %s", s.Type)
	}

	return json.Unmarshal(obj.Spec, s.Spec)
}
