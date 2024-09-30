package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConvertFileOpsDelete(t *testing.T) {
	// Get the working directory
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	// Update file path according to the location of fileOps_test.json
	filePath := filepath.Join(workingDir, "../convertTestFiles/fileOps/fileOpsDelete/fileOpsDelete_snippet.json")
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
				SpanId:  "128b1701d569ef2b",
				TraceId: "a9082386e210d442a5c28f7693d5a695",
				Parent:  "op-filedeleteoperations",
				Name:    "op-filedeleteoperations #3",
				AttributesMap: map[string]string{
					"harness-attribute-extra-pip: io.jenkins.plugins.opentelemetry.MigrateHarnessUrlChildAction@6a202719": "io.jenkins.plugins.opentelemetry.MigrateHarnessUrlChildAction@6a202719",
					"harness-others":             "-delegate-field org.jenkinsci.plugins.workflow.steps.CoreStep delegate-org.jenkinsci.plugins.workflow.steps.CoreStep.delegate-interface jenkins.tasks.SimpleBuildStep",
					"jenkins.pipeline.step.name": "File Operations",
					"ci.pipeline.run.user":       "SYSTEM",
					"jenkins.pipeline.step.id":   "16",
					"harness-attribute-extra-pip: org.jenkinsci.plugins.workflow.actions.TimingAction@d4d078e": "org.jenkinsci.plugins.workflow.actions.TimingAction@d4d078e",
					"jenkins.pipeline.step.type": "fileOperations",
					"harness-attribute-extra-pip: com.cloudbees.workflow.rest.endpoints.FlowNodeAPI@65a1849c": "com.cloudbees.workflow.rest.endpoints.FlowNodeAPI@65a1849c",
					"harness-attribute":                    "{\n  \"delegate\" : {\n    \"symbol\" : \"fileOperations\",\n    \"klass\" : null,\n    \"arguments\" : {\n      \"<anonymous>\" : [ {\n        \"symbol\" : \"fileDeleteOperation\",\n        \"klass\" : null,\n        \"arguments\" : {\n          \"includes\" : \"**/old-files/*.log\",\n          \"excludes\" : \"\"\n        },\n        \"model\" : null,\n        \"interpolatedStrings\" : [ ]\n      } ]\n    },\n    \"model\" : null\n  }\n}",
					"jenkins.pipeline.step.plugin.name":    "file-operations",
					"jenkins.pipeline.step.plugin.version": "266.v9d4e1eb_235b_a_",
				},
				ParentSpanId: "65cb4661ced5c4fe",
				SpanName:     "fileOperations",
				Type:         "Run Phase Span",
				ParameterMap: map[string]interface{}{
					"delegate": map[string]interface{}{
						"symbol": "fileOperations",
						"klass":  nil,
						"arguments": map[string]interface{}{
							"<anonymous>": []interface{}{
								map[string]interface{}{
									"symbol":              "fileDeleteOperation",
									"klass":               nil,
									"interpolatedStrings": []interface{}{},
									"arguments": map[string]interface{}{
										"excludes": "",
										"includes": "**/old-files/*.log",
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
