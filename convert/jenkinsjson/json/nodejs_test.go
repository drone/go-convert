package json

import (
	"testing"

	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
)

func TestConvertNodejs(t *testing.T) {

	var tests []runner
	tests = append(tests, prepare(t, "/nodejs/nodejs_snippet", &harness.Step{
		Id:   "Stage_null99086b",
		Name: "Nodejs",
		Type: "script",
		Spec: &harness.StepExec{
			Connector: "<+input_docker_hub_connector>",
			Image:     "node:latest",
			Shell:     "sh",
			Run:       "npm --version\nnpm install\n",
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertNodejs(tt.input)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertNodejs() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
