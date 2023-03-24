package yaml

import (
	"errors"
)

type OnConditions struct {
	BranchProtectionRule     *BranchProtectionRuleCondition     `yaml:"branch_protection_rule,omitempty"`
	CheckRunCondition        *CheckRunCondition                 `yaml:"check_run,omitempty"`
	CheckSuiteCondition      *CheckSuiteCondition               `yaml:"check_suite,omitempty"`
	CreateCondition          *CreateCondition                   `yaml:"create,omitempty"`
	DeleteCondition          *DeleteCondition                   `yaml:"delete,omitempty"`
	DiscussionCondition      *DiscussionCondition               `yaml:"discussion,omitempty"`
	ForkCondition            *ForkCondition                     `yaml:"fork,omitempty"`
	IssueCommentCondition    *IssueCommentCondition             `yaml:"issue_comment,omitempty"`
	IssuesCondition          *IssuesCondition                   `yaml:"issues,omitempty"`
	Label                    *LabelCondition                    `yaml:"label,omitempty"`
	DiscussionComment        *DiscussionCommentCondition        `yaml:"discussion_comment,omitempty"`
	GollumCondition          *GollumCondition                   `yaml:"gollum,omitempty"`
	Milestone                *MilestoneCondition                `yaml:"milestone,omitempty"`
	PageBuild                *PageBuildCondition                `yaml:"page_build,omitempty"`
	Project                  *ProjectCondition                  `yaml:"project,omitempty"`
	ProjectCard              *ProjectCardCondition              `yaml:"project_card,omitempty"`
	ProjectColumn            *ProjectColumnCondition            `yaml:"project_column,omitempty"`
	Public                   *PublicCondition                   `yaml:"public,omitempty"`
	PullRequestReview        *PullRequestReviewCondition        `yaml:"pull_request_review,omitempty"`
	PullRequestReviewComment *PullRequestReviewCommentCondition `yaml:"pull_request_review_comment,omitempty"`
	PullRequestTarget        *PullRequestTargetCondition        `yaml:"pull_request_target,omitempty"`
	Push                     *PushCondition                     `yaml:"push,omitempty"`
	PullRequest              *PullRequestCondition              `yaml:"pull_request,omitempty"`
	RegistryPackage          *RegistryPackageCondition          `yaml:"registry_package,omitempty"`
	Release                  *ReleaseCondition                  `yaml:"release,omitempty"`
	RepositoryDispatch       *RepositoryDispatchCondition       `yaml:"repository_dispatch,omitempty"`
	Schedule                 *ScheduleCondition                 `yaml:"schedule,omitempty"`
	Status                   *StatusCondition                   `yaml:"status,omitempty"`
	Watch                    *WatchCondition                    `yaml:"watch,omitempty"`
	WorkflowCall             *WorkflowCallCondition             `yaml:"workflow_call,omitempty"`
	WorkflowDispatch         *WorkflowDispatchCondition         `yaml:"workflow_dispatch,omitempty"`
	WorkflowRun              *WorkflowRunCondition              `yaml:"workflow_run,omitempty"`
}

