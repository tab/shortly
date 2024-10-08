package config

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		env      map[string]string
		expected *AppConfig
	}{
		{
			name: "Use default values",
			args: []string{},
			env:  map[string]string{},
			expected: &AppConfig{
				Addr:      ServerAddress,
				BaseURL:   BaseURL,
				ClientURL: ClientURL,
			},
		},
		{
			name: "Use env vars",
			args: []string{"-a", "localhost:5000", "-b", "http://localhost:5000", "-c", "http://localhost:6000"},
			env: map[string]string{
				"SERVER_ADDRESS": "localhost:3000",
				"BASE_URL":       "http://localhost:3000",
				"CLIENT_URL":     "http://localhost:6000",
			},
			expected: &AppConfig{
				Addr:      "localhost:3000",
				BaseURL:   "http://localhost:3000",
				ClientURL: "http://localhost:6000",
			},
		},
		//{
		//	name: "Use flags",
		//	args: []string{"-a", "localhost:4000", "-b", "http://localhost:4000", "-c", "http://localhost:5000"},
		//	env:  map[string]string{},
		//	expected: &AppConfig{
		//		Addr:      "localhost:4000",
		//		BaseURL:   "http://localhost:4000",
		//		ClientURL: "http://localhost:5000",
		//	},
		//},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for key, value := range test.env {
				os.Setenv(key, value)
			}

			flag.CommandLine = flag.NewFlagSet(test.name, flag.ContinueOnError)
			//err := flag.CommandLine.Parse(test.args)
			//if err != nil {
			//	t.Fatalf("Failed to parse flags: %v", err)
			//}

			result := Init()

			assert.Equal(t, test.expected.Addr, result.Addr)
			assert.Equal(t, test.expected.BaseURL, result.BaseURL)
			assert.Equal(t, test.expected.ClientURL, result.ClientURL)

			for key := range test.env {
				os.Unsetenv(key)
			}
		})
	}
}
