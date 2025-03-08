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

func Test_HealthHandler_HandleLiveness(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockHealthChecker(ctrl)
	handler := NewHealthHandler(mockService)

	type result struct {
		response dto.HealthResponse
		code     int
		status   string
	}

	tests := []struct {
		name     string
		expected result
	}{
		{
			name: "Success",
			expected: result{
				response: dto.HealthResponse{Result: "alive"},
				code:     http.StatusOK,
				status:   "200 OK",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/live", nil)
			w := httptest.NewRecorder()

			handler.HandleLiveness(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			var actual dto.HealthResponse
			err := json.NewDecoder(resp.Body).Decode(&actual)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.response.Result, actual.Result)
			assert.Equal(t, tt.expected.status, resp.Status)
			assert.Equal(t, tt.expected.code, resp.StatusCode)
		})
	}
}

func Test_HealthHandler_HandleReadiness(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockHealthChecker(ctrl)
	handler := NewHealthHandler(mockService)

	type result struct {
		response dto.HealthResponse
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
				mockService.EXPECT().Ping(gomock.Any()).Return(nil)
			},
			expected: result{
				response: dto.HealthResponse{Result: "ready"},
				error:    dto.ErrorResponse{},
				code:     http.StatusOK,
				status:   "200 OK",
			},
		},
		{
			name: "Service Error",
			before: func() {
				mockService.EXPECT().Ping(gomock.Any()).Return(assert.AnError)
			},
			expected: result{
				response: dto.HealthResponse{},
				error:    dto.ErrorResponse{Error: assert.AnError.Error()},
				code:     http.StatusInternalServerError,
				status:   "500 Internal Server Error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			req := httptest.NewRequest("GET", "/ready", nil)
			w := httptest.NewRecorder()

			handler.HandleReadiness(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.expected.error.Error != "" {
				var actual dto.ErrorResponse
				err := json.NewDecoder(resp.Body).Decode(&actual)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error.Error, actual.Error)
			} else {
				var actual dto.HealthResponse
				err := json.NewDecoder(resp.Body).Decode(&actual)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response.Result, actual.Result)
			}
			assert.Equal(t, tt.expected.status, resp.Status)
			assert.Equal(t, tt.expected.code, resp.StatusCode)
		})
	}
}

func Test_HealthHandler_HandlePing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockHealthChecker(ctrl)
	handler := NewHealthHandler(mockService)

	type result struct {
		response dto.HealthResponse
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
				mockService.EXPECT().Ping(gomock.Any()).Return(nil)
			},
			expected: result{
				response: dto.HealthResponse{Result: "pong"},
				error:    dto.ErrorResponse{},
				code:     http.StatusOK,
				status:   "200 OK",
			},
		},
		{
			name: "Service Error",
			before: func() {
				mockService.EXPECT().Ping(gomock.Any()).Return(assert.AnError)
			},
			expected: result{
				response: dto.HealthResponse{},
				error:    dto.ErrorResponse{Error: assert.AnError.Error()},
				code:     http.StatusInternalServerError,
				status:   "500 Internal Server Error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			req := httptest.NewRequest("GET", "/ping", nil)
			w := httptest.NewRecorder()

			handler.HandlePing(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.expected.error.Error != "" {
				var actual dto.ErrorResponse
				err := json.NewDecoder(resp.Body).Decode(&actual)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error.Error, actual.Error)
			} else {
				var actual dto.HealthResponse
				err := json.NewDecoder(resp.Body).Decode(&actual)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response.Result, actual.Result)
			}
			assert.Equal(t, tt.expected.status, resp.Status)
			assert.Equal(t, tt.expected.code, resp.StatusCode)
		})
	}
}
