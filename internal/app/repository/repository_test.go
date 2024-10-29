package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewRepository(t *testing.T) {
	tests := []struct {
		name     string
		expected Repository
	}{
		{
			name:     "Success",
			expected: NewInMemoryRepository(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewRepository()
			assert.NotNil(t, repo)
			assert.Equal(t, tt.expected, repo)
		})
	}
}
