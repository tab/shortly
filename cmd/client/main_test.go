package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFlags(t *testing.T) {
	args := os.Args
	defer func() { os.Args = args }()

	os.Args = []string{"cmd", "-a", "http://testserver:8080"}

	cfg := parseFlags()

	assert.Equal(t, "http://testserver:8080", cfg.BaseURL)
}
