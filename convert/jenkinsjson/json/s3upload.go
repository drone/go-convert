package json

import (
	"fmt"

	harness "github.com/drone/spec/dist/go"
)

// ExtractEntries retrieves the "entries" from the given node's ParameterMap.
// It expects a "delegate" map containing "arguments" with an "entries" array.
// If any of these are missing or invalid, it returns nil.
//
// Parameters:
//  - node: The Node containing a ParameterMap.
//
// Returns:
//  - A slice of maps representing the entries, or nil if extraction fails.

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

// Converts3Upload creates a Harness plugin step for uploading files to S3.
// It uses data from the provided node and entryMap to configure the step and generates a unique ID for each step.
//
// Parameters:
//   - node: The Node containing context for the step.
//   - entryMap: A map containing key-value pairs used to customize the step.
//   - index: An incremental value used to ensure each step has a unique ID.
//
// Returns:
//   - harness.Step representing the configured S3 upload plugin step.
func Converts3Upload(node Node, entryMap map[string]interface{}, index int) *harness.Step {

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
	// Generate a unique ID using SanitizeForId with SpanName and index.
	sanitizedID := SanitizeForId(node.SpanName, node.SpanName)
	stepID := fmt.Sprintf("%s_%d", sanitizedID, index)

	step := &harness.Step{
		Name: stepID,
		Id:   stepID,
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Connector: "<+input>", // Using input expression for the connector
			Image:     "plugins/s3",
			With:      withProperties, // Using the current entry for the step
		},
	}
	return step
}

// Converts3Archive creates a Harness plugin step for archiving files and uploading them to S3.
// It uses data from the provided node and entryMap to configure the step, and generates a unique ID for each step.
//
// Parameters:
//   - node: The Node containing context for the step.
//   - entryMap: A map containing key-value pairs used to customize the step, such as excluded files.
//   - index: An incremental value used to ensure each step has a unique ID.
//
// Returns:
//   - harness.Step representing the configured S3 archive plugin step.
func Converts3Archive(node Node, entryMap map[string]interface{}, index int) *harness.Step {
	// Generate a unique ID using SanitizeForId and index
	sanitizedID := SanitizeForId(node.SpanName, node.SpanId)
	stepID := fmt.Sprintf("%s_%d", sanitizedID, index)

	step := &harness.Step{
		Name: stepID,
		Id:   stepID, // Use the sanitized ID
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
