package api

import (
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

func NewURLHandler(cfg *config.Config, service *service.URLService) *URLHandler {
	return &URLHandler{cfg: cfg, service: service}
}

func (h *URLHandler) HandleCreateShortLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	shortURL, err := h.service.CreateShortLink(r)
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

func (h *URLHandler) HandleGetShortLink(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "id")

	url, found := h.service.GetShortLink(shortCode)
	if !found {
		http.Error(w, errors.ErrShortLinkNotFound.Error(), http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url.LongURL, http.StatusTemporaryRedirect)
}
