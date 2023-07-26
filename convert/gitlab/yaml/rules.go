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

type Rule struct {
	If            string            `yaml:"if,omitempty"`
	Changes       Change            `yaml:"changes,omitempty"`
	Exists        []string          `yaml:"exists,omitempty"`
	AllowFailures bool              `yaml:"allow_failure,omitempty"`
	Variables     map[string]string `yaml:"variables,omitempty"`
	When          string            `yaml:"when,omitempty"`
	Needs         []string          `yaml:"needs,omitempty"`
}

type Change struct {
	Paths     []string `yaml:"paths"`
	CompareTo string   `yaml:"compare_to,omitempty"`
}

func (v *Change) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 []string
	var out2 string
	var out3 = struct {
		Paths     []string `yaml:"paths"`
		CompareTo string   `yaml:"compare_to,omitempty"`
	}{}

	if err := unmarshal(&out1); err == nil {
		v.Paths = out1
		return nil
	}

	if err := unmarshal(&out2); err == nil {
		v.Paths = []string{out2}
		return nil
	}

	if err := unmarshal(&out3); err == nil {
		v.Paths = out3.Paths
		v.CompareTo = out3.CompareTo
		return nil
	}

	return errors.New("failed to unmarshal rules:changes")
}

func (v *Change) MarshalYAML() (interface{}, error) {
	if v.CompareTo != "" {
		return map[string]interface{}{
			"paths":      v.Paths,
			"compare_to": v.CompareTo,
		}, nil
	}
	return v.Paths, nil
}
