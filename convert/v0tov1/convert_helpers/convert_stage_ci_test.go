package converthelpers

import (
	"testing"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/drone/go-convert/internal/flexible"
	"github.com/google/go-cmp/cmp"
)

func TestConvertRuntime(t *testing.T) {
	tests := []struct {
		name     string
		input    *v0.Runtime
		expected *v1.Runtime
	}{
		{
			name:     "nil runtime",
			input:    nil,
			expected: nil,
		},
		{
			name: "Cloud runtime with image and size",
			input: &v0.Runtime{
				Type: "Cloud",
				Spec: &v0.RuntimeCloudSpec{
					Size: "flex",
					ImageSpec: &v0.ImageSpec{
						ImageName: "ubuntu-latest",
					},
				},
			},
			expected: &v1.Runtime{
				Cloud: &v1.RuntimeCloud{
					Image: v1.MachineImage("ubuntu-latest"),
					Size:  v1.MachineSize("flex"),
				},
			},
		},
		{
			name: "Cloud runtime with empty image spec",
			input: &v0.Runtime{
				Type: "Cloud",
				Spec: &v0.RuntimeCloudSpec{
					Size: "large",
				},
			},
			expected: &v1.Runtime{
				Cloud: &v1.RuntimeCloud{
					Size: v1.MachineSize("large"),
				},
			},
		},
		{
			name: "Docker runtime converts to shell true",
			input: &v0.Runtime{
				Type: "Docker",
				Spec: &v0.RuntimeDockerSpec{},
			},
			expected: &v1.Runtime{
				Shell: true,
			},
		},
		{
			name: "unknown runtime type returns nil",
			input: &v0.Runtime{
				Type: "UnknownType",
				Spec: nil,
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertRuntime(tt.input)
			if diff := cmp.Diff(tt.expected, result); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertInfrastructureToRuntime(t *testing.T) {
	tests := []struct {
		name     string
		input    *v0.Infrastructure
		expected *v1.Runtime
	}{
		{
			name:     "nil infrastructure",
			input:    nil,
			expected: nil,
		},
		{
			name: "KubernetesDirect with full spec",
			input: &v0.Infrastructure{
				Type: "KubernetesDirect",
				Spec: &v0.InfrastructureKubernetesDirectSpec{
					Conn:               "k8s-connector",
					Namespace:          "ci-builds",
					ServiceAccountName: "ci-sa",
					InitTimeout:        "10m",
					OS:                 "Linux",
					HarnessImageConnectorRef: "harness-docker",
					ImagePullPolicy:    "Always",
					PodSpecOverlay:     "overlay-spec",
					Annotations: &flexible.Field[map[string]string]{Value: map[string]string{
						"app": "ci",
					}},
					Labels: &flexible.Field[map[string]string]{Value: map[string]string{
						"team": "platform",
					}},
					NodeSelector: &flexible.Field[map[string]string]{Value: map[string]string{
						"disktype": "ssd",
					}},
					AutomountServiceAccountToken: &flexible.Field[bool]{Value: true},
					PriorityClassName: "high-priority",
					HostNames: &flexible.Field[[]string]{Value: []string{"host1", "host2"}},
				},
			},
			expected: &v1.Runtime{
				Kubernetes: &v1.RuntimeKubernetes{
					Namespace:             "ci-builds",
					Connector:             "k8s-connector",
					ServiceAccount:        "ci-sa",
					Timeout:               "10m",
					OS:                    "Linux",
					HarnessImageConnector: "harness-docker",
					ImagePullPolicy:       "always",
					PodSpecOverlay:        "overlay-spec",
					PriorityClass:         "high-priority",
					Annotations: &flexible.Field[map[string]string]{Value: map[string]string{
						"app": "ci",
					}},
					Labels: &flexible.Field[map[string]string]{Value: map[string]string{
						"team": "platform",
					}},
					Node: &flexible.Field[map[string]string]{Value: map[string]string{
						"disktype": "ssd",
					}},
					ServiceToken: &flexible.Field[bool]{Value: true},
					Host:         &flexible.Field[[]string]{Value: []string{"host1", "host2"}},
				},
			},
		},
		{
			name: "VM with pool",
			input: &v0.Infrastructure{
				Type: "VM",
				Spec: &v0.InfrastructureVMSpec{
					Type: "Pool",
					Spec: &v0.InfrastructureVMPool{
						PoolName:                 "my-pool",
						OS:                       "Linux",
						HarnessImageConnectorRef: "vm-connector",
						Timeout:                  "30m",
					},
				},
			},
			expected: &v1.Runtime{
				VM: &v1.RuntimeInstance{
					Pool:                  "my-pool",
					Os:                    "Linux",
					HarnessImageConnector: "vm-connector",
					Timeout:               "30m",
				},
			},
		},
		{
			name: "VM with pool using identifier fallback",
			input: &v0.Infrastructure{
				Type: "VM",
				Spec: &v0.InfrastructureVMSpec{
					Type: "Pool",
					Spec: &v0.InfrastructureVMPool{
						Identifier: "pool-id-fallback",
						OS:         "Windows",
					},
				},
			},
			expected: &v1.Runtime{
				VM: &v1.RuntimeInstance{
					Pool: "pool-id-fallback",
					Os:   "Windows",
				},
			},
		},
		{
			name: "unknown infrastructure type returns nil",
			input: &v0.Infrastructure{
				Type: "UnknownInfra",
				Spec: nil,
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertInfrastructureToRuntime(tt.input)
			if diff := cmp.Diff(tt.expected, result); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertServiceDependencyToBackgroundStep(t *testing.T) {
	tests := []struct {
		name     string
		input    *v0.Service
		expected *v1.Step
	}{
		{
			name:     "nil service",
			input:    nil,
			expected: nil,
		},
		{
			name: "nil spec",
			input: &v0.Service{
				ID:   "svc1",
				Name: "Service 1",
				Spec: nil,
			},
			expected: nil,
		},
		{
			name: "service with image and connector",
			input: &v0.Service{
				ID:   "redis_svc",
				Name: "Redis Service",
				Spec: &v0.ServiceSpec{
					Image: "redis:7",
					Conn:  "docker-conn",
					Env: &flexible.Field[map[string]interface{}]{Value: map[string]interface{}{
						"REDIS_PORT": "6379",
					}},
					Entrypoint: &flexible.Field[[]string]{Value: []string{"redis-server"}},
					Args:       []string{"--maxmemory", "256mb"},
					Privileged: &flexible.Field[bool]{Value: false},
				},
			},
			expected: &v1.Step{
				Id:   "redis_svc",
				Name: "Redis Service",
				Background: &v1.StepRun{
					Container: &v1.Container{
						Image:      "redis:7",
						Connector:  "docker-conn",
						Entrypoint: &flexible.Field[[]string]{Value: []string{"redis-server"}},
						Args:       v1.Stringorslice{"--maxmemory", "256mb"},
						Privileged: &flexible.Field[bool]{Value: false},
					},
					Env: &flexible.Field[map[string]interface{}]{Value: map[string]interface{}{
						"REDIS_PORT": "6379",
					}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertServiceDependencyToBackgroundStep(tt.input)
			if diff := cmp.Diff(tt.expected, result); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertServiceDependenciesToBackgroundSteps(t *testing.T) {
	tests := []struct {
		name     string
		input    []*v0.Service
		expected []*v1.Step
	}{
		{
			name:     "nil services",
			input:    nil,
			expected: nil,
		},
		{
			name:     "empty services",
			input:    []*v0.Service{},
			expected: nil,
		},
		{
			name: "multiple services",
			input: []*v0.Service{
				{
					ID:   "redis",
					Name: "Redis",
					Spec: &v0.ServiceSpec{
						Image: "redis:7",
					},
				},
				{
					ID:   "postgres",
					Name: "Postgres",
					Spec: &v0.ServiceSpec{
						Image: "postgres:14",
					},
				},
			},
			expected: []*v1.Step{
				{
					Id:   "redis",
					Name: "Redis",
					Background: &v1.StepRun{
						Container: &v1.Container{
							Image: "redis:7",
						},
					},
				},
				{
					Id:   "postgres",
					Name: "Postgres",
					Background: &v1.StepRun{
						Container: &v1.Container{
							Image: "postgres:14",
						},
					},
				},
			},
		},
		{
			name: "nil services in slice are skipped",
			input: []*v0.Service{
				{
					ID:   "redis",
					Name: "Redis",
					Spec: &v0.ServiceSpec{Image: "redis:7"},
				},
				nil,
				{
					ID:   "no-spec",
					Name: "No Spec",
					Spec: nil,
				},
			},
			expected: []*v1.Step{
				{
					Id:   "redis",
					Name: "Redis",
					Background: &v1.StepRun{
						Container: &v1.Container{
							Image: "redis:7",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertServiceDependenciesToBackgroundSteps(tt.input)
			if diff := cmp.Diff(tt.expected, result); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertSharedPaths(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []*v1.Volume
	}{
		{
			name:     "nil paths",
			input:    nil,
			expected: nil,
		},
		{
			name:     "empty paths",
			input:    []string{},
			expected: nil,
		},
		{
			name:  "single path",
			input: []string{"/tmp/shared"},
			expected: []*v1.Volume{
				{
					Name: "shared-0",
					Uses: "empty-dir",
					With: &v1.VolumeEmptyDir{MountPath: "/tmp/shared"},
				},
			},
		},
		{
			name:  "multiple paths",
			input: []string{"/tmp/a", "/var/data", "/opt/cache"},
			expected: []*v1.Volume{
				{Name: "shared-0", Uses: "empty-dir", With: &v1.VolumeEmptyDir{MountPath: "/tmp/a"}},
				{Name: "shared-1", Uses: "empty-dir", With: &v1.VolumeEmptyDir{MountPath: "/var/data"}},
				{Name: "shared-2", Uses: "empty-dir", With: &v1.VolumeEmptyDir{MountPath: "/opt/cache"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertSharedPaths(tt.input)
			if diff := cmp.Diff(tt.expected, result); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertInfrastructureToVolumes(t *testing.T) {
	boolTrue := &flexible.Field[bool]{Value: true}
	boolFalse := &flexible.Field[bool]{Value: false}

	tests := []struct {
		name     string
		input    *v0.Infrastructure
		expected []*v1.Volume
	}{
		{
			name:     "nil infrastructure",
			input:    nil,
			expected: nil,
		},
		{
			name: "non-k8s infrastructure",
			input: &v0.Infrastructure{
				Type: "VM",
				Spec: &v0.InfrastructureVMSpec{},
			},
			expected: nil,
		},
		{
			name: "k8s with no volumes",
			input: &v0.Infrastructure{
				Type: "KubernetesDirect",
				Spec: &v0.InfrastructureKubernetesDirectSpec{
					Namespace: "default",
					Volumes:   nil,
				},
			},
			expected: nil,
		},
		{
			name: "EmptyDir volume",
			input: &v0.Infrastructure{
				Type: "KubernetesDirect",
				Spec: &v0.InfrastructureKubernetesDirectSpec{
					Volumes: []*v0.Volume{
						{
							Type:      "EmptyDir",
							MountPath: "/tmp/data",
							Spec:      v0.EmptyDirVolumeSpec{Medium: "Memory", Size: "1Gi"},
						},
					},
				},
			},
			expected: []*v1.Volume{
				{
					Name: "infra-0",
					Uses: "empty-dir",
					With: &v1.VolumeEmptyDir{MountPath: "/tmp/data", Medium: "Memory", Size: "1Gi"},
				},
			},
		},
		{
			name: "PersistentVolumeClaim volume",
			input: &v0.Infrastructure{
				Type: "KubernetesDirect",
				Spec: &v0.InfrastructureKubernetesDirectSpec{
					Volumes: []*v0.Volume{
						{
							Type:      "PersistentVolumeClaim",
							MountPath: "/data",
							Spec:      v0.PersistentVolumeClaimVolumeSpec{ClaimName: "my-pvc", ReadOnly: boolTrue},
						},
					},
				},
			},
			expected: []*v1.Volume{
				{
					Name: "infra-0",
					Uses: "persistent-volume-claim",
					With: &v1.VolumeClaim{Name: "my-pvc", MountPath: "/data", ReadOnly: boolTrue},
				},
			},
		},
		{
			name: "HostPath volume",
			input: &v0.Infrastructure{
				Type: "KubernetesDirect",
				Spec: &v0.InfrastructureKubernetesDirectSpec{
					Volumes: []*v0.Volume{
						{
							Type:      "HostPath",
							MountPath: "/host/data",
							Spec:      v0.HostPathVolumeSpec{Path: "/mnt/data", Type: "DirectoryOrCreate"},
						},
					},
				},
			},
			expected: []*v1.Volume{
				{
					Name: "infra-0",
					Uses: "host-path",
					With: &v1.VolumeHostPath{Path: "/mnt/data", MountPath: "/host/data", Type: "DirectoryOrCreate"},
				},
			},
		},
		{
			name: "ConfigMap volume",
			input: &v0.Infrastructure{
				Type: "KubernetesDirect",
				Spec: &v0.InfrastructureKubernetesDirectSpec{
					Volumes: []*v0.Volume{
						{
							Type:      "ConfigMap",
							MountPath: "/etc/config",
							Spec:      v0.ConfigMapVolumeSpec{Name: "my-config", Optional: boolFalse},
						},
					},
				},
			},
			expected: []*v1.Volume{
				{
					Name: "infra-0",
					Uses: "config-map",
					With: &v1.VolumeConfigMap{Name: "my-config", MountPath: "/etc/config", Optional: boolFalse},
				},
			},
		},
		{
			name: "Secret volume",
			input: &v0.Infrastructure{
				Type: "KubernetesDirect",
				Spec: &v0.InfrastructureKubernetesDirectSpec{
					Volumes: []*v0.Volume{
						{
							Type:      "Secret",
							MountPath: "/etc/secrets",
							Spec:      v0.SecretVolumeSpec{Name: "my-secret", Optional: boolTrue},
						},
					},
				},
			},
			expected: []*v1.Volume{
				{
					Name: "infra-0",
					Uses: "secret",
					With: &v1.VolumeSecret{Name: "my-secret", MountPath: "/etc/secrets", Optional: boolTrue},
				},
			},
		},
		{
			name: "multiple volume types",
			input: &v0.Infrastructure{
				Type: "KubernetesDirect",
				Spec: &v0.InfrastructureKubernetesDirectSpec{
					Volumes: []*v0.Volume{
						{
							Type:      "EmptyDir",
							MountPath: "/tmp",
							Spec:      v0.EmptyDirVolumeSpec{},
						},
						{
							Type:      "Secret",
							MountPath: "/secrets",
							Spec:      v0.SecretVolumeSpec{Name: "creds"},
						},
					},
				},
			},
			expected: []*v1.Volume{
				{Name: "infra-0", Uses: "empty-dir", With: &v1.VolumeEmptyDir{MountPath: "/tmp"}},
				{Name: "infra-1", Uses: "secret", With: &v1.VolumeSecret{Name: "creds", MountPath: "/secrets"}},
			},
		},
		{
			name: "nil volume entries are skipped",
			input: &v0.Infrastructure{
				Type: "KubernetesDirect",
				Spec: &v0.InfrastructureKubernetesDirectSpec{
					Volumes: []*v0.Volume{
						nil,
						{
							Type:      "EmptyDir",
							MountPath: "/data",
							Spec:      v0.EmptyDirVolumeSpec{},
						},
					},
				},
			},
			expected: []*v1.Volume{
				{Name: "infra-1", Uses: "empty-dir", With: &v1.VolumeEmptyDir{MountPath: "/data"}},
			},
		},
		{
			name: "unknown volume type is skipped",
			input: &v0.Infrastructure{
				Type: "KubernetesDirect",
				Spec: &v0.InfrastructureKubernetesDirectSpec{
					Volumes: []*v0.Volume{
						{
							Type:      "UnknownVolumeType",
							MountPath: "/data",
							Spec:      nil,
						},
					},
				},
			},
			expected: []*v1.Volume{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertInfrastructureToVolumes(tt.input)
			if diff := cmp.Diff(tt.expected, result); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
