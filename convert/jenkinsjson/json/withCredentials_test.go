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
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestConvertWithCredentials(t *testing.T) {
	tests := []struct {
		name      string
		symbol    string
		arguments map[string]interface{}
		want      map[string]string
	}{
		{
			name:      "ReturnsStringVariable",
			symbol:    "string",
			arguments: map[string]interface{}{"variable": "ABC"},
			want:      map[string]string{"ABC": "<+pipeline.variables.ABC>"},
		},
		{
			name:      "ReturnsUserPassVariables",
			symbol:    "usernamePassword",
			arguments: map[string]interface{}{"usernameVariable": "USER", "passwordVariable": "PASS"},
			want:      map[string]string{"USER": "<+pipeline.variables.USER>", "PASS": "<+pipeline.variables.PASS>"},
		},
		{
			name:      "LogsErrorWhenUnknownSymbol",
			symbol:    "not_real",
			arguments: map[string]interface{}{"variable": "XYZ"},
			want:      map[string]string{},
		},
		{
			name:      "ReturnsEmptyWhenArgumentsIsNull",
			symbol:    "string",
			arguments: nil,
			want:      map[string]string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			node := Node{
				ParameterMap: map[string]any{
					"bindings": []any{
						map[string]any{
							"symbol":    tc.symbol,
							"arguments": tc.arguments,
						},
					},
				},
			}
			got := ConvertWithCredentials(node)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("ConvertWithCredentials mismatch (-want +got):\n%s", diff)
			}
		})
	}

}

func TestConvertWithCredentialsArray(t *testing.T) {
	tests := []struct {
		name string
		want map[string]string
	}{
		{
			name: "ReturnsMultipleVariables",
			want: map[string]string{"A": "<+pipeline.variables.A>", "C": "<+pipeline.variables.C>", "D": "<+pipeline.variables.D>", "B": "<+pipeline.variables.B>"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			node := Node{
				ParameterMap: map[string]any{
					"bindings": []any{
						map[string]any{
							"symbol":    "string",
							"arguments": map[string]interface{}{"variable": "A"},
						},
						map[string]any{
							"symbol":    "usernamePassword",
							"arguments": map[string]interface{}{"usernameVariable": "C", "passwordVariable": "D"},
						},
						map[string]any{
							"symbol":    "string",
							"arguments": map[string]interface{}{"variable": "B"},
						},
					},
				},
			}
			got := ConvertWithCredentials(node)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("ConvertWithCredentials mismatch (-want +got):\n%s", diff)
			}
		})
	}

}

func TestConvertWithCredentialsErrorCases(t *testing.T) {
	tests := []struct {
		name         string
		parameterMap map[string]any
		want         map[string]string
	}{
		{
			name:         "HandlesNullParameterMap",
			parameterMap: nil,
			want:         map[string]string{},
		},
		{
			name:         "HandlesMissingBindings",
			parameterMap: map[string]any{},
			want:         map[string]string{},
		},
		{
			name: "HandlesBindingsAsNull",
			parameterMap: map[string]any{
				"bindings": nil,
			},
			want: map[string]string{},
		},
		{
			name: "HandlesNullBindings",
			parameterMap: map[string]any{
				"bindings": []any{nil},
			},
			want: map[string]string{},
		},
		{
			name: "HandlesMissingSymbol",
			parameterMap: map[string]any{
				"bindings": []any{
					map[string]any{
						"arguments": map[string]interface{}{"variable": "ABC"},
					},
				},
			},
			want: map[string]string{},
		},
		{
			name: "HandlesSymbolAsNull",
			parameterMap: map[string]any{
				"bindings": []any{
					map[string]any{
						"symbol":    nil,
						"arguments": map[string]interface{}{"variable": "ABC"},
					},
				},
			},
			want: map[string]string{},
		},
		{
			name: "HandlesSymbolAsWrongType",
			parameterMap: map[string]any{
				"bindings": []any{
					map[string]any{
						"symbol":    1,
						"arguments": map[string]interface{}{"variable": "ABC"},
					},
				},
			},
			want: map[string]string{},
		},
		{
			name: "HandlesMissingArguments",
			parameterMap: map[string]any{
				"bindings": []any{
					map[string]any{
						"symbol": "string",
					},
				},
			},
			want: map[string]string{},
		},
		{
			name: "HandlesArgumentsAsNull",
			parameterMap: map[string]any{
				"bindings": []any{
					map[string]any{
						"symbol":    "string",
						"arguments": nil,
					},
				},
			},
			want: map[string]string{},
		},
		{
			name: "HandlesArgumentsAsWrongType",
			parameterMap: map[string]any{
				"bindings": []any{
					map[string]any{
						"symbol":    "string",
						"arguments": "bad",
					},
				},
			},
			want: map[string]string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			node := Node{
				ParameterMap: tc.parameterMap,
			}
			got := ConvertWithCredentials(node)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("ConvertWithCredentials mismatch (-want +got):\n%s", diff)
			}
		})
	}

}
