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
)

func TestNormalize(t *testing.T) {
	tests := []struct {
		before []*Steps
		after  []*Steps
	}{
		// group single step into stage
		{
			before: []*Steps{
				{Step: &Step{Name: "test"}},
			},
			after: []*Steps{
				{
					Stage: &Stage{
						Steps: []*Steps{
							{Step: &Step{Name: "test"}},
						},
					},
				},
			},
		},
		// group multiple steps into stage
		{
			before: []*Steps{
				{Step: &Step{Name: "build"}},
				{Step: &Step{Name: "test"}},
			},
			after: []*Steps{
				{
					Stage: &Stage{
						Steps: []*Steps{
							{Step: &Step{Name: "build"}},
							{Step: &Step{Name: "test"}},
						},
					},
				},
			},
		},
		// no change when steps already grouped
		// into top-level stages
		{
			before: []*Steps{
				{Stage: &Stage{Name: "stage1"}},
				{Stage: &Stage{Name: "stage2"}},
			},
			after: []*Steps{
				{Stage: &Stage{Name: "stage1"}},
				{Stage: &Stage{Name: "stage2"}},
			},
		},
		// handle a mix of steps and stages
		{
			before: []*Steps{
				{Step: &Step{Name: "step1"}},
				{Stage: &Stage{Name: "stage1"}},
				{Step: &Step{Name: "step2"}},
				{Step: &Step{Name: "step3"}},
				{Stage: &Stage{Name: "stage2"}},
				{Step: &Step{Name: "step4"}},
				{Step: &Step{Name: "step5"}},
			},
			after: []*Steps{
				{
					Stage: &Stage{
						Steps: []*Steps{
							{Step: &Step{Name: "step1"}},
						},
					},
				},
				{Stage: &Stage{Name: "stage1"}},
				{
					Stage: &Stage{
						Steps: []*Steps{
							{Step: &Step{Name: "step2"}},
							{Step: &Step{Name: "step3"}},
						},
					},
				},
				{Stage: &Stage{Name: "stage2"}},
				{
					Stage: &Stage{
						Steps: []*Steps{
							{Step: &Step{Name: "step4"}},
							{Step: &Step{Name: "step5"}},
						},
					},
				},
			},
		},
	}

	for i, test := range tests {
		got := new(Config)
		got.Pipelines.Default = test.before

		want := new(Config)
		want.Pipelines.Default = test.after

		Normalize(got)

		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("Unexpected parsing results for test %v", i)
			t.Log(diff)
		}
	}
}
