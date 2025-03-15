package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"shortly/internal/app/api/pagination"
	"shortly/internal/app/config"
	"shortly/internal/app/dto"
	"shortly/internal/app/errors"
	"shortly/internal/app/repository"
	"shortly/internal/app/worker"
)

// Shortener is a service for URL shortening
type Shortener interface {
	CreateShortLink(ctx context.Context, longURL string) (string, error)
	CreateShortLinks(ctx context.Context, params []dto.BatchCreateShortLinkParams) ([]dto.BatchCreateShortLinkResponse, error)
	GetShortLink(ctx context.Context, shortCode string) (*repository.URL, bool)
	GetUserURLs(ctx context.Context, pagination *pagination.Pagination) ([]dto.GetUserURLsResponse, int, error)
	DeleteUserURLs(ctx context.Context, params dto.BatchDeleteShortLinkRequest) error
}

// URLService is a service for URL operations
type URLService struct {
	cfg    *config.Config
	repo   repository.Repository
	rand   SecureRandomGenerator
	worker worker.Worker
}

// NewURLService creates a new URL service instance
func NewURLService(cfg *config.Config, repo repository.Repository, rand SecureRandomGenerator, worker worker.Worker) *URLService {
	return &URLService{
		cfg:    cfg,
		repo:   repo,
		rand:   rand,
		worker: worker,
	}
}

// CreateShortLink creates a new short link
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

// CreateShortLinks creates new short links
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

// GetShortLink returns a short link by short code
func (s *URLService) GetShortLink(ctx context.Context, shortCode string) (*repository.URL, bool) {
	return s.repo.GetURLByShortCode(ctx, shortCode)
}

// GetUserURLs returns user URLs
func (s *URLService) GetUserURLs(ctx context.Context, pagination *pagination.Pagination) ([]dto.GetUserURLsResponse, int, error) {
	currentUserID, ok := (ctx.Value(dto.CurrentUser)).(uuid.UUID)
	if !ok {
		return nil, 0, errors.ErrInvalidUserID
	}

	urls, total, err := s.repo.GetURLsByUserID(ctx, currentUserID, pagination.Per, pagination.Offset())
	if err != nil {
		return nil, 0, errors.ErrFailedToLoadUserUrls
	}

	results := make([]dto.GetUserURLsResponse, len(urls))
	for i, url := range urls {
		results[i] = dto.GetUserURLsResponse{
			ShortURL:    fmt.Sprintf("%s/%s", s.cfg.BaseURL, url.ShortCode),
			OriginalURL: url.LongURL,
		}
	}

	return results, total, nil
}

// DeleteUserURLs deletes user URLs
func (s *URLService) DeleteUserURLs(ctx context.Context, params dto.BatchDeleteShortLinkRequest) error {
	currentUserID, ok := (ctx.Value(dto.CurrentUser)).(uuid.UUID)
	if !ok {
		return errors.ErrInvalidUserID
	}

	s.worker.Add(dto.BatchDeleteParams{
		UserID:     currentUserID,
		ShortCodes: params,
	})

	return nil
}

// generateUniqueShortCode generates a unique short code
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
