package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConvertUnzipFile(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, "../convertTestFiles/unzip/unzipSnippet.json")
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
					"jenkins.pipeline.step.id":             "17",
					"jenkins.pipeline.step.name":           "Extract Zip file",
					"jenkins.pipeline.step.plugin.name":    "pipeline-utility-steps",
					"jenkins.pipeline.step.plugin.version": "2.17.0",
					"jenkins.pipeline.step.type":           "unzip",
					"harness-attribute":                    "{\n  \"dir\" : \"./archive_new2\",\n  \"zipFile\" : \"test.zip\"\n}",
					"harness-others":                       "",
				},
				Name:         "zipTest_AntStyle #1",
				Parent:       "zipTest_AntStyle",
				ParentSpanId: "ab763408b4f2d7c0",
				SpanId:       "c9e2eb9f12367d6c",
				SpanName:     "unzip",
				TraceId:      "fb302d811edeb970820eee05e1ae3d08",
				Type:         "Run Phase Span",
				ParameterMap: map[string]any{
					"dir":     "./archive_new2",
					"zipFile": "test.zip",
				},
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
