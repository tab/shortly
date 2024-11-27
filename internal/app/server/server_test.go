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
	"shortly/internal/app/worker"
	"shortly/internal/logger"
)

func Test_NewServer(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		Addr: "localhost:8080",
	}
	appLogger := logger.NewLogger()
	repo, _ := repository.NewRepository(ctx, &repository.Factory{
		DSN:    cfg.DatabaseDSN,
		Logger: appLogger,
	})
	appWorker := worker.NewDeleteWorker(ctx, cfg, repo, appLogger)
	appRouter := router.NewRouter(cfg, repo, appWorker, appLogger)

	srv := NewServer(cfg, appRouter)
	assert.NotNil(t, srv)

	s, ok := srv.(*server)
	assert.True(t, ok)

	assert.Equal(t, cfg.Addr, s.httpServer.Addr)
	assert.Equal(t, appRouter, s.httpServer.Handler)
	assert.Equal(t, 5*time.Second, s.httpServer.ReadTimeout)
	assert.Equal(t, 10*time.Second, s.httpServer.WriteTimeout)
	assert.Equal(t, 120*time.Second, s.httpServer.IdleTimeout)
}

func Test_Server_RunAndShutdown(t *testing.T) {
	cfg := &config.Config{
		Addr: "localhost:5000",
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
