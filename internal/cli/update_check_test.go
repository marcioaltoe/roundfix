// Suite: read-only Update Check public command
// Invariant: reports preserve exact streams and cached lookup semantics without changing the repository.
// Boundary IN: command dispatch, configuration, cache, comparison, and real in-process Baseline preview.
// Boundary OUT: GitHub transport, replaced by the command-context release lister.
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"testing/synctest"
	"time"

	"roundfix/internal/app"
	"roundfix/internal/gittest"
)

func TestUpdateCheckTranscripts(t *testing.T) {
	t.Parallel()
	const outside = "Baseline: not checked (outside a Git repository).\n"
	const notify = "Next: apply it with the Roundfix skill's --update mode; update.mode notify applies nothing on its own.\n"
	const auto = "Next: update.mode auto: apply it now with the Roundfix skill's --update mode.\n"
	releases := []app.Release{
		{Version: "99.0.0", Titles: []string{"The first fixture release"}, URL: "https://github.test/99"},
		{Version: "98.0.0", Titles: []string{"The second fixture release", "Its second title"}, URL: "https://github.test/98"},
	}
	changes := "Changes since 1.0.0:\n- 99.0.0: The first fixture release\n- 98.0.0: The second fixture release; Its second title\n"
	truncated := make([]app.Release, app.ReleaseListingPageSize)
	truncatedText := "Roundfix 1.0.0 → 99.0.0 is available (update.mode notify).\nChanges since 1.0.0:\n"
	for i := range truncated {
		version := fmt.Sprintf("%d.0.0", 99-i)
		truncated[i] = app.Release{Version: version, Titles: []string{"Title"}}
		truncatedText += fmt.Sprintf("- %s: Title\n", version)
	}
	truncatedText += "- earlier releases: https://github.com/marcioaltoe/roundfix/releases\n" + outside + notify
	for _, tt := range []struct {
		name, installed, mode string
		releases              []app.Release
		err                   error
		args                  []string
		want                  string
	}{
		{name: "available notify", installed: "v1.0.0", releases: releases, want: "Roundfix 1.0.0 → 99.0.0 is available (update.mode notify).\n" + changes + outside + notify},
		{name: "available auto", installed: "1.0.0", mode: "auto", releases: releases, want: "Roundfix 1.0.0 → 99.0.0 is available (update.mode auto).\n" + changes + outside + auto},
		{name: "current", installed: "99.0.0", releases: releases, want: "Roundfix 99.0.0 is current (latest release 99.0.0; update.mode notify).\n" + outside},
		{name: "ahead of latest", installed: "100.0.0", releases: releases, want: "Roundfix 100.0.0 is current (latest release 99.0.0; update.mode notify).\n" + outside},
		{name: "unknown", installed: "1.0.0", err: errors.New("error connecting\n  to api.github.com"), want: "Roundfix 1.0.0: latest release not reached (error connecting to api.github.com); the session continues.\n" + outside},
		{name: "no releases", installed: "1.0.0", err: fmt.Errorf("listing: %w", app.ErrNoReleases), want: "Roundfix 1.0.0: latest release not reached (no releases published); the session continues.\n" + outside},
		{name: "empty listing", installed: "1.0.0", want: "Roundfix 1.0.0: latest release not reached (no releases published); the session continues.\n" + outside},
		{name: "bounded error", installed: "1.0.0", err: errors.New(strings.Repeat("é", 301)), want: "Roundfix 1.0.0: latest release not reached (" + strings.Repeat("é", 300) + "); the session continues.\n" + outside},
		{name: "development", installed: "0.0.0-dev", want: "Roundfix 0.0.0-dev is a development build; no release comparison (update.mode notify).\n" + outside},
		{name: "truncated", installed: "1.0.0", releases: truncated, want: truncatedText},
		{name: "untitled release", installed: "98.0.0", releases: []app.Release{{Version: "99.0.0", URL: "https://github.test/99"}}, want: "Roundfix 98.0.0 → 99.0.0 is available (update.mode notify).\nChanges since 98.0.0:\n- 99.0.0: https://github.test/99\n" + outside + notify},
		{name: "json", installed: "1.0.0", releases: releases, args: []string{"--format", "json"}, want: `{"schemaVersion":"roundfix/update-check/v1","mode":"notify","installed":"1.0.0","latest":"99.0.0","state":"update_available","lookup":{"source":"network","checkedAt":"2026-11-06T15:00:00Z","error":""},"releases":[{"version":"99.0.0","titles":["The first fixture release"],"url":"https://github.test/99"},{"version":"98.0.0","titles":["The second fixture release","Its second title"],"url":"https://github.test/98"}],"releasesTruncated":false,"baseline":{"state":"not_checked","fileChanges":0,"historyUnits":0,"skillChanges":0,"message":""},"next":"apply"}` + "\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			home := t.TempDir()
			setCommandEnvironmentForTest(t, home, t.TempDir())
			if tt.mode != "" {
				writeUpdateCheckUserConfig(t, home, "update: {mode: "+tt.mode+"}\n")
			}
			calls := 0
			configureUpdateCheckTest(t, tt.installed, func(ctx context.Context) ([]app.Release, error) {
				calls++
				assertUpdateCheckDeadline(t, ctx)
				return tt.releases, tt.err
			})
			stdout := runUpdateCheckTest(t, append([]string{"update", "check"}, tt.args...)...)
			if stdout != tt.want {
				t.Fatalf("stdout:\n%s\nwant:\n%s", stdout, tt.want)
			}
			if tt.name == "development" {
				if calls != 0 {
					t.Fatalf("development lookup calls = %d", calls)
				}
				if _, err := os.Stat(filepath.Join(home, ".roundfix")); !errors.Is(err, os.ErrNotExist) {
					t.Fatalf("development wrote cache: %v", err)
				}
			} else if calls != 1 {
				t.Fatalf("lookup calls = %d", calls)
			}
		})
	}
}

