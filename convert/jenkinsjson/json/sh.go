package json

import (
	"fmt"
	"strings"

	harness "github.com/drone/spec/dist/go"
)

func ConvertSh(node Node, variables map[string]string, timeout string, dockerImage string, label string) *harness.Step {
	// The JD variables are Jenkins Docker built in variable name
	if node.ParameterMap["script"] == "docker push \"$JD_TAGGED_IMAGE_NAME\"" {
		step := &harness.Step{
			Name:    node.SpanName,
			Timeout: timeout,
			Id:      SanitizeForId(node.SpanName, node.SpanId),
			Type:    "plugin",
			Spec: &harness.StepPlugin{
				Image: "plugins/kaniko:latest",
				With: map[string]interface{}{
					"repo": variables["DOCKER_IMAGE"],
					"tags": variables["DOCKER_TAG"],
				},
			},
		}
		if len(variables) > 0 {
			step.Spec.(*harness.StepPlugin).Envs = variables
		}
		return step
	} else if node.ParameterMap["script"] == "docker build -t \"$JD_IMAGE\" ." || node.ParameterMap["script"] == "docker tag \"$JD_ID\" \"$JD_TAGGED_IMAGE_NAME\"" {
		return nil
	}

	trivy_prefix := "trivy image"
	if strings.HasPrefix(fmt.Sprintf("%v", node.ParameterMap["script"]), trivy_prefix) {
		image_raw := fmt.Sprintf("%v", node.ParameterMap["script"])
		image := strings.TrimSpace(strings.TrimPrefix(image_raw, trivy_prefix))
		image = strings.TrimSpace(image)
		parts := strings.Split(image, ":")
		step := &harness.Step{
			Name:    node.SpanName,
			Timeout: timeout,
			Id:      SanitizeForId(node.SpanName, node.SpanId),
			Type:    "plugin",
			Spec: &harness.StepPlugin{
				Image: "plugins/trivy:latest", // placeholder for sto
				With: map[string]interface{}{
					"image": parts[0],
					"tag":   parts[1],
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
			Image:     dockerImage,
			Connector: "account.harnessImage",
			Shell:     "sh",
			Run:       fmt.Sprintf("%v", node.ParameterMap["script"]),
		},
	}
	if len(variables) > 0 {
		shStep.Spec.(*harness.StepExec).Envs = variables
	}
	return shStep
}
