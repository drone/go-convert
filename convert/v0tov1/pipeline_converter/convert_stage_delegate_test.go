package pipelineconverter

import (
	"testing"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	"github.com/drone/go-convert/internal/flexible"
)

// PIPE-37434: stage-level delegate selectors move under runtime.delegate
// only when the stage's runtime is shell, so that delegate/shell live
// together under `runtime`. All other stages (no runtime, or a non-shell
// runtime such as cloud/kubernetes/vm) omit the delegate field entirely —
// it is no longer set at the top level of the stage.
func TestConvertStage_DelegateMovesUnderRuntime_ShellRuntime(t *testing.T) {
	converter := NewPipelineConverter()

	// Custom stages always get a fixed Shell runtime.
	stage := &v0.Stage{
		ID:                "Deploy",
		Name:              "Deploy",
		Type:              v0.StageTypeCustom,
		DelegateSelectors: &flexible.Field[[]string]{Value: []string{"del1"}},
		Spec: &v0.StageCustom{
			Execution: &v0.Execution{},
		},
	}

	v1Stages := converter.convertStages([]*v0.Stages{{Stage: stage}}, "pipeline")
	if len(v1Stages) != 1 {
		t.Fatalf("expected 1 stage, got %d", len(v1Stages))
	}
	got := v1Stages[0]

	if got.Runtime == nil || got.Runtime.Shell == nil {
		t.Fatal("expected stage.Runtime.Shell to be non-nil")
	}
	if got.Delegate != nil {
		t.Errorf("expected top-level stage.Delegate to be nil, got %v", got.Delegate)
	}
	if got.Runtime.Delegate == nil {
		t.Fatal("expected stage.Runtime.Delegate to be non-nil")
	}
	delegate, ok := got.Runtime.Delegate.AsStruct()
	if !ok || delegate == nil {
		t.Fatal("expected runtime.Delegate to resolve to a struct")
	}
	if len(delegate.Filter) != 1 || delegate.Filter[0] != "del1" {
		t.Errorf("expected delegate filter [del1], got %v", delegate.Filter)
	}
}

func TestConvertStage_DelegateOmitted_NonShellRuntime(t *testing.T) {
	converter := NewPipelineConverter()

	// CI stage with a Kubernetes runtime (not shell).
	stage := &v0.Stage{
		ID:                "Build",
		Name:              "Build",
		Type:              v0.StageTypeCI,
		DelegateSelectors: &flexible.Field[[]string]{Value: []string{"del1"}},
		Spec: &v0.StageCI{
			Infrastructure: &v0.Infrastructure{
				Type: "KubernetesDirect",
				Spec: &v0.InfrastructureKubernetesDirectSpec{
					Namespace: "ci",
					Conn:      "k8s-connector",
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

	if got.Runtime == nil || got.Runtime.Kubernetes == nil {
		t.Fatal("expected stage.Runtime.Kubernetes to be non-nil")
	}
	if got.Runtime.Delegate != nil {
		t.Errorf("expected stage.Runtime.Delegate to be nil for non-shell runtime, got %v", got.Runtime.Delegate)
	}
	if got.Delegate != nil {
		t.Errorf("expected top-level stage.Delegate to be nil (dropped for non-shell runtime), got %v", got.Delegate)
	}
}

func TestConvertStage_DelegateOmitted_WhenNoRuntime(t *testing.T) {
	converter := NewPipelineConverter()

	stage := &v0.Stage{
		ID:                "Approve",
		Name:              "Approve",
		Type:              v0.StageTypeApproval,
		DelegateSelectors: &flexible.Field[[]string]{Value: []string{"del1"}},
		Spec:              &v0.StageApproval{},
	}

	v1Stages := converter.convertStages([]*v0.Stages{{Stage: stage}}, "pipeline")
	if len(v1Stages) != 1 {
		t.Fatalf("expected 1 stage, got %d", len(v1Stages))
	}
	got := v1Stages[0]

	if got.Runtime != nil {
		t.Errorf("expected stage.Runtime to be nil, got %v", got.Runtime)
	}
	if got.Delegate != nil {
		t.Errorf("expected top-level stage.Delegate to be nil (dropped when no runtime), got %v", got.Delegate)
	}
}
