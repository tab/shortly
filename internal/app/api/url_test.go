package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"shortly/internal/app/config"
	"shortly/internal/app/dto"
	"shortly/internal/app/errors"
	"shortly/internal/app/repository"
	"shortly/internal/app/service"
	"shortly/internal/app/worker"
)

func Test_HandleCreateShortLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := gomock.Any()
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
	}
	repo := repository.NewMockDatabase(ctrl)
	rand := service.NewMockSecureRandomGenerator(ctrl)
	appWorker := worker.NewMockWorker(ctrl)
	srv := service.NewURLService(cfg, repo, rand, appWorker)
	handler := NewURLHandler(cfg, srv)

	UUID := uuid.MustParse("6455bd07-e431-4851-af3c-4f703f726639")

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
				repo.EXPECT().GetURLByShortCode(ctx, "abcd1234").Return(nil, false)

				repo.EXPECT().CreateURL(ctx, repository.URL{
					UUID:      UUID,
					LongURL:   "https://example.com",
					ShortCode: "abcd1234",
				}).Return(&repository.URL{
					UUID:      UUID,
					LongURL:   "https://example.com",
					ShortCode: "abcd1234",
				}, nil)
			},
			expected: result{
				response: dto.CreateShortLinkResponse{Result: "http://localhost:8080/abcd1234"},
				status:   "201 Created",
				code:     http.StatusCreated,
			},
		},
		{
			name:   "URL already exists",
			method: http.MethodPost,
			body:   strings.NewReader(`{"url":"https://example.com"}`),
			before: func() {
				rand.EXPECT().UUID().Return(UUID, nil)
				rand.EXPECT().Hex().Return("abcd1234", nil)
				repo.EXPECT().GetURLByShortCode(ctx, "abcd1234").Return(nil, false)

				repo.EXPECT().CreateURL(ctx, repository.URL{
					UUID:      UUID,
					LongURL:   "https://example.com",
					ShortCode: "abcd1234",
				}).Return(&repository.URL{
					UUID:      UUID,
					LongURL:   "https://example.com",
					ShortCode: "abab0001",
				}, nil)
			},
			expected: result{
				response: dto.CreateShortLinkResponse{Result: "http://localhost:8080/abab0001"},
				status:   "409 Conflict",
				code:     http.StatusConflict,
			},
		},
		{
			name:   "Empty body",
			method: http.MethodPost,
			body:   strings.NewReader("{}"),
			before: func() {},
			expected: result{
				error:  dto.ErrorResponse{Error: "original URL is required"},
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
				error:  dto.ErrorResponse{Error: "original URL is required"},
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

func Benchmark_HandleCreateShortLink(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
	}
	repo := repository.NewMockDatabase(ctrl)
	rand := service.NewMockSecureRandomGenerator(ctrl)
	appWorker := worker.NewMockWorker(ctrl)
	srv := service.NewURLService(cfg, repo, rand, appWorker)
	handler := NewURLHandler(cfg, srv)

	rand.EXPECT().UUID().Return(uuid.Must(uuid.NewRandom()), nil).AnyTimes()
	rand.EXPECT().Hex().Return("abcd1234", nil).AnyTimes()
	repo.EXPECT().GetURLByShortCode(gomock.Any(), "abcd1234").Return(nil, false).AnyTimes()
	repo.EXPECT().CreateURL(gomock.Any(), gomock.Any()).Return(&repository.URL{
		LongURL:   "https://example.com",
		ShortCode: "abcd1234",
	}, nil).AnyTimes()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(`{"url":"https://example.com"}`))
		w := httptest.NewRecorder()

		handler.HandleCreateShortLink(w, req)
	}
}

