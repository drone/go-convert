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
	"testing"

	"gopkg.in/yaml.v3"
)

func TestSize_String(t *testing.T) {
	tests := []struct {
		size Size
		want string
	}{
		{
			size: Size1x,
			want: "1x",
		},
		{
			size: Size2x,
			want: "2x",
		},
		{
			size: Size4x,
			want: "4x",
		},
		{
			size: Size8x,
			want: "8x",
		},
		{
			size: SizeNone,
			want: "",
		},
	}
	for _, test := range tests {
		if got, want := test.size.String(), test.want; got != want {
			t.Errorf("Want Size %s, got %s", want, got)
		}
	}
}

func TestSize_Marshal(t *testing.T) {
	tests := []struct {
		size Size
		want string
	}{
		{
			size: Size1x,
			want: "1x\n",
		},
		{
			size: Size2x,
			want: "2x\n",
		},
		{
			size: Size4x,
			want: "4x\n",
		},
		{
			size: Size8x,
			want: "8x\n",
		},
		{
			size: SizeNone,
			want: "null\n",
		},
	}
	for _, test := range tests {
		got, err := yaml.Marshal(test.size)
		if err != nil {
			t.Error(err)
			return
		}
		if got, want := string(got), test.want; got != want {
			t.Errorf("Want Size %q, got %q", want, got)
		}
	}
}

func TestSize_Unmarshal(t *testing.T) {
	tests := []struct {
		want Size
		yaml string
	}{
		{
			want: Size1x,
			yaml: "1x\n",
		},
		{
			want: Size2x,
			yaml: "2x\n",
		},
		{
			want: Size4x,
			yaml: "4x\n",
		},
		{
			want: Size8x,
			yaml: "8x\n",
		},
		{
			want: SizeNone,
			yaml: "99x\n", // ignore unknown values
		},
		{
			want: SizeNone,
			yaml: "\n",
		},
	}
	for _, test := range tests {
		var in Size
		err := yaml.Unmarshal([]byte(test.yaml), &in)
		if err != nil {
			t.Error(err)
			return
		}
		if got, want := in, test.want; got != want {
			t.Errorf("Want Size %v, got %v", want, got)
		}
	}
}
