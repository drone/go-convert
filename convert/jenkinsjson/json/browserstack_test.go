package json

import (
	"strings"
	"testing"

	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
)

func TestConvertBrowserstack(t *testing.T) {
	var tests []runner

	tests = append(tests, prepare(t, "/browserstack/browserstackSnippet", &harness.Step{
		Id:   "Stage_null9dcf46",
		Name: "Browserstack",
		Type: "script",
		Spec: &harness.StepExec{
			Connector: "browserstack_connector",
			Image:     "harnesscommunity/browserstack",
			Shell:     "sh",
			Run: strings.Join([]string{
				"npm install selenium-webdriver",
				"node test.js",
			}, "\n"),
			Envs: map[string]string{
				"BROWSERSTACK_USERNAME":         "<+secrets.getValue(\"browserstack_username\")>",
				"BROWSERSTACK_ACCESS_KEY":       "<+secrets.getValue(\"browserstack_accesskey\")>",
				"BROWSERSTACK_BUILD_NAME":       "<+pipeline.sequenceId>-<+pipeline.executionId>",
				"BROWSERSTACK_BUILD_IDENTIFIER": "<+pipeline.executionId>",
			},
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertBrowserStack(tt.input)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertBrowserstack() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
