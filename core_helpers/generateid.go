package core_helpers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"

	"github.com/google/uuid"
)

// GenerateCondensedUUID generates a condensed UUID by encoding a standard UUID to Base64 format and removing padding.
func GenerateCondensedUUID() (string, error) {
	// Generate a new UUID
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	// Encode the UUID as a URL-safe Base64 string
	encodedUUID := base64.URLEncoding.EncodeToString(newUUID[:])

	// Remove padding characters (like '=')
	condensedUUID := strings.TrimRight(encodedUUID, "=")

	return condensedUUID, nil
}

// GenerateNumericID generates a random, URL-safe, compact numeric ID.
func GenerateNumericID() (string, error) {
	// Define the maximum value for the numeric ID (e.g., 10^12 for a 12-digit ID)
	max := big.NewInt(1e8) // Adjust as needed for desired numeric range

	// Generate a random number in the range [0, max)
	randomNum, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", fmt.Errorf("failed to generate random number: %w", err)
	}

	// Convert the random number to a string
	return randomNum.String(), nil
}

// for use when creating directories
// GenerateDirectoryID generates a random, URL-safe, compact numeric ID.
func GenerateDirectoryID() (string, error) {
	// Define the maximum value for the numeric ID (e.g., 10^12 for a 12-digit ID)
	max := big.NewInt(1e6) // Adjust as needed for desired numeric range

	// Generate a random number in the range [0, max)
	randomNum, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", fmt.Errorf("failed to generate random number: %w", err)
	}

	// Convert the random number to a string
	return randomNum.String(), nil
}

// used when creating an organisation
// GenerateAlphanumericID generates a random, URL-safe alphanumeric ID of the given length.
func GenerateAlphanumericID(length int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	if length <= 0 {
		return "", fmt.Errorf("length must be greater than 0")
	}

	// Create a byte slice to hold the random characters
	randomBytes := make([]byte, length)

	// Generate random characters
	for i := range randomBytes {
		// Randomly select a character from the charset
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random character: %w", err)
		}
		randomBytes[i] = charset[num.Int64()]
	}

	// Convert the byte slice to a string and return
	return string(randomBytes), nil
}
