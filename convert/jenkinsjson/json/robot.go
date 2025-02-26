package json

import (
	harness "github.com/drone/spec/dist/go"
)

// ConvertRobot creates a Harness step for the Robot Framework plugin.
func ConvertRobot(node Node, parameterMap map[string]interface{}) *harness.Step {
	// Directly use the parameterMap to fetch the values
	reportDirectory, _ := parameterMap["outputPath"].(string)
	reportFileNamePattern, _ := parameterMap["outputFileName"].(string)
	if reportFileNamePattern == "" {
		reportFileNamePattern = "*.xml" // default value
	}
	passThreshold, _ := parameterMap["passThreshold"].(float64)
	unstableThreshold, _ := parameterMap["unstableThreshold"].(float64)
	countSkippedTests, _ := parameterMap["countSkippedTests"].(bool)
	onlyCritical, _ := parameterMap["onlyCritical"].(bool)

	// Build the parameter map dynamically, only including non-default values
	stepParams := map[string]interface{}{
		"level":               "info",
		"count_skipped_tests": false, // Default value
		"only_critical":       false, // Default value
	}

	if reportDirectory != "" {
		stepParams["report_directory"] = reportDirectory
	}
	if reportFileNamePattern != "" {
		stepParams["report_file_name_pattern"] = reportFileNamePattern
	}
	if passThreshold > 0 {
		stepParams["pass_threshold"] = int(passThreshold)
	}
	if unstableThreshold > 0 {
		stepParams["unstable_threshold"] = int(unstableThreshold)
	}
	if countSkippedTests {
		stepParams["count_skipped_tests"] = countSkippedTests
	}
	if onlyCritical {
		stepParams["only_critical"] = onlyCritical
	}

	// Create the Harness step for the Robot Framework plugin
	convertRobot := &harness.Step{
		Name: "robot",
		Type: "plugin",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepPlugin{
			Image: "plugins/robot",
			With:  stepParams,
		},
	}

	return convertRobot
}
