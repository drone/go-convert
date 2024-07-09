package json

import (
	"fmt"
	harness "github.com/drone/spec/dist/go"
)

func ConvertSleep(node Node, variables map[string]string) *harness.Step {
	sleepStep := &harness.Step{
		Name: node.SpanName,
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Type: "script",
		Spec: &harness.StepExec{
			Shell: "sh",
			Run:   fmt.Sprintf("sleep %v", node.ParameterMap["time"].(float64)),
		},
	}
	if len(variables) > 0 {
		sleepStep.Spec.(*harness.StepExec).Envs = variables
	}
	return sleepStep
}
