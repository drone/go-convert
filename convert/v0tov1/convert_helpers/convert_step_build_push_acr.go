package converthelpers

import (
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

	if spec.Repository != "" {
		with["image_name"] = spec.Repository
	}

	if spec.SubscriptionId != "" {
		with["subscriptionid"] = spec.SubscriptionId
	}

	if spec.Tags != nil {
		with["tags"] = spec.Tags
	}

	if spec.Caching != nil {
		with["caching"] = spec.Caching
	}

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

	return &v1.StepTemplate{
		Uses: "buildAndPushToACR",
		With: with,
	}
}
