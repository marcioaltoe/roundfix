package cli

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"roundfix/internal/agent"
	"roundfix/internal/codex"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/daemon"
	roundnotify "roundfix/internal/notify"
	"roundfix/internal/preflight"
	"roundfix/internal/reviewsource"
	"roundfix/internal/rounds"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/store"
	roundtui "roundfix/internal/tui"
	"roundfix/internal/watch"
	runworktree "roundfix/internal/worktree"
)

func init() {
	newOutcomeNotifier = func(roundconfig.Config) roundnotify.Notifier {
		return testNoopOutcomeNotifier{}
	}
}

type testNoopOutcomeNotifier struct{}

func (testNoopOutcomeNotifier) Notify(context.Context, roundnotify.Outcome) error {
	return nil
}

func TestRunHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"--help"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !strings.Contains(stdout.String(), "roundfix watch") {
		t.Fatalf("expected help output, got %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "roundfix runs\n") {
		t.Fatalf("expected help output to list bare runs command, got %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "roundfix runs list") {
		t.Fatalf("expected help output to list runs command, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunVersion(t *testing.T) {
	for _, args := range [][]string{{"--version"}, {"version"}, {"-v"}} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run(args, &stdout, &stderr)

			if code != 0 {
				t.Fatalf("expected exit code 0, got %d", code)
			}
			if !strings.HasPrefix(stdout.String(), "roundfix ") {
				t.Fatalf("expected version output, got %q", stdout.String())
			}
			if stderr.Len() != 0 {
				t.Fatalf("expected no stderr, got %q", stderr.String())
			}
		})
	}
}

func TestRunInitCreatesProjectConfigWithExplicitScope(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"init", "--scope", "project"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
	configPath := filepath.Join(repoDir, ".roundfixrc.yml")
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("expected Project Config %s: %v", configPath, err)
	}
	if !strings.Contains(stdout.String(), "Roundfix config created") || !strings.Contains(stdout.String(), "roundfix fetch --pr <number>") {
		t.Fatalf("expected init success output, got %q", stdout.String())
	}
}

func TestRunInitCreatesUserConfigWithExplicitScope(t *testing.T) {
	homeDir, _ := withCLIWorkspace(t)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"init", "--scope", "user"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	configPath := filepath.Join(homeDir, ".roundfix", "config.yml")
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("expected User Config %s: %v", configPath, err)
	}
	if !strings.Contains(stdout.String(), configPath) {
		t.Fatalf("expected user config path in stdout, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunInitPromptsForScopeAndDefaultsProject(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	prompted := false
	withInitScopePrompt(t, func(context.Context, io.Writer) (string, error) {
		prompted = true
		return roundconfig.InitScopeProject, nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"init"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !prompted {
		t.Fatal("expected init to prompt for missing scope")
	}
	if _, err := os.Stat(filepath.Join(repoDir, ".roundfixrc.yml")); err != nil {
		t.Fatalf("expected default Project Config: %v", err)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr from stubbed prompt, got %q", stderr.String())
	}
}

func TestReadInitScopeDefaultsProjectOnBlankInput(t *testing.T) {
	var stderr bytes.Buffer

	scope, err := readInitScope(context.Background(), strings.NewReader("\n"), &stderr)

	if err != nil {
		t.Fatalf("expected blank scope to default, got %v", err)
	}
	if scope != roundconfig.InitScopeProject {
		t.Fatalf("expected project default, got %q", scope)
	}
	if !strings.Contains(stderr.String(), "Scope [project]") {
		t.Fatalf("expected prompt output, got %q", stderr.String())
	}
}

func TestRunInitRejectsExistingConfigWithoutForce(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	configPath := filepath.Join(repoDir, ".roundfixrc.yml")
	mustWrite(t, configPath, "defaults:\n  agent: claude\n")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"init", "--scope", "project"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "--force") {
		t.Fatalf("expected force guidance, got %q", stderr.String())
	}
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	if !strings.Contains(string(content), "agent: claude") {
		t.Fatalf("expected existing config to remain, got %s", string(content))
	}
}

func TestRunInitForceOverwritesExistingConfig(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	configPath := filepath.Join(repoDir, ".roundfixrc.yml")
	mustWrite(t, configPath, "defaults:\n  agent: claude\n")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"init", "--scope", "project", "--force"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	if !strings.Contains(string(content), "agent: codex") || strings.Contains(string(content), "agent: claude") {
		t.Fatalf("expected generated config to replace old content, got %s", string(content))
	}
	if !strings.Contains(stdout.String(), "Roundfix config updated") {
		t.Fatalf("expected updated output, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunCommandHelp(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		contains []string
	}{
		{
			name:     "fetch",
			args:     []string{"fetch", "--help"},
			contains: []string{"roundfix fetch --source coderabbit --pr <number>"},
		},
		{
			name:     "resolve",
			args:     []string{"resolve", "--help"},
			contains: []string{"roundfix resolve --pr <number> --agent <agent>", "--reasoning-effort", "--no-agent-console", "--detach"},
		},
		{
			name:     "watch",
			args:     []string{"watch", "--help"},
			contains: []string{"roundfix watch --source coderabbit --pr <number> --agent <agent>", "--reasoning-effort", "--until-clean", "Review Source check succeeds", "--no-agent-console", "--detach"},
		},
		{
			name:     "setup",
			args:     []string{"setup", "--help"},
			contains: []string{"roundfix setup [--yes] [--no-input]", "--yes", "--no-input"},
		},
		{
			name:     "doctor",
			args:     []string{"doctor", "--help"},
			contains: []string{"roundfix doctor", "codex runtime hygiene", "Doctor mutates nothing"},
		},
		{
			name:     "upgrade",
			args:     []string{"upgrade", "--help"},
			contains: []string{"roundfix upgrade [--check]", "--check", "atomically"},
		},
		{
			name:     "runs",
			args:     []string{"runs", "--help"},
			contains: []string{"roundfix runs list [--all] [--state <active|terminal|all>] [--limit N]", "--all", "--state", "--limit"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run(tt.args, &stdout, &stderr)

			if code != 0 {
				t.Fatalf("expected exit code 0, got %d", code)
			}
			for _, want := range tt.contains {
				if !strings.Contains(stdout.String(), want) {
					t.Fatalf("expected help output to contain %q, got %q", want, stdout.String())
				}
			}
			if stderr.Len() != 0 {
				t.Fatalf("expected no stderr, got %q", stderr.String())
			}
		})
	}
}

// seedRunsForListColumns seeds one terminal and one Active Run in the
// current repository plus one Active Run in another repository, and pins
// the listing clock so durations are byte-stable.
func seedRunsForListColumns(t *testing.T, homeDir, repoDir, otherRepo string) []store.Run {
	t.Helper()
	base := time.Date(2026, 7, 6, 12, 0, 0, 0, time.UTC)
	runs := seedRunsForList(t, homeDir, []runListSeed{
		{
			kind:        store.KindResolve,
			state:       store.StateClean,
			gitRoot:     repoDir,
			branch:      "feature/review",
			prNumber:    "123",
			createdAt:   base,
			completedAt: base.Add(42 * time.Minute),
		},
		{
			kind:      store.KindImplement,
			state:     store.StateActive,
			gitRoot:   repoDir,
			branch:    "ma/spec-run",
			specSlug:  "0017-run-discovery",
			createdAt: base.Add(2 * time.Minute),
		},
		{
			kind:      store.KindWatch,
			state:     store.StateActive,
			gitRoot:   otherRepo,
			branch:    "feature/other",
			prNumber:  "999",
			createdAt: base.Add(4 * time.Minute),
		},
	})
	withRunsListNow(t, base.Add(14*time.Minute))
	return runs
}

func TestRunRunsListPrintsStableColumnsNewestFirst(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	otherRepo := filepath.Join(t.TempDir(), "other-repo")
	mustMkdir(t, filepath.Join(otherRepo, ".git"))
	runs := seedRunsForListColumns(t, homeDir, repoDir, otherRepo)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"runs", "list"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected exit code 0, got %d stderr=%q", code, stderr.String())
	}
	want := fmt.Sprintf(
		"%s  Active  implement  spec:0017-run-discovery  codex  2026-07-06T12:02:00Z  running 12m  ma/spec-run\n",
		runs[1].ID,
	)
	if stdout.String() != want {
		t.Fatalf("unexpected runs list output:\n got: %q\nwant: %q", stdout.String(), want)
	}
	if stderr.String() != "(1 terminal Run(s) hidden; use --state all)\n" {
		t.Fatalf("expected terminal-hidden note, got %q", stderr.String())
	}
}

func TestRunRunsListStateFlagFiltersAndNotes(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	otherRepo := filepath.Join(t.TempDir(), "other-repo")
	mustMkdir(t, filepath.Join(otherRepo, ".git"))
	runs := seedRunsForListColumns(t, homeDir, repoDir, otherRepo)
	activeLine := fmt.Sprintf(
		"%s  Active  implement  spec:0017-run-discovery  codex  2026-07-06T12:02:00Z  running 12m  ma/spec-run\n",
		runs[1].ID,
	)
	terminalLine := fmt.Sprintf(
		"%s  Clean  resolve  pr:123  codex  2026-07-06T12:00:00Z  42m  feature/review\n",
		runs[0].ID,
	)
	tests := []struct {
		name       string
		args       []string
		wantStdout string
		wantStderr string
	}{
		{
			name:       "state all includes terminal Runs without a note",
			args:       []string{"runs", "list", "--state", "all"},
			wantStdout: activeLine + terminalLine,
			wantStderr: "",
		},
		{
			name:       "state terminal excludes Active Runs and notes them",
			args:       []string{"runs", "list", "--state", "terminal"},
			wantStdout: terminalLine,
			wantStderr: "(1 active Run(s) hidden; use --state all)\n",
		},
		{
			name: "all flag appends the repository column",
			args: []string{"runs", "list", "--all", "--state", "all"},
			wantStdout: fmt.Sprintf(
				"%s  Active  watch  pr:999  codex  2026-07-06T12:04:00Z  running 10m  feature/other  %s\n",
				runs[2].ID, otherRepo,
			) + fmt.Sprintf(
				"%s  Active  implement  spec:0017-run-discovery  codex  2026-07-06T12:02:00Z  running 12m  ma/spec-run  %s\n",
				runs[1].ID, repoDir,
			) + fmt.Sprintf(
				"%s  Clean  resolve  pr:123  codex  2026-07-06T12:00:00Z  42m  feature/review  %s\n",
				runs[0].ID, repoDir,
			),
			wantStderr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run(tt.args, &stdout, &stderr)

			if code != exitOK {
				t.Fatalf("expected exit code 0, got %d stderr=%q", code, stderr.String())
			}
			if stdout.String() != tt.wantStdout {
				t.Fatalf("unexpected stdout:\n got: %q\nwant: %q", stdout.String(), tt.wantStdout)
			}
			if stderr.String() != tt.wantStderr {
				t.Fatalf("unexpected stderr:\n got: %q\nwant: %q", stderr.String(), tt.wantStderr)
			}
		})
	}
}

func TestRunRunsListLimitBoundsNewestMatches(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	base := time.Date(2026, 7, 6, 12, 0, 0, 0, time.UTC)
	seeds := make([]runListSeed, 0, 25)
	for index := 0; index < 25; index++ {
		seeds = append(seeds, runListSeed{
			kind:      store.KindResolve,
			state:     store.StateActive,
			gitRoot:   repoDir,
			branch:    fmt.Sprintf("feature/run-%02d", index),
			prNumber:  fmt.Sprintf("%d", 500+index),
			createdAt: base.Add(time.Duration(index) * time.Minute),
		})
	}
	runs := seedRunsForList(t, homeDir, seeds)
	withRunsListNow(t, base.Add(40*time.Minute))
	tests := []struct {
		name       string
		args       []string
		wantLines  int
		wantFirst  string
		wantLast   string
		wantStderr string
	}{
		{
			name:       "default bounds to the 20 newest",
			args:       []string{"runs", "list"},
			wantLines:  20,
			wantFirst:  runs[24].ID,
			wantLast:   runs[5].ID,
			wantStderr: "(5 older Run(s) hidden; use --limit 0)\n",
		},
		{
			name:       "limit zero is unbounded",
			args:       []string{"runs", "list", "--limit", "0"},
			wantLines:  25,
			wantFirst:  runs[24].ID,
			wantLast:   runs[0].ID,
			wantStderr: "",
		},
		{
			name:       "explicit limit bounds the newest",
			args:       []string{"runs", "list", "--limit", "3"},
			wantLines:  3,
			wantFirst:  runs[24].ID,
			wantLast:   runs[22].ID,
			wantStderr: "(22 older Run(s) hidden; use --limit 0)\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run(tt.args, &stdout, &stderr)

			if code != exitOK {
				t.Fatalf("expected exit code 0, got %d stderr=%q", code, stderr.String())
			}
			lines := strings.Split(strings.TrimSuffix(stdout.String(), "\n"), "\n")
			if len(lines) != tt.wantLines {
				t.Fatalf("expected %d lines, got %d:\n%s", tt.wantLines, len(lines), stdout.String())
			}
			if !strings.HasPrefix(lines[0], tt.wantFirst+"  ") {
				t.Fatalf("expected first line for Run %s, got %q", tt.wantFirst, lines[0])
			}
			if !strings.HasPrefix(lines[len(lines)-1], tt.wantLast+"  ") {
				t.Fatalf("expected last line for Run %s, got %q", tt.wantLast, lines[len(lines)-1])
			}
			if stderr.String() != tt.wantStderr {
				t.Fatalf("unexpected stderr:\n got: %q\nwant: %q", stderr.String(), tt.wantStderr)
			}
		})
	}
}

func TestRunRunsListAllRowsHiddenKeepsSingleEmptyLine(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	base := time.Date(2026, 7, 6, 12, 0, 0, 0, time.UTC)
	seedRunsForList(t, homeDir, []runListSeed{
		{
			kind:        store.KindResolve,
			state:       store.StateClean,
			gitRoot:     repoDir,
			branch:      "feature/only-terminal",
			prNumber:    "321",
			createdAt:   base,
			completedAt: base.Add(time.Minute),
		},
	})
	withRunsListNow(t, base.Add(10*time.Minute))
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"runs", "list"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected exit code 0, got %d stderr=%q", code, stderr.String())
	}
	if stdout.String() != "No Runs found.\n" {
		t.Fatalf("expected single empty-result line, got %q", stdout.String())
	}
	if stderr.String() != "(1 terminal Run(s) hidden; use --state all)\n" {
		t.Fatalf("expected terminal-hidden note, got %q", stderr.String())
	}
}

func TestRunRunsWithoutSubcommandHonorsInteractivity(t *testing.T) {
	t.Run("non-interactive exits 2 naming runs list", func(t *testing.T) {
		withCLIWorkspace(t)
		withRunsInteractiveInput(t, false)
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := Run([]string{"runs"}, &stdout, &stderr)

		if code != exitPreflight {
			t.Fatalf("expected exit code 2, got %d", code)
		}
		if stdout.Len() != 0 {
			t.Fatalf("expected no stdout, got %q", stdout.String())
		}
		if !strings.Contains(stderr.String(), "runs requires a subcommand in non-interactive mode; use 'roundfix runs list'") {
			t.Fatalf("expected non-interactive runs guidance, got %q", stderr.String())
		}
	})
	t.Run("interactive stdin without a TTY exits 2 naming runs list", func(t *testing.T) {
		withCLIWorkspace(t)
		withRunsInteractiveInput(t, true)
		t.Setenv("ROUNDFIX_TUI", "")
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := Run([]string{"runs"}, &stdout, &stderr)

		if code != exitPreflight {
			t.Fatalf("expected exit code 2, got %d", code)
		}
		if stdout.Len() != 0 {
			t.Fatalf("expected no stdout, got %q", stdout.String())
		}
		if !strings.Contains(stderr.String(), "runs requires a subcommand in non-interactive mode; use 'roundfix runs list'") {
			t.Fatalf("expected non-interactive runs guidance, got %q", stderr.String())
		}
	})
	t.Run("interactive terminal opens the Run Browser over the repository's Runs", func(t *testing.T) {
		homeDir, repoDir := withCLIWorkspace(t)
		base := time.Date(2026, 7, 6, 12, 0, 0, 0, time.UTC)
		runs := seedRunsForList(t, homeDir, []runListSeed{
			{
				kind:        store.KindResolve,
				state:       store.StateClean,
				gitRoot:     repoDir,
				branch:      "feature/terminal",
				prNumber:    "201",
				createdAt:   base,
				completedAt: base.Add(time.Minute),
			},
			{
				kind:      store.KindResolve,
				state:     store.StateActive,
				gitRoot:   repoDir,
				branch:    "feature/active",
				prNumber:  "202",
				createdAt: base.Add(2 * time.Minute),
			},
		})
		withRunsInteractiveInput(t, true)
		t.Setenv("ROUNDFIX_TUI", "always")
		sessionCalls := withRunBrowserSession(t, roundtui.BrowserOutcome{Cancelled: true})
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := Run([]string{"runs"}, &stdout, &stderr)

		if code != exitOK {
			t.Fatalf("expected exit code 0, got %d stderr=%q", code, stderr.String())
		}
		if stdout.Len() != 0 {
			t.Fatalf("expected no stdout after browser cancel, got %q", stdout.String())
		}
		calls := *sessionCalls
		if len(calls) != 1 {
			t.Fatalf("expected one Run Browser session, got %d", len(calls))
		}
		if len(calls[0].active) != 1 || calls[0].active[0].ID != runs[1].ID {
			t.Fatalf("expected the Active Run listed first, got %+v", calls[0].active)
		}
		if len(calls[0].all) != 2 || calls[0].all[0].ID != runs[1].ID || calls[0].all[1].ID != runs[0].ID {
			t.Fatalf("expected all Runs newest first, got %+v", calls[0].all)
		}
	})
	t.Run("interactive terminal without a Run Database opens the empty browser", func(t *testing.T) {
		withCLIWorkspace(t)
		withRunsInteractiveInput(t, true)
		t.Setenv("ROUNDFIX_TUI", "always")
		sessionCalls := withRunBrowserSession(t, roundtui.BrowserOutcome{Cancelled: true})
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := Run([]string{"runs"}, &stdout, &stderr)

		if code != exitOK {
			t.Fatalf("expected exit code 0, got %d stderr=%q", code, stderr.String())
		}
		calls := *sessionCalls
		if len(calls) != 1 {
			t.Fatalf("expected one Run Browser session, got %d", len(calls))
		}
		if len(calls[0].active) != 0 || len(calls[0].all) != 0 {
			t.Fatalf("expected empty listings without a Run Database, got %+v", calls[0])
		}
	})
}

func TestRunRunsListEmptyResultExitsZero(t *testing.T) {
	withCLIWorkspace(t)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"runs", "list"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected empty list exit 0, got %d stderr=%q", code, stderr.String())
	}
	if stdout.String() != "No Runs found.\n" {
		t.Fatalf("expected single empty-result line, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunRunsListUsageErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "unknown subcommand", args: []string{"runs", "bogus"}, want: "Run 'roundfix runs --help' for usage."},
		{name: "unexpected list argument", args: []string{"runs", "list", "extra"}, want: "Run 'roundfix runs list --help' for usage."},
		{name: "unknown state value", args: []string{"runs", "list", "--state", "bogus"}, want: `unknown --state "bogus"; use active, terminal, or all`},
		{name: "negative limit", args: []string{"runs", "list", "--limit", "-1"}, want: "--limit must be 0 or a positive count, got -1"},
		{name: "removed active flag", args: []string{"runs", "list", "--active"}, want: "Run 'roundfix runs list --help' for usage."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run(tt.args, &stdout, &stderr)

			if code != exitPreflight {
				t.Fatalf("expected exit code 2, got %d", code)
			}
			if stdout.Len() != 0 {
				t.Fatalf("expected no stdout, got %q", stdout.String())
			}
			if !strings.Contains(stderr.String(), tt.want) {
				t.Fatalf("expected usage pointer %q, got %q", tt.want, stderr.String())
			}
		})
	}
}

func TestRunRunsListOutsideRepositoryRequiresAll(t *testing.T) {
	homeDir := t.TempDir()
	outsideRepo := t.TempDir()
	t.Setenv("HOME", homeDir)
	t.Chdir(outsideRepo)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"runs", "list"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("expected exit code 2 outside repository, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "--all") {
		t.Fatalf("expected --all guidance, got %q", stderr.String())
	}
}

func TestRunDoctorReportsReadinessChecks(t *testing.T) {
	tests := []struct {
		name       string
		checker    *doctorFakeHealthChecker
		wantCode   int
		wantStdout string
	}{
		{
			name: "all checks pass",
			checker: newDoctorFakeHealthChecker(
				CheckResult{Name: HealthCheckNode, Status: CheckStatusOK, Detail: "v25.6.1 >= " + setupNodeMinimumVersion},
				CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK, Detail: agent.PinnedACPXVersion},
				CheckResult{Name: HealthCheckAgent, Status: CheckStatusOK, Detail: "codex"},
				CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK, Detail: "/home/roundfix/.local/bin/codex accepted"},
			),
			wantCode: exitOK,
			wantStdout: "node: ok (v25.6.1 >= " + setupNodeMinimumVersion + ")\n" +
				"acpx: ok (" + agent.PinnedACPXVersion + ")\n" +
				"agent: ok (codex)\n" +
				"codex: ok (/home/roundfix/.local/bin/codex accepted)\n",
		},
		{
			name: "quarantined codex fails with reinstall action",
			checker: newDoctorFakeHealthChecker(
				CheckResult{Name: HealthCheckNode, Status: CheckStatusOK, Detail: "v25.6.1 >= " + setupNodeMinimumVersion},
				CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK, Detail: agent.PinnedACPXVersion},
				CheckResult{Name: HealthCheckAgent, Status: CheckStatusOK, Detail: "codex"},
				CheckResult{Name: HealthCheckCodex, Status: CheckStatusFailed, Detail: "/tmp/codex is quarantined", NextAction: codex.ReinstallNextAction},
			),
			wantCode: exitRunFailed,
			wantStdout: "node: ok (v25.6.1 >= " + setupNodeMinimumVersion + ")\n" +
				"acpx: ok (" + agent.PinnedACPXVersion + ")\n" +
				"agent: ok (codex)\n" +
				"codex: failed (/tmp/codex is quarantined; next: " + codex.ReinstallNextAction + ")\n",
		},
		{
			name: "codex not applicable does not fail",
			checker: newDoctorFakeHealthChecker(
				CheckResult{Name: HealthCheckNode, Status: CheckStatusOK, Detail: "v25.6.1 >= " + setupNodeMinimumVersion},
				CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK, Detail: agent.PinnedACPXVersion},
				CheckResult{Name: HealthCheckAgent, Status: CheckStatusOK, Detail: "codex"},
				CheckResult{Name: HealthCheckCodex, Status: CheckStatusSkipped, Detail: "not-applicable on linux"},
			),
			wantCode: exitOK,
			wantStdout: "node: ok (v25.6.1 >= " + setupNodeMinimumVersion + ")\n" +
				"acpx: ok (" + agent.PinnedACPXVersion + ")\n" +
				"agent: ok (codex)\n" +
				"codex: skipped (not-applicable on linux)\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, repoDir := withCLIWorkspace(t)
			withDoctorFakeDeps(t, tt.checker)
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run([]string{"doctor"}, &stdout, &stderr)

			if code != tt.wantCode {
				t.Fatalf("expected exit code %d, got %d", tt.wantCode, code)
			}
			if got := stdout.String(); got != tt.wantStdout {
				t.Fatalf("unexpected stdout:\n got: %q\nwant: %q", got, tt.wantStdout)
			}
			if stderr.Len() != 0 {
				t.Fatalf("expected no stderr, got %q", stderr.String())
			}
			if len(tt.checker.agentRequests) != 1 {
				t.Fatalf("expected one configured codex probe, got %#v", tt.checker.agentRequests)
			}
			gotProbe := tt.checker.agentRequests[0]
			if gotProbe.WorkDir != "/repo/project" ||
				gotProbe.Runtime.ID != "codex" ||
				gotProbe.Runtime.Model != "gpt-5.5" ||
				gotProbe.Runtime.ReasoningEffort != "xhigh" {
				t.Fatalf("expected configured codex selection probe in repo workdir, got %#v", gotProbe)
			}
			assertDoctorPathMissing(t, filepath.Join(homeDir, ".acpx"))
			assertDoctorPathMissing(t, filepath.Join(homeDir, ".roundfix"))
			assertDoctorPathMissing(t, filepath.Join(repoDir, ".roundfixrc.yml"))
		})
	}
}

func TestRunDoctorRejectsMissingConfiguredAgentModel(t *testing.T) {
	config := roundconfig.Builtin()
	config.Runtimes.Codex.Model = ""
	checker := newDoctorFakeHealthChecker(
		CheckResult{Name: HealthCheckNode, Status: CheckStatusOK, Detail: "v25.6.1 >= " + setupNodeMinimumVersion},
		CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK, Detail: agent.PinnedACPXVersion},
		CheckResult{Name: HealthCheckAgent, Status: CheckStatusOK, Detail: "should not run"},
		CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK, Detail: "/home/roundfix/.local/bin/codex accepted"},
	)
	withDoctorFakeLoaded(t, checker, roundconfig.Loaded{
		Config:  config,
		GitRoot: "/repo/project",
		HomeDir: "/home/roundfix-test",
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"doctor"}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("expected doctor selection failure exit %d, got %d", exitRunFailed, code)
	}
	if !strings.Contains(stdout.String(), `agent: failed (agent selection for runtime "codex" missing model`) {
		t.Fatalf("expected missing selection diagnostic, got %q", stdout.String())
	}
	if len(checker.agentRequests) != 0 {
		t.Fatalf("expected invalid selection to skip readiness probe, got %#v", checker.agentRequests)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunDoctorRejectsArguments(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"doctor", "extra"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("expected exit code %d, got %d", exitPreflight, code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "unexpected argument \"extra\"") {
		t.Fatalf("expected argument diagnostic, got %q", stderr.String())
	}
}

func newDoctorFakeHealthChecker(node, acpx, agentResult, codex CheckResult) *doctorFakeHealthChecker {
	return &doctorFakeHealthChecker{
		node:        node,
		acpx:        acpx,
		agentResult: agentResult,
		codex:       codex,
	}
}

type doctorFakeHealthChecker struct {
	node          CheckResult
	acpx          CheckResult
	agentResult   CheckResult
	codex         CheckResult
	agentRequests []agent.ProbeRequest
}

func (checker *doctorFakeHealthChecker) Node(context.Context) CheckResult {
	return checker.node
}

func (checker *doctorFakeHealthChecker) ACPX(context.Context) CheckResult {
	return checker.acpx
}

func (checker *doctorFakeHealthChecker) Agent(_ context.Context, req agent.ProbeRequest) CheckResult {
	checker.agentRequests = append(checker.agentRequests, req)
	return checker.agentResult
}

func (checker *doctorFakeHealthChecker) Codex(context.Context) CheckResult {
	return checker.codex
}

func withDoctorFakeDeps(t *testing.T, checker HealthChecker) {
	t.Helper()
	withDoctorFakeLoaded(t, checker, roundconfig.Loaded{
		Config:  roundconfig.Builtin(),
		GitRoot: "/repo/project",
		HomeDir: "/home/roundfix-test",
	})
}

func withDoctorFakeLoaded(t *testing.T, checker HealthChecker, loaded roundconfig.Loaded) {
	t.Helper()
	old := doctorDeps
	doctorDeps = doctorDependencies{
		loadConfig: func(roundconfig.LoadOptions) (roundconfig.Loaded, error) {
			return loaded, nil
		},
		healthChecker: func(roundconfig.Loaded) HealthChecker {
			return checker
		},
	}
	t.Cleanup(func() {
		doctorDeps = old
	})
}

func assertDoctorPathMissing(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected %s to be absent, got error %v", path, err)
	}
}

func TestRunDetachRejectsInteractiveWithExistingConflictShape(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "resolve", args: []string{"resolve", "--detach", "--interactive"}},
		{name: "watch", args: []string{"watch", "--detach", "--interactive"}},
		{name: "implement", args: []string{"implement", "--detach", "--interactive"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, _ := withCLIWorkspace(t)
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := RunContext(context.Background(), tt.args, &stdout, &stderr)

			if code != exitPreflight {
				t.Fatalf("expected exit code 2, got %d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
			}
			if stdout.Len() != 0 {
				t.Fatalf("expected no stdout, got %q", stdout.String())
			}
			if !strings.Contains(stderr.String(), "--interactive cannot be used with --no-input") {
				t.Fatalf("expected existing no-input conflict shape, got %q", stderr.String())
			}
			assertNoRunDatabase(t, homeDir)
		})
	}
}

func TestRunSetupFreshMachineAcceptsOffers(t *testing.T) {
	fake := newSetupFakeDeps()
	fake.acpxErr = errors.New("acpx not found")
	fake.paths["codex-acp"] = "/bin/codex-acp"
	withSetupFakeDeps(t, fake)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"setup", "--yes"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected setup exit 0, got %d (stderr %q)", code, stderr.String())
	}
	assertSetupLineOrder(t, stdout.String(), []string{
		"node: ok",
		"acpx: installed",
		"agent probe: ok",
		"acpx agents override: installed",
		"User Config: installed",
		"Project Config: installed",
	})
	if got := strings.Join(fake.installCalls, "\n"); got != "npm install -g acpx@"+agent.PinnedACPXVersion {
		t.Fatalf("expected pinned acpx install command, got %q", got)
	}
	if got := strings.Join(fake.initScopes, ","); got != "user,project" {
		t.Fatalf("expected User and Project Config init flows, got %q", got)
	}
	if fake.acpxInitCalls != 1 {
		t.Fatalf("expected one acpx config init call, got %d", fake.acpxInitCalls)
	}
	acpxConfig := fake.files[fake.acpxConfigPath]
	if !strings.Contains(acpxConfig, `"codex"`) || !strings.Contains(acpxConfig, `"command": "codex-acp"`) {
		t.Fatalf("expected codex direct adapter override, got %s", acpxConfig)
	}
	for _, want := range []string{"acpx agents override diff:", "--- ", "+++ ", "+  \"agents\""} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("expected setup output to contain %q, got %q", want, stdout.String())
		}
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected --yes to avoid prompts, got stderr %q", stderr.String())
	}
}

func TestRunSetupHealthyMachineIsIdempotent(t *testing.T) {
	fake := newSetupFakeDeps()
	fake.paths["codex-acp"] = "/bin/codex-acp"
	fake.files[fake.userConfigPath] = "defaults:\n  agent: codex\n"
	fake.files[fake.projectConfigPath] = "defaults:\n  verification: make verify\n"
	fake.files[fake.acpxConfigPath] = "{\n  \"agents\": {\n    \"codex\": {\n      \"command\": \"codex-acp\"\n    }\n  }\n}\n"
	withSetupFakeDeps(t, fake)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"setup"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected setup exit 0, got %d (stderr %q)", code, stderr.String())
	}
	assertSetupLineOrder(t, stdout.String(), []string{
		"node: ok",
		"acpx: ok",
		"agent probe: ok",
		"acpx agents override: ok",
		"User Config: ok",
		"Project Config: ok",
	})
	if len(fake.installCalls) != 0 || len(fake.initScopes) != 0 || len(fake.writeCalls) != 0 || len(fake.prompts) != 0 {
		t.Fatalf("expected idempotent setup to avoid side effects, installs=%v init=%v writes=%v prompts=%v", fake.installCalls, fake.initScopes, fake.writeCalls, fake.prompts)
	}
	if len(fake.probeRequests) != 1 {
		t.Fatalf("expected one Agent readiness probe, got %#v", fake.probeRequests)
	}
	gotProbe := fake.probeRequests[0]
	if gotProbe.WorkDir != fake.gitRoot ||
		gotProbe.Runtime.ID != "codex" ||
		gotProbe.Runtime.Model != "gpt-5.5" ||
		gotProbe.Runtime.ReasoningEffort != "xhigh" {
		t.Fatalf("expected setup to probe the effective Codex selection in the repo workdir, got %#v", gotProbe)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunSetupAgentProbeUsesConfiguredSelectionAndWorkDir(t *testing.T) {
	fake := newSetupFakeDeps()
	fake.config.Runtimes.Codex.Model = "repo-codex"
	fake.config.Runtimes.Codex.ReasoningEffort = "repo-xhigh"
	withSetupFakeDeps(t, fake)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"setup", "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected setup exit 0, got %d (stdout %q stderr %q)", code, stdout.String(), stderr.String())
	}
	if len(fake.probeRequests) != 1 {
		t.Fatalf("expected one Agent readiness probe, got %#v", fake.probeRequests)
	}
	gotProbe := fake.probeRequests[0]
	if gotProbe.WorkDir != fake.gitRoot ||
		gotProbe.Runtime.ID != "codex" ||
		gotProbe.Runtime.Model != "repo-codex" ||
		gotProbe.Runtime.ReasoningEffort != "repo-xhigh" {
		t.Fatalf("expected setup probe to use configured Agent selection, got %#v", gotProbe)
	}
}

