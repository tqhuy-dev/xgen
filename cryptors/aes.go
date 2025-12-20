package cryptors

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

const (
	// AES256KeySize is the required key size for AES-256 encryption (32 bytes)
	AES256KeySize = 32
	// NonceSize is the standard nonce size for GCM mode (12 bytes)
	NonceSize = 12
)

// EncryptAES256 encrypts plaintext using AES-256-GCM.
// The key must be exactly 32 bytes (256 bits).
// Returns the encrypted data as base64 encoded string, or an error if encryption fails.
// The nonce is prepended to the ciphertext for decryption.
func EncryptAES256(plaintext []byte, key []byte) (string, error) {
	if len(key) != AES256KeySize {
		return "", fmt.Errorf("invalid key size: expected %d bytes, got %d bytes", AES256KeySize, len(key))
	}

	if len(plaintext) == 0 {
		return "", fmt.Errorf("plaintext cannot be empty")
	}

	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate a random nonce
	nonce, err := GenerateRandomBytes(NonceSize)
	if err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt the plaintext
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	// Encode to base64 for safe storage/transmission
	encoded := base64.StdEncoding.EncodeToString(ciphertext)

	return encoded, nil
}

// DecryptAES256 decrypts a base64 encoded ciphertext using AES-256-GCM.
// The key must be exactly 32 bytes (256 bits).
// Returns the decrypted plaintext, or an error if decryption fails.
func DecryptAES256(encodedCiphertext string, key []byte) ([]byte, error) {
	if len(key) != AES256KeySize {
		return nil, fmt.Errorf("invalid key size: expected %d bytes, got %d bytes", AES256KeySize, len(key))
	}

	if encodedCiphertext == "" {
		return nil, fmt.Errorf("ciphertext cannot be empty")
	}

	// Decode from base64
	ciphertext, err := base64.StdEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Verify minimum ciphertext length
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short: minimum length is %d bytes", nonceSize)
	}

	// Extract nonce and ciphertext
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt the ciphertext
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// EncryptAES256String is a convenience function that encrypts a string.
// Returns the encrypted data as base64 encoded string, or an error if encryption fails.
func EncryptAES256String(plaintext string, key []byte) (string, error) {
	return EncryptAES256([]byte(plaintext), key)
}

// DecryptAES256String is a convenience function that decrypts to a string.
// Returns the decrypted plaintext as string, or an error if decryption fails.
func DecryptAES256String(encodedCiphertext string, key []byte) (string, error) {
	plaintext, err := DecryptAES256(encodedCiphertext, key)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// GenerateAES256Key generates a cryptographically secure random 32-byte key for AES-256.
// Returns the key bytes or an error if generation fails.
func GenerateAES256Key() ([]byte, error) {
	return GenerateRandomBytes(AES256KeySize)
}

// GenerateAES256KeyBase64 generates a cryptographically secure random 32-byte key for AES-256.
// Returns the key as base64 encoded string for easy storage, or an error if generation fails.
func GenerateAES256KeyBase64() (string, error) {
	key, err := GenerateAES256Key()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

// DecodeAES256Key decodes a base64 encoded key and validates it's the correct size.
// Returns the decoded key bytes or an error if decoding or validation fails.
func DecodeAES256Key(encodedKey string) ([]byte, error) {
	key, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 key: %w", err)
	}

	if len(key) != AES256KeySize {
		return nil, fmt.Errorf("invalid key size: expected %d bytes, got %d bytes", AES256KeySize, len(key))
	}

	return key, nil
}
