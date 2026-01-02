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
	
	disabled := !clone
	if !disabled {
		return nil
	}
	return &v1.Clone{
		Disabled: disabled,
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
func ConvertBuildIntelligence(bi *v0.BuildIntelligence) bool {
	if bi == nil {
		return false
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
			var user *flexible.Field[int]
			if k8sSpec.ContainerSecurityContext != nil {
				if ctx, ok := k8sSpec.ContainerSecurityContext.AsStruct(); ok && ctx != nil {
					user = ctx.RunAsUser
				}
			}
			var pull string
			if k8sSpec.ImagePullPolicy == "Always" {
				pull = "always"
			} else if k8sSpec.ImagePullPolicy == "Never" {
				pull = "never"
			} else if k8sSpec.ImagePullPolicy == "IfNotPresent" {
				pull = "if-not-exists"
			} else {
				pull = k8sSpec.ImagePullPolicy
			}
			return &v1.Runtime{
				Kubernetes: &v1.RuntimeKubernetes{
					Namespace:      k8sSpec.Namespace,
					Connector:      k8sSpec.Conn,
					Annotations:    k8sSpec.Annotations,
					Labels:         k8sSpec.Labels,
					Node:           k8sSpec.NodeSelector,
					OS:             k8sSpec.OS,
					ServiceAccount: k8sSpec.ServiceAccountName,
					ServiceToken:   k8sSpec.AutomountServiceAccountToken,
					Tolerations:    convertTolerations(k8sSpec.Tolerations),
					Timeout:        k8sSpec.InitTimeout,
					Host:           k8sSpec.HostNames,
					PriorityClass:  k8sSpec.PriorityClassName,
					User:           user,
					ImagePullPolicy: pull,
				},
			}
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
					},
				}
			}
		}
	}

	return nil
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
            Uses: "temp",
            With: &v1.VolumeTemp{
                Target: path,
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

// convertInfraVolume converts a single v0 infrastructure volume to v1 format
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
			v1Vol.Uses = "temp"
			v1Vol.With = &v1.VolumeTemp{
				Target: v0Vol.MountPath,
				Medium: spec.Medium,
				Limit: spec.Size,
			}
		}

	case "PersistentVolumeClaim":
		if spec, ok := v0Vol.Spec.(v0.PersistentVolumeClaimVolumeSpec); ok {
			v1Vol.Uses = "claim"
			v1Vol.With = &v1.VolumeClaim{
				Name: spec.ClaimName,
			}
		}

	case "HostPath":
		if spec, ok := v0Vol.Spec.(v0.HostPathVolumeSpec); ok {
			v1Vol.Uses = "bind"
			v1Vol.With = &v1.VolumeBind{
				Path: spec.Path,
			}
		}
	default:
		return nil
	}

	return v1Vol
}
