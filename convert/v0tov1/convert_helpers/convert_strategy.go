package converthelpers

import (
	"fmt"
	"strings"
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
	"github.com/drone/go-convert/internal/flexible"
)

func ConvertStrategy(src *v0.Strategy) *v1.Strategy {
	if src == nil {
		return nil
	}
	dst := &v1.Strategy{}
	if src.Matrix != nil {
		matrix, maxParallel := convertMatrix(src.Matrix)
		dst.Matrix = matrix
		dst.MaxParallel = maxParallel
	}
	if src.Repeat != nil {
		repeat, maxParallel := convertRepeat(src.Repeat)
		dst.Repeat = repeat
		if maxParallel != nil {
			dst.MaxParallel = maxParallel
		}
	}
    if src.Parallelism != nil {
        dst.For = &v1.For{
            Iterations: src.Parallelism,
        }
    }
	return dst
}

func convertMatrix(src map[string]interface{}) (*v1.Matrix, *flexible.Field[int64]) {
    if src == nil {
        return nil, nil
    }

    axis := make(map[string]interface{})
    exclude := make([]map[string]string, 0)
    var maxParallel *flexible.Field[int64]
    
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
            maxParallel = &flexible.Field[int64]{}
            if vNumber, ok := v.(float64); ok {
                maxParallel.Set(int64(vNumber))
            } else if vString, ok := v.(string); ok {
                maxParallel.SetString(vString)
            }
        default:
            // Handle different value types for matrix axis
            if vString, ok := v.(string); ok {
                // Single string value
                axis[k] = []string{vString}
            } else if vStringSlice, ok := v.([]string); ok {
                // Already a string slice
                axis[k] = vStringSlice
            } else if vInterfaceSlice, ok := v.([]interface{}); ok {
                // Convert []interface{} to []string (common from YAML parsing)
                stringSlice := make([]string, len(vInterfaceSlice))
                for i, item := range vInterfaceSlice {
                    stringSlice[i] = fmt.Sprintf("%v", item)
                }
                axis[k] = stringSlice
            } else {
                // Single value of any type, convert to string and wrap in array
                axis[k] = []string{fmt.Sprintf("%v", v)}
            }
        }
    }

    return &v1.Matrix{
        Axis:    axis,
        Exclude: exclude,
    }, maxParallel
}

// returns the converted repeat and max parallel
func convertRepeat(src *v0.Repeat) (*v1.Repeat, *flexible.Field[int64]) {
	if src == nil {
		return nil, nil
	}

	var maxParallel *flexible.Field[int64]
	if src.MaxConcurrency != nil && !src.MaxConcurrency.IsNil() {
		maxParallel = src.MaxConcurrency
	}
	dst := &v1.Repeat{
		Iterations:         src.Times,
		Items:         src.Items,
		Start:         src.Start,
		End:           src.End,
		Unit:          strings.ToLower(src.Unit),
		NodeName:      src.NodeName,
		PartitionSize: src.PartitionSize,
	}

	return dst, maxParallel
}