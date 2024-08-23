package json

import (
	"encoding/json"
	"fmt"
	"strings"

	harness "github.com/drone/spec/dist/go"
)

func ConvertWriteYaml(node Node, variables map[string]string) *harness.Step {
	var file string
	var data map[string]interface{}

	if attr, ok := node.AttributesMap["harness-attribute"]; ok {
		var attrMap map[string]interface{}
		if err := json.Unmarshal([]byte(attr), &attrMap); err == nil {
			if f, ok := attrMap["file"].(string); ok {
				file = f
			}
			if d, ok := attrMap["data"].(map[string]interface{}); ok {
				data = d
			}
		}
	}

	var yamlContent strings.Builder
	for key, value := range data {
		yamlContent.WriteString(fmt.Sprintf("%s: %v\n", key, value))
	}

	step := &harness.Step{
		Name: node.SpanName,
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Type: "script",
		Spec: &harness.StepExec{
			Shell: "sh",
			Run:   fmt.Sprintf("cat <<EOF > %s\n%sEOF", file, yamlContent.String()),
		},
	}

	if len(variables) > 0 {
		step.Spec.(*harness.StepExec).Envs = variables
	}

	return step
}