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

func TestConvertFileOpsDownlaodFunction(t *testing.T) {

	var tests []testRunner
	tests = append(tests, prepareFileOpsTest(t, "fileOpsDownload_snippet", "fileOpsDownload", &harness.Step{
		Id:   "fileOperationsd6f89b",
		Name: "fileDownloadOperation",
		Type: "script",
		Spec: &harness.StepExec{
			Image: "alpine/curl",
			Run:   string("wget -P downloads/ https://github.com/git/git/archive/refs/heads/master.zip"),
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operation := extractAnanymousOperation(tt.input)
			got := ConvertFileDownload(tt.input, operation)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertFileDownload() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertFileOpsJoinFunction(t *testing.T) {

	var tests []testRunner
	tests = append(tests, prepareFileOpsTest(t, "fileOpsJoin_snippet", "fileOpsJoin", &harness.Step{
		Id:   "fileOperationsc8844b",
		Name: "fileJoinOperation",
		Type: "script",
		Spec: &harness.StepExec{
			Image: "alpine",
			Run:   string("cat file1.txt >> file2.txt"),
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operation := extractAnanymousOperation(tt.input)
			got := ConvertFileJoin(tt.input, operation)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertFileJoin() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertFileOpsJsonFunction(t *testing.T) {

	var tests []testRunner
	tests = append(tests, prepareFileOpsTest(t, "fileOpsJson_snippet", "fileOpsJson", &harness.Step{
		Id:   "fileOperations9397c0",
		Name: "filePropertiesToJsonOperation",
		Type: "script",
		Spec: &harness.StepExec{
			Image: "alpine",
			Run:   string("stat -c '{\"size\": %s, \"permissions\": \"%A\", \"owner\": %U, \"group\": %G, \"last_modified\": \"%y\"}' newfile.properties > property.json"),
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operation := extractAnanymousOperation(tt.input)
			got := ConvertFileJson(tt.input, operation)
			fmt.Sprintln(got)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertFileJson() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertFileOpsRenameFunction(t *testing.T) {

	var tests []testRunner
	tests = append(tests, prepareFileOpsTest(t, "fileOpsRename_snippet", "fileOpsRename", &harness.Step{
		Id:   "fileOperationse0b0a3",
		Name: "fileRenameOperation",
		Type: "script",
		Spec: &harness.StepExec{
			Image: "alpine",
			Run:   string("mv newfile.txt renamedfile.txt"),
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operation := extractAnanymousOperation(tt.input)
			got := ConvertFileRename(tt.input, operation)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertFileRename() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertFileOpsTransformFunction(t *testing.T) {

	var tests []testRunner
	tests = append(tests, prepareFileOpsTest(t, "fileOpsTransform_snippet", "fileOpsTransform", &harness.Step{
		Id:   "fileOperationsc8fb57",
		Name: "fileTransformOperation",
		Type: "script",
		Spec: &harness.StepExec{
			Image: "alpine",
			Run:   string("find . -type f -name 'newfile*.txt' ! -name 'newfile2.txt' -exec sh -c 'iconv -f <source_encoding> -t UTF-8 \"$0\" -o \"${0%.txt}.utf8\"' {} \\;"),
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operation := extractAnanymousOperation(tt.input)
			got := ConvertFileTranform(tt.input, operation)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertFileTranform() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertFileOpsFolderCopyFunction(t *testing.T) {

	var tests []testRunner
	tests = append(tests, prepareFileOpsTest(t, "fileOpsFolderCopy_snippet", "fileOpsFolderCopy", &harness.Step{
		Id:   "fileOperations8a2e5a",
		Name: "folderCopyOperation",
		Type: "script",
		Spec: &harness.StepExec{
			Image: "alpine",
			Run:   string("cp -r src/ dest"),
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operation := extractAnanymousOperation(tt.input)
			got := ConvertFolderCopy(tt.input, operation)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertFolderCopy() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertFileOpsFolderCreateFunction(t *testing.T) {

	var tests []testRunner
	tests = append(tests, prepareFileOpsTest(t, "fileOpsFolderCreate_snippet", "fileOpsFolderCreate", &harness.Step{
		Id:   "fileOperations7d2ecf",
		Name: "folderCreateOperation",
		Type: "script",
		Spec: &harness.StepExec{
			Image: "alpine",
			Run:   string("mkdir -p src/create"),
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operation := extractAnanymousOperation(tt.input)
			got := ConvertFolderCreate(tt.input, operation)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertFolderCreate() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertFileOpsFolderDeleteFunction(t *testing.T) {

	var tests []testRunner
	tests = append(tests, prepareFileOpsTest(t, "fileOpsFolderDelete_snippet", "fileOpsFolderDelete", &harness.Step{
		Id:   "fileOperations978830",
		Name: "folderDeleteOperation",
		Type: "script",
		Spec: &harness.StepExec{
			Image: "alpine",
			Run:   string("rm -rf src/create"),
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operation := extractAnanymousOperation(tt.input)
			got := ConvertFolderDelete(tt.input, operation)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertFolderDelete() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertFileOpsFolderRenameFunction(t *testing.T) {

	var tests []testRunner
	tests = append(tests, prepareFileOpsTest(t, "fileOpsFolderRename_snippet", "fileOpsFolderRename", &harness.Step{
		Id:   "fileOperations12ada0",
		Name: "folderRenameOperation",
		Type: "script",
		Spec: &harness.StepExec{
			Image: "alpine",
			Run:   string("mv src/ dest/"),
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operation := extractAnanymousOperation(tt.input)
			got := ConvertFolderRename(tt.input, operation)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertFolderRename() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertFileOpsUntarFunction(t *testing.T) {

	var tests []testRunner
	tests = append(tests, prepareFileOpsTest(t, "fileOpsUntar_snippet", "fileOpsUnTar", &harness.Step{
		Id:   "fileOperationsbf0ce1",
		Name: "fileUnTarOperation",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "plugins/archive:latest",
			With: map[string]interface{}{
				"action": "extract",
				"format": "gzip",
				"source": "src.tar.gz",
				"target": "dest/",
			},
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operation := extractAnanymousOperation(tt.input)
			got := ConvertFileUntar(tt.input, operation)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertFileUntar() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertFileOpsUnZipFunction(t *testing.T) {

	var tests []testRunner
	tests = append(tests, prepareFileOpsTest(t, "fileOpsUnZip_snippet", "fileOpsUnZip", &harness.Step{
		Id:   "fileOperations0d571f",
		Name: "fileUnZipOperation",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "plugins/archive:latest",
			With: map[string]interface{}{
				"action": "extract",
				"format": "zip",
				"source": "src.zip",
				"target": "dest/",
			},
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operation := extractAnanymousOperation(tt.input)
			got := ConvertFileUnzip(tt.input, operation)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertFileUnzip() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertFileOpsZipFunction(t *testing.T) {

	var tests []testRunner
	tests = append(tests, prepareFileOpsTest(t, "fileOpsZip_snippet", "fileOpsZip", &harness.Step{
		Id:   "fileOperations4ca6ab",
		Name: "fileZipOperation",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "plugins/archive:latest",
			With: map[string]interface{}{
				"action":    "archive",
				"format":    "zip",
				"overwrite": "true",
				"source":    "src/",
				"target":    "dest/",
			},
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operation := extractAnanymousOperation(tt.input)
			got := ConvertFileZip(tt.input, operation)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertFileZip() mismatch (-want +got):\n%s", diff)
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
