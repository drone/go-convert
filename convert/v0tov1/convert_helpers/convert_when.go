package converthelpers

import (
	"strings"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	flexible "github.com/drone/go-convert/internal/flexible"
)

// ConvertStepWhen converts v0 step when condition to v1 if expression
func ConvertStepWhen(src *flexible.Field[v0.StepWhen]) string {
	if src == nil {
		return ""
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
		if when.Condition != "" {
			parts = append(parts, "("+when.Condition+")")
		}

		// Combine with &&
		if len(parts) == 0 {
			return ""
		}
		if len(parts) == 1 {
			return parts[0]
		}
		return strings.Join(parts, " && ")
	} else if when_expression, ok := src.AsString(); ok && when_expression != "<+input>" {
		return when_expression
	}
	return ""
}

// ConvertStageWhen converts v0 stage when condition to v1 if expression
func ConvertStageWhen(src *flexible.Field[v0.StageWhen]) string {
	if src == nil {
		return ""
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
		if when.Condition != "" {
			parts = append(parts, "("+when.Condition+")")
		}

		// Combine with &&
		if len(parts) == 0 {
			return ""
		}
		if len(parts) == 1 {
			return parts[0]
		}
		return strings.Join(parts, " && ")
	} else if when_expression, ok := src.AsString(); ok && when_expression != "<+input>" {
		return when_expression
	}
	return ""
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
