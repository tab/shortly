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
			expected: errors.ErrOriginalURLEmpty,
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
			var params CreateShortLinkRequest
			err := params.Validate(tt.body)

			assert.Equal(t, tt.expected, err)
		})
	}
}

func Test_BatchValidate(t *testing.T) {
	tests := []struct {
		name     string
		body     io.Reader
		expected error
	}{
		{
			name:     "Success",
			body:     strings.NewReader(`[{"correlation_id": "1234", "original_url": "https://www.google.com"}]`),
			expected: nil,
		},
		{
			name:     "No correlation ID",
			body:     strings.NewReader(`[{"original_url": "https://www.google.com"}]`),
			expected: errors.ErrCorrelationIDEmpty,
		},
		{
			name:     "No original URL",
			body:     strings.NewReader(`[{"correlation_id": "1234"}]`),
			expected: errors.ErrInvalidURL,
		},
		{
			name:     "URL without scheme",
			body:     strings.NewReader(`[{"correlation_id": "1234", "original_url": "www.google.com"}]`),
			expected: errors.ErrInvalidURL,
		},
		{
			name:     "Invalid URL",
			body:     strings.NewReader(`[{"correlation_id": "1234", "original_url": "not-a-url"}]`),
			expected: errors.ErrInvalidURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var params BatchCreateShortLinkRequest
			err := params.Validate(tt.body)

			assert.Equal(t, tt.expected, err)
		})
	}
}

func Test_DeprecatedValidate(t *testing.T) {
	tests := []struct {
		name     string
		body     io.Reader
		expected error
	}{
		{
			name:     "Success",
			body:     strings.NewReader("https://www.google.com"),
			expected: nil,
		},
		{
			name:     "Success (URL with extra spaces)",
			body:     strings.NewReader(" https://www.google.com "),
			expected: nil,
		},
		{
			name:     "Empty",
			body:     strings.NewReader(""),
			expected: errors.ErrOriginalURLEmpty,
		},
		{
			name:     "URL without scheme",
			body:     strings.NewReader("www.google.com"),
			expected: errors.ErrInvalidURL,
		},
		{
			name:     "Invalid URL",
			body:     strings.NewReader("not-a-url"),
			expected: errors.ErrInvalidURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var params CreateShortLinkRequest
			err := params.DeprecatedValidate(tt.body)

			assert.Equal(t, tt.expected, err)
		})
	}
}
