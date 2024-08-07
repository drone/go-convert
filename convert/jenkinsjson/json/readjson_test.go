package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	harness "github.com/drone/spec/dist/go"
)

func TestConvertReadJson(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, "../convertTestFiles/readjson/readjsonSnippet.json")
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read JSON file: %v", err)
	}

	var node Node
	if err := json.Unmarshal(jsonData, &node); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	tests := []struct {
		name      string
		input     Node
		variables map[string]string
		want      *harness.Step
	}{
		{
			name:      "Test ReadJSON conversion",
			input:     node,
			variables: map[string]string{},
			want: &harness.Step{
				Name: "readJSON",
				Id:   "readJSON721ffe",
				Type: "script",
				Spec: &harness.StepExec{
					Image:   "alpine",
					Shell:   "sh",
					Run:     "jsonObj='$(cat /Users/rakshith/Downloads/IntermediateJson/BasicPipe.json | tr -d '\\n')'", 
					Outputs: []string{"jsonObj"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertReadJson(tt.input, tt.variables)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertReadJson() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}