package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

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
	id, err := s.rand.UUID()
	if err != nil {
		return "", errors.ErrFailedToGenerateUUID
	}

	shortCode, err := s.generateUniqueShortCode(ctx)
	if err != nil {
		return "", errors.ErrFailedToGenerateCode
	}

	currentUserID, ok := (ctx.Value(dto.CurrentUser)).(uuid.UUID)
	if !ok {
		currentUserID = uuid.Nil
	}

	url := repository.URL{
		UUID:      id,
		LongURL:   longURL,
		ShortCode: shortCode,
		UserUUID:  currentUserID,
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

	currentUserID, ok := (ctx.Value(dto.CurrentUser)).(uuid.UUID)
	if !ok {
		currentUserID = uuid.Nil
	}

	for _, param := range params {
		id, err := s.rand.UUID()
		if err != nil {
			return nil, errors.ErrFailedToGenerateUUID
		}

		shortCode, err := s.generateUniqueShortCode(ctx)
		if err != nil {
			return nil, errors.ErrFailedToGenerateCode
		}

		url := repository.URL{
			UUID:      id,
			LongURL:   param.OriginalURL,
			ShortCode: shortCode,
			UserUUID:  currentUserID,
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

func (s *URLService) generateUniqueShortCode(ctx context.Context) (string, error) {
	for {
		shortCode, err := s.rand.Hex()
		if err != nil {
			return "", err
		}

		if _, exists := s.repo.GetURLByShortCode(ctx, shortCode); !exists {
			return shortCode, nil
		}
	}
}

func (s *URLService) GetUserURLs(ctx context.Context) ([]dto.GetUserURLsResponse, error) {
	currentUserID, ok := (ctx.Value(dto.CurrentUser)).(uuid.UUID)
	if !ok {
		return nil, errors.ErrInvalidUserID
	}

	urls, err := s.repo.GetURLsByUserID(ctx, currentUserID)
	if err != nil {
		return nil, errors.ErrFailedToLoadUserUrls
	}

	results := make([]dto.GetUserURLsResponse, len(urls))
	for i, url := range urls {
		results[i] = dto.GetUserURLsResponse{
			ShortURL:    fmt.Sprintf("%s/%s", s.cfg.BaseURL, url.ShortCode),
			OriginalURL: url.LongURL,
		}
	}

	return results, nil
}
