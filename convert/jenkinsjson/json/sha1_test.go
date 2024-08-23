package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConvertSHA1(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, "../convertTestFiles/sha1/sha1Snippet.json")
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
					"jenkins.pipeline.step.name":           "Compute the SHA1 of a given file",
					"jenkins.pipeline.step.plugin.name":    "pipeline-utility-steps",
					"jenkins.pipeline.step.plugin.version": "2.17.0",
					"jenkins.pipeline.step.type":           "sha1",
					"harness-attribute":                    "{\n  \"file\" : \"hello.txt\"\n}",
					"harness-others":                       "",
				},
				Name:         "ag-readJSON #68",
				Parent:       "ag-readJSON",
				ParentSpanId: "38aeb40ddb2c07cf",
				SpanId:       "6c54047ff96cfa72",
				SpanName:     "sha1",
				TraceId:      "3f96d48045efd7afb65973ab88d00fa8",
				Type:         "Run Phase Span",
				ParameterMap: map[string]any{"file": "hello.txt"},
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
