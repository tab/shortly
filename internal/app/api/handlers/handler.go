package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"shortly/internal/app/config"
	"shortly/internal/app/errors"
	"shortly/internal/app/helpers"
	"shortly/internal/app/store"
)

type Handler struct {
	AppConfig    *config.AppConfig
	SecureRandom helpers.SecureRandomGenerator
	Store        store.URLStore
}

func (h *Handler) HandleCreateShortLink(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		httpError(res, &errors.MethodNotAllowedError{}, http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil || len(body) == 0 {
		httpError(res, &errors.InvalidRequestBodyError{}, http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	longURL := strings.TrimSpace(string(body))
	longURL = strings.Trim(longURL, "\"")

	if helpers.IsInvalidURL(longURL) {
		httpError(res, &errors.InvalidURLError{URL: longURL}, http.StatusBadRequest)
		return
	}

	shortCode, err := h.SecureRandom.Hex()
	if err != nil {
		httpError(res, &errors.ShortCodeGenerationError{}, http.StatusInternalServerError)
		return
	}

	shortURL := fmt.Sprintf("%s/%s", h.AppConfig.BaseURL, shortCode)
	h.Store.Set(shortCode, longURL)

	httpResponse(res, http.StatusCreated, []byte(shortURL), "")
}

func (h *Handler) HandleGetShortLink(res http.ResponseWriter, req *http.Request) {
	shortCode := strings.TrimPrefix(req.URL.Path, "/")

	longURL, found := h.Store.Get(shortCode)
	if !found {
		httpError(res, &errors.ShortCodeNotFoundError{Code: shortCode}, http.StatusNotFound)
		return
	}

	httpResponse(res, http.StatusTemporaryRedirect, nil, longURL)
}
