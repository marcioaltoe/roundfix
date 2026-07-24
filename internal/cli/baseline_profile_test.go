// Suite: Baseline Profile CLI
// Invariant: the public profile commands expose deterministic repository-only Baseline Profile behavior without changing Agent Selection Profiles.
// Boundary IN: command dispatch, flags, help, text and JSON output, exit categories, and repository files.
// Boundary OUT: Baseline planning/apply and Agent Selection Profile configuration or proof.

package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBaselineProfileCommandInitShowAndValidate(t *testing.T) {
	repo := newBaselineProfileTestRepository(t)
	t.Chdir(repo)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := RunContext(context.Background(), []string{
		"baseline", "profile", "init",
		"--id", "team-go",
		"--from", "go-cli-tui",
	}, &stdout, &stderr)
	if code != exitOK {
		t.Fatalf("profile init exit = %d, want 0 stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("profile init stderr = %q, want empty", stderr.String())
	}
	if !strings.Contains(stdout.String(), "Created repository-owned Baseline Profile team-go") ||
		!strings.Contains(stdout.String(), ".roundfix/baseline/profiles/team-go.json") {
		t.Fatalf("profile init stdout = %q", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	code = RunContext(context.Background(), []string{
		"baseline", "profile", "show", "team-go", "--format", "text",
	}, &stdout, &stderr)
	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("profile show text exit = %d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	text := stdout.String()
	for _, want := range []string{
		"Baseline Profile: team-go",
		"Source: repository",
		"Catalog schema: roundfix/baseline-catalog/v1",
		"Modules:",
		"Decisions:",
		"Templates:",
		"Digest: sha256:",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("profile show text missing %q:\n%s", want, text)
		}
	}

	stdout.Reset()
	stderr.Reset()
	code = RunContext(context.Background(), []string{
		"baseline", "profile", "show", "--format=json", "team-go",
	}, &stdout, &stderr)
	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("profile show JSON exit = %d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	var show baselineProfileResult
	if err := json.Unmarshal(stdout.Bytes(), &show); err != nil {
		t.Fatalf("decode profile show JSON: %v\n%s", err, stdout.String())
	}
	if show.SchemaVersion != baselineResultSchema || show.Operation != "profile.show" ||
		show.State != "valid" || show.Profile == nil || show.Profile.ID != "team-go" {
		t.Fatalf("profile show JSON = %#v", show)
	}
	if show.Profile.Path != ".roundfix/baseline/profiles/team-go.json" {
		t.Fatalf("profile show JSON path = %q", show.Profile.Path)
	}
	if !strings.Contains(text, show.Profile.Digest) {
		t.Fatalf("text and JSON resolved different profile digests: text=%q JSON=%#v", text, show.Profile)
	}

	stdout.Reset()
	stderr.Reset()
	code = RunContext(context.Background(), []string{
		"baseline", "profile", "validate",
		".roundfix/baseline/profiles/team-go.json",
		"--format", "json",
	}, &stdout, &stderr)
	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("profile validate exit = %d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	var validation baselineProfileResult
	if err := json.Unmarshal(stdout.Bytes(), &validation); err != nil {
		t.Fatalf("decode profile validation JSON: %v\n%s", err, stdout.String())
	}
	if validation.SchemaVersion != baselineResultSchema || validation.Operation != "profile.validate" ||
		validation.State != "valid" || len(validation.Profiles) != 1 ||
		validation.Profiles[0].Digest != show.Profile.Digest {
		t.Fatalf("profile validation JSON = %#v", validation)
	}
}

func TestBaselineProfileCommandRejectsInvalidAndNonRepositoryProfiles(t *testing.T) {
	repo := newBaselineProfileTestRepository(t)
	t.Chdir(repo)

	userHome := t.TempDir()
	t.Setenv("HOME", userHome)
	userProfileDir := filepath.Join(userHome, ".roundfix", "baseline", "profiles")
	if err := os.MkdirAll(userProfileDir, 0o755); err != nil {
		t.Fatalf("create user profile directory: %v", err)
	}
	userProfilePath := filepath.Join(userProfileDir, "user-only.json")
	if err := os.WriteFile(userProfilePath, []byte(`{}`), 0o644); err != nil {
		t.Fatalf("write user profile: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := RunContext(context.Background(), []string{
		"baseline", "profile", "validate", "--format", "json",
	}, &stdout, &stderr)
	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("empty repository validation exit = %d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	var empty baselineProfileResult
	if err := json.Unmarshal(stdout.Bytes(), &empty); err != nil {
		t.Fatalf("decode empty validation JSON: %v", err)
	}
	if len(empty.Profiles) != 0 {
		t.Fatalf("user-scoped profile participated in discovery: %#v", empty.Profiles)
	}

	stdout.Reset()
	stderr.Reset()
	code = RunContext(context.Background(), []string{
		"baseline", "profile", "validate", userProfilePath, "--format", "json",
	}, &stdout, &stderr)
	if code != exitPreflight {
		t.Fatalf("outside profile validation exit = %d, want %d", code, exitPreflight)
	}
	if !strings.Contains(stderr.String(), "unsafe custom Baseline Profile path") {
		t.Fatalf("outside profile stderr = %q", stderr.String())
	}
	var failure baselineProfileResult
	if err := json.Unmarshal(stdout.Bytes(), &failure); err != nil {
		t.Fatalf("decode profile failure JSON: %v\n%s", err, stdout.String())
	}
	if failure.SchemaVersion != baselineResultSchema || failure.State != "failed" ||
		failure.Category != "preflight" || failure.NextAction == "" {
		t.Fatalf("profile failure JSON = %#v", failure)
	}

	stdout.Reset()
	stderr.Reset()
	code = RunContext(context.Background(), []string{
		"baseline", "profile", "init", "--id", "team", "--from", "go-cli-tui",
	}, &stdout, &stderr)
	if code != exitOK {
		t.Fatalf("profile init exit = %d stderr=%q", code, stderr.String())
	}
	profilePath := filepath.Join(repo, ".roundfix", "baseline", "profiles", "team.json")
	content, err := os.ReadFile(profilePath)
	if err != nil {
		t.Fatalf("read initialized profile: %v", err)
	}
	invalid := strings.Replace(string(content), `"values":`, `"remote": "https://example.test/profile.json", "values":`, 1)
	if err := os.WriteFile(profilePath, []byte(invalid), 0o644); err != nil {
		t.Fatalf("write invalid profile: %v", err)
	}

	stdout.Reset()
	stderr.Reset()
	code = RunContext(context.Background(), []string{
		"baseline", "profile", "validate", "team", "--format=json",
	}, &stdout, &stderr)
	if code != exitPreflight || !strings.Contains(stderr.String(), "custom.profile.field.unknown") {
		t.Fatalf("invalid profile exit = %d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
}

func TestBaselineProfileHelpContract(t *testing.T) {
	tests := []struct {
		args []string
		want []string
	}{
		{
			args: []string{"baseline", "--help"},
			want: []string{
				"roundfix baseline profile init --id <id> [--from <built-in-id>]",
				"roundfix baseline profile show <id> [--format <text|json>]",
				"roundfix baseline profile validate [<id>|<path>] [--format <text|json>]",
			},
		},
		{
			args: []string{"baseline", "profile", "--help"},
			want: []string{"init", "show", "validate"},
		},
		{
			args: []string{"baseline", "profile", "init", "--help"},
			want: []string{"--id", "--from", "default go-cli-tui"},
		},
		{
			args: []string{"baseline", "profile", "show", "--help"},
			want: []string{"<id>", "--format", "text or json"},
		},
		{
			args: []string{"baseline", "profile", "validate", "--help"},
			want: []string{"[<id>|<path>]", "--format", "user-scoped profile catalog"},
		},
	}
	for _, tt := range tests {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			code := RunContext(context.Background(), tt.args, &stdout, &stderr)
			if code != exitOK || stderr.Len() != 0 {
				t.Fatalf("help exit = %d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
			}
			for _, want := range tt.want {
				if !strings.Contains(stdout.String(), want) {
					t.Fatalf("help missing %q:\n%s", want, stdout.String())
				}
			}
		})
	}

	profileHelp := commandUsage("baseline profile")
	for _, forbidden := range []string{"configure", "plan", "apply", "assets sync", "skills restore"} {
		if strings.Contains(profileHelp, forbidden) {
			t.Fatalf("baseline profile help advertised unimplemented command %q:\n%s", forbidden, profileHelp)
		}
	}

	agentProfilesHelp := commandUsage("profiles")
	for _, want := range []string{"Agent Selection Profiles", "configure", "validate"} {
		if !strings.Contains(agentProfilesHelp, want) {
			t.Fatalf("Agent Selection Profile help changed, missing %q:\n%s", want, agentProfilesHelp)
		}
	}
	if strings.Contains(agentProfilesHelp, ".roundfix/baseline/profiles") {
		t.Fatalf("Agent Selection Profile help was mixed with Baseline Profiles:\n%s", agentProfilesHelp)
	}
}

func newBaselineProfileTestRepository(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("create Git marker: %v", err)
	}
	return repo
}
