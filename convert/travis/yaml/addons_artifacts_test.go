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

func TestArtifacts(t *testing.T) {
	tests := []struct {
		yaml string
		want Artifacts
	}{
		// bool value
		{
			yaml: `true`,
			want: Artifacts{
				Enabled: true,
			},
		},
		// struct value with key
		{
			yaml: `{ enabled: true, key: AWS1234566789 }`,
			want: Artifacts{
				Enabled: true,
				Key: &Secure{
					Decrypted: "AWS1234566789",
				},
			},
		},
		// struct value with key as secret
		{
			yaml: `{ enabled: true, key: { secure: AWS1234566789 } }`,
			want: Artifacts{
				Enabled: true,
				Key: &Secure{
					Encrypted: "AWS1234566789",
				},
			},
		},
		// struct value with alias aws_access_key_id, aws_secret_access_key
		{
			yaml: `{ enabled: true, aws_access_key_id: AWS1234566789, aws_secret_access_key: da39a3ee5 }`,
			want: Artifacts{
				Enabled: true,
				Key: &Secure{
					Decrypted: "AWS1234566789",
				},
				Secret: &Secure{
					Decrypted: "da39a3ee5",
				},
			},
		},
		// struct value with alias aws_access_key, aws_secret_key
		{
			yaml: `{ enabled: true, aws_access_key: AWS1234566789, aws_secret_key: da39a3ee5 }`,
			want: Artifacts{
				Enabled: true,
				Key: &Secure{
					Decrypted: "AWS1234566789",
				},
				Secret: &Secure{
					Decrypted: "da39a3ee5",
				},
			},
		},
		// struct value with alias access_key_id
		{
			yaml: `{ enabled: true, access_key_id: AWS1234566789, secret_access_key: da39a3ee5 }`,
			want: Artifacts{
				Enabled: true,
				Key: &Secure{
					Decrypted: "AWS1234566789",
				},
				Secret: &Secure{
					Decrypted: "da39a3ee5",
				},
			},
		},
		// struct value with alias access_key
		{
			yaml: `{ enabled: true, access_key: AWS1234566789, secret_key: da39a3ee5 }`,
			want: Artifacts{
				Enabled: true,
				Key: &Secure{
					Decrypted: "AWS1234566789",
				},
				Secret: &Secure{
					Decrypted: "da39a3ee5",
				},
			},
		},
	}

	for i, test := range tests {
		got := new(Artifacts)
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

func TestArtifacts_Error(t *testing.T) {
	err := yaml.Unmarshal([]byte("[[]]"), new(Artifacts))
	if err == nil || err.Error() != "failed to unmarshal artifacts" {
		t.Errorf("Expect error, got %s", err)
	}
}
