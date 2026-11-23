package version

import "fmt"

var (
	Version     = "dev"
	Commit      = "none"
	Date        = "unknown"
	DefaultRepo = "https://github.com/uttam-li/rutter.git"
)

// Info is the machine-readable build description.
type Info struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
	Repo    string `json:"repo"`
}

// BuildInfo returns the values baked in at build time.
func BuildInfo() Info {
	return Info{
		Name:    "rutter",
		Version: Version,
		Commit:  Commit,
		Date:    Date,
		Repo:    DefaultRepo,
	}
}

// String is the human-readable one-line version.
func String() string {
	return fmt.Sprintf("rutter %s (commit %s, built %s)", Version, Commit, Date)
}