func Test_HandleBatchCreateShortLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := gomock.Any()
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
	}
	repo := repository.NewMockDatabase(ctrl)
	rand := service.NewMockSecureRandomGenerator(ctrl)
	appWorker := worker.NewMockWorker(ctrl)
	srv := service.NewURLService(cfg, repo, rand, appWorker)
	handler := NewURLHandler(cfg, srv)

	UUID1, _ := uuid.Parse("6455bd07-e431-4851-af3c-4f703f720001")
	UUID2, _ := uuid.Parse("6455bd07-e431-4851-af3c-4f703f720002")

	type result struct {
		response dto.BatchCreateShortLinkResponses
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
			body: strings.NewReader(`[
				{"correlation_id": "0001", "original_url": "https://github.com"},
				{"correlation_id": "0002", "original_url": "https://google.com"}
			]`),
			before: func() {
				rand.EXPECT().UUID().Return(UUID1, nil)
				rand.EXPECT().UUID().Return(UUID2, nil)
				rand.EXPECT().Hex().Return("abcd0001", nil)
				rand.EXPECT().Hex().Return("abcd0002", nil)
				repo.EXPECT().GetURLByShortCode(ctx, "abcd0001").Return(nil, false)
				repo.EXPECT().GetURLByShortCode(ctx, "abcd0002").Return(nil, false)

				urls := []repository.URL{
					{
						UUID:      UUID1,
						LongURL:   "https://github.com",
						ShortCode: "abcd0001",
					},
					{
						UUID:      UUID2,
						LongURL:   "https://google.com",
						ShortCode: "abcd0002",
					},
				}
				repo.EXPECT().CreateURLs(ctx, urls)
			},
			expected: result{
				response: dto.BatchCreateShortLinkResponses{
					{CorrelationID: "0001", ShortURL: "http://localhost:8080/abcd0001"},
					{CorrelationID: "0002", ShortURL: "http://localhost:8080/abcd0002"},
				},
				status: "201 Created",
				code:   http.StatusCreated,
			},
		},
		{
			name:   "No correlation ID",
			method: http.MethodPost,
			body:   strings.NewReader(`[{"original_url": "not-a-url"}]`),
			before: func() {},
			expected: result{
				error:  dto.ErrorResponse{Error: "correlation id is required"},
				status: "400 Bad Request",
				code:   http.StatusBadRequest,
			},
		},
		{
			name:   "No original URL",
			method: http.MethodPost,
			body:   strings.NewReader(`[{"correlation_id": "0001"}]`),
			before: func() {},
			expected: result{
				error:  dto.ErrorResponse{Error: "invalid URL"},
				status: "400 Bad Request",
				code:   http.StatusBadRequest,
			},
		},
		{
			name:   "Invalid URL",
			method: http.MethodPost,
			body:   strings.NewReader(`[{"correlation_id": "0001", "original_url": "not-a-url"}]`),
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
			body:   strings.NewReader(`[{"correlation_id": "0001", "original_url": "https://github.com"}]`),
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
			body:   strings.NewReader(`[{"correlation_id": "0001", "original_url": "https://github.com"}]`),
			before: func() {
				rand.EXPECT().UUID().Return(UUID1, nil)
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

			req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", tt.body)
			w := httptest.NewRecorder()

			handler.HandleBatchCreateShortLink(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if tt.expected.error.Error != "" {
				var actual dto.ErrorResponse
				err := json.NewDecoder(resp.Body).Decode(&actual)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error.Error, actual.Error)
			} else {
				var actual dto.BatchCreateShortLinkResponses
				err := json.NewDecoder(resp.Body).Decode(&actual)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response, actual)
			}
			assert.Equal(t, tt.expected.status, resp.Status)
			assert.Equal(t, tt.expected.code, resp.StatusCode)
		})
	}
}

func Test_HandleGetShortLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := gomock.Any()
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
	}
	repo := repository.NewMockDatabase(ctrl)
	rand := service.NewMockSecureRandomGenerator(ctrl)
	appWorker := worker.NewMockWorker(ctrl)
	srv := service.NewURLService(cfg, repo, rand, appWorker)
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
			name: "Deleted",
			path: "/api/shorten/abcd1234",
			before: func() {
				repo.EXPECT().GetURLByShortCode(ctx, "abcd1234").Return(&repository.URL{
					LongURL:   "https://example.com",
					ShortCode: "abcd1234",
					DeletedAt: time.Now(),
				}, true)
			},
			expected: result{
				error:  dto.ErrorResponse{Error: errors.ErrShortLinkDeleted.Error()},
				status: "410 Gone",
				code:   http.StatusGone,
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

func Benchmark_HandleGetShortLink(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	ctx := gomock.Any()
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
	}
	repo := repository.NewMockDatabase(ctrl)
	rand := service.NewMockSecureRandomGenerator(ctrl)
	appWorker := worker.NewMockWorker(ctrl)
	srv := service.NewURLService(cfg, repo, rand, appWorker)
	handler := NewURLHandler(cfg, srv)

	repo.EXPECT().GetURLByShortCode(ctx, "abcd1234").Return(&repository.URL{
		LongURL:   "https://example.com",
		ShortCode: "abcd1234",
	}, true).AnyTimes()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/shorten/abcd1234", nil)
		w := httptest.NewRecorder()

		r := chi.NewRouter()
		r.Get("/api/shorten/{id}", handler.HandleGetShortLink)
		r.ServeHTTP(w, req)
	}
}