func TestRunSetupRejectsMissingConfiguredAgentSelection(t *testing.T) {
	fake := newSetupFakeDeps()
	fake.config.Runtimes.Codex.Model = ""
	withSetupFakeDeps(t, fake)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"setup", "--no-input"}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("expected setup selection failure exit %d, got %d", exitRunFailed, code)
	}
	if !strings.Contains(stdout.String(), `agent probe: failed (agent selection for runtime "codex" missing model`) {
		t.Fatalf("expected missing selection diagnostic, got %q", stdout.String())
	}
	if len(fake.probeRequests) != 0 {
		t.Fatalf("expected invalid selection to skip readiness probe, got %#v", fake.probeRequests)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunSetupMismatchedACPXUpgradeOffer(t *testing.T) {
	fake := newSetupFakeDeps()
	fake.acpxVersion = "0.11.0"
	fake.files[fake.userConfigPath] = "defaults:\n  agent: codex\n"
	fake.files[fake.projectConfigPath] = "defaults:\n  verification: make verify\n"
	withSetupFakeDeps(t, fake)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"setup", "--yes"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected setup exit 0, got %d (stderr %q)", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "acpx: installed") || !strings.Contains(stdout.String(), "found 0.11.0") {
		t.Fatalf("expected mismatched acpx upgrade report, got %q", stdout.String())
	}
	if got := strings.Join(fake.installCalls, "\n"); got != "npm install -g acpx@"+agent.PinnedACPXVersion {
		t.Fatalf("expected pinned acpx install command, got %q", got)
	}
}

func TestRunSetupMergesACPXAgentsOverridePreservingUnrelatedBytes(t *testing.T) {
	fake := newSetupFakeDeps()
	fake.paths["codex-acp"] = "/bin/codex-acp"
	fake.files[fake.userConfigPath] = "defaults:\n  agent: codex\n"
	fake.files[fake.projectConfigPath] = "defaults:\n  verification: make verify\n"
	unrelated := "  \"theme\": {\n    \"color\": \"blue\",\n    \"nested\": [1, 2, 3]\n  }"
	existingAgent := "    \"claude\": {\n      \"command\": \"existing-claude\"\n    }"
	fake.files[fake.acpxConfigPath] = "{\n" + unrelated + ",\n  \"agents\": {\n" + existingAgent + "\n  }\n}\n"
	withSetupFakeDeps(t, fake)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"setup", "--yes"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected setup exit 0, got %d (stderr %q)", code, stderr.String())
	}
	acpxConfig := fake.files[fake.acpxConfigPath]
	for _, want := range []string{unrelated, existingAgent, `"codex"`, `"command": "codex-acp"`} {
		if !strings.Contains(acpxConfig, want) {
			t.Fatalf("expected merged config to preserve/include %q, got %s", want, acpxConfig)
		}
	}
	if !strings.Contains(stdout.String(), "acpx agents override diff:") || !strings.Contains(stdout.String(), "+    \"codex\"") {
		t.Fatalf("expected before/after diff for override, got %q", stdout.String())
	}
}

func TestRunSetupDeclinedOffersReportNoWrites(t *testing.T) {
	fake := newSetupFakeDeps()
	fake.acpxErr = errors.New("acpx not found")
	fake.paths["codex-acp"] = "/bin/codex-acp"
	fake.confirm = func(context.Context, io.Writer, string) (bool, error) {
		return false, nil
	}
	withSetupFakeDeps(t, fake)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"setup"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected declined setup offers to exit 0, got %d (stderr %q)", code, stderr.String())
	}
	assertSetupLineOrder(t, stdout.String(), []string{
		"node: ok",
		"acpx: offered: declined",
		"agent probe: skipped",
		"acpx agents override: offered: declined",
		"User Config: offered: declined",
		"Project Config: offered: declined",
	})
	if len(fake.installCalls) != 0 || len(fake.initScopes) != 0 || len(fake.writeCalls) != 0 || fake.acpxInitCalls != 0 {
		t.Fatalf("expected no writes after declined offers, installs=%v init=%v writes=%v acpxInit=%d", fake.installCalls, fake.initScopes, fake.writeCalls, fake.acpxInitCalls)
	}
	if len(fake.prompts) != 4 {
		t.Fatalf("expected four confirmation prompts, got %d: %v", len(fake.prompts), fake.prompts)
	}
}

func TestRunSetupNoInputSkipsOffers(t *testing.T) {
	fake := newSetupFakeDeps()
	fake.acpxErr = errors.New("acpx not found")
	fake.paths["codex-acp"] = "/bin/codex-acp"
	withSetupFakeDeps(t, fake)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"setup", "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected no-input skipped offers to exit 0, got %d (stderr %q)", code, stderr.String())
	}
	assertSetupLineOrder(t, stdout.String(), []string{
		"node: ok",
		"acpx: skipped",
		"agent probe: skipped",
		"acpx agents override: skipped",
		"User Config: skipped",
		"Project Config: skipped",
	})
	if len(fake.installCalls) != 0 || len(fake.initScopes) != 0 || len(fake.writeCalls) != 0 || len(fake.prompts) != 0 {
		t.Fatalf("expected --no-input to avoid side effects, installs=%v init=%v writes=%v prompts=%v", fake.installCalls, fake.initScopes, fake.writeCalls, fake.prompts)
	}
}

func TestRunSetupExitCodes(t *testing.T) {
	t.Run("check failure exits one", func(t *testing.T) {
		fake := newSetupFakeDeps()
		fake.nodeVersion = "v20.0.0"
		fake.files[fake.userConfigPath] = "defaults:\n  agent: codex\n"
		fake.files[fake.projectConfigPath] = "defaults:\n  verification: make verify\n"
		withSetupFakeDeps(t, fake)
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := Run([]string{"setup"}, &stdout, &stderr)

		if code != exitRunFailed {
			t.Fatalf("expected setup check failure exit 1, got %d", code)
		}
		if !strings.Contains(stdout.String(), "node: failed") {
			t.Fatalf("expected node failure report, got %q", stdout.String())
		}
	})

	t.Run("usage error exits two", func(t *testing.T) {
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := Run([]string{"setup", "--yes", "--no-input"}, &stdout, &stderr)

		if code != exitPreflight {
			t.Fatalf("expected setup usage error exit 2, got %d", code)
		}
		if stdout.Len() != 0 {
			t.Fatalf("expected no stdout, got %q", stdout.String())
		}
		if !strings.Contains(stderr.String(), "--yes cannot be used with --no-input") {
			t.Fatalf("expected usage diagnostic, got %q", stderr.String())
		}
	})
}

func TestHealthCheckerReturnsStructuredReadOnlyResults(t *testing.T) {
	ctx := context.Background()
	req := agent.ProbeRequest{
		Runtime: agent.RuntimeSpec{ID: "codex", Model: "gpt-5.5", ReasoningEffort: "xhigh"},
		WorkDir: "/repo/project",
	}
	probed := false
	checker := newHealthChecker(healthCheckDependencies{
		nodeVersion: func(context.Context) (string, error) {
			return " v25.6.1 \n", nil
		},
		acpxVersion: func(context.Context) (string, error) {
			return agent.PinnedACPXVersion + "\n", nil
		},
		probeAgent: func(_ context.Context, got agent.ProbeRequest) error {
			probed = true
			if !reflect.DeepEqual(got, req) {
				t.Fatalf("expected probe request %#v, got %#v", req, got)
			}
			return nil
		},
	})

	assertCheckResult(t, checker.Node(ctx), CheckResult{
		Name:   HealthCheckNode,
		Status: CheckStatusOK,
		Detail: "v25.6.1 >= " + setupNodeMinimumVersion,
	})
	assertCheckResult(t, checker.ACPX(ctx), CheckResult{
		Name:   HealthCheckACPX,
		Status: CheckStatusOK,
		Detail: agent.PinnedACPXVersion,
	})
	assertCheckResult(t, checker.Agent(ctx, req), CheckResult{
		Name:   HealthCheckAgent,
		Status: CheckStatusOK,
		Detail: "codex",
	})
	if !probed {
		t.Fatal("expected agent probe to run")
	}
}

func TestHealthCheckerReportsFailedPrerequisitesWithNextActions(t *testing.T) {
	ctx := context.Background()
	checker := newHealthChecker(healthCheckDependencies{
		nodeVersion: func(context.Context) (string, error) {
			return "v20.0.0", nil
		},
		acpxVersion: func(context.Context) (string, error) {
			return "0.11.0", nil
		},
		probeAgent: func(context.Context, agent.ProbeRequest) error {
			return errors.New("probe denied")
		},
	})

	assertCheckResult(t, checker.Node(ctx), CheckResult{
		Name:       HealthCheckNode,
		Status:     CheckStatusFailed,
		Detail:     "found v20.0.0; Node.js " + setupNodeMinimumVersion + " or newer is required",
		NextAction: "install Node.js " + setupNodeMinimumVersion + " or newer",
	})
	assertCheckResult(t, checker.ACPX(ctx), CheckResult{
		Name:       HealthCheckACPX,
		Status:     CheckStatusFailed,
		Detail:     "found 0.11.0; required " + agent.PinnedACPXVersion + "; run " + setupACPXInstallCommand(),
		NextAction: setupACPXInstallCommand(),
	})
	assertCheckResult(t, checker.Agent(ctx, agent.ProbeRequest{Runtime: agent.RuntimeSpec{ID: "codex"}}), CheckResult{
		Name:   HealthCheckAgent,
		Status: CheckStatusFailed,
		Detail: "probe denied",
	})
}

func TestHealthCheckerReportsCodexResult(t *testing.T) {
	checker := newHealthChecker(healthCheckDependencies{
		codexInspector: fakeCodexInspector{
			result: codex.Result{
				Status:     codex.StatusFailed,
				Detail:     "/path/codex has an invalid or missing code signature",
				NextAction: codex.ReinstallNextAction,
			},
		},
	})

	assertCheckResult(t, checker.Codex(context.Background()), CheckResult{
		Name:       HealthCheckCodex,
		Status:     CheckStatusFailed,
		Detail:     "/path/codex has an invalid or missing code signature",
		NextAction: codex.ReinstallNextAction,
	})
}

func assertCheckResult(t *testing.T, got CheckResult, want CheckResult) {
	t.Helper()
	if got != want {
		t.Fatalf("unexpected check result:\n got: %#v\nwant: %#v", got, want)
	}
}

type fakeCodexInspector struct {
	result codex.Result
}

func (fake fakeCodexInspector) Inspect(context.Context) codex.Result {
	return fake.result
}

func TestRunSkillsCheck(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"skills", "check"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !strings.Contains(stdout.String(), "Roundfix skill check passed") {
		t.Fatalf("expected skills check output, got %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "roundfix") || strings.Contains(stdout.String(), "roundfix-watch") || strings.Contains(stdout.String(), "roundfix-resolve-round") {
		t.Fatalf("expected consolidated skill name only, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunSkillsListSeparatesOwnedFromExternal(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"skills", "list"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	out := stdout.String()
	for _, want := range []string{
		"Bundled skills",
		"roundfix",
		"write-prd",
		"evidence-gate",
		"Recommended skills",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected skills list to contain %q, got %q", want, out)
		}
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunSkillsInstallCopiesArtifactsToProjectByDefault(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"skills", "install"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !strings.Contains(stdout.String(), "Installed Roundfix skill for project") {
		t.Fatalf("expected project install output, got %q", stdout.String())
	}
	targetDir := filepath.Join(repoDir, ".agents", "skills")
	for _, path := range []string{
		"roundfix/SKILL.md",
		"roundfix/agents/openai.yaml",
	} {
		if _, err := os.Stat(filepath.Join(targetDir, path)); err != nil {
			t.Fatalf("expected installed file %s: %v", path, err)
		}
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunSkillsInstallCopiesArtifactsToExplicitTarget(t *testing.T) {
	targetDir := t.TempDir()
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"skills", "install", "--target", "codex", "--dir", targetDir}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !strings.Contains(stdout.String(), "Installed Roundfix skill for codex") {
		t.Fatalf("expected install output, got %q", stdout.String())
	}
	for _, path := range []string{
		"roundfix/SKILL.md",
		"roundfix/agents/openai.yaml",
	} {
		if _, err := os.Stat(filepath.Join(targetDir, path)); err != nil {
			t.Fatalf("expected installed file %s: %v", path, err)
		}
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunSkillsInstallCreatesClaudeProjectSymlink(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	mustMkdir(t, filepath.Join(repoDir, ".claude", "skills"))
	prompted := false
	withClaudeSkillSymlinkPrompt(t, func(ctx context.Context, stderr io.Writer, linkPath string, target string) (bool, error) {
		prompted = true
		expectedLink := filepath.Join(repoDir, ".claude", "skills", "roundfix")
		if linkPath != expectedLink {
			t.Fatalf("expected prompt link path %q, got %q", expectedLink, linkPath)
		}
		expectedTarget := filepath.Join("..", "..", ".agents", "skills", "roundfix")
		if target != expectedTarget {
			t.Fatalf("expected prompt target %q, got %q", expectedTarget, target)
		}
		return true, nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"skills", "install"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !prompted {
		t.Fatal("expected Claude symlink prompt")
	}
	linkPath := filepath.Join(repoDir, ".claude", "skills", "roundfix")
	target, err := os.Readlink(linkPath)
	if err != nil {
		t.Fatalf("expected Claude skill symlink: %v", err)
	}
	expectedTarget := filepath.Join("..", "..", ".agents", "skills", "roundfix")
	if target != expectedTarget {
		t.Fatalf("expected symlink target %q, got %q", expectedTarget, target)
	}
	if !strings.Contains(stdout.String(), "Created Claude skill symlink") {
		t.Fatalf("expected symlink output, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunSkillsInstallRejectsUnsupportedTarget(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"skills", "install", "--target", "other"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "unsupported skill install target") {
		t.Fatalf("expected unsupported target error, got %q", stderr.String())
	}
}

func TestRunOperationalCommandAcceptsMVPFlags(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedCode   int
		expectedOutput string
	}{
		{
			name:           "fetch",
			args:           []string{"fetch", "--source", "coderabbit", "--pr", "123", "--round", "auto", "--no-input"},
			expectedCode:   0,
			expectedOutput: "reached Fetched",
		},
		{
			name:           "resolve",
			args:           []string{"resolve", "--pr", "123", "--agent", "codex", "--round", "all", "--no-input"},
			expectedCode:   0,
			expectedOutput: "resolve selected 1 downloaded Unresolved Review Issue",
		},
		{
			name:           "watch",
			args:           []string{"watch", "--source", "coderabbit", "--pr", "123", "--agent", "codex", "--until-clean", "--max-rounds", "6", "--no-input"},
			expectedCode:   0,
			expectedOutput: "Watch Run",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, repoDir := withCLIWorkspace(t)
			withSuccessfulPreflight(t, repoDir)
			if tt.name == "resolve" {
				persistCLIReviewIssue(t, repoDir, 1, "feature/review")
			}
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run(tt.args, &stdout, &stderr)

			if code != tt.expectedCode {
				t.Fatalf("expected exit code %d, got %d", tt.expectedCode, code)
			}
			output := stdout.String() + stderr.String()
			if !strings.Contains(output, tt.expectedOutput) {
				t.Fatalf("expected output to contain %q, got stdout=%q stderr=%q", tt.expectedOutput, stdout.String(), stderr.String())
			}
			if !strings.Contains(output, builtinArtifactDirForRepo(t, repoDir)) {
				t.Fatalf("expected artifact dir in output, got stdout=%q stderr=%q", stdout.String(), stderr.String())
			}
			if strings.Contains(output, "not implemented yet") {
				t.Fatalf("did not expect scaffold message, got stdout=%q stderr=%q", stdout.String(), stderr.String())
			}
			if tt.name == "fetch" {
				if !strings.Contains(stdout.String(), "did not start an Agent, commit, or push") {
					t.Fatalf("expected fetch safety confirmation, got %q", stdout.String())
				}
				if !strings.Contains(stdout.String(), "Review Issues: 1") {
					t.Fatalf("expected fetched issue count, got %q", stdout.String())
				}
				if !strings.Contains(stdout.String(), "Artifacts: created new Round") {
					t.Fatalf("expected new artifact confirmation, got %q", stdout.String())
				}
				issuePath := filepath.Join(defaultReviewRootForRepo(repoDir, "123"), "round-001", "issue_001.md")
				issueContent, err := os.ReadFile(issuePath)
				if err != nil {
					t.Fatalf("expected Review Issue artifact %s: %v", issuePath, err)
				}
				if !strings.Contains(string(issueContent), "source_ref: thread:PRRT_test,comment:PRRC_test") {
					t.Fatalf("expected source ref in issue artifact, got %s", string(issueContent))
				}
				assertRunCount(t, filepath.Join(os.Getenv("HOME"), ".roundfix", "roundfix.db"), 1)
			} else if tt.name == "resolve" {
				if !strings.Contains(stderr.String(), "fake agent output") {
					t.Fatalf("expected fake Agent output, got %q", stderr.String())
				}
				if !strings.Contains(stderr.String(), "fake verification output") {
					t.Fatalf("expected fake verification output, got %q", stderr.String())
				}
				if !strings.Contains(stderr.String(), "Batch commit created") {
					t.Fatalf("expected Batch commit confirmation, got %q", stderr.String())
				}
				if !strings.Contains(stderr.String(), "Resolved 1 Review Source thread") {
					t.Fatalf("expected Review Source resolution confirmation, got %q", stderr.String())
				}
				if !strings.Contains(stderr.String(), "Final Push completed") {
					t.Fatalf("expected Final Push confirmation, got %q", stderr.String())
				}
				assertRunCount(t, filepath.Join(os.Getenv("HOME"), ".roundfix", "roundfix.db"), 1)
				assertNoAgentLogs(t, repoDir)
			} else if tt.name == "watch" {
				if !strings.Contains(stderr.String(), "Review Source status: settled") {
					t.Fatalf("expected fake Review Source status output, got %q", stderr.String())
				}
				if !strings.Contains(stderr.String(), "Fetched Round 001 with 1 Review Issue") {
					t.Fatalf("expected watch fetch output, got %q", stderr.String())
				}
				if !strings.Contains(stderr.String(), "fake agent output") {
					t.Fatalf("expected fake Agent output, got %q", stderr.String())
				}
				if !strings.Contains(stderr.String(), "fake verification output") {
					t.Fatalf("expected fake verification output, got %q", stderr.String())
				}
				if !strings.Contains(stderr.String(), "Batch commit created") {
					t.Fatalf("expected Batch commit confirmation, got %q", stderr.String())
				}
				if !strings.Contains(stderr.String(), "Resolved 1 Review Source thread") {
					t.Fatalf("expected Review Source resolution confirmation, got %q", stderr.String())
				}
				if !strings.Contains(stderr.String(), "Final Push completed") {
					t.Fatalf("expected Final Push confirmation, got %q", stderr.String())
				}
				if !strings.Contains(stderr.String(), "reached Clean") {
					t.Fatalf("expected watch Clean terminal outcome, got %q", stderr.String())
				}
				assertRunCount(t, filepath.Join(os.Getenv("HOME"), ".roundfix", "roundfix.db"), 1)
				assertNoAgentLogs(t, repoDir)
			}
		})
	}
}

func TestRunFetchWritesReviewArtifactsUnderSpecSelector(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	specSlug := "0001-widget-flow"
	mustMkdir(t, filepath.Join(repoDir, "docs", "specs", specSlug))
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"fetch", "--pr", "123", "--spec", specSlug, "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected successful fetch exit code 0, got %d; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	issuePath := filepath.Join(repoDir, "docs", "specs", specSlug, "reviews", "round-001", "issue_001.md")
	if _, err := os.Stat(issuePath); err != nil {
		t.Fatalf("expected spec-associated Review Issue artifact %s: %v", issuePath, err)
	}
	oldPath := filepath.Join(builtinArtifactDirForRepo(t, repoDir), "reviews", "pr-123", "round-001", "issue_001.md")
	if _, err := os.Stat(oldPath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected no default Artifact Directory Review Issue at %s, got err %v", oldPath, err)
	}
}

func TestRunFetchWritesReviewArtifactsUnderTrailerSpec(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	specSlug := "0001-widget-flow"
	mustMkdir(t, filepath.Join(repoDir, "docs", "specs", specSlug))
	withReviewSpecGitRunner(t, &fakeReviewSpecGitRunner{message: "subject\n\nRoundfix-Spec: " + specSlug})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"fetch", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected successful fetch exit code 0, got %d; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	issuePath := filepath.Join(repoDir, "docs", "specs", specSlug, "reviews", "round-001", "issue_001.md")
	if _, err := os.Stat(issuePath); err != nil {
		t.Fatalf("expected trailer-associated Review Issue artifact %s: %v", issuePath, err)
	}
}

func TestRunFetchWritesReviewArtifactsUnderSpeclessRoot(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"fetch", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected successful fetch exit code 0, got %d; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	issuePath := filepath.Join(repoDir, "docs", "specs", "_reviews", "pr-123", "round-001", "issue_001.md")
	if _, err := os.Stat(issuePath); err != nil {
		t.Fatalf("expected spec-less Review Issue artifact %s: %v", issuePath, err)
	}
}

func TestReviewSpecSlugAssociation(t *testing.T) {
	repoRoot := t.TempDir()
	explicitSlug := "0001-explicit"
	oldSlug := "0002-old"
	newSlug := "0003-new"
	for _, slug := range []string{explicitSlug, oldSlug, newSlug} {
		mustMkdir(t, filepath.Join(repoRoot, "docs", "specs", slug))
	}

	tests := []struct {
		name     string
		req      reviewSpecRequest
		message  string
		want     string
		wantCall bool
	}{
		{
			name: "explicit selector wins without reading git",
			req: reviewSpecRequest{
				ExplicitSlug: explicitSlug,
				RepoRoot:     repoRoot,
				HeadSHA:      "abc123",
			},
			message:  "subject\n\nRoundfix-Spec: " + newSlug,
			want:     explicitSlug,
			wantCall: false,
		},
		{
			name: "newest trailer wins",
			req: reviewSpecRequest{
				RepoRoot: repoRoot,
				HeadSHA:  "abc123",
			},
			message:  "subject\n\nRoundfix-Spec: " + oldSlug + "\nRoundfix-Spec: " + newSlug,
			want:     newSlug,
			wantCall: true,
		},
		{
			name: "unknown trailer slug is no association",
			req: reviewSpecRequest{
				RepoRoot: repoRoot,
				HeadSHA:  "abc123",
			},
			message:  "subject\n\nRoundfix-Spec: 9999-missing",
			want:     "",
			wantCall: true,
		},
		{
			name: "missing trailer is no association",
			req: reviewSpecRequest{
				RepoRoot: repoRoot,
				HeadSHA:  "abc123",
			},
			message:  "subject only",
			want:     "",
			wantCall: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := &fakeReviewSpecGitRunner{message: tt.message}

			got, err := reviewSpecSlug(context.Background(), tt.req, runner)
			if err != nil {
				t.Fatalf("reviewSpecSlug() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("reviewSpecSlug() = %q, want %q", got, tt.want)
			}
			if gotCall := runner.calls > 0; gotCall != tt.wantCall {
				t.Fatalf("git calls = %v, want %v", gotCall, tt.wantCall)
			}
		})
	}
}

func TestReviewArtifactUsesDefaultSpecsRootMatchesConfigString(t *testing.T) {
	if !reviewArtifactUsesDefaultSpecsRoot("docs/specs") {
		t.Fatal("expected built-in default Spec Root to match")
	}
	if !reviewArtifactUsesDefaultSpecsRoot("  docs/specs  ") {
		t.Fatal("expected whitespace-trimmed built-in default Spec Root to match")
	}
	if reviewArtifactUsesDefaultSpecsRoot(`docs\specs`) {
		t.Fatal("expected OS-specific path spelling to be treated as non-default config")
	}
}

type fakeReviewSpecGitRunner struct {
	message string
	calls   int
}

func (runner *fakeReviewSpecGitRunner) RunGit(_ context.Context, _ string, _ ...string) (string, error) {
	runner.calls++
	return runner.message, nil
}

func TestRunFetchWarnsAndIgnoresDeprecatedUserConfig(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
	mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), `
resolve:
  concurrent: 1
`)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"fetch", "--source", "coderabbit", "--pr", "123", "--round", "auto", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected fetch exit 0, got %d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if strings.Contains(stdout.String(), "deprecated") {
		t.Fatalf("expected warning to stay off stdout, got %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "Review Issues: 1") || !strings.Contains(stdout.String(), "Artifacts: created new Round") {
		t.Fatalf("expected normal fetch stdout, got %q", stdout.String())
	}
	warning := "config: resolve.concurrent is deprecated and ignored; use worktree.concurrency\n"
	if strings.Count(stderr.String(), warning) != 1 {
		t.Fatalf("expected one deprecated config warning %q, got %q", warning, stderr.String())
	}
	if strings.Contains(stderr.String(), "Preflight failed") {
		t.Fatalf("expected fetch to proceed past Preflight, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
}

func TestRunWatchPrintsDeterministicStdoutReport(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T, repoDir string)
		args       []string
		wantCode   int
		wantStdout string
	}{
		{
			name: "clean",
			args: []string{"watch", "--source", "coderabbit", "--pr", "123", "--agent", "codex", "--until-clean", "--max-rounds", "6", "--no-input"},
			wantStdout: "" +
				"issue 001 resolved — major: handle test issue\n" +
				"Clean after 1 Round(s): 1 resolved, 0 invalid, 0 failed, 0 unresolved.\n",
		},
		{
			name: "unresolved",
			setup: func(t *testing.T, _ string) {
				t.Helper()
				withAgentRunner(t, &fakeAgentRunner{status: rounds.StatusFailed})
			},
			args:     []string{"watch", "--source", "coderabbit", "--pr", "123", "--agent", "codex", "--until-clean", "--max-rounds", "6", "--no-input"},
			wantCode: exitRunFailed,
			wantStdout: "" +
				"issue 001 failed — major: handle test issue\n" +
				"Unresolved after 1 Round(s): 0 resolved, 0 invalid, 1 failed, 0 unresolved.\n",
		},
		{
			name: "max rounds reached",
			setup: func(t *testing.T, repoDir string) {
				t.Helper()
				mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), `
resolve:
  batch_size: 1
`)
				withAgentRunner(t, &fakeAgentRunner{statuses: []string{rounds.StatusResolved, rounds.StatusFailed}})
				withFetchReviewItems(t, []reviewsource.ReviewItem{
					{
						Title:                   "major: handle first issue",
						File:                    "internal/first.go",
						Line:                    12,
						Severity:                "major",
						Author:                  "coderabbitai[bot]",
						Body:                    "First issue.",
						SourceRef:               "thread:PRRT_first,comment:PRRC_first",
						ReviewHash:              "review-hash-first",
						SourceReviewID:          "9001",
						SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC),
					},
					{
						Title:                   "major: handle second issue",
						File:                    "internal/second.go",
						Line:                    24,
						Severity:                "major",
						Author:                  "coderabbitai[bot]",
						Body:                    "Second issue.",
						SourceRef:               "thread:PRRT_second,comment:PRRC_second",
						ReviewHash:              "review-hash-second",
						SourceReviewID:          "9002",
						SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 1, 0, 0, time.UTC),
					},
				})
			},
			args: []string{"watch", "--source", "coderabbit", "--pr", "123", "--agent", "codex", "--until-clean", "--max-rounds", "1", "--no-input"},
			wantStdout: "" +
				"issue 001 resolved — major: handle first issue\n" +
				"issue 002 failed — major: handle second issue\n" +
				"MaxRoundsReached after 1 Round(s): 1 resolved, 0 invalid, 1 failed, 0 unresolved.\n",
		},
		{
			name: "stopped after fetch",
			setup: func(t *testing.T, _ string) {
				t.Helper()
				withAgentRunner(t, &fakeStoppingAgentRunner{})
				withChangedPaths(t, nil)
			},
			args: []string{"watch", "--source", "coderabbit", "--pr", "123", "--agent", "codex", "--until-clean", "--max-rounds", "6", "--no-input"},
			wantStdout: "" +
				"issue 001 unresolved — major: handle test issue\n" +
				"Stopped after 1 Round(s): 0 resolved, 0 invalid, 0 failed, 1 unresolved.\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, repoDir := withCLIWorkspace(t)
			withSuccessfulPreflight(t, repoDir)
			if tt.setup != nil {
				tt.setup(t, repoDir)
			}
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := RunContext(context.Background(), tt.args, &stdout, &stderr)

			if code != tt.wantCode {
				t.Fatalf("expected exit code %d, got %d stderr=%q", tt.wantCode, code, stderr.String())
			}
			if stdout.String() != tt.wantStdout {
				t.Fatalf("stdout mismatch\nwant:\n%q\ngot:\n%q\nstderr:\n%s", tt.wantStdout, stdout.String(), stderr.String())
			}
		})
	}
}

func TestRunResolvePrintsDeterministicStdoutReport(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--agent", "codex", "--round", "all", "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected clean resolve exit code 0, got %d stderr=%q", code, stderr.String())
	}
	wantStdout := "" +
		"issue 001 resolved — major: handle test issue\n" +
		"Clean after 1 Round(s): 1 resolved, 0 invalid, 0 failed, 0 unresolved.\n"
	if stdout.String() != wantStdout {
		t.Fatalf("stdout mismatch\nwant:\n%q\ngot:\n%q\nstderr:\n%s", wantStdout, stdout.String(), stderr.String())
	}
}

func TestRunResolvePersistsEffectiveSelection(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{
		"resolve",
		"--pr", "123",
		"--agent", "codex",
		"--model", "stored-resolve-model",
		"--reasoning-effort", "stored-resolve-reasoning",
		"--round", "all",
		"--no-input",
	}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected clean resolve exit code 0, got %d stderr=%q", code, stderr.String())
	}
	run := runFromStore(t, homeDir, reviewRunIDFromStderr(t, stderr.String()))
	if run.Model != "stored-resolve-model" || run.ReasoningEffort != "stored-resolve-reasoning" {
		t.Fatalf("expected stored resolve selection, got %#v", run)
	}
	if !strings.Contains(stderr.String(), "Agent Model: stored-resolve-model") ||
		!strings.Contains(stderr.String(), "Default Reasoning Effort: stored-resolve-reasoning") {
		t.Fatalf("expected resolve progress to show concrete selection, got %q", stderr.String())
	}
	if strings.Contains(stderr.String(), "auto") {
		t.Fatalf("expected no auto selection placeholder, got %q", stderr.String())
	}
}