func TestUpdateCheckReadsTheNetworkAtMostOncePerDay(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name            string
		advance         time.Duration
		refresh, failed bool
		wantCalls       int
	}{
		{name: "second check within a day", advance: 24*time.Hour - time.Nanosecond, wantCalls: 1},
		{name: "a day later", advance: 24 * time.Hour, wantCalls: 2},
		{name: "refresh", advance: time.Hour, refresh: true, wantCalls: 2},
		{name: "failed lookup", advance: time.Hour, failed: true, wantCalls: 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			home := t.TempDir()
			setCommandEnvironmentForTest(t, home, t.TempDir())
			now := updateCheckTestTime()
			calls := 0
			configureUpdateCheckTest(t, "1.0.0", func(ctx context.Context) ([]app.Release, error) {
				calls++
				assertUpdateCheckDeadline(t, ctx)
				if tt.failed {
					return nil, errors.New("offline")
				}
				return []app.Release{{Version: "2.0.0", Titles: []string{"A change"}}}, nil
			})
			updateCommandDependenciesForTest(t, func(deps *commandDependencies) { deps.updateCheck.now = func() time.Time { return now } })
			first := runUpdateCheckTest(t, "update", "check", "--format", "json")
			assertUpdateCheckLookup(t, first, "network", tt.failed)
			cachePath := filepath.Join(home, ".roundfix", updateCheckCacheFile)
			content, err := os.ReadFile(cachePath)
			if err != nil {
				t.Fatal(err)
			}
			var cache updateCheckCache
			if err := json.Unmarshal(content, &cache); err != nil {
				t.Fatal(err)
			}
			if cache.SchemaVersion != updateCheckCacheSchema || !cache.CheckedAt.Equal(now) {
				t.Fatalf("cache = %+v", cache)
			}
			if info, err := os.Stat(cachePath); err != nil || info.Mode().Perm() != 0o644 {
				t.Fatalf("cache permissions: %v, %v", info, err)
			}
			now = now.Add(tt.advance)
			args := []string{"update", "check", "--format", "json"}
			if tt.refresh {
				args = append(args, "--refresh")
			}
			second := runUpdateCheckTest(t, args...)
			source := "cache"
			if tt.wantCalls == 2 {
				source = "network"
			}
			assertUpdateCheckLookup(t, second, source, tt.failed)
			if calls != tt.wantCalls {
				t.Fatalf("lookup calls = %d, want %d", calls, tt.wantCalls)
			}
			if tt.failed && len(cache.Releases) != 0 {
				t.Fatalf("failed cache releases = %+v", cache.Releases)
			}
		})
	}
}

