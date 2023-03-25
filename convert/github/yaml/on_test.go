// Copyright 2022 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package yaml

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
)

func TestOnConditions(t *testing.T) {
	tests := []struct {
		yaml string
		want On
	}{
		//
		// test string
		//
		{
			yaml: `push`,
			want: On{
				Push: &Push{},
			},
		},
		//
		// test slice
		//
		{
			yaml: `[push, fork]`,
			want: On{
				Push: &Push{},
				Fork: &struct{}{},
			},
		},
		//
		// test object
		//
		{
			yaml: `{ push: { branches: [ main, 'releases/**' ] } }`,
			want: On{
				Push: &Push{
					Branches: []string{"main", "releases/**"},
				},
			},
		},
		//
		// test individual keywords to ensure non-nil structs
		// and default event types
		//
		{
			yaml: `branch_protection_rule`,
			want: On{BranchProtectionRule: &Event{}},
		},
		{
			yaml: `check_run`,
			want: On{CheckRun: &Event{}},
		},
		{
			yaml: `check_suite`,
			want: On{CheckSuite: &Event{}},
		},
		{
			yaml: `create`,
			want: On{Create: &struct{}{}},
		},
		{
			yaml: `delete`,
			want: On{Delete: &struct{}{}},
		},
		{
			yaml: `deployment`,
			want: On{Deployment: &struct{}{}},
		},
		{
			yaml: `deployment_status`,
			want: On{DeploymentStatus: &struct{}{}},
		},
		{
			yaml: `discussion`,
			want: On{Discussion: &Event{}},
		},
		{
			yaml: `discussion_comment`,
			want: On{DiscussionComment: &Event{}},
		},
		{
			yaml: `fork`,
			want: On{Fork: &struct{}{}},
		},
		{
			yaml: `gollum`,
			want: On{Gollum: &struct{}{}},
		},
		{
			yaml: `issue_comment`,
			want: On{IssueComment: &Event{}},
		},
		{
			yaml: `issues`,
			want: On{Issues: &Event{}},
		},
		{
			yaml: `label`,
			want: On{Label: &Event{}},
		},
		{
			yaml: `member`,
			want: On{Member: &Event{}},
		},
		{
			yaml: `merge_group`,
			want: On{MergeGroup: &Event{}},
		},
		{
			yaml: `milestone`,
			want: On{Milestone: &Event{}},
		},
		{
			yaml: `page_build`,
			want: On{PageBuild: &struct{}{}},
		},
		{
			yaml: `project`,
			want: On{Project: &Event{}},
		},
		{
			yaml: `project_card`,
			want: On{ProjectCard: &Event{}},
		},
		{
			yaml: `project_column`,
			want: On{ProjectColumn: &Event{}},
		},
		{
			yaml: `public`,
			want: On{Public: &struct{}{}},
		},
		{
			yaml: `pull_request`,
			want: On{PullRequest: &PullRequest{}},
		},
		{
			yaml: `pull_request_review`,
			want: On{PullRequestReview: &Event{}},
		},
		{
			yaml: `pull_request_review_comment`,
			want: On{PullRequestReviewComment: &Event{}},
		},
		{
			yaml: `pull_request_target`,
			want: On{PullRequestTarget: &PullRequestTarget{}},
		},
		{
			yaml: `push`,
			want: On{Push: &Push{}},
		},
		{
			yaml: `registry_package`,
			want: On{RegistryPackage: &Event{}},
		},
		{
			yaml: `repository_dispatch`,
			want: On{RepositoryDispatch: &Event{}},
		},
		{
			yaml: `release`,
			want: On{Release: &Event{}},
		},
		{
			yaml: `schedule`,
			want: On{Schedule: &Schedule{}},
		},
		{
			yaml: `status`,
			want: On{Status: &struct{}{}},
		},
		{
			yaml: `watch`,
			want: On{Watch: &Event{}},
		},
		{
			yaml: `workflow_call`,
			want: On{WorkflowCall: &WorkflowCall{}},
		},
		{
			yaml: `workflow_dispatch`,
			want: On{WorkflowDispatch: &WorkflowDispatch{}},
		},
		{
			yaml: `workflow_run`,
			want: On{WorkflowRun: &WorkflowRun{}},
		},
	}

	for i, test := range tests {
		got := new(On)
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
	err := yaml.Unmarshal([]byte("[[]]"), new(On))
	if err == nil || err.Error() != "failed to unmarshal on" {
		t.Errorf("Expect error, got %s", err)
	}
}
