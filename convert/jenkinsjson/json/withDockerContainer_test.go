package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConvertWithDockerContainer(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, "../convertTestFiles/withDockerContainer/withDockerContainerSnippet.json")
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read JSON file: %v", err)
	}

	var node1 Node
	if err := json.Unmarshal(jsonData, &node1); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	tests := []struct {
		json Node
		want Node
	}{
		{
			json: node1,
			want: Node{
				AttributesMap: map[string]string{
					"ci.pipeline.run.user":                 "SYSTEM",
					"jenkins.pipeline.step.id":             "3",
					"jenkins.pipeline.step.name":           "agent.allocate",
					"jenkins.pipeline.step.plugin.name":    "workflow-durable-task-step",
					"jenkins.pipeline.step.plugin.version": "1360.v82d13453da_a_f",
					"jenkins.pipeline.step.type":           "withDockerContainer",
					"harness-attribute":                    "{\n  \"args\" : \"\",\n  \"image\" : \"node:20.16.0-alpine3.20\",\n  \"toolName\" : null\n}",
					"harness-others":                       "",
				},
				Name:         "docker declarative #4",
				Parent:       "docker declarative",
				ParentSpanId: "3dbf046c2e8a9822",
				SpanId:       "77809b66e716b761",
				SpanName:     "Stage: null",
				TraceId:      "9a3549687043a2c1f88c2796a279a7b7",
				Type:         "Run Phase Span",
				ParameterMap: map[string]any{
					"args":     "",
					"image":    "node:20.16.0-alpine3.20",
					"toolName": ""},
			},
		},
	}

	for i, test := range tests {
		got := test.json
		if diff := cmp.Diff(got, test.want); diff != "" {
			t.Errorf("Unexpected parsing results for test %v", i)
			t.Log(diff)
		}
	}
}
