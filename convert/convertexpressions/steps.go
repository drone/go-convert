package convertexpressions

var StepsConversionRules = []ConversionRule{
	{"identifier", "id"},
}

// inside step.spec
//
// Rule format: {v0 field (relative to step.spec), v1 path (relative to step)}.
// For run-like v1 steps, the v1 root is `spec.*` (e.g. "spec.script",
// "spec.container.image", "spec.env.*", "spec.output"). For template-backed
// v1 steps the v1 root is the resolved action sub-step:
// "steps.<actionId>.spec.env.<KEY>" or "steps.<actionId>.spec.with.<KEY>".
//
// Sources: go-convert/v0_to_v1_field_mappings_all_steps.txt and
// template-library/v0_to_v1_template_resolved_field_mappings.txt.
var StepSpecConversionRules = map[string][]ConversionRule{

	// ============================================================
	// RUN-LIKE STEPS (v1: step.run / step.background / step.run-test)
	// Reference rule set: StepTypeRun (kept intact).
	// ============================================================

	StepTypeRun: {
		{"command", "spec.script"},
		{"envVariables", "spec.env"},
		{"connectorRef", "spec.container.connector"},
		{"image", "spec.container.image"},
		{"imagePullPolicy", "spec.container.pull"},
		{"privileged", "spec.container.privileged"},
		{"shell", "spec.shell"},
		{"outputVariables", "spec.output"},
		{"outputVariables[i].value", "(spec.output)[i].name"},
		{"outputVariables[i].name", "(spec.output)[i].alias"},
		{"reports", "spec.report"},
		{"resources.limits.cpu", "spec.container.cpu"},
		{"resources.limits.memory", "spec.container.memory"},
		{"runAsUser", "spec.container.user"},
	},

	StepTypeBackground: {
		{"command", "spec.script"},
		{"image", "spec.container.image"},
		{"connectorRef", "spec.container.connector"},
		{"imagePullPolicy", "spec.container.pull"},
		{"privileged", "spec.container.privileged"},
		{"runAsUser", "spec.container.user"},
		{"envVariables", "spec.env"},
		{"shell", "spec.shell"},
		{"entrypoint", "spec.container.entrypoint"},
		{"portBindings", "spec.container.ports"},
		{"resources.limits.cpu", "spec.container.cpu"},
		{"resources.limits.memory", "spec.container.memory"},
		{"reports", "spec.report"},
	},

	// Plugin -> v1 step.run (container image + settings)
	StepTypePlugin: {
		{"envVariables", "spec.env"},
		{"connectorRef", "spec.container.connector"},
		{"image", "spec.container.image"},
		{"imagePullPolicy", "spec.container.pull"},
		{"privileged", "spec.container.privileged"},
		{"runAsUser", "spec.container.user"},
		{"entrypoint", "spec.container.entrypoint"},
		{"resources.limits.cpu", "spec.container.cpu"},
		{"resources.limits.memory", "spec.container.memory"},
		{"settings", "spec.with"},
		{"reports", "spec.report"},
	},

	// Action -> v1 step.run emitting a plugin command.
	// spec.uses is derived into the run script; spec.with/env merge into env.
	StepTypeAction: {
		{"env", "spec.env"},
	},

	// Bitrise -> v1 step.run emitting a plugin command. Same shape as Action.
	StepTypeBitrise: {
		{"env", "spec.env"},
	},

	// Container -> v1 step.run
	StepTypeContainer: {
		{"command", "spec.script"},
		{"envVariables", "spec.env"},
		{"image", "spec.container.image"},
		{"connectorRef", "spec.container.connector"},
		{"imagePullPolicy", "spec.container.pull"},
		{"entrypoint", "spec.container.entrypoint"},
		{"privileged", "spec.container.privileged"},
		{"runAsUser", "spec.container.user"},
		{"shell", "spec.shell"},
		{"outputVariables", "spec.output"},
		{"outputVariables[i].value", "(spec.output)[i].name"},
		{"outputVariables[i].name", "(spec.output)[i].alias"},
	},

	// ShellScript -> v1 step.run (shell script source)
	StepTypeShellScript: {
		{"shell", "spec.shell"},
		{"source.spec.script", "spec.script"},
		{"environmentVariables.*", "spec.env.*"},
		{"outputVariables", "spec.output"},
		{"outputVariables[i].value", "(spec.output)[i].name"},
		{"outputVariables[i].name", "(spec.output)[i].alias"},
	},

	// ShellScriptProvision -> v1 step.run (shell provisioner)
	StepTypeShellScriptProvision: {
		{"source.spec.script", "spec.script"},
		{"environmentVariables.*", "spec.env.*"},
	},

	// Test (TestV2) -> v1 step.run-test.
	// Expression paths resolve under `spec.*` like run.
	StepTypeTest: {
		{"command", "spec.script"},
		{"envVariables", "spec.env"},
		{"image", "spec.container.image"},
		{"connectorRef", "spec.container.connector"},
		{"imagePullPolicy", "spec.container.pull"},
		{"privileged", "spec.container.privileged"},
		{"runAsUser", "spec.container.user"},
		{"resources.limits.cpu", "spec.container.cpu"},
		{"resources.limits.memory", "spec.container.memory"},
		{"shell", "spec.shell"},
		{"reports", "spec.report"},
		{"outputVariables", "spec.output"},
		{"outputVariables[i].value", "(spec.output)[i].name"},
		{"outputVariables[i].name", "(spec.output)[i].alias"},
		{"globs", "spec.match"},
	},

	// RunTests -> v1 step.run-test. The command is DERIVED from buildTool/args;
	// we keep the common container/env mappings only.
	StepTypeRunTests: {
		{"envVariables", "spec.env"},
		{"image", "spec.container.image"},
		{"connectorRef", "spec.container.connector"},
		{"imagePullPolicy", "spec.container.pull"},
		{"privileged", "spec.container.privileged"},
		{"runAsUser", "spec.container.user"},
		{"resources.limits.cpu", "spec.container.cpu"},
		{"resources.limits.memory", "spec.container.memory"},
		{"shell", "spec.shell"},
		{"reports", "spec.report"},
		{"outputVariables", "spec.output"},
		{"outputVariables[i].value", "(spec.output)[i].name"},
		{"outputVariables[i].name", "(spec.output)[i].alias"},
		{"testGlobs", "spec.match"},
	},

	// HTTP -> v1 step.run emitting the http plugin. Fields land on run env.
	StepTypeHTTP: {
		{"url", "spec.env.PLUGIN_URL"},
		{"method", "spec.env.PLUGIN_METHOD"},
		{"headers", "spec.env.PLUGIN_HEADERS"},
		{"requestBody", "spec.env.PLUGIN_BODY"},
		{"assertion", "spec.env.PLUGIN_ASSERTION"},
		{"inputVariables", "spec.env.PLUGIN_ENV_VARS"},
		{"outputVariables", "spec.env.PLUGIN_OUTPUT_VARS"},
		{"certificate", "spec.env.PLUGIN_CLIENT_CERT"},
		{"certificateKey", "spec.env.PLUGIN_CLIENT_KEY"},
	},

	// ============================================================
	// CONTROL FLOW / UTILITY STEPS
	// ============================================================

	// Wait -> v1 step.wait
	StepTypeWait: {
		{"duration", "spec.duration"},
	},

	// Barrier -> v1 step.barrier
	StepTypeBarrier: {
		{"barrierRef", "spec.name"},
	},

	// Queue -> v1 step.queue
	StepTypeQueue: {
		{"key", "spec.key"},
		{"scope", "spec.scope"},
	},

	// ============================================================
	// APPROVAL STEPS (v1: step.approval)
	// ============================================================

	StepTypeHarnessApproval: {
		{"approvalMessage", "spec.with.message"},
		{"includePipelineExecutionHistory", "spec.with.execution-details"},
		{"isAutoRejectEnabled", "spec.with.auto-reject"},
		{"approvers.minimumCount", "spec.with.approvers-min-count"},
		{"approvers.disallowPipelineExecutor", "spec.with.block-executor"},
		{"approvers.userGroups", "spec.with.user-groups"},
		{"autoApproval.comments", "spec.with.comments"},
		{"autoApproval.scheduledDeadline.time", "spec.with.deadline"},
		{"autoApproval.scheduledDeadline.timeZone", "spec.with.timezone"},
	},

	StepTypeCustomApproval: {
		{"scriptTimeout", "spec.with.script-timeout"},
		{"retryInterval", "spec.with.retry"},
		{"shell", "spec.with.run.shell"},
		{"source.spec.script", "spec.with.run.script"},
		{"environmentVariables.*", "spec.with.run.env.*"},
		{"outputVariables", "spec.with.run.outputs"},
	},

	StepTypeJiraApproval: {
		{"retryInterval", "spec.with.retry"},
		{"connectorRef", "spec.with.run.env.PLUGIN_HARNESS_CONNECTOR"},
		{"issueKey", "spec.with.run.env.PLUGIN_ISSUE_KEY"},
	},

	StepTypeServiceNowApproval: {
		{"retryInterval", "spec.with.retry"},
		{"changeWindow.startField", "spec.with.change-window.start"},
		{"changeWindow.endField", "spec.with.change-window.end"},
		{"connectorRef", "spec.with.run.env.PLUGIN_HARNESS_CONNECTOR"},
		{"ticketType", "spec.with.run.env.PLUGIN_TICKET_TYPE"},
		{"ticketNumber", "spec.with.run.env.PLUGIN_TICKET_NUMBER"},
	},

	// ============================================================
	// INTEGRATION STEPS (template-backed: step.template.with.*)
	// Resolved via sub-step spec.env.PLUGIN_* keys.
	// ============================================================

	// JiraCreate -> uses="jiraCreate", action id = jiraCreate
	StepTypeJiraCreate: {
		{"connectorRef", "steps.jiraCreate.spec.env.PLUGIN_HARNESS_CONNECTOR"},
		{"projectKey", "steps.jiraCreate.spec.env.PLUGIN_PROJECT"},
		{"issueType", "steps.jiraCreate.spec.env.PLUGIN_ISSUE_TYPE"},
		{"fields", "steps.jiraCreate.spec.env.PLUGIN_FIELDS"},
	},

	// JiraUpdate -> uses="jiraUpdate", action id = jiraUpdate
	StepTypeJiraUpdate: {
		{"connectorRef", "steps.jiraUpdate.spec.env.PLUGIN_HARNESS_CONNECTOR"},
		{"issueKey", "steps.jiraUpdate.spec.env.PLUGIN_ISSUE_KEY"},
		{"fields", "steps.jiraUpdate.spec.env.PLUGIN_FIELDS"},
		{"transitionTo.status", "steps.jiraUpdate.spec.env.PLUGIN_TRANSITION_STATUS"},
		{"transitionTo.transitionName", "steps.jiraUpdate.spec.env.PLUGIN_TRANSITION"},
	},

	// ServiceNowCreate -> uses="serviceNowCreate", action id = create
	StepTypeServiceNowCreate: {
		{"connectorRef", "steps.create.spec.env.PLUGIN_HARNESS_CONNECTOR"},
		{"ticketType", "steps.create.spec.env.PLUGIN_TICKET_TYPE"},
		{"fields", "steps.create.spec.env.PLUGIN_FIELDS"},
	},

	// ServiceNowUpdate -> uses="serviceNowUpdate", action id = update
	StepTypeServiceNowUpdate: {
		{"connectorRef", "steps.update.spec.env.PLUGIN_HARNESS_CONNECTOR"},
		{"ticketType", "steps.update.spec.env.PLUGIN_TICKET_TYPE"},
		{"ticketNumber", "steps.update.spec.env.PLUGIN_TICKET_NUMBER"},
		{"fields", "steps.update.spec.env.PLUGIN_FIELDS"},
	},

	// Email -> uses="email", action id = email
	StepTypeEmail: {
		{"to", "steps.email.spec.env.PLUGIN_EMAIL_IDS"},
		{"cc", "steps.email.spec.env.PLUGIN_CC_EMAIL_IDS"},
		{"subject", "steps.email.spec.env.PLUGIN_SUBJECT"},
		{"body", "steps.email.spec.env.PLUGIN_BODY"},
		{"toUserGroups", "steps.email.spec.env.PLUGIN_TO_USER_GROUPS"},
		{"ccUserGroups", "steps.email.spec.env.PLUGIN_CC_USER_GROUPS"},
	},

	// ============================================================
	// ARTIFACT UPLOAD STEPS (template-backed; sub-step uses run.with)
	// ============================================================

	// ArtifactoryUpload -> uses="uploadArtifactsToJfrogArtifactory", id = jfrogArtifactory
	StepTypeArtifactoryUpload: {
		{"sourcePath", "steps.jfrogArtifactory.spec.with.SOURCE"},
		{"target", "steps.jfrogArtifactory.spec.with.TARGET"},
	},

	// GCSUpload -> uses="uploadArtifactsToGCS", id = gcsUpload
	StepTypeGCSUpload: {
		{"sourcePath", "steps.gcsUpload.spec.with.SOURCE"},
		{"bucket", "steps.gcsUpload.spec.with.TARGET"},
		{"target", "steps.gcsUpload.spec.with.TARGET"},
	},

	// S3Upload -> uses="uploadArtifactsToS3", id = s3Upload
	StepTypeS3Upload: {
		{"bucket", "steps.s3Upload.spec.with.BUCKET"},
		{"region", "steps.s3Upload.spec.with.REGION"},
		{"endpoint", "steps.s3Upload.spec.with.ENDPOINT"},
		{"sourcePath", "steps.s3Upload.spec.with.SOURCE"},
		{"target", "steps.s3Upload.spec.with.TARGET"},
	},

	// ============================================================
	// CACHE STEPS (template-backed; sub-step uses run.with)
	// ============================================================

	// SaveCacheS3 -> uses="saveCacheToS3", id = saveCacheS3
	StepTypeSaveCacheS3: {
		{"bucket", "steps.saveCacheS3.spec.with.BUCKET"},
		{"key", "steps.saveCacheS3.spec.with.CACHE_KEY"},
		{"sourcePaths", "steps.saveCacheS3.spec.with.MOUNT"},
		{"region", "steps.saveCacheS3.spec.with.REGION"},
		{"endpoint", "steps.saveCacheS3.spec.with.ENDPOINT"},
		{"archiveFormat", "steps.saveCacheS3.spec.with.ARCHIVE_FORMAT"},
		{"pathStyle", "steps.saveCacheS3.spec.with.PATH_STYLE"},
		{"override", "steps.saveCacheS3.spec.with.OVERRIDE"},
	},

	// SaveCacheGCS -> uses="saveCacheToGCS", id = saveCacheGCS
	StepTypeSaveCacheGCS: {
		{"bucket", "steps.saveCacheGCS.spec.with.BUCKET"},
		{"key", "steps.saveCacheGCS.spec.with.CACHE_KEY"},
		{"sourcePaths", "steps.saveCacheGCS.spec.with.MOUNT"},
		{"archiveFormat", "steps.saveCacheGCS.spec.with.ARCHIVE_FORMAT"},
		{"override", "steps.saveCacheGCS.spec.with.OVERRIDE"},
	},

	// RestoreCacheS3 -> uses="restoreCacheFromS3", id = restoreCacheS3
	StepTypeRestoreCacheS3: {
		{"bucket", "steps.restoreCacheS3.spec.with.BUCKET"},
		{"key", "steps.restoreCacheS3.spec.with.CACHE_KEY"},
		{"region", "steps.restoreCacheS3.spec.with.REGION"},
		{"endpoint", "steps.restoreCacheS3.spec.with.ENDPOINT"},
		{"archiveFormat", "steps.restoreCacheS3.spec.with.ARCHIVE_FORMAT"},
		{"pathStyle", "steps.restoreCacheS3.spec.with.PATH_STYLE"},
		{"failIfKeyNotFound", "steps.restoreCacheS3.spec.with.FAIL_RESTORE_IF_KEY_NOT_PRESENT"},
	},

	// RestoreCacheGCS -> uses="restoreCacheFromGCS", id = restoreCacheGCS
	StepTypeRestoreCacheGCS: {
		{"bucket", "steps.restoreCacheGCS.spec.with.BUCKET"},
		{"key", "steps.restoreCacheGCS.spec.with.CACHE_KEY"},
		{"archiveFormat", "steps.restoreCacheGCS.spec.with.ARCHIVE_FORMAT"},
		{"failIfKeyNotFound", "steps.restoreCacheGCS.spec.with.FAIL_RESTORE_IF_KEY_NOT_PRESENT"},
	},

	// ============================================================
	// GITCLONE (template-backed; most fields land on env.DRONE_*)
	// ============================================================

	// GitClone -> uses="gitCloneStep", id = gitClone
	StepTypeGitClone: {
		{"build.spec.branch", "steps.gitClone.spec.env.DRONE_COMMIT_BRANCH"},
		{"build.spec.tag", "steps.gitClone.spec.env.DRONE_TAG"},
		{"build.spec.commitSha", "steps.gitClone.spec.env.DRONE_COMMIT_SHA"},
		{"cloneDirectory", "steps.gitClone.spec.env.DRONE_WORKSPACE"},
		{"depth", "steps.gitClone.spec.with.DEPTH"},
		{"sparseCheckout", "steps.gitClone.spec.env.DRONE_NETRC_SPARSE_CHECKOUT"},
		{"preFetchCommand", "steps.gitClone.spec.env.DRONE_NETRC_PRE_FETCH"},
		{"submoduleStrategy", "steps.gitClone.spec.env.DRONE_NETRC_SUBMODULE_STRATEGY"},
		{"outputFilePathsContent", "steps.gitClone.spec.with.OUTPUT_FILE_PATHS_CONTENT"},
		{"lfs", "steps.gitClone.spec.env.DRONE_NETRC_LFS_ENABLED"},
		{"debug", "steps.gitClone.spec.env.DRONE_NETRC_DEBUG"},
		{"fetchTags", "steps.gitClone.spec.env.DRONE_NETRC_FETCH_TAGS"},
	},

	// ============================================================
	// IACM STEPS (template-backed; template root is a single run step,
	// so resolved path is `spec.with.*` / `spec.env.*`)
	// ============================================================

	StepTypeIACMTerraformPlugin: {
		{"command", "spec.with.COMMAND"},
		{"target", "spec.with.PLAN_TARGET"},
		{"replace", "spec.with.PLAN_REPLACE"},
		{"importVars", "spec.with.IMPORT"},
	},

	StepTypeIACMOpenTofuPlugin: {
		{"command", "spec.with.COMMAND"},
		{"target", "spec.with.PLAN_TARGET"},
		{"replace", "spec.with.PLAN_REPLACE"},
		{"importVars", "spec.with.IMPORT"},
	},

	// ============================================================
	// K8S STEPS
	// ============================================================

	// K8sRollingDeploy (v1: k8sRollingDeployStep)
	// Step group: k8sRollingPrepareAction, k8sApplyAction, k8sSteadyStateCheckAction
	StepTypeK8sRollingDeploy: {
		{"skipDryRun", "steps.k8sApplyAction.spec.env.PLUGIN_SKIP_DRY_RUN"},
		{"pruningEnabled", "steps.k8sRollingPrepareAction.spec.env.PLUGIN_PRUNING_ENABLED"},
		// pruningEnabled also maps to: steps.k8sApplyAction.spec.env.PLUGIN_RELEASE_PRUNING_ENABLED
		{"flags", "steps.k8sApplyAction.spec.env.PLUGIN_FLAGS_JSON"},
	},

	// K8sRollingRollback (v1: k8sRollingRollbackStep)
	// Step group: k8sRollingRollbackAction
	StepTypeK8sRollingRollback: {
		{"pruningEnabled", "steps.k8sRollingRollbackAction.spec.env.PLUGIN_PRUNING"},
		{"flags", "steps.k8sRollingRollbackAction.spec.env.PLUGIN_FLAGS_JSON"},
	},

	// K8sApply (v1: k8sApplyStep)
	// Step group: k8sApplyAction, k8sSteadyStateCheckAction
	StepTypeK8sApply: {
		{"filePaths", "steps.k8sApplyAction.spec.env.PLUGIN_MANIFEST_PATH"},
		// COMPLEX: each path prefixed with "<+runtime.manifestPath>/"
		{"skipDryRun", "steps.k8sApplyAction.spec.env.PLUGIN_SKIP_DRY_RUN"},
		// skipSteadyStateCheck -> NO_ENV_MAPPING (controls conditional execution)
		// skipRendering -> NO_TEMPLATE_MAPPING
		{"flags", "steps.k8sApplyAction.spec.env.PLUGIN_FLAGS_JSON"},
	},

	// K8sBGSwapServices (v1: k8sBlueGreenSwapServicesSelectorsStep)
	// Step group: k8sBlueGreenSwapServicesSelectorsAction
	// All fields are COMPLEX (populated from exported variables, not from v0 spec)
	StepTypeK8sBGSwapServices: {
		{"stable_service", "steps.k8sBlueGreenSwapServicesSelectorsAction.spec.env.PLUGIN_STABLE_SERVICE"},
		{"stage_service", "steps.k8sBlueGreenSwapServicesSelectorsAction.spec.env.PLUGIN_STAGE_SERVICE"},
		{"is_openshift", "steps.k8sBlueGreenSwapServicesSelectorsAction.spec.env.HARNESS_IS_OPENSHIFT"},
	},

	// K8sBlueGreenStageScaleDown (v1: k8sBlueGreenStageScaleDownStep)
	// Step group: k8sBlueGreenStageScaleDownAction, (conditional) k8sDeleteAction
	// deleteResources -> NO_ENV_MAPPING (controls conditional execution)
	StepTypeK8sBlueGreenStageScaleDown: {},

	// K8sCanaryDelete (v1: k8sCanaryDeleteStep)
	// Step group: k8sCanaryDeleteAction (uses k8sDeleteAction)
	// All fields are COMPLEX (populated from exported variables)
	StepTypeK8sCanaryDelete: {
		{"resources", "steps.k8sCanaryDeleteAction.spec.env.PLUGIN_RESOURCES"},
		{"is_openshift", "steps.k8sCanaryDeleteAction.spec.env.HARNESS_IS_OPENSHIFT"},
		// select_delete_resources (Rollback) -> NO_ENV_MAPPING
	},

	// K8sDiff (v1: k8sDiffStep)
	// Step group: k8sDiffAction
	// No field mappings (v0 spec and v1 with are both empty)
	StepTypeK8sDiff: {},

	// K8sRollout (v1: k8sRolloutStep)
	// Step group: k8sRolloutAction (id: k8sRolloutStep in step group)
	StepTypeK8sRollout: {
		{"command", "steps.k8sRolloutStep.spec.env.PLUGIN_COMMAND"},
		// resources.type -> NO_ENV_MAPPING (controls UI visibility)
		{"resources.spec.resourceNames", "steps.k8sRolloutStep.spec.env.PLUGIN_RESOURCES"},
		// only when type=ResourceName
		{"resources.spec.manifestPaths", "steps.k8sRolloutStep.spec.env.PLUGIN_MANIFESTS"},
		// only when type=ManifestPath
		{"flags", "steps.k8sRolloutStep.spec.env.PLUGIN_FLAGS_JSON"},
	},

	// K8sScale (v1: k8sScaleStep)
	// Step group: k8sScaleAction, (conditional) k8sSteadyStateCheckAction
	StepTypeK8sScale: {
		{"instanceSelection.type", "steps.k8sScaleAction.spec.env.PLUGIN_INSTANCES_UNIT_TYPE"},
		// COMPLEX: "Count" -> "count", "Percentage" -> "percentage"
		{"instanceSelection.spec.count", "steps.k8sScaleAction.spec.env.PLUGIN_INSTANCES"},
		// only when type=Count
		{"instanceSelection.spec.percentage", "steps.k8sScaleAction.spec.env.PLUGIN_INSTANCES"},
		// only when type=Percentage
		{"workload", "steps.k8sScaleAction.spec.env.PLUGIN_WORKLOAD"},
		// skipSteadyStateCheck -> NO_ENV_MAPPING (controls conditional execution)
	},

	// K8sDryRun (v1: k8sDryRunStep)
	// Step group: k8sDryRunAction (uses k8sApplyAction with only_dry_run=true)
	StepTypeK8sDryRun: {
		{"encryptYamlOutput", "steps.k8sDryRunAction.spec.env.PLUGIN_ENCRYPT_YAML_OUTPUT"},
	},

	// K8sDelete (v1: k8sDeleteStep)
	// Step group: k8sDeleteAction
	StepTypeK8sDelete: {
		// deleteResources.type -> NO_ENV_MAPPING (controls UI visibility)
		{"deleteResources.spec.resourceNames", "steps.k8sDeleteAction.spec.env.PLUGIN_RESOURCES"},
		// only when type=ResourceName
		{"deleteResources.spec.manifestPaths", "steps.k8sDeleteAction.spec.env.PLUGIN_MANIFESTS"},
		// COMPLEX: each path prefixed with "<+runtime.manifestPath>/"; only when type=ManifestPath
		{"deleteResources.spec.deleteNamespace", "steps.k8sDeleteAction.spec.env.PLUGIN_INCLUDE_NAMESPACES"},
		// only when type=ReleaseName
		{"flags", "steps.k8sDeleteAction.spec.env.PLUGIN_FLAGS_JSON"},
	},

	// K8sTrafficRouting (v1: k8sTrafficRoutingStep)
	// Step group: k8sTrafficShiftAction
	StepTypeK8sTrafficRouting: {
		{"trafficRouting.provider", "steps.k8sTrafficShiftAction.spec.env.PLUGIN_PROVIDER"},
		{"trafficRouting.spec.name", "steps.k8sTrafficShiftAction.spec.env.PLUGIN_RESOURCE_NAME"},
		{"trafficRouting.spec.hosts", "steps.k8sTrafficShiftAction.spec.env.PLUGIN_HOSTNAMES"},
		{"trafficRouting.spec.gateways", "steps.k8sTrafficShiftAction.spec.env.PLUGIN_GATEWAYS"},
		{"trafficRouting.spec.routes", "steps.k8sTrafficShiftAction.spec.env.PLUGIN_ROUTES"},
		// COMPLEX: converted to JSON string
		// trafficRouting.spec.rootService -> NO_TEMPLATE_MAPPING
		// type -> NO_TEMPLATE_MAPPING (v1 always sets config="new")
	},

	// K8sCanaryDeploy (v1: k8sCanaryDeployStep)
	// Step group: k8sCanaryPrepareAction, k8sApplyAction, k8sSteadyStateCheckAction, (conditional) k8sTrafficRoutingAction
	StepTypeK8sCanaryDeploy: {
		{"instanceSelection.type", "steps.k8sCanaryPrepareAction.spec.env.PLUGIN_INSTANCES_UNIT_TYPE"},
		// COMPLEX: "Count" -> "count", "Percentage" -> "percentage"
		{"instanceSelection.spec.count", "steps.k8sCanaryPrepareAction.spec.env.PLUGIN_INSTANCES"},
		// only when type=Count
		{"instanceSelection.spec.percentage", "steps.k8sCanaryPrepareAction.spec.env.PLUGIN_INSTANCES"},
		// only when type=Percentage
		{"skipDryRun", "steps.k8sApplyAction.spec.env.PLUGIN_SKIP_DRY_RUN"},
		{"trafficRouting.provider", "steps.k8sTrafficRoutingAction.spec.env.PLUGIN_PROVIDER"},
		{"trafficRouting.spec.name", "steps.k8sTrafficRoutingAction.spec.env.PLUGIN_RESOURCE_NAME"},
		{"trafficRouting.spec.hosts", "steps.k8sTrafficRoutingAction.spec.env.PLUGIN_HOSTNAMES"},
		{"trafficRouting.spec.gateways", "steps.k8sTrafficRoutingAction.spec.env.PLUGIN_GATEWAYS"},
		// {"trafficRouting.spec.routes", "steps.k8sTrafficRoutingAction.spec.env.PLUGIN_ROUTES"},
		// COMPLEX: converted to JSON string
		// trafficRouting.spec.rootService -> NO_TEMPLATE_MAPPING
		{"flags", "steps.k8sApplyAction.spec.env.PLUGIN_FLAGS_JSON"},
	},

	// K8sBlueGreenDeploy (v1: k8sBlueGreenDeployStep)
	// Step group: k8sBlueGreenPrepareAction, k8sApplyAction, k8sSteadyStateCheckAction, (conditional) k8sTrafficRoutingAction
	StepTypeK8sBlueGreenDeploy: {
		{"skipDryRun", "steps.k8sApplyAction.spec.env.PLUGIN_SKIP_DRY_RUN"},
		{"pruningEnabled", "steps.k8sApplyAction.spec.env.PLUGIN_RELEASE_PRUNING_ENABLED"},
		{"skipUnchangedManifest", "steps.k8sBlueGreenPrepareAction.spec.env.PLUGIN_SKIP_UNCHANGED_MANIFEST"},
		{"trafficRouting.provider", "steps.k8sTrafficRoutingAction.spec.env.PLUGIN_PROVIDER"},
		{"trafficRouting.spec.name", "steps.k8sTrafficRoutingAction.spec.env.PLUGIN_RESOURCE_NAME"},
		{"trafficRouting.spec.hosts", "steps.k8sTrafficRoutingAction.spec.env.PLUGIN_HOSTNAMES"},
		{"trafficRouting.spec.gateways", "steps.k8sTrafficRoutingAction.spec.env.PLUGIN_GATEWAYS"},
		{"trafficRouting.spec.routes", "steps.k8sTrafficRoutingAction.spec.env.PLUGIN_ROUTES"},
		// COMPLEX: converted to JSON string
		// trafficRouting.spec.rootService -> NO_TEMPLATE_MAPPING
		{"flags", "steps.k8sApplyAction.spec.env.PLUGIN_FLAGS_JSON"},
	},

	// K8sPatch (v1: k8sPatchStep)
	// Step group: k8sPatchAction, (conditional) k8sSteadyStateCheckAction
	StepTypeK8sPatch: {
		{"workload", "steps.k8sPatchAction.spec.env.PLUGIN_WORKLOAD"},
		// skipSteadyStateCheck -> NO_ENV_MAPPING (controls conditional execution)
		{"mergeStrategyType", "steps.k8sPatchAction.spec.env.PLUGIN_MERGE_STRATEGY"},
		// COMPLEX: value is lowercased: "Json"->"json", "Strategic"->"strategic", "Merge"->"merge"
		// source.type -> NO_TEMPLATE_MAPPING (used to determine how to extract content)
		{"source.spec.content", "steps.k8sPatchAction.spec.env.PLUGIN_CONTENT"},
		// only when source.type=Inline
		// source.spec (Git/GitLab/etc) -> NO_TEMPLATE_MAPPING (remote source not supported)
		// recordChangeCause -> NO_TEMPLATE_MAPPING
	},

	// ============================================================
	// HELM STEPS
	// ============================================================

	// HelmBGDeploy (v1: helmBlueGreenDeployStep)
	// Step group: helmBluegreenDeployAction (uses helmDeployAction with strategy=blue-green), (conditional) helmTestAction
	StepTypeHelmBGDeploy: {
		{"ignoreReleaseHistFailStatus", "steps.helmBluegreenDeployAction.spec.env.PLUGIN_IGNORE_HISTORY_FAILURE"},
		{"skipSteadyStateCheck", "steps.helmBluegreenDeployAction.spec.env.PLUGIN_SKIP_STEADY_STATE_CHECK"},
		// useUpgradeInstall -> NO_ENV_MAPPING
		// runChartTests -> NO_ENV_MAPPING (controls conditional execution of helmTestAction)
		{"environmentVariables", "steps.helmBluegreenDeployAction.spec.env.PLUGIN_ENV_VARS"},
		// COMPLEX: map[string]string converted to []map with key/value entries
		{"flags", "steps.helmBluegreenDeployAction.spec.env.PLUGIN_FLAGS"},
	},

	// HelmBlueGreenSwapStep (v1: helmBlueGreenSwapStep)
	// Step group: helmBluegreenSwapAction
	StepTypeHelmBlueGreenSwapStep: {
		{"flags", "steps.helmBluegreenSwapAction.spec.env.PLUGIN_FLAGS"},
	},

	// HelmCanaryDeploy (v1: helmCanaryDeployStep)
	// Step group: helmBasicDeployAction (uses helmDeployAction with strategy=canary), (conditional) helmTestAction
	StepTypeHelmCanaryDeploy: {
		{"ignoreReleaseHistFailStatus", "steps.helmBasicDeployAction.spec.env.PLUGIN_IGNORE_HISTORY_FAILURE"},
		{"skipSteadyStateCheck", "steps.helmBasicDeployAction.spec.env.PLUGIN_SKIP_STEADY_STATE_CHECK"},
		// useUpgradeInstall -> NO_ENV_MAPPING
		// runChartTests -> NO_ENV_MAPPING (controls conditional execution of helmTestAction)
		{"environmentVariables", "steps.helmBasicDeployAction.spec.env.PLUGIN_ENV_VARS"},
		// COMPLEX: map[string]string converted to []map with key/value entries
		{"instanceSelection.type", "steps.helmBasicDeployAction.spec.env.PLUGIN_INSTANCES_UNIT_TYPE"},
		// COMPLEX: "Count" -> "count", "Percentage" -> "percentage"
		{"instanceSelection.spec.count", "steps.helmBasicDeployAction.spec.env.PLUGIN_INSTANCES"},
		// only when type=Count
		{"instanceSelection.spec.percentage", "steps.helmBasicDeployAction.spec.env.PLUGIN_INSTANCES"},
		// only when type=Percentage
		{"flags", "steps.helmBasicDeployAction.spec.env.PLUGIN_FLAGS"},
	},

	// HelmCanaryDelete (v1: helmCanaryDeleteStep)
	// Step group: helmUninstallAction (uses helmDeployAction — no with params from v0)
	// No field mappings (v0 spec is empty)
	StepTypeHelmCanaryDelete: {},

	// HelmDelete (v1: helmDeleteStep)
	// Step group: helmUninstallAction
	StepTypeHelmDelete: {
		{"releaseName", "steps.helmUninstallAction.spec.env.PLUGIN_RELEASE_NAME"},
		{"dryRun", "steps.helmUninstallAction.spec.env.PLUGIN_DRY_RUN"},
		{"commandFlags", "steps.helmUninstallAction.spec.env.PLUGIN_FLAGS"},
		{"environmentVariables", "steps.helmUninstallAction.spec.env.PLUGIN_ENV_VARS"},
		// COMPLEX: map[string]string converted to []map with key/value entries
	},

	// HelmDeploy (v1: helmBasicDeployStep)
	// Step group: helmBasicDeployAction (uses helmDeployAction with strategy=basic), (conditional) helmTestAction
	StepTypeHelmDeploy: {
		{"ignoreReleaseHistFailStatus", "steps.helmBasicDeployAction.spec.env.PLUGIN_IGNORE_HISTORY_FAILURE"},
		{"skipSteadyStateCheck", "steps.helmBasicDeployAction.spec.env.PLUGIN_SKIP_STEADY_STATE_CHECK"},
		// useUpgradeInstall -> NO_ENV_MAPPING
		// runChartTests -> NO_ENV_MAPPING (controls conditional execution of helmTestAction)
		{"environmentVariables", "steps.helmBasicDeployAction.spec.env.PLUGIN_ENV_VARS"},
		// COMPLEX: map[string]string converted to []map with key/value entries
		// skipDryRun -> NO_TEMPLATE_MAPPING (intentionally omitted in v1)
		// skipCleanup -> NO_TEMPLATE_MAPPING (intentionally omitted in v1)
		{"flags", "steps.helmBasicDeployAction.spec.env.PLUGIN_FLAGS"},
	},

	// HelmRollback (v1: helmRollbackStep)
	// Step group: helmRollbackAction, (conditional) helmTestAction
	StepTypeHelmRollback: {
		{"skipSteadyStateCheck", "steps.helmRollbackAction.spec.env.PLUGIN_SKIP_STEADY_STATE_CHECK"},
		// runChartTests -> NO_ENV_MAPPING (controls conditional execution of helmTestAction)
		{"environmentVariables", "steps.helmRollbackAction.spec.env.PLUGIN_ENV_VARS"},
		// COMPLEX: map[string]string converted to []map with key/value entries
		// skipDryRun -> NO_TEMPLATE_MAPPING (not used in conversion)
		{"flags", "steps.helmRollbackAction.spec.env.PLUGIN_FLAGS"},
	},

	// ============================================================
	// BUILD & PUSH STEPS
	// ============================================================
	// Note: Build & Push steps use conditional plugin selection (buildx/docker/kaniko).
	// Env var names are identical across all plugin variants.
	// "pushWithBuildx" is used as the representative action ID.

	// BuildAndPushDockerRegistry (v1: buildAndPushToDocker)
	// Step group: pushWithBuildx / pushWithDocker / pushWithKaniko (conditional)
	StepTypeBuildAndPushDockerRegistry: {
		// connectorRef -> COMPLEX (resolves to USERNAME, PASSWORD, REGISTRY via connector)
		{"repo", "steps.pushWithBuildx.spec.with.REPO"},
		{"tags", "steps.pushWithBuildx.spec.with.TAGS"},
		// caching -> NO_ENV_MAPPING (controls which plugin variant is selected)
		// baseImageConnectorRefs -> COMPLEX (resolves to BASE_IMAGE_REGISTRY, BASE_IMAGE_USERNAME, BASE_IMAGE_PASSWORD via connector)
		// envVariables -> NO_ENV_MAPPING (merged into plugin env, not a single env var)
		{"dockerfile", "steps.pushWithBuildx.spec.with.DOCKERFILE"},
		{"context", "steps.pushWithBuildx.spec.with.CONTEXT"},
		{"labels", "steps.pushWithBuildx.spec.with.CUSTOM_LABELS"},
		{"buildArgs", "steps.pushWithBuildx.spec.with.BUILD_ARGS"},
		{"target", "steps.pushWithBuildx.spec.with.TARGET"},
		{"runAsUser", "spec.container.user"},
		{"resources.limits.cpu", "spec.container.cpu"},
		{"resources.limits.memory", "spec.container.memory"},
		// optimize -> NO_TEMPLATE_MAPPING
		// privileged -> NO_TEMPLATE_MAPPING
		// remoteCacheRepo -> NO_TEMPLATE_MAPPING
		// reports -> NO_TEMPLATE_MAPPING
	},

	// BuildAndPushECR (v1: buildAndPushToECR)
	// Step group: pushWithBuildx / pushWithECR / pushWithKaniko (conditional)
	StepTypeBuildAndPushECR: {
		// connectorRef -> COMPLEX (resolves to ACCESS_KEY, SECRET_KEY, ASSUME_ROLE, EXTERNAL_ID, OIDC_TOKEN_ID via connector)
		{"region", "steps.pushWithBuildx.spec.with.REGION"},
		// account + region -> COMPLEX (constructed as "<account>.dkr.ecr.<region>.amazonaws.com" for REGISTRY)
		{"imageName", "steps.pushWithBuildx.spec.with.REPO"},
		{"tags", "steps.pushWithBuildx.spec.with.TAGS"},
		// caching -> NO_ENV_MAPPING (controls which plugin variant is selected)
		// envVariables -> NO_ENV_MAPPING
		{"labels", "steps.pushWithBuildx.spec.with.CUSTOM_LABELS"},
		{"buildArgs", "steps.pushWithBuildx.spec.with.BUILD_ARGS"},
		// baseImageConnectorRefs -> COMPLEX (resolves via connector)
		{"dockerfile", "steps.pushWithBuildx.spec.with.DOCKERFILE"},
		{"context", "steps.pushWithBuildx.spec.with.CONTEXT"},
		{"target", "steps.pushWithBuildx.spec.with.TARGET"},
		// account -> NO_TEMPLATE_MAPPING (consumed to build registry URL)
		// runAsUser -> NO_TEMPLATE_MAPPING
	},

	// BuildAndPushGAR (v1: buildAndPushToGAR)
	// Step group: pushWithBuildx / pushWithGAR / pushWithKaniko (conditional)
	StepTypeBuildAndPushGAR: {
		// connectorRef -> COMPLEX (resolves to JSON_KEY, OIDC_TOKEN_ID, PROJECT_NUMBER, POOL_ID, PROVIDER_ID, SERVICE_ACCOUNT_EMAIL via connector)
		// host + projectID -> COMPLEX (constructed as "<host>/<projectID>" for REGISTRY)
		{"imageName", "steps.pushWithBuildx.spec.with.REPO"},
		{"tags", "steps.pushWithBuildx.spec.with.TAGS"},
		// caching -> NO_ENV_MAPPING (controls which plugin variant is selected)
		// envVariables -> NO_ENV_MAPPING
		{"labels", "steps.pushWithBuildx.spec.with.CUSTOM_LABELS"},
		{"buildArgs", "steps.pushWithBuildx.spec.with.BUILD_ARGS"},
		// baseImageConnectorRefs -> COMPLEX (resolves via connector)
		{"dockerfile", "steps.pushWithBuildx.spec.with.DOCKERFILE"},
		{"context", "steps.pushWithBuildx.spec.with.CONTEXT"},
		{"target", "steps.pushWithBuildx.spec.with.TARGET"},
		// host -> NO_TEMPLATE_MAPPING (consumed to build registry URL)
		// projectID -> NO_TEMPLATE_MAPPING (consumed to build registry URL)
		// runAsUser -> NO_TEMPLATE_MAPPING
	},

	// BuildAndPushACR (v1: buildAndPushToACR)
	// Step group: pushWithBuildx / pushWithACR / pushWithKaniko (conditional)
	StepTypeBuildAndPushACR: {
		// connectorRef -> COMPLEX (resolves to CLIENT_ID, TENANT_ID, CLIENT_SECRET, CLIENT_CERTIFICATE, OIDC_TOKEN_ID via connector)
		{"registry", "steps.pushWithBuildx.spec.with.REGISTRY"},
		{"imageName", "steps.pushWithBuildx.spec.with.REPO"},
		{"tags", "steps.pushWithBuildx.spec.with.TAGS"},
		// caching -> NO_ENV_MAPPING (controls which plugin variant is selected)
		// envVariables -> NO_ENV_MAPPING
		{"labels", "steps.pushWithBuildx.spec.with.CUSTOM_LABELS"},
		{"buildArgs", "steps.pushWithBuildx.spec.with.BUILD_ARGS"},
		// baseImageConnectorRefs -> COMPLEX (resolves via connector)
		{"dockerfile", "steps.pushWithBuildx.spec.with.DOCKERFILE"},
		{"context", "steps.pushWithBuildx.spec.with.CONTEXT"},
		{"target", "steps.pushWithBuildx.spec.with.TARGET"},
		{"subscriptionId", "steps.pushWithBuildx.spec.with.SUBSCRIPTION_ID"},
		// runAsUser -> NO_TEMPLATE_MAPPING
	},
}

