package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"shortly/internal/app/config"
	"shortly/internal/app/store"
)

type MockSecureRandom struct{}

func (MockSecureRandom) Hex() (string, error) {
	return "abcd1234", nil
}

func TestHandleCreateShortLink(t *testing.T) {
	appConfig := &config.AppConfig{
		Addr:      config.ServerAddress,
		BaseURL:   config.BaseURL,
		ClientURL: config.ClientURL,
	}

	handler := &Handler{
		AppConfig:    appConfig,
		SecureRandom: MockSecureRandom{},
		Store:        *store.NewURLStore(),
	}

	type result struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name     string
		method   string
		body     string
		expected struct {
			code        int
			response    string
			contentType string
		}
	}{
		{
			name:   "Success",
			method: http.MethodPost,
			body:   "https://example.com",
			expected: result{
				code:        http.StatusCreated,
				response:    "http://localhost:8080/abcd1234",
				contentType: "text/plain",
			},
		},
		{
			name:   "Wrong HTTP method",
			method: http.MethodGet,
			body:   "",
			expected: result{
				code:        http.StatusBadRequest,
				response:    "Wrong HTTP method\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "Empty body",
			method: http.MethodPost,
			body:   "",
			expected: result{
				code:        http.StatusBadRequest,
				response:    "Unable to process request\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "Invalid URL: not-a-url\n",
			method: http.MethodPost,
			body:   "not-a-url",
			expected: result{
				code:        http.StatusBadRequest,
				response:    "Invalid URL: not-a-url\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, "/", strings.NewReader(test.body))
			recorder := httptest.NewRecorder()

			handler.HandleCreateShortLink()(recorder, request)

			response := recorder.Result()
			defer response.Body.Close()

			body, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			assert.Equal(t, test.expected.code, response.StatusCode)
			assert.Equal(t, test.expected.contentType, response.Header.Get("Content-Type"))
			assert.Equal(t, test.expected.response, string(body))
		})
	}
}

func TestHandleGetShortLink(t *testing.T) {
	appConfig := &config.AppConfig{
		Addr:      config.ServerAddress,
		BaseURL:   config.BaseURL,
		ClientURL: config.ClientURL,
	}

	handler := &Handler{
		AppConfig:    appConfig,
		SecureRandom: MockSecureRandom{},
		Store:        *store.NewURLStore(),
	}

	handler.Store.Set("abcd1234", "https://example.com")

	type result struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name     string
		path     string
		expected result
	}{
		{
			name: "Redirect to long URL",
			path: "/abcd1234",
			expected: result{
				code:        http.StatusTemporaryRedirect,
				response:    "",
				contentType: "text/plain",
			},
		},
		{
			name: "Not found",
			path: "/not-valid-code",
			expected: result{
				code:        http.StatusNotFound,
				response:    "Short code not found: not-valid-code\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.path, nil)
			recorder := httptest.NewRecorder()

			handler.HandleGetShortLink()(recorder, request)

			response := recorder.Result()
			defer response.Body.Close()

			body, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			assert.Equal(t, test.expected.code, response.StatusCode)
			assert.Equal(t, test.expected.response, string(body))
			assert.Equal(t, test.expected.contentType, response.Header.Get("Content-Type"))
		})
	}
}
