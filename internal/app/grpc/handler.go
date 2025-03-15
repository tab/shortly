package grpc

import (
	"context"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"shortly/internal/app/api/pagination"
	"shortly/internal/app/config"
	"shortly/internal/app/dto"
	"shortly/internal/app/errors"
	"shortly/internal/app/grpc/proto"
	"shortly/internal/app/service"
)

// Shortener is a handler for URL operations
type Shortener struct {
	cfg     *config.Config
	service service.Shortener
	proto.UnimplementedURLShortenerServer
}

// NewShortener creates a new Shortener instance
func NewShortener(cfg *config.Config, service service.Shortener) *Shortener {
	return &Shortener{cfg: cfg, service: service}
}

// CreateShortLink handles short link creation
func (s *Shortener) CreateShortLink(ctx context.Context, req *proto.CreateShortLinkV1Request) (*proto.CreateShortLinkV1Response, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, errors.ErrInvalidURL.Error())
	}

	shortURL, err := s.service.CreateShortLink(ctx, req.Url)
	if err != nil {
		if errors.Is(err, errors.ErrURLAlreadyExists) {
			return &proto.CreateShortLinkV1Response{
				ShortUrl: shortURL,
				Status:   codes.AlreadyExists.String(),
				Code:     int32(codes.AlreadyExists),
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

	return &proto.CreateShortLinkV1Response{
		ShortUrl: shortURL,
		Status:   codes.OK.String(),
		Code:     int32(codes.OK),
	}, nil
}

// CreateShortLinks handles batch short link creation
func (s *Shortener) CreateShortLinks(ctx context.Context, req *proto.BatchCreateShortLinksV1Request) (*proto.BatchCreateShortLinksV1Response, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, errors.ErrInvalidBatchParams.Error())
	}

	params := make([]dto.BatchCreateShortLinkParams, 0, len(req.Items))
	for _, item := range req.Items {
		params = append(params, dto.BatchCreateShortLinkParams{
			CorrelationID: item.CorrelationId,
			OriginalURL:   item.OriginalUrl,
		})
	}

	results, err := s.service.CreateShortLinks(ctx, params)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	items := make([]*proto.BatchCreateResult, len(results))
	for i, res := range results {
		items[i] = &proto.BatchCreateResult{
			CorrelationId: res.CorrelationID,
			ShortUrl:      res.ShortURL,
		}
	}

	return &proto.BatchCreateShortLinksV1Response{
		Items:  items,
		Status: codes.OK.String(),
		Code:   int32(codes.OK),
	}, nil
}

// GetShortLink handles short link retrieval
func (s *Shortener) GetShortLink(ctx context.Context, req *proto.GetShortLinkV1Request) (*proto.GetShortLinkV1Response, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, errors.ErrInvalidShortCode.Error())
	}

	url, ok := s.service.GetShortLink(ctx, req.ShortCode)
	if !ok {
		return nil, status.Error(codes.NotFound, errors.ErrShortLinkNotFound.Error())
	}

	if !url.DeletedAt.IsZero() {
		return nil, status.Error(codes.NotFound, errors.ErrShortLinkDeleted.Error())
	}

	return &proto.GetShortLinkV1Response{
		RedirectUrl: url.LongURL,
		Status:      codes.OK.String(),
		Code:        int32(codes.OK),
	}, nil
}

// GetUserURLs handles user URLs retrieval
func (s *Shortener) GetUserURLs(ctx context.Context, req *proto.GetUserURLsV1Request) (*proto.GetUserURLsV1Response, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	paginator := &pagination.Pagination{
		Page: int64(req.Page), //nolint:gosec
		Per:  int64(req.Per),  //nolint:gosec
	}

	results, total, err := s.service.GetUserURLs(ctx, paginator)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	items := make([]*proto.UserURL, len(results))
	for i, res := range results {
		items[i] = &proto.UserURL{
			ShortUrl:    res.ShortURL,
			OriginalUrl: res.OriginalURL,
		}
	}

	return &proto.GetUserURLsV1Response{
		Items: items,
		Total: uint64(total), //nolint:gosec
	}, nil
}

// DeleteUserURLs handles short link deletion
func (s *Shortener) DeleteUserURLs(ctx context.Context, req *proto.DeleteUserURLsV1Request) (*proto.DeleteUserURLsV1Response, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, errors.ErrShortCodeEmpty.Error())
	}

	if err := s.service.DeleteUserURLs(ctx, req.ShortCodes); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.DeleteUserURLsV1Response{
		Status: codes.OK.String(),
		Code:   int32(codes.OK),
	}, nil
}
