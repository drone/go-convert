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

package downgrade

type (
	// Args defines conversion args.
	Args struct {
		// Pipeline identifier.
		ID string

		// Pipeline name.
		Name string

		// Harness organization identifier.
		Organization string

		// Harness project identifier.
		Project string

		// Docker connector reference.
		Docker Docker

		// Docker connector reference.
		Kubernetes Kubernetes

		// Codebase details.
		Codebase Codebase
	}

	// Docker connector reference.
	Docker struct {
		Connector string
	}

	// Kubernetes connector reference.
	Kubernetes struct {
		Connector string
		Namespace string
	}

	// Codebase details.
	Codebase struct {
		Connector string
		Repo      string
		Build     string
	}
)
