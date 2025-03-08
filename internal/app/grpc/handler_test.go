package grpc

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"shortly/internal/app/config"
	"shortly/internal/app/dto"
	"shortly/internal/app/errors"
	"shortly/internal/app/grpc/proto"
	"shortly/internal/app/repository"
	"shortly/internal/app/service"
	"shortly/internal/app/worker"
)

func Test_NewShortener(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
	}
	repo := repository.NewMockDatabase(ctrl)
	rand := service.NewMockSecureRandomGenerator(ctrl)
	appWorker := worker.NewMockWorker(ctrl)
	srv := service.NewURLService(cfg, repo, rand, appWorker)

	handler := NewShortener(cfg, srv)

	assert.NotNil(t, handler)
	assert.Equal(t, cfg, handler.cfg)
	assert.Equal(t, srv, handler.service)
}

func Test_shortener_CreateShortLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
	}
	repo := repository.NewMockDatabase(ctrl)
	rand := service.NewMockSecureRandomGenerator(ctrl)
	appWorker := worker.NewMockWorker(ctrl)
	srv := service.NewURLService(cfg, repo, rand, appWorker)

	handler := NewShortener(cfg, srv)

	UUID := uuid.MustParse("6455bd07-e431-4851-af3c-4f703f726639")

	type result struct {
		response *proto.CreateShortLinkResponse
		err      error
	}

	tests := []struct {
		name     string
		request  *proto.CreateShortLinkRequest
		before   func()
		expected result
	}{
		{
			name: "Success",
			request: &proto.CreateShortLinkRequest{
				Url: "https://example.com",
			},
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
				response: &proto.CreateShortLinkResponse{
					Result: "http://localhost:8080/abcd1234",
					Status: codes.OK.String(),
					Code:   int32(codes.OK),
				},
				err: nil,
			},
		},
		{
			name: "URL already exists",
			request: &proto.CreateShortLinkRequest{
				Url: "https://example.com",
			},
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
				response: &proto.CreateShortLinkResponse{
					Result: "http://localhost:8080/abab0001",
					Status: codes.AlreadyExists.String(),
					Code:   int32(codes.AlreadyExists),
				},
				err: nil,
			},
		},
		{
			name: "Empty URL",
			request: &proto.CreateShortLinkRequest{
				Url: "",
			},
			before: func() {},
			expected: result{
				response: nil,
				err:      status.Error(codes.InvalidArgument, errors.ErrInvalidURL.Error()),
			},
		},
		{
			name: "URL with whitespace only",
			request: &proto.CreateShortLinkRequest{
				Url: "   ",
			},
			before: func() {},
			expected: result{
				response: nil,
				err:      status.Error(codes.InvalidArgument, errors.ErrInvalidURL.Error()),
			},
		},
		{
			name: "Invalid URL",
			request: &proto.CreateShortLinkRequest{
				Url: "not-a-valid-url",
			},
			before: func() {},
			expected: result{
				response: nil,
				err:      status.Error(codes.InvalidArgument, errors.ErrInvalidURL.Error()),
			},
		},
		{
			name: "Error generating UUID",
			request: &proto.CreateShortLinkRequest{
				Url: "https://example.com",
			},
			before: func() {
				rand.EXPECT().UUID().Return(uuid.UUID{}, errors.ErrFailedToGenerateUUID)
			},
			expected: result{
				response: nil,
				err:      status.Error(codes.Internal, errors.ErrFailedToGenerateUUID.Error()),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			response, err := handler.CreateShortLink(ctx, tt.request)

			if tt.expected.err != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expected.err.Error(), err.Error())
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response.Result, response.Result)
				assert.Equal(t, tt.expected.response.Status, response.Status)
				assert.Equal(t, tt.expected.response.Code, response.Code)
			}
		})
	}
}

func Test_Shortener_GetShortLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
	}
	repo := repository.NewMockDatabase(ctrl)
	rand := service.NewMockSecureRandomGenerator(ctrl)
	appWorker := worker.NewMockWorker(ctrl)
	srv := service.NewURLService(cfg, repo, rand, appWorker)

	handler := NewShortener(cfg, srv)

	type result struct {
		response *proto.GetShortLinkResponse
		err      error
	}

	tests := []struct {
		name     string
		request  *proto.GetShortLinkRequest
		before   func()
		expected result
	}{
		{
			name: "Success",
			request: &proto.GetShortLinkRequest{
				ShortCode: "abcd1234",
			},
			before: func() {
				repo.EXPECT().GetURLByShortCode(ctx, "abcd1234").Return(&repository.URL{
					LongURL:   "https://example.com",
					ShortCode: "abcd1234",
				}, true)
			},
			expected: result{
				response: &proto.GetShortLinkResponse{
					Result: "https://example.com",
					Status: codes.OK.String(),
					Code:   int32(codes.OK),
				},
				err: nil,
			},
		},
		{
			name: "Empty Short Code",
			request: &proto.GetShortLinkRequest{
				ShortCode: "",
			},
			before: func() {},
			expected: result{
				response: nil,
				err:      status.Error(codes.InvalidArgument, errors.ErrInvalidShortCode.Error()),
			},
		},
		{
			name: "Not Found",
			request: &proto.GetShortLinkRequest{
				ShortCode: "notfound",
			},
			before: func() {
				repo.EXPECT().GetURLByShortCode(ctx, "notfound").Return(nil, false)
			},
			expected: result{
				response: nil,
				err:      status.Error(codes.NotFound, errors.ErrShortLinkNotFound.Error()),
			},
		},
		{
			name: "Deleted",
			request: &proto.GetShortLinkRequest{
				ShortCode: "abcd1234",
			},
			before: func() {
				repo.EXPECT().GetURLByShortCode(ctx, "abcd1234").Return(&repository.URL{
					LongURL:   "https://example.com",
					ShortCode: "abcd1234",
					DeletedAt: time.Now(),
				}, true)
			},
			expected: result{
				response: nil,
				err:      status.Error(codes.NotFound, errors.ErrShortLinkDeleted.Error()),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			response, err := handler.GetShortLink(ctx, tt.request)

			if tt.expected.err != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expected.err.Error(), err.Error())
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response.Result, response.Result)
				assert.Equal(t, tt.expected.response.Status, response.Status)
				assert.Equal(t, tt.expected.response.Code, response.Code)
			}
		})
	}
}

