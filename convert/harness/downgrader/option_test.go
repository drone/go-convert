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

package downgrader

import "testing"

func TestOptions(t *testing.T) {
	p := New(
		WithCodebase("drone", "connector.github"),
		WithDockerhub("account.docker"),
		WithKubernetes("namespace", "connector.kubernetes"),
		WithIdentifier("foo"),
		WithName("bar"),
		WithOrganization("baz"),
		WithProject("qux"),
	)
	if got, want := p.codebaseConn, "connector.github"; got != want {
		t.Errorf("Want codebase connector %q, got %q", want, got)
	}
	if got, want := p.codebaseName, "drone"; got != want {
		t.Errorf("Want codebase name %q, got %q", want, got)
	}
	if got, want := p.kubeConnector, "connector.kubernetes"; got != want {
		t.Errorf("Want kubernetes connector %q, got %q", want, got)
	}
	if got, want := p.kubeNamespace, "namespace"; got != want {
		t.Errorf("Want kubernetes namespace %q, got %q", want, got)
	}
	if got, want := p.kubeEnabled, true; got != want {
		t.Errorf("Want kubernetes enabled %v, got %v", want, got)
	}
	if got, want := p.dockerhubConn, "account.docker"; got != want {
		t.Errorf("Want docker connector %q, got %q", want, got)
	}
	if got, want := p.pipelineId, "foo"; got != want {
		t.Errorf("Want pipeline id %q, got %q", want, got)
	}
	if got, want := p.pipelineName, "bar"; got != want {
		t.Errorf("Want pipeline name %q, got %q", want, got)
	}
	if got, want := p.pipelineOrg, "baz"; got != want {
		t.Errorf("Want pipeline org %q, got %q", want, got)
	}
	if got, want := p.pipelineProj, "qux"; got != want {
		t.Errorf("Want pipeline project %q, got %q", want, got)
	}
}

func TestOptions_Defaults(t *testing.T) {
	p := New()

	if got, want := p.kubeConnector, ""; got != want {
		t.Errorf("Want kubernetes connector %q, got %q", want, got)
	}
	if got, want := p.kubeNamespace, "default"; got != want {
		t.Errorf("Want kubernetes namespace %q, got %q", want, got)
	}
	if got, want := p.kubeEnabled, false; got != want {
		t.Errorf("Want kubernetes enabled %v, got %v", want, got)
	}
	if got, want := p.dockerhubConn, ""; got != want {
		t.Errorf("Want docker connector %q, got %q", want, got)
	}
	if got, want := p.pipelineName, "default"; got != want {
		t.Errorf("Want pipeline name %q, got %q", want, got)
	}
	if got, want := p.pipelineId, "default"; got != want {
		t.Errorf("Want pipeline id %q, got %q", want, got)
	}
	if got, want := p.pipelineOrg, "default"; got != want {
		t.Errorf("Want pipeline org %q, got %q", want, got)
	}
	if got, want := p.pipelineProj, "default"; got != want {
		t.Errorf("Want pipeline project %q, got %q", want, got)
	}
}