func TestUpdateCheckSummarizesTheBaselinePreview(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"current", "plan ready", "no manifest", "outside repository", "writes nothing", "history and skill counts"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var repo string
			switch name {
			case "current":
				repo = newBaselineUpdateRepository(t)
			case "plan ready", "writes nothing", "history and skill counts":
				repo = staleBaselineUpdateRepository(t)
			case "no manifest":
				repo = newBaselinePlanTestRepository(t)
			default:
				repo = t.TempDir()
			}
			if name == "history and skill counts" {
				required := installRequiredExternalSkills(t, repo)
				removeInstalledSkill(t, repo, required[0], true)
				writeBaselinePlanTestFile(t, repo, ".agents/skills/roundfix/SKILL.md", "---\nname: roundfix\nversion: 0.0.0\n---\n# Roundfix\n")
				writeBaselinePlanTestFile(t, repo, ".agents/skills/council/SKILL.md", "# Retired fixture\n")
				historyLegacyFolder(t, historyFixture{repo: repo}, "0001-fixture")
				commitBaselinePlanTestRepository(t, repo)
			}
			// Resolve from a subdirectory so the preview must use the Git root.
			workDir := filepath.Join(repo, "nested")
			if err := os.Mkdir(workDir, 0o755); err != nil {
				t.Fatal(err)
			}
			setCommandEnvironmentForTest(t, t.TempDir(), workDir)
			configureUpdateCheckTest(t, "99.0.0", func(context.Context) ([]app.Release, error) {
				return []app.Release{{Version: "99.0.0", Titles: []string{}}}, nil
			})
			before := baselinePlanTestTree(t, repo)
			var statusBefore string
			if name != "outside repository" {
				statusBefore = gittest.Run(t, repo, "status", "--porcelain=v1", "--untracked-files=all")
			}
			expected := updateCheckBaseline{State: "not_checked"}
			var preview baselineUpdateResult
			if name != "outside repository" {
				environment := commandEnvironmentForTest(t)
				ctx := contextWithCommandDependencies(t.Context(), environment.dependencies)
				var output bytes.Buffer
				runBaselineUpdateCommand(ctx, []string{"--repo", repo, "--format", "json"}, &output, io.Discard, environment)
				if err := json.Unmarshal(output.Bytes(), &preview); err != nil {
					t.Fatal(err)
				}
				expected = updateCheckBaseline{State: preview.State, Message: strings.Join(strings.Fields(preview.Message), " "), FileChanges: len(preview.FileChanges), HistoryUnits: len(preview.History.Units), SkillChanges: len(preview.Skills.Outdated) + len(preview.Skills.Missing) + len(preview.Skills.Retired)}
			}
			var report updateCheckReport
			output := runUpdateCheckTest(t, "update", "check", "--format", "json")
			if err := json.Unmarshal([]byte(output), &report); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(report.Baseline, expected) {
				t.Fatalf("baseline = %+v, want %+v", report.Baseline, expected)
			}
			line := "Baseline: not checked (outside a Git repository).\n"
			switch name {
			case "current":
				if report.Baseline.State != "current" {
					t.Fatalf("current baseline: %+v", report.Baseline)
				}
				line = "Baseline: current with the installed Roundfix.\n"
			case "plan ready", "writes nothing", "history and skill counts":
				if report.Baseline.State != "plan_ready" || report.Baseline.FileChanges == 0 {
					t.Fatalf("stale baseline: %+v", report.Baseline)
				}
				if name == "history and skill counts" && (len(preview.History.Units) != 1 || len(preview.Skills.Outdated) != 1 || len(preview.Skills.Missing) != 1 || len(preview.Skills.Retired) != 1 || report.Baseline.SkillChanges != 3) {
					t.Fatalf("fixture must exercise all counts: baseline=%+v history=%+v skills=%+v", report.Baseline, preview.History, preview.Skills)
				}
				line = fmt.Sprintf("Baseline: plan ready with the installed Roundfix: %d file change(s), %d history unit(s), %d skill change(s).\n", expected.FileChanges, expected.HistoryUnits, expected.SkillChanges)
			case "no manifest":
				if report.Baseline.State != "action_required" {
					t.Fatalf("no manifest: %+v", report.Baseline)
				}
				line = "Baseline: action_required: the repository has no Setup Manifest and cannot be updated\n"
			}
			want := "Roundfix 99.0.0 is current (latest release 99.0.0; update.mode notify).\n" + line
			if report.Baseline.State == "plan_ready" {
				if report.Next != "apply" {
					t.Fatalf("next = %q", report.Next)
				}
				want += "Next: apply it with the Roundfix skill's --update mode; update.mode notify applies nothing on its own.\n"
			} else if report.Next != "none" {
				t.Fatalf("next = %q", report.Next)
			}
			if text := runUpdateCheckTest(t, "update", "check"); text != want {
				t.Fatalf("text = %s, want %s", text, want)
			}
			if after := baselinePlanTestTree(t, repo); after != before {
				t.Fatal("check changed repository files or modes")
			}
			if name != "outside repository" {
				if after := gittest.Run(t, repo, "status", "--porcelain=v1", "--untracked-files=all"); after != statusBefore {
					t.Fatalf("status changed: %s -> %s", statusBefore, after)
				}
			}
		})
	}
}

