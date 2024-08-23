package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConvertVerifySha1(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, "../convertTestFiles/verifySha256/verifySha256Snippet.json")
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
					"jenkins.pipeline.step.id":             "9",
					"jenkins.pipeline.step.name":           "Verify the SHA256 of a given file",
					"jenkins.pipeline.step.plugin.name":    "pipeline-utility-steps",
					"jenkins.pipeline.step.plugin.version": "2.17.0",
					"jenkins.pipeline.step.type":           "verifySha256",
					"harness-attribute":                    "{\n  \"hash\" : \"a948904f2f0f479b8f8197694b30184b0d2ed1c1cd2a1ec0fb85d299a192a447\",\n  \"file\" : \"file.txt\"\n}",
					"harness-others":                       "",
				},
				Name:         "ag-readJSON #77",
				Parent:       "ag-readJSON",
				ParentSpanId: "5fba4997f08f9cf7",
				SpanId:       "b2b3687b5cc73816",
				SpanName:     "verifySha256",
				TraceId:      "7808bb3531a5d71e8e4468b16911eb96",
				Type:         "Run Phase Span",
				ParameterMap: map[string]any{
					"file": "file.txt",
					"hash": "a948904f2f0f479b8f8197694b30184b0d2ed1c1cd2a1ec0fb85d299a192a447"},
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
