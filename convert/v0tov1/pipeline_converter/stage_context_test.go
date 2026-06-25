package pipelineconverter

import (
	"testing"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	convert_helpers "github.com/drone/go-convert/convert/v0tov1/convert_helpers"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/drone/go-convert/internal/flexible"
)

func TestStageConversionContext_SetAndGet(t *testing.T) {
	ctx := convert_helpers.NewStageConversionContext()

	svc := &v1.ServiceRef{Items: []*v1.ServiceItem{{Id: "my-service"}}}
	env := &v1.EnvironmentRef{Items: []*v1.EnvironmentItem{{Id: "my-env"}}}
	rt := &v1.Runtime{Kubernetes: &v1.RuntimeKubernetes{Namespace: "default"}}

	ctx.Set("Stage1", &convert_helpers.StageConvertedData{
		Service:     svc,
		Environment: env,
		Runtime:     rt,
	})

	if got := ctx.GetService("Stage1"); got != svc {
		t.Errorf("GetService: expected %v, got %v", svc, got)
	}
	if got := ctx.GetEnvironment("Stage1"); got != env {
		t.Errorf("GetEnvironment: expected %v, got %v", env, got)
	}
	if got := ctx.GetRuntime("Stage1"); got != rt {
		t.Errorf("GetRuntime: expected %v, got %v", rt, got)
	}
}

func TestStageConversionContext_MissingStage(t *testing.T) {
	ctx := convert_helpers.NewStageConversionContext()

	if got := ctx.GetService("NonExistent"); got != nil {
		t.Errorf("expected nil for missing stage, got %v", got)
	}
	if got := ctx.GetEnvironment("NonExistent"); got != nil {
		t.Errorf("expected nil for missing stage, got %v", got)
	}
	if got := ctx.GetRuntime("NonExistent"); got != nil {
		t.Errorf("expected nil for missing stage, got %v", got)
	}
}

func TestUseFromStage_DeploymentServiceAndEnvironment(t *testing.T) {
	converter := NewPipelineConverter()

	// First stage: Dashboards - defines service and environment
	dashboardsStage := &v0.Stage{
		ID:   "Dashboards",
		Name: "Dashboards",
		Type: v0.StageTypeDeployment,
		Spec: &v0.StageDeployment{
			Service: &v0.DeploymentService{
				ServiceRef: "dashboard-svc",
			},
			Environment: &v0.Environment{
				EnvironmentRef: "prod-env",
				DeployToAll:    &flexible.Field[bool]{Value: true},
			},
			Execution: &v0.DeploymentExecution{},
		},
	}

	// Second stage: uses service and environment from Dashboards
	secondStage := &v0.Stage{
		ID:   "SecondDeploy",
		Name: "SecondDeploy",
		Type: v0.StageTypeDeployment,
		Spec: &v0.StageDeployment{
			Service: &v0.DeploymentService{
				UseFromStage: &v0.UseFromStage{Stage: "Dashboards"},
			},
			Environment: &v0.Environment{
				UseFromStage: &v0.UseFromStage{Stage: "Dashboards"},
			},
			Execution: &v0.DeploymentExecution{},
		},
	}

	v0Stages := []*v0.Stages{
		{Stage: dashboardsStage},
		{Stage: secondStage},
	}

	v1Stages := converter.convertStages(v0Stages, "pipeline")

	if len(v1Stages) != 2 {
		t.Fatalf("expected 2 stages, got %d", len(v1Stages))
	}

	// Verify first stage has its own service and environment
	first := v1Stages[0]
	if first.Service == nil {
		t.Fatal("first stage: service should not be nil")
	}
	if len(first.Service.Items) != 1 || first.Service.Items[0].Id != "dashboard-svc" {
		t.Errorf("first stage: expected service 'dashboard-svc', got %v", first.Service.Items)
	}
	if first.Environment == nil {
		t.Fatal("first stage: environment should not be nil")
	}
	if len(first.Environment.Items) != 1 || first.Environment.Items[0].Id != "prod-env" {
		t.Errorf("first stage: expected environment 'prod-env', got %v", first.Environment.Items)
	}

	// Verify second stage copied service and environment from first
	second := v1Stages[1]
	if second.Service == nil {
		t.Fatal("second stage: service should not be nil (useFromStage)")
	}
	if len(second.Service.Items) != 1 || second.Service.Items[0].Id != "dashboard-svc" {
		t.Errorf("second stage: expected service 'dashboard-svc' from useFromStage, got %v", second.Service.Items)
	}
	if second.Environment == nil {
		t.Fatal("second stage: environment should not be nil (useFromStage)")
	}
	if len(second.Environment.Items) != 1 || second.Environment.Items[0].Id != "prod-env" {
		t.Errorf("second stage: expected environment 'prod-env' from useFromStage, got %v", second.Environment.Items)
	}
}

