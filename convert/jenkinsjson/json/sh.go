package json

import (
	"fmt"

	harness "github.com/drone/spec/dist/go"
)

func ConvertSh(node Node) *harness.Step {
	return &harness.Step{
		Name: node.SpanName,
		Id:   node.SpanId,
		Spec: &harness.StepRun{
			Shell:  "sh",
			Script: []string{fmt.Sprintf("%v", node.ParameterMap["script"])},
		},
	}
}
