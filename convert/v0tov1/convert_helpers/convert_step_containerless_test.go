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
	"strings"
	"testing"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	"github.com/drone/go-convert/convert/v0tov1/messagelog"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/drone/go-convert/internal/flexible"
)

// Cloud runtime context: container-block should be suppressed when the v0
// spec has no image. This test suite covers each of the seven step helpers
// that emit container blocks, plus the ConvertServiceDependency helper.

const testLogFile = "containerless_test"

// cloudCtx builds a StepConvertContext whose runtime is Cloud (hosted VM).
func cloudCtx() *StepConvertContext {
	return &StepConvertContext{Runtime: &v1.Runtime{Cloud: &v1.RuntimeCloud{}}}
}

// k8sCtx builds a StepConvertContext whose runtime is Kubernetes (containers required).
func k8sCtx() *StepConvertContext {
	return &StepConvertContext{Runtime: &v1.Runtime{Kubernetes: &v1.RuntimeKubernetes{}}}
}

// resetLog wipes and re-enables the message logger, returning a cleanup func.
func resetLog(t *testing.T) func() {
	t.Helper()
	messagelog.ResetMessageLogger()
	messagelog.GetMessageLogger().Enable("")
	messagelog.GetMessageLogger().SetCurrentFile(testLogFile)
	return func() { messagelog.ResetMessageLogger() }
}

// assertContainerlessWarn asserts exactly one WARN with the containerless code
// was recorded and that its message names the expected dropped fields.
func assertContainerlessWarn(t *testing.T, wantFields ...string) {
	t.Helper()
	fl := messagelog.GetMessageLogger().GetFileLog(testLogFile)
	if fl == nil {
		t.Fatalf("no file log recorded")
	}
	var warns []messagelog.Message
	for _, m := range fl.Messages {
		if m.Code == "CLOUD_CONTAINERLESS_FIELDS_DROPPED" {
			warns = append(warns, m)
		}
	}
	if len(warns) != 1 {
		t.Fatalf("expected 1 CLOUD_CONTAINERLESS_FIELDS_DROPPED warn, got %d (messages: %+v)", len(warns), fl.Messages)
	}
	for _, f := range wantFields {
		if !strings.Contains(warns[0].Message, f) {
			t.Errorf("warn message %q missing expected field %q", warns[0].Message, f)
		}
	}
}

// assertNoContainerlessWarn asserts no containerless-drop warning was recorded.
func assertNoContainerlessWarn(t *testing.T) {
	t.Helper()
	fl := messagelog.GetMessageLogger().GetFileLog(testLogFile)
	if fl == nil {
		return
	}
	for _, m := range fl.Messages {
		if m.Code == "CLOUD_CONTAINERLESS_FIELDS_DROPPED" {
			t.Fatalf("did not expect CLOUD_CONTAINERLESS_FIELDS_DROPPED, got %+v", m)
		}
	}
}

// --- ConvertStepRun ----------------------------------------------------------

func TestConvertStepRun_CloudContainerlessDropsBlock(t *testing.T) {
	defer resetLog(t)()
	step := &v0.Step{
		ID:   "run1",
		Type: v0.StepTypeRun,
		Spec: &v0.StepRun{
			Command:    "echo hi",
			Resources:  &v0.Resources{Limits: &v0.ResourceSpec{CPU: &flexible.Field[*v0.MilliSize]{Value: "1"}}},
			ConnRef:    "docker-connector",
			Privileged: &flexible.Field[bool]{Value: true},
			RunAsUser:  &flexible.Field[int]{Value: 1000},
		},
	}
	got := ConvertStepRun(step, cloudCtx())
	if got == nil {
		t.Fatal("expected non-nil result")
	}
	if got.Container != nil {
		t.Errorf("expected Container=nil on Cloud+no-image, got %+v", got.Container)
	}
	if len(got.Script) == 0 || got.Script[0] != "echo hi" {
		t.Errorf("script not preserved: %+v", got.Script)
	}
	assertContainerlessWarn(t, "resources", "connectorRef", "privileged", "runAsUser")
}

