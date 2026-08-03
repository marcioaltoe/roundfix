package app

import "testing"

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
