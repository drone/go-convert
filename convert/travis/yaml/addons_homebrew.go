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

type Homebrew struct {
	Update   bool     `yaml:"update,omitempty"`
	Packages []string `yaml:"packages,omitempty"`
	Casks    []string `yaml:"casks,omitempty"`
	Taps     []string `yaml:"taps,omitempty"`
	Brewfile string   `yaml:"brewfile,omitempty"` // string or bool
}

// UnmarshalYAML implements the unmarshal interface.
func (v *Homebrew) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 bool
	var out2 string
	var out3 []string
	var out4 = struct {
		Update   bool          `yaml:"update,omitempty"`
		Packages Stringorslice `yaml:"packages,omitempty"`
		Casks    Stringorslice `yaml:"casks,omitempty"`
		Taps     Stringorslice `yaml:"taps,omitempty"`
		Brewfile interface{}   `yaml:"brewfile,omitempty"` // string or bool
	}{}

	if err := unmarshal(&out1); err == nil {
		v.Update = out1
		return nil
	}
	if err := unmarshal(&out2); err == nil {
		v.Update = true
		v.Packages = []string{out2}
		return nil
	}
	if err := unmarshal(&out3); err == nil {
		v.Update = true
		v.Packages = out3
		return nil
	}
	if err := unmarshal(&out4); err == nil {
		v.Update = out4.Update
		v.Packages = out4.Packages
		v.Casks = out4.Casks
		v.Taps = out4.Taps
		switch vv := out4.Brewfile.(type) {
		case bool:
			v.Brewfile = "Brewfile"
		case string:
			v.Brewfile = vv
		}
		return nil
	}

	return errors.New("failed to unmarshal homebrew")
}
