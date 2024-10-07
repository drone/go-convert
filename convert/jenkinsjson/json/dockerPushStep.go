package json

import (
	"encoding/json"
	harness "github.com/drone/spec/dist/go"
	"strings"
)

func ConvertDockerPushStep(node Node, variables map[string]string, timeout string) *harness.Step {
	var image string
	var tag string
	var targetRepo string

	if attr, ok := node.AttributesMap["harness-attribute"]; ok {
		var attrMap map[string]interface{}
		if err := json.Unmarshal([]byte(attr), &attrMap); err == nil {
			if img, ok := attrMap["image"].(string); ok {
				imageParts := strings.Split(img, ":")
				tag = imageParts[len(imageParts)-1]
				image = strings.Join(imageParts[:len(imageParts)-1], ":")
			}
			if tarRepo, ok := attrMap["targetRepo"].(string); ok {
				targetRepo = tarRepo
			}
		}
	}
	step := &harness.Step{
		Name:    node.SpanName,
		Timeout: timeout,
		Id:      SanitizeForId(node.SpanName, node.SpanId),
		Type:    "plugin",
		Spec: &harness.StepPlugin{
			Image: "plugins/kaniko:latest",
			With: map[string]interface{}{
				"repo":   image,
				"target": targetRepo,
				"tags":   tag,
			},
		},
	}
	if len(variables) > 0 {
		step.Spec.(*harness.StepPlugin).Envs = variables
	}
	return step
}
