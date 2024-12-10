package json

import (
	"testing"

	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
)

func TestConvertTestNG(t *testing.T) {

	var tests []runner
	tests = append(tests, prepare(t, "/testng/testng_snippet", &harness.Step{
		Id:   "testNG80776e",
		Name: "testng",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "plugins/testng",
			With: map[string]interface{}{
				"failed_fails":                  100,
				"failed_skips":                  100,
				"failure_on_failed_test_config": false,
				"level":                         "info",
				"report_filename_pattern":       "**/testng-results.xml",
				"threshold_mode":                "absolute",
			},
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertTestng(tt.input, tt.input.ParameterMap)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertTestng() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
