// Suite: Roundfix build identity
// Invariant: the auditing binary identifies itself and never infers current from missing age evidence.
// Boundary IN: build-stamp assembly and pure ancestry/version comparison.
// Boundary OUT: Git probing and QA Report serialization.
package app

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestAuditorAssemblesBuildIdentity(t *testing.T) {
	tests := []struct {
		name    string
		version string
		commit  string
		built   string
		want    AuditingBinary
		line    string
	}{
		{
			name:    "release build remains readable without local stamps",
			version: "1.2.3",
			want:    AuditingBinary{Version: "1.2.3"},
			line:    "1.2.3",
		},
		{
			name:    "local build includes every stamp",
			version: "1.2.3",
			commit:  " a1b2c3d ",
			built:   " 2026-07-15 14:32:05 -0300 ",
			want: AuditingBinary{
				Version: "1.2.3",
				Commit:  "a1b2c3d",
				Built:   "2026-07-15 14:32:05 -0300",
			},
			line: "1.2.3 (a1b2c3d, built 2026-07-15 14:32:05 -0300)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldVersion, oldCommit, oldTime := Version, BuildCommit, BuildTime
			t.Cleanup(func() {
				Version, BuildCommit, BuildTime = oldVersion, oldCommit, oldTime
			})
			Version, BuildCommit, BuildTime = tt.version, tt.commit, tt.built

			got := Auditor()
			if got != tt.want {
				t.Fatalf("Auditor() = %#v, want %#v", got, tt.want)
			}
			if got.String() != tt.line {
				t.Fatalf("Auditor().String() = %q, want %q", got.String(), tt.line)
			}
		})
	}
}

func TestCompareToTree(t *testing.T) {
	tests := []struct {
		name        string
		binary      AuditingBinary
		treeVersion string
		ancestry    AncestryResult
		want        Staleness
		wantReason  string
	}{
		{
			name:       "stale from commit ancestry",
			binary:     AuditingBinary{Version: "1.2.3", Commit: "a1b2c3d"},
			ancestry:   AncestryOlder,
			want:       StalenessStale,
			wantReason: "commit ancestry: build commit predates audited tree",
		},
		{
			name:       "current from commit ancestry",
			binary:     AuditingBinary{Version: "1.2.3", Commit: "a1b2c3d"},
			ancestry:   AncestryNotOlder,
			want:       StalenessCurrent,
			wantReason: "commit ancestry: build commit does not predate audited tree",
		},
		{
			name:        "stale released binary from declared tree version",
			binary:      AuditingBinary{Version: "1.2.2"},
			treeVersion: "1.2.3",
			want:        StalenessStale,
			wantReason:  "declared tree version: auditing binary 1.2.2 predates 1.2.3",
		},
		{
			name:        "current released binary from declared tree version",
			binary:      AuditingBinary{Version: "1.2.3"},
			treeVersion: "v1.2.3",
			want:        StalenessCurrent,
			wantReason:  "declared tree version: auditing binary 1.2.3 does not predate v1.2.3",
		},
		{
			name:        "unknown stamped build without ancestry answer",
			binary:      AuditingBinary{Version: "1.2.3", Commit: "a1b2c3d"},
			treeVersion: "1.2.3",
			ancestry:    AncestryUnknown,
			want:        StalenessUnknown,
			wantReason:  "unknown: commit ancestry is unavailable for stamped build",
		},
		{
			name:       "unknown without build stamp or declared tree version",
			binary:     AuditingBinary{Version: "1.2.3"},
			want:       StalenessUnknown,
			wantReason: "unknown: auditing binary has no build commit and audited tree declares no version",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, reason := tt.binary.CompareToTree(tt.treeVersion, tt.ancestry)
			if got != tt.want {
				t.Fatalf("CompareToTree() staleness = %q, want %q", got, tt.want)
			}
			if reason != tt.wantReason {
				t.Fatalf("CompareToTree() reason = %q, want %q", reason, tt.wantReason)
			}
		})
	}
}

func TestVersionLineIncludesBuildIdentityWhenStamped(t *testing.T) {
	tests := []struct {
		name   string
		commit string
		built  string
		want   string
	}{
		// The expectations compose the product version rather than repeating
		// it: this test characterizes how build identity is appended, not
		// which version is checked in, and a release bump must not have to
		// edit it.
		{name: "release build stays plain", commit: "", built: "", want: Version},
		{name: "commit and time", commit: "a1b2c3d", built: "2026-07-15 14:32:05 -0300", want: Version + " (a1b2c3d, built 2026-07-15 14:32:05 -0300)"},
		{name: "dirty commit only", commit: "a1b2c3d-dirty", built: "", want: Version + " (a1b2c3d-dirty)"},
		{name: "time only", commit: "", built: "2026-07-15 14:32:05 -0300", want: Version + " (built 2026-07-15 14:32:05 -0300)"},
		{name: "whitespace stamps stay plain", commit: "  ", built: "\t", want: Version},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldCommit, oldTime := BuildCommit, BuildTime
			t.Cleanup(func() {
				BuildCommit, BuildTime = oldCommit, oldTime
			})
			BuildCommit, BuildTime = tt.commit, tt.built
			if got := VersionLine(); got != tt.want {
				t.Fatalf("VersionLine() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestVersionMatchesTheReleaseManifest keeps the local constant in step with
// the manifest the release workflow validates the tag against. The 0.3.1
// release failed its first step because only the constant had been bumped;
// this test turns that into a red suite instead of a red release.
func TestVersionMatchesTheReleaseManifest(t *testing.T) {
	t.Parallel()
	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller could not locate the repository")
	}
	manifestPath := filepath.Join(filepath.Dir(testFile), "..", "..", "dist", "npm", "roundfix", "package.json")
	content, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read release manifest: %v", err)
	}
	var manifest struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(content, &manifest); err != nil {
		t.Fatalf("parse release manifest: %v", err)
	}
	if manifest.Version != Version {
		t.Fatalf("dist/npm/roundfix/package.json version = %q, but internal/app.Version = %q; the release workflow validates the tag against the manifest, so both must match", manifest.Version, Version)
	}
}
