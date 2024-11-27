package json

import (
	"strings"

	harness "github.com/drone/spec/dist/go"
)

// ConvertUnarchive creates a Harness step for plugin step to handle unarchiving.
func ConvertUnarchive(node Node, paramMap map[string]interface{}) *harness.Step {

	// Extract the mapping from the parameterMap
	mapping := paramMap["mapping"].(map[string]interface{})

	// Prepare plugin step fields for dynamic input values
	var source, target string
	for s, t := range mapping {
		source = s
		target = t.(string)
	}

	// Determine the format based on file extension
	var format string
	switch {
	case strings.HasSuffix(source, ".zip"):
		format = "zip"
	case strings.HasSuffix(source, ".tar"):
		format = "tar"
	case strings.HasSuffix(source, ".tar.gz"):
		format = "tar.gz"
	}

	// Create a Harness plugin step for unarchiving
	convertPlugin := &harness.Step{
		Name: "Unarchive",
		Type: "plugin",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepPlugin{
			Image: "plugins/archive",
			With: map[string]interface{}{
				"source":    source,
				"target":    target,
				"format":    format,
				"action":    "extract",
				"glob":      "**/*",
				"overwrite": "true",
				"exclude":   "",
			},
		},
	}

	return convertPlugin
}
