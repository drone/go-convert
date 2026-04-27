package convertexpressions

import (
	"strings"
)

// codebasePrefix is the v0 prefix for codebase expressions
const codebasePrefix = "pipeline.properties.ci.codebase."

// ConvertExpressionWithTrie converts expressions using trie-based matching with context.
// Uses a cached singleton trie for performance.
//
// When context.UseFQN is true and context.CurrentStepPath is set, the trie will
// automatically build FQN paths for step expressions during matching. This happens
// when the trie detects it's inside a step node (after capturing step ID).
func ConvertExpressionWithTrie(expr string, context *ConversionContext, inner bool) string {
	trie := GetPipelineTrie()

	if !inner {
		return replaceHarnessExprs(expr, func(innerContent string) string {
			return ConvertExpressionWithTrie(innerContent, context, true)
		})
	}

	// Process nested <+...> expressions first
	expr = replaceHarnessExprs(expr, func(innerContent string) string {
		return ConvertExpressionWithTrie(innerContent, context, true)
	})

	// Apply trie-based matching to dotted path segments
	// FQN building happens inside the trie when UseFQN is enabled
	return replaceDottedPaths(expr, func(m string) string {
		// Handle codebase path conversion: pipeline.properties.ci.codebase.* -> codebase.*
		if strings.HasPrefix(m, codebasePrefix) {
			m = "codebase." + m[len(codebasePrefix):]
		}

		if converted, matched := trie.Match(m, context); matched {
			return converted
		}
		return m
	})
}
