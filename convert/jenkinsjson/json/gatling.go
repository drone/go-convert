package json

import (
	harness "github.com/drone/spec/dist/go"
)

// ConvertGatling creates a Harness step for Gatling plugin.
func ConvertGatling(node Node) *harness.Step {
	convertGatling := &harness.Step{
		Name: "Gatling_Publish",
		Type: "plugin",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepPlugin{
			Image: "harnesscommunity/drone-s3-upload-publish",
			With: map[string]interface{}{
				"aws_access_key_id":     "<+input>",
				"aws_secret_access_key": "<+input>",
				"aws_default_region":    "<+input>",
				"aws_bucket":            "<+input>",
				"source":                "<+input>",
				"target":                "<+input>",
				"artifact_file":         "artifact.txt",
				"glob":                  "**/*.html, **/*.css",
			},
		},
	}

	return convertGatling
}
