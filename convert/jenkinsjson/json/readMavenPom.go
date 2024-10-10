package json

import (
	harness "github.com/drone/spec/dist/go"
)

func ConvertReadMavenPom(node Node) *harness.Step {

	settings := map[string]interface{}{
		"pom_path": node.ParameterMap["file"],
	}

	step := &harness.Step{
		Name: node.SpanName,
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "harnesscommunity/drone-get-maven-version:latest",
			With:  settings,
		},
	}
	return step
}
