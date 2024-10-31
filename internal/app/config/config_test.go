package config

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
				Addr:            ServerAddress,
				BaseURL:         BaseURL,
				ClientURL:       ClientURL,
				FileStoragePath: FileStoragePath,
			},
		},
		{
			name: "Use env vars",
			args: []string{"-a", "localhost:5000", "-b", "http://localhost:5000", "-c", "http://localhost:6000"},
			env: map[string]string{
				"SERVER_ADDRESS":    "localhost:3000",
				"BASE_URL":          "http://localhost:3000",
				"CLIENT_URL":        "http://localhost:6000",
				"FILE_STORAGE_PATH": "store-test.json",
			},
			expected: &Config{
				Addr:            "localhost:3000",
				BaseURL:         "http://localhost:3000",
				ClientURL:       "http://localhost:6000",
				FileStoragePath: "store-test.json",
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

			assert.Equal(t, tt.expected.Addr, result.Addr)
			assert.Equal(t, tt.expected.BaseURL, result.BaseURL)
			assert.Equal(t, tt.expected.ClientURL, result.ClientURL)

			t.Cleanup(func() {
				for key := range tt.env {
					os.Unsetenv(key)
				}
			})
		})
	}
}
