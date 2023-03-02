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
	"encoding/json"
	"fmt"
)

type (
	// Stage defines a pipeline stage.
	Stage struct {
		ID          string      `json:"identifier,omitempty"   yaml:"identifier,omitempty"`
		Description string      `json:"description,omitempty"  yaml:"description,omitempty"`
		Name        string      `json:"name,omitempty"         yaml:"name,omitempty"`
		Spec        interface{} `json:"spec,omitempty"         yaml:"spec,omitempty"`
		Type        string      `json:"type,omitempty"         yaml:"type,omitempty"`
		Vars        []Variable  `json:"variables,omitempty"    yaml:"variables,omitempty"`
	}

	// StageApproval defines an approval stage.
	StageApproval struct {
		// TODO
	}

	// StageCI defines a continuous integration stage.
	StageCI struct {
		Cache          Cache          `json:"cache,omitempty"              yaml:"cache,omitempty"`
		Clone          bool           `json:"cloneCodebase,omitempty"       yaml:"cloneCodebase,omitempty"`
		Execution      Execution      `json:"execution,omitempty"           yaml:"execution,omitempty"`
		Infrastructure Infrastructure `json:"infrastructure,omitempty"      yaml:"infrastructure,omitempty"`
		Services       []Service      `json:"serviceDependencies,omitempty" yaml:"serviceDependencies,omitempty"`
		SharedPaths    []string       `json:"sharedPaths,omitempty"         yaml:"sharedPaths,omitempty"`
	}

	// StageDeployment defines a deployment stage.
	StageDeployment struct {
		// TODO
	}

	// StageFeatureFlag defines a feature flag stage.
	StageFeatureFlag struct {
		Execution Execution `json:"execution,omitempty" yaml:"execution,omitempty"`
	}
)

// UnmarshalJSON implement the json.Unmarshaler interface.
func (s *Stage) UnmarshalJSON(data []byte) error {
	type S Stage
	type T struct {
		*S
		Spec json.RawMessage `json:"spec"`
	}

	obj := &T{S: (*S)(s)}
	if err := json.Unmarshal(data, obj); err != nil {
		return err
	}

	switch s.Type {
	case StageTypeCI:
		s.Spec = new(StageCI)
	case StageTypeFeatureFlag:
		s.Spec = new(StageFeatureFlag)
	default:
		return fmt.Errorf("unknown stage type %s", s.Type)
	}
	return json.Unmarshal(obj.Spec, s.Spec)
}
