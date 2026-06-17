package converthelpers

import (
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/drone/go-convert/internal/flexible"
)

// ConvertTemplateContainer builds a v1.Container for template steps
// from v0 RunAsUser and Resources fields. Returns nil if both are empty.
func ConvertTemplateContainer(runAsUser *flexible.Field[int], resources *v0.Resources) *v1.Container {
	res := ConvertContainerResources(resources)
	if runAsUser == nil && res == nil {
		return nil
	}
	return &v1.Container{
		User:      runAsUser,
		Resources: res,
	}
}

func ConvertContainerResources(resources *v0.Resources) *v1.ContainerResources {
	if resources == nil {
		return nil
	}

	result := &v1.ContainerResources{
		Limits:   ConvertResourceSpec(resources.Limits),
		Requests: ConvertResourceSpec(resources.Requests),
	}

	if result.Limits == nil && result.Requests == nil {
		return nil
	}

	return result
}

func ConvertResourceSpec(spec *v0.ResourceSpec) *v1.ContainerResourcesSpec {
	if spec == nil {
		return nil
	}

	cpu := spec.GetCPUString()
	memory := spec.GetMemoryString()

	if cpu == "" && memory == "" {
		return nil
	}

	return &v1.ContainerResourcesSpec{
		Cpu:    cpu,
		Memory: memory,
	}
}
