package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConvertArtifactUploadJfrog(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, "../convertTestFiles/artifactUploadJfrog/artifactuploadSnippet.json")
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
					"jenkins.pipeline.step.id":             "23",
					"jenkins.pipeline.step.name":           "Upload artifacts",
					"jenkins.pipeline.step.plugin.name":    "artifactory",
					"jenkins.pipeline.step.plugin.version": "4.0.8",
					"jenkins.pipeline.step.type":           "artifactoryUpload",
					"harness-attribute":                    "{\n  \"failNoOp\" : false,\n  \"server\" : \"UNSERIALIZABLE\",\n  \"buildInfo\" : \"UNSERIALIZABLE\",\n  \"spec\" : \"{\\n                    \\\"files\\\": [{\\n                        \\\"pattern\\\": \\\"output1.txt\\\",\\n                        \\\"target\\\": \\\"/artifactory/generic-local/op-ci/\\\"\\n                    },{\\n                        \\\"pattern\\\": \\\"output2.txt\\\",\\n                        \\\"target\\\": \\\"/artifactory/generic-local/op-ci2/\\\"\\n                    },{\\n                        \\\"pattern\\\": \\\"output3.txt\\\",\\n                        \\\"target\\\": \\\"/artifactory/generic-local/op-ci3/\\\"\\n                    }]\\n                 }\"\n}",
					"harness-others":                       "",
				},
				Name:         "artifactuploadjfrog #8",
				Parent:       "artifactuploadjfrog",
				ParentSpanId: "d4a2b65115153102",
				SpanId:       "c722199b42be60f3",
				SpanName:     "artifactoryUpload",
				TraceId:      "f01a23a1d020b2a9c7434cbd060f8835",
				Type:         "Run Phase Span",
				ParameterMap: map[string]any{
					"server":    "UNSERIALIZABLE",
					"buildInfo": "UNSERIALIZABLE",
					"failNoOp":  false,
					"spec":      "{\n                    \"files\": [{\n                        \"pattern\": \"output1.txt\",\n                        \"target\": \"/artifactory/generic-local/\"\n                    },{\n                        \"pattern\": \"output2.txt\",\n                        \"target\": \"/artifactory/generic-local/\"\n                    },{\n                        \"pattern\": \"output3.txt\",\n                        \"target\": \"/artifactory/generic-local/\"\n                    }]\n                 }",
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
