package cryptors

import (
	"testing"
)

// TestGenerateRandomBytes tests the GenerateRandomBytes function with various input lengths
// including valid lengths, zero, and negative values
func TestGenerateRandomBytes(t *testing.T) {
	tests := []struct {
		name        string
		length      int
		expectError bool
	}{
		{
			name:        "valid length 16 bytes",
			length:      16,
			expectError: false,
		},
		{
			name:        "valid length 32 bytes",
			length:      32,
			expectError: false,
		},
		{
			name:        "valid length 1 byte",
			length:      1,
			expectError: false,
		},
		{
			name:        "valid length 64 bytes",
			length:      64,
			expectError: false,
		},
		{
			name:        "zero length should error",
			length:      0,
			expectError: true,
		},
		{
			name:        "negative length should error",
			length:      -1,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: execute the function
			result, err := GenerateRandomBytes(tt.length)

			// Assert: verify error expectation
			if tt.expectError {
				if err == nil {
					t.Errorf("GenerateRandomBytes(%d) expected error but got nil", tt.length)
				}
				return
			}

			if err != nil {
				t.Errorf("GenerateRandomBytes(%d) unexpected error: %v", tt.length, err)
				return
			}

			// Verify the length of returned bytes
			if len(result) != tt.length {
				t.Errorf("GenerateRandomBytes(%d) returned %d bytes; expected %d bytes",
					tt.length, len(result), tt.length)
			}
		})
	}
}

// TestGenerateRandomBytesUniqueness tests that GenerateRandomBytes produces unique values
// on consecutive calls (randomness check)
func TestGenerateRandomBytesUniqueness(t *testing.T) {
	// Generate multiple random byte slices and verify they are different
	const iterations = 100
	const length = 32

	seen := make(map[string]bool)

	for i := 0; i < iterations; i++ {
		bytes, err := GenerateRandomBytes(length)
		if err != nil {
			t.Fatalf("GenerateRandomBytes(%d) unexpected error on iteration %d: %v", length, i, err)
		}

		key := string(bytes)
		if seen[key] {
			t.Errorf("GenerateRandomBytes produced duplicate value on iteration %d", i)
		}
		seen[key] = true
	}
}

// TestGenerateRandomHex tests the GenerateRandomHex function with various input lengths
// including valid lengths, zero, and negative values
func TestGenerateRandomHex(t *testing.T) {
	tests := []struct {
		name           string
		length         int
		expectedHexLen int
		expectError    bool
	}{
		{
			name:           "valid length 16 bytes produces 32 hex chars",
			length:         16,
			expectedHexLen: 32,
			expectError:    false,
		},
		{
			name:           "valid length 1 byte produces 2 hex chars",
			length:         1,
			expectedHexLen: 2,
			expectError:    false,
		},
		{
			name:           "valid length 8 bytes produces 16 hex chars",
			length:         8,
			expectedHexLen: 16,
			expectError:    false,
		},
		{
			name:           "zero length should error",
			length:         0,
			expectedHexLen: 0,
			expectError:    true,
		},
		{
			name:           "negative length should error",
			length:         -5,
			expectedHexLen: 0,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: execute the function
			result, err := GenerateRandomHex(tt.length)

			// Assert: verify error expectation
			if tt.expectError {
				if err == nil {
					t.Errorf("GenerateRandomHex(%d) expected error but got nil", tt.length)
				}
				return
			}

			if err != nil {
				t.Errorf("GenerateRandomHex(%d) unexpected error: %v", tt.length, err)
				return
			}

			// Verify the length of returned hex string
			if len(result) != tt.expectedHexLen {
				t.Errorf("GenerateRandomHex(%d) returned %d chars; expected %d chars",
					tt.length, len(result), tt.expectedHexLen)
			}

			// Verify that result contains only valid hex characters
			for i, c := range result {
				if !isHexChar(c) {
					t.Errorf("GenerateRandomHex(%d) returned invalid hex char '%c' at position %d",
						tt.length, c, i)
				}
			}
		})
	}
}

// isHexChar checks if a rune is a valid hexadecimal character
func isHexChar(c rune) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

