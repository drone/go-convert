package json

import (
	"testing"

	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
)

func TestConvertCucumber(t *testing.T) {
	var tests []runner

	tests = append(tests, prepare(t, "../convertTestFiles/cucumber/cucumber_snippet", &harness.Step{
		Id:   "cucumberf61234",
		Name: "cucumber",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "plugins/cucumber",
			With: map[string]interface{}{
				"file_include_pattern":  "**/*.json",
				"json_report_directory": "target",
				"sorting_method":        "ALPHABETICAL",
				"skip_empty_json_files": true,
				"level":                 "info",
			},
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertCucumber(tt.input, tt.input.ParameterMap)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertCucumber() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
