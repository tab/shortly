package api

import (
	"encoding/json"
	"net/http"

	"shortly/internal/app/dto"
	"shortly/internal/app/service"
)

// StatsHandler is a handler for stats
type StatsHandler struct {
	service service.StatsReporter
}

// NewStatsHandler creates a new StatsHandler
func NewStatsHandler(service service.StatsReporter) *StatsHandler {
	return &StatsHandler{service: service}
}

// HandleStats handles stats request
func (h *StatsHandler) HandleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	urls, users, err := h.service.Counters(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(dto.StatsResponse{URLs: urls, Users: users})
}
