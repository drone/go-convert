package converthelpers

import (
	"testing"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/google/go-cmp/cmp"
)

func TestConvertOutputVariables(t *testing.T) {
	tests := []struct {
		name     string
		input    []*v0.Output
		expected []*v1.Output
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: []*v1.Output{},
		},
		{
			name:     "empty input",
			input:    []*v0.Output{},
			expected: []*v1.Output{},
		},
		{
			name: "output with different shell variable name",
			input: []*v0.Output{
				{Name: "approvalStatus", Value: "STATUS"},
			},
			expected: []*v1.Output{
				{Name: "STATUS", Alias: "approvalStatus", Mask: false},
			},
		},
		{
			name: "output with same name and value (alias set)",
			input: []*v0.Output{
				{Name: "STATUS", Value: "STATUS"},
			},
			expected: []*v1.Output{
				{Name: "STATUS", Alias: "STATUS", Mask: false},
			},
		},
		{
			name: "output with empty value (uses name for both)",
			input: []*v0.Output{
				{Name: "MY_VAR", Value: ""},
			},
			expected: []*v1.Output{
				{Name: "MY_VAR", Alias: "MY_VAR", Mask: false},
			},
		},
		{
			name: "multiple outputs with different mappings",
			input: []*v0.Output{
				{Name: "buildNumber", Value: "BUILD_NUM"},
				{Name: "commitHash", Value: "GIT_COMMIT"},
				{Name: "version", Value: "APP_VERSION"},
			},
			expected: []*v1.Output{
				{Name: "BUILD_NUM", Alias: "buildNumber", Mask: false},
				{Name: "GIT_COMMIT", Alias: "commitHash", Mask: false},
				{Name: "APP_VERSION", Alias: "version", Mask: false},
			},
		},
		{
			name: "secret output variable",
			input: []*v0.Output{
				{Name: "apiToken", Value: "API_TOKEN", Type: "Secret"},
			},
			expected: []*v1.Output{
				{Name: "API_TOKEN", Alias: "apiToken", Mask: true},
			},
		},
		{
			name: "output with non-Secret type",
			input: []*v0.Output{
				{Name: "result", Value: "RESULT", Type: "String"},
			},
			expected: []*v1.Output{
				{Name: "RESULT", Alias: "result", Mask: false},
			},
		},
		{
			name: "nil entries in slice are skipped",
			input: []*v0.Output{
				{Name: "first", Value: "FIRST"},
				nil,
				{Name: "third", Value: "THIRD"},
			},
			expected: []*v1.Output{
				{Name: "FIRST", Alias: "first", Mask: false},
				{Name: "THIRD", Alias: "third", Mask: false},
			},
		},
		{
			name: "mixed secret and non-secret outputs",
			input: []*v0.Output{
				{Name: "publicKey", Value: "PUBLIC_KEY", Type: "String"},
				{Name: "privateKey", Value: "PRIVATE_KEY", Type: "Secret"},
				{Name: "status", Value: ""},
			},
			expected: []*v1.Output{
				{Name: "PUBLIC_KEY", Alias: "publicKey", Mask: false},
				{Name: "PRIVATE_KEY", Alias: "privateKey", Mask: true},
				{Name: "status", Alias: "status", Mask: false},
			},
		},
		{
			name: "realistic example from documentation",
			input: []*v0.Output{
				{Name: "approvalStatus", Value: "STATUS", Type: "String"},
				{Name: "deploymentId", Value: "DEPLOY_ID", Type: "String"},
			},
			expected: []*v1.Output{
				{Name: "STATUS", Alias: "approvalStatus", Mask: false},
				{Name: "DEPLOY_ID", Alias: "deploymentId", Mask: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertOutputVariables(tt.input)
			if diff := cmp.Diff(tt.expected, result); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
