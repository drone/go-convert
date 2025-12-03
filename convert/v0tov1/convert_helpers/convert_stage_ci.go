package converthelpers

import (
	"strings"

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
