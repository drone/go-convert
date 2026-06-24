package convertexpressions

import (
	"testing"
)

func TestTrieConvert_CodebasePaths(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Codebase repoName",
			input:    "<+pipeline.properties.ci.codebase.repoName>",
			expected: "<+codebase.repoName>",
		},
		{
			name:     "Codebase branch",
			input:    "<+pipeline.properties.ci.codebase.branch>",
			expected: "<+codebase.branch>",
		},
		{
			name:     "Codebase commitSha",
			input:    "<+pipeline.properties.ci.codebase.commitSha>",
			expected: "<+codebase.commitSha>",
		},
		{
			name:     "Codebase targetBranch",
			input:    "<+pipeline.properties.ci.codebase.targetBranch>",
			expected: "<+codebase.targetBranch>",
		},
		{
			name:     "Codebase sourceBranch",
			input:    "<+pipeline.properties.ci.codebase.sourceBranch>",
			expected: "<+codebase.sourceBranch>",
		},
		{
			name:     "Codebase prNumber",
			input:    "<+pipeline.properties.ci.codebase.prNumber>",
			expected: "<+codebase.prNumber>",
		},
		{
			name:     "Codebase prTitle",
			input:    "<+pipeline.properties.ci.codebase.prTitle>",
			expected: "<+codebase.prTitle>",
		},
		{
			name:     "Codebase build type",
			input:    "<+pipeline.properties.ci.codebase.build.type>",
			expected: "<+codebase.build.type>",
		},
		{
			name:     "Codebase build type branch",
			input:    "<+pipeline.properties.ci.codebase.build.spec.branch>",
			expected: "<+codebase.branch>",
		},
		{
			name:     "Codebase gitUserId",
			input:    "<+pipeline.properties.ci.codebase.gitUserId>",
			expected: "<+codebase.gitUserId>",
		},
		{
			name:     "Codebase gitUserEmail",
			input:    "<+pipeline.properties.ci.codebase.gitUserEmail>",
			expected: "<+codebase.gitUserEmail>",
		},
		{
			name:     "Codebase path in mixed text",
			input:    `echo "Repo: <+pipeline.properties.ci.codebase.repoName> Branch: <+pipeline.properties.ci.codebase.branch>"`,
			expected: `echo "Repo: <+codebase.repoName> Branch: <+codebase.branch>"`,
		},
		{
			name:     "Codebase path with function call",
			input:    "<+pipeline.properties.ci.codebase.repoName.toUpperCase()>",
			expected: "<+codebase.repoName.toUpperCase()>",
		},
		{
			name:     "Relative codebase alias",
			input:    "<+codebase.repoName>",
			expected: "<+codebase.repoName>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertExpressionWithTrie(tt.input, nil, false)
			if got != tt.expected {
				t.Errorf("ConvertExpressionWithTrie() failed\ninput:    %s\ngot:      %s\nexpected: %s", tt.input, got, tt.expected)
			}
		})
	}
}

