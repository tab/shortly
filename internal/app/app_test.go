package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewApplicaton(t *testing.T) {
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
			app, err := NewApplication()
			assert.NoError(t, err)

			assert.NotNil(t, app)
			assert.NotNil(t, app.cfg)
			assert.NotNil(t, app.logger)
			assert.NotNil(t, app.server)
		})
	}
}
