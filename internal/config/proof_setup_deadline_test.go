// Suite: Narrow profile proof setup deadline reader
// Invariant: the deadline agrees with Load while unrelated invalid keys cannot prevent a repair.
// Boundary IN: isolated User Config and Project Config at a temporary Git root.
// Boundary OUT: effective duration or file-qualified parse/read errors and the shared range error.
package config

import (
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"roundfix/internal/gittest"
)

func TestTheProofSetupDeadlineReaderAgreesWithLoad(t *testing.T) {
	t.Parallel()
	const rangeError = "defaults.profile_proof_setup_timeout must be greater than 0 and at most 10m"
	for _, tt := range []struct {
		name, user, project, wantError string
		want                           time.Duration
		loadRefuses                    bool
	}{
		{name: "built_in", want: DefaultProfileProofSetupTimeout},
		{name: "user", user: "defaults:\n  profile_proof_setup_timeout: 75s\n", want: 75 * time.Second},
		{name: "project", project: "defaults:\n  profile_proof_setup_timeout: 90s\n", want: 90 * time.Second},
		{name: "project_wins", user: "defaults:\n  profile_proof_setup_timeout: 75s\n", project: "defaults:\n  profile_proof_setup_timeout: 90s\n", want: 90 * time.Second},
		{name: "another_key_load_refuses", user: "defaults:\n  profile_proof_setup_timeout: 75s\nprofiles:\n  backend:\n    preferred: {runtime: stale-adapter, model: stale, reasoning_effort: high}\n    fallbacks:\n      - {runtime: codex, model: backup, reasoning_effort: high}\n", want: 75 * time.Second, loadRefuses: true},
		{name: "zero", project: "defaults:\n  profile_proof_setup_timeout: 0s\n", wantError: rangeError, loadRefuses: true},
		{name: "above_ten_minutes", project: "defaults:\n  profile_proof_setup_timeout: 11m\n", wantError: rangeError, loadRefuses: true},
		{name: "not_a_duration", project: "defaults:\n  profile_proof_setup_timeout: soon\n", wantError: `invalid duration "soon"`, loadRefuses: true},
		{name: "another_default_load_refuses", user: "defaults:\n  profile_proof_setup_timeout: 75s\n  agent_full_access: invalid\n", want: 75 * time.Second, loadRefuses: true},
		{name: "maximum", project: "defaults:\n  profile_proof_setup_timeout: 10m\n", want: MaxProfileProofSetupTimeout},
		{name: "negative", project: "defaults:\n  profile_proof_setup_timeout: -1s\n", wantError: rangeError, loadRefuses: true},
		{name: "project_overrides_out_of_range_user", user: "defaults:\n  profile_proof_setup_timeout: 0s\n", project: "defaults:\n  profile_proof_setup_timeout: 90s\n", want: 90 * time.Second},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			home, repo := t.TempDir(), t.TempDir()
			gittest.MarkRepository(t, repo)
			if tt.user != "" {
				mustMkdir(t, filepath.Join(home, ".roundfix"))
				mustWrite(t, filepath.Join(home, userConfigRelPath), tt.user)
			}
			if tt.project != "" {
				mustWrite(t, filepath.Join(repo, projectConfigName), tt.project)
			}
			// Start below the root to check the same Git-root lookup as Load.
			work := filepath.Join(repo, "nested")
			mustMkdir(t, work)
			got, err := ProfileProofSetupTimeout(home, work)
			if tt.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantError) {
					t.Fatalf("reader error = %v, want %q", err, tt.wantError)
				}
				if tt.wantError == rangeError && err.Error() != rangeError {
					t.Fatalf("reader error = %q, want exactly %q", err, rangeError)
				}
				if tt.name == "not_a_duration" && !strings.Contains(err.Error(), filepath.Join(repo, projectConfigName)) {
					t.Fatalf("duration error must name its file: %v", err)
				}
			} else if err != nil || got != tt.want {
				t.Fatalf("reader = %s, %v; want %s, nil", got, err, tt.want)
			}
			loaded, loadErr := Load(LoadOptions{HomeDir: home, WorkDir: work, WorktreeLocationOptional: true})
			if (loadErr != nil) != tt.loadRefuses {
				t.Fatalf("Load error = %v, want refusal = %t", loadErr, tt.loadRefuses)
			}
			if loadErr == nil && loaded.Config.Defaults.ProfileProofSetupTimeout != got {
				t.Fatalf("Load = %s, reader = %s", loaded.Config.Defaults.ProfileProofSetupTimeout, got)
			}
		})
	}
}

func TestTheProofSetupDeadlineReaderNamesUnreadableFiles(t *testing.T) {
	t.Parallel()
	for _, scope := range []string{"user", "project"} {
		t.Run(scope, func(t *testing.T) {
			t.Parallel()
			home, repo := t.TempDir(), t.TempDir()
			gittest.MarkRepository(t, repo)
			path := filepath.Join(home, userConfigRelPath)
			if scope == "project" {
				path = filepath.Join(repo, projectConfigName)
			}
			// A directory cannot be read as YAML, regardless of host permissions.
			mustMkdir(t, path)
			_, err := ProfileProofSetupTimeout(home, repo)
			if err == nil || !strings.HasPrefix(err.Error(), "read config "+strconv.Quote(path)+": ") {
				t.Fatalf("reader error = %v, want unreadable config path %q", err, path)
			}
		})
	}
}
