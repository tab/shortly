package api

import (
	"net/http"

	"github.com/goccy/go-json"

	"shortly/internal/app/dto"
	"shortly/internal/app/service"
)

// HealthHandler is a handler for health check
type HealthHandler struct {
	service service.HealthChecker
}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler(service service.HealthChecker) *HealthHandler {
	return &HealthHandler{service: service}
}

// HandlePing handles ping request
func (h *HealthHandler) HandlePing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := h.service.Ping(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(dto.PingResponse{Result: "pong"})
}
