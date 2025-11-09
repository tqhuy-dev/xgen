package main

import (
	"testing"
	"time"
)

// TestToBool tests the ToBool wrapper function
func TestToBool(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected bool
	}{
		{
			name:     "bool true",
			input:    true,
			expected: true,
		},
		{
			name:     "bool false",
			input:    false,
			expected: false,
		},
		{
			name:     "int non-zero",
			input:    1,
			expected: true,
		},
		{
			name:     "int zero",
			input:    0,
			expected: false,
		},
		{
			name:     "string true",
			input:    "true",
			expected: true,
		},
		{
			name:     "invalid string returns false",
			input:    "invalid",
			expected: false,
		},
		{
			name:     "nil returns false",
			input:    nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := ToBool(tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("ToBool(%v) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestToTime tests the ToTime wrapper function
func TestToTime(t *testing.T) {
	testTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		input         interface{}
		expectNonZero bool
	}{
		{
			name:          "time.Time",
			input:         testTime,
			expectNonZero: true,
		},
		{
			name:          "unix timestamp",
			input:         int64(1672531200),
			expectNonZero: true,
		},
		{
			name:          "date string",
			input:         "2023-01-01",
			expectNonZero: true,
		},
		{
			name:          "invalid input returns zero",
			input:         "invalid",
			expectNonZero: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := ToTime(tt.input)

			// Assert
			if tt.expectNonZero && result.IsZero() {
				t.Errorf("ToTime(%v) returned zero time, expected non-zero", tt.input)
			}
			if !tt.expectNonZero && !result.IsZero() {
				t.Errorf("ToTime(%v) returned non-zero time, expected zero", tt.input)
			}
		})
	}
}

// TestToTimeInDefaultLocation tests the ToTimeInDefaultLocation wrapper function
func TestToTimeInDefaultLocation(t *testing.T) {
	location := time.FixedZone("TEST", 3600)
	testTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		input         interface{}
		location      *time.Location
		expectNonZero bool
	}{
		{
			name:          "time.Time with custom location",
			input:         testTime,
			location:      location,
			expectNonZero: true,
		},
		{
			name:          "unix timestamp with custom location",
			input:         int64(1672531200),
			location:      location,
			expectNonZero: true,
		},
		{
			name:          "date string with custom location",
			input:         "2023-01-01",
			location:      location,
			expectNonZero: true,
		},
		{
			name:          "invalid input returns zero",
			input:         []int{1, 2},
			location:      time.UTC,
			expectNonZero: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := ToTimeInDefaultLocation(tt.input, tt.location)

			// Assert
			if tt.expectNonZero && result.IsZero() {
				t.Errorf("ToTimeInDefaultLocation(%v) returned zero time, expected non-zero", tt.input)
			}
			if !tt.expectNonZero && !result.IsZero() {
				t.Errorf("ToTimeInDefaultLocation(%v) returned non-zero time, expected zero", tt.input)
			}
		})
	}
}

