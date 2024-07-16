package json

import (
	"encoding/json"
	"fmt"
	"log"

	harness "github.com/drone/spec/dist/go"
)

func ConvertReadYaml(node Node, variables map[string]string) *harness.Step {
	var file string
	if attr, ok := node.AttributesMap["harness-attribute"]; ok {
		var attrMap map[string]interface{}
		if err := json.Unmarshal([]byte(attr), &attrMap); err == nil {
			if f, ok := attrMap["file"].(string); ok {
				file = f
			}
		} else {
			log.Printf("failed to unmarshal harness-attribute for node %s: %v", node.SpanName, err)
		}
	} else {
		log.Printf("harness-attribute missing for node %s", node.SpanName)
	}

	var runCommand string
	if file != "" {
		runCommand = fmt.Sprintf("cat %s", file)
	} else {
		log.Printf("no valid attribute found for node %s", node.SpanName)
		return nil
	}

	step := &harness.Step{
		Name: node.SpanName,
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Type: "script",
		Spec: &harness.StepExec{
			Shell: "sh",
			Run:   runCommand,
		},
	}

	if len(variables) > 0 {
		step.Spec.(*harness.StepExec).Envs = variables
	}
	return step
}
