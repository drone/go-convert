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
	StepTypeArtifactoryUpdload         = "ArtifactoryUpload"
	StepTypeBarrier                    = "Barrier"
	StepTypeBackground                 = "Background"
	StepTypeBitrise                    = "Bitrise"
	StepTypeBuildAndPushDockerRegistry = "BuildAndPushDockerRegistry"
	StepTypeBuildAndPushECR            = "BuildAndPushECR"
	StepTypeBuildAndPushGCR            = "BuildAndPushGCR"
	StepTypeFlagConfiguration          = "FlagConfiguration"
	StepTypeGCSUpload                  = "GCSUpload"
	StepTypeHarnessApproval            = "HarnessApproval"
	StepTypePlugin                     = "Plugin"
	StepTypeRestoreCacheGCS            = "RestoreCacheGCS"
	StepTypeRestoreCacheS3             = "RestoreCacheS3"
	StepTypeRun                        = "Run"
	StepTypeRunTests                   = "RunTests"
	StepTypeS3Upload                   = "S3Upload"
	StepTypeSaveCacheGCS               = "SaveCacheGCS"
	StepTypeSaveCacheS3                = "SaveCacheS3"

	// Feature Flags
	StepTypeHTTP               = "Http"
	StepTypeJiraApproval       = "JiraApproval"
	StepTypeJiraCreate         = "JiraCreate"
	StepTypeJiraUpdate         = "JiraUpdate"
	StepTypeServiceNowApproval = "ServiceNowApproval"
	StepTypeShellScript        = "ShellScript"

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
