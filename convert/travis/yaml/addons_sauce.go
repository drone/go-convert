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

type Sauce struct {
	Enabled          bool    `yaml:"enabled,omitempty"`
	Username         *Secure `yaml:"username,omitempty"`
	AccessKey        *Secure `yaml:"access_key,omitempty"`
	DirectDomains    string  `yaml:"direct_domains,omitempty"`
	TunnelDomains    string  `yaml:"tunnel_domains,omitempty"`
	NoSSLBumpDomains string  `yaml:"no_ssl_bump_domains,omitempty"`
}

// UnmarshalYAML implements the unmarshal interface.
func (v *Sauce) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 bool
	var out2 = struct {
		Enabled          *bool   `yaml:"enabled"`
		Username         *Secure `yaml:"username"`
		AccessKey        *Secure `yaml:"access_key"`
		DirectDomains    string  `yaml:"direct_domains"`
		TunnelDomains    string  `yaml:"tunnel_domains"`
		NoSSLBumpDomains string  `yaml:"no_ssl_bump_domains"`
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
		v.Username = out2.Username
		v.AccessKey = out2.AccessKey
		v.DirectDomains = out2.DirectDomains
		v.TunnelDomains = out2.TunnelDomains
		v.NoSSLBumpDomains = out2.NoSSLBumpDomains
		return nil
	}
	return errors.New("failed to unmarshal sauce_connect")
}
