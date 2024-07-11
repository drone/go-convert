package json

import (
	"encoding/json"
	"fmt"

	harness "github.com/drone/spec/dist/go"
)

func ConvertWriteFile(node Node, variables map[string]string) *harness.Step {
	var text string
	var file string
	if attr, ok := node.AttributesMap["harness-attribute"]; ok {
		var attrMap map[string]interface{}
		if err := json.Unmarshal([]byte(attr), &attrMap); err == nil {
			if f, ok := attrMap["file"].(string); ok {
				file = f
			}
			if t, ok := attrMap["text"].(string); ok {
				text = t
			}
		}
	}
	step := &harness.Step{
		Name: node.SpanName,
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Type: "script",
		Spec: &harness.StepExec{
			Shell: "sh",
			Run:   fmt.Sprintf("printf '%s' > %s", text, file),
		},
	}
	if len(variables) > 0 {
		step.Spec.(*harness.StepExec).Envs = variables
	}
	return step
}
