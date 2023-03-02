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

type Jobs struct {
	Include       []map[string]string `yaml:"include,omitempty"`
	Exclude       []map[string]string `yaml:"exclude,omitempty"`
	AllowFailures []map[string]string `yaml:"allow_failures,omitempty"`
	FastFinish    bool                `yaml:"fast_finish,omitempty"`
}

// UnmarshalYAML implements the unmarshal interface.
func (v *Jobs) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 map[string]string
	var out2 []map[string]string
	var out3 = struct {
		Include            []map[string]string `yaml:"include,omitempty"`
		Exclude            []map[string]string `yaml:"exclude,omitempty"`
		AllowFailures      []map[string]string `yaml:"allow_failures,omitempty"`
		AllowFailuresAlias []map[string]string `yaml:"allowed_failures,omitempty"`
		FastFinish         bool                `yaml:"fast_finish,omitempty"`
		FastFinishAlias    bool                `yaml:"fast_failure,omitempty"`
	}{}
	if err := unmarshal(&out1); err == nil {
		v.Include = append(v.Include, out1)
		return nil
	}
	if err := unmarshal(&out2); err == nil {
		v.Include = out2
		return nil
	}
	if err := unmarshal(&out3); err == nil {
		v.Include = out3.Include
		v.Exclude = out3.Exclude
		v.AllowFailures = out3.AllowFailures
		v.FastFinish = out3.FastFinish
		// map allowed_failures alias to allow_failures
		v.AllowFailures = append(v.AllowFailures, out3.AllowFailuresAlias...)
		// map fast_failure alias to fast_finish
		if out3.FastFinishAlias {
			v.FastFinish = true
		}
		return nil
	}
	return errors.New("failed to unmarshal jobs")
}
