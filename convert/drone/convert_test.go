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

package drone

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
)

func TestConvert(t *testing.T) {
	tests, err := filepath.Glob("testdata/examples/*.yaml")
	if err != nil {
		t.Error(err)
		return
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			orgSecrets := []string{
				"FIRST_ORG_SECRET",
				"SECOND_ORG_SECRET",
			}
			// convert the yaml file from drone to harness
			converter := New(
				WithOrgSecrets(orgSecrets...),
			)
			tmp1, err := converter.ConvertFile(test)
			if err != nil {
				t.Error(err)
				return
			}

			// unmarshal the converted yaml file to a map
			got := map[string]interface{}{}
			if err := yaml.Unmarshal(tmp1, &got); err != nil {
				t.Error(err)
				return
			}

			// parse the golden yaml file
			data, err := ioutil.ReadFile(test + ".golden")
			if err != nil {
				t.Error(err)
				return
			}

			// unmarshal the golden yaml file to a map
			want := map[string]interface{}{}
			if err := yaml.Unmarshal(data, &want); err != nil {
				t.Error(err)
				return
			}

			// compare the converted yaml to the golden file
			if diff := cmp.Diff(got, want); diff != "" {
				t.Errorf("Unexpected conversion result")
				t.Log(diff)
			}
		})
	}
}

func TestReplaceVars(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "DRONE_COMMIT_SHA",
			input:    "${DRONE_COMMIT_SHA}",
			expected: "<+codebase.commitSha>",
		},
		{
			name:     "DRONE_COMMIT_SHA without braces",
			input:    "$DRONE_COMMIT_SHA",
			expected: "<+codebase.commitSha>",
		},
		{
			name:     "CI_COMMIT_SHA",
			input:    "${CI_COMMIT_SHA}",
			expected: "<+codebase.commitSha>",
		},
		{
			name:     "CI_COMMIT_SHA without braces",
			input:    "$CI_COMMIT_SHA",
			expected: "<+codebase.commitSha>",
		},
		{
			name:     "Escaped DRONE_COMMIT_SHA",
			input:    "$${DRONE_COMMIT_SHA}",
			expected: "<+codebase.commitSha>",
		},
		{
			name:     "Escaped DRONE_COMMIT_SHA without braces",
			input:    "$$DRONE_COMMIT_SHA",
			expected: "<+codebase.commitSha>",
		},
		{
			name:     "Non-mapped variable",
			input:    "${DRONE_BUILD_ACTION}",
			expected: "${DRONE_BUILD_ACTION}",
		},
		// Known Mappings
		{
			name:     "DRONE_BRANCH",
			input:    "${DRONE_BRANCH}",
			expected: "<+codebase.branch>",
		},
		{
			name:     "DRONE_BUILD_NUMBER",
			input:    "${DRONE_BUILD_NUMBER}",
			expected: "<+pipeline.sequenceId>",
		},
		{
			name:     "DRONE_COMMIT_AUTHOR",
			input:    "${DRONE_COMMIT_AUTHOR}",
			expected: "<+codebase.gitUserId>",
		},
		{
			name:     "DRONE_COMMIT_BRANCH",
			input:    "${DRONE_COMMIT_BRANCH}",
			expected: "<+codebase.branch>",
		},
		{
			name:     "DRONE_COMMIT_SHA",
			input:    "${DRONE_COMMIT_SHA}",
			expected: "<+codebase.commitSha>",
		},
		{
			name:     "DRONE_PULL_REQUEST",
			input:    "${DRONE_PULL_REQUEST}",
			expected: "<+codebase.prNumber>",
		},
		{
			name:     "DRONE_PULL_REQUEST_TITLE",
			input:    "${DRONE_PULL_REQUEST_TITLE}",
			expected: "<+codebase.prTitle>",
		},
		{
			name:     "DRONE_REMOTE_URL",
			input:    "${DRONE_REMOTE_URL}",
			expected: "<+codebase.repoUrl>",
		},
		{
			name:     "DRONE_REPO_NAME",
			input:    "${DRONE_REPO_NAME}",
			expected: "<+<+codebase.repoUrl>.substring(<+codebase.repoUrl>.lastIndexOf('/') + 1)>",
		},
		// Unknown Mappings
		{
			name:     "DRONE_BUILD_ACTION",
			input:    "${DRONE_BUILD_ACTION}",
			expected: "${DRONE_BUILD_ACTION}", // Expect same input string for unknown mappings
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := replaceVars(tt.input)
			if output != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, output)
			}
		})
	}
}
