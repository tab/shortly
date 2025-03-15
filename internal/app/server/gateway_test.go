package server

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"shortly/internal/app/config"
)

func Test_NewGRPCGateway(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		GRPCGatewayAddr: "localhost:8181",
	}
	srv := NewGRPCGateway(ctx, cfg)
	assert.NotNil(t, srv)

	s, ok := srv.(*grpcGateway)
	assert.True(t, ok)
	assert.Equal(t, cfg, s.cfg)
}

func Test_GRPCGateway_RunAndShutdown(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		GRPCGatewayAddr: "localhost:8181",
	}
	srv := NewGRPCGateway(ctx, cfg)

	startupTimeout := 100 * time.Millisecond
	runErrCh := make(chan error, 1)

	go func() {
		err := srv.Run()
		if err != nil && err != http.ErrServerClosed {
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
