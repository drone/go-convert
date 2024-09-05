package json

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
)

func TestConvertPowerShell(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}

	filePath := filepath.Join(workingDir, "../convertTestFiles/powershell/powershellSnippet.json")

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
			name:      "Test ReadJSON conversion",
			input:     node,
			variables: map[string]string{},
			want: &harness.Step{
				Name:    "powershell",
				Timeout: "10m",
				Id:      "powershell1d3281",
				Type:    "script",
				Spec: &harness.StepExec{
					Shell: "Powershell",
					Run:   "\n                Write-Output \"Running a PowerShell command\"\n                ",
				},
			},
		},
		{
			name:  "Test ReadJSON conversion with variables",
			input: node,
			variables: map[string]string{
				"ENV_A": "ABC",
				"ENV_B": "DEF",
			},
			want: &harness.Step{
				Name:    "powershell",
				Timeout: "10m",
				Id:      "powershell1d3281",
				Type:    "script",
				Spec: &harness.StepExec{
					Shell: "Powershell",
					Run:   "\n                Write-Output \"Running a PowerShell command\"\n                ",
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
			got := ConvertPowerShell(tt.input, tt.variables, "10m")
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertPowerShell() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
