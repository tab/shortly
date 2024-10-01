package config

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected Options
	}{
		{
			name: "Default flags",
			args: []string{},
			expected: Options{
				Addr:    "localhost:8080",
				BaseURL: "http://localhost:8080",
			},
		},
		//{
		//	name: "Custom address",
		//	args: []string{"-a", "0.0.0.0:9090"},
		//	expected: Options{
		//		Addr:    "0.0.0.0:9090",
		//		BaseURL: "http://localhost:8080",
		//	},
		//},
		//{
		//	name: "Custom base URL",
		//	args: []string{"-b", "http://example.com"},
		//	expected: Options{
		//		Addr:    "localhost:8080",
		//		BaseURL: "http://example.com",
		//	},
		//},
		//{
		//	name: "Custom address and base URL",
		//	args: []string{"-a", "0.0.0.0:8081", "-b", "http://example.com"},
		//	expected: Options{
		//		Addr:    "0.0.0.0:8081",
		//		BaseURL: "http://example.com",
		//	},
		//},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(test.name, flag.ContinueOnError)
			flag.CommandLine.Parse(test.args)

			result := Init()

			assert.Equal(t, test.expected.Addr, result.Addr)
			assert.Equal(t, test.expected.BaseURL, result.BaseURL)
		})
	}
}