func TestConvertStepRun_CloudWithImagePreservesContainer(t *testing.T) {
	defer resetLog(t)()
	step := &v0.Step{
		ID:   "run2",
		Type: v0.StepTypeRun,
		Spec: &v0.StepRun{
			Image:   "alpine",
			Command: "echo hi",
		},
	}
	got := ConvertStepRun(step, cloudCtx())
	if got == nil || got.Container == nil || got.Container.Image != "alpine" {
		t.Fatalf("expected Container.Image=alpine on Cloud+image, got %+v", got)
	}
	assertNoContainerlessWarn(t)
}

func TestConvertStepRun_NonCloudPreservesTodayBehavior(t *testing.T) {
	defer resetLog(t)()
	step := &v0.Step{
		ID:   "run3",
		Type: v0.StepTypeRun,
		Spec: &v0.StepRun{
			Command:   "echo hi",
			Resources: &v0.Resources{Limits: &v0.ResourceSpec{CPU: &flexible.Field[*v0.MilliSize]{Value: "1"}}},
		},
	}
	// nil ctx and k8s ctx both must preserve today's behavior: container emitted.
	for _, name := range []string{"nil", "k8s"} {
		var ctx *StepConvertContext
		if name == "k8s" {
			ctx = k8sCtx()
		}
		t.Run(name, func(t *testing.T) {
			got := ConvertStepRun(step, ctx)
			if got == nil || got.Container == nil {
				t.Fatalf("expected Container preserved on non-Cloud, got %+v", got)
			}
			if got.Container.Image != "" {
				t.Errorf("expected empty image in container, got %q", got.Container.Image)
			}
		})
	}
	assertNoContainerlessWarn(t)
}

// --- ConvertStepBackground ---------------------------------------------------

func TestConvertStepBackground_CloudContainerlessDropsBlock(t *testing.T) {
	defer resetLog(t)()
	step := &v0.Step{
		ID:   "bg1",
		Type: v0.StepTypeBackground,
		Spec: &v0.StepBackground{
			Command:      "redis-server",
			Resources:    &v0.Resources{Limits: &v0.ResourceSpec{CPU: &flexible.Field[*v0.MilliSize]{Value: "500m"}}},
			Entrypoint:   &flexible.Field[[]string]{Value: []string{"/bin/sh"}},
			PortBindings: &flexible.Field[map[string]string]{Value: map[string]string{"6379": "6379"}},
		},
	}
	got := ConvertStepBackground(step, cloudCtx())
	if got == nil || got.Container != nil {
		t.Fatalf("expected Container=nil on Cloud+no-image, got %+v", got)
	}
	assertContainerlessWarn(t, "resources", "entrypoint", "portBindings")
}

func TestConvertStepBackground_NonCloudPreservesTodayBehavior(t *testing.T) {
	defer resetLog(t)()
	step := &v0.Step{
		ID:   "bg2",
		Type: v0.StepTypeBackground,
		Spec: &v0.StepBackground{
			Resources: &v0.Resources{Limits: &v0.ResourceSpec{CPU: &flexible.Field[*v0.MilliSize]{Value: "500m"}}},
		},
	}
	got := ConvertStepBackground(step, nil)
	if got == nil || got.Container == nil {
		t.Fatalf("expected Container preserved on nil-ctx, got %+v", got)
	}
	assertNoContainerlessWarn(t)
}

// --- ConvertStepPlugin -------------------------------------------------------

func TestConvertStepPlugin_CloudContainerlessDropsBlock(t *testing.T) {
	defer resetLog(t)()
	step := &v0.Step{
		ID:   "plugin1",
		Type: v0.StepTypePlugin,
		Spec: &v0.StepPlugin{
			ConnRef:    "connector",
			Resources:  &v0.Resources{Limits: &v0.ResourceSpec{CPU: &flexible.Field[*v0.MilliSize]{Value: "1"}}},
			Entrypoint: &flexible.Field[[]string]{Value: []string{"/plugin"}},
			RunAsUser:  &flexible.Field[int]{Value: 999},
		},
	}
	got := ConvertStepPlugin(step, cloudCtx())
	if got == nil || got.Container != nil {
		t.Fatalf("expected Container=nil on Cloud+no-image, got %+v", got)
	}
	assertContainerlessWarn(t, "connectorRef", "resources", "entrypoint", "runAsUser")
}

