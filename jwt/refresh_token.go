package jwt

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRefreshToken() (string, error) {
	return generateRandomString(32)
}

// GenerateRandomString generates a random string of length n
func generateRandomString(n int) (string, error) {

	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err

	}
	return hex.EncodeToString(b), nil
}
