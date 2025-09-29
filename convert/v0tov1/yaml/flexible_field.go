package yaml

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// FlexibleField is a generic type that can hold either a struct of type T or a value (int, bool, float, string, struct)
// It can also detect Harness expressions containing <+ anywhere in the string
type FlexibleField[T any] struct {
	Value interface{}
}

// UnmarshalJSON implements json.Unmarshaler for automatic handling of multiple types
func (f *FlexibleField[T]) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as different types in order of preference

	// Try string first
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		f.Value = str
		return nil
	}

	// Try int
	var intValue int
	if err := json.Unmarshal(data, &intValue); err == nil {
		f.Value = intValue
		return nil
	}

	// Try float64
	var floatValue float64
	if err := json.Unmarshal(data, &floatValue); err == nil {
		f.Value = floatValue
		return nil
	}

	// Try bool
	var boolValue bool
	if err := json.Unmarshal(data, &boolValue); err == nil {
		f.Value = boolValue
		return nil
	}

	// Finally try to unmarshal as struct T
	var structValue T
	if err := json.Unmarshal(data, &structValue); err != nil {
		return fmt.Errorf("failed to unmarshal as string, int, float, bool, or struct: %v", err)
	}

	f.Value = structValue
	return nil
}

// MarshalJSON implements json.Marshaler for proper serialization
func (f FlexibleField[T]) MarshalJSON() ([]byte, error) {
	if f.Value == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(f.Value)
}

// MarshalYAML implements yaml.Marshaler for proper YAML serialization
func (f FlexibleField[T]) MarshalYAML() (interface{}, error) {
	if f.Value == nil {
		return nil, nil
	}
	return f.Value, nil
}

// UnmarshalYAML implements yaml.Unmarshaler for YAML deserialization
func (f *FlexibleField[T]) UnmarshalYAML(node *yaml.Node) error {
	// Try to unmarshal as different types in order of preference

	// Try string first
	var str string
	if err := node.Decode(&str); err == nil {
		f.Value = str
		return nil
	}

	// Try int
	var intValue int
	if err := node.Decode(&intValue); err == nil {
		f.Value = intValue
		return nil
	}

	// Try float64
	var floatValue float64
	if err := node.Decode(&floatValue); err == nil {
		f.Value = floatValue
		return nil
	}

	// Try bool
	var boolValue bool
	if err := node.Decode(&boolValue); err == nil {
		f.Value = boolValue
		return nil
	}

	// Finally try to unmarshal as struct T
	var structValue T
	if err := node.Decode(&structValue); err != nil {
		return fmt.Errorf("failed to unmarshal as string, int, float, bool, or struct: %v", err)
	}

	f.Value = structValue
	return nil
}

// IsExpression returns true if the field contains a Harness expression string (contains <+ anywhere)
func (f *FlexibleField[T]) IsExpression() bool {
	if f.Value == nil {
		return false
	}

	str, ok := f.Value.(string)
	if !ok {
		return false
	}

	// Check for <+ anywhere in the string
	return strings.Contains(str, "<+")
}

// AsInt returns the value as int, or zero and false if it's not an int
func (f *FlexibleField[T]) AsInt() (int, bool) {
	if f.Value == nil {
		return 0, false
	}

	if intValue, ok := f.Value.(int); ok {
		return intValue, true
	}

	return 0, false
}

// AsFloat returns the value as float64, or zero and false if it's not a float
func (f *FlexibleField[T]) AsFloat() (float64, bool) {
	if f.Value == nil {
		return 0, false
	}

	if floatValue, ok := f.Value.(float64); ok {
		return floatValue, true
	}

	return 0, false
}

// AsBool returns the value as bool, or false and false if it's not a bool
func (f *FlexibleField[T]) AsBool() (bool, bool) {
	if f.Value == nil {
		return false, false
	}

	if boolValue, ok := f.Value.(bool); ok {
		return boolValue, true
	}

	return false, false
}

// AsStruct returns the value as struct T, or zero value and false if it's a string
func (f *FlexibleField[T]) AsStruct() (T, bool) {
	var zero T
	if f.Value == nil {
		return zero, false
	}

	if structValue, ok := f.Value.(T); ok {
		return structValue, true
	}

	return zero, false
}

// AsString returns the value as string, or empty string if it's a struct
func (f *FlexibleField[T]) AsString() string {
	if f.Value == nil {
		return ""
	}

	if str, ok := f.Value.(string); ok {
		return str
	}

	return ""
}

// Set sets the field to a struct value
func (f *FlexibleField[T]) Set(value T) {
	f.Value = value
}

// SetInt sets the field to an int value
func (f *FlexibleField[T]) SetInt(value int) {
	f.Value = value
}

// SetFloat sets the field to a float64 value
func (f *FlexibleField[T]) SetFloat(value float64) {
	f.Value = value
}

// SetBool sets the field to a bool value
func (f *FlexibleField[T]) SetBool(value bool) {
	f.Value = value
}

// SetString sets the field to a string value
func (f *FlexibleField[T]) SetString(value string) {
	f.Value = value
}

// SetExpression sets the field to a Harness expression string
func (f *FlexibleField[T]) SetExpression(expr string) {
	f.Value = expr
}

// IsNil returns true if the field is nil/empty
func (f *FlexibleField[T]) IsNil() bool {
	return f.Value == nil
}
