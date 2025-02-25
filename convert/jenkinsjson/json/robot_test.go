package json

import (
	"testing"

	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
)

func TestConvertRobot(t *testing.T) {
	var tests []runner

	tests = append(tests, prepare(t, "../convertTestFiles/robot/robot_snippet", &harness.Step{
		Id:   "robot27f6e9",
		Name: "robot",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "plugins/robot",
			With: map[string]interface{}{
				"count_skipped_tests":      true,
				"only_critical":            false,
				"pass_threshold":           90,
				"report_directory":         "results",
				"report_file_name_pattern": "output.xml",
				"unstable_threshold":       70,
				"level":                    "info",
			},
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertRobot(tt.input, tt.input.ParameterMap)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertCucumber() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
