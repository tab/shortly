package trusted

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Trusted_Middleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"urls": 10, "users": 2}`))
		assert.NoError(t, err)
	})

	tests := []struct {
		name          string
		trustedSubnet string
		clientIP      string
		code          int
	}{
		{
			name:          "IP in trusted subnet",
			trustedSubnet: "192.168.1.0/24",
			clientIP:      "192.168.1.1",
			code:          http.StatusOK,
		},
		{
			name:          "IP not in trusted subnet",
			trustedSubnet: "10.0.0.0/24",
			clientIP:      "192.168.1.1",
			code:          http.StatusForbidden,
		},
		{
			name:          "No trusted subnet defined",
			trustedSubnet: "",
			clientIP:      "192.168.1.1",
			code:          http.StatusForbidden,
		},
		{
			name:          "Missing X-Real-IP header",
			trustedSubnet: "192.168.1.0/24",
			clientIP:      "",
			code:          http.StatusForbidden,
		},
		{
			name:          "Invalid CIDR format",
			trustedSubnet: "invalid-format",
			clientIP:      "192.168.1.1",
			code:          http.StatusForbidden,
		},
		{
			name:          "Valid IPv6 subnet",
			trustedSubnet: "2001:db8::/32",
			clientIP:      "2001:db8::1",
			code:          http.StatusOK,
		},
		{
			name:          "IPv6 not in trusted subnet",
			trustedSubnet: "2001:db8::/32",
			clientIP:      "2001:db9::1",
			code:          http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := Middleware(tt.trustedSubnet)
			wrappedHandler := middleware(handler)

			req := httptest.NewRequest("GET", "/api/internal/stats", nil)
			if tt.clientIP != "" {
				req.Header.Set("X-Real-IP", tt.clientIP)
			}

			rr := httptest.NewRecorder()
			wrappedHandler.ServeHTTP(rr, req)

			assert.Equal(t, tt.code, rr.Code)
		})
	}
}
