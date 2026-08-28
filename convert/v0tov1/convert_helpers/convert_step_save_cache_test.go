package converthelpers

import (
	"testing"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	"github.com/drone/go-convert/internal/flexible"
	"github.com/google/go-cmp/cmp"
)

func TestConvertStepSaveCache(t *testing.T) {
	tests := []struct {
		name     string
		step     *v0.Step
		expected map[string]interface{}
	}{
		{
			name: "all required fields",
			step: &v0.Step{
				Spec: &v0.StepSaveCache{
					Key:         "myApp-cache",
					SourcePaths: &flexible.Field[[]string]{Value: []string{"/root/.m2"}},
				},
			},
			expected: map[string]interface{}{
				"cachekey":    "myApp-cache",
				"source_path": &flexible.Field[[]string]{Value: []string{"/root/.m2"}},
			},
		},
		{
			name: "with optional archive format and override",
			step: &v0.Step{
				Spec: &v0.StepSaveCache{
					Key:           "myApp-{{ checksum \"pom.xml\" }}",
					SourcePaths:   &flexible.Field[[]string]{Value: []string{"/root/.m2", "/root/.gradle"}},
					ArchiveFormat: "Gzip",
					Override:      &flexible.Field[bool]{Value: true},
				},
			},
			expected: map[string]interface{}{
				"cachekey":      "myApp-{{ checksum \"pom.xml\" }}",
				"source_path":   &flexible.Field[[]string]{Value: []string{"/root/.m2", "/root/.gradle"}},
				"archiveformat": "gzip",
				"override":      &flexible.Field[bool]{Value: true},
			},
		},
		{
			name: "minimal empty spec",
			step: &v0.Step{
				Spec: &v0.StepSaveCache{},
			},
			expected: map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertStepSaveCache(tt.step)
			if result == nil {
				t.Fatal("expected non-nil result")
			}

			if result.Uses != "saveCache" {
				t.Errorf("expected Uses to be saveCache, got %s", result.Uses)
			}

			if diff := cmp.Diff(tt.expected, result.With); diff != "" {
				t.Errorf("With mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertStepSaveCache_Container(t *testing.T) {
	step := &v0.Step{
		Spec: &v0.StepSaveCache{
			Key:       "cache-key",
			RunAsUser: &flexible.Field[int]{Value: 1000},
			Resources: &v0.Resources{
				Limits: &v0.ResourceSpec{
					CPU:    &flexible.Field[*v0.MilliSize]{Value: "500m"},
					Memory: &flexible.Field[*v0.BytesSize]{Value: "512Mi"},
				},
			},
		},
	}

	result := ConvertStepSaveCache(step)
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.Container == nil {
		t.Fatal("expected non-nil container")
	}

	if result.Container.User == nil {
		t.Fatal("expected container.user to be set")
	}

	if result.Container.Resources == nil || result.Container.Resources.Limits == nil {
		t.Fatal("expected container resource limits to be set")
	}

	if result.Container.Resources.Limits.Cpu != "500m" {
		t.Errorf("expected cpu limit 500m, got %s", result.Container.Resources.Limits.Cpu)
	}

	if result.Container.Resources.Limits.Memory != "512Mi" {
		t.Errorf("expected memory limit 512Mi, got %s", result.Container.Resources.Limits.Memory)
	}
}

func TestConvertStepSaveCache_NilCases(t *testing.T) {
	tests := []struct {
		name string
		step *v0.Step
	}{
		{
			name: "nil step",
			step: nil,
		},
		{
			name: "nil spec",
			step: &v0.Step{Spec: nil},
		},
		{
			name: "wrong spec type",
			step: &v0.Step{Spec: &v0.StepRun{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertStepSaveCache(tt.step)
			if result != nil {
				t.Errorf("expected nil result, got %v", result)
			}
		})
	}
}
