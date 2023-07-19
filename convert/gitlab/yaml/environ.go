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

type Environment struct {
	Name           string      `yaml:"name,omitempty"`
	Url            string      `yaml:"url,omitempty"`
	OnStop         string      `yaml:"on_stop,omitempty"`
	Action         string      `yaml:"action,omitempty"` // start, prepare, stop, verify, access
	AutoStopIn     string      `yaml:"auto_stop_in,omitempty"`
	DeploymentTier string      `yaml:"deployment_tier,omitempty"` // production, staging, testing, development, other
	Kubernetes     *Kubernetes `yaml:"kubernetes,omitempty"`
}

type Kubernetes struct {
	Namespace string `yaml:"namespace,omitempty"`
}

// UnmarshalYAML implements the unmarshal interface.
func (v *Environment) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 string
	var out2 = struct {
		Name           string      `yaml:"name,omitempty"`
		Url            string      `yaml:"url,omitempty"`
		OnStop         string      `yaml:"on_stop,omitempty"`
		Action         string      `yaml:"action,omitempty"`
		AutoStopIn     string      `yaml:"auto_stop_in,omitempty"`
		DeploymentTier string      `yaml:"deployment_tier,omitempty"`
		Kubernetes     *Kubernetes `yaml:"kubernetes,omitempty"`
	}{}

	if err := unmarshal(&out1); err == nil {
		v.Name = out1
		return nil
	}

	if err := unmarshal(&out2); err == nil {
		v.Name = out2.Name
		v.Url = out2.Url
		v.OnStop = out2.OnStop
		v.Action = out2.Action
		v.AutoStopIn = out2.AutoStopIn
		v.DeploymentTier = out2.DeploymentTier
		v.Kubernetes = out2.Kubernetes
		return nil
	}

	return errors.New("failed to unmarshal environment")
}
