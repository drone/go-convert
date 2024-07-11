package json

import (
	"encoding/json"
	"fmt"

	harness "github.com/drone/spec/dist/go"
)

func ConvertReadFile(node Node, variables map[string]string) *harness.Step {
	var file string
	if attr, ok := node.AttributesMap["harness-attribute"]; ok {
		var attrMap map[string]interface{}
		if err := json.Unmarshal([]byte(attr), &attrMap); err == nil {
			if f, ok := attrMap["file"].(string); ok {
				file = f
			}
		}
	}
	step := &harness.Step{
		Name: node.SpanName,
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Type: "script",
		Spec: &harness.StepExec{
			Shell: "sh",
			Run:   fmt.Sprintf("cat %s", file),
		},
	}
	if len(variables) > 0 {
		step.Spec.(*harness.StepExec).Envs = variables
	}
	return step
}
