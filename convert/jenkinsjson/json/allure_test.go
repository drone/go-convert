package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
)

// Test function for both Allure and S3 steps
func TestConvertSteps(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	// Define the test cases in a struct
	tests := []struct {
		fileName       string
		expectedAllure harness.StepExec
		expectedS3     harness.StepPlugin
	}{
		{
			fileName: "../convertTestFiles/allure/allureSnippet.json",
			expectedAllure: harness.StepExec{
				Run:   "allure generate build/allure-results --clean --output allure-report",
				Image: "harnesscommunity/allure:jdk21",
				Shell: "Sh",
			},
			expectedS3: harness.StepPlugin{
				With: map[string]interface{}{
					"aws_bucket":            "<+input>",
					"source":                "allure-report",
					"artifact_file":         "artifact.txt",
					"aws_access_key_id":     "<+input>",
					"aws_secret_access_key": "<+input>",
					"target":                "<+pipeline.sequenceId>",
				},
			},
		},
		{
			fileName: "../convertTestFiles/allure/allureMultiSnippet.json",
			expectedAllure: harness.StepExec{
				Run:   "allure generate build/allure-results1,build/allure-results2,build/allure-results3 --clean --output allure-report",
				Image: "harnesscommunity/allure:jdk21",
				Shell: "Sh",
			},
			expectedS3: harness.StepPlugin{
				With: map[string]interface{}{
					"aws_bucket":            "<+input>",
					"source":                "allure-report",
					"artifact_file":         "artifact.txt",
					"aws_access_key_id":     "<+input>",
					"aws_secret_access_key": "<+input>",
					"target":                "<+pipeline.sequenceId>",
				},
			},
		},
	}

	for i, test := range tests {
		filePath := filepath.Join(workingDir, test.fileName)
		jsonData, err := ioutil.ReadFile(filePath)
		if err != nil {
			t.Fatalf("failed to read JSON file: %v", err)
		}

		var rawNode map[string]interface{}
		if err := json.Unmarshal(jsonData, &rawNode); err != nil {
			t.Fatalf("failed to decode JSON: %v", err)
		}

		node1 := Node{}
		if err := mapToNodes(rawNode, &node1); err != nil {
			t.Fatalf("failed to convert raw node to Node: %v", err)
		}

		// Allure step test
		allureStep := ConvertAllureReport(node1)

		allureStepExec, ok := allureStep.Spec.(*harness.StepExec)
		if !ok {
			t.Fatalf("Expected Spec to be of type StepExec, but got %T", allureStep.Spec)
		}

		if diff := cmp.Diff(allureStepExec, &test.expectedAllure); diff != "" {
			t.Errorf("Allure step execution mismatch (-want +got) for test %d: %s", i, diff)
		}

		// S3 step test
		s3Step := Converts3UploadStep(node1)
		s3StepPlugin, ok := s3Step.Spec.(*harness.StepPlugin)
		if !ok {
			t.Fatalf("Expected Spec to be of type StepPlugin, but got %T", s3Step.Spec)
		}

		if diff := cmp.Diff(s3StepPlugin.With, test.expectedS3.With); diff != "" {
			t.Errorf("S3 step properties mismatch (-want +got) for test %d: %s", i, diff)
		}

	}
}