func TestRunWatchPersistsEffectiveSelection(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{
		"watch",
		"--source", "coderabbit",
		"--pr", "123",
		"--agent", "codex",
		"--model", "stored-watch-model",
		"--reasoning-effort", "stored-watch-reasoning",
		"--until-clean",
		"--max-rounds", "6",
		"--no-input",
	}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected clean watch exit code 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	run := runFromStore(t, homeDir, reviewRunIDFromStderr(t, stderr.String()))
	if run.Model != "stored-watch-model" || run.ReasoningEffort != "stored-watch-reasoning" {
		t.Fatalf("expected stored watch selection, got %#v", run)
	}
	if !strings.Contains(stderr.String(), "Agent Model: stored-watch-model") ||
		!strings.Contains(stderr.String(), "Default Reasoning Effort: stored-watch-reasoning") {
		t.Fatalf("expected watch progress to show concrete selection, got %q", stderr.String())
	}
	if strings.Contains(stderr.String(), "auto") {
		t.Fatalf("expected no auto selection placeholder, got %q", stderr.String())
	}
}

func TestRunOutcomeNotificationsCaptureTerminalResolveWatchAndImplement(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		setup        func(t *testing.T) string
		runID        func(t *testing.T, stderr string) string
		wantState    string
		wantKind     string
		wantTarget   string
		wantExitCode int
	}{
		{
			name: "resolve",
			args: []string{"resolve", "--pr", "123", "--agent", "codex", "--round", "all", "--no-input"},
			setup: func(t *testing.T) string {
				t.Helper()
				homeDir, repoDir := withCLIWorkspace(t)
				withSuccessfulPreflight(t, repoDir)
				persistCLIReviewIssue(t, repoDir, 1, "feature/review")
				return homeDir
			},
			runID:      reviewRunIDFromStderr,
			wantState:  store.StateClean,
			wantKind:   store.KindResolve,
			wantTarget: "pr:123",
		},
		{
			name: "watch",
			args: []string{"watch", "--source", "coderabbit", "--pr", "123", "--agent", "codex", "--until-clean", "--max-rounds", "6", "--no-input"},
			setup: func(t *testing.T) string {
				t.Helper()
				homeDir, repoDir := withCLIWorkspace(t)
				withSuccessfulPreflight(t, repoDir)
				return homeDir
			},
			runID:      reviewRunIDFromStderr,
			wantState:  store.StateClean,
			wantKind:   store.KindWatch,
			wantTarget: "pr:123",
		},
		{
			name: "implement",
			args: []string{"implement", "--spec", implementTestSlug, "--agent", "codex", "--no-input"},
			setup: func(t *testing.T) string {
				t.Helper()
				homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
					{id: "task_01", title: "Wire notification outcome"},
				})
				withImplementCollaborators(t, &implementFakeRunner{
					gitRoot: repoDir,
					statusByTask: map[string]spec.Status{
						"task_01": spec.StatusCompleted,
					},
				})
				return homeDir
			},
			runID:      implementRunIDFromStderr,
			wantState:  store.StateClean,
			wantKind:   store.KindImplement,
			wantTarget: "spec:" + implementTestSlug,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt.setup(t)
			notifier := &recordingOutcomeNotifier{}
			withOutcomeNotifier(t, notifier)
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := RunContext(context.Background(), tt.args, &stdout, &stderr)

			if code != tt.wantExitCode {
				t.Fatalf("expected exit code %d, got %d stderr=%q stdout=%q", tt.wantExitCode, code, stderr.String(), stdout.String())
			}
			runID := tt.runID(t, stderr.String())
			assertRecordedOutcomes(t, notifier, []roundnotify.Outcome{{
				RunID:  runID,
				State:  tt.wantState,
				Kind:   tt.wantKind,
				Target: tt.wantTarget,
			}})
		})
	}
}

func TestRunOutcomeNotificationsSkipFetch(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	notifier := &recordingOutcomeNotifier{}
	withOutcomeNotifier(t, notifier)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"fetch", "--source", "coderabbit", "--pr", "123", "--round", "auto", "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected fetch exit code 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	assertRecordedOutcomes(t, notifier, nil)
}

func TestNotifyTerminalOutcomeDetachesCanceledParentWithBoundedDeadline(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())
	cancel()
	notifier := &deadlineRecordingOutcomeNotifier{}
	run := store.Run{
		ID:       "run-123",
		Kind:     store.KindResolve,
		State:    store.StateClean,
		PRNumber: "123",
	}
	var stderr bytes.Buffer
	start := time.Now()

	notifyTerminalOutcome(parent, nil, notifier, &stderr, run)

	if notifier.wasCanceled {
		t.Fatal("expected terminal outcome notification context to ignore parent cancellation")
	}
	if !notifier.hadDeadline {
		t.Fatal("expected terminal outcome notification context to have a deadline")
	}
	if !notifier.deadline.After(start) || notifier.deadline.After(start.Add(outcomeNotificationTimeout+time.Second)) {
		t.Fatalf("expected deadline within notification timeout, got %s from start %s", notifier.deadline, start)
	}
	want := roundnotify.Outcome{
		RunID:  "run-123",
		State:  store.StateClean,
		Kind:   store.KindResolve,
		Target: "pr:123",
	}
	if !reflect.DeepEqual(notifier.outcome, want) {
		t.Fatalf("notification outcome mismatch\nwant: %#v\ngot:  %#v", want, notifier.outcome)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunOutcomeNotificationFailureWarnsAndJournalsWithoutChangingReportOrExit(t *testing.T) {
	wantStdout, _, wantCode, _ := runCleanResolveForOutcomeNotification(t, nil)
	gotStdout, gotStderr, gotCode, homeDir := runCleanResolveForOutcomeNotification(t, errors.New("forced notifier failure"))

	if gotCode != wantCode {
		t.Fatalf("expected exit code to stay %d, got %d stderr=%q", wantCode, gotCode, gotStderr)
	}
	if gotStdout != wantStdout {
		t.Fatalf("stdout changed after notification failure\nwant:\n%q\ngot:\n%q", wantStdout, gotStdout)
	}
	warning := "roundfix: outcome notification failed: forced notifier failure\n"
	if strings.Count(gotStderr, warning) != 1 {
		t.Fatalf("expected one notification warning %q, got stderr=%q", warning, gotStderr)
	}
	runID, events := journaledRunEvents(t, homeDir, gotStderr)
	for _, entry := range events {
		event := entry.Event
		if event.Source == runevent.SourceDaemon &&
			event.Kind == runevent.KindDaemonStatus &&
			strings.Contains(event.Summary, "outcome notification failed") &&
			strings.Contains(string(event.Payload), "forced notifier failure") {
			return
		}
	}
	t.Fatalf("expected Daemon-source notification failure event for Run %s, got %+v", runID, events)
}

func TestRunCleanCleanupFailureWarnsAndJournalsWithoutChangingReportOrExit(t *testing.T) {
	wantStdout, _, wantCode, _, _ := runCleanResolveForCleanup(t, nil)
	gotStdout, gotStderr, gotCode, homeDir, keptPath := runCleanResolveForCleanup(t, errors.New("forced cleanup failure"))

	if gotCode != wantCode {
		t.Fatalf("expected exit code to stay %d, got %d stderr=%q", wantCode, gotCode, gotStderr)
	}
	if gotStdout != wantStdout {
		t.Fatalf("stdout changed after cleanup failure\nwant:\n%q\ngot:\n%q", wantStdout, gotStdout)
	}
	warning := fmt.Sprintf("roundfix: Run Worktree cleanup failed; kept %s: forced cleanup failure\n", keptPath)
	if strings.Count(gotStderr, warning) != 1 {
		t.Fatalf("expected one cleanup warning %q, got stderr=%q", warning, gotStderr)
	}
	assertCleanCleanupWarningEvent(t, homeDir, gotStderr, keptPath, "forced cleanup failure")
	assertRunWorktreeExists(t, keptPath)
}

func TestRunWatchCleanCleanupFailureWarnsAndJournalsWithoutChangingReportOrExit(t *testing.T) {
	wantStdout, _, wantCode, _, _ := runCleanWatchForCleanup(t, nil)
	gotStdout, gotStderr, gotCode, homeDir, keptPath := runCleanWatchForCleanup(t, errors.New("forced cleanup failure"))

	if gotCode != wantCode {
		t.Fatalf("expected exit code to stay %d, got %d stderr=%q", wantCode, gotCode, gotStderr)
	}
	if gotStdout != wantStdout {
		t.Fatalf("stdout changed after cleanup failure\nwant:\n%q\ngot:\n%q", wantStdout, gotStdout)
	}
	warning := fmt.Sprintf("roundfix: Run Worktree cleanup failed; kept %s: forced cleanup failure\n", keptPath)
	if strings.Count(gotStderr, warning) != 1 {
		t.Fatalf("expected one cleanup warning %q, got stderr=%q", warning, gotStderr)
	}
	assertCleanCleanupWarningEvent(t, homeDir, gotStderr, keptPath, "forced cleanup failure")
	assertRunWorktreeExists(t, keptPath)
}

func assertCleanCleanupWarningEvent(t *testing.T, homeDir string, stderr string, keptPath string, reason string) {
	t.Helper()
	if strings.Contains(stderr, "failed after Run start") || strings.Contains(stderr, "Run Worktree kept:") {
		t.Fatalf("cleanup warning must not turn Clean into failed/kept diagnostics, got stderr=%q", stderr)
	}
	runID, events := journaledRunEvents(t, homeDir, stderr)
	run := runFromStore(t, homeDir, runID)
	if run.State != store.StateClean {
		t.Fatalf("expected cleanup failure to leave Run Clean, got %s", run.State)
	}
	for _, entry := range events {
		event := entry.Event
		if event.Source == runevent.SourceDaemon &&
			event.Kind == runevent.KindDaemonStatus &&
			strings.Contains(event.Summary, "Run Worktree cleanup failed") &&
			strings.Contains(string(event.Payload), keptPath) &&
			strings.Contains(string(event.Payload), reason) {
			return
		}
	}
	t.Fatalf("expected Daemon-source cleanup failure event for Run %s, got %+v", runID, events)
}

func TestRunOutcomeNotificationsDisabledSkipsNotifier(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), "notify:\n  enabled: false\n")
	withSuccessfulPreflight(t, repoDir)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	withOutcomeNotifierFactory(t, func(roundconfig.Config) roundnotify.Notifier {
		t.Fatal("disabled notifications must not construct a notifier")
		return nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--agent", "codex", "--round", "all", "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected resolve exit code 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	wantStdout := "" +
		"issue 001 resolved — major: handle test issue\n" +
		"Clean after 1 Round(s): 1 resolved, 0 invalid, 0 failed, 0 unresolved.\n"
	if stdout.String() != wantStdout {
		t.Fatalf("stdout mismatch\nwant:\n%q\ngot:\n%q\nstderr:\n%s", wantStdout, stdout.String(), stderr.String())
	}
	if strings.Contains(stderr.String(), "outcome notification failed") {
		t.Fatalf("disabled notifications must be silent, got stderr=%q", stderr.String())
	}
}

func TestRunResolveWorktreeIsolationExcludesConcurrentUserCommit(t *testing.T) {
	homeDir, repoDir := withReviewGitWorkspace(t)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	withRealReviewPreflight(t, repoDir, true)
	source := &fakeSourceResolver{}
	withSourceResolver(t, source)
	pusher := &checkingPusher{}
	withPusher(t, pusher)
	runner := &fakeAgentRunner{
		onRun: func(req agent.ExecuteRequest) error {
			if err := os.WriteFile(filepath.Join(req.GitRoot, "agent.txt"), []byte("agent work\n"), 0o644); err != nil {
				return err
			}
			mustWrite(t, filepath.Join(repoDir, "user.txt"), "user work\n")
			gitImplement(t, repoDir, "add", "user.txt")
			gitImplement(t, repoDir, "commit", "-m", "user work during resolve")
			return nil
		},
	}
	withAgentRunner(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--agent", "codex", "--round", "all", "--no-input"}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("expected IntegrationPending exit, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	runID := reviewRunIDFromStderr(t, stderr.String())
	run := runFromStore(t, homeDir, runID)
	if run.State != store.StateIntegrationPending {
		t.Fatalf("expected IntegrationPending Run, got %s", run.State)
	}
	if pusher.calls != 0 {
		t.Fatalf("expected Final Push skipped on IntegrationPending, got %d calls", pusher.calls)
	}
	if len(runner.gitRoots) != 1 || runner.gitRoots[0] != run.WorkDir {
		t.Fatalf("expected Agent to run in Run Worktree %q, got %v", run.WorkDir, runner.gitRoots)
	}
	branch := runworktree.BranchName(runID)
	assertRunWorktreeExists(t, run.WorkDir)
	assertRunBranchExists(t, repoDir, branch)
	taskCommitFiles := gitImplementOutput(t, repoDir, "diff-tree", "--no-commit-id", "--name-only", "-r", branch)
	if !strings.Contains(taskCommitFiles, "agent.txt") {
		t.Fatalf("expected Run Branch commit to contain agent work, got %q", taskCommitFiles)
	}
	if strings.Contains(taskCommitFiles, "user.txt") {
		t.Fatalf("expected Run Branch commit to exclude concurrent user file, got %q", taskCommitFiles)
	}
	userFiles := gitImplementOutput(t, repoDir, "ls-tree", "-r", "--name-only", "HEAD")
	if !strings.Contains(userFiles, "user.txt") {
		t.Fatalf("expected user commit to survive on user branch, got %q", userFiles)
	}
	if strings.Contains(userFiles, "agent.txt") {
		t.Fatalf("expected unintegrated agent work to stay off user branch, got %q", userFiles)
	}
	if !strings.Contains(stderr.String(), "Integration command: git merge --ff-only "+branch) {
		t.Fatalf("expected integration command on stderr, got %q", stderr.String())
	}
}

func TestRunResolvePushRunsAfterSuccessfulIntegrationAndCleansWorktree(t *testing.T) {
	homeDir, repoDir := withReviewGitWorkspace(t)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	withRealReviewPreflight(t, repoDir, true)
	source := &fakeSourceResolver{}
	withSourceResolver(t, source)
	pusher := &checkingPusher{
		check: func(req daemon.PushRequest) error {
			// ADR-0036: Final Push runs from the user checkout after
			// integration so the review artifact docs commit rides it.
			if req.WorkDir != repoDir {
				return fmt.Errorf("push WorkDir left the user checkout: %q", req.WorkDir)
			}
			if got := gitImplementOutput(t, repoDir, "show", "HEAD:agent.txt"); got != "agent work\n" {
				return fmt.Errorf("push ran before integrated agent work reached HEAD, got %q", got)
			}
			return nil
		},
	}
	withPusher(t, pusher)
	runner := &fakeAgentRunner{
		onRun: func(req agent.ExecuteRequest) error {
			return os.WriteFile(filepath.Join(req.GitRoot, "agent.txt"), []byte("agent work\n"), 0o644)
		},
	}
	withAgentRunner(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--agent", "codex", "--round", "all", "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected clean resolve exit, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if pusher.calls != 1 {
		t.Fatalf("expected one Final Push after integration, got %d", pusher.calls)
	}
	runID := reviewRunIDFromStderr(t, stderr.String())
	run := runFromStore(t, homeDir, runID)
	if run.State != store.StateClean {
		t.Fatalf("expected Clean Run, got %s", run.State)
	}
	if len(pusher.workDir) != 1 || pusher.workDir[0] != repoDir {
		t.Fatalf("expected Final Push from the user checkout %q, got %v", repoDir, pusher.workDir)
	}
	assertRunWorktreeRemoved(t, run.WorkDir)
	assertRunBranchRemoved(t, repoDir, runworktree.BranchName(runID))
	if got := mustRead(t, filepath.Join(repoDir, "agent.txt")); got != "agent work\n" {
		t.Fatalf("expected integrated agent work in user checkout, got %q", got)
	}
	subject := gitImplementOutput(t, repoDir, "log", "-1", "--format=%s")
	if !strings.HasPrefix(subject, "docs: review round") {
		t.Fatalf("expected review artifact docs commit at HEAD, got %q", subject)
	}
}

func TestRunWatchReusesOneRealWorktreeAcrossRoundsAndCleansOnIntegratedClean(t *testing.T) {
	homeDir, repoDir := withReviewGitWorkspace(t)
	withRealReviewPreflight(t, repoDir, true)
	withWatchStatus(t, (&fakeWatchStatus{
		statuses: []reviewsource.WatchStatus{
			{State: watch.StatusSettled},
			{State: watch.StatusSettled},
		},
	}).Status)
	fetchCalls := 0
	withFetchReviewItemsFunc(t, func(context.Context, reviewsource.FetchRequest) ([]reviewsource.ReviewItem, error) {
		fetchCalls++
		return []reviewsource.ReviewItem{
			{
				Title:                   fmt.Sprintf("major: handle watch issue %d", fetchCalls),
				File:                    fmt.Sprintf("internal/watch_%d.go", fetchCalls),
				Line:                    12,
				Severity:                "major",
				Author:                  "coderabbitai[bot]",
				Body:                    "Watch issue.",
				SourceRef:               fmt.Sprintf("thread:PRRT_watch_%d,comment:PRRC_watch_%d", fetchCalls, fetchCalls),
				ReviewHash:              fmt.Sprintf("review-hash-watch-%d", fetchCalls),
				SourceReviewID:          fmt.Sprintf("90%02d", fetchCalls),
				SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, fetchCalls, 0, 0, time.UTC),
			},
		}, nil
	})
	headCheck := &fakeWatchHeadCheck{states: []watch.HeadCheckState{watch.CheckFailure, watch.CheckSuccess}}
	withWatchHeadCheck(t, headCheck.Check)
	source := &fakeSourceResolver{}
	withSourceResolver(t, source)
	pusher := &checkingPusher{}
	withPusher(t, pusher)
	runCount := 0
	runner := &fakeAgentRunner{
		onRun: func(req agent.ExecuteRequest) error {
			runCount++
			return os.WriteFile(filepath.Join(req.GitRoot, "agent.txt"), []byte(fmt.Sprintf("agent round %d\n", runCount)), 0o644)
		},
	}
	withAgentRunner(t, runner)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"watch", "--source", "coderabbit", "--pr", "123", "--agent", "codex", "--until-clean", "--max-rounds", "2", "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected clean watch exit, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if runner.calls != 2 {
		t.Fatalf("expected two watched resolve rounds, got %d", runner.calls)
	}
	if len(runner.gitRoots) != 2 || runner.gitRoots[0] == "" || runner.gitRoots[0] != runner.gitRoots[1] {
		t.Fatalf("expected one Run Worktree reused across rounds, got %v", runner.gitRoots)
	}
	if pusher.calls != 2 {
		t.Fatalf("expected Final Push after each integrated clean round, got %d", pusher.calls)
	}
	runID := reviewRunIDFromStderr(t, stderr.String())
	run := runFromStore(t, homeDir, runID)
	if run.State != store.StateClean {
		t.Fatalf("expected Clean watch Run, got %s", run.State)
	}
	if run.WorkDir != runner.gitRoots[0] {
		t.Fatalf("expected recorded work_dir %q to match Agent WorkDir %q", run.WorkDir, runner.gitRoots[0])
	}
	assertRunWorktreeRemoved(t, run.WorkDir)
	assertRunBranchRemoved(t, repoDir, runworktree.BranchName(runID))
	if got := mustRead(t, filepath.Join(repoDir, "agent.txt")); got != "agent round 2\n" {
		t.Fatalf("expected latest integrated watch work in user checkout, got %q", got)
	}
	if !strings.Contains(stdout.String(), "Clean after 2 Round(s):") {
		t.Fatalf("expected clean two-round stdout report, got %q", stdout.String())
	}
}

func TestRunWatchPrintsBudgetExceededStdoutReport(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	clock := &fakeWatchClock{now: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC)}
	withWatchTiming(t, clock, &fakeWatchSleeper{clock: clock})
	withWatchStatus(t, (&fakeWatchStatus{
		statuses: []reviewsource.WatchStatus{
			{State: watch.StatusPending},
			{State: watch.StatusSettled},
		},
	}).Status)
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), `
watch:
  poll_interval: 1s
  review_timeout: 10s
  quiet_period: 2s
budget:
  enabled: true
  max_run_duration: 2s
`)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"watch", "--source", "coderabbit", "--pr", "123", "--agent", "codex", "--until-clean", "--max-rounds", "6", "--no-input"}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("expected BudgetExceeded exit %d, got %d stderr=%q", exitRunFailed, code, stderr.String())
	}
	wantStdout := "BudgetExceeded after 0 Round(s): 0 resolved, 0 invalid, 0 failed, 0 unresolved.\n"
	if stdout.String() != wantStdout {
		t.Fatalf("expected BudgetExceeded stdout report %q, got %q", wantStdout, stdout.String())
	}
	if !strings.Contains(stderr.String(), "reached BudgetExceeded") {
		t.Fatalf("expected BudgetExceeded terminal outcome on stderr, got %q", stderr.String())
	}
	if strings.Contains(stderr.String(), "Fetched Round") || strings.Contains(stderr.String(), "fake agent output") {
		t.Fatalf("BudgetExceeded before fetch must not fetch or run Agent, got %q", stderr.String())
	}
}

func TestOperationalStdoutReportStartsAfterTerminalRunLine(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		setup          func(t *testing.T, repoDir string)
		terminalNeedle string
		wantStdout     string
	}{
		{
			name:           "watch",
			args:           []string{"watch", "--source", "coderabbit", "--pr", "123", "--agent", "codex", "--until-clean", "--max-rounds", "6", "--no-input"},
			terminalNeedle: "reached Clean",
			wantStdout: "" +
				"issue 001 resolved — major: handle test issue\n" +
				"Clean after 1 Round(s): 1 resolved, 0 invalid, 0 failed, 0 unresolved.\n",
		},
		{
			name: "resolve",
			args: []string{"resolve", "--pr", "123", "--agent", "codex", "--round", "all", "--no-input"},
			setup: func(t *testing.T, repoDir string) {
				t.Helper()
				persistCLIReviewIssue(t, repoDir, 1, "feature/review")
			},
			terminalNeedle: "reached Clean",
			wantStdout: "" +
				"issue 001 resolved — major: handle test issue\n" +
				"Clean after 1 Round(s): 1 resolved, 0 invalid, 0 failed, 0 unresolved.\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, repoDir := withCLIWorkspace(t)
			withSuccessfulPreflight(t, repoDir)
			if tt.setup != nil {
				tt.setup(t, repoDir)
			}
			recorder := &orderedWriteRecorder{}

			code := RunContext(context.Background(), tt.args, recorder.writer("stdout"), recorder.writer("stderr"))

			if code != exitOK {
				t.Fatalf("expected exit code 0, got %d stderr=%q", code, recorder.content("stderr"))
			}
			if got := recorder.content("stdout"); got != tt.wantStdout {
				t.Fatalf("stdout mismatch\nwant:\n%q\ngot:\n%q\nstderr:\n%s", tt.wantStdout, got, recorder.content("stderr"))
			}
			terminalIndex := recorder.firstIndexContaining("stderr", tt.terminalNeedle)
			if terminalIndex < 0 {
				t.Fatalf("expected terminal stderr line containing %q, got %q", tt.terminalNeedle, recorder.content("stderr"))
			}
			stdoutIndex := recorder.firstIndex("stdout")
			if stdoutIndex < 0 {
				t.Fatal("expected stdout report")
			}
			if stdoutIndex < terminalIndex {
				t.Fatalf("stdout report began before terminal Run line: events=%#v", recorder.events)
			}
		})
	}
}

func TestRunResolveAgentFullAccessIsExplicitOptIn(t *testing.T) {
	tests := []struct {
		name       string
		config     string
		args       []string
		wantAccess string
	}{
		{
			name:       "default sandboxed",
			args:       []string{"resolve", "--pr", "123", "--agent", "codex", "--round", "all", "--no-input"},
			wantAccess: "",
		},
		{
			name:       "command flag opts in",
			args:       []string{"resolve", "--pr", "123", "--agent", "codex", "--agent-full-access", "--round", "all", "--no-input"},
			wantAccess: "full-access",
		},
		{
			name: "config opts in",
			config: `
defaults:
  agent_full_access: true
`,
			args:       []string{"resolve", "--pr", "123", "--agent", "codex", "--round", "all", "--no-input"},
			wantAccess: "full-access",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, repoDir := withCLIWorkspace(t)
			withSuccessfulPreflight(t, repoDir)
			runner := &fakeAgentRunner{}
			withAgentRunner(t, runner)
			if tt.config != "" {
				mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), tt.config)
			}
			persistCLIReviewIssue(t, repoDir, 1, "feature/review")
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run(tt.args, &stdout, &stderr)

			if code != 0 {
				t.Fatalf("expected exit code 0, got %d, stdout=%q stderr=%q", code, stdout.String(), stderr.String())
			}
			if len(runner.probedRuntimes) != 1 {
				t.Fatalf("expected one runtime probe, got %#v", runner.probedRuntimes)
			}
			if runner.probedRuntimes[0].FullAccessMode != tt.wantAccess {
				t.Fatalf("expected probe full-access mode %q, got %q", tt.wantAccess, runner.probedRuntimes[0].FullAccessMode)
			}
			if len(runner.runRuntimes) != 1 {
				t.Fatalf("expected one Agent run, got %#v", runner.runRuntimes)
			}
			if runner.runRuntimes[0].FullAccessMode != tt.wantAccess {
				t.Fatalf("expected run full-access mode %q, got %q", tt.wantAccess, runner.runRuntimes[0].FullAccessMode)
			}
		})
	}
}

func TestRunFetchReusesMatchingAutoRound(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	args := []string{"fetch", "--source", "coderabbit", "--pr", "123", "--round", "auto", "--no-input"}

	var firstStdout bytes.Buffer
	var firstStderr bytes.Buffer
	firstCode := Run(args, &firstStdout, &firstStderr)
	if firstCode != 0 {
		t.Fatalf("expected first fetch exit 0, got %d stdout=%q stderr=%q", firstCode, firstStdout.String(), firstStderr.String())
	}

	var secondStdout bytes.Buffer
	var secondStderr bytes.Buffer
	secondCode := Run(args, &secondStdout, &secondStderr)
	if secondCode != 0 {
		t.Fatalf("expected second fetch exit 0, got %d stdout=%q stderr=%q", secondCode, secondStdout.String(), secondStderr.String())
	}

	if !strings.Contains(secondStdout.String(), "Artifacts: reused existing matching Round") {
		t.Fatalf("expected second fetch to report reused Round, got %q", secondStdout.String())
	}
	if _, err := os.Stat(filepath.Join(defaultReviewRootForRepo(repoDir, "123"), "round-002")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected no duplicate round-002, got err %v", err)
	}
	assertRunCount(t, filepath.Join(homeDir, ".roundfix", "roundfix.db"), 2)
}

func TestRunWatchTimeoutOffersManualReviewWithoutFetching(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withFetchReviewItemsFunc(t, func(context.Context, reviewsource.FetchRequest) ([]reviewsource.ReviewItem, error) {
		t.Fatal("watch timeout must not fetch Review Source issues")
		return nil, nil
	})
	withWatchStatus(t, (&fakeWatchStatus{
		statuses: []reviewsource.WatchStatus{
			{State: watch.StatusPending},
			{State: watch.StatusPending},
			{State: watch.StatusPending},
		},
	}).Status)
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), `
watch:
  poll_interval: 1s
  review_timeout: 2s
  quiet_period: 1s
`)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"watch", "--source", "coderabbit", "--pr", "123", "--agent", "codex", "--until-clean", "--max-rounds", "6", "--no-input"}, &stdout, &stderr)

	if code != 1 {
		t.Fatalf("expected watch timeout exit 1, got %d", code)
	}
	wantStdout := "TimedOut after 0 Round(s): 0 resolved, 0 invalid, 0 failed, 0 unresolved.\n"
	if stdout.String() != wantStdout {
		t.Fatalf("expected timeout stdout report %q, got %q", wantStdout, stdout.String())
	}
	if !strings.Contains(stderr.String(), "Review Source status: pending") {
		t.Fatalf("expected pending Review Source status output, got %q", stderr.String())
	}
	if count := strings.Count(stderr.String(), "Review Source status: pending\n"); count != 1 {
		t.Fatalf("expected one unchanged status-poll line, got %d in %q", count, stderr.String())
	}
	if !strings.Contains(stderr.String(), "reached TimedOut") {
		t.Fatalf("expected TimedOut terminal outcome, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "@coderabbitai review") {
		t.Fatalf("expected manual review trigger guidance, got %q", stderr.String())
	}
	if strings.Contains(stderr.String(), "Fetched Round") || strings.Contains(stderr.String(), "fake agent output") {
		t.Fatalf("timeout must not fetch or run Agent, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
	assertNoActiveRun(t, homeDir, "owner/project", "feature/review")
}

func TestRunWatchMissingHeadCheckPrintsCleanNote(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withWatchHeadCheck(t, (&fakeWatchHeadCheck{
		states: []watch.HeadCheckState{watch.CheckMissing},
	}).Check)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"watch", "--source", "coderabbit", "--pr", "123", "--agent", "codex", "--until-clean", "--max-rounds", "6", "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected clean watch exit 0, got %d stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "Review Source check missing for the pushed HEAD; treating Run as Clean.") {
		t.Fatalf("expected missing-check stderr note, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Expected: Watch Run Clean normally means the Review Source check on the pushed HEAD reports success.") {
		t.Fatalf("expected missing-check documentation expectation, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Next: confirm the PR's Review Source check before merging.") {
		t.Fatalf("expected missing-check next action, got %q", stderr.String())
	}
	if strings.Count(stderr.String(), "Review Source check missing for the pushed HEAD") != 1 {
		t.Fatalf("expected one missing-check stderr note, got %q", stderr.String())
	}
	wantStdout := "" +
		"issue 001 resolved — major: handle test issue\n" +
		"Clean after 1 Round(s): 1 resolved, 0 invalid, 0 failed, 0 unresolved.\n"
	if stdout.String() != wantStdout {
		t.Fatalf("expected Clean stdout report %q, got %q", wantStdout, stdout.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
	assertNoActiveRun(t, homeDir, "owner/project", "feature/review")
}

func TestRunWatchReusesOneAgentSessionAcrossRoundsAndCloses(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	inner := &fakeAgentRunner{statuses: []string{rounds.StatusResolved, rounds.StatusFailed, rounds.StatusResolved}}
	runner := &sessionRecordingRunner{inner: inner}
	withAgentRunner(t, runner)
	withFakeWorktree(t)
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), `
resolve:
  batch_size: 1
`)
	fetchCalls := 0
	withFetchReviewItemsFunc(t, func(context.Context, reviewsource.FetchRequest) ([]reviewsource.ReviewItem, error) {
		fetchCalls++
		switch fetchCalls {
		case 1:
			return []reviewsource.ReviewItem{
				{
					Title:                   "major: handle first issue",
					File:                    "internal/first.go",
					Line:                    12,
					Severity:                "major",
					Author:                  "coderabbitai[bot]",
					Body:                    "First issue.",
					SourceRef:               "thread:PRRT_first,comment:PRRC_first",
					ReviewHash:              "review-hash-first",
					SourceReviewID:          "9001",
					SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC),
				},
				{
					Title:                   "major: handle retry issue",
					File:                    "internal/retry.go",
					Line:                    24,
					Severity:                "major",
					Author:                  "coderabbitai[bot]",
					Body:                    "Retry issue.",
					SourceRef:               "thread:PRRT_retry_1,comment:PRRC_retry_1",
					ReviewHash:              "review-hash-retry",
					SourceReviewID:          "9002",
					SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 1, 0, 0, time.UTC),
				},
			}, nil
		default:
			return []reviewsource.ReviewItem{
				{
					Title:                   "major: handle retry issue",
					File:                    "internal/retry.go",
					Line:                    24,
					Severity:                "major",
					Author:                  "coderabbitai[bot]",
					Body:                    "Retry issue.",
					SourceRef:               "thread:PRRT_retry_2,comment:PRRC_retry_2",
					ReviewHash:              "review-hash-retry",
					SourceReviewID:          "9003",
					SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 2, 0, 0, time.UTC),
				},
			}, nil
		}
	})
	withWatchStatus(t, (&fakeWatchStatus{
		statuses: []reviewsource.WatchStatus{
			{State: watch.StatusSettled},
			{State: watch.StatusSettled},
		},
	}).Status)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"watch", "--source", "coderabbit", "--pr", "123", "--agent", "codex", "--until-clean", "--max-rounds", "2", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean watch exit 0, got %d stderr=%q", code, stderr.String())
	}
	if fetchCalls != 2 {
		t.Fatalf("expected two fetched Rounds, got %d", fetchCalls)
	}
	if !strings.Contains(stderr.String(), "after 2 Round(s)") {
		t.Fatalf("expected Watch Run to span two Rounds, got %q", stderr.String())
	}
	runID, events := journaledRunEvents(t, homeDir, stderr.String())
	assertRecordedOneSessionForRun(t, runner, runID, inner.calls)
	assertSessionLifecycleEvents(t, events)
	assertNoActiveRun(t, homeDir, "owner/project", "feature/review")
}

