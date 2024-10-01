package json

import (
	"encoding/json"
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

	var node1 Node
	if err := json.Unmarshal(jsonData, &node1); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	pmStr, err := ToJsonStringFromMap[map[string]interface{}](node1.ParameterMap)
	if err != nil {
		t.Fatalf("failed to unmarshal harness-attribute for node %s: %v", node1.SpanName, err)
	}

	toTestParameterMap, err := ToStructFromJsonString[ParameterMap](pmStr)
	if err != nil {
		t.Fatalf("failed to unmarshal harness-attribute for node %s: %v", node1.SpanName, err)
	}

	diffs := cmp.Diff(toTestParameterMap, want)

	if len(diffs) != 0 {
		t.Fatalf("failed to convert JSON to struct: %v", diffs)
	}

}

var want = ParameterMap{
	ValidResponseCodes:     "200:299",
	HttpMode:               "POST",
	WrapAsMultipart:        false,
	Url:                    "https://jsonplaceholder.typicode.com/posts",
	Timeout:                60,
	ValidResponseContent:   "\"id\":",
	OutputFile:             "response.json",
	IgnoreSslErrors:        true,
	RequestBody:            "{\"title\": \"foo\", \"body\": \"bar\", \"userId\": 1}",
	ConsoleLogResponseBody: true,
	Quiet:                  false,
	ContentType:            "APPLICATION_JSON",
	CustomHeaders: []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}{
		{
			Name:  "Authorization",
			Value: "Bearer <token>",
		},
		{
			Name:  "X-Custom-Header",
			Value: "example-header",
		},
	},
	Authentication: "httprequest",
	AcceptType:     "APPLICATION_JSON",
	UploadFile:     "/tmp/sample_file.txt",
	MultipartName:  "sample_upload_file",
}

type ParameterMap struct {
	ValidResponseCodes     string `json:"validResponseCodes"`
	HttpMode               string `json:"httpMode"`
	WrapAsMultipart        bool   `json:"wrapAsMultipart"`
	Url                    string `json:"url"`
	Timeout                int    `json:"timeout"`
	ValidResponseContent   string `json:"validResponseContent"`
	OutputFile             string `json:"outputFile"`
	IgnoreSslErrors        bool   `json:"ignoreSslErrors"`
	RequestBody            string `json:"requestBody"`
	ConsoleLogResponseBody bool   `json:"consoleLogResponseBody"`
	Quiet                  bool   `json:"quiet"`
	ContentType            string `json:"contentType"`
	CustomHeaders          []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"customHeaders"`
	Authentication string `json:"authentication"`
	AcceptType     string `json:"acceptType"`
	UploadFile     string `json:"uploadFile"`
	MultipartName  string `json:"multipartName"`
}
