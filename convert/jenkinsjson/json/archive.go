package json

import (
	"fmt"
	"strings"

	harness "github.com/drone/spec/dist/go"
)

func ConvertArchive(node Node) []*harness.Step {

	var steps []*harness.Step

	var artifacts string
	fmt.Print(node.ParameterMap)
	if artifactsValue, ok := node.ParameterMap["delegate"].(map[string]interface{})["arguments"].(map[string]interface{})["artifacts"].(string); ok {
		artifacts = artifactsValue
	} else if anonymousValue, ok := node.ParameterMap["delegate"].(map[string]interface{})["arguments"].(map[string]interface{})["<anonymous>"].(string); ok {
		artifacts = anonymousValue
	} else {
		fmt.Println("Neither 'artifacts' nor '<anonymous>' key found or both have non-string values.")
		return steps
	}
	artifactsArray := strings.Split(artifacts, ",")

	
	var excludes string
	if excludesValue, ok := node.ParameterMap["delegate"].(map[string]interface{})["arguments"].(map[string]interface{})["excludes"].(string); ok {
		excludes = excludesValue
	}


	for i, artifact := range artifactsArray {
		step := &harness.Step{
			Name: node.SpanName+fmt.Sprintf("_%d", i),
			Id:   SanitizeForId(node.SpanName, node.SpanId)+fmt.Sprintf("_%d", i),
			Type: "plugin",
			Spec: &harness.StepPlugin{
				Image: "plugins/s3",
				Name: node.SpanName+fmt.Sprintf("_%d", i),
				Connector: "harnessImage",
				With: map[string]interface{}{
					"source":     artifact,
					"access_key": "access-key",
					"secret_key": "secret-key",
					"bucket":     "bucket-name",
				},
				Envs: map[string]string{
					"EXCLUDE": excludes,
				},
			},
		}
		steps = append(steps, step)
	}

	return steps

}
