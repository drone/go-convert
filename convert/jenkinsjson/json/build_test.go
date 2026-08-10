package json

import (
	"testing"

	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
)

func TestConvertBuild(t *testing.T) {
	var tests []runner

	tests = append(tests, prepare(t, "../convertTestFiles/build/build_snippet", &harness.Step{
		Id:      "build_JobB0e03d3",
		Name:    "build JobB",
		Type:    "script",
		Timeout: "10",
		Spec: &harness.StepExec{
			Image: "immich",
			Shell: "sh",
			Run:   "curl https://replace_with_webhook_url",
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertBuild(tt.input, map[string]string{}, "10", "immich")

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ConvertBuild() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