func Test_Shortener_GetUserURLs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
	}
	repo := repository.NewMockDatabase(ctrl)
	rand := service.NewMockSecureRandomGenerator(ctrl)
	appWorker := worker.NewMockWorker(ctrl)
	srv := service.NewURLService(cfg, repo, rand, appWorker)

	handler := NewShortener(cfg, srv)

	limit := int64(25)
	offset := int64(0)

	UUID1, _ := uuid.Parse("6455bd07-e431-4851-af3c-4f703f720001")
	UUID2, _ := uuid.Parse("6455bd07-e431-4851-af3c-4f703f720002")
	UserUUID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174001")

	ctx := context.WithValue(context.Background(), dto.CurrentUser, UserUUID)

	type result struct {
		response *proto.GetUserURLsResponse
		err      error
	}

	tests := []struct {
		name     string
		request  *proto.GetUserURLsRequest
		before   func()
		expected result
	}{
		{
			name: "Success",
			request: &proto.GetUserURLsRequest{
				Page: 1,
				Per:  25,
			},
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
				response: &proto.GetUserURLsResponse{
					Items: []*proto.UserURL{
						{
							ShortUrl:    fmt.Sprintf("%s/%s", cfg.BaseURL, "abcd0001"),
							OriginalUrl: "https://google.com",
						},
						{
							ShortUrl:    fmt.Sprintf("%s/%s", cfg.BaseURL, "abcd0002"),
							OriginalUrl: "https://github.com",
						},
					},
					Total: 2,
				},
			},
		},
		{
			name: "No URLs",
			request: &proto.GetUserURLsRequest{
				Page: 1,
				Per:  25,
			},
			before: func() {
				repo.EXPECT().GetURLsByUserID(ctx, UserUUID, limit, offset).Return([]repository.URL{}, 0, nil)
			},
			expected: result{
				response: &proto.GetUserURLsResponse{
					Items: []*proto.UserURL{},
					Total: 0,
				},
			},
		},
		{
			name: "Error",
			request: &proto.GetUserURLsRequest{
				Page: 1,
				Per:  25,
			},
			before: func() {
				repo.EXPECT().GetURLsByUserID(ctx, UserUUID, limit, offset).Return(nil, 0, errors.ErrFailedToLoadUserUrls)
			},
			expected: result{
				response: nil,
				err:      status.Error(codes.Internal, errors.ErrFailedToLoadUserUrls.Error()),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			response, err := handler.GetUserURLs(ctx, tt.request)

			if tt.expected.err != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expected.err.Error(), err.Error())
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response.Items, response.Items)
				assert.Equal(t, tt.expected.response.Total, response.Total)
			}
		})
	}
}

func Test_Shortener_DeleteUserURLs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
	}
	repo := repository.NewMockDatabase(ctrl)
	rand := service.NewMockSecureRandomGenerator(ctrl)
	appWorker := worker.NewMockWorker(ctrl)
	srv := service.NewURLService(cfg, repo, rand, appWorker)

	handler := NewShortener(cfg, srv)

	UserUUID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174001")

	type result struct {
		response *proto.DeleteUserURLsResponse
		err      error
	}

	tests := []struct {
		name     string
		request  *proto.DeleteUserURLsRequest
		ctx      context.Context
		before   func()
		expected result
	}{
		{
			name: "Success",
			request: &proto.DeleteUserURLsRequest{
				ShortCodes: []string{"abcd1234", "efgh5678"},
			},
			ctx: context.WithValue(context.Background(), dto.CurrentUser, UserUUID),
			before: func() {
				appWorker.EXPECT().Add(dto.BatchDeleteParams{
					UserID:     UserUUID,
					ShortCodes: []string{"abcd1234", "efgh5678"},
				})
			},
			expected: result{
				response: &proto.DeleteUserURLsResponse{
					Success: true,
					Status:  codes.OK.String(),
					Code:    int32(codes.OK),
				},
				err: nil,
			},
		},
		{
			name: "Empty",
			request: &proto.DeleteUserURLsRequest{
				ShortCodes: []string{},
			},
			ctx:    context.WithValue(context.Background(), dto.CurrentUser, UserUUID),
			before: func() {},
			expected: result{
				response: nil,
				err:      status.Error(codes.InvalidArgument, errors.ErrShortCodeEmpty.Error()),
			},
		},
		{
			name: "Error",
			request: &proto.DeleteUserURLsRequest{
				ShortCodes: []string{"abcd1234", "efgh5678"},
			},
			ctx:    context.WithValue(context.Background(), dto.CurrentUser, nil),
			before: func() {},
			expected: result{
				response: nil,
				err:      status.Error(codes.Internal, errors.ErrInvalidUserID.Error()),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			response, err := handler.DeleteUserURLs(tt.ctx, tt.request)

			if tt.expected.err != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expected.err.Error(), err.Error())
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.response.Success, response.Success)
				assert.Equal(t, tt.expected.response.Status, response.Status)
				assert.Equal(t, tt.expected.response.Code, response.Code)
			}
		})
	}
}
