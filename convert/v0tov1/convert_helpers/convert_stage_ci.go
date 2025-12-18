package converthelpers

import (
	"strings"
	"fmt"
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
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

// ConvertInfrastructure converts v0 Infrastructure to v1 Runtime
func ConvertInfrastructure(infra *v0.Infrastructure) *v1.Runtime {
	if infra == nil {
		return nil
	}

	// Handle KubernetesDirect type
	if strings.EqualFold(infra.Type, "KubernetesDirect") && infra.Spec != nil {
		return &v1.Runtime{
			Kubernetes: &v1.RuntimeKubernetes{
				Namespace: infra.Spec.Namespace,
				Connector: infra.Spec.Conn,
			},
		}
	}

	// Add other infrastructure types as needed
	return nil
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
            if src.Spec.Resources.Limits.CPU != nil {
                cpu = src.Spec.Resources.Limits.CPU.String() + "m"
            }
            if src.Spec.Resources.Limits.Memory != nil {
                memory = src.Spec.Resources.Limits.Memory.String()
            }
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