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

func TestSlackSend(t *testing.T) {

	jsonFilePath := "../convertTestFiles/slackSend/slackSend.json"
	pluginImageName := SlackPluginImage
	jenkinsToDroneMapperList := JenkinsToDroneSlackSendParamMapperList
	expectedStepJson := expectedSlackSendStepJSON

	isDebug := func() bool {
		return true
	}

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
	variables := make(map[string]string)

	tmpTestStep := ConvertToStepWithProperties(&node, variables, jenkinsToDroneMapperList,
		pluginImageName)

	if isDebug() {
		js, er := ToJsonStringFromStruct[harness.Step](*tmpTestStep)
		if er != nil {
			t.Fatalf("Failed to convert struct to JSON: %v", er)
		}
		fmt.Println(js)
	}

	wantStep, err := ToStructFromJsonString[harness.Step](expectedStepJson)
	if err != nil {
		t.Fatalf("want step : %v", err)
	}

	diffs := cmp.Diff(wantStep, *tmpTestStep)

	if len(diffs) != 0 {
		t.Fatalf("Failed to convert JSON to struct: %v", diffs)
	}

}

var expectedSlackSendStepJSON = `{
    "id": "slackSend8962b0",
    "name": "slackSend",
    "type": "plugin",
    "spec": {
        "image": "plugins/slack",
        "with": {
            "access_token": "slackcreds03",
            "channel": "dev",
            "color": "good",
            "committers_slack_id": true,
            "icon_emoji": ":rocket:",
            "message": "hi, this is senthil-slack-file-send-03_dec10",
            "username": "jenkins-test-user"
        },
        "inputs": {
            "access_token": "slackcreds03",
            "channel": "dev",
            "color": "good",
            "committers_slack_id": true,
            "icon_emoji": ":rocket:",
            "message": "hi, this is senthil-slack-file-send-03_dec10",
            "username": "jenkins-test-user"
        }
    }
}`
