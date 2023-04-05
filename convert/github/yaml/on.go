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

import "errors"

type On struct {
	BranchProtectionRule     *Event             `yaml:"branch_protection_rule,omitempty"`
	CheckRun                 *Event             `yaml:"check_run,omitempty"`
	CheckSuite               *Event             `yaml:"check_suite,omitempty"`
	Create                   *struct{}          `yaml:"create,omitempty"`
	Delete                   *struct{}          `yaml:"delete,omitempty"`
	Deployment               *struct{}          `yaml:"deployment,omitempty"`
	DeploymentStatus         *struct{}          `yaml:"deployment_status,omitempty"`
	Discussion               *Event             `yaml:"discussion,omitempty"`
	DiscussionComment        *Event             `yaml:"discussion_comment,omitempty"`
	Fork                     *struct{}          `yaml:"fork,omitempty"`
	Gollum                   *struct{}          `yaml:"gollum,omitempty"`
	IssueComment             *Event             `yaml:"issue_comment,omitempty"`
	Issues                   *Event             `yaml:"issues,omitempty"`
	Label                    *Event             `yaml:"label,omitempty"`
	Member                   *Event             `yaml:"member,omitempty"`
	MergeGroup               *Event             `yaml:"merge_group,omitempty"`
	Milestone                *Event             `yaml:"milestone,omitempty"`
	PageBuild                *struct{}          `yaml:"page_build,omitempty"`
	Project                  *Event             `yaml:"project,omitempty"`
	ProjectCard              *Event             `yaml:"project_card,omitempty"`
	ProjectColumn            *Event             `yaml:"project_column,omitempty"`
	Public                   *struct{}          `yaml:"public,omitempty"`
	PullRequest              *PullRequest       `yaml:"pull_request,omitempty"`
	PullRequestReview        *Event             `yaml:"pull_request_review,omitempty"`
	PullRequestReviewComment *Event             `yaml:"pull_request_review_comment,omitempty"`
	PullRequestTarget        *PullRequestTarget `yaml:"pull_request_target,omitempty"`
	Push                     *Push              `yaml:"push,omitempty"`
	RegistryPackage          *Event             `yaml:"registry_package,omitempty"`
	RepositoryDispatch       *Event             `yaml:"repository_dispatch,omitempty"`
	Release                  *Event             `yaml:"release,omitempty"`
	Schedule                 *Schedule          `yaml:"schedule,omitempty"`
	Status                   *struct{}          `yaml:"status,omitempty"`
	Watch                    *Event             `yaml:"watch,omitempty"`
	WorkflowCall             *WorkflowCall      `yaml:"workflow_call,omitempty"`
	WorkflowDispatch         *WorkflowDispatch  `yaml:"workflow_dispatch,omitempty"`
	WorkflowRun              *WorkflowRun       `yaml:"workflow_run,omitempty"`
}

// UnmarshalYAML implements the unmarshal interface for WorkflowTriggers.
func (v *On) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 string
	var out2 []string

	if err := unmarshal(&out1); err == nil {
		v.setEvent(out1)
		return nil
	}
	if err := unmarshal(&out2); err == nil {
		for _, s := range out2 {
			v.setEvent(s)
		}
		return nil
	}

	out3 := struct {
		BranchProtectionRule     *Event             `yaml:"branch_protection_rule,omitempty"`
		CheckRun                 *Event             `yaml:"check_run,omitempty"`
		CheckSuite               *Event             `yaml:"check_suite,omitempty"`
		Create                   *struct{}          `yaml:"create,omitempty"`
		Delete                   *struct{}          `yaml:"delete,omitempty"`
		Deployment               *struct{}          `yaml:"deployment,omitempty"`
		DeploymentStatus         *struct{}          `yaml:"deployment_status,omitempty"`
		Discussion               *Event             `yaml:"discussion,omitempty"`
		DiscussionComment        *Event             `yaml:"discussion_comment,omitempty"`
		Fork                     *struct{}          `yaml:"fork,omitempty"`
		Gollum                   *struct{}          `yaml:"gollum,omitempty"`
		IssueComment             *Event             `yaml:"issue_comment,omitempty"`
		Issues                   *Event             `yaml:"issues,omitempty"`
		Label                    *Event             `yaml:"label,omitempty"`
		Member                   *Event             `yaml:"member,omitempty"`
		MergeGroup               *Event             `yaml:"merge_group,omitempty"`
		Milestone                *Event             `yaml:"milestone,omitempty"`
		PageBuild                *struct{}          `yaml:"page_build,omitempty"`
		Project                  *Event             `yaml:"project,omitempty"`
		ProjectCard              *Event             `yaml:"project_card,omitempty"`
		ProjectColumn            *Event             `yaml:"project_column,omitempty"`
		Public                   *struct{}          `yaml:"public,omitempty"`
		PullRequest              *PullRequest       `yaml:"pull_request,omitempty"`
		PullRequestReview        *Event             `yaml:"pull_request_review,omitempty"`
		PullRequestReviewComment *Event             `yaml:"pull_request_review_comment,omitempty"`
		PullRequestTarget        *PullRequestTarget `yaml:"pull_request_target,omitempty"`
		Push                     *Push              `yaml:"push,omitempty"`
		RegistryPackage          *Event             `yaml:"registry_package,omitempty"`
		RepositoryDispatch       *Event             `yaml:"repository_dispatch,omitempty"`
		Release                  *Event             `yaml:"release,omitempty"`
		Schedule                 *Schedule          `yaml:"schedule,omitempty"`
		Status                   *struct{}          `yaml:"status,omitempty"`
		Watch                    *Event             `yaml:"watch,omitempty"`
		WorkflowCall             *WorkflowCall      `yaml:"workflow_call,omitempty"`
		WorkflowDispatch         *WorkflowDispatch  `yaml:"workflow_dispatch,omitempty"`
		WorkflowRun              *WorkflowRun       `yaml:"workflow_run,omitempty"`
	}{}
	if err := unmarshal(&out3); err == nil {
		*v = out3
		return nil
	}

	return errors.New("failed to unmarshal on")
}

