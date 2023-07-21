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

type AllowFailure struct {
	Value     bool  `yaml:"-"`
	ExitCodes []int `yaml:"exit_codes"`
}

// UnmarshalYAML implements the unmarshal interface.
func (v *AllowFailure) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 bool
	var out2 = struct {
		ExitCode int `yaml:"exit_codes"`
	}{}
	var out3 = struct {
		ExitCodes []int `yaml:"exit_codes"`
	}{}

	if err := unmarshal(&out3); err == nil {
		v.Value = true
		v.ExitCodes = out3.ExitCodes
		return nil
	}

	if err := unmarshal(&out2); err == nil {
		v.Value = true
		v.ExitCodes = []int{out2.ExitCode}
		return nil
	}

	if err := unmarshal(&out1); err == nil {
		v.Value = out1
		return nil
	}

	return errors.New("failed to unmarshal allow_failure")
}

// MarshalYAML implements the Marshal interface.
func (v *AllowFailure) MarshalYAML() (interface{}, error) {
	type allowFailureMarshal AllowFailure
	if !v.Value {
		return nil, nil
	}

	if len(v.ExitCodes) == 0 {
		return v.Value, nil
	}

	return allowFailureMarshal(*v), nil
}