func TestTrieConvert_InputSetPaths(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "InputSet pipeline variable",
			input:    "<+inputSet.pipeline.variables.var1>",
			expected: "<+inputSet.overlay.pipeline.variables.var1>",
		},
		{
			name:     "InputSet pipeline identifier",
			input:    "<+inputSet.pipeline.identifier>",
			expected: "<+inputSet.overlay.pipeline.identifier>",
		},
		{
			name:     "InputSet deep pipeline path",
			input:    "<+inputSet.pipeline.stages.build.spec.execution.steps.s1.name>",
			expected: "<+inputSet.overlay.pipeline.stages.build.spec.execution.steps.s1.name>",
		},
		{
			name:     "InputSet path with function call",
			input:    "<+inputSet.pipeline.variables.var1.toUpperCase()>",
			expected: "<+inputSet.overlay.pipeline.variables.var1.toUpperCase()>",
		},
		{
			name:     "InputSet path in mixed text",
			input:    `echo "Var: <+inputSet.pipeline.variables.var1>"`,
			expected: `echo "Var: <+inputSet.overlay.pipeline.variables.var1>"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertExpressionWithTrie(tt.input, nil, false)
			if got != tt.expected {
				t.Errorf("ConvertExpressionWithTrie() failed\ninput:    %s\ngot:      %s\nexpected: %s", tt.input, got, tt.expected)
			}
		})
	}
}

func TestTrieConvert_PipelineLevel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple pipeline variable - no trie rule",
			input:    "<+pipeline.variables.var1>",
			expected: "<+pipeline.variables.var1>",
		},
		{
			name:     "Stage identifier",
			input:    "<+pipeline.stages.build.identifier>",
			expected: "<+pipeline.stages.build.id>",
		},
		{
			name:     "Stage env identifier",
			input:    "<+pipeline.stages.deploy.spec.env.identifier>",
			expected: "<+pipeline.stages.deploy.steps.env.id>",
		},
		{
			name:     "Stage env group ref",
			input:    "<+pipeline.stages.deploy.spec.env.envGroupRef>",
			expected: "<+pipeline.stages.deploy.steps.env.group.id>",
		},
		{
			name:     "Stage env group name",
			input:    "<+pipeline.stages.deploy.spec.env.envGroupName>",
			expected: "<+pipeline.stages.deploy.steps.env.group.name>",
		},
		{
			name:     "Stage infra connector",
			input:    "<+pipeline.stages.deploy.spec.infra.connectorRef>",
			expected: "<+pipeline.stages.deploy.steps.infra.connector>",
		},
		{
			name:     "Relative stage env identifier",
			input:    "<+stage.spec.env.identifier>",
			expected: "<+stage.steps.env.id>",
		},
		{
			name:     "Relative spec env group ref",
			input:    "<+spec.env.envGroupRef>",
			expected: "<+stage.steps.env.group.id>",
		},
		{
			name:     "Direct env identifier",
			input:    "<+env.identifier>",
			expected: "<+env.id>",
		},
		// Deployment stage fields move under stage.steps (full/spec-prefixed),
		// while bare relative forms stay at root.
		{
			name:     "Stage service identifier",
			input:    "<+pipeline.stages.deploy.spec.service.identifier>",
			expected: "<+pipeline.stages.deploy.steps.service.id>",
		},
		{
			name:     "Stage service serviceInputs",
			input:    "<+pipeline.stages.deploy.spec.service.serviceInputs>",
			expected: "<+pipeline.stages.deploy.steps.service.with.overlay>",
		},
		{
			name:     "Relative stage service identifier",
			input:    "<+stage.spec.service.identifier>",
			expected: "<+stage.steps.service.id>",
		},
		{
			name:     "Direct service identifier",
			input:    "<+service.identifier>",
			expected: "<+service.id>",
		},
		{
			name:     "Stage infra connector full",
			input:    "<+stage.spec.infra.connectorRef>",
			expected: "<+stage.steps.infra.connector>",
		},
		{
			name:     "Direct infra connector",
			input:    "<+infra.connectorRef>",
			expected: "<+infra.connector>",
		},
		{
			name:     "Stage manifests passthrough field",
			input:    "<+pipeline.stages.deploy.spec.manifests.myManifest.identifier>",
			expected: "<+pipeline.stages.deploy.steps.manifests.myManifest.identifier>",
		},
		{
			name:     "Relative stage manifests",
			input:    "<+stage.spec.manifests.myManifest.store>",
			expected: "<+stage.steps.manifests.myManifest.store>",
		},
		{
			name:     "Direct manifests passthrough",
			input:    "<+manifests.myManifest.store>",
			expected: "<+manifests.myManifest.store>",
		},
		{
			name:     "Stage configFiles field",
			input:    "<+pipeline.stages.deploy.spec.configFiles.cf1.content>",
			expected: "<+pipeline.stages.deploy.steps.configFiles.cf1.content>",
		},
		{
			name:     "Direct configFiles passthrough",
			input:    "<+configFiles.cf1.content>",
			expected: "<+configFiles.cf1.content>",
		},
		{
			name:     "Stage artifacts field",
			input:    "<+pipeline.stages.deploy.spec.artifacts.primary.tag>",
			expected: "<+pipeline.stages.deploy.steps.artifacts.primary.tag>",
		},
		{
			name:     "Direct artifacts passthrough",
			input:    "<+artifacts.primary.tag>",
			expected: "<+artifacts.primary.tag>",
		},
		{
			name:     "Spec execution steps removal",
			input:    "<+pipeline.stages.build.spec.execution.steps.step1.output.outputVariables.var1>",
			expected: "<+pipeline.stages.build.steps.step1.output.outputVariables.var1>",
		},
		{
			name:     "No conversion needed - plain text outside expression",
			input:    "pipeline.stages.build.identifier",
			expected: "pipeline.stages.build.identifier",
		},
		{
			name:     "No conversion needed - unmatched path",
			input:    "<+pipeline.name>",
			expected: "<+pipeline.name>",
		},
		{
			name:     "Partial match - converts what matches",
			input:    "<+pipeline.stages.build.spec.execution.steps.step1.unknown.field>",
			expected: "<+pipeline.stages.build.steps.step1.unknown.field>",
		},
		{
			name:     "Env with function call",
			input:    "<+pipeline.stages.deploy.spec.env.envGroupName.toUpperCase()>",
			expected: "<+pipeline.stages.deploy.steps.env.group.name.toUpperCase()>",
		},
		{
			name:     "Step with function call",
			input:    "<+pipeline.stages.build.spec.execution.steps.step1.output.outputVariables.var1.toLowerCase()>",
			expected: "<+pipeline.stages.build.steps.step1.output.outputVariables.var1.toLowerCase()>",
		},
		{
			name:     "Mixed text with expressions",
			input:    `stage: <+pipeline.stages.build.identifier> env: <+pipeline.stages.deploy.spec.env.identifier>`,
			expected: `stage: <+pipeline.stages.build.id> env: <+pipeline.stages.deploy.steps.env.id>`,
		},
		{
			name:     "JSON array with expressions",
			input:    `["<+pipeline.stages.build.identifier>", "<+pipeline.stages.deploy.spec.env.envGroupName>"]`,
			expected: `["<+pipeline.stages.build.id>", "<+pipeline.stages.deploy.steps.env.group.name>"]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertExpressionWithTrie(tt.input, nil, false)
			if got != tt.expected {
				t.Errorf("ConvertExpressionWithTrie() failed\ninput:    %s\ngot:      %s\nexpected: %s", tt.input, got, tt.expected)
			}
		})
	}
}

