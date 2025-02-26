package json

import (
	"encoding/json"
	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
	"os"
	"path/filepath"
	"testing"
)

const (
	CoberturaNodeDef = "../convertTestFiles/cobertura/cobertura.json"
)

func TestCoberturaCoverage(t *testing.T) {

	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, CoberturaNodeDef)

	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read JSON file: %v", err)
	}

	var node Node
	if err := json.Unmarshal(jsonData, &node); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}
	variables := make(map[string]string)

	tmpTestStep := ConvertToStepUsingParameterMapDelegate(&node, variables, CorberturaJenkinsToDroneParamMapperList,
		CoverageReportImage)

	wantStep, err := ToStructFromJsonString[harness.Step](expectedCoberturaStepJSON)
	if err != nil {
		t.Fatalf("want step : %v", err)
	}

	diffs := cmp.Diff(wantStep, *tmpTestStep)

	if len(diffs) != 0 {
		t.Fatalf("failed to convert JSON to struct: %v", diffs)
	}

}

var expectedCoberturaStepJSON = `{
    "id": "cobertura5133f5",
    "name": "cobertura",
    "type": "plugin",
    "spec": {
        "image": "plugins/coverage-report",
        "with": {
            "fail_on_threshold": false,
            "reports_path_pattern": "**/coverage*.xml",
            "threshold_branch": 75,
            "threshold_class": 75,
            "threshold_file": 75,
            "threshold_line": 75,
            "threshold_method": 75,
            "threshold_package": 75,
            "tool": "cobertura"
        },
        "inputs": {
            "fail_on_threshold": false,
            "reports_path_pattern": "**/coverage*.xml",
            "threshold_branch": 75,
            "threshold_class": 75,
            "threshold_file": 75,
            "threshold_line": 75,
            "threshold_method": 75,
            "threshold_package": 75,
            "tool": "cobertura"
        }
    }
}`
