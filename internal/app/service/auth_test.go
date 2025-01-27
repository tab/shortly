package service

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"shortly/internal/app/config"
	"shortly/internal/app/errors"
)

func Test_NewJWTService(t *testing.T) {
	cfg := &config.Config{
		SecretKey: "jwt-secret-key",
	}
	service := NewAuthService(cfg)

	assert.NotNil(t, service)
}

func Test_JWTService_Generate(t *testing.T) {
	cfg := &config.Config{
		SecretKey: "jwt-secret-key",
	}
	service := NewAuthService(cfg)

	UUID, _ := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")

	type result struct {
		header string
	}

	tests := []struct {
		name     string
		id       uuid.UUID
		expected result
	}{
		{
			name: "Success",
			id:   UUID,
			expected: result{
				header: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			},
		},
		{
			name: "Empty id",
			id:   uuid.UUID{},
			expected: result{
				header: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			},
		},
		{
			name: "Nil id",
			id:   uuid.Nil,
			expected: result{
				header: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := service.Generate(tt.id)
			assert.NoError(t, err)
			assert.NotEmpty(t, token)
			assert.Equal(t, tt.expected.header, token[:36])
		})
	}
}

func Test_JWTService_Verify(t *testing.T) {
	cfg := &config.Config{
		SecretKey: "jwt-secret-key",
	}
	service := NewAuthService(cfg)

	UUID, _ := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")

	validToken, err := service.Generate(UUID)
	assert.NoError(t, err)

	wrongSecretService := NewAuthService(&config.Config{
		SecretKey: "some-other-key",
	})
	wrongToken, err := wrongSecretService.Generate(UUID)
	assert.NoError(t, err)

	invalidToken := validToken
	invalidToken = invalidToken[:len(invalidToken)-5] + "XYZ"

	tests := []struct {
		name     string
		token    string
		expected uuid.UUID
		error    error
	}{
		{
			name:     "Success",
			token:    validToken,
			expected: UUID,
			error:    nil,
		},
		{
			name:     "Wrong token",
			token:    wrongToken,
			expected: uuid.Nil,
			error:    jwt.ErrSignatureInvalid,
		},
		{
			name:     "Invalid token",
			token:    invalidToken,
			expected: uuid.Nil,
			error:    errors.ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := service.Verify(tt.token)

			if tt.error != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, id)
		})
	}
}