func TestTrieConvert_StepLevel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		context  *ConversionContext
		expected string
	}{
		// General step rules
		{
			name:     "Step identifier - no context",
			input:    "<+pipeline.stages.build.spec.execution.steps.step1.identifier>",
			context:  nil,
			expected: "<+pipeline.stages.build.steps.step1.id>",
		},
		{
			name:     "Step in stepGroup",
			input:    "<+pipeline.stages.build.spec.execution.steps.group1.steps.step1.identifier>",
			context:  nil,
			expected: "<+pipeline.stages.build.steps.group1.steps.step1.id>",
		},
		{
			name:     "Relative stepGroup reference",
			input:    "<+stepGroup.steps.step1.identifier>",
			context:  nil,
			expected: "<+group.steps.step1.id>",
		},
		{
			name:     "Step in nested stepGroup",
			input:    "<+pipeline.stages.build.spec.execution.steps.group1.steps.group2.steps.step1.identifier>",
			context:  nil,
			expected: "<+pipeline.stages.build.steps.group1.steps.group2.steps.step1.id>",
		},
		// Run step
		{
			name:     "Run step command",
			input:    "<+pipeline.stages.build.spec.execution.steps.step1.spec.command>",
			context:  &ConversionContext{StepType: StepTypeRun},
			expected: "<+pipeline.stages.build.steps.step1.spec.script>",
		},
		{
			name:     "Run step image",
			input:    "<+pipeline.stages.build.spec.execution.steps.step1.spec.image>",
			context:  &ConversionContext{StepType: StepTypeRun},
			expected: "<+pipeline.stages.build.steps.step1.spec.container.image>",
		},
		{
			name:     "Run step envVariables",
			input:    "<+pipeline.stages.build.spec.execution.steps.step1.spec.envVariables>",
			context:  &ConversionContext{StepType: StepTypeRun},
			expected: "<+pipeline.stages.build.steps.step1.spec.env>",
		},
		{
			name:     "Run step outputVariables",
			input:    "<+pipeline.stages.build.spec.execution.steps.step1.spec.outputVariables>",
			context:  &ConversionContext{StepType: StepTypeRun},
			expected: "<+pipeline.stages.build.steps.step1.spec.output>",
		},
		// Parenthesized array rules - outputVariables[i].field -> spec.output[i].field
		{
			name:     "output variables spec name",
			input:    "<+pipeline.stages.build.spec.execution.steps.step1.spec.outputVariables[1].name>",
			context:  &ConversionContext{StepType: StepTypeRun},
			expected: "<+pipeline.stages.build.steps.step1.spec.output[1].alias>",
		},
		{
			name:     "output variables spec value",
			input:    "<+pipeline.stages.build.spec.execution.steps.step1.spec.outputVariables[0].value>",
			context:  &ConversionContext{StepType: StepTypeRun},
			expected: "<+pipeline.stages.build.steps.step1.spec.output[0].name>",
		},
		{
			name:     "Relative Run step command",
			input:    "<+execution.steps.step1.spec.command>",
			context:  &ConversionContext{StepType: StepTypeRun},
			expected: "<+stage.steps.step1.spec.script>",
		},
		// Background step
		{
			name:     "Background step command",
			input:    "<+pipeline.stages.build.spec.execution.steps.bg1.spec.command>",
			context:  &ConversionContext{StepType: StepTypeBackground},
			expected: "<+pipeline.stages.build.steps.bg1.spec.script>",
		},
		// HTTP step - spec rules
		{
			name:     "HTTP step url",
			input:    "<+pipeline.stages.build.spec.execution.steps.http1.spec.url>",
			context:  &ConversionContext{StepType: StepTypeHTTP},
			expected: "<+pipeline.stages.build.steps.http1.spec.env.PLUGIN_URL>",
		},
		{
			name:     "HTTP step method",
			input:    "<+pipeline.stages.build.spec.execution.steps.http1.spec.method>",
			context:  &ConversionContext{StepType: StepTypeHTTP},
			expected: "<+pipeline.stages.build.steps.http1.spec.env.PLUGIN_METHOD>",
		},
		// HTTP step - output rules
		{
			name:     "HTTP output httpResponseCode",
			input:    "<+pipeline.stages.build.spec.execution.steps.http1.output.httpResponseCode>",
			context:  &ConversionContext{StepType: StepTypeHTTP},
			expected: "<+pipeline.stages.build.steps.http1.steps.httpStep.output.outputVariables.PLUGIN_HTTP_RESPONSE_CODE>",
		},
		// RestoreCacheS3
		{
			name:     "RestoreCacheS3 bucket",
			input:    "<+pipeline.stages.build.spec.execution.steps.restore1.spec.bucket>",
			context:  &ConversionContext{StepType: StepTypeRestoreCacheS3},
			expected: "<+pipeline.stages.build.steps.restore1.steps.restoreCacheS3.spec.with.BUCKET>",
		},
		{
			name:     "RestoreCacheS3 key",
			input:    "<+pipeline.stages.build.spec.execution.steps.restore1.spec.key>",
			context:  &ConversionContext{StepType: StepTypeRestoreCacheS3},
			expected: "<+pipeline.stages.build.steps.restore1.steps.restoreCacheS3.spec.with.CACHE_KEY>",
		},
		{
			name:     "Relative RestoreCacheS3",
			input:    "<+execution.steps.restore1.spec.bucket>",
			context:  &ConversionContext{StepType: StepTypeRestoreCacheS3},
			expected: "<+stage.steps.restore1.steps.restoreCacheS3.spec.with.BUCKET>",
		},
		{
			name:     "StepGroup RestoreCacheS3",
			input:    "<+stepGroup.steps.restore1.spec.bucket>",
			context:  &ConversionContext{StepType: StepTypeRestoreCacheS3},
			expected: "<+group.steps.restore1.steps.restoreCacheS3.spec.with.BUCKET>",
		},
		// SaveCacheGCS
		{
			name:     "SaveCacheGCS bucket",
			input:    "<+pipeline.stages.build.spec.execution.steps.save1.spec.bucket>",
			context:  &ConversionContext{StepType: StepTypeSaveCacheGCS},
			expected: "<+pipeline.stages.build.steps.save1.steps.saveCacheGCS.spec.with.BUCKET>",
		},
		{
			name:     "SaveCacheGCS sourcePaths",
			input:    "<+pipeline.stages.build.spec.execution.steps.save1.spec.sourcePaths>",
			context:  &ConversionContext{StepType: StepTypeSaveCacheGCS},
			expected: "<+pipeline.stages.build.steps.save1.steps.saveCacheGCS.spec.with.MOUNT>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertExpressionWithTrie(tt.input, tt.context, false)
			if got != tt.expected {
				t.Errorf("ConvertExpressionWithTrie() failed\ninput:    %s\ncontext:  %+v\ngot:      %s\nexpected: %s", tt.input, tt.context, got, tt.expected)
			}
		})
	}
}

