package json

import (
	"encoding/json"
	"log"

	harness "github.com/drone/spec/dist/go"
)

func ConvertUntar(node Node, variables map[string]string) *harness.Step {
	var dir, file string
	var glob string

	if attr, ok := node.AttributesMap["harness-attribute"]; ok {
		var attrMap map[string]interface{}
		if err := json.Unmarshal([]byte(attr), &attrMap); err == nil {
			if d, ok := attrMap["dir"].(string); ok {
				dir = d
			}
			if f, ok := attrMap["file"].(string); ok {
				file = f
			}
			if gl, ok := attrMap["glob"].(string); ok {
				glob = gl
			}
		} else {
			log.Printf("failed to unmarshal harness-attribute for node %s: %v", node.SpanName, err)
		}
	} else {
		log.Printf("harness-attribute missing for node %s", node.SpanName)
	}

	withProperties := make(map[string]interface{})
	withProperties["format"] = "tar"
	withProperties["action"] = "extract"
	if dir != "" {
		withProperties["target"] = dir
	}
	if file != "" {
		withProperties["source"] = file
	}
	if glob != "" {
		withProperties["glob"] = glob
	}

	step := &harness.Step{
		Name: node.SpanName,
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Connector: "c.docker",
			Image:     "plugins/archive:latest",
			With:      withProperties,
		},
	}

	if len(variables) > 0 {
		step.Spec.(*harness.StepPlugin).Envs = variables
	}

	return step
}
