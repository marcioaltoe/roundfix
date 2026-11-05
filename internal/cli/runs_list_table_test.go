// Suite: Runs List formats
// Invariant: table output is readable while plain output keeps the stable bytes.
// Boundary IN: CLI flags, private environment, clock, zone, and disposable Run Database.
// Boundary OUT: terminal programs and Run execution.
package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"roundfix/internal/gittest"
	"roundfix/internal/store"
	roundtui "roundfix/internal/tui"
)

func withRunsListLocation(t *testing.T, loc *time.Location) {
	t.Helper()
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.runsListLocation = func() *time.Location { return loc }
	})
}

func runsListFormatEnvironment(t *testing.T) {
	t.Helper()
	setCommandEnvForTest(t, "ROUNDFIX_TUI", "never")
	setCommandEnvForTest(t, "TERM", "xterm")
	setCommandEnvForTest(t, "COLUMNS", "")
}

func assertRunsListTranscript(t *testing.T, args []string, wantOut, wantErr string, wantCode int) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := runCLI(t, args, &stdout, &stderr)
	if code != wantCode || stdout.String() != wantOut || stderr.String() != wantErr {
		t.Fatalf("%v: exit=%d stdout=%q stderr=%q; want exit=%d stdout=%q stderr=%q", args,
			code, stdout.String(), stderr.String(), wantCode, wantOut, wantErr)
	}
}

func TestRunsListPlainFormatKeepsTheStableLine(t *testing.T) {
	t.Parallel()
	home, repo := withCLIWorkspace(t)
	runsListFormatEnvironment(t)
	other := filepath.Join(t.TempDir(), "other-repo")
	gittest.MarkRepository(t, other)
	runs := seedRunsForListColumns(t, home, repo, other)
	for _, test := range []struct {
		name  string
		flags []string
	}{
		{name: "default"},
		{name: "all repositories", flags: []string{"--all"}},
		{name: "all states", flags: []string{"--state", "all"}},
		{name: "all repositories and states", flags: []string{"--all", "--state", "all"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			args := append([]string{"runs", "list"}, test.flags...)
			var stdout, stderr bytes.Buffer
			code := runCLI(t, args, &stdout, &stderr)
			if code != exitOK {
				t.Fatalf("stable line: exit=%d stderr=%q", code, stderr.String())
			}
			if len(test.flags) == 0 {
				want := fmt.Sprintf("%s  Active  implement  spec:0017-run-discovery  codex  2026-07-06T12:02:00Z  running 12m  ma/spec-run\n", runs[1].ID)
				if stdout.String() != want || stderr.String() != "(1 terminal Run(s) hidden; use --state all)\n" {
					t.Fatalf("stable line changed: stdout=%q stderr=%q", stdout.String(), stderr.String())
				}
			}
			args = append(args, "--format", "plain")
			assertRunsListTranscript(t, args, stdout.String(), stderr.String(), exitOK)
			setCommandEnvForTest(t, "ROUNDFIX_TUI", "always")
			assertRunsListTranscript(t, args, stdout.String(), stderr.String(), exitOK)
		})
	}
}

func seedRunsForTable(t *testing.T) (string, string, []store.Run) {
	t.Helper()
	home, _ := withCLIWorkspace(t)
	roots := t.TempDir()
	fiscus, quaestor := filepath.Join(roots, "fiscus"), filepath.Join(roots, "quaestor")
	gittest.MarkRepository(t, fiscus)
	gittest.MarkRepository(t, quaestor)
	setCommandEnvironmentForTest(t, home, fiscus)
	runsListFormatEnvironment(t)
	base := time.Date(2026, 7, 6, 12, 0, 0, 0, time.UTC)
	runs := seedRunsForList(t, home, []runListSeed{
		{kind: store.KindResolve, state: store.StateClean, gitRoot: fiscus, branch: "feature/review", prNumber: "123", createdAt: base, completedAt: base.Add(42 * time.Minute)},
		{kind: store.KindImplement, state: store.StateActive, gitRoot: fiscus, branch: "feat/fiscus", specSlug: "0063-um-gtin-por-mestre-e-a-pendencia-de-identidade", createdAt: base.Add(2 * time.Minute)},
		{kind: store.KindImplement, state: store.StateActive, gitRoot: quaestor, branch: "feat/quaestor", specSlug: "0063-sequencia-pela-emissao-propria-e-pdf-ao-cliente", createdAt: base.Add(4 * time.Minute)},
	})
	withRunsListNow(t, base.Add(14*time.Minute))
	withRunsListLocation(t, time.FixedZone("-03", -3*60*60))
	return home, fiscus, runs
}

