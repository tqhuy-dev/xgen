package utilities

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"testing"
	"time"
)

// TestToBoolE tests the ToBoolE function with various input types
// including bool, integers, floats, strings, nil, and json.Number
func TestToBoolE(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    bool
		expectError bool
	}{
		{
			name:        "bool true",
			input:       true,
			expected:    true,
			expectError: false,
		},
		{
			name:        "bool false",
			input:       false,
			expected:    false,
			expectError: false,
		},
		{
			name:        "nil value",
			input:       nil,
			expected:    false,
			expectError: false,
		},
		{
			name:        "int non-zero",
			input:       1,
			expected:    true,
			expectError: false,
		},
		{
			name:        "int zero",
			input:       0,
			expected:    false,
			expectError: false,
		},
		{
			name:        "int64 non-zero",
			input:       int64(100),
			expected:    true,
			expectError: false,
		},
		{
			name:        "int32 zero",
			input:       int32(0),
			expected:    false,
			expectError: false,
		},
		{
			name:        "uint non-zero",
			input:       uint(5),
			expected:    true,
			expectError: false,
		},
		{
			name:        "float64 non-zero",
			input:       float64(1.5),
			expected:    true,
			expectError: false,
		},
		{
			name:        "float32 zero",
			input:       float32(0.0),
			expected:    false,
			expectError: false,
		},
		{
			name:        "string true",
			input:       "true",
			expected:    true,
			expectError: false,
		},
		{
			name:        "string false",
			input:       "false",
			expected:    false,
			expectError: false,
		},
		{
			name:        "string 1",
			input:       "1",
			expected:    true,
			expectError: false,
		},
		{
			name:        "string 0",
			input:       "0",
			expected:    false,
			expectError: false,
		},
		{
			name:        "invalid string",
			input:       "invalid",
			expected:    false,
			expectError: true,
		},
		{
			name:        "json.Number valid",
			input:       json.Number("1"),
			expected:    true,
			expectError: false,
		},
		{
			name:        "json.Number zero",
			input:       json.Number("0"),
			expected:    false,
			expectError: false,
		},
		{
			name:        "time.Duration non-zero",
			input:       time.Second,
			expected:    true,
			expectError: false,
		},
		{
			name:        "unsupported type",
			input:       []int{1, 2, 3},
			expected:    false,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := ToBoolE(tt.input)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("ToBoolE(%v) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ToBoolE(%v) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("ToBoolE(%v) = %v; expected %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}

// TestToFloat64E tests the ToFloat64E function with various numeric types,
// strings, booleans, and edge cases
func TestToFloat64E(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    float64
		expectError bool
	}{
		{
			name:        "float64 positive",
			input:       float64(123.456),
			expected:    123.456,
			expectError: false,
		},
		{
			name:        "float32",
			input:       float32(45.67),
			expected:    float64(float32(45.67)),
			expectError: false,
		},
		{
			name:        "int",
			input:       42,
			expected:    42.0,
			expectError: false,
		},
		{
			name:        "int64",
			input:       int64(100),
			expected:    100.0,
			expectError: false,
		},
		{
			name:        "uint",
			input:       uint(50),
			expected:    50.0,
			expectError: false,
		},
		{
			name:        "string valid",
			input:       "123.45",
			expected:    123.45,
			expectError: false,
		},
		{
			name:        "string invalid",
			input:       "not a number",
			expected:    0,
			expectError: true,
		},
		{
			name:        "bool true",
			input:       true,
			expected:    1.0,
			expectError: false,
		},
		{
			name:        "bool false",
			input:       false,
			expected:    0.0,
			expectError: false,
		},
		{
			name:        "nil",
			input:       nil,
			expected:    0.0,
			expectError: false,
		},
		{
			name:        "negative float",
			input:       -99.99,
			expected:    -99.99,
			expectError: false,
		},
		{
			name:        "unsupported type",
			input:       []int{1, 2},
			expected:    0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := ToFloat64E(tt.input)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("ToFloat64E(%v) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ToFloat64E(%v) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("ToFloat64E(%v) = %v; expected %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}

// TestToFloat32E tests the ToFloat32E function with various numeric types
func TestToFloat32E(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    float32
		expectError bool
	}{
		{
			name:        "float32 positive",
			input:       float32(123.456),
			expected:    123.456,
			expectError: false,
		},
		{
			name:        "float64",
			input:       float64(45.67),
			expected:    float32(45.67),
			expectError: false,
		},
		{
			name:        "int",
			input:       42,
			expected:    42.0,
			expectError: false,
		},
		{
			name:        "string valid",
			input:       "123.45",
			expected:    123.45,
			expectError: false,
		},
		{
			name:        "string invalid",
			input:       "invalid",
			expected:    0,
			expectError: true,
		},
		{
			name:        "bool true",
			input:       true,
			expected:    1.0,
			expectError: false,
		},
		{
			name:        "bool false",
			input:       false,
			expected:    0.0,
			expectError: false,
		},
		{
			name:        "nil",
			input:       nil,
			expected:    0.0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := ToFloat32E(tt.input)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("ToFloat32E(%v) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ToFloat32E(%v) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("ToFloat32E(%v) = %v; expected %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}

// TestToInt64E tests the ToInt64E function with various integer types,
// floats, strings, and edge cases
func TestToInt64E(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    int64
		expectError bool
	}{
		{
			name:        "int64 positive",
			input:       int64(12345),
			expected:    12345,
			expectError: false,
		},
		{
			name:        "int",
			input:       42,
			expected:    42,
			expectError: false,
		},
		{
			name:        "int32",
			input:       int32(100),
			expected:    100,
			expectError: false,
		},
		{
			name:        "int16",
			input:       int16(50),
			expected:    50,
			expectError: false,
		},
		{
			name:        "int8",
			input:       int8(10),
			expected:    10,
			expectError: false,
		},
		{
			name:        "uint",
			input:       uint(25),
			expected:    25,
			expectError: false,
		},
		{
			name:        "uint64",
			input:       uint64(999),
			expected:    999,
			expectError: false,
		},
		{
			name:        "float64",
			input:       float64(123.99),
			expected:    123,
			expectError: false,
		},
		{
			name:        "string valid",
			input:       "12345",
			expected:    12345,
			expectError: false,
		},
		{
			name:        "string with decimal zeros",
			input:       "100.00",
			expected:    100,
			expectError: false,
		},
		{
			name:        "string invalid",
			input:       "not a number",
			expected:    0,
			expectError: true,
		},
		{
			name:        "json.Number",
			input:       json.Number("456"),
			expected:    456,
			expectError: false,
		},
		{
			name:        "bool true",
			input:       true,
			expected:    1,
			expectError: false,
		},
		{
			name:        "bool false",
			input:       false,
			expected:    0,
			expectError: false,
		},
		{
			name:        "nil",
			input:       nil,
			expected:    0,
			expectError: false,
		},
		{
			name:        "negative int",
			input:       -100,
			expected:    -100,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := ToInt64E(tt.input)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("ToInt64E(%v) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ToInt64E(%v) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("ToInt64E(%v) = %v; expected %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}

// TestToInt32E tests the ToInt32E function
func TestToInt32E(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    int32
		expectError bool
	}{
		{
			name:        "int32 positive",
			input:       int32(12345),
			expected:    12345,
			expectError: false,
		},
		{
			name:        "int",
			input:       42,
			expected:    42,
			expectError: false,
		},
		{
			name:        "string valid",
			input:       "999",
			expected:    999,
			expectError: false,
		},
		{
			name:        "float64",
			input:       float64(100.5),
			expected:    100,
			expectError: false,
		},
		{
			name:        "bool true",
			input:       true,
			expected:    1,
			expectError: false,
		},
		{
			name:        "nil",
			input:       nil,
			expected:    0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := ToInt32E(tt.input)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("ToInt32E(%v) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ToInt32E(%v) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("ToInt32E(%v) = %v; expected %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}

// TestToInt16E tests the ToInt16E function
func TestToInt16E(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    int16
		expectError bool
	}{
		{
			name:        "int16 positive",
			input:       int16(1234),
			expected:    1234,
			expectError: false,
		},
		{
			name:        "int",
			input:       42,
			expected:    42,
			expectError: false,
		},
		{
			name:        "string valid",
			input:       "999",
			expected:    999,
			expectError: false,
		},
		{
			name:        "nil",
			input:       nil,
			expected:    0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := ToInt16E(tt.input)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("ToInt16E(%v) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ToInt16E(%v) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("ToInt16E(%v) = %v; expected %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}

// TestToInt8E tests the ToInt8E function
func TestToInt8E(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    int8
		expectError bool
	}{
		{
			name:        "int8 positive",
			input:       int8(123),
			expected:    123,
			expectError: false,
		},
		{
			name:        "int",
			input:       42,
			expected:    42,
			expectError: false,
		},
		{
			name:        "string valid",
			input:       "99",
			expected:    99,
			expectError: false,
		},
		{
			name:        "nil",
			input:       nil,
			expected:    0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := ToInt8E(tt.input)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("ToInt8E(%v) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ToInt8E(%v) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("ToInt8E(%v) = %v; expected %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}

// TestToIntE tests the ToIntE function
func TestToIntE(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    int
		expectError bool
	}{
		{
			name:        "int positive",
			input:       12345,
			expected:    12345,
			expectError: false,
		},
		{
			name:        "int64",
			input:       int64(999),
			expected:    999,
			expectError: false,
		},
		{
			name:        "string valid",
			input:       "777",
			expected:    777,
			expectError: false,
		},
		{
			name:        "string invalid",
			input:       "abc",
			expected:    0,
			expectError: true,
		},
		{
			name:        "float64",
			input:       float64(123.99),
			expected:    123,
			expectError: false,
		},
		{
			name:        "bool true",
			input:       true,
			expected:    1,
			expectError: false,
		},
		{
			name:        "bool false",
			input:       false,
			expected:    0,
			expectError: false,
		},
		{
			name:        "nil",
			input:       nil,
			expected:    0,
			expectError: false,
		},
		{
			name:        "time.Weekday",
			input:       time.Monday,
			expected:    1,
			expectError: false,
		},
		{
			name:        "time.Month",
			input:       time.March,
			expected:    3,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := ToIntE(tt.input)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("ToIntE(%v) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ToIntE(%v) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("ToIntE(%v) = %v; expected %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}

// TestToUintE tests the ToUintE function with positive values and negative value handling
func TestToUintE(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    uint
		expectError bool
	}{
		{
			name:        "uint positive",
			input:       uint(12345),
			expected:    12345,
			expectError: false,
		},
		{
			name:        "int positive",
			input:       42,
			expected:    42,
			expectError: false,
		},
		{
			name:        "int negative",
			input:       -10,
			expected:    0,
			expectError: true,
		},
		{
			name:        "int64 negative",
			input:       int64(-100),
			expected:    0,
			expectError: true,
		},
		{
			name:        "string positive",
			input:       "123",
			expected:    123,
			expectError: false,
		},
		{
			name:        "string negative",
			input:       "-50",
			expected:    0,
			expectError: true,
		},
		{
			name:        "float64 positive",
			input:       float64(99.5),
			expected:    99,
			expectError: false,
		},
		{
			name:        "float64 negative",
			input:       float64(-10.5),
			expected:    0,
			expectError: true,
		},
		{
			name:        "bool true",
			input:       true,
			expected:    1,
			expectError: false,
		},
		{
			name:        "bool false",
			input:       false,
			expected:    0,
			expectError: false,
		},
		{
			name:        "nil",
			input:       nil,
			expected:    0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := ToUintE(tt.input)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("ToUintE(%v) expected error but got none", tt.input)
				}
				if !errors.Is(err, errNegativeNotAllowed) && err.Error() != errNegativeNotAllowed.Error() {
					// Check if it's a parse error or negative error
					if result == 0 && err != nil {
						// This is acceptable
					}
				}
			} else {
				if err != nil {
					t.Errorf("ToUintE(%v) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("ToUintE(%v) = %v; expected %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}

// TestToUint64E tests the ToUint64E function
func TestToUint64E(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    uint64
		expectError bool
	}{
		{
			name:        "uint64 positive",
			input:       uint64(12345),
			expected:    12345,
			expectError: false,
		},
		{
			name:        "int positive",
			input:       42,
			expected:    42,
			expectError: false,
		},
		{
			name:        "int negative",
			input:       -10,
			expected:    0,
			expectError: true,
		},
		{
			name:        "string positive",
			input:       "999",
			expected:    999,
			expectError: false,
		},
		{
			name:        "bool true",
			input:       true,
			expected:    1,
			expectError: false,
		},
		{
			name:        "nil",
			input:       nil,
			expected:    0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := ToUint64E(tt.input)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("ToUint64E(%v) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ToUint64E(%v) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("ToUint64E(%v) = %v; expected %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}

// TestToUint32E tests the ToUint32E function
func TestToUint32E(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    uint32
		expectError bool
	}{
		{
			name:        "uint32 positive",
			input:       uint32(12345),
			expected:    12345,
			expectError: false,
		},
		{
			name:        "int positive",
			input:       42,
			expected:    42,
			expectError: false,
		},
		{
			name:        "int negative",
			input:       -10,
			expected:    0,
			expectError: true,
		},
		{
			name:        "nil",
			input:       nil,
			expected:    0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := ToUint32E(tt.input)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("ToUint32E(%v) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ToUint32E(%v) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("ToUint32E(%v) = %v; expected %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}

// TestToUint16E tests the ToUint16E function
func TestToUint16E(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    uint16
		expectError bool
	}{
		{
			name:        "uint16 positive",
			input:       uint16(1234),
			expected:    1234,
			expectError: false,
		},
		{
			name:        "int positive",
			input:       42,
			expected:    42,
			expectError: false,
		},
		{
			name:        "int negative",
			input:       -10,
			expected:    0,
			expectError: true,
		},
		{
			name:        "nil",
			input:       nil,
			expected:    0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := ToUint16E(tt.input)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("ToUint16E(%v) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ToUint16E(%v) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("ToUint16E(%v) = %v; expected %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}

// TestToUint8E tests the ToUint8E function
func TestToUint8E(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    uint8
		expectError bool
	}{
		{
			name:        "uint8 positive",
			input:       uint8(123),
			expected:    123,
			expectError: false,
		},
		{
			name:        "int positive",
			input:       42,
			expected:    42,
			expectError: false,
		},
		{
			name:        "int negative",
			input:       -10,
			expected:    0,
			expectError: true,
		},
		{
			name:        "nil",
			input:       nil,
			expected:    0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := ToUint8E(tt.input)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("ToUint8E(%v) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ToUint8E(%v) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("ToUint8E(%v) = %v; expected %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}

// TestToStringE tests the ToStringE function with various types including
// primitives, template types, and types implementing fmt.Stringer
func TestToStringE(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    string
		expectError bool
	}{
		{
			name:        "string",
			input:       "hello",
			expected:    "hello",
			expectError: false,
		},
		{
			name:        "bool true",
			input:       true,
			expected:    "true",
			expectError: false,
		},
		{
			name:        "bool false",
			input:       false,
			expected:    "false",
			expectError: false,
		},
		{
			name:        "int",
			input:       42,
			expected:    "42",
			expectError: false,
		},
		{
			name:        "int64",
			input:       int64(12345),
			expected:    "12345",
			expectError: false,
		},
		{
			name:        "float64",
			input:       123.456,
			expected:    "123.456",
			expectError: false,
		},
		{
			name:        "float32",
			input:       float32(12.34),
			expected:    "12.34",
			expectError: false,
		},
		{
			name:        "uint",
			input:       uint(100),
			expected:    "100",
			expectError: false,
		},
		{
			name:        "byte slice",
			input:       []byte("hello"),
			expected:    "hello",
			expectError: false,
		},
		{
			name:        "template.HTML",
			input:       template.HTML("<h1>Title</h1>"),
			expected:    "<h1>Title</h1>",
			expectError: false,
		},
		{
			name:        "template.URL",
			input:       template.URL("http://example.com"),
			expected:    "http://example.com",
			expectError: false,
		},
		{
			name:        "nil",
			input:       nil,
			expected:    "",
			expectError: false,
		},
		{
			name:        "json.Number",
			input:       json.Number("123.45"),
			expected:    "123.45",
			expectError: false,
		},
		{
			name:        "error type",
			input:       fmt.Errorf("test error"),
			expected:    "test error",
			expectError: false,
		},
		{
			name:        "unsupported type",
			input:       []int{1, 2, 3},
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := ToStringE(tt.input)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("ToStringE(%v) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ToStringE(%v) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("ToStringE(%v) = %q; expected %q", tt.input, result, tt.expected)
				}
			}
		})
	}
}

// TestToDurationE tests the ToDurationE function with various input types
// including time.Duration, integers, floats, and strings
func TestToDurationE(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    time.Duration
		expectError bool
	}{
		{
			name:        "time.Duration",
			input:       time.Second * 5,
			expected:    time.Second * 5,
			expectError: false,
		},
		{
			name:        "int nanoseconds",
			input:       1000000000,
			expected:    time.Second,
			expectError: false,
		},
		{
			name:        "int64",
			input:       int64(2000000000),
			expected:    time.Second * 2,
			expectError: false,
		},
		{
			name:        "float64",
			input:       float64(3000000000),
			expected:    time.Second * 3,
			expectError: false,
		},
		{
			name:        "string with unit",
			input:       "5s",
			expected:    time.Second * 5,
			expectError: false,
		},
		{
			name:        "string without unit",
			input:       "1000",
			expected:    time.Nanosecond * 1000,
			expectError: false,
		},
		{
			name:        "string minutes",
			input:       "2m",
			expected:    time.Minute * 2,
			expectError: false,
		},
		{
			name:        "string hours",
			input:       "1h",
			expected:    time.Hour,
			expectError: false,
		},
		{
			name:        "invalid string",
			input:       "invalid",
			expected:    0,
			expectError: true,
		},
		{
			name:        "unsupported type",
			input:       []int{1, 2},
			expected:    0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := ToDurationE(tt.input)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("ToDurationE(%v) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ToDurationE(%v) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("ToDurationE(%v) = %v; expected %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}

// TestToTimeE tests the ToTimeE function with various time formats and types
func TestToTimeE(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectError bool
	}{
		{
			name:        "time.Time",
			input:       time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			expectError: false,
		},
		{
			name:        "unix timestamp int",
			input:       int(1672531200),
			expectError: false,
		},
		{
			name:        "unix timestamp int64",
			input:       int64(1672531200),
			expectError: false,
		},
		{
			name:        "RFC3339 string",
			input:       "2023-01-01T00:00:00Z",
			expectError: false,
		},
		{
			name:        "date string",
			input:       "2023-01-01",
			expectError: false,
		},
		{
			name:        "json.Number",
			input:       json.Number("1672531200"),
			expectError: false,
		},
		{
			name:        "invalid string",
			input:       "not a date",
			expectError: true,
		},
		{
			name:        "unsupported type",
			input:       []int{1, 2},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := ToTimeE(tt.input)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("ToTimeE(%v) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ToTimeE(%v) unexpected error: %v", tt.input, err)
				}
				if result.IsZero() && !tt.expectError {
					t.Errorf("ToTimeE(%v) returned zero time unexpectedly", tt.input)
				}
			}
		})
	}
}

// TestToTimeInDefaultLocationE tests the ToTimeInDefaultLocationE function
func TestToTimeInDefaultLocationE(t *testing.T) {
	location := time.FixedZone("TEST", 3600)

	tests := []struct {
		name        string
		input       interface{}
		location    *time.Location
		expectError bool
	}{
		{
			name:        "time.Time with custom location",
			input:       time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			location:    location,
			expectError: false,
		},
		{
			name:        "unix timestamp with custom location",
			input:       int64(1672531200),
			location:    location,
			expectError: false,
		},
		{
			name:        "date string with custom location",
			input:       "2023-01-01",
			location:    location,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := ToTimeInDefaultLocationE(tt.input, tt.location)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("ToTimeInDefaultLocationE(%v) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ToTimeInDefaultLocationE(%v) unexpected error: %v", tt.input, err)
				}
				if result.IsZero() {
					t.Errorf("ToTimeInDefaultLocationE(%v) returned zero time", tt.input)
				}
			}
		})
	}
}

// TestToStringMapStringE tests the ToStringMapStringE function with various map types
func TestToStringMapStringE(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    map[string]string
		expectError bool
	}{
		{
			name:        "map[string]string",
			input:       map[string]string{"key1": "value1", "key2": "value2"},
			expected:    map[string]string{"key1": "value1", "key2": "value2"},
			expectError: false,
		},
		{
			name:        "map[string]interface{}",
			input:       map[string]interface{}{"key1": "value1", "key2": 123},
			expected:    map[string]string{"key1": "value1", "key2": "123"},
			expectError: false,
		},
		{
			name:        "map[interface{}]string",
			input:       map[interface{}]string{"key1": "value1", 2: "value2"},
			expected:    map[string]string{"key1": "value1", "2": "value2"},
			expectError: false,
		},
		{
			name:        "map[interface{}]interface{}",
			input:       map[interface{}]interface{}{"key1": "value1", "key2": 456},
			expected:    map[string]string{"key1": "value1", "key2": "456"},
			expectError: false,
		},
		{
			name:        "JSON string",
			input:       `{"key1":"value1","key2":"value2"}`,
			expected:    map[string]string{"key1": "value1", "key2": "value2"},
			expectError: false,
		},
		{
			name:        "invalid type",
			input:       []int{1, 2, 3},
			expected:    map[string]string{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := ToStringMapStringE(tt.input)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("ToStringMapStringE(%v) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ToStringMapStringE(%v) unexpected error: %v", tt.input, err)
				}
				if len(result) != len(tt.expected) {
					t.Errorf("ToStringMapStringE(%v) returned map with %d elements; expected %d", tt.input, len(result), len(tt.expected))
				}
				for k, v := range tt.expected {
					if result[k] != v {
						t.Errorf("ToStringMapStringE(%v)[%s] = %v; expected %v", tt.input, k, result[k], v)
					}
				}
			}
		})
	}
}

// TestToStringMapStringSliceE tests the ToStringMapStringSliceE function
func TestToStringMapStringSliceE(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectError bool
	}{
		{
			name:        "map[string][]string",
			input:       map[string][]string{"key1": {"a", "b"}, "key2": {"c"}},
			expectError: false,
		},
		{
			name:        "map[string]string",
			input:       map[string]string{"key1": "value1"},
			expectError: false,
		},
		{
			name:        "map[string]interface{} with slice",
			input:       map[string]interface{}{"key1": []string{"a", "b"}},
			expectError: false,
		},
		{
			name:        "invalid type",
			input:       42,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := ToStringMapStringSliceE(tt.input)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("ToStringMapStringSliceE(%v) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ToStringMapStringSliceE(%v) unexpected error: %v", tt.input, err)
				}
				if result == nil {
					t.Errorf("ToStringMapStringSliceE(%v) returned nil", tt.input)
				}
			}
		})
	}
}

// TestToStringMapBoolE tests the ToStringMapBoolE function
func TestToStringMapBoolE(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectError bool
	}{
		{
			name:        "map[string]bool",
			input:       map[string]bool{"key1": true, "key2": false},
			expectError: false,
		},
		{
			name:        "map[string]interface{}",
			input:       map[string]interface{}{"key1": true, "key2": 1},
			expectError: false,
		},
		{
			name:        "map[interface{}]interface{}",
			input:       map[interface{}]interface{}{"key1": true, "key2": false},
			expectError: false,
		},
		{
			name:        "JSON string",
			input:       `{"key1":true,"key2":false}`,
			expectError: false,
		},
		{
			name:        "invalid type",
			input:       []int{1, 2},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := ToStringMapBoolE(tt.input)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("ToStringMapBoolE(%v) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ToStringMapBoolE(%v) unexpected error: %v", tt.input, err)
				}
				if result == nil {
					t.Errorf("ToStringMapBoolE(%v) returned nil", tt.input)
				}
			}
		})
	}
}

// TestToStringMapE tests the ToStringMapE function
func TestToStringMapE(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectError bool
	}{
		{
			name:        "map[string]interface{}",
			input:       map[string]interface{}{"key1": "value1", "key2": 123},
			expectError: false,
		},
		{
			name:        "map[interface{}]interface{}",
			input:       map[interface{}]interface{}{"key1": "value1", 2: 456},
			expectError: false,
		},
		{
			name:        "JSON string",
			input:       `{"key1":"value1","key2":123}`,
			expectError: false,
		},
		{
			name:        "invalid type",
			input:       []int{1, 2},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := ToStringMapE(tt.input)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("ToStringMapE(%v) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ToStringMapE(%v) unexpected error: %v", tt.input, err)
				}
				if result == nil {
					t.Errorf("ToStringMapE(%v) returned nil", tt.input)
				}
			}
		})
	}
}

// TestToStringMapIntE tests the ToStringMapIntE function
func TestToStringMapIntE(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectError bool
	}{
		{
			name:        "map[string]int",
			input:       map[string]int{"key1": 1, "key2": 2},
			expectError: false,
		},
		{
			name:        "map[string]interface{}",
			input:       map[string]interface{}{"key1": 1, "key2": "2"},
			expectError: false,
		},
		{
			name:        "map[interface{}]interface{}",
			input:       map[interface{}]interface{}{"key1": 1, "key2": 2},
			expectError: false,
		},
		{
			name:        "nil",
			input:       nil,
			expectError: true,
		},
		{
			name:        "invalid type",
			input:       "not a map",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := ToStringMapIntE(tt.input)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("ToStringMapIntE(%v) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ToStringMapIntE(%v) unexpected error: %v", tt.input, err)
				}
				if result == nil {
					t.Errorf("ToStringMapIntE(%v) returned nil", tt.input)
				}
			}
		})
	}
}

