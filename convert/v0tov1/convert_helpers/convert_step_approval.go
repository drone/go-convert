package converthelpers

import (
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"fmt"
	"strings"
)

func ConvertStepCustomApproval(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}

	// Type assert the spec to StepCustomApproval
	spec, ok := src.Spec.(*v0.StepCustomApproval)
	if !ok {
		return nil
	}

	dst := &v1.StepTemplate{
		Uses: "customApproval",
		With: make(map[string]interface{}),
	}

	with := dst.With.(map[string]interface{})

	if spec.Source != nil && spec.Source.Spec.Script != "" {
		with["script"] = spec.Source.Spec.Script
	}
	if spec.Shell != "" {
		with["shell"] = strings.ToLower(spec.Shell)
	}
	if spec.ScriptTimeout != "" {
		with["script_timeout"] = spec.ScriptTimeout
	}
	if spec.RetryInterval != "" {
		with["retry"] = spec.RetryInterval
	}

	if approve := convertCriteria(spec.ApprovalCriteria); approve != nil {
		with["approve"] = approve
	}
	if reject := convertCriteria(spec.RejectionCriteria); reject != nil {
		with["reject"] = reject
	}

	if outputs := ConvertOutputVariables(spec.OutputVariables); len(outputs) > 0 {
		with["output_variables"] = outputs
	}

	env_map := make(map[string]interface{})
	for _, envVar := range spec.EnvironmentVariables {
		if envVar == nil || envVar.Name == "" || envVar.Value == "" {
			continue
		}
		env_map[envVar.Name] = envVar.Value
	}
	if len(env_map) > 0 {
		with["environment_variables"] = env_map
	}

	return dst
}

func ConvertStepJiraApproval(src *v0.Step) *v1.StepApproval {
	if src == nil || src.Spec == nil {
		return nil
	}

	// Type assert the spec to StepJiraApproval
	spec, ok := src.Spec.(*v0.StepJiraApproval)
	if !ok {
		return nil
	}

	dst := &v1.StepApproval{
		Uses: "jira",
		With: make(map[string]interface{}),
	}
    if spec.RetryInterval != "" {
        dst.With["retry"] = spec.RetryInterval
    }
    approve := convertCriteria(spec.ApprovalCriteria)
    reject := convertCriteria(spec.RejectionCriteria)
    if approve != nil {
        dst.With["approve"] = approve
    }
    if reject != nil {
        dst.With["reject"] = reject
    }

    download := map[string]interface{}{
        "url": "https://storage.googleapis.com/unified-plugins/jira-approval/v0.0.1/",
        "target": "$PLUGIN_PATH",
    }
    env := map[string]interface{}{
        "PLUGIN_HARNESS_CONNECTOR": spec.ConnectorRef,
        "PLUGIN_LOG_LEVEL": "error",
        "PLUGIN_ISSUE_KEY": spec.IssueKey,
    }

    dst.With["run"] = map[string]interface{}{
        "download": download,
        "script": "$PLUGIN_PATH",
        "env": env,

    }

	return dst
}

func ConvertStepServiceNowApproval(src *v0.Step) *v1.StepApproval {
	if src == nil || src.Spec == nil {
		return nil
	}

	// Type assert the spec to StepServiceNowApproval
	spec, ok := src.Spec.(*v0.StepServiceNowApproval)
	if !ok {
		return nil
	}

	dst := &v1.StepApproval{
		Uses: "servicenow",
		With: make(map[string]interface{}),
	}
    if spec.RetryInterval != "" {
        dst.With["retry"] = spec.RetryInterval
    }
    approve := convertCriteria(spec.ApprovalCriteria)
    reject := convertCriteria(spec.RejectionCriteria)
    if approve != nil {
        dst.With["approve"] = approve
    }
    if reject != nil {
        dst.With["reject"] = reject
    }
    if spec.ChangeWindow != nil {
        changeWindow := map[string]interface{}{
            "start": spec.ChangeWindow.StartField,
            "end": spec.ChangeWindow.EndField,
        }
        dst.With["change-window"] = changeWindow
    }
    download := map[string]interface{}{
        "source": "https://storage.googleapis.com/unified-plugins/servicenow-approval/v0.0.1/",
        "target": "$PLUGIN_PATH",
    }
    env := map[string]interface{}{
        "PLUGIN_HARNESS_CONNECTOR": spec.ConnectorRef,
        "PLUGIN_LOG_LEVEL": "error",
        "PLUGIN_TICKET_TYPE": spec.TicketType,
        "PLUGIN_TICKET_NUMBER": spec.TicketNumber,
    }

    dst.With["run"] = map[string]interface{}{
        "download": download,
        "script": "$PLUGIN_PATH",
        "env": env,

    }

	return dst
}

