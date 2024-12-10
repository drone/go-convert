package json

import (
	"fmt"
	"strings"

	harness "github.com/drone/spec/dist/go"
)

// ConvertNotification creates a Harness step for nunit plugin.
func ConvertNotification(node Node, parameterMap map[string]interface{}) *harness.Step {

	// Extract values from parameterMap
	endpoints, ok := parameterMap["endpoints"].([]interface{})
	urls := []string{}
	headers := ""
	contentType := "application/json" // Default content type

	if ok && len(endpoints) > 0 {
		// Assuming only one endpoint in the array for simplicity
		if endpoint, ok := endpoints[0].(map[string]interface{}); ok {
			if url, ok := endpoint["url"].(string); ok {
				urls = append(urls, url)
			}
			if headerMap, ok := endpoint["headers"].(map[string]interface{}); ok {
				headerParts := []string{}
				for key, value := range headerMap {
					headerParts = append(headerParts, fmt.Sprintf("%s=%v", key, value))
				}
				headers = strings.Join(headerParts, ",")
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

	loglines, _ := parameterMap["loglines"].(string)
	phase, _ := parameterMap["phase"].(string)
	notes, _ := parameterMap["notes"].(string)

	// Create the template
	template := fmt.Sprintf(`{
  		"loglines": "%s",
  		"phase": "%s",
  		"status": "Success",
  		"notes": "%s",
  		"timestamp": "${time.now()}"
	}`, loglines, phase, notes)

	convertNotification := &harness.Step{
		Name: "Notification",
		Type: "plugin",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepPlugin{
			Connector: "<+input>",
			Image:     "plugins/webhook",
			With: map[string]interface{}{
				"urls":         urls,
				"method":       "POST",
				"username":     "<+input>",
				"password":     "<+input>",
				"token-value":  "<+input>",
				"token-type":   "<+input>",
				"content-type": contentType,
				"headers":      headers,
				"template":     template,
			},
		},
	}

	return convertNotification
}