func TestRunWatchStopRequestBeforeAgentMarksStopped(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	withSuccessfulPreflight(t, repoDir)
	withWatchStatus(t, func(context.Context, reviewsource.WatchStatusRequest) (reviewsource.WatchStatus, error) {
		cancel()
		return reviewsource.WatchStatus{State: watch.StatusPending}, nil
	})
	withFetchReviewItemsFunc(t, func(context.Context, reviewsource.FetchRequest) ([]reviewsource.ReviewItem, error) {
		t.Fatal("Stop Request before Agent must not fetch Review Source issues")
		return nil, nil
	})
	withChangedPaths(t, nil)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(ctx, []string{"watch", "--source", "coderabbit", "--pr", "123", "--agent", "codex", "--until-clean", "--max-rounds", "6", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean Stop Request exit 0, got %d", code)
	}
	wantStdout := "Stopped after 0 Round(s): 0 resolved, 0 invalid, 0 failed, 0 unresolved.\n"
	if stdout.String() != wantStdout {
		t.Fatalf("expected stopped stdout report %q, got %q", wantStdout, stdout.String())
	}
	if !strings.Contains(stderr.String(), "Watch Run") || !strings.Contains(stderr.String(), "reached Stopped") {
		t.Fatalf("expected stopped Watch Run output, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Changed paths after Stop Request: none") {
		t.Fatalf("expected changed-path report, got %q", stderr.String())
	}
	if strings.Contains(stderr.String(), "Fetched Round") || strings.Contains(stderr.String(), "fake agent output") {
		t.Fatalf("Stop Request before Agent must not fetch or run Agent, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
	assertNoActiveRun(t, homeDir, "owner/project", "feature/review")
}

func TestRunResolveStopRequestDuringAgentPreservesWorkAndSkipsDaemonMutations(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	verifier := &fakeVerifier{}
	committer := &fakeCommitter{}
	sourceResolver := &fakeSourceResolver{}
	pusher := &fakePusher{}
	withVerifier(t, verifier)
	withCommitter(t, committer)
	withFakeWorktree(t)
	withSourceResolver(t, sourceResolver)
	withPusher(t, pusher)
	withAgentRunner(t, &fakeStoppingAgentRunner{})
	withFakeWorktree(t)
	withChangedPaths(t, []preflight.ChangedPath{{Status: "M", Path: "src/app.go"}})
	result := persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean Stop Request exit 0, got %d", code)
	}
	wantStdout := "" +
		"issue 001 unresolved — major: handle test issue\n" +
		"Stopped after 1 Round(s): 0 resolved, 0 invalid, 0 failed, 1 unresolved.\n"
	if stdout.String() != wantStdout {
		t.Fatalf("expected stopped stdout report %q, got %q", wantStdout, stdout.String())
	}
	if !strings.Contains(stderr.String(), "Resolve Run") || !strings.Contains(stderr.String(), "reached Stopped") {
		t.Fatalf("expected stopped Resolve Run output, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "partial agent output") {
		t.Fatalf("expected persisted Agent output in stream, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "M src/app.go") {
		t.Fatalf("expected changed path report, got %q", stderr.String())
	}
	if verifier.calls != 0 {
		t.Fatalf("Stop Request must not run verification, got %d calls", verifier.calls)
	}
	if committer.calls != 0 {
		t.Fatalf("Stop Request must not create commits, got %d calls", committer.calls)
	}
	if sourceResolver.calls != 0 {
		t.Fatalf("Stop Request must not resolve Review Source threads, got %d calls", sourceResolver.calls)
	}
	if pusher.calls != 0 {
		t.Fatalf("Stop Request must not push, got %d calls", pusher.calls)
	}
	issue, err := rounds.ParseIssue(result.IssuePaths[0])
	if err != nil {
		t.Fatalf("parse issue: %v", err)
	}
	if issue.Status != rounds.StatusPending {
		t.Fatalf("Stop Request must preserve assigned issue status, got %q", issue.Status)
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
	assertNoActiveRun(t, homeDir, "owner/project", "feature/review")

}

func TestRunResolveSIGINTStopReturns130(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withAgentRunner(t, &fakeStoppingAgentRunner{})
	withFakeWorktree(t)
	withChangedPaths(t, nil)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	stoppedCode := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)
	code := exitForInterrupt(stoppedCode, true)

	if code != 130 {
		t.Fatalf("expected SIGINT exit 130, got %d", code)
	}
	if !strings.Contains(stderr.String(), "reached Stopped") {
		t.Fatalf("expected Stopped output before SIGINT exit, got %q", stderr.String())
	}
}

func TestExitForWatchOutcome(t *testing.T) {
	tests := []struct {
		outcome string
		code    int
	}{
		{outcome: store.StateClean, code: 0},
		{outcome: store.StateMaxRoundsReached, code: 0},
		{outcome: store.StateStopped, code: 0},
		{outcome: store.StateBudgetExceeded, code: 1},
		{outcome: store.StateTimedOut, code: 1},
		{outcome: store.StateFailed, code: 1},
		{outcome: store.StateUnresolved, code: 1},
	}

	for _, tt := range tests {
		t.Run(tt.outcome, func(t *testing.T) {
			if code := exitForWatchOutcome(tt.outcome); code != tt.code {
				t.Fatalf("expected exit code %d for %s, got %d", tt.code, tt.outcome, code)
			}
		})
	}
}

func TestRunResolveRejectsMissingCompatibleArtifactsBeforeRun(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withFetchReviewItemsFunc(t, func(context.Context, reviewsource.FetchRequest) ([]reviewsource.ReviewItem, error) {
		t.Fatal("resolve must not fetch Review Source issues")
		return nil, nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "no downloaded Unresolved Review Issues in Compatible Artifacts") {
		t.Fatalf("expected missing Compatible Artifacts error, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "roundfix fetch --source coderabbit --pr 123") {
		t.Fatalf("expected fetch guidance, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "did not create a Run") {
		t.Fatalf("expected no side-effect confirmation, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 0)
}

func TestRunResolveHonorsRoundSelector(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withFetchReviewItemsFunc(t, func(context.Context, reviewsource.FetchRequest) ([]reviewsource.ReviewItem, error) {
		t.Fatal("resolve must not fetch Review Source issues")
		return nil, nil
	})
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	persistCLIReviewIssue(t, repoDir, 2, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--agent", "codex", "--round", "2", "--no-input"}, &stdout, &stderr)

	if code != 1 {
		t.Fatalf("expected Unresolved resolve exit code 1, got %d", code)
	}
	wantStdout := "" +
		"issue 001 unresolved — major: handle test issue\n" +
		"issue 002 resolved — major: handle test issue\n" +
		"Unresolved after 1 Round(s): 1 resolved, 0 invalid, 0 failed, 1 unresolved.\n"
	if stdout.String() != wantStdout {
		t.Fatalf("expected Unresolved stdout report %q, got %q", wantStdout, stdout.String())
	}
	if !strings.Contains(stderr.String(), "selected 1 downloaded Unresolved Review Issue") {
		t.Fatalf("expected one selected issue, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Round scope: 002") {
		t.Fatalf("expected round scope in stderr, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Final Push blocked: 1 Unresolved Review Issue") {
		t.Fatalf("expected Final Push to be blocked by unselected unresolved issue, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "reached Unresolved") {
		t.Fatalf("expected Unresolved outcome for out-of-scope unresolved issue, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
}

func TestRunResolveDeduplicatesBeforeBatching(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	sourceResolver := &fakeSourceResolver{}
	withSourceResolver(t, sourceResolver)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	persistCLIReviewIssue(t, repoDir, 2, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected successful resolve exit code 0, got %d", code)
	}
	wantStdout := "" +
		"issue 001 duplicated — major: handle test issue\n" +
		"issue 002 resolved — major: handle test issue\n" +
		"Clean after 1 Round(s): 1 resolved, 0 invalid, 0 failed, 0 unresolved.\n"
	if stdout.String() != wantStdout {
		t.Fatalf("expected duplicate stdout report %q, got %q", wantStdout, stdout.String())
	}
	if !strings.Contains(stderr.String(), "selected 2 downloaded Unresolved Review Issue(s)") {
		t.Fatalf("expected selected issue count, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "assigned 1 newest occurrence(s) into 1 Batch(es)") {
		t.Fatalf("expected deduped assignment count, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "associated 1 older duplicate occurrence(s)") {
		t.Fatalf("expected duplicate association count, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Marked 1 older duplicate Review Issue occurrence") {
		t.Fatalf("expected older duplicate marker, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Resolved 1 Review Source thread") {
		t.Fatalf("expected only newest source thread resolution, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Final Push completed") {
		t.Fatalf("expected Final Push after all duplicates terminal, got %q", stderr.String())
	}
	if sourceResolver.calls != 1 {
		t.Fatalf("expected one Review Source resolution call, got %d", sourceResolver.calls)
	}
	if got := len(sourceResolver.requests[0].Issues); got != 1 {
		t.Fatalf("expected only newest issue to resolve source thread, got %d", got)
	}
	if sourceResolver.requests[0].Issues[0].FilePath != filepath.Join(defaultReviewRootForRepo(repoDir, "123"), "round-002", "issue_001.md") {
		t.Fatalf("expected newest duplicate source resolution, got %#v", sourceResolver.requests[0].Issues[0])
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
}

func TestRunResolveVerificationFailureDoesNotCommit(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	verifier := &fakeVerifier{err: errors.New("tests failed")}
	committer := &fakeCommitter{}
	sourceResolver := &fakeSourceResolver{}
	pusher := &fakePusher{}
	withVerifier(t, verifier)
	withCommitter(t, committer)
	withFakeWorktree(t)
	withSourceResolver(t, sourceResolver)
	withPusher(t, pusher)
	result := persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != 1 {
		t.Fatalf("expected Run failure exit 1, got %d", code)
	}
	wantStdout := "" +
		"issue 001 failed — major: handle test issue\n" +
		"Unresolved after 1 Round(s): 0 resolved, 0 invalid, 1 failed, 0 unresolved.\n"
	if stdout.String() != wantStdout {
		t.Fatalf("expected failed stdout report %q, got %q", wantStdout, stdout.String())
	}
	if verifier.calls != 1 {
		t.Fatalf("expected one verification call, got %d", verifier.calls)
	}
	if committer.calls != 0 {
		t.Fatalf("expected no Batch commit after verification failure, got %d", committer.calls)
	}
	if sourceResolver.calls != 0 {
		t.Fatalf("expected no Review Source resolution after verification failure, got %d", sourceResolver.calls)
	}
	if pusher.calls != 0 {
		t.Fatalf("expected no Final Push after verification failure, got %d", pusher.calls)
	}
	if !strings.Contains(stderr.String(), "tests failed") {
		t.Fatalf("expected verification failure, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "reached Unresolved") {
		t.Fatalf("expected Unresolved outcome after failed verification, got %q", stderr.String())
	}
	issue, err := rounds.ParseIssue(result.IssuePaths[0])
	if err != nil {
		t.Fatalf("parse issue: %v", err)
	}
	if issue.Status != rounds.StatusFailed {
		t.Fatalf("expected failed issue status, got %q", issue.Status)
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
	assertNoActiveRun(t, homeDir, "owner/project", "feature/review")
}

func TestRunResolveProcessesAllBatchesBeforeFinalPush(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	verifier := &fakeVerifier{}
	committer := &fakeCommitter{}
	sourceResolver := &fakeSourceResolver{}
	pusher := &fakePusher{}
	withVerifier(t, verifier)
	withCommitter(t, committer)
	withFakeWorktree(t)
	withSourceResolver(t, sourceResolver)
	withPusher(t, pusher)
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), `
resolve:
  batch_size: 1
`)
	result := persistCLIReviewItems(t, repoDir, 1, "feature/review", []reviewsource.ReviewItem{
		{
			Title:                   "major: handle first issue",
			File:                    "internal/first.go",
			Line:                    12,
			Severity:                "major",
			Author:                  "coderabbitai[bot]",
			Body:                    "First issue.",
			SourceRef:               "thread:PRRT_first,comment:PRRC_first",
			ReviewHash:              "review-hash-first",
			SourceReviewID:          "9001",
			SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC),
		},
		{
			Title:                   "major: handle second issue",
			File:                    "internal/second.go",
			Line:                    24,
			Severity:                "major",
			Author:                  "coderabbitai[bot]",
			Body:                    "Second issue.",
			SourceRef:               "thread:PRRT_second,comment:PRRC_second",
			ReviewHash:              "review-hash-second",
			SourceReviewID:          "9002",
			SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 1, 0, 0, time.UTC),
		},
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected successful resolve exit 0, got %d", code)
	}
	wantStdout := "" +
		"issue 001 resolved — major: handle first issue\n" +
		"issue 002 resolved — major: handle second issue\n" +
		"Clean after 1 Round(s): 2 resolved, 0 invalid, 0 failed, 0 unresolved.\n"
	if stdout.String() != wantStdout {
		t.Fatalf("expected clean stdout report %q, got %q", wantStdout, stdout.String())
	}
	if verifier.calls != 2 {
		t.Fatalf("expected one verification call per Batch, got %d", verifier.calls)
	}
	if committer.calls != 2 {
		t.Fatalf("expected one commit per Batch, got %d", committer.calls)
	}
	if sourceResolver.calls != 2 {
		t.Fatalf("expected one Review Source resolution call per Batch, got %d", sourceResolver.calls)
	}
	if got := len(sourceResolver.requests[0].Issues); got != 1 {
		t.Fatalf("expected first Batch source resolution to include one issue, got %d", got)
	}
	if got := len(sourceResolver.requests[1].Issues); got != 1 {
		t.Fatalf("expected second Batch source resolution to include one issue, got %d", got)
	}
	if pusher.calls != 1 {
		t.Fatalf("expected Final Push after all Batches complete, got %d", pusher.calls)
	}
	if committer.messages[0] != daemon.BatchCommitMessage(1) {
		t.Fatalf("expected Batch commit message %q, got %q", daemon.BatchCommitMessage(1), committer.messages[0])
	}
	if committer.messages[1] != daemon.BatchCommitMessage(2) {
		t.Fatalf("expected Batch commit message %q, got %q", daemon.BatchCommitMessage(2), committer.messages[1])
	}
	if !strings.Contains(stderr.String(), "assigned 2 newest occurrence(s) into 2 Batch(es)") {
		t.Fatalf("expected two planned Batches, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Batch: 001/002") || !strings.Contains(stderr.String(), "Batch: 002/002") {
		t.Fatalf("expected both Batch progress lines, got %q", stderr.String())
	}
	first, err := rounds.ParseIssue(result.IssuePaths[0])
	if err != nil {
		t.Fatalf("parse first issue: %v", err)
	}
	second, err := rounds.ParseIssue(result.IssuePaths[1])
	if err != nil {
		t.Fatalf("parse second issue: %v", err)
	}
	if first.Status != rounds.StatusResolved {
		t.Fatalf("expected first Batch issue resolved, got %q", first.Status)
	}
	if second.Status != rounds.StatusResolved {
		t.Fatalf("expected second Batch issue resolved, got %q", second.Status)
	}
	if strings.Contains(stderr.String(), "Final Push blocked") {
		t.Fatalf("did not expect Final Push blocked after all Batches complete, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Final Push completed: git push origin HEAD:feature/review") {
		t.Fatalf("expected Final Push output, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
	assertNoActiveRun(t, homeDir, "owner/project", "feature/review")
}

func TestRunResolveFinalPushRunsOnceAfterAllUnresolvedTerminal(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	verifier := &fakeVerifier{}
	committer := &fakeCommitter{}
	sourceResolver := &fakeSourceResolver{}
	pusher := &fakePusher{}
	withVerifier(t, verifier)
	withCommitter(t, committer)
	withFakeWorktree(t)
	withSourceResolver(t, sourceResolver)
	withPusher(t, pusher)
	persistCLIReviewItems(t, repoDir, 1, "feature/review", []reviewsource.ReviewItem{
		{
			Title:                   "major: handle first issue",
			File:                    "internal/first.go",
			Line:                    12,
			Severity:                "major",
			Author:                  "coderabbitai[bot]",
			Body:                    "First issue.",
			SourceRef:               "thread:PRRT_first,comment:PRRC_first",
			ReviewHash:              "review-hash-first",
			SourceReviewID:          "9001",
			SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC),
		},
	})
	persistCLIReviewItems(t, repoDir, 2, "feature/review", []reviewsource.ReviewItem{
		{
			Title:                   "major: handle second issue",
			File:                    "internal/second.go",
			Line:                    24,
			Severity:                "major",
			Author:                  "coderabbitai[bot]",
			Body:                    "Second issue.",
			SourceRef:               "thread:PRRT_second,comment:PRRC_second",
			ReviewHash:              "review-hash-second",
			SourceReviewID:          "9002",
			SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 1, 0, 0, time.UTC),
		},
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected successful resolve exit 0, got %d", code)
	}
	if verifier.calls != 1 {
		t.Fatalf("expected one verification call for one Batch, got %d", verifier.calls)
	}
	if committer.calls != 1 {
		t.Fatalf("expected one Batch commit, got %d", committer.calls)
	}
	if sourceResolver.calls != 1 {
		t.Fatalf("expected one Review Source resolution call, got %d", sourceResolver.calls)
	}
	if got := len(sourceResolver.requests[0].Issues); got != 2 {
		t.Fatalf("expected both assigned terminal issues to be source-resolved together, got %d", got)
	}
	if pusher.calls != 1 {
		t.Fatalf("expected one Final Push, got %d", pusher.calls)
	}
	if pusher.remotes[0] != "origin" || pusher.branches[0] != "feature/review" {
		t.Fatalf("expected push to origin feature/review, got %s %s", pusher.remotes[0], pusher.branches[0])
	}
	if strings.Contains(strings.Join(pusher.args, " "), "--force") {
		t.Fatalf("Final Push must not force-push, got args %#v", pusher.args)
	}
	if !strings.Contains(stderr.String(), "Final Push completed: git push origin HEAD:feature/review") {
		t.Fatalf("expected Final Push output, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
	assertNoActiveRun(t, homeDir, "owner/project", "feature/review")
}

func TestRunResolveResolvesInvalidAssignedIssueSourceThread(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	sourceResolver := &fakeSourceResolver{}
	withAgentRunner(t, &fakeAgentRunner{status: rounds.StatusInvalid})
	withFakeWorktree(t)
	withSourceResolver(t, sourceResolver)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected successful resolve exit 0, got %d", code)
	}
	if sourceResolver.calls != 1 {
		t.Fatalf("expected one Review Source resolution call, got %d", sourceResolver.calls)
	}
	if got := sourceResolver.requests[0].Issues[0].Status; got != rounds.StatusInvalid {
		t.Fatalf("expected invalid issue to resolve source thread, got %q", got)
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
}

func TestRunResolveProbeFailureDoesNotCreateRun(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withAgentRunner(t, &fakeAgentRunner{probeErr: errors.New("adapter missing")})
	withFakeWorktree(t)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected probe preflight exit 2, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "adapter missing") {
		t.Fatalf("expected probe diagnostic, got %q", stderr.String())
	}
	assertNoRunDatabase(t, homeDir)
}

func TestRunResolveSelectionPreflightRejectionReportsTupleAndCreatesNoRun(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	probeErr := &agent.SelectionPreflightError{
		Runtime:         "codex",
		Model:           "gpt-5.6-sol",
		ReasoningEffort: "xhigh",
		Err:             errors.New("missing model metadata"),
	}
	runner := &fakeAgentRunner{probeErr: probeErr}
	withAgentRunner(t, runner)
	withFakeWorktree(t)
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), `
runtimes:
  codex:
    model: gpt-5.6-sol
    reasoning_effort: xhigh
`)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("expected selection preflight exit %d, got %d", exitPreflight, code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	for _, want := range []string{
		"codex",
		"gpt-5.6-sol",
		"xhigh",
		"missing model metadata",
		"update the ACP Runtime or adapter",
		"choose supported Agent Model and Default Reasoning Effort values",
		`set runtimes.codex.reasoning_effort "" when the model manages reasoning`,
	} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("expected stderr to contain %q, got %q", want, stderr.String())
		}
	}
	if runner.calls != 0 {
		t.Fatalf("expected no durable Agent Session run, got %d run call(s)", runner.calls)
	}
	if len(runner.probeRequests) != 1 {
		t.Fatalf("expected one selection preflight attempt, got %#v", runner.probeRequests)
	}
	if runner.probeRequests[0].WorkDir != repoDir {
		t.Fatalf("expected preflight workdir %q, got %q", repoDir, runner.probeRequests[0].WorkDir)
	}
	if got := runner.probeRequests[0].Runtime; got.Model != "gpt-5.6-sol" || got.ReasoningEffort != "xhigh" {
		t.Fatalf("expected exact selection without fallback, got %#v", got)
	}
	assertNoRunDatabase(t, homeDir)
}

func TestRunWatchSelectionPreflightFailureCreatesNoRun(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	runner := &fakeAgentRunner{probeErr: &agent.SelectionPreflightError{
		Runtime:         "codex",
		Model:           "gpt-5.5",
		ReasoningEffort: "unsupported",
		Err:             errors.New("reasoning rejected"),
	}}
	withAgentRunner(t, runner)
	withFakeWorktree(t)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"watch", "--source", "coderabbit", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("expected selection preflight exit %d, got %d", exitPreflight, code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "reasoning rejected") || !strings.Contains(stderr.String(), "unsupported") {
		t.Fatalf("expected reasoning rejection diagnostic, got %q", stderr.String())
	}
	if runner.calls != 0 {
		t.Fatalf("expected no durable Agent Session run, got %d run call(s)", runner.calls)
	}
	if len(runner.probeRequests) != 1 || runner.probeRequests[0].WorkDir != repoDir {
		t.Fatalf("expected one selection preflight in git root %q, got %#v", repoDir, runner.probeRequests)
	}
	assertNoRunDatabase(t, homeDir)
}

func TestRunReviewAgentCommandsPassOneRunSelectionOverridesToPreflight(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "resolve",
			args: []string{"resolve", "--pr", "123", "--agent", "codex", "--model", "future-model", "--reasoning-effort", "experimental-reasoning", "--no-input"},
		},
		{
			name: "watch",
			args: []string{"watch", "--source", "coderabbit", "--pr", "123", "--agent", "codex", "--model", "future-model", "--reasoning-effort", "experimental-reasoning", "--no-input"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, repoDir := withCLIWorkspace(t)
			withSuccessfulPreflight(t, repoDir)
			runner := &fakeAgentRunner{probeErr: errors.New("stop after selection preflight")}
			withAgentRunner(t, runner)
			if tt.name == "resolve" {
				persistCLIReviewIssue(t, repoDir, 1, "feature/review")
			}
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run(tt.args, &stdout, &stderr)

			if code != exitPreflight {
				t.Fatalf("expected preflight exit %d, got %d stderr=%q", exitPreflight, code, stderr.String())
			}
			if stdout.Len() != 0 {
				t.Fatalf("expected no stdout diagnostics, got %q", stdout.String())
			}
			if len(runner.probeRequests) != 1 {
				t.Fatalf("expected one selection preflight, got %#v", runner.probeRequests)
			}
			got := runner.probeRequests[0].Runtime
			if got.Model != "future-model" || got.ReasoningEffort != "experimental-reasoning" {
				t.Fatalf("expected one-Run selection overrides at preflight, got %#v", got)
			}
			assertNoRunDatabase(t, homeDir)
		})
	}
}

func TestRunReviewAgentCommandsRejectExplicitEmptySelectionOverrides(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "resolve model",
			args: []string{"resolve", "--pr", "123", "--agent", "codex", "--model=", "--no-input"},
			want: "--model cannot be empty",
		},
		{
			name: "resolve reasoning",
			args: []string{"resolve", "--pr", "123", "--agent", "codex", "--reasoning-effort=", "--no-input"},
			want: "--reasoning-effort cannot be empty",
		},
		{
			name: "watch model",
			args: []string{"watch", "--source", "coderabbit", "--pr", "123", "--agent", "codex", "--model=", "--no-input"},
			want: "--model cannot be empty",
		},
		{
			name: "watch reasoning",
			args: []string{"watch", "--source", "coderabbit", "--pr", "123", "--agent", "codex", "--reasoning-effort=", "--no-input"},
			want: "--reasoning-effort cannot be empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, _ := withCLIWorkspace(t)
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run(tt.args, &stdout, &stderr)

			if code != exitPreflight {
				t.Fatalf("expected preflight exit %d, got %d stderr=%q", exitPreflight, code, stderr.String())
			}
			if stdout.Len() != 0 {
				t.Fatalf("expected no stdout diagnostics, got %q", stdout.String())
			}
			if !strings.Contains(stderr.String(), tt.want) {
				t.Fatalf("expected stderr to contain %q, got %q", tt.want, stderr.String())
			}
			assertNoRunDatabase(t, homeDir)
		})
	}
}

func TestRunResolveACPXProbeFailureReportsActionablePreflight(t *testing.T) {
	installCommand := "npm install -g acpx@" + agent.PinnedACPXVersion
	tests := []struct {
		name     string
		runner   agent.Runner
		contains []string
	}{
		{
			name:   "missing binary",
			runner: &agent.ACPXRunner{Command: filepath.Join(t.TempDir(), "missing-acpx")},
			contains: []string{
				"acpx is required but was not found on PATH",
				installCommand,
			},
		},
		{
			name:   "version mismatch",
			runner: &agent.ACPXRunner{Command: fakeACPXVersionCommand(t, "0.11.0")},
			contains: []string{
				"found 0.11.0",
				"requires " + agent.PinnedACPXVersion,
				"upgrade or downgrade",
				installCommand,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, repoDir := withCLIWorkspace(t)
			withSuccessfulPreflight(t, repoDir)
			withAgentRunner(t, tt.runner)
			withFakeWorktree(t)
			persistCLIReviewIssue(t, repoDir, 1, "feature/review")
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run([]string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

			if code != 2 {
				t.Fatalf("expected probe preflight exit 2, got %d", code)
			}
			if stdout.Len() != 0 {
				t.Fatalf("expected no stdout, got %q", stdout.String())
			}
			for _, want := range tt.contains {
				if !strings.Contains(stderr.String(), want) {
					t.Fatalf("expected stderr to contain %q, got %q", want, stderr.String())
				}
			}
			if strings.Count(stderr.String(), installCommand) != 1 {
				t.Fatalf("expected one actionable install command, got %q", stderr.String())
			}
			if !strings.Contains(stderr.String(), "Roundfix did not create a Run") {
				t.Fatalf("expected no-side-effects line, got %q", stderr.String())
			}
			assertNoRunDatabase(t, homeDir)
		})
	}
}

func TestRunResolveAgentFailureMarksBatchFailed(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withAgentRunner(t, &fakeAgentRunner{runErr: errors.New("agent crashed")})
	withFakeWorktree(t)
	result := persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != 1 {
		t.Fatalf("expected Run failure exit 1, got %d", code)
	}
	wantStdout := "" +
		"issue 001 failed — major: handle test issue\n" +
		"Unresolved after 1 Round(s): 0 resolved, 0 invalid, 1 failed, 0 unresolved.\n"
	if stdout.String() != wantStdout {
		t.Fatalf("expected failed stdout report %q, got %q", wantStdout, stdout.String())
	}
	if !strings.Contains(stderr.String(), "agent crashed") {
		t.Fatalf("expected Agent failure, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "reached Unresolved") {
		t.Fatalf("expected Unresolved outcome after Agent failure, got %q", stderr.String())
	}
	issue, err := rounds.ParseIssue(result.IssuePaths[0])
	if err != nil {
		t.Fatalf("parse issue: %v", err)
	}
	if issue.Status != rounds.StatusFailed {
		t.Fatalf("expected failed issue status, got %q", issue.Status)
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
	assertNoActiveRun(t, homeDir, "owner/project", "feature/review")
}

func TestRunResolveAgentFailureContinuesWithLaterBatches(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	runner := &fakeAgentRunner{runErr: errors.New("agent crashed")}
	withAgentRunner(t, runner)
	withFakeWorktree(t)
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), `
resolve:
  batch_size: 1
`)
	result := persistCLIReviewItems(t, repoDir, 1, "feature/review", []reviewsource.ReviewItem{
		{
			Title:                   "major: handle first issue",
			File:                    "internal/first.go",
			Line:                    12,
			Severity:                "major",
			Author:                  "coderabbitai[bot]",
			Body:                    "First issue.",
			SourceRef:               "thread:PRRT_first,comment:PRRC_first",
			ReviewHash:              "review-hash-first",
			SourceReviewID:          "9001",
			SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC),
		},
		{
			Title:                   "major: handle second issue",
			File:                    "internal/second.go",
			Line:                    24,
			Severity:                "major",
			Author:                  "coderabbitai[bot]",
			Body:                    "Second issue.",
			SourceRef:               "thread:PRRT_second,comment:PRRC_second",
			ReviewHash:              "review-hash-second",
			SourceReviewID:          "9002",
			SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 1, 0, 0, time.UTC),
		},
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != 1 {
		t.Fatalf("expected Unresolved exit 1, got %d", code)
	}
	if runner.calls != 2 {
		t.Fatalf("expected later Batches to run after the first failed, got %d Agent calls", runner.calls)
	}
	first, err := rounds.ParseIssue(result.IssuePaths[0])
	if err != nil {
		t.Fatalf("parse first issue: %v", err)
	}
	second, err := rounds.ParseIssue(result.IssuePaths[1])
	if err != nil {
		t.Fatalf("parse second issue: %v", err)
	}
	if first.Status != rounds.StatusFailed {
		t.Fatalf("expected first Batch issue failed, got %q", first.Status)
	}
	if second.Status != rounds.StatusFailed {
		t.Fatalf("expected second Batch issue failed after its own Agent failure, got %q", second.Status)
	}
	if !strings.Contains(stderr.String(), "Batch 001 failed") || !strings.Contains(stderr.String(), "Batch 002 failed") {
		t.Fatalf("expected both Batch failures reported, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "reached Unresolved") {
		t.Fatalf("expected Unresolved outcome, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
	assertNoActiveRun(t, homeDir, "owner/project", "feature/review")
}

func TestRunResolveUsesOneAgentSessionPerRunAndCloses(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	inner := &fakeAgentRunner{}
	runner := &sessionRecordingRunner{inner: inner}
	withAgentRunner(t, runner)
	withFakeWorktree(t)
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), `
resolve:
  batch_size: 1
`)
	persistCLIReviewItems(t, repoDir, 1, "feature/review", []reviewsource.ReviewItem{
		{
			Title:                   "major: handle first issue",
			File:                    "internal/first.go",
			Line:                    12,
			Severity:                "major",
			Author:                  "coderabbitai[bot]",
			Body:                    "First issue.",
			SourceRef:               "thread:PRRT_first,comment:PRRC_first",
			ReviewHash:              "review-hash-first",
			SourceReviewID:          "9001",
			SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC),
		},
		{
			Title:                   "major: handle second issue",
			File:                    "internal/second.go",
			Line:                    24,
			Severity:                "major",
			Author:                  "coderabbitai[bot]",
			Body:                    "Second issue.",
			SourceRef:               "thread:PRRT_second,comment:PRRC_second",
			ReviewHash:              "review-hash-second",
			SourceReviewID:          "9002",
			SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 1, 0, 0, time.UTC),
		},
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean resolve exit 0, got %d stderr=%q", code, stderr.String())
	}
	if inner.calls != 2 {
		t.Fatalf("expected two Agent calls for two Batches, got %d", inner.calls)
	}
	runID, events := journaledRunEvents(t, homeDir, stderr.String())
	assertRecordedOneSessionForRun(t, runner, runID, inner.calls)
	assertSessionLifecycleEvents(t, events)
}

func TestRunResolveClosesAgentSessionForTerminalOutcomes(t *testing.T) {
	tests := []struct {
		name      string
		inner     agent.Runner
		pusherErr error
		wantCode  int
		wantState string
		closeErr  error
	}{
		{
			name:      "clean",
			inner:     &fakeAgentRunner{},
			wantCode:  0,
			wantState: store.StateClean,
			closeErr:  errors.New("close failed"),
		},
		{
			name:      "unresolved",
			inner:     &fakeAgentRunner{runErr: errors.New("agent crashed")},
			wantCode:  1,
			wantState: store.StateUnresolved,
		},
		{
			name:      "failed",
			inner:     &fakeAgentRunner{},
			pusherErr: errors.New("push failed"),
			wantCode:  1,
			wantState: store.StateFailed,
		},
		{
			name:      "stopped",
			inner:     &fakeStoppingAgentRunner{},
			wantCode:  0,
			wantState: store.StateStopped,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, repoDir := withCLIWorkspace(t)
			withSuccessfulPreflight(t, repoDir)
			runner := &sessionRecordingRunner{inner: tt.inner, closeErr: tt.closeErr}
			withAgentRunner(t, runner)
			withFakeWorktree(t)
			if tt.pusherErr != nil {
				withPusher(t, &fakePusher{err: tt.pusherErr})
			}
			if tt.wantState == store.StateStopped {
				withChangedPaths(t, nil)
			}
			persistCLIReviewIssue(t, repoDir, 1, "feature/review")
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

			if code != tt.wantCode {
				t.Fatalf("expected exit code %d, got %d stderr=%q", tt.wantCode, code, stderr.String())
			}
			runID, _ := journaledRunEvents(t, homeDir, stderr.String())
			run := runFromStore(t, homeDir, runID)
			if run.State != tt.wantState {
				t.Fatalf("expected Run state %q, got %q", tt.wantState, run.State)
			}
			assertRecordedOneSessionForRun(t, runner, runID, len(runner.runSessions))
		})
	}
}

func TestRunResolveRejectsIncompatibleArtifacts(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	persistCLIReviewIssue(t, repoDir, 1, "other-branch")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
	if !strings.Contains(stderr.String(), "Head Repository \"owner/project\", PR Head Branch \"feature/review\"") {
		t.Fatalf("expected incompatible artifact context, got %q", stderr.String())
	}
	assertNoRunDatabase(t, homeDir)
}

func TestRunOperationalCommandRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		contains string
	}{
		{
			name:     "missing pull request with no input",
			args:     []string{"fetch", "--source", "coderabbit", "--no-input"},
			contains: "missing required --pr because --no-input disables Interactive Input",
		},
		{
			name:     "invalid pull request number",
			args:     []string{"fetch", "--source", "coderabbit", "--pr", "abc", "--no-input"},
			contains: "--pr must be a positive integer",
		},
		{
			name:     "unsupported source",
			args:     []string{"fetch", "--source", "other", "--pr", "123", "--no-input"},
			contains: "unsupported Review Source",
		},
		{
			name:     "unsupported agent",
			args:     []string{"resolve", "--pr", "123", "--agent", "other", "--no-input"},
			contains: "unsupported Agent",
		},
		{
			name:     "invalid max rounds",
			args:     []string{"watch", "--source", "coderabbit", "--pr", "123", "--agent", "codex", "--max-rounds", "0", "--no-input"},
			contains: "--max-rounds must be greater than 0",
		},
		{
			name:     "resolve agent console suppression conflicts with interactive input",
			args:     []string{"resolve", "--pr", "123", "--agent", "codex", "--interactive", "--no-agent-console"},
			contains: "--interactive cannot be used with --no-agent-console",
		},
		{
			name:     "watch agent console suppression conflicts with interactive input",
			args:     []string{"watch", "--source", "coderabbit", "--pr", "123", "--agent", "codex", "--interactive", "--no-agent-console"},
			contains: "--interactive cannot be used with --no-agent-console",
		},
		{
			name:     "unknown flag",
			args:     []string{"fetch", "--unknown"},
			contains: "flag provided but not defined",
		},
		{
			name:     "unexpected argument",
			args:     []string{"fetch", "--source", "coderabbit", "--pr", "123", "extra"},
			contains: "unexpected argument",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, _ := withCLIWorkspace(t)
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run(tt.args, &stdout, &stderr)

			if code != 2 {
				t.Fatalf("expected exit code 2, got %d", code)
			}
			if stdout.Len() != 0 {
				t.Fatalf("expected no stdout, got %q", stdout.String())
			}
			if !strings.Contains(stderr.String(), tt.contains) {
				t.Fatalf("expected stderr to contain %q, got %q", tt.contains, stderr.String())
			}
			if !strings.Contains(stderr.String(), "did not create a Run") {
				t.Fatalf("expected no side-effect confirmation, got %q", stderr.String())
			}
			assertNoRunDatabase(t, homeDir)
		})
	}
}

func TestRunNoAgentConsoleRejectsInteractiveCockpit(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "resolve",
			args: []string{"resolve", "--pr", "123", "--agent", "codex", "--no-agent-console", "--no-input"},
		},
		{
			name: "watch",
			args: []string{"watch", "--source", "coderabbit", "--pr", "123", "--agent", "codex", "--no-agent-console", "--no-input"},
		},
		{
			name: "implement",
			args: []string{"implement", "--spec", "0001-widget-flow", "--agent", "codex", "--no-agent-console", "--no-input"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, _ := withCLIWorkspace(t)
			t.Setenv("ROUNDFIX_TUI", "always")
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run(tt.args, &stdout, &stderr)

			if code != exitPreflight {
				t.Fatalf("expected preflight exit %d, got %d stderr=%q", exitPreflight, code, stderr.String())
			}
			if stdout.Len() != 0 {
				t.Fatalf("expected no stdout, got %q", stdout.String())
			}
			if !strings.Contains(stderr.String(), "--no-agent-console cannot be used with the interactive cockpit") {
				t.Fatalf("expected interactive cockpit conflict, got %q", stderr.String())
			}
			assertNoRunDatabase(t, homeDir)
		})
	}
}

func TestRunFetchCollectsMissingPullRequestWithInteractiveInput(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withCurrentPullRequestSuggestion(t, "321")
	var inputReq roundtui.InputRequest
	withInteractiveInput(t, func(_ context.Context, req roundtui.InputRequest) (roundtui.CommandValues, error) {
		inputReq = req
		values := req.Values
		values.PRNumber = req.PRSuggestion.Value
		return values, nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"fetch", "--source", "coderabbit"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected fetch exit 0, got %d", code)
	}
	if inputReq.Command != "fetch" {
		t.Fatalf("expected fetch Interactive Input request, got %#v", inputReq)
	}
	if inputReq.PRSuggestion.Value != "321" || inputReq.PRSuggestion.Source != "current" {
		t.Fatalf("expected current PR suggestion, got %#v", inputReq.PRSuggestion)
	}
	if !strings.Contains(stderr.String(), "Interactive Input collected command parameters.") {
		t.Fatalf("expected Interactive Input confirmation, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Roundfix fetch") || !strings.Contains(stderr.String(), "State: FetchingIssues") {
		t.Fatalf("expected Live Run View output, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
	defaults := readInteractiveDefaults(t, homeDir)
	if defaults.PRNumber != "321" {
		t.Fatalf("expected remembered PR 321, got %#v", defaults)
	}
}

func TestRunFetchReportsBlankInteractivePullRequestWithoutSuggestingInteractiveAgain(t *testing.T) {
	homeDir, _ := withCLIWorkspace(t)
	withCurrentPullRequestSuggestion(t, "")
	withInteractiveInput(t, func(_ context.Context, req roundtui.InputRequest) (roundtui.CommandValues, error) {
		return req.Values, nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"fetch"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "Interactive Input did not collect required --pr") {
		t.Fatalf("expected blank Interactive Input guidance, got %q", stderr.String())
	}
	if strings.Contains(stderr.String(), "use --interactive") {
		t.Fatalf("did not expect guidance to reopen Interactive Input, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "did not create a Run") {
		t.Fatalf("expected no side-effect confirmation, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 0)
}

func TestRunResolveInteractiveInputSuggestsConfiguredAgentAndRememberedPullRequest(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	if err := runStore.RememberInteractiveDefaults(context.Background(), store.InteractiveDefaults{PRNumber: "123", Agent: "claude"}); err != nil {
		t.Fatalf("remember defaults: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}
	withSuccessfulPreflight(t, repoDir)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	withCurrentPullRequestSuggestion(t, "")
	var inputReq roundtui.InputRequest
	withInteractiveInput(t, func(_ context.Context, req roundtui.InputRequest) (roundtui.CommandValues, error) {
		inputReq = req
		values := req.Values
		values.PRNumber = req.PRSuggestion.Value
		return values, nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"resolve", "--interactive"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected resolve exit 0, got %d", code)
	}
	if inputReq.PRSuggestion.Value != "123" || inputReq.PRSuggestion.Source != "remembered" {
		t.Fatalf("expected remembered PR suggestion, got %#v", inputReq.PRSuggestion)
	}
	if inputReq.AgentSuggestion.Value != "codex" || inputReq.AgentSuggestion.Source != "config" {
		t.Fatalf("expected configured Agent suggestion before remembered value, got %#v", inputReq.AgentSuggestion)
	}
	if !strings.Contains(stderr.String(), "State: ResolvingWithAgent") {
		t.Fatalf("expected resolving Live Run View, got %q", stderr.String())
	}
	defaults := readInteractiveDefaults(t, homeDir)
	if defaults.PRNumber != "123" || defaults.Agent != "codex" {
		t.Fatalf("expected remembered PR and configured Agent after run, got %#v", defaults)
	}
}

func TestRunInteractiveInputRunsPreflightBeforeSideEffects(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withFetchReviewItemsFunc(t, func(context.Context, reviewsource.FetchRequest) ([]reviewsource.ReviewItem, error) {
		t.Fatal("fetch must not run after post-input Preflight Validation failure")
		return nil, nil
	})
	withCurrentPullRequestSuggestion(t, "123")
	withInteractiveInput(t, func(_ context.Context, req roundtui.InputRequest) (roundtui.CommandValues, error) {
		values := req.Values
		values.PRNumber = req.PRSuggestion.Value
		return values, nil
	})
	withPreflight(t, func(context.Context, commandRequest, roundconfig.Loaded) (preflight.Result, error) {
		return preflight.Result{}, errors.New("post-input preflight failed")
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"fetch", "--source", "coderabbit"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected post-input preflight exit 2, got %d", code)
	}
	if !strings.Contains(stderr.String(), "post-input preflight failed") {
		t.Fatalf("expected preflight failure, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "did not create a Run") {
		t.Fatalf("expected no side-effect confirmation, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 0)
	_ = repoDir
}

func TestRunOperationalCommandAppliesConfigAndCLIArtifactDirPrecedence(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
	mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), `
defaults:
  artifact_dir: user-artifacts
`)
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), `
defaults:
  artifact_dir: project-artifacts
`)

	var projectStdout bytes.Buffer
	var projectStderr bytes.Buffer
	projectCode := Run([]string{"fetch", "--source", "coderabbit", "--pr", "123", "--no-input"}, &projectStdout, &projectStderr)

	if projectCode != 0 {
		t.Fatalf("expected tracked Fetch Run exit 0, got %d", projectCode)
	}
	if !strings.Contains(projectStdout.String(), filepath.Join(repoDir, "project-artifacts")) {
		t.Fatalf("expected project artifact dir to win over user config, got %q", projectStdout.String())
	}

	var cliStdout bytes.Buffer
	var cliStderr bytes.Buffer
	cliCode := Run([]string{"fetch", "--source", "coderabbit", "--pr", "123", "--artifact-dir", "cli-artifacts", "--no-input"}, &cliStdout, &cliStderr)

	if cliCode != 0 {
		t.Fatalf("expected tracked Fetch Run exit 0, got %d", cliCode)
	}
	if !strings.Contains(cliStdout.String(), filepath.Join(repoDir, "cli-artifacts")) {
		t.Fatalf("expected CLI artifact dir to win over project config, got %q", cliStdout.String())
	}
}

func TestRunFetchRejectsDuplicateActiveRun(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)

	ctx := context.Background()
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	active, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:           store.KindWatch,
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
		BaseRepository: "owner/project",
		PRNumber:       "123",
		GitRoot:        repoDir,
		LocalBranch:    "feature/review",
		HeadSHA:        "abc123",
		ArtifactDir:    filepath.Join(repoDir, ".roundfix"),
	})
	if err != nil {
		t.Fatalf("create active run: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}
	subdir := filepath.Join(repoDir, "nested")
	mustMkdir(t, subdir)
	t.Chdir(subdir)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"fetch", "--source", "coderabbit", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected exit code 2 for duplicate Active Run, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "existing run_id="+active.ID) {
		t.Fatalf("expected existing run id in stderr, got %q", stderr.String())
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
}

func TestRunStopByRunIDRecordsStopRequest(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)

	ctx := context.Background()
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	active, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:           store.KindResolve,
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
		BaseRepository: "owner/project",
		PRNumber:       "123",
		GitRoot:        repoDir,
		LocalBranch:    "feature/review",
		HeadSHA:        "abc123",
		ArtifactDir:    filepath.Join(repoDir, ".roundfix"),
	})
	if err != nil {
		t.Fatalf("create active run: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"stop", active.ID}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected stop exit 0, got %d stderr=%q", code, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
	for _, expected := range []string{active.ID, "Stop Request recorded; the Run stops after the current Work Item settles.", "--force"} {
		if !strings.Contains(stdout.String(), expected) {
			t.Fatalf("expected stop output to contain %q, got %q", expected, stdout.String())
		}
	}
	assertStopRequested(t, homeDir, active.ID)
	assertActiveRun(t, homeDir, "owner/project", "feature/review", active.ID)
}

func TestRunStopByPullRequestRecordsStopRequestForMatchingActiveRun(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withStopPullRequestResolver(t, preflight.PullRequest{
		Number:         "123",
		State:          "OPEN",
		BaseRepository: "owner/project",
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
	})

	ctx := context.Background()
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	active, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:           store.KindWatch,
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
		BaseRepository: "owner/project",
		PRNumber:       "123",
		GitRoot:        repoDir,
		LocalBranch:    "feature/review",
		HeadSHA:        "abc123",
		ArtifactDir:    filepath.Join(repoDir, ".roundfix"),
	})
	if err != nil {
		t.Fatalf("create active run: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"stop", "--pr", "123"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected stop exit 0, got %d stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), active.ID) {
		t.Fatalf("expected stopped run id in stdout, got %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "Stop Request recorded; the Run stops after the current Work Item settles.") {
		t.Fatalf("expected Stop Request report in stdout, got %q", stdout.String())
	}
	assertStopRequested(t, homeDir, active.ID)
	assertActiveRun(t, homeDir, "owner/project", "feature/review", active.ID)
}

func TestRunStopBySpecRecordsStopRequestForMatchingActiveRun(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)

	ctx := context.Background()
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	request := store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     repoDir,
		LocalBranch: "ma/implement-spec",
		SpecSlug:    "0001-widget-flow",
	}
	active, err := runStore.CreateRun(ctx, request)
	if err != nil {
		t.Fatalf("create active implement run: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"stop", "--spec", "0001-widget-flow"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected stop exit 0, got %d stderr=%q", code, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
	if !strings.Contains(stdout.String(), active.ID) ||
		!strings.Contains(stdout.String(), "Stop Request recorded; the Run stops after the current Work Item settles.") ||
		!strings.Contains(stdout.String(), "--force") {
		t.Fatalf("expected Stop Request report for implement run in stdout, got %q", stdout.String())
	}
	runStore, err = store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("reopen store: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close reopened store: %v", err)
		}
	}()
	current, found, err := runStore.Run(ctx, active.ID)
	if err != nil || !found {
		t.Fatalf("lookup Run: found=%v err=%v", found, err)
	}
	if current.State != store.StateActive {
		t.Fatalf("expected Run to stay Active until the engine settles, got %q", current.State)
	}
	requested, err := runStore.StopRequested(ctx, active.ID)
	if err != nil {
		t.Fatalf("read Stop Request flag: %v", err)
	}
	if !requested {
		t.Fatal("expected Stop Request flag recorded")
	}
	if _, err := runStore.CreateRun(ctx, request); err == nil {
		t.Fatal("expected Stop Request to keep the Active Run lock until engine completion")
	}
}

func TestRunStopForceCancelsImplementRunAndReleasesLocks(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	active, request := createActiveImplementRunForStop(t, homeDir, repoDir, "0001-widget-flow", "codex")
	var calls []struct {
		runtime agent.RuntimeSpec
		session agent.SessionRef
	}
	withStopAgentSessionCanceler(t, func(_ context.Context, runtime agent.RuntimeSpec, session agent.SessionRef) error {
		calls = append(calls, struct {
			runtime agent.RuntimeSpec
			session agent.SessionRef
		}{runtime: runtime, session: session})
		return nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"stop", "--force", "--spec", "0001-widget-flow"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected force stop exit 0, got %d stderr=%q", code, stderr.String())
	}
	if want := "roundfix: closed session roundfix-" + active.ID; !strings.Contains(stderr.String(), want) {
		t.Fatalf("expected force close report %q, got %q", want, stderr.String())
	}
	if len(calls) != 1 {
		t.Fatalf("expected one Agent Session cancel, got %#v", calls)
	}
	if calls[0].runtime.ID != "codex" || calls[0].session.Name != "roundfix-"+active.ID || calls[0].session.WorkDir != repoDir {
		t.Fatalf("unexpected cancel call: %#v", calls[0])
	}
	if !strings.Contains(stdout.String(), "Roundfix Run force-stopped") || !strings.Contains(stdout.String(), active.ID) {
		t.Fatalf("expected force stop report with Run ID, got %q", stdout.String())
	}
	assertRunState(t, homeDir, active.ID, store.StateStopped)

	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store: %v", err)
		}
	}()
	next, err := runStore.CreateRun(context.Background(), request)
	if err != nil {
		t.Fatalf("expected force stop to release Spec lock, got %v", err)
	}
	if next.ID == active.ID {
		t.Fatal("expected a new Run after force stop")
	}
}

func TestRunStopForceCancelsListsAndClosesRunAndTaskSessions(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{
		id:     "task_01",
		title:  "Stopped task",
		status: string(spec.StatusPending),
	}})
	location := configureSettleWorktreeLocation(t, repoDir, filepath.Join(homeDir, "worktrees"))
	active, _, taskRef := createImplementRunWorktreeFixture(t, homeDir, repoDir, location, implementTestSlug, "task_01", "")
	calls := []string{}
	withStopAgentSessionCanceler(t, func(_ context.Context, runtime agent.RuntimeSpec, session agent.SessionRef) error {
		calls = append(calls, "cancel "+runtime.ID+" "+session.Name+" "+session.WorkDir)
		return nil
	})
	withRoundfixSessionLister(t, func(_ context.Context, runtime agent.RuntimeSpec, workDir string) ([]agent.RoundfixSession, error) {
		calls = append(calls, "list "+runtime.ID+" "+workDir)
		return []agent.RoundfixSession{
			{Name: "roundfix-" + active.ID + "-task_01", RunID: active.ID, TaskID: "task_01"},
			{Name: "roundfix-other-run-task_01", RunID: "other-run", TaskID: "task_01"},
		}, nil
	})
	withStopAgentSessionCloser(t, func(_ context.Context, runtime agent.RuntimeSpec, session agent.SessionRef) error {
		calls = append(calls, "close "+runtime.ID+" "+session.Name+" "+session.WorkDir)
		return nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"stop", "--force", active.ID}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected force stop exit 0, got %d stderr=%q", code, stderr.String())
	}
	wantCalls := []string{
		"cancel codex roundfix-" + active.ID + " " + active.WorkDir,
		"list codex " + active.WorkDir,
		"close codex roundfix-" + active.ID + " " + active.WorkDir,
		"close codex roundfix-" + active.ID + "-task_01 " + taskRef.Path,
	}
	if !reflect.DeepEqual(calls, wantCalls) {
		t.Fatalf("unexpected session invocations\nwant: %#v\ngot:  %#v", wantCalls, calls)
	}
	for _, want := range []string{
		"roundfix: closed session roundfix-" + active.ID,
		"roundfix: closed session roundfix-" + active.ID + "-task_01",
	} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("expected stderr to contain %q, got %q", want, stderr.String())
		}
	}
	if strings.Contains(stderr.String(), "other-run") {
		t.Fatalf("foreign Run session must not be reported or closed, got %q", stderr.String())
	}
	if !strings.Contains(stdout.String(), "Roundfix Run force-stopped") {
		t.Fatalf("expected force stop report, got %q", stdout.String())
	}
	assertRunState(t, homeDir, active.ID, store.StateStopped)
}

func TestRunStopForceCloseFailureStillCompletesStopped(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	active, request := createActiveImplementRunForStop(t, homeDir, repoDir, "0001-widget-flow", "codex")
	withStopAgentSessionCanceler(t, func(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
		return nil
	})
	withStopAgentSessionCloser(t, func(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
		return errors.New("close denied")
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"stop", "--force", active.ID}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected force stop exit 0, got %d stderr=%q", code, stderr.String())
	}
	if want := "roundfix: could not close session roundfix-" + active.ID + ": close denied"; !strings.Contains(stderr.String(), want) {
		t.Fatalf("expected close failure note %q, got %q", want, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Roundfix Run force-stopped") {
		t.Fatalf("expected force stop report, got %q", stdout.String())
	}
	assertRunState(t, homeDir, active.ID, store.StateStopped)

	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	_, err = runStore.CreateRun(context.Background(), request)
	closeErr := runStore.Close()
	if err != nil {
		t.Fatalf("expected force stop to release Spec lock after close failure, got %v", err)
	}
	if closeErr != nil {
		t.Fatalf("close store: %v", closeErr)
	}
}

func TestRunStopForceReapsEmptyRunAndTaskWorktrees(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{
		id:     "task_01",
		title:  "Stopped task",
		status: string(spec.StatusPending),
	}})
	location := configureSettleWorktreeLocation(t, repoDir, filepath.Join(homeDir, "worktrees"))
	active, _, taskRef := createImplementRunWorktreeFixture(t, homeDir, repoDir, location, implementTestSlug, "task_01", "")
	withStopAgentSessionCanceler(t, func(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
		return nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"stop", "--force", active.ID}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected force stop exit 0, got %d stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Roundfix Run force-stopped") {
		t.Fatalf("expected force stop report, got %q", stdout.String())
	}
	for _, expected := range []string{
		"roundfix: reaped terminal Worktree path=" + active.WorkDir + " branch=" + runworktree.BranchName(active.ID),
		"roundfix: reaped terminal Worktree path=" + taskRef.Path + " branch=" + taskRef.Branch,
	} {
		if !strings.Contains(stderr.String(), expected) {
			t.Fatalf("expected stderr to contain %q, got %q", expected, stderr.String())
		}
	}
	assertRunState(t, homeDir, active.ID, store.StateStopped)
	assertRunWorktreeRemoved(t, active.WorkDir)
	assertRunBranchRemoved(t, repoDir, runworktree.BranchName(active.ID))
	assertRunWorktreeRemoved(t, taskRef.Path)
	assertRunBranchRemoved(t, repoDir, taskRef.Branch)
}

func TestRunStopForceKeepsRunWorktreeWithCommits(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{
		id:     "task_01",
		title:  "Stopped task",
		status: string(spec.StatusPending),
	}})
	location := configureSettleWorktreeLocation(t, repoDir, filepath.Join(homeDir, "worktrees"))
	active, runRef, _ := createImplementRunWorktreeFixture(t, homeDir, repoDir, location, implementTestSlug, "", "")
	if err := os.WriteFile(filepath.Join(runRef.Path, "kept.txt"), []byte("settled work\n"), 0o644); err != nil {
		t.Fatalf("write kept work: %v", err)
	}
	gitImplement(t, runRef.Path, "add", "kept.txt")
	gitImplement(t, runRef.Path, "commit", "-m", "settled work")
	withStopAgentSessionCanceler(t, func(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
		return nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"stop", "--force", active.ID}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected force stop exit 0, got %d stderr=%q", code, stderr.String())
	}
	if strings.Contains(stderr.String(), "reaped terminal Worktree") {
		t.Fatalf("expected no debris reap report for a Run Branch with commits, got %q", stderr.String())
	}
	assertRunState(t, homeDir, active.ID, store.StateStopped)
	assertRunWorktreeExists(t, active.WorkDir)
	assertRunBranchExists(t, repoDir, runworktree.BranchName(active.ID))
}

func TestRunStopForceReportsCancelFailuresButCompletes(t *testing.T) {
	tests := []struct {
		name          string
		agentID       string
		cancelErr     error
		wantStderr    []string
		wantCanceller bool
	}{
		{
			name:          "missing acpx",
			agentID:       "codex",
			cancelErr:     errors.New("acpx not found"),
			wantStderr:    []string{"Agent Session cancel failed", "acpx not found"},
			wantCanceller: true,
		},
		{
			name:          "unknown runtime",
			agentID:       "gemini",
			wantStderr:    []string{"Agent Session cancel skipped", `unsupported Agent "gemini"`},
			wantCanceller: false,
		},
		{
			name:          "no recorded runtime",
			agentID:       "",
			wantStderr:    []string{"Agent Session cancel skipped", "no Agent runtime recorded"},
			wantCanceller: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, repoDir := withCLIWorkspace(t)
			active, request := createActiveImplementRunForStop(t, homeDir, repoDir, "0001-widget-flow", tt.agentID)
			cancelCalls := 0
			withStopAgentSessionCanceler(t, func(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
				cancelCalls++
				return tt.cancelErr
			})
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run([]string{"stop", "--force", active.ID}, &stdout, &stderr)

			if code != 0 {
				t.Fatalf("expected force stop exit 0, got %d stderr=%q", code, stderr.String())
			}
			for _, want := range tt.wantStderr {
				if !strings.Contains(stderr.String(), want) {
					t.Fatalf("expected stderr to contain %q, got %q", want, stderr.String())
				}
			}
			if tt.wantCanceller && cancelCalls != 1 {
				t.Fatalf("expected one cancel attempt, got %d", cancelCalls)
			}
			if !tt.wantCanceller && cancelCalls != 0 {
				t.Fatalf("expected no cancel attempt, got %d", cancelCalls)
			}
			if !strings.Contains(stdout.String(), "Roundfix Run force-stopped") {
				t.Fatalf("expected force stop report, got %q", stdout.String())
			}
			assertRunState(t, homeDir, active.ID, store.StateStopped)
			runStore, err := store.Open(context.Background(), homeDir)
			if err != nil {
				t.Fatalf("open store: %v", err)
			}
			_, err = runStore.CreateRun(context.Background(), request)
			closeErr := runStore.Close()
			if err != nil {
				t.Fatalf("expected force stop to release Spec lock after cancel failure, got %v", err)
			}
			if closeErr != nil {
				t.Fatalf("close store: %v", closeErr)
			}
		})
	}
}

func TestRunStopGracefulThenForceCompletesImmediately(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	active, request := createActiveImplementRunForStop(t, homeDir, repoDir, "0001-widget-flow", "codex")
	var gracefulStdout bytes.Buffer
	var gracefulStderr bytes.Buffer

	gracefulCode := Run([]string{"stop", "--spec", "0001-widget-flow"}, &gracefulStdout, &gracefulStderr)

	if gracefulCode != 0 {
		t.Fatalf("expected graceful stop exit 0, got %d stderr=%q", gracefulCode, gracefulStderr.String())
	}
	if !strings.Contains(gracefulStdout.String(), "Stop Request recorded") {
		t.Fatalf("expected graceful Stop Request report, got %q", gracefulStdout.String())
	}
	assertStopRequested(t, homeDir, active.ID)
	cancelCalls := 0
	withStopAgentSessionCanceler(t, func(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
		cancelCalls++
		return nil
	})
	var forceStdout bytes.Buffer
	var forceStderr bytes.Buffer

	forceCode := Run([]string{"stop", "--force", "--spec", "0001-widget-flow"}, &forceStdout, &forceStderr)

	if forceCode != 0 {
		t.Fatalf("expected force stop exit 0, got %d stderr=%q", forceCode, forceStderr.String())
	}
	if want := "roundfix: closed session roundfix-" + active.ID; !strings.Contains(forceStderr.String(), want) {
		t.Fatalf("expected force close report %q, got %q", want, forceStderr.String())
	}
	if cancelCalls != 1 {
		t.Fatalf("expected one cancel after graceful request, got %d", cancelCalls)
	}
	if !strings.Contains(forceStdout.String(), "Roundfix Run force-stopped") {
		t.Fatalf("expected force stop report, got %q", forceStdout.String())
	}
	assertRunState(t, homeDir, active.ID, store.StateStopped)

	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store: %v", err)
		}
	}()
	if _, err := runStore.CreateRun(context.Background(), request); err != nil {
		t.Fatalf("expected force after graceful request to release Spec lock, got %v", err)
	}
}

func TestRunStopSpecSelectorRejectsOtherSelectors(t *testing.T) {
	withCLIWorkspace(t)
	tests := []struct {
		name string
		args []string
	}{
		{name: "with pull request", args: []string{"stop", "--spec", "0001-widget-flow", "--pr", "123"}},
		{name: "with run id", args: []string{"stop", "--spec", "0001-widget-flow", "run_existing"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run(tt.args, &stdout, &stderr)

			if code != 2 {
				t.Fatalf("expected stop validation exit 2, got %d", code)
			}
			if stdout.Len() != 0 {
				t.Fatalf("expected no stdout, got %q", stdout.String())
			}
			if !strings.Contains(stderr.String(), "cannot be combined") {
				t.Fatalf("expected mutual-exclusion message, got %q", stderr.String())
			}
		})
	}
}

func TestRunStopBySpecReportsMissingActiveRun(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"stop", "--spec", "0002-missing"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected stop validation exit 2, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	want := fmt.Sprintf("no Active Run exists for repository %q and Spec %q", repoDir, "0002-missing")
	if !strings.Contains(stderr.String(), want) {
		t.Fatalf("expected missing spec target %q, got %q", want, stderr.String())
	}
}

func TestRunStopHelpListsSpecSelector(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"stop", "--help"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected stop help exit 0, got %d", code)
	}
	for _, want := range []string{"roundfix stop --spec <slug>", "--spec", "--force", "dead or runaway", "graceful"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("expected stop help to contain %q, got %q", want, stdout.String())
		}
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunStopRejectsAlreadyTerminalRun(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)

	ctx := context.Background()
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	active, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:           store.KindResolve,
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
		BaseRepository: "owner/project",
		PRNumber:       "123",
		GitRoot:        repoDir,
		LocalBranch:    "feature/review",
		HeadSHA:        "abc123",
		ArtifactDir:    filepath.Join(repoDir, ".roundfix"),
	})
	if err != nil {
		t.Fatalf("create active run: %v", err)
	}
	if _, err := runStore.CompleteRun(ctx, active.ID, store.StateClean); err != nil {
		t.Fatalf("complete run: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := Run([]string{"stop", active.ID}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected stop validation exit 2, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "cannot record Stop Request for terminal Run") ||
		!strings.Contains(stderr.String(), "Clean") {
		t.Fatalf("expected named terminal Stop Request refusal, got %q", stderr.String())
	}
}

func TestRunPreflightFailureLeavesBufferOutputPlainByDefault(t *testing.T) {
	withCLIWorkspace(t)
	withPreflight(t, func(context.Context, commandRequest, roundconfig.Loaded) (preflight.Result, error) {
		return preflight.Result{}, errors.New("plain preflight failure")
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"fetch", "--source", "coderabbit", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
	if strings.Contains(stderr.String(), "\x1b[") {
		t.Fatalf("did not expect ANSI color in buffer output, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "Preflight failed") || !strings.Contains(stderr.String(), "plain preflight failure") {
		t.Fatalf("expected formatted preflight failure, got %q", stderr.String())
	}
}

func TestRunPreflightFailureColorsOutputWhenForced(t *testing.T) {
	t.Setenv("ROUNDFIX_COLOR", "always")
	withCLIWorkspace(t)
	withPreflight(t, func(context.Context, commandRequest, roundconfig.Loaded) (preflight.Result, error) {
		return preflight.Result{}, errors.New("colored preflight failure")
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"fetch", "--source", "coderabbit", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
	if !strings.Contains(stderr.String(), "\x1b[31m") {
		t.Fatalf("expected forced ANSI color in stderr, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "colored preflight failure") {
		t.Fatalf("expected original error text, got %q", stderr.String())
	}
}

func TestRunOperationalCommandRejectsInvalidConfigAndArtifactDirectory(t *testing.T) {
	t.Run("invalid user config YAML", func(t *testing.T) {
		homeDir, _ := withCLIWorkspace(t)
		mustMkdir(t, filepath.Join(homeDir, ".roundfix"))
		mustWrite(t, filepath.Join(homeDir, ".roundfix", "config.yml"), "defaults:\n  agent: [")

		var stdout bytes.Buffer
		var stderr bytes.Buffer
		code := Run([]string{"fetch", "--source", "coderabbit", "--pr", "123", "--no-input"}, &stdout, &stderr)

		if code != 2 {
			t.Fatalf("expected exit code 2, got %d", code)
		}
		if stdout.Len() != 0 {
			t.Fatalf("expected no stdout, got %q", stdout.String())
		}
		if !strings.Contains(stderr.String(), "parse config") {
			t.Fatalf("expected config parse error, got %q", stderr.String())
		}
		if !strings.Contains(stderr.String(), "did not create a Run") {
			t.Fatalf("expected no side-effect confirmation, got %q", stderr.String())
		}
	})

	t.Run("artifact path is a file", func(t *testing.T) {
		_, repoDir := withCLIWorkspace(t)
		artifactFile := filepath.Join(repoDir, "artifact-file")
		mustWrite(t, artifactFile, "not a directory")

		var stdout bytes.Buffer
		var stderr bytes.Buffer
		code := Run([]string{"fetch", "--source", "coderabbit", "--pr", "123", "--artifact-dir", artifactFile, "--no-input"}, &stdout, &stderr)

		if code != 2 {
			t.Fatalf("expected exit code 2, got %d", code)
		}
		if stdout.Len() != 0 {
			t.Fatalf("expected no stdout, got %q", stdout.String())
		}
		if !strings.Contains(stderr.String(), "is not a directory") {
			t.Fatalf("expected artifact dir error, got %q", stderr.String())
		}
		if !strings.Contains(stderr.String(), "did not create a Run") {
			t.Fatalf("expected no side-effect confirmation, got %q", stderr.String())
		}
	})
}

type setupFakeDeps struct {
	homeDir           string
	gitRoot           string
	config            roundconfig.Config
	userConfigPath    string
	projectConfigPath string
	acpxConfigPath    string
	nodeVersion       string
	nodeErr           error
	acpxVersion       string
	acpxErr           error
	probeErr          error
	paths             map[string]string
	files             map[string]string
	installCalls      []string
	initScopes        []string
	acpxInitCalls     int
	writeCalls        []string
	probeRequests     []agent.ProbeRequest
	prompts           []string
	confirm           func(context.Context, io.Writer, string) (bool, error)
}

func newSetupFakeDeps() *setupFakeDeps {
	homeDir := "/home/roundfix-test"
	gitRoot := "/repo/project"
	return &setupFakeDeps{
		homeDir:           homeDir,
		gitRoot:           gitRoot,
		config:            roundconfig.Builtin(),
		userConfigPath:    filepath.Join(homeDir, ".roundfix", "config.yml"),
		projectConfigPath: filepath.Join(gitRoot, ".roundfixrc.yml"),
		acpxConfigPath:    filepath.Join(homeDir, ".acpx", "config.json"),
		nodeVersion:       "v25.6.1",
		acpxVersion:       agent.PinnedACPXVersion,
		paths:             map[string]string{},
		files:             map[string]string{},
	}
}

func withSetupFakeDeps(t *testing.T, fake *setupFakeDeps) {
	t.Helper()
	old := setupDeps
	setupDeps = setupDependencies{
		loadConfig: func(roundconfig.LoadOptions) (roundconfig.Loaded, error) {
			return roundconfig.Loaded{
				Config:            fake.config,
				GitRoot:           fake.gitRoot,
				HomeDir:           fake.homeDir,
				UserConfigPath:    fake.userConfigPath,
				ProjectConfigPath: fake.projectConfigPath,
			}, nil
		},
		nodeVersion: func(context.Context) (string, error) {
			return fake.nodeVersion, fake.nodeErr
		},
		acpxVersion: func(context.Context) (string, error) {
			return fake.acpxVersion, fake.acpxErr
		},
		installACPX: func(context.Context) error {
			fake.installCalls = append(fake.installCalls, "npm install -g acpx@"+agent.PinnedACPXVersion)
			return nil
		},
		probeAgent: func(_ context.Context, req agent.ProbeRequest) error {
			fake.probeRequests = append(fake.probeRequests, req)
			return fake.probeErr
		},
		lookPath: func(command string) (string, error) {
			path := fake.paths[command]
			if path == "" {
				return "", os.ErrNotExist
			}
			return path, nil
		},
		exists: func(path string) (bool, error) {
			_, ok := fake.files[path]
			return ok, nil
		},
		readFile: func(path string) ([]byte, error) {
			content, ok := fake.files[path]
			if !ok {
				return nil, os.ErrNotExist
			}
			return []byte(content), nil
		},
		writeFile: func(path string, content []byte) error {
			fake.files[path] = string(content)
			fake.writeCalls = append(fake.writeCalls, path)
			return nil
		},
		mkdirAll: func(string) error {
			return nil
		},
		initACPXConfig: func(context.Context, string) error {
			fake.acpxInitCalls++
			fake.files[fake.acpxConfigPath] = "{\n}\n"
			return nil
		},
		initConfig: func(_ context.Context, opts roundconfig.InitOptions) (roundconfig.InitResult, error) {
			fake.initScopes = append(fake.initScopes, opts.Scope)
			path := fake.userConfigPath
			if opts.Scope == roundconfig.InitScopeProject {
				path = fake.projectConfigPath
			}
			fake.files[path] = roundconfig.DefaultConfigYAML()
			return roundconfig.InitResult{Scope: opts.Scope, Path: path}, nil
		},
		confirm: func(ctx context.Context, stderr io.Writer, prompt string) (bool, error) {
			fake.prompts = append(fake.prompts, prompt)
			if fake.confirm != nil {
				return fake.confirm(ctx, stderr, prompt)
			}
			return true, nil
		},
	}
	t.Cleanup(func() {
		setupDeps = old
	})
}

func assertSetupLineOrder(t *testing.T, output string, lines []string) {
	t.Helper()
	last := -1
	for _, line := range lines {
		index := strings.Index(output, line)
		if index < 0 {
			t.Fatalf("expected setup output to contain %q, got %q", line, output)
		}
		if index < last {
			t.Fatalf("expected setup output line %q after previous lines, got %q", line, output)
		}
		last = index
	}
}

type recordingOutcomeNotifier struct {
	mu       sync.Mutex
	err      error
	outcomes []roundnotify.Outcome
}

func (notifier *recordingOutcomeNotifier) Notify(ctx context.Context, outcome roundnotify.Outcome) error {
	if ctx == nil {
		return errors.New("notification context is required")
	}
	notifier.mu.Lock()
	defer notifier.mu.Unlock()
	notifier.outcomes = append(notifier.outcomes, outcome)
	return notifier.err
}

func (notifier *recordingOutcomeNotifier) recorded() []roundnotify.Outcome {
	notifier.mu.Lock()
	defer notifier.mu.Unlock()
	return append([]roundnotify.Outcome(nil), notifier.outcomes...)
}

type deadlineRecordingOutcomeNotifier struct {
	hadDeadline bool
	deadline    time.Time
	wasCanceled bool
	outcome     roundnotify.Outcome
}

func (notifier *deadlineRecordingOutcomeNotifier) Notify(ctx context.Context, outcome roundnotify.Outcome) error {
	deadline, ok := ctx.Deadline()
	notifier.hadDeadline = ok
	notifier.deadline = deadline
	select {
	case <-ctx.Done():
		notifier.wasCanceled = true
	default:
	}
	notifier.outcome = outcome
	return nil
}

func withOutcomeNotifier(t *testing.T, notifier roundnotify.Notifier) {
	t.Helper()
	withOutcomeNotifierFactory(t, func(roundconfig.Config) roundnotify.Notifier {
		return notifier
	})
}

func withOutcomeNotifierFactory(t *testing.T, factory func(roundconfig.Config) roundnotify.Notifier) {
	t.Helper()
	old := newOutcomeNotifier
	newOutcomeNotifier = factory
	t.Cleanup(func() {
		newOutcomeNotifier = old
	})
}

func assertRecordedOutcomes(t *testing.T, notifier *recordingOutcomeNotifier, want []roundnotify.Outcome) {
	t.Helper()
	got := notifier.recorded()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("recorded notification outcomes mismatch\nwant: %#v\ngot:  %#v", want, got)
	}
}

func runCleanResolveForOutcomeNotification(t *testing.T, notifyErr error) (string, string, int, string) {
	t.Helper()
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	notifier := &recordingOutcomeNotifier{err: notifyErr}
	withOutcomeNotifier(t, notifier)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--agent", "codex", "--round", "all", "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected clean resolve exit code 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if len(notifier.recorded()) != 1 {
		t.Fatalf("expected one notification attempt, got %#v", notifier.recorded())
	}
	return stdout.String(), stderr.String(), code, homeDir
}

func runCleanResolveForCleanup(t *testing.T, cleanupErr error) (string, string, int, string, string) {
	t.Helper()
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	cleanupPath := ""
	if cleanupErr != nil {
		oldCleanup := cleanupCleanRunWorktree
		cleanupCleanRunWorktree = func(_ context.Context, ref runworktree.Ref) error {
			cleanupPath = ref.Path
			return cleanupErr
		}
		t.Cleanup(func() {
			cleanupCleanRunWorktree = oldCleanup
		})
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--agent", "codex", "--round", "all", "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected clean resolve exit code 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if cleanupErr != nil && strings.TrimSpace(cleanupPath) == "" {
		t.Fatal("expected cleanup failure to capture the Run Worktree path")
	}
	return stdout.String(), stderr.String(), code, homeDir, cleanupPath
}

func runCleanWatchForCleanup(t *testing.T, cleanupErr error) (string, string, int, string, string) {
	t.Helper()
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	cleanupPath := ""
	if cleanupErr != nil {
		oldCleanup := cleanupCleanRunWorktree
		cleanupCleanRunWorktree = func(_ context.Context, ref runworktree.Ref) error {
			cleanupPath = ref.Path
			return cleanupErr
		}
		t.Cleanup(func() {
			cleanupCleanRunWorktree = oldCleanup
		})
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"watch", "--source", "coderabbit", "--pr", "123", "--agent", "codex", "--until-clean", "--max-rounds", "1", "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected clean watch exit code 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if cleanupErr != nil && strings.TrimSpace(cleanupPath) == "" {
		t.Fatal("expected cleanup failure to capture the Run Worktree path")
	}
	return stdout.String(), stderr.String(), code, homeDir, cleanupPath
}

func withCLIWorkspace(t *testing.T) (string, string) {
	t.Helper()
	homeDir := t.TempDir()
	repoDir := t.TempDir()
	mustMkdir(t, filepath.Join(repoDir, ".git"))
	t.Setenv("HOME", homeDir)
	t.Chdir(repoDir)
	return homeDir, repoDir
}

type runListSeed struct {
	kind        string
	state       string
	gitRoot     string
	branch      string
	prNumber    string
	specSlug    string
	createdAt   time.Time
	completedAt time.Time
}

func seedRunsForList(t *testing.T, homeDir string, seeds []runListSeed) []store.Run {
	t.Helper()
	ctx := context.Background()
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open run store: %v", err)
	}
	runs := make([]store.Run, 0, len(seeds))
	for _, seed := range seeds {
		run, err := runStore.CreateRun(ctx, createRunRequestForList(seed))
		if err != nil {
			t.Fatalf("create listed Run: %v", err)
		}
		switch {
		case seed.state == "" || seed.state == store.StateActive:
		case store.IsTerminalState(seed.state):
			run, err = runStore.CompleteRun(ctx, run.ID, seed.state)
			if err != nil {
				t.Fatalf("complete listed Run as %s: %v", seed.state, err)
			}
		default:
			if err := runStore.UpdateRunState(ctx, run.ID, seed.state); err != nil {
				t.Fatalf("update listed Run state to %s: %v", seed.state, err)
			}
			var found bool
			run, found, err = runStore.Run(ctx, run.ID)
			if err != nil {
				t.Fatalf("read updated listed Run: %v", err)
			}
			if !found {
				t.Fatalf("expected listed Run %s to exist", run.ID)
			}
		}
		runs = append(runs, run)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close run store: %v", err)
	}
	for index, seed := range seeds {
		setListedRunCreatedAt(t, homeDir, runs[index].ID, seed.createdAt)
		if !seed.completedAt.IsZero() {
			setListedRunCompletedAt(t, homeDir, runs[index].ID, seed.completedAt)
		}
	}
	return runs
}

func createRunRequestForList(seed runListSeed) store.CreateRunRequest {
	req := store.CreateRunRequest{
		Kind:            seed.kind,
		GitRoot:         seed.gitRoot,
		LocalBranch:     seed.branch,
		HeadSHA:         "abc123",
		ArtifactDir:     filepath.Join(seed.gitRoot, ".roundfix"),
		Agent:           "codex",
		Model:           "gpt-5.5",
		ReasoningEffort: "xhigh",
	}
	if req.Kind == store.KindImplement {
		req.SpecSlug = seed.specSlug
		return req
	}
	req.HeadRepository = "owner/project"
	req.HeadBranch = seed.branch
	req.BaseRepository = "owner/project"
	req.PRNumber = seed.prNumber
	return req
}

func setListedRunCreatedAt(t *testing.T, homeDir string, runID string, createdAt time.Time) {
	t.Helper()
	ctx := t.Context()
	db, err := sql.Open("sqlite", "file:"+store.DatabasePath(homeDir)+"?_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatalf("open raw Run Database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("close raw Run Database: %v", err)
		}
	}()
	result, err := db.ExecContext(
		ctx,
		`UPDATE runs SET created_at = ?, updated_at = ? WHERE id = ?`,
		createdAt.UTC().Format(time.RFC3339Nano),
		createdAt.UTC().Format(time.RFC3339Nano),
		runID,
	)
	if err != nil {
		t.Fatalf("set listed Run creation time: %v", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("read listed Run timestamp update result: %v", err)
	}
	if affected != 1 {
		t.Fatalf("expected to update one listed Run timestamp, updated %d", affected)
	}
}

func setListedRunCompletedAt(t *testing.T, homeDir string, runID string, completedAt time.Time) {
	t.Helper()
	ctx := t.Context()
	db, err := sql.Open("sqlite", "file:"+store.DatabasePath(homeDir)+"?_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatalf("open raw Run Database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("close raw Run Database: %v", err)
		}
	}()
	result, err := db.ExecContext(
		ctx,
		`UPDATE runs SET completed_at = ? WHERE id = ?`,
		completedAt.UTC().Format(time.RFC3339Nano),
		runID,
	)
	if err != nil {
		t.Fatalf("set listed Run completion time: %v", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("read listed Run completion update result: %v", err)
	}
	if affected != 1 {
		t.Fatalf("expected to update one listed Run completion time, updated %d", affected)
	}
}

func withRunsListNow(t *testing.T, now time.Time) {
	t.Helper()
	old := runsListNow
	runsListNow = func() time.Time { return now }
	t.Cleanup(func() {
		runsListNow = old
	})
}

func withRunsInteractiveInput(t *testing.T, available bool) {
	t.Helper()
	old := runsInteractiveInputAvailable
	runsInteractiveInputAvailable = func() bool {
		return available
	}
	t.Cleanup(func() {
		runsInteractiveInputAvailable = old
	})
}

func withReviewGitWorkspace(t *testing.T) (string, string) {
	t.Helper()
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	repoDir := t.TempDir()
	gitImplement(t, repoDir, "init", "--initial-branch=main")
	gitImplement(t, repoDir, "config", "user.name", "Roundfix Test")
	gitImplement(t, repoDir, "config", "user.email", "roundfix-test@example.com")
	gitImplement(t, repoDir, "config", "commit.gpgsign", "false")
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), `
defaults:
  verification: test -f agent.txt
watch:
  poll_interval: 1s
  review_timeout: 5s
  quiet_period: 1ms
`)
	mustWrite(t, filepath.Join(repoDir, "README.md"), "seed\n")
	gitImplement(t, repoDir, "add", "-A")
	gitImplement(t, repoDir, "commit", "-m", "seed review repo")
	gitImplement(t, repoDir, "checkout", "-b", "feature/review")
	resolved, err := filepath.EvalSymlinks(repoDir)
	if err != nil {
		t.Fatalf("resolve review repo dir: %v", err)
	}
	t.Chdir(resolved)
	return homeDir, resolved
}

func withRealReviewPreflight(t *testing.T, repoDir string, pushEnabled bool) {
	t.Helper()
	withPreflight(t, func(_ context.Context, req commandRequest, _ roundconfig.Loaded) (preflight.Result, error) {
		if req.pr == "" {
			return preflight.Result{}, errors.New("missing pr in test preflight")
		}
		head := strings.TrimSpace(gitImplementOutput(t, repoDir, "rev-parse", "HEAD"))
		return preflight.Result{
			Git: preflight.GitState{
				Root:           repoDir,
				Branch:         "feature/review",
				HEAD:           head,
				UpstreamRemote: "origin",
				UpstreamBranch: "feature/review",
			},
			PullRequest: preflight.PullRequest{
				Number:         req.pr,
				State:          "OPEN",
				BaseRepository: "owner/project",
				HeadBranch:     "feature/review",
				HeadRepository: "owner/project",
			},
			PushPlan: preflight.PushPlan{
				Enabled: pushEnabled,
				Remote:  "origin",
				Branch:  "feature/review",
				Command: []string{"git", "push", "origin", "HEAD:feature/review"},
			},
		}, nil
	})
}

func mustMkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
}

func mustWrite(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func withSuccessfulPreflight(t *testing.T, repoDir string) {
	t.Helper()
	withAgentRunner(t, &fakeAgentRunner{})
	withFakeReviewRunWorktrees(t)
	withFakeWorktree(t)
	withVerifier(t, &fakeVerifier{})
	withCommitter(t, &fakeCommitter{})
	withFakeWorktree(t)
	withSourceResolver(t, &fakeSourceResolver{})
	withPusher(t, &fakePusher{})
	clock := &fakeWatchClock{now: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC)}
	withWatchTiming(t, clock, &fakeWatchSleeper{clock: clock})
	withWatchStatus(t, (&fakeWatchStatus{
		statuses: []reviewsource.WatchStatus{{State: watch.StatusSettled}},
	}).Status)
	withWatchHeadCheck(t, (&fakeWatchHeadCheck{
		states: []watch.HeadCheckState{watch.CheckSuccess},
	}).Check)
	withWatchHeadSHA(t, func(context.Context, string) (string, error) {
		return "abc123", nil
	})
	withReviewSpecGitRunner(t, &fakeReviewSpecGitRunner{})
	withFetchReviewItems(t, []reviewsource.ReviewItem{
		{
			Title:                   "major: handle test issue",
			File:                    "internal/test.go",
			Line:                    12,
			Severity:                "major",
			Author:                  "coderabbitai[bot]",
			Body:                    "Keep this reviewer text literal: `$(rm -rf /)`.",
			SourceRef:               "thread:PRRT_test,comment:PRRC_test",
			ReviewHash:              "review-hash",
			SourceReviewID:          "9001",
			SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC),
		},
	})
	withPreflight(t, func(_ context.Context, req commandRequest, _ roundconfig.Loaded) (preflight.Result, error) {
		if req.pr == "" {
			return preflight.Result{}, errors.New("missing pr in test preflight")
		}
		return preflight.Result{
			Git: preflight.GitState{
				Root:            repoDir,
				Branch:          "feature/review",
				HEAD:            "abc123",
				UpstreamRemote:  "origin",
				UpstreamBranch:  "feature/review",
				UnpushedCommits: 1,
			},
			PullRequest: preflight.PullRequest{
				Number:         req.pr,
				State:          "OPEN",
				BaseRepository: "owner/project",
				HeadBranch:     "feature/review",
				HeadRepository: "owner/project",
			},
			PushPlan: preflight.PushPlan{
				Enabled: req.name != "fetch",
				Remote:  "origin",
				Branch:  "feature/review",
				Command: []string{"git", "push", "origin", "HEAD:feature/review"},
			},
		}, nil
	})
}

func withFakeReviewRunWorktrees(t *testing.T) {
	t.Helper()
	oldCreate := createRunWorktree
	oldIntegrate := integrateRunWorktree
	oldCleanup := cleanupCleanRunWorktree
	oldPrune := pruneTerminalRunWorktrees
	createRunWorktree = func(_ context.Context, opts runworktree.CreateOptions) (runworktree.Ref, error) {
		userRoot := opts.UserRoot
		runID := opts.RunID
		path := filepath.Join(os.Getenv("HOME"), ".roundfix", "worktrees", "fake-review", runID)
		if err := os.MkdirAll(path, 0o755); err != nil {
			return runworktree.Ref{}, err
		}
		gitImplement(t, path, "init", "--initial-branch=main")
		gitImplement(t, path, "commit", "--allow-empty", "-m", "seed fake review run worktree")
		gitImplement(t, path, "branch", "-m", runworktree.BranchName(runID))
		return runworktree.Ref{
			RunID:    runID,
			Path:     path,
			Branch:   runworktree.BranchName(runID),
			UserRoot: userRoot,
		}, nil
	}
	integrateRunWorktree = func(context.Context, runworktree.Ref, string, string) (runworktree.IntegrationResult, error) {
		return runworktree.IntegrationResult{Mode: runworktree.ModeFastForwardMerge}, nil
	}
	cleanupCleanRunWorktree = func(_ context.Context, ref runworktree.Ref) error {
		return os.RemoveAll(ref.Path)
	}
	pruneTerminalRunWorktrees = func(context.Context, string, string, func(string) bool) ([]runworktree.PrunedRef, error) {
		return nil, nil
	}
	t.Cleanup(func() {
		createRunWorktree = oldCreate
		integrateRunWorktree = oldIntegrate
		cleanupCleanRunWorktree = oldCleanup
		pruneTerminalRunWorktrees = oldPrune
	})
}

func withFetchReviewItems(t *testing.T, items []reviewsource.ReviewItem) {
	t.Helper()
	withFetchReviewItemsFunc(t, func(context.Context, reviewsource.FetchRequest) ([]reviewsource.ReviewItem, error) {
		return items, nil
	})
}

func withFetchReviewItemsFunc(t *testing.T, fn func(context.Context, reviewsource.FetchRequest) ([]reviewsource.ReviewItem, error)) {
	t.Helper()
	old := fetchReviewItems
	fetchReviewItems = fn
	t.Cleanup(func() {
		fetchReviewItems = old
	})
}

func withPreflight(t *testing.T, fn func(context.Context, commandRequest, roundconfig.Loaded) (preflight.Result, error)) {
	t.Helper()
	old := runCommandPreflight
	runCommandPreflight = fn
	t.Cleanup(func() {
		runCommandPreflight = old
	})
}

func withReviewSpecGitRunner(t *testing.T, runner preflight.GitRunner) {
	t.Helper()
	old := reviewSpecGitRunner
	reviewSpecGitRunner = runner
	t.Cleanup(func() {
		reviewSpecGitRunner = old
	})
}

func overrideCollaborators(t *testing.T, mutate func(*engineCollaborators)) {
	t.Helper()
	old := newEngineCollaborators
	current := old()
	mutate(&current)
	newEngineCollaborators = func() engineCollaborators { return current }
	t.Cleanup(func() {
		newEngineCollaborators = old
	})
}

func withAgentRunner(t *testing.T, runner agent.Runner) {
	overrideCollaborators(t, func(collaborators *engineCollaborators) {
		collaborators.runner = runner
	})
	withRoundfixSessionLister(t, func(context.Context, agent.RuntimeSpec, string) ([]agent.RoundfixSession, error) {
		return nil, nil
	})
	withStopAgentSessionCloser(t, func(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
		return nil
	})
}

func withStopAgentSessionCanceler(t *testing.T, cancel func(context.Context, agent.RuntimeSpec, agent.SessionRef) error) {
	t.Helper()
	old := cancelStopAgentSession
	cancelStopAgentSession = cancel
	t.Cleanup(func() {
		cancelStopAgentSession = old
	})
	withRoundfixSessionLister(t, func(context.Context, agent.RuntimeSpec, string) ([]agent.RoundfixSession, error) {
		return nil, nil
	})
	withStopAgentSessionCloser(t, func(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
		return nil
	})
}

func withRoundfixSessionLister(t *testing.T, list func(context.Context, agent.RuntimeSpec, string) ([]agent.RoundfixSession, error)) {
	t.Helper()
	old := listRoundfixAgentSessions
	listRoundfixAgentSessions = list
	t.Cleanup(func() {
		listRoundfixAgentSessions = old
	})
}

func withStopAgentSessionCloser(t *testing.T, closeSession func(context.Context, agent.RuntimeSpec, agent.SessionRef) error) {
	t.Helper()
	old := closeStopAgentSession
	closeStopAgentSession = closeSession
	t.Cleanup(func() {
		closeStopAgentSession = old
	})
}

func fakeACPXVersionCommand(t *testing.T, version string) string {
	t.Helper()
	script := filepath.Join(t.TempDir(), "acpx")
	body := "#!/bin/sh\n" +
		"if [ \"$1\" = \"--version\" ]; then\n" +
		"  printf '%s\\n' '" + version + "'\n" +
		"  exit 0\n" +
		"fi\n" +
		"printf 'unexpected acpx invocation: %s\\n' \"$*\" >&2\n" +
		"exit 2\n"
	if err := os.WriteFile(script, []byte(body), 0o755); err != nil {
		t.Fatalf("write fake acpx command: %v", err)
	}
	return script
}

func withVerifier(t *testing.T, verifier daemon.Verifier) {
	overrideCollaborators(t, func(collaborators *engineCollaborators) {
		collaborators.verifier = verifier
	})
}

func withCommitter(t *testing.T, committer daemon.Committer) {
	overrideCollaborators(t, func(collaborators *engineCollaborators) {
		collaborators.committer = committer
	})
}

func withSourceResolver(t *testing.T, resolver *fakeSourceResolver) {
	overrideCollaborators(t, func(collaborators *engineCollaborators) {
		collaborators.source = daemon.ReviewSourceResolverFunc(resolver.Resolve)
	})
}

func withPusher(t *testing.T, pusher daemon.Pusher) {
	overrideCollaborators(t, func(collaborators *engineCollaborators) {
		collaborators.pusher = pusher
	})
}

// fakeWorktree reports an Agent-made change on every post-Batch snapshot so
// CLI tests exercise the commit path without a real git repository.
type fakeWorktree struct {
	mu    sync.Mutex
	calls int
}

func (worktree *fakeWorktree) Snapshot(context.Context, string) ([]string, error) {
	worktree.mu.Lock()
	defer worktree.mu.Unlock()
	worktree.calls++
	if worktree.calls%2 == 0 {
		return []string{"src/agent-change.go"}, nil
	}
	return nil, nil
}

func withFakeWorktree(t *testing.T) {
	overrideCollaborators(t, func(collaborators *engineCollaborators) {
		collaborators.worktree = &fakeWorktree{}
	})
}

func withWatchStatus(t *testing.T, fn func(context.Context, reviewsource.WatchStatusRequest) (reviewsource.WatchStatus, error)) {
	t.Helper()
	old := watchReviewStatus
	watchReviewStatus = fn
	t.Cleanup(func() {
		watchReviewStatus = old
	})
}

func withWatchHeadCheck(t *testing.T, fn func(context.Context, reviewsource.HeadCheckRequest) (watch.HeadCheckState, error)) {
	t.Helper()
	old := watchHeadCheck
	watchHeadCheck = fn
	t.Cleanup(func() {
		watchHeadCheck = old
	})
}

func withWatchHeadSHA(t *testing.T, fn func(context.Context, string) (string, error)) {
	t.Helper()
	old := watchHeadSHA
	watchHeadSHA = fn
	t.Cleanup(func() {
		watchHeadSHA = old
	})
}

func withWatchTiming(t *testing.T, clock watch.Clock, sleeper watch.Sleeper) {
	t.Helper()
	oldClock := watchClock
	oldSleeper := watchSleeper
	watchClock = clock
	watchSleeper = sleeper
	t.Cleanup(func() {
		watchClock = oldClock
		watchSleeper = oldSleeper
	})
}

func withInteractiveInput(t *testing.T, fn func(context.Context, roundtui.InputRequest) (roundtui.CommandValues, error)) {
	t.Helper()
	old := collectInteractiveInput
	collectInteractiveInput = fn
	t.Cleanup(func() {
		collectInteractiveInput = old
	})
}

func withInitScopePrompt(t *testing.T, fn func(context.Context, io.Writer) (string, error)) {
	t.Helper()
	old := promptInitScope
	promptInitScope = fn
	t.Cleanup(func() {
		promptInitScope = old
	})
}

func withClaudeSkillSymlinkPrompt(t *testing.T, fn func(context.Context, io.Writer, string, string) (bool, error)) {
	t.Helper()
	old := promptProjectClaudeSkillSymlink
	promptProjectClaudeSkillSymlink = fn
	t.Cleanup(func() {
		promptProjectClaudeSkillSymlink = old
	})
}

func withCurrentPullRequestSuggestion(t *testing.T, value string) {
	t.Helper()
	old := suggestCurrentPullRequest
	suggestCurrentPullRequest = func(context.Context, string) (string, error) {
		return value, nil
	}
	t.Cleanup(func() {
		suggestCurrentPullRequest = old
	})
}

func withStopPullRequestResolver(t *testing.T, result preflight.PullRequest) {
	t.Helper()
	old := resolvePullRequestForStop
	resolvePullRequestForStop = func(context.Context, string, string) (preflight.PullRequest, error) {
		return result, nil
	}
	t.Cleanup(func() {
		resolvePullRequestForStop = old
	})
}

func withChangedPaths(t *testing.T, changes []preflight.ChangedPath) {
	t.Helper()
	old := inspectChangedPaths
	inspectChangedPaths = func(context.Context, string) ([]preflight.ChangedPath, error) {
		return changes, nil
	}
	t.Cleanup(func() {
		inspectChangedPaths = old
	})
}

func readInteractiveDefaults(t *testing.T, homeDir string) store.InteractiveDefaults {
	t.Helper()
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store for defaults: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store for defaults: %v", err)
		}
	}()
	defaults, err := runStore.InteractiveDefaults(context.Background())
	if err != nil {
		t.Fatalf("read defaults: %v", err)
	}
	return defaults
}

func assertRunCount(t *testing.T, dbPath string, expected int) {
	t.Helper()
	ctx := context.Background()
	runStore, err := store.Open(ctx, filepath.Dir(filepath.Dir(dbPath)))
	if err != nil {
		t.Fatalf("open store for count: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store for count: %v", err)
		}
	}()
	count, err := runStore.RunCount(ctx)
	if err != nil {
		t.Fatalf("count runs: %v", err)
	}
	if count != expected {
		t.Fatalf("expected %d Run record(s), got %d", expected, count)
	}
}

func assertNoRunDatabase(t *testing.T, homeDir string) {
	t.Helper()
	if _, err := os.Stat(store.DatabasePath(homeDir)); err == nil {
		t.Fatalf("expected no Run Database at %s", store.DatabasePath(homeDir))
	} else if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("stat Run Database: %v", err)
	}
}

func createActiveImplementRunForStop(t *testing.T, homeDir string, repoDir string, specSlug string, agentID string) (store.Run, store.CreateRunRequest) {
	t.Helper()
	request := store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     repoDir,
		LocalBranch: "ma/implement-spec",
		HeadSHA:     "abc123",
		SpecSlug:    specSlug,
		Agent:       agentID,
	}
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store: %v", err)
		}
	}()
	active, err := runStore.CreateRun(context.Background(), request)
	if err != nil {
		t.Fatalf("create active implement run: %v", err)
	}
	return active, request
}

func createActiveImplementRunForStopInGitRoot(t *testing.T, homeDir string, gitRoot string, specSlug string) store.Run {
	t.Helper()
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store: %v", err)
		}
	}()
	active, err := runStore.CreateRun(context.Background(), store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     gitRoot,
		LocalBranch: "ma/other-work",
		HeadSHA:     "abc123",
		SpecSlug:    specSlug,
		Agent:       "codex",
	})
	if err != nil {
		t.Fatalf("create active implement run: %v", err)
	}
	return active
}

