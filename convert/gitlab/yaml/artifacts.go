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

type Artifacts struct {
	Paths     Stringorslice          `yaml:"paths,omitempty"`
	Exclude   Stringorslice          `yaml:"exclude,omitempty"`
	ExpireIn  string                 `yaml:"expire_in,omitempty"`
	ExposeAs  string                 `yaml:"expose_as,omitempty"`
	Name      string                 `yaml:"name,omitempty"`
	Public    *bool                  `yaml:"public,omitempty"`
	Reports   map[string]interface{} `yaml:"reports,omitempty"`
	Untracked bool                   `yaml:"untracked,omitempty"`
	When      string                 `yaml:"when,omitempty"` // on_success, on_failure, always
}
