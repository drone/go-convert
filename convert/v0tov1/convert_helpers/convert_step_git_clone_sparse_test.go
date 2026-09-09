package converthelpers

import (
	"testing"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	"github.com/drone/go-convert/internal/flexible"
)

func TestConvertStepGitClone_SparseCheckoutPathsWithSpaces(t *testing.T) {
	result := ConvertStepGitClone(&v0.Step{
		Spec: &v0.StepGitClone{
			ConnRef: "gitlab-connector",
			SparseCheckout: &flexible.Field[[]string]{
				Value: []string{"folder2", "folder with space"},
			},
		},
	})
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	with, ok := result.With.(map[string]interface{})
	if !ok {
		t.Fatalf("expected with map, got %T", result.With)
	}
	got, ok := with["sparse_checkout"].(string)
	if !ok {
		t.Fatalf("expected sparse_checkout string, got %T", with["sparse_checkout"])
	}
	want := "folder2\nfolder with space"
	if got != want {
		t.Fatalf("sparse_checkout = %q, want %q", got, want)
	}
}
