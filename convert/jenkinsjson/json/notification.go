package json

import (
	"fmt"

	harness "github.com/drone/spec/dist/go"
)

// ConvertNotification creates a Harness step for Notification plugin.
func ConvertNotification(node Node, parameterMap map[string]interface{}) *harness.Step {

	// Extract values from parameterMap
	url, urlOk := parameterMap["endpoint"].(string)
	urls := []string{}

	if urlOk {
		urls = append(urls, url)
	} else {
		urls = append(urls, "<+input>")
	}

	phase, _ := parameterMap["phase"].(string)
	notes, _ := parameterMap["notes"].(string)

	// Determine the status based on phase
	status := ""
	if phase == "COMPLETED" {
		status = "SUCCESSFUL"
	}

	// Create the template
	template := fmt.Sprintf(`{
    "status": "%s",
    "notes": "%s"
}`, status, notes)

	convertNotification := &harness.Step{
		Name: "Notification",
		Type: "plugin",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepPlugin{
			Image: "plugins/webhook",
			With: map[string]interface{}{
				"urls":         urls,
				"method":       "POST",
				"token-value":  "<+input>",
				"content-type": "application/json",
				"template":     template,
			},
		},
	}

	return convertNotification
}
