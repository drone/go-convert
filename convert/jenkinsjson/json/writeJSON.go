package json

import (
	"encoding/json"
	"fmt"

	harness "github.com/drone/spec/dist/go"
)

func ConvertWriteJSON(node Node, variables map[string]string) *harness.Step {
	var file string
	var jsonData map[string]interface{}

	if attr, ok := node.AttributesMap["harness-attribute"]; ok {
		var attrMap map[string]interface{}
		if err := json.Unmarshal([]byte(attr), &attrMap); err == nil {
			if f, ok := attrMap["file"].(string); ok {
				file = f
			}
			if j, ok := attrMap["json"].(map[string]interface{}); ok {
				jsonData = j
			}
		}
	}

	jsonString, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return nil
	}

	step := &harness.Step{
		Name: node.SpanName,
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Type: "script",
		Spec: &harness.StepExec{
			Shell: "sh",
			Run:   fmt.Sprintf("cat <<EOF > %s\n%s\nEOF", file, string(jsonString)),
		},
	}

	if len(variables) > 0 {
		step.Spec.(*harness.StepExec).Envs = variables
	}

	return step
}