package flexible

import (
	"encoding/json"
	"testing"
	"gopkg.in/yaml.v3"
)

// Test struct for testing Field with struct types
type TestStruct struct {
	Name  string `json:"name" yaml:"name"`
	Value int    `json:"value" yaml:"value"`
}

type TestDelegate struct {
	Filter  []string `json:"filter,omitempty" yaml:"filter,omitempty"`
	Inherit bool     `json:"inherit-from-infrastructure,omitempty" yaml:"inherit-from-infrastructure,omitempty"`
}

// TestNilValueHandling tests behavior with nil/empty fields
func TestNilValueHandling(t *testing.T) {
	t.Run("IsNil returns true for uninitialized field", func(t *testing.T) {
		var f Field[TestStruct]
		if !f.IsNil() {
			t.Error("Expected IsNil() to return true for uninitialized field")
		}
	})

	t.Run("IsNil returns true for explicitly nil field", func(t *testing.T) {
		f := Field[TestStruct]{Value: nil}
		if !f.IsNil() {
			t.Error("Expected IsNil() to return true for nil field")
		}
	})

	t.Run("IsNil returns false after setting value", func(t *testing.T) {
		f := Field[TestStruct]{}
		f.Set(TestStruct{Name: "test", Value: 42})
		if f.IsNil() {
			t.Error("Expected IsNil() to return false after setting value")
		}
	})

	t.Run("IsNil returns false after setting string", func(t *testing.T) {
		f := Field[TestStruct]{}
		f.SetString("test")
		if f.IsNil() {
			t.Error("Expected IsNil() to return false after setting string")
		}
	})

	t.Run("AsStruct returns zero value and false for nil field", func(t *testing.T) {
		var f Field[TestStruct]
		val, ok := f.AsStruct()
		if ok {
			t.Error("Expected AsStruct() to return false for nil field")
		}
		if val.Name != "" || val.Value != 0 {
			t.Error("Expected AsStruct() to return zero value for nil field")
		}
	})

	t.Run("AsString returns empty and false for nil field", func(t *testing.T) {
		var f Field[TestStruct]
		val, ok := f.AsString()
		if ok {
			t.Error("Expected AsString() to return false for nil field")
		}
		if val != "" {
			t.Error("Expected AsString() to return empty string for nil field")
		}
	})

	t.Run("IsExpression returns false for nil field", func(t *testing.T) {
		var f Field[TestStruct]
		if f.IsExpression() {
			t.Error("Expected IsExpression() to return false for nil field")
		}
	})
}

// TestJSONMarshaling tests JSON marshaling and unmarshaling
func TestJSONMarshaling(t *testing.T) {
	t.Run("Marshal nil field to null", func(t *testing.T) {
		var f Field[TestStruct]
		data, err := json.Marshal(f)
		if err != nil {
			t.Fatalf("Failed to marshal nil field: %v", err)
		}
		if string(data) != "null" {
			t.Errorf("Expected 'null', got '%s'", string(data))
		}
	})

	t.Run("Marshal struct value", func(t *testing.T) {
		f := Field[TestStruct]{}
		f.Set(TestStruct{Name: "test", Value: 42})
		data, err := json.Marshal(f)
		if err != nil {
			t.Fatalf("Failed to marshal struct: %v", err)
		}
		expected := `{"name":"test","value":42}`
		if string(data) != expected {
			t.Errorf("Expected '%s', got '%s'", expected, string(data))
		}
	})

	t.Run("Marshal string value", func(t *testing.T) {
		f := Field[TestStruct]{}
		f.SetString("test-string")
		data, err := json.Marshal(f)
		if err != nil {
			t.Fatalf("Failed to marshal string: %v", err)
		}
		expected := `"test-string"`
		if string(data) != expected {
			t.Errorf("Expected '%s', got '%s'", expected, string(data))
		}
	})

	t.Run("Unmarshal struct value", func(t *testing.T) {
		jsonData := `{"name":"test","value":42}`
		var f Field[TestStruct]
		err := json.Unmarshal([]byte(jsonData), &f)
		if err != nil {
			t.Fatalf("Failed to unmarshal struct: %v", err)
		}
		val, ok := f.AsStruct()
		if !ok {
			t.Error("Expected AsStruct() to return true")
		}
		if val.Name != "test" || val.Value != 42 {
			t.Errorf("Expected {Name:test Value:42}, got %+v", val)
		}
	})

	t.Run("Unmarshal string value", func(t *testing.T) {
		jsonData := `"test-string"`
		var f Field[TestStruct]
		err := json.Unmarshal([]byte(jsonData), &f)
		if err != nil {
			t.Fatalf("Failed to unmarshal string: %v", err)
		}
		val, ok := f.AsString()
		if !ok {
			t.Error("Expected AsString() to return true")
		}
		if val != "test-string" {
			t.Errorf("Expected 'test-string', got '%s'", val)
		}
	})
}

