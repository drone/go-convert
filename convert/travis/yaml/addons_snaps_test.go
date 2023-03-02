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

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
)

func TestSnap(t *testing.T) {
	tests := []struct {
		yaml string
		want Snaps
	}{
		{
			yaml: `"mysql.snap"`,
			want: Snaps{
				Items: []*Snap{
					{Name: "mysql.snap"},
				},
			},
		},
		{
			yaml: `[ "mysql.snap" ]`,
			want: Snaps{
				Items: []*Snap{
					{Name: "mysql.snap"},
				},
			},
		},
		{
			yaml: `{ name: mysql.snap, classic: true, channel: edge }`,
			want: Snaps{
				Items: []*Snap{
					{Name: "mysql.snap", Channel: "edge", Classic: true},
				},
			},
		},
		{
			yaml: `[{ name: mysql.snap, classic: true, channel: edge }]`,
			want: Snaps{
				Items: []*Snap{
					{Name: "mysql.snap", Channel: "edge", Classic: true},
				},
			},
		},
	}

	for i, test := range tests {
		got := new(Snaps)
		if err := yaml.Unmarshal([]byte(test.yaml), got); err != nil {
			t.Log(test.yaml)
			t.Error(err)
			return
		}
		if diff := cmp.Diff(got, &test.want); diff != "" {
			t.Log(test.yaml)
			t.Errorf("Unexpected parsing results for test %v", i)
			t.Log(diff)
		}
	}
}

func TestSnap_Error(t *testing.T) {
	err := yaml.Unmarshal([]byte("[[]]"), new(Snap))
	if err == nil || err.Error() != "failed to unmarshal snap" {
		t.Errorf("Expect error, got %s", err)
	}
}

func TestSnaps_Error(t *testing.T) {
	err := yaml.Unmarshal([]byte("[[]]"), new(Snaps))
	if err == nil || err.Error() != "failed to unmarshal snaps" {
		t.Errorf("Expect error, got %s", err)
	}
}

func TestSnaps_Marshal(t *testing.T) {
	tests := []struct {
		before Snaps
		after  string
	}{
		{
			before: Snaps{Items: []*Snap{{Name: "mysql"}}},
			after:  "- name: mysql\n",
		},
		{
			before: Snaps{Items: []*Snap{{Name: "mysql", Classic: true, Channel: "edge"}}},
			after:  "- name: mysql\n  classic: true\n  channel: edge\n",
		},
	}

	for _, test := range tests {
		after, err := yaml.Marshal(&test.before)
		if err != nil {
			t.Error(err)
			return
		}
		if got, want := string(after), test.after; got != want {
			t.Errorf("want yaml %q, got %q", want, got)
		}
	}
}
