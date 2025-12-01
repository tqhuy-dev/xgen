package utilities

import (
	"testing"
)

// TestFormatNumber tests the FormatNumber function
// which formats integers using Vietnamese locale number formatting
func TestFormatNumber(t *testing.T) {
	tests := []struct {
		name     string
		number   int
		expected string
	}{
		{
			name:     "zero",
			number:   0,
			expected: "0",
		},
		{
			name:     "positive single digit",
			number:   5,
			expected: "5",
		},
		{
			name:     "positive two digits",
			number:   42,
			expected: "42",
		},
		{
			name:     "positive three digits",
			number:   123,
			expected: "123",
		},
		{
			name:     "positive with thousands",
			number:   1234,
			expected: "1.234",
		},
		{
			name:     "positive with millions",
			number:   1234567,
			expected: "1.234.567",
		},
		{
			name:     "negative single digit",
			number:   -5,
			expected: "-5",
		},
		{
			name:     "negative with thousands",
			number:   -1234,
			expected: "-1.234",
		},
		{
			name:     "negative with millions",
			number:   -1234567,
			expected: "-1.234.567",
		},
		{
			name:     "max int32",
			number:   2147483647,
			expected: "2.147.483.647",
		},
		{
			name:     "min int32",
			number:   -2147483648,
			expected: "-2.147.483.648",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := FormatNumber(tt.number)

			// Assert
			if result != tt.expected {
				t.Errorf("FormatNumber(%d) = %q; expected %q", tt.number, result, tt.expected)
			}
		})
	}
}

// TestRoundToInteger tests the RoundToInteger function
// which rounds float64 values to the nearest integer
func TestRoundToInteger(t *testing.T) {
	tests := []struct {
		name     string
		number   float64
		expected int
	}{
		{
			name:     "zero",
			number:   0.0,
			expected: 0,
		},
		{
			name:     "positive integer",
			number:   5.0,
			expected: 5,
		},
		{
			name:     "positive round down below 0.5",
			number:   5.4,
			expected: 5,
		},
		{
			name:     "positive round up at 0.5",
			number:   5.5,
			expected: 6,
		},
		{
			name:     "positive round up above 0.5",
			number:   5.6,
			expected: 6,
		},
		{
			name:     "positive round up at 0.9",
			number:   5.9,
			expected: 6,
		},
		{
			name:     "negative integer",
			number:   -5.0,
			expected: -5,
		},
		{
			name:     "negative round up (toward zero) below 0.5",
			number:   -5.4,
			expected: -5,
		},
		{
			name:     "negative round down (away from zero) at 0.5",
			number:   -5.5,
			expected: -6,
		},
		{
			name:     "negative round down at 0.6",
			number:   -5.6,
			expected: -6,
		},
		{
			name:     "small positive decimal",
			number:   0.1,
			expected: 0,
		},
		{
			name:     "small negative decimal",
			number:   -0.1,
			expected: 0,
		},
		{
			name:     "large positive number",
			number:   123456.7,
			expected: 123457,
		},
		{
			name:     "large negative number",
			number:   -123456.7,
			expected: -123457,
		},
		{
			name:     "exactly 0.5 rounds away from zero",
			number:   0.5,
			expected: 1,
		},
		{
			name:     "exactly -0.5 rounds away from zero",
			number:   -0.5,
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := RoundToInteger(tt.number)

			// Assert
			if result != tt.expected {
				t.Errorf("RoundToInteger(%f) = %d; expected %d", tt.number, result, tt.expected)
			}
		})
	}
}

// TestCeilToInt tests the CeilToInt function
// which rounds float64 values up to the nearest integer (ceiling)
func TestCeilToInt(t *testing.T) {
	tests := []struct {
		name     string
		number   float64
		expected int
	}{
		{
			name:     "zero",
			number:   0.0,
			expected: 0,
		},
		{
			name:     "positive integer",
			number:   5.0,
			expected: 5,
		},
		{
			name:     "positive with small decimal",
			number:   5.1,
			expected: 6,
		},
		{
			name:     "positive with large decimal",
			number:   5.9,
			expected: 6,
		},
		{
			name:     "negative integer",
			number:   -5.0,
			expected: -5,
		},
		{
			name:     "negative with small decimal rounds toward zero",
			number:   -5.1,
			expected: -5,
		},
		{
			name:     "negative with large decimal rounds toward zero",
			number:   -5.9,
			expected: -5,
		},
		{
			name:     "very small positive",
			number:   0.001,
			expected: 1,
		},
		{
			name:     "very small negative",
			number:   -0.001,
			expected: 0,
		},
		{
			name:     "large positive",
			number:   999999.1,
			expected: 1000000,
		},
		{
			name:     "large negative",
			number:   -999999.1,
			expected: -999999,
		},
		{
			name:     "positive at 0.5",
			number:   5.5,
			expected: 6,
		},
		{
			name:     "negative at 0.5",
			number:   -5.5,
			expected: -5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := CeilToInt(tt.number)

			// Assert
			if result != tt.expected {
				t.Errorf("CeilToInt(%f) = %d; expected %d", tt.number, result, tt.expected)
			}
		})
	}
}

