// Copyright 2023 Harness, Inc.
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

package jenkinsjson

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

// WithInfrastructure returns an option to set the infrastructure type.
func WithInfrastructure(infra string) Option {
	return func(d *Converter) {
		d.infrastructure = infra
	}
}

// WithOS returns an option to set the operating system.
func WithOS(os string) Option {
	return func(d *Converter) {
		d.os = os
	}
}

// WithArch returns an option to set the CPU architecture.
func WithArch(arch string) Option {
	return func(d *Converter) {
		d.arch = arch
	}
}

func WithUseIntelligence(useIntelligence bool) Option {
	return func(d *Converter) {
		d.useIntelligence = useIntelligence
	}
}

func WithConfigFile(configFile string) Option {
	return func(d *Converter) {
		d.configFile = configFile
	}
}
