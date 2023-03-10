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

// List of machine sizes.
type Size int

// Size enumeration.
const (
	SizeNone Size = iota
	Size1x
	Size2x
	Size4x
	Size8x
)

// String returns the Size as a string.
func (e Size) String() string {
	switch e {
	case Size1x:
		return "1x"
	case Size2x:
		return "2x"
	case Size4x:
		return "4x"
	case Size8x:
		return "8x"
	default:
		return ""
	}
}

// UnmarshalYAML implements the unmarshal interface.
func (e *Size) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v string
	unmarshal(&v)
	switch v {
	case "1x":
		*e = Size1x
	case "2x":
		*e = Size2x
	case "4x":
		*e = Size4x
	case "8x":
		*e = Size8x
	default:
		*e = SizeNone
	}
	return nil
}

// MarshalJSON marshals the enum as a quoted json string
func (e Size) MarshalYAML() (interface{}, error) {
	if e == SizeNone {
		return nil, nil
	} else {
		return e.String(), nil
	}
}
