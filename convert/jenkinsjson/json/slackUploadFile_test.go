package json

import (
	"encoding/json"
	"fmt"
	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
	"os"
	"path/filepath"
	"testing"
)

func TestSlackUploadFile(t *testing.T) {

	jsonFilePath := "../convertTestFiles/slackUploadFile/slackUploadFile.json"
	pluginImageName := SlackPluginImage
	jenkinsToDroneMapperList := JenkinsToDroneSlackUploadParamMapperList
	expectedStepJson := expectedSlackFileUploadStepJSON

	isDebug := func() bool {
		return false
	}

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
	variables := make(map[string]string)

	tmpTestStep := ConvertToStepWithProperties(&node, variables, jenkinsToDroneMapperList,
		pluginImageName)

	if isDebug() {
		js, er := ToJsonStringFromStruct[harness.Step](*tmpTestStep)
		if er != nil {
			t.Fatalf("failed to convert struct to JSON: %v", er)
		}
		fmt.Println(js)
	}

	wantStep, err := ToStructFromJsonString[harness.Step](expectedStepJson)
	if err != nil {
		t.Fatalf("want step : %v", err)
	}

	diffs := cmp.Diff(wantStep, *tmpTestStep)

	if len(diffs) != 0 {
		t.Fatalf("failed to convert JSON to struct: %v", diffs)
	}

}

var expectedSlackFileUploadStepJSON = `{
  "id": "slackUploadFileb5f507",
  "name": "slackUploadFile",
  "type": "plugin",
  "spec": {
    "image": "plugins/slack",
    "with": {
      "channel": "C07TL1KNV8Q",
      "access_token": "SlackToken01",
      "fail_on_error": true,
      "file_path": "b.txt",
      "initial_comment": "jenkins file upload test"
    },
    "inputs": {
      "channel": "C07TL1KNV8Q",
      "credential_id": "SlackToken01",
      "fail_on_error": true,
      "file_path": "b.txt",
      "initial_comment": "jenkins file upload test"
    }
  }
}`
