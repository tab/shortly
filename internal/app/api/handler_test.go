package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"shortly/internal/app/config"
	"shortly/internal/app/errors"
	"shortly/internal/app/repository"
	"shortly/internal/app/service"
)

func Test_HandleCreateShortLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
	}
	repo := repository.NewMockURLRepository(ctrl)
	rand := service.NewMockSecureRandomGenerator(ctrl)
	srv := service.NewURLService(cfg, repo, rand)
	handler := NewURLHandler(cfg, srv)

	tests := []struct {
		name     string
		method   string
		body     io.Reader
		before   func()
		expected Response
	}{
		{
			name:   "Success",
			method: http.MethodPost,
			body:   strings.NewReader(`{"url":"https://example.com"}`),
			before: func() {
				rand.EXPECT().Hex().Return("abcd1234", nil)
				repo.EXPECT().Set(repository.URL{
					LongURL:   "https://example.com",
					ShortCode: "abcd1234",
				})
			},
			expected: Response{
				Result: "http://localhost:8080/abcd1234",
				Status: http.StatusText(http.StatusCreated),
				Code:   http.StatusCreated,
			},
		},
		{
			name:   "Invalid request method",
			method: http.MethodGet,
			body:   strings.NewReader(`{"url":"https://example.com"}`),
			before: func() {},
			expected: Response{
				Error:  "Invalid request method",
				Status: http.StatusText(http.StatusBadRequest),
				Code:   http.StatusBadRequest,
			},
		},
		{
			name:   "Empty body",
			method: http.MethodPost,
			body:   strings.NewReader(`{"url":""}`),
			before: func() {},
			expected: Response{
				Error:  "Request body is empty",
				Status: http.StatusText(http.StatusBadRequest),
				Code:   http.StatusBadRequest,
			},
		},
		{
			name:   "Invalid JSON",
			method: http.MethodPost,
			body:   strings.NewReader(`{"url"}`),
			before: func() {},
			expected: Response{
				Error:  "Invalid JSON",
				Status: http.StatusText(http.StatusInternalServerError),
				Code:   http.StatusInternalServerError,
			},
		},
		{
			name:   "Invalid URL",
			method: http.MethodPost,
			body:   strings.NewReader(`{"url":"not-a-url"}`),
			before: func() {},
			expected: Response{
				Error:  "Invalid URL",
				Status: http.StatusText(http.StatusBadRequest),
				Code:   http.StatusBadRequest,
			},
		},
		{
			name:   "Error generating short code",
			method: http.MethodPost,
			body:   strings.NewReader(`{"url":"https://example.com"}`),
			before: func() {
				rand.EXPECT().Hex().Return("", errors.ErrFailedToReadRandomBytes)
			},
			expected: Response{
				Error:  "Could not generate short code",
				Status: http.StatusText(http.StatusInternalServerError),
				Code:   http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			req := httptest.NewRequest(tt.method, "/api/shorten", tt.body)
			w := httptest.NewRecorder()

			handler.HandleCreateShortLink(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			var actual Response
			err := json.NewDecoder(resp.Body).Decode(&actual)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Result, actual.Result)
			assert.Equal(t, tt.expected.Status, actual.Status)
			assert.Equal(t, tt.expected.Code, actual.Code)
		})
	}
}

func Test_HandleGetShortLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
	}
	repo := repository.NewMockURLRepository(ctrl)
	rand := service.NewMockSecureRandomGenerator(ctrl)
	srv := service.NewURLService(cfg, repo, rand)
	handler := NewURLHandler(cfg, srv)

	tests := []struct {
		name     string
		path     string
		before   func()
		expected Response
	}{
		{
			name: "Success",
			path: "/api/shorten/abcd1234",
			before: func() {
				repo.EXPECT().Get("abcd1234").Return(&repository.URL{
					LongURL:   "https://example.com",
					ShortCode: "abcd1234",
				}, true)
			},
			expected: Response{
				Result: "https://example.com",
				Status: http.StatusText(http.StatusOK),
				Code:   http.StatusOK,
			},
		},
		{
			name: "Not Found",
			path: "/api/shorten/not-a-short-code",
			before: func() {
				repo.EXPECT().Get("not-a-short-code").Return(nil, false)
			},
			expected: Response{
				Error:  errors.ErrShortLinkNotFound.Error(),
				Status: http.StatusText(http.StatusNotFound),
				Code:   http.StatusNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Get("/api/shorten/{id}", handler.HandleGetShortLink)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			var actual Response
			err := json.NewDecoder(resp.Body).Decode(&actual)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Result, actual.Result)
			assert.Equal(t, tt.expected.Status, actual.Status)
			assert.Equal(t, tt.expected.Code, actual.Code)
		})
	}
}

func Test_DeprecatedHandleCreateShortLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
	}
	repo := repository.NewMockURLRepository(ctrl)
	rand := service.NewMockSecureRandomGenerator(ctrl)
	srv := service.NewURLService(cfg, repo, rand)
	handler := NewURLHandler(cfg, srv)

	type result struct {
		status   int
		response string
	}

	tests := []struct {
		name     string
		method   string
		body     io.Reader
		before   func()
		expected result
	}{
		{
			name:   "Success",
			method: http.MethodPost,
			body:   strings.NewReader("https://example.com"),
			before: func() {
				rand.EXPECT().Hex().Return("abcd1234", nil)
				repo.EXPECT().Set(repository.URL{
					LongURL:   "https://example.com",
					ShortCode: "abcd1234",
				})
			},
			expected: result{
				status:   http.StatusCreated,
				response: "http://localhost:8080/abcd1234",
			},
		},
		{
			name:   "Invalid request method",
			method: http.MethodGet,
			body:   strings.NewReader("https://example.com"),
			before: func() {},
			expected: result{
				status:   http.StatusBadRequest,
				response: "Invalid request method",
			},
		},
		{
			name:   "Empty body",
			method: http.MethodPost,
			body:   strings.NewReader(""),
			before: func() {},
			expected: result{
				status:   http.StatusBadRequest,
				response: errors.ErrRequestBodyEmpty.Error(),
			},
		},
		{
			name:   "Invalid URL",
			method: http.MethodPost,
			body:   strings.NewReader("not-a-url"),
			before: func() {},
			expected: result{
				status:   http.StatusBadRequest,
				response: errors.ErrInvalidURL.Error(),
			},
		},
		{
			name:   "Error generating short code",
			method: http.MethodPost,
			body:   strings.NewReader("https://example.com"),
			before: func() {
				rand.EXPECT().Hex().Return("", errors.ErrFailedToReadRandomBytes)
			},
			expected: result{
				status:   http.StatusInternalServerError,
				response: errors.ErrCouldNotGenerateCode.Error(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			req := httptest.NewRequest(tt.method, "/", tt.body)
			w := httptest.NewRecorder()

			handler.DeprecatedHandleCreateShortLink(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.expected.status, resp.StatusCode)
			assert.Equal(t, tt.expected.response, strings.TrimSpace(w.Body.String()))
		})
	}
}

func Test_DeprecatedHandleGetShortLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
	}
	repo := repository.NewMockURLRepository(ctrl)
	rand := service.NewMockSecureRandomGenerator(ctrl)
	srv := service.NewURLService(cfg, repo, rand)
	handler := NewURLHandler(cfg, srv)

	type result struct {
		status   int
		header   string
		response string
	}

	tests := []struct {
		name     string
		path     string
		before   func()
		expected result
	}{
		{
			name: "Success",
			path: "/abcd1234",
			before: func() {
				repo.EXPECT().Get("abcd1234").Return(&repository.URL{
					LongURL:   "https://example.com",
					ShortCode: "abcd1234",
				}, true)
			},
			expected: result{
				status: http.StatusTemporaryRedirect,
				header: "https://example.com",
			},
		},
		{
			name: "Not Found",
			path: "/not-a-short-code",
			before: func() {
				repo.EXPECT().Get("not-a-short-code").Return(nil, false)
			},
			expected: result{
				status:   http.StatusNotFound,
				response: errors.ErrShortLinkNotFound.Error(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Get("/{id}", handler.DeprecatedHandleGetShortLink)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.expected.status, resp.StatusCode)
			if tt.expected.status == http.StatusTemporaryRedirect {
				assert.Equal(t, tt.expected.header, w.Header().Get("Location"))
			} else {
				assert.Equal(t, tt.expected.response, strings.TrimSpace(w.Body.String()))
			}
		})
	}
}
