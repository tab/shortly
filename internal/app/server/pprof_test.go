package server

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"shortly/internal/app/config"
	"shortly/internal/spec"
)

func Test_NewPprofServer(t *testing.T) {
	cfg := &config.Config{
		ProfilerAddr: "localhost:2080",
	}

	srv := NewPprofServer(cfg)
	assert.NotNil(t, srv)

	s, ok := srv.(*pprofServer)
	assert.True(t, ok)

	assert.Equal(t, cfg.ProfilerAddr, s.httpServer.Addr)
	assert.Equal(t, 60*time.Second, s.httpServer.ReadTimeout)
}

func Test_PprofServer_RunAndShutdown(t *testing.T) {
	cfg := &config.Config{
		ProfilerAddr: "localhost:6000",
	}
	srv := NewPprofServer(cfg)

	runErrCh := make(chan error, 1)
	go func() {
		err := srv.Run()
		if err != nil && err != http.ErrServerClosed {
			runErrCh <- err
		}
		close(runErrCh)
	}()

	spec.WaitForServerStart(t, fmt.Sprintf("http://%s/debug/pprof/", cfg.ProfilerAddr))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	assert.NoError(t, err)

	err = <-runErrCh
	assert.NoError(t, err)
}
