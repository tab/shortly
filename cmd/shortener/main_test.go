package main

import (
	"flag"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"shortly/internal/app/config"
)

func Test_Main(t *testing.T) {
	cfg := &config.Config{
		Addr:    "localhost:8080",
		BaseURL: "http://localhost:8080",
	}

	tests := []struct {
		name   string
		before func()
		signal os.Signal
	}{
		{
			name:   "SIGTERM",
			signal: syscall.SIGTERM,
		},
		{
			name:   "SIGINT",
			signal: syscall.SIGINT,
		},
		{
			name:   "SIGQUIT",
			signal: syscall.SIGQUIT,
		},
	}

	for _, tt := range tests {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = []string{oldArgs[0], cfg.Addr}

		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		done := make(chan struct{})

		go func() {
			main()
			close(done)
		}()

		require.Eventually(t, func() bool {
			resp, err := http.Get(cfg.BaseURL + "/live")
			if err != nil {
				return false
			}
			defer resp.Body.Close()
			return resp.StatusCode == http.StatusOK
		}, 1*time.Second, 50*time.Millisecond, "timeout: server did not start")

		p, err := os.FindProcess(os.Getpid())
		require.NoError(t, err)
		require.NotNil(t, p)

		require.NoError(t, p.Signal(tt.signal))

		select {
		case <-done:
			// main() exited successfully
		case <-time.After(1 * time.Second):
			t.Fatal("timeout: main() did not exit after signal")
		}
	}
}
