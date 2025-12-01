package converthelpers

import (
	"fmt"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// ConvertStepBuildAndPushECR converts a v0 BuildAndPushECR step to v1 template format
func ConvertStepBuildAndPushECR(src *v0.Step) *v1.StepTemplate {
	if src == nil {
		return nil
	}

	spec, ok := src.Spec.(*v0.StepBuildAndPushECR)
	if !ok {
		return nil
	}

	// Create the with parameters map
	with := make(map[string]interface{})

	if spec.ConnectorRef != "" {
		with["connector"] = spec.ConnectorRef
	}

	if spec.Region != "" {
		with["region"] = spec.Region
	}

	// Construct registry URL: <account_id>.dkr.ecr.<region>.amazonaws.com
	if spec.Account != "" && spec.Region != "" {
		registry := fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com", spec.Account, spec.Region)
		with["registry"] = registry
	}

	if spec.ImageName != "" {
		with["repo"] = spec.ImageName
	}

	if len(spec.Tags) > 0 {
		with["tags"] = spec.Tags
	}

	if spec.Caching {
		with["caching"] = spec.Caching
	}

	if len(spec.Env) > 0 {
		with["envvars"] = spec.Env
	}

	if len(spec.Labels) > 0 {
		with["labels"] = spec.Labels
	}

	if len(spec.BuildArgs) > 0 {
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
		Uses: "buildAndPushToECR",
		With: with,
	}
}
