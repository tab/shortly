package service

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"shortly/internal/app/api/pagination"
	"shortly/internal/app/config"
	"shortly/internal/app/dto"
	"shortly/internal/app/errors"
	"shortly/internal/app/repository"
	"shortly/internal/app/worker"
)

func Test_CreateShortLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Addr:      "localhost:8080",
		BaseURL:   "http://localhost:8080",
		ClientURL: "http://localhost:8080",
	}

	ctx := context.Background()
	repo := repository.NewMockRepository(ctrl)
	rand := NewMockSecureRandomGenerator(ctrl)
	appWorker := worker.NewMockWorker(ctrl)
	service := NewURLService(cfg, repo, rand, appWorker)

	UUID1, _ := uuid.Parse("6455bd07-e431-4851-af3c-4f703f720001")
	UUID2, _ := uuid.Parse("6455bd07-e431-4851-af3c-4f703f720002")

	type result struct {
		shortCode string
		shortURL  string
		error     error
	}

	tests := []struct {
		name     string
		body     io.Reader
		before   func()
		expected result
	}{
		{
			name: "Success",
			body: strings.NewReader(`{"url":"https://example.com"}`),
			before: func() {
				rand.EXPECT().UUID().Return(UUID1, nil)
				rand.EXPECT().Hex().Return("abcd1234", nil)
				repo.EXPECT().GetURLByShortCode(ctx, "abcd1234").Return(nil, false)

				url := repository.URL{
					UUID:      UUID1,
					LongURL:   "https://example.com",
					ShortCode: "abcd1234",
				}
				repo.EXPECT().CreateURL(ctx, url).Return(&url, nil)
			},
			expected: result{
				shortCode: "abcd1234",
				shortURL:  "http://localhost:8080/abcd1234",
				error:     nil,
			},
		},
		{
			name: "URL already exists",
			body: strings.NewReader(`{"url":"https://example.com"}`),
			before: func() {
				rand.EXPECT().UUID().Return(UUID2, nil)
				rand.EXPECT().Hex().Return("abcd1234", nil)
				repo.EXPECT().GetURLByShortCode(ctx, "abcd1234").Return(nil, false)

				existingURL := repository.URL{
					UUID:      UUID1,
					LongURL:   "https://example.com",
					ShortCode: "abab0001",
				}

				repo.EXPECT().CreateURL(ctx, repository.URL{
					UUID:      UUID2,
					LongURL:   "https://example.com",
					ShortCode: "abcd1234",
				}).Return(&existingURL, nil)
			},
			expected: result{
				shortCode: "abab0001",
				shortURL:  "http://localhost:8080/abab0001",
				error:     errors.ErrURLAlreadyExists,
			},
		},
		{
			name: "Error generating UUID",
			body: strings.NewReader(`{"url":"https://example.com"}`),
			before: func() {
				rand.EXPECT().UUID().Return(uuid.UUID{}, errors.ErrFailedToGenerateUUID)
			},
			expected: result{
				shortCode: "",
				shortURL:  "",
				error:     errors.ErrFailedToGenerateUUID,
			},
		},
		{
			name: "Error generating short code",
			body: strings.NewReader(`{"url":"https://example.com"}`),
			before: func() {
				rand.EXPECT().UUID().Return(UUID1, nil)
				rand.EXPECT().Hex().Return("", errors.ErrFailedToReadRandomBytes)
			},
			expected: result{
				shortCode: "",
				shortURL:  "",
				error:     errors.ErrFailedToGenerateCode,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			r, _ := http.NewRequest(http.MethodPost, "/", tt.body)
			var req dto.CreateShortLinkRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			assert.NoError(t, err)

			shortURL, err := service.CreateShortLink(ctx, req.URL)

			assert.Equal(t, tt.expected.shortURL, shortURL)
			assert.Equal(t, tt.expected.error, err)
		})
	}
}

