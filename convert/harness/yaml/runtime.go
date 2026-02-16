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
)

type Runtime struct {
	Type string      `json:"type,omitempty"   yaml:"type,omitempty"`
	Spec interface{} `json:"spec,omitempty"   yaml:"spec,omitempty"`
}

// RuntimeCloudSpec defines the Cloud runtime spec
type RuntimeCloudSpec struct {
	Size      string         `json:"size,omitempty"      yaml:"size,omitempty"`
	ImageSpec *ImageSpec     `json:"imageSpec,omitempty" yaml:"imageSpec,omitempty"`
}

// ImageSpec defines the image specification for Cloud runtime
type ImageSpec struct {
	ImageName string `json:"imageName,omitempty" yaml:"imageName,omitempty"`
}

// RuntimeDockerSpec defines the Docker runtime spec
type RuntimeDockerSpec struct {
	HarnessImageConnectorRef string `json:"harnessImageConnectorRef,omitempty" yaml:"harnessImageConnectorRef,omitempty"`
}

// UnmarshalJSON implements custom unmarshalling for Runtime based on type
func (r *Runtime) UnmarshalJSON(data []byte) error {
	type Alias Runtime
	aux := &struct {
		Type string          `json:"type"`
		Spec json.RawMessage `json:"spec"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	r.Type = aux.Type

	// Unmarshal spec based on type
	switch aux.Type {
	case "Cloud":
		var spec RuntimeCloudSpec
		if err := json.Unmarshal(aux.Spec, &spec); err != nil {
			return err
		}
		r.Spec = &spec
	case "Docker":
		var spec RuntimeDockerSpec
		if err := json.Unmarshal(aux.Spec, &spec); err != nil {
			return err
		}
		r.Spec = &spec
	default:
		// For unknown types, keep as interface{}
		var spec interface{}
		if err := json.Unmarshal(aux.Spec, &spec); err != nil {
			return err
		}
		r.Spec = spec
	}

	return nil
}
