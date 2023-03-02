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

type Cache struct {
	Directories []string `yaml:"directories,omitempty"`
	Apt         bool     `yaml:"apt,omitempty"`
	Bundler     bool     `yaml:"bundler,omitempty"`
	Cargo       bool     `yaml:"cargo,omitempty"`
	Ccache      bool     `yaml:"ccache,omitempty"`
	Cocoapods   bool     `yaml:"cocoapods,omitempty"`
	Npm         bool     `yaml:"npm,omitempty"`
	Packages    bool     `yaml:"packages,omitempty"`
	Pip         bool     `yaml:"pip,omitempty"`
	Yarn        bool     `yaml:"yarn,omitempty"`
	Edge        bool     `yaml:"edge,omitempty"`
	Branch      string   `yaml:"branch,omitempty"`
	Timeout     int      `yaml:"timeout,omitempty"`
}

// UnmarshalYAML implements the unmarshal interface.
func (v *Cache) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 bool
	var out2 = struct {
		Directories Stringorslice `yaml:"directories"`
		Apt         bool          `yaml:"apt"`
		Bundler     bool          `yaml:"bundler"`
		Cargo       bool          `yaml:"cargo"`
		Ccache      bool          `yaml:"ccache"`
		Cocoapods   bool          `yaml:"cocoapods"`
		Npm         bool          `yaml:"npm"`
		Packages    bool          `yaml:"packages"`
		Pip         bool          `yaml:"pip"`
		Yarn        bool          `yaml:"yarn"`
		Edge        bool          `yaml:"edge"`
		Branch      string        `yaml:"branch"`
		Timeout     int           `yaml:"timeout"`
	}{}

	if err := unmarshal(&out1); err == nil {
		v.Timeout = 3
		return nil
	}

	if err := unmarshal(&out2); err == nil {
		v.Directories = out2.Directories
		v.Apt = out2.Apt
		v.Bundler = out2.Bundler
		v.Cargo = out2.Cargo
		v.Ccache = out2.Ccache
		v.Cocoapods = out2.Cocoapods
		v.Npm = out2.Npm
		v.Packages = out2.Packages
		v.Pip = out2.Pip
		v.Edge = out2.Edge
		v.Branch = out2.Branch
		v.Timeout = out2.Timeout
		if v.Timeout == 0 {
			v.Timeout = 3
		}
		return nil
	}

	return errors.New("failed to unmarshal cache")
}
