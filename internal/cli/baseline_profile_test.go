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

	"roundfix/internal/baseline"
)

func TestCapabilityRecheck(t *testing.T) {
	t.Run("requires no decisions and writes nothing", func(t *testing.T) {
		repository := newHumanBaselineRepository(t)
		writeBaselinePlanTestFile(t, repository, ".roundfix/run-journal.jsonl", "existing journal entry\n")
		t.Setenv("PATH", t.TempDir())
		before := baselinePlanTestTree(t, repository)

		var stdout bytes.Buffer
		var stderr bytes.Buffer
		code := RunContext(context.Background(), []string{
			"baseline", "capabilities", "check",
			"--profile", "standard-typescript-monorepo",
			"--repo", repository,
			"--format", "json",
		}, &stdout, &stderr)
		if code != exitUnverified {
			t.Fatalf("capability re-check exit = %d, want %d stdout=%s stderr=%s",
				code, exitUnverified, stdout.String(), stderr.String())
		}
		if stderr.Len() != 0 {
			t.Fatalf("capability re-check stderr = %q, want empty", stderr.String())
		}
		var result baseline.CapabilityRecheckResult
		if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
			t.Fatalf("decode capability re-check: %v\n%s", err, stdout.String())
		}
		if result.SchemaVersion != baseline.CapabilityRecheckSchemaVersion ||
			result.Operation != "capabilities.check" || result.Profile == nil ||
			result.Profile.ID != "standard-typescript-monorepo" ||
			len(result.Capabilities) == 0 || len(result.Divergences) == 0 {
			t.Fatalf("capability re-check result = %#v", result)
		}
		var topLevel map[string]json.RawMessage
		if err := json.Unmarshal(stdout.Bytes(), &topLevel); err != nil {
			t.Fatalf("decode capability re-check fields: %v", err)
		}
		if _, resolvedDecisions := topLevel["decisions"]; resolvedDecisions {
			t.Fatalf("capability re-check unexpectedly emitted resolved decisions:\n%s", stdout.String())
		}
		if after := baselinePlanTestTree(t, repository); after != before {
			t.Fatalf("capability re-check changed repository or journal bytes: before=%s after=%s", before, after)
		}

		stdout.Reset()
		stderr.Reset()
		code = RunContext(context.Background(), []string{
			"baseline", "capabilities", "check",
			"--profile", "standard-typescript-monorepo",
			"--repo", repository,
		}, &stdout, &stderr)
		if code != exitUnverified || stderr.Len() != 0 {
			t.Fatalf("text capability re-check exit = %d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if rendered := baseline.RenderProfileDivergences(result.Divergences); !strings.Contains(stdout.String(), rendered) {
			t.Fatalf("text capability re-check did not share divergence rendering:\n%s", stdout.String())
		}
		if after := baselinePlanTestTree(t, repository); after != before {
			t.Fatalf("text capability re-check changed repository or journal bytes: before=%s after=%s", before, after)
		}
	})

	t.Run("names an unresolvable Profile", func(t *testing.T) {
		repository := newBaselineProfileTestRepository(t)
		before := baselinePlanTestTree(t, repository)
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		code := RunContext(context.Background(), []string{
			"baseline", "capabilities", "check",
			"--repo", repository,
			"--format", "json",
		}, &stdout, &stderr)
		if code != exitPreflight {
			t.Fatalf("missing-Profile re-check exit = %d, want %d stdout=%s stderr=%s",
				code, exitPreflight, stdout.String(), stderr.String())
		}
		if !strings.Contains(stderr.String(), "no resolvable Baseline Profile") {
			t.Fatalf("missing-Profile stderr = %q", stderr.String())
		}
		var result baseline.CapabilityRecheckResult
		if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
			t.Fatalf("decode missing-Profile result: %v\n%s", err, stdout.String())
		}
		if result.State != "failed" || result.Category != "profile" {
			t.Fatalf("missing-Profile result = %#v", result)
		}
		if after := baselinePlanTestTree(t, repository); after != before {
			t.Fatalf("missing-Profile re-check changed repository bytes: before=%s after=%s", before, after)
		}
	})

	t.Run("preserves JSON failures for every format flag spelling", func(t *testing.T) {
		for _, formatArgs := range [][]string{
			{"--format", "json"},
			{"--format=json"},
			{"-format", "json"},
			{"-format=json"},
		} {
			t.Run(strings.Join(formatArgs, "_"), func(t *testing.T) {
				args := []string{"baseline", "capabilities", "check"}
				args = append(args, formatArgs...)
				args = append(args, "--unknown")
				var stdout bytes.Buffer
				var stderr bytes.Buffer
				code := RunContext(context.Background(), args, &stdout, &stderr)
				if code != exitPreflight {
					t.Fatalf("invalid capability re-check exit = %d, want %d stdout=%q stderr=%q",
						code, exitPreflight, stdout.String(), stderr.String())
				}
				var result baseline.CapabilityRecheckResult
				if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
					t.Fatalf("decode structured capability failure: %v\nstdout=%q\nstderr=%q",
						err, stdout.String(), stderr.String())
				}
				if result.State != "failed" || result.Category != "preflight" {
					t.Fatalf("structured capability failure = %#v", result)
				}
			})
		}
	})
}

func TestCapabilityTextRendersProbe(t *testing.T) {
	repository := newHumanBaselineRepository(t)
	bin := t.TempDir()
	candidate := filepath.Join(bin, "rtk")
	if err := os.Symlink("missing-rtk-target", candidate); err != nil {
		t.Fatalf("create broken rtk candidate: %v", err)
	}
	t.Setenv("PATH", bin)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := RunContext(context.Background(), []string{
		"baseline", "capabilities", "check",
		"--profile", "go-cli-tui",
		"--repo", repository,
		"--format", "text",
	}, &stdout, &stderr)
	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("capability text exit = %d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	for _, want := range []string{
		"Inspected candidate: " + filepath.ToSlash(candidate) + " (broken-link)",
		"Repair the inspected candidate " + filepath.ToSlash(candidate) + " (broken-link)",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Errorf("public capability text missing %q:\n%s", want, stdout.String())
		}
	}
	if strings.Contains(stdout.String(), "Install rtk") {
		t.Errorf("public capability text recommends installing rejected candidate:\n%s", stdout.String())
	}
}

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
			args: []string{"baseline", "capabilities", "check", "--help"},
			want: []string{
				"Exit codes:",
				"1  output failure",
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