func Test_HandleGetUserURLs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
	}
	repo := repository.NewMockDatabase(ctrl)
	rand := service.NewMockSecureRandomGenerator(ctrl)
	appWorker := worker.NewMockWorker(ctrl)
	srv := service.NewURLService(cfg, repo, rand, appWorker)
	handler := NewURLHandler(cfg, srv)

	limit := int64(25)
	offset := int64(0)

	UUID1, _ := uuid.Parse("6455bd07-e431-4851-af3c-4f703f720001")
	UUID2, _ := uuid.Parse("6455bd07-e431-4851-af3c-4f703f720002")
	UserUUID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174001")

	ctx := context.WithValue(context.Background(), dto.CurrentUser, UserUUID)

	type result struct {
		response []dto.GetUserURLsResponse
		error    dto.ErrorResponse
		code     int
		status   string
	}

	tests := []struct {
		name     string
		before   func()
		expected result
	}{
		{
			name: "Success",
			before: func() {
				repo.EXPECT().GetURLsByUserID(ctx, UserUUID, limit, offset).Return([]repository.URL{
					{
						UUID:      UUID1,
						LongURL:   "https://google.com",
						ShortCode: "abcd0001",
					},
					{
						UUID:      UUID2,
						LongURL:   "https://github.com",
						ShortCode: "abcd0002",
					},
				}, 2, nil)
			},
			expected: result{
				response: []dto.GetUserURLsResponse{
					{
						ShortURL:    "http://localhost:8080/abcd0001",
						OriginalURL: "https://google.com",
					},
					{
						ShortURL:    "http://localhost:8080/abcd0002",
						OriginalURL: "https://github.com",
					},
				},
				status: "200 OK",
				code:   http.StatusOK,
			},
		},
		{
			name: "No URLs",
			before: func() {
				repo.EXPECT().GetURLsByUserID(ctx, UserUUID, limit, offset).Return(nil, 0, nil)
			},
			expected: result{
				response: []dto.GetUserURLsResponse(nil),
				status:   "204 No Content",
				code:     http.StatusNoContent,
			},
		},
		{
			name: "Error",
			before: func() {
				repo.EXPECT().GetURLsByUserID(ctx, UserUUID, limit, offset).Return(nil, 0, errors.ErrFailedToLoadUserUrls)
			},
			expected: result{
				error:  dto.ErrorResponse{Error: errors.ErrFailedToLoadUserUrls.Error()},
				status: "500 Internal Server Error",
				code:   http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.HandleGetUserURLs(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			switch tt.expected.code {
			case http.StatusNoContent:
				body, _ := io.ReadAll(resp.Body)
				assert.Empty(t, body)

			case http.StatusInternalServerError, http.StatusBadRequest, http.StatusNotFound:
				var actual dto.ErrorResponse
				err := json.NewDecoder(resp.Body).Decode(&actual)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error.Error, actual.Error)

			default:
				var actual []dto.GetUserURLsResponse
				err := json.NewDecoder(resp.Body).Decode(&actual)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response, actual)
			}

			assert.Equal(t, tt.expected.status, resp.Status)
			assert.Equal(t, tt.expected.code, resp.StatusCode)
		})
	}
}

