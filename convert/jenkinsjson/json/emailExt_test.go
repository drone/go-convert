package json

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
)

func TestConvertEmailext(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, "../convertTestFiles/emailext/emailextSnippet.json")
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read JSON file: %v", err)
	}

	var rawNode map[string]interface{}
	if err := json.Unmarshal(jsonData, &rawNode); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	node1 := Node{}
	if err := mapToNodes(rawNode, &node1); err != nil {
		t.Fatalf("failed to convert raw node to Node: %v", err)
	}

	tests := []struct {
		name string
		json Node
		vars map[string]string
		want *harness.Step
	}{
		{
			name: "Basic Conversion",
			json: node1,
			want: &harness.Step{
				Name: "emailext",
				Id:   "emailext2b1b55",
				Type: "plugin",
				Spec: &harness.StepPlugin{
					Image: "plugins/email",
					With: map[string]interface{}{
						"subject": "Email <+pipeline.status>",
						"body":    "Email <+pipeline.sequenceId>",
						"to":      "first@example.com second@example.com",
						"from":    nil,
						"replyTo": nil,
						"host":    "<+input>",
					},
				},
			},
		},
		{
			name: "Missing Fields",
			json: Node{
				SpanName: "emailext",
				SpanId:   "2b1b55",
			},
			want: &harness.Step{
				Name: "emailext",
				Id:   "emailext2b1b55",
				Type: "plugin",
				Spec: &harness.StepPlugin{
					Image: "plugins/email",
					With: map[string]interface{}{
						"subject": "",
						"body":    "",
						"to":      nil,
						"from":    nil,
						"replyTo": nil,
						"host":    "<+input>",
					},
				},
			},
		},
	}

	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := ConvertEmailext(test.json, test.vars, test.want.Timeout)
			if diff := cmp.Diff(got, test.want); diff != "" {
				t.Errorf("Unexpected conversion results for test %v", i)
				t.Log(diff)
			}
		})
	}
}

func mapToNodes(raw map[string]interface{}, node *Node) error {
	node.AttributesMap = make(map[string]string)
	for k, v := range raw["attributesMap"].(map[string]interface{}) {
		node.AttributesMap[k] = fmt.Sprintf("%v", v)
	}

	node.Name = raw["name"].(string)
	node.Parent = raw["parent"].(string)
	node.ParentSpanId = raw["parentSpanId"].(string)
	node.SpanId = raw["spanId"].(string)
	node.SpanName = raw["spanName"].(string)
	node.TraceId = raw["traceId"].(string)
	node.Type = raw["type"].(string)
	node.ParameterMap = raw["parameterMap"].(map[string]interface{})

	return nil
}
