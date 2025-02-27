package main

import (
	"flag"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"shortly/internal/app/config"
	"shortly/internal/spec"
)

const shutdownTimeout = 1 * time.Second

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

		spec.WaitForServerStart(t, cfg.BaseURL+"/live")

		p, err := os.FindProcess(os.Getpid())
		require.NoError(t, err)
		require.NotNil(t, p)

		require.NoError(t, p.Signal(tt.signal))

		select {
		case <-done:
			// main() exited successfully
		case <-time.After(shutdownTimeout):
			t.Fatal("timeout: main() did not exit after signal")
		}
	}
}
