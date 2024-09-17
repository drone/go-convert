package json

import (
	"encoding/json"
	"log"

	harness "github.com/drone/spec/dist/go"
)

func ConvertZip(node Node, variables map[string]string) *harness.Step {
	var dir, zipFile string
	var exclude, glob string
	var overwrite bool

	if attr, ok := node.AttributesMap["harness-attribute"]; ok {
		var attrMap map[string]interface{}
		if err := json.Unmarshal([]byte(attr), &attrMap); err == nil {
			if d, ok := attrMap["dir"].(string); ok {
				dir = d
			}
			if zf, ok := attrMap["zipFile"].(string); ok {
				zipFile = zf
			}
			if e, ok := attrMap["exclude"].(string); ok {
				exclude = e
			}
			if gl, ok := attrMap["glob"].(string); ok {
				glob = gl
			}
			if ow, ok := attrMap["overwrite"].(bool); ok {
				overwrite = ow
			}
		} else {
			log.Printf("failed to unmarshal harness-attribute for node %s: %v", node.SpanName, err)
		}
	} else {
		log.Printf("harness-attribute missing for node %s", node.SpanName)
	}

	withProperties := make(map[string]interface{})
	withProperties["format"] = "zip"
	withProperties["action"] = "archive"
	if dir != "" {
		withProperties["target"] = dir
	}
	if zipFile != "" {
		withProperties["source"] = zipFile
	}
	if exclude != "" {
		withProperties["exclude"] = exclude
	}
	if glob != "" {
		withProperties["glob"] = glob
	}
	if overwrite != false {
		withProperties["overwrite"] = overwrite
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
