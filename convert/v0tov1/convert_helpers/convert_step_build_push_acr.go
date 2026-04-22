package converthelpers

import (
	"strings"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// ConvertStepBuildAndPushACR converts a v0 BuildAndPushACR step to v1 template format
func ConvertStepBuildAndPushACR(src *v0.Step) *v1.StepTemplate {
	if src == nil {
		return nil
	}

	spec, ok := src.Spec.(*v0.StepBuildAndPushACR)
	if !ok {
		return nil
	}

	// Create the with parameters map
	with := make(map[string]interface{})

	if spec.ConnectorRef != "" {
		with["connector"] = spec.ConnectorRef
	}

	// v0 Repository contains the full image path (e.g., myregistry.azurecr.io/myapp)
	// v1 splits it into registry (myregistry.azurecr.io) + image_name (myapp)
	if spec.Repository != "" {
		parts := strings.SplitN(spec.Repository, "/", 2)
		if len(parts) == 2 {
			// Format: registry/image_name
			with["registry"] = parts[0]
			with["image_name"] = parts[1]
		} else {
			// No slash found, treat entire value as image_name
			with["image_name"] = spec.Repository
		}
	}

	if spec.SubscriptionId != "" {
		with["subscriptionid"] = spec.SubscriptionId
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

	if spec.Env != nil {
		with["envvars"] = spec.Env
	}

	if spec.Labels != nil {
		with["labels"] = spec.Labels
	}

	if spec.BuildArgs != nil {
		with["buildargs"] = spec.BuildArgs
	}

	if spec.BaseImageConnectorRefs != nil {
		with["baseimageconnector"] = spec.BaseImageConnectorRefs
	}

	if spec.Dockerfile != "" {
		with["dockerfile"] = spec.Dockerfile
	}

	if spec.Context != "" {
		with["context"] = spec.Context
	}

	if spec.Target != "" {
		with["target"] = spec.Target
	}

	if spec.RemoteCacheImage != "" {
		with["remotecacheimage"] = spec.RemoteCacheImage
	}

	return &v1.StepTemplate{
		Uses: "buildAndPushToACR",
		With: with,
	}
}