// TestToDuration tests the ToDuration wrapper function
func TestToDuration(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected time.Duration
	}{
		{
			name:     "time.Duration",
			input:    time.Second * 5,
			expected: time.Second * 5,
		},
		{
			name:     "int nanoseconds",
			input:    1000000000,
			expected: time.Second,
		},
		{
			name:     "string with unit",
			input:    "5s",
			expected: time.Second * 5,
		},
		{
			name:     "invalid input returns zero",
			input:    "invalid",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := ToDuration(tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("ToDuration(%v) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestToFloat64 tests the ToFloat64 wrapper function
func TestToFloat64(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected float64
	}{
		{
			name:     "float64",
			input:    123.456,
			expected: 123.456,
		},
		{
			name:     "int",
			input:    42,
			expected: 42.0,
		},
		{
			name:     "string",
			input:    "123.45",
			expected: 123.45,
		},
		{
			name:     "bool true",
			input:    true,
			expected: 1.0,
		},
		{
			name:     "bool false",
			input:    false,
			expected: 0.0,
		},
		{
			name:     "invalid input returns zero",
			input:    "invalid",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := ToFloat64(tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("ToFloat64(%v) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestToFloat32 tests the ToFloat32 wrapper function
func TestToFloat32(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected float32
	}{
		{
			name:     "float32",
			input:    float32(123.456),
			expected: 123.456,
		},
		{
			name:     "int",
			input:    42,
			expected: 42.0,
		},
		{
			name:     "string",
			input:    "123.45",
			expected: 123.45,
		},
		{
			name:     "invalid input returns zero",
			input:    "invalid",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := ToFloat32(tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("ToFloat32(%v) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestToInt64 tests the ToInt64 wrapper function
func TestToInt64(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected int64
	}{
		{
			name:     "int64",
			input:    int64(12345),
			expected: 12345,
		},
		{
			name:     "int",
			input:    42,
			expected: 42,
		},
		{
			name:     "string",
			input:    "999",
			expected: 999,
		},
		{
			name:     "float64",
			input:    float64(100.99),
			expected: 100,
		},
		{
			name:     "bool true",
			input:    true,
			expected: 1,
		},
		{
			name:     "invalid input returns zero",
			input:    "invalid",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := ToInt64(tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("ToInt64(%v) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestToInt32 tests the ToInt32 wrapper function
func TestToInt32(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected int32
	}{
		{
			name:     "int32",
			input:    int32(12345),
			expected: 12345,
		},
		{
			name:     "int",
			input:    42,
			expected: 42,
		},
		{
			name:     "string",
			input:    "999",
			expected: 999,
		},
		{
			name:     "invalid input returns zero",
			input:    "invalid",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := ToInt32(tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("ToInt32(%v) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestToInt16 tests the ToInt16 wrapper function
func TestToInt16(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected int16
	}{
		{
			name:     "int16",
			input:    int16(1234),
			expected: 1234,
		},
		{
			name:     "int",
			input:    42,
			expected: 42,
		},
		{
			name:     "string",
			input:    "999",
			expected: 999,
		},
		{
			name:     "invalid input returns zero",
			input:    "invalid",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := ToInt16(tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("ToInt16(%v) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestToInt8 tests the ToInt8 wrapper function
func TestToInt8(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected int8
	}{
		{
			name:     "int8",
			input:    int8(123),
			expected: 123,
		},
		{
			name:     "int",
			input:    42,
			expected: 42,
		},
		{
			name:     "string",
			input:    "99",
			expected: 99,
		},
		{
			name:     "invalid input returns zero",
			input:    "invalid",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := ToInt8(tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("ToInt8(%v) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestToInt tests the ToInt wrapper function
func TestToInt(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected int
	}{
		{
			name:     "int",
			input:    12345,
			expected: 12345,
		},
		{
			name:     "int64",
			input:    int64(999),
			expected: 999,
		},
		{
			name:     "string",
			input:    "777",
			expected: 777,
		},
		{
			name:     "float64",
			input:    float64(123.99),
			expected: 123,
		},
		{
			name:     "bool true",
			input:    true,
			expected: 1,
		},
		{
			name:     "time.Weekday",
			input:    time.Monday,
			expected: 1,
		},
		{
			name:     "invalid input returns zero",
			input:    "invalid",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := ToInt(tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("ToInt(%v) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestToUint tests the ToUint wrapper function
func TestToUint(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected uint
	}{
		{
			name:     "uint",
			input:    uint(12345),
			expected: 12345,
		},
		{
			name:     "int positive",
			input:    42,
			expected: 42,
		},
		{
			name:     "int negative returns zero",
			input:    -10,
			expected: 0,
		},
		{
			name:     "string",
			input:    "999",
			expected: 999,
		},
		{
			name:     "bool true",
			input:    true,
			expected: 1,
		},
		{
			name:     "invalid input returns zero",
			input:    "invalid",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := ToUint(tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("ToUint(%v) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestToUint64 tests the ToUint64 wrapper function
func TestToUint64(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected uint64
	}{
		{
			name:     "uint64",
			input:    uint64(12345),
			expected: 12345,
		},
		{
			name:     "int positive",
			input:    42,
			expected: 42,
		},
		{
			name:     "int negative returns zero",
			input:    -10,
			expected: 0,
		},
		{
			name:     "string",
			input:    "999",
			expected: 999,
		},
		{
			name:     "invalid input returns zero",
			input:    "invalid",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := ToUint64(tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("ToUint64(%v) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestToUint32 tests the ToUint32 wrapper function
func TestToUint32(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected uint32
	}{
		{
			name:     "uint32",
			input:    uint32(12345),
			expected: 12345,
		},
		{
			name:     "int positive",
			input:    42,
			expected: 42,
		},
		{
			name:     "int negative returns zero",
			input:    -10,
			expected: 0,
		},
		{
			name:     "invalid input returns zero",
			input:    "invalid",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := ToUint32(tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("ToUint32(%v) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestToUint16 tests the ToUint16 wrapper function
func TestToUint16(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected uint16
	}{
		{
			name:     "uint16",
			input:    uint16(1234),
			expected: 1234,
		},
		{
			name:     "int positive",
			input:    42,
			expected: 42,
		},
		{
			name:     "int negative returns zero",
			input:    -10,
			expected: 0,
		},
		{
			name:     "invalid input returns zero",
			input:    "invalid",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := ToUint16(tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("ToUint16(%v) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestToUint8 tests the ToUint8 wrapper function
func TestToUint8(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected uint8
	}{
		{
			name:     "uint8",
			input:    uint8(123),
			expected: 123,
		},
		{
			name:     "int positive",
			input:    42,
			expected: 42,
		},
		{
			name:     "int negative returns zero",
			input:    -10,
			expected: 0,
		},
		{
			name:     "invalid input returns zero",
			input:    "invalid",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := ToUint8(tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("ToUint8(%v) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestToString tests the ToString wrapper function
func TestToString(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "string",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "int",
			input:    42,
			expected: "42",
		},
		{
			name:     "float64",
			input:    123.456,
			expected: "123.456",
		},
		{
			name:     "bool true",
			input:    true,
			expected: "true",
		},
		{
			name:     "bool false",
			input:    false,
			expected: "false",
		},
		{
			name:     "nil",
			input:    nil,
			expected: "",
		},
		{
			name:     "invalid input returns empty",
			input:    []int{1, 2},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := ToString(tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("ToString(%v) = %q; expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestToStringMapString tests the ToStringMapString wrapper function
func TestToStringMapString(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectEmpty bool
	}{
		{
			name:        "map[string]string",
			input:       map[string]string{"key1": "value1"},
			expectEmpty: false,
		},
		{
			name:        "map[string]interface{}",
			input:       map[string]interface{}{"key1": "value1"},
			expectEmpty: false,
		},
		{
			name:        "invalid input returns empty map",
			input:       []int{1, 2},
			expectEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := ToStringMapString(tt.input)

			// Assert
			if tt.expectEmpty && len(result) != 0 {
				t.Errorf("ToStringMapString(%v) returned non-empty map; expected empty", tt.input)
			}
			if !tt.expectEmpty && len(result) == 0 {
				t.Errorf("ToStringMapString(%v) returned empty map; expected non-empty", tt.input)
			}
		})
	}
}

// TestToStringMapStringSlice tests the ToStringMapStringSlice wrapper function
func TestToStringMapStringSlice(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectEmpty bool
	}{
		{
			name:        "map[string][]string",
			input:       map[string][]string{"key1": {"a", "b"}},
			expectEmpty: false,
		},
		{
			name:        "map[string]string",
			input:       map[string]string{"key1": "value1"},
			expectEmpty: false,
		},
		{
			name:        "invalid input returns empty map",
			input:       42,
			expectEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := ToStringMapStringSlice(tt.input)

			// Assert
			if tt.expectEmpty && len(result) != 0 {
				t.Errorf("ToStringMapStringSlice(%v) returned non-empty map; expected empty", tt.input)
			}
			if !tt.expectEmpty && len(result) == 0 {
				t.Errorf("ToStringMapStringSlice(%v) returned empty map; expected non-empty", tt.input)
			}
		})
	}
}

// TestToStringMapBool tests the ToStringMapBool wrapper function
func TestToStringMapBool(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectEmpty bool
	}{
		{
			name:        "map[string]bool",
			input:       map[string]bool{"key1": true},
			expectEmpty: false,
		},
		{
			name:        "map[string]interface{}",
			input:       map[string]interface{}{"key1": true},
			expectEmpty: false,
		},
		{
			name:        "invalid input returns empty map",
			input:       []int{1, 2},
			expectEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := ToStringMapBool(tt.input)

			// Assert
			if tt.expectEmpty && len(result) != 0 {
				t.Errorf("ToStringMapBool(%v) returned non-empty map; expected empty", tt.input)
			}
			if !tt.expectEmpty && len(result) == 0 {
				t.Errorf("ToStringMapBool(%v) returned empty map; expected non-empty", tt.input)
			}
		})
	}
}

// TestToStringMapInt tests the ToStringMapInt wrapper function
func TestToStringMapInt(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectEmpty bool
	}{
		{
			name:        "map[string]int",
			input:       map[string]int{"key1": 1},
			expectEmpty: false,
		},
		{
			name:        "map[string]interface{}",
			input:       map[string]interface{}{"key1": 1},
			expectEmpty: false,
		},
		{
			name:        "invalid input returns empty map",
			input:       "not a map",
			expectEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := ToStringMapInt(tt.input)

			// Assert
			if tt.expectEmpty && len(result) != 0 {
				t.Errorf("ToStringMapInt(%v) returned non-empty map; expected empty", tt.input)
			}
			if !tt.expectEmpty && len(result) == 0 {
				t.Errorf("ToStringMapInt(%v) returned empty map; expected non-empty", tt.input)
			}
		})
	}
}

// TestToStringMapInt64 tests the ToStringMapInt64 wrapper function
func TestToStringMapInt64(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectEmpty bool
	}{
		{
			name:        "map[string]int64",
			input:       map[string]int64{"key1": 1},
			expectEmpty: false,
		},
		{
			name:        "map[string]interface{}",
			input:       map[string]interface{}{"key1": int64(1)},
			expectEmpty: false,
		},
		{
			name:        "invalid input returns empty map",
			input:       "not a map",
			expectEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := ToStringMapInt64(tt.input)

			// Assert
			if tt.expectEmpty && len(result) != 0 {
				t.Errorf("ToStringMapInt64(%v) returned non-empty map; expected empty", tt.input)
			}
			if !tt.expectEmpty && len(result) == 0 {
				t.Errorf("ToStringMapInt64(%v) returned empty map; expected non-empty", tt.input)
			}
		})
	}
}

// TestToStringMap tests the ToStringMap wrapper function
func TestToStringMap(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectEmpty bool
	}{
		{
			name:        "map[string]interface{}",
			input:       map[string]interface{}{"key1": "value1"},
			expectEmpty: false,
		},
		{
			name:        "map[interface{}]interface{}",
			input:       map[interface{}]interface{}{"key1": "value1"},
			expectEmpty: false,
		},
		{
			name:        "invalid input returns empty map",
			input:       []int{1, 2},
			expectEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := ToStringMap(tt.input)

			// Assert
			if tt.expectEmpty && len(result) != 0 {
				t.Errorf("ToStringMap(%v) returned non-empty map; expected empty", tt.input)
			}
			if !tt.expectEmpty && len(result) == 0 {
				t.Errorf("ToStringMap(%v) returned empty map; expected non-empty", tt.input)
			}
		})
	}
}

// TestToSlice tests the ToSlice wrapper function
func TestToSlice(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectEmpty bool
	}{
		{
			name:        "[]interface{}",
			input:       []interface{}{1, "two", 3.0},
			expectEmpty: false,
		},
		{
			name:        "[]map[string]interface{}",
			input:       []map[string]interface{}{{"key": "value"}},
			expectEmpty: false,
		},
		{
			name:        "invalid input returns empty slice",
			input:       "not a slice",
			expectEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := ToSlice(tt.input)

			// Assert
			if tt.expectEmpty && len(result) != 0 {
				t.Errorf("ToSlice(%v) returned non-empty slice; expected empty", tt.input)
			}
			if !tt.expectEmpty && len(result) == 0 {
				t.Errorf("ToSlice(%v) returned empty slice; expected non-empty", tt.input)
			}
		})
	}
}

// TestToBoolSlice tests the ToBoolSlice wrapper function
func TestToBoolSlice(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectEmpty bool
	}{
		{
			name:        "[]bool",
			input:       []bool{true, false, true},
			expectEmpty: false,
		},
		{
			name:        "[]interface{} with bools",
			input:       []interface{}{true, 1, 0},
			expectEmpty: false,
		},
		{
			name:        "invalid input returns empty slice",
			input:       "not a slice",
			expectEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := ToBoolSlice(tt.input)

			// Assert
			if tt.expectEmpty && len(result) != 0 {
				t.Errorf("ToBoolSlice(%v) returned non-empty slice; expected empty", tt.input)
			}
			if !tt.expectEmpty && len(result) == 0 {
				t.Errorf("ToBoolSlice(%v) returned empty slice; expected non-empty", tt.input)
			}
		})
	}
}

// TestToStringSlice tests the ToStringSlice wrapper function
func TestToStringSlice(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectEmpty bool
	}{
		{
			name:        "[]string",
			input:       []string{"a", "b", "c"},
			expectEmpty: false,
		},
		{
			name:        "[]interface{}",
			input:       []interface{}{"a", 1, true},
			expectEmpty: false,
		},
		{
			name:        "string with fields",
			input:       "hello world",
			expectEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := ToStringSlice(tt.input)

			// Assert
			if tt.expectEmpty && len(result) != 0 {
				t.Errorf("ToStringSlice(%v) returned non-empty slice; expected empty", tt.input)
			}
			if !tt.expectEmpty && len(result) == 0 {
				t.Errorf("ToStringSlice(%v) returned empty slice; expected non-empty", tt.input)
			}
		})
	}
}

// TestToIntSlice tests the ToIntSlice wrapper function
func TestToIntSlice(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectEmpty bool
	}{
		{
			name:        "[]int",
			input:       []int{1, 2, 3},
			expectEmpty: false,
		},
		{
			name:        "[]interface{} with numbers",
			input:       []interface{}{1, "2", 3.0},
			expectEmpty: false,
		},
		{
			name:        "invalid input returns empty slice",
			input:       "not a slice",
			expectEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := ToIntSlice(tt.input)

			// Assert
			if tt.expectEmpty && len(result) != 0 {
				t.Errorf("ToIntSlice(%v) returned non-empty slice; expected empty", tt.input)
			}
			if !tt.expectEmpty && len(result) == 0 {
				t.Errorf("ToIntSlice(%v) returned empty slice; expected non-empty", tt.input)
			}
		})
	}
}

// TestToDurationSlice tests the ToDurationSlice wrapper function
func TestToDurationSlice(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectEmpty bool
	}{
		{
			name:        "[]time.Duration",
			input:       []time.Duration{time.Second, time.Minute},
			expectEmpty: false,
		},
		{
			name:        "[]interface{} with durations",
			input:       []interface{}{time.Second, "1m"},
			expectEmpty: false,
		},
		{
			name:        "invalid input returns empty slice",
			input:       "not a slice",
			expectEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := ToDurationSlice(tt.input)

			// Assert
			if tt.expectEmpty && len(result) != 0 {
				t.Errorf("ToDurationSlice(%v) returned non-empty slice; expected empty", tt.input)
			}
			if !tt.expectEmpty && len(result) == 0 {
				t.Errorf("ToDurationSlice(%v) returned empty slice; expected non-empty", tt.input)
			}
		})
	}
}

// TestStringParseE tests the generic StringParseE function with various types
func TestStringParseE(t *testing.T) {
	t.Run("bool type", func(t *testing.T) {
		tests := []struct {
			name        string
			input       string
			expected    bool
			expectError bool
		}{
			{
				name:        "true lowercase",
				input:       "true",
				expected:    true,
				expectError: false,
			},
			{
				name:        "TRUE uppercase",
				input:       "TRUE",
				expected:    true,
				expectError: false,
			},
			{
				name:        "false lowercase",
				input:       "false",
				expected:    false,
				expectError: false,
			},
			{
				name:        "invalid bool",
				input:       "invalid",
				expected:    false,
				expectError: true,
			},
			{
				name:        "whitespace trimmed",
				input:       "  true  ",
				expected:    true,
				expectError: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Act
				result, err := StringParseE[bool](tt.input)

				// Assert
				if tt.expectError {
					if err == nil {
						t.Errorf("StringParseE[bool](%q) expected error but got none", tt.input)
					}
				} else {
					if err != nil {
						t.Errorf("StringParseE[bool](%q) unexpected error: %v", tt.input, err)
					}
					if result != tt.expected {
						t.Errorf("StringParseE[bool](%q) = %v; expected %v", tt.input, result, tt.expected)
					}
				}
			})
		}
	})

	t.Run("int64 type", func(t *testing.T) {
		tests := []struct {
			name        string
			input       string
			expected    int64
			expectError bool
		}{
			{
				name:        "positive integer",
				input:       "123",
				expected:    123,
				expectError: false,
			},
			{
				name:        "negative integer",
				input:       "-456",
				expected:    -456,
				expectError: false,
			},
			{
				name:        "zero",
				input:       "0",
				expected:    0,
				expectError: false,
			},
			{
				name:        "whitespace trimmed",
				input:       "  789  ",
				expected:    789,
				expectError: false,
			},
			{
				name:        "invalid integer",
				input:       "abc",
				expected:    0,
				expectError: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Act
				result, err := StringParseE[int64](tt.input)

				// Assert
				if tt.expectError {
					if err == nil {
						t.Errorf("StringParseE[int64](%q) expected error but got none", tt.input)
					}
				} else {
					if err != nil {
						t.Errorf("StringParseE[int64](%q) unexpected error: %v", tt.input, err)
					}
					if result != tt.expected {
						t.Errorf("StringParseE[int64](%q) = %v; expected %v", tt.input, result, tt.expected)
					}
				}
			})
		}
	})

	t.Run("uint64 type", func(t *testing.T) {
		tests := []struct {
			name        string
			input       string
			expected    uint64
			expectError bool
		}{
			{
				name:        "positive unsigned integer",
				input:       "123",
				expected:    123,
				expectError: false,
			},
			{
				name:        "zero",
				input:       "0",
				expected:    0,
				expectError: false,
			},
			{
				name:        "negative number",
				input:       "-123",
				expected:    0,
				expectError: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Act
				result, err := StringParseE[uint64](tt.input)

				// Assert
				if tt.expectError {
					if err == nil {
						t.Errorf("StringParseE[uint64](%q) expected error but got none", tt.input)
					}
				} else {
					if err != nil {
						t.Errorf("StringParseE[uint64](%q) unexpected error: %v", tt.input, err)
					}
					if result != tt.expected {
						t.Errorf("StringParseE[uint64](%q) = %v; expected %v", tt.input, result, tt.expected)
					}
				}
			})
		}
	})

	t.Run("float64 type", func(t *testing.T) {
		tests := []struct {
			name        string
			input       string
			expected    float64
			expectError bool
		}{
			{
				name:        "positive float",
				input:       "123.456",
				expected:    123.456,
				expectError: false,
			},
			{
				name:        "negative float",
				input:       "-78.9",
				expected:    -78.9,
				expectError: false,
			},
			{
				name:        "integer as float",
				input:       "100",
				expected:    100.0,
				expectError: false,
			},
			{
				name:        "invalid float",
				input:       "abc",
				expected:    0,
				expectError: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Act
				result, err := StringParseE[float64](tt.input)

				// Assert
				if tt.expectError {
					if err == nil {
						t.Errorf("StringParseE[float64](%q) expected error but got none", tt.input)
					}
				} else {
					if err != nil {
						t.Errorf("StringParseE[float64](%q) unexpected error: %v", tt.input, err)
					}
					if result != tt.expected {
						t.Errorf("StringParseE[float64](%q) = %v; expected %v", tt.input, result, tt.expected)
					}
				}
			})
		}
	})

	t.Run("string type", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected string
		}{
			{
				name:     "simple string",
				input:    "hello",
				expected: "hello",
			},
			{
				name:     "whitespace trimmed",
				input:    "  world  ",
				expected: "world",
			},
			{
				name:     "empty string",
				input:    "   ",
				expected: "",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Act
				result, err := StringParseE[string](tt.input)

				// Assert
				if err != nil {
					t.Errorf("StringParseE[string](%q) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("StringParseE[string](%q) = %q; expected %q", tt.input, result, tt.expected)
				}
			})
		}
	})
}

// TestStringParse tests the generic StringParse function (non-error version)
func TestStringParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{
			name:     "valid integer",
			input:    "42",
			expected: 42,
		},
		{
			name:     "invalid integer returns zero",
			input:    "invalid",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := StringParse[int64](tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("StringParse[int64](%q) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}
