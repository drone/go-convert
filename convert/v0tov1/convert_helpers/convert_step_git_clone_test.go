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
				"connector":    "github-connector",
				"build_target": "Git Branch",
				"branch":       "main",
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
				"connector":    "github-connector",
				"build_target": "Tag",
				"tag":          "v1.0.0",
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
				"connector":    "github-connector",
				"build_target": "Pull Request",
				"pr":           "42",
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
				"connector":    "github-connector",
				"build_target": "Commit",
				"commit_sha":   "abc123def456",
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
				"repo_name": "my-repo",
			},
		},
		{
			name: "all optional fields",
			step: &v0.Step{
				Spec: &v0.StepGitClone{
					ConnRef:           "github-connector",
					Repository:        "my-repo",
					CloneDirectory:    "/workspace/src",
					Depth:             &flexible.Field[int]{Value: 50},
					Lfs:               &flexible.Field[bool]{Value: true},
					Debug:             &flexible.Field[bool]{Value: true},
					FetchTags:         &flexible.Field[bool]{Value: true},
					SparseCheckout:    &flexible.Field[[]string]{Value: []string{"src", "lib"}},
					SubmoduleStrategy: &flexible.Field[bool]{Value: "recursive"},
					PreFetchCommand:   "git config --global user.email test@test.com",
					BuildType: &flexible.Field[v0.Build]{Value: v0.Build{
						Type: "branch",
						Spec: v0.BuildSpec{Branch: "develop"},
					}},
				},
			},
			expected: map[string]interface{}{
				"connector":          "github-connector",
				"repo_name":          "my-repo",
				"build_target":       "Git Branch",
				"branch":             "develop",
				"clone_directory":    "/workspace/src",
				"depth":              "50",
				"lfs_enabled":        true,
				"debug":              true,
				"fetch_tags":         true,
				"sparse_checkout":    "src,lib",
				"submodule_strategy": "recursive",
				"pre_fetch":          "git config --global user.email test@test.com",
			},
		},
		{
			name: "minimal step with no optional fields",
			step: &v0.Step{
				Spec: &v0.StepGitClone{},
			},
			expected: map[string]interface{}{},
		},
		{
			name: "expression in build type",
			step: &v0.Step{
				Spec: &v0.StepGitClone{
					ConnRef:   "github-connector",
					BuildType: &flexible.Field[v0.Build]{Value: "<+input>"},
				},
			},
			expected: map[string]interface{}{
				"connector":    "github-connector",
			},
		},
		{
			name: "expression in depth",
			step: &v0.Step{
				Spec: &v0.StepGitClone{
					ConnRef: "github-connector",
					Depth:   &flexible.Field[int]{Value: "<+pipeline.variables.clone_depth>"},
				},
			},
			expected: map[string]interface{}{
				"connector": "github-connector",
				"depth":     "<+pipeline.variables.clone_depth>",
			},
		},
		{
			name: "output file paths content",
			step: &v0.Step{
				Spec: &v0.StepGitClone{
					ConnRef:                "github-connector",
					OutputFilePathsContent: &flexible.Field[[]string]{Value: []string{"file1.txt", "file2.txt"}},
				},
			},
			expected: map[string]interface{}{
				"connector":          "github-connector",
				"file_paths_content": "file1.txt,file2.txt",
			},
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
