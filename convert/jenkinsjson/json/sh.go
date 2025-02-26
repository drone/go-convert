package json

import (
	"fmt"

	harness "github.com/drone/spec/dist/go"
)

func ConvertSh(node Node, variables map[string]string, timeout string, dockerImage string, label string) *harness.Step {
	shStep := &harness.Step{
		Name:    SanitizeForName(node.SpanName) + label,
		Timeout: timeout,
		Id:      SanitizeForId(node.SpanName, node.SpanId),
		Type:    "script",
		Spec: &harness.StepExec{
			Image: dockerImage,
			Shell: "sh",
			Run:   fmt.Sprintf("%v", node.ParameterMap["script"]),
		},
	}
	if len(variables) > 0 {
		shStep.Spec.(*harness.StepExec).Envs = variables
	}
	return shStep
}
