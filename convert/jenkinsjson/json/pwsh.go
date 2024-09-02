package json

import (
	"fmt"

	harness "github.com/drone/spec/dist/go"
)

func ConvertPwsh(node Node, variables map[string]string, timeout string, dockerImage string) *harness.Step {
	pwshStep := &harness.Step{
		Name:    node.SpanName,
		Timeout: timeout,
		Id:      SanitizeForId(node.SpanName, node.SpanId),
		Type:    "script",
		Spec: &harness.StepExec{
			Image: dockerImage,
			Shell: "pwsh",
			Run:   fmt.Sprintf("%v", node.ParameterMap["script"]),
		},
	}
	if len(variables) > 0 {
		pwshStep.Spec.(*harness.StepExec).Envs = variables
	}
	return pwshStep
}