// TestGenerateRandomBinaryString tests the GenerateRandomBinaryString function
// with various input lengths including valid lengths, zero, and negative values
func TestGenerateRandomBinaryString(t *testing.T) {
	tests := []struct {
		name              string
		length            int
		expectedBinaryLen int
		expectError       bool
	}{
		{
			name:              "valid length 1 byte produces 8 binary chars",
			length:            1,
			expectedBinaryLen: 8,
			expectError:       false,
		},
		{
			name:              "valid length 2 bytes produces 16 binary chars",
			length:            2,
			expectedBinaryLen: 16,
			expectError:       false,
		},
		{
			name:              "valid length 4 bytes produces 32 binary chars",
			length:            4,
			expectedBinaryLen: 32,
			expectError:       false,
		},
		{
			name:              "zero length should error",
			length:            0,
			expectedBinaryLen: 0,
			expectError:       true,
		},
		{
			name:              "negative length should error",
			length:            -1,
			expectedBinaryLen: 0,
			expectError:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: execute the function
			result, err := GenerateRandomBinaryString(tt.length)

			// Assert: verify error expectation
			if tt.expectError {
				if err == nil {
					t.Errorf("GenerateRandomBinaryString(%d) expected error but got nil", tt.length)
				}
				return
			}

			if err != nil {
				t.Errorf("GenerateRandomBinaryString(%d) unexpected error: %v", tt.length, err)
				return
			}

			// Verify the length of returned binary string
			if len(result) != tt.expectedBinaryLen {
				t.Errorf("GenerateRandomBinaryString(%d) returned %d chars; expected %d chars",
					tt.length, len(result), tt.expectedBinaryLen)
			}

			// Verify that result contains only 0s and 1s
			for i, c := range result {
				if c != '0' && c != '1' {
					t.Errorf("GenerateRandomBinaryString(%d) returned invalid binary char '%c' at position %d",
						tt.length, c, i)
				}
			}
		})
	}
}

// TestGenerateRandomBinaryStringWithBitLength tests the GenerateRandomBinaryStringWithBitLength function
// with various bit lengths including valid lengths, zero, and negative values
func TestGenerateRandomBinaryStringWithBitLength(t *testing.T) {
	tests := []struct {
		name        string
		bitLength   int
		expectError bool
	}{
		{
			name:        "valid bit length 8",
			bitLength:   8,
			expectError: false,
		},
		{
			name:        "valid bit length 16",
			bitLength:   16,
			expectError: false,
		},
		{
			name:        "valid bit length 5 (not multiple of 8)",
			bitLength:   5,
			expectError: false,
		},
		{
			name:        "valid bit length 13 (not multiple of 8)",
			bitLength:   13,
			expectError: false,
		},
		{
			name:        "valid bit length 1",
			bitLength:   1,
			expectError: false,
		},
		{
			name:        "zero bit length should error",
			bitLength:   0,
			expectError: true,
		},
		{
			name:        "negative bit length should error",
			bitLength:   -1,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: execute the function
			result, err := GenerateRandomBinaryStringWithBitLength(tt.bitLength)

			// Assert: verify error expectation
			if tt.expectError {
				if err == nil {
					t.Errorf("GenerateRandomBinaryStringWithBitLength(%d) expected error but got nil", tt.bitLength)
				}
				return
			}

			if err != nil {
				t.Errorf("GenerateRandomBinaryStringWithBitLength(%d) unexpected error: %v", tt.bitLength, err)
				return
			}

			// Verify the exact bit length
			if len(result) != tt.bitLength {
				t.Errorf("GenerateRandomBinaryStringWithBitLength(%d) returned %d bits; expected %d bits",
					tt.bitLength, len(result), tt.bitLength)
			}

			// Verify that result contains only 0s and 1s
			for i, c := range result {
				if c != '0' && c != '1' {
					t.Errorf("GenerateRandomBinaryStringWithBitLength(%d) returned invalid binary char '%c' at position %d",
						tt.bitLength, c, i)
				}
			}
		})
	}
}