func assertNoActiveRun(t *testing.T, homeDir string, headRepository string, headBranch string) {
	t.Helper()
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store for active run: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store for active run: %v", err)
		}
	}()
	if active, found, err := runStore.ActiveRun(context.Background(), headRepository, headBranch); err != nil {
		t.Fatalf("lookup active run: %v", err)
	} else if found {
		t.Fatalf("expected no Active Run, got %#v", active)
	}
}

func assertActiveRun(t *testing.T, homeDir string, headRepository string, headBranch string, runID string) {
	t.Helper()
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store for active run: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store for active run: %v", err)
		}
	}()
	active, found, err := runStore.ActiveRun(context.Background(), headRepository, headBranch)
	if err != nil {
		t.Fatalf("lookup active run: %v", err)
	}
	if !found || active.ID != runID {
		t.Fatalf("expected Active Run %s, found=%v active=%#v", runID, found, active)
	}
}

func assertRunState(t *testing.T, homeDir string, runID string, want string) {
	t.Helper()
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store for Run state: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store for Run state: %v", err)
		}
	}()
	run, found, err := runStore.Run(context.Background(), runID)
	if err != nil {
		t.Fatalf("lookup Run state: %v", err)
	}
	if !found {
		t.Fatalf("expected Run %s to exist", runID)
	}
	if run.State != want {
		t.Fatalf("expected Run %s state %s, got %s", runID, want, run.State)
	}
}

