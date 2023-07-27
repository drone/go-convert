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
	"strings"
)

type Secret struct {
	Vault *Vault `yaml:"vault,omitempty"`
	File  *bool  `yaml:"file,omitempty"`
	Token string `yaml:"token,omitempty"`
}

type Vault struct {
	Engine *VaultEngine `yaml:"engine,omitempty"`
	Path   string       `yaml:"path,omitempty"`
	Field  string       `yaml:"field,omitempty"`
}

type VaultEngine struct {
	Name string `yaml:"name,omitempty"`
	Path string `yaml:"path,omitempty"`
}

// UnmarshalYAML implements the unmarshal interface.
func (v *Vault) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 string
	var out2 = struct {
		Engine *VaultEngine `yaml:"engine,omitempty"`
		Path   string       `yaml:"path,omitempty"`
		Field  string       `yaml:"field,omitempty"`
	}{}

	if err := unmarshal(&out1); err == nil {
		parts := strings.SplitN(out1, "/", 3)
		if len(parts) == 3 {
			v.Path = parts[0] + "/" + parts[1]
			v.Field = parts[2]
			engineParts := strings.SplitN(parts[2], "@", 2)
			if len(engineParts) == 2 {
				v.Field = engineParts[0]
				v.Engine = &VaultEngine{
					Path: engineParts[1],
					Name: "kv-v2",
				}
			}
		} else {
			v.Path = out1
		}
		return nil
	}

	if err := unmarshal(&out2); err == nil {
		v.Path = out2.Path
		v.Field = out2.Field
		v.Engine = out2.Engine
		return nil
	}

	return errors.New("failed to unmarshal vault")
}
