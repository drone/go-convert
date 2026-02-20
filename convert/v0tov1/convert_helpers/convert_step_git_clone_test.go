package converthelpers

import (
	"testing"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	"github.com/drone/go-convert/internal/flexible"
	"github.com/google/go-cmp/cmp"
)

func TestConvertStepGitClone(t *testing.T) {
	tests := []struct {
		name     string
		step     *v0.Step
		expected map[string]interface{}
	}{
		{
			name: "branch build type",
			step: &v0.Step{
				Spec: &v0.StepGitClone{
					ConnRef: "github-connector",
					BuildType: &flexible.Field[v0.Build]{Value: v0.Build{
						Type: "branch",
						Spec: v0.BuildSpec{Branch: "main"},
					}},
				},
			},
			expected: map[string]interface{}{
				"connector": "github-connector",
				"branch":    "main",
			},
		},
		{
			name: "tag build type",
			step: &v0.Step{
				Spec: &v0.StepGitClone{
					ConnRef: "github-connector",
					BuildType: &flexible.Field[v0.Build]{Value: v0.Build{
						Type: "tag",
						Spec: v0.BuildSpec{Tag: "v1.0.0"},
					}},
				},
			},
			expected: map[string]interface{}{
				"connector": "github-connector",
				"tag":       "v1.0.0",
			},
		},
		{
			name: "PR build type",
			step: &v0.Step{
				Spec: &v0.StepGitClone{
					ConnRef: "github-connector",
					BuildType: &flexible.Field[v0.Build]{Value: v0.Build{
						Type: "PR",
						Spec: v0.BuildSpec{Number: &flexible.Field[int]{Value: 42}},
					}},
				},
			},
			expected: map[string]interface{}{
				"connector": "github-connector",
				"pr":        &flexible.Field[int]{Value: 42},
			},
		},
		{
			name: "commitSha build type",
			step: &v0.Step{
				Spec: &v0.StepGitClone{
					ConnRef: "github-connector",
					BuildType: &flexible.Field[v0.Build]{Value: v0.Build{
						Type: "commitSha",
						Spec: v0.BuildSpec{CommitSha: "abc123def456"},
					}},
				},
			},
			expected: map[string]interface{}{
				"connector": "github-connector",
				"commitSha": "abc123def456",
			},
		},
		{
			name: "nil build spec",
			step: &v0.Step{
				Spec: &v0.StepGitClone{
					ConnRef:    "github-connector",
					Repository: "my-repo",
				},
			},
			expected: map[string]interface{}{
				"connector": "github-connector",
				"repoName":  "my-repo",
			},
		},
		{
			name: "all optional fields",
			step: &v0.Step{
				Spec: &v0.StepGitClone{
					ConnRef:        "github-connector",
					Repository:     "my-repo",
					CloneDirectory: "/workspace/src",
					Depth:          &flexible.Field[int]{Value: 50},
					SSLVerify:      "false",
					BuildType: &flexible.Field[v0.Build]{Value: v0.Build{
						Type: "branch",
						Spec: v0.BuildSpec{Branch: "develop"},
					}},
				},
			},
			expected: map[string]interface{}{
				"connector":      "github-connector",
				"repoName":       "my-repo",
				"branch":         "develop",
				"cloneDirectory": "/workspace/src",
				"depth":          &flexible.Field[int]{Value: 50},
				"sslVerify":      "false",
			},
		},
		{
			name: "minimal step with no optional fields",
			step: &v0.Step{
				Spec: &v0.StepGitClone{},
			},
			expected: map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertStepGitClone(tt.step)
			if result == nil {
				t.Fatal("expected non-nil result")
			}

			if result.Uses != "gitCloneStep" {
				t.Errorf("expected Uses to be gitCloneStep, got %s", result.Uses)
			}

			if diff := cmp.Diff(tt.expected, result.With); diff != "" {
				t.Errorf("With mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertStepGitClone_NilCases(t *testing.T) {
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
			step: &v0.Step{
				Spec: nil,
			},
		},
		{
			name: "wrong spec type",
			step: &v0.Step{
				Spec: &v0.StepRun{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertStepGitClone(tt.step)
			if result != nil {
				t.Errorf("expected nil result, got %v", result)
			}
		})
	}
}
