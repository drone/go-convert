package json

import (
	"fmt"

	harness "github.com/drone/spec/dist/go"
)

// ConvertJavadoc converts a Jenkins Node to a Harness Step with the Javadoc plugin spec.
func ConvertJavadoc(node Node, variables map[string]string) *harness.Step {
    // Extract the Javadoc directory path
    javadocDir := extractJavadocDir(node)
    if javadocDir == "" {
        fmt.Println("Warning: javadocDir not found, using default directory")
        javadocDir = "build/docs/javadoc"
    }

    javadocStep := &harness.Step{
        Id:   node.SpanId,
        Name: "Upload and Publish Javadoc",
        Type: "plugin",
        Spec: &harness.StepPlugin{
            Connector: "<+input>", // Placeholder for dynamic input
            Image:     "harnesscommunity/drone-s3-upload-publish",
            With: map[string]interface{}{
                "aws_access_key_id":     "<+input>",
                "aws_secret_access_key": "<+input>",
                "aws_bucket":            "<+input>",
                "source":                javadocDir,
                "target":                "<+pipeline.name>/<+pipeline.sequenceId>",
            },
        },
    }
	if len(variables) > 0 {
		javadocStep.Spec.(*harness.StepExec).Envs = variables
	}
    return javadocStep
}

func extractJavadocDir(node Node) string {
    if delegate, ok := node.ParameterMap["delegate"].(map[string]interface{}); ok {
        if arguments, ok := delegate["arguments"].(map[string]interface{}); ok {
            if dir, ok := arguments["javadocDir"].(string); ok {
                return dir
            }
        }
    }
    return ""
}