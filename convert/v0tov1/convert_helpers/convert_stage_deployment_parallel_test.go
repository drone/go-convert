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

package converthelpers

import (
	"testing"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	"github.com/drone/go-convert/internal/flexible"
)

// v0 defaulted to parallel execution, while v1 now defaults to serial. The converter
// must therefore emit an explicit `parallel` field so a migrated pipeline keeps its
// original v0 behaviour:
//   - v0 parallel: true       -> v1 parallel: true
//   - v0 parallel: false      -> v1 parallel: false
//   - v0 parallel unspecified -> v1 parallel: true (preserve v0's parallel default)

func boolPtr(b bool) *bool { return &b }

// assertParallel checks that the emitted *flexible.Field[bool] carries the expected value.
func assertParallel(t *testing.T, got *flexible.Field[bool], want *bool) {
	t.Helper()
	if want == nil {
		if got != nil {
			t.Fatalf("expected no parallel field, got %+v", got)
		}
		return
	}
	if got == nil {
		t.Fatalf("expected parallel=%v, got nil field", *want)
	}
	v, ok := got.AsStruct()
	if !ok {
		t.Fatalf("expected parallel bool value, field was not a plain bool: %+v", got)
	}
	if v != *want {
		t.Fatalf("expected parallel=%v, got parallel=%v", *want, v)
	}
}

func TestConvertDeploymentServices_ParallelMapping(t *testing.T) {
	tests := []struct {
		name         string
		metadata     *v0.ServicesMetadata
		wantParallel *bool
	}{
		{
			name:         "v0 unspecified -> parallel true (preserve v0 default)",
			metadata:     nil,
			wantParallel: boolPtr(true),
		},
		{
			name:         "v0 parallel true -> parallel true",
			metadata:     &v0.ServicesMetadata{Parallel: &flexible.Field[bool]{Value: true}},
			wantParallel: boolPtr(true),
		},
		{
			name:         "v0 parallel false -> parallel false",
			metadata:     &v0.ServicesMetadata{Parallel: &flexible.Field[bool]{Value: false}},
			wantParallel: boolPtr(false),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := &v0.DeploymentServices{
				Metadata: tt.metadata,
				Values: &flexible.Field[[]*v0.DeploymentService]{
					Value: []*v0.DeploymentService{
						{ServiceRef: "svc1"},
						{ServiceRef: "svc2"},
					},
				},
			}
			got := ConvertDeploymentServices(src, NewStageConversionContext())
			if got == nil {
				t.Fatalf("expected a ServiceRef, got nil")
			}
			assertParallel(t, got.Parallel, tt.wantParallel)
		})
	}
}

func TestConvertEnvironments_ParallelMapping(t *testing.T) {
	tests := []struct {
		name         string
		metadata     *v0.EnvironmentMetadata
		wantParallel *bool
	}{
		{
			name:         "v0 unspecified -> parallel true (preserve v0 default)",
			metadata:     nil,
			wantParallel: boolPtr(true),
		},
		{
			name:         "v0 parallel true -> parallel true",
			metadata:     &v0.EnvironmentMetadata{Parallel: &flexible.Field[bool]{Value: true}},
			wantParallel: boolPtr(true),
		},
		{
			name:         "v0 parallel false -> parallel false",
			metadata:     &v0.EnvironmentMetadata{Parallel: &flexible.Field[bool]{Value: false}},
			wantParallel: boolPtr(false),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := &v0.Environments{
				Metadata: tt.metadata,
				Values: &flexible.Field[[]*v0.Environment]{
					Value: []*v0.Environment{
						{EnvironmentRef: "env1"},
						{EnvironmentRef: "env2"},
					},
				},
			}
			got := ConvertEnvironments(src, NewStageConversionContext())
			if got == nil {
				t.Fatalf("expected an EnvironmentRef, got nil")
			}
			assertParallel(t, got.Parallel, tt.wantParallel)
		})
	}
}
