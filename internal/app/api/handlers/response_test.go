package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"shortly/internal/app/errors"
)

type MockResponseWriter struct {
	httptest.ResponseRecorder
}

func TestHttpResponse(t *testing.T) {
	type params struct {
		code        int
		body        []byte
		redirectURL string
	}

	type result struct {
		code        int
		contentType string
	}

	tests := []struct {
		name     string
		params   params
		expected result
	}{
		{
			name: "Success",
			params: params{
				code: http.StatusOK,
				body: []byte("Success"),
			},
			expected: result{
				code:        http.StatusOK,
				contentType: "text/plain",
			},
		},
		{
			name: "Redirect",
			params: params{
				code:        http.StatusMovedPermanently,
				body:        nil,
				redirectURL: "http://example.com",
			},
			expected: result{
				code:        http.StatusMovedPermanently,
				contentType: "text/plain",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := &MockResponseWriter{}
			httpResponse(response, test.params.code, test.params.body, test.params.redirectURL)

			assert.Equal(t, test.expected.code, response.Code)
			assert.Equal(t, test.expected.contentType, response.Header().Get("Content-Type"))
		})
	}
}

func TestHttpError(t *testing.T) {
	type params struct {
		err  error
		code int
	}

	type result struct {
		code        int
		contentType string
	}

	tests := []struct {
		name     string
		params   params
		expected result
	}{
		{
			name: "NotFound",
			params: params{
				err:  &errors.ShortCodeNotFoundError{Code: "not-valid-code"},
				code: http.StatusNotFound,
			},
			expected: result{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "InternalServerError",
			params: params{
				err:  &errors.ResponseWriteError{},
				code: http.StatusInternalServerError,
			},
			expected: result{
				code:        http.StatusInternalServerError,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := &MockResponseWriter{}
			httpError(response, test.params.err, test.params.code)

			assert.Equal(t, test.expected.code, response.Code)
			assert.Equal(t, test.expected.contentType, response.Header().Get("Content-Type"))
		})
	}
}