func reviewRunIDFromStderr(t *testing.T, stderr string) string {
	t.Helper()
	for _, line := range strings.Split(stderr, "\n") {
		trimmed := strings.TrimSpace(line)
		if value, ok := strings.CutPrefix(trimmed, "Run: "); ok {
			return strings.TrimSpace(value)
		}
		if value, ok := strings.CutPrefix(trimmed, "Watch Run: "); ok {
			return strings.TrimSpace(value)
		}
	}
	t.Fatalf("expected review Run id in stderr, got %q", stderr)
	return ""
}

func assertStopRequested(t *testing.T, homeDir string, runID string) {
	t.Helper()
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store for Stop Request flag: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store for Stop Request flag: %v", err)
		}
	}()
	requested, err := runStore.StopRequested(context.Background(), runID)
	if err != nil {
		t.Fatalf("read Stop Request flag: %v", err)
	}
	if !requested {
		t.Fatalf("expected Stop Request recorded for Run %s", runID)
	}
}

func assertAgentLogContains(t *testing.T, repoDir string, expected string) {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join(builtinArtifactDirForRepo(t, repoDir), "runs", "*", "agent", "batch-001.log"))
	if err != nil {
		t.Fatalf("glob Agent logs: %v", err)
	}
	if len(matches) != 1 {
		t.Fatalf("expected one Agent log, got %#v", matches)
	}
	content, err := os.ReadFile(matches[0])
	if err != nil {
		t.Fatalf("read Agent log %s: %v", matches[0], err)
	}
	if !strings.Contains(string(content), expected) {
		t.Fatalf("expected Agent log to contain %q, got %s", expected, string(content))
	}
}