// TestToStringMapInt64E tests the ToStringMapInt64E function
func TestToStringMapInt64E(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectError bool
	}{
		{
			name:        "map[string]int64",
			input:       map[string]int64{"key1": 1, "key2": 2},
			expectError: false,
		},
		{
			name:        "map[string]interface{}",
			input:       map[string]interface{}{"key1": int64(1), "key2": "2"},
			expectError: false,
		},
		{
			name:        "nil",
			input:       nil,
			expectError: true,
		},
		{
			name:        "invalid type",
			input:       "not a map",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := ToStringMapInt64E(tt.input)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("ToStringMapInt64E(%v) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ToStringMapInt64E(%v) unexpected error: %v", tt.input, err)
				}
				if result == nil {
					t.Errorf("ToStringMapInt64E(%v) returned nil", tt.input)
				}
			}
		})
	}
}

// TestToSliceE tests the ToSliceE function
func TestToSliceE(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectError bool
	}{
		{
			name:        "[]interface{}",
			input:       []interface{}{1, "two", 3.0},
			expectError: false,
		},
		{
			name:        "[]map[string]interface{}",
			input:       []map[string]interface{}{{"key": "value"}},
			expectError: false,
		},
		{
			name:        "invalid type",
			input:       "not a slice",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := ToSliceE(tt.input)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("ToSliceE(%v) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ToSliceE(%v) unexpected error: %v", tt.input, err)
				}
				if result == nil {
					t.Errorf("ToSliceE(%v) returned nil", tt.input)
				}
			}
		})
	}
}

