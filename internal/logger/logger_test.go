package logger

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func Test_NewLogger(t *testing.T) {
	tests := []struct {
		name          string
		expectedLevel zerolog.Level
	}{
		{
			name:          "Default configuration",
			expectedLevel: zerolog.Level(LogLevel),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger()

			assert.Equal(t, tt.expectedLevel, logger.log.GetLevel())
			assert.NotNil(t, logger)
		})
	}
}

func Test_LoggerMiddleware(t *testing.T) {
	type result struct {
		path       string
		method     string
		duration   string
		requestURI string
	}

	tests := []struct {
		name     string
		method   string
		path     string
		expected result
	}{
		{
			name:   "GET request logging",
			method: http.MethodGet,
			path:   "/test-get",
			expected: result{
				method:     http.MethodGet,
				path:       "/test-get",
				duration:   "duration",
				requestURI: "/test-get",
			},
		},
		{
			name:   "POST request logging",
			method: http.MethodPost,
			path:   "/test-post",
			expected: result{
				method:     http.MethodPost,
				path:       "/test-post",
				duration:   "duration",
				requestURI: "/test-post",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewLogger()
			logger.log = logger.log.Output(&buf)

			req := httptest.NewRequest(tt.method, tt.path, nil)
			req.RemoteAddr = "127.0.0.1:1234"
			req.Host = "example.com"
			w := httptest.NewRecorder()

			handler := logger.Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			handler.ServeHTTP(w, req)

			result := buf.String()

			assert.Contains(t, result, tt.expected.method)
			assert.Contains(t, result, tt.expected.path)
			assert.Contains(t, result, tt.expected.duration)
			assert.Contains(t, result, tt.expected.requestURI)
		})
	}
}

func Test_Logger_Info(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "Success",
			expected: "info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewLogger()
			logger.log = logger.log.Output(&buf)

			logger.Info().Msg(tt.expected)

			result := buf.String()

			assert.Contains(t, result, tt.expected)
		})
	}
}

func Test_Logger_Warn(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "Success",
			expected: "warn",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewLogger()
			logger.log = logger.log.Output(&buf)

			logger.Warn().Msg(tt.expected)

			result := buf.String()

			assert.Contains(t, result, tt.expected)
		})
	}
}

func Test_Logger_Error(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "Success",
			expected: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewLogger()
			logger.log = logger.log.Output(&buf)

			logger.Error().Msg(tt.expected)

			result := buf.String()

			assert.Contains(t, result, tt.expected)
		})
	}
}
