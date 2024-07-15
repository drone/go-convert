package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConvertReadFile(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, "../convertTestFiles/readfile/readfileSnippet.json")
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
					"jenkins.pipeline.step.id":             "14",
					"jenkins.pipeline.step.name":           "Read file from workspace",
					"jenkins.pipeline.step.plugin.name":    "workflow-basic-steps",
					"jenkins.pipeline.step.plugin.version": "1058.vcb_fc1e3a_21a_9",
					"jenkins.pipeline.step.type":           "readFile",
					"harness-attribute":                    "{\n  \"file\" : \"output1.txt\"\n}",
					"harness-others":                       "",
				},
				Name:         "readfile #1",
				Parent:       "readfile",
				ParentSpanId: "5af3d9cada8b7ebe",
				SpanId:       "484461c194e15443",
				SpanName:     "readFile",
				TraceId:      "900da4b8783a6df372c284b2fbdcf8a8",
				Type:         "Run Phase Span",
				ParameterMap: map[string]any{"file": string("output1.txt")},
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
