package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConvertFileOpsCreate(t *testing.T) {
	// Get the working directory
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	// Update file path according to the location of fileOps_test.json
	filePath := filepath.Join(workingDir, "../convertTestFiles/fileOps/fileOpsCreate/fileOpsCreate_snippet.json")
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read JSON file: %v", err)
	}

	// Unmarshal the JSON into a Node struct
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
				SpanId:  "9b3c555cf009005f",
				TraceId: "5c4a042d6d48afeba2d56a38cb497c1e",
				Parent:  "op-filecreateoperations",
				Name:    "op-filecreateoperations #2",
				AttributesMap: map[string]string{
					"harness-others":             "-delegate-field org.jenkinsci.plugins.workflow.steps.CoreStep delegate-org.jenkinsci.plugins.workflow.steps.CoreStep.delegate-interface jenkins.tasks.SimpleBuildStep",
					"jenkins.pipeline.step.name": "File Operations",
					"ci.pipeline.run.user":       "SYSTEM",
					"jenkins.pipeline.step.id":   "7",
					"harness-attribute-extra-pip: io.jenkins.plugins.opentelemetry.MigrateHarnessUrlChildAction@5fad4a0b": "io.jenkins.plugins.opentelemetry.MigrateHarnessUrlChildAction@5fad4a0b",
					"harness-attribute-extra-pip: org.jenkinsci.plugins.workflow.actions.TimingAction@1e9f64fb":           "org.jenkinsci.plugins.workflow.actions.TimingAction@1e9f64fb",
					"jenkins.pipeline.step.type": "fileOperations",
					"harness-attribute-extra-pip: com.cloudbees.workflow.rest.endpoints.FlowNodeAPI@bd3676d": "com.cloudbees.workflow.rest.endpoints.FlowNodeAPI@bd3676d",
					"harness-attribute":                    "{\n  \"delegate\" : {\n    \"symbol\" : \"fileOperations\",\n    \"klass\" : null,\n    \"arguments\" : {\n      \"<anonymous>\" : [ {\n        \"symbol\" : \"fileCreateOperation\",\n        \"klass\" : null,\n        \"arguments\" : {\n          \"fileName\" : \"newfile.txt\",\n          \"fileContent\" : \"Hello, World!\"\n        },\n        \"model\" : null,\n        \"interpolatedStrings\" : [ ]\n      } ]\n    },\n    \"model\" : null\n  }\n}",
					"jenkins.pipeline.step.plugin.name":    "file-operations",
					"jenkins.pipeline.step.plugin.version": "266.v9d4e1eb_235b_a_",
				},
				ParentSpanId: "5fc14ef50c3e4cda",
				SpanName:     "fileOperations",
				Type:         "Run Phase Span",
				ParameterMap: map[string]interface{}{
					"delegate": map[string]interface{}{
						"symbol": "fileOperations",
						"klass":  nil,
						"arguments": map[string]interface{}{
							"<anonymous>": []interface{}{
								map[string]interface{}{
									"symbol":              "fileCreateOperation",
									"klass":               nil,
									"interpolatedStrings": []interface{}{},
									"arguments": map[string]interface{}{
										"fileName":    "newfile.txt",
										"fileContent": "Hello, World!",
									},
									"model": nil,
								},
							},
						},
						"model": nil,
					},
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
