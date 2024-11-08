package spec

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoadEnv(t *testing.T) {
	type env struct {
		BaseURL         string
		ClientURL       string
		ServerAddress   string
		FileStoragePath string
		DatabaseDSN     string
	}

	tests := []struct {
		name     string
		expected env
	}{
		{
			name: "Success",
			expected: env{
				BaseURL:         "http://localhost:8080",
				ClientURL:       "http://localhost:3000",
				ServerAddress:   "localhost:8080",
				FileStoragePath: "store-test.json",
				DatabaseDSN:     "postgres://postgres:postgres@localhost:5432/shortly-test?sslmode=disable",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := LoadEnv()
			assert.NoError(t, err)

			hash := []struct{ key, value string }{
				{"BASE_URL", tt.expected.BaseURL},
				{"CLIENT_URL", tt.expected.ClientURL},
				{"SERVER_ADDRESS", tt.expected.ServerAddress},
				{"FILE_STORAGE_PATH", tt.expected.FileStoragePath},
				{"DATABASE_DSN", tt.expected.DatabaseDSN},
			}

			for _, h := range hash {
				envValue := os.Getenv(h.key)
				assert.Equal(t, h.value, envValue)
			}
		})
	}
}