func TestRunsListTableShowsAlignedColumns(t *testing.T) {
	t.Parallel()
	home, fiscus, runs := seedRunsForTable(t)
	want := fmt.Sprintf("RUN       REPO      STATE   TARGET                                                AGENT  STARTED (-03)  DURATION\n%s  quaestor  Active  0063-sequencia-pela-emissao-propria-e-pdf-ao-cliente  codex  09:04          10m\n%s  fiscus    Active  0063-um-gtin-por-mestre-e-a-pendencia-de-identidade   codex  09:02          12m\n",
		roundtui.ShortRunID(runs[2].ID), roundtui.ShortRunID(runs[1].ID))
	const note = "(1 terminal Run(s) hidden; use --state all)\n"
	assertRunsListTranscript(t, []string{"runs", "list", "--all", "--format", "table"}, want, note, exitOK)
	t.Run("one repository", func(t *testing.T) {
		want := fmt.Sprintf("RUN       STATE   TARGET                                               AGENT  STARTED (-03)  DURATION\n%s  Active  0063-um-gtin-por-mestre-e-a-pendencia-de-identidade  codex  09:02          12m\n", roundtui.ShortRunID(runs[1].ID))
		assertRunsListTranscript(t, []string{"runs", "list", "--format", "table"}, want, note, exitOK)
	})
	t.Run("terminal detection", func(t *testing.T) {
		setCommandEnvForTest(t, "ROUNDFIX_TUI", "always")
		assertRunsListTranscript(t, []string{"runs", "list", "--all"}, want, note, exitOK)
	})
	t.Run("outside repository", func(t *testing.T) {
		setCommandEnvironmentForTest(t, home, t.TempDir())
		assertRunsListTranscript(t, []string{"runs", "list", "--format", "table"}, want, note, exitOK)
	})
	t.Run("invalid width falls back to no budget", func(t *testing.T) {
		for _, columns := range []string{"invalid", "39", "0"} {
			t.Run(columns, func(t *testing.T) {
				setCommandEnvForTest(t, "COLUMNS", columns)
				assertRunsListTranscript(t, []string{"runs", "list", "--all", "--format", "table"}, want, note, exitOK)
			})
		}
	})
	t.Run("another day", func(t *testing.T) {
		withRunsListNow(t, time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC))
		want := fmt.Sprintf("RUN       REPO      STATE   TARGET                                                AGENT  STARTED (-03)  DURATION\n%s  quaestor  Active  0063-sequencia-pela-emissao-propria-e-pdf-ao-cliente  codex  Jul 06 09:04   47h56m\n%s  fiscus    Active  0063-um-gtin-por-mestre-e-a-pendencia-de-identidade   codex  Jul 06 09:02   47h58m\n", roundtui.ShortRunID(runs[2].ID), roundtui.ShortRunID(runs[1].ID))
		assertRunsListTranscript(t, []string{"runs", "list", "--all", "--format", "table"}, want, note, exitOK)
	})
	t.Run("terminal duration and target", func(t *testing.T) {
		want := fmt.Sprintf("RUN       STATE  TARGET  AGENT  STARTED (-03)  DURATION\n%s  Clean  pr:123  codex  09:00          42m\n", roundtui.ShortRunID(runs[0].ID))
		assertRunsListTranscript(t, []string{"runs", "list", "--state", "terminal", "--format", "table"}, want, "(1 active Run(s) hidden; use --state all)\n", exitOK)
	})
	t.Run("owner identity unproven", func(t *testing.T) {
		home := t.TempDir()
		setCommandEnvironmentForTest(t, home, fiscus)
		base := time.Date(2026, 7, 6, 12, 0, 0, 0, time.UTC)
		owner := seedRunsForList(t, home, []runListSeed{{kind: store.KindImplement, state: store.StateActive,
			gitRoot: fiscus, branch: "ma/spec-run", specSlug: "0055-owner-identity-without-fork", ownerPID: os.Getpid(), createdAt: base}})
		withRunsListNow(t, base.Add(3*time.Minute))
		want := fmt.Sprintf("RUN       STATE   TARGET                            AGENT  STARTED (-03)  DURATION  NOTE\n%s  Active  0055-owner-identity-without-fork  codex  09:00          3m        owner identity unproven\n", roundtui.ShortRunID(owner[0].ID))
		assertRunsListTranscript(t, []string{"runs", "list", "--format", "table"}, want, "", exitOK)
	})
}

func TestRunsListTableFitsTheWidthByCuttingTheTarget(t *testing.T) {
	t.Parallel()
	_, _, runs := seedRunsForTable(t)
	setCommandEnvForTest(t, "COLUMNS", "80")
	want := fmt.Sprintf("RUN       REPO      STATE   TARGET                AGENT  STARTED (-03)  DURATION\n%s  quaestor  Active  0063-sequencia-pela…  codex  09:04          10m\n%s  fiscus    Active  0063-um-gtin-por-me…  codex  09:02          12m\n", roundtui.ShortRunID(runs[2].ID), roundtui.ShortRunID(runs[1].ID))
	assertRunsListTranscript(t, []string{"runs", "list", "--all", "--format", "table"}, want, "(1 terminal Run(s) hidden; use --state all)\n", exitOK)
	for _, line := range strings.Split(want, "\n") {
		if utf8.RuneCountInString(line) > 80 {
			t.Fatalf("line wider than 80: %q", line)
		}
	}
}

