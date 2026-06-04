package converter

import (
	"strings"

	convertexpressions "github.com/drone/go-convert/convert/convertexpressions"
)

// ExpressionContext is the FQN-only conversion context: a v1 pipeline YAML
// plus an optional call-site FQN.
type ExpressionContext struct {
	// ContextPipelineYAML is an optional v1 pipeline YAML. When set, it builds
	// the FQN-keyed step lookup (StepInfoByFQN) and enables FQN mode.
	ContextPipelineYAML string `json:"context_pipeline_yaml,omitempty"`

	// CurrentFQN is the v1 FQN of the step the expression lives in (e.g.
	// "pipeline.stages.prod.steps.G1.steps.deploy"). The stage, step-group
	// chain, and current step type are derived from it.
	CurrentFQN string `json:"current_fqn,omitempty"`
}

// ConvertExpression converts a single Harness v0 expression to v1 format.
// The expression should include the <+...> / ${{...}} delimiters.
func ConvertExpression(expression string, ctx *ExpressionContext) string {
	out, _ := ConvertExpressionWithWarnings(expression, ctx)
	return out
}

// ConvertExpressionWithWarnings is ConvertExpression plus any diagnostics
// collected during conversion (e.g. ambiguous step types resolved via
// best-match fallback, unmapped template/approval uses).
func ConvertExpressionWithWarnings(expression string, ctx *ExpressionContext) (string, []string) {
	if expression == "" {
		return "", nil
	}
	if !hasExpressionDelimiter(expression) {
		return expression, nil
	}
	convCtx := buildConversionContext(ctx)
	out := convertexpressions.ConvertExpressionWithTrie(expression, convCtx, false)
	return out, convCtx.Warnings
}

// ConvertExpressions converts multiple Harness v0 expressions to v1 format.
// Returns a map of original expression to converted expression.
func ConvertExpressions(expressions []string, ctx *ExpressionContext) map[string]string {
	out, _ := ConvertExpressionsWithWarnings(expressions, ctx)
	return out
}

// ConvertExpressionsWithWarnings is ConvertExpressions plus the aggregated
// diagnostics collected across all expressions. Each warning embeds the
// offending path, so a flat list stays attributable.
func ConvertExpressionsWithWarnings(expressions []string, ctx *ExpressionContext) (map[string]string, []string) {
	result := make(map[string]string, len(expressions))
	convCtx := buildConversionContext(ctx)

	for _, expr := range expressions {
		if expr == "" {
			result[expr] = ""
			continue
		}
		if !hasExpressionDelimiter(expr) {
			result[expr] = expr
			continue
		}
		result[expr] = convertexpressions.ConvertExpressionWithTrie(expr, convCtx, false)
	}

	return result, convCtx.Warnings
}

// hasExpressionDelimiter reports whether s contains a Harness expression
// delimiter (<+...> or ${{...}}).
func hasExpressionDelimiter(s string) bool {
	return strings.Contains(s, "<+") || strings.Contains(s, "${{")
}

// ConvertExpressionWithPipeline converts a single expression using context
// automatically derived from a v1 pipeline YAML (see BuildContextFromPipeline).
func ConvertExpressionWithPipeline(expression string, pipelineYAML string) string {
	ctx := BuildContextFromPipeline(pipelineYAML)
	return ConvertExpression(expression, ctx)
}

// ConvertExpressionsWithPipeline converts multiple expressions using context
// automatically derived from a v1 pipeline YAML.
func ConvertExpressionsWithPipeline(expressions []string, pipelineYAML string) map[string]string {
	ctx := BuildContextFromPipeline(pipelineYAML)
	return ConvertExpressions(expressions, ctx)
}

// BuildContextFromPipeline returns an ExpressionContext resolving step
// references against the v1 pipelineYAML. No call-site is supplied, so set
// ExpressionContext.CurrentFQN to scope a reference to its enclosing group.
func BuildContextFromPipeline(pipelineYAML string) *ExpressionContext {
	return &ExpressionContext{ContextPipelineYAML: pipelineYAML}
}

// buildConversionContext maps the public ExpressionContext to the internal
// ConversionContext, building StepInfoByFQN from the v1 YAML and deriving the
// call-site (stage / step-group chain / step type) from CurrentFQN.
func buildConversionContext(ctx *ExpressionContext) *convertexpressions.ConversionContext {
	if ctx == nil {
		return &convertexpressions.ConversionContext{}
	}

	cc := &convertexpressions.ConversionContext{}

	// FQN lookup: build StepInfoByFQN from the v1 context pipeline.
	if strings.TrimSpace(ctx.ContextPipelineYAML) != "" {
		if m, ok := buildContextFromPipelineYAML(ctx.ContextPipelineYAML); ok {
			cc.StepInfoByFQN = m
			cc.UseFQN = true
		}
	}

	// Derive call-site context from the current step's v1 FQN.
	if ctx.CurrentFQN != "" {
		if stage, chain, _, ok := convertexpressions.ParseFQN(ctx.CurrentFQN); ok {
			cc.CurrentStageID = stage
			cc.CurrentStepGroupChain = chain
			cc.CurrentFQN = ctx.CurrentFQN
			if info, found := cc.StepInfoByFQN[ctx.CurrentFQN]; found && info != nil {
				cc.CurrentStepType = info.Type
				cc.CurrentStepTypes = info.Types
			}
		}
	}

	return cc
}