func TestUpdateCheckUsage(t *testing.T) {
	t.Parallel()
	const help = `Usage:
  roundfix update check [--refresh] [--format <text|json>]

Reports, read-only, whether a newer Roundfix release exists, the title of each
release since the installed version, and this repository's Baseline plan with
the installed binary. The release listing comes from the GitHub CLI and is
cached for 24 hours in Roundfix Home; a lookup that fails is reported in one
line. Every report exits 0. update.mode in Project Config or User Config
selects notify (the default) or auto, which the Next: line names.

Commands:
  check  Report the latest release, what it changes, and this repository's Baseline plan

Options:
  --refresh  Read the release listing now instead of the cached one
  --format   Output format: text (default) or json
`
	for _, tt := range []struct {
		name            string
		args            []string
		reason, command string
		invalidConfig   bool
	}{
		{name: "bare group", args: []string{"update"}, reason: `expected subcommand "check"`, command: "update"},
		{name: "unknown subcommand", args: []string{"update", "apply"}, reason: `expected subcommand "check"`, command: "update"},
		{name: "unexpected argument", args: []string{"update", "check", "extra"}, reason: `unexpected argument "extra"`, command: "update check"},
		{name: "invalid format", args: []string{"update", "check", "--format", "yaml"}, reason: `format must be text or json, got "yaml"`, command: "update check"},
		{name: "invalid configuration", args: []string{"update", "check"}, reason: `update.mode must be notify or auto, got "always"`, command: "update check", invalidConfig: true},
		{name: "help", args: []string{"update", "check", "--help"}},
		{name: "group help", args: []string{"update", "--help"}},
		{name: "short help", args: []string{"update", "check", "-h"}},
		{name: "short group help", args: []string{"update", "-h"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			home := t.TempDir()
			setCommandEnvironmentForTest(t, home, t.TempDir())
			configureUpdateCheckTest(t, "1.0.0", func(context.Context) ([]app.Release, error) { t.Fatal("usage reached release lookup"); return nil, nil })
			if tt.invalidConfig {
				writeUpdateCheckUserConfig(t, home, "update: {mode: always}\n")
			}
			var stdout, stderr bytes.Buffer
			code := runCLI(t, tt.args, &stdout, &stderr)
			if tt.reason == "" {
				if code != 0 || stderr.Len() != 0 || stdout.String() != help {
					t.Fatalf("help exit=%d stdout=%s stderr=%s", code, &stdout, &stderr)
				}
				return
			}
			if tt.invalidConfig {
				tt.reason = fmt.Sprintf("parse config %q: %s", filepath.Join(home, ".roundfix", "config.yml"), tt.reason)
			}
			want := "Preflight failed\n\nReason:\n  " + tt.reason + "\n\nNo side effects:\n  Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.\n\nUsage:\n  Run 'roundfix " + tt.command + " --help' for usage.\n"
			if code != exitPreflight || stdout.Len() != 0 || stderr.String() != want {
				t.Fatalf("usage exit=%d stdout=%s stderr=%s\nwant stderr=%s", code, &stdout, &stderr, want)
			}
			if _, err := os.Stat(filepath.Join(home, ".roundfix", updateCheckCacheFile)); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("usage wrote cache: %v", err)
			}
		})
	}
}

