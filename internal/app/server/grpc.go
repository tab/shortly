package server

import (
	"context"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"shortly/internal/app/config"
)

// GRPCServer is an interface for gRPC server
type GRPCServer interface {
	Run() error
	Shutdown(ctx context.Context) error
}

type grpcServer struct {
	cfg        *config.Config
	grpcServer *grpc.Server
}

// NewGRPCServer creates a new gRPC server instance
func NewGRPCServer(cfg *config.Config) GRPCServer {
	var opts []grpc.ServerOption

	params := keepalive.ServerParameters{
		MaxConnectionIdle:     15 * time.Second,
		MaxConnectionAge:      30 * time.Second,
		MaxConnectionAgeGrace: 5 * time.Second,
		Time:                  5 * time.Second,
		Timeout:               1 * time.Second,
	}
	opts = append(opts, grpc.KeepaliveParams(params))

	if config.IsTLSEnabled(cfg) {
		tls, err := credentials.NewServerTLSFromFile(cfg.Certificate, cfg.PrivateKey)
		if err == nil {
			opts = append(opts, grpc.Creds(tls))
		}
	} else {
		opts = append(opts, grpc.Creds(insecure.NewCredentials()))
	}

	if cfg.GRPCSecretKey != "" {
		opts = append(opts, grpc.UnaryInterceptor(authInterceptor(cfg.GRPCSecretKey)))
	}

	srv := grpc.NewServer(opts...)

	return &grpcServer{
		cfg:        cfg,
		grpcServer: srv,
	}
}

// Run starts the gRPC server
func (g *grpcServer) Run() error {
	listener, err := net.Listen("tcp", g.cfg.GRPCServerAddr)
	if err != nil {
		return err
	}

	return g.grpcServer.Serve(listener)
}

// Shutdown stops the gRPC server
func (g *grpcServer) Shutdown(ctx context.Context) error {
	if g.grpcServer == nil {
		return nil
	}

	done := make(chan struct{})
	go func() {
		g.grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-ctx.Done():
		g.grpcServer.Stop()
		return ctx.Err()
	case <-done:
		return nil
	}
}

// authInterceptor is a gRPC interceptor for authentication
func authInterceptor(secretKey string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if info.FullMethod == "/grpc.health.v1.Health/Check" {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
		}

		auth, ok := md["authorization"]
		if !ok || len(auth) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "missing authorization header")
		}

		if auth[0] != secretKey {
			return nil, status.Errorf(codes.Unauthenticated, "invalid authorization secret")
		}

		return handler(ctx, req)
	}
}
