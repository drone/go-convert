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

package normalize

import (
	"testing"

	"github.com/bradrydzewski/spec/yaml"
	"github.com/google/go-cmp/cmp"
)

func TestNormalizer(t *testing.T) {
	got := &yaml.Schema{
		Pipeline: &yaml.Pipeline{
			Stages: []*yaml.Stage{
				{Id: "stage1"},
				{Id: "stage2"},
				{Id: "stage3", Group: &yaml.StageGroup{
					Stages: []*yaml.Stage{
						{Id: "stage1"},
						{Id: "stage2"},
					},
				}},
				{Id: "stage4", Steps: []*yaml.Step{
					{Id: "step1"},
					{Id: "step2"},
					{Id: "step3", Group: &yaml.StepGroup{
						Steps: []*yaml.Step{
							{Id: "step1"},
							{Id: "step2"},
						},
					}},
				}},
			},
		},
	}

	Normalize(got)

	want := &yaml.Schema{
		Pipeline: &yaml.Pipeline{
			Stages: []*yaml.Stage{
				{Id: "stage1"},
				{Id: "stage2"},
				{Id: "stage3", Group: &yaml.StageGroup{
					Stages: []*yaml.Stage{
						{Id: "stage11"},
						{Id: "stage21"},
					},
				}},
				{Id: "stage4", Steps: []*yaml.Step{
					{Id: "step1"},
					{Id: "step2"},
					{Id: "step3", Group: &yaml.StepGroup{
						Steps: []*yaml.Step{
							{Id: "step11"},
							{Id: "step21"},
						},
					}},
				}},
			},
		},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Log(diff)
		t.Fail()
	}
}
