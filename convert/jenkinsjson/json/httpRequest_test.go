package json

import (
	"encoding/json"
	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
	"os"
	"path/filepath"
	"testing"
)

func TestHttpRequest(t *testing.T) {

	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, "../convertTestFiles/httpRequest/httpRequest.json")

	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read JSON file: %v", err)
	}

	var node Node
	if err := json.Unmarshal(jsonData, &node); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	tmpTestStep := GetStepWithProperties(&node, JenkinsToDroneParamMapperList, HttpRequestPluginImage)

	wantStep, err := ToStructFromJsonString[harness.Step](expectedStepJSON)
	if err != nil {
		t.Fatalf("want step : %v", err)
	}

	diffs := cmp.Diff(wantStep, *tmpTestStep)

	if len(diffs) != 0 {
		t.Fatalf("failed to convert JSON to struct: %v", diffs)
	}

}

var expectedStepJSON = `{
    "id": "httpRequestbf3ce3",
    "name": "httpRequest",
    "type": "plugin",
    "spec": {
        "image": "plugins/httpRequest",
        "with": {
            "accept_type": "APPLICATION_JSON",
            "auth_basic": "httprequest",
            "content_type": "APPLICATION_JSON",
            "headers": "Authorization:Bearer \u003ctoken\u003e,X-Custom-Header:example-header",
            "http_method": "POST",
            "ignore_ssl": true,
            "log_response": true,
            "multipart_name": "sample_upload_file",
            "output_file": "response.json",
            "quiet": false,
            "request_body": "{\"title\": \"foo\", \"body\": \"bar\", \"userId\": 1}",
            "timeout": 60,
            "upload_file": "sample_file.txt",
            "url": "https://jsonplaceholder.typicode.com/posts",
            "valid_response_body": "\"id\":",
            "valid_response_codes": "200:299",
            "wrap_as_multipart": false
        },
        "inputs": {
            "accept_type": "APPLICATION_JSON",
            "auth_basic": "httprequest",
            "content_type": "APPLICATION_JSON",
            "headers": "Authorization:Bearer \u003ctoken\u003e,X-Custom-Header:example-header",
            "http_method": "POST",
            "ignore_ssl": true,
            "log_response": true,
            "multipart_name": "sample_upload_file",
            "output_file": "response.json",
            "quiet": false,
            "request_body": "{\"title\": \"foo\", \"body\": \"bar\", \"userId\": 1}",
            "timeout": 60,
            "upload_file": "sample_file.txt",
            "url": "https://jsonplaceholder.typicode.com/posts",
            "valid_response_body": "\"id\":",
            "valid_response_codes": "200:299",
            "wrap_as_multipart": false
        }
    }
}`