func TestUseFromStage_CIInfrastructure(t *testing.T) {
	converter := NewPipelineConverter()

	// First stage: Lint_Test - defines infrastructure
	lintStage := &v0.Stage{
		ID:   "Lint_Test",
		Name: "Lint Test",
		Type: v0.StageTypeCI,
		Spec: &v0.StageCI{
			Infrastructure: &v0.Infrastructure{
				Type: "KubernetesDirect",
				Spec: &v0.InfrastructureKubernetesDirectSpec{
					Namespace: "harness-build",
					Conn:      "k8s-connector",
				},
			},
			Execution: v0.Execution{},
		},
	}

	// Second stage: uses infrastructure from Lint_Test
	buildStage := &v0.Stage{
		ID:   "Build",
		Name: "Build",
		Type: v0.StageTypeCI,
		Spec: &v0.StageCI{
			Infrastructure: &v0.Infrastructure{
				Type: "UseFromStage",
				From: "Lint_Test",
			},
			Execution: v0.Execution{},
		},
	}

	v0Stages := []*v0.Stages{
		{Stage: lintStage},
		{Stage: buildStage},
	}

	v1Stages := converter.convertStages(v0Stages, "pipeline")

	if len(v1Stages) != 2 {
		t.Fatalf("expected 2 stages, got %d", len(v1Stages))
	}

	// Verify first stage has its own runtime
	first := v1Stages[0]
	if first.Runtime == nil {
		t.Fatal("first stage: runtime should not be nil")
	}
	if first.Runtime.Kubernetes == nil {
		t.Fatal("first stage: runtime.kubernetes should not be nil")
	}
	if first.Runtime.Kubernetes.Namespace != "harness-build" {
		t.Errorf("first stage: expected namespace 'harness-build', got %s", first.Runtime.Kubernetes.Namespace)
	}

	// Verify second stage copied runtime from first
	second := v1Stages[1]
	if second.Runtime == nil {
		t.Fatal("second stage: runtime should not be nil (useFromStage)")
	}
	if second.Runtime.Kubernetes == nil {
		t.Fatal("second stage: runtime.kubernetes should not be nil (useFromStage)")
	}
	if second.Runtime.Kubernetes.Namespace != "harness-build" {
		t.Errorf("second stage: expected namespace 'harness-build' from useFromStage, got %s", second.Runtime.Kubernetes.Namespace)
	}
	if second.Runtime.Kubernetes.Connector != "k8s-connector" {
		t.Errorf("second stage: expected connector 'k8s-connector' from useFromStage, got %s", second.Runtime.Kubernetes.Connector)
	}
}

func TestUseFromStage_MissingReferenceDoesNotPanic(t *testing.T) {
	converter := NewPipelineConverter()

	// Stage referencing a non-existent stage - should produce a warning, not panic
	stage := &v0.Stage{
		ID:   "Deploy",
		Name: "Deploy",
		Type: v0.StageTypeDeployment,
		Spec: &v0.StageDeployment{
			Service: &v0.DeploymentService{
				UseFromStage: &v0.UseFromStage{Stage: "NonExistent"},
			},
			Environment: &v0.Environment{
				UseFromStage: &v0.UseFromStage{Stage: "NonExistent"},
			},
			Execution: &v0.DeploymentExecution{},
		},
	}

	v0Stages := []*v0.Stages{{Stage: stage}}
	v1Stages := converter.convertStages(v0Stages, "pipeline")

	if len(v1Stages) != 1 {
		t.Fatalf("expected 1 stage, got %d", len(v1Stages))
	}

	// Service and environment should be nil since the referenced stage doesn't exist
	if v1Stages[0].Service != nil {
		t.Errorf("expected nil service for missing useFromStage, got %v", v1Stages[0].Service)
	}
	if v1Stages[0].Environment != nil {
		t.Errorf("expected nil environment for missing useFromStage, got %v", v1Stages[0].Environment)
	}
}