// TestFloorToInt tests the FloorToInt function
// which rounds float64 values down to the nearest integer (floor)
func TestFloorToInt(t *testing.T) {
	tests := []struct {
		name     string
		number   float64
		expected int
	}{
		{
			name:     "zero",
			number:   0.0,
			expected: 0,
		},
		{
			name:     "positive integer",
			number:   5.0,
			expected: 5,
		},
		{
			name:     "positive with small decimal",
			number:   5.1,
			expected: 5,
		},
		{
			name:     "positive with large decimal",
			number:   5.9,
			expected: 5,
		},
		{
			name:     "negative integer",
			number:   -5.0,
			expected: -5,
		},
		{
			name:     "negative with small decimal rounds away from zero",
			number:   -5.1,
			expected: -6,
		},
		{
			name:     "negative with large decimal rounds away from zero",
			number:   -5.9,
			expected: -6,
		},
		{
			name:     "very small positive",
			number:   0.001,
			expected: 0,
		},
		{
			name:     "very small negative",
			number:   -0.001,
			expected: -1,
		},
		{
			name:     "large positive",
			number:   999999.9,
			expected: 999999,
		},
		{
			name:     "large negative",
			number:   -999999.1,
			expected: -1000000,
		},
		{
			name:     "positive at 0.5",
			number:   5.5,
			expected: 5,
		},
		{
			name:     "negative at 0.5",
			number:   -5.5,
			expected: -6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := FloorToInt(tt.number)

			// Assert
			if result != tt.expected {
				t.Errorf("FloorToInt(%f) = %d; expected %d", tt.number, result, tt.expected)
			}
		})
	}
}

// TestRoundByThreshold tests the RoundByThreshold function
// which rounds a number based on a custom threshold within a given unit
func TestRoundByThreshold(t *testing.T) {
	tests := []struct {
		name      string
		n         int
		unit      int
		threshold int
		expected  int
	}{
		{
			name:      "round down when below threshold",
			n:         123,
			unit:      10,
			threshold: 5,
			expected:  120,
		},
		{
			name:      "round up when at threshold",
			n:         125,
			unit:      10,
			threshold: 5,
			expected:  130,
		},
		{
			name:      "round up when above threshold",
			n:         127,
			unit:      10,
			threshold: 5,
			expected:  130,
		},
		{
			name:      "exact multiple no rounding",
			n:         120,
			unit:      10,
			threshold: 5,
			expected:  120,
		},
		{
			name:      "round to nearest hundred",
			n:         1234,
			unit:      100,
			threshold: 50,
			expected:  1200,
		},
		{
			name:      "round up to nearest hundred",
			n:         1250,
			unit:      100,
			threshold: 50,
			expected:  1300,
		},
		{
			name:      "zero value",
			n:         0,
			unit:      10,
			threshold: 5,
			expected:  0,
		},
		{
			name:      "negative number behavior",
			n:         -123,
			unit:      10,
			threshold: 5,
			expected:  -120, // remainder is -3, which is < 5, so rounds down to -120
		},
		{
			name:      "negative number with different remainder",
			n:         -127,
			unit:      10,
			threshold: 5,
			expected:  -120, // remainder is -7, which is < 5, so rounds down to -120
		},
		{
			name:      "unit of 1 always returns n",
			n:         123,
			unit:      1,
			threshold: 1,
			expected:  123,
		},
		{
			name:      "zero unit returns original (avoid division by zero)",
			n:         123,
			unit:      0,
			threshold: 5,
			expected:  123,
		},
		{
			name:      "negative unit returns original (avoid issues)",
			n:         123,
			unit:      -10,
			threshold: 5,
			expected:  123,
		},
		{
			name:      "threshold of 0 always rounds down",
			n:         125,
			unit:      10,
			threshold: 0,
			expected:  130,
		},
		{
			name:      "threshold equals unit",
			n:         121,
			unit:      10,
			threshold: 10,
			expected:  120, // remainder is 1, which is < 10, so rounds down to 120
		},
		{
			name:      "large numbers",
			n:         987654,
			unit:      1000,
			threshold: 500,
			expected:  988000,
		},
		{
			name:      "medium numbers",
			n:         12301,
			unit:      1000,
			threshold: 300,
			expected:  13000,
		},
		{
			name:      "medium numbers",
			n:         12298,
			unit:      1000,
			threshold: 300,
			expected:  12000,
		},
		{
			name:      "round to nearest 5",
			n:         23,
			unit:      5,
			threshold: 3,
			expected:  25,
		},
		{
			name:      "round to nearest 5 down",
			n:         22,
			unit:      5,
			threshold: 3,
			expected:  20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := RoundByThreshold(tt.n, tt.unit, tt.threshold)

			// Assert
			if result != tt.expected {
				t.Errorf("RoundByThreshold(%d, %d, %d) = %d; expected %d",
					tt.n, tt.unit, tt.threshold, result, tt.expected)
			}
		})
	}
}

