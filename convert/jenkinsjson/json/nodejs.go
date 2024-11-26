package json

import (
	harness "github.com/drone/spec/dist/go"
)

// ConvertNodejs creates a Harness step for nunit plugin.
func ConvertNodejs(node Node) *harness.Step {
	// Recursively process children
	var script string
	processNodejsScript(node, &script)
	convertNodeJs := &harness.Step{
		Name: "Nodejs",
		Type: "script",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepExec{
			Connector: "<+input_docker_hub_connector>",
			Image:     "node:latest",
			Shell:     "sh",
			Run:       script,
		},
	}

	return convertNodeJs
}

func processNodejsScript(currentNode Node, result *string) {
	for _, child := range currentNode.Children {
		processNode(child, result)
	}
}

func processNode(node Node, result *string) {
	if stepType, exists := node.AttributesMap["jenkins.pipeline.step.type"]; exists {
		switch stepType {
		case "wrap":
			for _, child := range node.Children {
				processNode(child, result)
			}
		case "bat", "sh":
			if script, ok := node.ParameterMap["script"].(string); ok {
				*result += script + "\n"
			}
		}
	}
}
