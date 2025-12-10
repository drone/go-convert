package converthelpers

import (
	"testing"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/google/go-cmp/cmp"
)

func TestConvertStepServiceNowCreate(t *testing.T) {
	tests := []struct {
		name     string
		step     *v0.Step
		expected map[string]interface{}
	}{
		{
			name: "basic ServiceNow create",
			step: &v0.Step{
				Spec: &v0.StepServiceNowCreate{
					ConnectorRef: "servicenow-connector",
					TicketType:   "incident",
					Fields: []*v0.ServiceNowField{
						{Name: "short_description", Value: "Test incident"},
						{Name: "priority", Value: "1"},
					},
				},
			},
			expected: map[string]interface{}{
				"connector":             "servicenow-connector",
				"ticket_type":           "incident",
				"create_ticket_options": "Fields",
				"fields": map[string]string{
					"short_description": "Test incident",
					"priority":          "1",
				},
			},
		},
		{
			name: "ServiceNow create with empty fields array",
			step: &v0.Step{
				Spec: &v0.StepServiceNowCreate{
					ConnectorRef: "servicenow-connector",
					TicketType:   "change_request",
					Fields:       []*v0.ServiceNowField{},
				},
			},
			expected: map[string]interface{}{
				"connector":             "servicenow-connector",
				"ticket_type":           "change_request",
				"create_ticket_options": "Fields",
			},
		},
		{
			name: "ServiceNow create with nil fields",
			step: &v0.Step{
				Spec: &v0.StepServiceNowCreate{
					ConnectorRef: "servicenow-connector",
					TicketType:   "problem",
					Fields:       nil,
				},
			},
			expected: map[string]interface{}{
				"connector":             "servicenow-connector",
				"ticket_type":           "problem",
				"create_ticket_options": "Fields",
			},
		},
		{
			name: "ServiceNow create with nil field entries",
			step: &v0.Step{
				Spec: &v0.StepServiceNowCreate{
					ConnectorRef: "servicenow-connector",
					TicketType:   "incident",
					Fields: []*v0.ServiceNowField{
						{Name: "description", Value: "Valid field"},
						nil,
						{Name: "urgency", Value: "2"},
						nil,
					},
				},
			},
			expected: map[string]interface{}{
				"connector":             "servicenow-connector",
				"ticket_type":           "incident",
				"create_ticket_options": "Fields",
				"fields": map[string]string{
					"description": "Valid field",
					"urgency":     "2",
				},
			},
		},
		{
			name: "ServiceNow create with template type",
			step: &v0.Step{
				Spec: &v0.StepServiceNowCreate{
					ConnectorRef: "servicenow-connector",
					TicketType:   "incident",
					CreateType:   "Template",
					Fields: []*v0.ServiceNowField{
						{Name: "template_name", Value: "standard_incident"},
					},
				},
			},
			expected: map[string]interface{}{
				"connector":             "servicenow-connector",
				"ticket_type":           "incident",
				"create_ticket_options": "Template",
				"fields": map[string]string{
					"template_name": "standard_incident",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertStepServiceNowCreate(tt.step)
			if result == nil {
				t.Fatal("expected non-nil result")
			}

			if result.Uses != "serviceNowCreate" {
				t.Errorf("expected Uses to be serviceNowCreate, got %s", result.Uses)
			}

			if diff := cmp.Diff(tt.expected, result.With); diff != "" {
				t.Errorf("With mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertStepServiceNowUpdate(t *testing.T) {
	tests := []struct {
		name     string
		step     *v0.Step
		expected map[string]interface{}
	}{
		{
			name: "basic ServiceNow update",
			step: &v0.Step{
				Spec: &v0.StepServiceNowUpdate{
					ConnectorRef: "servicenow-connector",
					TicketType:   "incident",
					TicketNumber: "INC0010001",
					Fields: []*v0.ServiceNowField{
						{Name: "state", Value: "2"},
						{Name: "comments", Value: "Updated by automation"},
					},
				},
			},
			expected: map[string]interface{}{
				"connector":            "servicenow-connector",
				"ticket_type":          "incident",
				"ticket_number":        "INC0010001",
				"update_ticket_option": "Fields",
				"fields": map[string]string{
					"state":    "2",
					"comments": "Updated by automation",
				},
			},
		},
		{
			name: "ServiceNow update with empty fields array",
			step: &v0.Step{
				Spec: &v0.StepServiceNowUpdate{
					ConnectorRef: "servicenow-connector",
					TicketType:   "change_request",
					TicketNumber: "CHG0030001",
					Fields:       []*v0.ServiceNowField{},
				},
			},
			expected: map[string]interface{}{
				"connector":            "servicenow-connector",
				"ticket_type":          "change_request",
				"ticket_number":        "CHG0030001",
				"update_ticket_option": "Fields",
			},
		},
		{
			name: "ServiceNow update with nil fields",
			step: &v0.Step{
				Spec: &v0.StepServiceNowUpdate{
					ConnectorRef: "servicenow-connector",
					TicketType:   "problem",
					TicketNumber: "PRB0040001",
					Fields:       nil,
				},
			},
			expected: map[string]interface{}{
				"connector":            "servicenow-connector",
				"ticket_type":          "problem",
				"ticket_number":        "PRB0040001",
				"update_ticket_option": "Fields",
			},
		},
		{
			name: "ServiceNow update with nil field entries",
			step: &v0.Step{
				Spec: &v0.StepServiceNowUpdate{
					ConnectorRef: "servicenow-connector",
					TicketType:   "incident",
					TicketNumber: "INC0010002",
					Fields: []*v0.ServiceNowField{
						nil,
						{Name: "work_notes", Value: "Investigation complete"},
						nil,
						{Name: "assigned_to", Value: "john.doe"},
					},
				},
			},
			expected: map[string]interface{}{
				"connector":            "servicenow-connector",
				"ticket_type":          "incident",
				"ticket_number":        "INC0010002",
				"update_ticket_option": "Fields",
				"fields": map[string]string{
					"work_notes":  "Investigation complete",
					"assigned_to": "john.doe",
				},
			},
		},
		{
			name: "ServiceNow update with template",
			step: &v0.Step{
				Spec: &v0.StepServiceNowUpdate{
					ConnectorRef:          "servicenow-connector",
					TicketType:            "incident",
					TicketNumber:          "INC0010003",
					UseServiceNowTemplate: true,
					Fields: []*v0.ServiceNowField{
						{Name: "template_name", Value: "resolution_template"},
					},
				},
			},
			expected: map[string]interface{}{
				"connector":            "servicenow-connector",
				"ticket_type":          "incident",
				"ticket_number":        "INC0010003",
				"update_ticket_option": "Template",
				"fields": map[string]string{
					"template_name": "resolution_template",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertStepServiceNowUpdate(tt.step)
			if result == nil {
				t.Fatal("expected non-nil result")
			}

			if result.Uses != "serviceNowUpdate" {
				t.Errorf("expected Uses to be serviceNowUpdate, got %s", result.Uses)
			}

			if diff := cmp.Diff(tt.expected, result.With); diff != "" {
				t.Errorf("With mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertStepServiceNow_NilCases(t *testing.T) {
	tests := []struct {
		name      string
		step      *v0.Step
		converter func(*v0.Step) *v1.StepTemplate
	}{
		{
			name:      "create with nil step",
			step:      nil,
			converter: ConvertStepServiceNowCreate,
		},
		{
			name: "create with nil spec",
			step: &v0.Step{
				Spec: nil,
			},
			converter: ConvertStepServiceNowCreate,
		},
		{
			name: "create with wrong spec type",
			step: &v0.Step{
				Spec: &v0.StepRun{},
			},
			converter: ConvertStepServiceNowCreate,
		},
		{
			name:      "update with nil step",
			step:      nil,
			converter: ConvertStepServiceNowUpdate,
		},
		{
			name: "update with nil spec",
			step: &v0.Step{
				Spec: nil,
			},
			converter: ConvertStepServiceNowUpdate,
		},
		{
			name: "update with wrong spec type",
			step: &v0.Step{
				Spec: &v0.StepRun{},
			},
			converter: ConvertStepServiceNowUpdate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.converter(tt.step)
			if result != nil {
				t.Errorf("expected nil result, got %v", result)
			}
		})
	}
}
