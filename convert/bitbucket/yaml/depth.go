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

// Depth configures the clone depth.
type Depth struct {
	Full  bool
	Value int
}

// MarshalYAML implements the marshal interface.
func (v *Depth) MarshalYAML() (interface{}, error) {
	if v.Full {
		return "full", nil
	} else if v.Value > 0 {
		return v.Value, nil
	} else {
		return nil, nil
	}
}

// UnmarshalYAML implements the unmarshal interface.
func (v *Depth) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 string
	var out2 int
	if err := unmarshal(&out1); err == nil {
		if out1 == "full" {
			v.Full = true
			return nil
		}
	}
	if err := unmarshal(&out2); err == nil {
		v.Value = out2
		return nil
	}
	return errors.New("failed to unmarshal depth")
}
