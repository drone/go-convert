// Copyright 2023 Harness, Inc.
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


package json

import (
	"testing"

	harness "github.com/drone/spec/dist/go"
)

func TestConvertReadTrusted(t *testing.T) {
	tests := []struct {
		name     string
		node     Node
		expected *harness.Step
		wantErr  bool
	}{
		{
			name: "Basic readTrusted",
			node: Node{
				ParameterMap: map[string]interface{}{
					"path": "path/to/file.txt",
				},
			},
			expected: &harness.Step{
				Name: "read_trusted",
				Type: "plugin",
				Spec: &harness.StepPlugin{
					Image: "plugins/read-trusted",
					With: map[string]interface{}{
						"file_path":     "path/to/file.txt",
						"trusted_branch": "main",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertReadTrusted(tt.node)
			if tt.wantErr {
				if result != nil {
					t.Errorf("Expected nil result but got %v", result)
				}
			} else {
				if result == nil {
					t.Error("Expected non-nil result but got nil")
					return
				}
				// Compare the important fields
				if result.Name != tt.expected.Name ||
					result.Type != tt.expected.Type ||
					result.Spec.(*harness.StepPlugin).Image != tt.expected.Spec.(*harness.StepPlugin).Image {
					t.Errorf("ConvertReadTrusted() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}
