package json

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
)

type runner struct {
	name  string
	input Node
	want  *harness.Step
}

func prepare(t *testing.T, filename string, step *harness.Step) runner {

	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	jsonData, err := os.ReadFile(filepath.Join(workingDir, "../convertTestFiles/findFiles", filename+".json"))
	if err != nil {
		t.Fatalf("failed to read JSON file: %v", err)
	}

	var inputNode Node
	if err := json.Unmarshal(jsonData, &inputNode); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	return runner{
		name:  filename,
		input: inputNode,
		want:  step,
	}
}

func TestConvertFindFiles(t *testing.T) {

	var tests []runner
	tests = append(tests, prepare(t, "findFiles_GlobOnly", &harness.Step{
		Id:   "findFiles4d4efe",
		Name: "findFiles",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "harness-community/drone-findfiles:latest",
			With: map[string]interface{}{
				"glob": string("**/*.txt"),
			},
		},
	}))
	tests = append(tests, prepare(t, "findFiles_Excludes", &harness.Step{
		Id:   "findFilescc58db",
		Name: "findFiles",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "harness-community/drone-findfiles:latest",
			With: map[string]interface{}{
				"glob":     string("**/*.txt"),
				"excludes": string("**/1.txt"),
			},
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertFindFiles(tt.input)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertFindFiles() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
