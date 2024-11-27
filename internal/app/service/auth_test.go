package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"shortly/internal/app/config"
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

func Test_JWTService_Decode(t *testing.T) {
	cfg := &config.Config{
		SecretKey: "jwt-secret-key",
	}
	service := NewAuthService(cfg)

	UUID, _ := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")

	tests := []struct {
		name     string
		id       uuid.UUID
		expected string
	}{
		{
			name:     "Success",
			id:       UUID,
			expected: "123e4567-e89b-12d3-a456-426614174000",
		},
		{
			name:     "Empty id",
			id:       uuid.UUID{},
			expected: "00000000-0000-0000-0000-000000000000",
		},
		{
			name:     "Nil id",
			id:       uuid.Nil,
			expected: "00000000-0000-0000-0000-000000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := service.Generate(tt.id)
			assert.NoError(t, err)

			id, err := service.Verify(token)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, id.String())
		})
	}
}