type (
	GollumCondition struct{}

	ForkCondition struct{}

	LabelCondition struct {
		Types []string `yaml:"types,omitempty"`
	}

	MilestoneCondition struct {
		Types []string `yaml:"types,omitempty"`
	}

	PageBuildCondition struct{}

	ProjectCondition struct {
		Types []string `yaml:"types,omitempty"`
	}

	ProjectCardCondition struct {
		Types []string `yaml:"types,omitempty"`
	}

	ProjectColumnCondition struct {
		Types []string `yaml:"types,omitempty"`
	}

	PublicCondition struct{}

	PullRequestCondition struct {
		Branches        []string `yaml:"branches,omitempty"`
		BranchesIgnore  []string `yaml:"branches-ignore,omitempty"`
		Paths           []string `yaml:"paths,omitempty"`
		PathsIgnore     []string `yaml:"paths-ignore,omitempty"`
		Tags            []string `yaml:"tags,omitempty"`
		TagsIgnore      []string `yaml:"tags-ignore,omitempty"`
		ReviewApproved  bool     `yaml:"review-approved,omitempty"`
		ReviewDismissed bool     `yaml:"review-dismissed,omitempty"`
	}

	PullRequestReviewCondition struct {
		Types []string `yaml:"types,omitempty"`
	}

	PullRequestReviewCommentCondition struct {
		Types []string `yaml:"types,omitempty"`
	}

	PullRequestTargetCondition struct {
		Branches       []string `yaml:"branches,omitempty"`
		BranchesIgnore []string `yaml:"branches-ignore,omitempty"`
		Types          []string `yaml:"types,omitempty"`
	}

	RegistryPackageCondition struct {
		Types []string `yaml:"types,omitempty"`
	}

	ReleaseCondition struct {
		Types []string `yaml:"types,omitempty"`
	}

	RepositoryDispatchCondition struct {
		Types []string `yaml:"types,omitempty"`
	}

	ScheduleCondition struct {
		Cron []string `yaml:"cron,omitempty"`
	}

	StatusCondition struct{}

	WatchCondition struct {
		Types []string `yaml:"types,omitempty"`
	}

	WorkflowCallSecrets map[string]struct {
		Description string `yaml:"description,omitempty"`
		Required    bool   `yaml:"required,omitempty"`
	}
)

// UnmarshalYAML implements the unmarshal interface for WorkflowTriggers.
func (v *OnConditions) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out0 string
	var out1 []string
	var out2 = struct {
		BranchProtectionRule     *BranchProtectionRuleCondition     `yaml:"branch_protection_rule,omitempty"`
		CheckRunCondition        *CheckRunCondition                 `yaml:"check_run,omitempty"`
		CheckSuiteCondition      *CheckSuiteCondition               `yaml:"check_suite,omitempty"`
		CreateCondition          *CreateCondition                   `yaml:"create,omitempty"`
		DeleteCondition          *DeleteCondition                   `yaml:"delete,omitempty"`
		DiscussionCondition      *DiscussionCondition               `yaml:"discussion,omitempty"`
		ForkCondition            *ForkCondition                     `yaml:"fork,omitempty"`
		IssueCommentCondition    *IssueCommentCondition             `yaml:"issue_comment,omitempty"`
		IssuesCondition          *IssuesCondition                   `yaml:"issues,omitempty"`
		Label                    *LabelCondition                    `yaml:"label,omitempty"`
		DiscussionComment        *DiscussionCommentCondition        `yaml:"discussion_comment,omitempty"`
		GollumCondition          *GollumCondition                   `yaml:"gollum,omitempty"`
		Milestone                *MilestoneCondition                `yaml:"milestone,omitempty"`
		PageBuild                *PageBuildCondition                `yaml:"page_build,omitempty"`
		Project                  *ProjectCondition                  `yaml:"project,omitempty"`
		ProjectCard              *ProjectCardCondition              `yaml:"project_card,omitempty"`
		ProjectColumn            *ProjectColumnCondition            `yaml:"project_column,omitempty"`
		Public                   *PublicCondition                   `yaml:"public,omitempty"`
		PullRequestReview        *PullRequestReviewCondition        `yaml:"pull_request_review,omitempty"`
		PullRequestReviewComment *PullRequestReviewCommentCondition `yaml:"pull_request_review_comment,omitempty"`
		PullRequestTarget        *PullRequestTargetCondition        `yaml:"pull_request_target,omitempty"`
		Push                     *PushCondition                     `yaml:"push,omitempty"`
		PullRequest              *PullRequestCondition              `yaml:"pull_request,omitempty"`
		RegistryPackage          *RegistryPackageCondition          `yaml:"registry_package,omitempty"`
		Release                  *ReleaseCondition                  `yaml:"release,omitempty"`
		RepositoryDispatch       *RepositoryDispatchCondition       `yaml:"repository_dispatch,omitempty"`
		Schedule                 *ScheduleCondition                 `yaml:"schedule,omitempty"`
		Status                   *StatusCondition                   `yaml:"status,omitempty"`
		Watch                    *WatchCondition                    `yaml:"watch,omitempty"`
		WorkflowCall             *WorkflowCallCondition             `yaml:"workflow_call,omitempty"`
		WorkflowDispatch         *WorkflowDispatchCondition         `yaml:"workflow_dispatch,omitempty"`
		WorkflowRun              *WorkflowRunCondition              `yaml:"workflow_run,omitempty"`
	}{}

	if err := unmarshal(&out0); err == nil {
		v.setEvent(out0)
		return nil
	}

	if err := unmarshal(&out1); err == nil {
		for _, event := range out1 {
			v.setEvent(event)
		}
		return nil
	}

	if err := unmarshal(&out2); err == nil {
		*v = out2
		return nil
	}

	return errors.New("failed to unmarshal on conditions")
}

