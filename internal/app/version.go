package app

import "strings"

const Name = "roundfix"

// Version is the Roundfix semantic version. Local builds keep this product
// identity; the release workflow overrides it with the matching pushed tag via
// -ldflags "-X roundfix/internal/app.Version=<version>".
var Version = "0.0.1"

// BuildCommit and BuildTime identify a local build. The Makefile build and
// install targets stamp them via -ldflags with the short commit SHA (plus
// "-dirty" when the tree has changes) and the local build time; release
// binaries built from the tag workflow leave them empty.
var (
	BuildCommit = ""
	BuildTime   = ""
)

// VersionLine is the complete `--version` line body: the plain version for
// release builds, and the version plus build identity for stamped local
// builds, shaped like `0.0.1 (a1b2c3d, built 2026-07-15 14:32:05 -0300)`.
func VersionLine() string {
	commit := strings.TrimSpace(BuildCommit)
	builtAt := strings.TrimSpace(BuildTime)
	if commit == "" && builtAt == "" {
		return Version
	}
	details := make([]string, 0, 2)
	if commit != "" {
		details = append(details, commit)
	}
	if builtAt != "" {
		details = append(details, "built "+builtAt)
	}
	return Version + " (" + strings.Join(details, ", ") + ")"
}
