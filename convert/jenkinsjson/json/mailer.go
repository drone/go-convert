package json

import (
	harness "github.com/drone/spec/dist/go"
)

// ConvertMailer creates a Harness step for nunit plugin.
func ConvertMailer(node Node, arguments map[string]interface{}) *harness.Step {
	subject, _ := arguments["subject"].(string)
	to, _ := arguments["to"].(string)
	body, _ := arguments["body"].(string)

	convertMailer := &harness.Step{
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Name: "Mailer",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "plugins/email",
			With: map[string]interface{}{
				"host":         "<+input>",
				"port":         "<+input>",
				"username":     "<+input>",
				"password":     "<+input>",
				"subject":      subject,
				"body":         body,
				"recipients":   to,
				"from.address": "<+input>",
			},
		},
	}

	return convertMailer
}
