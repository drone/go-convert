package yaml

import (
	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestOnConditions(t *testing.T) {
	tests := []struct {
		yaml string
		want OnConditions
	}{
		{
			yaml: `push`,
			want: OnConditions{
				Push: &PushCondition{},
			},
		},
		{
			yaml: `[push, fork]`,
			want: OnConditions{
				Push: &PushCondition{},
				Fork: &ForkCondition{},
			},
		},
		{
			yaml: `
push:
  branches:
    - main
    - 'releases/**'
`,
			want: OnConditions{
				Push: &PushCondition{
					Branches: []string{"main", "releases/**"},
				},
			},
		},
		{
			yaml: `
label:
  types:
    - created
`,
			want: OnConditions{
				Label: &LabelCondition{
					Types: []string{"created"},
				},
			},
		},
		{
			yaml: `
issues:
  types:
    - opened
    - labeled
`,
			want: OnConditions{
				Issues: &IssuesCondition{
					Types: []string{"opened", "labeled"},
				},
			},
		},
		{
			yaml: `
  push:
    branches:
      - main
      - 'releases/**'
`,
			want: OnConditions{
				Push: &PushCondition{
					Branches: []string{"main", "releases/**"},
				},
			},
		},
		{
			yaml: `
push:
  branches:
    - main
label:
  types:
    - created
page_build: {}
`,
			want: OnConditions{
				Label: &LabelCondition{
					Types: []string{"created"},
				},
				Push: &PushCondition{
					Branches: []string{"main"},
				},
				PageBuild: &PageBuildCondition{},
			},
		},
		{
			yaml: `
label:
  types:
    - created
    - edited
`,
			want: OnConditions{
				Label: &LabelCondition{
					Types: []string{"created", "edited"},
				},
			},
		},
		{
			yaml: `
  pull_request:
    branches:    
      - main
      - 'mona/octocat'
      - 'releases/**'
`,
			want: OnConditions{
				PullRequest: &PullRequestCondition{
					Branches: []string{"main", "mona/octocat", "releases/**"},
				},
			},
		},
		{
			yaml: `
 pull_request:
    branches-ignore:    
      - 'mona/octocat'
      - 'releases/**-alpha'
`,
			want: OnConditions{
				PullRequest: &PullRequestCondition{
					BranchesIgnore: []string{"mona/octocat", "releases/**-alpha"},
				},
			},
		},
		{
			yaml: `
  pull_request:
    branches:    
      - 'releases/**'
      - '!releases/**-alpha'
`,
			want: OnConditions{
				PullRequest: &PullRequestCondition{
					Branches: []string{"releases/**", "!releases/**-alpha"},
				},
			},
		},
		{
			yaml: `
  push:
    # Sequence of patterns matched against refs/heads
    branches:    
      - main
      - 'mona/octocat'
      - 'releases/**'
    # Sequence of patterns matched against refs/tags
    tags:        
      - v2
      - v1.*
`,
			want: OnConditions{
				Push: &PushCondition{
					Branches: []string{"main", "mona/octocat", "releases/**"},
					Tags:     []string{"v2", "v1.*"},
				},
			},
		},
		{
			yaml: `
  push:
    # Sequence of patterns matched against refs/heads
    branches-ignore:    
      - 'mona/octocat'
      - 'releases/**-alpha'
    # Sequence of patterns matched against refs/tags
    tags-ignore:        
      - v2
      - v1.*
`,
			want: OnConditions{
				Push: &PushCondition{
					BranchesIgnore: []string{"mona/octocat", "releases/**-alpha"},
					TagsIgnore:     []string{"v2", "v1.*"},
				},
			},
		},
		{
			yaml: `
push:
  paths:
    - '**.js'
  paths-ignore:
    - '**.md'
`,
			want: OnConditions{
				Push: &PushCondition{
					Paths:       []string{"**.js"},
					PathsIgnore: []string{"**.md"},
				},
			},
		},
		{
			yaml: `
schedule:
  cron: '30 5,17 * * *'
`,
			want: OnConditions{
				Schedule: &ScheduleCondition{
					Cron: []string{"30 5,17 * * *"},
				},
			},
		},
		{
			yaml: `
schedule:
  - cron: '30 5 * * 1,3'
  - cron: '30 5 * * 2,4'
`,
			want: OnConditions{
				Schedule: &ScheduleCondition{
					Cron: []string{"30 5 * * 1,3", "30 5 * * 2,4"},
				},
			},
		},
	}

	for i, test := range tests {
		got := new(OnConditions)
		if err := yaml.Unmarshal([]byte(test.yaml), got); err != nil {
			t.Log(test.yaml)
			t.Error(err)
			return
		}
		if diff := cmp.Diff(got, &test.want); diff != "" {
			t.Log(test.yaml)
			t.Errorf("Unexpected parsing results for test %v", i)
			t.Log(diff)
		}
	}
}

func TestOn_Error(t *testing.T) {
	err := yaml.Unmarshal([]byte("[[]]"), new(Concurrency))
	if err == nil || err.Error() != "failed to unmarshal concurrency" {
		t.Errorf("Expect error, got %s", err)
	}
}
