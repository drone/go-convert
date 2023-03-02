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

// // List of operating systems.
// type OS int

// // OS enumeration.
// const (
// 	OSNone OS = iota
// 	OSLinux
// 	OSWindows
// 	OSMacos
// )

// // String returns the OS as a string.
// func (e OS) String() string {
// 	switch e {
// 	case OSLinux:
// 		return "linux"
// 	case OSWindows:
// 		return "windows"
// 	case OSMacos:
// 		return "osx"
// 	default:
// 		return ""
// 	}
// }

// // MarshalJSON marshals the type as a JSON string.
// func (e OS) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(e.String())
// }

// // UnmarshalJSON unmashals a quoted json string to the enum value.
// func (e *OS) UnmarshalJSON(b []byte) error {
// 	var v string
// 	json.Unmarshal(b, &v)
// 	switch v {
// 	case "linux":
// 		*e = OSLinux
// 	case "windows":
// 		*e = OSWindows
// 	case "macos", "mac", "osx":
// 		*e = OSMacos
// 	default:
// 		*e = OSNone
// 	}
// 	return nil
// }

// // UnmarshalJSON unmashals a quoted json string to the enum value.
// // UnmarshalYAML implements the unmarshal interface.
// func (e *OS) UnmarshalYAML(unmarshal func(interface{}) error) error {
// 	var v string
// 	unmarshal(v)
// 	switch v {
// 	case "linux":
// 		*e = OSLinux
// 	case "windows":
// 		*e = OSWindows
// 	case "macos", "mac", "osx":
// 		*e = OSMacos
// 	default:
// 		*e = OSNone
// 	}
// 	return nil
// }
