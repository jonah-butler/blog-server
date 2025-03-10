package user

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

func generateToken() (string, string, error) {
	// Generate 32 random bytes (256 bits)
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", "", err
	}

	// Encode the random bytes to a URL-safe base64 string (without padding)
	token := base64.RawURLEncoding.EncodeToString(tokenBytes)

	// Compute the SHA-256 hash of the token for storage
	hash := sha256.Sum256([]byte(token))
	tokenHash := base64.RawURLEncoding.EncodeToString(hash[:])

	return token, tokenHash, nil
}
