package converthelpers

import (
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// ConvertStepJiraCreate converts a v0 JiraCreate step to a v1 action step
func ConvertStepJiraCreate(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}

	sp, ok := src.Spec.(*v0.StepJiraCreate)
	if !ok || sp == nil {
		return nil
	}

	// fields mapping
	fields := map[string]string{}
	fields["project"] = sp.ProjectKey
	fields["summary"] = ""

	for _, f := range sp.Fields {
		if f == nil {
			continue
		}
		if f.Name == "Summary" {
			fields["summary"] = f.Value
			continue
		}
		fields[f.Name] = f.Value
	}

	with := map[string]interface{}{
		"connector": sp.ConnectorRef,
		"project":   sp.ProjectKey,
		"issue_type": sp.IssueType,
	}

	with["fields"] = fields

	return &v1.StepTemplate{
		Uses: "jiraCreate",
		With: with,
	}	
}

// ConvertStepJiraUpdate converts a v0 JiraUpdate step to a v1 action step
func ConvertStepJiraUpdate(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}

	sp, ok := src.Spec.(*v0.StepJiraUpdate)
	if !ok || sp == nil {
		return nil
	}

	// fields mapping
	fields := map[string]string{}
	for _, f := range sp.Fields {
		if f == nil {
			continue
		}
		fields[f.Name] = f.Value
	}

	with := map[string]interface{}{
		"connector": sp.ConnectorRef,
		"issue":       sp.IssueKey,
	}

	with["fields"] = fields

	if sp.TransitionTo != nil {
		if sp.TransitionTo.Status != "" {
			with["status"] = sp.TransitionTo.Status
		}
		if sp.TransitionTo.TransitionName != "" {
			with["transition"] = sp.TransitionTo.TransitionName
		}
	}

	return &v1.StepTemplate{
		Uses: "jiraUpdate",
		With: with,
	}
}
