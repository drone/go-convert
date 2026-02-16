package flexible

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Field is a generic type that can hold either a value of type T or a string expression
type Field[T any] struct {
	Value interface{}
}

// UnmarshalJSON implements json.Unmarshaler for automatic handling of multiple types
func (f *Field[T]) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as struct T first
	var structValue T
	if err := json.Unmarshal(data, &structValue); err == nil {
		f.Value = structValue
		return nil
	}

	// Fall back to string
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("failed to unmarshal as struct or string: %v", err)
	}

	f.Value = str
	return nil
}

// MarshalJSON implements json.Marshaler for proper serialization
func (f Field[T]) MarshalJSON() ([]byte, error) {
	if f.Value == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(f.Value)
}

// MarshalYAML implements yaml.Marshaler for proper YAML serialization
func (f Field[T]) MarshalYAML() (interface{}, error) {
	if f.Value == nil {
		return nil, nil
	}
	return f.Value, nil
}

// UnmarshalYAML implements yaml.Unmarshaler for YAML deserialization
func (f *Field[T]) UnmarshalYAML(node *yaml.Node) error {
	// Try to unmarshal as struct T first
	var structValue T
	if err := node.Decode(&structValue); err == nil {
		f.Value = structValue
		return nil
	}

	// Fall back to string
	var str string
	if err := node.Decode(&str); err != nil {
		return fmt.Errorf("failed to unmarshal as struct or string: %v", err)
	}

	f.Value = str
	return nil
}

// IsExpression returns true if the field contains a Harness expression string (contains <+ anywhere)
func (f *Field[T]) IsExpression() bool {
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

// AsStruct returns the value as struct T, or zero value and false if it's a string
func (f *Field[T]) AsStruct() (T, bool) {
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
func (f *Field[T]) AsString() (string, bool) {
	if f.Value == nil {
		return "", false
	}

	if str, ok := f.Value.(string); ok {
		return str, true
	}

	return "", false
}

// Set sets the field to a struct value
func (f *Field[T]) Set(value T) {
	f.Value = value
}

// SetString sets the field to a string value
func (f *Field[T]) SetString(value string) {
	f.Value = value
}

// SetExpression sets the field to a Harness expression string
func (f *Field[T]) SetExpression(expr string) {
	f.Value = expr
}

// IsNil returns true if the field is nil/empty
func (f *Field[T]) IsNil() bool {
	return f.Value == nil
}

// NegateBool negates a boolean flexible field.
// For struct values, it negates the boolean directly.
// For expressions, it wraps the expression with <+!...> to negate it.
// Returns nil if the input field is nil.
func NegateBool(field *Field[bool]) *Field[bool] {
	if field == nil {
		return nil
	}

	result := &Field[bool]{}

	if val, ok := field.AsStruct(); ok {
		// It's a boolean value - negate it
		result.Set(!val)
		return result
	}

	if expr, ok := field.AsString(); ok {
		// It's an expression - wrap with negation
		modifiedExpr := "<+!" + expr + ">"
		result.SetExpression(modifiedExpr)
		return result
	}

	return nil
}

// // Convert applies a conversion function directly to the field, preserving expressions
// func Convert[From, To any](field *Field[From], converter func(From) To) *Field[To] {
// 	if field == nil {
// 		return nil
// 	}

// 	result := &Field[To]{}

// 	if field.IsExpression() {
// 		result.SetExpression(field.AsString())
// 		return result
// 	}

// 	if structValue, ok := field.AsStruct(); ok {
// 		converted := converter(structValue)
// 		result.Set(converted)
// 		return result
// 	}

// 	// For other types, pass through unchanged if compatible
// 	result.Value = field.Value
// 	return result
// }

// // ConvertWithError applies a conversion function that returns an error
// func ConvertWithError[From, To any](field *Field[From], converter func(From) (To, error)) (*Field[To], error) {
// 	if field == nil {
// 		return nil, nil
// 	}

// 	result := &Field[To]{}

// 	if field.IsExpression() {
// 		result.SetExpression(field.AsString())
// 		return result, nil
// 	}

// 	if structValue, ok := field.AsStruct(); ok {
// 		converted, err := converter(structValue)
// 		if err != nil {
// 			return nil, err
// 		}
// 		result.Set(converted)
// 		return result, nil
// 	}

// 	// For other types, pass through unchanged
// 	result.Value = field.Value
// 	return result, nil
// }

// // CallWith calls an existing function with unwrapped values, then rewraps the result
// func CallWith[From, To any](field *Field[From], fn func(From) To) *Field[To] {
// 	if field == nil {
// 		return nil
// 	}

// 	result := &Field[To]{}

// 	if field.IsExpression() {
// 		result.SetExpression(field.AsString())
// 		return result
// 	}

// 	if structValue, ok := field.AsStruct(); ok {
// 		converted := fn(structValue)
// 		result.Set(converted)
// 		return result
// 	}

// 	// For other types, pass through unchanged
// 	result.Value = field.Value
// 	return result
// }

// // CallWithError calls an existing function that returns an error
// func CallWithError[From, To any](field *Field[From], fn func(From) (To, error)) (*Field[To], error) {
// 	if field == nil {
// 		return nil, nil
// 	}

// 	result := &Field[To]{}

// 	if field.IsExpression() {
// 		result.SetExpression(field.AsString())
// 		return result, nil
// 	}

// 	if structValue, ok := field.AsStruct(); ok {
// 		converted, err := fn(structValue)
// 		if err != nil {
// 			return nil, err
// 		}
// 		result.Set(converted)
// 		return result, nil
// 	}

// 	// For other types, pass through unchanged
// 	result.Value = field.Value
// 	return result, nil
// }

// // Wrap creates a FlexibleField from a regular value
// func Wrap[T any](value T) *Field[T] {
// 	field := &Field[T]{}
// 	field.Set(value)
// 	return field
// }

// // WrapExpression creates a FlexibleField from an expression string
// func WrapExpression[T any](expr string) *Field[T] {
// 	field := &Field[T]{}
// 	field.SetExpression(expr)
// 	return field
// }

// // ConvertSlice converts a slice of FlexibleFields using a conversion function
// func ConvertSlice[From, To any](fields []*Field[From], converter func(From) To) []*Field[To] {
// 	if fields == nil {
// 		return nil
// 	}

// 	result := make([]*Field[To], len(fields))
// 	for i, field := range fields {
// 		result[i] = Convert(field, converter)
// 	}
// 	return result
// }

// // ConvertMap converts a map of FlexibleFields using a conversion function
// func ConvertMap[From, To any](fields map[string]*Field[From], converter func(From) To) map[string]*Field[To] {
// 	if fields == nil {
// 		return nil
// 	}

// 	result := make(map[string]*Field[To])
// 	for key, field := range fields {
// 		result[key] = Convert(field, converter)
// 	}
// 	return result
// }
