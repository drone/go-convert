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
	StepTypeK8sRollingDeploy           = "k8sRollingDeployStep"
	StepTypeK8sRollingRollback         = "k8sRollingRollbackStep"
	StepTypeK8sApply                   = "k8sApplyStep"
	StepTypeK8sBGSwapServices          = "k8sBlueGreenSwapServicesStep"
	StepTypeK8sBlueGreenStageScaleDown = "k8sBlueGreenStageScaleDownStep"
	StepTypeK8sCanaryDelete            = "k8sCanaryDeleteStep"
	StepTypeK8sDelete                  = "k8sDeleteStep"
	StepTypeK8sDiff                    = "k8sDiffStep"
	StepTypeK8sRollout                 = "k8sRolloutStep"
	StepTypeK8sScale                   = "k8sScaleStep"
	StepTypeK8sDryRun                  = "k8sDryRunStep"
	StepTypeK8sTrafficRouting          = "k8sTrafficRoutingStep"
	StepTypeK8sCanaryDeploy            = "k8sCanaryDeployStep"
	StepTypeK8sBlueGreenDeploy         = "k8sBlueGreenDeployStep"

	// Helm
	StepTypeHelmBGDeploy          = "helmDeployBluegreenStep"
	StepTypeHelmBlueGreenSwapStep = "helmBluegreenSwapStep"
	StepTypeHelmCanaryDeploy      = "helmDeployCanaryStep"
	StepTypeHelmDelete            = "helmDeleteStep"
	StepTypeHelmDeploy            = "helmDeployBasicStep"
	StepTypeHelmRollback          = "helmRollbackStep"

	// Approval
	StepTypeVerify = "Verify"
)
