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
	"errors"
)

type (
	// Secure represents a secure variable.
	Secure struct {
		Decrypted string
		Encrypted string
	}

	// secure is a temporary structure used for
	// encoding and decoding.
	secure struct {
		Secure string `yaml:"secure,omitempty" json:"secure,omitempty"`
	}
)

// UnmarshalYAML implements the unmarshal interface.
func (v *Secure) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 string
	var out2 *secure
	if err := unmarshal(&out1); err == nil {
		v.Decrypted = out1
		return nil
	}
	if err := unmarshal(&out2); err == nil {
		v.Encrypted = out2.Secure
		return nil
	}
	return errors.New("failed to unmarshal secure variable")
}

// MarshalYAML implements the marshal interface.
func (v *Secure) MarshalYAML() (interface{}, error) {
	if v.Encrypted != "" {
		return secure{v.Encrypted}, nil
	}
	return v.Decrypted, nil
}

// UnmarshalJSON implements the unmarshal interface.
func (v *Secure) UnmarshalJSON(data []byte) error {
	var out1 string
	var out2 *secure
	if err := json.Unmarshal(data, &out1); err == nil {
		v.Decrypted = out1
		return nil
	}
	if err := json.Unmarshal(data, &out2); err == nil {
		v.Encrypted = out2.Secure
		return nil
	}
	return errors.New("failed to unmarshal secure variable")
}

// MarshalJSON implements the marshal interface.
func (v *Secure) MarshalJSON() ([]byte, error) {
	if v.Encrypted != "" {
		return json.Marshal(
			&secure{v.Encrypted},
		)
	}
	return json.Marshal(v.Decrypted)
}
