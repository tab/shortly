package grpc

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"shortly/internal/app/config"
	"shortly/internal/app/errors"
	"shortly/internal/app/grpc/proto"
	"shortly/internal/app/service"
	"shortly/internal/app/validator"
)

type Shortener struct {
	cfg     *config.Config
	service service.Shortener
	proto.UnimplementedURLShortenerServer
}

func NewShortener(cfg *config.Config, service service.Shortener) *Shortener {
	return &Shortener{cfg: cfg, service: service}
}

func (s *Shortener) CreateShortLink(ctx context.Context, req *proto.CreateShortLinkRequest) (*proto.CreateShortLinkResponse, error) {
	url := strings.TrimSpace(req.Url)
	if url == "" {
		return nil, status.Error(codes.InvalidArgument, errors.ErrOriginalURLEmpty.Error())
	}

	if err := validator.Validate(url); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	shortURL, err := s.service.CreateShortLink(ctx, url)
	if err != nil {
		if errors.Is(err, errors.ErrURLAlreadyExists) {
			return &proto.CreateShortLinkResponse{
				Result: shortURL,
				Status: codes.AlreadyExists.String(),
				Code:   int32(codes.AlreadyExists),
			}, nil
		}

		switch {
		case errors.Is(err, errors.ErrOriginalURLEmpty),
			errors.Is(err, errors.ErrInvalidURL):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, errors.ErrFailedToGenerateUUID),
			errors.Is(err, errors.ErrFailedToGenerateCode),
			errors.Is(err, errors.ErrFailedToSaveURL):
			return nil, status.Error(codes.Internal, err.Error())
		default:
			return nil, status.Error(codes.Unknown, err.Error())
		}
	}

	return &proto.CreateShortLinkResponse{
		Result: shortURL,
		Status: codes.OK.String(),
		Code:   int32(codes.OK),
	}, nil
}

func (s *Shortener) GetShortLink(ctx context.Context, req *proto.GetShortLinkRequest) (*proto.GetShortLinkResponse, error) {
	shortCode := strings.TrimSpace(req.ShortCode)
	if shortCode == "" {
		return nil, status.Error(codes.InvalidArgument, errors.ErrShortCodeEmpty.Error())
	}

	url, ok := s.service.GetShortLink(ctx, shortCode)
	if !ok {
		return nil, status.Error(codes.NotFound, errors.ErrShortLinkNotFound.Error())
	}

	if !url.DeletedAt.IsZero() {
		return nil, status.Error(codes.NotFound, errors.ErrShortLinkDeleted.Error())
	}

	return &proto.GetShortLinkResponse{
		Result: url.LongURL,
		Status: codes.OK.String(),
		Code:   int32(codes.OK),
	}, nil
}