func TestURLService_CreateShortLinks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Addr:      "localhost:8080",
		BaseURL:   "http://localhost:8080",
		ClientURL: "http://localhost:8080",
	}

	ctx := context.Background()
	repo := repository.NewMockRepository(ctrl)
	rand := NewMockSecureRandomGenerator(ctrl)
	appWorker := worker.NewMockWorker(ctrl)
	service := NewURLService(cfg, repo, rand, appWorker)

	UUID1, _ := uuid.Parse("6455bd07-e431-4851-af3c-4f703f720001")
	UUID2, _ := uuid.Parse("6455bd07-e431-4851-af3c-4f703f720002")

	type result struct {
		shortCode string
		shortURL  string
		error     error
	}

	tests := []struct {
		name     string
		body     io.Reader
		before   func()
		expected []result
	}{
		{
			name: "Success",
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
			expected: []result{
				{
					shortCode: "abcd0001",
					shortURL:  "http://localhost:8080/abcd0001",
					error:     nil,
				},
				{
					shortCode: "abcd0002",
					shortURL:  "http://localhost:8080/abcd0002",
					error:     nil,
				},
			},
		},
		{
			name: "Error generating UUID",
			body: strings.NewReader(`[{"correlation_id": "0001", "original_url": "https://github.com"}]`),
			before: func() {
				rand.EXPECT().UUID().Return(uuid.UUID{}, errors.ErrFailedToGenerateUUID)
			},
			expected: []result{
				{
					shortCode: "",
					shortURL:  "",
					error:     errors.ErrFailedToGenerateUUID,
				},
			},
		},
		{
			name: "Error generating short code",
			body: strings.NewReader(`[{"correlation_id": "0001", "original_url": "https://github.com"}]`),
			before: func() {
				rand.EXPECT().UUID().Return(UUID1, nil)
				rand.EXPECT().Hex().Return("", errors.ErrFailedToReadRandomBytes)
			},
			expected: []result{
				{
					shortCode: "",
					shortURL:  "",
					error:     errors.ErrFailedToGenerateCode,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			r, _ := http.NewRequest(http.MethodPost, "/", tt.body)
			var req dto.BatchCreateShortLinkRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			assert.NoError(t, err)

			shortURLs, err := service.CreateShortLinks(ctx, req)

			for i, shortURL := range shortURLs {
				assert.Equal(t, tt.expected[i].shortURL, shortURL.ShortURL)
				assert.Equal(t, tt.expected[i].error, err)
			}
		})
	}
}

func Test_GetShortLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Addr:      "localhost:8080",
		BaseURL:   "http://localhost:8080",
		ClientURL: "http://localhost:8080",
	}

	ctx := context.Background()
	repo := repository.NewMockRepository(ctrl)
	rand := NewMockSecureRandomGenerator(ctrl)
	appWorker := worker.NewMockWorker(ctrl)
	service := NewURLService(cfg, repo, rand, appWorker)

	type result struct {
		url   *repository.URL
		found bool
	}

	tests := []struct {
		name      string
		shortCode string
		expected  result
	}{
		{
			name:      "Success",
			shortCode: "abcd1234",
			expected: result{
				url: &repository.URL{
					LongURL:   "https://example.com",
					ShortCode: "abcd1234",
				},
				found: true,
			},
		},
		{
			name:      "Not Found",
			shortCode: "1234abcd",
			expected: result{
				url:   nil,
				found: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.EXPECT().GetURLByShortCode(ctx, tt.shortCode).Return(tt.expected.url, tt.expected.found)

			url, found := service.GetShortLink(ctx, tt.shortCode)

			assert.Equal(t, tt.expected.url, url)
			assert.Equal(t, tt.expected.found, found)
		})
	}
}

