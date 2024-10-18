package main

import (
	"net/http"
	"net/http/httptest"
	"shortly/internal/app/repository"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"shortly/internal/app/config"
)

func Test_Run(t *testing.T) {
	errCh := make(chan error, 1)

	go func() {
		err := run()

		if err != nil && err != http.ErrServerClosed {
			errCh <- err
		}

		close(errCh)
	}()

	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:8080/health")
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_HealthCheck(t *testing.T) {
	cfg := &config.Config{
		ClientURL: "http://localhost:8080",
	}
	router := setupRouter(cfg)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	resp := w.Result()
	defer resp.Body.Close()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_CreateShortLink(t *testing.T) {
	cfg := &config.Config{
		Addr: "localhost:8080",
	}
	router := setupRouter(cfg)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://example.com"))
	w := httptest.NewRecorder()

	resp := w.Result()
	defer resp.Body.Close()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_GetShortLink(t *testing.T) {
	cfg := &config.Config{
		Addr: "localhost:8080",
	}
	repo := repository.NewInMemoryRepository()
	router := setupRouter(cfg)

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
