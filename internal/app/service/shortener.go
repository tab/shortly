package service

import (
	"context"
	"fmt"

	"shortly/internal/app/config"
	"shortly/internal/app/dto"
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

	record, err := s.repo.CreateURL(ctx, url)
	if err != nil {
		return "", errors.ErrFailedToSaveURL
	}

	if record.ShortCode != shortCode {
		return fmt.Sprintf("%s/%s", s.cfg.BaseURL, record.ShortCode), errors.ErrURLAlreadyExists
	}

	return fmt.Sprintf("%s/%s", s.cfg.BaseURL, shortCode), nil
}

func (s *URLService) CreateShortLinks(ctx context.Context, params []dto.BatchCreateShortLinkParams) ([]dto.BatchCreateShortLinkResponse, error) {
	longURLs := make([]repository.URL, 0, len(params))
	results := make([]dto.BatchCreateShortLinkResponse, 0, len(params))

	for _, param := range params {
		uuid, err := s.rand.UUID()
		if err != nil {
			return nil, errors.ErrFailedToGenerateUUID
		}

		shortCode, err := s.rand.Hex()
		if err != nil {
			return nil, errors.ErrFailedToGenerateCode
		}

		url := repository.URL{
			UUID:      uuid,
			LongURL:   param.OriginalURL,
			ShortCode: shortCode,
		}
		longURLs = append(longURLs, url)

		results = append(results, dto.BatchCreateShortLinkResponse{
			CorrelationID: param.CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", s.cfg.BaseURL, shortCode),
		})
	}

	if err := s.repo.CreateURLs(ctx, longURLs); err != nil {
		return nil, errors.ErrFailedToSaveURL
	}

	return results, nil
}

func (s *URLService) GetShortLink(ctx context.Context, shortCode string) (*repository.URL, bool) {
	return s.repo.GetURLByShortCode(ctx, shortCode)
}
