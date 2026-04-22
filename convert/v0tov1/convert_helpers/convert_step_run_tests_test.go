package converthelpers

import (
	"testing"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/drone/go-convert/internal/flexible"
	"github.com/google/go-cmp/cmp"
)

func TestConvertStepRunTests(t *testing.T) {
	tests := []struct {
		name     string
		input    *v0.Step
		expected *v1.StepTest
	}{
		{
			name: "Java Maven with full configuration",
			input: &v0.Step{
				ID:   "run_tests",
				Name: "Run Tests",
				Type: v0.StepTypeRunTests,
				Spec: &v0.StepRunTests{
					Language:             "Java",
					BuildTool:            "Maven",
					Args:                 "-Dmaven.test.failure.ignore=true",
					PreCommand:           "echo \"Starting tests\"",
					PostCommand:          "echo \"Tests complete\"",
					RunOnlySelectedTests: &flexible.Field[bool]{Value: true},
					Image:                "maven:3.8-openjdk-11",
					ConnectorRef:         "docker_hub",
					ImagePullPolicy:      "Always",
					Privileged:           &flexible.Field[bool]{Value: false},
					RunAsUser:            &flexible.Field[int]{Value: 1000},
					TestGlobs: "**/*Test.java, **/*Tests.java, **/*IT.java",
					Shell:                "Bash",
					EnvVariables: &flexible.Field[map[string]string]{Value: map[string]string{
						"MAVEN_OPTS": "-Xmx512m",
					}},
					OutputVariables: []*v0.Output{
						{Name: "TEST_RESULT", Type: "String"},
					},
					Reports: &v0.Report{
						Type: "JUnit",
						Spec: &v0.ReportJunit{
							Paths: &flexible.Field[[]string]{Value: []string{"**/target/surefire-reports/*.xml"}},
						},
					},
				},
			},
			expected: &v1.StepTest{
				Script: v1.Stringorslice{"echo \"Starting tests\"\nmvn test -Dmaven.test.failure.ignore=true\necho \"Tests complete\""},
				Shell:  "bash",
				Intelligence: &v1.TestIntelligence{
					Disabled: &flexible.Field[bool]{Value: false},
				},
				Container: &v1.Container{
					Image:      "maven:3.8-openjdk-11",
					Connector:  "docker_hub",
					Pull:       "always",
					Privileged: &flexible.Field[bool]{Value: false},
					User:       &flexible.Field[int]{Value: 1000},
				},
				Env: &flexible.Field[map[string]string]{Value: map[string]string{
					"MAVEN_OPTS": "-Xmx512m",
				}},
				Match:  v1.Stringorslice{"**/*Test.java", "**/*Tests.java", "**/*IT.java"},
				Outputs: []*v1.Output{
					{Name: "TEST_RESULT", Alias: "TEST_RESULT"},
				},
				Report: &v1.Reports{
					Type:  "junit",
					Paths: &flexible.Field[[]string]{Value: []string{"**/target/surefire-reports/*.xml"}},
				},
			},
		},
		{
			name: "Python Pytest basic",
			input: &v0.Step{
				ID:   "pytest_tests",
				Name: "Python Tests",
				Type: v0.StepTypeRunTests,
				Spec: &v0.StepRunTests{
					Language:   "Python",
					BuildTool:  "Pytest",
					Args:       "-v --tb=short",
					PreCommand: "pip install -r requirements.txt",
					Image:      "python:3.9",
					Shell:      "Bash",
				},
			},
			expected: &v1.StepTest{
				Script: v1.Stringorslice{"pip install -r requirements.txt\npytest -v --tb=short"},
				Shell:  "bash",
				Container: &v1.Container{
					Image: "python:3.9",
				},
				Outputs: []*v1.Output{},
			},
		},
		{
			name: "C# Dotnet",
			input: &v0.Step{
				ID:   "dotnet_tests",
				Name: "C# Tests",
				Type: v0.StepTypeRunTests,
				Spec: &v0.StepRunTests{
					Language:   "Csharp",
					BuildTool:  "Dotnet",
					Args:       "--configuration Release --no-build",
					PreCommand: "dotnet restore",
				},
			},
			expected: &v1.StepTest{
				Script: v1.Stringorslice{"dotnet restore\ndotnet test --configuration Release --no-build"},
				Outputs: []*v1.Output{},
			},
		},
		{
			name: "Ruby Rspec",
			input: &v0.Step{
				ID:   "rspec_tests",
				Name: "Ruby Tests",
				Type: v0.StepTypeRunTests,
				Spec: &v0.StepRunTests{
					Language:   "Ruby",
					BuildTool:  "Rspec",
					Args:       "--format documentation --color",
					PreCommand: "bundle install",
				},
			},
			expected: &v1.StepTest{
				Script: v1.Stringorslice{"bundle install\nrspec --format documentation --color"},
				Outputs: []*v1.Output{},
			},
		},
		{
			name: "Scala SBT",
			input: &v0.Step{
				ID:   "sbt_tests",
				Name: "Scala Tests",
				Type: v0.StepTypeRunTests,
				Spec: &v0.StepRunTests{
					Language:   "Scala",
					BuildTool:  "SBT",
					Args:       "-v",
					PreCommand: "export SBT_OPTS='-Xmx2G'",
				},
			},
			expected: &v1.StepTest{
				Script: v1.Stringorslice{"export SBT_OPTS='-Xmx2G'\nsbt test -v"},
				Outputs: []*v1.Output{},
			},
		},
		{
			name: "Java Gradle",
			input: &v0.Step{
				ID:   "gradle_tests",
				Name: "Gradle Tests",
				Type: v0.StepTypeRunTests,
				Spec: &v0.StepRunTests{
					Language:   "Java",
					BuildTool:  "Gradle",
					Args:       "--info --stacktrace",
					PreCommand: "chmod +x gradlew",
				},
			},
			expected: &v1.StepTest{
				Script: v1.Stringorslice{"chmod +x gradlew\n./gradlew test --info --stacktrace"},
				Outputs: []*v1.Output{},
			},
		},
		{
			name: "Java Bazel",
			input: &v0.Step{
				ID:   "bazel_tests",
				Name: "Bazel Tests",
				Type: v0.StepTypeRunTests,
				Spec: &v0.StepRunTests{
					Language:  "Java",
					BuildTool: "Bazel",
					Args:      "//src/test/... --test_output=all",
				},
			},
			expected: &v1.StepTest{
				Script: v1.Stringorslice{"bazel test //src/test/... --test_output=all"},
				Outputs: []*v1.Output{},
			},
		},
		{
			name: "Python Unittest",
			input: &v0.Step{
				ID:   "unittest_tests",
				Name: "Unittest Tests",
				Type: v0.StepTypeRunTests,
				Spec: &v0.StepRunTests{
					Language:  "Python",
					BuildTool: "Unittest",
					Args:      "discover -s tests -p 'test_*.py'",
				},
			},
			expected: &v1.StepTest{
				Script: v1.Stringorslice{"python -m unittest discover -s tests -p 'test_*.py'"},
				Outputs: []*v1.Output{},
			},
		},
		{
			name: "C# NUnit Console",
			input: &v0.Step{
				ID:   "nunit_tests",
				Name: "NUnit Tests",
				Type: v0.StepTypeRunTests,
				Spec: &v0.StepRunTests{
					Language:  "Csharp",
					BuildTool: "Nunitconsole",
					Args:      "--workers=4 --timeout=30000",
					RunOnlySelectedTests: &flexible.Field[bool]{Value: false},
				},
			},
			expected: &v1.StepTest{
				Script: v1.Stringorslice{"nunit3-console --workers=4 --timeout=30000"},
				Outputs: []*v1.Output{},
				Intelligence: &v1.TestIntelligence{
					Disabled: &flexible.Field[bool]{Value: true},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertStepRunTests(tt.input)

			if diff := cmp.Diff(tt.expected, result); diff != "" {
				t.Errorf("ConvertStepRunTests() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertStepRunTests_NilCases(t *testing.T) {
	tests := []struct {
		name  string
		input *v0.Step
	}{
		{
			name:  "nil step",
			input: nil,
		},
		{
			name: "nil spec",
			input: &v0.Step{
				ID:   "nil_spec",
				Name: "Nil Spec",
				Type: v0.StepTypeRunTests,
				Spec: nil,
			},
		},
		{
			name: "wrong spec type",
			input: &v0.Step{
				ID:   "wrong_spec",
				Name: "Wrong Spec",
				Type: v0.StepTypeRunTests,
				Spec: &v0.StepRun{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertStepRunTests(tt.input)
			if result != nil {
				t.Errorf("expected nil result, got %v", result)
			}
		})
	}
}

func TestGenerateTestCommand(t *testing.T) {
	tests := []struct {
		name      string
		language  string
		buildTool string
		args      string
		expected  string
	}{
		{"Maven", "Java", "Maven", "-DskipITs", "mvn test -DskipITs"},
		{"Gradle", "Java", "Gradle", "--info", "./gradlew test --info"},
		{"Bazel", "Java", "Bazel", "//src/test/...", "bazel test //src/test/..."},
		{"SBT", "Scala", "SBT", "-v", "sbt test -v"},
		{"Dotnet", "Csharp", "Dotnet", "--no-build", "dotnet test --no-build"},
		{"NUnit", "Csharp", "Nunitconsole", "--workers=4", "nunit3-console --workers=4"},
		{"Pytest", "Python", "Pytest", "-v", "pytest -v"},
		{"Unittest", "Python", "Unittest", "discover -s tests", "python -m unittest discover -s tests"},
		{"Rspec", "Ruby", "Rspec", "--format doc", "rspec --format doc"},
		{"Unknown", "Unknown", "Unknown", "some args", ""},
		{"Empty args", "Java", "Maven", "", "mvn test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateTestCommand(tt.language, tt.buildTool, tt.args)
			if result != tt.expected {
				t.Errorf("generateTestCommand(%q, %q, %q) = %q, want %q",
					tt.language, tt.buildTool, tt.args, result, tt.expected)
			}
		})
	}
}

func TestBuildTestScript(t *testing.T) {
	tests := []struct {
		name        string
		language    string
		buildTool   string
		preCommand  string
		args        string
		postCommand string
		expected    string
	}{
		{
			name:        "Full script",
			language:    "Java",
			buildTool:   "Maven",
			preCommand:  "echo start",
			args:        "-DskipITs",
			postCommand: "echo done",
			expected:    "echo start\nmvn test -DskipITs\necho done",
		},
		{
			name:        "No pre/post commands",
			language:    "Python",
			buildTool:   "Pytest",
			preCommand:  "",
			args:        "-v",
			postCommand: "",
			expected:    "pytest -v",
		},
		{
			name:        "Only preCommand",
			language:    "Ruby",
			buildTool:   "Rspec",
			preCommand:  "bundle install",
			args:        "",
			postCommand: "",
			expected:    "bundle install\nrspec",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildTestScript(tt.language, tt.buildTool, tt.preCommand, tt.args, tt.postCommand)
			if result != tt.expected {
				t.Errorf("buildTestScript() = %q, want %q", result, tt.expected)
			}
		})
	}
}
