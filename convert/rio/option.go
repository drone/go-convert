package rio

// Option configures a Converter option.
type Option func(*Converter)

// WithDockerhub returns an option to set the default
// dockerhub registry connector.
func WithDockerhub(connector string) Option {
	return func(d *Converter) {
		d.dockerhubConn = connector
	}
}

// WithKubernetes returns an option to set the default
// runtime to Kubernetes.
func WithKubernetes(namespace, connector string) Option {
	return func(d *Converter) {
		d.kubeNamespace = namespace
		d.kubeConnector = connector
	}
}

// WithIdentifier returns an option to set the pipeline
// identifier.
func WithIdentifier(identifier string) Option {
	return func(d *Converter) {
		d.pipelineId = identifier
	}
}

// WithName returns an option to set the pipeline name.
func WithName(name string) Option {
	return func(d *Converter) {
		d.pipelineName = name
	}
}

// WithOrganization returns an option to set the
// harness organization.
func WithOrganization(organization string) Option {
	return func(d *Converter) {
		d.pipelineOrg = organization
	}
}

// WithProject returns an option to set the harness
// project name.
func WithProject(project string) Option {
	return func(d *Converter) {
		d.pipelineProj = project
	}
}

// WithNotifyUserGroup returns an option to set the notification email.
func WithNotifyUserGroup(notifyUserGroup string) Option {
	return func(d *Converter) {
		d.notifyUserGroup = notifyUserGroup
	}
}

// WithGithubConnector returns an option to set the github connector.
func WithGithubConnector(githubConnector string) Option {
	return func(d *Converter) {
		d.githubConnector = githubConnector
	}
}
