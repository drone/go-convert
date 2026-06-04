package convertexpressions

// RunCandidateTypes are the v0 types that all convert to v1 step.run (HTTP is
// excluded — it is detected separately by its hardcoded plugin image).
var RunCandidateTypes = []string{
	StepTypeRun,
	StepTypePlugin,
	StepTypeAction,
	StepTypeBitrise,
	StepTypeContainer,
	StepTypeShellScript,
	StepTypeShellScriptProvision,
}

// RunTestCandidateTypes are the v0 types that convert to v1 step.run-test.
var RunTestCandidateTypes = []string{
	StepTypeRunTests,
	StepTypeTest,
}

// TemplateUsesToV0Type maps a v1 step.template `uses` to the v0 step type that
// keys the conversion rules.
var TemplateUsesToV0Type = map[string]string{
	// Cache
	"saveCacheToGCS":      StepTypeSaveCacheGCS,
	"saveCacheToS3":       StepTypeSaveCacheS3,
	"restoreCacheFromGCS": StepTypeRestoreCacheGCS,
	"restoreCacheFromS3":  StepTypeRestoreCacheS3,

	// Artifact upload
	"uploadArtifactsToGCS":              StepTypeGCSUpload,
	"uploadArtifactsToS3":               StepTypeS3Upload,
	"uploadArtifactsToJfrogArtifactory": StepTypeArtifactoryUpload,

	// Build & push
	"buildAndPushToECR":    StepTypeBuildAndPushECR,
	"buildAndPushToGAR":    StepTypeBuildAndPushGAR,
	"buildAndPushToACR":    StepTypeBuildAndPushACR,
	"buildAndPushToDocker": StepTypeBuildAndPushDockerRegistry,

	// Jira / ServiceNow / Email
	"jiraCreate":       StepTypeJiraCreate,
	"jiraUpdate":       StepTypeJiraUpdate,
	"serviceNowCreate": StepTypeServiceNowCreate,
	"serviceNowUpdate": StepTypeServiceNowUpdate,
	"email":            StepTypeEmail,

	// Git clone
	"gitCloneStep": StepTypeGitClone,

	// Kubernetes
	"k8sRollingDeployStep":                  StepTypeK8sRollingDeploy,
	"k8sRollingRollbackStep":                StepTypeK8sRollingRollback,
	"k8sApplyStep":                          StepTypeK8sApply,
	"k8sBlueGreenSwapServicesSelectorsStep": StepTypeK8sBGSwapServices,
	"k8sBlueGreenStageScaleDownStep":        StepTypeK8sBlueGreenStageScaleDown,
	"k8sCanaryDeleteStep":                   StepTypeK8sCanaryDelete,
	"k8sDeleteStep":                         StepTypeK8sDelete,
	"k8sDiffStep":                           StepTypeK8sDiff,
	"k8sRolloutStep":                        StepTypeK8sRollout,
	"k8sScaleStep":                          StepTypeK8sScale,
	"k8sDryRunStep":                         StepTypeK8sDryRun,
	"k8sTrafficRoutingStep":                 StepTypeK8sTrafficRouting,
	"k8sCanaryDeployStep":                   StepTypeK8sCanaryDeploy,
	"k8sBlueGreenDeployStep":                StepTypeK8sBlueGreenDeploy,
	"k8sPatchStep":                          StepTypeK8sPatch,

	// Helm
	"helmBlueGreenDeployStep": StepTypeHelmBGDeploy,
	"helmBlueGreenSwapStep":   StepTypeHelmBlueGreenSwapStep,
	"helmCanaryDeployStep":    StepTypeHelmCanaryDeploy,
	"helmCanaryDeleteStep":    StepTypeHelmCanaryDelete,
	"helmDeleteStep":          StepTypeHelmDelete,
	"helmBasicDeployStep":     StepTypeHelmDeploy,
	"helmRollbackStep":        StepTypeHelmRollback,

	// IACM
	"terraformStep": StepTypeIACMTerraformPlugin,
	"openTofuStep":  StepTypeIACMOpenTofuPlugin,
}

// ApprovalUsesToV0Type maps a v1 step.approval `uses` to the v0 step type.
var ApprovalUsesToV0Type = map[string]string{
	"harness":    StepTypeHarnessApproval,
	"custom":     StepTypeCustomApproval,
	"jira":       StepTypeJiraApproval,
	"servicenow": StepTypeServiceNowApproval,
}
