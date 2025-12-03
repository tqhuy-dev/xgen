package mongo_db

import (
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// TestObjectIDToString tests converting ObjectID to string
func TestObjectIDToString(t *testing.T) {
	tests := []struct {
		name     string
		objectID bson.ObjectID
		expected string
	}{
		{
			name:     "valid ObjectID",
			objectID: bson.NewObjectID(),
			expected: "", // Will check length instead
		},
		{
			name:     "zero ObjectID",
			objectID: bson.ObjectID{},
			expected: "000000000000000000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := ObjectIDToString(tt.objectID)

			// Assert
			if tt.expected != "" {
				if result != tt.expected {
					t.Errorf("ObjectIDToString() = %v; expected %v", result, tt.expected)
				}
			} else {
				// For generated ObjectIDs, just check the length
				if len(result) != 24 {
					t.Errorf("ObjectIDToString() length = %d; expected 24", len(result))
				}
			}
		})
	}
}

// TestStringToObjectID tests converting string to ObjectID
func TestStringToObjectID(t *testing.T) {
	validObjectID := bson.NewObjectID()
	validHex := validObjectID.Hex()

	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "valid ObjectID hex string",
			input:       validHex,
			expectError: false,
		},
		{
			name:        "valid zero ObjectID",
			input:       "000000000000000000000000",
			expectError: false,
		},
		{
			name:        "invalid hex string - too short",
			input:       "507f1f77bcf86cd79943901",
			expectError: true,
		},
		{
			name:        "invalid hex string - too long",
			input:       "507f1f77bcf86cd7994390112",
			expectError: true,
		},
		{
			name:        "invalid hex string - non-hex characters",
			input:       "507f1f77bcf86cd79943901g",
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expectError: true,
		},
		{
			name:        "invalid hex string - special characters",
			input:       "507f1f77bcf86cd79943901!",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := StringToObjectID(tt.input)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("StringToObjectID(%q) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("StringToObjectID(%q) unexpected error: %v", tt.input, err)
				}
				// Verify we can convert back
				if result.Hex() != tt.input {
					t.Errorf("StringToObjectID(%q).Hex() = %v; expected %v", tt.input, result.Hex(), tt.input)
				}
			}
		})
	}
}

// TestAnyToObjectID tests converting any type to ObjectID
func TestAnyToObjectID(t *testing.T) {
	validObjectID := bson.NewObjectID()
	validHex := validObjectID.Hex()
	validBytes := []byte{0x50, 0x7f, 0x1f, 0x77, 0xbc, 0xf8, 0x6c, 0xd7, 0x99, 0x43, 0x90, 0x11}
	validByteArray := [12]byte{0x50, 0x7f, 0x1f, 0x77, 0xbc, 0xf8, 0x6c, 0xd7, 0x99, 0x43, 0x90, 0x11}

	tests := []struct {
		name        string
		input       any
		expectError bool
		checkValue  bool
		expected    bson.ObjectID
	}{
		{
			name:        "ObjectID input",
			input:       validObjectID,
			expectError: false,
			checkValue:  true,
			expected:    validObjectID,
		},
		{
			name:        "string hex input",
			input:       validHex,
			expectError: false,
			checkValue:  true,
			expected:    validObjectID,
		},
		{
			name:        "byte slice input - 12 bytes",
			input:       validBytes,
			expectError: false,
			checkValue:  false,
		},
		{
			name:        "byte array input - [12]byte",
			input:       validByteArray,
			expectError: false,
			checkValue:  false,
		},
		{
			name:        "nil input",
			input:       nil,
			expectError: true,
		},
		{
			name:        "invalid string - not hex",
			input:       "invalid-hex-string",
			expectError: true,
		},
		{
			name:        "invalid byte slice - wrong length",
			input:       []byte{0x50, 0x7f, 0x1f},
			expectError: true,
		},
		{
			name:        "unsupported type - int",
			input:       12345,
			expectError: true,
		},
		{
			name:        "unsupported type - bool",
			input:       true,
			expectError: true,
		},
		{
			name:        "unsupported type - struct",
			input:       struct{ ID string }{ID: "test"},
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
			result, err := AnyToObjectID(tt.input)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("AnyToObjectID(%v) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("AnyToObjectID(%v) unexpected error: %v", tt.input, err)
				}
				if tt.checkValue && result != tt.expected {
					t.Errorf("AnyToObjectID(%v) = %v; expected %v", tt.input, result, tt.expected)
				}
				// Verify result is not zero ObjectID (unless input was zero)
				if !tt.checkValue && result.IsZero() && tt.input != nil {
					zeroBytes := make([]byte, 12)
					if _, ok := tt.input.([]byte); ok && string(tt.input.([]byte)) != string(zeroBytes) {
						t.Errorf("AnyToObjectID(%v) returned zero ObjectID", tt.input)
					}
				}
			}
		})
	}
}

