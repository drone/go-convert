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

type Git struct {
	Strategy        string `yaml:"strategy,omitempty"` // enum: clone, tarball
	Depth           *Depth `yaml:"depth,omitempty"`
	Quiet           bool   `yaml:"quiet,omitempty"`
	Submodules      bool   `yaml:"submodules,omitempty"`
	SubmodulesDepth int    `yaml:"submodules_depth,omitempty"`
	LFSSkipSmudge   bool   `yaml:"lfs_skip_smudge,omitempty"`
	SparseCheckout  string `yaml:"sparse_checkout,omitempty"`
	Autocrlf        string `yaml:"autocrlf,omitempty"`
}

// Depth sets the clone depth.
type Depth struct {
	Value int
}

// UnmarshalYAML implements the unmarshal interface.
func (v *Depth) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 int
	var out2 bool
	if err := unmarshal(&out1); err == nil {
		v.Value = out1
		return nil
	}
	if err := unmarshal(&out2); err == nil {
		if out2 == true {
			v.Value = 50 // default depth for travis pipelines
		}
		return nil
	}
	return errors.New("failed to unmarshal depth")
}
