package server

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"shortly/internal/app/helpers"
	"shortly/internal/app/store"
)

type HandlerConfig struct {
	SecureRandom helpers.SecureRandomGenerator
	Store        *store.URLStore
}

func NewHandlerConfig(SecureRandom helpers.SecureRandomGenerator, Store *store.URLStore) *HandlerConfig {
	return &HandlerConfig{
		SecureRandom: SecureRandom,
		Store:        Store,
	}
}

func (config *HandlerConfig) HandleCreateShortLink(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Wrong HTTP method", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Unable to read request body", http.StatusBadRequest)
		return
	}
	if len(body) == 0 {
		http.Error(res, "Unable to process request", http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	longURL := strings.TrimSpace(string(body))
	longURL = strings.Trim(longURL, "\"")

	if helpers.IsInvalidURL(longURL) {
		http.Error(res, "Invalid URL", http.StatusBadRequest)
		return
	}

	shortCode, err := config.SecureRandom.Hex()
	if err != nil {
		http.Error(res, "Failed to generate short code", http.StatusInternalServerError)
		return
	}

	shortURL := fmt.Sprintf("%s/%s", options.BaseURL, shortCode)

	config.Store.Set(shortCode, longURL)

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)

	_, err = res.Write([]byte(shortURL))
	if err != nil {
		http.Error(res, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (config *HandlerConfig) HandleGetShortLink(res http.ResponseWriter, req *http.Request) {
	shortCode := strings.TrimPrefix(req.URL.Path, "/")

	longURL, found := config.Store.Get(shortCode)
	if !found {
		http.Error(res, "Short code not found", http.StatusNotFound)
		return
	}

	res.Header().Set("Content-Type", "text/plain")
	http.Redirect(res, req, longURL, http.StatusTemporaryRedirect)
}
