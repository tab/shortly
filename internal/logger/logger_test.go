package logger

import (
	"bytes"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_GetLogger(t *testing.T) {
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
			logger := GetLogger()

			assert.Equal(t, tt.expectedLevel, logger.GetLevel())
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
			logger := GetLogger()
			logger = logger.Output(&buf)

			log = logger

			req := httptest.NewRequest(tt.method, tt.path, nil)
			req.RemoteAddr = "127.0.0.1:1234"
			req.Host = "example.com"
			w := httptest.NewRecorder()

			handler := Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
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
