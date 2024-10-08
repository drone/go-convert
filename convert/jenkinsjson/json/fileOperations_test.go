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

func prepareFileOpsTest(t *testing.T, filename string, folderName string, step *harness.Step) testRunner {

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

	return testRunner{
		name:  filename,
		input: inputNode,
		want:  step,
	}
}

func TestConvertFileOpsCreateFunction(t *testing.T) {

	var tests []testRunner
	tests = append(tests, prepareFileOpsTest(t, "fileOpsCreate_snippet", "fileOpsCreate", &harness.Step{
		Id:   "fileOperations9b3c55",
		Name: "fileCreateOperation",
		Type: "script",
		Spec: &harness.StepExec{
			Image: "alpine",
			Run:   string("echo 'Hello, World!' > newfile.txt"),
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operation := extractAnanymousOperation(tt.input)
			got := ConvertFileCreate(tt.input, operation)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertFileCreate() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
func TestConvertFileOpsDeleteFunction(t *testing.T) {

	var tests []testRunner
	tests = append(tests, prepareFileOpsTest(t, "fileOpsDelete_snippet", "fileOpsDelete", &harness.Step{
		Id:   "fileOperations128b17",
		Name: "fileDeleteOperation",
		Type: "script",
		Spec: &harness.StepExec{
			Image: "alpine",
			Run:   string("rm -rf **/old-files/*.log"),
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operation := extractAnanymousOperation(tt.input)
			got := ConvertFileDelete(tt.input, operation)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertFileDelete() mismatch (-want +got):\n%s", diff)
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
