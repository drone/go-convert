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

type Codeclimate struct {
	Enabled   bool    `yaml:"enabled,omitempty"`
	RepoToken *Secure `yaml:"repo_token,omitempty"`
}

// UnmarshalYAML implements the unmarshal interface.
func (v *Codeclimate) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 bool
	var out2 *Secure
	var out3 = struct {
		Enabled   *bool   `yaml:"enabled"`
		RepoToken *Secure `yaml:"repo_token,omitempty"`
	}{}
	if err := unmarshal(&out1); err == nil {
		v.Enabled = out1
		return nil
	}
	if err := unmarshal(&out2); err == nil && out2 != nil && out2.Encrypted != "" {
		v.Enabled = true
		v.RepoToken = out2
		return nil
	}
	if err := unmarshal(&out3); err == nil {
		v.Enabled = true
		if out3.Enabled != nil {
			v.Enabled = *out3.Enabled
		}
		v.RepoToken = out3.RepoToken
		return nil
	}
	return errors.New("failed to unmarshal codeclimate")
}
