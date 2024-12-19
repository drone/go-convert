package json

import (
	"encoding/json"
	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
	"os"
	"path/filepath"
	"testing"
)

func TestFlywayRunner(t *testing.T) {
	variables := make(map[string]string)
	jsonFilePath := "../convertTestFiles/flywayrunner/flywayrunnerSnippet.json"

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

	tmpTestStep := ConvertFlywayRunner(node, variables)

	wantStep, err := ToStructFromJsonString[harness.Step](expectedFywayRunnerStep)
	if err != nil {
		t.Fatalf("want step : %v", err)
	}

	diffs := cmp.Diff(wantStep, *tmpTestStep)

	if len(diffs) != 0 {
		t.Fatalf("Failed to convert JSON to struct: %v", diffs)
	}
}

var expectedFywayRunnerStep = `{
    "id": "flywayrunner092504",
    "name": "flywayrunner",
    "type": "plugin",
    "spec": {
        "image": "plugins/flyway",
        "with": {
            "command_line_args": "-X",
            "flyway_command": "migrate",
            "locations": "/opt/hns/harness-plugins/flyway-test-files/migration_files",
            "password": "\u003c+input\u003e",
            "url": "jdbc:mysql://43.204.190.241:3306/flyway_test",
            "username": "\u003c+input\u003e"
        },
        "inputs": {
            "command_line_args": "-X",
            "flyway_command": "migrate",
            "locations": "/opt/hns/harness-plugins/flyway-test-files/migration_files",
            "password": "\u003c+input\u003e",
            "url": "jdbc:mysql://43.204.190.241:3306/flyway_test",
            "username": "\u003c+input\u003e"
        }
    }
}`
