package json

import (
	"fmt"
	"strings"

	harness "github.com/drone/spec/dist/go"
)

func Converts3UploadStep(node Node) *harness.Step {
	step := &harness.Step{
		Id:   SanitizeForId("UploadPublish", node.SpanId),
		Name: "Upload and Publish",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Connector: "<+input>", // Using input expression for the connector
			Image:     "harnesscommunity/drone-s3-upload-publish",
			With: map[string]interface{}{
				"aws_access_key_id":     "<+input>",
				"aws_secret_access_key": "<+input>",
				"aws_bucket":            "<+input>",
				"artifact_file":         "artifact.txt",
				"source":                "allure-report",
				"target":                "<+pipeline.sequenceId>",
			}, // Using the current entry for the step
		},
	}
	return step
}

func ConvertAllureReport(node Node) *harness.Step {
	// Extract the "results" paths as a comma-separated list of paths
	results := extractResultsPaths(node)

	step := &harness.Step{
		Name: node.SpanName,
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Type: "script",
		Spec: &harness.StepExec{
			Image: "harnesscommunity/allure:jdk21",
			Shell: "Sh",
			Run:   fmt.Sprintf("allure generate %s --clean --output allure-report", results),
		},
	}
	return step
}

func extractResultsPaths(node Node) string {
	if delegate, ok := node.ParameterMap["delegate"].(map[string]interface{}); ok {
		return getResultsPathsFromDelegate(delegate)
	}
	return ""
}

func getResultsPathsFromDelegate(delegate map[string]interface{}) string {
	if arguments, ok := delegate["arguments"].(map[string]interface{}); ok {
		return getResultsPathsFromArguments(arguments)
	}
	return ""
}

func getResultsPathsFromArguments(arguments map[string]interface{}) string {
	if resultPaths, ok := arguments["results"].([]interface{}); ok {
		return buildCommaSeparatedPaths(resultPaths)
	}
	return ""
}

func buildCommaSeparatedPaths(resultPaths []interface{}) string {
	paths := []string{}
	for _, result := range resultPaths {
		if pathMap, ok := result.(map[string]interface{}); ok {
			if path, ok := pathMap["path"].(string); ok {
				paths = append(paths, path)
			}
		}
	}
	return strings.Join(paths, ",")
}
