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

package downgrader

// Option configures a Downgrade option.
type Option func(*Downgrader)

// WithCodebase returns an option to set the codebase.
func WithCodebase(name, connector string) Option {
	return func(d *Downgrader) {
		d.codebaseName = name
		d.codebaseConn = connector
	}
}

// WithDockerhub returns an option to set the default
// dockerhub registry connector.
func WithDockerhub(connector string) Option {
	return func(d *Downgrader) {
		d.dockerhubConn = connector
	}
}

// WithKubernetes returns an option to set the default
// runtime to Kubernetes.
func WithKubernetes(namespace, connector string) Option {
	return func(d *Downgrader) {
		d.kubeNamespace = namespace
		d.kubeConnector = connector
	}
}

// WithIdentifier returns an option to set the pipeline
// identifier.
func WithIdentifier(identifier string) Option {
	return func(d *Downgrader) {
		d.pipelineId = identifier
	}
}

// WithName returns an option to set the pipeline name.
func WithName(name string) Option {
	return func(d *Downgrader) {
		d.pipelineName = name
	}
}

// WithOrganization returns an option to set the
// harness organization.
func WithOrganization(organization string) Option {
	return func(d *Downgrader) {
		d.pipelineOrg = organization
	}
}

// WithProject returns an option to set the harness
// project name..
func WithProject(project string) Option {
	return func(d *Downgrader) {
		d.pipelineProj = project
	}
}
