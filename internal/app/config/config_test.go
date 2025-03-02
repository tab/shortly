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
				GRPCServerAddr:  "localhost:9090",
				GRPCSecretKey:   "grpc-secret-key",
				ProfilerAddr:    "localhost:2080",
				FileStoragePath: "store-test.json",
				DatabaseDSN:     "postgres://postgres:postgres@localhost:5432/shortly-test?sslmode=disable",
				SecretKey:       "jwt-secret-key",
				TrustedSubnet:   "",
			},
		},
		{
			name: "Use env vars",
			args: []string{
				"-a", "localhost:5000",
				"-b", "http://localhost:5000",
				"-g", "localhost:9091",
				"-s", "grpc-secret-key",
				"-p", "localhost:2080",
				"-f", "store-test.json",
				"-d", "postgres://postgres:postgres@localhost:5432/shortly-test?sslmode=disable",
				"-k", "jwt-secret-key",
				"-c", "config.json",
				"-t", "10.0.0.0/24",
			},
			env: map[string]string{
				"SERVER_ADDRESS":      "localhost:3000",
				"BASE_URL":            "http://localhost:3000",
				"GRPC_SERVER_ADDRESS": "localhost:9090",
				"GRPC_SECRET_KEY":     "grpc-secret-key",
				"PROFILER_ADDRESS":    "localhost:2080",
				"FILE_STORAGE_PATH":   "store-test.json",
				"DATABASE_DSN":        "postgres://postgres:postgres@localhost:5432/shortly-test?sslmode=disable",
				"SECRET_KEY":          "jwt-secret-key",
				"CONFIG":              "config.json",
				"TRUSTED_SUBNET":      "10.0.0.0/24",
			},
			expected: &Config{
				AppEnv:          "test",
				Addr:            "localhost:3000",
				BaseURL:         "http://localhost:3000",
				GRPCServerAddr:  "localhost:9090",
				GRPCSecretKey:   "grpc-secret-key",
				ProfilerAddr:    "localhost:2080",
				FileStoragePath: "store-test.json",
				DatabaseDSN:     "postgres://postgres:postgres@localhost:5432/shortly-test?sslmode=disable",
				SecretKey:       "jwt-secret-key",
				ConfigFilePath:  "config.json",
				TrustedSubnet:   "10.0.0.0/24",
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
			assert.Equal(t, tt.expected.GRPCServerAddr, result.GRPCServerAddr)
			assert.Equal(t, tt.expected.GRPCSecretKey, result.GRPCSecretKey)
			assert.Equal(t, tt.expected.ProfilerAddr, result.ProfilerAddr)
			assert.Equal(t, tt.expected.FileStoragePath, result.FileStoragePath)
			assert.Equal(t, tt.expected.DatabaseDSN, result.DatabaseDSN)
			assert.Equal(t, tt.expected.SecretKey, result.SecretKey)
			assert.Equal(t, tt.expected.ConfigFilePath, result.ConfigFilePath)
			assert.Equal(t, tt.expected.TrustedSubnet, result.TrustedSubnet)

			t.Cleanup(func() {
				for key := range tt.env {
					os.Unsetenv(key)
				}
			})
		})
	}
}

