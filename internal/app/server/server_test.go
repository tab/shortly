package server

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"shortly/internal/app/config"
	"shortly/internal/app/repository"
	"shortly/internal/app/router"
	"shortly/internal/app/worker"
	"shortly/internal/logger"
	"shortly/internal/spec"
)

func Test_NewServer(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		Addr: "localhost:8080",
	}
	appLogger := logger.NewLogger()
	repo, _ := repository.NewRepository(ctx, &repository.Factory{
		DSN:    cfg.DatabaseDSN,
		Logger: appLogger,
	})
	appWorker := worker.NewDeleteWorker(ctx, cfg, repo, appLogger)
	appRouter := router.NewRouter(cfg, repo, appWorker, appLogger)

	srv := NewServer(cfg, appRouter)
	assert.NotNil(t, srv)

	s, ok := srv.(*server)
	assert.True(t, ok)

	assert.Equal(t, cfg.Addr, s.httpServer.Addr)
	assert.Equal(t, appRouter, s.httpServer.Handler)
	assert.Equal(t, 5*time.Second, s.httpServer.ReadTimeout)
	assert.Equal(t, 10*time.Second, s.httpServer.WriteTimeout)
	assert.Equal(t, 120*time.Second, s.httpServer.IdleTimeout)
}

func Test_Server_RunAndShutdown(t *testing.T) {
	cfg := &config.Config{
		Addr:    "localhost:8181",
		BaseURL: "http://localhost:8181",
	}
	handler := http.NewServeMux()
	handler.HandleFunc("/live", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	srv := NewServer(cfg, handler)

	runErrCh := make(chan error, 1)
	go func() {
		err := srv.Run()
		if err != nil && err != http.ErrServerClosed {
			runErrCh <- err
		}
		close(runErrCh)
	}()

	spec.WaitForServerStart(t, cfg.BaseURL+"/live")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	assert.NoError(t, err)

	err = <-runErrCh
	assert.NoError(t, err)
}

func Test_Server_RunAndShutdownTLS(t *testing.T) {
	certPath, keyPath := generateTestCertificates(t)

	cfg := &config.Config{
		Addr:        "localhost:10443",
		BaseURL:     "https://localhost:10443",
		EnableHTTPS: true,
		Certificate: certPath,
		PrivateKey:  keyPath,
	}
	handler := http.NewServeMux()
	handler.HandleFunc("/live", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	srv := NewServer(cfg, handler)

	runErrCh := make(chan error, 1)
	go func() {
		err := srv.Run()
		if err != nil && err != http.ErrServerClosed {
			runErrCh <- err
		}
		close(runErrCh)
	}()

	spec.WaitForServerStart(t, cfg.BaseURL+"/live")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	assert.NoError(t, err)

	err = <-runErrCh
	assert.NoError(t, err)
}

func generateTestCertificates(t *testing.T) (certPath, keyPath string) {
	tempDir, err := os.MkdirTemp("", "tls-test-*")
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(tempDir) })

	certPath = filepath.Join(tempDir, "cert.pem")
	keyPath = filepath.Join(tempDir, "key.pem")

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Acme"},
			CommonName:   "localhost",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	require.NoError(t, err)

	certFile, err := os.Create(certPath)
	require.NoError(t, err)
	defer certFile.Close()

	err = pem.Encode(certFile, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})
	require.NoError(t, err)

	keyFile, err := os.Create(keyPath)
	require.NoError(t, err)
	defer keyFile.Close()

	err = pem.Encode(keyFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	require.NoError(t, err)

	return certPath, keyPath
}
