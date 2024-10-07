package helpers

import (
	"crypto/rand"
	"errors"
)

const ShortCodeLength = 8

type SecureRandom interface {
	Code() (string, error)
}

type Generator struct{}

var generator SecureRandom = Generator{}

func (Generator) Code() (string, error) { // Changed method name to Code
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, ShortCodeLength)

	if _, err := rand.Read(bytes); err != nil {
		return "", errors.New("failed to generate short code")
	}

	for i, b := range bytes {
		bytes[i] = chars[b%byte(len(chars))]
	}

	return string(bytes), nil
}

func SetShortCodeGenerator(gen SecureRandom) {
	generator = gen
}

func ShortCode() (string, error) {
	return generator.Code()
}
