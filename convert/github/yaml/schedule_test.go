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

func TestSchedule(t *testing.T) {
	tests := []struct {
		yaml string
		want Schedule
	}{
		{
			yaml: `{ cron: '30 5,17 * * *' }`,
			want: Schedule{
				Items: []*ScheduleItem{{Cron: "30 5,17 * * *"}},
			},
		},
		{
			yaml: `[ { cron: '30 5 * * 1,3' }, { cron: '30 5 * * 2,4' } ]`,
			want: Schedule{
				Items: []*ScheduleItem{
					{Cron: "30 5 * * 1,3"},
					{Cron: "30 5 * * 2,4"},
				},
			},
		},
	}

	for i, test := range tests {
		got := new(Schedule)
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

func TestSchedule_Marshal(t *testing.T) {
	tests := []struct {
		before Schedule
		after  string
	}{
		{
			before: Schedule{Items: []*ScheduleItem{
				{Cron: "30 5,17 * * *"},
			}},
			after: "- cron: 30 5,17 * * *\n",
		},
		{
			before: Schedule{Items: []*ScheduleItem{
				{Cron: "30 5 * * 1,3"},
				{Cron: "30 5 * * 2,4"},
			}},
			after: "- cron: 30 5 * * 1,3\n- cron: 30 5 * * 2,4\n",
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

func TestSchedule_Error(t *testing.T) {
	err := yaml.Unmarshal([]byte("[[]]"), new(Schedule))
	if err == nil || err.Error() != "failed to unmarshal on.schedule" {
		t.Errorf("Expect error, got %s", err)
	}
}
