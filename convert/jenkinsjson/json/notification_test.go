package json

import (
	"testing"

	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
)

func TestConvertNotification(t *testing.T) { // Convert bytes to string

	var tests []runner
	tests = append(tests, prepare(t, "/notification/notification_snippet", &harness.Step{
		Id:   "notifyEndpointsca9fc8",
		Name: "Notification",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Connector: "<+input>",
			Image:     "plugins/webhook",
			With: map[string]interface{}{
				"urls":         []string{"https://httpbin.org/post"},
				"method":       "POST",
				"username":     "<+input>",
				"password":     "<+input>",
				"token-value":  "<+input>",
				"token-type":   "<+input>",
				"content-type": "application/json",
				"headers":      "Authorization=abcd",
				"template": `{
  		"loglines": "10",
  		"phase": "COMPLETED",
  		"status": "Success",
  		"notes": "Build metrics for analysis",
  		"timestamp": "${time.now()}"
	}`,
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
