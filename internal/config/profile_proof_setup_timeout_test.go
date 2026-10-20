// Suite: Profile proof setup deadline configuration
// Invariant: User and Project Config override the built-in bounded setup deadline.
// Boundary IN: strict config loading and validation.
// Boundary OUT: the Agent setup context, covered in internal/agent.
package config

import (
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestProfileProofSetupTimeoutConfig(t *testing.T) {
	t.Parallel()
	const rangeError = "defaults.profile_proof_setup_timeout must be greater than 0 and at most 10m"
	for _, tt := range []struct {
		name, user, project, wantError string
		want                           time.Duration
	}{
		{name: "built-in", want: DefaultProfileProofSetupTimeout},
		{name: "user", user: "45s", want: 45 * time.Second},
		{name: "project overrides user", user: "45s", project: "90s", want: 90 * time.Second},
		{name: "maximum", project: "10m", want: MaxProfileProofSetupTimeout},
		{name: "zero", project: "0s", wantError: rangeError},
		{name: "negative", project: "-1s", wantError: rangeError},
		{name: "too long", project: "11m", wantError: rangeError},
		{name: "invalid duration", project: "soon", wantError: "invalid duration"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			home, work := t.TempDir(), t.TempDir()
			mustMkdir(t, filepath.Join(work, ".git"))
			if tt.user != "" {
				mustMkdir(t, filepath.Join(home, ".roundfix"))
				mustWrite(t, filepath.Join(home, ".roundfix", "config.yml"), "defaults:\n  profile_proof_setup_timeout: "+tt.user+"\n")
			}
			if tt.project != "" {
				mustWrite(t, filepath.Join(work, ".roundfixrc.yml"), "defaults:\n  profile_proof_setup_timeout: "+tt.project+"\n")
			}
			loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: work})
			if tt.wantError != "" {
				if err == nil {
					t.Fatal("expected invalid deadline to be refused")
				}
				if tt.wantError == rangeError && err.Error() != rangeError {
					t.Fatalf("error = %q, want %q", err, rangeError)
				}
				if !strings.Contains(err.Error(), tt.wantError) {
					t.Fatalf("error = %q, want %q", err, tt.wantError)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got := loaded.Config.Defaults.ProfileProofSetupTimeout; got != tt.want {
				t.Fatalf("timeout = %s, want %s", got, tt.want)
			}
		})
	}
}
