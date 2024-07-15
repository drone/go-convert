package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConvertWriteFile(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, "../convertTestFiles/writefile/writefileSnippet.json")
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
					"jenkins.pipeline.step.id":             "7",
					"jenkins.pipeline.step.name":           "Write file to workspace",
					"jenkins.pipeline.step.plugin.name":    "workflow-basic-steps",
					"jenkins.pipeline.step.plugin.version": "1058.vcb_fc1e3a_21a_9",
					"jenkins.pipeline.step.type":           "writeFile",
					"harness-attribute":                    "{\n  \"file\" : \"output1.txt\",\n  \"text\" : \"line1 \\n\\nline2\\n\\nline3\"\n}",
					"harness-others":                       "",
				},
				Name:         "writefile #1",
				Parent:       "writefile",
				ParentSpanId: "ffaea7a0704eaf64",
				SpanId:       "705399e59b21cf47",
				SpanName:     "writeFile",
				TraceId:      "900da4b8783a6df372c284b2fbdcf8a8",
				Type:         "Run Phase Span",
				ParameterMap: map[string]any{"file": string("output1.txt"), "text": string("line1 \n\nline2\n\nline3")},
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
