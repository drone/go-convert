package json

import (
	"log"

	harness "github.com/drone/spec/dist/go"
)

// creates a Harness step for S3Upload
func Converts3Upload(node Node, variables map[string]string) *harness.Step {
	var withPropertiesList []map[string]interface{}
	// Check the "parameterMap" and extract "delegate"
	if delegate, ok := node.ParameterMap["delegate"].(map[string]interface{}); ok {
		if arguments, ok := delegate["arguments"].(map[string]interface{}); ok {
			// Extract values from the "entries" in the parameterMap
			if entries, ok := arguments["entries"].([]interface{}); ok {
				for _, entry := range entries {
					if entryMap, ok := entry.(map[string]interface{}); ok {
						withProperties := make(map[string]interface{})

						// Extract values for each key from the current entry in the list
						if r, ok := entryMap["selectedRegion"].(string); ok {
							withProperties["region"] = r
						}
						if b, ok := entryMap["bucket"].(string); ok {
							withProperties["bucket"] = b
						}
						if s, ok := entryMap["sourceFile"].(string); ok {
							withProperties["source"] = s
						}
						if e, ok := entryMap["excludedFile"].(string); ok {
							withProperties["exclude"] = e
						}
						withProperties["access_key"] = "<+input>"
						withProperties["secret_key"] = "<+input>"
						withProperties["target"] = "<+input>"
						// Add the current entry's properties to the list
						withPropertiesList = append(withPropertiesList, withProperties)
					}
				}
			}
		} else {
			log.Println("Missing 'arguments' in delegate map")
			return nil
		}
	} else {
		log.Println("Missing 'delegate' in parameterMap")
		return nil
	}

	// Create the step based on the first entry in the list (or modify this as per your requirements)
	if len(withPropertiesList) > 0 {
		step := &harness.Step{
			Name: node.SpanName,
			Id:   SanitizeForId(node.SpanName, node.SpanId),
			Type: "plugin",
			Spec: &harness.StepPlugin{
				Connector: "<+input>", // Using input expression for the connector
				Image:     "plugins/s3",
				With:      withPropertiesList[0], // Using the first entry for the step
			},
		}

		// Add any environment variables to the step
		if len(variables) > 0 {
			step.Spec.(*harness.StepPlugin).Envs = variables
		}

		return step
	}

	return nil
}
