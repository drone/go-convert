package converthelpers

import (
	"testing"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/drone/go-convert/internal/flexible"
	"github.com/google/go-cmp/cmp"
)

func TestConvertStepBackground(t *testing.T) {
	tests := []struct {
		name     string
		step     *v0.Step
		expected *v1.StepRun
	}{
		{
			name: "minimal step with image only",
			step: &v0.Step{
				Spec: &v0.StepBackground{
					Image: "redis:latest",
				},
			},
			expected: &v1.StepRun{
				Container: &v1.Container{
					Image: "redis:latest",
					Ports: []string{},
				},
				Shell: "sh",
			},
		},
		{
			name: "full step with ports, reports, resources, shell, env, entrypoint",
			step: &v0.Step{
				Spec: &v0.StepBackground{
					Image:           "postgres:14",
					ConnRef:         "docker-connector",
					Command:         "pg_isready",
					Shell:           "Bash",
					ImagePullPolicy: "Always",
					Privileged:      &flexible.Field[bool]{Value: true},
					Entrypoint:      &flexible.Field[[]string]{Value: []string{"/bin/sh", "-c"}},
					Env: &flexible.Field[map[string]interface{}]{Value: map[string]interface{}{
						"POSTGRES_USER": "admin",
						"POSTGRES_DB":   "testdb",
					}},
					PortBindings: map[string]string{
						"5432": "5432",
						"8080": "80",
					},
					Reports: &v0.Report{
						Type: "JUnit",
						Spec: &v0.ReportJunit{
							Paths: []string{"report1.xml", "report2.xml"},
						},
					},
					RunAsUser: &flexible.Field[int]{Value: 1000},
				},
			},
			expected: &v1.StepRun{
				Container: &v1.Container{
					Image:      "postgres:14",
					Connector:  "docker-connector",
					Privileged: &flexible.Field[bool]{Value: true},
					Pull:       "always",
					Entrypoint: &flexible.Field[[]string]{Value: []string{"/bin/sh", "-c"}},
					User:       &flexible.Field[int]{Value: 1000},
					Ports:      []string{"5432:5432", "8080:80"},
				},
				Script: v1.Stringorslice{"pg_isready"},
				Shell:  "bash",
				Env: &flexible.Field[map[string]interface{}]{Value: map[string]interface{}{
					"POSTGRES_USER": "admin",
					"POSTGRES_DB":   "testdb",
				}},
				Report: &v1.ReportList{
					&v1.Report{Type: "junit", Path: "report1.xml"},
					&v1.Report{Type: "junit", Path: "report2.xml"},
				},
			},
		},
		{
			name: "empty command produces no script",
			step: &v0.Step{
				Spec: &v0.StepBackground{
					Image: "nginx:latest",
				},
			},
			expected: &v1.StepRun{
				Container: &v1.Container{
					Image: "nginx:latest",
					Ports: []string{},
				},
				Shell: "sh",
			},
		},
		{
			name: "unknown shell value is lowercased",
			step: &v0.Step{
				Spec: &v0.StepBackground{
					Image: "alpine",
					Shell: "Zsh",
				},
			},
			expected: &v1.StepRun{
				Container: &v1.Container{
					Image: "alpine",
					Ports: []string{},
				},
				Shell: "zsh",
			},
		},
		{
			name: "ImagePullPolicy Never",
			step: &v0.Step{
				Spec: &v0.StepBackground{
					Image:           "myimage",
					ImagePullPolicy: "Never",
				},
			},
			expected: &v1.StepRun{
				Container: &v1.Container{
					Image: "myimage",
					Pull:  "never",
					Ports: []string{},
				},
				Shell: "sh",
			},
		},
		{
			name: "ImagePullPolicy IfNotPresent",
			step: &v0.Step{
				Spec: &v0.StepBackground{
					Image:           "myimage",
					ImagePullPolicy: "IfNotPresent",
				},
			},
			expected: &v1.StepRun{
				Container: &v1.Container{
					Image: "myimage",
					Pull:  "if-not-present",
					Ports: []string{},
				},
				Shell: "sh",
			},
		},
		{
			name: "reports with empty paths are skipped",
			step: &v0.Step{
				Spec: &v0.StepBackground{
					Image: "test",
					Reports: &v0.Report{
						Type: "JUnit",
						Spec: &v0.ReportJunit{
							Paths: []string{"valid.xml", "  ", ""},
						},
					},
				},
			},
			expected: &v1.StepRun{
				Container: &v1.Container{
					Image: "test",
					Ports: []string{},
				},
				Shell: "sh",
				Report: &v1.ReportList{
					&v1.Report{Type: "junit", Path: "valid.xml"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertStepBackground(tt.step)
			if result == nil {
				t.Fatal("expected non-nil result")
			}

			if diff := cmp.Diff(tt.expected, result); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertStepBackground_NilCases(t *testing.T) {
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
			step: &v0.Step{
				Spec: nil,
			},
		},
		{
			name: "wrong spec type",
			step: &v0.Step{
				Spec: &v0.StepRun{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertStepBackground(tt.step)
			if result != nil {
				t.Errorf("expected nil result, got %v", result)
			}
		})
	}
}