// TestYAMLMarshaling tests YAML marshaling and unmarshaling
func TestYAMLMarshaling(t *testing.T) {
	t.Run("Marshal nil field to null", func(t *testing.T) {
		var f Field[TestStruct]
		data, err := yaml.Marshal(f)
		if err != nil {
			t.Fatalf("Failed to marshal nil field: %v", err)
		}
		if string(data) != "null\n" {
			t.Errorf("Expected 'null\\n', got '%s'", string(data))
		}
	})

	t.Run("Marshal struct value", func(t *testing.T) {
		f := Field[TestStruct]{}
		f.Set(TestStruct{Name: "test", Value: 42})
		data, err := yaml.Marshal(f)
		if err != nil {
			t.Fatalf("Failed to marshal struct: %v", err)
		}
		// Unmarshal to verify structure
		var result TestStruct
		err = yaml.Unmarshal(data, &result)
		if err != nil {
			t.Fatalf("Failed to unmarshal result: %v", err)
		}
		if result.Name != "test" || result.Value != 42 {
			t.Errorf("Expected {Name:test Value:42}, got %+v", result)
		}
	})

	t.Run("Marshal string value", func(t *testing.T) {
		f := Field[TestStruct]{}
		f.SetString("test-string")
		data, err := yaml.Marshal(f)
		if err != nil {
			t.Fatalf("Failed to marshal string: %v", err)
		}
		expected := "test-string\n"
		if string(data) != expected {
			t.Errorf("Expected '%s', got '%s'", expected, string(data))
		}
	})

	t.Run("Unmarshal struct value", func(t *testing.T) {
		yamlData := "name: test\nvalue: 42\n"
		var f Field[TestStruct]
		err := yaml.Unmarshal([]byte(yamlData), &f)
		if err != nil {
			t.Fatalf("Failed to unmarshal struct: %v", err)
		}
		val, ok := f.AsStruct()
		if !ok {
			t.Error("Expected AsStruct() to return true")
		}
		if val.Name != "test" || val.Value != 42 {
			t.Errorf("Expected {Name:test Value:42}, got %+v", val)
		}
	})

	t.Run("Unmarshal string value", func(t *testing.T) {
		yamlData := "test-string\n"
		var f Field[TestStruct]
		err := yaml.Unmarshal([]byte(yamlData), &f)
		if err != nil {
			t.Fatalf("Failed to unmarshal string: %v", err)
		}
		val, ok := f.AsString()
		if !ok {
			t.Error("Expected AsString() to return true")
		}
		if val != "test-string" {
			t.Errorf("Expected 'test-string', got '%s'", val)
		}
	})
}

