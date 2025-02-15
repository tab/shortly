package config

import (
	"crypto/tls"
	"os"
)

// IsTLSEnabled checks can TLS be enabled
func IsTLSEnabled(cfg *Config) bool {
	if !cfg.EnableHTTPS {
		return false
	}

	if !fileExists(cfg.Certificate) || !fileExists(cfg.PrivateKey) {
		return false
	}

	if err := validateTLSFiles(cfg.Certificate, cfg.PrivateKey); err != nil {
		return false
	}

	return true
}

// fileExists checks if a file present at the given path
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// validateTLSFiles checks if the given files are valid certificate and key pair
func validateTLSFiles(certPath, keyPath string) error {
	_, err := tls.LoadX509KeyPair(certPath, keyPath)
	return err
}
