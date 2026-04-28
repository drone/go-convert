package converthelpers

import (
	v0 "github.com/drone/go-convert/convert/harness/yaml"
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

	if spec.ConnectorRef != "" {
		with["connector"] = spec.ConnectorRef
	}

	if spec.Repo != "" {
		with["repo"] = spec.Repo
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

	return &v1.StepTemplate{
		Uses: "buildAndPushToDocker",
		With: with,
	}
}
