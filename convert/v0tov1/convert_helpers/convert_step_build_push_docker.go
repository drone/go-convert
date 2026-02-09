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

	// Create the with parameters map
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
	}

	// Handle baseImageConnectorRefs - can be string, []string, or expression
	if spec.BaseImageConnectorRefs != nil {
		// If it's a slice, take the first element
		switch v := spec.BaseImageConnectorRefs.(type) {
		case []interface{}:
			if len(v) > 0 {
				with["baseimageconnector"] = v[0]
			}
		case []string:
			if len(v) > 0 {
				with["baseimageconnector"] = v[0]
			}
		default:
			// Single value or expression
			with["baseimageconnector"] = spec.BaseImageConnectorRefs
		}
	}

	if len(spec.Env) > 0 {
		with["envvars"] = spec.Env
	}

	if spec.Dockerfile != "" {
		with["dockerfile"] = spec.Dockerfile
	}

	if spec.Context != "" {
		with["context"] = spec.Context
	}

	if len(spec.Labels) > 0 {
		with["labels"] = spec.Labels
	}

	if len(spec.BuildArgs) > 0 {
		with["buildargs"] = spec.BuildArgs
	}

	if spec.Target != "" {
		with["target"] = spec.Target
	}

	return &v1.StepTemplate{
		Uses: "buildAndPushToDocker",
		With: with,
	}
}
