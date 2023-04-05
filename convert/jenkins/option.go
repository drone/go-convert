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

package jenkins

import "strings"

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

// WithAttempts returns an option to set the number
// of retry attempts.
func WithAttempts(attempts int) Option {
	return func(d *Converter) {
		d.attempts = attempts
	}
}

// WithToken returns an option to set the API token
// for Chat GPT.
func WithToken(token string) Option {
	return func(d *Converter) {
		d.token = token
	}
}

// WithDebug returns an option to use debug mode.
func WithDebug() Option {
	return func(d *Converter) {
		d.debug = true
	}
}

// WithFormat returns an option to customize
// the intermediate format.
func WithFormat(format Format) Option {
	return func(d *Converter) {
		d.format = format
	}
}

// WithFormat returns an option to customize
// the intermediate format.
func WithFormatString(format string) Option {
	return func(d *Converter) {
		format = strings.TrimSpace(format)
		format = strings.ToLower(format)
		switch format {
		case "github":
			d.format = FromGithub
		case "gitlab":
			d.format = FromGitlab
		case "drone":
			d.format = FromDrone
		}
	}
}
