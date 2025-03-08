package trusted

import (
	"net"
	"net/http"
)

// Middleware restricts access to endpoints based on client IP
func Middleware(trustedSubnet string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var ipNet *net.IPNet
			if trustedSubnet != "" {
				_, network, _ := net.ParseCIDR(trustedSubnet)
				ipNet = network
			}

			if ipNet == nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			clientIP := net.ParseIP(r.Header.Get("X-Real-IP"))
			if clientIP == nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			if !ipNet.Contains(clientIP) {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
