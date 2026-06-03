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
)

func TestRuntimeMarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		runtime  Runtime
		expected string
	}{
		{
			name:     "shell without connector marshals to string",
			runtime:  Runtime{Shell: &RuntimeShell{}},
			expected: `"shell"`,
		},
		{
			name:     "shell with connector marshals to object",
			runtime:  Runtime{Shell: &RuntimeShell{Connector: "gt_docker"}},
			expected: `{"shell":{"connector":"gt_docker"}}`,
		},
		{
			name:     "cloud runtime marshals to object",
			runtime:  Runtime{Cloud: &RuntimeCloud{Size: "large"}},
			expected: `{"cloud":{"size":"large"}}`,
		},
		{
			name:     "kubernetes runtime marshals to object",
			runtime:  Runtime{Kubernetes: &RuntimeKubernetes{Namespace: "ci"}},
			expected: `{"kubernetes":{"namespace":"ci"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Marshal(tt.runtime)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(result) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(result))
			}
		})
	}
}

func TestRuntimeUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectShell    bool
		expectConnector string
		expectCloud    bool
		expectK8s      bool
	}{
		{
			name:        "string 'shell' unmarshals to shell runtime",
			input:       `"shell"`,
			expectShell: true,
		},
		{
			name:            "object with shell.connector unmarshals correctly",
			input:           `{"shell":{"connector":"gt_docker"}}`,
			expectShell:     true,
			expectConnector: "gt_docker",
		},
		{
			name:        "string 'cloud' unmarshals to cloud runtime",
			input:       `"cloud"`,
			expectCloud: true,
		},
		{
			name:      "string 'kubernetes' unmarshals to kubernetes runtime",
			input:     `"kubernetes"`,
			expectK8s: true,
		},
		{
			name:        "object with cloud unmarshals correctly",
			input:       `{"cloud":{"size":"large"}}`,
			expectCloud: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var runtime Runtime
			err := json.Unmarshal([]byte(tt.input), &runtime)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.expectShell {
				if runtime.Shell == nil {
					t.Error("expected Shell to be non-nil")
				} else if tt.expectConnector != "" && runtime.Shell.Connector != tt.expectConnector {
					t.Errorf("expected connector %q, got %q", tt.expectConnector, runtime.Shell.Connector)
				}
			}

			if tt.expectCloud && runtime.Cloud == nil {
				t.Error("expected Cloud to be non-nil")
			}

			if tt.expectK8s && runtime.Kubernetes == nil {
				t.Error("expected Kubernetes to be non-nil")
			}
		})
	}
}

func TestRuntimeIsShell(t *testing.T) {
	tests := []struct {
		name     string
		runtime  *Runtime
		expected bool
	}{
		{
			name:     "nil runtime returns false",
			runtime:  nil,
			expected: false,
		},
		{
			name:     "shell runtime returns true",
			runtime:  &Runtime{Shell: &RuntimeShell{}},
			expected: true,
		},
		{
			name:     "cloud runtime returns false",
			runtime:  &Runtime{Cloud: &RuntimeCloud{}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.runtime.IsShell()
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestRuntimeRoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		runtime Runtime
	}{
		{
			name:    "shell without connector",
			runtime: Runtime{Shell: &RuntimeShell{}},
		},
		{
			name:    "shell with connector",
			runtime: Runtime{Shell: &RuntimeShell{Connector: "my-connector"}},
		},
		{
			name:    "cloud runtime",
			runtime: Runtime{Cloud: &RuntimeCloud{Size: "large"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.runtime)
			if err != nil {
				t.Fatalf("marshal error: %v", err)
			}

			var result Runtime
			err = json.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("unmarshal error: %v", err)
			}

			// For shell runtime, verify connector is preserved
			if tt.runtime.Shell != nil {
				if result.Shell == nil {
					t.Error("expected Shell to be non-nil after round-trip")
				} else if result.Shell.Connector != tt.runtime.Shell.Connector {
					t.Errorf("connector mismatch: expected %q, got %q",
						tt.runtime.Shell.Connector, result.Shell.Connector)
				}
			}

			// For cloud runtime, verify size and image are preserved
			if tt.runtime.Cloud != nil {
				if result.Cloud == nil {
					t.Error("expected Cloud to be non-nil after round-trip")
				} else {
					if result.Cloud.Size != tt.runtime.Cloud.Size {
						t.Errorf("size mismatch: expected %q, got %q",
							tt.runtime.Cloud.Size, result.Cloud.Size)
					}
					if result.Cloud.Image != tt.runtime.Cloud.Image {
						t.Errorf("image mismatch: expected %q, got %q",
							tt.runtime.Cloud.Image, result.Cloud.Image)
					}
				}
			}
		})
	}
}
