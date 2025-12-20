package cryptors

import (
	"encoding/base64"
	"strings"
	"testing"
)

// TestEncryptAES256 tests the EncryptAES256 function with various inputs
// including valid data, invalid keys, and edge cases
func TestEncryptAES256(t *testing.T) {
	// Generate a valid AES256 key for testing
	validKey, err := GenerateAES256Key()
	if err != nil {
		t.Fatalf("Failed to generate test key: %v", err)
	}

	tests := []struct {
		name        string
		plaintext   []byte
		key         []byte
		expectError bool
	}{
		{
			name:        "valid encryption with text",
			plaintext:   []byte("Hello, World!"),
			key:         validKey,
			expectError: false,
		},
		{
			name:        "valid encryption with long text",
			plaintext:   []byte("This is a longer text to test AES encryption with more data"),
			key:         validKey,
			expectError: false,
		},
		{
			name:        "valid encryption with binary data",
			plaintext:   []byte{0x00, 0x01, 0x02, 0xFF, 0xFE},
			key:         validKey,
			expectError: false,
		},
		{
			name:        "invalid key size - too short",
			plaintext:   []byte("test"),
			key:         []byte("shortkey"),
			expectError: true,
		},
		{
			name:        "invalid key size - too long",
			plaintext:   []byte("test"),
			key:         make([]byte, 64),
			expectError: true,
		},
		{
			name:        "empty plaintext should error",
			plaintext:   []byte{},
			key:         validKey,
			expectError: true,
		},
		{
			name:        "nil plaintext should error",
			plaintext:   nil,
			key:         validKey,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: execute the function
			result, err := EncryptAES256(tt.plaintext, tt.key)

			// Assert: verify error expectation
			if tt.expectError {
				if err == nil {
					t.Errorf("EncryptAES256() expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("EncryptAES256() unexpected error: %v", err)
				return
			}

			// Verify result is not empty
			if result == "" {
				t.Errorf("EncryptAES256() returned empty result")
			}

			// Verify result is valid base64
			_, err = base64.StdEncoding.DecodeString(result)
			if err != nil {
				t.Errorf("EncryptAES256() result is not valid base64: %v", err)
			}
		})
	}
}

// TestDecryptAES256 tests the DecryptAES256 function with various inputs
// including valid encrypted data, invalid keys, and edge cases
func TestDecryptAES256(t *testing.T) {
	// Generate a valid AES256 key for testing
	validKey, err := GenerateAES256Key()
	if err != nil {
		t.Fatalf("Failed to generate test key: %v", err)
	}

	// Encrypt some test data
	plaintext := []byte("Test message for decryption")
	encrypted, err := EncryptAES256(plaintext, validKey)
	if err != nil {
		t.Fatalf("Failed to encrypt test data: %v", err)
	}

	tests := []struct {
		name              string
		encodedCiphertext string
		key               []byte
		expectedPlaintext []byte
		expectError       bool
	}{
		{
			name:              "valid decryption",
			encodedCiphertext: encrypted,
			key:               validKey,
			expectedPlaintext: plaintext,
			expectError:       false,
		},
		{
			name:              "invalid key - wrong key",
			encodedCiphertext: encrypted,
			key:               make([]byte, AES256KeySize),
			expectedPlaintext: nil,
			expectError:       true,
		},
		{
			name:              "invalid key size - too short",
			encodedCiphertext: encrypted,
			key:               []byte("shortkey"),
			expectedPlaintext: nil,
			expectError:       true,
		},
		{
			name:              "empty ciphertext should error",
			encodedCiphertext: "",
			key:               validKey,
			expectedPlaintext: nil,
			expectError:       true,
		},
		{
			name:              "invalid base64",
			encodedCiphertext: "not-valid-base64!!!",
			key:               validKey,
			expectedPlaintext: nil,
			expectError:       true,
		},
		{
			name:              "ciphertext too short",
			encodedCiphertext: base64.StdEncoding.EncodeToString([]byte("short")),
			key:               validKey,
			expectedPlaintext: nil,
			expectError:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: execute the function
			result, err := DecryptAES256(tt.encodedCiphertext, tt.key)

			// Assert: verify error expectation
			if tt.expectError {
				if err == nil {
					t.Errorf("DecryptAES256() expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("DecryptAES256() unexpected error: %v", err)
				return
			}

			// Verify the decrypted plaintext matches expected
			if string(result) != string(tt.expectedPlaintext) {
				t.Errorf("DecryptAES256() = %q; expected %q", string(result), string(tt.expectedPlaintext))
			}
		})
	}
}

// TestEncryptDecryptAES256RoundTrip tests the encryption and decryption process
// to ensure data can be encrypted and then decrypted back to original form
func TestEncryptDecryptAES256RoundTrip(t *testing.T) {
	tests := []struct {
		name      string
		plaintext []byte
	}{
		{
			name:      "simple text",
			plaintext: []byte("Hello, World!"),
		},
		{
			name:      "unicode text",
			plaintext: []byte("こんにちは世界 🌍"),
		},
		{
			name:      "long text",
			plaintext: []byte(strings.Repeat("This is a test message. ", 100)),
		},
		{
			name:      "binary data",
			plaintext: []byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD},
		},
		{
			name:      "special characters",
			plaintext: []byte("!@#$%^&*()_+-=[]{}|;':\",./<>?"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: generate a new key for each test
			key, err := GenerateAES256Key()
			if err != nil {
				t.Fatalf("Failed to generate key: %v", err)
			}

			// Act: encrypt the plaintext
			encrypted, err := EncryptAES256(tt.plaintext, key)
			if err != nil {
				t.Fatalf("EncryptAES256() failed: %v", err)
			}

			// Act: decrypt the ciphertext
			decrypted, err := DecryptAES256(encrypted, key)
			if err != nil {
				t.Fatalf("DecryptAES256() failed: %v", err)
			}

			// Assert: verify decrypted matches original plaintext
			if string(decrypted) != string(tt.plaintext) {
				t.Errorf("Round trip failed: got %q; expected %q", string(decrypted), string(tt.plaintext))
			}
		})
	}
}

// TestEncryptAES256String tests the string convenience function for encryption
func TestEncryptAES256String(t *testing.T) {
	validKey, err := GenerateAES256Key()
	if err != nil {
		t.Fatalf("Failed to generate test key: %v", err)
	}

	tests := []struct {
		name        string
		plaintext   string
		key         []byte
		expectError bool
	}{
		{
			name:        "valid string encryption",
			plaintext:   "Hello, World!",
			key:         validKey,
			expectError: false,
		},
		{
			name:        "empty string should error",
			plaintext:   "",
			key:         validKey,
			expectError: true,
		},
		{
			name:        "unicode string",
			plaintext:   "こんにちは 🌍",
			key:         validKey,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: execute the function
			result, err := EncryptAES256String(tt.plaintext, tt.key)

			// Assert: verify error expectation
			if tt.expectError {
				if err == nil {
					t.Errorf("EncryptAES256String() expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("EncryptAES256String() unexpected error: %v", err)
				return
			}

			// Verify result is not empty
			if result == "" {
				t.Errorf("EncryptAES256String() returned empty result")
			}
		})
	}
}

// TestDecryptAES256String tests the string convenience function for decryption
func TestDecryptAES256String(t *testing.T) {
	validKey, err := GenerateAES256Key()
	if err != nil {
		t.Fatalf("Failed to generate test key: %v", err)
	}

	plaintext := "Test message"
	encrypted, err := EncryptAES256String(plaintext, validKey)
	if err != nil {
		t.Fatalf("Failed to encrypt test data: %v", err)
	}

	tests := []struct {
		name              string
		encodedCiphertext string
		key               []byte
		expectedPlaintext string
		expectError       bool
	}{
		{
			name:              "valid string decryption",
			encodedCiphertext: encrypted,
			key:               validKey,
			expectedPlaintext: plaintext,
			expectError:       false,
		},
		{
			name:              "invalid ciphertext",
			encodedCiphertext: "invalid",
			key:               validKey,
			expectedPlaintext: "",
			expectError:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: execute the function
			result, err := DecryptAES256String(tt.encodedCiphertext, tt.key)

			// Assert: verify error expectation
			if tt.expectError {
				if err == nil {
					t.Errorf("DecryptAES256String() expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("DecryptAES256String() unexpected error: %v", err)
				return
			}

			// Verify the decrypted plaintext matches expected
			if result != tt.expectedPlaintext {
				t.Errorf("DecryptAES256String() = %q; expected %q", result, tt.expectedPlaintext)
			}
		})
	}
}

// TestGenerateAES256Key tests the AES256 key generation function
func TestGenerateAES256Key(t *testing.T) {
	// Act: generate multiple keys
	key1, err := GenerateAES256Key()
	if err != nil {
		t.Fatalf("GenerateAES256Key() failed: %v", err)
	}

	key2, err := GenerateAES256Key()
	if err != nil {
		t.Fatalf("GenerateAES256Key() failed on second call: %v", err)
	}

	// Assert: verify key length
	if len(key1) != AES256KeySize {
		t.Errorf("GenerateAES256Key() returned %d bytes; expected %d bytes", len(key1), AES256KeySize)
	}

	// Assert: verify keys are different (randomness check)
	if string(key1) == string(key2) {
		t.Errorf("GenerateAES256Key() produced duplicate keys")
	}
}

// TestGenerateAES256KeyBase64 tests the base64 key generation function
func TestGenerateAES256KeyBase64(t *testing.T) {
	// Act: generate key
	encodedKey, err := GenerateAES256KeyBase64()
	if err != nil {
		t.Fatalf("GenerateAES256KeyBase64() failed: %v", err)
	}

	// Assert: verify result is not empty
	if encodedKey == "" {
		t.Errorf("GenerateAES256KeyBase64() returned empty string")
	}

	// Assert: verify result is valid base64
	decodedKey, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		t.Errorf("GenerateAES256KeyBase64() result is not valid base64: %v", err)
	}

	// Assert: verify decoded key has correct length
	if len(decodedKey) != AES256KeySize {
		t.Errorf("GenerateAES256KeyBase64() decoded key has %d bytes; expected %d bytes",
			len(decodedKey), AES256KeySize)
	}
}

// TestDecodeAES256Key tests the key decoding and validation function
func TestDecodeAES256Key(t *testing.T) {
	// Generate a valid key for testing
	validKey, err := GenerateAES256Key()
	if err != nil {
		t.Fatalf("Failed to generate test key: %v", err)
	}
	validEncodedKey := base64.StdEncoding.EncodeToString(validKey)

	tests := []struct {
		name        string
		encodedKey  string
		expectError bool
	}{
		{
			name:        "valid base64 encoded key",
			encodedKey:  validEncodedKey,
			expectError: false,
		},
		{
			name:        "invalid base64",
			encodedKey:  "not-valid-base64!!!",
			expectError: true,
		},
		{
			name:        "valid base64 but wrong key size",
			encodedKey:  base64.StdEncoding.EncodeToString([]byte("short")),
			expectError: true,
		},
		{
			name:        "empty string should error",
			encodedKey:  "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: execute the function
			result, err := DecodeAES256Key(tt.encodedKey)

			// Assert: verify error expectation
			if tt.expectError {
				if err == nil {
					t.Errorf("DecodeAES256Key() expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("DecodeAES256Key() unexpected error: %v", err)
				return
			}

			// Verify the decoded key has correct length
			if len(result) != AES256KeySize {
				t.Errorf("DecodeAES256Key() returned %d bytes; expected %d bytes",
					len(result), AES256KeySize)
			}
		})
	}
}

// TestAES256DifferentNonces tests that encrypting the same plaintext twice
// produces different ciphertexts due to different nonces
func TestAES256DifferentNonces(t *testing.T) {
	// Arrange: generate key and plaintext
	key, err := GenerateAES256Key()
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	plaintext := []byte("Same message encrypted twice")

	// Act: encrypt the same plaintext twice
	encrypted1, err := EncryptAES256(plaintext, key)
	if err != nil {
		t.Fatalf("First encryption failed: %v", err)
	}

	encrypted2, err := EncryptAES256(plaintext, key)
	if err != nil {
		t.Fatalf("Second encryption failed: %v", err)
	}

	// Assert: verify the ciphertexts are different
	if encrypted1 == encrypted2 {
		t.Errorf("Encrypting same plaintext twice produced identical ciphertext")
	}

	// Assert: verify both decrypt to the same plaintext
	decrypted1, err := DecryptAES256(encrypted1, key)
	if err != nil {
		t.Fatalf("First decryption failed: %v", err)
	}

	decrypted2, err := DecryptAES256(encrypted2, key)
	if err != nil {
		t.Fatalf("Second decryption failed: %v", err)
	}

	if string(decrypted1) != string(plaintext) || string(decrypted2) != string(plaintext) {
		t.Errorf("Decrypted text does not match original plaintext")
	}
}

