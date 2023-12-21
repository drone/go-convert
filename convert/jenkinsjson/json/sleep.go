package json

import (
	"fmt"
	harness "github.com/drone/spec/dist/go"
)

func ConvertSleep(node Node) *harness.Step {
	return &harness.Step{
		Name: node.SpanName,
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Type: "script",
		Spec: &harness.StepExec{
			Shell: "sh",
			Run:   fmt.Sprintf("sleep %v", node.ParameterMap["time"].(float64)),
		},
	}
}
