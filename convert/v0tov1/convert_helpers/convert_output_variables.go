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
		alias := outputVar.Name
		name := outputVar.Value
		if name == "" {
			name = outputVar.Name
		}
		mask := false
		if outputVar.Type == "Secret" {
			mask = true
		}

		output := &v1.Output{
			Alias:  alias,
			Mask:   mask,
			Name: name,
		}
		outputs = append(outputs, output)
	}
	return outputs
}