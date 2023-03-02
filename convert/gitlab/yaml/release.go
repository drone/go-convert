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

type Release struct {
	Tag        string        `yaml:"tag_name,omitempty"`
	TagMessage string        `yaml:"tag_message,omitempty"`
	Name       string        `yaml:"name,omitempty"`
	Desc       string        `yaml:"description,omitempty"`
	Ref        string        `yaml:"ref,omitempty"`
	Milestones Stringorslice `yaml:"milestones,omitempty"`
	ReleasedAt string        `yaml:"released_at,omitempty"`
	Assets     *Assets       `yaml:"assets,omitempty"`
}

type Assets struct {
	Links []*AssetLink `yaml:"links,omitempty"`
}

type AssetLink struct {
	Name     string `yaml:"name,omitempty"`
	Url      string `yaml:"url,omitempty"`
	FilePath string `yaml:"filepath,omitempty"`
	LinkType string `yaml:"link_type,omitempty"`
}
