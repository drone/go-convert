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

package travis

import "testing"

func TestOptions(t *testing.T) {
	p := New(
		WithDockerhub("account.docker"),
		WithKubernetes("namespace", "connector.kubernetes"),
	)

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
}
