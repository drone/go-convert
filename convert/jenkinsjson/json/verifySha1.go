package json

import (
	"encoding/json"
	"fmt"

	harness "github.com/drone/spec/dist/go"
)

func ConvertVerifySha1(node Node, variables map[string]string, dockerImage string) *harness.Step {
	var file string
	var hash string
	if attr, ok := node.AttributesMap["harness-attribute"]; ok {
		var attrMap map[string]interface{}
		if err := json.Unmarshal([]byte(attr), &attrMap); err == nil {
			if f, ok := attrMap["file"].(string); ok {
				file = f
			}
			if h, ok := attrMap["hash"].(string); ok {
				hash = h
			}
		}
	}

	var runCommand string = fmt.Sprintf("echo %v %v | sha1sum -c", hash,file)
	verifySha1step := &harness.Step{
		Name: node.SpanName,
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Type: "script",
		Spec: &harness.StepExec{
			Image:   dockerImage,
			Shell:   "sh",
			Run:     runCommand,
		},
	}

	if len(variables) > 0 {
		verifySha1step.Spec.(*harness.StepExec).Envs = variables
	}

	return verifySha1step
}
