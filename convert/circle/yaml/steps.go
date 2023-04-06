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

type (
	Step struct {
		AddSSHKeys         *AddSSHKeys
		AttachWorkspace    *AttachWorkspace
		Checkout           *Checkout
		PersistToWorkspace *PersistToWorkspace
		RestoreCache       *RestoreCache
		Run                *Run
		SaveCache          *SaveCache
		SetupRemoteDocker  *SetupRemoteDocker
		StoreArtifacts     *StoreArtifacts
		StoreTestResults   *StoreTestResults
		Unless             *Unless
		When               *When
		Custom             *Custom
	}

	step struct {
		AddSSHKeys         *AddSSHKeys            `yaml:"add_ssh_keys,omitempty"`
		AttachWorkspace    *AttachWorkspace       `yaml:"attach_workspace,omitempty"`
		Checkout           *Checkout              `yaml:"checkout,omitempty"`
		PersistToWorkspace *PersistToWorkspace    `yaml:"persist_to_workspace,omitempty"`
		RestoreCache       *RestoreCache          `yaml:"restore_cache,omitempty"`
		Run                *Run                   `yaml:"run,omitempty"`
		SaveCache          *SaveCache             `yaml:"save_cache,omitempty"`
		SetupRemoteDocker  *SetupRemoteDocker     `yaml:"setup_remote_docker,omitempty"`
		StoreArtifacts     *StoreArtifacts        `yaml:"store_artifacts,omitempty"`
		StoreTestResults   *StoreTestResults      `yaml:"store_test_results,omitempty"`
		Unless             *Unless                `yaml:"unless,omitempty"`
		When               *When                  `yaml:"when,omitempty"`
		Custom             map[string]interface{} `yaml:",inline"`
	}
)

// UnmarshalYAML implements the unmarshal interface.
func (v *Step) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 string
	var out2 *step

	if err := unmarshal(&out1); err == nil {
		switch out1 {
		case "add_ssh_keys":
			v.AddSSHKeys = new(AddSSHKeys)
		case "attach_workspace":
			v.AttachWorkspace = new(AttachWorkspace)
		case "checkout":
			v.Checkout = new(Checkout)
		case "persist_to_workspace":
			v.PersistToWorkspace = new(PersistToWorkspace)
		case "restore_cache":
			v.RestoreCache = new(RestoreCache)
		case "run":
			v.Run = new(Run)
		case "save_cache":
			v.SaveCache = new(SaveCache)
		case "setup_remote_docker":
			v.SetupRemoteDocker = new(SetupRemoteDocker)
		case "store_artifacts":
			v.StoreArtifacts = new(StoreArtifacts)
		case "store_test_results":
			v.StoreTestResults = new(StoreTestResults)
		case "unless":
			v.Unless = new(Unless)
		case "when":
			v.When = new(When)
		default:
			v.Custom = new(Custom)
			v.Custom.Name = out1
			v.Custom.Params = map[string]interface{}{}
		}
		return nil
	}

	if err := unmarshal(&out2); err == nil {
		v.AddSSHKeys = out2.AddSSHKeys
		v.AttachWorkspace = out2.AttachWorkspace
		v.Checkout = out2.Checkout
		v.PersistToWorkspace = out2.PersistToWorkspace
		v.RestoreCache = out2.RestoreCache
		v.Run = out2.Run
		v.SaveCache = out2.SaveCache
		v.SetupRemoteDocker = out2.SetupRemoteDocker
		v.StoreArtifacts = out2.StoreArtifacts
		v.StoreTestResults = out2.StoreTestResults
		v.Unless = out2.Unless
		v.When = out2.When
		for name, params := range out2.Custom {
			v.Custom = new(Custom)
			v.Custom.Name = name
			if vv, ok := params.(map[string]interface{}); ok {
				v.Custom.Params = vv
			}
		}

		return nil
	}

	return errors.New("failed to unmarshal step")
}

// MarshalYAML implements the marshal interface.
func (v *Step) MarshalYAML() (interface{}, error) {
	if v.Custom != nil {
		return map[string]interface{}{
			v.Custom.Name: v.Custom.Params,
		}, nil
	}
	return &step{
		AddSSHKeys:         v.AddSSHKeys,
		AttachWorkspace:    v.AttachWorkspace,
		Checkout:           v.Checkout,
		PersistToWorkspace: v.PersistToWorkspace,
		RestoreCache:       v.RestoreCache,
		Run:                v.Run,
		SaveCache:          v.SaveCache,
		SetupRemoteDocker:  v.SetupRemoteDocker,
		StoreArtifacts:     v.StoreArtifacts,
		StoreTestResults:   v.StoreTestResults,
		Unless:             v.Unless,
		When:               v.When,
	}, nil
}

//
// Step Types
//

type (
	AddSSHKeys struct {
		Fingerprints []string `yaml:"fingerprints,omitempty"`
		Name         string   `yaml:"name,omitempty"`
	}

	AttachWorkspace struct {
		At   string  `yaml:"at,omitempty"`
		Name *string `yaml:"name,omitempty"`
	}

	Checkout struct {
		Name string `yaml:"name,omitempty"`
		Path string `yaml:"path,omitempty"`
	}

	Custom struct {
		Name   string
		Params map[string]interface{}
	}

	PersistToWorkspace struct {
		Name  string        `yaml:"name,omitempty"`
		Paths Stringorslice `yaml:"paths,omitempty"`
		Root  string        `yaml:"root,omitempty"`
	}

	RestoreCache struct {
		Key  string        `yaml:"key,omitempty"`
		Name string        `yaml:"name,omitempty"`
		Keys Stringorslice `yaml:"keys,omitempty"`
	}

	SaveCache struct {
		Key   string        `yaml:"key,omitempty"`
		Name  string        `yaml:"name,omitempty"`
		Paths Stringorslice `yaml:"paths,omitempty"`
		When  string        `yaml:"when,omitempty"`
	}

	SetupRemoteDocker struct {
		DockerLayerCaching bool   `yaml:"docker_layer_caching,omitempty"`
		Name               string `yaml:"name,omitempty"`
		Version            string `yaml:"version,omitempty"`
	}

	StoreArtifacts struct {
		Destination string `yaml:"destination,omitempty"`
		Name        string `yaml:"name,omitempty"`
		Path        string `yaml:"path,omitempty"`
	}

	StoreTestResults struct {
		Name string `yaml:"name,omitempty"`
		Path string `yaml:"path,omitempty"`
	}

	Unless struct {
		Condition *Logical `json:"condition,omitempty"`
		Steps     []*Step  `json:"steps,omitempty"`
	}

	When struct {
		Condition *Logical `json:"condition,omitempty"`
		Steps     []*Step  `json:"steps,omitempty"`
	}
)
