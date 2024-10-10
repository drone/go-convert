package json

import (
	"testing"

	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
)

func TestConvertFindFiles(t *testing.T) {

	var tests []runner
	tests = append(tests, prepare(t, "findFiles/findFiles_GlobOnly", &harness.Step{
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
	tests = append(tests, prepare(t, "findFiles/findFiles_Excludes", &harness.Step{
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
