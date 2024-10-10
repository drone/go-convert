package json

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	harness "github.com/drone/spec/dist/go"
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

	jsonData, err := os.ReadFile(filepath.Join(workingDir, "../convertTestFiles/", filename+".json"))
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
