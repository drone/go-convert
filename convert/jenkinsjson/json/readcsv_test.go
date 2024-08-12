package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConvertReadCsv(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, "../convertTestFiles/readcsv/readcsvSnippet.json")
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read CSV file: %v", err)
	}

	var node1 Node
	if err := json.Unmarshal(jsonData, &node1); err != nil {
		t.Fatalf("failed to decode CSV: %v", err)
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
					"jenkins.pipeline.step.id":             "9",
					"jenkins.pipeline.step.name":           "Read file from workspace",
					"jenkins.pipeline.step.plugin.name":    "workflow-basic-steps",
					"jenkins.pipeline.step.plugin.version": "1058.vcb_fc1e3a_21a_9",
					"jenkins.pipeline.step.type":           "readFile",
					"harness-attribute":                    "{\n  \"file\" : \"ex.csv\"\n}",
					"harness-others":                       "",
				},
				Name:         "ag-readcsv #1",
				Parent:       "ag-readcsv",
				ParentSpanId: "b5e4d8633a9f3810",
				SpanId:       "a48446f1eccfe72e",
				SpanName:     "readFile",
				TraceId:      "3ef97a9c6004ba651d9e4c35e9649ba3",
				Type:         "Run Phase Span",
				ParameterMap: map[string]any{"file": "ex.csv"},
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
