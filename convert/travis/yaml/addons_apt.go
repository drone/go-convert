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

// the alias details come from here:
// https://github.com/travis-ci/apt-source-safelist/blob/master/ubuntu.json

type Apt struct {
	Enabled  bool         `yaml:"enabled,omitempty"`
	Packages []string     `yaml:"packages,omitempty"`
	Sources  []*AptSource `yaml:"sources,omitempty"`
	Dist     string       `yaml:"dist,omitempty"`
	Update   bool         `yaml:"bool,omitempty"`
}

// UnmarshalYAML implements the unmarshal interface.
func (v *Apt) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 bool
	var out2 string
	var out3 []string
	var out4 = struct {
		Enabled  *bool         `yaml:"enabled"`
		Packages Stringorslice `yaml:"packages,omitempty"`
		Package  Stringorslice `yaml:"package,omitempty"` // alias
		Sources  AptSources    `yaml:"sources,omitempty"`
		Source   AptSources    `yaml:"source,omitempty"` // alias
		Dist     string        `yaml:"dist,omitempty"`
		Update   bool          `yaml:"bool,omitempty"`
	}{}
	if err := unmarshal(&out1); err == nil {
		v.Enabled = out1
		return nil
	}
	if err := unmarshal(&out2); err == nil {
		v.Enabled = true
		v.Packages = append(v.Packages, out2)
		return nil
	}
	if err := unmarshal(&out3); err == nil {
		v.Enabled = true
		v.Packages = append(v.Packages, out3...)
		return nil
	}
	if err := unmarshal(&out4); err == nil {
		v.Enabled = true
		if out4.Enabled != nil {
			v.Enabled = *out4.Enabled
		}
		v.Packages = out4.Packages
		v.Sources = append(v.Sources, out4.Sources.Items...)
		v.Dist = out4.Dist
		v.Update = out4.Update
		v.Packages = append(v.Packages, out4.Package...)
		v.Sources = append(v.Sources, out4.Source.Items...)
		return nil
	}
	return errors.New("failed to unmarshal apt")
}

type AptSources struct {
	Items []*AptSource
}

// UnmarshalYAML implements the unmarshal interface.
func (v *AptSources) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 string
	var out2 *AptSource
	var out3 []*AptSource
	if err := unmarshal(&out1); err == nil {
		v.Items = append(v.Items, &AptSource{Alias: out1})
		return nil
	}
	if err := unmarshal(&out2); err == nil {
		v.Items = append(v.Items, out2)
		return nil
	}
	if err := unmarshal(&out3); err == nil {
		v.Items = append(v.Items, out3...)
		return nil
	}
	return errors.New("failed to unmarshal apt source list")
}

type AptSource struct {
	Alias           string `yaml:"alias,omitempty"`
	Sourceline      string `yaml:"sourceline,omitempty"`
	KeyURL          string `yaml:"key_url,omitempty"`
	CanonicalKeyURL string `yaml:"canonical_key_url,omitempty"`
}

// UnmarshalYAML implements the unmarshal interface.
func (v *AptSource) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 string
	var out2 = struct {
		Alias           string `yaml:"alias,omitempty"`
		Sourceline      string `yaml:"sourceline,omitempty"`
		KeyURL          string `yaml:"key_url,omitempty"`
		CanonicalKeyURL string `yaml:"canonical_key_url,omitempty"`
	}{}
	if err := unmarshal(&out1); err == nil {
		v.Alias = out1
		return nil
	}
	if err := unmarshal(&out2); err == nil {
		v.Alias = out2.Alias
		v.Sourceline = out2.Sourceline
		v.KeyURL = out2.KeyURL
		v.CanonicalKeyURL = out2.CanonicalKeyURL
		return nil
	}
	return errors.New("failed to unmarshal apt source")
}
