package json

import (
	"encoding/json"
	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
	"os"
	"path/filepath"
	"testing"
)

const expectedTestResultsAggregatorStep = `{
    "id": "testResultsAggregator01c08d",
    "name": "testResultsAggregator",
    "type": "plugin",
    "spec": {
        "image": "plugins/test-results-aggregator",
        "with": {
            "group": "\u003c+input\u003e",
            "include_pattern": "\u003c+input\u003e",
            "reports_dir": "\u003c+input\u003e",
            "tool": "\u003c+input\u003e"
        },
        "inputs": {
            "group": "\u003c+input\u003e",
            "include_pattern": "\u003c+input\u003e",
            "reports_dir": "\u003c+input\u003e",
            "tool": "\u003c+input\u003e"
        }
    }
}`

func TestResultsAggregatorStep(t *testing.T) {
	jsonFilePath := "../convertTestFiles/testResultsAggregator/testResultsAggregator_test.json"

	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, jsonFilePath)

	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	var node Node
	if err := json.Unmarshal(jsonData, &node); err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	tmpTestStep := ConvertTestResultsAggregator(node, nil)

	wantStep, err := ToStructFromJsonString[harness.Step](expectedTestResultsAggregatorStep)
	if err != nil {
		t.Fatalf("want step : %v", err)
	}

	diffs := cmp.Diff(wantStep, *tmpTestStep)

	if len(diffs) != 0 {
		t.Fatalf("Failed to convert JSON to struct: %v", diffs)
	}
}
