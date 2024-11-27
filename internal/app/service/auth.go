package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"shortly/internal/app/config"
	"shortly/internal/app/errors"
)

const TokenExp = time.Hour * 24

type Authenticator interface {
	Generate(id uuid.UUID) (string, error)
	Verify(token string) (uuid.UUID, error)
}

type authService struct {
	cfg *config.Config
}

type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
}

func NewAuthService(cfg *config.Config) Authenticator {
	return &authService{cfg: cfg}
}

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
