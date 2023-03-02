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

const mockPipe = `name: Send alert to Opsgenie
pipe: atlassian/opsgenie-send-alert:latest
variables:
    GENIE_KEY: $GENIE_KEY
    MESSAGE: Wake up!
`

func TestScript(t *testing.T) {
	tests := []struct {
		yaml string
		want Script
	}{
		{
			yaml: `"echo hello world"`,
			want: Script{
				Text: "echo hello world",
			},
		},
		{
			yaml: mockPipe,
			want: Script{
				Pipe: &Pipe{
					Image: "atlassian/opsgenie-send-alert:latest",
					Name:  "Send alert to Opsgenie",
					Variables: map[string]string{
						"GENIE_KEY": "$GENIE_KEY",
						"MESSAGE":   "Wake up!",
					},
				},
			},
		},
	}

	for i, test := range tests {
		got := new(Script)
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

func TestScript_Error(t *testing.T) {
	err := yaml.Unmarshal([]byte("[[]]"), new(Script))
	if err == nil || err.Error() != "failed to unmarshal script" {
		t.Errorf("Expect error, got %s", err)
	}
}

func TestScript_Marshal(t *testing.T) {
	tests := []struct {
		before Script
		after  string
	}{
		{
			before: Script{Text: "echo hello world"},
			after:  "echo hello world\n",
		},
		{
			before: Script{
				Pipe: &Pipe{
					Image: "atlassian/opsgenie-send-alert:latest",
					Name:  "Send alert to Opsgenie",
					Variables: map[string]string{
						"GENIE_KEY": "$GENIE_KEY",
						"MESSAGE":   "Wake up!",
					},
				},
			},
			after: mockPipe,
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
