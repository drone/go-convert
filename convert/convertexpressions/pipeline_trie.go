package convertexpressions

import "sync"

var (
	pipelineTrieOnce sync.Once
	pipelineTrie     *Trie
)

// GetPipelineTrie returns the singleton pipeline trie instance.
// It is built once and cached for all subsequent calls.
func GetPipelineTrie() *Trie {
	pipelineTrieOnce.Do(func() {
		pipelineTrie = buildPipelineTrie()
	})
	return pipelineTrie
}

func buildPipelineTrie() *Trie {
	trie := NewTrie()

	// Build main pipeline structure with proper hierarchy. The stage node, its
	// ancestors (pipeline./stages.), and the execution roots (spec./execution.)
	// are flagged WithNoFQNOverride(): they already carry stage context, so FQN
	// expansion is skipped and the remainder converts structurally.
	trie.AddPath().
		Node("pipeline").WithAlias("pipeline").WithV1Name("pipeline").WithID("pipeline_node").WithNoFQNOverride().
		Node("stages").WithAlias("stages").WithV1Name("stages").WithID("stages_node").WithNoFQNOverride().
		Node("*").WithAlias("stage").WithID("stage_node").WithV1Name("*").WithNoFQNOverride().
		Node("spec").WithAlias("spec").WithV1Name("-").WithID("stage_spec_node").WithNoFQNOverride().
		Node("execution").WithAlias("execution").WithV1Name("-").WithID("stage_execution_node").WithNoFQNOverride().
		Node("steps").WithAlias("steps").WithV1Name("steps").WithID("steps_node").
		Node("*").WithAlias("step").WithV1Name("*").WithID("step_node").
		Node("spec").WithAlias("spec").WithV1Name("-").WithID("step_spec_node")

	// codebase expressions to seperate path,
	// pipeline.properties.ci.codebase... expressions will get cut down to codebase.
	trie.AddPath().
		Node("codebase").WithAlias("codebase").WithID("codebase_node").WithV1Name("codebase")

	// inputSet expressions rename their root to inputs.overlay, preserving the
	// remainder structurally: inputSet.pipeline... -> inputs.overlay.pipeline...
	trie.AddPath().
		Node("inputSet").WithAlias("inputSet").WithID("input_set_node").WithV1Name("inputSet.overlay")

	trie.AddPathFromID("step_node").
		Node("output").WithV1Name("-").WithID("step_output_node")

	trie.AddPathFromID("step_node").
		LinkToNodeByID("steps", "steps_node")

	// Build stepGroup structure
	trie.AddPath().
		Node("stepGroup").WithAlias("stepGroup").WithV1Name("group").WithID("step_group_node").
		LinkToNodeByID("steps", "steps_node")

	// getParentStepGroup edge: a dedicated intermediate node so the v1 output
	// keeps the edge name "getParentStepGroup" rather than the target's "group".
	trie.AddPathFromID("step_group_node").
		Node("getParentStepGroup").WithV1Name("getParentStepGroup").WithID("parent_step_group_node").
		LinkToNodeByID("steps", "steps_node")

	// Allow chaining: getParentStepGroup.getParentStepGroup...
	trie.AddPathFromID("parent_step_group_node").
		LinkToNodeByID("getParentStepGroup", "parent_step_group_node")

	trie.AddPathFromID("stage_execution_node").
		Node("rollbackSteps").WithAlias("rollbackSteps").WithV1Name("rollback").WithID("rollback_node").
		LinkToNodeByID("*", "step_node")

	// Deployment stage spec fields (service/manifests/configFiles/infra/env/
	// artifacts) move under stage.steps in v1. Each field gets two entry points
	// that share the same node-relative field rules (see stage.go):
	//   - a spec-child node under stage_spec_node with v1Name "steps.<field>"
	//     (full / spec-prefixed paths: stage.spec.env.* -> stage.steps.env.*)
	//   - a standalone alias root with v1Name "<field>"
	//     (bare relative forms: env.* stay at root)
	// The partial-match fallback handles unknown subfields, so no wildcard child
	// is needed. infra (deployment) and CI's infrastructure/runtime use distinct
	// node names, so there is no clash.
	deploymentStageFields := []string{
		"service", "manifests", "configFiles", "infra", "env", "artifacts",
	}
	for _, field := range deploymentStageFields {
		trie.AddPathFromID("stage_spec_node").
			Node(field).WithV1Name("steps." + field).WithID(field + "_steps_node").WithNoFQNOverride()
		trie.AddPath().
			Node(field).WithAlias(field).WithV1Name(field).WithID(field + "_alias_node").WithNoFQNOverride()
	}

	// Field rule sets attached to BOTH the spec-child and alias nodes per field.
	deploymentFieldRules := map[string][]ConversionRule{
		"env":     EnvFieldRules,
		"service": ServiceFieldRules,
		"infra":   InfraFieldRules,
	}
	for field, rules := range deploymentFieldRules {
		trie.AttachRulesAt(field+"_steps_node", rules)
		trie.AttachRulesAt(field+"_alias_node", rules)
	}

	// Template references (template.templateInputs -> template.with.overlay) can
	// appear at the pipeline, stage, stepGroup and step levels. Each parent gets a
	// "template" child carrying TemplateFieldRules, so both full/relative forms
	// convert (e.g. step.template.*, stage.template.*, and the FQN form
	// pipeline.stages.S...steps.X.template.*). template never targets spec/output,
	// so FQN step expansion does not apply and the remainder converts structurally.
	templateParentIDs := []string{"pipeline_node", "stage_node", "step_group_node", "step_node"}
	for _, parentID := range templateParentIDs {
		templateNodeID := parentID + "_template_node"
		trie.AddPathFromID(parentID).
			Node("template").WithV1Name("template").WithID(templateNodeID).WithNoFQNOverride()
		trie.AttachRulesAt(templateNodeID, TemplateFieldRules)
	}

	// Standalone alias entry so bare relative forms (template.templateInputs) match.
	trie.AddPath().
		Node("template").WithAlias("template").WithV1Name("template").WithID("template_alias_node").WithNoFQNOverride()
	trie.AttachRulesAt("template_alias_node", TemplateFieldRules)

	// Table-driven general rule registration: nodeID → list of rule sets.
	generalRules := map[string][][]ConversionRule{
		"step_node": {
			StepsConversionRules,
			FailureStrategiesConversionRules,
		},
		"stage_node": {
			StageConversionRules,
			FailureStrategiesConversionRules,
		},
		"stage_spec_node": {
			DeploymentStageSpecConversionRules,
			CIStageSpecConversionRules,
		},
		"step_group_node": {
			StepsConversionRules,
		},
		"codebase_node": {
			CodebaseConversionRules,
		},
		"pipeline_node": {
			PipelineConversionRules,
			NotificationRulesConversionRules,
		},
	}
	for nodeID, ruleSets := range generalRules {
		for _, rules := range ruleSets {
			trie.AttachRulesAt(nodeID, rules)
		}
	}

	// Table-driven context-aware rule registration: nodeID → (stepType → rules).
	contextRules := map[string]map[string][]ConversionRule{
		"step_spec_node":   StepSpecConversionRules,
		"step_output_node": StepOutputConversionRules,
	}
	for nodeID, ruleMap := range contextRules {
		for stepType, rules := range ruleMap {
			trie.AttachRulesWithContextAt(nodeID, stepType, rules)
		}
	}

	return trie
}