func TestUseFromStage_ServiceOnlyFromPreviousStage(t *testing.T) {
	converter := NewPipelineConverter()

	// First stage defines both service and environment
	firstStage := &v0.Stage{
		ID:   "Stage1",
		Name: "Stage1",
		Type: v0.StageTypeDeployment,
		Spec: &v0.StageDeployment{
			Service: &v0.DeploymentService{
				ServiceRef: "svc-alpha",
			},
			Environment: &v0.Environment{
				EnvironmentRef: "env-beta",
			},
			Execution: &v0.DeploymentExecution{},
		},
	}

	// Second stage: useFromStage only for service, has own environment
	secondStage := &v0.Stage{
		ID:   "Stage2",
		Name: "Stage2",
		Type: v0.StageTypeDeployment,
		Spec: &v0.StageDeployment{
			Service: &v0.DeploymentService{
				UseFromStage: &v0.UseFromStage{Stage: "Stage1"},
			},
			Environment: &v0.Environment{
				EnvironmentRef: "env-gamma",
			},
			Execution: &v0.DeploymentExecution{},
		},
	}

	v0Stages := []*v0.Stages{
		{Stage: firstStage},
		{Stage: secondStage},
	}

	v1Stages := converter.convertStages(v0Stages, "pipeline")

	second := v1Stages[1]
	// Service from useFromStage
	if second.Service == nil || second.Service.Items[0].Id != "svc-alpha" {
		t.Errorf("expected service 'svc-alpha' from useFromStage, got %v", second.Service)
	}
	// Environment is its own
	if second.Environment == nil || len(second.Environment.Items) != 1 || second.Environment.Items[0].Id != "env-gamma" {
		t.Errorf("expected environment 'env-gamma' (own), got %v", second.Environment)
	}
}

// V0 K8s stages carry OS on infrastructure.spec.os. V1 sources K8s stage OS from platform.os
// (runtime.kubernetes.os is no longer a valid V1 field). The converter must lift infra.spec.os
// up to stage.Platform.Os so the converted pipeline runs on the correct OS.
func TestConvertCIStage_K8sInfraOsLiftedToPlatform_Windows(t *testing.T) {
	converter := NewPipelineConverter()

	stage := &v0.Stage{
		ID:   "Build",
		Name: "Build",
		Type: v0.StageTypeCI,
		Spec: &v0.StageCI{
			Infrastructure: &v0.Infrastructure{
				Type: "KubernetesDirect",
				Spec: &v0.InfrastructureKubernetesDirectSpec{
					Namespace: "ci",
					Conn:      "k8s-connector",
					OS:        "Windows",
				},
			},
			Execution: v0.Execution{},
		},
	}

	v1Stages := converter.convertStages([]*v0.Stages{{Stage: stage}}, "pipeline")
	if len(v1Stages) != 1 {
		t.Fatalf("expected 1 stage, got %d", len(v1Stages))
	}
	got := v1Stages[0]

	// platform.os must be set, lowercased, to match the V1 enum casing.
	if got.Platform == nil {
		t.Fatal("expected stage.Platform to be non-nil after lifting infra.spec.os")
	}
	if got.Platform.Os != "windows" {
		t.Errorf("expected platform.os=\"windows\", got %q", got.Platform.Os)
	}
	// Arch must be populated even though V0 K8s infra has no Arch field — V0 K8s pipelines
	// implicitly defaulted to amd64, AND the V1 backend's PlatformV1.toPlatform() NPEs if
	// platform is present but arch is missing.
	if got.Platform.Arch != "amd64" {
		t.Errorf("expected platform.arch=\"amd64\" (V0 K8s implicit default), got %q", got.Platform.Arch)
	}
	// (RuntimeKubernetes no longer has an OS field — the V1 schema removed it, so we
	// dropped it from the struct in runtime_kubernetes.go. Nothing to assert here.)
	if got.Runtime == nil || got.Runtime.Kubernetes == nil {
		t.Fatal("expected runtime.kubernetes to be non-nil")
	}
}

