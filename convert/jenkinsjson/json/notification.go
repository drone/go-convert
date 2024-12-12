package json

import (
	"fmt"

	harness "github.com/drone/spec/dist/go"
)

// ConvertNotification creates a Harness step for Notification plugin.
func ConvertNotification(node Node, parameterMap map[string]interface{}) *harness.Step {

	// Extract values from parameterMap
	endpoints, ok := parameterMap["endpoints"].([]interface{})
	urls := []string{}
	contentType := "application/json" // Default content type

	if ok && len(endpoints) > 0 {
		// Assuming only one endpoint in the array for simplicity
		if endpoint, ok := endpoints[0].(map[string]interface{}); ok {
			if url, ok := endpoint["url"].(string); ok {
				urls = append(urls, url)
			}

			if format, ok := endpoint["format"].(string); ok {
				switch format {
				case "JSON":
					contentType = "application/json"
				case "XML":
					contentType = "application/xml"
				default:
					contentType = "application/json" // Fallback
				}
			}
		}
	}

	// If no URLs were found, use <+input>
	if len(urls) == 0 {
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
				"content-type": contentType,
				"template":     template,
			},
		},
	}

	return convertNotification
}
