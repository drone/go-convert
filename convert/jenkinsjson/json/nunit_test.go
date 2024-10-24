package json

import (
	"fmt"
	"testing"

	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
)

func extractArguments(currentNode Node) map[string]interface{} {
	// Step 1: Extract the 'delegate' map from the 'parameterMap'
	delegate, ok := currentNode.ParameterMap["delegate"].(map[string]interface{})
	if !ok {
		fmt.Println("Missing 'delegate' in parameterMap")
	}

	// Step 2: Extract the 'arguments' map from the 'delegate'
	arguments, ok := delegate["arguments"].(map[string]interface{})
	if !ok {
		fmt.Println("Missing 'arguments' in delegate map")
	}

	return arguments
}

func TestConvertNunit(t *testing.T) {

	var tests []runner
	tests = append(tests, prepare(t, "/nunit/nunit_snippet", &harness.Step{
		Id:   "nunit0d6d45",
		Name: "nunit",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image:   "plugins/nunit",
			Reports: []*harness.Report{{Path: harness.Stringorslice{"TestResult.xml"}, Type: "junit"}},
			With: map[string]interface{}{
				"test_report_path":        string("TestResult.xml"),
				"fail_if_no_results":      bool(true),
				"failed_tests_fail_build": bool(true),
			},
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operation := extractArguments(tt.input)
			got := ConvertNunit(tt.input, operation)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertNunit() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
