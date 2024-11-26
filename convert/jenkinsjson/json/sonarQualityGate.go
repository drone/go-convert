package json

import (
	harness "github.com/drone/spec/dist/go"
)

// ConvertSonarQualityGate creates a Harness step for nunit plugin.
func ConvertSonarQualityGate(node Node) *harness.Step {

	convertSonarQualityGate := &harness.Step{
		Name: "Sonarqube_Quality_Gate",
		Type: "plugin",
		Id:   SanitizeForId(node.SpanName, node.SpanId),
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
	}

	return convertSonarQualityGate
}
