package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Run(t *testing.T) {
	errCh := make(chan error, 1)

	go func() {
		err := run()

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
}