func TestConvertStepPlugin_NonCloudPreservesTodayBehavior(t *testing.T) {
	defer resetLog(t)()
	step := &v0.Step{
		ID:   "plugin2",
		Type: v0.StepTypePlugin,
		Spec: &v0.StepPlugin{
			Resources: &v0.Resources{Limits: &v0.ResourceSpec{CPU: &flexible.Field[*v0.MilliSize]{Value: "1"}}},
		},
	}
	got := ConvertStepPlugin(step, k8sCtx())
	if got == nil || got.Container == nil {
		t.Fatalf("expected Container preserved on K8s, got %+v", got)
	}
	assertNoContainerlessWarn(t)
}

// --- ConvertStepRunTests -----------------------------------------------------

func TestConvertStepRunTests_CloudContainerlessDropsBlock(t *testing.T) {
	defer resetLog(t)()
	step := &v0.Step{
		ID:   "runtests1",
		Type: v0.StepTypeRunTests,
		Spec: &v0.StepRunTests{
			Language:     "java",
			BuildTool:    "maven",
			ConnectorRef: "docker-conn",
			Resources:    &v0.Resources{Limits: &v0.ResourceSpec{CPU: &flexible.Field[*v0.MilliSize]{Value: "2"}}},
			RunAsUser:    &flexible.Field[int]{Value: 1001},
		},
	}
	got := ConvertStepRunTests(step, cloudCtx())
	if got == nil || got.Container != nil {
		t.Fatalf("expected Container=nil on Cloud+no-image, got %+v", got)
	}
	assertContainerlessWarn(t, "connectorRef", "resources", "runAsUser")
}

func TestConvertStepRunTests_NonCloudPreservesTodayBehavior(t *testing.T) {
	defer resetLog(t)()
	step := &v0.Step{
		ID:   "runtests2",
		Type: v0.StepTypeRunTests,
		Spec: &v0.StepRunTests{
			Language:  "java",
			BuildTool: "maven",
			Resources: &v0.Resources{Limits: &v0.ResourceSpec{CPU: &flexible.Field[*v0.MilliSize]{Value: "2"}}},
		},
	}
	got := ConvertStepRunTests(step, nil)
	if got == nil || got.Container == nil {
		t.Fatalf("expected Container preserved on nil-ctx, got %+v", got)
	}
	assertNoContainerlessWarn(t)
}

// --- ConvertStepTestIntelligence ---------------------------------------------

func TestConvertStepTestIntelligence_CloudContainerlessDropsBlock(t *testing.T) {
	defer resetLog(t)()
	step := &v0.Step{
		ID:   "ti1",
		Type: v0.StepTypeTest,
		Spec: &v0.StepTestIntelligence{
			Command:   "pytest",
			ConnRef:   "conn",
			Resources: &v0.Resources{Limits: &v0.ResourceSpec{CPU: &flexible.Field[*v0.MilliSize]{Value: "1"}}},
		},
	}
	got := ConvertStepTestIntelligence(step, cloudCtx())
	if got == nil || got.Container != nil {
		t.Fatalf("expected Container=nil on Cloud+no-image, got %+v", got)
	}
	assertContainerlessWarn(t, "connectorRef", "resources")
}

func TestConvertStepTestIntelligence_NonCloudPreservesTodayBehavior(t *testing.T) {
	defer resetLog(t)()
	step := &v0.Step{
		ID:   "ti2",
		Type: v0.StepTypeTest,
		Spec: &v0.StepTestIntelligence{
			Command:   "pytest",
			Resources: &v0.Resources{Limits: &v0.ResourceSpec{CPU: &flexible.Field[*v0.MilliSize]{Value: "1"}}},
		},
	}
	got := ConvertStepTestIntelligence(step, nil)
	if got == nil || got.Container == nil {
		t.Fatalf("expected Container preserved on nil-ctx, got %+v", got)
	}
	assertNoContainerlessWarn(t)
}

