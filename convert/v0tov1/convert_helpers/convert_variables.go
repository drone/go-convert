package converthelpers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	v0 "github.com/drone/go-convert/convert/harness/yaml"
	v1 "github.com/drone/go-convert/convert/v0tov1/yaml"
)

// ConvertVariables converts a list of v0 Variables to v1 Inputs.
func ConvertVariables(src []*v0.Variable) map[string]*v1.Input {
	if len(src) == 0 {
		return nil
	}

	dst := make(map[string]*v1.Input)

	for _, variable := range src {
		if variable == nil || variable.Name == "" {
			continue
		}

		v1Type := ConvertVariableType(variable.Type)

		input := &v1.Input{
			Type:        v1Type,
			Required:    variable.Required,
			Description: variable.Description,
		}

		if !IsEmptyValue(variable.Value) {
			valueStr, ok := variable.Value.(string)
			if ok {
				parsed := ParseInputExpression(valueStr)

				if parsed.Value != "" {
					input.Value = FormatValueForV1Type(parsed.Value, v1Type)
				}
				if parsed.DefaultVal != "" {
					input.Default = FormatValueForV1Type(parsed.DefaultVal, v1Type)
				}
				if len(parsed.Enum) > 0 {
					input.Enum = FormatEnumForV1Type(parsed.Enum, v1Type)
				}
				if parsed.Regex != "" {
					input.Pattern = parsed.Regex
				}
				if parsed.ExecutionInput {
					input.ExecutionInput = true
				}
			} else {
				input.Value = FormatValueForV1Type(variable.Value, v1Type)
			}
		}

		// YAML-level default is a fallback when inline .default() wasn't found
		if !IsEmptyValue(variable.Default) && IsEmptyValue(input.Default) {
			input.Default = FormatValueForV1Type(variable.Default, v1Type)
		}

		if variable.Type == "Secret" {
			input.Mask = true
		}

		dst[variable.Name] = input
	}

	return dst
}

type ParsedExpression struct {
	Value          string
	DefaultVal     string
	Enum           []string
	Regex          string
	ExecutionInput bool
}

var lastFuncPattern = regexp.MustCompile(`^(.*)\.(default|allowedValues|selectOneFrom|selectManyFrom|regex)\((.+)\)$`)

func ParseInputExpression(valueStr string) *ParsedExpression {
	result := &ParsedExpression{}

	for {
		// Try .executionInput() (no args)
		if strings.HasSuffix(valueStr, ".executionInput()") {
			result.ExecutionInput = true
			valueStr = strings.TrimSuffix(valueStr, ".executionInput()")
			continue
		}

		matches := lastFuncPattern.FindStringSubmatch(valueStr)
		if matches == nil {
			break
		}

		prefix := matches[1]
		funcName := matches[2]
		args := matches[3]

		switch funcName {
		case "default":
			result.DefaultVal = args
		case "allowedValues", "selectOneFrom", "selectManyFrom":
			result.Enum = ParseCommaSeparatedParams(args)
		case "regex":
			result.Regex = args
		}

		valueStr = prefix
	}

	result.Value = valueStr
	return result
}

func ParseCommaSeparatedParams(paramsStr string) []string {
	var params []string
	var current strings.Builder

	for i := 0; i < len(paramsStr); i++ {
		char := paramsStr[i]

		if char == ',' {
			param := strings.TrimSpace(current.String())
			if param != "" {
				params = append(params, param)
			}
			current.Reset()
		} else {
			current.WriteByte(char)
		}
	}

	param := strings.TrimSpace(current.String())
	if param != "" {
		params = append(params, param)
	}

	return params
}

func IsEmptyValue(v interface{}) bool {
	if v == nil {
		return true
	}
	if s, ok := v.(string); ok && s == "" {
		return true
	}
	return false
}

// FormatValueForV1Type coerces a value to match the declared v1 type.
// For string/secret types, ensures the result is always a string.
// For number type, attempts numeric conversion.
func FormatValueForV1Type(value interface{}, v1Type string) interface{} {
	if value == nil {
		return nil
	}
	switch v1Type {
	case "string", "secret":
		switch v := value.(type) {
		case string:
			return v
		case bool:
			return fmt.Sprintf("%v", v)
		case int:
			return strconv.Itoa(v)
		case int64:
			return strconv.FormatInt(v, 10)
		case float64:
			if v == float64(int64(v)) {
				return strconv.FormatInt(int64(v), 10)
			}
			return strconv.FormatFloat(v, 'f', -1, 64)
		default:
			return fmt.Sprintf("%v", v)
		}
	case "number":
		switch v := value.(type) {
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				return f
			}
			return v
		default:
			return v
		}
	default:
		return value
	}
}

// FormatEnumForV1Type coerces enum values to match the declared v1 type.
func FormatEnumForV1Type(params []string, v1Type string) interface{} {
	switch v1Type {
	case "number":
		floatParams := make([]float64, 0, len(params))
		for _, p := range params {
			if f, err := strconv.ParseFloat(p, 64); err == nil {
				floatParams = append(floatParams, f)
			}
		}
		if len(floatParams) > 0 {
			return floatParams
		}
		return params
	default:
		return params
	}
}

func ConvertVariableType(v0Type string) string {
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
		return "string"
	}
}