// TestRoundByThreshold_EdgeCases tests edge cases and boundary conditions
// for the RoundByThreshold function
func TestRoundByThreshold_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		n         int
		unit      int
		threshold int
		expected  int
	}{
		{
			name:      "remainder exactly at threshold",
			n:         155,
			unit:      100,
			threshold: 55,
			expected:  200,
		},
		{
			name:      "remainder one below threshold",
			n:         154,
			unit:      100,
			threshold: 55,
			expected:  100,
		},
		{
			name:      "remainder one above threshold",
			n:         156,
			unit:      100,
			threshold: 55,
			expected:  200,
		},
		{
			name:      "very large unit",
			n:         50000,
			unit:      100000,
			threshold: 60000,
			expected:  0,
		},
		{
			name:      "threshold larger than unit",
			n:         125,
			unit:      10,
			threshold: 15,
			expected:  120, // remainder is 5, which is < 15, so rounds down to 120
		},
		{
			name:      "negative threshold",
			n:         125,
			unit:      10,
			threshold: -5,
			expected:  130,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := RoundByThreshold(tt.n, tt.unit, tt.threshold)

			// Assert
			if result != tt.expected {
				t.Errorf("RoundByThreshold(%d, %d, %d) = %d; expected %d",
					tt.n, tt.unit, tt.threshold, result, tt.expected)
			}
		})
	}
}

// TestRoundingConsistency tests that different rounding functions
// produce consistent results for known values
func TestRoundingConsistency(t *testing.T) {
	tests := []struct {
		name          string
		value         float64
		expectedRound int
		expectedCeil  int
		expectedFloor int
	}{
		{
			name:          "5.5 consistency",
			value:         5.5,
			expectedRound: 6,
			expectedCeil:  6,
			expectedFloor: 5,
		},
		{
			name:          "5.4 consistency",
			value:         5.4,
			expectedRound: 5,
			expectedCeil:  6,
			expectedFloor: 5,
		},
		{
			name:          "5.6 consistency",
			value:         5.6,
			expectedRound: 6,
			expectedCeil:  6,
			expectedFloor: 5,
		},
		{
			name:          "negative 5.5 consistency",
			value:         -5.5,
			expectedRound: -6,
			expectedCeil:  -5,
			expectedFloor: -6,
		},
		{
			name:          "zero consistency",
			value:         0.0,
			expectedRound: 0,
			expectedCeil:  0,
			expectedFloor: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			round := RoundToInteger(tt.value)
			ceil := CeilToInt(tt.value)
			floor := FloorToInt(tt.value)

			// Assert
			if round != tt.expectedRound {
				t.Errorf("RoundToInteger(%f) = %d; expected %d", tt.value, round, tt.expectedRound)
			}
			if ceil != tt.expectedCeil {
				t.Errorf("CeilToInt(%f) = %d; expected %d", tt.value, ceil, tt.expectedCeil)
			}
			if floor != tt.expectedFloor {
				t.Errorf("FloorToInt(%f) = %d; expected %d", tt.value, floor, tt.expectedFloor)
			}

			// Additional consistency checks
			if ceil < floor {
				t.Errorf("CeilToInt(%f) = %d is less than FloorToInt(%f) = %d",
					tt.value, ceil, tt.value, floor)
			}
		})
	}
}

// TestFormatNumber_LargeNumbers tests FormatNumber with very large numbers
// to ensure proper formatting of thousands separators
func TestFormatNumber_LargeNumbers(t *testing.T) {
	tests := []struct {
		name     string
		number   int
		expected string
	}{
		{
			name:     "one billion",
			number:   1000000000,
			expected: "1.000.000.000",
		},
		{
			name:     "complex large number",
			number:   123456789,
			expected: "123.456.789",
		},
		{
			name:     "negative billion",
			number:   -1000000000,
			expected: "-1.000.000.000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := FormatNumber(tt.number)

			// Assert
			if result != tt.expected {
				t.Errorf("FormatNumber(%d) = %q; expected %q", tt.number, result, tt.expected)
			}
		})
	}
}
