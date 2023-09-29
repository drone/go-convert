package yaml

type (
	// Config defines a rio pipeline.
	Config struct {
		SchemaVersion float32    `yaml:"schemaVersion"`
		Timeout       int        `yaml:"timeout,omitempty"`
		Pipelines     []Pipeline `yaml:"pipelines"`
		Notify        Notify     `yaml:"notify,omitempty"`
	}

	Pipeline struct {
		Name       string `yaml:"name,omitempty"`
		BranchName string `yaml:"branchName,omitempty"`
		Machine    struct {
			BaseImage       string            `yaml:"baseImage,omitempty"`
			Env             map[string]string `yaml:"env,omitempty"`
			TargetPlatforms []string          `yaml:"targetPlatforms,omitempty"`
		} `yaml:"machine,omitempty"`
		Build   Build   `yaml:"build,omitempty"`
		Reports Reports `yaml:"reports,omitempty"`
		Pkg     Pkg     `yaml:"package,omitempty"`
		Finally Finally `yaml:"finally,omitempty"`
		Trigger struct {
			GitPush bool `yaml:"gitPush,omitempty"`
		} `yaml:"trigger,omitempty"`
		Checkout struct {
			FetchTags bool `yaml:"fetchTags,omitempty"`
		} `yaml:"checkout,omitempty"`
	}

	Build struct {
		Template string   `yaml:"template,omitempty"`
		Steps    []string `yaml:"steps,omitempty"`
	}

	Reports struct {
		Findbugs bool              `yaml:"findbugs,omitempty"`
		Jacoco   map[string]string `yaml:"jacoco,omitempty"`
	}

	Pkg struct {
		Release    bool         `yaml:"release,omitempty"`
		Dockerfile []Dockerfile `yaml:"dockerfile,omitempty"`
	}

	Dockerfile struct {
		Context        string            `yaml:"context,omitempty"`
		Version        string            `yaml:"version,omitempty"`
		PerApplication bool              `yaml:"perApplication,omitempty"`
		DockerfilePath string            `yaml:"dockerfilePath,omitempty"`
		Env            map[string]string `yaml:"env,omitempty"`
		Publish        []Publish         `yaml:"publish,omitempty"`
		ExtraTags      []string          `yaml:"extraTags,omitempty"`
	}

	Publish struct {
		Repo string `yaml:"repo,omitempty"`
	}

	Finally struct {
		Tag struct {
			Enabled bool `yaml:"enabled,omitempty"`
		} `yaml:"tag,omitempty"`
	}

	Notify struct {
		Email struct {
			Enabled bool `yaml:"enabled,omitempty"`
		} `yaml:"email,omitempty"`
		PullRequestComment struct {
			PostOnSuccess bool `yaml:"postOnSuccess,omitempty"`
			PostOnFailure bool `yaml:"postOnFailure,omitempty"`
		} `yaml:"pullRequestComment,omitempty"`
	}
)
