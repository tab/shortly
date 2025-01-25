package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"shortly/internal/app/config"
	"shortly/internal/app/errors"
)

// TokenExp is a token expiration time
const TokenExp = time.Hour * 24

// Authenticator is an interface for authentication operations
type Authenticator interface {
	Generate(id uuid.UUID) (string, error)
	Verify(token string) (uuid.UUID, error)
}

// authService is a service for authentication operations
type authService struct {
	cfg *config.Config
}

// Claims is a type for JWT claims
type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
}

// NewAuthService creates a new authentication service instance
func NewAuthService(cfg *config.Config) Authenticator {
	return &authService{cfg: cfg}
}

// Generate generates a new JWT token
func (s *authService) Generate(id uuid.UUID) (string, error) {
	result := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: id,
	})

	token, err := result.SignedString([]byte(s.cfg.SecretKey))
	if err != nil {
		return "", err
	}

	return token, nil
}

// Verify verifies a JWT token
func (s *authService) Verify(token string) (uuid.UUID, error) {
	claims := &Claims{}

	result, err := jwt.ParseWithClaims(token, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.ErrInvalidSigningMethod
			}
			return []byte(s.cfg.SecretKey), nil
		})

	if err != nil {
		return uuid.Nil, err
	}

	if !result.Valid {
		return uuid.Nil, errors.ErrInvalidToken
	}

	return claims.UserID, nil
}
