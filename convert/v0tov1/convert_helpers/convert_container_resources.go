package converthelpers

import (
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/drone/go-convert/internal/flexible"
)

// containerConfig holds the optional container fields configured via
// ContainerOption. runAsUser and resources remain required positional args.
type containerConfig struct {
	privileged      *flexible.Field[bool]
	imagePullPolicy string
}

// ContainerOption configures optional fields on the v1.Container built by
// ConvertTemplateContainer.
type ContainerOption func(*containerConfig)

// WithPrivileged sets the container privileged flag.
func WithPrivileged(v *flexible.Field[bool]) ContainerOption {
	return func(c *containerConfig) { c.privileged = v }
}

// WithImagePullPolicy sets the v0 image pull policy; it is converted to the v1
// pull format internally by ConvertTemplateContainer.
func WithImagePullPolicy(v string) ContainerOption {
	return func(c *containerConfig) { c.imagePullPolicy = v }
}

// ConvertTemplateContainer builds a v1.Container for template steps from v0
// RunAsUser and Resources fields plus any optional fields supplied via opts.
// Returns nil if no fields are set.
func ConvertTemplateContainer(runAsUser *flexible.Field[int], resources *v0.Resources, opts ...ContainerOption) *v1.Container {
	cfg := &containerConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	res := ConvertContainerResources(resources)
	pull := ConvertImagePullPolicy(cfg.imagePullPolicy)
	if runAsUser == nil && res == nil && cfg.privileged == nil && pull == "" {
		return nil
	}
	return &v1.Container{
		User:       runAsUser,
		Resources:  res,
		Privileged: cfg.privileged,
		Pull:       pull,
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
