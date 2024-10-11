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
	wantStep  *harness.Step
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

	if err := json.Unmarshal(jsonData, &inputNode); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}
	return s3runner{
		name:      filename,
		inputNode: inputNode,
		wantStep:  step,
	}
}

// Test function for Converts3Upload
func TestConverts3Upload(t *testing.T) {
	// Define test cases for Converts3Upload
	var tests []s3runner

	tests = append(tests, s3prepare(t, "s3upload/s3upload_snippet", &harness.Step{
		Id:   "s3UploadPlugin",
		Name: "s3Upload",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Connector: "<+input>",
			Image:     "plugins/s3",
			With: map[string]interface{}{
				"region":     "us-west-1",
				"bucket":     "bucket-1",
				"source":     "*.txt",
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
			for index, entry := range entries {
				// Convert the entryMap to a map for easy access
				entryMap, ok := entry.(map[string]interface{})
				if !ok {
					continue
				}

				got := Converts3Upload(tt.inputNode, entryMap, index)
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
		Id:   "s3UploadPlugin",
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

	for index, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Converts3Archive(tt.inputNode, map[string]interface{}{
				"excludedFile": "*.log",
			}, index)
			if diff := cmp.Diff(got, tt.wantStep); diff != "" {
				t.Errorf("Converts3Archive() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

type s3parserunner struct {
	name      string
	inputNode Node
	want      []map[string]interface{}
}

// s3prepare helper function to prepare test cases
func s3parse(name string, input map[string]interface{}, want []map[string]interface{}) s3parserunner {
	return s3parserunner{
		name:      name,
		inputNode: Node{ParameterMap: input},
		want:      want,
	}
}

// Test function for ExtractEntries
func TestExtractEntries(t *testing.T) {
	// Define test cases for ExtractEntries
	var tests []s3parserunner
	tests = append(tests, s3parse("Valid entries", map[string]interface{}{
		"delegate": map[string]interface{}{
			"arguments": map[string]interface{}{
				"entries": []interface{}{
					map[string]interface{}{"key1": "value1"},
					map[string]interface{}{"key2": "value2"},
				},
			},
		},
	}, []map[string]interface{}{
		{"key1": "value1"},
		{"key2": "value2"},
	}))
	tests = append(tests, s3parse("Missing delegate", map[string]interface{}{
		// no "delegate" key
	}, nil))
	tests = append(tests, s3parse("Invalid entries format", map[string]interface{}{
		"delegate": map[string]interface{}{
			"arguments": map[string]interface{}{
				"entries": "invalidFormat", // Not a slice of interface{}
			},
		},
	}, nil))

	// Execute each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractEntries(tt.inputNode)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ExtractEntries() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