func (s *ScheduleCondition) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var items []map[string]string
	if err := unmarshal(&items); err != nil {
		var m map[string]string
		if err2 := unmarshal(&m); err2 != nil {
			return err
		}
		if cron, ok := m["cron"]; ok {
			s.Cron = []string{cron}
		}
		return nil
	}
	s.Cron = make([]string, len(items))
	for i, item := range items {
		if cron, ok := item["cron"]; ok {
			s.Cron[i] = cron
		}
	}
	return nil
}

func (v *OnConditions) setEvent(event string) {
	switch event {
	case "branch_protection_rule":
		v.BranchProtectionRule = &BranchProtectionRuleCondition{}
	case "check_run":
		v.CheckRunCondition = &CheckRunCondition{}
	case "check_suite":
		v.CheckSuiteCondition = &CheckSuiteCondition{}
	case "create":
		v.CreateCondition = &CreateCondition{}
	case "delete":
		v.DeleteCondition = &DeleteCondition{}
	case "discussion":
		v.DiscussionCondition = &DiscussionCondition{}
	case "push":
		v.Push = &PushCondition{}
	case "pull_request":
		v.PullRequest = &PullRequestCondition{}
	case "fork":
		v.ForkCondition = &ForkCondition{}
	case "issue_comment":
		v.IssueCommentCondition = &IssueCommentCondition{}
	case "issues":
		v.IssuesCondition = &IssuesCondition{}
	case "label":
		v.Label = &LabelCondition{}
	case "discussion_comment":
		v.DiscussionComment = &DiscussionCommentCondition{}
	case "gollum":
		v.GollumCondition = &GollumCondition{}
	case "milestone":
		v.Milestone = &MilestoneCondition{}
	case "page_build":
		v.PageBuild = &PageBuildCondition{}
	case "project":
		v.Project = &ProjectCondition{}
	case "project_card":
		v.ProjectCard = &ProjectCardCondition{}
	case "project_column":
		v.ProjectColumn = &ProjectColumnCondition{}
	case "public":
		v.Public = &PublicCondition{}
	case "pull_request_review":
		v.PullRequestReview = &PullRequestReviewCondition{}
	case "pull_request_review_comment":
		v.PullRequestReviewComment = &PullRequestReviewCommentCondition{}
	case "pull_request_target":
		v.PullRequestTarget = &PullRequestTargetCondition{}
	case "registry_package":
		v.RegistryPackage = &RegistryPackageCondition{}
	case "release":
		v.Release = &ReleaseCondition{}
	case "repository_dispatch":
		v.RepositoryDispatch = &RepositoryDispatchCondition{}
	case "schedule":
		v.Schedule = &ScheduleCondition{}
	case "status":
		v.Status = &StatusCondition{}
	case "watch":
		v.Watch = &WatchCondition{}
	case "workflow_call":
		v.WorkflowCall = &WorkflowCallCondition{}
	case "workflow_dispatch":
		v.WorkflowDispatch = &WorkflowDispatchCondition{}
	case "workflow_run":
		v.WorkflowRun = &WorkflowRunCondition{}
	}
}
