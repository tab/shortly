package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"shortly/internal/app/config"
	"shortly/internal/app/helpers"
	"shortly/internal/app/store"
)

type Handler struct {
	AppConfig    *config.AppConfig
	SecureRandom helpers.SecureRandomGenerator
	Store        store.URLStore
}

func (h *Handler) HandleCreateShortLink() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(res, "Wrong HTTP method", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(req.Body)
		if err != nil || len(body) == 0 {
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

		shortCode, err := h.SecureRandom.Hex()
		if err != nil {
			http.Error(res, "Failed to generate short code", http.StatusInternalServerError)
			return
		}

		shortURL := fmt.Sprintf("%s/%s", h.AppConfig.BaseURL, shortCode)
		h.Store.Set(shortCode, longURL)

		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusCreated)

		_, err = res.Write([]byte(shortURL))
		if err != nil {
			http.Error(res, "Failed to write response", http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) HandleGetShortLink() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		shortCode := strings.TrimPrefix(req.URL.Path, "/")

		longURL, found := h.Store.Get(shortCode)
		if !found {
			http.Error(res, "Short code not found", http.StatusNotFound)
			return
		}

		res.Header().Set("Content-Type", "text/plain")
		http.Redirect(res, req, longURL, http.StatusTemporaryRedirect)
	}
}
