package convertexpressions

import (
	"testing"
)

func TestTrieConvert_ManifestConfig(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Primary manifest id",
			input:    "<+manifestConfig.primaryManifestId>",
			expected: "${{manifests.primary.id}}",
		},
		{
			name:     "Primary manifest id with method call",
			input:    "<+manifestConfig.primaryManifestId.toUpperCase()>",
			expected: "<+manifests.primary.id.toUpperCase()>",
		},
		{
			name:     "Primary manifest id in mixed text",
			input:    `manifestId: <+manifestConfig.primaryManifestId>`,
			expected: "manifestId: ${{manifests.primary.id}}",
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
