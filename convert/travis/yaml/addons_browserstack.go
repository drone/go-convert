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

type Browserstack struct {
	Enabled    bool    `yaml:"enabled,omitempty"`
	Username   *Secure `yaml:"username,omitempty"`
	AccessKey  *Secure `yaml:"access_key,omitempty"`
	ForceLocal bool    `yaml:"forcelocal,omitempty"`
	Only       string  `yaml:"only,omitempty"`
	AppPath    string  `yaml:"app_path,omitempty"`
	ProxyHost  string  `yaml:"proxyHost,omitempty"`
	ProxyPort  string  `yaml:"proxyPort,omitempty"`
	ProxyUser  string  `yaml:"proxyUser,omitempty"`
	ProxyPass  *Secure `yaml:"proxyPass,omitempty"`
}

// UnmarshalYAML implements the unmarshal interface.
func (v *Browserstack) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 bool
	var out2 = struct {
		Enabled    *bool   `yaml:"enabled"`
		Username   *Secure `yaml:"username,omitempty"`
		AccessKey  *Secure `yaml:"access_key,omitempty"`
		ForceLocal bool    `yaml:"forcelocal,omitempty"`
		Only       string  `yaml:"only,omitempty"`
		AppPath    string  `yaml:"app_path,omitempty"`
		ProxyHost  string  `yaml:"proxyHost,omitempty"`
		ProxyPort  string  `yaml:"proxyPort,omitempty"`
		ProxyUser  string  `yaml:"proxyUser,omitempty"`
		ProxyPass  *Secure `yaml:"proxyPass,omitempty"`
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
		v.ForceLocal = out2.ForceLocal
		v.Only = out2.Only
		v.AppPath = out2.AppPath
		v.ProxyHost = out2.ProxyHost
		v.ProxyPort = out2.ProxyPort
		v.ProxyUser = out2.ProxyUser
		v.ProxyPass = out2.ProxyPass
		return nil
	}
	return errors.New("failed to unmarshal browserstack")
}
