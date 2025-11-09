package main

import (
	"errors"
	"testing"
)

// TestTernary tests the Ternary function with various types and conditions
// including true and false conditions with different data types
func TestTernary(t *testing.T) {
	tests := []struct {
		name      string
		condition bool
		trueValue interface{}
		elseValue interface{}
		expected  interface{}
	}{
		{
			name:      "condition true with integers",
			condition: true,
			trueValue: 10,
			elseValue: 20,
			expected:  10,
		},
		{
			name:      "condition false with integers",
			condition: false,
			trueValue: 10,
			elseValue: 20,
			expected:  20,
		},
		{
			name:      "condition true with strings",
			condition: true,
			trueValue: "hello",
			elseValue: "world",
			expected:  "hello",
		},
		{
			name:      "condition false with strings",
			condition: false,
			trueValue: "hello",
			elseValue: "world",
			expected:  "world",
		},
		{
			name:      "condition true with zero values",
			condition: true,
			trueValue: 0,
			elseValue: 1,
			expected:  0,
		},
		{
			name:      "condition false with empty string",
			condition: false,
			trueValue: "non-empty",
			elseValue: "",
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result interface{}

			// Act: execute the function based on the type
			switch v := tt.trueValue.(type) {
			case int:
				result = Ternary(tt.condition, v, tt.elseValue.(int))
			case string:
				result = Ternary(tt.condition, v, tt.elseValue.(string))
			}

			// Assert: verify the result
			if result != tt.expected {
				t.Errorf("Ternary(%v, %v, %v) = %v; expected %v",
					tt.condition, tt.trueValue, tt.elseValue, result, tt.expected)
			}
		})
	}
}

// TestTernaryF tests the TernaryF function with function evaluation
// ensuring lazy evaluation and correct function selection based on condition
func TestTernaryF(t *testing.T) {
	tests := []struct {
		name       string
		condition  bool
		trueFunc   func() int
		elseFunc   func() int
		expected   int
		trueCalled bool
		elseCalled bool
	}{
		{
			name:      "condition true calls true function",
			condition: true,
			trueFunc: func() int {
				return 100
			},
			elseFunc: func() int {
				return 200
			},
			expected:   100,
			trueCalled: true,
		},
		{
			name:      "condition false calls else function",
			condition: false,
			trueFunc: func() int {
				return 100
			},
			elseFunc: func() int {
				return 200
			},
			expected:   200,
			elseCalled: true,
		},
		{
			name:      "condition true with zero return value",
			condition: true,
			trueFunc: func() int {
				return 0
			},
			elseFunc: func() int {
				return 999
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: execute the function
			result := TernaryF(tt.condition, tt.trueFunc, tt.elseFunc)

			// Assert: verify the result
			if result != tt.expected {
				t.Errorf("TernaryF(%v, trueFunc, elseFunc) = %v; expected %v",
					tt.condition, result, tt.expected)
			}
		})
	}
}

// TestTernaryFLazyEvaluation tests that TernaryF only evaluates the selected function
func TestTernaryFLazyEvaluation(t *testing.T) {
	tests := []struct {
		name      string
		condition bool
		callCount *int
	}{
		{
			name:      "condition true only evaluates true function",
			condition: true,
			callCount: new(int),
		},
		{
			name:      "condition false only evaluates else function",
			condition: false,
			callCount: new(int),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: create functions that track if they were called
			trueCalled := false
			elseCalled := false

			trueFunc := func() int {
				trueCalled = true
				return 1
			}
			elseFunc := func() int {
				elseCalled = true
				return 2
			}

			// Act: execute the function
			_ = TernaryF(tt.condition, trueFunc, elseFunc)

			// Assert: verify only one function was called
			if tt.condition && !trueCalled {
				t.Error("Expected true function to be called when condition is true")
			}
			if tt.condition && elseCalled {
				t.Error("Expected else function NOT to be called when condition is true")
			}
			if !tt.condition && trueCalled {
				t.Error("Expected true function NOT to be called when condition is false")
			}
			if !tt.condition && !elseCalled {
				t.Error("Expected else function to be called when condition is false")
			}
		})
	}
}

// TestLogicalOrInt tests the LogicalOrInt function with various signed integer types
// including zero values, positive numbers, and negative numbers
func TestLogicalOrInt(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		elValue  int
		expected int
	}{
		{
			name:     "zero value returns elValue",
			value:    0,
			elValue:  42,
			expected: 42,
		},
		{
			name:     "positive value returns value",
			value:    10,
			elValue:  42,
			expected: 10,
		},
		{
			name:     "negative value returns value",
			value:    -5,
			elValue:  42,
			expected: -5,
		},
		{
			name:     "zero value with zero elValue",
			value:    0,
			elValue:  0,
			expected: 0,
		},
		{
			name:     "negative value with negative elValue",
			value:    -10,
			elValue:  -20,
			expected: -10,
		},
		{
			name:     "large positive value",
			value:    999999,
			elValue:  1,
			expected: 999999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: execute the function
			result := LogicalOrInt(tt.value, tt.elValue)

			// Assert: verify the result
			if result != tt.expected {
				t.Errorf("LogicalOrInt(%d, %d) = %d; expected %d",
					tt.value, tt.elValue, result, tt.expected)
			}
		})
	}
}