func TestTrieConvert_ComplexNested(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		context  *ConversionContext
		expected string
	}{
		// Nested expressions with comparisons
		{
			name:     "Nested comparison with stage identifier",
			input:    `<+<+pipeline.stages.build.identifier> == "build">`,
			context:  nil,
			expected: `<+<+pipeline.stages.build.id> == "build">`,
		},
		{
			name:     "Nested comparison with step output",
			input:    `<+<+pipeline.stages.build.spec.execution.steps.step1.output.outputVariables.url>.contains("test")>`,
			context:  nil,
			expected: `<+<+pipeline.stages.build.steps.step1.output.outputVariables.url>.contains("test")>`,
		},
		// Multiple nested expressions
		{
			name:     "Concatenation with multiple conversions",
			input:    `<+<+pipeline.stages.build.spec.execution.steps.step1.identifier> + "-" + <+pipeline.stages.deploy.identifier>>`,
			context:  nil,
			expected: `<+<+pipeline.stages.build.steps.step1.id> + "-" + <+pipeline.stages.deploy.id>>`,
		},
		// Ternary expressions
		{
			name:     "Ternary with stage and env conversions",
			input:    `<+<+pipeline.stages.build.identifier>=="build"?<+pipeline.stages.deploy.spec.env.identifier>:"default">`,
			context:  nil,
			expected: `<+<+pipeline.stages.build.id>=="build"?<+pipeline.stages.deploy.steps.env.id>:"default">`,
		},
		{
			name:     "Ternary with step context",
			input:    `<+<+pipeline.stages.build.spec.execution.steps.run1.spec.command>.contains("test")?<+pipeline.stages.deploy.spec.env.identifier>:"default">`,
			context:  &ConversionContext{StepType: StepTypeRun},
			expected: `<+<+pipeline.stages.build.steps.run1.spec.script>.contains("test")?<+pipeline.stages.deploy.steps.env.id>:"default">`,
		},
		// Function chaining
		{
			name:     "Function chain on env group name",
			input:    `<+<+pipeline.stages.deploy.spec.env.envGroupName>.toLowerCase().contains("prod")>`,
			context:  nil,
			expected: `<+<+pipeline.stages.deploy.steps.env.group.name>.toLowerCase().contains("prod")>`,
		},
		{
			name:     "Nested step image with function call",
			input:    `<+<+pipeline.stages.build.spec.execution.steps.run1.spec.image>.contains("alpine")>`,
			context:  &ConversionContext{StepType: StepTypeRun},
			expected: `<+<+pipeline.stages.build.steps.run1.spec.container.image>.contains("alpine")>`,
		},
		// Multi-level nesting
		{
			name:     "Three levels of nesting",
			input:    `<+ <+<+pipeline.stages.build.spec.execution.steps.step1.identifier>.toLowerCase()> == "step1" >`,
			context:  nil,
			expected: `<+ <+<+pipeline.stages.build.steps.step1.id>.toLowerCase()> == "step1" >`,
		},
		// Dynamic stage reference with step conversion
		{
			name:     "Dynamic stage with nested expression and step conversion",
			input:    `<+pipeline.stages.<+<+pipeline.stages.build.identifier> == "build"?"stage1":"stage2">.spec.execution.steps.step1.output.outputVariables.url>`,
			context:  nil,
			expected: `<+pipeline.stages.<+<+pipeline.stages.build.id> == "build"?"stage1":"stage2">.steps.step1.output.outputVariables.url>`,
		},
		// Mixed text with multiple expressions
		{
			name:     "Mixed text with step and env expressions",
			input:    `step: <+pipeline.stages.build.spec.execution.steps.step1.identifier> group: <+pipeline.stages.deploy.spec.env.envGroupRef>`,
			context:  nil,
			expected: `step: <+pipeline.stages.build.steps.step1.id> group: <+pipeline.stages.deploy.steps.env.group.id>`,
		},
		{
			name:     "JSON with mixed conversions",
			input:    `{"stage": "<+pipeline.stages.build.identifier>", "env": "<+pipeline.stages.deploy.spec.env.envGroupName>", "step": "<+pipeline.stages.build.spec.execution.steps.step1.identifier>"}`,
			context:  nil,
			expected: `{"stage": "<+pipeline.stages.build.id>", "env": "<+pipeline.stages.deploy.steps.env.group.name>", "step": "<+pipeline.stages.build.steps.step1.id>"}`,
		},
		// OR expression with multiple conversions
		{
			name:     "OR expression with step output and env",
			input:    `<+ <+<+pipeline.stages.build.spec.execution.steps.step1.identifier>.contains("test")> || <+pipeline.stages.deploy.spec.env.identifier> == "prod" >`,
			context:  nil,
			expected: `<+ <+<+pipeline.stages.build.steps.step1.id>.contains("test")> || <+pipeline.stages.deploy.steps.env.id> == "prod" >`,
		},
		// Step context with nested expressions
		{
			name:     "HTTP output nested in comparison",
			input:    `<+<+pipeline.stages.build.spec.execution.steps.http1.output.httpResponseCode> == "200">`,
			context:  &ConversionContext{StepType: StepTypeHTTP},
			expected: `<+<+pipeline.stages.build.steps.http1.steps.httpStep.output.outputVariables.PLUGIN_HTTP_RESPONSE_CODE> == "200">`,
		},
		{
			name:     "Run step output in ternary",
			input:    `<+<+pipeline.stages.build.spec.execution.steps.run1.spec.command>.contains("deploy")?<+pipeline.stages.build.spec.execution.steps.run1.spec.image>:"default">`,
			context:  &ConversionContext{StepType: StepTypeRun},
			expected: `<+<+pipeline.stages.build.steps.run1.spec.script>.contains("deploy")?<+pipeline.stages.build.steps.run1.spec.container.image>:"default">`,
		},
		// String concatenation with conversions
		{
			name:     "URL concatenation with step output",
			input:    `<+"https://api.example.com/"+<+pipeline.stages.build.spec.execution.steps.step1.identifier>>`,
			context:  nil,
			expected: `<+"https://api.example.com/"+<+pipeline.stages.build.steps.step1.id>>`,
		},
		// Complex cache step in expression
		{
			name:     "Cache step bucket in comparison",
			input:    `<+<+pipeline.stages.build.spec.execution.steps.save1.spec.bucket> == "my-bucket">`,
			context:  &ConversionContext{StepType: StepTypeSaveCacheGCS},
			expected: `<+<+pipeline.stages.build.steps.save1.steps.saveCacheGCS.spec.with.BUCKET> == "my-bucket">`,
		},
		{
			name:     "Nested expression with mixed segment containing underscores and nested expressions",
			input:    `<+<+pipeline.stages.deployment_<+infra.environment.environmentRef>_<+infra.infraIdentifier>.spec.execution.steps.evaluation_release.steps.setup_envs.ContainerStep.output.outputVariables.apigee_cdev_enable>>`,
			expected: `<+<+pipeline.stages.deployment_<+infra.environment.environmentRef>_<+infra.infraIdentifier>.steps.evaluation_release.steps.setup_envs.ContainerStep.output.outputVariables.apigee_cdev_enable>>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertExpressionWithTrie(tt.input, tt.context, false)
			if got != tt.expected {
				t.Errorf("ConvertExpressionWithTrie() failed\ninput:    %s\ncontext:  %+v\ngot:      %s\nexpected: %s", tt.input, tt.context, got, tt.expected)
			}
		})
	}
}

// TestTrieConvert_LazyStepTypeResolution tests that step type resolution happens
// lazily inside the trie when it encounters step.spec expressions.
// The trie captures step IDs during path traversal and resolves step type via
// the FQN-keyed StepInfoByFQN lookup.
// TestTrieConvert_FQNMode tests that when UseFQN is enabled, relative step expressions
// are converted to fully qualified names.
// - "step." prefix uses CurrentFQN and CurrentStepType (the step we're inside)
// - "steps.STEPID" prefix uses StepInfoByFQN to look up the referenced step
func TestTrieConvert_FQNMode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		context  *ConversionContext
		expected string
	}{
		// ==========================================
		// "step." prefix - uses CurrentStepV1Path/Type
		// ==========================================
		{
			name:  "step.spec.bucket with FQN - RestoreCacheGCS",
			input: "<+step.spec.bucket>",
			context: &ConversionContext{
				UseFQN:          true,
				CurrentFQN:      "pipeline.stages.build.steps.restoreCache",
				CurrentStepType: StepTypeRestoreCacheGCS,
			},
			expected: "<+pipeline.stages.build.steps.restoreCache.steps.restoreCacheGCS.spec.with.BUCKET>",
		},
		{
			name:  "step.spec.key with FQN - RestoreCacheS3",
			input: "<+step.spec.key>",
			context: &ConversionContext{
				UseFQN:          true,
				CurrentFQN:      "pipeline.stages.build.steps.restoreS3",
				CurrentStepType: StepTypeRestoreCacheS3,
			},
			expected: "<+pipeline.stages.build.steps.restoreS3.steps.restoreCacheS3.spec.with.CACHE_KEY>",
		},
		{
			name:  "step.spec.command with FQN - Run step",
			input: "<+step.spec.command>",
			context: &ConversionContext{
				UseFQN:          true,
				CurrentFQN:      "pipeline.stages.build.steps.runStep",
				CurrentStepType: StepTypeRun,
			},
			expected: "<+pipeline.stages.build.steps.runStep.spec.script>",
		},
		{
			// .identifier is not a spec/output node, so FQN substitution is
			// suppressed and the relative "step" alias is preserved.
			name:  "step.identifier with FQN",
			input: "<+step.identifier>",
			context: &ConversionContext{
				UseFQN:          true,
				CurrentFQN:      "pipeline.stages.build.steps.myStep",
				CurrentStepType: StepTypeRun,
			},
			expected: "<+step.id>",
		},
		// Step inside step group with FQN
		{
			name:  "step.spec.bucket in step group with FQN",
			input: "<+step.spec.bucket>",
			context: &ConversionContext{
				UseFQN:          true,
				CurrentFQN:      "pipeline.stages.build.steps.cacheGroup.steps.restoreCache",
				CurrentStepType: StepTypeRestoreCacheGCS,
			},
			expected: "<+pipeline.stages.build.steps.cacheGroup.steps.restoreCache.steps.restoreCacheGCS.spec.with.BUCKET>",
		},
		// Nested step group with FQN
		{
			name:  "step.spec.bucket in nested step group with FQN",
			input: "<+step.spec.bucket>",
			context: &ConversionContext{
				UseFQN:          true,
				CurrentFQN:      "pipeline.stages.build.steps.outerGroup.steps.innerGroup.steps.restoreCache",
				CurrentStepType: StepTypeRestoreCacheGCS,
			},
			expected: "<+pipeline.stages.build.steps.outerGroup.steps.innerGroup.steps.restoreCache.steps.restoreCacheGCS.spec.with.BUCKET>",
		},
		// HTTP step with FQN
		{
			name:  "step.spec.url with FQN - HTTP step",
			input: "<+step.spec.url>",
			context: &ConversionContext{
				UseFQN:          true,
				CurrentFQN:      "pipeline.stages.build.steps.httpCall",
				CurrentStepType: StepTypeHTTP,
			},
			expected: "<+pipeline.stages.build.steps.httpCall.spec.env.PLUGIN_URL>",
		},
		// HTTP output with FQN
		{
			name:  "step.output.httpResponseCode with FQN - HTTP step",
			input: "<+step.output.httpResponseCode>",
			context: &ConversionContext{
				UseFQN:          true,
				CurrentFQN:      "pipeline.stages.build.steps.httpCall",
				CurrentStepType: StepTypeHTTP,
			},
			expected: "<+pipeline.stages.build.steps.httpCall.steps.httpStep.output.outputVariables.PLUGIN_HTTP_RESPONSE_CODE>",
		},
		// Multiple expressions in same string with FQN
		{
			name:  "multiple step expressions with FQN",
			input: `bucket: <+step.spec.bucket> key: <+step.spec.key>`,
			context: &ConversionContext{
				UseFQN:          true,
				CurrentFQN:      "pipeline.stages.build.steps.restoreCache",
				CurrentStepType: StepTypeRestoreCacheS3,
			},
			expected: `bucket: <+pipeline.stages.build.steps.restoreCache.steps.restoreCacheS3.spec.with.BUCKET> key: <+pipeline.stages.build.steps.restoreCache.steps.restoreCacheS3.spec.with.CACHE_KEY>`,
		},

		// ==========================================
		// FQN disabled - should use normal conversion
		// ==========================================
		{
			name:  "step.spec.command without FQN flag",
			input: "<+step.spec.command>",
			context: &ConversionContext{
				UseFQN:          false,
				CurrentFQN:      "pipeline.stages.build.steps.runStep",
				CurrentStepType: StepTypeRun,
			},
			expected: "<+step.spec.script>", // Normal relative conversion
		},
		// FQN enabled but no CurrentStepV1Path - should use normal conversion
		{
			name:  "step.spec.command with FQN but no path",
			input: "<+step.spec.command>",
			context: &ConversionContext{
				UseFQN:          true,
				CurrentFQN:      "",
				CurrentStepType: StepTypeRun,
			},
			expected: "<+step.spec.script>", // Normal relative conversion
		},
		// Already FQN path - should convert normally (not relative)
		{
			name:  "already FQN path with FQN enabled",
			input: "<+pipeline.stages.build.spec.execution.steps.step1.spec.command>",
			context: &ConversionContext{
				UseFQN:          true,
				CurrentFQN:      "pipeline.stages.build.steps.step1",
				CurrentStepType: StepTypeRun,
				StepInfoByFQN: map[string]*StepInfoFQN{
					"pipeline.stages.build.steps.step1": {Type: StepTypeRun, StageID: "build", StepID: "step1"},
				},
			},
			expected: "<+pipeline.stages.build.steps.step1.spec.script>",
		},

		// ==========================================
		// "steps.STEPID" prefix - uses StepV1PathMap/TypeMap
		// ==========================================
		{
			name:  "steps.STEPID.spec.command with FQN - Run step",
			input: "<+steps.runStep.spec.command>",
			context: &ConversionContext{
				UseFQN:         true,
				CurrentStageID: "build",
				StepInfoByFQN: map[string]*StepInfoFQN{
					"pipeline.stages.build.steps.runStep": {Type: StepTypeRun, StageID: "build", StepID: "runStep"},
				},
			},
			expected: "<+pipeline.stages.build.steps.runStep.spec.script>",
		},
		{
			name:  "steps.STEPID.spec.bucket with FQN - RestoreCacheGCS",
			input: "<+steps.restoreCache.spec.bucket>",
			context: &ConversionContext{
				UseFQN:         true,
				CurrentStageID: "build",
				StepInfoByFQN: map[string]*StepInfoFQN{
					"pipeline.stages.build.steps.restoreCache": {Type: StepTypeRestoreCacheGCS, StageID: "build", StepID: "restoreCache"},
				},
			},
			expected: "<+pipeline.stages.build.steps.restoreCache.steps.restoreCacheGCS.spec.with.BUCKET>",
		},
		// execution.steps.STEPID.spec.bucket
		{
			name:  "execution.steps.STEPID.spec.bucket with FQN",
			input: "<+execution.steps.restoreCache.spec.bucket>",
			context: &ConversionContext{
				UseFQN:         true,
				CurrentStageID: "build",
				StepInfoByFQN: map[string]*StepInfoFQN{
					"pipeline.stages.build.steps.restoreCache": {Type: StepTypeRestoreCacheGCS, StageID: "build", StepID: "restoreCache"},
				},
			},
			expected: "<+stage.steps.restoreCache.steps.restoreCacheGCS.spec.with.BUCKET>",
		},
		// spec.execution.steps.STEPID.spec.bucket
		{
			name:  "spec.execution.steps.STEPID.spec.bucket with FQN",
			input: "<+spec.execution.steps.restoreCache.spec.bucket>",
			context: &ConversionContext{
				UseFQN:         true,
				CurrentStageID: "build",
				StepInfoByFQN: map[string]*StepInfoFQN{
					"pipeline.stages.build.steps.restoreCache": {Type: StepTypeRestoreCacheGCS, StageID: "build", StepID: "restoreCache"},
				},
			},
			expected: "<+stage.steps.restoreCache.steps.restoreCacheGCS.spec.with.BUCKET>",
		},
		// stepGroup.steps.STEPID.spec.bucket
		{
			name:  "stepGroup.steps.STEPID.spec.bucket with FQN",
			input: "<+stepGroup.steps.restoreCache.spec.bucket>",
			context: &ConversionContext{
				UseFQN:         true,
				CurrentStageID: "build",
				StepInfoByFQN: map[string]*StepInfoFQN{
					"pipeline.stages.build.steps.cacheGroup.steps.restoreCache": {Type: StepTypeRestoreCacheGCS, StageID: "build", Chain: []string{"cacheGroup"}, StepID: "restoreCache"},
				},
			},
			expected: "<+pipeline.stages.build.steps.cacheGroup.steps.restoreCache.steps.restoreCacheGCS.spec.with.BUCKET>",
		},
		// Reference a different step than the current one
		{
			name:  "steps.otherStep reference from inside currentStep",
			input: "<+steps.otherStep.output.result>",
			context: &ConversionContext{
				UseFQN:          true,
				CurrentFQN:      "pipeline.stages.build.steps.currentStep",
				CurrentStepType: StepTypeRun,
				StepInfoByFQN: map[string]*StepInfoFQN{
					"pipeline.stages.build.steps.otherStep":   {Type: StepTypeRun, StageID: "build", StepID: "otherStep"},
					"pipeline.stages.build.steps.currentStep": {Type: StepTypeRun, StageID: "build", StepID: "currentStep"},
				},
			},
			expected: "<+pipeline.stages.build.steps.otherStep.output.result>",
		},
		// steps.STEPID without StepV1PathMap - should fall back to normal conversion
		{
			name:  "steps.STEPID without StepV1PathMap",
			input: "<+steps.unknownStep.spec.command>",
			context: &ConversionContext{
				UseFQN:        true,
				StepInfoByFQN: map[string]*StepInfoFQN{},
			},
			expected: "<+steps.unknownStep.spec.script>", // Falls back to normal conversion
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertExpressionWithTrie(tt.input, tt.context, false)
			if got != tt.expected {
				t.Errorf("ConvertExpressionWithTrie() FQN mode failed\ninput:    %s\ncontext:  %+v\ngot:      %s\nexpected: %s", tt.input, tt.context, got, tt.expected)
			}
		})
	}
}

func TestTrieConvert_LazyStepTypeResolution(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		context  *ConversionContext
		expected string
	}{
		// Lazy resolution from StepTypeMap - step ID captured from path
		{
			name:  "Lazy resolution from StepTypeMap for Run step",
			input: "<+pipeline.stages.build.spec.execution.steps.runStep1.spec.command>",
			context: &ConversionContext{
				StepInfoByFQN: map[string]*StepInfoFQN{
					"pipeline.stages.build.steps.runStep1": {Type: StepTypeRun, StageID: "build", StepID: "runStep1"},
				},
				// Note: StepType is NOT set - should be resolved lazily
			},
			expected: "<+pipeline.stages.build.steps.runStep1.spec.script>",
		},
		{
			name:  "Lazy resolution from StepTypeMap for Http step",
			input: "<+pipeline.stages.build.spec.execution.steps.httpStep.spec.url>",
			context: &ConversionContext{
				StepInfoByFQN: map[string]*StepInfoFQN{
					"pipeline.stages.build.steps.httpStep": {Type: StepTypeHTTP, StageID: "build", StepID: "httpStep"},
				},
			},
			expected: "<+pipeline.stages.build.steps.httpStep.spec.env.PLUGIN_URL>",
		},
		{
			name:  "Lazy resolution for nested step group",
			input: "<+pipeline.stages.build.spec.execution.steps.myGroup.steps.innerRun.spec.command>",
			context: &ConversionContext{
				StepInfoByFQN: map[string]*StepInfoFQN{
					"pipeline.stages.build.steps.myGroup.steps.innerRun": {Type: StepTypeRun, StageID: "build", Chain: []string{"myGroup"}, StepID: "innerRun"},
				},
			},
			expected: "<+pipeline.stages.build.steps.myGroup.steps.innerRun.spec.script>",
		},
		// CurrentStepType fallback for "step.*" expressions
		{
			name:  "CurrentStepType fallback for step.spec expressions",
			input: "<+step.spec.command>",
			context: &ConversionContext{
				CurrentStepType: StepTypeRun,
			},
			expected: "<+step.spec.script>",
		},
		// Multiple steps with different types in same StepTypeMap
		{
			name:  "Multiple step types - resolves correct type for each",
			input: "<+pipeline.stages.build.spec.execution.steps.gcsStep.spec.env.bucket>",
			context: &ConversionContext{
				StepInfoByFQN: map[string]*StepInfoFQN{
					"pipeline.stages.build.steps.runStep1": {Type: StepTypeRun, StageID: "build", StepID: "runStep1"},
					"pipeline.stages.build.steps.gcsStep":  {Type: StepTypeGCSUpload, StageID: "build", StepID: "gcsStep"},
					"pipeline.stages.build.steps.httpStep": {Type: StepTypeHTTP, StageID: "build", StepID: "httpStep"},
				},
			},
			expected: "<+pipeline.stages.build.steps.gcsStep.spec.env.bucket>",
		},
		// Step ID not in map - deterministic fallback applies first matching context rule
		{
			name:  "Unknown step ID - deterministic fallback to first matching context",
			input: "<+pipeline.stages.build.spec.execution.steps.unknownStep.spec.command>",
			context: &ConversionContext{
				StepInfoByFQN: map[string]*StepInfoFQN{
					"pipeline.stages.build.steps.runStep1": {Type: StepTypeRun, StageID: "build", StepID: "runStep1"},
				},
			},
			// When step ID not found, trie uses deterministic fallback (alphabetically first matching context)
			// "Run" context has command->script rule, so it applies
			expected: "<+pipeline.stages.build.steps.unknownStep.spec.script>",
		},
		// Relative path with lazy resolution
		{
			name:  "Relative steps path with lazy resolution",
			input: "<+execution.steps.runStep1.spec.command>",
			context: &ConversionContext{
				StepInfoByFQN: map[string]*StepInfoFQN{
					"pipeline.stages.build.steps.runStep1": {Type: StepTypeRun, StageID: "build", StepID: "runStep1"},
				},
			},
			expected: "<+stage.steps.runStep1.spec.script>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertExpressionWithTrie(tt.input, tt.context, false)
			if got != tt.expected {
				t.Errorf("ConvertExpressionWithTrie() lazy resolution failed\ninput:    %s\ncontext:  %+v\ngot:      %s\nexpected: %s", tt.input, tt.context, got, tt.expected)
			}
		})
	}
}

// TestTrieConvert_DollarDelimiter verifies that ${{ ... }} delimited expressions
// are detected and converted with the same path-conversion logic as <+ ... >,
// and that the ${{ }} delimiter style is preserved on output.
func TestTrieConvert_DollarDelimiter(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		context  *ConversionContext
		expected string
	}{
		{
			name:     "Dollar step output path collapse",
			input:    "${{pipeline.stages.cs1.spec.execution.steps.ShellScript_1.output.outputVariables.rishi1}}",
			expected: "${{pipeline.stages.cs1.steps.ShellScript_1.output.outputVariables.rishi1}}",
		},
		{
			name:     "Dollar stage identifier",
			input:    "${{pipeline.stages.build.identifier}}",
			expected: "${{pipeline.stages.build.id}}",
		},
		{
			name:     "Dollar codebase prefix",
			input:    "${{pipeline.properties.ci.codebase.branch}}",
			expected: "${{codebase.branch}}",
		},
		{
			name:     "Dollar mixed with text",
			input:    "echo ${{pipeline.stages.build.spec.execution.steps.step1.output}} done",
			expected: "echo ${{pipeline.stages.build.steps.step1.output}} done",
		},
		{
			name:     "Dollar and angle delimiters in same string",
			input:    "${{pipeline.stages.build.spec.execution.steps.step1.output}} and <+pipeline.stages.deploy.spec.execution.steps.step2.output>",
			expected: "${{pipeline.stages.build.steps.step1.output}} and <+pipeline.stages.deploy.steps.step2.output>",
		},
		{
			name:     "Dollar plain string passthrough",
			input:    "${{just.a.plain.path.no.conversion}}",
			expected: "${{just.a.plain.path.no.conversion}}",
		},
		{
			name:  "Dollar step self-reference in FQN mode",
			input: "${{step.spec.command}}",
			context: &ConversionContext{
				UseFQN:          true,
				CurrentFQN:      "pipeline.stages.build.steps.compile",
				CurrentStepType: StepTypeRun,
			},
			expected: "${{pipeline.stages.build.steps.compile.spec.script}}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertExpressionWithTrie(tt.input, tt.context, false)
			if got != tt.expected {
				t.Errorf("ConvertExpressionWithTrie() dollar delimiter failed\ninput:    %s\ngot:      %s\nexpected: %s", tt.input, got, tt.expected)
			}
		})
	}
}

// TestTrieConvert_DeploymentStageFields verifies that deployment-stage spec
// fields (service/manifests/configFiles/infra/env/artifacts) move under
// stage.steps for FQN, stage-relative, and spec-relative entry points, while
// bare alias-relative forms stay at root. Fields with no rename
// (manifests/configFiles/artifacts) pass their subfields through unchanged.
func TestTrieConvert_DeploymentStageFields(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// ---------------- env (identifier/envGroupName/envGroupRef) -----------
		{
			name:     "env FQN identifier",
			input:    "<+pipeline.stages.deploy.spec.env.identifier>",
			expected: "<+pipeline.stages.deploy.steps.env.id>",
		},
		{
			name:     "env FQN group name",
			input:    "<+pipeline.stages.deploy.spec.env.envGroupName>",
			expected: "<+pipeline.stages.deploy.steps.env.group.name>",
		},
		{
			name:     "env stage-relative group ref",
			input:    "<+stage.spec.env.envGroupRef>",
			expected: "<+stage.steps.env.group.id>",
		},
		{
			name:     "env spec-relative identifier",
			input:    "<+spec.env.identifier>",
			expected: "<+stage.steps.env.id>",
		},
		{
			name:     "env alias-relative identifier",
			input:    "<+env.identifier>",
			expected: "<+env.id>",
		},
		{
			name:     "env alias-relative group name",
			input:    "<+env.envGroupName>",
			expected: "<+env.group.name>",
		},
		{
			name:     "env FQN unknown subfield passthrough under steps",
			input:    "<+pipeline.stages.deploy.spec.env.variables.MY_VAR>",
			expected: "<+pipeline.stages.deploy.steps.env.variables.MY_VAR>",
		},
		// ---------------- service (identifier/serviceInputs) ------------------
		{
			name:     "service FQN identifier",
			input:    "<+pipeline.stages.deploy.spec.service.identifier>",
			expected: "<+pipeline.stages.deploy.steps.service.id>",
		},
		{
			name:     "service FQN serviceInputs",
			input:    "<+pipeline.stages.deploy.spec.service.serviceInputs>",
			expected: "<+pipeline.stages.deploy.steps.service.with.overlay>",
		},
		{
			name:     "service stage-relative serviceInputs",
			input:    "<+stage.spec.service.serviceInputs>",
			expected: "<+stage.steps.service.with.overlay>",
		},
		{
			name:     "service spec-relative identifier",
			input:    "<+spec.service.identifier>",
			expected: "<+stage.steps.service.id>",
		},
		{
			name:     "service alias-relative identifier",
			input:    "<+service.identifier>",
			expected: "<+service.id>",
		},
		{
			name:     "service alias-relative serviceInputs",
			input:    "<+service.serviceInputs>",
			expected: "<+service.with.overlay>",
		},
		// ---------------- infra (connectorRef) --------------------------------
		{
			name:     "infra FQN connector",
			input:    "<+pipeline.stages.deploy.spec.infra.connectorRef>",
			expected: "<+pipeline.stages.deploy.steps.infra.connector>",
		},
		{
			name:     "infra stage-relative connector",
			input:    "<+stage.spec.infra.connectorRef>",
			expected: "<+stage.steps.infra.connector>",
		},
		{
			name:     "infra spec-relative connector",
			input:    "<+spec.infra.connectorRef>",
			expected: "<+stage.steps.infra.connector>",
		},
		{
			name:     "infra alias-relative connector",
			input:    "<+infra.connectorRef>",
			expected: "<+infra.connector>",
		},
		{
			name:     "infra alias-relative passthrough subfield",
			input:    "<+infra.infraIdentifier>",
			expected: "<+infra.infraIdentifier>",
		},
		// ---------------- infra  -------------------------
		{
			name:     "infra alias-relative infraInputs",
			input:    "<+infra.infraInputs>",
			expected: "<+infra.with.overlay>",
		},
		{
			name:     "infra stage-relative infraInputs",
			input:    "<+stage.spec.infra.infraInputs>",
			expected: "<+stage.steps.infra.with.overlay>",
		},
		{
			name:     "infra alias-relative passthrough subfield",
			input:    "<+infra.spec.environmentRef>",
			expected: "<+infra.spec.environmentRef>",
		},
		// ---------------- manifests (passthrough) -----------------------------
		{
			name:     "manifests FQN passthrough",
			input:    "<+pipeline.stages.deploy.spec.manifests.m1.store.spec.connectorRef>",
			expected: "<+pipeline.stages.deploy.steps.manifests.m1.store.spec.connectorRef>",
		},
		{
			name:     "manifests stage-relative passthrough",
			input:    "<+stage.spec.manifests.m1.valuesPaths>",
			expected: "<+stage.steps.manifests.m1.valuesPaths>",
		},
		{
			name:     "manifests spec-relative passthrough",
			input:    "<+spec.manifests.m1.store>",
			expected: "<+stage.steps.manifests.m1.store>",
		},
		{
			name:     "manifests alias-relative passthrough",
			input:    "<+manifests.m1.store>",
			expected: "<+manifests.m1.store>",
		},
		// ---------------- configFiles (passthrough) ---------------------------
		{
			name:     "configFiles FQN passthrough",
			input:    "<+pipeline.stages.deploy.spec.configFiles.cf1.files>",
			expected: "<+pipeline.stages.deploy.steps.configFiles.cf1.files>",
		},
		{
			name:     "configFiles stage-relative passthrough",
			input:    "<+stage.spec.configFiles.cf1.content>",
			expected: "<+stage.steps.configFiles.cf1.content>",
		},
		{
			name:     "configFiles spec-relative passthrough",
			input:    "<+spec.configFiles.cf1.content>",
			expected: "<+stage.steps.configFiles.cf1.content>",
		},
		{
			name:     "configFiles alias-relative passthrough",
			input:    "<+configFiles.cf1.content>",
			expected: "<+configFiles.cf1.content>",
		},
		// ---------------- artifacts (passthrough) -----------------------------
		{
			name:     "artifacts FQN passthrough",
			input:    "<+pipeline.stages.deploy.spec.artifacts.primary.tag>",
			expected: "<+pipeline.stages.deploy.steps.artifacts.primary.tag>",
		},
		{
			name:     "artifacts stage-relative passthrough",
			input:    "<+stage.spec.artifacts.primary.image>",
			expected: "<+stage.steps.artifacts.primary.image>",
		},
		{
			name:     "artifacts spec-relative passthrough",
			input:    "<+spec.artifacts.primary.tag>",
			expected: "<+stage.steps.artifacts.primary.tag>",
		},
		{
			name:     "artifacts alias-relative passthrough",
			input:    "<+artifacts.primary.tag>",
			expected: "<+artifacts.primary.tag>",
		},
		// ---------------- mixed / nested expressions --------------------------
		{
			name:     "env field with function call (FQN)",
			input:    "<+pipeline.stages.deploy.spec.env.envGroupName.toUpperCase()>",
			expected: "<+pipeline.stages.deploy.steps.env.group.name.toUpperCase()>",
		},
		{
			name:     "mixed service and env (FQN + alias)",
			input:    `svc: <+pipeline.stages.deploy.spec.service.identifier> env: <+env.identifier>`,
			expected: `svc: <+pipeline.stages.deploy.steps.service.id> env: <+env.id>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertExpressionWithTrie(tt.input, nil, false)
			if got != tt.expected {
				t.Errorf("ConvertExpressionWithTrie() deployment stage field failed\ninput:    %s\ngot:      %s\nexpected: %s", tt.input, got, tt.expected)
			}
		})
	}
}