func assertNoAgentLogs(t *testing.T, repoDir string) {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join(builtinArtifactDirForRepo(t, repoDir), "runs", "*", "agent", "batch-*.log"))
	if err != nil {
		t.Fatalf("glob Agent logs: %v", err)
	}
	if len(matches) != 0 {
		t.Fatalf("expected no Agent logs, got %#v", matches)
	}
}

func persistCLIReviewIssue(t *testing.T, repoDir string, roundNumber int, headBranch string) rounds.PersistResult {
	t.Helper()
	return persistCLIReviewItems(t, repoDir, roundNumber, headBranch, []reviewsource.ReviewItem{
		{
			Title:                   "major: handle test issue",
			File:                    "internal/test.go",
			Line:                    12,
			Severity:                "major",
			Author:                  "coderabbitai[bot]",
			Body:                    "Keep this reviewer text literal.",
			SourceRef:               "thread:PRRT_test,comment:PRRC_test",
			ReviewHash:              "review-hash",
			SourceReviewID:          "9001",
			SourceReviewSubmittedAt: time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC).Add(time.Duration(roundNumber) * time.Minute),
		},
	})
}

func persistCLIReviewItems(t *testing.T, repoDir string, roundNumber int, headBranch string, items []reviewsource.ReviewItem) rounds.PersistResult {
	t.Helper()
	result, err := rounds.PersistRound(context.Background(), rounds.PersistRequest{
		ArtifactDir:    builtinArtifactDirForRepo(t, repoDir),
		ReviewRoot:     defaultReviewRootForRepo(repoDir, "123"),
		Source:         reviewsource.SourceCodeRabbit,
		PRNumber:       "123",
		HeadRepository: "owner/project",
		HeadBranch:     headBranch,
		HeadSHA:        "abc123",
		Round:          roundNumber,
		CreatedAt:      time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC).Add(time.Duration(roundNumber) * time.Minute),
		Items:          items,
	})
	if err != nil {
		t.Fatalf("persist CLI Review Issue artifact: %v", err)
	}
	return result
}

func defaultReviewRootForRepo(repoDir string, prNumber string) string {
	return filepath.Join(repoDir, "docs", "specs", "_reviews", "pr-"+prNumber)
}

func builtinArtifactDirForRepo(t *testing.T, repoDir string) string {
	t.Helper()
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		t.Fatal("HOME is required for builtin Artifact Directory")
	}
	sum := sha256.Sum256([]byte(filepath.Clean(repoDir)))
	repoID := hex.EncodeToString(sum[:])[:16]
	return filepath.Join(homeDir, ".roundfix", "artifacts", repoID)
}

type fakeVerifier struct {
	err      error
	calls    int
	commands []string
}

func (verifier *fakeVerifier) Verify(_ context.Context, req daemon.VerifyRequest) error {
	verifier.calls++
	verifier.commands = append(verifier.commands, req.Command)
	if req.Stream != nil {
		if _, err := io.WriteString(req.Stream, "fake verification output\n"); err != nil {
			return err
		}
	}
	return verifier.err
}

type fakeCommitter struct {
	err         error
	afterCommit func(context.Context, daemon.CommitRequest) error
	calls       int
	workDirs    []string
	messages    []string
}

func (committer *fakeCommitter) Commit(ctx context.Context, req daemon.CommitRequest) error {
	committer.calls++
	committer.workDirs = append(committer.workDirs, req.WorkDir)
	committer.messages = append(committer.messages, req.Message)
	if committer.err != nil {
		return committer.err
	}
	if committer.afterCommit != nil {
		return committer.afterCommit(ctx, req)
	}
	return nil
}

type fakeSourceResolver struct {
	err      error
	calls    int
	requests []reviewsource.ResolveRequest
}

func (resolver *fakeSourceResolver) Resolve(_ context.Context, req reviewsource.ResolveRequest) error {
	resolver.calls++
	resolver.requests = append(resolver.requests, req)
	return resolver.err
}

type fakePusher struct {
	err      error
	calls    int
	workDirs []string
	remotes  []string
	branches []string
	args     []string
}

func (pusher *fakePusher) Push(_ context.Context, req daemon.PushRequest) error {
	pusher.calls++
	pusher.workDirs = append(pusher.workDirs, req.WorkDir)
	pusher.remotes = append(pusher.remotes, req.Remote)
	pusher.branches = append(pusher.branches, req.Branch)
	pusher.args = append(pusher.args, "push", req.Remote, "HEAD:"+req.Branch)
	return pusher.err
}

type checkingPusher struct {
	calls   int
	workDir []string
	check   func(daemon.PushRequest) error
}

func (pusher *checkingPusher) Push(_ context.Context, req daemon.PushRequest) error {
	pusher.calls++
	pusher.workDir = append(pusher.workDir, req.WorkDir)
	if pusher.check != nil {
		return pusher.check(req)
	}
	return nil
}

type fakeWatchStatus struct {
	err      error
	calls    int
	statuses []reviewsource.WatchStatus
}

func (source *fakeWatchStatus) Status(context.Context, reviewsource.WatchStatusRequest) (reviewsource.WatchStatus, error) {
	source.calls++
	if source.err != nil {
		return reviewsource.WatchStatus{}, source.err
	}
	if len(source.statuses) == 0 {
		return reviewsource.WatchStatus{State: watch.StatusSettled}, nil
	}
	status := source.statuses[0]
	if len(source.statuses) > 1 {
		source.statuses = source.statuses[1:]
	}
	return status, nil
}

type fakeWatchHeadCheck struct {
	err      error
	calls    int
	states   []watch.HeadCheckState
	headSHAs []string
}

func (source *fakeWatchHeadCheck) Check(_ context.Context, req reviewsource.HeadCheckRequest) (watch.HeadCheckState, error) {
	source.calls++
	source.headSHAs = append(source.headSHAs, req.HeadSHA)
	if source.err != nil {
		return "", source.err
	}
	if len(source.states) == 0 {
		return watch.CheckSuccess, nil
	}
	state := source.states[0]
	if len(source.states) > 1 {
		source.states = source.states[1:]
	}
	return state, nil
}

type fakeWatchClock struct {
	now time.Time
}

func (clock *fakeWatchClock) Now() time.Time {
	return clock.now
}

func (clock *fakeWatchClock) Advance(duration time.Duration) {
	clock.now = clock.now.Add(duration)
}

type fakeWatchSleeper struct {
	clock *fakeWatchClock
}

func (sleeper *fakeWatchSleeper) Sleep(_ context.Context, duration time.Duration) error {
	if sleeper.clock != nil {
		sleeper.clock.Advance(duration)
	}
	return nil
}

type orderedWriteEvent struct {
	stream string
	text   string
}

type orderedWriteRecorder struct {
	events []orderedWriteEvent
}

func (recorder *orderedWriteRecorder) writer(stream string) io.Writer {
	return orderedStreamWriter{recorder: recorder, stream: stream}
}

func (recorder *orderedWriteRecorder) content(stream string) string {
	var builder strings.Builder
	for _, event := range recorder.events {
		if event.stream == stream {
			builder.WriteString(event.text)
		}
	}
	return builder.String()
}

func (recorder *orderedWriteRecorder) firstIndex(stream string) int {
	for index, event := range recorder.events {
		if event.stream == stream {
			return index
		}
	}
	return -1
}

func (recorder *orderedWriteRecorder) firstIndexContaining(stream string, needle string) int {
	for index, event := range recorder.events {
		if event.stream == stream && strings.Contains(event.text, needle) {
			return index
		}
	}
	return -1
}

type orderedStreamWriter struct {
	recorder *orderedWriteRecorder
	stream   string
}

func (writer orderedStreamWriter) Write(data []byte) (int, error) {
	writer.recorder.events = append(writer.recorder.events, orderedWriteEvent{
		stream: writer.stream,
		text:   string(data),
	})
	return len(data), nil
}

func publishFakeAgentOutput(ctx context.Context, sink runevent.Sink, req agent.ExecuteRequest, text string) error {
	if sink == nil {
		return nil
	}
	payload, err := json.Marshal(struct {
		Text string `json:"text"`
	}{Text: text})
	if err != nil {
		return err
	}
	return sink.Publish(ctx, runevent.RunEvent{
		RunID:   req.RunID,
		Batch:   req.Batch.Number,
		Source:  runevent.SourceAgent,
		Kind:    runevent.KindAgentRaw,
		Summary: text,
		Payload: payload,
	})
}

type fakeAgentRunner struct {
	probeErr       error
	runErr         error
	status         string
	statuses       []string
	onRun          func(agent.ExecuteRequest) error
	calls          int
	gitRoots       []string
	probeRequests  []agent.ProbeRequest
	probedRuntimes []agent.RuntimeSpec
	runRuntimes    []agent.RuntimeSpec
}

func (runner *fakeAgentRunner) Probe(_ context.Context, req agent.ProbeRequest) error {
	runner.probeRequests = append(runner.probeRequests, req)
	runner.probedRuntimes = append(runner.probedRuntimes, req.Runtime)
	return runner.probeErr
}

func (runner *fakeAgentRunner) Run(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	runner.calls++
	runner.gitRoots = append(runner.gitRoots, req.GitRoot)
	runner.runRuntimes = append(runner.runRuntimes, req.Runtime)
	status := runner.status
	if len(runner.statuses) >= runner.calls {
		status = runner.statuses[runner.calls-1]
	}
	if status == "" {
		status = rounds.StatusResolved
	}
	if runner.onRun != nil {
		if err := runner.onRun(req); err != nil {
			return agent.ExecuteResult{}, err
		}
	}
	for _, issue := range req.Batch.Issues {
		if err := rounds.SetIssueStatus(issue.Path, status, ""); err != nil {
			return agent.ExecuteResult{}, err
		}
	}
	output := "fake agent output\n"
	if err := publishFakeAgentOutput(ctx, sink, req, output); err != nil {
		return agent.ExecuteResult{}, err
	}
	if strings.TrimSpace(req.LogPath) != "" {
		if err := os.MkdirAll(filepath.Dir(req.LogPath), 0o755); err != nil {
			return agent.ExecuteResult{}, err
		}
		if err := os.WriteFile(req.LogPath, []byte(output), 0o644); err != nil {
			return agent.ExecuteResult{}, err
		}
	}
	if runner.runErr != nil {
		return agent.ExecuteResult{LogPath: req.LogPath, Output: output}, runner.runErr
	}
	return agent.ExecuteResult{LogPath: req.LogPath, Output: output}, nil
}

func (runner *fakeAgentRunner) EndSession(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
	return nil
}

type fakeStoppingAgentRunner struct{}

func (runner *fakeStoppingAgentRunner) Probe(context.Context, agent.ProbeRequest) error {
	return nil
}

func (runner *fakeStoppingAgentRunner) Run(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	output := "partial agent output\n"
	if err := publishFakeAgentOutput(ctx, sink, req, output); err != nil {
		return agent.ExecuteResult{}, err
	}
	if sink != nil {
		// Mirror the real runner: a stopped Agent publishes its status
		// event before returning.
		if err := sink.Publish(ctx, runevent.RunEvent{
			RunID:   req.RunID,
			Batch:   req.Batch.Number,
			Source:  runevent.SourceAgent,
			Kind:    runevent.KindAgentStatus,
			Summary: "SESSION STOPPED\n",
			Payload: []byte(`{"status":"stopped"}`),
		}); err != nil {
			return agent.ExecuteResult{}, err
		}
	}
	if strings.TrimSpace(req.LogPath) != "" {
		if err := os.MkdirAll(filepath.Dir(req.LogPath), 0o755); err != nil {
			return agent.ExecuteResult{}, err
		}
		if err := os.WriteFile(req.LogPath, []byte(output), 0o644); err != nil {
			return agent.ExecuteResult{}, err
		}
	}
	return agent.ExecuteResult{LogPath: req.LogPath, Output: output}, agent.StopError{
		LogPath: req.LogPath,
		Output:  output,
		Err:     context.Canceled,
	}
}

func (runner *fakeStoppingAgentRunner) EndSession(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
	return nil
}

type sessionRecordingRunner struct {
	inner          agent.Runner
	closeErr       error
	runSessions    []agent.SessionRef
	ensureSessions []agent.SessionRef
	closeSessions  []agent.SessionRef
	seenSessions   map[string]bool
}

func (runner *sessionRecordingRunner) Probe(ctx context.Context, req agent.ProbeRequest) error {
	return runner.inner.Probe(ctx, req)
}

func (runner *sessionRecordingRunner) Run(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	runner.runSessions = append(runner.runSessions, req.Session)
	sessionName := strings.TrimSpace(req.Session.Name)
	if sessionName != "" {
		if runner.seenSessions == nil {
			runner.seenSessions = map[string]bool{}
		}
		if !runner.seenSessions[sessionName] {
			runner.seenSessions[sessionName] = true
			runner.ensureSessions = append(runner.ensureSessions, req.Session)
			if err := publishFakeSessionStatus(ctx, sink, req, "session_started"); err != nil {
				return agent.ExecuteResult{}, err
			}
		}
	}
	return runner.inner.Run(ctx, req, sink)
}

func (runner *sessionRecordingRunner) EndSession(ctx context.Context, runtime agent.RuntimeSpec, session agent.SessionRef) error {
	runner.closeSessions = append(runner.closeSessions, session)
	if err := runner.inner.EndSession(ctx, runtime, session); err != nil {
		return err
	}
	return runner.closeErr
}

func publishFakeSessionStatus(ctx context.Context, sink runevent.Sink, req agent.ExecuteRequest, status string) error {
	if sink == nil {
		return nil
	}
	payload, err := json.Marshal(struct {
		Status string `json:"status"`
	}{Status: status})
	if err != nil {
		return err
	}
	return sink.Publish(ctx, runevent.RunEvent{
		RunID:   req.RunID,
		Batch:   req.Batch.Number,
		Source:  runevent.SourceAgent,
		Kind:    runevent.KindAgentStatus,
		Summary: "SESSION " + strings.ToUpper(status) + "\n",
		Payload: payload,
	})
}

func assertRecordedOneSessionForRun(t *testing.T, runner *sessionRecordingRunner, runID string, wantRunCalls int) {
	t.Helper()
	wantName := "roundfix-" + runID
	if len(runner.ensureSessions) != 1 || runner.ensureSessions[0].Name != wantName || strings.TrimSpace(runner.ensureSessions[0].WorkDir) == "" {
		t.Fatalf("expected exactly one ensure named %q with a working directory, got %#v", wantName, runner.ensureSessions)
	}
	want := runner.ensureSessions[0]
	if len(runner.closeSessions) != 1 || runner.closeSessions[0] != want {
		t.Fatalf("expected exactly one close for %#v, got %#v", want, runner.closeSessions)
	}
	if len(runner.runSessions) != wantRunCalls {
		t.Fatalf("expected %d Agent run session records, got %#v", wantRunCalls, runner.runSessions)
	}
	for _, session := range runner.runSessions {
		if session != want {
			t.Fatalf("expected every Agent run to use %#v, got %#v", want, runner.runSessions)
		}
	}
}

func assertSessionLifecycleEvents(t *testing.T, events []store.JournalEvent) {
	t.Helper()
	started := false
	closed := false
	for _, entry := range events {
		if entry.Event.Kind != runevent.KindAgentStatus {
			continue
		}
		payload := string(entry.Event.Payload)
		started = started || strings.Contains(payload, `"session_started"`)
		closed = closed || strings.Contains(payload, `"session_closed"`)
	}
	if !started || !closed {
		t.Fatalf("expected session_started and session_closed agent.status events, got %+v", events)
	}
}

func runFromStore(t *testing.T, homeDir string, runID string) store.Run {
	t.Helper()
	ctx := context.Background()
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store: %v", err)
		}
	}()
	run, found, err := runStore.Run(ctx, runID)
	if err != nil {
		t.Fatalf("read run %s: %v", runID, err)
	}
	if !found {
		t.Fatalf("run %s not found", runID)
	}
	return run
}

func journaledRunEvents(t *testing.T, homeDir string, stderr string) (string, []store.JournalEvent) {
	t.Helper()
	runID := ""
	for _, line := range strings.Split(stderr, "\n") {
		trimmed := strings.TrimSpace(line)
		if value, ok := strings.CutPrefix(trimmed, "Run: "); ok {
			runID = strings.TrimSpace(value)
			break
		}
		if value, ok := strings.CutPrefix(trimmed, "Watch Run: "); ok {
			runID = strings.TrimSpace(value)
			break
		}
		if value, ok := strings.CutPrefix(trimmed, "Implement Run: "); ok {
			runID = strings.TrimSpace(value)
			break
		}
	}
	if runID == "" {
		t.Fatalf("expected Run id in stderr, got %q", stderr)
	}
	reader, err := store.OpenReader(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open Run Database reader: %v", err)
	}
	defer func() {
		_ = reader.Close()
	}()
	events, err := reader.RunEventsAfter(context.Background(), runID, 0, 100)
	if err != nil {
		t.Fatalf("list journaled Run Events: %v", err)
	}
	return runID, events
}

func TestAgentConsoleDisplaySinkKeepsWriterBytesByDefault(t *testing.T) {
	event := runevent.RunEvent{
		Source:  runevent.SourceAgent,
		Kind:    runevent.KindAgentRaw,
		Payload: []byte(`{"text":"fake agent output\n"}`),
	}
	var direct bytes.Buffer
	var wrapped bytes.Buffer

	if err := (agent.WriterSink{Writer: &direct}).Publish(context.Background(), event); err != nil {
		t.Fatalf("publish through direct writer sink: %v", err)
	}
	if err := agentConsoleDisplaySink(&wrapped, false).Publish(context.Background(), event); err != nil {
		t.Fatalf("publish through default display sink: %v", err)
	}

	if wrapped.String() != direct.String() {
		t.Fatalf("expected default display bytes %q, got %q", direct.String(), wrapped.String())
	}
}

func assertJournalContainsAgentAndDaemonEvents(t *testing.T, events []store.JournalEvent, agentText string) {
	t.Helper()
	agentSeen := false
	daemonSeen := false
	for _, entry := range events {
		switch entry.Event.Source {
		case runevent.SourceAgent:
			if strings.Contains(entry.Event.Summary, agentText) || strings.Contains(string(entry.Event.Payload), agentText) {
				agentSeen = true
			}
		case runevent.SourceDaemon:
			daemonSeen = true
		}
	}
	if !agentSeen {
		t.Fatalf("expected Agent-source event containing %q in the Run Event Journal, got %+v", agentText, events)
	}
	if !daemonSeen {
		t.Fatalf("expected Daemon-source events in the Run Event Journal, got %+v", events)
	}
}

func TestResolveJournalsAgentRunEventsDurably(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withVerifier(t, &fakeVerifier{})
	withCommitter(t, &fakeCommitter{})
	withFakeWorktree(t)
	withSourceResolver(t, &fakeSourceResolver{})
	withPusher(t, &fakePusher{})
	withAgentRunner(t, &fakeAgentRunner{})
	withFakeWorktree(t)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean resolve exit, got %d stderr=%q", code, stderr.String())
	}
	runID, events := journaledRunEvents(t, homeDir, stderr.String())
	if len(events) == 0 {
		t.Fatal("expected journaled Run Events after the Resolve Run")
	}
	var event runevent.RunEvent
	for _, entry := range events {
		if entry.Event.Kind == runevent.KindAgentRaw {
			event = entry.Event
			break
		}
	}
	if event.Kind != runevent.KindAgentRaw || event.Source != runevent.SourceAgent {
		t.Fatalf("expected agent.raw journal entry, got %+v", events)
	}
	if event.RunID != runID || event.Batch != 1 {
		t.Fatalf("expected Run identity on journaled event, got %+v", event)
	}
	if !strings.Contains(event.Summary, "fake agent output") {
		t.Fatalf("expected bounded summary, got %q", event.Summary)
	}
	if !strings.Contains(string(event.Payload), "fake agent output") {
		t.Fatalf("expected raw payload preserved, got %s", event.Payload)
	}
	if !strings.Contains(stderr.String(), "fake agent output") {
		t.Fatalf("expected console output unchanged, got %q", stderr.String())
	}
}

func TestRunResolveSkipsAgentLogFilesByDefaultAndStillJournals(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean resolve exit, got %d stderr=%q", code, stderr.String())
	}
	assertNoAgentLogs(t, repoDir)
	_, events := journaledRunEvents(t, homeDir, stderr.String())
	assertJournalContainsAgentAndDaemonEvents(t, events, "fake agent output")
}

func TestRunResolveWritesAgentLogFilesWhenEnabledAndStillJournals(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), `
logs:
  agent: true
`)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean resolve exit, got %d stderr=%q", code, stderr.String())
	}
	assertAgentLogContains(t, repoDir, "fake agent output")
	_, events := journaledRunEvents(t, homeDir, stderr.String())
	assertJournalContainsAgentAndDaemonEvents(t, events, "fake agent output")
}

func TestRunResolveNoAgentConsoleSuppressesAgentDisplayOnly(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withVerifier(t, &fakeVerifier{})
	withCommitter(t, &fakeCommitter{})
	withFakeWorktree(t)
	withSourceResolver(t, &fakeSourceResolver{})
	withPusher(t, &fakePusher{})
	withAgentRunner(t, &fakeAgentRunner{})
	withFakeWorktree(t)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--agent", "codex", "--no-agent-console", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean resolve exit, got %d stderr=%q", code, stderr.String())
	}
	if strings.Contains(stderr.String(), "fake agent output") {
		t.Fatalf("expected Agent console hidden from stderr, got %q", stderr.String())
	}
	for _, want := range []string{
		"resolve selected 1 downloaded Unresolved Review Issue",
		"Batch: 001/001 (1 Review Issue(s))",
		"Verification command passed",
		"Batch commit created",
		"Resolved 1 Review Source thread",
		"Resolve Run",
	} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("expected stderr to keep daemon/progress line %q, got %q", want, stderr.String())
		}
	}
	_, events := journaledRunEvents(t, homeDir, stderr.String())
	assertJournalContainsAgentAndDaemonEvents(t, events, "fake agent output")
}

func TestRunWatchNoAgentConsoleSuppressesAgentDisplayOnly(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"watch", "--source", "coderabbit", "--pr", "123", "--agent", "codex", "--until-clean", "--max-rounds", "6", "--no-agent-console", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean watch exit, got %d stderr=%q", code, stderr.String())
	}
	if strings.Contains(stderr.String(), "fake agent output") {
		t.Fatalf("expected Agent console hidden from stderr, got %q", stderr.String())
	}
	for _, want := range []string{
		"Watch Run:",
		"Review Source status: settled",
		"Fetched Round 001 with 1 Review Issue",
		"Batch: 001/001 (1 Review Issue(s))",
		"Verification command passed",
		"Final Push completed",
		"Watch Run",
	} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("expected stderr to keep daemon/progress line %q, got %q", want, stderr.String())
		}
	}
	_, events := journaledRunEvents(t, homeDir, stderr.String())
	assertJournalContainsAgentAndDaemonEvents(t, events, "fake agent output")
}

func TestStoppedResolveJournalsStoppedEventBeforeReturning(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withVerifier(t, &fakeVerifier{})
	withCommitter(t, &fakeCommitter{})
	withFakeWorktree(t)
	withSourceResolver(t, &fakeSourceResolver{})
	withPusher(t, &fakePusher{})
	withAgentRunner(t, &fakeStoppingAgentRunner{})
	withFakeWorktree(t)
	withChangedPaths(t, []preflight.ChangedPath{{Status: "M", Path: "src/app.go"}})
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean Stop Request exit, got %d stderr=%q", code, stderr.String())
	}
	_, events := journaledRunEvents(t, homeDir, stderr.String())
	stoppedSeen := false
	for _, entry := range events {
		if entry.Event.Kind != runevent.KindAgentStatus {
			continue
		}
		if strings.Contains(string(entry.Event.Payload), `"stopped"`) {
			stoppedSeen = true
		}
	}
	if !stoppedSeen {
		t.Fatalf("expected stopped status event journaled, got %+v", events)
	}
}

func TestWatchSkipsFinalPushWhenAutoPushDisabled(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	pusher := &fakePusher{}
	withPusher(t, pusher)
	withPreflight(t, func(_ context.Context, req commandRequest, _ roundconfig.Loaded) (preflight.Result, error) {
		return preflight.Result{
			Git: preflight.GitState{
				Root:            repoDir,
				Branch:          "feature/review",
				HEAD:            "abc123",
				UpstreamRemote:  "origin",
				UpstreamBranch:  "feature/review",
				UnpushedCommits: 1,
			},
			PullRequest: preflight.PullRequest{
				Number:         req.pr,
				State:          "OPEN",
				BaseRepository: "owner/project",
				HeadBranch:     "feature/review",
				HeadRepository: "owner/project",
			},
			PushPlan: preflight.PushPlan{Enabled: false},
		}, nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"watch", "--source", "coderabbit", "--pr", "123", "--agent", "codex", "--until-clean", "--max-rounds", "6", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean watch exit, got %d stderr=%q", code, stderr.String())
	}
	if pusher.calls != 0 {
		t.Fatalf("expected Final Push gated off when auto-push is disabled, got %d push(es)", pusher.calls)
	}
	if !strings.Contains(stderr.String(), "Final Push skipped: auto-push disabled or no push target configured.") {
		t.Fatalf("expected gating message, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "reached Clean") {
		t.Fatalf("expected Clean watch outcome without push, got %q", stderr.String())
	}
}

func TestWatchFinalPushRunsOncePerCleanRoundThroughEngine(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	pusher := &fakePusher{}
	withPusher(t, pusher)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"watch", "--source", "coderabbit", "--pr", "123", "--agent", "codex", "--until-clean", "--max-rounds", "6", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean watch exit, got %d stderr=%q", code, stderr.String())
	}
	if pusher.calls != 1 {
		t.Fatalf("expected exactly one Final Push for the clean Round, got %d", pusher.calls)
	}
	if pusher.remotes[0] != "origin" || pusher.branches[0] != "feature/review" {
		t.Fatalf("expected push target from the push plan, got %v %v", pusher.remotes, pusher.branches)
	}
	if !strings.Contains(stderr.String(), "Batch commit created") {
		t.Fatalf("expected engine-driven Batch commit in watch Round, got %q", stderr.String())
	}
}

func runResolveForAttachTest(t *testing.T, repoDir string) (string, *bytes.Buffer) {
	t.Helper()
	withSuccessfulPreflight(t, repoDir)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("seed resolve run failed: %d stderr=%q", code, stderr.String())
	}
	runID := ""
	for _, line := range strings.Split(stderr.String(), "\n") {
		if value, ok := strings.CutPrefix(strings.TrimSpace(line), "Run: "); ok {
			runID = strings.TrimSpace(value)
			break
		}
	}
	if runID == "" {
		t.Fatalf("expected Run id in resolve output, got %q", stderr.String())
	}
	return runID, &stderr
}

func TestAttachReplaysCompletedRunReadOnly(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	runID, _ := runResolveForAttachTest(t, repoDir)
	// Attach must never probe or start an Agent; a probing attach would
	// fail loudly here.
	withAgentRunner(t, &fakeAgentRunner{probeErr: errors.New("attach must not probe Agents")})
	dbPath := filepath.Join(homeDir, ".roundfix", "roundfix.db")
	assertRunCount(t, dbPath, 1)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"attach", runID}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean attach exit, got %d stderr=%q", code, stderr.String())
	}
	output := stdout.String()
	for _, expected := range []string{
		"Roundfix attach",
		"ID: " + runID,
		"State: Clean",
		"fake agent output",
		"Run " + runID + " reached Clean; timeline replayed read-only.",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected attach output to contain %q, got:\n%s", expected, output)
		}
	}
	// Non-mutating: no new Runs, and the Run's terminal state is untouched.
	assertRunCount(t, dbPath, 1)
	reader, err := store.OpenReader(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open reader: %v", err)
	}
	defer func() { _ = reader.Close() }()
	run, found, err := reader.Run(context.Background(), runID)
	if err != nil || !found {
		t.Fatalf("lookup run after attach: found=%v err=%v", found, err)
	}
	if run.State != store.StateClean {
		t.Fatalf("expected attach to leave Run state untouched, got %q", run.State)
	}
}

