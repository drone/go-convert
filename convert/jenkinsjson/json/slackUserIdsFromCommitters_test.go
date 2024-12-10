package json

import (
	"encoding/json"
	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
	"os"
	"path/filepath"
	"testing"
)

func TestSlackUserIdsFromCommitters(t *testing.T) {
	variables := make(map[string]string)
	jsonFilePath := "../convertTestFiles/slackUserIdsFromCommitters/slackUserIdsFromCommittersSnippet.json"

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

	tmpTestStep := ConvertSlackUserIdFromEmail(node, variables)
	wantStep, err := ToStructFromJsonString[harness.Step](expectedSlackUserIdsFromCommitters)
	if err != nil {
		t.Fatalf("want step : %v", err)
	}

	diffs := cmp.Diff(wantStep, *tmpTestStep)

	if len(diffs) != 0 {
		t.Fatalf("Failed to convert JSON to struct: %v", diffs)
	}

}

var expectedSlackUserIdsFromCommitters = `{
    "id": "slackUserIdsFromCommitters0bf200",
    "name": "slackUserIdsFromCommitters",
    "type": "plugin",
    "spec": {
        "image": "plugins/slack",
        "with": {
            "access_token": "\u003c+input\u003e"
        },
        "inputs": {
            "access_token": "\u003c+input\u003e"
        }
    }
}`
