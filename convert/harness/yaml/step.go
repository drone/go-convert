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
		ID                string             `json:"identifier,omitempty"        yaml:"identifier,omitempty"`
		Description       string             `json:"description,omitempty"       yaml:"description,omitempty"`
		Name              string             `json:"name,omitempty"              yaml:"name,omitempty"`
		Skip              string             `json:"skipCondition,omitempty"     yaml:"skipCondition,omitempty"`
		Spec              interface{}        `json:"spec,omitempty"              yaml:"spec,omitempty"`
		Timeout           Duration           `json:"timeout,omitempty"           yaml:"timeout,omitempty"`
		Type              string             `json:"type,omitempty"              yaml:"type,omitempty"`
		When              *StepWhen          `json:"when,omitempty"              yaml:"when,omitempty"`
		Env               map[string]string  `json:"envVariables,omitempty"      yaml:"envVariables,omitempty"`
		Strategy          *Strategy          `json:"strategy,omitempty"     yaml:"strategy,omitempty"`
		FailureStrategies []*FailureStrategy `json:"failureStrategies,omitempty" yaml:"failureStrategies,omitempty"`
	}

	StepGroup struct { // TODO missing failure strategies
		ID          string            `json:"identifier,omitempty"        yaml:"identifier,omitempty"`
		Description string            `json:"description,omitempty"       yaml:"description,omitempty"`
		Name        string            `json:"name,omitempty"              yaml:"name,omitempty"`
		Skip        string            `json:"skipCondition,omitempty"     yaml:"skipCondition,omitempty"`
		Steps       []*Steps          `json:"steps,omitempty"              yaml:"steps,omitempty"`
		Timeout     Duration          `json:"timeout,omitempty"           yaml:"timeout,omitempty"`
		When        *StepWhen         `json:"when,omitempty"              yaml:"when,omitempty"`
		Env         map[string]string `json:"envVariables,omitempty"      yaml:"envVariables,omitempty"`
		Strategy    *Strategy         `json:"strategy,omitempty"     yaml:"strategy,omitempty"`
	}

	//
	// Step specifications
	//

	CommonStepSpec struct {
		IncludeInfraSelectors bool     `json:"includeInfraSelectors,omitempty" yaml:"includeInfraSelectors,omitempty"`
		DelegateSelectors     []string `json:"delegateSelectors,omitempty" yaml:"delegateSelectors,omitempty"`
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
		CommonStepSpec
		ApprovalMessage                 string           `json:"approvalMessage,omitempty"                 yaml:"approvalMessage,omitempty"`
		Approvers                       *Approvers       `json:"approvers,omitempty"                       yaml:"approvers,omitempty"`
		ApproverInputs                  []*ApproverInput `json:"approverInputs,omitempty"                  yaml:"approverInputs,omitempty"`
		IncludePipelineExecutionHistory bool             `json:"includePipelineExecutionHistory,omitempty" yaml:"includePipelineExecutionHistory,omitempty"`
		IsAutoRejectEnabled             bool             `json:"isAutoRejectEnabled,omitempty"             yaml:"isAutoRejectEnabled,omitempty"`
		AutoApproval                    *AutoApproval    `json:"autoApproval,omitempty"                    yaml:"autoApproval,omitempty"`
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
		CommonStepSpec
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

	StepGitClone struct {
		CommonStepSpec
		Repository     string     `json:"repoName,omitempty"    yaml:"repoNama,omitempty"`
		ConnRef        string     `json:"connectorRef,omitempty"    yaml:"connectorRef,omitempty"`
		BuildType      string     `json:"build,omitempty"    yaml:"build,omitempty"`
		CloneDirectory string     `json:"cloneDirectory,omitempty"    yaml:"cloneDirectory,omitempty"`
		Privileged     bool       `json:"privileged,omitempty"      yaml:"privileged,omitempty"`
		Depth          string     `json:"depth,omitempty"    yaml:"cloneDirectory,omitempty"`
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
		Env             map[string]string `json:"envVariables,omitempty"    yaml:"envVariables,omitempty"`
		Command         string            `json:"command,omitempty"         yaml:"command,omitempty"`
		ConnRef         string            `json:"connectorRef,omitempty"    yaml:"connectorRef,omitempty"`
		Image           string            `json:"image,omitempty"           yaml:"image,omitempty"`
		ImagePullPolicy string            `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty"`
		Outputs         []*Output         `json:"outputVariables,omitempty" yaml:"outputVariables,omitempty"`
		Privileged      bool              `json:"privileged,omitempty"      yaml:"privileged,omitempty"`
		Resources       *Resources        `json:"resources,omitempty"       yaml:"resources,omitempty"`
		RunAsUser       string            `json:"runAsUser,omitempty"       yaml:"runAsUser,omitempty"`
		Reports         *Report           `json:"reports,omitempty"         yaml:"reports,omitempty"`
		Shell           string            `json:"shell,omitempty"           yaml:"shell,omitempty"`
	}

	StepBackground struct {
		CommonStepSpec
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

	StepShellScript struct {
		CommonStepSpec
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

	StepK8sRollingDeploy struct {
		CommonStepSpec
		SkipDryRun     bool `json:"skipDryRun,omitempty" yaml:"skipDryRun,omitempty"`
		PruningEnabled bool `json:"pruningEnabled,omitempty" yaml:"pruningEnabled,omitempty"`
	}

	StepK8sRollingRollback struct {
		CommonStepSpec
		PruningEnabled bool `json:"pruningEnabled,omitempty" yaml:"pruningEnabled,omitempty"`
	}

	StepK8sApply struct {
		CommonStepSpec
		SkipDryRun           bool          `json:"skipDryRun,omitempty" yaml:"skipDryRun,omitempty"`
		SkipSteadyStateCheck bool          `json:"skipSteadyStateCheck,omitempty" yaml:"skipSteadyStateCheck,omitempty"`
		SkipRendering        bool          `json:"skipRendering,omitempty" yaml:"skipRendering,omitempty"`
		Overrides            []interface{} `json:"overrides,omitempty" yaml:"overrides,omitempty"`
		FilePaths            []string      `json:"filePaths,omitempty" yaml:"filePaths,omitempty"`
	}

	StepK8sBGSwapServices struct {
		CommonStepSpec
		// Spec intentionally empty per v0 example
	}

	StepK8sBlueGreenStageScaleDown struct {
		CommonStepSpec
		DeleteResources bool `json:"deleteResources,omitempty" yaml:"deleteResources,omitempty"`
	}

	// CD: K8s Delete
	StepK8sDelete struct {
		CommonStepSpec
		DeleteResources *K8sDeleteResources `json:"deleteResources,omitempty" yaml:"deleteResources,omitempty"`
	}

	// CD: K8s Canary Delete (no spec fields in the provided example)
	StepK8sCanaryDelete struct {
		CommonStepSpec
		// empty
	}

	// CD: K8s Diff (no spec fields in the provided example)
	StepK8sDiff struct {
		CommonStepSpec
		// empty
	}

	// CD: K8s Rollout
	StepK8sRollout struct {
		CommonStepSpec
		Command   string               `json:"command,omitempty" yaml:"command,omitempty"`
		Resources *K8sRolloutResources `json:"resources,omitempty" yaml:"resources,omitempty"`
	}

	K8sRolloutResources struct {
		Type string                   `json:"type,omitempty" yaml:"type,omitempty"`
		Spec *K8sRolloutResourcesSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
	}

	K8sRolloutResourcesSpec struct {
		ResourceNames []string `json:"resourceNames,omitempty" yaml:"resourceNames,omitempty"`
		ManifestPaths []string `json:"manifestPaths,omitempty" yaml:"manifestPaths,omitempty"`
	}

	// K8sDeleteResources captures the delete selection and its spec.
	// Type is one of: ResourceName | ManifestPath | ReleaseName
	K8sDeleteResources struct {
		Type string                  `json:"type,omitempty" yaml:"type,omitempty"`
		Spec *K8sDeleteResourcesSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
	}

	// K8sDeleteResourcesSpec holds the possible selectors. Only one list is expected
	// to be populated depending on the Type above.
	K8sDeleteResourcesSpec struct {
		ResourceNames []string `json:"resourceNames,omitempty" yaml:"resourceNames,omitempty"`
		ManifestPaths []string `json:"manifestPaths,omitempty" yaml:"manifestPaths,omitempty"`
		ReleaseNames  []string `json:"releaseNames,omitempty" yaml:"releaseNames,omitempty"`
	}

	// CD: K8s Scale
	StepK8sScale struct {
		CommonStepSpec
		InstanceSelection    *K8sScaleInstanceSelection `json:"instanceSelection,omitempty" yaml:"instanceSelection,omitempty"`
		SkipSteadyStateCheck bool                       `json:"skipSteadyStateCheck,omitempty" yaml:"skipSteadyStateCheck,omitempty"`
		Workload             string                     `json:"workload,omitempty" yaml:"workload,omitempty"`
	}

	K8sScaleInstanceSelection struct {
		Type string                         `json:"type,omitempty" yaml:"type,omitempty"`
		Spec *K8sScaleInstanceSelectionSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
	}

	K8sScaleInstanceSelectionSpec struct {
		Count      int `json:"count,omitempty" yaml:"count,omitempty"`
		Percentage int `json:"percentage,omitempty" yaml:"percentage,omitempty"`
	}

	// CD: K8s Dry Run
	StepK8sDryRun struct {
		CommonStepSpec
		EncryptYamlOutput bool `json:"encryptYamlOutput,omitempty" yaml:"encryptYamlOutput,omitempty"`
	}

	// CD: K8s Traffic Routing
	StepK8sTrafficRouting struct {
		CommonStepSpec
		Type           string                   `json:"type,omitempty" yaml:"type,omitempty"`
		TrafficRouting *K8sTrafficRoutingConfig `json:"trafficRouting,omitempty" yaml:"trafficRouting,omitempty"`
	}

	K8sTrafficRoutingConfig struct {
		Provider string                 `json:"provider,omitempty" yaml:"provider,omitempty"`
		Spec     *K8sTrafficRoutingSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
	}

	K8sTrafficRoutingSpec struct {
		Name        string                    `json:"name,omitempty" yaml:"name,omitempty"`
		RootService string                    `json:"rootService,omitempty" yaml:"rootService,omitempty"`
		Hosts       interface{}               `json:"hosts,omitempty" yaml:"hosts,omitempty"`
		Gateways    interface{}               `json:"gateways,omitempty" yaml:"gateways,omitempty"`
		Routes      []*K8sTrafficRoutingRoute `json:"routes,omitempty" yaml:"routes,omitempty"`
	}

	K8sTrafficRoutingRoute struct {
		Route *K8sTrafficRoutingRouteSpec `json:"route,omitempty" yaml:"route,omitempty"`
	}

	K8sTrafficRoutingRouteSpec struct {
		Type         string                          `json:"type,omitempty" yaml:"type,omitempty"`
		Name         string                          `json:"name,omitempty" yaml:"name,omitempty"`
		Destinations []*K8sTrafficRoutingDestination `json:"destinations,omitempty" yaml:"destinations,omitempty"`
	}

	K8sTrafficRoutingDestination struct {
		Destination *K8sTrafficRoutingDestinationSpec `json:"destination,omitempty" yaml:"destination,omitempty"`
	}

	K8sTrafficRoutingDestinationSpec struct {
		Host   string `json:"host,omitempty" yaml:"host,omitempty"`
		Weight int    `json:"weight,omitempty" yaml:"weight,omitempty"`
	}

	// CD: K8s Canary Deploy
	StepK8sCanaryDeploy struct {
		CommonStepSpec
		SkipDryRun        bool                     `json:"skipDryRun,omitempty" yaml:"skipDryRun,omitempty"`
		InstanceSelection *K8sInstanceSelection    `json:"instanceSelection,omitempty" yaml:"instanceSelection,omitempty"`
		TrafficRouting    *K8sTrafficRoutingConfig `json:"trafficRouting,omitempty" yaml:"trafficRouting,omitempty"`
	}

	K8sInstanceSelection struct {
		Type string                    `json:"type,omitempty" yaml:"type,omitempty"`
		Spec *K8sInstanceSelectionSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
	}

	K8sInstanceSelectionSpec struct {
		Count      int `json:"count,omitempty" yaml:"count,omitempty"`
		Percentage int `json:"percentage,omitempty" yaml:"percentage,omitempty"`
	}

	// CD: K8s Blue Green Deploy
	StepK8sBlueGreenDeploy struct {
		CommonStepSpec
		SkipDryRun            bool                     `json:"skipDryRun,omitempty" yaml:"skipDryRun,omitempty"`
		PruningEnabled        bool                     `json:"pruningEnabled,omitempty" yaml:"pruningEnabled,omitempty"`
		SkipUnchangedManifest bool                     `json:"skipUnchangedManifest,omitempty" yaml:"skipUnchangedManifest,omitempty"`
		TrafficRouting        *K8sTrafficRoutingConfig `json:"trafficRouting,omitempty" yaml:"trafficRouting,omitempty"`
	}

	// CD: Helm Blue Green Deploy
	StepHelmBGDeploy struct {
		CommonStepSpec
		EnvironmentVariables        map[string]string `json:"environmentVariables,omitempty" yaml:"environmentVariables,omitempty"`
		IgnoreReleaseHistFailStatus bool              `json:"ignoreReleaseHistFailStatus,omitempty" yaml:"ignoreReleaseHistFailStatus,omitempty"`
		SkipSteadyStateCheck        bool              `json:"skipSteadyStateCheck,omitempty" yaml:"skipSteadyStateCheck,omitempty"`
		UseUpgradeInstall           bool              `json:"useUpgradeInstall,omitempty" yaml:"useUpgradeInstall,omitempty"`
	}

	// CD: Helm Blue Green Swap (no spec fields in provided example)
	StepHelmBlueGreenSwapStep struct {
		CommonStepSpec
	}

	// CD: Helm Canary Deploy
	StepHelmCanaryDeploy struct {
		CommonStepSpec
		EnvironmentVariables        map[string]string      `json:"environmentVariables,omitempty" yaml:"environmentVariables,omitempty"`
		IgnoreReleaseHistFailStatus bool                   `json:"ignoreReleaseHistFailStatus,omitempty" yaml:"ignoreReleaseHistFailStatus,omitempty"`
		SkipSteadyStateCheck        bool                   `json:"skipSteadyStateCheck,omitempty" yaml:"skipSteadyStateCheck,omitempty"`
		UseUpgradeInstall           bool                   `json:"useUpgradeInstall,omitempty" yaml:"useUpgradeInstall,omitempty"`
		InstanceSelection           *HelmInstanceSelection `json:"instanceSelection,omitempty" yaml:"instanceSelection,omitempty"`
	}

	HelmInstanceSelection struct {
		Type string                     `json:"type,omitempty" yaml:"type,omitempty"`
		Spec *HelmInstanceSelectionSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
	}

	HelmInstanceSelectionSpec struct {
		Count      int `json:"count,omitempty" yaml:"count,omitempty"`
		Percentage int `json:"percentage,omitempty" yaml:"percentage,omitempty"`
	}

	// CD: Helm Delete
	StepHelmDelete struct {
		CommonStepSpec
		DryRun               bool              `json:"dryRun,omitempty" yaml:"dryRun,omitempty"`
		CommandFlags         []string          `json:"commandFlags,omitempty" yaml:"commandFlags,omitempty"`
		EnvironmentVariables map[string]string `json:"environmentVariables,omitempty" yaml:"environmentVariables,omitempty"`
		ReleaseName          string            `json:"releaseName,omitempty" yaml:"releaseName,omitempty"`
	}

	// CD: Helm Deploy (Basic)
	StepHelmDeploy struct {
		CommonStepSpec
		SkipDryRun                  bool              `json:"skipDryRun,omitempty" yaml:"skipDryRun,omitempty"`
		IgnoreReleaseHistFailStatus bool              `json:"ignoreReleaseHistFailStatus,omitempty" yaml:"ignoreReleaseHistFailStatus,omitempty"`
		EnvironmentVariables        map[string]string `json:"environmentVariables,omitempty" yaml:"environmentVariables,omitempty"`
		SkipCleanup                 bool              `json:"skipCleanup,omitempty" yaml:"skipCleanup,omitempty"`
		SkipSteadyStateCheck        bool              `json:"skipSteadyStateCheck,omitempty" yaml:"skipSteadyStateCheck,omitempty"`
		UseUpgradeInstall           bool              `json:"useUpgradeInstall,omitempty" yaml:"useUpgradeInstall,omitempty"`
	}

	// CD: Helm Rollback
	StepHelmRollback struct {
		CommonStepSpec
		SkipDryRun           bool              `json:"skipDryRun,omitempty" yaml:"skipDryRun,omitempty"`
		EnvironmentVariables map[string]string `json:"environmentVariables,omitempty" yaml:"environmentVariables,omitempty"`
		SkipSteadyStateCheck bool              `json:"skipSteadyStateCheck,omitempty" yaml:"skipSteadyStateCheck,omitempty"`
	}

	// Feature: Jira Create
	StepJiraCreate struct {
		CommonStepSpec
		ConnectorRef string      `json:"connectorRef,omitempty" yaml:"connectorRef,omitempty"`
		ProjectKey   string      `json:"projectKey,omitempty" yaml:"projectKey,omitempty"`
		IssueType    string      `json:"issueType,omitempty" yaml:"issueType,omitempty"`
		Fields       []*Variable `json:"fields,omitempty" yaml:"fields,omitempty"`
	}

	// Feature: Jira Update
	StepJiraUpdate struct {
		CommonStepSpec
		ConnectorRef string          `json:"connectorRef,omitempty" yaml:"connectorRef,omitempty"`
		ProjectKey   string          `json:"projectKey,omitempty" yaml:"projectKey,omitempty"`
		IssueType    string          `json:"issueType,omitempty" yaml:"issueType,omitempty"`
		IssueKey     string          `json:"issueKey,omitempty" yaml:"issueKey,omitempty"`
		Fields       []*Variable     `json:"fields,omitempty" yaml:"fields,omitempty"`
		TransitionTo *JiraTransition `json:"transitionTo,omitempty" yaml:"transitionTo,omitempty"`
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
	case StepTypeAction:
		s.Spec = new(StepAction)
	case StepTypeRun:
		s.Spec = new(StepRun)
	case StepTypeBarrier:
		s.Spec = new(StepBarrier)
	case StepTypeQueue:
		s.Spec = new(StepQueue)
	case StepTypePlugin:
		s.Spec = new(StepPlugin)
	case StepTypeBuildAndPushDockerRegistry:
		s.Spec = new(StepDocker)
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
	case StepTypeJiraCreate:
		s.Spec = new(StepJiraCreate)
	case StepTypeJiraUpdate:
		s.Spec = new(StepJiraUpdate)
	case StepTypeHelmBGDeploy:
		s.Spec = new(StepHelmBGDeploy)
	case StepTypeHelmBlueGreenSwapStep:
		s.Spec = new(StepHelmBlueGreenSwapStep)
	case StepTypeHelmCanaryDeploy:
		s.Spec = new(StepHelmCanaryDeploy)
	case StepTypeHelmDelete:
		s.Spec = new(StepHelmDelete)
	case StepTypeHelmDeploy:
		s.Spec = new(StepHelmDeploy)
	case StepTypeHelmRollback:
		s.Spec = new(StepHelmRollback)
	case StepTypeWait:
		s.Spec = new(StepWait)
	default:
		fmt.Printf("unknown step type while unmarshalling %s", s.Type)
		return nil
	}

	return json.Unmarshal(obj.Spec, s.Spec)
}
