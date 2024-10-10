package service

import (
	"crypto/rand"
	"encoding/hex"

	"shortly/internal/app/errors"
)

const BytesLength = 4 // 4 bytes = 8 characters

type SecureRandomGenerator interface {
	Hex() (string, error)
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
		return "", errors.ErrorFailedToReadRandomBytes
	}

	return hex.EncodeToString(bytes), nil
}
