package config

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_IsTLSEnabled(t *testing.T) {
	certFilePath, privateKeyFilePath, err := createTempCertFiles()
	assert.NoError(t, err)

	tests := []struct {
		name     string
		cfg      *Config
		expected bool
	}{
		{
			name: "TLS enabled, files are present",
			cfg: &Config{
				EnableHTTPS: true,
				Certificate: certFilePath,
				PrivateKey:  privateKeyFilePath,
			},
			expected: true,
		},
		{
			name: "TLS enabled, certificate and key files are missing",
			cfg: &Config{
				EnableHTTPS: true,
			},
			expected: false,
		},
		{
			name: "TLS disabled",
			cfg: &Config{
				EnableHTTPS: false,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsTLSEnabled(tt.cfg)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func createTempCertFiles() (string, string, error) {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1234567890),
		Subject: pkix.Name{
			Organization: []string{"Shortly"},
			Country:      []string{"IO"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(1, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatal(err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	var certPEM bytes.Buffer
	pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	var privateKeyPEM bytes.Buffer
	pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	certFile, err := os.CreateTemp("", "certificate-*.pem")
	if err != nil {
		return "", "", err
	}
	defer certFile.Close()

	privateKeyFile, err := os.CreateTemp("", "private-key-*.pem")
	if err != nil {
		return "", "", err
	}
	defer privateKeyFile.Close()

	_, err = certFile.WriteString(certPEM.String())
	if err != nil {
		return "", "", err
	}
	_, err = privateKeyFile.WriteString(privateKeyPEM.String())
	if err != nil {
		return "", "", err
	}

	return certFile.Name(), privateKeyFile.Name(), nil
}