// --- ConvertStepContainer ----------------------------------------------------

func TestConvertStepContainer_CloudContainerlessDropsBlock(t *testing.T) {
	defer resetLog(t)()
	step := &v0.Step{
		ID:   "container1",
		Type: v0.StepTypeContainer,
		Spec: &v0.StepContainer{
			Command:    "echo hi",
			ConnRef:    "conn",
			Privileged: &flexible.Field[bool]{Value: true},
			Entrypoint: &flexible.Field[[]string]{Value: []string{"/bin/sh"}},
		},
	}
	got := ConvertStepContainer(step, cloudCtx())
	if got == nil || got.Container != nil {
		t.Fatalf("expected Container=nil on Cloud+no-image, got %+v", got)
	}
	assertContainerlessWarn(t, "connectorRef", "privileged", "entrypoint")
}

func TestConvertStepContainer_NonCloudPreservesTodayBehavior(t *testing.T) {
	defer resetLog(t)()
	step := &v0.Step{
		ID:   "container2",
		Type: v0.StepTypeContainer,
		Spec: &v0.StepContainer{
			Command:    "echo hi",
			Privileged: &flexible.Field[bool]{Value: true},
		},
	}
	got := ConvertStepContainer(step, nil)
	if got == nil || got.Container == nil {
		t.Fatalf("expected Container preserved on nil-ctx, got %+v", got)
	}
	assertNoContainerlessWarn(t)
}

// --- ConvertServiceDependencyToBackgroundStep --------------------------------

func TestConvertServiceDependency_CloudContainerlessDropsBlock(t *testing.T) {
	defer resetLog(t)()
	svc := &v0.Service{
		ID: "svc1",
		Spec: &v0.ServiceSpec{
			Conn:      "connector",
			Resources: &v0.Resources{Limits: &v0.ResourceSpec{CPU: &flexible.Field[*v0.MilliSize]{Value: "1"}}},
		},
	}
	got := ConvertServiceDependencyToBackgroundStep(svc, cloudCtx())
	if got == nil {
		t.Fatal("expected non-nil result")
	}
	if got.Background == nil || got.Background.Container != nil {
		t.Fatalf("expected Background.Container=nil on Cloud+no-image, got %+v", got.Background)
	}
	assertContainerlessWarn(t, "connector", "resources")
}

func TestConvertServiceDependency_NonCloudPreservesTodayBehavior(t *testing.T) {
	defer resetLog(t)()
	svc := &v0.Service{
		ID: "svc2",
		Spec: &v0.ServiceSpec{
			Conn:      "connector",
			Resources: &v0.Resources{Limits: &v0.ResourceSpec{CPU: &flexible.Field[*v0.MilliSize]{Value: "1"}}},
		},
	}
	got := ConvertServiceDependencyToBackgroundStep(svc, nil)
	if got == nil || got.Background == nil || got.Background.Container == nil {
		t.Fatalf("expected Background.Container preserved on nil-ctx, got %+v", got)
	}
	assertNoContainerlessWarn(t)
}

// --- StepConvertContext.IsCloud (nil-safety) ---------------------------------

func TestStepConvertContext_IsCloud_NilSafe(t *testing.T) {
	var ctx *StepConvertContext
	if ctx.IsCloud() {
		t.Fatal("nil ctx must not be Cloud")
	}
	if (&StepConvertContext{}).IsCloud() {
		t.Fatal("ctx with nil runtime must not be Cloud")
	}
	if (&StepConvertContext{Runtime: &v1.Runtime{}}).IsCloud() {
		t.Fatal("ctx with runtime but no Cloud subfield must not be Cloud")
	}
	if !cloudCtx().IsCloud() {
		t.Fatal("cloudCtx() must be Cloud")
	}
}
