package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"shortly/internal/app/api/pagination"
	"shortly/internal/app/config"
	"shortly/internal/app/dto"
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
	w.Header().Set("Content-Type", "application/json")

	var params dto.CreateShortLinkRequest

	if err := params.Validate(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	shortURL, err := h.service.CreateShortLink(r.Context(), params.URL)
	if err != nil {
		if errors.Is(err, errors.ErrURLAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(dto.CreateShortLinkResponse{Result: shortURL})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.CreateShortLinkResponse{Result: shortURL})
}

func (h *URLHandler) HandleBatchCreateShortLink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params dto.BatchCreateShortLinkRequest

	if err := params.Validate(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	results, err := h.service.CreateShortLinks(r.Context(), params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(results)
}

func (h *URLHandler) HandleGetShortLink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	shortCode := chi.URLParam(r, "id")

	result, found := h.service.GetShortLink(r.Context(), shortCode)
	if !found {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: errors.ErrShortLinkNotFound.Error()})
		return
	}

	if !result.DeletedAt.IsZero() {
		w.WriteHeader(http.StatusGone)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: errors.ErrShortLinkDeleted.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.GetShortLinkResponse{Result: result.LongURL})
}

func (h *URLHandler) HandleGetUserURLs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	paginator := pagination.NewPagination(r)

	urls, _, err := h.service.GetUserURLs(r.Context(), paginator)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// TODO: use paginated response
	// response := dto.PaginatedResponse{
	//	Data:  urls,
	//	Page:  paginator.Page,
	//	Per:   paginator.Per,
	//	Total: total,
	// }

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(urls)
}

func (h *URLHandler) HandleBatchDeleteUserURLs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params dto.BatchDeleteShortLinkRequest

	if err := params.Validate(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	err := h.service.DeleteUserURLs(r.Context(), params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// NOTE: text/plain request is deprecated
func (h *URLHandler) DeprecatedHandleCreateShortLink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	var params dto.CreateShortLinkRequest

	if err := params.DeprecatedValidate(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL, err := h.service.CreateShortLink(r.Context(), params.URL)
	if err != nil {
		if errors.Is(err, errors.ErrURLAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			_, err = w.Write([]byte(shortURL))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortURL))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// NOTE: text/plain request is deprecated
func (h *URLHandler) DeprecatedHandleGetShortLink(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "id")

	result, found := h.service.GetShortLink(r.Context(), shortCode)
	if !found {
		http.Error(w, errors.ErrShortLinkNotFound.Error(), http.StatusNotFound)
		return
	}

	if !result.DeletedAt.IsZero() {
		http.Error(w, errors.ErrShortLinkDeleted.Error(), http.StatusGone)
		return
	}

	http.Redirect(w, r, result.LongURL, http.StatusTemporaryRedirect)
}