func TestUpdateCheckCacheFailuresAndDevelopmentIsolation(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"corrupt", "schema mismatch", "future timestamp", "read and write failure", "development ignores cache"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			home := t.TempDir()
			setCommandEnvironmentForTest(t, home, t.TempDir())
			path := filepath.Join(home, ".roundfix", updateCheckCacheFile)
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			cache := updateCheckCache{SchemaVersion: updateCheckCacheSchema, CheckedAt: updateCheckTestTime(), Releases: []app.Release{{Version: "7.0.0", Titles: []string{}}}}
			installed := "1.0.0"
			switch name {
			case "schema mismatch":
				cache.SchemaVersion = "old"
			case "future timestamp":
				cache.CheckedAt = cache.CheckedAt.Add(time.Nanosecond)
			case "development ignores cache":
				installed = "0.0.0-dev"
			}
			content, err := json.Marshal(cache)
			if err != nil {
				t.Fatal(err)
			}
			switch name {
			case "corrupt":
				content = []byte("invalid json")
			case "read and write failure":
				if err := os.Mkdir(path, 0o755); err != nil {
					t.Fatal(err)
				}
			}
			if name != "read and write failure" {
				if err := os.WriteFile(path, content, 0o644); err != nil {
					t.Fatal(err)
				}
			}
			calls := 0
			configureUpdateCheckTest(t, installed, func(context.Context) ([]app.Release, error) {
				calls++
				return []app.Release{{Version: "2.0.0", Titles: []string{}}}, nil
			})
			output := runUpdateCheckTest(t, "update", "check", "--format", "json")
			var report updateCheckReport
			if err := json.Unmarshal([]byte(output), &report); err != nil {
				t.Fatal(err)
			}
			if name == "development ignores cache" {
				if calls != 0 || report.State != "development" || report.Latest != "" || report.Lookup.Source != "skipped" || strings.Contains(output, `"checkedAt"`) || !reflect.DeepEqual(report.Releases, []app.Release{}) {
					t.Fatalf("development read cached comparison: calls=%d report=%s", calls, output)
				}
				after, err := os.ReadFile(path)
				if err != nil || !bytes.Equal(after, content) {
					t.Fatalf("development changed cache: %s %v", after, err)
				}
			} else if calls != 1 || report.Latest != "2.0.0" || report.Lookup.Source != "network" {
				t.Fatalf("bad cache blocked network report: calls=%d report=%s", calls, output)
			}
		})
	}
}

func TestUpdateCheckLookupTimeoutContinues(t *testing.T) {
	t.Parallel()
	home, workDir := t.TempDir(), t.TempDir()
	setCommandEnvironmentForTest(t, home, workDir)
	synctest.Test(t, func(t *testing.T) {
		started := time.Now()
		configureUpdateCheckTest(t, "1.0.0", func(ctx context.Context) ([]app.Release, error) {
			<-ctx.Done()
			return nil, ctx.Err()
		})
		output := runUpdateCheckTest(t, "update", "check")
		want := "Roundfix 1.0.0: latest release not reached (context deadline exceeded); the session continues.\nBaseline: not checked (outside a Git repository).\n"
		if output != want || time.Since(started) != 5*time.Second {
			t.Fatalf("timeout report after %s: %s", time.Since(started), output)
		}
		cached := runUpdateCheckTest(t, "update", "check")
		if cached != want || time.Since(started) != 5*time.Second {
			t.Fatalf("timed-out lookup retried: %s", cached)
		}
	})
}

func updateCheckTestTime() time.Time { return time.Date(2026, 11, 6, 15, 0, 0, 0, time.UTC) }

func configureUpdateCheckTest(t *testing.T, installed string, lister func(context.Context) ([]app.Release, error)) {
	t.Helper()
	updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
		deps.updateCheck = updateCheckDependencies{now: updateCheckTestTime, currentVersion: func() string { return installed }, listReleases: lister}
	})
}

func runUpdateCheckTest(t *testing.T, args ...string) string {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := runCLIContext(t, t.Context(), args, &stdout, &stderr)
	if code != 0 || stderr.Len() != 0 {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, &stdout, &stderr)
	}
	return stdout.String()
}

func writeUpdateCheckUserConfig(t *testing.T, home, content string) {
	t.Helper()
	path := filepath.Join(home, ".roundfix", "config.yml")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func assertUpdateCheckDeadline(t *testing.T, ctx context.Context) {
	t.Helper()
	deadline, ok := ctx.Deadline()
	if !ok || time.Until(deadline) > 5*time.Second {
		t.Fatalf("lookup deadline = %v, present %t", deadline, ok)
	}
}

func assertUpdateCheckLookup(t *testing.T, output, source string, failed bool) {
	t.Helper()
	var report updateCheckReport
	if err := json.Unmarshal([]byte(output), &report); err != nil {
		t.Fatal(err)
	}
	if report.Lookup.Source != source || (report.Lookup.Error != "") != failed {
		t.Fatalf("lookup = %+v", report.Lookup)
	}
	if failed && (report.State != "unknown" || report.Latest != "" || len(report.Releases) != 0 || report.Next != "none") {
		t.Fatalf("failed report = %+v", report)
	}
}
