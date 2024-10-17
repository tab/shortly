package app

import (
	"net/http"
)

type Server interface {
	Serve(addr string, handler http.Handler) error
}

type HTTPServer struct{}

func (s *HTTPServer) Serve(addr string, handler http.Handler) error {
	return http.ListenAndServe(addr, handler)
}