func TestTrieConvert_TemplateInputs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// ---------------- step.template --------------------------------------
		{
			name:     "step template FQN templateInputs",
			input:    "<+pipeline.stages.deploy.spec.execution.steps.s1.template.templateInputs>",
			expected: "<+pipeline.stages.deploy.steps.s1.template.with.overlay>",
		},
		{
			name:     "step template alias-relative templateInputs",
			input:    "<+step.template.templateInputs>",
			expected: "<+step.template.with.overlay>",
		},
		{
			name:     "step template FQN nested subfield passthrough",
			input:    "<+pipeline.stages.deploy.spec.execution.steps.s1.template.templateInputs.spec.command>",
			expected: "<+pipeline.stages.deploy.steps.s1.template.with.overlay.spec.command>",
		},
		// ---------------- stage.template -------------------------------------
		{
			name:     "stage template FQN templateInputs",
			input:    "<+pipeline.stages.deploy.template.templateInputs>",
			expected: "<+pipeline.stages.deploy.template.with.overlay>",
		},
		{
			name:     "stage template stage-relative templateInputs",
			input:    "<+stage.template.templateInputs>",
			expected: "<+stage.template.with.overlay>",
		},
		// ---------------- stepGroup.template ---------------------------------
		{
			name:     "stepGroup template alias-relative templateInputs",
			input:    "<+stepGroup.template.templateInputs>",
			expected: "<+group.template.with.overlay>",
		},
		// ---------------- pipeline.template ----------------------------------
		{
			name:     "pipeline template templateInputs",
			input:    "<+pipeline.template.templateInputs>",
			expected: "<+pipeline.template.with.overlay>",
		},
		// ---------------- standalone alias -----------------------------------
		{
			name:     "template alias-relative templateInputs",
			input:    "<+template.templateInputs>",
			expected: "<+template.with.overlay>",
		},
		{
			name:     "template alias-relative passthrough subfield",
			input:    "<+template.versionLabel>",
			expected: "<+template.versionLabel>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertExpressionWithTrie(tt.input, nil, false)
			if got != tt.expected {
				t.Errorf("ConvertExpressionWithTrie() template inputs failed\ninput:    %s\ngot:      %s\nexpected: %s", tt.input, got, tt.expected)
			}
		})
	}
}
