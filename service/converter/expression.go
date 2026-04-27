package converter

import (
	"strings"

	convertexpressions "github.com/drone/go-convert/convert/convertexpressions"
)

// ExpressionContext holds the context needed for expression conversion.
// This mirrors the ConversionContext from the convertexpressions package
// but is exposed as a public API type for the service layer.
type ExpressionContext struct {
	// CurrentStepID is the ID of the current step we're inside (optional)
	CurrentStepID string `json:"current_step_id,omitempty"`

	// CurrentStepType is the type of the step we're currently inside (e.g., "Run", "Action", "Plugin")
	// Used when expression is "step.spec.*" (no explicit step ID)
	CurrentStepType string `json:"current_step_type,omitempty"`

	// CurrentStepV1Path is the v1 FQN base path to the current step
	// Example: "pipeline.stages.build.steps.restoreCache"
	CurrentStepV1Path string `json:"current_step_v1_path,omitempty"`

	// StepTypeMap maps step ID to step type for all steps in the pipeline
	// Example: {"step1": "Run", "step2": "Action"}
	StepTypeMap map[string]string `json:"step_type_map,omitempty"`

	// StepV1PathMap maps step ID to its v1 FQN base path
	// Example: {"restoreCache": "pipeline.stages.build.steps.restoreCache"}
	StepV1PathMap map[string]string `json:"step_v1_path_map,omitempty"`

	// UseFQN enables fully qualified name mode for step expressions
	// When true, step expressions are converted to their full v1 paths
	UseFQN bool `json:"use_fqn,omitempty"`
}

// ConvertExpression converts a single Harness v0 expression to v1 format.
// The expression should include the <+...> delimiters.
// Returns the converted expression string.
func ConvertExpression(expression string, ctx *ExpressionContext) string {
	if expression == "" {
		return ""
	}

	// Check if the expression contains Harness expression markers
	if !strings.Contains(expression, "<+") {
		return expression
	}

	// Build the internal ConversionContext
	convCtx := buildConversionContext(ctx)

	// Use the trie-based expression converter
	return convertexpressions.ConvertExpressionWithTrie(expression, convCtx, false)
}

// ConvertExpressions converts multiple Harness v0 expressions to v1 format.
// Returns a map of original expression to converted expression.
func ConvertExpressions(expressions []string, ctx *ExpressionContext) map[string]string {
	result := make(map[string]string, len(expressions))
	convCtx := buildConversionContext(ctx)

	for _, expr := range expressions {
		if expr == "" {
			result[expr] = ""
			continue
		}
		if !strings.Contains(expr, "<+") {
			result[expr] = expr
			continue
		}
		result[expr] = convertexpressions.ConvertExpressionWithTrie(expr, convCtx, false)
	}

	return result
}

// buildConversionContext converts the public ExpressionContext to the internal ConversionContext
func buildConversionContext(ctx *ExpressionContext) *convertexpressions.ConversionContext {
	if ctx == nil {
		return &convertexpressions.ConversionContext{}
	}

	return &convertexpressions.ConversionContext{
		StepID:            ctx.CurrentStepID,
		CurrentStepType:   ctx.CurrentStepType,
		CurrentStepV1Path: ctx.CurrentStepV1Path,
		StepTypeMap:       ctx.StepTypeMap,
		StepV1PathMap:     ctx.StepV1PathMap,
		UseFQN:            ctx.UseFQN,
	}
}
