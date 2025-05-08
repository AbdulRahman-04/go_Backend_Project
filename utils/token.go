package utils

import (
	"math/rand"
	"time"
)

// GenerateRandomToken returns a random alphanumeric string.
func GenerateRandomToken() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const tokenLength = 16

	rand.Seed(time.Now().UnixNano()) // Ensure randomness
	tokenBytes := make([]byte, tokenLength)
	for i := range tokenBytes {
		tokenBytes[i] = charset[rand.Intn(len(charset))]
	}
	return string(tokenBytes)
}