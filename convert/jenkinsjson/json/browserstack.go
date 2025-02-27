package json

import (
	"strings"

	harness "github.com/drone/spec/dist/go"
)

// ConvertBrowserStack creates a Harness step for BrowserStack execution.
func ConvertBrowserStack(node Node) *harness.Step {

	// Get all shell scripts combined
	script := findAllBrowserStackScripts(node)

	browserstackStep := &harness.Step{
		Name: "Browserstack",
		Type: "script",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Spec: &harness.StepExec{
			Connector: "browserstack_connector",
			Image:     "harnesscommunity/browserstack",
			Shell:     "sh",
			Run:       script,
			Envs: map[string]string{
				"BROWSERSTACK_USERNAME":         "<+secrets.getValue(\"browserstack_username\")>",
				"BROWSERSTACK_ACCESS_KEY":       "<+secrets.getValue(\"browserstack_accesskey\")>",
				"BROWSERSTACK_BUILD_NAME":       "<+pipeline.sequenceId>-<+pipeline.executionId>",
				"BROWSERSTACK_BUILD_IDENTIFIER": "<+pipeline.executionId>",
			},
		},
	}

	return browserstackStep
}

// Function to collect all "sh" scripts from child nodes and concatenate them
func findAllBrowserStackScripts(node Node) string {
	var scripts []string

	// Check if the parent node is a BrowserStack step
	if node.AttributesMap["jenkins.pipeline.step.type"] == "browserstack" {
		// Iterate over child nodes
		for _, child := range node.Children {
			// If the child node is an "sh" step, extract its script
			if child.AttributesMap["jenkins.pipeline.step.type"] == "sh" {
				if scriptValue, ok := child.ParameterMap["script"]; ok {
					scripts = append(scripts, scriptValue.(string))
				}
			} else {
				// Recursively collect scripts from further descendants
				if result := findAllBrowserStackScripts(child); result != "" {
					scripts = append(scripts, result)
				}
			}
		}
	}

	// Join multiple scripts with a newline separator to ensure correct execution order
	return strings.Join(scripts, "\n")
}
