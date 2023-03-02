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
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
)

func TestSecure_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		yaml string
		want Secure
	}{
		// test string value
		{
			yaml: `"pa55word"`,
			want: Secure{
				Decrypted: "pa55word",
			},
		},
		// test encrypted value
		{
			yaml: `{secure: mcUCykGm4bUZ3CaW6AxrIMFzuAYjA98VIz6YmYTmM0}`,
			want: Secure{
				Encrypted: "mcUCykGm4bUZ3CaW6AxrIMFzuAYjA98VIz6YmYTmM0",
			},
		},
	}

	for i, test := range tests {
		got := new(Secure)
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

func TestSecure_UnmarshalYAML_Error(t *testing.T) {
	err := yaml.Unmarshal([]byte("[[]]"), new(Secure))
	if err == nil || err.Error() != "failed to unmarshal secure variable" {
		t.Errorf("Expect error, got %s", err)
	}
}

func TestSecure_MarshalYAML(t *testing.T) {
	tests := []struct {
		before Secure
		after  string
	}{
		{
			before: Secure{Decrypted: "pa55word"},
			after:  "pa55word\n",
		},
		{
			before: Secure{Encrypted: "pa55word"},
			after:  "secure: pa55word\n",
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

func TestSecure_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		json string
		want Secure
	}{
		// test string value
		{
			json: `"pa55word"`,
			want: Secure{
				Decrypted: "pa55word",
			},
		},
		// test encrypted value
		{
			json: `{"secure":"mcUCykGm4bUZ3CaW6AxrIMFzuAYjA98VIz6YmYTmM0"}`,
			want: Secure{
				Encrypted: "mcUCykGm4bUZ3CaW6AxrIMFzuAYjA98VIz6YmYTmM0",
			},
		},
	}

	for i, test := range tests {
		got := new(Secure)
		if err := json.Unmarshal([]byte(test.json), got); err != nil {
			t.Log(test.json)
			t.Error(err)
			return
		}
		if diff := cmp.Diff(got, &test.want); diff != "" {
			t.Log(test.json)
			t.Errorf("Unexpected parsing results for test %v", i)
			t.Log(diff)
		}
	}
}

func TestSecure_UnmarshalJSON_Error(t *testing.T) {
	err := json.Unmarshal([]byte("[[]]"), new(Secure))
	if err == nil || err.Error() != "failed to unmarshal secure variable" {
		t.Errorf("Expect error, got %s", err)
	}
}

func TestSecure_MarshalJSON(t *testing.T) {
	tests := []struct {
		before Secure
		after  string
	}{
		{
			before: Secure{Decrypted: "pa55word"},
			after:  `"pa55word"`,
		},
		{
			before: Secure{Encrypted: "pa55word"},
			after:  `{"secure":"pa55word"}`,
		},
	}

	for _, test := range tests {
		after, err := json.Marshal(&test.before)
		if err != nil {
			t.Error(err)
			return
		}
		if got, want := string(after), test.after; got != want {
			t.Errorf("want json %s, got %s", want, got)
		}
	}
}
