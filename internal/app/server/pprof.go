package server

import (
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/go-chi/chi/v5"

	"shortly/internal/app/config"
)

func NewPprofServer(cfg *config.Config) *http.Server {
	handler := chi.NewRouter()

	handler.HandleFunc("/debug/pprof/", pprof.Index)
	handler.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	handler.HandleFunc("/debug/pprof/profile", pprof.Profile)
	handler.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	handler.HandleFunc("/debug/pprof/trace", pprof.Trace)
	handler.HandleFunc("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)

	return &http.Server{
		Addr:        cfg.ProfilerAddr,
		Handler:     handler,
		ReadTimeout: 60 * time.Second,
	}
}
