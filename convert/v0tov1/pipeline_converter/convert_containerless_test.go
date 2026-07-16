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

package pipelineconverter

import (
	"testing"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	"github.com/drone/go-convert/internal/flexible"
)

// Integration-level assertion: on Cloud runtime, a v0 Run step without an
// image must produce a v1 Run step with no container block, preserving the
// containerless execution mode. This exercises the full stage → step wiring
// (runtime resolved before steps, ctx threaded, guard applied).
func TestConvertCIStage_CloudRuntime_ContainerlessRunStep(t *testing.T) {
	converter := NewPipelineConverter()

	pipeline := &v0.Pipeline{
		ID:   "p",
		Name: "p",
		Stages: []*v0.Stages{
			{Stage: &v0.Stage{
				ID:   "ci",
				Name: "ci",
				Type: v0.StageTypeCI,
				Spec: &v0.StageCI{
					Runtime: &v0.Runtime{
						Type: "Cloud",
						Spec: &v0.RuntimeCloudSpec{Size: "Standard"},
					},
					Execution: v0.Execution{
						Steps: []*v0.Steps{
							{Step: &v0.Step{
								ID:   "r1",
								Name: "r1",
								Type: v0.StepTypeRun,
								Spec: &v0.StepRun{
									Command:    "echo hi",
									Resources:  &v0.Resources{Limits: &v0.ResourceSpec{CPU: &flexible.Field[*v0.MilliSize]{Value: "1"}}},
									RunAsUser:  &flexible.Field[int]{Value: 1000},
								},
							}},
						},
					},
				},
			}},
		},
	}

	got := converter.ConvertPipeline(pipeline)
	if got == nil || len(got.Stages) != 1 {
		t.Fatalf("expected 1 stage, got %+v", got)
	}
	stage := got.Stages[0]
	if stage.Runtime == nil || stage.Runtime.Cloud == nil {
		t.Fatalf("expected runtime.cloud, got %+v", stage.Runtime)
	}
	if len(stage.Steps) != 1 {
		t.Fatalf("expected 1 step, got %d", len(stage.Steps))
	}
	run := stage.Steps[0].Run
	if run == nil {
		t.Fatalf("expected Run spec, got %+v", stage.Steps[0])
	}
	if run.Container != nil {
		t.Errorf("expected containerless (nil Container) on Cloud+no-image, got %+v", run.Container)
	}
	if len(run.Script) == 0 || run.Script[0] != "echo hi" {
		t.Errorf("expected script preserved, got %+v", run.Script)
	}
}

// Cloud pipelines with platform.os set (a common shape — see CI-23632
// example YAML) must emit platform.arch=amd64 by default, otherwise
// PlatformV1.toPlatform() NPEs at plan-creation time on the V1 backend.
func TestConvertCIStage_CloudRuntime_PlatformArchDefaultsToAmd64(t *testing.T) {
	converter := NewPipelineConverter()

	pipeline := &v0.Pipeline{
		ID:   "p",
		Name: "p",
		Stages: []*v0.Stages{
			{Stage: &v0.Stage{
				ID:   "ci",
				Name: "ci",
				Type: v0.StageTypeCI,
				Spec: &v0.StageCI{
					Platform: &v0.Platform{OS: "Linux"},
					Runtime:  &v0.Runtime{Type: "Cloud", Spec: &v0.RuntimeCloudSpec{}},
					Execution: v0.Execution{
						Steps: []*v0.Steps{
							{Step: &v0.Step{
								ID:   "r1",
								Name: "r1",
								Type: v0.StepTypeRun,
								Spec: &v0.StepRun{Command: "echo hi"},
							}},
						},
					},
				},
			}},
		},
	}

	got := converter.ConvertPipeline(pipeline)
	if got == nil || len(got.Stages) != 1 {
		t.Fatalf("expected 1 stage, got %+v", got)
	}
	stage := got.Stages[0]
	if stage.Platform == nil {
		t.Fatalf("expected platform to be set, got nil")
	}
	if stage.Platform.Os != "linux" {
		t.Errorf("expected os=linux, got %q", stage.Platform.Os)
	}
	if stage.Platform.Arch != "amd64" {
		t.Errorf("expected arch=amd64 default, got %q", stage.Platform.Arch)
	}
}

// Regression sentinel: same step body on Kubernetes infra must still emit a
// container block (K8s behavior is byte-for-byte preserved).
func TestConvertCIStage_KubernetesDirect_ContainerFieldsPreserved(t *testing.T) {
	converter := NewPipelineConverter()

	pipeline := &v0.Pipeline{
		ID:   "p",
		Name: "p",
		Stages: []*v0.Stages{
			{Stage: &v0.Stage{
				ID:   "ci",
				Name: "ci",
				Type: v0.StageTypeCI,
				Spec: &v0.StageCI{
					Infrastructure: &v0.Infrastructure{
						Type: "KubernetesDirect",
						Spec: &v0.InfrastructureKubernetesDirectSpec{
							Conn:      "k8s-conn",
							Namespace: "default",
						},
					},
					Execution: v0.Execution{
						Steps: []*v0.Steps{
							{Step: &v0.Step{
								ID:   "r1",
								Name: "r1",
								Type: v0.StepTypeRun,
								Spec: &v0.StepRun{
									Command:   "echo hi",
									Resources: &v0.Resources{Limits: &v0.ResourceSpec{CPU: &flexible.Field[*v0.MilliSize]{Value: "1"}}},
									RunAsUser: &flexible.Field[int]{Value: 1000},
								},
							}},
						},
					},
				},
			}},
		},
	}

	got := converter.ConvertPipeline(pipeline)
	if got == nil || len(got.Stages) != 1 {
		t.Fatalf("expected 1 stage, got %+v", got)
	}
	stage := got.Stages[0]
	if stage.Runtime == nil || stage.Runtime.Kubernetes == nil {
		t.Fatalf("expected runtime.kubernetes, got %+v", stage.Runtime)
	}
	run := stage.Steps[0].Run
	if run == nil || run.Container == nil {
		t.Fatalf("expected Container preserved on K8s, got %+v", run)
	}
	if run.Container.Image != "" {
		t.Errorf("expected empty image (today's behavior) in container, got %q", run.Container.Image)
	}
}
