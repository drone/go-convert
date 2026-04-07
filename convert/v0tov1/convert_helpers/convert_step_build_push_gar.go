package converthelpers

import (
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// ConvertStepBuildAndPushGCR converts a v0 BuildAndPushGCR/BuildAndPushGAR step to v1 template format
func ConvertStepBuildAndPushGAR(src *v0.Step) *v1.StepTemplate {
	if src == nil {
		return nil
	}

	spec, ok := src.Spec.(*v0.StepBuildAndPushGAR)
	if !ok {
		return nil
	}

	// Create the with parameters map
	with := make(map[string]interface{})

	if spec.ConnectorRef != "" {
		with["connector"] = spec.ConnectorRef
	}

	// v1 uses separate host + project_id fields instead of computed registry URL
	if spec.Host != "" {
		with["host"] = spec.Host
	}

	if spec.ProjectID != "" {
		with["project_id"] = spec.ProjectID
	}

	// v1 renamed repo → image_name
	if spec.ImageName != "" {
		with["image_name"] = spec.ImageName
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
		Uses: "buildAndPushToGAR",
		With: with,
	}
}
