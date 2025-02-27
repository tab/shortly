package spec

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_WaitForServerStart(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tests := []struct {
		name string
		path string
		err  bool
	}{
		{
			name: "Success",
			path: server.URL,
			err:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WaitForServerStart(t, tt.path)
		})
	}
}
