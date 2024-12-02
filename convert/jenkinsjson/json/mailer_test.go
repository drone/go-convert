package json

import (
	"testing"

	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
)

func TestConvertMailer(t *testing.T) {

	var tests []runner
	tests = append(tests, prepare(t, "/mailer/mailer_snippet", &harness.Step{
		Id:   "mail0fb620",
		Name: "Mailer",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "plugins/email",
			With: map[string]interface{}{
				"body":         string("\n                    Build Status: FAILURE\n                    Check console output at: http://localhost:8080/job/Mail_Pipeline/15/console\n                    "),
				"from.address": string("<+input>"),
				"host":         string("<+input>"),
				"password":     string("<+input>"),
				"port":         string("<+input>"),
				"recipients":   string("test.user@testmail.com"),
				"subject":      string("Jenkins Build - Mail_Pipeline #15"),
				"username":     string("<+input>"),
			},
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertMailer(tt.input, tt.input.ParameterMap)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertMailer() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
