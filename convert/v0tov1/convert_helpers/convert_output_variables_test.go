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
			name: "single output with value",
			input: []*v0.Output{
				{Name: "MY_VAR", Value: "output_value"},
			},
			expected: []*v1.Output{
				{Alias: "MY_VAR", Name: "output_value", Mask: false},
			},
		},
		{
			name: "single output name fallback when value is empty",
			input: []*v0.Output{
				{Name: "MY_VAR", Value: ""},
			},
			expected: []*v1.Output{
				{Alias: "MY_VAR", Name: "MY_VAR", Mask: false},
			},
		},
		{
			name: "multiple outputs",
			input: []*v0.Output{
				{Name: "VAR_A", Value: "val_a"},
				{Name: "VAR_B", Value: "val_b"},
				{Name: "VAR_C", Value: "val_c"},
			},
			expected: []*v1.Output{
				{Alias: "VAR_A", Name: "val_a", Mask: false},
				{Alias: "VAR_B", Name: "val_b", Mask: false},
				{Alias: "VAR_C", Name: "val_c", Mask: false},
			},
		},
		{
			name: "output with Type Secret sets mask true",
			input: []*v0.Output{
				{Name: "SECRET_VAR", Value: "secret_val", Type: "Secret"},
			},
			expected: []*v1.Output{
				{Alias: "SECRET_VAR", Name: "secret_val", Mask: true},
			},
		},
		{
			name: "output with non-Secret type sets mask false",
			input: []*v0.Output{
				{Name: "NORMAL_VAR", Value: "normal_val", Type: "String"},
			},
			expected: []*v1.Output{
				{Alias: "NORMAL_VAR", Name: "normal_val", Mask: false},
			},
		},
		{
			name: "nil entries in slice are skipped",
			input: []*v0.Output{
				{Name: "VAR_A", Value: "val_a"},
				nil,
				{Name: "VAR_C", Value: "val_c"},
			},
			expected: []*v1.Output{
				{Alias: "VAR_A", Name: "val_a", Mask: false},
				{Alias: "VAR_C", Name: "val_c", Mask: false},
			},
		},
		{
			name: "mixed secret and non-secret outputs",
			input: []*v0.Output{
				{Name: "PUBLIC", Value: "pub_val", Type: "String"},
				{Name: "PRIVATE", Value: "priv_val", Type: "Secret"},
				{Name: "FALLBACK", Value: ""},
			},
			expected: []*v1.Output{
				{Alias: "PUBLIC", Name: "pub_val", Mask: false},
				{Alias: "PRIVATE", Name: "priv_val", Mask: true},
				{Alias: "FALLBACK", Name: "FALLBACK", Mask: false},
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
