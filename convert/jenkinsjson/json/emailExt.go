package json

import (
	"strings"

	harness "github.com/drone/spec/dist/go"
)

func ConvertEmailext(node Node, variables map[string]string) *harness.Step {
	tokenReplacements := map[string]string{
		"${BUILD_NUMBER}": "<+pipeline.sequenceId>",
		"${BUILD_STATUS}": "<+pipeline.status>",
		"${PROJECT_NAME}": "<+project.name>",
		"${BUILD_URL}":    "<+pipeline.executionUrl>",
		"${BUILD_USER}":   "<+pipeline.triggeredBy.email>",
	}

	replaceTokens := func(input string) string {
		for token, replacement := range tokenReplacements {
			input = strings.ReplaceAll(input, token, replacement)
		}
		return input
	}
	getStringValue := func(key string) string {
		if value, ok := node.ParameterMap[key].(string); ok {
			return replaceTokens(value)
		}
		return ""
	}

	step := &harness.Step{
		Name: node.SpanName,
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "plugins/email",
			With: map[string]interface{}{
				"subject": getStringValue("subject"),
				"body":    getStringValue("body"),
				"to":      node.ParameterMap["to"],
				"from":    node.ParameterMap["from"],
				"replyTo": node.ParameterMap["replyTo"],
				"host":    "<+input>",
			},
		},
	}
	if len(variables) > 0 {
		step.Spec.(*harness.StepExec).Envs = variables
	}
	return step
}
