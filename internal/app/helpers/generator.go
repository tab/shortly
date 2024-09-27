package helpers

import (
	"crypto/rand"
)

func GenerateShortCode() string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "randomID"
	}
	for i, b := range bytes {
		bytes[i] = chars[b%byte(len(chars))]
	}
	return string(bytes)
}