// TestMustStringToObjectID tests the Must variant of StringToObjectID
func TestMustStringToObjectID(t *testing.T) {
	validObjectID := bson.NewObjectID()
	validHex := validObjectID.Hex()

	tests := []struct {
		name        string
		input       string
		shouldPanic bool
	}{
		{
			name:        "valid hex string - no panic",
			input:       validHex,
			shouldPanic: false,
		},
		{
			name:        "invalid hex string - should panic",
			input:       "invalid",
			shouldPanic: true,
		},
		{
			name:        "empty string - should panic",
			input:       "",
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act & Assert
			defer func() {
				r := recover()
				if tt.shouldPanic && r == nil {
					t.Errorf("MustStringToObjectID(%q) expected panic but got none", tt.input)
				}
				if !tt.shouldPanic && r != nil {
					t.Errorf("MustStringToObjectID(%q) unexpected panic: %v", tt.input, r)
				}
			}()

			result := MustStringToObjectID(tt.input)
			if !tt.shouldPanic && result.Hex() != tt.input {
				t.Errorf("MustStringToObjectID(%q) = %v; expected %v", tt.input, result.Hex(), tt.input)
			}
		})
	}
}

// TestMustAnyToObjectID tests the Must variant of AnyToObjectID
func TestMustAnyToObjectID(t *testing.T) {
	validObjectID := bson.NewObjectID()
	validHex := validObjectID.Hex()

	tests := []struct {
		name        string
		input       any
		shouldPanic bool
	}{
		{
			name:        "valid ObjectID - no panic",
			input:       validObjectID,
			shouldPanic: false,
		},
		{
			name:        "valid hex string - no panic",
			input:       validHex,
			shouldPanic: false,
		},
		{
			name:        "nil input - should panic",
			input:       nil,
			shouldPanic: true,
		},
		{
			name:        "invalid type - should panic",
			input:       12345,
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act & Assert
			defer func() {
				r := recover()
				if tt.shouldPanic && r == nil {
					t.Errorf("MustAnyToObjectID(%v) expected panic but got none", tt.input)
				}
				if !tt.shouldPanic && r != nil {
					t.Errorf("MustAnyToObjectID(%v) unexpected panic: %v", tt.input, r)
				}
			}()

			_ = MustAnyToObjectID(tt.input)
		})
	}
}

// TestIsValidObjectIDHex tests checking if a string is a valid ObjectID hex
func TestIsValidObjectIDHex(t *testing.T) {
	validObjectID := bson.NewObjectID()
	validHex := validObjectID.Hex()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid ObjectID hex",
			input:    validHex,
			expected: true,
		},
		{
			name:     "valid zero ObjectID",
			input:    "000000000000000000000000",
			expected: true,
		},
		{
			name:     "invalid - too short",
			input:    "507f1f77bcf86cd79943901",
			expected: false,
		},
		{
			name:     "invalid - too long",
			input:    "507f1f77bcf86cd7994390112",
			expected: false,
		},
		{
			name:     "invalid - non-hex characters",
			input:    "507f1f77bcf86cd79943901g",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "invalid - special characters",
			input:    "507f1f77bcf86cd79943901!",
			expected: false,
		},
		{
			name:     "invalid - spaces",
			input:    "507f1f77 bcf86cd799439011",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := IsValidObjectIDHex(tt.input)

			// Assert
			if result != tt.expected {
				t.Errorf("IsValidObjectIDHex(%q) = %v; expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestRoundTrip tests converting ObjectID -> String -> ObjectID
func TestRoundTrip(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "round trip conversion",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: create original ObjectID
			original := bson.NewObjectID()

			// Act: convert to string
			str := ObjectIDToString(original)

			// Act: convert back to ObjectID
			result, err := StringToObjectID(str)

			// Assert: should match original
			if err != nil {
				t.Errorf("Round trip conversion failed: %v", err)
			}
			if result != original {
				t.Errorf("Round trip conversion: got %v; expected %v", result, original)
			}
		})
	}
}

// TestAnyToObjectID_EdgeCases tests edge cases for AnyToObjectID
func TestAnyToObjectID_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       any
		expectError bool
		description string
	}{
		{
			name:        "zero ObjectID",
			input:       bson.ObjectID{},
			expectError: false,
			description: "should accept zero ObjectID",
		},
		{
			name:        "zero ObjectID string",
			input:       "000000000000000000000000",
			expectError: false,
			description: "should accept zero ObjectID hex string",
		},
		{
			name:        "zero byte slice",
			input:       make([]byte, 12),
			expectError: false,
			description: "should accept 12 zero bytes",
		},
		{
			name:        "pointer to ObjectID",
			input:       &bson.ObjectID{},
			expectError: true,
			description: "should not accept pointer to ObjectID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			_, err := AnyToObjectID(tt.input)

			// Assert
			if tt.expectError && err == nil {
				t.Errorf("AnyToObjectID(%v) expected error but got none: %s", tt.input, tt.description)
			}
			if !tt.expectError && err != nil {
				t.Errorf("AnyToObjectID(%v) unexpected error: %v (%s)", tt.input, err, tt.description)
			}
		})
	}
}

