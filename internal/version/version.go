package version

import "fmt"

var (
	Version     = "dev"
	Commit      = "none"
	Date        = "unknown"
	DefaultRepo = "https://github.com/uttam-li/rutter.git"
)

func String() string {
	return fmt.Sprintf("rutter %s (commit %s, built %s)", Version, Commit, Date)
}
