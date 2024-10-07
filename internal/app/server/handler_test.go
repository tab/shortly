package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleCreateShortLink(t *testing.T) {
	type result struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name     string
		method   string
		body     string
		expected result
	}{
		{
			name:   "Success",
			method: http.MethodPost,
			body:   "https://example.com",
			expected: result{
				code:        201,
				response:    "http://localhost:8080/abcd1234",
				contentType: "text/plain",
			},
		},
		{
			name:   "Wrong HTTP method",
			method: http.MethodGet,
			body:   "",
			expected: result{
				code:        400,
				response:    "Wrong HTTP method\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "Empty body",
			method: http.MethodPost,
			body:   "",
			expected: result{
				code:        400,
				response:    "Unable to process request\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "Invalid body",
			method: http.MethodPost,
			body:   "not-a-url",
			expected: result{
				code:        400,
				response:    "invalid body\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "Read body error",
			method: http.MethodPost,
			body:   "invalid-body",
			expected: result{
				code:        400,
				response:    "Unable to read request body\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, "/", strings.NewReader(test.body))
			recorder := httptest.NewRecorder()

			if test.method == http.MethodPost && test.body == "invalid-body" {
				request.Body = io.NopCloser(io.LimitReader(strings.NewReader(test.body), -1))
			}

			HandleCreateShortLink(recorder, request)

			response := recorder.Result()

			assert.Equal(t, test.expected.code, response.StatusCode)

			defer response.Body.Close()
			responseBody, err := io.ReadAll(response.Body)

			require.NoError(t, err)
			assert.NotEmpty(t, string(responseBody))
			assert.Equal(t, test.expected.contentType, response.Header.Get("Content-Type"))
		})
	}
}

func TestHandleGetShortLink(t *testing.T) {
	storage.Set("abcd1234", "https://example.com")

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
			name: "Redirect",
			path: "/abcd1234",
			expected: result{
				code:        307,
				response:    "",
				contentType: "text/plain",
			},
		},
		{
			name: "Not found code",
			path: "/not-valid-code",
			expected: result{
				code:        404,
				response:    "Short code not found\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Empty code",
			path: "/",
			expected: result{
				code:        404,
				response:    "Short code not found\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.path, nil)
			recorder := httptest.NewRecorder()
			HandleGetShortLink(recorder, request)

			response := recorder.Result()

			assert.Equal(t, test.expected.code, response.StatusCode)

			defer response.Body.Close()
			responseBody, err := io.ReadAll(response.Body)

			require.NoError(t, err)
			assert.Equal(t, test.expected.response, string(responseBody))
			assert.Equal(t, test.expected.contentType, response.Header.Get("Content-Type"))
		})
	}
}
