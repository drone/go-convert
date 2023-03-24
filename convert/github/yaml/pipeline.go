package yaml

type Pipeline struct {
	Concurrency *Concurrency `yaml:"concurrency"`
	Name        string       `yaml:"name,omitempty"`
	//On          *WorkflowTriggers `yaml:"on,omitempty"`
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
	If          string              `yaml:"if,omitempty"`
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
	Image    string            `yaml:"image,omitempty"`
	Env      map[string]string `yaml:"env,omitempty"`
	Ports    []string          `yaml:"ports,omitempty"`
	Options  []string          `yaml:"options,omitempty"`
	Volumes  []string          `yaml:"volumes,omitempty"`
	Networks []string          `yaml:"networks,omitempty"`
}