func Test_HandleBatchDeleteUserURLs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
	}
	repo := repository.NewMockDatabase(ctrl)
	rand := service.NewMockSecureRandomGenerator(ctrl)
	appWorker := worker.NewMockWorker(ctrl)
	srv := service.NewURLService(cfg, repo, rand, appWorker)
	handler := NewURLHandler(cfg, srv)

	UserUUID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174001")

	type result struct {
		error  dto.ErrorResponse
		code   int
		status string
	}

	tests := []struct {
		name     string
		ctx      context.Context
		body     io.Reader
		before   func()
		expected result
	}{
		{
			name: "Success",
			ctx:  context.WithValue(context.Background(), dto.CurrentUser, UserUUID),
			body: strings.NewReader(`["abcd0001", "abcd0002"]`),
			before: func() {
				appWorker.EXPECT().Add(dto.BatchDeleteParams{
					UserID:     UserUUID,
					ShortCodes: []string{"abcd0001", "abcd0002"},
				})
			},
			expected: result{
				status: "202 Accepted",
				code:   http.StatusAccepted,
			},
		},
		{
			name:   "Empty",
			ctx:    context.WithValue(context.Background(), dto.CurrentUser, UserUUID),
			body:   strings.NewReader("[]"),
			before: func() {},
			expected: result{
				error:  dto.ErrorResponse{Error: errors.ErrShortCodeEmpty.Error()},
				status: "400 Bad Request",
				code:   http.StatusBadRequest,
			},
		},
		{
			name: "Error",
			ctx:  context.WithValue(context.Background(), dto.CurrentUser, nil),
			body: strings.NewReader(`["abcd0001", "abcd0002"]`),
			before: func() {
			},
			expected: result{
				error:  dto.ErrorResponse{Error: errors.ErrInvalidUserID.Error()},
				status: "500 Internal Server Error",
				code:   http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", tt.body)
			req = req.WithContext(tt.ctx)
			w := httptest.NewRecorder()

			handler.HandleBatchDeleteUserURLs(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			switch tt.expected.code {
			case http.StatusInternalServerError, http.StatusBadRequest, http.StatusNotFound:
				var actual dto.ErrorResponse
				err := json.NewDecoder(resp.Body).Decode(&actual)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.error.Error, actual.Error)

			default:
				assert.Equal(t, tt.expected.status, resp.Status)
				assert.Equal(t, tt.expected.code, resp.StatusCode)
			}
		})
	}
}

func Test_DeprecatedHandleCreateShortLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := gomock.Any()
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
	}
	repo := repository.NewMockDatabase(ctrl)
	rand := service.NewMockSecureRandomGenerator(ctrl)
	appWorker := worker.NewMockWorker(ctrl)
	srv := service.NewURLService(cfg, repo, rand, appWorker)
	handler := NewURLHandler(cfg, srv)

	UUID := uuid.MustParse("6455bd07-e431-4851-af3c-4f703f726639")

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
				repo.EXPECT().GetURLByShortCode(ctx, "abcd1234").Return(nil, false)

				repo.EXPECT().CreateURL(ctx, repository.URL{
					UUID:      UUID,
					LongURL:   "https://example.com",
					ShortCode: "abcd1234",
				}).Return(&repository.URL{
					UUID:      UUID,
					LongURL:   "https://example.com",
					ShortCode: "abcd1234",
				}, nil)
			},
			expected: result{
				status:   http.StatusCreated,
				response: "http://localhost:8080/abcd1234",
			},
		},
		{
			name:   "URL already exists",
			method: http.MethodPost,
			body:   strings.NewReader("https://example.com"),
			before: func() {
				rand.EXPECT().UUID().Return(UUID, nil)
				rand.EXPECT().Hex().Return("abcd1234", nil)
				repo.EXPECT().GetURLByShortCode(ctx, "abcd1234").Return(nil, false)

				repo.EXPECT().CreateURL(ctx, repository.URL{
					UUID:      UUID,
					LongURL:   "https://example.com",
					ShortCode: "abcd1234",
				}).Return(&repository.URL{
					UUID:      UUID,
					LongURL:   "https://example.com",
					ShortCode: "abab0001",
				}, nil)
			},
			expected: result{
				status:   http.StatusConflict,
				response: "http://localhost:8080/abab0001",
			},
		},
		{
			name:   "Empty body",
			method: http.MethodPost,
			body:   strings.NewReader(""),
			before: func() {},
			expected: result{
				status:   http.StatusBadRequest,
				response: errors.ErrOriginalURLEmpty.Error(),
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

	ctx := gomock.Any()
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
	}
	repo := repository.NewMockDatabase(ctrl)
	rand := service.NewMockSecureRandomGenerator(ctrl)
	appWorker := worker.NewMockWorker(ctrl)
	srv := service.NewURLService(cfg, repo, rand, appWorker)
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
			name: "Deleted",
			path: "/abcd1234",
			before: func() {
				repo.EXPECT().GetURLByShortCode(ctx, "abcd1234").Return(&repository.URL{
					LongURL:   "https://example.com",
					ShortCode: "abcd1234",
					DeletedAt: time.Now(),
				}, true)
			},
			expected: result{
				status:   http.StatusGone,
				response: errors.ErrShortLinkDeleted.Error(),
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
