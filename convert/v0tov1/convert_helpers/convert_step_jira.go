package converthelpers

import (
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// ConvertStepJiraCreate converts a v0 JiraCreate step to a v1 action step
func ConvertStepJiraCreate(src *v0.Step) *v1.StepAction {
	if src == nil || src.Spec == nil {
		return nil
	}

	sp, ok := src.Spec.(*v0.StepJiraCreate)
	if !ok || sp == nil {
		return nil
	}

	// Build fields array as []map[string]interface{}{ {name:..., value:...}, ... }
	fields := make([]map[string]interface{}, 0, len(sp.Fields))
	for _, f := range sp.Fields {
		if f == nil {
			continue
		}
		m := map[string]interface{}{}
		if f.Name != "" {
			m["name"] = f.Name
		}
		// Prefer Value when present
		if f.Value != "" {
			m["value"] = f.Value
		}
		if len(m) > 0 {
			fields = append(fields, m)
		}
	}

	with := map[string]interface{}{
		"connector": sp.ConnectorRef,
		"project":   sp.ProjectKey,
		"type":      sp.IssueType,
	}
	if len(fields) > 0 {
		with["fields"] = fields
	}
	// add timeout inside with as per expected output
	if src.Timeout.String() != "" {
		with["timeout"] = src.Timeout.String()
	}

	return &v1.StepAction{
		Uses: "jira-create",
		With: with,
	}
}

// ConvertStepJiraUpdate converts a v0 JiraUpdate step to a v1 action step
func ConvertStepJiraUpdate(src *v0.Step) *v1.StepAction {
	if src == nil || src.Spec == nil {
		return nil
	}

	sp, ok := src.Spec.(*v0.StepJiraUpdate)
	if !ok || sp == nil {
		return nil
	}

	// fields mapping
	fields := make([]map[string]interface{}, 0, len(sp.Fields))
	for _, f := range sp.Fields {
		if f == nil {
			continue
		}
		m := map[string]interface{}{}
		if f.Name != "" {
			m["name"] = f.Name
		}
		if f.Value != "" {
			m["value"] = f.Value
		}
		if len(m) > 0 {
			fields = append(fields, m)
		}
	}

	with := map[string]interface{}{
		"connector": sp.ConnectorRef,
		"project":   sp.ProjectKey,
		"type":      sp.IssueType,
		"key":       sp.IssueKey,
	}
	if len(fields) > 0 {
		with["fields"] = fields
	}
	if sp.TransitionTo != nil {
		if sp.TransitionTo.Status != "" {
			with["status"] = sp.TransitionTo.Status
		}
		if sp.TransitionTo.TransitionName != "" {
			with["status-transition"] = sp.TransitionTo.TransitionName
		}
	}
	if src.Timeout.String() != "" {
		with["timeout"] = src.Timeout.String()
	}

	return &v1.StepAction{
		Uses: "jira-update",
		With: with,
	}
}
