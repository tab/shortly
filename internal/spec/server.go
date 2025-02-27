package spec

import (
	"crypto/tls"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	// serverStartTimeout is the maximum time to wait for server to start
	serverStartTimeout = 1 * time.Second
	// serverStartPollInterval is how frequently to check if server has started
	serverStartPollInterval = 50 * time.Millisecond
	// clientPollTimeout is the maximum time to wait for client to connect to server
	clientPollTimeout = 500 * time.Millisecond
)

// WaitForServerStart polls the specified URL until it returns a 200 OK response or times out
func WaitForServerStart(t *testing.T, url string) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				//nolint:gosec
				InsecureSkipVerify: true,
			},
		},
		Timeout: clientPollTimeout,
	}

	require.Eventually(t, func() bool {
		resp, err := client.Get(url)

		if err != nil {
			return false
		}

		defer resp.Body.Close()
		return resp.StatusCode == http.StatusOK
	}, serverStartTimeout, serverStartPollInterval, "timeout: server did not start")
}
