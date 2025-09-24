package converthelpers

import (
	"fmt"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

func ConvertStrategy(src *v0.Strategy) *v1.Strategy {
	if src == nil {
		return nil
	}
	dst := &v1.Strategy{}
	if src.Matrix != nil {
		matrix, maxParallel := convertMatrix(src.Matrix)
		dst.Matrix = matrix
		dst.MaxParallel = int64(maxParallel)
	} else if src.Repeat != nil {
		repeat, maxParallel := convertRepeat(src.Repeat)
		forRepeat, ok := repeat.(*v1.For)
		if ok {
			dst.For = forRepeat
			dst.MaxParallel = int64(maxParallel)
		}
	}
	return dst
}

func convertMatrix(src map[string]interface{}) (*v1.Matrix, int) {
	if src == nil {
		return nil, 0
	}

	axis := make(map[string]interface{})
	exclude := make([]map[string]string, 0)
	maxParallel := 0
	for k, v := range src {
		switch k {
		case "exclude":
			// Handle exclude configurations
			if excludeList, ok := v.([]interface{}); ok {
				for _, excludeItem := range excludeList {
					if excludeMap, ok := excludeItem.(map[string]interface{}); ok {
						convertedExclude := make(map[string]string)
						for excludeKey, excludeValue := range excludeMap {
							convertedExclude[excludeKey] = fmt.Sprintf("%v", excludeValue)
						}
						exclude = append(exclude, convertedExclude)
					} else if excludeMap, ok := excludeItem.(map[string]string); ok {
						// Handle case where it's already map[string]string
						exclude = append(exclude, excludeMap)
					}
				}
			} else if excludeList, ok := v.([]map[string]string); ok {
				// Handle case where exclude is already []map[string]string
				exclude = excludeList
			}
		case "maxConcurrency":
			maxParallel = int(v.(float64))
		default:
			axis[k] = v
		}
	}

	return &v1.Matrix{
		Axis:    axis,
		Exclude: exclude,
	}, maxParallel
}

func convertRepeat(src *v0.Repeat) (interface{}, int) {
	if src == nil {
		return nil, 0
	}
	if src.Times != 0 {
		return v1.For{
			Iterations: int64(src.Times),
		}, 0
	}
	// TODO: Handle for repeat over items
	return src, 0
}