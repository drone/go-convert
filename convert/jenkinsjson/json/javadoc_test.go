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

// Test function for dynamically testing ConvertJavadoc
func TestConvertJavadoc(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, "../convertTestFiles/javadoc/javadocSnippet.json")
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read JSON file: %v", err)
	}

	var rawNode map[string]interface{}
	if err := json.Unmarshal(jsonData, &rawNode); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	node := Node{}
	if err := mapToNodes(rawNode, &node); err != nil {
		t.Fatalf("failed to convert raw node to Node: %v", err)
	}

	javadocStep := ConvertJavadoc(node, nil)
	javadocStepPlugin, ok := javadocStep.Spec.(*harness.StepPlugin)
	if !ok {
		t.Fatalf("Expected Spec to be of type StepPlugin, but got %T", javadocStep.Spec)
	}

	expectedSource := extractJavadocDir(node)
	expectedStepPlugin := &harness.StepPlugin{
		With: map[string]interface{}{
			"aws_access_key_id":     "<+input>",
			"aws_secret_access_key": "<+input>",
			"aws_bucket":            "<+input>",
			"source":                expectedSource,
			"target":                "<+pipeline.name>/<+pipeline.sequenceId>",
		},
	}

	if diff := cmp.Diff(javadocStepPlugin.With, expectedStepPlugin.With); diff != "" {
		t.Errorf("Javadoc step properties mismatch (-want +got): %s", diff)
	}
}
