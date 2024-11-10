package service

import (
	"context"
	"fmt"

	"shortly/internal/app/config"
	"shortly/internal/app/errors"
	"shortly/internal/app/repository"
)

type URLService struct {
	cfg  *config.Config
	repo repository.Repository
	rand SecureRandomGenerator
}

func NewURLService(cfg *config.Config, repo repository.Repository, rand SecureRandomGenerator) *URLService {
	return &URLService{
		cfg:  cfg,
		repo: repo,
		rand: rand,
	}
}

func (s *URLService) CreateShortLink(ctx context.Context, longURL string) (string, error) {
	uuid, err := s.rand.UUID()
	if err != nil {
		return "", errors.ErrFailedToGenerateUUID
	}

	shortCode, err := s.rand.Hex()
	if err != nil {
		return "", errors.ErrFailedToGenerateCode
	}

	url := repository.URL{
		UUID:      uuid,
		LongURL:   longURL,
		ShortCode: shortCode,
	}
	err = s.repo.CreateURL(ctx, url)
	if err != nil {
		return "", errors.ErrFailedToSaveURL
	}

	return fmt.Sprintf("%s/%s", s.cfg.BaseURL, shortCode), nil
}

func (s *URLService) GetShortLink(ctx context.Context, shortCode string) (*repository.URL, bool) {
	return s.repo.GetURLByShortCode(ctx, shortCode)
}
