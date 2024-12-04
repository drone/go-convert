package json

import (
	harness "github.com/drone/spec/dist/go"
)

func ConvertCheckout(node Node, variables map[string]string) *harness.Step {
	step := &harness.Step{
		Name: node.SpanName,
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "plugins/drone-git:latest",
			With: map[string]interface{}{
				"platform": node.AttributesMap["peer.service"],
				"git_url":  extractGitUrl(node),
				"branch":   extractGitBranch(node),
				"depth":    node.AttributesMap["git.clone.depth"],
			},
		},
	}
	return step
}

func extractGitUrl(node Node) string {
	return extractGitField(node, "http.url", "userRemoteConfigs", "url")
}

func extractGitBranch(node Node) string {
	return extractGitField(node, "git.branch", "branches", "name")
}

func extractGitField(node Node, attributesMapFieldName, pmArgumentsField, fieldName string) string {
	defaultValue, ok := node.AttributesMap[attributesMapFieldName]

	if !ok {
		scm, ok := node.ParameterMap["scm"].(map[string]interface{})
		if !ok {
			return ""
		}
		arguments, ok := scm["arguments"].(map[string]interface{})
		if !ok {
			return ""
		}
		listOfArgs, ok := arguments[pmArgumentsField].([]interface{})
		if !ok {
			return ""
		}

		if value, ok := readFieldValue(listOfArgs, fieldName); ok {
			return value
		}
	}

	return defaultValue
}

func readFieldValue(arguments []interface{}, fieldName string) (string, bool) {
	for _, element := range arguments {
		elementMap := element.(map[string]interface{})
		if subArguments, ok := elementMap["arguments"]; ok {
			elementMap = subArguments.(map[string]interface{})
		}

		value, ok := elementMap[fieldName]
		if ok {
			return value.(string), true
		}
	}
	return "", false
}
