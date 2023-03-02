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

import (
	"errors"
)

type Cache struct {
	Paths     Stringorslice `yaml:"paths,omitempty"`
	Key       string        `yaml:"key,omitempty"` // TODO complex item (files []string, prefix)
	Untracked bool          `yaml:"untracked,omitempty"`
	Unprotect bool          `yaml:"unprotect,omitempty"`
	When      string        `yaml:"when,omitempty"`   // on_success, on_failure, always
	Policy    string        `yaml:"policy,omitempty"` // pull, push, pull-push
}

type CacheKey struct {
	Value  string        `yaml:"-"`
	Files  Stringorslice `yaml:"files,omitempty"`
	Prefix string        `yaml:"prefix,omitempty"`
}

// UnmarshalYAML implements the unmarshal interface.
func (v *CacheKey) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 string
	var out2 = struct {
		Files  Stringorslice `yaml:"files"`
		Prefix string        `yaml:"prefix"`
	}{}

	if err := unmarshal(&out1); err == nil {
		v.Value = out1
		return nil
	}

	if err := unmarshal(&out2); err == nil {
		v.Files = out2.Files
		v.Prefix = out2.Prefix
		return nil
	}

	return errors.New("failed to unmarshal cache key")
}
