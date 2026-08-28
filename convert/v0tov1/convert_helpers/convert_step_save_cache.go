package converthelpers

import (
	"strings"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// ConvertStepSaveCache converts a v0 SaveCache step (generic Harness Cloud
// cache backend) to v1 template format.
func ConvertStepSaveCache(src *v0.Step) *v1.StepTemplate {
	if src == nil {
		return nil
	}

	spec, ok := src.Spec.(*v0.StepSaveCache)
	if !ok {
		return nil
	}

	// Create the with parameters map
	with := make(map[string]interface{})

	if spec.Key != "" {
		with["cachekey"] = spec.Key
	}

	if spec.SourcePaths != nil {
		with["source_path"] = spec.SourcePaths
	}

	if spec.ArchiveFormat != "" {
		with["archiveformat"] = strings.ToLower(spec.ArchiveFormat)
	}

	if spec.Override != nil {
		with["override"] = spec.Override
	}

	return &v1.StepTemplate{
		Uses:      "saveCache",
		With:      with,
		Container: ConvertTemplateContainer(spec.RunAsUser, spec.Resources),
	}
}
