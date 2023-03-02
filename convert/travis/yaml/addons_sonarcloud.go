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

type Sonarcloud struct {
	Enabled      bool     `yaml:"enabled,omitempty"`
	Organization string   `yaml:"organization,omitempty"`
	Token        *Secure  `yaml:"token,omitempty"`
	GithubToken  *Secure  `yaml:"github_token,omitempty"`
	Branches     []string `yaml:"branches,omitempty"`
}

// UnmarshalYAML implements the unmarshal interface.
func (v *Sonarcloud) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 bool
	var out2 = struct {
		Enabled      *bool         `yaml:"enabled"`
		Organization string        `yaml:"organization"`
		Token        *Secure       `yaml:"token"`
		GithubToken  *Secure       `yaml:"github_token"`
		Branches     Stringorslice `yaml:"branches"`
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
		v.Organization = out2.Organization
		v.Token = out2.Token
		v.GithubToken = out2.GithubToken
		v.Branches = out2.Branches
		return nil
	}
	return errors.New("failed to unmarshal sonarcloud")
}
