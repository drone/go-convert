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

type Coverity struct {
	Enabled             bool             `yaml:"enabled,omitempty"`
	Project             *CoverityProject `yaml:"project,omitempty"`
	BuildScriptURL      string           `yaml:"build_script_url,omitempty"`
	BranchPattern       string           `yaml:"branch_pattern,omitempty"`
	NotificationEmail   *Secure          `yaml:"notification_email,omitempty"`
	BuildCommand        string           `yaml:"build_command,omitempty"`
	BuildCommandPrepend string           `yaml:"build_command_prepend,omitempty"`
}

// UnmarshalYAML implements the unmarshal interface.
func (v *Coverity) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 bool
	var out2 = struct {
		Enabled             *bool            `yaml:"enabled"`
		Project             *CoverityProject `yaml:"project,omitempty"`
		BuildScriptURL      string           `yaml:"build_script_url,omitempty"`
		BranchPattern       string           `yaml:"branch_pattern,omitempty"`
		NotificationEmail   *Secure          `yaml:"notification_email,omitempty"`
		BuildCommand        string           `yaml:"build_command,omitempty"`
		BuildCommandPrepend string           `yaml:"build_command_prepend,omitempty"`
	}{}
	if err := unmarshal(&out1); err == nil {
		v.Enabled = out1
		return nil
	}
	if err := unmarshal(&out2); err == nil {
		v.Enabled = true
		if out2.Enabled != nil {
			v.Enabled = *out2.Enabled
		}
		v.Project = out2.Project
		v.BuildScriptURL = out2.BuildScriptURL
		v.BranchPattern = out2.BranchPattern
		v.NotificationEmail = out2.NotificationEmail
		v.BuildCommand = out2.BuildCommand
		v.BuildCommandPrepend = out2.BuildCommandPrepend
		return nil
	}
	return errors.New("failed to unmarshal coverity_scan")
}

type CoverityProject struct {
	Name        string `yaml:"name,omitempty"`
	Version     string `yaml:"version,omitempty"`
	Description string `yaml:"description,omitempty"`
}

// UnmarshalYAML implements the unmarshal interface.
func (v *CoverityProject) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 string
	var out2 = struct {
		Name        string `yaml:"name,omitempty"`
		Version     string `yaml:"version,omitempty"`
		Description string `yaml:"description,omitempty"`
	}{}
	if err := unmarshal(&out1); err == nil {
		v.Name = out1
		return nil
	}
	if err := unmarshal(&out2); err == nil {
		v.Name = out2.Name
		v.Version = out2.Version
		v.Description = out2.Description
		return nil
	}
	return errors.New("failed to unmarshal coverity project")
}
