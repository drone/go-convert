package converthelpers

import (
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

type EmailStepWith struct {
	EmailIds string `json:"email_ids,omitempty"`
	CCEmailIds string `json:"cc_email_ids,omitempty"`
	Subject string `json:"subject,omitempty"`
	Body string `json:"body,omitempty"`
	ToUserGroups string `json:"to_user_groups,omitempty"`
	CcUserGroups string `json:"cc_user_groups,omitempty"`
}

func ConvertStepEmail(src *v0.Step) *v1.StepTemplate {
	if src == nil || src.Spec == nil {
		return nil
	}

	// Extract the typed spec
	spec, ok := src.Spec.(*v0.StepEmail)
	if !ok {
		return nil
	}

	with := EmailStepWith{
		EmailIds: spec.To,
		CCEmailIds: spec.Cc,
		Subject: spec.Subject,
		Body: spec.Body,
	}
	if len(spec.ToUserGroups) > 0 {
		toUgs := ""
		for _, ug := range spec.ToUserGroups {
			toUgs += ug + ","
		}
		with.ToUserGroups = toUgs[:len(toUgs)-1]
	}
	if len(spec.CcUserGroups) > 0 {
		ccUgs := ""
		for _, ug := range spec.CcUserGroups {
			ccUgs += ug + ","
		}
		with.CcUserGroups = ccUgs[:len(ccUgs)-1]
	}

	return &v1.StepTemplate{
		Uses: "email",
		With: with,
	}
}
	