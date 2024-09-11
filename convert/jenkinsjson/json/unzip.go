package json

import (
	"encoding/json"
	"log"

	harness "github.com/drone/spec/dist/go"
)

// TODO: this logic needs to be updated according to the plugin
func ConvertUnzip(node Node, variables map[string]string) *harness.Step {
	var dir, zipFile string

	if attr, ok := node.AttributesMap["harness-attribute"]; ok {
		var attrMap map[string]interface{}
		if err := json.Unmarshal([]byte(attr), &attrMap); err == nil {
			if d, ok := attrMap["dir"].(string); ok {
				dir = d
			}
			if zf, ok := attrMap["zipFile"].(string); ok {
				zipFile = zf
			}
		} else {
			log.Printf("failed to unmarshal harness-attribute for node %s: %v", node.SpanName, err)
		}
	} else {
		log.Printf("harness-attribute missing for node %s", node.SpanName)
	}

	withProperties := make(map[string]interface{})
	if dir != "" {
		withProperties["dir"] = dir
	}
	if zipFile != "" {
		withProperties["zipFile"] = zipFile
	}

	step := &harness.Step{
		Name: node.SpanName,
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Connector: "c.docker",
			Image:     "harnesscommunitytest/archive-plugin:latest",
			With:      withProperties,
		},
	}

	if len(variables) > 0 {
		step.Spec.(*harness.StepPlugin).Envs = variables
	}

	return step
}
