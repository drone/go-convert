package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConvertFileOpsCopy(t *testing.T) {
	// Get the working directory
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	// Update file path according to the location of fileOps_test.json
	filePath := filepath.Join(workingDir, "../convertTestFiles/fileOps/fileOpsCopy/fileOpsCopy_snippet.json")
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
				SpanId:  "a39de5c02c3bb8ff",
				TraceId: "939f2e2d83a1c32f40f7a86f4e5cacef",
				Parent:  "File_Ops_Copy_Test",
				Name:    "File_Ops_Copy_Test #2",
				AttributesMap: map[string]string{
					"harness-attribute-extra-pip: io.jenkins.plugins.opentelemetry.MigrateHarnessUrlChildAction@30d15708": "io.jenkins.plugins.opentelemetry.MigrateHarnessUrlChildAction@30d15708",
					"harness-others":                       "-delegate-field org.jenkinsci.plugins.workflow.steps.CoreStep delegate-org.jenkinsci.plugins.workflow.steps.CoreStep.delegate-interface jenkins.tasks.SimpleBuildStep",
					"jenkins.pipeline.step.name":           "File Operations",
					"ci.pipeline.run.user":                 "SYSTEM",
					"jenkins.pipeline.step.id":             "27",
					"jenkins.pipeline.step.type":           "fileOperations",
					"harness-attribute":                    "{\r\n  \"delegate\" : {\r\n    \"symbol\" : \"fileOperations\",\r\n    \"klass\" : null,\r\n    \"arguments\" : {\r\n      \"<anonymous>\" : [ {\r\n        \"symbol\" : \"fileCopyOperation\",\r\n        \"klass\" : null,\r\n        \"arguments\" : {\r\n          \"includes\" : \"src/*.txt\",\r\n          \"excludes\" : \"\",\r\n          \"targetLocation\" : \"dest/\",\r\n          \"flattenFiles\" : false,\r\n          \"renameFiles\" : false,\r\n          \"sourceCaptureExpression\" : \"\",\r\n          \"targetNameExpression\" : \"\",\r\n          \"useDefaultExcludes\" : true\r\n        },\r\n        \"model\" : null,\r\n        \"interpolatedStrings\" : [ ]\r\n      } ]\r\n    },\r\n    \"model\" : null\r\n  }\r\n}",
					"jenkins.pipeline.step.plugin.name":    "file-operations",
					"jenkins.pipeline.step.plugin.version": "266.v9d4e1eb_235b_a_",
					"harness-attribute-extra-pip: org.jenkinsci.plugins.workflow.actions.TimingAction@26e009c8": "org.jenkinsci.plugins.workflow.actions.TimingAction@26e009c8",
				},
				ParentSpanId: "35737607f38c09af",
				SpanName:     "fileOperations",
				Type:         "Run Phase Span",
				ParameterMap: map[string]interface{}{
					"delegate": map[string]interface{}{
						"symbol": "fileOperations",
						"klass":  nil,
						"arguments": map[string]interface{}{
							"<anonymous>": []interface{}{
								map[string]interface{}{
									"symbol":              "fileCopyOperation",
									"klass":               nil,
									"interpolatedStrings": []interface{}{},
									"arguments": map[string]interface{}{
										"excludes":                "",
										"targetNameExpression":    "",
										"includes":                "src/*.txt",
										"flattenFiles":            false,
										"targetLocation":          "dest/",
										"renameFiles":             false,
										"sourceCaptureExpression": "",
										"useDefaultExcludes":      true,
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
