package json

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
)

type testRunner struct {
	name  string
	input Node
	want  *harness.Step
}

func TestConvertFileOpsCopyTrace(t *testing.T) {
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

func prepareFileOpsTest(t *testing.T, filename string, folderName string, step *harness.Step) runner {

	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	jsonData, err := os.ReadFile(filepath.Join(workingDir, "../convertTestFiles/fileOps/"+folderName, filename+".json"))
	if err != nil {
		t.Fatalf("failed to read JSON file: %v", err)
	}

	var inputNode Node
	if err := json.Unmarshal(jsonData, &inputNode); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	return runner{
		name:  filename,
		input: inputNode,
		want:  step,
	}
}

func TestConvertFileOpsCopyFunction(t *testing.T) {

	var tests []runner
	tests = append(tests, prepareFileOpsTest(t, "fileOpsCopy_snippet", "fileOpsCopy", &harness.Step{
		Id:   "fileOperationsa39de5",
		Name: "fileCopyOperation",
		Type: "script",
		Spec: &harness.StepExec{
			Image: "alpine",
			Run:   string("cp -r src/*.txt dest/"),
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operation := extractAnanymousOperation(tt.input)
			got := ConvertFileCopy(tt.input, operation)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertFileCopy() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func extractAnanymousOperation(currentNode Node) map[string]interface{} {
	// Step 1: Extract the 'delegate' map from the 'parameterMap'
	delegate, ok := currentNode.ParameterMap["delegate"].(map[string]interface{})
	if !ok {
		fmt.Println("Missing 'delegate' in parameterMap")
	}

	// Step 2: Extract the 'arguments' map from the 'delegate'
	arguments, ok := delegate["arguments"].(map[string]interface{})
	if !ok {
		fmt.Println("Missing 'arguments' in delegate map")
	}

	// Step 3: Extract the list of anonymous operations
	anonymousOps, ok := arguments["<anonymous>"].([]interface{})
	if !ok {
		fmt.Println("No anonymous operations found in arguments")
	}
	var extractedOperation map[string]interface{}
	// Step 4: Iterate over each operation and handle based on the 'symbol' type
	for _, op := range anonymousOps {
		// Convert the operation to a map for easy access
		operation, ok := op.(map[string]interface{})
		if !ok {
			fmt.Println("Invalid operation format")
			continue
		}
		extractedOperation = operation
	}
	return extractedOperation
}