func TestAttachDisplaysStoredSelectionAfterConfigChanges(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var resolveStdout bytes.Buffer
	var resolveStderr bytes.Buffer

	code := RunContext(context.Background(), []string{
		"resolve",
		"--pr", "123",
		"--agent", "codex",
		"--model", "historical-model",
		"--reasoning-effort", "historical-reasoning",
		"--no-input",
	}, &resolveStdout, &resolveStderr)
	if code != exitOK {
		t.Fatalf("seed resolve run failed: %d stderr=%q", code, resolveStderr.String())
	}
	runID := reviewRunIDFromStderr(t, resolveStderr.String())
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), `
runtimes:
  codex:
    model: changed-config-model
    reasoning_effort: changed-config-reasoning
`)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code = RunContext(context.Background(), []string{"attach", runID}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected attach exit 0, got %d stderr=%q", code, stderr.String())
	}
	output := stdout.String()
	for _, expected := range []string{
		"Agent: Codex",
		"Agent Model: historical-model",
		"Default Reasoning Effort: historical-reasoning",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected attach output to contain stored value %q, got:\n%s", expected, output)
		}
	}
	for _, unexpected := range []string{"changed-config-model", "changed-config-reasoning", "auto"} {
		if strings.Contains(output, unexpected) {
			t.Fatalf("expected attach output not to contain %q, got:\n%s", unexpected, output)
		}
	}
	if run := runFromStore(t, homeDir, runID); run.Model != "historical-model" || run.ReasoningEffort != "historical-reasoning" {
		t.Fatalf("expected stored historical selection to remain unchanged, got %#v", run)
	}
}

func TestAttachRunBrowserLoopOpensCockpitAndRefreshes(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	otherRepo := filepath.Join(t.TempDir(), "other-repo")
	mustMkdir(t, filepath.Join(otherRepo, ".git"))
	base := time.Date(2026, 7, 6, 14, 0, 0, 0, time.UTC)
	runs := seedRunsForList(t, homeDir, []runListSeed{
		{
			kind:      store.KindResolve,
			state:     store.StateClean,
			gitRoot:   repoDir,
			branch:    "feature/terminal",
			prNumber:  "201",
			createdAt: base.Add(time.Minute),
		},
		{
			kind:      store.KindImplement,
			state:     store.StateVerifying,
			gitRoot:   repoDir,
			branch:    "ma/spec-active",
			specSlug:  "0017-run-discovery",
			createdAt: base.Add(2 * time.Minute),
		},
		{
			kind:      store.KindResolve,
			state:     store.StateActive,
			gitRoot:   otherRepo,
			branch:    "feature/other",
			prNumber:  "999",
			createdAt: base.Add(3 * time.Minute),
		},
	})
	terminalRun, activeRun, otherRepoRun := runs[0], runs[1], runs[2]
	withAttachInteractiveInput(t, true)
	t.Setenv("ROUNDFIX_TUI", "always")
	var createdBehindCockpit store.Run
	cockpitCalls := withBrowserAttachCockpit(t, func(run store.Run, concurrency int) int {
		// Mutating the store while the cockpit is open proves the next
		// browser pass re-queries instead of reusing stale listings.
		createdBehindCockpit = seedRunsForList(t, homeDir, []runListSeed{{
			kind:      store.KindResolve,
			state:     store.StateActive,
			gitRoot:   repoDir,
			branch:    "feature/started-behind-cockpit",
			prNumber:  "202",
			createdAt: base.Add(4 * time.Minute),
		}})[0]
		return exitOK
	})
	sessionCalls := withRunBrowserSession(t,
		roundtui.BrowserOutcome{RunID: activeRun.ID},
		roundtui.BrowserOutcome{Cancelled: true},
	)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"attach"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected browser loop exit 0, got %d stderr=%q", code, stderr.String())
	}
	calls := *sessionCalls
	if len(calls) != 2 {
		t.Fatalf("expected browser sessions before and after the cockpit, got %d", len(calls))
	}
	first := calls[0]
	// The browser is machine-wide: Active Runs from every repository,
	// newest first.
	if len(first.active) != 2 || first.active[0].ID != otherRepoRun.ID || first.active[1].ID != activeRun.ID {
		t.Fatalf("expected every repository's Active Runs newest first, got %+v", first.active)
	}
	if len(first.all) != 3 || first.all[0].ID != otherRepoRun.ID || first.all[1].ID != activeRun.ID || first.all[2].ID != terminalRun.ID {
		t.Fatalf("expected every repository's Runs newest first, got %+v", first.all)
	}
	cockpit := *cockpitCalls
	if len(cockpit) != 1 || cockpit[0].runID != activeRun.ID {
		t.Fatalf("expected one cockpit pass for the selected Run, got %+v", cockpit)
	}
	if cockpit[0].concurrency != 2 {
		t.Fatalf("expected the explicit-attach concurrency fallback, got %d", cockpit[0].concurrency)
	}
	second := calls[1]
	if len(second.all) != 4 || second.all[0].ID != createdBehindCockpit.ID {
		t.Fatalf("expected the refreshed browser to list the Run created behind the cockpit, got %+v", second.all)
	}
}

func TestBrowserAttachCockpitIsTheExplicitAttachCockpit(t *testing.T) {
	if reflect.ValueOf(browserAttachCockpit).Pointer() != reflect.ValueOf(runAttachCockpit).Pointer() {
		t.Fatal("expected the browser loop to open the same attach cockpit as explicit attach")
	}
}

func TestAttachRunBrowserCancelExitsZeroWithoutAttaching(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	runID, _ := runResolveForAttachTest(t, repoDir)
	dbPath := filepath.Join(homeDir, ".roundfix", "roundfix.db")
	assertRunCount(t, dbPath, 1)
	withAttachInteractiveInput(t, true)
	t.Setenv("ROUNDFIX_TUI", "always")
	cockpitCalls := withBrowserAttachCockpit(t, nil)
	withRunBrowserSession(t, roundtui.BrowserOutcome{Cancelled: true})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"attach"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected cancelled browser exit 0, got %d stderr=%q", code, stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no attach output after browser cancel, got %q", stdout.String())
	}
	if len(*cockpitCalls) != 0 {
		t.Fatalf("expected no cockpit pass after cancel, got %+v", *cockpitCalls)
	}
	assertRunCount(t, dbPath, 1)
	reader, err := store.OpenReader(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open reader: %v", err)
	}
	defer func() { _ = reader.Close() }()
	run, found, err := reader.Run(context.Background(), runID)
	if err != nil || !found {
		t.Fatalf("lookup run after browser cancel: found=%v err=%v", found, err)
	}
	if run.State != store.StateClean {
		t.Fatalf("expected browser cancel to leave Run untouched, got %q", run.State)
	}
}

func TestAttachWithoutRunIDNonInteractiveNamesAllRunsList(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		interactive bool
	}{
		{name: "no tty", args: []string{"attach"}, interactive: false},
		{name: "no input flag", args: []string{"attach", "--no-input"}, interactive: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withCLIWorkspace(t)
			withAttachInteractiveInput(t, tt.interactive)
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := RunContext(context.Background(), tt.args, &stdout, &stderr)

			if code != exitPreflight {
				t.Fatalf("expected exit 2, got %d", code)
			}
			if stdout.Len() != 0 {
				t.Fatalf("expected no stdout, got %q", stdout.String())
			}
			if !strings.Contains(stderr.String(), "roundfix runs list --all") {
				t.Fatalf("expected discovery command guidance, got %q", stderr.String())
			}
			if !strings.Contains(stderr.String(), "pass a run id") {
				t.Fatalf("expected run id requirement, got %q", stderr.String())
			}
		})
	}
}

func TestAttachSpecRunReadsTasksFromKeptWorkDir(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	workDir := t.TempDir()
	writeImplementSpec(t, repoDir, implementTestSlug, []implementSeed{{id: "task_01", title: "Read state", status: "pending"}})
	writeImplementSpec(t, workDir, implementTestSlug, []implementSeed{{id: "task_01", title: "Read state", status: "completed"}})
	run := createTerminalAttachSpecRun(t, homeDir, repoDir, workDir, store.StateUnresolved)
	appendAttachConcurrencyEvent(t, homeDir, run.ID, 4)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"attach", run.ID}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean attach exit, got %d stderr=%q", code, stderr.String())
	}
	output := stdout.String()
	for _, expected := range []string{
		"Run Worktree: " + workDir,
		"Concurrency: 4",
		"task_01 completed — Read state",
		"Run " + run.ID + " reached Unresolved; timeline replayed read-only.",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected attach output to contain %q, got:\n%s", expected, output)
		}
	}
	if strings.Contains(output, "task_01 pending") {
		t.Fatalf("expected attach to ignore stale GitRoot task status, got:\n%s", output)
	}
}

func TestAttachSpecRunFallsBackToGitRootWhenCleanWorkDirIsPruned(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	workDir := filepath.Join(t.TempDir(), "pruned-worktree")
	writeImplementSpec(t, repoDir, implementTestSlug, []implementSeed{{id: "task_01", title: "Integrated state", status: "completed"}})
	run := createTerminalAttachSpecRun(t, homeDir, repoDir, workDir, store.StateClean)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"attach", run.ID}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean attach exit for pruned worktree, got %d stderr=%q", code, stderr.String())
	}
	output := stdout.String()
	for _, expected := range []string{
		"Run Worktree: " + workDir,
		"task_01 completed — Integrated state",
		"Run " + run.ID + " reached Clean; timeline replayed read-only.",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected attach output to contain %q, got:\n%s", expected, output)
		}
	}
}

func TestAttachUnknownRunFailsBeforeTUIStart(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	runResolveForAttachTest(t, repoDir)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"attach", "41"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("expected preflight exit for unknown Run, got %d", code)
	}
	want := `Run "41" does not exist; picker numbers are not stable Run ids — pass a run id or run 'roundfix attach' to pick interactively`
	if !strings.Contains(stderr.String(), want) {
		t.Fatalf("expected picker-number error %q, got %q", want, stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no replay output before failure, got %q", stdout.String())
	}
}

func TestAttachWithoutRunDatabaseFailsAsCLIError(t *testing.T) {
	withCLIWorkspace(t)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"attach", "run_x"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("expected preflight exit without a Run Database, got %d", code)
	}
	if !strings.Contains(stderr.String(), "roundfix attach failed:") {
		t.Fatalf("expected CLI error before TUI start, got %q", stderr.String())
	}
}

func TestAttachSkipsUnknownEventKindsOnReplay(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	runID, _ := runResolveForAttachTest(t, repoDir)
	writer, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open writer: %v", err)
	}
	if _, err := writer.AppendRunEvent(context.Background(), runevent.RunEvent{
		RunID:   runID,
		Source:  runevent.SourceDaemon,
		Kind:    "future.unknown",
		Summary: "future event",
		Time:    time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC),
		Payload: []byte(`{"new":"shape"}`),
	}); err != nil {
		t.Fatalf("append future event: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"attach", runID}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected unknown kinds skipped, got exit %d stderr=%q", code, stderr.String())
	}
	if strings.Contains(stdout.String(), "future event") {
		t.Fatalf("expected unknown kind not rendered, got %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "fake agent output") {
		t.Fatalf("expected known events still replayed, got %q", stdout.String())
	}
}

func createTerminalAttachSpecRun(t *testing.T, homeDir string, repoDir string, workDir string, state string) store.Run {
	t.Helper()
	ctx := context.Background()
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close store: %v", err)
		}
	}()
	run, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     repoDir,
		LocalBranch: "ma/widget-flow",
		HeadSHA:     "abc123",
		WorkDir:     workDir,
		SpecSlug:    implementTestSlug,
		Agent:       "codex",
	})
	if err != nil {
		t.Fatalf("create implement run: %v", err)
	}
	completed, err := runStore.CompleteRun(ctx, run.ID, state)
	if err != nil {
		t.Fatalf("complete implement run: %v", err)
	}
	return completed
}

func appendAttachConcurrencyEvent(t *testing.T, homeDir string, runID string, concurrency int) {
	t.Helper()
	writer, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open writer: %v", err)
	}
	defer func() {
		if err := writer.Close(); err != nil {
			t.Fatalf("close writer: %v", err)
		}
	}()
	_, err = writer.AppendRunEvent(context.Background(), runevent.RunEvent{
		RunID:   runID,
		Source:  runevent.SourceDaemon,
		Kind:    runevent.KindDaemonStatus,
		Summary: fmt.Sprintf("Task cycle started with concurrency %d.", concurrency),
		Payload: []byte(fmt.Sprintf(`{"concurrency":%d}`, concurrency)),
	})
	if err != nil {
		t.Fatalf("append concurrency event: %v", err)
	}
}

type fakeAttachSource struct {
	mu          sync.Mutex
	run         store.Run
	events      []store.JournalEvent
	version     int64
	eventReads  int
	versionRead int
}

func (source *fakeAttachSource) Run(context.Context, string) (store.Run, bool, error) {
	source.mu.Lock()
	defer source.mu.Unlock()
	return source.run, true, nil
}

func (source *fakeAttachSource) RunEventsAfter(_ context.Context, _ string, cursor int64, limit int) ([]store.JournalEvent, error) {
	source.mu.Lock()
	defer source.mu.Unlock()
	source.eventReads++
	page := []store.JournalEvent{}
	for _, entry := range source.events {
		if entry.Cursor > cursor && len(page) < limit {
			page = append(page, entry)
		}
	}
	return page, nil
}

func (source *fakeAttachSource) DataVersion(context.Context) (int64, error) {
	source.mu.Lock()
	defer source.mu.Unlock()
	source.versionRead++
	return source.version, nil
}

func (source *fakeAttachSource) appendEvent(summary string) {
	source.mu.Lock()
	defer source.mu.Unlock()
	cursor := int64(len(source.events) + 1)
	source.events = append(source.events, store.JournalEvent{
		Cursor: cursor,
		Event: runevent.RunEvent{
			RunID:   source.run.ID,
			Source:  runevent.SourceAgent,
			Kind:    runevent.KindAgentRaw,
			Summary: summary,
			Payload: []byte(`{"text":` + strconv.Quote(summary) + `}`),
		},
	})
	source.version++
}

func (source *fakeAttachSource) complete(state string) {
	source.mu.Lock()
	defer source.mu.Unlock()
	source.run.State = state
	source.version++
}

func TestAttachFollowerAppendsOnlyNewerEventsWithoutDuplicates(t *testing.T) {
	source := &fakeAttachSource{run: store.Run{ID: "run-1", State: store.StateActive}, version: 1}
	source.appendEvent("backlog one\n")
	source.appendEvent("backlog two\n")
	accepted := []string{}
	steps := 0
	follower := attachFollower{
		source: source,
		sleep: func(context.Context) error {
			steps++
			switch steps {
			case 1:
				source.appendEvent("live three\n")
			case 2:
				source.complete(store.StateClean)
			}
			return nil
		},
		accept: func(entry store.JournalEvent) {
			accepted = append(accepted, entry.Event.Summary)
		},
	}

	final, cursor, err := follower.follow(context.Background(), "run-1", 2)

	if err != nil {
		t.Fatalf("follow: %v", err)
	}
	if final.State != store.StateClean {
		t.Fatalf("expected terminal run returned, got %+v", final)
	}
	if strings.Join(accepted, "") != "live three\n" {
		t.Fatalf("expected only events newer than the replay cursor, got %v", accepted)
	}
	if cursor != 3 {
		t.Fatalf("expected cursor advanced to 3, got %d", cursor)
	}

	// Reconnect from the returned cursor: nothing replays twice.
	reaccepted := []string{}
	again := attachFollower{
		source: source,
		sleep:  func(context.Context) error { return nil },
		accept: func(entry store.JournalEvent) {
			reaccepted = append(reaccepted, entry.Event.Summary)
		},
	}
	if _, _, err := again.follow(context.Background(), "run-1", cursor); err != nil {
		t.Fatalf("reconnect follow: %v", err)
	}
	if len(reaccepted) != 0 {
		t.Fatalf("expected no duplicates across reconnects, got %v", reaccepted)
	}
}

func TestAttachFollowerIdlePollsReadNoEventRows(t *testing.T) {
	source := &fakeAttachSource{run: store.Run{ID: "run-1", State: store.StateActive}, version: 7}
	steps := 0
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	follower := attachFollower{
		source: source,
		sleep: func(context.Context) error {
			steps++
			if steps >= 5 {
				cancel()
			}
			return ctx.Err()
		},
		accept: func(store.JournalEvent) {},
	}

	_, _, err := follower.follow(ctx, "run-1", 0)

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected detach via cancellation, got %v", err)
	}
	if source.eventReads != 1 {
		t.Fatalf("expected only the initial drain to read rows; idle polls must not, got %d reads", source.eventReads)
	}
	if source.versionRead < 5 {
		t.Fatalf("expected change-detection polls to continue, got %d", source.versionRead)
	}
}

func TestAttachDetachLeavesRunActive(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	runResolveForAttachTest(t, repoDir)
	writer, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open writer: %v", err)
	}
	active, err := writer.CreateRun(context.Background(), store.CreateRunRequest{
		Kind:           store.KindResolve,
		HeadRepository: "owner/project",
		HeadBranch:     "feature/other",
		BaseRepository: "owner/project",
		PRNumber:       "456",
		GitRoot:        repoDir,
		LocalBranch:    "feature/other",
		HeadSHA:        "def456",
		ArtifactDir:    filepath.Join(repoDir, ".roundfix"),
	})
	if err != nil {
		t.Fatalf("create active run: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Detach on the first follow poll, after backlog replay.
	withAttachSleep(t, func(sleepCtx context.Context) error {
		cancel()
		return sleepCtx.Err()
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(ctx, []string{"attach", active.ID}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean detach exit, got %d stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Following live events.") {
		t.Fatalf("expected live-follow status visible, got %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "Detached; Run "+active.ID+" keeps going.") {
		t.Fatalf("expected detach message, got %q", stdout.String())
	}
	reader, err := store.OpenReader(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open reader: %v", err)
	}
	defer func() { _ = reader.Close() }()
	after, found, err := reader.Run(context.Background(), active.ID)
	if err != nil || !found {
		t.Fatalf("lookup run after detach: found=%v err=%v", found, err)
	}
	if store.IsTerminalState(after.State) {
		t.Fatalf("expected detach to leave the Run active, got %q", after.State)
	}
}

func TestAttachFollowsLiveRunToTerminalState(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	runResolveForAttachTest(t, repoDir)
	writer, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open writer: %v", err)
	}
	t.Cleanup(func() { _ = writer.Close() })
	active, err := writer.CreateRun(context.Background(), store.CreateRunRequest{
		Kind:           store.KindResolve,
		HeadRepository: "owner/project",
		HeadBranch:     "feature/other",
		BaseRepository: "owner/project",
		PRNumber:       "456",
		GitRoot:        repoDir,
		LocalBranch:    "feature/other",
		HeadSHA:        "def456",
		ArtifactDir:    filepath.Join(repoDir, ".roundfix"),
	})
	if err != nil {
		t.Fatalf("create active run: %v", err)
	}
	withAttachSleep(t, func(ctx context.Context) error { return ctx.Err() })

	writerDone := make(chan error, 1)
	go func() {
		for index := 0; index < 5; index++ {
			if _, err := writer.AppendRunEvent(context.Background(), runevent.RunEvent{
				RunID:   active.ID,
				Source:  runevent.SourceAgent,
				Kind:    runevent.KindAgentRaw,
				Summary: fmt.Sprintf("live line %d\n", index),
				Payload: []byte(fmt.Sprintf(`{"text":"live line %d\n"}`, index)),
			}); err != nil {
				writerDone <- err
				return
			}
		}
		_, err := writer.CompleteRun(context.Background(), active.ID, store.StateStopped)
		writerDone <- err
	}()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := RunContext(context.Background(), []string{"attach", active.ID}, &stdout, &stderr)

	if err := <-writerDone; err != nil {
		t.Fatalf("writer goroutine: %v", err)
	}
	if code != 0 {
		t.Fatalf("expected clean follow exit, got %d stderr=%q", code, stderr.String())
	}
	output := stdout.String()
	for index := 0; index < 5; index++ {
		needle := fmt.Sprintf("live line %d", index)
		if count := strings.Count(output, needle); count != 1 {
			t.Fatalf("expected %q exactly once, got %d in:\n%s", needle, count, output)
		}
	}
	if !strings.Contains(output, "Run "+active.ID+" reached Stopped.") {
		t.Fatalf("expected terminal follow exit message, got %q", output)
	}
}

func withAttachSleep(t *testing.T, sleep func(ctx context.Context) error) {
	t.Helper()
	old := attachSleep
	attachSleep = sleep
	t.Cleanup(func() {
		attachSleep = old
	})
}

type browserSessionCall struct {
	active []store.Run
	all    []store.Run
}

// withRunBrowserSession scripts Run Browser outcomes so entry-point tests
// drive the browser loop through the model seam, not terminal emulation.
func withRunBrowserSession(t *testing.T, outcomes ...roundtui.BrowserOutcome) *[]browserSessionCall {
	t.Helper()
	calls := &[]browserSessionCall{}
	old := runBrowserSession
	runBrowserSession = func(_ context.Context, _ io.Writer, active, all []store.Run) (roundtui.BrowserOutcome, error) {
		index := len(*calls)
		*calls = append(*calls, browserSessionCall{active: active, all: all})
		if index >= len(outcomes) {
			t.Fatalf("unexpected Run Browser session call %d", index+1)
		}
		return outcomes[index], nil
	}
	t.Cleanup(func() {
		runBrowserSession = old
	})
	return calls
}

type browserCockpitCall struct {
	runID       string
	concurrency int
}

// withBrowserAttachCockpit replaces the browser loop's cockpit step and
// records each pass; a nil handler just reports success.
func withBrowserAttachCockpit(t *testing.T, handler func(run store.Run, concurrency int) int) *[]browserCockpitCall {
	t.Helper()
	calls := &[]browserCockpitCall{}
	old := browserAttachCockpit
	browserAttachCockpit = func(_ context.Context, _ roundconfig.Loaded, _ *store.Store, run store.Run, concurrency int, _ io.Writer, _ io.Writer) int {
		*calls = append(*calls, browserCockpitCall{runID: run.ID, concurrency: concurrency})
		if handler == nil {
			return exitOK
		}
		return handler(run, concurrency)
	}
	t.Cleanup(func() {
		browserAttachCockpit = old
	})
	return calls
}

func withAttachInteractiveInput(t *testing.T, available bool) {
	t.Helper()
	old := attachInteractiveInputAvailable
	attachInteractiveInputAvailable = func() bool {
		return available
	}
	t.Cleanup(func() {
		attachInteractiveInputAvailable = old
	})
}

func daemonKindSequence(events []store.JournalEvent) []string {
	kinds := []string{}
	for _, entry := range events {
		if entry.Event.Source != runevent.SourceDaemon {
			continue
		}
		kinds = append(kinds, string(entry.Event.Kind))
	}
	return kinds
}

func assertOrderedSubsequence(t *testing.T, haystack []string, expected []string) {
	t.Helper()
	index := 0
	for _, item := range haystack {
		if index < len(expected) && item == expected[index] {
			index++
		}
	}
	if index != len(expected) {
		t.Fatalf("expected ordered subsequence %v in %v (matched %d)", expected, haystack, index)
	}
}

func TestWatchRunJournalsOrderedLoopNarrative(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"watch", "--source", "coderabbit", "--pr", "123", "--agent", "codex", "--until-clean", "--max-rounds", "6", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean watch exit, got %d stderr=%q", code, stderr.String())
	}
	_, events := journaledRunEvents(t, homeDir, stderr.String())
	kinds := daemonKindSequence(events)
	assertOrderedSubsequence(t, kinds, []string{
		string(runevent.KindDaemonReviewStatus),
		string(runevent.KindDaemonFetch),
		string(runevent.KindDaemonFetch),
		string(runevent.KindDaemonSelection),
		string(runevent.KindDaemonBatch),
		string(runevent.KindDaemonVerification),
		string(runevent.KindDaemonVerification),
		string(runevent.KindDaemonCommit),
		string(runevent.KindDaemonSourceResolution),
		string(runevent.KindDaemonBatch),
		string(runevent.KindDaemonPush),
		string(runevent.KindDaemonOutcome),
	})
	pushed := 0
	for _, entry := range events {
		if entry.Event.Kind == runevent.KindDaemonPush && strings.Contains(string(entry.Event.Payload), `"pushed"`) {
			pushed++
		}
	}
	if pushed != 1 {
		t.Fatalf("expected exactly one pushed decision on the clean Round, got %d", pushed)
	}
	last := events[len(events)-1].Event
	if last.Kind != runevent.KindDaemonOutcome || !strings.Contains(string(last.Payload), `"Clean"`) {
		t.Fatalf("expected Clean outcome event last, got %+v", last)
	}
}

func TestStoppedRunJournalsStopWithoutLaterUnsafeDaemonEvents(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withAgentRunner(t, &fakeStoppingAgentRunner{})
	withChangedPaths(t, []preflight.ChangedPath{{Status: "M", Path: "src/app.go"}})
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean Stop Request exit, got %d", code)
	}
	_, events := journaledRunEvents(t, homeDir, stderr.String())
	stopIndex := -1
	for index, entry := range events {
		if entry.Event.Kind == runevent.KindAgentStatus && strings.Contains(string(entry.Event.Payload), `"stopped"`) {
			stopIndex = index
		}
	}
	if stopIndex < 0 {
		t.Fatalf("expected stop event journaled, got %+v", events)
	}
	for _, entry := range events[stopIndex+1:] {
		switch entry.Event.Kind {
		case runevent.KindDaemonVerification, runevent.KindDaemonCommit, runevent.KindDaemonPush, runevent.KindDaemonSourceResolution, runevent.KindDaemonFetch:
			t.Fatalf("expected no unsafe daemon events after stop, got %+v", entry.Event)
		}
	}
	last := events[len(events)-1].Event
	if last.Kind != runevent.KindDaemonOutcome || !strings.Contains(string(last.Payload), `"Stopped"`) {
		t.Fatalf("expected Stopped outcome event last, got %+v", last)
	}
}

func TestFailedVerificationJournalsFailureWithoutCommitEvents(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withVerifier(t, &fakeVerifier{err: errors.New("tests failed")})
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code == 0 {
		t.Fatal("expected failed verification to fail the Run")
	}
	_, events := journaledRunEvents(t, homeDir, stderr.String())
	failed := false
	for _, entry := range events {
		if entry.Event.Kind == runevent.KindDaemonVerification && strings.Contains(string(entry.Event.Payload), `"failed"`) {
			failed = true
		}
		if entry.Event.Kind == runevent.KindDaemonCommit {
			t.Fatalf("expected no commit events after failed verification, got %+v", entry.Event)
		}
		if entry.Event.Kind == runevent.KindDaemonPush && strings.Contains(string(entry.Event.Payload), `"pushed"`) {
			t.Fatalf("expected no push after failed verification, got %+v", entry.Event)
		}
	}
	if !failed {
		t.Fatalf("expected verification failure event journaled, got %+v", events)
	}
	last := events[len(events)-1].Event
	if last.Kind != runevent.KindDaemonOutcome || !strings.Contains(string(last.Payload), `"Unresolved"`) {
		t.Fatalf("expected Unresolved outcome event last, got %+v", last)
	}
}

func TestTriageOnlyBatchJournalsCommitSkipDecision(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	// Identical snapshots: the Agent triaged without touching the worktree.
	overrideCollaborators(t, func(collaborators *engineCollaborators) {
		collaborators.worktree = staticWorktree{}
	})
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected triage-only Batch to succeed, got %d stderr=%q", code, stderr.String())
	}
	_, events := journaledRunEvents(t, homeDir, stderr.String())
	skipped := false
	for _, entry := range events {
		if entry.Event.Kind == runevent.KindDaemonCommit && strings.Contains(string(entry.Event.Payload), `"skipped"`) {
			skipped = true
		}
	}
	if !skipped {
		t.Fatalf("expected commit-skip decision journaled, got %v", daemonKindSequence(events))
	}
}

type staticWorktree struct{}

func (staticWorktree) Snapshot(context.Context, string) ([]string, error) {
	return []string{"user-wip.txt"}, nil
}

func TestAttachRendersWatchDaemonEventsInTimeline(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	var watchStdout bytes.Buffer
	var watchStderr bytes.Buffer
	code := Run([]string{"watch", "--source", "coderabbit", "--pr", "123", "--agent", "codex", "--until-clean", "--max-rounds", "6", "--no-input"}, &watchStdout, &watchStderr)
	if code != 0 {
		t.Fatalf("seed watch run failed: %d stderr=%q", code, watchStderr.String())
	}
	runID := ""
	for _, line := range strings.Split(watchStderr.String(), "\n") {
		if value, ok := strings.CutPrefix(strings.TrimSpace(line), "Watch Run: "); ok {
			runID = strings.TrimSpace(value)
			break
		}
	}
	if runID == "" {
		t.Fatalf("expected Watch Run id, got %q", watchStderr.String())
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	attachCode := RunContext(context.Background(), []string{"attach", runID}, &stdout, &stderr)

	if attachCode != 0 {
		t.Fatalf("expected clean attach exit, got %d stderr=%q", attachCode, stderr.String())
	}
	for _, expected := range []string{
		"Review Source status: settled",
		"Verification command passed:",
		"Batch commit created:",
		"Final Push completed:",
		"Run reached Clean.",
	} {
		if !strings.Contains(stdout.String(), expected) {
			t.Fatalf("expected attach timeline to narrate %q, got:\n%s", expected, stdout.String())
		}
	}
}

func TestResolvePrintsIssueSummaryAfterCompletion(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	persistCLIReviewIssue(t, repoDir, 1, "feature/review")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := RunContext(context.Background(), []string{"resolve", "--pr", "123", "--agent", "codex", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected clean resolve exit, got %d stderr=%q", code, stderr.String())
	}
	wantStdout := "" +
		"issue 001 resolved — major: handle test issue\n" +
		"Clean after 1 Round(s): 1 resolved, 0 invalid, 0 failed, 0 unresolved.\n"
	if stdout.String() != wantStdout {
		t.Fatalf("expected stdout report %q, got %q", wantStdout, stdout.String())
	}
	if !strings.Contains(stderr.String(), "Resolve Run") || !strings.Contains(stderr.String(), "reached Clean") {
		t.Fatalf("expected terminal daemon diagnostics on stderr, got %q", stderr.String())
	}
}

func TestStageableReviewRootClassifiesInsideOutsideAndSymlink(t *testing.T) {
	repoDir := t.TempDir()
	external := t.TempDir()
	mustMkdir(t, filepath.Join(repoDir, "docs", "specs", "_reviews", "pr-9"))
	mustMkdir(t, filepath.Join(external, "specs", "_reviews", "pr-9"))

	relative, ok := stageableReviewRoot(repoDir, filepath.Join(repoDir, "docs", "specs", "_reviews", "pr-9"))
	if !ok || relative != "docs/specs/_reviews/pr-9" {
		t.Fatalf("expected inside root to stage as docs/specs/_reviews/pr-9, got %q ok=%v", relative, ok)
	}
	if _, ok := stageableReviewRoot(repoDir, filepath.Join(external, "specs", "_reviews", "pr-9")); ok {
		t.Fatal("expected external root to be unstageable")
	}
	linkRepo := t.TempDir()
	mustMkdir(t, filepath.Join(linkRepo, "docs"))
	if err := os.Symlink(filepath.Join(external, "specs"), filepath.Join(linkRepo, "docs", "specs")); err != nil {
		t.Fatalf("create specs symlink: %v", err)
	}
	if _, ok := stageableReviewRoot(linkRepo, filepath.Join(linkRepo, "docs", "specs", "_reviews", "pr-9")); ok {
		t.Fatal("expected symlink-crossing root to be unstageable")
	}
}
