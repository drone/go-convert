package yaml

type Permissions struct {
	Actions            string `yaml:"actions,omitempty"`
	Checks             string `yaml:"checks,omitempty"`
	Contents           string `yaml:"contents,omitempty"`
	Deployments        string `yaml:"deployments,omitempty"`
	IDToken            string `yaml:"id-token,omitempty"`
	Issues             string `yaml:"issues,omitempty"`
	Discussions        string `yaml:"discussions,omitempty"`
	Packages           string `yaml:"packages,omitempty"`
	Pages              string `yaml:"pages,omitempty"`
	PullRequests       string `yaml:"pull-requests,omitempty"`
	RepositoryProjects string `yaml:"repository-projects,omitempty"`
	SecurityEvents     string `yaml:"security-events,omitempty"`
	Statuses           string `yaml:"statuses,omitempty"`
}
