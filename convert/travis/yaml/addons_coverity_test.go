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

func TestCoverity(t *testing.T) {
	tests := []struct {
		yaml string
		want Coverity
	}{
		// bool value
		{
			yaml: `true`,
			want: Coverity{
				Enabled: true,
			},
		},
		// struct value
		{
			yaml: `{ enabled: true, notification_email: janecitizen@mail.com }`,
			want: Coverity{
				Enabled: true,
				NotificationEmail: &Secure{
					Decrypted: "janecitizen@mail.com",
				},
			},
		},
		// struct value
		{
			yaml: `{ enabled: true, notification_email: { secure: mcUCykGm4bUZ3CaW6AxrIMFzuAYjA98VIz6YmYTmM0 } }`,
			want: Coverity{
				Enabled: true,
				NotificationEmail: &Secure{
					Encrypted: "mcUCykGm4bUZ3CaW6AxrIMFzuAYjA98VIz6YmYTmM0",
				},
			},
		},
		// project string
		{
			yaml: `{ enabled: true, project: my_github/my_project }`,
			want: Coverity{
				Enabled: true,
				Project: &CoverityProject{
					Name: "my_github/my_project",
				},
			},
		},
		// project struct
		{
			yaml: `{ enabled: true, project: { name: my_github/my_project, version: 1.0, description: my project } }`,
			want: Coverity{
				Enabled: true,
				Project: &CoverityProject{
					Name:        "my_github/my_project",
					Version:     "1.0",
					Description: "my project",
				},
			},
		},
	}

	for i, test := range tests {
		got := new(Coverity)
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

func TestCoverity_Error(t *testing.T) {
	err := yaml.Unmarshal([]byte("[[]]"), new(Coverity))
	if err == nil || err.Error() != "failed to unmarshal coverity_scan" {
		t.Errorf("Expect error, got %s", err)
	}
}

func TestCoveritProject_Error(t *testing.T) {
	err := yaml.Unmarshal([]byte("[[]]"), new(CoverityProject))
	if err == nil || err.Error() != "failed to unmarshal coverity project" {
		t.Errorf("Expect error, got %s", err)
	}
}
