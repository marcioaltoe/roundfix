package app

import "strings"

const Name = "roundfix"

// Version is the Roundfix semantic version for local builds only. The
// released value comes from the tag, stamped at build time with
// -ldflags "-X roundfix/internal/app.Version=<version>".
//
// This constant is NOT the release source of truth: the release workflow
// validates the tag against dist/npm/roundfix/package.json. Bumping only
// this one fails the release at its first step with "tag does not match the
// checked-in Roundfix version" — keep the two in step, and treat the
// package manifest as the authority.
//
// Version is the Roundfix semantic version. Local builds keep this product
// identity; the release workflow overrides it with the matching pushed tag via
// -ldflags "-X roundfix/internal/app.Version=<version>".
var Version = "0.12.0"

// BuildCommit and BuildTime identify a local build. The Makefile build and
// install targets stamp them via -ldflags with the short commit SHA (plus
// "-dirty" when the tree has changes) and the local build time; release
// binaries built from the tag workflow leave them empty.
var (
	BuildCommit = ""
	BuildTime   = ""
)

// AuditingBinary identifies the Roundfix that produced a verdict, as distinct
// from the tree it audited.
type AuditingBinary struct {
	Version string
	Commit  string
	Built   string
}

// Auditor returns the running binary's build identity.
func Auditor() AuditingBinary {
	return AuditingBinary{
		Version: strings.TrimSpace(Version),
		Commit:  strings.TrimSpace(BuildCommit),
		Built:   strings.TrimSpace(BuildTime),
	}
}

// String renders the complete build identity. Released binaries remain
// readable as their version even though their commit and build time are empty.
func (binary AuditingBinary) String() string {
	version := strings.TrimSpace(binary.Version)
	commit := strings.TrimSpace(binary.Commit)
	builtAt := strings.TrimSpace(binary.Built)
	if commit == "" && builtAt == "" {
		return version
	}
	details := make([]string, 0, 2)
	if commit != "" {
		details = append(details, commit)
	}
	if builtAt != "" {
		details = append(details, "built "+builtAt)
	}
	return version + " (" + strings.Join(details, ", ") + ")"
}

// AncestryResult states what commit ancestry established about whether the
// auditing binary predates the audited tree.
type AncestryResult uint8

const (
	AncestryUnknown AncestryResult = iota
	AncestryNotOlder
	AncestryOlder
)

// Staleness is what a report says about the auditing binary's age.
type Staleness string

const (
	StalenessCurrent Staleness = "current"
	StalenessStale   Staleness = "stale"
	StalenessUnknown Staleness = "unknown"
)

// CompareToTree answers whether this binary predates the tree it audits. A
// stamped build prefers commit ancestry; otherwise the declared tree version
// is used. Missing evidence always produces StalenessUnknown.
//
// The reason never repeats the state it accompanies. A caller composes the two
// as "<state>: <reason>", and a reason that prefixed its own state produced
// "unknown: unknown: ..." in one path and a bare reason in another, leaving the
// recorded field with no stable shape to read.
// StalenessLine composes the recorded field: the state, then the reason that
// established it. Both report writers use this rather than composing their own,
// because two compositions are how one writer came to record a bare reason with
// no state while the other recorded both.
func (binary AuditingBinary) StalenessLine(treeVersion string, ancestry AncestryResult) string {
	state, reason := binary.CompareToTree(treeVersion, ancestry)
	return string(state) + ": " + reason
}

func (binary AuditingBinary) CompareToTree(treeVersion string, ancestry AncestryResult) (Staleness, string) {
	commit := strings.TrimSpace(binary.Commit)
	if commit != "" {
		switch ancestry {
		case AncestryOlder:
			return StalenessStale, "commit ancestry: build commit predates audited tree"
		case AncestryNotOlder:
			return StalenessCurrent, "commit ancestry: build commit does not predate audited tree"
		default:
			return StalenessUnknown, "commit ancestry could not be resolved for the recorded build commit"
		}
	}

	version := strings.TrimSpace(binary.Version)
	treeVersion = strings.TrimSpace(treeVersion)
	if version != "" && treeVersion != "" {
		if CompareVersions(version, treeVersion) < 0 {
			return StalenessStale, "declared tree version: auditing binary " + version + " predates " + treeVersion
		}
		return StalenessCurrent, "declared tree version: auditing binary " + version + " does not predate " + treeVersion
	}

	if treeVersion == "" {
		return StalenessUnknown, "auditing binary has no build commit and audited tree declares no version"
	}
	return StalenessUnknown, "auditing binary declares no version"
}

// VersionLine is the complete `--version` line body: the plain version for
// release builds, and the version plus build identity for stamped local
// builds, shaped like `0.0.1 (a1b2c3d, built 2026-07-15 14:32:05 -0300)`.
func VersionLine() string {
	return Auditor().String()
}
