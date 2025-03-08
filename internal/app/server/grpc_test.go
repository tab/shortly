package server

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"shortly/internal/app/config"
	"shortly/internal/app/service"
)

const startupTimeout = 100 * time.Millisecond

func Test_NewGRPCServer(t *testing.T) {
	cfg := &config.Config{
		GRPCServerAddr: "localhost:9090",
		GRPCSecretKey:  "grpc-secret-key",
	}
	shortener := &service.URLService{}
	srv := NewGRPCServer(cfg, shortener)
	assert.NotNil(t, srv)

	s, ok := srv.(*grpcServer)
	assert.True(t, ok)
	assert.Equal(t, cfg, s.cfg)
	assert.NotNil(t, s.grpcServer)
}

func Test_GRPCServer_RunAndShutdown(t *testing.T) {
	cfg := &config.Config{
		GRPCServerAddr: "localhost:9090",
		GRPCSecretKey:  "grpc-secret-key",
	}
	shortener := &service.URLService{}
	srv := NewGRPCServer(cfg, shortener)

	runErrCh := make(chan error, 1)

	go func() {
		err := srv.Run()
		if err != nil {
			runErrCh <- err
		}
		close(runErrCh)
	}()

	time.Sleep(startupTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	assert.NoError(t, err)

	err = <-runErrCh
	assert.NoError(t, err)
}

func Test_GRPCServer_RunWithTLS(t *testing.T) {
	certPath, keyPath := generateTestCertificates(t)

	cfg := &config.Config{
		GRPCServerAddr: "localhost:9090",
		GRPCSecretKey:  "grpc-secret-key",
		EnableHTTPS:    true,
		Certificate:    certPath,
		PrivateKey:     keyPath,
	}
	shortener := &service.URLService{}
	srv := NewGRPCServer(cfg, shortener)

	runErrCh := make(chan error, 1)

	go func() {
		err := srv.Run()
		if err != nil {
			runErrCh <- err
		}
		close(runErrCh)
	}()

	time.Sleep(startupTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := srv.Shutdown(ctx)
	assert.NoError(t, err)

	err = <-runErrCh
	assert.NoError(t, err)
}

func Test_GRPCServer_AuthInterceptor(t *testing.T) {
	ctx := context.Background()

	const secretKey = "test-secret-key"
	interceptor := authInterceptor(secretKey)

	tests := []struct {
		name     string
		method   string
		metadata metadata.MD
		expect   codes.Code
		err      bool
	}{
		{
			name:     "Success health check",
			method:   "/grpc.health.v1.Health/Check",
			metadata: metadata.MD{},
			expect:   codes.OK,
			err:      false,
		},
		{
			name:     "Success: valid authorization",
			method:   "/api.URLShortener/CreateShortLink",
			metadata: metadata.MD{"authorization": []string{secretKey}},
			expect:   codes.OK,
			err:      false,
		},
		{
			name:     "Missing metadata",
			method:   "/api.URLShortener/CreateShortLink",
			metadata: nil,
			expect:   codes.Unauthenticated,
			err:      true,
		},
		{
			name:     "Missing authorization header",
			method:   "/api.URLShortener/CreateShortLink",
			metadata: metadata.MD{},
			expect:   codes.Unauthenticated,
			err:      true,
		},
		{
			name:     "Invalid authorization secret",
			method:   "/api.URLShortener/CreateShortLink",
			metadata: metadata.MD{"authorization": []string{"wrong-secret"}},
			expect:   codes.Unauthenticated,
			err:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx = metadata.NewIncomingContext(ctx, tt.metadata)

			info := &grpc.UnaryServerInfo{
				FullMethod: tt.method,
			}

			handler := func(_ context.Context, _ interface{}) (interface{}, error) {
				return "test-response", nil
			}
			resp, err := interceptor(ctx, "test-request", info, handler)

			if tt.err {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.expect, st.Code())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "test-response", resp)
			}
		})
	}
}
