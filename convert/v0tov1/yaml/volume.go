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

type Volume struct {
	Name string      `json:"name,omitempty"`
	Uses string      `json:"uses,omitempty"`
	With interface{} `json:"with,omitempty"`
}

// type VolumeWith struct {
// 	Bind   *VolumeBind
// 	Claim  *VolumeClaim
// 	Config *VolumeConfigMap
// 	Temp   *VolumeTemp
// }

// UnmarshalJSON implement the json.Unmarshaler interface.
func (v *Volume) UnmarshalJSON(data []byte) error {
	type S Volume
	type T struct {
		*S
		With json.RawMessage `json:"with"`
	}

	obj := &T{S: (*S)(v)}
	if err := json.Unmarshal(data, obj); err != nil {
		return err
	}

	switch v.Uses {
	case "bind":
		v.With = new(VolumeBind)
	case "persistent-volume-claim":
		v.With = new(VolumeClaim)
	case "config-map":
		v.With = new(VolumeConfigMap)
	case "temp":
		v.With = new(VolumeTemp)
	case "secret":
		v.With = new(VolumeSecret)
	case "host-path":
		v.With = new(VolumeHostPath)
	case "empty-dir":
		v.With = new(VolumeEmptyDir)
	default:
		return fmt.Errorf("unknown uses %s", v.Uses)
	}

	return json.Unmarshal(obj.With, v.With)
}
