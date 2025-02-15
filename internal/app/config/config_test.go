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
				EnableHTTPS:     false,
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
				"-k", "jwt-secret-key",
				"-s", "true",
			},
			env: map[string]string{
				"SERVER_ADDRESS":    "localhost:3000",
				"BASE_URL":          "http://localhost:3000",
				"CLIENT_URL":        "http://localhost:6000",
				"PROFILER_ADDRESS":  "localhost:2080",
				"FILE_STORAGE_PATH": "store-test.json",
				"DATABASE_DSN":      "postgres://postgres:postgres@localhost:5432/shortly-test?sslmode=disable",
				"SECRET_KEY":        "jwt-secret-key",
				"ENABLE_HTTPS":      "true",
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
				EnableHTTPS:     true,
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
			assert.Equal(t, tt.expected.EnableHTTPS, result.EnableHTTPS)

			t.Cleanup(func() {
				for key := range tt.env {
					os.Unsetenv(key)
				}
			})
		})
	}
}

func Test_getEnvOrFlag(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		envKey   string
		flagVal  string
		expected string
	}{
		{
			name: "Flag value is empty",
			env: map[string]string{
				"TEST_ENV": "some-env-value",
			},
			envKey:   "TEST_ENV",
			flagVal:  "",
			expected: "some-env-value",
		},
		{
			name: "Flag value is present",
			env: map[string]string{
				"TEST_ENV": "some-env-value",
			},
			envKey:   "TEST_ENV",
			flagVal:  "some-flag-value",
			expected: "some-env-value",
		},
		{
			name:     "Env value is empty",
			env:      map[string]string{},
			envKey:   "TEST_ENV",
			flagVal:  "some-flag-value",
			expected: "some-flag-value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.env {
				os.Setenv(key, value)
			}

			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ContinueOnError)

			result := getEnvOrFlag(tt.envKey, tt.flagVal)
			assert.Equal(t, tt.expected, result)

			t.Cleanup(func() {
				for key := range tt.env {
					os.Unsetenv(key)
				}
			})
		})
	}
}

func Test_getEnvOrBoolFlag(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		envKey   string
		flagVal  bool
		expected bool
	}{
		{
			name: "Flag value is false",
			env: map[string]string{
				"TEST_ENV": "true",
			},
			envKey:   "TEST_ENV",
			flagVal:  false,
			expected: true,
		},
		{
			name: "Flag value is true",
			env: map[string]string{
				"TEST_ENV": "false",
			},
			envKey:   "TEST_ENV",
			flagVal:  true,
			expected: false,
		},
		{
			name:     "Env value is false",
			env:      map[string]string{},
			envKey:   "TEST_ENV",
			flagVal:  true,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.env {
				os.Setenv(key, value)
			}

			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ContinueOnError)

			result := getEnvOrBoolFlag(tt.envKey, tt.flagVal)
			assert.Equal(t, tt.expected, result)

			t.Cleanup(func() {
				for key := range tt.env {
					os.Unsetenv(key)
				}
			})
		})
	}
}
