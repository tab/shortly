package version

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Print(t *testing.T) {
	tests := []struct {
		name     string
		version  version
		expected string
	}{
		{
			name: "Success",
			version: version{
				buildVersion: "v1.0.0",
				buildDate:    "01.02.2025",
				buildCommit:  "abcd1234",
			},
			expected: "Build version: v1.0.0\nBuild date: 01.02.2025\nBuild commit: abcd1234\n",
		},
		{
			name: "Empty",
			version: version{
				buildVersion: "N/A",
				buildDate:    "N/A",
				buildCommit:  "N/A",
			},
			expected: "Build version: N/A\nBuild date: N/A\nBuild commit: N/A\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buildVersion = tt.version.buildVersion
			buildDate = tt.version.buildDate
			buildCommit = tt.version.buildCommit

			appVersion := NewVersion()

			var buf bytes.Buffer
			appVersion.Print(&buf)

			assert.Equal(t, tt.expected, buf.String())
		})
	}
}
