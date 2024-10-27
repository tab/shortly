package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Run(t *testing.T) {
	os.Setenv("FILE_STORAGE_PATH", "store-test.json")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)

	go func() {
		err := run(ctx)

		if err != nil && err != http.ErrServerClosed {
			errCh <- err
		}

		close(errCh)
	}()

	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:8080/health")
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	t.Cleanup(func() {
		os.Unsetenv("FILE_STORAGE_PATH")
		os.RemoveAll("store-test.json")
	})
}
