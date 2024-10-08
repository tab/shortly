package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name  string
		url   string
		valid bool
	}{
		{
			name:  "Valid URL",
			url:   "https://www.google.com",
			valid: true,
		},
		{
			name:  "Invalid URL",
			url:   "not-a-url",
			valid: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := Validate(test.url)

			if test.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
			assert.Equal(t, test.valid, result)
		})
	}
}
