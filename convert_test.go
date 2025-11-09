package main

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

// TestBytesToStruct tests the BytesToStruct function with various input types
// including valid JSON, invalid JSON, and different struct types
func TestBytesToStruct(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	type Address struct {
		Street string `json:"street"`
		City   string `json:"city"`
		Zip    int    `json:"zip"`
	}

	tests := []struct {
		name        string
		jsonBytes   []byte
		target      interface{}
		expected    interface{}
		expectError bool
	}{
		{
			name:      "valid JSON to simple struct",
			jsonBytes: []byte(`{"name":"John","age":30}`),
			target:    &Person{},
			expected: &Person{
				Name: "John",
				Age:  30,
			},
			expectError: false,
		},
		{
			name:      "valid JSON to address struct",
			jsonBytes: []byte(`{"street":"123 Main St","city":"New York","zip":10001}`),
			target:    &Address{},
			expected: &Address{
				Street: "123 Main St",
				City:   "New York",
				Zip:    10001,
			},
			expectError: false,
		},
		{
			name:        "invalid JSON",
			jsonBytes:   []byte(`{"name":"John","age":}`),
			target:      &Person{},
			expected:    &Person{},
			expectError: true,
		},
		{
			name:        "empty JSON bytes",
			jsonBytes:   []byte(``),
			target:      &Person{},
			expected:    &Person{},
			expectError: true,
		},
		{
			name:      "JSON with missing fields",
			jsonBytes: []byte(`{"name":"Jane"}`),
			target:    &Person{},
			expected: &Person{
				Name: "Jane",
				Age:  0, // default value
			},
			expectError: false,
		},
		{
			name:      "JSON with extra fields",
			jsonBytes: []byte(`{"name":"Bob","age":25,"extra":"field"}`),
			target:    &Person{},
			expected: &Person{
				Name: "Bob",
				Age:  25,
			},
			expectError: false,
		},
		{
			name:      "empty JSON object",
			jsonBytes: []byte(`{}`),
			target:    &Person{},
			expected: &Person{
				Name: "",
				Age:  0,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: execute the function
			err := BytesToStruct(tt.jsonBytes, tt.target)

			// Assert: verify error expectation
			if tt.expectError {
				if err == nil {
					t.Errorf("BytesToStruct() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("BytesToStruct() unexpected error: %v", err)
				}
				// Verify the result matches expected
				if !reflect.DeepEqual(tt.target, tt.expected) {
					t.Errorf("BytesToStruct() got %+v; expected %+v", tt.target, tt.expected)
				}
			}
		})
	}
}

// TestMapToStruct tests the MapToStruct function with various map types
// including nested structures, empty maps, and type conversions
func TestMapToStruct(t *testing.T) {
	type User struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Age   int    `json:"age"`
	}

	type Product struct {
		ID    int     `json:"id"`
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}

	tests := []struct {
		name        string
		data        map[string]interface{}
		target      interface{}
		expected    interface{}
		expectError bool
	}{
		{
			name: "simple map to struct",
			data: map[string]interface{}{
				"name":  "Alice",
				"email": "alice@example.com",
				"age":   25,
			},
			target: &User{},
			expected: &User{
				Name:  "Alice",
				Email: "alice@example.com",
				Age:   25,
			},
			expectError: false,
		},
		{
			name: "map with numeric types",
			data: map[string]interface{}{
				"id":    1,
				"name":  "Widget",
				"price": 19.99,
			},
			target: &Product{},
			expected: &Product{
				ID:    1,
				Name:  "Widget",
				Price: 19.99,
			},
			expectError: false,
		},
		{
			name:   "empty map",
			data:   map[string]interface{}{},
			target: &User{},
			expected: &User{
				Name:  "",
				Email: "",
				Age:   0,
			},
			expectError: false,
		},
		{
			name: "map with missing fields",
			data: map[string]interface{}{
				"name": "Bob",
			},
			target: &User{},
			expected: &User{
				Name:  "Bob",
				Email: "",
				Age:   0,
			},
			expectError: false,
		},
		{
			name: "map with extra fields",
			data: map[string]interface{}{
				"name":  "Charlie",
				"email": "charlie@example.com",
				"age":   30,
				"extra": "field",
			},
			target: &User{},
			expected: &User{
				Name:  "Charlie",
				Email: "charlie@example.com",
				Age:   30,
			},
			expectError: false,
		},
		{
			name:   "nil map",
			data:   nil,
			target: &User{},
			expected: &User{
				Name:  "",
				Email: "",
				Age:   0,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: execute the function
			err := MapToStruct(tt.data, tt.target)

			// Assert: verify error expectation
			if tt.expectError {
				if err == nil {
					t.Errorf("MapToStruct() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("MapToStruct() unexpected error: %v", err)
				}
				// Verify the result matches expected
				if !reflect.DeepEqual(tt.target, tt.expected) {
					t.Errorf("MapToStruct() got %+v; expected %+v", tt.target, tt.expected)
				}
			}
		})
	}
}

// TestMapToStructWithComplexTypes tests MapToStruct with nested and complex data types
func TestMapToStructWithComplexTypes(t *testing.T) {
	type Address struct {
		Street string `json:"street"`
		City   string `json:"city"`
	}

	type Person struct {
		Name    string  `json:"name"`
		Address Address `json:"address"`
	}

	tests := []struct {
		name        string
		data        map[string]interface{}
		target      interface{}
		expected    interface{}
		expectError bool
	}{
		{
			name: "nested struct",
			data: map[string]interface{}{
				"name": "John",
				"address": map[string]interface{}{
					"street": "123 Main St",
					"city":   "Boston",
				},
			},
			target: &Person{},
			expected: &Person{
				Name: "John",
				Address: Address{
					Street: "123 Main St",
					City:   "Boston",
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: execute the function
			err := MapToStruct(tt.data, tt.target)

			// Assert: verify error expectation
			if tt.expectError {
				if err == nil {
					t.Errorf("MapToStruct() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("MapToStruct() unexpected error: %v", err)
				}
				// Verify the result matches expected
				if !reflect.DeepEqual(tt.target, tt.expected) {
					t.Errorf("MapToStruct() got %+v; expected %+v", tt.target, tt.expected)
				}
			}
		})
	}
}

// TestStructToJSONString tests the StructToJSONString function with various struct types
// including simple structs, nested structs, and edge cases
func TestStructToJSONString(t *testing.T) {
	type SimpleStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	type EmptyStruct struct{}

	type StructWithSlice struct {
		Items []string `json:"items"`
	}

	tests := []struct {
		name        string
		structData  interface{}
		expected    string
		expectError bool
	}{
		{
			name: "simple struct",
			structData: SimpleStruct{
				Name: "John",
				Age:  30,
			},
			expected:    `{"name":"John","age":30}`,
			expectError: false,
		},
		{
			name: "struct with zero values",
			structData: SimpleStruct{
				Name: "",
				Age:  0,
			},
			expected:    `{"name":"","age":0}`,
			expectError: false,
		},
		{
			name:        "empty struct",
			structData:  EmptyStruct{},
			expected:    `{}`,
			expectError: false,
		},
		{
			name: "struct with slice",
			structData: StructWithSlice{
				Items: []string{"apple", "banana", "cherry"},
			},
			expected:    `{"items":["apple","banana","cherry"]}`,
			expectError: false,
		},
		{
			name: "struct with empty slice",
			structData: StructWithSlice{
				Items: []string{},
			},
			expected:    `{"items":[]}`,
			expectError: false,
		},
		{
			name: "struct with nil slice",
			structData: StructWithSlice{
				Items: nil,
			},
			expected:    `{"items":null}`,
			expectError: false,
		},
		{
			name:        "pointer to struct",
			structData:  &SimpleStruct{Name: "Alice", Age: 25},
			expected:    `{"name":"Alice","age":25}`,
			expectError: false,
		},
		{
			name:        "nil pointer",
			structData:  (*SimpleStruct)(nil),
			expected:    `null`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: execute the function
			result, err := StructToJSONString(tt.structData)

			// Assert: verify error expectation
			if tt.expectError {
				if err == nil {
					t.Errorf("StructToJSONString() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("StructToJSONString() unexpected error: %v", err)
				}
				// Verify the result matches expected
				if result != tt.expected {
					t.Errorf("StructToJSONString() got %q; expected %q", result, tt.expected)
				}
			}
		})
	}
}

// TestStructToJSONStringWithComplexTypes tests StructToJSONString with nested and complex types
func TestStructToJSONStringWithComplexTypes(t *testing.T) {
	type Address struct {
		Street string `json:"street"`
		City   string `json:"city"`
	}

	type Person struct {
		Name    string  `json:"name"`
		Age     int     `json:"age"`
		Address Address `json:"address"`
	}

	tests := []struct {
		name        string
		structData  interface{}
		expected    string
		expectError bool
	}{
		{
			name: "nested struct",
			structData: Person{
				Name: "John",
				Age:  30,
				Address: Address{
					Street: "123 Main St",
					City:   "New York",
				},
			},
			expected:    `{"name":"John","age":30,"address":{"street":"123 Main St","city":"New York"}}`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: execute the function
			result, err := StructToJSONString(tt.structData)

			// Assert: verify error expectation
			if tt.expectError {
				if err == nil {
					t.Errorf("StructToJSONString() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("StructToJSONString() unexpected error: %v", err)
				}
				// Verify the result matches expected
				if result != tt.expected {
					t.Errorf("StructToJSONString() got %q; expected %q", result, tt.expected)
				}
			}
		})
	}
}

// TestStructToJSONStringWithUnmarshallableTypes tests error cases for StructToJSONString
func TestStructToJSONStringWithUnmarshallableTypes(t *testing.T) {
	tests := []struct {
		name        string
		structData  interface{}
		expectError bool
	}{
		{
			name:        "channel type should fail",
			structData:  make(chan int),
			expectError: true,
		},
		{
			name:        "function type should fail",
			structData:  func() {},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: execute the function
			result, err := StructToJSONString(tt.structData)

			// Assert: verify error expectation
			if tt.expectError {
				if err == nil {
					t.Errorf("StructToJSONString() expected error but got nil")
				}
				if result != "" {
					t.Errorf("StructToJSONString() expected empty string on error, got %q", result)
				}
			}
		})
	}
}

// TestInterfaceToStruct tests the InterfaceToStruct function with various interface types
// including maps, structs, and primitive types
func TestInterfaceToStruct(t *testing.T) {
	type Source struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	type Destination struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	type PartialDest struct {
		Name string `json:"name"`
	}

	tests := []struct {
		name        string
		source      interface{}
		target      interface{}
		expected    interface{}
		expectError bool
	}{
		{
			name: "struct to struct",
			source: Source{
				Name:  "test",
				Value: 42,
			},
			target: &Destination{},
			expected: &Destination{
				Name:  "test",
				Value: 42,
			},
			expectError: false,
		},
		{
			name: "map to struct",
			source: map[string]interface{}{
				"name":  "from map",
				"value": 100,
			},
			target: &Destination{},
			expected: &Destination{
				Name:  "from map",
				Value: 100,
			},
			expectError: false,
		},
		{
			name: "struct to partial struct",
			source: Source{
				Name:  "partial",
				Value: 99,
			},
			target: &PartialDest{},
			expected: &PartialDest{
				Name: "partial",
			},
			expectError: false,
		},
		{
			name:   "empty struct",
			source: Source{},
			target: &Destination{},
			expected: &Destination{
				Name:  "",
				Value: 0,
			},
			expectError: false,
		},
		{
			name: "pointer to struct",
			source: &Source{
				Name:  "pointer source",
				Value: 123,
			},
			target: &Destination{},
			expected: &Destination{
				Name:  "pointer source",
				Value: 123,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: execute the function
			err := InterfaceToStruct(tt.source, tt.target)

			// Assert: verify error expectation
			if tt.expectError {
				if err == nil {
					t.Errorf("InterfaceToStruct() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("InterfaceToStruct() unexpected error: %v", err)
				}
				// Verify the result matches expected
				if !reflect.DeepEqual(tt.target, tt.expected) {
					t.Errorf("InterfaceToStruct() got %+v; expected %+v", tt.target, tt.expected)
				}
			}
		})
	}
}

// TestInterfaceToStructWithComplexTypes tests InterfaceToStruct with nested structures
func TestInterfaceToStructWithComplexTypes(t *testing.T) {
	type Address struct {
		Street string `json:"street"`
		City   string `json:"city"`
	}

	type SourcePerson struct {
		Name    string  `json:"name"`
		Age     int     `json:"age"`
		Address Address `json:"address"`
	}

	type DestPerson struct {
		Name    string  `json:"name"`
		Age     int     `json:"age"`
		Address Address `json:"address"`
	}

	tests := []struct {
		name        string
		source      interface{}
		target      interface{}
		expected    interface{}
		expectError bool
	}{
		{
			name: "nested struct conversion",
			source: SourcePerson{
				Name: "John Doe",
				Age:  35,
				Address: Address{
					Street: "456 Elm St",
					City:   "Los Angeles",
				},
			},
			target: &DestPerson{},
			expected: &DestPerson{
				Name: "John Doe",
				Age:  35,
				Address: Address{
					Street: "456 Elm St",
					City:   "Los Angeles",
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: execute the function
			err := InterfaceToStruct(tt.source, tt.target)

			// Assert: verify error expectation
			if tt.expectError {
				if err == nil {
					t.Errorf("InterfaceToStruct() expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("InterfaceToStruct() unexpected error: %v", err)
				}
				// Verify the result matches expected
				if !reflect.DeepEqual(tt.target, tt.expected) {
					t.Errorf("InterfaceToStruct() got %+v; expected %+v", tt.target, tt.expected)
				}
			}
		})
	}
}

// TestInterfaceToStructErrorCases tests error scenarios for InterfaceToStruct
func TestInterfaceToStructErrorCases(t *testing.T) {
	type Destination struct {
		Name string `json:"name"`
	}

	tests := []struct {
		name        string
		source      interface{}
		target      interface{}
		expectError bool
	}{
		{
			name:        "channel as source",
			source:      make(chan int),
			target:      &Destination{},
			expectError: true,
		},
		{
			name:        "function as source",
			source:      func() {},
			target:      &Destination{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: execute the function
			err := InterfaceToStruct(tt.source, tt.target)

			// Assert: verify error expectation
			if tt.expectError && err == nil {
				t.Errorf("InterfaceToStruct() expected error but got nil")
			}
		})
	}
}

// TestAllFunctionsIntegration tests all conversion functions working together
func TestAllFunctionsIntegration(t *testing.T) {
	type TestStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	t.Run("full conversion cycle", func(t *testing.T) {
		// Arrange: create original struct
		original := TestStruct{
			Name:  "Integration Test",
			Value: 999,
		}

		// Act 1: struct to JSON string
		jsonString, err := StructToJSONString(original)
		if err != nil {
			t.Fatalf("StructToJSONString() failed: %v", err)
		}

		// Act 2: JSON string to bytes to struct
		var result1 TestStruct
		err = BytesToStruct([]byte(jsonString), &result1)
		if err != nil {
			t.Fatalf("BytesToStruct() failed: %v", err)
		}

		// Assert: verify first cycle
		if !reflect.DeepEqual(original, result1) {
			t.Errorf("First conversion cycle failed: got %+v; expected %+v", result1, original)
		}

		// Act 3: struct to interface to struct
		var result2 TestStruct
		err = InterfaceToStruct(original, &result2)
		if err != nil {
			t.Fatalf("InterfaceToStruct() failed: %v", err)
		}

		// Assert: verify second cycle
		if !reflect.DeepEqual(original, result2) {
			t.Errorf("Interface conversion failed: got %+v; expected %+v", result2, original)
		}

		// Act 4: map to struct
		dataMap := map[string]interface{}{
			"name":  "Integration Test",
			"value": 999,
		}
		var result3 TestStruct
		err = MapToStruct(dataMap, &result3)
		if err != nil {
			t.Fatalf("MapToStruct() failed: %v", err)
		}

		// Assert: verify map conversion
		if !reflect.DeepEqual(original, result3) {
			t.Errorf("Map conversion failed: got %+v; expected %+v", result3, original)
		}
	})
}

// TestErrorHandlingConsistency tests that all functions handle errors consistently
func TestErrorHandlingConsistency(t *testing.T) {
	t.Run("all functions return proper errors", func(t *testing.T) {
		type Target struct {
			Name string `json:"name"`
		}

		// Test BytesToStruct with invalid JSON
		err1 := BytesToStruct([]byte(`invalid json`), &Target{})
		if err1 == nil {
			t.Error("BytesToStruct should return error for invalid JSON")
		}

		// Test MapToStruct with unmarshalable value
		// Create a map with a channel (which can't be marshaled)
		badMap := map[string]interface{}{
			"channel": make(chan int),
		}
		err2 := MapToStruct(badMap, &Target{})
		if err2 == nil {
			t.Error("MapToStruct should return error for unmarshalable data")
		}

		// Test StructToJSONString with unmarshalable type
		_, err3 := StructToJSONString(make(chan int))
		if err3 == nil {
			t.Error("StructToJSONString should return error for unmarshalable type")
		}

		// Test InterfaceToStruct with unmarshalable source
		err4 := InterfaceToStruct(make(chan int), &Target{})
		if err4 == nil {
			t.Error("InterfaceToStruct should return error for unmarshalable source")
		}
	})
}

// TestNilPointerHandling tests how functions handle nil pointers
func TestNilPointerHandling(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name"`
	}

	t.Run("BytesToStruct with valid pointer", func(t *testing.T) {
		var target TestStruct
		err := BytesToStruct([]byte(`{"name":"test"}`), &target)
		if err != nil {
			t.Errorf("BytesToStruct() with valid pointer failed: %v", err)
		}
		if target.Name != "test" {
			t.Errorf("Expected name 'test', got '%s'", target.Name)
		}
	})

	t.Run("StructToJSONString with nil pointer", func(t *testing.T) {
		var nilPtr *TestStruct
		result, err := StructToJSONString(nilPtr)
		if err != nil {
			t.Errorf("StructToJSONString() with nil pointer failed: %v", err)
		}
		if result != "null" {
			t.Errorf("Expected 'null', got '%s'", result)
		}
	})
}

// TestJSONTagHandling tests that JSON tags are properly respected
func TestJSONTagHandling(t *testing.T) {
	type CustomTagStruct struct {
		PublicName  string `json:"public_name"`
		InternalAge int    `json:"age"`
		IgnoredField string `json:"-"`
	}

	t.Run("JSON tags are respected in conversion", func(t *testing.T) {
		// Arrange
		original := CustomTagStruct{
			PublicName:   "John",
			InternalAge:  30,
			IgnoredField: "should not appear",
		}

		// Act: convert to JSON string
		jsonString, err := StructToJSONString(original)
		if err != nil {
			t.Fatalf("StructToJSONString() failed: %v", err)
		}

		// Assert: verify ignored field is not in JSON
		var jsonMap map[string]interface{}
		if err := json.Unmarshal([]byte(jsonString), &jsonMap); err != nil {
			t.Fatalf("Failed to parse JSON: %v", err)
		}

		if _, exists := jsonMap["IgnoredField"]; exists {
			t.Error("Ignored field should not appear in JSON")
		}

		if jsonMap["public_name"] != "John" {
			t.Errorf("Expected public_name to be 'John', got %v", jsonMap["public_name"])
		}
	})
}

// BenchmarkBytesToStruct benchmarks the BytesToStruct function
func BenchmarkBytesToStruct(b *testing.B) {
	type TestStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	jsonBytes := []byte(`{"name":"benchmark","value":42}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result TestStruct
		_ = BytesToStruct(jsonBytes, &result)
	}
}

// BenchmarkStructToJSONString benchmarks the StructToJSONString function
func BenchmarkStructToJSONString(b *testing.B) {
	type TestStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	testData := TestStruct{
		Name:  "benchmark",
		Value: 42,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = StructToJSONString(testData)
	}
}

// Helper function to create test error
func createTestError(msg string) error {
	return errors.New(msg)
}

