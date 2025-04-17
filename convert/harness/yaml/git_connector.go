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

type (
	GitConnector struct {
		ID      string            `json:"identifier,omitempty"        yaml:"identifier,omitempty"`
		Name    string            `json:"name,omitempty"              yaml:"name,omitempty"`
		Desc    string            `json:"description,omitempty"       yaml:"description,omitempty"`
		Account string            `json:"accountIdentifier,omitempty" yaml:"accountIdentifier,omitempty"`
		Org     string            `json:"orgIdentifier,omitempty"     yaml:"orgIdentifier,omitempty"`
		Project string            `json:"projectIdentifier,omitempty" yaml:"projectIdentifier,omitempty"`
		Type    string            `json:"type,omitempty"              yaml:"type,omitempty"`
		Spec    *GitConnectorSpec `json:"spec,omitempty"              yaml:"spec,omitempty"`
		Tags    map[string]string `json:"tags,omitempty"              yaml:"tags,omitempty"`
	}

	GitConnectorSpec struct {
		Url                  string             `json:"url,omitempty"                  yaml:"url,omitempty"`
		ValidationRepo       string             `json:"validationRepo,omitempty"       yaml:"validationRepo,omitempty"`
		ExecuteOnDelegate    bool               `json:"executeOnDelegate,omitempty"    yaml:"executeOnDelegate,omitempty"`
		Proxy                bool               `json:"proxy,omitempty"                yaml:"proxy,omitempty"`
		IgnoreTestConnection bool               `json:"ignoreTestConnection,omitempty" yaml:"ignoreTestConnection,omitempty"`
		Type                 string             `json:"type,omitempty"                 yaml:"type,omitempty"`
		ConnectionType       string             `json:"connectionType,omitempty"       yaml:"connectionType,omitempty"`
		Spec                 *GitConnectionSpec `json:"spec,omitempty"                 yaml:"spec,omitempty"`
	}

	GitConnectionSpec struct {
		Username    string `json:"username"    yaml:"username"`
		PasswordRef string `json:"passwordRef,omitempty" yaml:"passwordRef,omitempty"`
	}
)
