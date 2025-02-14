package server

import (
	"context"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/go-chi/chi/v5"

	"shortly/internal/app/config"
)

// PprofServer is an interface for pprof server
type PprofServer interface {
	Run() error
	Shutdown(ctx context.Context) error
}

type pprofServer struct {
	httpServer *http.Server
}

// NewPprofServer creates a new pprof server instance
func NewPprofServer(cfg *config.Config) PprofServer {
	handler := chi.NewRouter()

	handler.HandleFunc("/debug/pprof/", pprof.Index)
	handler.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	handler.HandleFunc("/debug/pprof/profile", pprof.Profile)
	handler.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	handler.HandleFunc("/debug/pprof/trace", pprof.Trace)
	handler.HandleFunc("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)

	return &pprofServer{
		httpServer: &http.Server{
			Addr:        cfg.ProfilerAddr,
			Handler:     handler,
			ReadTimeout: 60 * time.Second,
		},
	}
}

// Run starts the pprof server
func (p *pprofServer) Run() error {
	return p.httpServer.ListenAndServe()
}

// Shutdown stops the pprof server
func (p *pprofServer) Shutdown(ctx context.Context) error {
	return p.httpServer.Shutdown(ctx)
}
