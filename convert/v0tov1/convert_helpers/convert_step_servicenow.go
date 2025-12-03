package converthelpers

import (
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// ConvertStepServiceNowCreate converts a v0 ServiceNowCreate step to a v1 template step
func ConvertStepServiceNowCreate(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}

	sp, ok := src.Spec.(*v0.StepServiceNowCreate)
	if !ok || sp == nil {
		return nil
	}

	// Build fields array as []map[string]interface{}{ {key:..., value:...}, ... }
	fields := make([]map[string]interface{}, 0, len(sp.Fields))
	for _, f := range sp.Fields {
		if f == nil {
			continue
		}
		m := map[string]interface{}{
			"key":   f.Name,
			"value": f.Value,
		}
		fields = append(fields, m)
	}

	// Determine create_ticket_options based on createType
	createTicketOptions := "Fields"
	if sp.CreateType != "" && sp.CreateType != "Normal" {
		createTicketOptions = sp.CreateType
	}

	with := map[string]interface{}{
		"connector":             sp.ConnectorRef,
		"ticket_type":           sp.TicketType,
		"create_ticket_options": createTicketOptions,
	}
	if len(fields) > 0 {
		with["fields"] = fields
	}

	return &v1.StepTemplate{
		Uses: "serviceNowCreate",
		With: with,
	}
}

// ConvertStepServiceNowUpdate converts a v0 ServiceNowUpdate step to a v1 template step
func ConvertStepServiceNowUpdate(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}

	sp, ok := src.Spec.(*v0.StepServiceNowUpdate)
	if !ok || sp == nil {
		return nil
	}

	// Build fields array as []map[string]interface{}{ {key:..., value:...}, ... }
	fields := make([]map[string]interface{}, 0, len(sp.Fields))
	for _, f := range sp.Fields {
		if f == nil {
			continue
		}
		m := map[string]interface{}{
			"key":   f.Name,
			"value": f.Value,
		}
		fields = append(fields, m)
	}

	// Determine update_ticket_option based on useServiceNowTemplate
	updateTicketOption := "Fields"
	if sp.UseServiceNowTemplate {
		updateTicketOption = "Template"
	}

	with := map[string]interface{}{
		"connector":            sp.ConnectorRef,
		"ticket_type":          sp.TicketType,
		"ticket_number":        sp.TicketNumber,
		"update_ticket_option": updateTicketOption,
	}
	if len(fields) > 0 {
		with["fields"] = fields
	}

	return &v1.StepTemplate{
		Uses: "serviceNowUpdate",
		With: with,
	}
}
