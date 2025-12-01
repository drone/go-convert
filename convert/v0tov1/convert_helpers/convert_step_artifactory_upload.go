package converthelpers

import (
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// ConvertStepArtifactoryUpload converts a v0 ArtifactoryUpload step to v1 template format
func ConvertStepArtifactoryUpload(src *v0.Step) *v1.StepTemplate {
	if src == nil {
		return nil
	}

	spec, ok := src.Spec.(*v0.StepArtifactoryUpload)
	if !ok {
		return nil
	}

	// Create the with parameters map
	with := make(map[string]interface{})

	if spec.ConnRef != "" {
		with["connector"] = spec.ConnRef
	}

	if spec.Target != "" {
		with["target"] = spec.Target
	}

	if spec.SourcePath != "" {
		with["source"] = spec.SourcePath
	}

	return &v1.StepTemplate{
		Uses: "uploadArtifactsToJfrogArtifactory",
		With: with,
	}
}