// TestToBoolSliceE tests the ToBoolSliceE function
func TestToBoolSliceE(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    []bool
		expectError bool
	}{
		{
			name:        "[]bool",
			input:       []bool{true, false, true},
			expected:    []bool{true, false, true},
			expectError: false,
		},
		{
			name:        "[]interface{} with bools",
			input:       []interface{}{true, 1, 0, "true"},
			expected:    []bool{true, true, false, true},
			expectError: false,
		},
		{
			name:        "nil",
			input:       nil,
			expected:    []bool{},
			expectError: true,
		},
		{
			name:        "invalid type",
			input:       "not a slice",
			expected:    []bool{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := ToBoolSliceE(tt.input)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("ToBoolSliceE(%v) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ToBoolSliceE(%v) unexpected error: %v", tt.input, err)
				}
				if len(result) != len(tt.expected) {
					t.Errorf("ToBoolSliceE(%v) returned slice with %d elements; expected %d", tt.input, len(result), len(tt.expected))
				}
			}
		})
	}
}

// TestToStringSliceE tests the ToStringSliceE function
func TestToStringSliceE(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    []string
		expectError bool
	}{
		{
			name:        "[]string",
			input:       []string{"a", "b", "c"},
			expected:    []string{"a", "b", "c"},
			expectError: false,
		},
		{
			name:        "[]interface{}",
			input:       []interface{}{"a", 1, true},
			expected:    []string{"a", "1", "true"},
			expectError: false,
		},
		{
			name:        "[]int",
			input:       []int{1, 2, 3},
			expected:    []string{"1", "2", "3"},
			expectError: false,
		},
		{
			name:        "string with fields",
			input:       "hello world test",
			expected:    []string{"hello", "world", "test"},
			expectError: false,
		},
		{
			name:        "[]error",
			input:       []error{fmt.Errorf("error1"), fmt.Errorf("error2")},
			expected:    []string{"error1", "error2"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := ToStringSliceE(tt.input)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("ToStringSliceE(%v) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ToStringSliceE(%v) unexpected error: %v", tt.input, err)
				}
				if len(result) != len(tt.expected) {
					t.Errorf("ToStringSliceE(%v) returned %v; expected %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}

// TestToIntSliceE tests the ToIntSliceE function
func TestToIntSliceE(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    []int
		expectError bool
	}{
		{
			name:        "[]int",
			input:       []int{1, 2, 3},
			expected:    []int{1, 2, 3},
			expectError: false,
		},
		{
			name:        "[]interface{} with numbers",
			input:       []interface{}{1, "2", 3.0},
			expected:    []int{1, 2, 3},
			expectError: false,
		},
		{
			name:        "nil",
			input:       nil,
			expected:    []int{},
			expectError: true,
		},
		{
			name:        "invalid type",
			input:       "not a slice",
			expected:    []int{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := ToIntSliceE(tt.input)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("ToIntSliceE(%v) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ToIntSliceE(%v) unexpected error: %v", tt.input, err)
				}
				if len(result) != len(tt.expected) {
					t.Errorf("ToIntSliceE(%v) returned slice with %d elements; expected %d", tt.input, len(result), len(tt.expected))
				}
			}
		})
	}
}

// TestToDurationSliceE tests the ToDurationSliceE function
func TestToDurationSliceE(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectError bool
	}{
		{
			name:        "[]time.Duration",
			input:       []time.Duration{time.Second, time.Minute},
			expectError: false,
		},
		{
			name:        "[]interface{} with durations",
			input:       []interface{}{time.Second, "1m", 1000000000},
			expectError: false,
		},
		{
			name:        "nil",
			input:       nil,
			expectError: true,
		},
		{
			name:        "invalid type",
			input:       "not a slice",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := ToDurationSliceE(tt.input)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("ToDurationSliceE(%v) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ToDurationSliceE(%v) unexpected error: %v", tt.input, err)
				}
				if result == nil {
					t.Errorf("ToDurationSliceE(%v) returned nil", tt.input)
				}
			}
		})
	}
}

