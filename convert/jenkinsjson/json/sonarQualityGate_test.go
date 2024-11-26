package json

import (
	"testing"

	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
)

func TestConvertSonarQualityGate(t *testing.T) {

	var tests []runner
	tests = append(tests, prepare(t, "/sonarQualityGate/sonarQualityGate_snippet", &harness.Step{
		Id:   "waitForQualityGate6278c9",
		Name: "Sonarqube_Quality_Gate",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "plugins/sonarqube-scanner:v2.4.2",
			With: map[string]interface{}{
				"sonar_key":             "<+input>",
				"sonar_name":            "<+input>",
				"sonar_host":            "<+input>",
				"sonar_token":           "<+input>",
				"timeout":               "300",
				"sources":               ".",
				"skip_scan":             "false",
				"sonar_qualitygate":     "OK",
				"sonar_quality_enabled": "true",
				"sonar_organization":    "<+input>",
				"scm_disabled":          "false",
			},
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertSonarQualityGate(tt.input)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertSonarQualityGate() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
