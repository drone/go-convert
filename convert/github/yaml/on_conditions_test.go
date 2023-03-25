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
				Push:          &PushCondition{},
				ForkCondition: &ForkCondition{},
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
				IssuesCondition: &IssuesCondition{
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
		{
			yaml: `
  gollum
`,
			want: OnConditions{
				GollumCondition: &GollumCondition{},
			},
		},
		{
			yaml: `
  fork
`,
			want: OnConditions{
				ForkCondition: &ForkCondition{},
			},
		},
		{
			yaml: `
  branch_protection_rule:
    types: [created, deleted]
`,
			want: OnConditions{
				BranchProtectionRule: &BranchProtectionRuleCondition{
					Types: []string{"created", "deleted"},
				},
			},
		},
		{
			yaml: `
  check_run:
    types: [requested, completed]
`,
			want: OnConditions{
				CheckRunCondition: &CheckRunCondition{
					Types: []string{"requested", "completed"},
				},
			},
		},
		{
			yaml: `
  check_suite:
    types: [completed]
`,
			want: OnConditions{
				CheckSuiteCondition: &CheckSuiteCondition{
					Types: []string{"completed"},
				},
			},
		},
		{
			yaml: `
  create:
    branches:
      - 'master'
    tags:
      - 'v*'
`,
			want: OnConditions{
				CreateCondition: &CreateCondition{
					Branches: []string{"master"},
					Tags:     []string{"v*"},
				},
			},
		},
		{
			yaml: `
  delete:
    branches: [ main ]
`,
			want: OnConditions{
				DeleteCondition: &DeleteCondition{
					Branches: []string{"main"},
				},
			},
		},
		{
			yaml: `
  discussion_comment:
    types: [created, deleted]
`,
			want: OnConditions{
				DiscussionComment: &DiscussionCommentCondition{
					Types: []string{"created", "deleted"},
				},
			},
		},
		{
			yaml: `
  discussion:
    types:
      - created
      - edited
      - answered
`,
			want: OnConditions{
				DiscussionCondition: &DiscussionCondition{
					Types: []string{"created", "edited", "answered"},
				},
			},
		},
		{
			yaml: `
  issue_comment:
    types: [created, deleted]
`,
			want: OnConditions{
				IssueCommentCondition: &IssueCommentCondition{
					Types: []string{"created", "deleted"},
				},
			},
		},
		{
			yaml: `
  issues:
    types: [opened, edited, milestoned]
`,
			want: OnConditions{
				IssuesCondition: &IssuesCondition{
					Types: []string{"opened", "edited", "milestoned"},
				},
			},
		},
		{
			yaml: `
workflow_call:
  inputs:
    username:
      description: A username passed from the caller workflow
      default: john-doe
      required: false
      type: string
`,
			want: OnConditions{
				WorkflowCall: &WorkflowCallCondition{
					Inputs: map[string]interface{}{
						"username": map[string]interface{}{
							"description": "A username passed from the caller workflow",
							"default":     "john-doe",
							"required":    false,
							"type":        "string",
						},
					},
				},
			},
		},
		{
			yaml: `
workflow_call:
  outputs:
    workflow_output1:
      description: The first job output
      value: ${{ jobs.my_job.outputs.job_output1 }}
    workflow_output2:
      description: The second job output
      value: ${{ jobs.my_job.outputs.job_output2 }}

`,
			want: OnConditions{
				WorkflowCall: &WorkflowCallCondition{
					Outputs: map[string]interface{}{
						"workflow_output1": map[string]interface{}{
							"description": "The first job output",
							"value":       "${{ jobs.my_job.outputs.job_output1 }}",
						},
						"workflow_output2": map[string]interface{}{
							"description": "The second job output",
							"value":       "${{ jobs.my_job.outputs.job_output2 }}",
						},
					},
				},
			},
		},
		{
			yaml: `
workflow_call:
  secrets:
    access-token:
      description: A token passed from the caller workflow
      required: false
`,
			want: OnConditions{
				WorkflowCall: &WorkflowCallCondition{
					Secrets: map[string]WorkflowSecrets{
						"access-token": {
							Description: "A token passed from the caller workflow",
							Required:    false,
						},
					},
				},
			},
		},
		{
			yaml: `
workflow_dispatch:
  inputs:
    logLevel:
      description: Log level
      required: true
      default: warning
      type: choice
      options:
        - info
        - warning
        - debug
    print_tags:
      description: True to print to STDOUT
      required: true
      type: boolean
    tags:
      description: Test scenario tags
      required: true
      type: string
    environment:
      description: Environment to run tests against
      type: environment
`,
			want: OnConditions{
				WorkflowDispatch: &WorkflowDispatchCondition{
					Inputs: map[string]InputDefinition{
						"logLevel": {
							Description: "Log level",
							Required:    true,
							Default:     "warning",
							Type:        "choice",
							Options: []interface{}{
								"info",
								"warning",
								"debug",
							},
						},
						"print_tags": {
							Description: "True to print to STDOUT",
							Required:    true,
							Type:        "boolean",
						},
						"tags": {
							Description: "Test scenario tags",
							Required:    true,
							Type:        "string",
						},
						"environment": {
							Description: "Environment to run tests against",
							Type:        "environment",
						},
					},
				},
			},
		},
		{
			yaml: `
workflow_dispatch:
  inputs:
    logLevel:
      description: Log level
      required: true
      default: warning
      type: choice
      options:
        - info
        - warning
        - debug
    print_tags:
      description: True to print to STDOUT
      required: true
      type: boolean
    tags:
      description: Test scenario tags
      required: true
      type: string
    environment:
      description: Environment to run tests against
      type: environment
      required: true

`,
			want: OnConditions{
				WorkflowDispatch: &WorkflowDispatchCondition{
					Inputs: map[string]InputDefinition{
						"logLevel": {
							Description: "Log level",
							Required:    true,
							Default:     "warning",
							Type:        "choice",
							Options: []interface{}{
								"info",
								"warning",
								"debug",
							},
						},
						"print_tags": {
							Description: "True to print to STDOUT",
							Required:    true,
							Type:        "boolean",
						},
						"tags": {
							Description: "Test scenario tags",
							Required:    true,
							Type:        "string",
						},
						"environment": {
							Description: "Environment to run tests against",
							Type:        "environment",
							Required:    true,
						},
					},
				},
			},
		},
		{
			yaml: `
workflow_run:
  workflows:
    - Build
  types:
    - requested
  branches:
    - releases/**
`,
			want: OnConditions{
				WorkflowRun: &WorkflowRunCondition{
					Workflows: []string{"Build"},
					Types:     []string{"requested"},
					Branches:  []string{"releases/**"},
				},
			},
		},
		{
			yaml: `
workflow_run:
  workflows:
    - Build
  types:
    - requested
  branches:
    - releases/**
    - "!releases/**-alpha"
`,
			want: OnConditions{
				WorkflowRun: &WorkflowRunCondition{
					Workflows: []string{"Build"},
					Types:     []string{"requested"},
					Branches:  []string{"releases/**", "!releases/**-alpha"},
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
