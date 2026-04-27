package convertexpressions

// rules for clone codebase expressions: <+codebase.>
var CodebaseConversionRules = []ConversionRule{
	{"(build.spec.branch)", "branch"},
	{"(build.spec.tag)", "tag"},
	{"(build.spec.number)", "prNumber"},
}
