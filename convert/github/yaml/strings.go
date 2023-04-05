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
	"fmt"
)

// Stringorslice represents a string or an array of strings.
type Stringorslice []string

// UnmarshalYAML implements the unmarshal interface.
func (s *Stringorslice) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var stringType string
	if err := unmarshal(&stringType); err == nil {
		*s = []string{stringType}
		return nil
	}

	var sliceType []interface{}
	if err := unmarshal(&sliceType); err == nil {
		parts, err := toStrings(sliceType)
		if err != nil {
			return err
		}
		*s = parts
		return nil
	}

	return errors.New("failed to unmarshal string or string array")
}

// helper function converts a slice of interfaces
// to a slice of strings.
func toStrings(s []interface{}) ([]string, error) {
	if len(s) == 0 {
		return nil, nil
	}
	r := make([]string, len(s))
	for k, v := range s {
		switch vv := v.(type) {
		case string:
			r[k] = vv
		case int:
			r[k] = fmt.Sprint(vv)
		case float64:
			r[k] = fmt.Sprint(vv)
		case bool:
			r[k] = fmt.Sprint(vv)
		default:
			return nil, fmt.Errorf("cannot unmarshal %v of type %T into a string value", v, v)
		}
	}
	return r, nil
}
