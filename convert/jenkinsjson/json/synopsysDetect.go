package json

import (
	"encoding/json"
	"strings"

	harness "github.com/drone/spec/dist/go"
)

func ConvertSynopsysDetect(node Node, variables map[string]string) *harness.Step {
	var detectProperties string
	if attr, ok := node.AttributesMap["harness-attribute"]; ok {
		var attrMap map[string]interface{}
		if err := json.Unmarshal([]byte(attr), &attrMap); err == nil {
			if props, ok := attrMap["detectProperties"].(string); ok {
				detectProperties = props
			}
		}
	}

	parsedProperties := parseDetectProperties(detectProperties)
	withProperties := make(map[string]interface{})
	for key, value := range parsedProperties {
		withProperties[key] = value
	}

	step := &harness.Step{
		Name: node.SpanName,
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Connector: "c.docker",
			Image:     "harnesscommunitytest/synopsys-detect:latest",
			With:      withProperties,
		},
	}
	if len(variables) > 0 {
		step.Spec.(*harness.StepPlugin).Envs = variables
	}
	return step
}

func parseDetectProperties(detectProperties string) map[string]string {
	properties := map[string]string{
		"blackduck_url":             "--blackduck.url=",
		"blackduck_token":           "--blackduck.api.token=",
		"blackduck_project":         "--detect.project.name=",
		"blackduck_offline_mode":    "--blackduck.offline.mode=",
		"blackduck_test_connection": "--detect.test.connection=",
		"blackduck_offline_bdio":    "--blackduck.offline.mode.force.bdio=",
		"blackduck_trust_certs":     "--blackduck.trust.cert=",
		"blackduck_timeout":         "--detect.timeout=",
		"blackduck_scan_mode":       "--detect.blackduck.scan.mode=",
	}

	parsedProperties := make(map[string]string)
	remainingProperties := detectProperties

	for key, prefix := range properties {
		startIndex := strings.Index(detectProperties, prefix)
		if startIndex != -1 {
			startIndex += len(prefix)
			endIndex := strings.Index(detectProperties[startIndex:], " ")
			if endIndex == -1 {
				endIndex = len(detectProperties)
			} else {
				endIndex += startIndex
			}
			parsedProperties[key] = detectProperties[startIndex:endIndex]
			remainingProperties = strings.Replace(remainingProperties, prefix+parsedProperties[key], "", 1)
		}
	}

	remainingProperties = strings.TrimSpace(remainingProperties)
	if remainingProperties != "" {
		parsedProperties["blackduck_properties"] = remainingProperties
	}

	return parsedProperties
}
