package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Version(t *testing.T) {
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
			expected: "v1.0.0",
		},
		{
			name: "Empty",
			version: version{
				buildVersion: "N/A",
				buildDate:    "N/A",
				buildCommit:  "N/A",
			},
			expected: "N/A",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buildVersion = tt.version.buildVersion
			buildDate = tt.version.buildDate
			buildCommit = tt.version.buildCommit

			appVersion := NewVersion()
			assert.Equal(t, tt.expected, appVersion.Version())
		})
	}
}

func Test_Date(t *testing.T) {
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
			expected: "01.02.2025",
		},
		{
			name: "Empty",
			version: version{
				buildVersion: "N/A",
				buildDate:    "N/A",
				buildCommit:  "N/A",
			},
			expected: "N/A",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buildVersion = tt.version.buildVersion
			buildDate = tt.version.buildDate
			buildCommit = tt.version.buildCommit

			appVersion := NewVersion()
			assert.Equal(t, tt.expected, appVersion.Date())
		})
	}
}

func Test_Commit(t *testing.T) {
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
			expected: "abcd1234",
		},
		{
			name: "Empty",
			version: version{
				buildVersion: "N/A",
				buildDate:    "N/A",
				buildCommit:  "N/A",
			},
			expected: "N/A",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buildVersion = tt.version.buildVersion
			buildDate = tt.version.buildDate
			buildCommit = tt.version.buildCommit

			appVersion := NewVersion()
			assert.Equal(t, tt.expected, appVersion.Commit())
		})
	}
}
