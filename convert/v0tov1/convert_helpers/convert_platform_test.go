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

package converthelpers

import (
	"testing"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/google/go-cmp/cmp"
)

func TestConvertPlatform(t *testing.T) {
	tests := []struct {
		name string
		in   *v0.Platform
		want *v1.Platform
	}{
		{
			name: "nil platform stays nil",
			in:   nil,
			want: nil,
		},
		{
			// The load-bearing case for CI-23632: V0 pipelines commonly declare
			// only platform.os. PlatformV1.toPlatform() NPEs on arch.getValue()
			// when arch is missing, so we default to amd64 here.
			name: "os set, arch empty defaults to amd64",
			in:   &v0.Platform{OS: "Linux"},
			want: &v1.Platform{Os: "linux", Arch: "amd64"},
		},
		{
			name: "explicit arch wins over default",
			in:   &v0.Platform{OS: "Linux", Arch: "arm64"},
			want: &v1.Platform{Os: "linux", Arch: "arm64"},
		},
		{
			name: "os and arch both lowercased",
			in:   &v0.Platform{OS: "MacOS", Arch: "ARM64"},
			want: &v1.Platform{Os: "macos", Arch: "arm64"},
		},
		{
			// Empty platform (both fields blank) stays empty — we don't default
			// arch when there's no OS either. The V1 backend does not require
			// platform at all; only when present with an OS.
			name: "empty platform stays empty",
			in:   &v0.Platform{},
			want: &v1.Platform{Os: "", Arch: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertPlatform(tt.in)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ConvertPlatform() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
