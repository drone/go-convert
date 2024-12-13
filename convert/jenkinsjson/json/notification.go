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

	// Create the template
	template := fmt.Sprintf(`{
  		"phase": "%s",
  		"notes": "%s",
	}`, phase, notes)

	convertNotification := &harness.Step{
		Name: "Notification",
		Type: "plugin",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepPlugin{
			Image: "plugins/webhook",
			With: map[string]interface{}{
				"urls":         urls,
				"method":       "POST",
				"username":     "<+input>",
				"password":     "<+input>",
				"token-value":  "<+input>",
				"token-type":   "<+input>",
				"content-type": "application/json",
				"template":     template,
			},
		},
	}

	return convertNotification
}
