package json

import (
	"strings"

	harness "github.com/drone/spec/dist/go"
)

// ConvertNunit creates a Harness step for nunit plugin.
func ConvertNunit(node Node, arguments map[string]interface{}) *harness.Step {
	testResultsPattern, _ := arguments["testResultsPattern"].(string)
	failIfNoResults, _ := arguments["failIfNoResults"].(bool)
	failedTestsFailBuild, _ := arguments["failedTestsFailBuild"].(bool)

	resultArray := strings.Split(testResultsPattern, ",")

	// Create a Report with the resultArray and type "junit"
	report := harness.Report{
		Type: "junit",
		Path: resultArray,
	}

	convertNunit := &harness.Step{
		Name: "nunit",
		Type: "plugin",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepPlugin{
			Image:   "plugins/nunit",
			Reports: []*harness.Report{&report},
			With: map[string]interface{}{
				"test_report_path":        testResultsPattern,
				"fail_if_no_results":      failIfNoResults,
				"failed_tests_fail_build": failedTestsFailBuild,
			},
		},
	}

	return convertNunit
}
