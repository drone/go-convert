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

// import "encoding/json"

// // List of cpu architectures.
// type Arch int

// // Arch enumeration.
// const (
// 	ArchNone Arch = iota
// 	ArchAmd64
// 	ArchArm64
// 	ArchPpc64le
// 	ArchS390x
// )

// // String returns the Arch as a string.
// func (e Arch) String() string {
// 	switch e {
// 	case ArchAmd64:
// 		return "amd64"
// 	case ArchArm64:
// 		return "arm64"
// 	case ArchPpc64le:
// 		return "ppc64le"
// 	case ArchS390x:
// 		return "s390x"
// 	default:
// 		return ""
// 	}
// }

// // MarshalJSON marshals the type as a JSON string.
// func (e Arch) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(e.String())
// }

// // UnmarshalJSON unmashals a quoted json string to the enum value.
// func (e *Arch) UnmarshalJSON(b []byte) error {
// 	var v string
// 	json.Unmarshal(b, &v)
// 	switch v {
// 	case "amd64":
// 		*e = ArchAmd64
// 	case "arm64", "arm", "arm32":
// 		*e = ArchArm64
// 	case "ppc64le", "ppc", "ppc64":
// 		*e = ArchPpc64le
// 	case "s390x", "s390":
// 		*e = ArchS390x
// 	default:
// 		*e = ArchNone
// 	}
// 	return nil
// }

// // UnmarshalJSON unmashals a quoted json string to the enum value.
// // UnmarshalYAML implements the unmarshal interface.
// func (e *Arch) UnmarshalYAML(unmarshal func(interface{}) error) error {
// 	var v string
// 	unmarshal(v)
// 	switch v {
// 	case "amd64":
// 		*e = ArchAmd64
// 	case "arm64", "arm", "arm32":
// 		*e = ArchArm64
// 	case "ppc64le", "ppc", "ppc64":
// 		*e = ArchPpc64le
// 	case "s390x", "s390":
// 		*e = ArchS390x
// 	default:
// 		*e = ArchNone
// 	}
// 	return nil
// }
