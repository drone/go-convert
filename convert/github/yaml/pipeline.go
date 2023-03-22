package yaml

type Pipeline struct {
	Name        string            `yaml:"name,omitempty"`
	On          WorkflowTriggers  `yaml:"on,omitempty"`
	Jobs        map[string]Job    `yaml:"jobs,omitempty"`
	Environment map[string]string `yaml:"env,omitempty"`
}

type Event struct {
	Branches []string `yaml:"branches,omitempty"`
	Tags     []string `yaml:"tags,omitempty"`
	Paths    []string `yaml:"paths,omitempty"`
}

type Job struct {
	RunsOn      string              `yaml:"runs-on,omitempty"`
	Container   string              `yaml:"container,omitempty"`
	Services    map[string]*Service `yaml:"services,omitempty"`
	Steps       []*Step             `yaml:"steps,omitempty"`
	Environment map[string]string   `yaml:"env,omitempty"`
	Strategy    *Strategy           `yaml:"strategy,omitempty"`
}

type Step struct {
	Name        string                 `yaml:"name,omitempty"`
	Uses        string                 `yaml:"uses,omitempty"`
	With        map[string]interface{} `yaml:"with,omitempty"`
	Run         string                 `yaml:"run,omitempty"`
	If          string                 `yaml:"if,omitempty"`
	Environment map[string]string      `yaml:"env,omitempty"`
}

type Strategy struct {
	Matrix *Matrix `yaml:"matrix,omitempty"`
}

type Matrix struct {
	Matrix  map[string][]string      `yaml:",inline"`
	Include []map[string]interface{} `yaml:"include,omitempty"`
	Exclude []map[string]interface{} `yaml:"exclude,omitempty"`
}

type Service struct {
	Image string `yaml:"image,omitempty"`
}

type WorkflowTriggers struct {
	Push                     *PushCondition                     `yaml:"push,omitempty"`
	PullRequest              *PullRequestCondition              `yaml:"pull_request,omitempty"`
	WorkflowDispatch         *WorkflowDispatchCondition         `yaml:"workflow_dispatch,omitempty"`
	Schedule                 []*ScheduleCondition               `yaml:"schedule,omitempty"`
	RepositoryDispatch       *RepositoryDispatchCondition       `yaml:"repository_dispatch,omitempty"`
	IssueComment             *IssueCommentCondition             `yaml:"issue_comment,omitempty"`
	Issues                   *IssuesCondition                   `yaml:"issues,omitempty"`
	PullRequestReview        *PullRequestReviewCondition        `yaml:"pull_request_review,omitempty"`
	PullRequestReviewComment *PullRequestReviewCommentCondition `yaml:"pull_request_review_comment,omitempty"`
	Label                    *LabelCondition                    `yaml:"label,omitempty"`
	Release                  *ReleaseCondition                  `yaml:"release,omitempty"`
}

type PushCondition struct {
	Branches []string `yaml:"branches,omitempty"`
	Tags     []string `yaml:"tags,omitempty"`
	Paths    []string `yaml:"paths,omitempty"`
}

type PullRequestCondition struct {
	Branches []string `yaml:"branches,omitempty"`
	Tags     []string `yaml:"tags,omitempty"`
	Paths    []string `yaml:"paths,omitempty"`
	Types    []string `yaml:"types,omitempty"`
}

type WorkflowDispatchCondition struct {
	Inputs map[string]InputDefinition `yaml:"inputs,omitempty"`
}

type InputDefinition struct {
	Description string `yaml:"description,omitempty"`
	Default     string `yaml:"default,omitempty"`
	Required    bool   `yaml:"required,omitempty"`
}

type ScheduleCondition struct {
	Cron string `yaml:"cron,omitempty"`
}

type RepositoryDispatchCondition struct {
	Types []string `yaml:"types,omitempty"`
}

type IssueCommentCondition struct {
	Types []string `yaml:"types,omitempty"`
}

type IssuesCondition struct {
	Types []string `yaml:"types,omitempty"`
}

type PullRequestReviewCondition struct {
	Types []string `yaml:"types,omitempty"`
}

type PullRequestReviewCommentCondition struct {
	Types []string `yaml:"types,omitempty"`
}

type LabelCondition struct {
	Types []string `yaml:"types,omitempty"`
}

type ReleaseCondition struct {
	Types []string `yaml:"types,omitempty"`
}
