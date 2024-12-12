package json

import (
	"testing"

	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
)

func TestConvertNotification(t *testing.T) { // Convert bytes to string

	var tests []runner
	tests = append(tests, prepare(t, "/notification/notification_snippet", &harness.Step{
		Id:   "notifyEndpoints2f61b6",
		Name: "Notification",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "plugins/webhook",
			With: map[string]interface{}{
				"urls":         []string{"https://httpbin.org/post"},
				"method":       "POST",
				"username":     "<+input>",
				"password":     "<+input>",
				"token-value":  "<+input>",
				"token-type":   "<+input>",
				"content-type": "application/json",
				"template": `{
  		"phase": "COMPLETED",
  		"notes": "Build metrics for analysis",
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
