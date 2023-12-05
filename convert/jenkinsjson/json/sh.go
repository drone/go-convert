package json

import (
	"fmt"

	harness "github.com/drone/spec/dist/go"
)

func ConvertSh(node Node) *harness.Step {
	return &harness.Step{
		Name: node.SpanName,
		Id:   node.SpanId,
		Type: "script",
		Spec: &harness.StepExec{
			Shell: "sh",
			Run:   fmt.Sprintf("%v", node.ParameterMap["script"]),
		},
	}
}
