package service

import (
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
	repo repository.URLRepository
	rand SecureRandomGenerator
	cfg  *config.Config
}

func NewURLService(repo repository.URLRepository, rand SecureRandomGenerator, cfg *config.Config) *URLService {
	return &URLService{
		repo: repo,
		rand: rand,
		cfg:  cfg,
	}
}

func (s *URLService) CreateShortLink(r *http.Request) (string, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		return "", errors.ErrRequestBodyEmpty
	}
	defer r.Body.Close()

	longURL := strings.TrimSpace(string(body))
	longURL = strings.Trim(longURL, "\"")

	if err := validator.Validate(longURL); err != nil {
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
