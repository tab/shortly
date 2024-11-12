package spec

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TruncateTables(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		dsn      string
		expected error
	}{
		{
			name:     "Success",
			dsn:      "postgres://postgres:postgres@localhost:5432/shortly-test?sslmode=disable",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := TruncateTables(ctx, tt.dsn)
			assert.NoError(t, err)
		})
	}
}
