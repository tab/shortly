package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"shortly/internal/app/config"
	"shortly/internal/app/errors"
	"shortly/internal/app/repository"
	"shortly/internal/app/validator"
)

type URLService struct {
	cfg  *config.Config
	repo repository.URLRepository
	rand SecureRandomGenerator
}

func NewURLService(cfg *config.Config, repo repository.URLRepository, rand SecureRandomGenerator) *URLService {
	return &URLService{
		cfg:  cfg,
		repo: repo,
		rand: rand,
	}
}

type RequestParams struct {
	URL string `json:"url"`
}

func (s *URLService) CreateShortLink(r *http.Request) (string, error) {
	if r.ContentLength == 0 {
		return "", errors.ErrRequestBodyEmpty
	}

	var params RequestParams

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		return "", err
	}

	longURL := strings.Trim(strings.TrimSpace(params.URL), "\"")

	if err = validator.Validate(longURL); err != nil {
		return "", err
	}

	shortCode, err := s.rand.Hex()
	if err != nil {
		return "", errors.ErrCouldNotGenerateCode
	}

	url := repository.URL{
		LongURL:   longURL,
		ShortCode: shortCode,
	}
	s.repo.Set(url)

	return fmt.Sprintf("%s/%s", s.cfg.BaseURL, shortCode), nil
}

// NOTE: text/plain request is deprecated
func (s *URLService) DeprecatedCreateShortLink(r *http.Request) (string, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		return "", errors.ErrRequestBodyEmpty
	}
	defer r.Body.Close()

	longURL := strings.Trim(strings.TrimSpace(string(body)), "\"")

	if err = validator.Validate(longURL); err != nil {
		return "", err
	}

	shortCode, err := s.rand.Hex()
	if err != nil {
		return "", errors.ErrCouldNotGenerateCode
	}

	url := repository.URL{
		LongURL:   longURL,
		ShortCode: shortCode,
	}
	s.repo.Set(url)

	return fmt.Sprintf("%s/%s", s.cfg.BaseURL, shortCode), nil
}

func (s *URLService) GetShortLink(shortCode string) (*repository.URL, bool) {
	return s.repo.Get(shortCode)
}
