package converthelpers

import (
	"strings"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// ConvertStepRestoreCacheS3 converts a v0 RestoreCacheS3 step to v1 template format
func ConvertStepRestoreCacheS3(src *v0.Step) *v1.StepTemplate {
	if src == nil {
		return nil
	}

	spec, ok := src.Spec.(*v0.StepRestoreCacheS3)
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

	if spec.Bucket != "" {
		with["bucket"] = spec.Bucket
	}

	if spec.Key != "" {
		with["cachekey"] = spec.Key
	}

	if spec.Endpoint != "" {
		with["endpoint"] = spec.Endpoint
	}

	if spec.ArchiveFormat != "" {
		// Convert to lowercase as per v1 spec
		with["archiveformat"] = strings.ToLower(spec.ArchiveFormat)
	}

	with["pathstyle"] = spec.PathStyle
	
	with["failIfKeyNotFound"] = spec.FailIfKeyNotFound

	return &v1.StepTemplate{
		Uses: "restoreCacheFromS3",
		With: with,
	}
}
