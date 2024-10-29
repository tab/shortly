package server

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"shortly/internal/app/config"
	"shortly/internal/app/repository"
	"shortly/internal/app/router"
	"shortly/internal/logger"
)

func Test_NewServer(t *testing.T) {
	cfg := &config.Config{
		ClientURL: "http://localhost:8080",
	}
	repo := repository.NewRepository()
	appLogger := logger.NewLogger()
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
