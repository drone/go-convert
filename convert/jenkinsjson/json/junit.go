package json

import (
	"strings"

	harness "github.com/drone/spec/dist/go"
)

func ConvertJunit(node Node, variables map[string]string) *harness.Step {
	// Extract the "testResults" parameter from the node
	pathsString, ok := node.ParameterMap["testResults"].(string)
	if !ok {
		// Handle the case where the assertion fails by returning an empty Step
		return &harness.Step{}
	}

	// Split the comma-separated paths string into an array
	resultArray := strings.Split(pathsString, ",")

	// Create a Report with the resultArray and type "junit"
	report := harness.Report{
		Path: resultArray,
		Type: "junit",
	}

	// Create and return a Step with the specified properties
	junitStep := &harness.Step{
		Name: node.SpanName,
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Type: "script",
		Spec: &harness.StepExec{
			Shell:   "sh",
			Run:     "echo 'This Step is to Upload JUNIT Reports'",
			Reports: []*harness.Report{&report},
		},
	}
	if len(variables) > 0 {
		junitStep.Spec.(*harness.StepExec).Envs = variables
	}
	return junitStep
}
