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

type Parallel struct {
	Count  int                   `yaml:"-,omitempty"`
	Matrix []map[string][]string `yaml:"matrix,omitempty"`
}

// UnmarshalYAML implements the unmarshal interface.
func (v *Parallel) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 int
	var out2 map[string][]map[string]interface{}

	if err := unmarshal(&out1); err == nil {
		v.Count = out1
		return nil
	}

	if err := unmarshal(&out2); err == nil {
		matrix := make([]map[string][]string, len(out2["matrix"]))
		for i, item := range out2["matrix"] {
			matrixItem := make(map[string][]string)
			for key, value := range item {
				switch v := value.(type) {
				case string:
					matrixItem[key] = []string{v}
				case []interface{}:
					strings := make([]string, len(v))
					for i, s := range v {
						strings[i] = s.(string)
					}
					matrixItem[key] = strings
				}
			}
			matrix[i] = matrixItem
		}
		v.Matrix = matrix
		return nil
	}

	return errors.New("failed to unmarshal parallel")
}

func (v *Parallel) MarshalYAML() (interface{}, error) {
	if v.Count > 0 {
		return v.Count, nil
	}

	// Convert the complex structure of v.Matrix back into a simpler form
	matrix := make([]map[string]interface{}, len(v.Matrix))
	for i, item := range v.Matrix {
		matrixItem := make(map[string]interface{})
		for key, value := range item {
			if len(value) == 1 {
				matrixItem[key] = value[0] // Single string
			} else {
				matrixItem[key] = value // Slice of strings
			}
		}
		matrix[i] = matrixItem
	}

	return map[string][]map[string]interface{}{"matrix": matrix}, nil
}
