package json

import (
	"fmt"

	harness "github.com/drone/spec/dist/go"
)

func ConvertBat(node Node, variables map[string]string, timeout string) *harness.Step {
	script := fmt.Sprintf("%v", node.ParameterMap["script"])
	stepId := SanitizeForId(node.SpanName, node.SpanId)
	cmd := getFileCreateCommand(stepId, script)
	batStep := &harness.Step{
		Name:    node.SpanName,
		Timeout: timeout,
		Id:      stepId,
		Type:    "script",
		Spec: &harness.StepExec{
			Shell: "pwsh",
			Run:   cmd,
		},
	}
	if len(variables) > 0 {
		batStep.Spec.(*harness.StepExec).Envs = variables
	}
	return batStep
}

func getFileCreateCommand(stepId string, script string) string {
	return fmt.Sprintf("echo @\"\n%v\n\"@ > %v.bat\n ./%v.bat", script, stepId, stepId)
}
