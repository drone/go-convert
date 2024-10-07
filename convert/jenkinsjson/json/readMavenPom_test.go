package json

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
)

type readMavenPomRunner struct {
	name  string
	input Node
	want  *harness.Step
}

func prepareReadMavenPom(t *testing.T, filename string, step *harness.Step) readMavenPomRunner {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	jsonData, err := os.ReadFile(filepath.Join(workingDir, "../convertTestFiles/readMavenPom", filename))
	if err != nil {
		t.Fatalf("failed to read JSON file: %v", err)
	}

	var inputNode Node
	if err := json.Unmarshal(jsonData, &inputNode); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	return readMavenPomRunner{
		name:  filename,
		input: inputNode,
		want:  step,
	}
}

func TestConvertReadMavenPom(t *testing.T) {
	// Preparing the test runner with `readMavenPomSnippet.json`
	tests := []readMavenPomRunner{
		prepareReadMavenPom(t, "readMavenPomSnippet.json", &harness.Step{
			Id:   "readMavenPomf2e306",
			Name: "readMavenPom",
			Type: "plugin",
			Spec: &harness.StepPlugin{
				Image: "harnesscommunity/drone-get-maven-version:latest",
				With: map[string]interface{}{
					"pom_path": "pom.xml",
				},
			},
		}),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertReadMavenPom(tt.input)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertReadMavenPom() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