// TestIndirect tests the indirect function for dereferencing pointers
func TestIndirect(t *testing.T) {
	value := 42
	ptr := &value
	ptrPtr := &ptr

	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "nil",
			input:    nil,
			expected: nil,
		},
		{
			name:     "non-pointer value",
			input:    42,
			expected: 42,
		},
		{
			name:     "single pointer",
			input:    ptr,
			expected: 42,
		},
		{
			name:     "double pointer",
			input:    ptrPtr,
			expected: 42,
		},
		{
			name:     "string value",
			input:    "hello",
			expected: "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := indirect(tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("indirect(%v) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestToIntHelper tests the toInt helper function
func TestToIntHelper(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected int
		ok       bool
	}{
		{
			name:     "int value",
			input:    42,
			expected: 42,
			ok:       true,
		},
		{
			name:     "time.Weekday",
			input:    time.Monday,
			expected: 1,
			ok:       true,
		},
		{
			name:     "time.Month",
			input:    time.December,
			expected: 12,
			ok:       true,
		},
		{
			name:     "string value",
			input:    "not an int",
			expected: 0,
			ok:       false,
		},
		{
			name:     "float value",
			input:    3.14,
			expected: 0,
			ok:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, ok := toInt(tt.input)

			// Assert
			if ok != tt.ok {
				t.Errorf("toInt(%v) ok = %v; expected %v", tt.input, ok, tt.ok)
			}
			if result != tt.expected {
				t.Errorf("toInt(%v) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestTrimZeroDecimal tests the trimZeroDecimal helper function
func TestTrimZeroDecimal(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "integer no decimal",
			input:    "123",
			expected: "123",
		},
		{
			name:     "decimal with zeros",
			input:    "123.00",
			expected: "123",
		},
		{
			name:     "decimal with non-zeros",
			input:    "123.45",
			expected: "123.45",
		},
		{
			name:     "decimal with trailing zeros",
			input:    "123.4500",
			expected: "123.4500",
		},
		{
			name:     "only zeros after decimal",
			input:    "0.00",
			expected: "0",
		},
		{
			name:     "single zero",
			input:    "0",
			expected: "0",
		},
		{
			name:     "no decimal point",
			input:    "999",
			expected: "999",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := trimZeroDecimal(tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("trimZeroDecimal(%q) = %q; expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestStringToDate tests the StringToDate function
func TestStringToDate(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "RFC3339 format",
			input:       "2023-01-15T14:30:00Z",
			expectError: false,
		},
		{
			name:        "date only format",
			input:       "2023-01-15",
			expectError: false,
		},
		{
			name:        "RFC1123 format",
			input:       "Sun, 15 Jan 2023 14:30:00 GMT",
			expectError: false,
		},
		{
			name:        "invalid format",
			input:       "not a date",
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := StringToDate(tt.input)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("StringToDate(%q) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("StringToDate(%q) unexpected error: %v", tt.input, err)
				}
				if result.IsZero() {
					t.Errorf("StringToDate(%q) returned zero time", tt.input)
				}
			}
		})
	}
}

// TestStringToDateInDefaultLocation tests the StringToDateInDefaultLocation function
func TestStringToDateInDefaultLocation(t *testing.T) {
	location := time.FixedZone("TEST", 3600)

	tests := []struct {
		name        string
		input       string
		location    *time.Location
		expectError bool
	}{
		{
			name:        "date with custom location",
			input:       "2023-01-15",
			location:    location,
			expectError: false,
		},
		{
			name:        "RFC3339 with custom location",
			input:       "2023-01-15T14:30:00Z",
			location:    location,
			expectError: false,
		},
		{
			name:        "invalid date",
			input:       "invalid",
			location:    time.UTC,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := StringToDateInDefaultLocation(tt.input, tt.location)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("StringToDateInDefaultLocation(%q) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("StringToDateInDefaultLocation(%q) unexpected error: %v", tt.input, err)
				}
				if result.IsZero() {
					t.Errorf("StringToDateInDefaultLocation(%q) returned zero time", tt.input)
				}
			}
		})
	}
}
