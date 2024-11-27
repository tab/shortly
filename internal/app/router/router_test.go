package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"shortly/internal/app/config"
	"shortly/internal/app/repository"
	"shortly/internal/app/worker"
	"shortly/internal/logger"
)

func Test_HealthCheck(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		DatabaseDSN: "postgres://postgres:postgres@localhost:5432/shortly-test?sslmode=disable",
	}
	appLogger := logger.NewLogger()
	repo, _ := repository.NewRepository(ctx, &repository.Factory{
		DSN:    cfg.DatabaseDSN,
		Logger: appLogger,
	})
	appWorker := worker.NewDeleteWorker(ctx, cfg, repo, appLogger)
	router := NewRouter(cfg, repo, appWorker, appLogger)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	resp := w.Result()
	defer resp.Body.Close()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_CreateShortLink(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
	}
	appLogger := logger.NewLogger()
	repo := repository.NewInMemoryRepository()
	appWorker := worker.NewDeleteWorker(ctx, cfg, repo, appLogger)
	router := NewRouter(cfg, repo, appWorker, appLogger)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://example.com"))
	w := httptest.NewRecorder()

	resp := w.Result()
	defer resp.Body.Close()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_GetShortLink(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
	}
	appLogger := logger.NewLogger()
	repo := repository.NewInMemoryRepository()
	appWorker := worker.NewDeleteWorker(ctx, cfg, repo, appLogger)
	router := NewRouter(cfg, repo, appWorker, appLogger)

	UUID, _ := uuid.Parse("6455bd07-e431-4851-af3c-4f703f726639")

	_, err := repo.CreateURL(ctx, repository.URL{
		UUID:      UUID,
		LongURL:   "https://example.com",
		ShortCode: "abcd1234",
	})
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/abcd1234", nil)
	w := httptest.NewRecorder()

	resp := w.Result()
	defer resp.Body.Close()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