// Explicit V0 top-level platform.os must win over the lifted infra.spec.os, since the user
// specified it explicitly. (ConvertPlatform already lowercases on its own.)
func TestConvertCIStage_ExplicitV0PlatformWinsOverInfraOs(t *testing.T) {
	converter := NewPipelineConverter()

	stage := &v0.Stage{
		ID:   "Build",
		Name: "Build",
		Type: v0.StageTypeCI,
		Spec: &v0.StageCI{
			Platform: &v0.Platform{OS: "Linux", Arch: "Amd64"},
			Infrastructure: &v0.Infrastructure{
				Type: "KubernetesDirect",
				Spec: &v0.InfrastructureKubernetesDirectSpec{
					Namespace: "ci",
					Conn:      "k8s-connector",
					OS:        "Windows", // different from platform.os; platform.os must win
				},
			},
			Execution: v0.Execution{},
		},
	}

	v1Stages := converter.convertStages([]*v0.Stages{{Stage: stage}}, "pipeline")
	got := v1Stages[0]

	if got.Platform == nil || got.Platform.Os != "linux" {
		t.Errorf("expected explicit platform.os=\"linux\" to win, got %+v", got.Platform)
	}
}

// When neither V0 platform.os nor V0 infra.spec.os is set, stage.Platform must remain nil so
// the V1 backend applies its own default (Linux). Importantly we must not synthesize a Platform
// with empty fields.
func TestConvertCIStage_NoOsAnywhere_PlatformStaysNil(t *testing.T) {
	converter := NewPipelineConverter()

	stage := &v0.Stage{
		ID:   "Build",
		Name: "Build",
		Type: v0.StageTypeCI,
		Spec: &v0.StageCI{
			Infrastructure: &v0.Infrastructure{
				Type: "KubernetesDirect",
				Spec: &v0.InfrastructureKubernetesDirectSpec{
					Namespace: "ci",
					Conn:      "k8s-connector",
					// no OS
				},
			},
			Execution: v0.Execution{},
		},
	}

	v1Stages := converter.convertStages([]*v0.Stages{{Stage: stage}}, "pipeline")
	got := v1Stages[0]

	if got.Platform != nil {
		t.Errorf("expected stage.Platform to remain nil when no OS is set anywhere, got %+v", got.Platform)
	}
}

// Linux is the dominant V0 K8s case in production. Same lift path as Windows: V0 carries OS
// on infra.spec.os (capitalized) and no top-level platform; converter must emit
// platform.{os: linux, arch: amd64}. This guards the common path against a silent regression.
func TestConvertCIStage_K8sInfraOsLiftedToPlatform_Linux(t *testing.T) {
	converter := NewPipelineConverter()

	stage := &v0.Stage{
		ID:   "Build",
		Name: "Build",
		Type: v0.StageTypeCI,
		Spec: &v0.StageCI{
			Infrastructure: &v0.Infrastructure{
				Type: "KubernetesDirect",
				Spec: &v0.InfrastructureKubernetesDirectSpec{
					Namespace: "ci",
					Conn:      "k8s-connector",
					OS:        "Linux",
				},
			},
			Execution: v0.Execution{},
		},
	}

	v1Stages := converter.convertStages([]*v0.Stages{{Stage: stage}}, "pipeline")
	if len(v1Stages) != 1 {
		t.Fatalf("expected 1 stage, got %d", len(v1Stages))
	}
	got := v1Stages[0]

	if got.Platform == nil {
		t.Fatal("expected stage.Platform to be non-nil after lifting infra.spec.os")
	}
	if got.Platform.Os != "linux" {
		t.Errorf("expected platform.os=\"linux\" (lowercased from V0 \"Linux\"), got %q", got.Platform.Os)
	}
	if got.Platform.Arch != "amd64" {
		t.Errorf("expected platform.arch=\"amd64\" (V0 K8s implicit default), got %q", got.Platform.Arch)
	}
	if got.Runtime == nil || got.Runtime.Kubernetes == nil {
		t.Fatal("expected runtime.kubernetes to be non-nil")
	}
}
