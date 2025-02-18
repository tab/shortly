package main

import (
	"flag"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Main(t *testing.T) {
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
		os.Args = []string{oldArgs[0], "-a=localhost:0"}

		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		done := make(chan struct{})

		go func() {
			main()
			close(done)
		}()

		time.Sleep(100 * time.Millisecond)

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
