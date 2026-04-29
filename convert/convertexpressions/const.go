package convertexpressions

// ConversionContext holds metadata for context-aware conversion
type ConversionContext struct {
	StepType    string            // Type of the step (e.g., "Run", "GCSUpload", "Http") - resolved lazily by trie
	StepID      string            // ID of the current step we're inside (from postprocess)
	StepTypeMap map[string]string // Map of step ID to step type for all steps in the pipeline

	// CurrentStepType is the type of the step we're currently inside (set by postprocess)
	// This is used as fallback when expression is "step.spec.*" (no explicit step ID)
	CurrentStepType string

	// UseFQN enables fully qualified name mode for step expressions.
	// When true, at step_node the v1Path is replaced with the step's v1 FQN base path.
	UseFQN bool

	// CurrentStepV1Path is the v1 FQN base path to the current step we're inside.
	// Example: "pipeline.stages.build.steps.restoreCache"
	// Used when UseFQN is true and expression starts with "step." (the step we're inside).
	CurrentStepV1Path string

	// StepV1PathMap maps step ID to its v1 FQN base path for all steps in the pipeline.
	// Example: {"restoreCache": "pipeline.stages.build.steps.restoreCache"}
	// Used when UseFQN is true and expression references a specific step ID (e.g., "steps.STEPID.spec.X").
	StepV1PathMap map[string]string
}

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
