package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"shortly/internal/app/config"
	"shortly/internal/app/errors"
	"shortly/internal/app/service"
)

type URLHandler struct {
	cfg     *config.Config
	service *service.URLService
}

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewURLHandler(cfg *config.Config, service *service.URLService) *URLHandler {
	return &URLHandler{cfg: cfg, service: service}
}

func (h *URLHandler) HandleCreateShortLink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	shortURL, err := h.service.CreateShortLink(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{Result: shortURL})
}

func (h *URLHandler) HandleGetShortLink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	shortCode := chi.URLParam(r, "id")

	url, found := h.service.GetShortLink(shortCode)
	if !found {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: errors.ErrShortLinkNotFound.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{Result: url.LongURL})
}

// NOTE: text/plain request is deprecated
func (h *URLHandler) DeprecatedHandleCreateShortLink(w http.ResponseWriter, r *http.Request) {
	shortURL, err := h.service.DeprecatedCreateShortLink(r)
	if err != nil {
		switch {
		case errors.Is(err, errors.ErrRequestBodyEmpty):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, errors.ErrInvalidURL):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, errors.ErrCouldNotGenerateCode):
			http.Error(w, err.Error(), http.StatusInternalServerError)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortURL))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// NOTE: text/plain request is deprecated
func (h *URLHandler) DeprecatedHandleGetShortLink(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "id")

	url, found := h.service.GetShortLink(shortCode)
	if !found {
		http.Error(w, errors.ErrShortLinkNotFound.Error(), http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url.LongURL, http.StatusTemporaryRedirect)
}
