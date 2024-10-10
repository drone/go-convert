package json

import (
	"fmt"

	harness "github.com/drone/spec/dist/go"
)

func ExtractEntries(node Node) []map[string]interface{} {
	delegate, ok := node.ParameterMap["delegate"].(map[string]interface{})
	if !ok {
		fmt.Println("Missing 'delegate' in parameterMap")
		return nil
	}
	// Extract the 'arguments' map from the 'delegate'
	arguments, ok := delegate["arguments"].(map[string]interface{})
	if !ok {
		fmt.Println("Missing 'arguments' in delegate map")
		return nil
	}
	// Extract values from the "entries" in the parameterMap
	entries, ok := arguments["entries"].([]interface{})
	if !ok {
		fmt.Println("No entries operations found in arguments")
		return nil
	}

	// Convert entries to a slice of maps for easier access
	var entryMaps []map[string]interface{}
	for _, entry := range entries {
		entryMap, ok := entry.(map[string]interface{})
		if ok {
			entryMaps = append(entryMaps, entryMap)
		} else {
			fmt.Println("Invalid entryMap format")
		}
	}

	return entryMaps
}

func Converts3Upload(node Node, entryMap map[string]interface{}) *harness.Step {
	var withPropertiesList []map[string]interface{}
	// Check the "parameterMap" and extract "delegate"

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

	// Create the step based on the first entry in the list
	if len(withPropertiesList) > 0 {
		step := &harness.Step{
			Name: node.SpanName,
			Id:   SanitizeForId(node.SpanName, node.SpanId),
			Type: "plugin",
			Spec: &harness.StepPlugin{
				Connector: "<+input>", // Using input expression for the connector
				Image:     "plugins/s3",
				With:      withPropertiesList[0], // Using the current entry for the step
			},
		}
		return step
	}

	return nil
}

// creates a Harness step for S3Upload
func Converts3Archive(node Node, entryMap map[string]interface{}) *harness.Step {
	step := &harness.Step{
		Name: "Plugin_0",
		Id:   "Plugin_0",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Connector: "<+input>",
			Image:     "plugins/archive",
			With: map[string]interface{}{
				"source":  ".",
				"target":  "s3Upload.gzip",
				"glob":    "*.txt",
				"exclude": entryMap["excludedFile"],
			},
		},
	}
	return step
}
