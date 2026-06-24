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

	var fields []map[string]string
	for _, f := range sp.Fields {
		if f == nil {
			continue
		}
		fields = append(fields, map[string]string{"key": f.Name, "value": f.Value})
	}

	// Determine create_ticket_options based on createType.
	// v0 createType (Normal|Form|Standard) maps to template options
	// (Fields|Form Template|Standard Template).
	createTicketOptions := "Fields"
	switch sp.CreateType {
	case "Form":
		createTicketOptions = "Form Template"
	case "Standard":
		createTicketOptions = "Standard Template"
	}

	with := map[string]interface{}{
		"connector":             sp.ConnectorRef,
		"ticket_type":           sp.TicketType,
		"create_ticket_options": createTicketOptions,
	}
	if len(fields) > 0 {
		with["fields"] = fields
	}

	// templateName maps to form_template or standard_template based on createType
	if sp.TemplateName != "" {
		switch sp.CreateType {
		case "Form":
			with["form_template"] = sp.TemplateName
		case "Standard":
			with["standard_template"] = sp.TemplateName
		}
	}

	// log_level: no v0 field; template default is "error"

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

	var fields []map[string]string
	for _, f := range sp.Fields {
		if f == nil {
			continue
		}
		fields = append(fields, map[string]string{"key": f.Name, "value": f.Value})
	}

	// Determine update_ticket_option based on useServiceNowTemplate
	updateTicketOption := "Fields"
	if sp.UseServiceNowTemplate != nil {
		if val, ok := sp.UseServiceNowTemplate.AsStruct(); ok && val {
			updateTicketOption = "Template"
		}
	}

	with := map[string]interface{}{
		"connector":            sp.ConnectorRef,
		"ticket_type":          sp.TicketType,
		"update_ticket_option": updateTicketOption,
	}
	if sp.TicketNumber != "" {
		with["ticket_number"] = sp.TicketNumber
	}
	if len(fields) > 0 {
		with["fields"] = fields
	}

	// templateName maps to template when useServiceNowTemplate is true
	if sp.TemplateName != "" {
		with["template"] = sp.TemplateName
	}

	// updateMultiple maps to update_multiple + change request details
	if sp.UpdateMultiple != nil {
		with["update_multiple"] = true
		if spec := sp.UpdateMultiple.Spec; spec != nil {
			if spec.ChangeRequestNumber != "" {
				with["change_request_number"] = spec.ChangeRequestNumber
			}
			if spec.ChangeTaskType != "" {
				with["change_task_type"] = spec.ChangeTaskType
			}
		}
	}

	// log_level: no v0 field; template default is "error"

	return &v1.StepTemplate{
		Uses: "serviceNowUpdate",
		With: with,
	}
}
