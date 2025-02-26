package json

import (
	"fmt"

	harness "github.com/drone/spec/dist/go"
)

// ConvertCucumber creates a Harness step for the Cucumber plugin.
func ConvertCucumber(node Node, parameterMap map[string]interface{}) *harness.Step {
	delegate, ok := parameterMap["delegate"].(map[string]interface{})
	if !ok {
		fmt.Println("Error: 'delegate' is missing or not a valid map")
		return nil
	}

	arguments, ok := delegate["arguments"].(map[string]interface{})
	if !ok {
		fmt.Println("Error: 'arguments' is missing or not a valid map")
		return nil
	}

	// Extract values from arguments with default values
	fileIncludePattern, _ := arguments["fileIncludePattern"].(string)

	fileExcludePattern, _ := arguments["fileExcludePattern"].(string)

	failedAsNotFailingStatus, _ := arguments["failedAsNotFailingStatus"].(bool)
	mergeFeaturesById, _ := arguments["mergeFeaturesById"].(bool)
	pendingAsNotFailingStatus, _ := arguments["pendingAsNotFailingStatus"].(bool)
	skipEmptyJSONFiles, _ := arguments["skipEmptyJSONFiles"].(bool)
	skippedAsNotFailingStatus, _ := arguments["skippedAsNotFailingStatus"].(bool)
	stopBuildOnFailedReport, _ := arguments["stopBuildOnFailedReport"].(bool)
	undefinedAsNotFailingStatus, _ := arguments["undefinedAsNotFailingStatus"].(bool)

	jsonReportDirectory, _ := arguments["jsonReportDirectory"].(string)
	sortingMethod, _ := arguments["sortingMethod"].(string)

	// Extract and convert numeric fields from arguments if they exist
	failedFeaturesNumber, _ := arguments["failedFeaturesNumber"].(float64)
	failedFeaturesPercentage, _ := arguments["failedFeaturesPercentage"].(float64)
	failedScenariosNumber, _ := arguments["failedScenariosNumber"].(float64)
	failedScenariosPercentage, _ := arguments["failedScenariosPercentage"].(float64)
	failedStepsNumber, _ := arguments["failedStepsNumber"].(float64)
	failedStepsPercentage, _ := arguments["failedStepsPercentage"].(float64)
	pendingStepsNumber, _ := arguments["pendingStepsNumber"].(float64)
	pendingStepsPercentage, _ := arguments["pendingStepsPercentage"].(float64)
	skippedStepsNumber, _ := arguments["skippedStepsNumber"].(float64)
	skippedStepsPercentage, _ := arguments["skippedStepsPercentage"].(float64)
	undefinedStepsNumber, _ := arguments["undefinedStepsNumber"].(float64)
	undefinedStepsPercentage, _ := arguments["undefinedStepsPercentage"].(float64)

	// Build the parameter map dynamically, only including non-default values
	stepParams := map[string]interface{}{
		"file_include_pattern": fileIncludePattern,
		"level":                "info", // Default
	}

	if fileExcludePattern != "" {
		stepParams["file_exclude_pattern"] = fileExcludePattern
	}
	if jsonReportDirectory != "" {
		stepParams["json_report_directory"] = jsonReportDirectory
	}
	if sortingMethod != "" {
		stepParams["sorting_method"] = sortingMethod
	}
	if failedAsNotFailingStatus {
		stepParams["failed_as_not_failing_status"] = failedAsNotFailingStatus
	}
	if mergeFeaturesById {
		stepParams["merge_features_by_id"] = mergeFeaturesById
	}
	if pendingAsNotFailingStatus {
		stepParams["pending_as_not_failing_status"] = pendingAsNotFailingStatus
	}
	if skipEmptyJSONFiles {
		stepParams["skip_empty_json_files"] = skipEmptyJSONFiles
	}
	if skippedAsNotFailingStatus {
		stepParams["skipped_as_not_failing_status"] = skippedAsNotFailingStatus
	}
	if stopBuildOnFailedReport {
		stepParams["stop_build_on_failed_report"] = stopBuildOnFailedReport
	}
	if undefinedAsNotFailingStatus {
		stepParams["undefined_as_not_failing_status"] = undefinedAsNotFailingStatus
	}

	// Include numeric values if non-zero
	if failedFeaturesNumber > 0 {
		stepParams["failed_features_number"] = failedFeaturesNumber
	}
	if failedFeaturesPercentage > 0 {
		stepParams["failed_features_percentage"] = failedFeaturesPercentage
	}
	if failedScenariosNumber > 0 {
		stepParams["failed_scenarios_number"] = failedScenariosNumber
	}
	if failedScenariosPercentage > 0 {
		stepParams["failed_scenarios_percentage"] = failedScenariosPercentage
	}
	if failedStepsNumber > 0 {
		stepParams["failed_steps_number"] = failedStepsNumber
	}
	if failedStepsPercentage > 0 {
		stepParams["failed_steps_percentage"] = failedStepsPercentage
	}
	if pendingStepsNumber > 0 {
		stepParams["pending_steps_number"] = pendingStepsNumber
	}
	if pendingStepsPercentage > 0 {
		stepParams["pending_steps_percentage"] = pendingStepsPercentage
	}
	if skippedStepsNumber > 0 {
		stepParams["skipped_steps_number"] = skippedStepsNumber
	}
	if skippedStepsPercentage > 0 {
		stepParams["skipped_steps_percentage"] = skippedStepsPercentage
	}
	if undefinedStepsNumber > 0 {
		stepParams["undefined_steps_number"] = undefinedStepsNumber
	}
	if undefinedStepsPercentage > 0 {
		stepParams["undefined_steps_percentage"] = undefinedStepsPercentage
	}

	convertCucumber := &harness.Step{
		Name: "cucumber",
		Type: "plugin",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepPlugin{
			Image: "plugins/cucumber",
			With:  stepParams,
		},
	}

	return convertCucumber
}
