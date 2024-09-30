package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConverts3Upload(t *testing.T) {
	// Get the working directory
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	// Update file path according to the location of fileOps_test.json
	filePath := filepath.Join(workingDir, "../convertTestFiles/s3publisher/s3upload/s3upload_snippet.json")
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read JSON file: %v", err)
	}

	// Unmarshal the JSON into a Node struct
	var node1 Node
	if err := json.Unmarshal(jsonData, &node1); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	tests := []struct {
		json Node
		want Node
	}{
		{
			json: node1,
			want: Node{
				SpanId:  "1ec902dbd864116c",
				TraceId: "f03b61a089e2ec4d5e81778ca44190e8",
				Parent:  "S3 Publisher-1",
				Name:    "S3 Publisher-1 #2",
				AttributesMap: map[string]string{
					"harness-attribute-extra-pip: org.jenkinsci.plugins.workflow.actions.TimingAction@593f0a9": "org.jenkinsci.plugins.workflow.actions.TimingAction@593f0a9",
					"harness-others":             "-delegate-field org.jenkinsci.plugins.workflow.steps.CoreStep delegate-org.jenkinsci.plugins.workflow.steps.CoreStep.delegate-interface jenkins.tasks.SimpleBuildStep",
					"jenkins.pipeline.step.name": "Publish artifacts to S3 Bucket",
					"ci.pipeline.run.user":       "SYSTEM",
					"jenkins.pipeline.step.id":   "7",
					"harness-attribute-extra-pip: io.jenkins.plugins.opentelemetry.MigrateHarnessUrlChildAction@2e8f827b": "io.jenkins.plugins.opentelemetry.MigrateHarnessUrlChildAction@2e8f827b",
					"jenkins.pipeline.step.type":           "s3Upload",
					"harness-attribute":                    "{\n  \"delegate\" : {\n    \"symbol\" : \"s3Upload\",\n    \"klass\" : null,\n    \"arguments\" : {\n      \"consoleLogLevel\" : \"INFO\",\n      \"dontSetBuildResultOnFailure\" : false,\n      \"dontWaitForConcurrentBuildCompletion\" : false,\n      \"entries\" : [ {\n        \"bucket\" : \"bucket-1\",\n        \"excludedFile\" : \"2.txxt\",\n        \"flatten\" : false,\n        \"gzipFiles\" : false,\n        \"keepForever\" : false,\n        \"managedArtifacts\" : false,\n        \"noUploadOnFailure\" : false,\n        \"selectedRegion\" : \"us-west-1\",\n        \"showDirectlyInBrowser\" : false,\n        \"sourceFile\" : \"*.txt\",\n        \"storageClass\" : \"STANDARD\",\n        \"uploadFromSlave\" : false,\n        \"useServerSideEncryption\" : false\n      } ],\n      \"pluginFailureResultConstraint\" : \"FAILURE\",\n      \"profileName\" : \"default\",\n      \"userMetadata\" : [ ]\n    },\n    \"model\" : null\n  }\n}",
					"jenkins.pipeline.step.plugin.name":    "s3",
					"jenkins.pipeline.step.plugin.version": "483.vcb_db_3dcee68f",
				},
				ParentSpanId: "31b6240c3095e9f3",
				SpanName:     "s3Upload",
				Type:         "Run Phase Span",
				ParameterMap: map[string]interface{}{
					"delegate": map[string]interface{}{
						"symbol": "s3Upload",
						"klass":  nil,
						"arguments": map[string]interface{}{
							"profileName":                          "default",
							"dontSetBuildResultOnFailure":          false,
							"dontWaitForConcurrentBuildCompletion": false,
							"entries": []interface{}{
								map[string]interface{}{
									"excludedFile":            "2.txxt",
									"uploadFromSlave":         false,
									"managedArtifacts":        false,
									"keepForever":             false,
									"gzipFiles":               false,
									"sourceFile":              "*.txt",
									"bucket":                  "bucket-1",
									"flatten":                 false,
									"noUploadOnFailure":       false,
									"storageClass":            "STANDARD",
									"selectedRegion":          "us-west-1",
									"showDirectlyInBrowser":   false,
									"useServerSideEncryption": false,
								},
							},
							"consoleLogLevel":               "INFO",
							"pluginFailureResultConstraint": "FAILURE",
							"userMetadata":                  []interface{}{},
						},
						"model": nil,
					},
				},
			},
		},
	}

	for i, test := range tests {
		got := test.json
		if diff := cmp.Diff(got, test.want); diff != "" {
			t.Errorf("Unexpected parsing results for test %v", i)
			t.Log(diff)
		}
	}
}