func TestRunsListWideShowsEveryColumn(t *testing.T) {
	t.Parallel()
	_, _, runs := seedRunsForTable(t)
	setCommandEnvForTest(t, "COLUMNS", "80")
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"runs", "list", "--all", "--wide"}, &stdout, &stderr)
	if code != exitOK || stderr.String() != "(1 terminal Run(s) hidden; use --state all)\n" {
		t.Fatalf("wide: exit=%d stderr=%q", code, stderr.String())
	}
	lines := strings.Split(strings.TrimSuffix(stdout.String(), "\n"), "\n")
	if len(lines) != 3 {
		t.Fatalf("wide rows=%q", stdout.String())
	}
	header := strings.Fields(lines[0])
	if strings.Join(header, " ") != "RUN REPO STATE KIND TARGET AGENT STARTED (-03) DURATION BRANCH WORKTREE" {
		t.Fatalf("wide header=%q", lines[0])
	}
	for i, run := range []store.Run{runs[2], runs[1]} {
		cells := strings.Fields(lines[i+1])
		if len(cells) != 10 || cells[0] != run.ID || cells[3] != run.Kind || cells[4] != run.SpecSlug || cells[8] != run.LocalBranch || cells[9] != run.GitRoot {
			t.Fatalf("wide row omitted a full field: %q; Run=%+v", lines[i+1], run)
		}
	}
	if strings.Contains(stdout.String(), "…") {
		t.Fatalf("wide output cut a cell: %q", stdout.String())
	}
}

func TestRunsListFormatTranscripts(t *testing.T) {
	t.Parallel()
	runsListFormatEnvironment(t)
	for _, test := range []struct {
		name           string
		flags          []string
		stdout, stderr string
		code           int
	}{
		{name: "unknown format", flags: []string{"--format", "yaml"}, stderr: "roundfix: runs list failed: unknown --format \"yaml\"; use table or plain\nRun 'roundfix runs list --help' for usage.\n", code: 2},
		{name: "wide plain", flags: []string{"--wide", "--format", "plain"}, stderr: "roundfix: runs list failed: --wide applies only to --format table\nRun 'roundfix runs list --help' for usage.\n", code: 2},
		{name: "no Run Database", flags: []string{"--format", "table"}, stdout: "No Runs found.\n"},
	} {
		t.Run(test.name, func(t *testing.T) {
			setCommandEnvironmentForTest(t, t.TempDir(), t.TempDir())
			assertRunsListTranscript(t, append([]string{"runs", "list"}, test.flags...), test.stdout, test.stderr, test.code)
		})
	}
}

func TestRunsHelpNamesTheFormatAndWideFlags(t *testing.T) {
	t.Parallel()
	setCommandEnvironmentForTest(t, t.TempDir(), t.TempDir())
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"runs", "--help"}, &stdout, &stderr)
	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("help: exit=%d stderr=%q", code, stderr.String())
	}
	for _, line := range []string{
		"  roundfix runs list [--all] [--state <active|terminal|all>] [--limit N] [--format <table|plain>] [--wide]\n",
		"  --format list: table or plain (default: table at an interactive terminal,\n",
		"  --wide   list: table with the full Run id, kind, branch and worktree\n",
	} {
		if !strings.Contains(stdout.String(), line) {
			t.Errorf("help missing %q", line)
		}
	}
}

func TestRunsListTableKeepsTheWarningAndHiddenNote(t *testing.T) {
	t.Parallel()
	home, _, _ := seedRunsForTable(t)
	outside := filepath.Join(home, "drive")
	mustMkdir(t, outside)
	gittest.MarkRepository(t, home)
	setCommandEnvironmentForTest(t, home, outside)
	warning := fmt.Sprintf("roundfix: warning: worktree.location must not be inside the repository tree: %q is inside the Git root %q (found %q); commands that create a Run Worktree refuse it\n", filepath.Join(home, ".roundfix", "worktrees"), home, filepath.Join(home, ".git"))
	wantErr := warning + "(1 terminal Run(s) hidden; use --state all)\n"
	for _, format := range []string{"plain", "table"} {
		t.Run(format, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := runCLI(t, []string{"runs", "list", "--all", "--format", format}, &stdout, &stderr)
			if code != exitOK || stdout.Len() == 0 || stderr.String() != wantErr {
				t.Fatalf("exit=%d stdout=%q stderr=%q; want exit=0, rows, stderr=%q", code, stdout.String(), stderr.String(), wantErr)
			}
		})
	}
}
