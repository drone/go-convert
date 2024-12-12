package json

import (
	"encoding/json"
	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
	"os"
	"path/filepath"
	"testing"
)

func TestSlackUserIdFromEmail(t *testing.T) {
	variables := make(map[string]string)
	jsonFilePath := "../convertTestFiles/slackUserIdFromEmail/slackUserIdFromEmailSnippet.json"

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

	wantStep, err := ToStructFromJsonString[harness.Step](expectedSlackUserIdFromEmail)
	if err != nil {
		t.Fatalf("want step : %v", err)
	}

	diffs := cmp.Diff(wantStep, *tmpTestStep)

	if len(diffs) != 0 {
		t.Fatalf("Failed to convert JSON to struct: %v", diffs)
	}

}

var expectedSlackUserIdFromEmail = `{
    "id": "slackUserIdFromEmail976043",
    "name": "slackUserIdFromEmail",
    "type": "plugin",
    "spec": {
        "image": "plugins/slack",
        "with": {
            "access_token": "\u003c+input\u003e",
            "slack_user_email_id": "slack_user_tst223@gmail.com"
        },
        "inputs": {
            "access_token": "\u003c+input\u003e",
            "slack_user_email_id": "slack_user_tst223@gmail.com"
        }
    }
}`
