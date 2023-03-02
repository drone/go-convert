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

func TestCodeclimate(t *testing.T) {
	tests := []struct {
		yaml string
		want Codeclimate
	}{
		// bool value
		{
			yaml: `true`,
			want: Codeclimate{
				Enabled: true,
			},
		},
		// struct value
		{
			yaml: `{ enabled: true, repo_token: da39a3ee5 }`,
			want: Codeclimate{
				Enabled: true,
				RepoToken: &Secure{
					Decrypted: "da39a3ee5",
				},
			},
		},
		// struct value with secure values
		{
			yaml: `{ enabled: true, repo_token: { secure: mcUCykGm4bUZ3CaW6AxrIMFzuAYjA98VIz6YmYTmM0 } }`,
			want: Codeclimate{
				Enabled: true,
				RepoToken: &Secure{
					Encrypted: "mcUCykGm4bUZ3CaW6AxrIMFzuAYjA98VIz6YmYTmM0",
				},
			},
		},
		// struct value with secure values
		{
			yaml: `{ secure: mcUCykGm4bUZ3CaW6AxrIMFzuAYjA98VIz6YmYTmM0 }`,
			want: Codeclimate{
				Enabled: true,
				RepoToken: &Secure{
					Encrypted: "mcUCykGm4bUZ3CaW6AxrIMFzuAYjA98VIz6YmYTmM0",
				},
			},
		},
	}

	for i, test := range tests {
		got := new(Codeclimate)
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

func TestCodeclimate_Error(t *testing.T) {
	err := yaml.Unmarshal([]byte("[[]]"), new(Codeclimate))
	if err == nil || err.Error() != "failed to unmarshal codeclimate" {
		t.Errorf("Expect error, got %s", err)
	}
}
