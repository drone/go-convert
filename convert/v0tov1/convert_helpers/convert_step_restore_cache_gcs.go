package converthelpers

import (
	"strings"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// ConvertStepRestoreCacheGCS converts a v0 RestoreCacheGCS step to v1 template format
func ConvertStepRestoreCacheGCS(src *v0.Step) *v1.StepTemplate {
	if src == nil {
		return nil
	}

	spec, ok := src.Spec.(*v0.StepRestoreCacheGCS)
	if !ok {
		return nil
	}

	// Create the with parameters map
	with := make(map[string]interface{})

	if spec.ConnectorRef != "" {
		with["connector"] = spec.ConnectorRef
	}

	if spec.Bucket != "" {
		with["bucket"] = spec.Bucket
	}

	if spec.Key != "" {
		with["cachekey"] = spec.Key
	}

	if spec.ArchiveFormat != "" {
		// Convert to lowercase as per v1 spec
		with["archiveformat"] = strings.ToLower(spec.ArchiveFormat)
	}

	
	with["failIfKeyNotFound"] = spec.FailIfKeyNotFound
	

	return &v1.StepTemplate{
		Uses: "restoreCacheFromGCS",
		With: with,
	}
}
