package json

import (
	"testing"

	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
)

func TestConvertUnarchive(t *testing.T) {

	targetStr := string("unarchived_files/zip/")
	source := string("archive.zip")
	var tests []runner
	tests = append(tests, prepare(t, "/unarchive/unarchive_snippet", &harness.Step{
		Id:   "unarchive4bf5e3",
		Name: "Unarchive",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "plugins/archive",
			With: map[string]interface{}{
				"source":    source,
				"target":    targetStr,
				"format":    "zip",
				"action":    "extract",
				"glob":      "**/*",
				"overwrite": "true",
				"exclude":   "",
			},
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertUnarchive(tt.input, tt.input.ParameterMap)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertUnarchive() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
