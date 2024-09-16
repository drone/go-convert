package json

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConvertFindFiles(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, "../convertTestFiles/findFiles/findFilesSnippet.json")

	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read JSON file: %v", err)
	}

	var node Node
	if err := json.Unmarshal(jsonData, &node); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	tests := []struct {
		json Node
		want Node
	}{
		{
			json: node,
			want: Node{
				AttributesMap: map[string]string{
					"ci.pipeline.run.user":                 "SYSTEM",
					"harness-attribute":                    "{\n  \"glob\" : \"**/*.txt\"\n}",
					"harness-others":                       "",
					"jenkins.pipeline.step.name":           "Find files in the workspace",
					"jenkins.pipeline.step.id":             "9",
					"jenkins.pipeline.step.type":           "findFiles",
					"jenkins.pipeline.step.plugin.name":    "pipeline-utility-steps",
					"jenkins.pipeline.step.plugin.version": "2.17.0",
				},
				Name:         "Find Files #16",
				Parent:       "Find Files",
				ParentSpanId: "3e4d10779bd33fa0",
				SpanId:       "4d4efecbe310473f",
				SpanName:     "findFiles",
				TraceId:      "d46fda25e54d4691bf39cbb5712a3225",
				Type:         "Run Phase Span",
				ParameterMap: map[string]any{
					"glob": "**/*.txt",
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
