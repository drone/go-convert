package converthelpers

import (
	"testing"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/google/go-cmp/cmp"
)

func TestConvertNotifications(t *testing.T) {
	tests := []struct {
		name     string
		rules    []*v0.NotificationRule
		expected []*v1.Notification
	}{
		{
			name: "basic webhook notification",
			rules: []*v0.NotificationRule{
				{
					Identifier: "webhook_notification",
					Name:       "Webhook Alert",
					Enabled:    true,
					NotificationMethod: &v0.NotificationMethod{
						Type: "Webhook",
						Spec: v0.WebhookNotificationSpec{
							WebhookUrl: "https://example.com/webhook",
							Headers: map[string]string{
								"Authorization": "Bearer token",
							},
						},
					},
					PipelineEvents: []*v0.PipelineEvent{
						{Type: "PipelineSuccess"},
					},
				},
			},
			expected: []*v1.Notification{
				{
					ID:       "webhook_notification",
					Name:     "Webhook Alert",
					Disabled: false,
					Uses:     "webhook",
					With: map[string]interface{}{
						"url": "https://example.com/webhook",
						"headers": map[string]string{
							"Authorization": "Bearer token",
						},
					},
					On: []*v1.NotificationOn{
						{Pipeline: "success"},
					},
				},
			},
		},
		{
			name: "slack notification with user groups",
			rules: []*v0.NotificationRule{
				{
					Identifier: "slack_alert",
					Name:       "Slack Notification",
					Enabled:    true,
					NotificationMethod: &v0.NotificationMethod{
						Type: "Slack",
						Spec: v0.SlackNotificationSpec{
							WebhookUrl: "https://hooks.slack.com/services/xxx",
							UserGroups: []string{"@channel", "@team"},
						},
					},
					PipelineEvents: []*v0.PipelineEvent{
						{Type: "PipelineFailed"},
					},
				},
			},
			expected: []*v1.Notification{
				{
					ID:       "slack_alert",
					Name:     "Slack Notification",
					Disabled: false,
					Uses:     "slack",
					With: map[string]interface{}{
						"webhook":     "https://hooks.slack.com/services/xxx",
						"user-groups": []string{"@channel", "@team"},
					},
					On: []*v1.NotificationOn{
						{Pipeline: "failed"},
					},
				},
			},
		},
		{
			name: "email notification with empty recipients",
			rules: []*v0.NotificationRule{
				{
					Identifier: "email_alert",
					Name:       "Email Alert",
					Enabled:    true,
					NotificationMethod: &v0.NotificationMethod{
						Type: "Email",
						Spec: v0.EmailNotificationSpec{
							Recipients: []string{},
							UserGroups: []string{"admins"},
						},
					},
					PipelineEvents: []*v0.PipelineEvent{
						{Type: "PipelineEnd"},
					},
				},
			},
			expected: []*v1.Notification{
				{
					ID:       "email_alert",
					Name:     "Email Alert",
					Disabled: false,
					Uses:     "email",
					With: map[string]interface{}{
						"user-groups": []string{"admins"},
					},
					On: []*v1.NotificationOn{
						{Pipeline: "end"},
					},
				},
			},
		},
		{
			name: "pagerduty with empty user groups",
			rules: []*v0.NotificationRule{
				{
					Identifier: "pd_alert",
					Name:       "PagerDuty Alert",
					Enabled:    false,
					NotificationMethod: &v0.NotificationMethod{
						Type: "PagerDuty",
						Spec: v0.PagerDutyNotificationSpec{
							IntegrationKey: "pd-key-123",
							UserGroups:     []string{},
						},
					},
					PipelineEvents: []*v0.PipelineEvent{
						{Type: "PipelineFailed"},
					},
				},
			},
			expected: []*v1.Notification{
				{
					ID:       "pd_alert",
					Name:     "PagerDuty Alert",
					Disabled: true,
					Uses:     "pagerduty",
					With: map[string]interface{}{
						"key": "pd-key-123",
					},
					On: []*v1.NotificationOn{
						{Pipeline: "failed"},
					},
				},
			},
		},
		{
			name: "msteams with empty keys",
			rules: []*v0.NotificationRule{
				{
					Identifier: "teams_alert",
					Name:       "Teams Alert",
					Enabled:    true,
					NotificationMethod: &v0.NotificationMethod{
						Type: "MsTeams",
						Spec: v0.MsTeamsNotificationSpec{
							MsTeamKeys: []string{},
							UserGroups: []string{"devops"},
						},
					},
					PipelineEvents: []*v0.PipelineEvent{
						{Type: "PipelineStart"},
					},
				},
			},
			expected: []*v1.Notification{
				{
					ID:       "teams_alert",
					Name:     "Teams Alert",
					Disabled: false,
					Uses:     "ms-teams",
					With: map[string]interface{}{
						"user-groups": []string{"devops"},
					},
					On: []*v1.NotificationOn{
						{Pipeline: "start"},
					},
				},
			},
		},
		{
			name: "datadog with empty headers",
			rules: []*v0.NotificationRule{
				{
					Identifier: "dd_alert",
					Name:       "Datadog Alert",
					Enabled:    true,
					NotificationMethod: &v0.NotificationMethod{
						Type: "Datadog",
						Spec: v0.DatadogNotificationSpec{
							ApiKey:  "dd-api-key",
							Url:     "https://api.datadoghq.com",
							Headers: map[string]string{},
						},
					},
					PipelineEvents: []*v0.PipelineEvent{
						{Type: "PipelineSuccess"},
					},
				},
			},
			expected: []*v1.Notification{
				{
					ID:       "dd_alert",
					Name:     "Datadog Alert",
					Disabled: false,
					Uses:     "datadog",
					With: map[string]interface{}{
						"api-key": "dd-api-key",
						"url":     "https://api.datadoghq.com",
					},
					On: []*v1.NotificationOn{
						{Pipeline: "success"},
					},
				},
			},
		},
		{
			name: "AllEvents notification",
			rules: []*v0.NotificationRule{
				{
					Identifier: "all_events",
					Name:       "All Events Alert",
					Enabled:    true,
					NotificationMethod: &v0.NotificationMethod{
						Type: "Slack",
						Spec: v0.SlackNotificationSpec{
							WebhookUrl: "https://hooks.slack.com/services/all",
						},
					},
					PipelineEvents: []*v0.PipelineEvent{
						{Type: "AllEvents"},
						{Type: "PipelineSuccess"},
					},
				},
			},
			expected: []*v1.Notification{
				{
					ID:       "all_events",
					Name:     "All Events Alert",
					Disabled: false,
					Uses:     "slack",
					With: map[string]interface{}{
						"webhook": "https://hooks.slack.com/services/all",
					},
					On: []*v1.NotificationOn{
						{Pipeline: "all"},
						{Stage: "all"},
						{Step: "all"},
					},
				},
			},
		},
		{
			name: "stage events with specific stages",
			rules: []*v0.NotificationRule{
				{
					Identifier: "stage_alert",
					Name:       "Stage Alert",
					Enabled:    true,
					NotificationMethod: &v0.NotificationMethod{
						Type: "Email",
						Spec: v0.EmailNotificationSpec{
							Recipients: []string{"dev@example.com"},
						},
					},
					PipelineEvents: []*v0.PipelineEvent{
						{Type: "StageSuccess", ForStages: []string{"build", "test"}},
						{Type: "StageFailed", ForStages: []string{"deploy"}},
					},
				},
			},
			expected: []*v1.Notification{
				{
					ID:       "stage_alert",
					Name:     "Stage Alert",
					Disabled: false,
					Uses:     "email",
					With: map[string]interface{}{
						"recipients": []string{"dev@example.com"},
					},
					On: []*v1.NotificationOn{
						{
							Stage: map[string]interface{}{
								"success": []string{"build", "test"},
								"failed":  "deploy",
							},
						},
					},
				},
			},
		},
		{
			name: "stage events with AllStages",
			rules: []*v0.NotificationRule{
				{
					Identifier: "all_stages",
					Name:       "All Stages Alert",
					Enabled:    true,
					NotificationMethod: &v0.NotificationMethod{
						Type: "Slack",
						Spec: v0.SlackNotificationSpec{
							WebhookUrl: "https://hooks.slack.com/services/stages",
						},
					},
					PipelineEvents: []*v0.PipelineEvent{
						{Type: "StageFailed", ForStages: []string{"AllStages"}},
					},
				},
			},
			expected: []*v1.Notification{
				{
					ID:       "all_stages",
					Name:     "All Stages Alert",
					Disabled: false,
					Uses:     "slack",
					With: map[string]interface{}{
						"webhook": "https://hooks.slack.com/services/stages",
					},
					On: []*v1.NotificationOn{
						{
							Stage: map[string]interface{}{
								"failed": "all",
							},
						},
					},
				},
			},
		},
		{
			name: "stage events with empty ForStages",
			rules: []*v0.NotificationRule{
				{
					Identifier: "empty_stages",
					Name:       "Empty Stages",
					Enabled:    true,
					NotificationMethod: &v0.NotificationMethod{
						Type: "Email",
						Spec: v0.EmailNotificationSpec{
							Recipients: []string{"test@example.com"},
						},
					},
					PipelineEvents: []*v0.PipelineEvent{
						{Type: "StageSuccess", ForStages: []string{}},
						{Type: "PipelineFailed"},
					},
				},
			},
			expected: []*v1.Notification{
				{
					ID:       "empty_stages",
					Name:     "Empty Stages",
					Disabled: false,
					Uses:     "email",
					With: map[string]interface{}{
						"recipients": []string{"test@example.com"},
					},
					On: []*v1.NotificationOn{
						{Pipeline: "failed"},
					},
				},
			},
		},
		{
			name: "notification with template",
			rules: []*v0.NotificationRule{
				{
					Identifier: "templated_notification",
					Name:       "Templated Alert",
					Enabled:    true,
					NotificationMethod: &v0.NotificationMethod{
						Type: "Slack",
						Spec: v0.SlackNotificationSpec{
							WebhookUrl: "https://hooks.slack.com/services/template",
						},
					},
					PipelineEvents: []*v0.PipelineEvent{
						{Type: "PipelineSuccess"},
					},
					NotificationTemplate: &v0.NotificationTemplate{
						TemplateRef:  "custom_template",
						VersionLabel: "v1.0",
					},
				},
			},
			expected: []*v1.Notification{
				{
					ID:       "templated_notification",
					Name:     "Templated Alert",
					Disabled: false,
					Uses:     "slack",
					With: map[string]interface{}{
						"webhook": "https://hooks.slack.com/services/template",
					},
					On: []*v1.NotificationOn{
						{Pipeline: "success"},
					},
					NotificationTemplate: &v1.NotificationTemplate{
						Uses: "custom_template",
						With: map[string]interface{}{
							"version": "v1.0",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertNotifications(tt.rules)

			if diff := cmp.Diff(tt.expected, result); diff != "" {
				t.Errorf("Notification mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertNotifications_NilAndEmpty(t *testing.T) {
	tests := []struct {
		name     string
		rules    []*v0.NotificationRule
		expected []*v1.Notification
	}{
		{
			name:     "nil rules",
			rules:    nil,
			expected: nil,
		},
		{
			name:     "empty rules",
			rules:    []*v0.NotificationRule{},
			expected: nil,
		},
		{
			name: "rules with nil entries",
			rules: []*v0.NotificationRule{
				nil,
				{
					Identifier: "valid",
					Name:       "Valid",
					Enabled:    true,
					NotificationMethod: &v0.NotificationMethod{
						Type: "Email",
						Spec: v0.EmailNotificationSpec{
							Recipients: []string{"test@example.com"},
						},
					},
					PipelineEvents: []*v0.PipelineEvent{
						{Type: "PipelineSuccess"},
					},
				},
				nil,
			},
			expected: []*v1.Notification{
				{
					ID:       "valid",
					Name:     "Valid",
					Disabled: false,
					Uses:     "email",
					With: map[string]interface{}{
						"recipients": []string{"test@example.com"},
					},
					On: []*v1.NotificationOn{
						{Pipeline: "success"},
					},
				},
			},
		},
		{
			name: "notification with nil method",
			rules: []*v0.NotificationRule{
				{
					Identifier:         "no_method",
					Name:               "No Method",
					Enabled:            true,
					NotificationMethod: nil,
					PipelineEvents: []*v0.PipelineEvent{
						{Type: "PipelineSuccess"},
					},
				},
			},
			expected: []*v1.Notification{
				{
					ID:       "no_method",
					Name:     "No Method",
					Disabled: false,
					Uses:     "",
					With:     nil,
					On: []*v1.NotificationOn{
						{Pipeline: "success"},
					},
				},
			},
		},
		{
			name: "notification with empty events",
			rules: []*v0.NotificationRule{
				{
					Identifier: "no_events",
					Name:       "No Events",
					Enabled:    true,
					NotificationMethod: &v0.NotificationMethod{
						Type: "Slack",
						Spec: v0.SlackNotificationSpec{
							WebhookUrl: "https://hooks.slack.com/services/test",
						},
					},
					PipelineEvents: []*v0.PipelineEvent{},
				},
			},
			expected: []*v1.Notification{
				{
					ID:       "no_events",
					Name:     "No Events",
					Disabled: false,
					Uses:     "slack",
					With: map[string]interface{}{
						"webhook": "https://hooks.slack.com/services/test",
					},
					On: nil,
				},
			},
		},
		{
			name: "notification with nil events in array",
			rules: []*v0.NotificationRule{
				{
					Identifier: "nil_events",
					Name:       "Nil Events",
					Enabled:    true,
					NotificationMethod: &v0.NotificationMethod{
						Type: "Email",
						Spec: v0.EmailNotificationSpec{
							Recipients: []string{"test@example.com"},
						},
					},
					PipelineEvents: []*v0.PipelineEvent{
						nil,
						{Type: "PipelineSuccess"},
						nil,
					},
				},
			},
			expected: []*v1.Notification{
				{
					ID:       "nil_events",
					Name:     "Nil Events",
					Disabled: false,
					Uses:     "email",
					With: map[string]interface{}{
						"recipients": []string{"test@example.com"},
					},
					On: []*v1.NotificationOn{
						{Pipeline: "success"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertNotifications(tt.rules)

			if diff := cmp.Diff(tt.expected, result); diff != "" {
				t.Errorf("Notification mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertNotificationMethodType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Webhook", "Webhook", "webhook"},
		{"Slack", "Slack", "slack"},
		{"PagerDuty", "PagerDuty", "pagerduty"},
		{"MsTeams", "MsTeams", "ms-teams"},
		{"Email", "Email", "email"},
		{"Datadog", "Datadog", "datadog"},
		{"Unknown", "CustomType", "customtype"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertNotificationMethodType(tt.input)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
