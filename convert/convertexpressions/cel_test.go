package convertexpressions

import "testing"

// TestConvertExpressionWithTrie_CELSimplePath verifies that a converted
// expression is emitted using CEL-style ${{...}} delimiters when the
// converted content is a simple dotted path (no function calls, ternaries,
// nested expressions, or other operators), and stays <+...> otherwise.
func TestConvertExpressionWithTrie_CELSimplePath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Plain dotted path becomes CEL",
			input:    "<+pipeline.variables.var1>",
			expected: "${{pipeline.variables.var1}}",
		},
		{
			name:     "Path with array index becomes CEL",
			input:    "<+pipeline.stages.build.spec.execution.steps.step1.spec.outputVariables[0].value>",
			expected: "${{pipeline.stages.build.steps.step1.spec.output[0].name}}",
		},
		{
			name:     "Function call stays angle-delimited",
			input:    "<+pipeline.variables.var1.toUpperCase()>",
			expected: "<+pipeline.variables.var1.toUpperCase()>",
		},
		{
			name:     "Ternary expression stays angle-delimited",
			input:    `<+<+pipeline.stages.build.identifier>=="build"?<+pipeline.stages.deploy.identifier>:"default">`,
			expected: `<+<+pipeline.stages.build.id>=="build"?<+pipeline.stages.deploy.id>:"default">`,
		},
		{
			name:     "Nested expression comparison stays angle-delimited",
			input:    `<+<+pipeline.stages.build.identifier> == "build">`,
			expected: `<+<+pipeline.stages.build.id> == "build">`,
		},
		{
			name:     "Concatenation stays angle-delimited",
			input:    `<+<+pipeline.stages.build.identifier> + "-" + <+pipeline.stages.deploy.identifier>>`,
			expected: `<+<+pipeline.stages.build.id> + "-" + <+pipeline.stages.deploy.id>>`,
		},
		{
			name:     "Single-segment expression stays angle-delimited",
			input:    "<+input>",
			expected: "<+input>",
		},
		{
			name:     "Single-segment serviceVariables stays angle-delimited",
			input:    "<+serviceVariables>",
			expected: "<+serviceVariables>",
		},
		{
			name:     "Expression inside larger script text becomes CEL",
			input:    `echo "Repo: <+pipeline.properties.ci.codebase.repoName>"`,
			expected: `echo "Repo: ${{codebase.repoName}}"`,
		},
		{
			name:     "Already CEL-delimited path stays CEL",
			input:    "${{pipeline.variables.var1}}",
			expected: "${{pipeline.variables.var1}}",
		},
		{
			name:     "Whitespace-padded path becomes CEL",
			input:    "<+ pipeline.variables.var1 >",
			expected: "${{ pipeline.variables.var1 }}",
		},
		{
			name:     "Unmatched passthrough path still becomes CEL",
			input:    "<+just.a.plain.path>",
			expected: "${{just.a.plain.path}}",
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

func TestIsSimplePathExpression(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"plain dotted path", "a.b.c", true},
		{"path with array index", "a.b[0].c", true},
		{"path with named index", "a.b[name].c", true},
		{"single segment", "a", false},
		{"empty string", "", false},
		{"function call", "a.b.toUpperCase()", false},
		{"ternary", `a=="b"?c:"default"`, false},
		{"nested expression", "<+a.b>", false},
		{"comparison", `a == "b"`, false},
		{"concatenation", `a + "-" + b`, false},
		{"whitespace padded", " a.b.c ", true},
		{"trailing dot invalid", "a.b.", false},
		{"unclosed bracket invalid", "a.b[0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSimplePathExpression(tt.input); got != tt.expected {
				t.Errorf("isSimplePathExpression(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}
