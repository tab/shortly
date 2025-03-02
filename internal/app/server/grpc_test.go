package server

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"shortly/internal/app/config"
)

const startupTimeout = 100 * time.Millisecond

func Test_NewGRPCServer(t *testing.T) {
	cfg := &config.Config{
		GRPCServerAddr: "localhost:9090",
		GRPCSecretKey:  "grpc-secret-key",
	}

	srv := NewGRPCServer(cfg)
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
	srv := NewGRPCServer(cfg)

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
	srv := NewGRPCServer(cfg)

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
