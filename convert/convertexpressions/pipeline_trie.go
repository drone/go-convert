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

	// Build main pipeline structure with proper hierarchy
	trie.AddPath().
		Node("pipeline").WithAlias("pipeline").WithV1Name("pipeline").
		Node("stages").WithAlias("stages").WithV1Name("stages").
		Node("*").WithAlias("stage").WithID("stage_node").WithV1Name("*").
		Node("spec").WithAlias("spec").WithV1Name("-").WithID("stage_spec_node").
		Node("execution").WithAlias("execution").WithV1Name("-").WithID("stage_execution_node").
		Node("steps").WithAlias("steps").WithV1Name("steps").WithID("steps_node").
		Node("*").WithAlias("step").WithV1Name("*").WithID("step_node").
		Node("spec").WithAlias("spec").WithV1Name("-").WithID("step_spec_node")

	trie.AddPath().
		Node("codebase").WithAlias("codebase").WithID("codebase_node").WithV1Name("codebase")

	trie.AddPathFromID("step_node").
		Node("output").WithV1Name("-").WithID("step_output_node")

	trie.AddPathFromID("step_node").
		LinkToNodeByID("steps", "steps_node")

	// Build stepGroup structure
	trie.AddPath().
		Node("stepGroup").WithAlias("stepGroup").WithV1Name("group").WithID("step_group_node").
		LinkToNodeByID("steps", "steps_node")

	// Add getParentStepGroup edge that loops back to stepGroup
	trie.AddPathFromID("step_group_node").
		LinkToNodeByID("getParentStepGroup", "step_group_node") // Self-reference

	trie.AddPathFromID("stage_execution_node").
		Node("rollbackSteps").WithAlias("rollbackSteps").WithV1Name("rollback").WithID("rollback_node").
		LinkToNodeByID("*", "step_node")

	// Table-driven general rule registration: nodeID → list of rule sets.
	generalRules := map[string][][]ConversionRule{
		"step_node": {
			StepsConversionRules,
			FailureStrategiesConversionRules,
			NotificationRulesConversionRules,
		},
		"stage_node": {
			StageConversionRules,
			FailureStrategiesConversionRules,
			NotificationRulesConversionRules,
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
