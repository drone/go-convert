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

type (
	// Image configures the container image.
	Image struct {
		Name      string
		Username  string
		Password  string
		Email     string
		RunAsUser int
		AWS       *AWS
	}

	// temporary data structure for unmarshaling and
	// marsaling images.
	image struct {
		Name      string `yaml:"name,omitempty"`
		Username  string `yaml:"username,omitempty"`
		Password  string `yaml:"password,omitempty"`
		Email     string `yaml:"email,omitempty"`
		RunAsUser int    `yaml:"run-as-user,omitempty"`
		AWS       *AWS   `yaml:"aws,omitempty"`
	}
)

// UnmarshalYAML implements the unmarshal interface.
func (v *Image) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 string
	var out2 *image
	if err := unmarshal(&out1); err == nil {
		v.Name = out1
		return nil
	}
	if err := unmarshal(&out2); err == nil {
		v.Name = out2.Name
		v.Username = out2.Username
		v.Password = out2.Password
		v.Email = out2.Email
		v.RunAsUser = out2.RunAsUser
		v.AWS = out2.AWS
		return nil
	}
	return errors.New("failed to unmarshal image")
}

// MarshalYAML implements the marshal interface.
func (v *Image) MarshalYAML() (interface{}, error) {
	// marshal the image using short syntax if only the image
	// name is provided.
	if v.Username == "" && v.Password == "" && v.Email == "" && v.RunAsUser == 0 && v.AWS == nil {
		return v.Name, nil
	}
	// else marshal the image using the long syntax.
	return &image{
		Name:      v.Name,
		Username:  v.Username,
		Password:  v.Password,
		Email:     v.Email,
		RunAsUser: v.RunAsUser,
		AWS:       v.AWS,
	}, nil
}
