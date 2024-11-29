package json

import (
	"testing"

	harness "github.com/drone/spec/dist/go"
	"github.com/google/go-cmp/cmp"
)

func TestConvertPagerDuty(t *testing.T) {

	var tests []runner
	tests = append(tests, prepare(t, "/pagerduty/pagerduty_snippet", &harness.Step{
		Id:   "pagerduty51d22b",
		Name: "Pagerduty",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "plugins/pagerduty",
			With: map[string]interface{}{
				"log_level":         string("info"),
				"routing_key":       string("a666ad1326f34605d06c0dbd4d87c1cb"),
				"incident_summary":  string("Build Failed for test-pager-duty"),
				"incident_source":   string("test-pager-duty"),
				"incident_severity": string("critical"),
				"resolve":           bool(false),
				"job_status":        string("<+pipeline.status>"),
				"dedup_key":         string("E54EC853A59A3815EF3632D5F854CF26"),
			},
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertPagerDuty(tt.input, tt.input.ParameterMap)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertPagerDuty() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertPagerDutyChangeEvent(t *testing.T) {

	var tests []runner
	tests = append(tests, prepare(t, "/pagerdutyChangeEvent/pagerdutyChangeEvent_snippet", &harness.Step{
		Id:   "pagerdutyChangeEvent499434",
		Name: "Pagerduty_Change_Event",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "plugins/pagerduty",
			With: map[string]interface{}{
				"log_level":           string("info"),
				"routing_key":         string("a666ad1326f34605d06c0dbd4d87c1cb"),
				"incident_summary":    string("Job test-pager-duty completed with status SUCCESS"),
				"incident_source":     string(""),
				"custom_details":      string(`{"buildNumber":"22","jobName":"test-pager-duty","jobURL":"http://localhost:8080/job/test-pager-duty/22/"}`),
				"create_change_event": bool(true),
			},
		},
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertPagerDutyChangeEvent(tt.input, tt.input.ParameterMap)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConvertPagerDutyChangeEvent() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
