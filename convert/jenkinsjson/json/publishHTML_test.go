package json

import (
	"encoding/json"
	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
	"os"
	"path/filepath"
	"testing"
)

func TestPublishHTML(t *testing.T) {

	jsonFilePath := "../convertTestFiles/publishHTML/publishHTML.json"

	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, jsonFilePath)

	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read JSON file: %v", err)
	}

	var node Node
	if err := json.Unmarshal(jsonData, &node); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	tmpTestStep := ConvertPublishHtml(node, nil)
	if tmpTestStep == nil {
		t.Fatalf("failed to convert JSON to struct")
	}

	wantStep, err := ToStructFromJsonString[harness.Step](expectedPublishHTMLStepJSON)
	if err != nil {
		t.Fatalf("want step : %v", err)
	}

	diffs := cmp.Diff(wantStep, *tmpTestStep)

	if len(diffs) != 0 {
		t.Fatalf("failed to convert JSON to struct: %v", diffs)
	}
}

var expectedPublishHTMLStepJSON = `{
    "id": "UploadPublish43d688",
    "name": "Upload and Publish",
    "type": "plugin",
    "spec": {
        "image": "harnesscommunity/drone-s3-upload-publish",
        "connector": "\u003c+input\u003e",
        "with": {
            "artifact_file": "artifact.txt",
            "aws_access_key_id": "\u003c+input\u003e",
            "aws_bucket": "\u003c+input\u003e",
            "aws_default_region": "\u003c+input\u003e",
            "aws_secret_access_key": "\u003c+input\u003e",
            "include": "**/*.html,**/*.css",
            "source": "reports",
            "target": "\u003c+pipeline.sequenceId\u003e"
        }
    }
}`
