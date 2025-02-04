package version

import (
	"fmt"
	"io"
)

var (
	// buildVersion is the application build version, example: v1.0.0
	buildVersion = "N/A"

	// buildDate is the application build date, example: 01.02.2025
	buildDate = "N/A"

	// buildCommit is the application build commit, example: abcd1234
	buildCommit = "N/A"
)

type Version interface {
	Print(w io.Writer)
}

type version struct {
	buildVersion string
	buildDate    string
	buildCommit  string
}

// NewVersion creates a new version instance
func NewVersion() Version {
	return &version{
		buildVersion: buildVersion,
		buildDate:    buildDate,
		buildCommit:  buildCommit,
	}
}

// Print prints the app version information
func (v *version) Print(w io.Writer) {
	_, _ = fmt.Fprintln(w, "Build version:", buildVersion)
	_, _ = fmt.Fprintln(w, "Build date:", buildDate)
	_, _ = fmt.Fprintln(w, "Build commit:", buildCommit)
}
