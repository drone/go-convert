package json

import (
	"fmt"

	harness "github.com/drone/spec/dist/go"
)

func ConvertSh(node Node, variables map[string]string, timeout string, dockerImage string, label string) *harness.Step {
	// "docker push \"$JD_TAGGED_IMAGE_NAME\""
	if node.ParameterMap["script"] == "docker push \"$JD_TAGGED_IMAGE_NAME\"" {
		step := &harness.Step{
			Name:    node.SpanName,
			Timeout: timeout,
			Id:      SanitizeForId(node.SpanName, node.SpanId),
			Type:    "plugin",
			Spec: &harness.StepPlugin{
				Image: "plugins/kaniko:latest",
				With: map[string]interface{}{
					"repo": "$DOCKER_IMAGE",
					"tags": "$DOCKER_TAG",
				},
			},
		}
		if len(variables) > 0 {
			step.Spec.(*harness.StepPlugin).Envs = variables
		}
		return step
	}

	shStep := &harness.Step{
		Name:    SanitizeForName(node.SpanName) + label,
		Timeout: timeout,
		Id:      SanitizeForId(node.SpanName, node.SpanId),
		Type:    "script",
		Spec: &harness.StepExec{
			Image: dockerImage,
			Shell: "sh",
			Run:   fmt.Sprintf("%v", node.ParameterMap["script"]),
		},
	}
	if len(variables) > 0 {
		shStep.Spec.(*harness.StepExec).Envs = variables
	}
	return shStep
}