// inside step.output
//
// Rule format: {v0 outcome field, v1 path relative to the step}.
// For template-backed steps the v1 path is `steps.<actionId>.output.outputVariables.<PLUGIN_KEY>`,
// where <actionId> is the wrapper template's internal sub-step id.
//
// DERIVED outputs (e.g. podIps from HARNESS_INSTANCES[*].PodIP) and
// NO_PLUGIN_OUTPUT fields are intentionally omitted — they cannot be
// expressed as a 1:1 path rewrite.
//
// Source: template-library/v0_to_v1_template_resolved_output_mappings.txt.
var StepOutputConversionRules = map[string][]ConversionRule{
 
	// ============================================================
	// HTTP (run plugin)
	// ============================================================
	StepTypeHTTP: {
		{"httpUrl", "steps.httpStep.output.outputVariables.PLUGIN_HTTP_URL"},
		{"httpMethod", "steps.httpStep.output.outputVariables.PLUGIN_HTTP_METHOD"},
		{"httpResponseCode", "steps.httpStep.output.outputVariables.PLUGIN_HTTP_RESPONSE_CODE"},
		{"httpResponseBody", "steps.httpStep.output.outputVariables.PLUGIN_HTTP_RESPONSE_BODY_BYTES"},
		{"status", "steps.httpStep.output.outputVariables.PLUGIN_EXECUTION_STATUS"},
		{"responseHeaders", "steps.httpStep.output.outputVariables.PLUGIN_RESPONSE_HEADERS"},
	},
 
	// ============================================================
	// K8S STEPS
	// ============================================================
 
	// K8sRollingDeploy (sub-steps: k8sRollingPrepareAction, k8sApplyAction, k8sSteadyStateCheckAction)
	StepTypeK8sRollingDeploy: {
		{"releaseName", "steps.k8sApplyAction.output.outputVariables.PLUGIN_RELEASE_NAME"},
		{"releaseNumber", "steps.k8sRollingPrepareAction.output.outputVariables.PLUGIN_RELEASE_NUMBER"},
		// also: steps.k8sApplyAction.output.outputVariables.PLUGIN_RELEASE_NUMBER
		{"prunedResourceIds", "steps.k8sRollingPrepareAction.output.outputVariables.PLUGIN_RESOURCES_FOR_PRUNING"},
		// COMPLEX: plugin emits []string of resource refs, not KubernetesResourceId structs
		// podIps -> DERIVED from steps.k8sSteadyStateCheckAction.output.outputVariables.HARNESS_INSTANCES[*].PodIP
		// manifest -> NO_PLUGIN_OUTPUT
	},
 
	// K8sRollingRollback (sub-steps: k8sRollingRollbackAction)
	StepTypeK8sRollingRollback: {
		{"recreatedResourceIds", "steps.k8sRollingRollbackAction.output.outputVariables.PLUGIN_RECREATED_RESOURCES"},
		// also: PLUGIN_RECREATED_WORKLOADS for workload subset
		// COMPLEX: plugin emits []string, not KubernetesResourceId
		// podIps -> DERIVED from steps.k8sRollingRollbackAction.output.outputVariables.HARNESS_INSTANCES[*].PodIP
	},
 
	// K8sCanaryDeploy (sub-steps: k8sCanaryPrepareAction, k8sApplyAction, k8sApplySteadyStateCheckAction, (cond) k8sTrafficRoutingAction)
	StepTypeK8sCanaryDeploy: {
		{"releaseName", "steps.k8sApplyAction.output.outputVariables.PLUGIN_RELEASE_NAME"},
		{"targetInstances", "steps.k8sCanaryPrepareAction.output.outputVariables.PLUGIN_TARGET_INSTANCES"},
		{"releaseNumber", "steps.k8sCanaryPrepareAction.output.outputVariables.PLUGIN_RELEASE_NUMBER"},
		{"canaryWorkload", "steps.k8sCanaryPrepareAction.output.outputVariables.PLUGIN_CANARY_WORKLOADS"},
		// COMPLEX: plugin emits []string; v0 was a single workload string
		// canaryWorkloadDeployed -> NO_PLUGIN_OUTPUT
		// podIps -> DERIVED from steps.k8sApplySteadyStateCheckAction.output.outputVariables.HARNESS_INSTANCES[*].PodIP
		// manifest -> NO_PLUGIN_OUTPUT
	},
 
	// K8sCanaryDelete (EMPTY_OUTCOME — no fields to map)
 
	// K8sBlueGreenDeploy (sub-steps: k8sBlueGreenPrepareAction, k8sApplyAction, k8sApplySteadyStateCheckAction, (cond) k8sTrafficRoutingAction)
	StepTypeK8sBlueGreenDeploy: {
		{"releaseNumber", "steps.k8sBlueGreenPrepareAction.output.outputVariables.PLUGIN_RELEASE_NUMBER"},
		{"releaseName", "steps.k8sApplyAction.output.outputVariables.PLUGIN_RELEASE_NAME"},
		{"primaryServiceName", "steps.k8sBlueGreenPrepareAction.output.outputVariables.PLUGIN_STABLE_SERVICE"},
		{"stageServiceName", "steps.k8sBlueGreenPrepareAction.output.outputVariables.PLUGIN_STAGE_SERVICE"},
		// stageColor / primaryColor -> NO_PLUGIN_OUTPUT (computed by swap step; see K8sBGSwapServices)
		{"stageDeploymentSkipped", "steps.k8sBlueGreenPrepareAction.output.outputVariables.PLUGIN_SKIPPED"},
		// podIps -> DERIVED from steps.k8sApplySteadyStateCheckAction.output.outputVariables.HARNESS_INSTANCES[*].PodIP
		{"prunedResourceIds", "steps.k8sBlueGreenPrepareAction.output.outputVariables.PLUGIN_RESOURCES_FOR_PRUNING"},
		// COMPLEX: plugin emits []string
		// manifest -> NO_PLUGIN_OUTPUT
	},
 
	// K8sBGSwapServices (sub-steps: k8sBlueGreenSwapServicesSelectorsAction)
	// v0 outcome was empty; v1 plugin publishes new outputs.
	StepTypeK8sBGSwapServices: {
		{"primaryColor", "steps.k8sBlueGreenSwapServicesSelectorsAction.output.outputVariables.PLUGIN_PRIMARY_COLOR"},
		{"stageColor", "steps.k8sBlueGreenSwapServicesSelectorsAction.output.outputVariables.PLUGIN_STAGE_COLOR"},
		{"stableService", "steps.k8sBlueGreenSwapServicesSelectorsAction.output.outputVariables.PLUGIN_STABLE_SERVICE"},
		{"stageService", "steps.k8sBlueGreenSwapServicesSelectorsAction.output.outputVariables.PLUGIN_STAGE_SERVICE"},
	},
 
	// K8sDryRun (sub-steps: k8sDryRunAction — k8s-apply with only_dry_run=true)
	StepTypeK8sDryRun: {
		{"manifestDryRun", "steps.k8sDryRunAction.output.outputVariables.PLUGIN_MANIFEST_OUTPUT_FILE_PATH"},
		// COMPLEX: v0 returned the rendered manifest as a string;
		// v1 returns the path to the rendered file
		// manifest -> NO_PLUGIN_OUTPUT
	},
 
	// K8sDiff -> NO_PLUGIN_OUTPUT for all v0 outcome fields (manifestDiff, exitValue)
	// K8sPatch -> NO_PLUGIN_OUTPUT (releaseName, releaseNumber)
	// K8sDelete -> EMPTY_OUTCOME
	// K8sRollout -> EMPTY_OUTCOME
	// K8sCanaryDelete -> EMPTY_OUTCOME
 
	// K8sApply (sub-steps: k8sApplyAction, k8sApplySteadyStateCheckAction)
	// v0 returns Status only; v1 plugins publish new outputs not previously available.
	StepTypeK8sApply: {
		{"manifestOutputFilePath", "steps.k8sApplyAction.output.outputVariables.PLUGIN_MANIFEST_OUTPUT_FILE_PATH"},
		{"releaseNumber", "steps.k8sApplyAction.output.outputVariables.PLUGIN_RELEASE_NUMBER"},
		{"releaseName", "steps.k8sApplyAction.output.outputVariables.PLUGIN_RELEASE_NAME"},
		{"managedWorkloads", "steps.k8sApplyAction.output.outputVariables.PLUGIN_MANAGED_WORKLOADS"},
		{"customWorkloads", "steps.k8sApplyAction.output.outputVariables.PLUGIN_CUSTOM_WORKLOADS"},
		{"resources", "steps.k8sApplyAction.output.outputVariables.PLUGIN_RESOURCES"},
		{"existingPods", "steps.k8sApplyAction.output.outputVariables.PLUGIN_EXISTING_PODS"},
		// COMPLEX: []K8sPod
		{"existingPodNames", "steps.k8sApplyAction.output.outputVariables.PLUGIN_EXISTING_POD_NAMES"},
		{"isOpenshift", "steps.k8sApplyAction.output.outputVariables.HARNESS_IS_OPENSHIFT"},
		{"pods", "steps.k8sApplySteadyStateCheckAction.output.outputVariables.PLUGIN_PODS"},
		// COMPLEX: []K8sPod
		{"instances", "steps.k8sApplySteadyStateCheckAction.output.outputVariables.HARNESS_INSTANCES"},
		// COMPLEX: []Instance
		{"skipped", "steps.k8sApplySteadyStateCheckAction.output.outputVariables.PLUGIN_SKIPPED"},
	},
 
	// K8sScale (sub-steps: k8sScaleAction, (cond) k8sScaleSteadyStateCheckAction)
	// v0 returns Status only; v1 plugins publish new outputs.
	StepTypeK8sScale: {
		{"pods", "steps.k8sScaleAction.output.outputVariables.PLUGIN_PODS"},
		// COMPLEX: []K8sPod
		{"instances", "steps.k8sScaleAction.output.outputVariables.HARNESS_INSTANCES"},
		// COMPLEX: []Instance — also published by k8sScaleSteadyStateCheckAction
	},
 
	// K8sBlueGreenStageScaleDown (no v0 outcome; v1 publishes scaled resources)
	StepTypeK8sBlueGreenStageScaleDown: {
		{"scaledResources", "steps.k8sBlueGreenStageScaleDownAction.output.outputVariables.PLUGIN_SCALED_RESOURCES"},
	},
 
	// K8sTrafficRouting (sub-steps: k8sTrafficShiftAction)
	// v0 has no outcome; v1 plugin publishes new outputs (config_type=new variant shown).
	StepTypeK8sTrafficRouting: {
		{"manifestOutputFilePath", "steps.k8sTrafficShiftAction.output.outputVariables.PLUGIN_MANIFEST_OUTPUT_FILE_PATH"},
		{"resources", "steps.k8sTrafficShiftAction.output.outputVariables.PLUGIN_RESOURCES"},
		{"releaseNumber", "steps.k8sTrafficShiftAction.output.outputVariables.PLUGIN_RELEASE_NUMBER"},
		{"resourcesForPruning", "steps.k8sTrafficShiftAction.output.outputVariables.PLUGIN_RESOURCES_FOR_PRUNING"},
	},
 
	// ============================================================
	// HELM STEPS (Native Helm)
	// ============================================================
 
	// HelmDeploy / HelmDeployBasic (sub-steps: helmBasicDeployAction, (cond) helmTestAction)
	StepTypeHelmDeploy: {
		{"releaseName", "steps.helmBasicDeployAction.output.outputVariables.PLUGIN_RELEASE_NAME"},
		{"prevReleaseVersion", "steps.helmBasicDeployAction.output.outputVariables.PLUGIN_PREVIOUS_RELEASE_REVISION"},
		{"newReleaseVersion", "steps.helmBasicDeployAction.output.outputVariables.PLUGIN_NEW_RELEASE_REVISION"},
		// hasInstallUpgradeStarted -> NO_PLUGIN_OUTPUT (CD-NG flag)
		// chartTestSucceeded -> NO_PLUGIN_OUTPUT (helm-test only sets step status)
		// podIps -> DERIVED from steps.helmBasicDeployAction.output.outputVariables.PLUGIN_NEW_PODS[*].PodIP
		// deltaForNewReleaseVersion -> DERIVED (PLUGIN_NEW_RELEASE_REVISION - PLUGIN_PREVIOUS_RELEASE_REVISION)
		{"containerInfoList", "steps.helmBasicDeployAction.output.outputVariables.PLUGIN_CONTAINER_INFO_LIST"},
		// COMPLEX: []*ContainerInfo
	},
 
	// HelmRollback (sub-steps: helmRollbackAction, (cond) helmTestAction)
	StepTypeHelmRollback: {
		{"releaseName", "steps.helmRollbackAction.output.outputVariables.PLUGIN_RELEASE_NAME"},
		{"newReleaseVersion", "steps.helmRollbackAction.output.outputVariables.PLUGIN_NEW_RELEASE_REVISION"},
		{"rollbackVersion", "steps.helmRollbackAction.output.outputVariables.PLUGIN_ROLLBACK_RELEASE_REVISION"},
		// chartTestSucceeded -> NO_PLUGIN_OUTPUT
		{"containerInfoList", "steps.helmRollbackAction.output.outputVariables.PLUGIN_CONTAINER_INFO_LIST"},
	},
 
	// HelmBlueGreenDeploy / HelmBGDeploy (sub-steps: helmBluegreenDeployAction, (cond) helmTestAction)
	StepTypeHelmBGDeploy: {
		{"releaseNumber", "steps.helmBluegreenDeployAction.output.outputVariables.PLUGIN_NEW_RELEASE_REVISION"},
		{"releaseName", "steps.helmBluegreenDeployAction.output.outputVariables.PLUGIN_RELEASE_NAME"},
		{"primaryServiceNames", "steps.helmBluegreenDeployAction.output.outputVariables.PLUGIN_STABLE_SERVICE"},
		// COMPLEX: plugin emits a single string; v0 was List<String>
		{"stageServiceNames", "steps.helmBluegreenDeployAction.output.outputVariables.PLUGIN_STAGE_SERVICE"},
		// COMPLEX: same — single string vs list
		{"stageColor", "steps.helmBluegreenDeployAction.output.outputVariables.PLUGIN_STAGE_COLOR"},
		{"primaryColor", "steps.helmBluegreenDeployAction.output.outputVariables.PLUGIN_STABLE_COLOR"},
		// hasInstallUpgradeStarted / chartTestSucceeded -> NO_PLUGIN_OUTPUT
		// podIps -> DERIVED from steps.helmBluegreenDeployAction.output.outputVariables.PLUGIN_NEW_PODS[*].PodIP
		// primaryResources / stageResources -> NO_PLUGIN_OUTPUT
		//   (helm-deploy emits PLUGIN_WORKLOADS as []string, not structured KubernetesResource)
	},
 
	// HelmBlueGreenSwapStep / HelmBGSwapServices -> EMPTY_OUTCOME
 
	// HelmCanaryDeploy (sub-steps: helmCanaryDeployAction (uses helmDeployAction), (cond) helmTestAction)
	StepTypeHelmCanaryDeploy: {
		{"releaseName", "steps.helmCanaryDeployAction.output.outputVariables.PLUGIN_CANARY_RELEASE_NAME"},
		// canary-strategy variant of PLUGIN_RELEASE_NAME
		// hasInstallUpgradeStarted / canaryWorkloadDeployed / chartTestSucceeded -> NO_PLUGIN_OUTPUT
		{"canaryWorkload", "steps.helmCanaryDeployAction.output.outputVariables.PLUGIN_WORKLOADS"},
		// COMPLEX: plugin emits []string, v0 was a single string
		{"previousReleaseVersion", "steps.helmCanaryDeployAction.output.outputVariables.PLUGIN_PREVIOUS_RELEASE_REVISION"},
		// podIps -> DERIVED from steps.helmCanaryDeployAction.output.outputVariables.PLUGIN_NEW_PODS[*].PodIP
	},
 
	// HelmCanaryDelete -> EMPTY_OUTCOME
	// HelmDelete -> EMPTY_OUTCOME
}