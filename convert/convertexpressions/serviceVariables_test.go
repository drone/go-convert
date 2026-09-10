package convertexpressions

import (
	"testing"
)

func TestTrieConvert_ServiceVariables(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Service variable single field",
			input:    "<+serviceConfig.serviceDefinition.spec.variables.myVar>",
			expected: "<+serviceVariables.myVar>",
		},
		{
			name:     "Service variable with method call",
			input:    "<+serviceConfig.serviceDefinition.spec.variables.myVar.toUpperCase()>",
			expected: "<+serviceVariables.myVar.toUpperCase()>",
		},
		{
			name:     "Service variable exact path, no trailing field",
			input:    "<+serviceConfig.serviceDefinition.spec.variables>",
			expected: "<+serviceVariables>",
		},
		{
			name:     "Service variable in mixed text",
			input:    `image: <+serviceConfig.serviceDefinition.spec.variables.image>:latest`,
			expected: `image: <+serviceVariables.image>:latest`,
		},
		{
			name:     "Already v1 serviceVariables passes through unchanged",
			input:    "<+serviceVariables.myVar>",
			expected: "<+serviceVariables.myVar>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertExpressionWithTrie(tt.input, nil, false)
			if got != tt.expected {
				t.Errorf("ConvertExpressionWithTrie() failed\ninput:    %s\ngot:      %s\nexpected: %s", tt.input, got, tt.expected)
			}
		})
	}
}
