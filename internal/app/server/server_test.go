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

func requestHelper(t *testing.T, server *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, server.URL+path, body)
	require.NoError(t, err)

	resp, err := server.Client().Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}

func TestAppRouter(t *testing.T) {
	storage.Set("abcd1234", "https://example.com")

	server := httptest.NewServer(AppRouter())
	defer server.Close()

	var tests = []struct {
		name   string
		path   string
		method string
		body   string
		code   int
	}{
		{
			name:   "Create short link",
			path:   "/",
			method: "POST",
			body:   "https://example.com",
			code:   http.StatusCreated,
		},
		{
			name:   "Get short link",
			path:   "/abcd1234",
			method: "GET",
			body:   "",
			code:   http.StatusOK,
		},
		{
			name:   "Short code not found",
			path:   "/not-valid-code",
			method: "GET",
			body:   "",
			code:   http.StatusNotFound,
		},
	}
	for _, test := range tests {
		resp, _ := requestHelper(t, server, test.method, test.path, strings.NewReader(test.body))
		defer resp.Body.Close()

		assert.Equal(t, test.code, resp.StatusCode)
	}
}