func Test_Config_WithFile(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		filePath string
		expected *Config
	}{
		{
			name: "Use default values",
			config: Config{
				AppEnv:          "test",
				Addr:            "localhost:8080",
				BaseURL:         "http://localhost:8080",
				GRPCServerAddr:  "localhost:9090",
				GRPCSecretKey:   "grpc-secret-key",
				ProfilerAddr:    "localhost:2080",
				FileStoragePath: "store-test.json",
				DatabaseDSN:     "postgres://postgres:postgres@localhost:5432/shortly-test?sslmode=disable",
				SecretKey:       "jwt-secret-key",
				EnableHTTPS:     false,
				TrustedSubnet:   "10.0.0.0/24",
			},
			expected: &Config{
				AppEnv:          "test",
				Addr:            "localhost:8080",
				BaseURL:         "http://localhost:8080",
				GRPCServerAddr:  "localhost:9090",
				GRPCSecretKey:   "grpc-secret-key",
				ProfilerAddr:    "localhost:2080",
				FileStoragePath: "store-test.json",
				DatabaseDSN:     "postgres://postgres:postgres@localhost:5432/shortly-test?sslmode=disable",
				SecretKey:       "jwt-secret-key",
				EnableHTTPS:     false,
				TrustedSubnet:   "10.0.0.0/24",
			},
		},
		{
			name: "Use config file",
			config: Config{
				AppEnv: "test",
			},
			filePath: "testdata/config.json",
			expected: &Config{
				AppEnv:          "test",
				Addr:            "localhost:9000",
				BaseURL:         "http://localhost:9000",
				GRPCServerAddr:  "localhost:9191",
				GRPCSecretKey:   "grpc-secret-key-test",
				ProfilerAddr:    "localhost:9080",
				FileStoragePath: "store-test.json",
				DatabaseDSN:     "postgres://postgres:postgres@localhost:5432/shortly-test?sslmode=disable",
				SecretKey:       "jwt-secret-key-test",
				EnableHTTPS:     false,
				TrustedSubnet:   "10.0.0.0/24",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := &Builder{
				cfg: &Config{
					AppEnv:          tt.config.AppEnv,
					Addr:            tt.config.Addr,
					BaseURL:         tt.config.BaseURL,
					GRPCServerAddr:  tt.config.GRPCServerAddr,
					GRPCSecretKey:   tt.config.GRPCSecretKey,
					ProfilerAddr:    tt.config.ProfilerAddr,
					FileStoragePath: tt.config.FileStoragePath,
					DatabaseDSN:     tt.config.DatabaseDSN,
					SecretKey:       tt.config.SecretKey,
					EnableHTTPS:     tt.config.EnableHTTPS,
					TrustedSubnet:   tt.config.TrustedSubnet,
				},
			}

			builder.cfg.ConfigFilePath = tt.filePath

			builder.WithFile()
			cfg := builder.Build()

			assert.Equal(t, tt.expected.AppEnv, cfg.AppEnv)
			assert.Equal(t, tt.expected.Addr, cfg.Addr)
			assert.Equal(t, tt.expected.BaseURL, cfg.BaseURL)
			assert.Equal(t, tt.expected.GRPCServerAddr, cfg.GRPCServerAddr)
			assert.Equal(t, tt.expected.GRPCSecretKey, cfg.GRPCSecretKey)
			assert.Equal(t, tt.expected.ProfilerAddr, cfg.ProfilerAddr)
			assert.Equal(t, tt.expected.FileStoragePath, cfg.FileStoragePath)
			assert.Equal(t, tt.expected.DatabaseDSN, cfg.DatabaseDSN)
			assert.Equal(t, tt.expected.SecretKey, cfg.SecretKey)
			assert.Equal(t, tt.expected.EnableHTTPS, cfg.EnableHTTPS)
			assert.Equal(t, tt.expected.TrustedSubnet, cfg.TrustedSubnet)
		})
	}
}

func Test_Config_WithFlags(t *testing.T) {
	tests := []struct {
		name     string
		flags    Flags
		config   Config
		expected Config
	}{
		{
			name: "Set all fields from flags",
			flags: Flags{
				ConfigFilePath:  "config.json",
				Addr:            "localhost:4000",
				BaseURL:         "http://localhost:4000",
				GRPCServerAddr:  "localhost:9999",
				GRPCSecretKey:   "secret",
				ProfilerAddr:    "localhost:2081",
				FileStoragePath: "store.json",
				DatabaseDSN:     "postgres://user:pass@localhost:5432/db",
				SecretKey:       "secret",
				EnableHTTPS:     false,
				TrustedSubnet:   "10.0.0.0/24",
			},
			config: Config{
				AppEnv: "test",
			},
			expected: Config{
				AppEnv:          "test",
				ConfigFilePath:  "config.json",
				Addr:            "localhost:4000",
				BaseURL:         "http://localhost:4000",
				GRPCServerAddr:  "localhost:9999",
				GRPCSecretKey:   "secret",
				ProfilerAddr:    "localhost:2081",
				FileStoragePath: "store.json",
				DatabaseDSN:     "postgres://user:pass@localhost:5432/db",
				SecretKey:       "secret",
				EnableHTTPS:     false,
				TrustedSubnet:   "10.0.0.0/24",
			},
		},
		{
			name: "Partial flags update",
			flags: Flags{
				Addr:        "localhost:5000",
				EnableHTTPS: false,
			},
			config: Config{
				AppEnv:          "test",
				Addr:            "default",
				BaseURL:         "http://default",
				GRPCServerAddr:  "default",
				GRPCSecretKey:   "default",
				ProfilerAddr:    "default",
				FileStoragePath: "default",
				DatabaseDSN:     "default",
				SecretKey:       "default",
				EnableHTTPS:     false,
				TrustedSubnet:   "10.0.0.0/24",
			},
			expected: Config{
				AppEnv:          "test",
				Addr:            "localhost:5000",
				BaseURL:         "http://default",
				GRPCServerAddr:  "default",
				GRPCSecretKey:   "default",
				ProfilerAddr:    "default",
				FileStoragePath: "default",
				DatabaseDSN:     "default",
				SecretKey:       "default",
				EnableHTTPS:     false,
				TrustedSubnet:   "10.0.0.0/24",
			},
		},
		{
			name:  "Empty flags do not override existing values",
			flags: Flags{},
			config: Config{
				AppEnv:          "test",
				Addr:            "initial",
				BaseURL:         "initial",
				GRPCServerAddr:  "initial",
				GRPCSecretKey:   "initial",
				ProfilerAddr:    "initial",
				FileStoragePath: "initial",
				DatabaseDSN:     "initial",
				SecretKey:       "initial",
				EnableHTTPS:     false,
				TrustedSubnet:   "10.0.0.0/24",
			},
			expected: Config{
				AppEnv:          "test",
				Addr:            "initial",
				BaseURL:         "initial",
				GRPCServerAddr:  "initial",
				GRPCSecretKey:   "initial",
				ProfilerAddr:    "initial",
				FileStoragePath: "initial",
				DatabaseDSN:     "initial",
				SecretKey:       "initial",
				EnableHTTPS:     false,
				TrustedSubnet:   "10.0.0.0/24",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := &Builder{
				cfg: &Config{
					AppEnv:          tt.config.AppEnv,
					Addr:            tt.config.Addr,
					BaseURL:         tt.config.BaseURL,
					GRPCServerAddr:  tt.config.GRPCServerAddr,
					GRPCSecretKey:   tt.config.GRPCSecretKey,
					ProfilerAddr:    tt.config.ProfilerAddr,
					FileStoragePath: tt.config.FileStoragePath,
					DatabaseDSN:     tt.config.DatabaseDSN,
					SecretKey:       tt.config.SecretKey,
					EnableHTTPS:     tt.config.EnableHTTPS,
					TrustedSubnet:   tt.config.TrustedSubnet,
				},
			}

			builder.WithFlags(tt.flags)
			cfg := builder.Build()

			assert.Equal(t, tt.expected.AppEnv, cfg.AppEnv)
			assert.Equal(t, tt.expected.ConfigFilePath, cfg.ConfigFilePath)
			assert.Equal(t, tt.expected.Addr, cfg.Addr)
			assert.Equal(t, tt.expected.BaseURL, cfg.BaseURL)
			assert.Equal(t, tt.expected.GRPCServerAddr, cfg.GRPCServerAddr)
			assert.Equal(t, tt.expected.GRPCSecretKey, cfg.GRPCSecretKey)
			assert.Equal(t, tt.expected.ProfilerAddr, cfg.ProfilerAddr)
			assert.Equal(t, tt.expected.FileStoragePath, cfg.FileStoragePath)
			assert.Equal(t, tt.expected.DatabaseDSN, cfg.DatabaseDSN)
			assert.Equal(t, tt.expected.SecretKey, cfg.SecretKey)
			assert.Equal(t, tt.expected.EnableHTTPS, cfg.EnableHTTPS)
			assert.Equal(t, tt.expected.TrustedSubnet, cfg.TrustedSubnet)
		})
	}
}