// convertCriteria converts v0 Criteria to v1 format
func convertCriteria(criteria *v0.Criteria) interface{} {
    if criteria == nil || criteria.Spec == nil {
        return nil
    }

    // Handle Jexl type - return expression string directly
    if criteria.Type == "Jexl" && criteria.Spec.Expression != "" {
        return criteria.Spec.Expression
    }

    // Handle KeyValues type - convert to conditions structure
    if criteria.Type == "KeyValues" {
        result := make(map[string]interface{})

        result["match-any-condition"] = criteria.Spec.MatchAnyCondition

        // Convert conditions
        conditions := make([]map[string]interface{}, 0)
        for _, condition := range criteria.Spec.Conditions {
            if condition == nil {
                continue
            }

            conditionMap := convertCondition(condition)
            if conditionMap != nil {
                conditions = append(conditions, conditionMap)
            }
        }

        result["conditions"] = conditions

        return result
    }

    return nil
}

// convertCondition converts a single v0 Condition to v1 format
func convertCondition(condition *v0.Condition) map[string]interface{} {
    if condition == nil {
        return nil
    }

    // Get the key as string
    key := fmt.Sprintf("%v", condition.Key)
    if key == "" {
        return nil
    }

    // Build the condition based on operator
    operator := strings.ToLower(condition.Operator)
    value := condition.Value

    var operatorMap map[string]interface{}

    switch operator {
    case "equals":
        // Simple equality: key: { eq: value }
        operatorMap = map[string]interface{}{
            "eq": value,
        }

    case "not equals":
        // Not equals: key: { not: { eq: value } }
        operatorMap = map[string]interface{}{
            "not": map[string]interface{}{
                "eq": value,
            },
        }

    case "in":
        // In operator: key: { in: [values] }
        // Value might be a comma-separated string or already a slice
        var inValues interface{}
        if valueStr, ok := value.(string); ok {
            // Split comma-separated values
            parts := strings.Split(valueStr, ",")
            trimmedParts := make([]string, 0, len(parts))
            for _, part := range parts {
                trimmedParts = append(trimmedParts, strings.TrimSpace(part))
            }
            inValues = trimmedParts
        } else {
            inValues = value
        }
        operatorMap = map[string]interface{}{
            "in": inValues,
        }

    case "not in":
        // Not in operator: key: { not: { in: [values] } }
        var inValues interface{}
        if valueStr, ok := value.(string); ok {
            // Split comma-separated values
            parts := strings.Split(valueStr, ",")
            trimmedParts := make([]string, 0, len(parts))
            for _, part := range parts {
                trimmedParts = append(trimmedParts, strings.TrimSpace(part))
            }
            inValues = trimmedParts
        } else {
            inValues = value
        }
        operatorMap = map[string]interface{}{
            "not": map[string]interface{}{
                "in": inValues,
            },
        }

    default:
        // Unknown operator, default to equals
        operatorMap = map[string]interface{}{
            "eq": value,
        }
    }

    // Return the condition as { key: { operator: value } }
    return map[string]interface{}{
        key: operatorMap,
    }
}