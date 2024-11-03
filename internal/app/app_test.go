package app

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewApplication(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		expected Application
	}{
		{
			name:     "Success",
			expected: Application{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := NewApplication(ctx)
			assert.NoError(t, err)

			assert.NotNil(t, app)
			assert.NotNil(t, app.cfg)
			assert.NotNil(t, app.logger)
			assert.NotNil(t, app.server)
		})
	}
}
