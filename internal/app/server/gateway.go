package server

import (
	"context"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"shortly/internal/app/config"
	"shortly/internal/app/grpc/proto"
)

// GRPCGateway is an interface for gRPC gateway
type GRPCGateway interface {
	Run() error
	Shutdown(ctx context.Context) error
}

// grpcGateway is a gRPC gateway implementation
type grpcGateway struct {
	cfg        *config.Config
	httpServer *http.Server
}

// NewGRPCGateway creates a new gRPC gateway instance
func NewGRPCGateway(ctx context.Context, cfg *config.Config) GRPCGateway {
	handler := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := proto.RegisterURLShortenerHandlerFromEndpoint(ctx, handler, cfg.GRPCServerAddr, opts)
	if err != nil {
		return nil
	}

	return &grpcGateway{
		cfg: cfg,
		httpServer: &http.Server{
			Addr:         cfg.GRPCGatewayAddr,
			Handler:      handler,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
	}
}

// Run starts the gRPC gateway server
func (g *grpcGateway) Run() error {
	return g.httpServer.ListenAndServe()
}

// Shutdown stops the gRPC gateway server
func (g *grpcGateway) Shutdown(ctx context.Context) error {
	if g.httpServer == nil {
		return nil
	}

	return g.httpServer.Shutdown(ctx)
}
