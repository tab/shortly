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
		expected Options
	}{
		{
			name: "Use default values",
			args: []string{},
			env:  map[string]string{},
			expected: Options{
				Addr:      "localhost:8080",
				BaseURL:   "http://localhost:8080",
				ClientURL: "http://localhost:3000",
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
			expected: Options{
				Addr:      "localhost:3000",
				BaseURL:   "http://localhost:3000",
				ClientURL: "http://localhost:6000",
			},
		},
		//{
		//	name: "Use flags",
		//	args: []string{"-a", "localhost:4000", "-b", "http://localhost:4000", "-c", "http://localhost:5000"},
		//	env:  map[string]string{},
		//	expected: Options{
		//		Addr:    "localhost:4000",
		//		BaseURL: "http://localhost:4000",
		//    ClientURL: "http://localhost:5000",
		//	},
		//},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for key, value := range test.env {
				os.Setenv(key, value)
			}

			flag.CommandLine = flag.NewFlagSet(test.name, flag.ContinueOnError)
			flag.CommandLine.Parse(test.args)

			result := Init()

			assert.Equal(t, test.expected.Addr, result.Addr)
			assert.Equal(t, test.expected.BaseURL, result.BaseURL)

			for key := range test.env {
				os.Unsetenv(key)
			}
		})
	}
}