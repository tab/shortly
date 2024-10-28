package dto

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"shortly/internal/app/errors"
)

func Test_Validate(t *testing.T) {
	tests := []struct {
		name     string
		body     io.Reader
		expected error
	}{
		{
			name:     "Success",
			body:     strings.NewReader(`{"url": "https://www.google.com"}`),
			expected: nil,
		},
		{
			name:     "Success (URL with extra spaces)",
			body:     strings.NewReader(`{"url": " https://www.google.com "}`),
			expected: nil,
		},
		{
			name:     "Empty URL",
			body:     strings.NewReader(`{"url": ""}`),
			expected: errors.ErrRequestBodyEmpty,
		},
		{
			name:     "URL without scheme",
			body:     strings.NewReader(`{"url": "www.google.com"}`),
			expected: errors.ErrInvalidURL,
		},
		{
			name:     "Invalid URL",
			body:     strings.NewReader(`{"url": "not-a-url"}`),
			expected: errors.ErrInvalidURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var params CreateShortLinkParams
			err := params.Validate(tt.body)

			assert.Equal(t, tt.expected, err)
		})
	}
}
