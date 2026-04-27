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
			expected: "<+pipeline.stages.deploy.env.id>",
		},
		{
			name:     "Stage env group ref",
			input:    "<+pipeline.stages.deploy.spec.env.envGroupRef>",
			expected: "<+pipeline.stages.deploy.env.group.id>",
		},
		{
			name:     "Stage env group name",
			input:    "<+pipeline.stages.deploy.spec.env.envGroupName>",
			expected: "<+pipeline.stages.deploy.env.group.name>",
		},
		{
			name:     "Stage infra connector",
			input:    "<+pipeline.stages.deploy.spec.infra.connectorRef>",
			expected: "<+pipeline.stages.deploy.infra.connector>",
		},
		{
			name:     "Relative stage env identifier",
			input:    "<+stage.spec.env.identifier>",
			expected: "<+stage.env.id>",
		},
		{
			name:     "Relative spec env group ref",
			input:    "<+spec.env.envGroupRef>",
			expected: "<+env.group.id>",
		},
		{
			name:     "Direct env identifier",
			input:    "<+env.identifier>",
			expected: "<+env.id>",
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
			expected: "<+pipeline.stages.deploy.env.group.name.toUpperCase()>",
		},
		{
			name:     "Step with function call",
			input:    "<+pipeline.stages.build.spec.execution.steps.step1.output.outputVariables.var1.toLowerCase()>",
			expected: "<+pipeline.stages.build.steps.step1.output.outputVariables.var1.toLowerCase()>",
		},
		{
			name:     "Mixed text with expressions",
			input:    `stage: <+pipeline.stages.build.identifier> env: <+pipeline.stages.deploy.spec.env.identifier>`,
			expected: `stage: <+pipeline.stages.build.id> env: <+pipeline.stages.deploy.env.id>`,
		},
		{
			name:     "JSON array with expressions",
			input:    `["<+pipeline.stages.build.identifier>", "<+pipeline.stages.deploy.spec.env.envGroupName>"]`,
			expected: `["<+pipeline.stages.build.id>", "<+pipeline.stages.deploy.env.group.name>"]`,
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
			expected: "<+steps.step1.spec.script>",
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
			expected: "<+steps.restore1.steps.restoreCacheS3.spec.with.BUCKET>",
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
			expected: `<+<+pipeline.stages.build.id>=="build"?<+pipeline.stages.deploy.env.id>:"default">`,
		},
		{
			name:     "Ternary with step context",
			input:    `<+<+pipeline.stages.build.spec.execution.steps.run1.spec.command>.contains("test")?<+pipeline.stages.deploy.spec.env.identifier>:"default">`,
			context:  &ConversionContext{StepType: StepTypeRun},
			expected: `<+<+pipeline.stages.build.steps.run1.spec.script>.contains("test")?<+pipeline.stages.deploy.env.id>:"default">`,
		},
		// Function chaining
		{
			name:     "Function chain on env group name",
			input:    `<+<+pipeline.stages.deploy.spec.env.envGroupName>.toLowerCase().contains("prod")>`,
			context:  nil,
			expected: `<+<+pipeline.stages.deploy.env.group.name>.toLowerCase().contains("prod")>`,
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
			expected: `step: <+pipeline.stages.build.steps.step1.id> group: <+pipeline.stages.deploy.env.group.id>`,
		},
		{
			name:     "JSON with mixed conversions",
			input:    `{"stage": "<+pipeline.stages.build.identifier>", "env": "<+pipeline.stages.deploy.spec.env.envGroupName>", "step": "<+pipeline.stages.build.spec.execution.steps.step1.identifier>"}`,
			context:  nil,
			expected: `{"stage": "<+pipeline.stages.build.id>", "env": "<+pipeline.stages.deploy.env.group.name>", "step": "<+pipeline.stages.build.steps.step1.id>"}`,
		},
		// OR expression with multiple conversions
		{
			name:     "OR expression with step output and env",
			input:    `<+ <+<+pipeline.stages.build.spec.execution.steps.step1.identifier>.contains("test")> || <+pipeline.stages.deploy.spec.env.identifier> == "prod" >`,
			context:  nil,
			expected: `<+ <+<+pipeline.stages.build.steps.step1.id>.contains("test")> || <+pipeline.stages.deploy.env.id> == "prod" >`,
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
// The trie captures step IDs during path traversal and resolves step type from StepTypeMap.
// TestTrieConvert_FQNMode tests that when UseFQN is enabled, relative step expressions
// are converted to fully qualified names.
// - "step." prefix uses CurrentStepV1Path and CurrentStepType (the step we're inside)
// - "steps.STEPID" prefix uses StepV1PathMap and StepTypeMap to look up the referenced step
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
				UseFQN:            true,
				CurrentStepV1Path: "pipeline.stages.build.steps.restoreCache",
				CurrentStepType:   StepTypeRestoreCacheGCS,
			},
			expected: "<+pipeline.stages.build.steps.restoreCache.steps.restoreCacheGCS.spec.with.BUCKET>",
		},
		{
			name:  "step.spec.key with FQN - RestoreCacheS3",
			input: "<+step.spec.key>",
			context: &ConversionContext{
				UseFQN:            true,
				CurrentStepV1Path: "pipeline.stages.build.steps.restoreS3",
				CurrentStepType:   StepTypeRestoreCacheS3,
			},
			expected: "<+pipeline.stages.build.steps.restoreS3.steps.restoreCacheS3.spec.with.CACHE_KEY>",
		},
		{
			name:  "step.spec.command with FQN - Run step",
			input: "<+step.spec.command>",
			context: &ConversionContext{
				UseFQN:            true,
				CurrentStepV1Path: "pipeline.stages.build.steps.runStep",
				CurrentStepType:   StepTypeRun,
			},
			expected: "<+pipeline.stages.build.steps.runStep.spec.script>",
		},
		{
			name:  "step.identifier with FQN",
			input: "<+step.identifier>",
			context: &ConversionContext{
				UseFQN:            true,
				CurrentStepV1Path: "pipeline.stages.build.steps.myStep",
				CurrentStepType:   StepTypeRun,
			},
			expected: "<+pipeline.stages.build.steps.myStep.id>",
		},
		// Step inside step group with FQN
		{
			name:  "step.spec.bucket in step group with FQN",
			input: "<+step.spec.bucket>",
			context: &ConversionContext{
				UseFQN:            true,
				CurrentStepV1Path: "pipeline.stages.build.steps.cacheGroup.steps.restoreCache",
				CurrentStepType:   StepTypeRestoreCacheGCS,
			},
			expected: "<+pipeline.stages.build.steps.cacheGroup.steps.restoreCache.steps.restoreCacheGCS.spec.with.BUCKET>",
		},
		// Nested step group with FQN
		{
			name:  "step.spec.bucket in nested step group with FQN",
			input: "<+step.spec.bucket>",
			context: &ConversionContext{
				UseFQN:            true,
				CurrentStepV1Path: "pipeline.stages.build.steps.outerGroup.steps.innerGroup.steps.restoreCache",
				CurrentStepType:   StepTypeRestoreCacheGCS,
			},
			expected: "<+pipeline.stages.build.steps.outerGroup.steps.innerGroup.steps.restoreCache.steps.restoreCacheGCS.spec.with.BUCKET>",
		},
		// HTTP step with FQN
		{
			name:  "step.spec.url with FQN - HTTP step",
			input: "<+step.spec.url>",
			context: &ConversionContext{
				UseFQN:            true,
				CurrentStepV1Path: "pipeline.stages.build.steps.httpCall",
				CurrentStepType:   StepTypeHTTP,
			},
			expected: "<+pipeline.stages.build.steps.httpCall.spec.env.PLUGIN_URL>",
		},
		// HTTP output with FQN
		{
			name:  "step.output.httpResponseCode with FQN - HTTP step",
			input: "<+step.output.httpResponseCode>",
			context: &ConversionContext{
				UseFQN:            true,
				CurrentStepV1Path: "pipeline.stages.build.steps.httpCall",
				CurrentStepType:   StepTypeHTTP,
			},
			expected: "<+pipeline.stages.build.steps.httpCall.steps.httpStep.output.outputVariables.PLUGIN_HTTP_RESPONSE_CODE>",
		},
		// Multiple expressions in same string with FQN
		{
			name:  "multiple step expressions with FQN",
			input: `bucket: <+step.spec.bucket> key: <+step.spec.key>`,
			context: &ConversionContext{
				UseFQN:            true,
				CurrentStepV1Path: "pipeline.stages.build.steps.restoreCache",
				CurrentStepType:   StepTypeRestoreCacheS3,
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
				UseFQN:            false,
				CurrentStepV1Path: "pipeline.stages.build.steps.runStep",
				CurrentStepType:   StepTypeRun,
			},
			expected: "<+step.spec.script>", // Normal relative conversion
		},
		// FQN enabled but no CurrentStepV1Path - should use normal conversion
		{
			name:  "step.spec.command with FQN but no path",
			input: "<+step.spec.command>",
			context: &ConversionContext{
				UseFQN:            true,
				CurrentStepV1Path: "",
				CurrentStepType:   StepTypeRun,
			},
			expected: "<+step.spec.script>", // Normal relative conversion
		},
		// Already FQN path - should convert normally (not relative)
		{
			name:  "already FQN path with FQN enabled",
			input: "<+pipeline.stages.build.spec.execution.steps.step1.spec.command>",
			context: &ConversionContext{
				UseFQN:            true,
				CurrentStepV1Path: "pipeline.stages.build.steps.step1",
				CurrentStepType:   StepTypeRun,
				StepTypeMap:       map[string]string{"step1": StepTypeRun},
				StepV1PathMap:     map[string]string{"step1": "pipeline.stages.build.steps.step1"},
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
				UseFQN:        true,
				StepTypeMap:   map[string]string{"runStep": StepTypeRun},
				StepV1PathMap: map[string]string{"runStep": "pipeline.stages.build.steps.runStep"},
			},
			expected: "<+pipeline.stages.build.steps.runStep.spec.script>",
		},
		{
			name:  "steps.STEPID.spec.bucket with FQN - RestoreCacheGCS",
			input: "<+steps.restoreCache.spec.bucket>",
			context: &ConversionContext{
				UseFQN:        true,
				StepTypeMap:   map[string]string{"restoreCache": StepTypeRestoreCacheGCS},
				StepV1PathMap: map[string]string{"restoreCache": "pipeline.stages.build.steps.restoreCache"},
			},
			expected: "<+pipeline.stages.build.steps.restoreCache.steps.restoreCacheGCS.spec.with.BUCKET>",
		},
		// execution.steps.STEPID.spec.bucket
		{
			name:  "execution.steps.STEPID.spec.bucket with FQN",
			input: "<+execution.steps.restoreCache.spec.bucket>",
			context: &ConversionContext{
				UseFQN:        true,
				StepTypeMap:   map[string]string{"restoreCache": StepTypeRestoreCacheGCS},
				StepV1PathMap: map[string]string{"restoreCache": "pipeline.stages.build.steps.restoreCache"},
			},
			expected: "<+pipeline.stages.build.steps.restoreCache.steps.restoreCacheGCS.spec.with.BUCKET>",
		},
		// spec.execution.steps.STEPID.spec.bucket
		{
			name:  "spec.execution.steps.STEPID.spec.bucket with FQN",
			input: "<+spec.execution.steps.restoreCache.spec.bucket>",
			context: &ConversionContext{
				UseFQN:        true,
				StepTypeMap:   map[string]string{"restoreCache": StepTypeRestoreCacheGCS},
				StepV1PathMap: map[string]string{"restoreCache": "pipeline.stages.build.steps.restoreCache"},
			},
			expected: "<+pipeline.stages.build.steps.restoreCache.steps.restoreCacheGCS.spec.with.BUCKET>",
		},
		// stepGroup.steps.STEPID.spec.bucket
		{
			name:  "stepGroup.steps.STEPID.spec.bucket with FQN",
			input: "<+stepGroup.steps.restoreCache.spec.bucket>",
			context: &ConversionContext{
				UseFQN:        true,
				StepTypeMap:   map[string]string{"restoreCache": StepTypeRestoreCacheGCS},
				StepV1PathMap: map[string]string{"restoreCache": "pipeline.stages.build.steps.cacheGroup.steps.restoreCache"},
			},
			expected: "<+pipeline.stages.build.steps.cacheGroup.steps.restoreCache.steps.restoreCacheGCS.spec.with.BUCKET>",
		},
		// Reference a different step than the current one
		{
			name:  "steps.otherStep reference from inside currentStep",
			input: "<+steps.otherStep.output.result>",
			context: &ConversionContext{
				UseFQN:            true,
				CurrentStepV1Path: "pipeline.stages.build.steps.currentStep",
				CurrentStepType:   StepTypeRun,
				StepTypeMap:       map[string]string{"otherStep": StepTypeRun, "currentStep": StepTypeRun},
				StepV1PathMap:     map[string]string{"otherStep": "pipeline.stages.build.steps.otherStep", "currentStep": "pipeline.stages.build.steps.currentStep"},
			},
			expected: "<+pipeline.stages.build.steps.otherStep.output.result>",
		},
		// steps.STEPID without StepV1PathMap - should fall back to normal conversion
		{
			name:  "steps.STEPID without StepV1PathMap",
			input: "<+steps.unknownStep.spec.command>",
			context: &ConversionContext{
				UseFQN:        true,
				StepTypeMap:   map[string]string{},
				StepV1PathMap: map[string]string{},
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
				StepTypeMap: map[string]string{"runStep1": StepTypeRun},
				// Note: StepType is NOT set - should be resolved lazily
			},
			expected: "<+pipeline.stages.build.steps.runStep1.spec.script>",
		},
		{
			name:  "Lazy resolution from StepTypeMap for Http step",
			input: "<+pipeline.stages.build.spec.execution.steps.httpStep.spec.url>",
			context: &ConversionContext{
				StepTypeMap: map[string]string{"httpStep": StepTypeHTTP},
			},
			expected: "<+pipeline.stages.build.steps.httpStep.spec.env.PLUGIN_URL>",
		},
		{
			name:  "Lazy resolution for nested step group",
			input: "<+pipeline.stages.build.spec.execution.steps.myGroup.steps.innerRun.spec.command>",
			context: &ConversionContext{
				StepTypeMap: map[string]string{"innerRun": StepTypeRun},
			},
			expected: "<+pipeline.stages.build.steps.myGroup.steps.innerRun.spec.script>",
		},
		// CurrentStepType fallback for "step.*" expressions
		{
			name:  "CurrentStepType fallback for step.spec expressions",
			input: "<+step.spec.command>",
			context: &ConversionContext{
				CurrentStepType: StepTypeRun,
				StepTypeMap:     map[string]string{},
			},
			expected: "<+step.spec.script>",
		},
		// Multiple steps with different types in same StepTypeMap
		{
			name:  "Multiple step types - resolves correct type for each",
			input: "<+pipeline.stages.build.spec.execution.steps.gcsStep.spec.env.bucket>",
			context: &ConversionContext{
				StepTypeMap: map[string]string{
					"runStep1": StepTypeRun,
					"gcsStep":  StepTypeGCSUpload,
					"httpStep": StepTypeHTTP,
				},
			},
			expected: "<+pipeline.stages.build.steps.gcsStep.spec.env.bucket>",
		},
		// Step ID not in map - deterministic fallback applies first matching context rule
		{
			name:  "Unknown step ID - deterministic fallback to first matching context",
			input: "<+pipeline.stages.build.spec.execution.steps.unknownStep.spec.command>",
			context: &ConversionContext{
				StepTypeMap: map[string]string{"runStep1": StepTypeRun},
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
				StepTypeMap: map[string]string{"runStep1": StepTypeRun},
			},
			expected: "<+steps.runStep1.spec.script>",
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
