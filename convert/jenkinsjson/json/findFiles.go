package json

import (
	harness "github.com/drone/spec/dist/go"
)

func ConvertFindFiles(node Node, variables map[string]string) *harness.Step {

	settings := map[string]interface{}{
		"glob": node.ParameterMap["glob"],
	}
	if value, ok := node.ParameterMap["excludes"]; ok {
		settings["excludes"] = value
	}

	step := &harness.Step{
		Name: node.SpanName,
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "harness-community/drone-findfiles:latest",
			With:  settings,
		},
	}
	if len(variables) > 0 {
		step.Spec.(*harness.StepPlugin).Envs = variables
	}
	return step
}
