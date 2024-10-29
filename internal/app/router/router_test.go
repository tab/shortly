package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"shortly/internal/app/config"
	"shortly/internal/app/repository"
	"shortly/internal/logger"
)

func Test_HealthCheck(t *testing.T) {
	cfg := &config.Config{
		ClientURL: "http://localhost:8080",
	}
	repo := repository.NewRepository()
	appLogger := logger.NewLogger()
	router := NewRouter(cfg, appLogger, repo)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	resp := w.Result()
	defer resp.Body.Close()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_CreateShortLink(t *testing.T) {
	cfg := &config.Config{
		ClientURL: "http://localhost:8080",
	}
	repo := repository.NewRepository()
	appLogger := logger.NewLogger()
	router := NewRouter(cfg, appLogger, repo)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://example.com"))
	w := httptest.NewRecorder()

	resp := w.Result()
	defer resp.Body.Close()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_GetShortLink(t *testing.T) {
	cfg := &config.Config{
		ClientURL: "http://localhost:8080",
	}
	repo := repository.NewRepository()
	appLogger := logger.NewLogger()
	router := NewRouter(cfg, appLogger, repo)

	repo.Set(repository.URL{
		LongURL:   "https://example.com",
		ShortCode: "abcd1234",
	})

	req := httptest.NewRequest(http.MethodGet, "/abcd1234", nil)
	w := httptest.NewRecorder()

	resp := w.Result()
	defer resp.Body.Close()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
