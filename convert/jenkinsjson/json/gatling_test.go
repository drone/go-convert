package json

import (
	"testing"

	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
)

func TestConvertGatling(t *testing.T) {

	var tests []runner
	tests = append(tests, prepare(t, "/gatling/gatling_snippet", &harness.Step{
		Id:   "gatlingArchivec8e1d8",
		Name: "Gatling_Publish",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "harnesscommunity/drone-s3-upload-publish",
			With: map[string]interface{}{
				"artifact_file":         "artifact.txt",
				"aws_access_key_id":     "<+input>",
				"aws_bucket":            "<+input>",
				"aws_default_region":    "<+input>",
				"aws_secret_access_key": "<+input>",
				"glob":                  "**/*.html, **/*.css",
				"source":                "<+input>",
				"target":                "<+input>",
			},
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertGatling(tt.input)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertGatling() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
