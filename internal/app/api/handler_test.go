package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"shortly/internal/app/config"
	"shortly/internal/app/dto"
	"shortly/internal/app/errors"
	"shortly/internal/app/repository"
	"shortly/internal/app/service"
)

func Test_HandleCreateShortLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
	}
	repo := repository.NewMockRepository(ctrl)
	rand := service.NewMockSecureRandomGenerator(ctrl)
	srv := service.NewURLService(cfg, repo, rand)
	handler := NewURLHandler(cfg, srv)

	UUID, _ := uuid.Parse("6455bd07-e431-4851-af3c-4f703f726639")

	type result struct {
		response dto.CreateShortLinkResponse
		error    dto.ErrorResponse
		code     int
		status   string
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
			body:   strings.NewReader(`{"url":"https://example.com"}`),
			before: func() {
				rand.EXPECT().UUID().Return(UUID, nil)
				rand.EXPECT().Hex().Return("abcd1234", nil)

				repo.EXPECT().CreateURL(ctx, repository.URL{
					UUID:      UUID,
					LongURL:   "https://example.com",
					ShortCode: "abcd1234",
				})
			},
			expected: result{
				response: dto.CreateShortLinkResponse{Result: "http://localhost:8080/abcd1234"},
				status:   "201 Created",
				code:     http.StatusCreated,
			},
		},
		{
			name:   "Empty body",
			method: http.MethodPost,
			body:   strings.NewReader("{}"),
			before: func() {},
			expected: result{
				error:  dto.ErrorResponse{Error: "request body is empty"},
				status: "400 Bad Request",
				code:   http.StatusBadRequest,
			},
		},
		{
			name:   "Empty URL",
			method: http.MethodPost,
			body:   strings.NewReader(`{"url":""}`),
			before: func() {},
			expected: result{
				error:  dto.ErrorResponse{Error: "request body is empty"},
				status: "400 Bad Request",
				code:   http.StatusBadRequest,
			},
		},
		{
			name:   "Invalid JSON",
			method: http.MethodPost,
			body:   strings.NewReader(`{"url"}`),
			before: func() {},
			expected: result{
				error:  dto.ErrorResponse{Error: "invalid character '}' after object key"},
				status: "400 Bad Request",
				code:   http.StatusBadRequest,
			},
		},
		{
			name:   "Invalid URL",
			method: http.MethodPost,
			body:   strings.NewReader(`{"url":"not-a-url"}`),
			before: func() {},
			expected: result{
				error:  dto.ErrorResponse{Error: "invalid URL"},
				status: "400 Bad Request",
				code:   http.StatusBadRequest,
			},
		},
		{
			name:   "Error generating UUID",
			method: http.MethodPost,
			body:   strings.NewReader(`{"url":"https://example.com"}`),
			before: func() {
				rand.EXPECT().UUID().Return(uuid.UUID{}, errors.ErrFailedToGenerateUUID)
			},
			expected: result{
				error:  dto.ErrorResponse{Error: "failed to generate UUID"},
				status: "500 Internal Server Error",
				code:   http.StatusInternalServerError,
			},
		},
		{
			name:   "Error generating short code",
			method: http.MethodPost,
			body:   strings.NewReader(`{"url":"https://example.com"}`),
			before: func() {
				rand.EXPECT().UUID().Return(UUID, nil)
				rand.EXPECT().Hex().Return("", errors.ErrFailedToReadRandomBytes)
			},
			expected: result{
				error:  dto.ErrorResponse{Error: "failed to generate short code"},
				status: "500 Internal Server Error",
				code:   http.StatusInternalServerError,
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

			if tt.expected.error.Error != "" {
				var actual dto.ErrorResponse
				err := json.NewDecoder(resp.Body).Decode(&actual)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error.Error, actual.Error)
			} else {
				var actual dto.CreateShortLinkResponse
				err := json.NewDecoder(resp.Body).Decode(&actual)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response.Result, actual.Result)
			}
			assert.Equal(t, tt.expected.status, resp.Status)
			assert.Equal(t, tt.expected.code, resp.StatusCode)
		})
	}
}

func Test_HandleGetShortLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
	}
	repo := repository.NewMockRepository(ctrl)
	rand := service.NewMockSecureRandomGenerator(ctrl)
	srv := service.NewURLService(cfg, repo, rand)
	handler := NewURLHandler(cfg, srv)

	type result struct {
		response dto.CreateShortLinkResponse
		error    dto.ErrorResponse
		code     int
		status   string
	}

	tests := []struct {
		name     string
		path     string
		before   func()
		expected result
	}{
		{
			name: "Success",
			path: "/api/shorten/abcd1234",
			before: func() {
				repo.EXPECT().GetURLByShortCode(ctx, "abcd1234").Return(&repository.URL{
					LongURL:   "https://example.com",
					ShortCode: "abcd1234",
				}, true)
			},
			expected: result{
				response: dto.CreateShortLinkResponse{Result: "https://example.com"},
				status:   "200 OK",
				code:     http.StatusOK,
			},
		},
		{
			name: "Not Found",
			path: "/api/shorten/not-a-short-code",
			before: func() {
				repo.EXPECT().GetURLByShortCode(ctx, "not-a-short-code").Return(nil, false)
			},
			expected: result{
				error:  dto.ErrorResponse{Error: errors.ErrShortLinkNotFound.Error()},
				status: "404 Not Found",
				code:   http.StatusNotFound,
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

			if tt.expected.error.Error != "" {
				var actual dto.ErrorResponse
				err := json.NewDecoder(resp.Body).Decode(&actual)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error.Error, actual.Error)
			} else {
				var actual dto.CreateShortLinkResponse
				err := json.NewDecoder(resp.Body).Decode(&actual)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response.Result, actual.Result)
			}
			assert.Equal(t, tt.expected.status, resp.Status)
			assert.Equal(t, tt.expected.code, resp.StatusCode)
		})
	}
}

func Test_DeprecatedHandleCreateShortLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
	}
	repo := repository.NewMockRepository(ctrl)
	rand := service.NewMockSecureRandomGenerator(ctrl)
	srv := service.NewURLService(cfg, repo, rand)
	handler := NewURLHandler(cfg, srv)

	UUID, _ := uuid.Parse("6455bd07-e431-4851-af3c-4f703f726639")

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
				rand.EXPECT().UUID().Return(UUID, nil)
				rand.EXPECT().Hex().Return("abcd1234", nil)

				repo.EXPECT().CreateURL(ctx, repository.URL{
					UUID:      UUID,
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
			name:   "Error generating UUID",
			method: http.MethodPost,
			body:   strings.NewReader("https://example.com"),
			before: func() {
				rand.EXPECT().UUID().Return(uuid.UUID{}, errors.ErrFailedToGenerateUUID)
			},
			expected: result{
				status:   http.StatusInternalServerError,
				response: errors.ErrFailedToGenerateUUID.Error(),
			},
		},
		{
			name:   "Error generating short code",
			method: http.MethodPost,
			body:   strings.NewReader("https://example.com"),
			before: func() {
				rand.EXPECT().UUID().Return(UUID, nil)
				rand.EXPECT().Hex().Return("", errors.ErrFailedToReadRandomBytes)
			},
			expected: result{
				status:   http.StatusInternalServerError,
				response: errors.ErrFailedToGenerateCode.Error(),
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

	ctx := context.Background()
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
	}
	repo := repository.NewMockRepository(ctrl)
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
				repo.EXPECT().GetURLByShortCode(ctx, "abcd1234").Return(&repository.URL{
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
				repo.EXPECT().GetURLByShortCode(ctx, "not-a-short-code").Return(nil, false)
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
