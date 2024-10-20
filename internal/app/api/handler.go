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
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
	Code   int    `json:"code"`
	Status string `json:"status"`
}

func NewURLHandler(cfg *config.Config, service *service.URLService) *URLHandler {
	return &URLHandler{cfg: cfg, service: service}
}

func (h *URLHandler) HandleCreateShortLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response := Response{
			Error:  "Invalid request method",
			Status: http.StatusText(http.StatusBadRequest),
			Code:   http.StatusBadRequest,
		}
		renderResponse(w, response)
		return
	}

	shortURL, err := h.service.CreateShortLink(r)
	if err != nil {
		response := mapErrorToResponse(err)
		renderResponse(w, response)
		return
	}

	response := Response{
		Result: shortURL,
		Status: http.StatusText(http.StatusCreated),
		Code:   http.StatusCreated,
	}
	renderResponse(w, response)
}

func (h *URLHandler) HandleGetShortLink(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "id")

	url, found := h.service.GetShortLink(shortCode)
	if !found {
		response := Response{
			Error:  errors.ErrShortLinkNotFound.Error(),
			Status: http.StatusText(http.StatusNotFound),
			Code:   http.StatusNotFound,
		}
		renderResponse(w, response)
		return
	}

	response := Response{
		Result: url.LongURL,
		Status: http.StatusText(http.StatusOK),
		Code:   http.StatusOK,
	}
	renderResponse(w, response)
}

func renderResponse(w http.ResponseWriter, response Response) {
	w.Header().Set("Content-Type", "application/json")

	resp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(response.Code)
	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func mapErrorToResponse(err error) Response {
	switch {
	case errors.Is(err, errors.ErrRequestBodyEmpty):
		return Response{
			Error:  err.Error(),
			Status: http.StatusText(http.StatusBadRequest),
			Code:   http.StatusBadRequest,
		}
	case errors.Is(err, errors.ErrInvalidURL):
		return Response{
			Error:  err.Error(),
			Status: http.StatusText(http.StatusBadRequest),
			Code:   http.StatusBadRequest,
		}
	case errors.Is(err, errors.ErrCouldNotGenerateCode):
		return Response{
			Error:  err.Error(),
			Status: http.StatusText(http.StatusInternalServerError),
			Code:   http.StatusInternalServerError,
		}
	default:
		return Response{
			Error:  "Internal server error",
			Status: http.StatusText(http.StatusInternalServerError),
			Code:   http.StatusInternalServerError,
		}
	}
}

// NOTE: text/plain request is deprecated
func (h *URLHandler) DeprecatedHandleCreateShortLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

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
