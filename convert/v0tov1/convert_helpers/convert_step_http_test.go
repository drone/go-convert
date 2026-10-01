package converthelpers

import (
	"testing"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	"github.com/google/go-cmp/cmp"
)

func TestConvertStepHTTP(t *testing.T) {
	tests := []struct {
		name     string
		step     *v0.Step
		expected map[string]interface{}
	}{
		{
			name: "basic url and method",
			step: &v0.Step{
				Spec: &v0.StepHTTP{
					URL:    "https://example.com",
					Method: "GET",
				},
			},
			expected: map[string]interface{}{
				"url":    "https://example.com",
				"method": "GET",
			},
		},
		{
			name: "headers, body and assertion",
			step: &v0.Step{
				Spec: &v0.StepHTTP{
					URL:    "https://example.com",
					Method: "POST",
					Headers: []*v0.KeyValuePair{
						{Key: "Content-Type", Value: "application/json"},
					},
					RequestBody: `{"name":"John"}`,
					Assertion:   "response.code == 200",
				},
			},
			expected: map[string]interface{}{
				"url":       "https://example.com",
				"method":    "POST",
				"headers":   `{"Content-Type":"application/json"}`,
				"body":      `{"name":"John"}`,
				"assertion": "response.code == 200",
			},
		},
		{
			name: "output variables and input variables",
			step: &v0.Step{
				Spec: &v0.StepHTTP{
					URL:    "https://example.com",
					Method: "GET",
					OutputVariables: []*v0.Variable{
						{Name: "user_id", Value: "response.body.id"},
					},
					InputVariables: []*v0.Variable{
						{Name: "TOKEN", Value: "abc123"},
					},
				},
			},
			expected: map[string]interface{}{
				"url":         "https://example.com",
				"method":      "GET",
				"output_vars": `{"user_id":"response.body.id"}`,
				"env_vars":    `{"TOKEN":"abc123"}`,
			},
		},
		{
			name: "client certificate fields",
			step: &v0.Step{
				Spec: &v0.StepHTTP{
					URL:            "https://example.com",
					Method:         "GET",
					Certificate:    "/path/to/client.crt",
					CertificateKey: "/path/to/client.key",
				},
			},
			expected: map[string]interface{}{
				"url":         "https://example.com",
				"method":      "GET",
				"client_cert": "/path/to/client.crt",
				"client_key":  "/path/to/client.key",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertStepHTTP(tt.step)
			if result == nil {
				t.Fatal("expected non-nil result")
			}

			if result.Uses != "httpStep" {
				t.Errorf("expected Uses to be httpStep, got %s", result.Uses)
			}

			if diff := cmp.Diff(tt.expected, result.With); diff != "" {
				t.Errorf("With mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertStepHTTP_NilCases(t *testing.T) {
	tests := []struct {
		name string
		step *v0.Step
	}{
		{
			name: "nil step",
			step: nil,
		},
		{
			name: "nil spec",
			step: &v0.Step{Spec: nil},
		},
		{
			name: "wrong spec type",
			step: &v0.Step{Spec: &v0.StepRun{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := ConvertStepHTTP(tt.step); result != nil {
				t.Errorf("expected nil result, got %+v", result)
			}
		})
	}
}
