package json

import (
	"encoding/json"
	"fmt"
	"log"

	harness "github.com/drone/spec/dist/go"
)

func ConvertReadYaml(node Node, variables map[string]string) *harness.Step {
	var file, text string
	attr, ok := node.AttributesMap["harness-attribute"]
	if !ok {
		log.Printf("harness-attribute missing for node, spanId=%s", node.SpanId)
		return nil
	}

	var attrMap map[string]interface{}
	err := json.Unmarshal([]byte(attr), &attrMap)
	if err != nil {
		log.Printf("failed to unmarshal harness-attribute for node, spanId=%s \n %v", node.SpanId, err)
		return nil
	}

	if f, ok := attrMap["file"].(string); ok {
		file = f
	}
	if t, ok := attrMap["text"].(string); ok {
		text = t
	}

	if file == "" && text == "" {
		log.Printf("No valid attributes found for node, spanId=%s", node.SpanId)
		return nil
	}

	var runCommand string
	if file != "" {
		runCommand = fmt.Sprintf("cat %s", file)
	} else {
		runCommand = fmt.Sprintf("echo '%s'", text)
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
