package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConvertDir(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, "../convertTestFiles/deleteDir/DirSnippet.json")
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
			// children are not considered here for unit tests. Refer TestConvertDeleteDir as it is dir's child node
			json: node1,
			want: Node{
				AttributesMap: map[string]string{
					"jenkins.pipeline.step.type": "dir",
					"harness-attribute":          "{\n  \"path\" : \"target\"\n}",
					"harness-others":             "",
				},
				Name:         "Dir #1",
				Parent:       "Dir",
				ParentSpanId: "9c5d0bbf29308002",
				SpanId:       "fd7dfe5971c096aa",
				SpanName:     "Stage: null",
				TraceId:      "473b5dc91e544902871080a25554e963",
				Type:         "Run Phase Span",
				ParameterMap: map[string]any{"path": string("target")},
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

func TestConvertDeleteDir(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, "../convertTestFiles/deleteDir/deleteDirSnippet.json")
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
					"jenkins.pipeline.step.id":             "1",
					"jenkins.pipeline.step.name":           "Recursively delete the current directory from the workspace",
					"jenkins.pipeline.step.plugin.name":    "workflow-basic-steps",
					"jenkins.pipeline.step.plugin.version": "1058.vcb_fc1e3a_21a_9",
					"jenkins.pipeline.step.type":           "deleteDir",
					"harness-others":                       "",
				},
				Name:         "deleteDir #1",
				Parent:       "deleteDir",
				ParentSpanId: "58a6a7b62dc5eb76",
				SpanId:       "0f283dfd620daa10",
				SpanName:     "deleteDir",
				TraceId:      "473b5dc91e544902871080a25554e963",
				Type:         "Run Phase Span",
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