// TestIsExpression tests expression detection
func TestIsExpression(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"Simple expression", "<+input>", true},
		{"Pipeline variable", "<+pipeline.variables.myvar>", true},
		{"Stage variable", "<+pipeline.stages.build.variables.version>", true},
		{"Expression in middle", "prefix-<+var>-suffix", true},
		{"Multiple expressions", "<+var1>-<+var2>", true},
		{"Not an expression - no brackets", "input", false},
		{"Not an expression - only <", "<input", false},
		{"Not an expression - only +", "+input", false},
		{"Not an expression - reversed", "+<input>", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Field[TestStruct]{}
			f.SetString(tt.value)
			result := f.IsExpression()
			if result != tt.expected {
				t.Errorf("IsExpression() for '%s': expected %v, got %v", tt.value, tt.expected, result)
			}
		})
	}

	t.Run("IsExpression returns false for struct value", func(t *testing.T) {
		f := Field[TestStruct]{}
		f.Set(TestStruct{Name: "test", Value: 42})
		if f.IsExpression() {
			t.Error("Expected IsExpression() to return false for struct value")
		}
	})

	t.Run("IsExpression returns false for nil value", func(t *testing.T) {
		var f Field[TestStruct]
		if f.IsExpression() {
			t.Error("Expected IsExpression() to return false for nil value")
		}
	})
}

// TestPointerFieldWithOmitempty tests pointer fields with omitempty in parent structs
func TestPointerFieldWithOmitempty(t *testing.T) {
	type ParentStruct struct {
		Name     string               `json:"name" yaml:"name"`
		Delegate *Field[TestDelegate] `json:"delegate,omitempty" yaml:"delegate,omitempty"`
		Config   *Field[TestStruct]   `json:"config,omitempty" yaml:"config,omitempty"`
	}

	t.Run("JSON omits nil pointer field", func(t *testing.T) {
		parent := ParentStruct{Name: "test"}
		data, err := json.Marshal(parent)
		if err != nil {
			t.Fatalf("Failed to marshal: %v", err)
		}
		expected := `{"name":"test"}`
		if string(data) != expected {
			t.Errorf("Expected '%s', got '%s'", expected, string(data))
		}
	})

	t.Run("JSON includes non-nil pointer field with struct", func(t *testing.T) {
		delegate := &Field[TestDelegate]{}
		delegate.Set(TestDelegate{Filter: []string{"selector1"}})
		parent := ParentStruct{Name: "test", Delegate: delegate}
		data, err := json.Marshal(parent)
		if err != nil {
			t.Fatalf("Failed to marshal: %v", err)
		}
		// Verify delegate is present
		expected := `{"name":"test","delegate":{"filter":["selector1"]}}`
		if string(data) != expected {
			t.Errorf("Expected '%s', got '%s'", expected, string(data))
		}
	})

	t.Run("YAML omits nil pointer field", func(t *testing.T) {
		parent := ParentStruct{Name: "test"}
		data, err := yaml.Marshal(parent)
		if err != nil {
			t.Fatalf("Failed to marshal: %v", err)
		}
		expected := "name: test\n"
		if string(data) != expected {
			t.Errorf("Expected '%s', got '%s'", expected, string(data))
		}
	})

	t.Run("YAML includes non-nil pointer field with struct", func(t *testing.T) {
		delegate := &Field[TestDelegate]{}
		delegate.Set(TestDelegate{Filter: []string{"selector1"}})
		parent := ParentStruct{Name: "test", Delegate: delegate}
		data, err := yaml.Marshal(parent)
		if err != nil {
			t.Fatalf("Failed to marshal: %v", err)
		}
		// Verify delegate is present
		var result map[string]interface{}
		err = yaml.Unmarshal(data, &result)
		if err != nil {
			t.Fatalf("Failed to unmarshal result: %v", err)
		}
		if _, ok := result["delegate"]; !ok {
			t.Error("Expected 'delegate' field to be present")
		}
	})
}

