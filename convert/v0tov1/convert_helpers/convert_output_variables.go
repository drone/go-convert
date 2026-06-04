package converthelpers

import (
	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

func ConvertOutputVariables(src []*v0.Output) []*v1.Output {
	outputs := make([]*v1.Output, 0)
	for _, outputVar := range src {
		if outputVar == nil {
			continue
		}

		alias := outputVar.Name   // Harness key (required)
		name := outputVar.Value // Shell variable to capture (optional)
		if name == "" {
			name = alias
		}

		mask := false
		if outputVar.Type == "Secret" {
			mask = true
		}

		output := &v1.Output{
			Name:  name,
			Alias: alias,
			Mask:  mask,
		}
		outputs = append(outputs, output)
	}
	return outputs
}
