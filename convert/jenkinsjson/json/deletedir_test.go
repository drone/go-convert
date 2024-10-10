package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	harness "github.com/drone/spec/dist/go"
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

	var tests []runner
	tests = append(tests, prepare(t, "deleteDir/deleteDir_SunnyDay", &harness.Step{
		Id:   "deleteDir57d7ab",
		Name: "DeleteDir",
		Type: "script",
		Spec: &harness.StepExec{
			Shell: "sh",
			Run:   "dir_to_delete=$(pwd) && cd .. && rm -rf $dir_to_delete",
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertDeleteDir(tt.input, map[string]string{})
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertDeleteDir() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