// TestLogicalOrIntDifferentTypes tests LogicalOrInt with different signed integer types
func TestLogicalOrIntDifferentTypes(t *testing.T) {
	t.Run("int8 type", func(t *testing.T) {
		var value int8 = 0
		var elValue int8 = 10
		result := LogicalOrInt(value, elValue)
		if result != elValue {
			t.Errorf("LogicalOrInt[int8](0, 10) = %d; expected %d", result, elValue)
		}
	})

	t.Run("int16 type", func(t *testing.T) {
		var value int16 = 100
		var elValue int16 = 10
		result := LogicalOrInt(value, elValue)
		if result != value {
			t.Errorf("LogicalOrInt[int16](100, 10) = %d; expected %d", result, value)
		}
	})

	t.Run("int32 type", func(t *testing.T) {
		var value int32 = 0
		var elValue int32 = 200
		result := LogicalOrInt(value, elValue)
		if result != elValue {
			t.Errorf("LogicalOrInt[int32](0, 200) = %d; expected %d", result, elValue)
		}
	})

	t.Run("int64 type", func(t *testing.T) {
		var value int64 = -50
		var elValue int64 = 100
		result := LogicalOrInt(value, elValue)
		if result != value {
			t.Errorf("LogicalOrInt[int64](-50, 100) = %d; expected %d", result, value)
		}
	})
}

// TestLogicalOrFloat tests the LogicalOrFloat function with various float types
// including zero values, positive numbers, negative numbers, and special values
func TestLogicalOrFloat(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		elValue  float64
		expected float64
	}{
		{
			name:     "zero value returns elValue",
			value:    0.0,
			elValue:  3.14,
			expected: 3.14,
		},
		{
			name:     "positive value returns value",
			value:    2.5,
			elValue:  3.14,
			expected: 2.5,
		},
		{
			name:     "negative value returns value",
			value:    -1.5,
			elValue:  3.14,
			expected: -1.5,
		},
		{
			name:     "zero value with zero elValue",
			value:    0.0,
			elValue:  0.0,
			expected: 0.0,
		},
		{
			name:     "small positive value",
			value:    0.001,
			elValue:  100.0,
			expected: 0.001,
		},
		{
			name:     "small negative value",
			value:    -0.001,
			elValue:  100.0,
			expected: -0.001,
		},
		{
			name:     "large value",
			value:    999999.999,
			elValue:  1.0,
			expected: 999999.999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: execute the function
			result := LogicalOrFloat(tt.value, tt.elValue)

			// Assert: verify the result
			if result != tt.expected {
				t.Errorf("LogicalOrFloat(%f, %f) = %f; expected %f",
					tt.value, tt.elValue, result, tt.expected)
			}
		})
	}
}

// TestLogicalOrFloatDifferentTypes tests LogicalOrFloat with different float types
func TestLogicalOrFloatDifferentTypes(t *testing.T) {
	t.Run("float32 type", func(t *testing.T) {
		var value float32 = 0.0
		var elValue float32 = 1.23
		result := LogicalOrFloat(value, elValue)
		if result != elValue {
			t.Errorf("LogicalOrFloat[float32](0.0, 1.23) = %f; expected %f", result, elValue)
		}
	})

	t.Run("float32 non-zero", func(t *testing.T) {
		var value float32 = 5.67
		var elValue float32 = 1.23
		result := LogicalOrFloat(value, elValue)
		if result != value {
			t.Errorf("LogicalOrFloat[float32](5.67, 1.23) = %f; expected %f", result, value)
		}
	})

	t.Run("float64 type", func(t *testing.T) {
		var value float64 = -2.5
		var elValue float64 = 10.0
		result := LogicalOrFloat(value, elValue)
		if result != value {
			t.Errorf("LogicalOrFloat[float64](-2.5, 10.0) = %f; expected %f", result, value)
		}
	})
}

// TestLogicalOrString tests the LogicalOrString function with various string values
// including empty strings, non-empty strings, and special characters
func TestLogicalOrString(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		elValue  string
		expected string
	}{
		{
			name:     "empty string returns elValue",
			value:    "",
			elValue:  "default",
			expected: "default",
		},
		{
			name:     "non-empty string returns value",
			value:    "hello",
			elValue:  "default",
			expected: "hello",
		},
		{
			name:     "empty string with empty elValue",
			value:    "",
			elValue:  "",
			expected: "",
		},
		{
			name:     "single character string",
			value:    "a",
			elValue:  "default",
			expected: "a",
		},
		{
			name:     "whitespace string returns value",
			value:    " ",
			elValue:  "default",
			expected: " ",
		},
		{
			name:     "string with special characters",
			value:    "hello\nworld",
			elValue:  "default",
			expected: "hello\nworld",
		},
		{
			name:     "long string",
			value:    "this is a very long string with many characters",
			elValue:  "default",
			expected: "this is a very long string with many characters",
		},
		{
			name:     "unicode string",
			value:    "你好世界",
			elValue:  "default",
			expected: "你好世界",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: execute the function
			result := LogicalOrString(tt.value, tt.elValue)

			// Assert: verify the result
			if result != tt.expected {
				t.Errorf("LogicalOrString(%q, %q) = %q; expected %q",
					tt.value, tt.elValue, result, tt.expected)
			}
		})
	}
}

