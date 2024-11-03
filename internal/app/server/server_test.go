package server

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"shortly/internal/app/config"
	"shortly/internal/app/repository"
	"shortly/internal/app/router"
	"shortly/internal/logger"
)

func Test_NewServer(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		ClientURL:   config.ClientURL,
		DatabaseDSN: config.DatabaseDSN,
	}
	appLogger := logger.NewLogger()
	repo, _ := repository.NewRepository(ctx, &repository.Factory{
		DSN:    cfg.DatabaseDSN,
		Logger: appLogger,
	})
	appRouter := router.NewRouter(cfg, appLogger, repo)

	tests := []struct {
		name     string
		expected Server
	}{
		{
			name:     "Success",
			expected: Server{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewServer(cfg, appRouter)

			assert.NotNil(t, srv)
			assert.Equal(t, cfg.Addr, srv.httpServer.Addr)
			assert.Equal(t, appRouter, srv.httpServer.Handler)
			assert.Equal(t, 5*time.Second, srv.httpServer.ReadTimeout)
			assert.Equal(t, 10*time.Second, srv.httpServer.WriteTimeout)
			assert.Equal(t, 120*time.Second, srv.httpServer.IdleTimeout)
		})
	}
}

func Test_Server_RunAndShutdown(t *testing.T) {
	cfg := &config.Config{
		Addr: config.ServerAddress,
	}
	handler := http.NewServeMux()
	srv := NewServer(cfg, handler)

	runErrCh := make(chan error, 1)
	go func() {
		err := srv.Run()
		if err != nil && err != http.ErrServerClosed {
			runErrCh <- err
		}
		close(runErrCh)
	}()

	time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	assert.NoError(t, err)

	err = <-runErrCh
	assert.NoError(t, err)
}
