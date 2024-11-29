package json

import (
	"encoding/json"

	harness "github.com/drone/spec/dist/go"
)

// ConvertNotification creates a Harness step for nunit plugin.
func ConvertNotification(node Node, arguments map[string]interface{}) *harness.Step {
	data, _ := arguments["data"].(string)

	// Remove line breaks and format as compact JSON
	var compactData map[string]interface{}
	_ = json.Unmarshal([]byte(data), &compactData) // Parse JSON
	compactBytes, _ := json.Marshal(compactData)   // Convert to compact JSON
	compactString := string(compactBytes)          // Convert bytes to string

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
				"template":     compactString,
			},
		},
	}

	return convertNotification
}