// TestLogicalOrStringCustomType tests LogicalOrString with custom string types
func TestLogicalOrStringCustomType(t *testing.T) {
	type CustomString string

	tests := []struct {
		name     string
		value    CustomString
		elValue  CustomString
		expected CustomString
	}{
		{
			name:     "empty custom string",
			value:    "",
			elValue:  "fallback",
			expected: "fallback",
		},
		{
			name:     "non-empty custom string",
			value:    "custom",
			elValue:  "fallback",
			expected: "custom",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: execute the function
			result := LogicalOrString(tt.value, tt.elValue)

			// Assert: verify the result
			if result != tt.expected {
				t.Errorf("LogicalOrString[CustomString](%q, %q) = %q; expected %q",
					tt.value, tt.elValue, result, tt.expected)
			}
		})
	}
}

// TestLogicalError tests the LogicalError function with various error conditions
// including nil errors, non-nil errors, and conditional errors
func TestLogicalError(t *testing.T) {
	errFirst := errors.New("first error")
	errCondition := errors.New("condition error")

	tests := []struct {
		name         string
		err          error
		condition    bool
		errCondition error
		expected     error
	}{
		{
			name:         "err is not nil returns err",
			err:          errFirst,
			condition:    false,
			errCondition: errCondition,
			expected:     errFirst,
		},
		{
			name:         "err is not nil returns err even when condition is true",
			err:          errFirst,
			condition:    true,
			errCondition: errCondition,
			expected:     errFirst,
		},
		{
			name:         "err is nil and condition false returns nil",
			err:          nil,
			condition:    false,
			errCondition: errCondition,
			expected:     nil,
		},
		{
			name:         "err is nil and condition true returns errCondition",
			err:          nil,
			condition:    true,
			errCondition: errCondition,
			expected:     errCondition,
		},
		{
			name:         "all nil returns nil",
			err:          nil,
			condition:    false,
			errCondition: nil,
			expected:     nil,
		},
		{
			name:         "err nil, condition true, errCondition nil returns nil",
			err:          nil,
			condition:    true,
			errCondition: nil,
			expected:     nil,
		},
		{
			name:         "err nil, condition false, errCondition not nil returns nil",
			err:          nil,
			condition:    false,
			errCondition: errCondition,
			expected:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: execute the function
			result := LogicalError(tt.err, tt.condition, tt.errCondition)

			// Assert: verify the result
			if result != tt.expected {
				t.Errorf("LogicalError(%v, %v, %v) = %v; expected %v",
					tt.err, tt.condition, tt.errCondition, result, tt.expected)
			}
		})
	}
}

// TestLogicalErrorPriority tests that LogicalError prioritizes the first error
func TestLogicalErrorPriority(t *testing.T) {
	t.Run("first error has priority over condition error", func(t *testing.T) {
		// Arrange: create distinct errors
		firstErr := errors.New("first error should be returned")
		conditionErr := errors.New("this should not be returned")

		// Act: execute the function with both errors present
		result := LogicalError(firstErr, true, conditionErr)

		// Assert: verify first error is returned
		if result != firstErr {
			t.Errorf("Expected first error to be returned, got %v", result)
		}
		if result == conditionErr {
			t.Error("Condition error should not be returned when first error exists")
		}
	})
}

// TestLogicalErrorMessages tests that error messages are preserved
func TestLogicalErrorMessages(t *testing.T) {
	tests := []struct {
		name            string
		err             error
		condition       bool
		errCondition    error
		expectedMessage string
		expectNil       bool
	}{
		{
			name:            "first error message preserved",
			err:             errors.New("first error message"),
			condition:       true,
			errCondition:    errors.New("condition error"),
			expectedMessage: "first error message",
			expectNil:       false,
		},
		{
			name:            "condition error message preserved",
			err:             nil,
			condition:       true,
			errCondition:    errors.New("condition error message"),
			expectedMessage: "condition error message",
			expectNil:       false,
		},
		{
			name:         "no error returns nil",
			err:          nil,
			condition:    false,
			errCondition: errors.New("should not be returned"),
			expectNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: execute the function
			result := LogicalError(tt.err, tt.condition, tt.errCondition)

			// Assert: verify the result
			if tt.expectNil {
				if result != nil {
					t.Errorf("Expected nil error, got %v", result)
				}
			} else {
				if result == nil {
					t.Error("Expected non-nil error, got nil")
				} else if result.Error() != tt.expectedMessage {
					t.Errorf("Expected error message %q, got %q", tt.expectedMessage, result.Error())
				}
			}
		})
	}
}
