package converthelpers

import (
	"fmt"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	"github.com/drone/go-convert/convert/v0tov1/messagelog"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// ConvertStepBuildAndPushDockerRegistry converts a v0 BuildAndPushDockerRegistry step to v1 template format
func ConvertStepBuildAndPushDockerRegistry(src *v0.Step) *v1.StepTemplate {
	if src == nil {
		return nil
	}

	spec, ok := src.Spec.(*v0.StepBuildAndPushDockerRegistry)
	if !ok {
		return nil
	}

	with := make(map[string]interface{})

	var isHarnessRegistry bool = false
	if spec.ConnectorRef != "" {
		with["connector"] = spec.ConnectorRef
	} else if spec.RegistryRef != "" {
		isHarnessRegistry = true
	} else {
		messagelog.GetMessageLogger().LogWarning(
			"NO_CONNECTOR_OR_REGISTRY_REF",
			fmt.Sprintf("Connector or registryRef not provided in BuildAndPushDockerRegistry step: %s", src.ID),
			messagelog.WithStep(src.ID, src.Type),
		)
	}

	if spec.Repo != "" && !isHarnessRegistry {
		with["repo"] = spec.Repo
	} else if spec.Repo != "" && isHarnessRegistry {
		with["repo"] = fmt.Sprintf("%s/%s", spec.RegistryRef, spec.Repo)
	} else {
		messagelog.GetMessageLogger().LogWarning(
			"REPO_NOT_PROVIDED",
			fmt.Sprintf("Repo not provided in BuildAndPushDockerRegistry step: %s", src.ID),
			messagelog.WithStep(src.ID, src.Type),
		)
	}

	if spec.Tags != nil {
		with["tags"] = spec.Tags
	}

	if spec.Caching != nil {
		with["caching"] = spec.Caching
	} else {
		// caching is required in v1, default: true
		with["caching"] = true
	}

	// build_mode is required in v1, default: build_and_push
	with["build_mode"] = "build_and_push"

	// Handle baseImageConnectorRefs - can be string, []string, or expression
	if spec.BaseImageConnectorRefs != nil {
		switch v := spec.BaseImageConnectorRefs.(type) {
		case []interface{}:
			if len(v) > 0 {
				with["base_image_connector"] = v[0]
			}
		case []string:
			if len(v) > 0 {
				with["base_image_connector"] = v[0]
			}
		default:
			with["base_image_connector"] = spec.BaseImageConnectorRefs
		}
	}

	if spec.Env != nil {
		with["env_vars"] = spec.Env
	}

	if spec.Dockerfile != "" {
		with["dockerfile"] = spec.Dockerfile
	}

	if spec.Context != "" {
		with["context"] = spec.Context
	}

	if spec.Labels != nil {
		with["labels"] = spec.Labels
	}

	if spec.BuildArgs != nil {
		with["build_args"] = spec.BuildArgs
	}

	if spec.Target != "" {
		with["target"] = spec.Target
	}

	if spec.Optimize != nil {
		with["optimize"] = spec.Optimize
	}

	if spec.RemoteCacheRepo != "" {
		with["cache_repo"] = spec.RemoteCacheRepo
	}

	if spec.CacheFrom != nil {
		with["cache_from"] = spec.CacheFrom
	}

	if spec.CacheTo != "" {
		with["cache_to"] = spec.CacheTo
	}

	var uses string
	if isHarnessRegistry {
		uses = "buildAndPushToHAR"
	} else {
		uses = "buildAndPushToDocker"
	}

	return &v1.StepTemplate{
		Uses:      uses,
		With:      with,
		Container: ConvertTemplateContainer(spec.RunAsUser, spec.Resources),
	}
}
