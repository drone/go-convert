package converthelpers

import (
	"strings"
	"fmt"
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	flexible "github.com/drone/go-convert/internal/flexible"
)

// ConvertCloneCodebase converts v0 cloneCodebase bool to v1 Clone
func ConvertCloneCodebase(clone bool) *v1.Clone {
	if !clone {
		return nil
	}
	return &v1.Clone{
		Enabled: clone,
	}
}

// ConvertCaching converts v0 Cache to v1 Cache
func ConvertCaching(cache *v0.Cache) *v1.Cache {
	if cache == nil {
		return nil
	}

	disabled := !cache.Enabled
	if !disabled {
		return nil
	}
	return &v1.Cache{
		Disabled: disabled,
	}
}

// ConvertBuildIntelligence converts v0 BuildIntelligence to v1 bool
func ConvertBuildIntelligence(bi *v0.BuildIntelligence) *flexible.Field[bool] {
	if bi == nil {
		return nil
	}
	return bi.Enabled
}

func ConvertRuntime(runtime *v0.Runtime) *v1.Runtime {
	if runtime == nil {
		return nil
	}

	switch runtime.Type {
	case "Cloud":
		if cloudSpec, ok := runtime.Spec.(*v0.RuntimeCloudSpec); ok {
			result := &v1.Runtime{
				Cloud: &v1.RuntimeCloud{},
			}

			// Convert image
			if cloudSpec.ImageSpec != nil && cloudSpec.ImageSpec.ImageName != "" {
				result.Cloud.Image = v1.MachineImage(cloudSpec.ImageSpec.ImageName)
			}

			// Convert size - map "flex" to "xlarge" as per your example
			if cloudSpec.Size != "" {
				result.Cloud.Size = v1.MachineSize(cloudSpec.Size)
			}

			return result
		}

	case "Docker":
		// Docker runtime converts to shell: true in v1
		return &v1.Runtime{
			Shell: true,
		}
	}

	return nil
}

// ConvertInfrastructureToRuntime converts v0 Infrastructure to v1 Runtime
func ConvertInfrastructureToRuntime(infra *v0.Infrastructure) *v1.Runtime {
	if infra == nil {
		return nil
	}

	// Handle KubernetesDirect type
	if strings.EqualFold(infra.Type, "KubernetesDirect") {
		if k8sSpec, ok := infra.Spec.(*v0.InfrastructureKubernetesDirectSpec); ok && k8sSpec != nil {
			runtime := &v1.Runtime{
				Kubernetes: &v1.RuntimeKubernetes{
					Namespace:             k8sSpec.Namespace,
					Connector:             k8sSpec.Conn,
					Annotations:           k8sSpec.Annotations,
					Labels:                k8sSpec.Labels,
					Node:                  k8sSpec.NodeSelector,
					OS:                    k8sSpec.OS,
					ServiceAccount:        k8sSpec.ServiceAccountName,
					ServiceToken:          k8sSpec.AutomountServiceAccountToken,
					Tolerations:           convertTolerations(k8sSpec.Tolerations),
					Timeout:               k8sSpec.InitTimeout,
					Host:                  k8sSpec.HostNames,
					PriorityClass:         k8sSpec.PriorityClassName,
					HarnessImageConnector: k8sSpec.HarnessImageConnectorRef,
					PodSpecOverlay:        k8sSpec.PodSpecOverlay,
					ImagePullPolicy:       convertImagePullPolicy(k8sSpec.ImagePullPolicy),
					User:                  k8sSpec.RunAsUser,
					SecurityContext:       convertSecurityContext(k8sSpec.ContainerSecurityContext),
					Volumes:               ConvertInfrastructureToVolumes(infra),
				},
			}
			return runtime
		}
	}

	// Handle VM type
	if strings.EqualFold(infra.Type, "VM") {
		if vmSpec, ok := infra.Spec.(*v0.InfrastructureVMSpec); ok && vmSpec != nil {
			if vmSpec.Type == "Pool" && vmSpec.Spec != nil {
				pool := vmSpec.Spec.PoolName
				if pool == "" {
					pool = vmSpec.Spec.Identifier
				}
				return &v1.Runtime{
					VM: &v1.RuntimeInstance{
						Pool: pool,
						Os: vmSpec.Spec.OS,
						HarnessImageConnector: vmSpec.Spec.HarnessImageConnectorRef,
						Timeout: vmSpec.Spec.Timeout,
					},
				}
			}
		}
	}

	return nil
}

