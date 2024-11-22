package json

import (
	"encoding/json"
	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
	"os"
	"path/filepath"
	"testing"
)

func TestNexusUploader(t *testing.T) {
	variables := make(map[string]string)
	jsonFilePath := "../convertTestFiles/nexusArtifactUploader/nexusArtifactUploader.json"

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

	tmpTestStep := ConvertNexusArtifactUploader(node, variables)

	wantStep, err := ToStructFromJsonString[harness.Step](expectedNexusUploaderStepJSON)
	if err != nil {
		t.Fatalf("want step : %v", err)
	}

	diffs := cmp.Diff(wantStep, *tmpTestStep)

	if len(diffs) != 0 {
		t.Fatalf("Failed to convert JSON to struct: %v", diffs)
	}

}

var expectedNexusUploaderStepJSON = `{
    "id": "nexusArtifactUploadera3af28",
    "name": "nexusArtifactUploader",
    "type": "plugin",
    "spec": {
        "image": "harnesscommunity/drone-nexus-publish",
        "with": {
            "artifacts": "[{\"artifactId\":\"fruit-names-v0.1\",\"classifier\":\"\",\"file\":\"fruit_names.txt\",\"type\":\"txt\"}]",
            "group_id": "test01",
            "nexus_version": "nexus3",
            "password": "add password",
            "protocol": "http",
            "repository": "plg01maven2",
            "server_url": "4.24.1.2:9192/",
            "username": "add user name"
        },
        "inputs": {
            "artifacts": "[{\"artifactId\":\"fruit-names-v0.1\",\"classifier\":\"\",\"file\":\"fruit_names.txt\",\"type\":\"txt\"}]",
            "group_id": "test01",
            "nexus_version": "nexus3",
            "password": "add password",
            "protocol": "http",
            "repository": "plg01maven2",
            "server_url": "4.24.1.2:9192/",
            "username": "add user name"
        }
    }
}`
