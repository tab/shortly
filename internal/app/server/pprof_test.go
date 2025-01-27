package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"shortly/internal/app/config"
)

func TestNewPprofServer(t *testing.T) {
	tests := []struct {
		name string
		cfg  *config.Config
	}{
		{
			name: "default",
			cfg:  &config.Config{ProfilerAddr: "localhost:2080"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewPprofServer(tt.cfg)
			assert.Equal(t, tt.cfg.ProfilerAddr, srv.Addr)
			assert.NotNil(t, srv.Handler)

			req, err := http.NewRequest(http.MethodGet, "/debug/pprof/", nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			srv.Handler.ServeHTTP(rr, req)
			assert.Equal(t, http.StatusOK, rr.Code)
		})
	}
}
