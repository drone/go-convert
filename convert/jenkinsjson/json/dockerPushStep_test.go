package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConvertDockerPushStep(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, "../convertTestFiles/dockerPushStep/dockerPushStepSnippet.json")
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
				SpanId:       "ed0750b60f975189",
				TraceId:      "de4a510a21792ddd94d610f1cbf9533b",
				Parent:       "dockerTest",
				Name:         "dockerTest #4",
				Type:         "Run Phase Span",
				ParentSpanId: "1fe7c9899ab7ec10",
				SpanName:     "dockerPushStep",
				AttributesMap: map[string]string{
					"harness-others":                       "",
					"jenkins.pipeline.step.name":           "Artifactory docker push",
					"jenkins.pipeline.step.id":             "32",
					"ci.pipeline.run.user":                 "SYSTEM",
					"jenkins.pipeline.step.type":           "dockerPushStep",
					"jenkins.pipeline.step.plugin.name":    "artifactory",
					"jenkins.pipeline.step.plugin.version": "4.0.7",
					"harness-attribute": `{
  "image" : "anshika/testimage:latest",
  "server" : "https://localhost:5000/",
  "targetRepo" : "testimage"
}`,
				},
				ParameterMap: map[string]any{
					"image":      "anshika/testimage:latest",
					"server":     "https://localhost:5000/",
					"targetRepo": "testimage",
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
