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

type (
	Artifacts struct {
		Download *bool
		Paths    []string
	}

	// temporary data structure for unmarshaling and
	// marshaling artifacts.
	artifacts struct {
		Download *bool    `yaml:"download,omitempty"`
		Paths    []string `yaml:"paths,omitempty"`
	}
)

// MarshalYAML implements the marshal interface.
func (v *Artifacts) MarshalYAML() (interface{}, error) {
	if len(v.Paths) == 0 && v.Download == nil {
		return nil, nil
	} else if v.Download == nil {
		// emit the short syntax when the download
		// value is true (default)
		return v.Paths, nil
	} else {
		// emit the short syntax when the download
		// value is false (non-default)
		return &artifacts{
			Download: v.Download,
			Paths:    v.Paths,
		}, nil
	}
}

// UnmarshalYAML implements the unmarshal interface.
func (v *Artifacts) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 string
	var out2 []string
	var out3 *artifacts
	if err := unmarshal(&out1); err == nil {
		v.Paths = append(v.Paths, out1)
		return nil
	}
	if err := unmarshal(&out2); err == nil {
		v.Paths = append(v.Paths, out2...)
		return nil
	}
	if err := unmarshal(&out3); err == nil {
		v.Download = out3.Download
		v.Paths = append(v.Paths, out3.Paths...)
		return nil
	}
	return errors.New("failed to unmarshal artifacts")
}
