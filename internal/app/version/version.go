package version

var (
	// buildVersion is the application build version, example: v1.0.0
	buildVersion = "N/A"

	// buildDate is the application build date, example: 01.02.2025
	buildDate = "N/A"

	// buildCommit is the application build commit, example: abcd1234
	buildCommit = "N/A"
)

// Version is the interface that provides application build
type Version interface {
	Version() string
	Date() string
	Commit() string
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

// BuildVersion returns the application build version
func (v *version) Version() string {
	return v.buildVersion
}

// BuildDate returns the application build date
func (v *version) Date() string {
	return v.buildDate
}

// BuildCommit returns the application build commit
func (v *version) Commit() string {
	return v.buildCommit
}