// TestAsStructAndAsString tests the accessor methods
func TestAsStructAndAsString(t *testing.T) {
	t.Run("AsStruct returns struct when set", func(t *testing.T) {
		f := Field[TestStruct]{}
		expected := TestStruct{Name: "test", Value: 42}
		f.Set(expected)
		val, ok := f.AsStruct()
		if !ok {
			t.Error("Expected AsStruct() to return true")
		}
		if val != expected {
			t.Errorf("Expected %+v, got %+v", expected, val)
		}
	})

	t.Run("AsStruct returns false for string value", func(t *testing.T) {
		f := Field[TestStruct]{}
		f.SetString("test")
		_, ok := f.AsStruct()
		if ok {
			t.Error("Expected AsStruct() to return false for string value")
		}
	})

	t.Run("AsString returns string when set", func(t *testing.T) {
		f := Field[TestStruct]{}
		expected := "test-string"
		f.SetString(expected)
		val, ok := f.AsString()
		if !ok {
			t.Error("Expected AsString() to return true")
		}
		if val != expected {
			t.Errorf("Expected '%s', got '%s'", expected, val)
		}
	})

	t.Run("AsString returns false for struct value", func(t *testing.T) {
		f := Field[TestStruct]{}
		f.Set(TestStruct{Name: "test", Value: 42})
		_, ok := f.AsString()
		if ok {
			t.Error("Expected AsString() to return false for struct value")
		}
	})

	t.Run("AsString works with expressions", func(t *testing.T) {
		f := Field[TestStruct]{}
		expected := "<+pipeline.variables.myvar>"
		f.SetExpression(expected)
		val, ok := f.AsString()
		if !ok {
			t.Error("Expected AsString() to return true for expression")
		}
		if val != expected {
			t.Errorf("Expected '%s', got '%s'", expected, val)
		}
	})
}

// TestSetMethods tests the setter methods
func TestSetMethods(t *testing.T) {
	t.Run("Set updates field with struct", func(t *testing.T) {
		f := Field[TestStruct]{}
		expected := TestStruct{Name: "test", Value: 42}
		f.Set(expected)
		val, ok := f.AsStruct()
		if !ok {
			t.Error("Expected AsStruct() to return true")
		}
		if val != expected {
			t.Errorf("Expected %+v, got %+v", expected, val)
		}
	})

	t.Run("SetString updates field with string", func(t *testing.T) {
		f := Field[TestStruct]{}
		expected := "test-string"
		f.SetString(expected)
		val, ok := f.AsString()
		if !ok {
			t.Error("Expected AsString() to return true")
		}
		if val != expected {
			t.Errorf("Expected '%s', got '%s'", expected, val)
		}
	})

	t.Run("SetExpression updates field with expression", func(t *testing.T) {
		f := Field[TestStruct]{}
		expected := "<+pipeline.variables.myvar>"
		f.SetExpression(expected)
		if !f.IsExpression() {
			t.Error("Expected IsExpression() to return true")
		}
		val, ok := f.AsString()
		if !ok {
			t.Error("Expected AsString() to return true")
		}
		if val != expected {
			t.Errorf("Expected '%s', got '%s'", expected, val)
		}
	})

	t.Run("Set overwrites previous string value", func(t *testing.T) {
		f := Field[TestStruct]{}
		f.SetString("old-value")
		f.Set(TestStruct{Name: "new", Value: 99})
		if f.IsExpression() {
			t.Error("Expected IsExpression() to return false after Set")
		}
		val, ok := f.AsStruct()
		if !ok {
			t.Error("Expected AsStruct() to return true")
		}
		if val.Name != "new" || val.Value != 99 {
			t.Errorf("Expected {Name:new Value:99}, got %+v", val)
		}
	})

	t.Run("SetString overwrites previous struct value", func(t *testing.T) {
		f := Field[TestStruct]{}
		f.Set(TestStruct{Name: "old", Value: 1})
		f.SetString("new-value")
		val, ok := f.AsString()
		if !ok {
			t.Error("Expected AsString() to return true")
		}
		if val != "new-value" {
			t.Errorf("Expected 'new-value', got '%s'", val)
		}
	})
}