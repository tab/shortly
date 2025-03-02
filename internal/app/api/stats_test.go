package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"shortly/internal/app/dto"
	"shortly/internal/app/service"
)

func Test_StatsHandler_HandleStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockStatsReporter(ctrl)
	handler := NewStatsHandler(mockService)

	type result struct {
		response dto.StatsResponse
		error    dto.ErrorResponse
		code     int
		status   string
	}

	tests := []struct {
		name     string
		before   func()
		expected result
	}{
		{
			name: "Success",
			before: func() {
				mockService.EXPECT().Counters(gomock.Any()).Return(10, 2, nil)
			},
			expected: result{
				response: dto.StatsResponse{URLs: 10, Users: 2},
				error:    dto.ErrorResponse{},
				code:     http.StatusOK,
				status:   "200 OK",
			},
		},
		{
			name: "Service Error",
			before: func() {
				mockService.EXPECT().Counters(gomock.Any()).Return(0, 0, assert.AnError)
			},
			expected: result{
				response: dto.StatsResponse{},
				error:    dto.ErrorResponse{Error: assert.AnError.Error()},
				code:     http.StatusInternalServerError,
				status:   "500 Internal Server Error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			req := httptest.NewRequest("GET", "/api/internal/stats", nil)
			w := httptest.NewRecorder()

			handler.HandleStats(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.expected.error.Error != "" {
				var actual dto.ErrorResponse
				err := json.NewDecoder(resp.Body).Decode(&actual)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error.Error, actual.Error)
			} else {
				var actual dto.StatsResponse
				err := json.NewDecoder(resp.Body).Decode(&actual)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response.URLs, actual.URLs)
				assert.Equal(t, tt.expected.response.Users, actual.Users)
			}
			assert.Equal(t, tt.expected.status, resp.Status)
			assert.Equal(t, tt.expected.code, resp.StatusCode)
		})
	}
}
