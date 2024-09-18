package json

import (
	"encoding/json"
	"log"

	harness "github.com/drone/spec/dist/go"
)

func ConvertTar(node Node, variables map[string]string) *harness.Step {
	var dir, file string
	var exclude, glob string
	var overwrite, tarcompress bool

	if attr, ok := node.AttributesMap["harness-attribute"]; ok {
		var attrMap map[string]interface{}
		if err := json.Unmarshal([]byte(attr), &attrMap); err == nil {
			if d, ok := attrMap["dir"].(string); ok {
				dir = d
			}
			if f, ok := attrMap["file"].(string); ok {
				file = f
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
			if tc, ok := attrMap["compress"].(bool); ok {
				tarcompress = tc
			}
		} else {
			log.Printf("failed to unmarshal harness-attribute for node %s: %v", node.SpanName, err)
		}
	} else {
		log.Printf("harness-attribute missing for node %s", node.SpanName)
	}

	withProperties := make(map[string]interface{})
	withProperties["format"] = "tar"
	withProperties["action"] = "archive"
	if dir != "" {
		withProperties["source"] = dir
	}
	if file != "" {
		withProperties["target"] = file
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
	if tarcompress != false {
		withProperties["tarcompress"] = tarcompress
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