// Helper function to convert image pull policy
func convertImagePullPolicy(policy string) string {
	switch policy {
	case "Always":
		return "always"
	case "Never":
		return "never"
	case "IfNotPresent":
		return "if-not-present"
	default:
		return policy
	}
}

// Convert security context
func convertSecurityContext(v0SecCtx *flexible.Field[*v0.SecurityContext]) *flexible.Field[*v1.SecurityContext] {
	if v0SecCtx == nil {
		return nil
	}

	// If it's an expression, pass it through
	if v0SecCtx.IsExpression() {
		result := &flexible.Field[*v1.SecurityContext]{}
		if expr, ok := v0SecCtx.AsString(); ok {
			result.SetExpression(expr)
		}
		return result
	}

	// Convert struct values
	v0Ctx, ok := v0SecCtx.AsStruct()
	if !ok || v0Ctx == nil {
		return nil
	}

	v1Ctx := &v1.SecurityContext{
		AllowPrivilegeEscalation: v0Ctx.AllowPrivilegeEscalation,
		ProcMount:                v0Ctx.ProcMount,
		Privileged:               v0Ctx.Privileged,
		ReadOnlyRootFilesystem:   v0Ctx.ReadOnlyRootFilesystem,
		RunAsNonRoot:             v0Ctx.RunAsNonRoot,
		RunAsGroup:               v0Ctx.RunAsGroup,
		User:                     v0Ctx.RunAsUser,
		Capabilities:             convertCapabilities(v0Ctx.Capabilities),
	}

	result := &flexible.Field[*v1.SecurityContext]{}
	result.Set(v1Ctx)
	return result
}

// Convert capabilities
func convertCapabilities(v0Caps *flexible.Field[*v0.Capabilities]) *flexible.Field[*v1.Capabilities] {
	if v0Caps == nil {
		return nil
	}

	// If it's an expression, pass it through
	if v0Caps.IsExpression() {
		result := &flexible.Field[*v1.Capabilities]{}
		if expr, ok := v0Caps.AsString(); ok {
			result.SetExpression(expr)
		}
		return result
	}

	// Convert struct values
	v0Cap, ok := v0Caps.AsStruct()
	if !ok || v0Cap == nil {
		return nil
	}

	v1Cap := &v1.Capabilities{
		Add:  v0Cap.Add,
		Drop: v0Cap.Drop,
	}

	result := &flexible.Field[*v1.Capabilities]{}
	result.Set(v1Cap)
	return result
}

func convertTolerations(tolerations *flexible.Field[[]*v0.Toleration]) *flexible.Field[[]*v1.Toleration] {
	if tolerations == nil {
		return nil
	}

	// If it's an expression, pass it through
	if tolerations.IsExpression() {
		result := &flexible.Field[[]*v1.Toleration]{}
		if expr, ok := tolerations.AsString(); ok {
			result.SetExpression(expr)
		}
		return result
	}

	// Convert struct values
	v0Tolerations, ok := tolerations.AsStruct()
	if !ok || v0Tolerations == nil {
		return nil
	}

	v1Tolerations := make([]*v1.Toleration, len(v0Tolerations))
	for i, v0Tol := range v0Tolerations {
		if v0Tol == nil {
			continue
		}
		v1Tolerations[i] = &v1.Toleration{
			Effect:            v0Tol.Effect,
			Key:               v0Tol.Key,
			Operator:          v0Tol.Operator,
			TolerationSeconds: v0Tol.TolerationSeconds,
			Value:             v0Tol.Value,
		}
	}

	result := &flexible.Field[[]*v1.Toleration]{}
	result.Set(v1Tolerations)
	return result
}

// ConvertPlatform converts v0 Platform to v1 Platform
func ConvertPlatform(platform *v0.Platform) *v1.Platform {
	if platform == nil {
		return nil
	}

	return &v1.Platform{
		Os:   strings.ToLower(platform.OS),
		Arch: strings.ToLower(platform.Arch),
	}
}

// ConvertServiceDependencyToBackgroundStep converts a v0 Service dependency to v1 background step
func ConvertServiceDependencyToBackgroundStep(src *v0.Service) *v1.Step {
    if src == nil || src.Spec == nil {
        return nil
    }

    // Container mapping
    var container *v1.Container
    if src.Spec.Image != "" || src.Spec.Conn != "" {
        cpu := ""
        memory := ""
        if src.Spec.Resources != nil && src.Spec.Resources.Limits != nil {
            memory = src.Spec.Resources.Limits.GetMemoryString()
			cpu = src.Spec.Resources.Limits.GetCPUString()
        }

        container = &v1.Container{
            Image:      src.Spec.Image,
            Connector:  src.Spec.Conn,
            Cpu:        cpu,
            Memory:     memory,
            Entrypoint: src.Spec.Entrypoint,
            Args:       src.Spec.Args,
            Privileged: src.Spec.Privileged,
        }
    }

    background := &v1.StepRun{
        Container: container,
        Env:       map[string]interface{}{},
    }

    // Add environment variables
    for k, v := range src.Spec.Env {
        background.Env[k] = v
    }

    // Create the step
    step := &v1.Step{
        Id:         src.ID,
        Name:       src.Name,
        Background: background,
    }

    return step
}

