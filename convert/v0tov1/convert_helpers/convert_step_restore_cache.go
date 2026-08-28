package converthelpers

import (
	"strings"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// ConvertStepRestoreCache converts a v0 RestoreCache step (generic Harness
// Cloud cache backend) to v1 template format.
func ConvertStepRestoreCache(src *v0.Step) *v1.StepTemplate {
	if src == nil {
		return nil
	}

	spec, ok := src.Spec.(*v0.StepRestoreCache)
	if !ok {
		return nil
	}

	// Create the with parameters map
	with := make(map[string]interface{})

	if spec.Key != "" {
		with["cachekey"] = spec.Key
	}

	if spec.ArchiveFormat != "" {
		with["archiveformat"] = strings.ToLower(spec.ArchiveFormat)
	}

	// v1 renamed failIfKeyNotFound -> fail_if_key_not_found
	if spec.FailIfKeyNotFound != nil {
		with["fail_if_key_not_found"] = spec.FailIfKeyNotFound
	}

	return &v1.StepTemplate{
		Uses:      "restoreCache",
		With:      with,
		Container: ConvertTemplateContainer(spec.RunAsUser, spec.Resources),
	}
}
