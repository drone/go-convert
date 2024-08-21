package json

import (
	"encoding/json"
	"fmt"

	harness "github.com/drone/spec/dist/go"
)

func ConvertSHA1(node Node, variables map[string]string, dockerImage string) *harness.Step {
	var file string
	if attr, ok := node.AttributesMap["harness-attribute"]; ok {
		var attrMap map[string]interface{}
		if err := json.Unmarshal([]byte(attr), &attrMap); err == nil {
			if f, ok := attrMap["file"].(string); ok {
				file = f
			}
		}
	}
	var runCommand string = fmt.Sprintf("checksum=$(sha1sum %s | awk '{print $1}')", file)
	sha1step := &harness.Step{
		Name: node.SpanName,
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Type: "script",
		Spec: &harness.StepExec{
			Image: dockerImage,
			Shell: "sh",
			Run:   runCommand,
			Outputs: []string{"checksum"},
		},
	}

	if len(variables) > 0 {
		sha1step.Spec.(*harness.StepExec).Envs = variables
	}

	return sha1step
}
