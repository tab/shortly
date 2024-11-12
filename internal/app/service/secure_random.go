package service

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/google/uuid"

	"shortly/internal/app/errors"
)

const BytesLength = 4 // 4 bytes = 8 characters

type SecureRandomGenerator interface {
	Hex() (string, error)
	UUID() (uuid.UUID, error)
}

type SecureRandom struct{}

func NewSecureRandom() *SecureRandom {
	return &SecureRandom{}
}

func (random *SecureRandom) Read(bytes []byte) (int, error) {
	return rand.Read(bytes)
}

func (random *SecureRandom) Hex() (string, error) {
	bytes := make([]byte, BytesLength)

	_, err := random.Read(bytes)
	if err != nil {
		return "", errors.ErrFailedToReadRandomBytes
	}

	return hex.EncodeToString(bytes), nil
}

func (random *SecureRandom) UUID() (uuid.UUID, error) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return uuid.UUID{}, errors.ErrFailedToGenerateUUID
	}

	return newUUID, nil
}
