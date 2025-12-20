package cryptors

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// GenerateRandomBytes generates cryptographically secure random bytes with the specified length.
// It uses crypto/rand which is suitable for security-sensitive operations.
// Returns the random bytes or an error if the random generation fails.
func GenerateRandomBytes(length int) ([]byte, error) {
	if length <= 0 {
		return nil, fmt.Errorf("length must be greater than 0, got %d", length)
	}

	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return bytes, nil
}

// GenerateRandomHex generates a cryptographically secure random hex string.
// The resulting hex string will have length * 2 characters (each byte = 2 hex chars).
// Returns the hex string or an error if the random generation fails.
func GenerateRandomHex(length int) (string, error) {
	bytes, err := GenerateRandomBytes(length)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

// GenerateRandomBinaryString generates a cryptographically secure random binary string.
// Each byte is represented as 8 binary digits (0s and 1s).
// The resulting binary string will have length * 8 characters.
// Returns the binary string or an error if the random generation fails.
func GenerateRandomBinaryString(length int) (string, error) {
	bytes, err := GenerateRandomBytes(length)
	if err != nil {
		return "", err
	}

	binary := ""
	for _, b := range bytes {
		binary += fmt.Sprintf("%08b", b)
	}

	return binary, nil
}

// GenerateRandomBinaryStringWithBitLength generates a cryptographically secure random binary string
// with exactly the specified number of bits.
// Returns the binary string or an error if the random generation fails.
func GenerateRandomBinaryStringWithBitLength(bitLength int) (string, error) {
	if bitLength <= 0 {
		return "", fmt.Errorf("bitLength must be greater than 0, got %d", bitLength)
	}

	// Calculate the number of bytes needed
	byteLength := (bitLength + 7) / 8

	bytes, err := GenerateRandomBytes(byteLength)
	if err != nil {
		return "", err
	}

	binary := ""
	for _, b := range bytes {
		binary += fmt.Sprintf("%08b", b)
	}

	// Trim to exact bit length
	return binary[:bitLength], nil
}

