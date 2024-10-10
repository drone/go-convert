package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConvertTarFile(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, "../convertTestFiles/tar/tarSnippet.json")
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
					"jenkins.pipeline.step.name":           "Create Tar file",
					"jenkins.pipeline.step.plugin.name":    "pipeline-utility-steps",
					"jenkins.pipeline.step.plugin.version": "2.17.0",
					"jenkins.pipeline.step.type":           "tar",
					"harness-attribute":                    "{\n  \"glob\" : \"**/*.txt\",\n  \"archive\" : true,\n  \"exclude\" : \"**/another.txt\",\n  \"file\" : \"hello3.tar\",\n  \"dir\" : \"hello\",\n  \"compress\" : false\n}",
					"harness-others":                       "",
				},
				Name:         "tarTest_AntStyle #1",
				Parent:       "tarTest_AntStyle",
				ParentSpanId: "6f5b865e9e5b447c",
				SpanId:       "7d1af0394d19d9b9",
				SpanName:     "tar",
				TraceId:      "4a4e917f696fb16e60967bba7bfc05e6",
				Type:         "Run Phase Span",
				ParameterMap: map[string]any{
					"file":     "hello3.tar",
					"compress": false,
					"glob":     "**/*.txt",
					"archive":  true,
					"exclude":  "**/another.txt",
					"dir":      "hello",
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
