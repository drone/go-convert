package json

import (
	"fmt"

	harness "github.com/drone/spec/dist/go"
)

// ConvertTestng creates a Harness step for TestNG plugin.
func ConvertTestng(node Node, parameterMap map[string]interface{}) *harness.Step {

	delegate, ok := parameterMap["delegate"].(map[string]interface{})
	if !ok {
		fmt.Println("Error: 'delegate' is missing or not a valid map")
		return nil
	}

	// Extract the "arguments" key from the delegate map
	arguments, ok := delegate["arguments"].(map[string]interface{})
	if !ok {
		fmt.Println("Error: 'arguments' is missing or not a valid map")
		return nil
	}

	reportFilenamePattern, _ := arguments["reportFilenamePattern"].(string)
	failureOnFailedTestConfig, _ := arguments["failureOnFailedTestConfig"].(bool)
	// Extract and convert numeric fields from arguments
	failedSkips := 0
	if val, ok := arguments["failedSkips"].(float64); ok {
		failedSkips = int(val)
	}

	failedFails := 0
	if val, ok := arguments["failedFails"].(float64); ok {
		failedFails = int(val)
	}

	unstableSkips := 0
	if val, ok := arguments["unstableSkips"].(float64); ok {
		unstableSkips = int(val)
	}

	unstableFails := 0
	if val, ok := arguments["unstableFails"].(float64); ok {
		unstableFails = int(val)
	}

	thresholdMode, ok := arguments["thresholdMode"]

	if !ok {
		thresholdMode = 1
	}

	convertTestng := &harness.Step{
		Name: "testng",
		Type: "plugin",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepPlugin{
			Image: "plugins/testng",
			With: map[string]interface{}{
				"report_filename_pattern":       reportFilenamePattern,
				"failed_fails":                  failedFails,
				"failed_skips":                  failedSkips,
				"unstable_fails":                unstableFails,
				"unstable_skips":                unstableSkips,
				"threshold_mode":                thresholdMode,
				"failure_on_failed_test_config": failureOnFailedTestConfig,
				"fail_if_no_results":            true,
				"job_status":                    "<+pipeline.status>",
				"level":                         "info",
			},
		},
	}

	return convertTestng
}
