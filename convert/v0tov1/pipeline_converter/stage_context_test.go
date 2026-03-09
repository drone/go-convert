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

	svc := &v1.ServiceRef{Items: []string{"my-service"}}
	env := &v1.EnvironmentRef{Name: "my-env", Id: "my-env"}
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

	v1Stages := converter.convertStages(v0Stages)

	if len(v1Stages) != 2 {
		t.Fatalf("expected 2 stages, got %d", len(v1Stages))
	}

	// Verify first stage has its own service and environment
	first := v1Stages[0]
	if first.Service == nil {
		t.Fatal("first stage: service should not be nil")
	}
	if len(first.Service.Items) != 1 || first.Service.Items[0] != "dashboard-svc" {
		t.Errorf("first stage: expected service 'dashboard-svc', got %v", first.Service.Items)
	}
	if first.Environment == nil {
		t.Fatal("first stage: environment should not be nil")
	}
	if first.Environment.Name != "prod-env" {
		t.Errorf("first stage: expected environment 'prod-env', got %s", first.Environment.Name)
	}

	// Verify second stage copied service and environment from first
	second := v1Stages[1]
	if second.Service == nil {
		t.Fatal("second stage: service should not be nil (useFromStage)")
	}
	if len(second.Service.Items) != 1 || second.Service.Items[0] != "dashboard-svc" {
		t.Errorf("second stage: expected service 'dashboard-svc' from useFromStage, got %v", second.Service.Items)
	}
	if second.Environment == nil {
		t.Fatal("second stage: environment should not be nil (useFromStage)")
	}
	if second.Environment.Name != "prod-env" {
		t.Errorf("second stage: expected environment 'prod-env' from useFromStage, got %s", second.Environment.Name)
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

	v1Stages := converter.convertStages(v0Stages)

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
	v1Stages := converter.convertStages(v0Stages)

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

	v1Stages := converter.convertStages(v0Stages)

	second := v1Stages[1]
	// Service from useFromStage
	if second.Service == nil || second.Service.Items[0] != "svc-alpha" {
		t.Errorf("expected service 'svc-alpha' from useFromStage, got %v", second.Service)
	}
	// Environment is its own
	if second.Environment == nil || second.Environment.Name != "env-gamma" {
		t.Errorf("expected environment 'env-gamma' (own), got %v", second.Environment)
	}
}
