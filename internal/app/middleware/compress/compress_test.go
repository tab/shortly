package compress

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CompressMiddleware_ResponseCompression(t *testing.T) {
	client := &http.Client{
		Transport: &http.Transport{
			DisableCompression: true,
		},
	}

	type result struct {
		compress bool
		code     int
	}

	tests := []struct {
		name           string
		acceptEncoding string
		body           string
		expected       result
	}{
		{
			name:           "Sends gzip",
			acceptEncoding: "gzip",
			body:           `{"test": "data"}`,
			expected: result{
				compress: true,
				code:     http.StatusOK,
			},
		},
		{
			name:           "No gzip support",
			acceptEncoding: "",
			body:           `{"test": "data"}`,
			expected: result{
				compress: false,
				code:     http.StatusOK,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(`{"response": "ok"}`))
				assert.NoError(t, err)
			})

			ts := httptest.NewServer(Middleware(handler))
			defer ts.Close()

			req, err := http.NewRequest("POST", ts.URL, strings.NewReader(tt.body))
			req.Header.Set("Accept-Encoding", tt.acceptEncoding)
			assert.NoError(t, err)

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()
			assert.Equal(t, tt.expected.code, resp.StatusCode)

			if tt.expected.compress {
				assert.Equal(t, "gzip", resp.Header.Get("Content-Encoding"))

				gzReader, err := gzip.NewReader(resp.Body)
				assert.NoError(t, err)
				defer gzReader.Close()

				unzippedBody, err := io.ReadAll(gzReader)
				assert.NoError(t, err)
				assert.Equal(t, `{"response": "ok"}`, string(unzippedBody))
			} else {
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)
				assert.Equal(t, `{"response": "ok"}`, string(body))
			}
		})
	}
}

func Test_CompressMiddleware_RequestDecompression(t *testing.T) {
	client := &http.Client{
		Transport: &http.Transport{
			DisableCompression: true,
		},
	}

	type result struct {
		body string
		code int
	}

	tests := []struct {
		name            string
		contentEncoding string
		body            string
		expected        result
	}{
		{
			name:            "Receives gzip",
			contentEncoding: "gzip",
			body:            `{"test": "data"}`,
			expected: result{
				body: `{"test": "data"}`,
				code: http.StatusOK,
			},
		},
		{
			name:            "Receives plain",
			contentEncoding: "",
			body:            `{"test": "data"}`,
			expected: result{
				body: `{"test": "data"}`,
				code: http.StatusOK,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.body, string(body))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, err = w.Write(body)
				assert.NoError(t, err)
			})

			ts := httptest.NewServer(Middleware(handler))
			defer ts.Close()

			var requestBody io.Reader
			if tt.contentEncoding == "gzip" {
				var buf bytes.Buffer
				gz := gzip.NewWriter(&buf)
				_, err := gz.Write([]byte(tt.body))
				assert.NoError(t, err)
				err = gz.Close()
				assert.NoError(t, err)
				requestBody = &buf
			} else {
				requestBody = strings.NewReader(tt.body)
			}

			req, err := http.NewRequest("POST", ts.URL, requestBody)
			req.Header.Set("Content-Encoding", tt.contentEncoding)
			assert.NoError(t, err)

			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.body, string(body))
			assert.Equal(t, tt.expected.code, resp.StatusCode)
		})
	}
}
