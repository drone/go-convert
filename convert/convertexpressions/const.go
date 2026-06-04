package convertexpressions

// Step type constants
const (
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
	StepTypeBuildAndPushACR            = "BuildAndPushACR"
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
	StepTypeContainer                  = "Container"
	StepTypeShellScriptProvision       = "ShellScriptProvision"

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

	// IACM
	StepTypeIACMTerraformPlugin = "IACMTerraformPlugin"
	StepTypeIACMOpenTofuPlugin  = "IACMOpenTofuPlugin"

	// Terraform
	StepTypeTerraformPlan  = "TerraformPlan"
	StepTypeTerraformApply = "TerraformApply"
)
