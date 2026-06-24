package converthelpers

import (
	"strings"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	flexible "github.com/drone/go-convert/internal/flexible"
)

// ConvertStepWhen converts v0 step when condition to v1 if expression.
// A non-empty skipCondition is combined with the when condition (skip means
// "skip if true", so v1 runs if "not skip") and never overwrites it.
func ConvertStepWhen(src *flexible.Field[v0.StepWhen], skip string) string {
	if src == nil {
		return skipExpr(skip)
	}

	var parts []string
	if when, ok := src.AsStruct(); ok {
		// Convert stageStatus to v1 format
		if when.StageStatus != "" {
			statusExpr := convertStageStatus(when.StageStatus)
			if statusExpr != "" {
				parts = append(parts, statusExpr)
			}
		}

		// Add custom condition if present
        if when.Condition != nil && !when.Condition.IsNil() {
            // Check if it's an expression string
            if condExpr, ok := when.Condition.AsString(); ok && condExpr != "" {
                parts = append(parts, "("+condExpr+")")
            } else if condBool, ok := when.Condition.AsStruct(); ok {
                // Convert boolean to string representation
                if condBool {
                    parts = append(parts, "(true)")
                } else {
                    parts = append(parts, "(false)")
                }
            }
        }

        // Combine with &&
		whenStr := ""
		if len(parts) == 1 {
			whenStr = parts[0]
		} else if len(parts) > 1 {
			whenStr = strings.Join(parts, " && ")
		}
		return combineWhenSkip(whenStr, skip)
	} else if when_expression, ok := src.AsString(); ok && when_expression != "<+input>" {
		return combineWhenSkip(when_expression, skip)
	}
	return skipExpr(skip)
}

// ConvertStageWhen converts v0 stage when condition to v1 if expression.
// A non-empty skipCondition is combined with the when condition (skip means
// "skip if true", so v1 runs if "not skip") and never overwrites it.
func ConvertStageWhen(src *flexible.Field[v0.StageWhen], skip string) string {
	if src == nil {
		return skipExpr(skip)
	}

	var parts []string
	if when, ok := src.AsStruct(); ok {
		// Convert pipelineStatus to v1 format
		if when.PipelineStatus != "" {
			statusExpr := convertPipelineStatus(when.PipelineStatus)
			if statusExpr != "" {
				parts = append(parts, statusExpr)
			}
		}

		// Add custom condition if present
        if when.Condition != nil && !when.Condition.IsNil() {
            // Check if it's an expression string
            if condExpr, ok := when.Condition.AsString(); ok && condExpr != "" {
                parts = append(parts, "("+condExpr+")")
            } else if condBool, ok := when.Condition.AsStruct(); ok {
                // Convert boolean to string representation
                if condBool {
                    parts = append(parts, "(true)")
                } else {
                    parts = append(parts, "(false)")
                }
            }
        }

        // Combine with &&
		whenStr := ""
		if len(parts) == 1 {
			whenStr = parts[0]
		} else if len(parts) > 1 {
			whenStr = strings.Join(parts, " && ")
		}
		return combineWhenSkip(whenStr, skip)
	} else if when_expression, ok := src.AsString(); ok && when_expression != "<+input>" {
		return combineWhenSkip(when_expression, skip)
	}
	return skipExpr(skip)
}

// skipExpr converts a v0 skipCondition into a v1 "run if not skipped" expression.
func skipExpr(skip string) string {
	if skip == "" {
		return ""
	}
	return "<+!" + skip + ">"
}

// combineWhenSkip joins a converted when expression with a skipCondition.
// Both are preserved: when && !skip. Either may be empty.
func combineWhenSkip(whenStr, skip string) string {
	skipStr := skipExpr(skip)
	switch {
	case whenStr == "":
		return skipStr
	case skipStr == "":
		return whenStr
	default:
		return whenStr + " && " + skipStr
	}
}

// convertPipelineStatus converts v0 pipelineStatus to v1 expression
func convertPipelineStatus(status string) string {
	switch strings.ToLower(status) {
	case "success":
		return "<+OnPipelineSuccess>"
	case "failure":
		return "<+OnPipelineFailure>"
	case "all":
		return "<+Always>"
	default:
		return ""
	}
}

// convertStageStatus converts v0 stageStatus to v1 expression
func convertStageStatus(status string) string {
	switch strings.ToLower(status) {
	case "success":
		return "<+OnStageSuccess>"
	case "failure":
		return "<+OnStageFailure>"
	case "all":
		return "<+Always>"
	default:
		return ""
	}
}
