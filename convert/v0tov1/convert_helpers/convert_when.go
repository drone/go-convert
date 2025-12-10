package converthelpers

import (
	"strings"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
)

// ConvertStepWhen converts v0 step when condition to v1 if expression
func ConvertStepWhen(src *v0.StepWhen) string {
	if src == nil {
		return ""
	}

	var parts []string

	// Convert stageStatus to v1 format
	if src.StageStatus != "" {
		statusExpr := convertStageStatus(src.StageStatus)
		if statusExpr != "" {
			parts = append(parts, statusExpr)
		}
	}

	// Add custom condition if present
	if src.Condition != "" {
		parts = append(parts, "("+src.Condition+")")
	}

	// Combine with &&
	if len(parts) == 0 {
		return ""
	}
	if len(parts) == 1 {
		return parts[0]
	}
	return strings.Join(parts, " && ")
}

// ConvertStageWhen converts v0 stage when condition to v1 if expression
func ConvertStageWhen(src *v0.StageWhen) string {
	if src == nil {
		return ""
	}

	var parts []string

	// Convert pipelineStatus to v1 format
	if src.PipelineStatus != "" {
		statusExpr := convertPipelineStatus(src.PipelineStatus)
		if statusExpr != "" {
			parts = append(parts, statusExpr)
		}
	}

	// Add custom condition if present
	if src.Condition != "" {
		parts = append(parts, "("+src.Condition+")")
	}

	// Combine with &&
	if len(parts) == 0 {
		return ""
	}
	if len(parts) == 1 {
		return parts[0]
	}
	return strings.Join(parts, " && ")
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