func Test_GetUserURLs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Addr:      "localhost:8080",
		BaseURL:   "http://localhost:8080",
		ClientURL: "http://localhost:8080",
	}

	repo := repository.NewMockRepository(ctrl)
	rand := NewMockSecureRandomGenerator(ctrl)
	appWorker := worker.NewMockWorker(ctrl)
	service := NewURLService(cfg, repo, rand, appWorker)
	paginator := pagination.Pagination{
		Page: 1,
		Per:  25,
	}
	limit := int64(25)
	offset := int64(0)

	UUID1, _ := uuid.Parse("6455bd07-e431-4851-af3c-4f703f720001")
	UUID2, _ := uuid.Parse("6455bd07-e431-4851-af3c-4f703f720002")
	UserUUID, _ := uuid.Parse("123e4567-e89b-12d3-a456-426614174001")

	type result struct {
		urls  []dto.GetUserURLsResponse
		total int
		error error
	}

	tests := []struct {
		name     string
		ctx      context.Context
		before   func(ctx context.Context)
		expected result
	}{
		{
			name: "Success",
			ctx:  context.WithValue(context.Background(), dto.CurrentUser, UserUUID),
			before: func(ctx context.Context) {
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
				urls: []dto.GetUserURLsResponse{
					{
						ShortURL:    "http://localhost:8080/abcd0001",
						OriginalURL: "https://google.com",
					},
					{
						ShortURL:    "http://localhost:8080/abcd0002",
						OriginalURL: "https://github.com",
					},
				},
				total: 2,
				error: nil,
			},
		},
		{
			name: "No URLs found",
			ctx:  context.WithValue(context.Background(), dto.CurrentUser, UserUUID),
			before: func(ctx context.Context) {
				repo.EXPECT().GetURLsByUserID(ctx, UserUUID, limit, offset).Return(nil, 0, nil)
			},
			expected: result{
				urls:  []dto.GetUserURLsResponse{},
				total: 0,
				error: nil,
			},
		},
		{
			name: "Error loading user URLs",
			ctx:  context.WithValue(context.Background(), dto.CurrentUser, UserUUID),
			before: func(ctx context.Context) {
				repo.EXPECT().GetURLsByUserID(ctx, UserUUID, limit, offset).Return(nil, 0, errors.ErrFailedToLoadUserUrls)
			},
			expected: result{
				urls:  nil,
				total: 0,
				error: errors.ErrFailedToLoadUserUrls,
			},
		},
		{
			name:   "Error invalid user ID",
			ctx:    context.Background(),
			before: func(_ context.Context) {},
			expected: result{
				urls:  nil,
				total: 0,
				error: errors.ErrInvalidUserID,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before(tt.ctx)

			urls, total, err := service.GetUserURLs(tt.ctx, &paginator)

			assert.Equal(t, tt.expected.urls, urls)
			assert.Equal(t, tt.expected.total, total)
			assert.Equal(t, tt.expected.error, err)
		})
	}
}

func Test_DeleteUserURLs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Addr:      "localhost:8080",
		BaseURL:   "http://localhost:8080",
		ClientURL: "http://localhost:8080",
	}

	repo := repository.NewMockRepository(ctrl)
	rand := NewMockSecureRandomGenerator(ctrl)
	appWorker := worker.NewMockWorker(ctrl)
	service := NewURLService(cfg, repo, rand, appWorker)

	UserUUID, _ := uuid.Parse("123e4567-e89b-12d3-a456-426614174001")

	tests := []struct {
		name     string
		ctx      context.Context
		before   func()
		params   dto.BatchDeleteShortLinkRequest
		expected error
	}{
		{
			name: "Success",
			ctx:  context.WithValue(context.Background(), dto.CurrentUser, UserUUID),
			before: func() {
				appWorker.EXPECT().Add(dto.BatchDeleteParams{
					UserID:     UserUUID,
					ShortCodes: []string{"abcd0001", "abcd0002"},
				})
			},
			params:   []string{"abcd0001", "abcd0002"},
			expected: nil,
		},
		{
			name:     "Error invalid user ID",
			ctx:      context.WithValue(context.Background(), dto.CurrentUser, nil),
			before:   func() {},
			params:   []string{"abcd0001", "abcd0002"},
			expected: errors.ErrInvalidUserID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			err := service.DeleteUserURLs(tt.ctx, tt.params)
			assert.Equal(t, tt.expected, err)
		})
	}
}
