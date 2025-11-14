package pipelineconverter

import (
	"regexp"
	"strconv"
	"strings"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	convert_helpers "github.com/drone/go-convert/convert/v0tov1/convert_helpers"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

type PipelineConverter struct {}

func NewPipelineConverter() *PipelineConverter {
	return &PipelineConverter{}
}

// ConvertPipeline converts a v0 Pipeline to a v1 Pipeline.
func (c *PipelineConverter) ConvertPipeline(src *v0.Pipeline) *v1.Pipeline {
	if src == nil {
		return nil
	}

	var barriers []string
	if src.FlowControl != nil {
		barriers = convertBarriers(src.FlowControl.Barriers)
	}

	inputs := c.convertVariables(src.Variables)
	stages := c.convertStages(src.Stages)

	repo,clone := c.convertCodebase(src.Props.CI.Codebase)
	dst := &v1.Pipeline{
		Id:            src.ID,
		Name:          src.Name,
		Inputs:        inputs,
		Stages:        stages,
		Barriers:      barriers,
		Clone:         clone,
		Notifications: convert_helpers.ConvertNotifications(src.NotificationRules),
		Repo:          repo,
	}

	return dst
}

func (c *PipelineConverter) convertCodebase(src v0.Codebase) (*v1.Repository, *v1.Clone) {
	if src.Conn == "" && src.Name == "" {
		return nil, &v1.Clone{
			Disabled: true,
		}
	}
	repo := &v1.Repository{
		Name:  src.Name,
		Connector:  src.Conn,
	}
	clone := &v1.Clone{
		Disabled: false,
	}
	return repo, clone
}

// convertVariables converts a list of v0 Variables to v1 Inputs.
func (c *PipelineConverter) convertVariables(src []*v0.Variable) map[string]*v1.Input {
	if len(src) == 0 {
		return nil
	}

	dst := make(map[string]*v1.Input)

	// Convert variables to inputs
	for _, variable := range src {
		if variable == nil || variable.Name == "" {
			continue
		}

		input := &v1.Input{
			Type:     convertVariableType(variable.Type),
			Required: variable.Required,
		}
		if variable.Default != "" {
			input.Default = variable.Default
		}

		if variable.Value != "" {
			parsedValue, enumValues := parseVariableAllowedValues(variable.Value, variable.Type)
			input.Value = parsedValue
			// Only set Enum if we have actual values
			if enumValues != nil {
				input.Enum = enumValues
			}
		}

		// Set mask to true for secret types
		if variable.Type == "Secret" {
			input.Mask = true
		}

		dst[variable.Name] = input
	}

	return dst
}

func parseVariableAllowedValues(value interface{}, var_type string) (interface{}, interface{}) {
	// Convert value to string for parsing
	valueStr, ok := value.(string)
	if !ok {
		// If it's already a number, return as-is
		if var_type == "Number" {
			if floatVal, ok := value.(float64); ok {
				return floatVal, nil
			}
		}
		return value, nil
	}

	// Check for special functions: allowedValues(), selectOneFrom(), selectManyFrom()
	funcPattern := regexp.MustCompile(`^(.+?)\.(allowedValues|selectOneFrom|selectManyFrom)\((.+)\)$`)
	matches := funcPattern.FindStringSubmatch(valueStr)

	if len(matches) != 4 {
		// No function found, return value as-is
		if var_type == "Number" {
			if floatVal, err := strconv.ParseFloat(valueStr, 64); err == nil {
				return floatVal, nil
			}
		}
		return valueStr, nil
	}

	// Extract the value before the function
	extractedValue := matches[1]
	// Extract the parameters string
	paramsStr := matches[3]

	// Parse the comma-separated parameters, handling nested parentheses
	params := parseCommaSeparatedParams(paramsStr)

	// Convert extracted value to appropriate type
	var finalValue interface{} = extractedValue
	if var_type == "Number" {
		if floatVal, err := strconv.ParseFloat(extractedValue, 64); err == nil {
			finalValue = floatVal
		}
	}

	// Convert params to appropriate type based on var_type
	if var_type == "Number" {
		// Convert params to array of floats
		floatParams := make([]float64, 0, len(params))
		for _, param := range params {
			if floatVal, err := strconv.ParseFloat(param, 64); err == nil {
				floatParams = append(floatParams, floatVal)
			}
		}
		if len(floatParams) > 0 {
			return finalValue, floatParams
		}
		return finalValue, nil
	}

	// Return string array for non-number types
	if len(params) > 0 {
		return finalValue, params
	}
	return finalValue, nil
}

// parseCommaSeparatedParams parses comma-separated parameters, handling nested parentheses
func parseCommaSeparatedParams(paramsStr string) []string {
	var params []string
	var current strings.Builder

	for i := 0; i < len(paramsStr); i++ {
		char := paramsStr[i]

		if char == ',' {
			// Comma is always a separator
			param := strings.TrimSpace(current.String())
			if param != "" {
				params = append(params, param)
			}
			current.Reset()
		} else {
			// Add any other character (including parentheses) to current parameter
			current.WriteByte(char)
		}
	}

	// Add the last parameter
	param := strings.TrimSpace(current.String())
	if param != "" {
		params = append(params, param)
	}

	return params
}

// convertVariableType converts v0 variable type to v1 input type.
func convertVariableType(v0Type string) string {
	switch v0Type {
	case "Secret":
		return "secret"
	case "Text":
		return "string"
	case "String":
		return "string"
	case "Number":
		return "number"
	default:
		return "string" // Default to string type
	}
}

// convertBarriers converts a list of v0 Barriers to v1 Barriers.
func convertBarriers(src []*v0.Barrier) []string {
	if len(src) == 0 {
		return nil
	}
	dst := make([]string, 0, len(src))
	for _, barrier := range src {
		if barrier == nil {
			continue
		}
		dst = append(dst, barrier.Name)
	}
	return dst
}