func (v *On) setEvent(event string) {
	switch event {
	case "branch_protection_rule":
		v.BranchProtectionRule = new(Event)
	case "check_run":
		v.CheckRun = new(Event)
	case "check_suite":
		v.CheckSuite = new(Event)
	case "create":
		v.Create = new(struct{})
	case "delete":
		v.Delete = new(struct{})
	case "deployment":
		v.Deployment = new(struct{})
	case "deployment_status":
		v.DeploymentStatus = new(struct{})
	case "discussion":
		v.Discussion = new(Event)
	case "discussion_comment":
		v.DiscussionComment = new(Event)
	case "fork":
		v.Fork = new(struct{})
	case "gollum":
		v.Gollum = new(struct{})
	case "issue_comment":
		v.IssueComment = new(Event)
	case "issues":
		v.Issues = new(Event)
	case "label":
		v.Label = new(Event)
	case "member":
		v.Member = new(Event)
	case "merge_group":
		v.MergeGroup = new(Event)
	case "milestone":
		v.Milestone = new(Event)
	case "page_build":
		v.PageBuild = new(struct{})
	case "project":
		v.Project = new(Event)
	case "project_card":
		v.ProjectCard = new(Event)
	case "project_column":
		v.ProjectColumn = new(Event)
	case "public":
		v.Public = new(struct{})
	case "pull_request":
		v.PullRequest = new(PullRequest)
	case "pull_request_review":
		v.PullRequestReview = new(Event)
	case "pull_request_review_comment":
		v.PullRequestReviewComment = new(Event)
	case "pull_request_target":
		v.PullRequestTarget = new(PullRequestTarget)
	case "push":
		v.Push = new(Push)
	case "registry_package":
		v.RegistryPackage = new(Event)
	case "repository_dispatch":
		v.RepositoryDispatch = new(Event)
	case "release":
		v.Release = new(Event)
	case "schedule":
		v.Schedule = new(Schedule)
	case "status":
		v.Status = new(struct{})
	case "watch":
		v.Watch = new(Event)
	case "workflow_call":
		v.WorkflowCall = new(WorkflowCall)
	case "workflow_dispatch":
		v.WorkflowDispatch = new(WorkflowDispatch)
	case "workflow_run":
		v.WorkflowRun = new(WorkflowRun)
	}
}

// // ApplyDefaults is a helper funciton to apply default actions
// // to events if the event is non-nil but empty.
// func (v *On) ApplyDefaults() {
// 	set := func(event *Event, actions []string) {
// 		if event != nil && len(event.Types) == 0 {
// 			event.Types = append(event.Types, actions...)
// 		}
// 	}

// 	set(v.BranchProtectionRule, []string{"created", "edited", "deleted"})
// 	set(v.CheckRun, []string{"created", "rerequested", "completed", "requested_action"})
// 	set(v.CheckSuite, []string{"completed", "requested", "rerequested"})
// 	set(v.Discussion, []string{"created", "edited", "deleted", "transferred", "pinned", "unpinned", "labeled", "unlabeled", "locked", "unlocked", "category_changed", "answered", "unanswered"})
// 	set(v.DiscussionComment, []string{"created", "edited", "deleted"})
// 	set(v.IssueComment, []string{"created", "edited", "deleted"})
// 	set(v.Issues, []string{"opened", "edited", "deleted", "transferred", "pinned", "unpinned", "closed", "reopened", "assigned", "unassigned", "labeled", "unlabeled", "locked", "unlocked", "milestoned", "demilestoned"})
// 	set(v.Label, []string{"created", "edited", "deleted"})
// 	set(v.Member, []string{"added", "edited", "deleted"})
// 	set(v.MergeGroup, []string{"checks_requested"})
// 	set(v.Milestone, []string{"created", "closed", "opened", "edited", "deleted"})
// 	set(v.Project, []string{"created", "updated", "closed", "reopened", "edited", "deleted"})
// 	set(v.ProjectCard, []string{"created", "moved", "converted", "edited", "deleted"})
// 	set(v.ProjectColumn, []string{"created", "updated", "moved", "deleted"})
// 	set(v.PullRequestReview, []string{"submitted", "edited", "dismissed"})
// 	set(v.PullRequestReviewComment, []string{"created", "edited", "deleted"})
// 	set(v.RegistryPackage, []string{"published", "updated"})
// 	set(v.Release, []string{"published", "unpublished", "created", "edited", "deleted", "prereleased", "released"})
// 	set(v.Watch, []string{"started"})

// 	if v.PullRequest != nil && len(v.PullRequest.Types) != 0 {
// 		v.PullRequest.Types = []string{"opened", "synchronize", "reopened"}
// 	}
// 	if v.PullRequestTarget != nil && len(v.PullRequestTarget.Types) != 0 {
// 		v.PullRequestTarget.Types = []string{"opened", "synchronize", "reopened"}
// 	}
// }
