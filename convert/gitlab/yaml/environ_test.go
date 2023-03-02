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

func TestEnvironment(t *testing.T) {
	tests := []struct {
		yaml string
		want Environment
	}{
		{
			yaml: `"production"`,
			want: Environment{
				Name: "production",
			},
		},
		{
			yaml: `{ "name": "production", "url": "https://prod.example.com" }`,
			want: Environment{
				Name: "production",
				Url:  "https://prod.example.com",
			},
		},
		{
			yaml: `{ "name": "production", "url": "https://prod.example.com", "kubernetes": { "namespace": "production" } }`,
			want: Environment{
				Name: "production",
				Url:  "https://prod.example.com",
				Kubernetes: &Kubernetes{
					Namespace: "production",
				},
			},
		},
	}

	for i, test := range tests {
		got := new(Environment)
		if err := yaml.Unmarshal([]byte(test.yaml), got); err != nil {
			t.Error(err)
			return
		}
		if diff := cmp.Diff(got, &test.want); diff != "" {
			t.Errorf("Unexpected parsing results for test %v", i)
			t.Log(diff)
		}
	}
}

func TestEnviron_Error(t *testing.T) {
	err := yaml.Unmarshal([]byte("[]"), new(Environment))
	if err == nil || err.Error() != "failed to unmarshal environment" {
		t.Errorf("Expect error, got %s", err)
	}
}

// deploy to production:
//   stage: deploy
//   script: git push production HEAD:main
//   environment: production

// deploy to production:
//   stage: deploy
//   script: git push production HEAD:main
//   environment:
//     name: production

// deploy to production:
//   stage: deploy
//   script: git push production HEAD:main
//   environment:
//     name: production
//     url: https://prod.example.com

// stop_review_app:
//   stage: deploy
//   variables:
//     GIT_STRATEGY: none
//   script: make delete-app
//   when: manual
//   environment:
//     name: review/$CI_COMMIT_REF_SLUG
//     action: stop

// review_app:
//   script: deploy-review-app
//   environment:
//     name: review/$CI_COMMIT_REF_SLUG
//     auto_stop_in: 1 day

// deploy:
//   stage: deploy
//   script: make deploy-app
//   environment:
//     name: production
//     kubernetes:
//       namespace: production

// deploy:
//   script: echo
//   environment:
//     name: customer-portal
//     deployment_tier: production
