package convertexpressions

import (
	"testing"
)

func TestTrieRules_PipelineLevel(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		context  *ConversionContext
		expected string
	}{
		{
			name:     "stage variables",
			path:     "stage.variables.var1",
			context:  nil,
			expected: "stage.variables.var1",
		},
		{
			name:     "stage identifier FQN",
			path:     "pipeline.stages.build.identifier",
			context:  nil,
			expected: "pipeline.stages.build.id",
		},
		{
			name:     "stage env identifier FQN",
			path:     "pipeline.stages.build.spec.env.identifier",
			context:  nil,
			expected: "pipeline.stages.build.env.id",
		},
		{
			name:     "stage env group identifier FQN",
			path:     "pipeline.stages.build.spec.env.envGroupRef",
			context:  nil,
			expected: "pipeline.stages.build.env.group.id",
		},
		{
			name:     "stage env group name FQN",
			path:     "pipeline.stages.build.spec.env.envGroupName",
			context:  nil,
			expected: "pipeline.stages.build.env.group.name",
		},
		{
			name:     "stage env identifier relative",
			path:     "stage.spec.env.identifier",
			context:  nil,
			expected: "stage.env.id",
		},
		{
			name:     "stage env group identifier relative",
			path:     "spec.env.envGroupRef",
			context:  nil,
			expected: "env.group.id",
		},
		{
			name:     "stage env group name relative",
			path:     "spec.env.envGroupName",
			context:  nil,
			expected: "env.group.name",
		},
		{
			name:     "stage env identifier direct",
			path:     "env.identifier",
			context:  nil,
			expected: "env.id",
		},
		{
			name:     "stage env group identifier direct",
			path:     "env.envGroupRef",
			context:  nil,
			expected: "env.group.id",
		},
		{
			name:     "stage env group name direct",
			path:     "env.envGroupName",
			context:  nil,
			expected: "env.group.name",
		},
		{
			name:     "spec.execution.steps removal FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.output.outputVariables.var1",
			context:  nil,
			expected: "pipeline.stages.build.steps.step1.output.outputVariables.var1",
		},
		{
			name:     "is not resolved",
			path:     "expression.isResolved(<+pipeline.variables.var1>)",
			context:  nil,
			expected: "expression.isResolved(<+pipeline.variables.var1>)",
		},
		{
			name:     "nested and function",
			path:     "<+pipeline.variables.var1>.some_func(\"param\")",
			context:  nil,
			expected: "<+pipeline.variables.var1>.some_func(\"param\")",
		},
	}

	trie := buildPipelineTrie()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := trie.Match(tt.path, tt.context)
			if result != tt.expected {
				t.Errorf("Match() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTrieRules_StepLevel(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		context  *ConversionContext
		expected string
	}{
		// General step rules (no context needed)
		{
			name:     "step identifier FQN - no context",
			path:     "pipeline.stages.build.spec.execution.steps.step1.identifier",
			context:  nil,
			expected: "pipeline.stages.build.steps.step1.id",
		},
		{
			name:     "step identifier FQN - with context",
			path:     "pipeline.stages.build.spec.execution.steps.step1.identifier",
			context:  &ConversionContext{StepType: StepTypeRun},
			expected: "pipeline.stages.build.steps.step1.id",
		},
		{
			name:     "output variables FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.output.outputVariables.var1",
			context:  &ConversionContext{StepType: StepTypeRun},
			expected: "pipeline.stages.build.steps.step1.output.outputVariables.var1",
		},
		{
			name:     "output variables relative",
			path:     "spec.execution.steps.step1.output.outputVariables.var1",
			context:  &ConversionContext{StepType: StepTypeRun},
			expected: "steps.step1.output.outputVariables.var1",
		},
		{
			name:     "output variables spec name",
			path:     "spec.execution.steps.step1.spec.outputVariables[1].name",
			context:  &ConversionContext{StepType: StepTypeRun},
			expected: "steps.step1.spec.output[1].alias",
		},
		{
			name:     "output variables spec",
			path:     "spec.execution.steps.step1.spec.outputVariables",
			context:  &ConversionContext{StepType: StepTypeRun},
			expected: "steps.step1.spec.output",
		},
		{
			name:     "output variables relative with function",
			path:     "spec.execution.steps.step1.output.outputVariables.var1.some_func(\"param\")",
			context:  &ConversionContext{StepType: StepTypeRun},
			expected: "steps.step1.output.outputVariables.var1.some_func(\"param\")",
		},
		// Failure strategies
		{
			name:     "failure strategy in step - errors",
			path:     "pipeline.stages.build.spec.execution.steps.step1.failureStrategies[0].onFailure.errors",
			context:  &ConversionContext{StepType: StepTypeRun},
			expected: "pipeline.stages.build.steps.step1.onFailure[0].errors",
		},
		{
			name:     "failure strategy in step - retry count",
			path:     "pipeline.stages.build.spec.execution.steps.step1.failureStrategies[0].onFailure.action.specConfig.retryCount",
			context:  &ConversionContext{StepType: StepTypeRun},
			expected: "pipeline.stages.build.steps.step1.onFailure[0].action.retry.attempts",
		},
		// Step-specific rules with context - SaveCacheGCS
		{
			name:     "save cache to gcs FQN",
			path:     "pipeline.stages.build.spec.execution.steps.STEPID.spec.bucket",
			context:  &ConversionContext{StepType: StepTypeSaveCacheGCS},
			expected: "pipeline.stages.build.steps.STEPID.steps.saveCacheGCS.spec.with.BUCKET",
		},
		{
			name:     "save cache to gcs relative",
			path:     "execution.steps.STEPID.spec.bucket",
			context:  &ConversionContext{StepType: StepTypeSaveCacheGCS},
			expected: "steps.STEPID.steps.saveCacheGCS.spec.with.BUCKET",
		},
		// Step-specific rules with context - RestoreCacheS3
		{
			name:     "restore cache from s3 FQN",
			path:     "pipeline.stages.build.spec.execution.steps.STEPID.spec.bucket",
			context:  &ConversionContext{StepType: StepTypeRestoreCacheS3},
			expected: "pipeline.stages.build.steps.STEPID.steps.restoreCacheS3.spec.with.BUCKET",
		},
		{
			name:     "restore cache from s3 relative",
			path:     "execution.steps.STEPID.spec.bucket",
			context:  &ConversionContext{StepType: StepTypeRestoreCacheS3},
			expected: "steps.STEPID.steps.restoreCacheS3.spec.with.BUCKET",
		},
		{
			name:     "restore cache from s3 inside step group",
			path:     "pipeline.stages.build.spec.execution.steps.stepGroupID.steps.STEPID.spec.bucket",
			context:  &ConversionContext{StepType: StepTypeRestoreCacheS3},
			expected: "pipeline.stages.build.steps.stepGroupID.steps.STEPID.steps.restoreCacheS3.spec.with.BUCKET",
		},
		{
			name:     "restore cache from s3 relative stepgroup",
			path:     "stepGroup.steps.STEPID.spec.bucket",
			context:  &ConversionContext{StepType: StepTypeRestoreCacheS3},
			expected: "group.steps.STEPID.steps.restoreCacheS3.spec.with.BUCKET",
		},
		// Step-specific rules with context - Run
		{
			name:     "run step command FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.command",
			context:  &ConversionContext{StepType: StepTypeRun},
			expected: "pipeline.stages.build.steps.step1.spec.script",
		},
		{
			name:     "run step image FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.image",
			context:  &ConversionContext{StepType: StepTypeRun},
			expected: "pipeline.stages.build.steps.step1.spec.container.image",
		},
		{
			name:     "run step envVariables FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.envVariables",
			context:  &ConversionContext{StepType: StepTypeRun},
			expected: "pipeline.stages.build.steps.step1.spec.env",
		},
		// Step-specific rules with context - HTTP
		{
			name:     "http step url FQN",
			path:     "pipeline.stages.build.spec.execution.steps.http1.spec.url",
			context:  &ConversionContext{StepType: StepTypeHTTP},
			expected: "pipeline.stages.build.steps.http1.spec.env.PLUGIN_URL",
		},
		{
			name:     "http step output httpResponseCode FQN",
			path:     "pipeline.stages.build.spec.execution.steps.http1.output.httpResponseCode",
			context:  &ConversionContext{StepType: StepTypeHTTP},
			expected: "pipeline.stages.build.steps.http1.steps.httpStep.output.outputVariables.PLUGIN_HTTP_RESPONSE_CODE",
		},
		{
			name:     "http step output status FQN",
			path:     "pipeline.stages.build.spec.execution.steps.http1.output.status",
			context:  &ConversionContext{StepType: StepTypeHTTP},
			expected: "pipeline.stages.build.steps.http1.steps.httpStep.output.outputVariables.PLUGIN_EXECUTION_STATUS",
		},
	}

	trie := buildPipelineTrie()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := trie.Match(tt.path, tt.context)
			if result != tt.expected {
				t.Errorf("Match() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTrieRules_K8sSteps(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		context  *ConversionContext
		expected string
	}{
		// K8sRollingDeploy
		{
			name:     "K8sRollingDeploy skipDryRun FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.skipDryRun",
			context:  &ConversionContext{StepType: StepTypeK8sRollingDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sApplyAction.spec.env.PLUGIN_SKIP_DRY_RUN",
		},
		{
			name:     "K8sRollingDeploy pruningEnabled FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.pruningEnabled",
			context:  &ConversionContext{StepType: StepTypeK8sRollingDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sRollingPrepareAction.spec.env.PLUGIN_PRUNING_ENABLED",
		},
		{
			name:     "K8sRollingDeploy flags FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.flags",
			context:  &ConversionContext{StepType: StepTypeK8sRollingDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sApplyAction.spec.env.PLUGIN_FLAGS_JSON",
		},
		// K8sRollingRollback
		{
			name:     "K8sRollingRollback pruningEnabled FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.pruningEnabled",
			context:  &ConversionContext{StepType: StepTypeK8sRollingRollback},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sRollingRollbackAction.spec.env.PLUGIN_PRUNING",
		},
		{
			name:     "K8sRollingRollback flags FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.flags",
			context:  &ConversionContext{StepType: StepTypeK8sRollingRollback},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sRollingRollbackAction.spec.env.PLUGIN_FLAGS_JSON",
		},
		// K8sApply
		{
			name:     "K8sApply filePaths FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.filePaths",
			context:  &ConversionContext{StepType: StepTypeK8sApply},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sApplyAction.spec.env.PLUGIN_MANIFEST_PATH",
		},
		{
			name:     "K8sApply skipDryRun FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.skipDryRun",
			context:  &ConversionContext{StepType: StepTypeK8sApply},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sApplyAction.spec.env.PLUGIN_SKIP_DRY_RUN",
		},
		{
			name:     "K8sApply flags FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.flags",
			context:  &ConversionContext{StepType: StepTypeK8sApply},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sApplyAction.spec.env.PLUGIN_FLAGS_JSON",
		},
		// K8sBGSwapServices
		{
			name:     "K8sBGSwapServices stable_service FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.stable_service",
			context:  &ConversionContext{StepType: StepTypeK8sBGSwapServices},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sBlueGreenSwapServicesSelectorsAction.spec.env.PLUGIN_STABLE_SERVICE",
		},
		{
			name:     "K8sBGSwapServices stage_service FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.stage_service",
			context:  &ConversionContext{StepType: StepTypeK8sBGSwapServices},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sBlueGreenSwapServicesSelectorsAction.spec.env.PLUGIN_STAGE_SERVICE",
		},
		{
			name:     "K8sBGSwapServices is_openshift FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.is_openshift",
			context:  &ConversionContext{StepType: StepTypeK8sBGSwapServices},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sBlueGreenSwapServicesSelectorsAction.spec.env.HARNESS_IS_OPENSHIFT",
		},
		// K8sCanaryDelete
		{
			name:     "K8sCanaryDelete resources FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.resources",
			context:  &ConversionContext{StepType: StepTypeK8sCanaryDelete},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sCanaryDeleteAction.spec.env.PLUGIN_RESOURCES",
		},
		{
			name:     "K8sCanaryDelete is_openshift FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.is_openshift",
			context:  &ConversionContext{StepType: StepTypeK8sCanaryDelete},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sCanaryDeleteAction.spec.env.HARNESS_IS_OPENSHIFT",
		},
		// K8sRollout
		{
			name:     "K8sRollout command FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.command",
			context:  &ConversionContext{StepType: StepTypeK8sRollout},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sRolloutStep.spec.env.PLUGIN_COMMAND",
		},
		{
			name:     "K8sRollout resources.spec.resourceNames FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.resources.spec.resourceNames",
			context:  &ConversionContext{StepType: StepTypeK8sRollout},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sRolloutStep.spec.env.PLUGIN_RESOURCES",
		},
		{
			name:     "K8sRollout resources.spec.manifestPaths FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.resources.spec.manifestPaths",
			context:  &ConversionContext{StepType: StepTypeK8sRollout},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sRolloutStep.spec.env.PLUGIN_MANIFESTS",
		},
		{
			name:     "K8sRollout flags FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.flags",
			context:  &ConversionContext{StepType: StepTypeK8sRollout},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sRolloutStep.spec.env.PLUGIN_FLAGS_JSON",
		},
		// K8sScale
		{
			name:     "K8sScale instanceSelection.type FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.instanceSelection.type",
			context:  &ConversionContext{StepType: StepTypeK8sScale},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sScaleAction.spec.env.PLUGIN_INSTANCES_UNIT_TYPE",
		},
		{
			name:     "K8sScale instanceSelection.spec.count FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.instanceSelection.spec.count",
			context:  &ConversionContext{StepType: StepTypeK8sScale},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sScaleAction.spec.env.PLUGIN_INSTANCES",
		},
		{
			name:     "K8sScale workload FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.workload",
			context:  &ConversionContext{StepType: StepTypeK8sScale},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sScaleAction.spec.env.PLUGIN_WORKLOAD",
		},
		// K8sDryRun
		{
			name:     "K8sDryRun encryptYamlOutput FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.encryptYamlOutput",
			context:  &ConversionContext{StepType: StepTypeK8sDryRun},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sDryRunAction.spec.env.PLUGIN_ENCRYPT_YAML_OUTPUT",
		},
		// K8sDelete
		{
			name:     "K8sDelete deleteResources.spec.resourceNames FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.deleteResources.spec.resourceNames",
			context:  &ConversionContext{StepType: StepTypeK8sDelete},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sDeleteAction.spec.env.PLUGIN_RESOURCES",
		},
		{
			name:     "K8sDelete deleteResources.spec.manifestPaths FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.deleteResources.spec.manifestPaths",
			context:  &ConversionContext{StepType: StepTypeK8sDelete},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sDeleteAction.spec.env.PLUGIN_MANIFESTS",
		},
		{
			name:     "K8sDelete deleteResources.spec.deleteNamespace FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.deleteResources.spec.deleteNamespace",
			context:  &ConversionContext{StepType: StepTypeK8sDelete},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sDeleteAction.spec.env.PLUGIN_INCLUDE_NAMESPACES",
		},
		{
			name:     "K8sDelete flags FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.flags",
			context:  &ConversionContext{StepType: StepTypeK8sDelete},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sDeleteAction.spec.env.PLUGIN_FLAGS_JSON",
		},
		// K8sTrafficRouting
		{
			name:     "K8sTrafficRouting trafficRouting.provider FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.trafficRouting.provider",
			context:  &ConversionContext{StepType: StepTypeK8sTrafficRouting},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sTrafficShiftAction.spec.env.PLUGIN_PROVIDER",
		},
		{
			name:     "K8sTrafficRouting trafficRouting.spec.name FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.trafficRouting.spec.name",
			context:  &ConversionContext{StepType: StepTypeK8sTrafficRouting},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sTrafficShiftAction.spec.env.PLUGIN_RESOURCE_NAME",
		},
		{
			name:     "K8sTrafficRouting trafficRouting.spec.hosts FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.trafficRouting.spec.hosts",
			context:  &ConversionContext{StepType: StepTypeK8sTrafficRouting},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sTrafficShiftAction.spec.env.PLUGIN_HOSTNAMES",
		},
		{
			name:     "K8sTrafficRouting trafficRouting.spec.routes FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.trafficRouting.spec.routes",
			context:  &ConversionContext{StepType: StepTypeK8sTrafficRouting},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sTrafficShiftAction.spec.env.PLUGIN_ROUTES",
		},
		// K8sCanaryDeploy
		{
			name:     "K8sCanaryDeploy instanceSelection.type FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.instanceSelection.type",
			context:  &ConversionContext{StepType: StepTypeK8sCanaryDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sCanaryPrepareAction.spec.env.PLUGIN_INSTANCES_UNIT_TYPE",
		},
		{
			name:     "K8sCanaryDeploy skipDryRun FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.skipDryRun",
			context:  &ConversionContext{StepType: StepTypeK8sCanaryDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sApplyAction.spec.env.PLUGIN_SKIP_DRY_RUN",
		},
		{
			name:     "K8sCanaryDeploy trafficRouting.provider FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.trafficRouting.provider",
			context:  &ConversionContext{StepType: StepTypeK8sCanaryDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sTrafficRoutingAction.spec.env.PLUGIN_PROVIDER",
		},
		{
			name:     "K8sCanaryDeploy flags FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.flags",
			context:  &ConversionContext{StepType: StepTypeK8sCanaryDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sApplyAction.spec.env.PLUGIN_FLAGS_JSON",
		},
		// K8sBlueGreenDeploy
		{
			name:     "K8sBlueGreenDeploy skipDryRun FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.skipDryRun",
			context:  &ConversionContext{StepType: StepTypeK8sBlueGreenDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sApplyAction.spec.env.PLUGIN_SKIP_DRY_RUN",
		},
		{
			name:     "K8sBlueGreenDeploy pruningEnabled FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.pruningEnabled",
			context:  &ConversionContext{StepType: StepTypeK8sBlueGreenDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sApplyAction.spec.env.PLUGIN_RELEASE_PRUNING_ENABLED",
		},
		{
			name:     "K8sBlueGreenDeploy skipUnchangedManifest FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.skipUnchangedManifest",
			context:  &ConversionContext{StepType: StepTypeK8sBlueGreenDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sBlueGreenPrepareAction.spec.env.PLUGIN_SKIP_UNCHANGED_MANIFEST",
		},
		{
			name:     "K8sBlueGreenDeploy trafficRouting.spec.gateways FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.trafficRouting.spec.gateways",
			context:  &ConversionContext{StepType: StepTypeK8sBlueGreenDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sTrafficRoutingAction.spec.env.PLUGIN_GATEWAYS",
		},
		{
			name:     "K8sBlueGreenDeploy flags FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.flags",
			context:  &ConversionContext{StepType: StepTypeK8sBlueGreenDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sApplyAction.spec.env.PLUGIN_FLAGS_JSON",
		},
		// K8sPatch
		{
			name:     "K8sPatch workload FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.workload",
			context:  &ConversionContext{StepType: StepTypeK8sPatch},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sPatchAction.spec.env.PLUGIN_WORKLOAD",
		},
		{
			name:     "K8sPatch mergeStrategyType FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.mergeStrategyType",
			context:  &ConversionContext{StepType: StepTypeK8sPatch},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sPatchAction.spec.env.PLUGIN_MERGE_STRATEGY",
		},
		{
			name:     "K8sPatch source.spec.content FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.source.spec.content",
			context:  &ConversionContext{StepType: StepTypeK8sPatch},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sPatchAction.spec.env.PLUGIN_CONTENT",
		},
		// K8sRollingDeploy relative path
		{
			name:     "K8sRollingDeploy skipDryRun relative",
			path:     "execution.steps.step1.spec.skipDryRun",
			context:  &ConversionContext{StepType: StepTypeK8sRollingDeploy},
			expected: "steps.step1.steps.k8sApplyAction.spec.env.PLUGIN_SKIP_DRY_RUN",
		},
	}

	trie := buildPipelineTrie()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := trie.Match(tt.path, tt.context)
			if result != tt.expected {
				t.Errorf("Match() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTrieRules_HelmSteps(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		context  *ConversionContext
		expected string
	}{
		// HelmBGDeploy
		{
			name:     "HelmBGDeploy ignoreReleaseHistFailStatus FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.ignoreReleaseHistFailStatus",
			context:  &ConversionContext{StepType: StepTypeHelmBGDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBluegreenDeployAction.spec.env.PLUGIN_IGNORE_HISTORY_FAILURE",
		},
		{
			name:     "HelmBGDeploy skipSteadyStateCheck FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.skipSteadyStateCheck",
			context:  &ConversionContext{StepType: StepTypeHelmBGDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBluegreenDeployAction.spec.env.PLUGIN_SKIP_STEADY_STATE_CHECK",
		},
		{
			name:     "HelmBGDeploy environmentVariables FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.environmentVariables",
			context:  &ConversionContext{StepType: StepTypeHelmBGDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBluegreenDeployAction.spec.env.PLUGIN_ENV_VARS",
		},
		{
			name:     "HelmBGDeploy flags FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.flags",
			context:  &ConversionContext{StepType: StepTypeHelmBGDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBluegreenDeployAction.spec.env.PLUGIN_FLAGS",
		},
		// HelmBlueGreenSwapStep
		{
			name:     "HelmBlueGreenSwapStep flags FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.flags",
			context:  &ConversionContext{StepType: StepTypeHelmBlueGreenSwapStep},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBluegreenSwapAction.spec.env.PLUGIN_FLAGS",
		},
		// HelmCanaryDeploy
		{
			name:     "HelmCanaryDeploy ignoreReleaseHistFailStatus FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.ignoreReleaseHistFailStatus",
			context:  &ConversionContext{StepType: StepTypeHelmCanaryDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBasicDeployAction.spec.env.PLUGIN_IGNORE_HISTORY_FAILURE",
		},
		{
			name:     "HelmCanaryDeploy instanceSelection.type FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.instanceSelection.type",
			context:  &ConversionContext{StepType: StepTypeHelmCanaryDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBasicDeployAction.spec.env.PLUGIN_INSTANCES_UNIT_TYPE",
		},
		{
			name:     "HelmCanaryDeploy instanceSelection.spec.count FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.instanceSelection.spec.count",
			context:  &ConversionContext{StepType: StepTypeHelmCanaryDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBasicDeployAction.spec.env.PLUGIN_INSTANCES",
		},
		{
			name:     "HelmCanaryDeploy flags FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.flags",
			context:  &ConversionContext{StepType: StepTypeHelmCanaryDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBasicDeployAction.spec.env.PLUGIN_FLAGS",
		},
		// HelmDelete
		{
			name:     "HelmDelete releaseName FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.releaseName",
			context:  &ConversionContext{StepType: StepTypeHelmDelete},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmUninstallAction.spec.env.PLUGIN_RELEASE_NAME",
		},
		{
			name:     "HelmDelete dryRun FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.dryRun",
			context:  &ConversionContext{StepType: StepTypeHelmDelete},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmUninstallAction.spec.env.PLUGIN_DRY_RUN",
		},
		{
			name:     "HelmDelete commandFlags FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.commandFlags",
			context:  &ConversionContext{StepType: StepTypeHelmDelete},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmUninstallAction.spec.env.PLUGIN_FLAGS",
		},
		// HelmDeploy
		{
			name:     "HelmDeploy ignoreReleaseHistFailStatus FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.ignoreReleaseHistFailStatus",
			context:  &ConversionContext{StepType: StepTypeHelmDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBasicDeployAction.spec.env.PLUGIN_IGNORE_HISTORY_FAILURE",
		},
		{
			name:     "HelmDeploy environmentVariables FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.environmentVariables",
			context:  &ConversionContext{StepType: StepTypeHelmDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBasicDeployAction.spec.env.PLUGIN_ENV_VARS",
		},
		{
			name:     "HelmDeploy flags FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.flags",
			context:  &ConversionContext{StepType: StepTypeHelmDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBasicDeployAction.spec.env.PLUGIN_FLAGS",
		},
		// HelmRollback
		{
			name:     "HelmRollback skipSteadyStateCheck FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.skipSteadyStateCheck",
			context:  &ConversionContext{StepType: StepTypeHelmRollback},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmRollbackAction.spec.env.PLUGIN_SKIP_STEADY_STATE_CHECK",
		},
		{
			name:     "HelmRollback environmentVariables FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.environmentVariables",
			context:  &ConversionContext{StepType: StepTypeHelmRollback},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmRollbackAction.spec.env.PLUGIN_ENV_VARS",
		},
		{
			name:     "HelmRollback flags FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.flags",
			context:  &ConversionContext{StepType: StepTypeHelmRollback},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmRollbackAction.spec.env.PLUGIN_FLAGS",
		},
		// HelmBGDeploy relative path
		{
			name:     "HelmBGDeploy ignoreReleaseHistFailStatus relative",
			path:     "execution.steps.step1.spec.ignoreReleaseHistFailStatus",
			context:  &ConversionContext{StepType: StepTypeHelmBGDeploy},
			expected: "steps.step1.steps.helmBluegreenDeployAction.spec.env.PLUGIN_IGNORE_HISTORY_FAILURE",
		},
	}

	trie := buildPipelineTrie()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := trie.Match(tt.path, tt.context)
			if result != tt.expected {
				t.Errorf("Match() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTrieRules_BuildAndPushSteps(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		context  *ConversionContext
		expected string
	}{
		// BuildAndPushDockerRegistry
		{
			name:     "BuildAndPushDockerRegistry repo FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.repo",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushDockerRegistry},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.REPO",
		},
		{
			name:     "BuildAndPushDockerRegistry tags FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.tags",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushDockerRegistry},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.TAGS",
		},
		{
			name:     "BuildAndPushDockerRegistry dockerfile FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.dockerfile",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushDockerRegistry},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.DOCKERFILE",
		},
		{
			name:     "BuildAndPushDockerRegistry context FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.context",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushDockerRegistry},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.CONTEXT",
		},
		{
			name:     "BuildAndPushDockerRegistry labels FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.labels",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushDockerRegistry},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.CUSTOM_LABELS",
		},
		{
			name:     "BuildAndPushDockerRegistry buildArgs FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.buildArgs",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushDockerRegistry},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.BUILD_ARGS",
		},
		{
			name:     "BuildAndPushDockerRegistry target FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.target",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushDockerRegistry},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.TARGET",
		},
		// BuildAndPushECR
		{
			name:     "BuildAndPushECR region FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.region",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushECR},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.REGION",
		},
		{
			name:     "BuildAndPushECR imageName FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.imageName",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushECR},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.REPO",
		},
		{
			name:     "BuildAndPushECR tags FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.tags",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushECR},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.TAGS",
		},
		{
			name:     "BuildAndPushECR dockerfile FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.dockerfile",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushECR},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.DOCKERFILE",
		},
		// BuildAndPushGAR
		{
			name:     "BuildAndPushGAR imageName FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.imageName",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushGAR},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.REPO",
		},
		{
			name:     "BuildAndPushGAR tags FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.tags",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushGAR},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.TAGS",
		},
		{
			name:     "BuildAndPushGAR labels FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.labels",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushGAR},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.CUSTOM_LABELS",
		},
		// BuildAndPushACR
		{
			name:     "BuildAndPushACR registry FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.registry",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushACR},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.REGISTRY",
		},
		{
			name:     "BuildAndPushACR imageName FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.imageName",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushACR},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.REPO",
		},
		{
			name:     "BuildAndPushACR tags FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.tags",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushACR},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.TAGS",
		},
		{
			name:     "BuildAndPushACR subscriptionId FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.subscriptionId",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushACR},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.SUBSCRIPTION_ID",
		},
		{
			name:     "BuildAndPushACR dockerfile FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.dockerfile",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushACR},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.DOCKERFILE",
		},
		// BuildAndPushDockerRegistry relative path
		{
			name:     "BuildAndPushDockerRegistry repo relative",
			path:     "execution.steps.step1.spec.repo",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushDockerRegistry},
			expected: "steps.step1.steps.pushWithBuildx.spec.with.REPO",
		},
		// BuildAndPushECR inside stepGroup
		{
			name:     "BuildAndPushECR region in stepGroup",
			path:     "pipeline.stages.build.spec.execution.steps.group1.steps.step1.spec.region",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushECR},
			expected: "pipeline.stages.build.steps.group1.steps.step1.steps.pushWithBuildx.spec.with.REGION",
		},
	}

	trie := buildPipelineTrie()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := trie.Match(tt.path, tt.context)
			if result != tt.expected {
				t.Errorf("Match() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTrieRules_CIInfrastructureAndRuntime(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		context  *ConversionContext
		expected string
	}{
		// KubernetesDirect Infrastructure -> runtime.kubernetes
		{
			name:     "K8s infrastructure connectorRef FQN",
			path:     "stage.spec.infrastructure.spec.connectorRef",
			context:  nil,
			expected: "stage.runtime.kubernetes.connector",
		},
		{
			name:     "K8s infrastructure namespace FQN",
			path:     "stage.spec.infrastructure.spec.namespace",
			context:  nil,
			expected: "stage.runtime.kubernetes.namespace",
		},
		{
			name:     "K8s infrastructure annotations FQN",
			path:     "stage.spec.infrastructure.spec.annotations",
			context:  nil,
			expected: "stage.runtime.kubernetes.annotations",
		},
		{
			name:     "K8s infrastructure labels FQN",
			path:     "stage.spec.infrastructure.spec.labels",
			context:  nil,
			expected: "stage.runtime.kubernetes.labels",
		},
		{
			name:     "K8s infrastructure serviceAccountName FQN",
			path:     "stage.spec.infrastructure.spec.serviceAccountName",
			context:  nil,
			expected: "stage.runtime.kubernetes.service-account",
		},
		{
			name:     "K8s infrastructure initTimeout FQN",
			path:     "stage.spec.infrastructure.spec.initTimeout",
			context:  nil,
			expected: "stage.runtime.kubernetes.timeout",
		},
		{
			name:     "K8s infrastructure nodeSelector FQN",
			path:     "stage.spec.infrastructure.spec.nodeSelector",
			context:  nil,
			expected: "stage.runtime.kubernetes.node",
		},
		{
			name:     "K8s infrastructure hostNames FQN",
			path:     "stage.spec.infrastructure.spec.hostNames",
			context:  nil,
			expected: "stage.runtime.kubernetes.host",
		},
		{
			name:     "K8s infrastructure tolerations FQN",
			path:     "stage.spec.infrastructure.spec.tolerations",
			context:  nil,
			expected: "stage.runtime.kubernetes.tolerations",
		},
		{
			name:     "K8s infrastructure automountServiceAccountToken FQN",
			path:     "stage.spec.infrastructure.spec.automountServiceAccountToken",
			context:  nil,
			expected: "stage.runtime.kubernetes.automount-service-token",
		},
		{
			name:     "K8s infrastructure containerSecurityContext FQN",
			path:     "stage.spec.infrastructure.spec.containerSecurityContext",
			context:  nil,
			expected: "stage.runtime.kubernetes.security-context",
		},
		{
			name:     "K8s infrastructure containerSecurityContext.privileged FQN",
			path:     "stage.spec.infrastructure.spec.containerSecurityContext.privileged",
			context:  nil,
			expected: "stage.runtime.kubernetes.security-context.privileged",
		},
		{
			name:     "K8s infrastructure containerSecurityContext.runAsUser FQN",
			path:     "stage.spec.infrastructure.spec.containerSecurityContext.runAsUser",
			context:  nil,
			expected: "stage.runtime.kubernetes.security-context.user",
		},
		{
			name:     "K8s infrastructure containerSecurityContext.capabilities.add FQN",
			path:     "stage.spec.infrastructure.spec.containerSecurityContext.capabilities.add",
			context:  nil,
			expected: "stage.runtime.kubernetes.security-context.capabilities.add",
		},
		{
			name:     "K8s infrastructure priorityClassName FQN",
			path:     "stage.spec.infrastructure.spec.priorityClassName",
			context:  nil,
			expected: "stage.runtime.kubernetes.priority-class",
		},
		{
			name:     "K8s infrastructure os FQN",
			path:     "stage.spec.infrastructure.spec.os",
			context:  nil,
			expected: "stage.runtime.kubernetes.os",
		},
		{
			name:     "K8s infrastructure harnessImageConnectorRef FQN",
			path:     "stage.spec.infrastructure.spec.harnessImageConnectorRef",
			context:  nil,
			expected: "stage.runtime.kubernetes.harness-image-connector",
		},
		{
			name:     "K8s infrastructure imagePullPolicy FQN",
			path:     "stage.spec.infrastructure.spec.imagePullPolicy",
			context:  nil,
			expected: "stage.runtime.kubernetes.pull",
		},
		{
			name:     "K8s infrastructure podSpecOverlay FQN",
			path:     "stage.spec.infrastructure.spec.podSpecOverlay",
			context:  nil,
			expected: "stage.runtime.kubernetes.pod-spec-overlay",
		},
		{
			name:     "K8s infrastructure runAsUser FQN",
			path:     "stage.spec.infrastructure.spec.runAsUser",
			context:  nil,
			expected: "stage.runtime.kubernetes.user",
		},
		{
			name:     "K8s infrastructure volumes FQN",
			path:     "stage.spec.infrastructure.spec.volumes",
			context:  nil,
			expected: "stage.runtime.kubernetes.volumes",
		},
		// Cloud Runtime -> runtime.cloud
		{
			name:     "Cloud runtime size FQN",
			path:     "stage.spec.runtime.spec.size",
			context:  nil,
			expected: "stage.runtime.cloud.size",
		},
		{
			name:     "Cloud runtime imageSpec.imageName FQN",
			path:     "stage.spec.runtime.spec.imageSpec.imageName",
			context:  nil,
			expected: "stage.runtime.cloud.image",
		},
		// Relative paths
		{
			name:     "K8s infrastructure connectorRef relative",
			path:     "spec.infrastructure.spec.connectorRef",
			context:  nil,
			expected: "runtime.kubernetes.connector",
		},
		{
			name:     "K8s infrastructure namespace relative",
			path:     "spec.infrastructure.spec.namespace",
			context:  nil,
			expected: "runtime.kubernetes.namespace",
		},
		{
			name:     "Cloud runtime size relative",
			path:     "spec.runtime.spec.size",
			context:  nil,
			expected: "runtime.cloud.size",
		},
		{
			name:     "Cloud runtime image relative",
			path:     "spec.runtime.spec.imageSpec.imageName",
			context:  nil,
			expected: "runtime.cloud.image",
		},
	}

	trie := buildPipelineTrie()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := trie.Match(tt.path, tt.context)
			if result != tt.expected {
				t.Errorf("Match() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTrieRules_StepOutputConversions(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		context  *ConversionContext
		expected string
	}{
		// ============================================================
		// HTTP Step Outputs
		// ============================================================
		{
			name:     "HTTP output httpUrl FQN",
			path:     "pipeline.stages.build.spec.execution.steps.http1.output.httpUrl",
			context:  &ConversionContext{StepType: StepTypeHTTP},
			expected: "pipeline.stages.build.steps.http1.steps.httpStep.output.outputVariables.PLUGIN_HTTP_URL",
		},
		{
			name:     "HTTP output httpMethod FQN",
			path:     "pipeline.stages.build.spec.execution.steps.http1.output.httpMethod",
			context:  &ConversionContext{StepType: StepTypeHTTP},
			expected: "pipeline.stages.build.steps.http1.steps.httpStep.output.outputVariables.PLUGIN_HTTP_METHOD",
		},
		{
			name:     "HTTP output httpResponseCode FQN",
			path:     "pipeline.stages.build.spec.execution.steps.http1.output.httpResponseCode",
			context:  &ConversionContext{StepType: StepTypeHTTP},
			expected: "pipeline.stages.build.steps.http1.steps.httpStep.output.outputVariables.PLUGIN_HTTP_RESPONSE_CODE",
		},
		{
			name:     "HTTP output status FQN",
			path:     "pipeline.stages.build.spec.execution.steps.http1.output.status",
			context:  &ConversionContext{StepType: StepTypeHTTP},
			expected: "pipeline.stages.build.steps.http1.steps.httpStep.output.outputVariables.PLUGIN_EXECUTION_STATUS",
		},
		{
			name:     "HTTP output responseHeaders relative",
			path:     "execution.steps.http1.output.responseHeaders",
			context:  &ConversionContext{StepType: StepTypeHTTP},
			expected: "steps.http1.steps.httpStep.output.outputVariables.PLUGIN_RESPONSE_HEADERS",
		},

		// ============================================================
		// K8sRollingDeploy Step Outputs
		// ============================================================
		{
			name:     "K8sRollingDeploy output releaseName FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.releaseName",
			context:  &ConversionContext{StepType: StepTypeK8sRollingDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sApplyAction.output.outputVariables.PLUGIN_RELEASE_NAME",
		},
		{
			name:     "K8sRollingDeploy output releaseNumber FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.releaseNumber",
			context:  &ConversionContext{StepType: StepTypeK8sRollingDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sRollingPrepareAction.output.outputVariables.PLUGIN_RELEASE_NUMBER",
		},
		{
			name:     "K8sRollingDeploy output prunedResourceIds FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.prunedResourceIds",
			context:  &ConversionContext{StepType: StepTypeK8sRollingDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sRollingPrepareAction.output.outputVariables.PLUGIN_RESOURCES_FOR_PRUNING",
		},

		// ============================================================
		// K8sRollingRollback Step Outputs
		// ============================================================
		{
			name:     "K8sRollingRollback output recreatedResourceIds FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.recreatedResourceIds",
			context:  &ConversionContext{StepType: StepTypeK8sRollingRollback},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sRollingRollbackAction.output.outputVariables.PLUGIN_RECREATED_RESOURCES",
		},

		// ============================================================
		// K8sCanaryDeploy Step Outputs
		// ============================================================
		{
			name:     "K8sCanaryDeploy output releaseName FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.releaseName",
			context:  &ConversionContext{StepType: StepTypeK8sCanaryDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sApplyAction.output.outputVariables.PLUGIN_RELEASE_NAME",
		},
		{
			name:     "K8sCanaryDeploy output targetInstances FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.targetInstances",
			context:  &ConversionContext{StepType: StepTypeK8sCanaryDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sCanaryPrepareAction.output.outputVariables.PLUGIN_TARGET_INSTANCES",
		},
		{
			name:     "K8sCanaryDeploy output releaseNumber FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.releaseNumber",
			context:  &ConversionContext{StepType: StepTypeK8sCanaryDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sCanaryPrepareAction.output.outputVariables.PLUGIN_RELEASE_NUMBER",
		},
		{
			name:     "K8sCanaryDeploy output canaryWorkload FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.canaryWorkload",
			context:  &ConversionContext{StepType: StepTypeK8sCanaryDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sCanaryPrepareAction.output.outputVariables.PLUGIN_CANARY_WORKLOADS",
		},

		// ============================================================
		// K8sBlueGreenDeploy Step Outputs
		// ============================================================
		{
			name:     "K8sBlueGreenDeploy output releaseNumber FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.releaseNumber",
			context:  &ConversionContext{StepType: StepTypeK8sBlueGreenDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sBlueGreenPrepareAction.output.outputVariables.PLUGIN_RELEASE_NUMBER",
		},
		{
			name:     "K8sBlueGreenDeploy output releaseName FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.releaseName",
			context:  &ConversionContext{StepType: StepTypeK8sBlueGreenDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sApplyAction.output.outputVariables.PLUGIN_RELEASE_NAME",
		},
		{
			name:     "K8sBlueGreenDeploy output primaryServiceName FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.primaryServiceName",
			context:  &ConversionContext{StepType: StepTypeK8sBlueGreenDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sBlueGreenPrepareAction.output.outputVariables.PLUGIN_STABLE_SERVICE",
		},
		{
			name:     "K8sBlueGreenDeploy output stageServiceName FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.stageServiceName",
			context:  &ConversionContext{StepType: StepTypeK8sBlueGreenDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sBlueGreenPrepareAction.output.outputVariables.PLUGIN_STAGE_SERVICE",
		},
		{
			name:     "K8sBlueGreenDeploy output stageDeploymentSkipped FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.stageDeploymentSkipped",
			context:  &ConversionContext{StepType: StepTypeK8sBlueGreenDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sBlueGreenPrepareAction.output.outputVariables.PLUGIN_SKIPPED",
		},
		{
			name:     "K8sBlueGreenDeploy output prunedResourceIds FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.prunedResourceIds",
			context:  &ConversionContext{StepType: StepTypeK8sBlueGreenDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sBlueGreenPrepareAction.output.outputVariables.PLUGIN_RESOURCES_FOR_PRUNING",
		},

		// ============================================================
		// K8sBGSwapServices Step Outputs
		// ============================================================
		{
			name:     "K8sBGSwapServices output primaryColor FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.primaryColor",
			context:  &ConversionContext{StepType: StepTypeK8sBGSwapServices},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sBlueGreenSwapServicesSelectorsAction.output.outputVariables.PLUGIN_PRIMARY_COLOR",
		},
		{
			name:     "K8sBGSwapServices output stageColor FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.stageColor",
			context:  &ConversionContext{StepType: StepTypeK8sBGSwapServices},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sBlueGreenSwapServicesSelectorsAction.output.outputVariables.PLUGIN_STAGE_COLOR",
		},
		{
			name:     "K8sBGSwapServices output stableService FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.stableService",
			context:  &ConversionContext{StepType: StepTypeK8sBGSwapServices},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sBlueGreenSwapServicesSelectorsAction.output.outputVariables.PLUGIN_STABLE_SERVICE",
		},
		{
			name:     "K8sBGSwapServices output stageService FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.stageService",
			context:  &ConversionContext{StepType: StepTypeK8sBGSwapServices},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sBlueGreenSwapServicesSelectorsAction.output.outputVariables.PLUGIN_STAGE_SERVICE",
		},

		// ============================================================
		// K8sDryRun Step Outputs
		// ============================================================
		{
			name:     "K8sDryRun output manifestDryRun FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.manifestDryRun",
			context:  &ConversionContext{StepType: StepTypeK8sDryRun},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sDryRunAction.output.outputVariables.PLUGIN_MANIFEST_OUTPUT_FILE_PATH",
		},

		// ============================================================
		// K8sApply Step Outputs
		// ============================================================
		{
			name:     "K8sApply output manifestOutputFilePath FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.manifestOutputFilePath",
			context:  &ConversionContext{StepType: StepTypeK8sApply},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sApplyAction.output.outputVariables.PLUGIN_MANIFEST_OUTPUT_FILE_PATH",
		},
		{
			name:     "K8sApply output releaseNumber FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.releaseNumber",
			context:  &ConversionContext{StepType: StepTypeK8sApply},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sApplyAction.output.outputVariables.PLUGIN_RELEASE_NUMBER",
		},
		{
			name:     "K8sApply output releaseName FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.releaseName",
			context:  &ConversionContext{StepType: StepTypeK8sApply},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sApplyAction.output.outputVariables.PLUGIN_RELEASE_NAME",
		},
		{
			name:     "K8sApply output managedWorkloads FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.managedWorkloads",
			context:  &ConversionContext{StepType: StepTypeK8sApply},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sApplyAction.output.outputVariables.PLUGIN_MANAGED_WORKLOADS",
		},
		{
			name:     "K8sApply output isOpenshift FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.isOpenshift",
			context:  &ConversionContext{StepType: StepTypeK8sApply},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sApplyAction.output.outputVariables.HARNESS_IS_OPENSHIFT",
		},
		{
			name:     "K8sApply output pods FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.pods",
			context:  &ConversionContext{StepType: StepTypeK8sApply},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sApplySteadyStateCheckAction.output.outputVariables.PLUGIN_PODS",
		},
		{
			name:     "K8sApply output instances FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.instances",
			context:  &ConversionContext{StepType: StepTypeK8sApply},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sApplySteadyStateCheckAction.output.outputVariables.HARNESS_INSTANCES",
		},

		// ============================================================
		// K8sScale Step Outputs
		// ============================================================
		{
			name:     "K8sScale output pods FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.pods",
			context:  &ConversionContext{StepType: StepTypeK8sScale},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sScaleAction.output.outputVariables.PLUGIN_PODS",
		},
		{
			name:     "K8sScale output instances FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.instances",
			context:  &ConversionContext{StepType: StepTypeK8sScale},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sScaleAction.output.outputVariables.HARNESS_INSTANCES",
		},

		// ============================================================
		// K8sBlueGreenStageScaleDown Step Outputs
		// ============================================================
		{
			name:     "K8sBlueGreenStageScaleDown output scaledResources FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.scaledResources",
			context:  &ConversionContext{StepType: StepTypeK8sBlueGreenStageScaleDown},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sBlueGreenStageScaleDownAction.output.outputVariables.PLUGIN_SCALED_RESOURCES",
		},

		// ============================================================
		// K8sTrafficRouting Step Outputs
		// ============================================================
		{
			name:     "K8sTrafficRouting output manifestOutputFilePath FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.manifestOutputFilePath",
			context:  &ConversionContext{StepType: StepTypeK8sTrafficRouting},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sTrafficShiftAction.output.outputVariables.PLUGIN_MANIFEST_OUTPUT_FILE_PATH",
		},
		{
			name:     "K8sTrafficRouting output resources FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.resources",
			context:  &ConversionContext{StepType: StepTypeK8sTrafficRouting},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sTrafficShiftAction.output.outputVariables.PLUGIN_RESOURCES",
		},
		{
			name:     "K8sTrafficRouting output releaseNumber FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.releaseNumber",
			context:  &ConversionContext{StepType: StepTypeK8sTrafficRouting},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sTrafficShiftAction.output.outputVariables.PLUGIN_RELEASE_NUMBER",
		},
		{
			name:     "K8sTrafficRouting output resourcesForPruning FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.resourcesForPruning",
			context:  &ConversionContext{StepType: StepTypeK8sTrafficRouting},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sTrafficShiftAction.output.outputVariables.PLUGIN_RESOURCES_FOR_PRUNING",
		},

		// ============================================================
		// HelmDeploy Step Outputs
		// ============================================================
		{
			name:     "HelmDeploy output releaseName FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.releaseName",
			context:  &ConversionContext{StepType: StepTypeHelmDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBasicDeployAction.output.outputVariables.PLUGIN_RELEASE_NAME",
		},
		{
			name:     "HelmDeploy output prevReleaseVersion FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.prevReleaseVersion",
			context:  &ConversionContext{StepType: StepTypeHelmDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBasicDeployAction.output.outputVariables.PLUGIN_PREVIOUS_RELEASE_REVISION",
		},
		{
			name:     "HelmDeploy output newReleaseVersion FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.newReleaseVersion",
			context:  &ConversionContext{StepType: StepTypeHelmDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBasicDeployAction.output.outputVariables.PLUGIN_NEW_RELEASE_REVISION",
		},
		{
			name:     "HelmDeploy output containerInfoList FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.containerInfoList",
			context:  &ConversionContext{StepType: StepTypeHelmDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBasicDeployAction.output.outputVariables.PLUGIN_CONTAINER_INFO_LIST",
		},

		// ============================================================
		// HelmRollback Step Outputs
		// ============================================================
		{
			name:     "HelmRollback output releaseName FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.releaseName",
			context:  &ConversionContext{StepType: StepTypeHelmRollback},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmRollbackAction.output.outputVariables.PLUGIN_RELEASE_NAME",
		},
		{
			name:     "HelmRollback output newReleaseVersion FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.newReleaseVersion",
			context:  &ConversionContext{StepType: StepTypeHelmRollback},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmRollbackAction.output.outputVariables.PLUGIN_NEW_RELEASE_REVISION",
		},
		{
			name:     "HelmRollback output rollbackVersion FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.rollbackVersion",
			context:  &ConversionContext{StepType: StepTypeHelmRollback},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmRollbackAction.output.outputVariables.PLUGIN_ROLLBACK_RELEASE_REVISION",
		},
		{
			name:     "HelmRollback output containerInfoList FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.containerInfoList",
			context:  &ConversionContext{StepType: StepTypeHelmRollback},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmRollbackAction.output.outputVariables.PLUGIN_CONTAINER_INFO_LIST",
		},

		// ============================================================
		// HelmBGDeploy Step Outputs
		// ============================================================
		{
			name:     "HelmBGDeploy output releaseNumber FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.releaseNumber",
			context:  &ConversionContext{StepType: StepTypeHelmBGDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBluegreenDeployAction.output.outputVariables.PLUGIN_NEW_RELEASE_REVISION",
		},
		{
			name:     "HelmBGDeploy output releaseName FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.releaseName",
			context:  &ConversionContext{StepType: StepTypeHelmBGDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBluegreenDeployAction.output.outputVariables.PLUGIN_RELEASE_NAME",
		},
		{
			name:     "HelmBGDeploy output primaryServiceNames FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.primaryServiceNames",
			context:  &ConversionContext{StepType: StepTypeHelmBGDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBluegreenDeployAction.output.outputVariables.PLUGIN_STABLE_SERVICE",
		},
		{
			name:     "HelmBGDeploy output stageServiceNames FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.stageServiceNames",
			context:  &ConversionContext{StepType: StepTypeHelmBGDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBluegreenDeployAction.output.outputVariables.PLUGIN_STAGE_SERVICE",
		},
		{
			name:     "HelmBGDeploy output stageColor FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.stageColor",
			context:  &ConversionContext{StepType: StepTypeHelmBGDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBluegreenDeployAction.output.outputVariables.PLUGIN_STAGE_COLOR",
		},
		{
			name:     "HelmBGDeploy output primaryColor FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.primaryColor",
			context:  &ConversionContext{StepType: StepTypeHelmBGDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBluegreenDeployAction.output.outputVariables.PLUGIN_STABLE_COLOR",
		},

		// ============================================================
		// HelmCanaryDeploy Step Outputs
		// ============================================================
		{
			name:     "HelmCanaryDeploy output releaseName FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.releaseName",
			context:  &ConversionContext{StepType: StepTypeHelmCanaryDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmCanaryDeployAction.output.outputVariables.PLUGIN_CANARY_RELEASE_NAME",
		},
		{
			name:     "HelmCanaryDeploy output canaryWorkload FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.canaryWorkload",
			context:  &ConversionContext{StepType: StepTypeHelmCanaryDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmCanaryDeployAction.output.outputVariables.PLUGIN_WORKLOADS",
		},
		{
			name:     "HelmCanaryDeploy output previousReleaseVersion FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.output.previousReleaseVersion",
			context:  &ConversionContext{StepType: StepTypeHelmCanaryDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmCanaryDeployAction.output.outputVariables.PLUGIN_PREVIOUS_RELEASE_REVISION",
		},

		// ============================================================
		// Relative Path Tests
		// ============================================================
		{
			name:     "K8sRollingDeploy output releaseName relative",
			path:     "execution.steps.step1.output.releaseName",
			context:  &ConversionContext{StepType: StepTypeK8sRollingDeploy},
			expected: "steps.step1.steps.k8sApplyAction.output.outputVariables.PLUGIN_RELEASE_NAME",
		},
		{
			name:     "HelmDeploy output releaseName relative",
			path:     "execution.steps.step1.output.releaseName",
			context:  &ConversionContext{StepType: StepTypeHelmDeploy},
			expected: "steps.step1.steps.helmBasicDeployAction.output.outputVariables.PLUGIN_RELEASE_NAME",
		},
		{
			name:     "K8sApply output pods relative",
			path:     "execution.steps.step1.output.pods",
			context:  &ConversionContext{StepType: StepTypeK8sApply},
			expected: "steps.step1.steps.k8sApplySteadyStateCheckAction.output.outputVariables.PLUGIN_PODS",
		},

		// ============================================================
		// Step Group Tests
		// ============================================================
		{
			name:     "K8sRollingDeploy output in step group FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.group1.steps.step1.output.releaseName",
			context:  &ConversionContext{StepType: StepTypeK8sRollingDeploy},
			expected: "pipeline.stages.deploy.steps.group1.steps.step1.steps.k8sApplyAction.output.outputVariables.PLUGIN_RELEASE_NAME",
		},
		{
			name:     "HelmBGDeploy output in step group FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.group1.steps.step1.output.stageColor",
			context:  &ConversionContext{StepType: StepTypeHelmBGDeploy},
			expected: "pipeline.stages.deploy.steps.group1.steps.step1.steps.helmBluegreenDeployAction.output.outputVariables.PLUGIN_STAGE_COLOR",
		},
	}

	trie := buildPipelineTrie()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := trie.Match(tt.path, tt.context)
			if result != tt.expected {
				t.Errorf("Match() = %v, want %v", result, tt.expected)
			}
		})
	}
}
