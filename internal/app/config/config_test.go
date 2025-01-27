package config

import (
	"flag"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"shortly/internal/spec"
)

func TestMain(m *testing.M) {
	if err := spec.LoadEnv(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	if os.Getenv("GO_ENV") == "ci" {
		os.Exit(0)
	}

	code := m.Run()
	os.Exit(code)
}

func Test_LoadConfig(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		env      map[string]string
		expected *Config
	}{
		{
			name: "Use default values",
			args: []string{},
			env:  map[string]string{},
			expected: &Config{
				AppEnv:          "test",
				Addr:            "localhost:8080",
				BaseURL:         "http://localhost:8080",
				ClientURL:       "http://localhost:3000",
				ProfilerAddr:    "localhost:2080",
				FileStoragePath: "store-test.json",
				DatabaseDSN:     "postgres://postgres:postgres@localhost:5432/shortly-test?sslmode=disable",
				SecretKey:       "jwt-secret-key",
			},
		},
		{
			name: "Use env vars",
			args: []string{
				"-a", "localhost:5000",
				"-b", "http://localhost:5000",
				"-c", "http://localhost:6000",
				"-p", "localhost:2080",
				"-f", "store-test.json",
				"-d", "postgres://postgres:postgres@localhost:5432/shortly-test?sslmode=disable",
				"-s", "jwt-secret-key",
			},
			env: map[string]string{
				"SERVER_ADDRESS":    "localhost:3000",
				"BASE_URL":          "http://localhost:3000",
				"CLIENT_URL":        "http://localhost:6000",
				"PROFILER_ADDRESS":  "localhost:2080",
				"FILE_STORAGE_PATH": "store-test.json",
				"DATABASE_DSN":      "postgres://postgres:postgres@localhost:5432/shortly-test?sslmode=disable",
				"SECRET_KEY":        "jwt-secret-key",
			},
			expected: &Config{
				AppEnv:          "test",
				Addr:            "localhost:3000",
				BaseURL:         "http://localhost:3000",
				ClientURL:       "http://localhost:6000",
				ProfilerAddr:    "localhost:2080",
				FileStoragePath: "store-test.json",
				DatabaseDSN:     "postgres://postgres:postgres@localhost:5432/shortly-test?sslmode=disable",
				SecretKey:       "jwt-secret-key",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.env {
				os.Setenv(key, value)
			}

			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ContinueOnError)
			result := LoadConfig()

			assert.Equal(t, tt.expected.AppEnv, result.AppEnv)
			assert.Equal(t, tt.expected.Addr, result.Addr)
			assert.Equal(t, tt.expected.BaseURL, result.BaseURL)
			assert.Equal(t, tt.expected.ClientURL, result.ClientURL)
			assert.Equal(t, tt.expected.ProfilerAddr, result.ProfilerAddr)
			assert.Equal(t, tt.expected.FileStoragePath, result.FileStoragePath)
			assert.Equal(t, tt.expected.DatabaseDSN, result.DatabaseDSN)
			assert.Equal(t, tt.expected.SecretKey, result.SecretKey)

			t.Cleanup(func() {
				for key := range tt.env {
					os.Unsetenv(key)
				}
			})
		})
	}
}
