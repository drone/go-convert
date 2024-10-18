package json

import (
	"encoding/json"
	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
	"os"
	"path/filepath"
	"testing"
)

func TestJacocoCoverage(t *testing.T) {

	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, "../convertTestFiles/jacoco/jacoco.json")

	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read JSON file: %v", err)
	}

	var node Node
	if err := json.Unmarshal(jsonData, &node); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}
	variables := make(map[string]string)

	tmpTestStep := ConvertToStepUsingParameterMapDelegate(&node, variables, JacocoJenkinsToDroneParamMapperList,
		CoverageReportImage)

	wantStep, err := ToStructFromJsonString[harness.Step](expectedJacocoStepJSON)
	if err != nil {
		t.Fatalf("want step : %v", err)
	}

	diffs := cmp.Diff(wantStep, *tmpTestStep)

	if len(diffs) != 0 {
		t.Fatalf("failed to convert JSON to struct: %v", diffs)
	}

}

var expectedJacocoStepJSON = `{
    "id": "jacoco66be31",
    "name": "jacoco",
    "type": "plugin",
    "spec": {
        "image": "plugins/coverage-report",
        "with": {
            "class_directories": "**/target/classes",
            "reports_path_pattern": "**/**.exec",
            "source_directories": "**/src/main",
            "tool": "jacoco"
        },
        "inputs": {
            "class_directories": "**/target/classes",
            "reports_path_pattern": "**/**.exec",
            "source_directories": "**/src/main",
            "tool": "jacoco"
        }
    }
}`
