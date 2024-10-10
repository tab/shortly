package api

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"shortly/internal/app/config"
	"shortly/internal/app/errors"
	"shortly/internal/app/repository"
	"shortly/internal/app/service"
)

func TestHandleCreateShortLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
	}
	repo := repository.NewMockURLRepository(ctrl)
	rand := service.NewMockSecureRandomGenerator(ctrl)
	srv := service.NewURLService(repo, rand, cfg)
	handler := NewURLHandler(cfg, srv)

	type result struct {
		status   int
		response string
	}

	tests := []struct {
		name     string
		method   string
		body     string
		before   func()
		expected result
	}{
		{
			name:   "Success",
			method: http.MethodPost,
			body:   "https://example.com",
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
			body:   "",
			before: func() {},
			expected: result{
				status:   http.StatusBadRequest,
				response: "Invalid request method",
			},
		},
		{
			name:   "Empty body",
			method: http.MethodPost,
			body:   "",
			before: func() {},
			expected: result{
				status:   http.StatusBadRequest,
				response: errors.ErrorRequestBodyEmpty.Error(),
			},
		},
		{
			name:   "Invalid URL",
			method: http.MethodPost,
			body:   "not-a-url",
			before: func() {},
			expected: result{
				status:   http.StatusBadRequest,
				response: errors.ErrorInvalidURL.Error(),
			},
		},
		{
			name:   "Error generating short code",
			method: http.MethodPost,
			body:   "https://example.com",
			before: func() {
				rand.EXPECT().Hex().Return("", errors.ErrorFailedToReadRandomBytes)
			},
			expected: result{
				status:   http.StatusInternalServerError,
				response: errors.ErrorCouldNotGenerateCode.Error(),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.before()

			req := httptest.NewRequest(test.method, "/", strings.NewReader(test.body))
			w := httptest.NewRecorder()

			handler.HandleCreateShortLink(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, test.expected.status, resp.StatusCode)
			assert.Equal(t, test.expected.response, strings.TrimSpace(w.Body.String()))
		})
	}
}

func TestHandleGetShortLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
	}
	repo := repository.NewMockURLRepository(ctrl)
	rand := service.NewMockSecureRandomGenerator(ctrl)
	srv := service.NewURLService(repo, rand, cfg)
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
				response: errors.ErrorShortLinkNotFound.Error(),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.before()

			req := httptest.NewRequest(http.MethodGet, test.path, nil)
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Get("/{id}", handler.HandleGetShortLink)
			r.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, test.expected.status, resp.StatusCode)
			if test.expected.status == http.StatusTemporaryRedirect {
				assert.Equal(t, test.expected.header, w.Header().Get("Location"))
			} else {
				assert.Equal(t, test.expected.response, strings.TrimSpace(w.Body.String()))
			}
		})
	}
}
