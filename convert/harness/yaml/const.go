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

type WhenStatus string

const (
	WhenStatusNone    WhenStatus = "None"
	WhenStatusAll                = "All"
	WhenStatusFailure            = "Failure"
	WhenStatusSuccess            = "Success"
)

type StageType string

const (
	StageTypeNone        StageType = "None"
	StageTypeApproval              = "Approval"
	StageTypeCI                    = "CI"
	StageTypeDeployment            = "Deployment"
	StageTypeFeatureFlag           = "FeatureFlag"
	StageTypeCustom                = "Custom"
	StageTypePipeline              = "Pipeline"
)

type InfraType string

const (
	InfraTypeNone             InfraType = "None"
	InfraTypeKubernetesDirect           = "KubernetesDirect"
	InfraTypeUseFromStage               = "UseFromStage"
	InfraTypeVM                         = "VM"
)

type StepType string

const (
	StepTypeNone StepType = "None"

	// Continuous Integration
	StepTypeAction                     = "Action"
	StepTypeArtifactoryUpload          = "ArtifactoryUpload"
	StepTypeBarrier                    = "Barrier"
	StepTypeBackground                 = "Background"
	StepTypeBitrise                    = "Bitrise"
	StepTypeBuildAndPushDockerRegistry = "BuildAndPushDockerRegistry"
	StepTypeAquaTrivy                  = "AquaTrivy"
	StepTypeBuildAndPushECR            = "BuildAndPushECR"
	StepTypeBuildAndPushGCR            = "BuildAndPushGCR"
	StepTypeBuildAndPushGAR            = "BuildAndPushGAR"
	StepTypeFlagConfiguration          = "FlagConfiguration"
	StepTypeGCSUpload                  = "GCSUpload"
	StepTypeHarnessApproval            = "HarnessApproval"
	StepTypePlugin                     = "Plugin"
	StepTypeQueue                      = "Queue"
	StepTypeRestoreCacheGCS            = "RestoreCacheGCS"
	StepTypeRestoreCacheS3             = "RestoreCacheS3"
	StepTypeRun                        = "Run"
	StepTypeRunTests                   = "RunTests"
	StepTypeS3Upload                   = "S3Upload"
	StepTypeSaveCacheGCS               = "SaveCacheGCS"
	StepTypeSaveCacheS3                = "SaveCacheS3"
	StepTypeTest                       = "Test"
	StepTypeGitClone                   = "GitClone"

	// Feature Flags
	StepTypeHTTP               = "Http"
	StepTypeEmail              = "Email"
	StepTypeJiraApproval       = "JiraApproval"
	StepTypeJiraCreate         = "JiraCreate"
	StepTypeJiraUpdate         = "JiraUpdate"
	StepTypeServiceNowApproval = "ServiceNowApproval"
	StepTypeServiceNowCreate   = "ServiceNowCreate"
	StepTypeServiceNowUpdate   = "ServiceNowUpdate"
	StepTypeShellScript        = "ShellScript"
	StepTypeWait               = "Wait"
	StepTypeCustomApproval     = "CustomApproval"

	// CD / Kubernetes
	StepTypeK8sRollingDeploy           = "K8sRollingDeploy"
	StepTypeK8sRollingRollback         = "K8sRollingRollback"
	StepTypeK8sApply                   = "K8sApply"
	StepTypeK8sBGSwapServices          = "K8sBGSwapServices"
	StepTypeK8sBlueGreenStageScaleDown = "K8sBlueGreenStageScaleDown"
	StepTypeK8sCanaryDelete            = "K8sCanaryDelete"
	StepTypeK8sDelete                  = "K8sDelete"
	StepTypeK8sDiff                    = "K8sDiff"
	StepTypeK8sRollout                 = "K8sRollout"
	StepTypeK8sScale                   = "K8sScale"
	StepTypeK8sDryRun                  = "K8sDryRun"
	StepTypeK8sTrafficRouting          = "K8sTrafficRouting"
	StepTypeK8sCanaryDeploy            = "K8sCanaryDeploy"
	StepTypeK8sBlueGreenDeploy         = "K8sBlueGreenDeploy"
	StepTypeK8sPatch                   = "K8sPatch"

	// Helm
	StepTypeHelmBGDeploy          = "HelmBGDeploy"
	StepTypeHelmBlueGreenSwapStep = "HelmBlueGreenSwapStep"
	StepTypeHelmCanaryDeploy      = "HelmCanaryDeploy"
	StepTypeHelmCanaryDelete      = "HelmCanaryDelete"
	StepTypeHelmDelete            = "HelmDelete"
	StepTypeHelmDeploy            = "HelmDeploy"
	StepTypeHelmRollback          = "HelmRollback"

	// Approval
	StepTypeVerify = "Verify"
)

type Shell string

const (
	ShellNone       Shell = "None"
	ShellBash             = "Bash"
	ShellPosix            = "Shell"
	ShellPowershell       = "Powershell"
)

type ImagePullPolicy string

const (
	ImagePullNone         ImagePullPolicy = "None"
	ImagePullAlways                       = "Always"
	ImagePullIfNotPresent                 = "IfNotPresent"
	ImagePullNever                        = "Never"
)
