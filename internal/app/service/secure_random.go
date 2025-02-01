package service

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/google/uuid"

	"shortly/internal/app/errors"
)

// BytesLength is a length short code: 4 bytes = 8 characters
const BytesLength = 4

// SecureRandomGenerator is an interface for secure random operations
type SecureRandomGenerator interface {
	Hex() (string, error)
	UUID() (uuid.UUID, error)
}

// SecureRandom is a service for secure random operations
type SecureRandom struct{}

// NewSecureRandom creates a new secure random service instance
func NewSecureRandom() *SecureRandom {
	return &SecureRandom{}
}

// Read reads random bytes
func (random *SecureRandom) Read(bytes []byte) (int, error) {
	return rand.Read(bytes)
}

// Hex generates a new hex string
func (random *SecureRandom) Hex() (string, error) {
	bytes := make([]byte, BytesLength)

	_, err := random.Read(bytes)
	if err != nil {
		return "", errors.ErrFailedToReadRandomBytes
	}

	return hex.EncodeToString(bytes), nil
}

// UUID generates a new UUID
func (random *SecureRandom) UUID() (uuid.UUID, error) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return uuid.UUID{}, errors.ErrFailedToGenerateUUID
	}

	return newUUID, nil
}
