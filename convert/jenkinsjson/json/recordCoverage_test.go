package json

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"os"
	"path/filepath"
	"testing"
)

func TestRecordCoverage(t *testing.T) {
	jsonFilePath := "../convertTestFiles/recordCoverage/recordCoverage.json"
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

	tmpTestStepsList := ConvertRecordCoverage(node, make(map[string]string))
	if len(tmpTestStepsList) == 0 {
		t.Fatalf("Failed to convert JSON to struct")
	}

	for i, tmpTestStep := range tmpTestStepsList {
		var expectedMap map[string]interface{}
		if err := json.Unmarshal([]byte(expectedRecordCoverageStepsList[i]), &expectedMap); err != nil {
			t.Fatalf("Failed to convert expected JSON to map: %v", err)
		}
		jsonBytes, err := json.Marshal(tmpTestStep)
		if err != nil {
			t.Fatalf("Failed to marshal struct: %v", err)
		}
		var actualMap map[string]interface{}
		if err := json.Unmarshal(jsonBytes, &actualMap); err != nil {
			t.Fatalf("Failed to unmarshal JSON to map: %v", err)
		}
		if diff := cmp.Diff(expectedMap, actualMap); diff != "" {
			t.Fatalf("JSON comparison failed:\n%s", diff)
		}
	}
}

var expectedRecordCoverageStepsList = []string{
	`{"id":"recordCoverage6ecee1","name":"recordCoverage","type":"plugin","spec":{"image":"plugins/coverage-report","with":{"fail_on_threshold":true,"source_code_encoding":"UTF-8","threshold_branch":60,"threshold_line":60,"tool":"jacoco-xml"}}}`,
	`{"id":"recordCoverage6ecee1","name":"recordCoverage","type":"plugin","spec":{"image":"plugins/coverage-report","with":{"fail_on_threshold":true,"reports_path_pattern":"**/coverage.xml","source_code_encoding":"UTF-8","threshold_branch":60,"threshold_line":60,"tool":"cobertura"}}}`,
}
