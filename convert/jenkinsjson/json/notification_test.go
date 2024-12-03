package json

import (
	"encoding/json"
	"testing"

	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
)

func TestConvertNotification(t *testing.T) {

	// Remove line breaks and format as compact JSON
	var compactData map[string]interface{}
	_ = json.Unmarshal([]byte("{\"status\": \"Build Successful\", \"job\": \"${env.JOB_NAME}\", \"buildNumber\": \"${env.BUILD_NUMBER}\"}"), &compactData) // Parse JSON
	compactBytes, _ := json.Marshal(compactData)                                                                                                           // Convert to compact JSON
	compactString := string(compactBytes)                                                                                                                  // Convert bytes to string

	var tests []runner
	tests = append(tests, prepare(t, "/notification/notification_snippet", &harness.Step{
		Id:   "notifyEndpointsca9fc8",
		Name: "Notification",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Connector: "<+input>",
			Image:     "plugins/webhook",
			With: map[string]interface{}{
				"urls":         "<+input>",
				"username":     "<+input>",
				"password":     "<+input>",
				"method":       "<+input>",
				"content_type": "application/json",
				"debug":        "true",
				"template":     compactString,
			},
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertNotification(tt.input, tt.input.ParameterMap)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertNotification() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
