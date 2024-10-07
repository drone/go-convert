package json

import (
	"encoding/json"
	"fmt"

	harness "github.com/drone/spec/dist/go"
	"github.com/tidwall/gjson"
)

type FileSpec struct {
	Pattern string `json:"pattern"`
	Target  string `json:"target"`
}

func ConvertArtifactUploadJfrog(node Node, variables map[string]string, timeout string) *harness.Step {
	step := &harness.Step{
		Name:    node.SpanName,
		Timeout: timeout,
		Id:      SanitizeForId(node.SpanName, node.SpanId),
		Type:    "plugin",
		Spec: &harness.StepPlugin{
			Image: "plugins/artifactory:latest",
			With: map[string]interface{}{
				"source": node.AttributesMap["pattern"],
				"target": node.AttributesMap["target"],
			},
		},
	}
	if len(variables) > 0 {
		step.Spec.(*harness.StepPlugin).Envs = variables
	}
	return step
}

func ParseHarnessAttribute(attr string) ([]FileSpec, error) {
	var attrMap map[string]interface{}
	if err := json.Unmarshal([]byte(attr), &attrMap); err != nil {
		return nil, err
	}

	specStr, ok := attrMap["spec"].(string)
	if !ok {
		return nil, fmt.Errorf("spec not found or not a string")
	}

	filesJson := gjson.Get(specStr, "files")
	if !filesJson.Exists() {
		return nil, fmt.Errorf("files not found in spec")
	}

	var fileSpecs []FileSpec
	if err := json.Unmarshal([]byte(filesJson.String()), &fileSpecs); err != nil {
		return nil, err
	}

	return fileSpecs, nil
}
