package json

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
)

type s3runner struct {
	name      string
	inputNode Node
	//entryMap map[string]interface{}
	wantStep *harness.Step
}

// Helper function to prepare test cases from JSON files
func s3prepare(t *testing.T, filename string, step *harness.Step) s3runner {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	jsonData, err := os.ReadFile(filepath.Join(workingDir, "../convertTestFiles/s3publisher", filename+".json"))
	if err != nil {
		t.Fatalf("failed to read JSON file: %v", err)
	}

	var inputNode Node
	//	entryMap map[string]interface{}

	if err := json.Unmarshal(jsonData, &inputNode); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}
	return s3runner{
		name:      filename,
		inputNode: inputNode,
		//entryMap map,

		wantStep: step,
	}
}

// Test function for Converts3Upload
func TestConverts3Upload(t *testing.T) {
	// Define test cases for Converts3Upload
	var tests []s3runner

	tests = append(tests, s3prepare(t, "s3upload/s3upload_snippet", &harness.Step{
		Id:   "s3Upload1ec902",
		Name: "s3Upload",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Connector: "<+input>",
			Image:     "plugins/s3",
			With: map[string]interface{}{
				"region": "us-west-1",
				"bucket": "bucket-1",
				"source": "*.txt",
				//"glob":       "*.txt",
				"exclude":    "2.txt",
				"access_key": "<+input>",
				"secret_key": "<+input>",
				"target":     "<+input>",
			},
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delegate := tt.inputNode.ParameterMap["delegate"].(map[string]interface{})
			// Extract the 'arguments' map from the 'delegate'
			arguments := delegate["arguments"].(map[string]interface{})
			//Extract values from the "entries" in the parameterMap
			entries := arguments["entries"].([]interface{})
			// Iterate over each entry and handle based on the 'symbol' type
			for _, entry := range entries {
				// Convert the entryMap to a map for easy access
				entryMap, ok := entry.(map[string]interface{})
				if !ok {
					continue
				}

				got := Converts3Upload(tt.inputNode, entryMap)
				if diff := cmp.Diff(got, tt.wantStep); diff != "" {
					t.Errorf("Converts3Upload() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

// Test function for Converts3Archive
func TestConverts3Archive(t *testing.T) {
	// Define test cases for Converts3Archive
	var tests []s3runner

	// Append a test case using the s3prepare helper function
	tests = append(tests, s3prepare(t, "s3upload/s3upload_snippet", &harness.Step{
		Id:   "Plugin_0",
		Name: "Plugin_0",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Connector: "<+input>",
			Image:     "plugins/archive",
			With: map[string]interface{}{
				"source":  ".",
				"target":  "s3Upload.gzip",
				"glob":    "*.txt",
				"exclude": "*.log",
			},
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Converts3Archive(tt.inputNode, map[string]interface{}{
				"excludedFile": "*.log",
			})
			if diff := cmp.Diff(got, tt.wantStep); diff != "" {
				t.Errorf("Converts3Archive() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
