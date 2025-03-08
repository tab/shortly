package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"shortly/internal/app/config"
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
