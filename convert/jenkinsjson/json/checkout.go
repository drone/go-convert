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
			Image: "checkout_plugin",
			With: map[string]interface{}{
				"platform": node.AttributesMap["peer.service"],
				"git_url":  node.AttributesMap["http.url"],
				"branch":   node.AttributesMap["git.branch"],
				"depth":    node.AttributesMap["git.clone.depth"],
			},
		},
	}
	if len(variables) > 0 {
		step.Spec.(*harness.StepExec).Envs = variables
	}
	return step
}