func Test_Config_WithENV(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		expected *Config
	}{
		{
			name: "Use default values",
			env:  map[string]string{},
			expected: &Config{
				AppEnv: "test",
			},
		},
		{
			name: "Use env vars",
			env: map[string]string{
				"SERVER_ADDRESS":      "localhost:3000",
				"BASE_URL":            "http://localhost:3000",
				"GRPC_SERVER_ADDRESS": "localhost:9090",
				"GRPC_SECRET_KEY":     "grpc-secret-key",
				"PROFILER_ADDRESS":    "localhost:2080",
				"FILE_STORAGE_PATH":   "store-test.json",
				"DATABASE_DSN":        "postgres://postgres:postgres@localhost:5432/shortly-test?sslmode=disable",
				"SECRET_KEY":          "jwt-secret-key",
				"ENABLE_HTTPS":        "false",
				"TRUSTED_SUBNET":      "10.0.0.0/24",
			},
			expected: &Config{
				AppEnv:          "test",
				Addr:            "localhost:3000",
				BaseURL:         "http://localhost:3000",
				GRPCServerAddr:  "localhost:9090",
				GRPCSecretKey:   "grpc-secret-key",
				ProfilerAddr:    "localhost:2080",
				FileStoragePath: "store-test.json",
				DatabaseDSN:     "postgres://postgres:postgres@localhost:5432/shortly-test?sslmode=disable",
				SecretKey:       "jwt-secret-key",
				EnableHTTPS:     false,
				TrustedSubnet:   "10.0.0.0/24",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.env {
				os.Setenv(key, value)
			}

			fs := flag.NewFlagSet(tt.name, flag.ContinueOnError)
			flag.CommandLine = fs

			cfg := NewConfigBuilder().WithEnv().Build()

			assert.Equal(t, tt.expected.AppEnv, cfg.AppEnv)
			assert.Equal(t, tt.expected.Addr, cfg.Addr)
			assert.Equal(t, tt.expected.BaseURL, cfg.BaseURL)
			assert.Equal(t, tt.expected.GRPCServerAddr, cfg.GRPCServerAddr)
			assert.Equal(t, tt.expected.GRPCSecretKey, cfg.GRPCSecretKey)
			assert.Equal(t, tt.expected.ProfilerAddr, cfg.ProfilerAddr)
			assert.Equal(t, tt.expected.FileStoragePath, cfg.FileStoragePath)
			assert.Equal(t, tt.expected.DatabaseDSN, cfg.DatabaseDSN)
			assert.Equal(t, tt.expected.SecretKey, cfg.SecretKey)
			assert.Equal(t, tt.expected.EnableHTTPS, cfg.EnableHTTPS)
			assert.Equal(t, tt.expected.TrustedSubnet, cfg.TrustedSubnet)

			t.Cleanup(func() {
				for key := range tt.env {
					os.Unsetenv(key)
				}
			})
		})
	}
}
