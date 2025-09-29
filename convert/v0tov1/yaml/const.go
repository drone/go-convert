package yaml

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
	StepTypeAquaTrivy                  = "AquaTrivy"
	StepTypeBuildAndPushECR            = "BuildAndPushECR"
	StepTypeBuildAndPushGCR            = "BuildAndPushGCR"
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
	StepTypeJiraApproval       = "JiraApproval"
	StepTypeJiraCreate         = "JiraCreate"
	StepTypeJiraUpdate         = "JiraUpdate"
	StepTypeServiceNowApproval = "ServiceNowApproval"
	StepTypeShellScript        = "ShellScript"
	StepTypeWait               = "Wait"

	// CD / Kubernetes
	StepTypeK8sRollingDeploy           = "k8sRollingDeployStep@1.0.0"
	StepTypeK8sRollingRollback         = "k8sRollingRollbackStep@1.0.0"
	StepTypeK8sApply                   = "k8sApplyStep@1.0.0"
	StepTypeK8sBGSwapServices          = "k8sBlueGreenSwapServicesStep"
	StepTypeK8sBlueGreenStageScaleDown = "k8sBlueGreenStageScaleDownStep@1.0.0"
	StepTypeK8sCanaryDelete            = "k8sCanaryDeleteStep@1.0.0"
	StepTypeK8sDelete                  = "k8sDeleteStep@1.0.0"
	StepTypeK8sDiff                    = "k8sDiffStep@1.0.0"
	StepTypeK8sRollout                 = "k8sRolloutStep@1.0.0"
	StepTypeK8sScale                   = "k8sScaleStep@1.0.0"
	StepTypeK8sDryRun                  = "k8sDryRunStep@1.0.0"
	StepTypeK8sTrafficRouting          = "k8sTrafficRoutingStep@1.0.0"
	StepTypeK8sCanaryDeploy            = "k8sCanaryStep@1.0.0"
	StepTypeK8sBlueGreenDeploy         = "k8sBlueGreenDeployStep@1.0.0"

	// Helm
	StepTypeHelmBGDeploy          = "helmDeployBluegreenStep@1.0.0"
	StepTypeHelmBlueGreenSwapStep = "helmBluegreenSwapStep@1.0.0"
	StepTypeHelmCanaryDeploy      = "helmDeployCanaryStep@1.0.0"
	StepTypeHelmDelete            = "helmDeleteStep@1.0.0"
	StepTypeHelmDeploy            = "helmDeployBasicStep@1.0.0"
	StepTypeHelmRollback          = "helmRollbackStep@1.0.0"

	// Approval
	StepTypeVerify = "Verify"
)
