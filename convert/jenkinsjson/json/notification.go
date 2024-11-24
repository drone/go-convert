package json

import (
	harness "github.com/drone/spec/dist/go"
)

// ConvertNotification creates a Harness step for nunit plugin.
func ConvertNotification(node Node, arguments map[string]interface{}) *harness.Step {
	data, _ := arguments["data"].(string)

	convertNotification := &harness.Step{
		Name: "Notification",
		Type: "plugin",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepPlugin{
			Connector: "<+input>",
			Image:     "plugins/webhook",
			With: map[string]interface{}{
				"urls":         "<+input>",
				"username":     "<+input>",
				"password":     "<+input>",
				"method":       "<+input>",
				"content_type": "application/json",
				"debug":        "true",
				"template":     data,
			},
		},
	}

	return convertNotification
}