// ConvertServiceDependenciesToBackgroundSteps converts v0 service dependencies to v1 background steps
func ConvertServiceDependenciesToBackgroundSteps(services []*v0.Service) []*v1.Step {
    if len(services) == 0 {
        return nil
    }

    steps := make([]*v1.Step, 0, len(services))
    for _, service := range services {
        if step := ConvertServiceDependencyToBackgroundStep(service); step != nil {
            steps = append(steps, step)
        }
    }

    return steps
}

// ConvertSharedPaths converts v0 SharedPaths to v1 Volumes
func ConvertSharedPaths(sharedPaths []string) []*v1.Volume {
    if len(sharedPaths) == 0 {
        return nil
    }

    volumes := make([]*v1.Volume, 0, len(sharedPaths))
    for i, path := range sharedPaths {
        volume := &v1.Volume{
            Name: fmt.Sprintf("shared-%d", i),
            Uses: "empty-dir",
            With: &v1.VolumeEmptyDir{
                MountPath: path,
            },
        }
        volumes = append(volumes, volume)
    }
    return volumes
}

// ConvertInfrastructureToVolumes converts v0 Infrastructure volumes to v1 Volumes
func ConvertInfrastructureToVolumes(infra *v0.Infrastructure) []*v1.Volume {
	if infra == nil {
		return nil
	}

	k8sSpec, ok := infra.Spec.(*v0.InfrastructureKubernetesDirectSpec)
	if !ok || k8sSpec == nil || len(k8sSpec.Volumes) == 0 {
		return nil
	}

	volumes := make([]*v1.Volume, 0, len(k8sSpec.Volumes))
	for i, v0Vol := range k8sSpec.Volumes {
		if v0Vol == nil {
			continue
		}

		v1Vol := convertInfraVolume(v0Vol, i)
		if v1Vol != nil {
			volumes = append(volumes, v1Vol)
		}
	}

	return volumes
}

func convertInfraVolume(v0Vol *v0.Volume, index int) *v1.Volume {
	if v0Vol == nil {
		return nil
	}

	v1Vol := &v1.Volume{
		Name: fmt.Sprintf("infra-%d", index),
	}

	switch v0Vol.Type {
	case "EmptyDir":
		if spec, ok := v0Vol.Spec.(v0.EmptyDirVolumeSpec); ok {
			v1Vol.Uses = "empty-dir"
			v1Vol.With = &v1.VolumeEmptyDir{
				MountPath: v0Vol.MountPath,
				Medium:    spec.Medium,
				Size:     spec.Size,
			}
		}

	case "PersistentVolumeClaim":
		if spec, ok := v0Vol.Spec.(v0.PersistentVolumeClaimVolumeSpec); ok {
			v1Vol.Uses = "persistent-volume-claim"
			v1Vol.With = &v1.VolumeClaim{
				Name:      spec.ClaimName,
				MountPath: v0Vol.MountPath,
				ReadOnly:  spec.ReadOnly,
			}
		}

	case "HostPath":
		if spec, ok := v0Vol.Spec.(v0.HostPathVolumeSpec); ok {
			v1Vol.Uses = "host-path"
			v1Vol.With = &v1.VolumeHostPath{
				Path:      spec.Path,
				MountPath: v0Vol.MountPath,
				Type:      spec.Type,
			}
		}

	case "ConfigMap":
		if spec, ok := v0Vol.Spec.(v0.ConfigMapVolumeSpec); ok {
			v1Vol.Uses = "config-map"
			v1Vol.With = &v1.VolumeConfigMap{
				Name:      spec.Name,
				MountPath: v0Vol.MountPath,
				Optional:  spec.Optional,
			}
		}

	case "Secret":
		if spec, ok := v0Vol.Spec.(v0.SecretVolumeSpec); ok {
			v1Vol.Uses = "secret"
			v1Vol.With = &v1.VolumeSecret{
				Name:      spec.Name,
				MountPath: v0Vol.MountPath,
				Optional:  spec.Optional,
			}
		}

	default:
		return nil
	}

	return v1Vol
}
