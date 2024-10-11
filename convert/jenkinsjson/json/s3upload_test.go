package json

import (
	"testing"

	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
)

// Test function for Converts3Upload
func TestConverts3Upload(t *testing.T) {
	// Define test cases for Converts3Upload
	var tests []runner

	tests = append(tests, prepare(t, "s3publisher/s3upload/s3upload_snippet", &harness.Step{
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
			delegate := tt.input.ParameterMap["delegate"].(map[string]interface{})
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

				got := Converts3Upload(tt.input, entryMap, index)
				if diff := cmp.Diff(got, tt.want); diff != "" {
					t.Errorf("Converts3Upload() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

// Test function for Converts3Archive
func TestConverts3Archive(t *testing.T) {
	// Define test cases for Converts3Archive
	var tests []runner

	// Append a test case using the s3prepare helper function
	tests = append(tests, prepare(t, "s3publisher/s3upload/s3upload_snippet", &harness.Step{
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
			got := Converts3Archive(tt.input, map[string]interface{}{
				"excludedFile": "*.log",
			}, index)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Converts3Archive() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// Extend runner to include want
type extendedRunner struct {
	runner
	want []map[string]interface{} // to verify the expected
}

// s3prepare helper function to prepare test cases
func s3prepare(name string, input map[string]interface{}, want []map[string]interface{}) extendedRunner {
	return extendedRunner{
		runner: runner{name: name, input: Node{ParameterMap: input}},
		want:   want,
	}
}

// Test function for ExtractEntries to make sure it has the parsing hierarchy
func TestExtractEntries(t *testing.T) {
	// Define test cases for ExtractEntries
	var tests []extendedRunner
	tests = append(tests, s3prepare("Valid entries", map[string]interface{}{
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
	tests = append(tests, s3prepare("Missing delegate", map[string]interface{}{
		// no "delegate" key
	}, nil))
	tests = append(tests, s3prepare("Invalid entries format", map[string]interface{}{
		"delegate": map[string]interface{}{
			"arguments": map[string]interface{}{
				"entries": "invalidFormat", // Not a slice of interface{}
			},
		},
	}, nil))

	// Execute each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractEntries(tt.input)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ExtractEntries() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
