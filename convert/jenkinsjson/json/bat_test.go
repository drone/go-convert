package json

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
)

func TestBat(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, "../convertTestFiles/bat/batSnippet.json")

	jsonData, err := os.ReadFile(filePath)

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
			name:      "Test Bat conversion",
			input:     node,
			variables: map[string]string{},
			want: &harness.Step{
				Name:    "bat",
				Timeout: "10m",
				Id:      "bat3b5f74",
				Type:    "script",
				Spec: &harness.StepExec{
					Shell: "pwsh",
					Run:   "echo @\"\necho hello world\n\"@ > bat3b5f74.bat\n ./bat3b5f74.bat",
				},
			},
		},
		{
			name:  "Test Bat conversion with variables",
			input: node,
			variables: map[string]string{
				"ENV_A": "ABC",
				"ENV_B": "DEF",
			},
			want: &harness.Step{
				Name:    "bat",
				Timeout: "10m",
				Id:      "bat3b5f74",
				Type:    "script",
				Spec: &harness.StepExec{
					Shell: "pwsh",
					Run:   "echo @\"\necho hello world\n\"@ > bat3b5f74.bat\n ./bat3b5f74.bat",
					Envs: map[string]string{
						"ENV_A": "ABC",
						"ENV_B": "DEF",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertBat(tt.input, tt.variables, "10m")
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertBat() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
