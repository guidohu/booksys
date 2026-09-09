// Package version holds the build information of the binary. Its values are
// populated at link time.
package version

import (
	"fmt"
)

// BuildDate, Commit, and Release are populated by go build -ldflags
var (
	// Release holds the version tag (e.g., "v1.2.3")
	Release = "unknown"

	// Commit holds the short git commit hash
	Commit = "unknown"

	// BuildDate holds the date the binary was built
	BuildDate = "unknown"
)

// PrintVersion writes the build information to standard output.
func PrintVersion() {
	fmt.Printf("Version: %s\n", Release)
	fmt.Printf("Commit: %s\n", Commit)
	fmt.Printf("Built: %s\n", BuildDate)
}
