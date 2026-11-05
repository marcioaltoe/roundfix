// Suite: Read-Only Run Command configuration and output
// Invariant: readers preserve output and warn once when a worktree location is refused.
// Boundary IN: CLI entry points, temporary homes, Run Database, and browser seams.
// Boundary OUT: Agent execution and terminal rendering, covered by their owning suites.
package cli

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	roundtui "roundfix/internal/tui"
)

func TestReadOnlyRunCommandsKeepTheirOutputInsideARepository(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := withCLIWorkspace(t)
	setCommandEnvForTest(t, "ROUNDFIX_TUI", "never")
	setCommandEnvForTest(t, "ROUNDFIX_COLOR", "never")
	setCommandEnvForTest(t, "COLUMNS", "100")
	runID, _ := runResolveForAttachTest(t, repoDir)
	outsideDir := t.TempDir()

	for _, args := range [][]string{
		{"runs", "show", runID},
		{"runs", "show", runID, "--json"},
		{"events", runID},
		{"attach", runID},
		{"runs", "list", "--all", "--state", "all"},
	} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			setCommandEnvironmentForTest(t, homeDir, repoDir)
			var inside, insideErr bytes.Buffer
			if code := runCLI(t, args, &inside, &insideErr); code != exitOK || insideErr.Len() != 0 {
				t.Fatalf("inside: exit=%d stderr=%q", code, insideErr.String())
			}
			if inside.Len() == 0 {
				t.Fatal("expected stored Run output")
			}
			setCommandEnvironmentForTest(t, homeDir, outsideDir)
			var outside, outsideErr bytes.Buffer
			if code := runCLI(t, args, &outside, &outsideErr); code != exitOK || outsideErr.Len() != 0 {
				t.Fatalf("outside: exit=%d stderr=%q", code, outsideErr.String())
			}
			if outside.String() != inside.String() {
				t.Fatalf("stdout changed outside repository:\ninside: %q\noutside: %q", inside.String(), outside.String())
			}
		})
	}
	t.Run("repository scoped list", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		if code := runCLI(t, []string{"runs", "list", "--state", "all"}, &stdout, &stderr); code != exitOK || stderr.Len() != 0 {
			t.Fatalf("exit=%d stderr=%q", code, stderr.String())
		}
		rows := strings.Split(strings.TrimSuffix(stdout.String(), "\n"), "\n")
		if len(rows) != 1 || !strings.HasPrefix(rows[0], runID+"  Clean  resolve  pr:123  ") {
			t.Fatalf("expected exactly the resolved Run's row, got %q", stdout.String())
		}
	})
	for _, args := range [][]string{{"attach"}, {"runs"}} {
		t.Run(strings.Join(args, " ")+" browser inside", func(t *testing.T) {
			assertRunReaderBrowser(t, args, runID, "")
		})
		t.Run(strings.Join(args, " ")+" browser outside", func(t *testing.T) {
			setCommandEnvironmentForTest(t, homeDir, outsideDir)
			assertRunReaderBrowser(t, args, runID, "")
		})
	}
}

func assertRunReaderBrowser(t *testing.T, args []string, runID, warning string) {
	t.Helper()
	withAttachInteractiveInput(t, true)
	withRunsInteractiveInput(t, true)
	setCommandEnvForTest(t, "ROUNDFIX_TUI", "always")
	cockpit := withBrowserAttachCockpit(t, nil)
	sessions := withRunBrowserSession(t,
		roundtui.BrowserOutcome{RunID: runID},
		roundtui.BrowserOutcome{Cancelled: true},
	)
	var stdout, stderr bytes.Buffer
	if code := runCLI(t, args, &stdout, &stderr); code != exitOK {
		t.Fatalf("browser exit=%d stderr=%q", code, stderr.String())
	}
	// The scripted cancel has no reason, so the closing line names none.
	wantStderr := warning + "roundfix: Run Browser closed; no Run selected.\n"
	if stdout.Len() != 0 || stderr.String() != wantStderr {
		t.Fatalf("browser stdout=%q stderr=%q; want empty stdout and stderr=%q", stdout.String(), stderr.String(), wantStderr)
	}
	if len(*sessions) != 2 || len(*cockpit) != 1 || (*cockpit)[0].runID != runID {
		t.Fatalf("expected select, cockpit, then cancel: sessions=%+v cockpit=%+v", *sessions, *cockpit)
	}
	for _, session := range *sessions {
		if len(session.active) != 0 || len(session.all) != 1 || session.all[0].ID != runID {
			t.Fatalf("expected only the resolved Run in browser: %+v", session)
		}
	}
}

func TestReadOnlyRunCommandsWarnOnAStrayGitDirectory(t *testing.T) {
	t.Parallel()
	for _, marker := range []string{"directory", "file"} {
		t.Run(marker, func(t *testing.T) {
			homeDir, repoDir := withCLIWorkspace(t)
			setCommandEnvForTest(t, "ROUNDFIX_TUI", "never")
			setCommandEnvForTest(t, "ROUNDFIX_COLOR", "never")
			setCommandEnvForTest(t, "COLUMNS", "100")
			runID, _ := runResolveForAttachTest(t, repoDir)
			drive := filepath.Join(homeDir, "drive")
			mustMkdir(t, drive)
			if marker == "directory" {
				mustMkdir(t, filepath.Join(homeDir, ".git"))
			} else {
				mustWrite(t, filepath.Join(homeDir, ".git"), "gitdir: "+t.TempDir()+"\n")
			}
			warning := fmt.Sprintf("roundfix: warning: worktree.location must not be inside the repository tree: %q is inside the Git root %q (found %q); commands that create a Run Worktree refuse it\n",
				filepath.Join(homeDir, ".roundfix", "worktrees"), homeDir, filepath.Join(homeDir, ".git"))
			for _, args := range [][]string{
				{"runs", "show", runID},
				{"runs", "show", runID, "--json"},
				{"events", runID},
				{"attach", runID},
				{"runs", "list", "--all", "--state", "all"},
				// A later hidden-Run diagnostic must follow the warning.
				{"runs", "list", "--all"},
			} {
				t.Run(strings.Join(args, " "), func(t *testing.T) {
					setCommandEnvironmentForTest(t, homeDir, repoDir)
					var inside, insideErr bytes.Buffer
					insideCode := runCLI(t, args, &inside, &insideErr)
					if insideCode != exitOK {
						t.Fatalf("inside exit=%d stderr=%q", insideCode, insideErr.String())
					}
					setCommandEnvironmentForTest(t, homeDir, drive)
					var stray, strayErr bytes.Buffer
					code := runCLI(t, args, &stray, &strayErr)
					if code != insideCode || stray.String() != inside.String() || strayErr.String() != warning+insideErr.String() {
						t.Fatalf("stray: exit=%d stdout=%q stderr=%q; want exit=%d stdout=%q stderr=%q",
							code, stray.String(), strayErr.String(), insideCode, inside.String(), warning+insideErr.String())
					}
				})
			}
			t.Run("repository scope stays the stray root", func(t *testing.T) {
				setCommandEnvironmentForTest(t, homeDir, drive)
				var stdout, stderr bytes.Buffer
				code := runCLI(t, []string{"runs", "list", "--state", "all"}, &stdout, &stderr)
				if code != exitOK || stdout.String() != "No Runs found.\n" || stderr.String() != warning {
					t.Fatalf("exit=%d stdout=%q stderr=%q; want exit=0, no Runs, and warning=%q", code, stdout.String(), stderr.String(), warning)
				}
			})
			for _, args := range [][]string{{"attach"}, {"runs"}} {
				t.Run(strings.Join(args, " ")+" browser", func(t *testing.T) {
					setCommandEnvironmentForTest(t, homeDir, drive)
					assertRunReaderBrowser(t, args, runID, warning)
				})
			}
		})
	}
}

func TestStrayGitDirectoryTranscripts(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name   string
		args   []string
		stdout string
		stderr string
		code   int
	}{
		{
			name:   "Surface Transcript 2",
			args:   []string{"runs", "list"},
			stdout: "No Runs found.\n",
			stderr: "roundfix: warning: worktree.location must not be inside the repository tree: \"<home>/.roundfix/worktrees\" is inside the Git root \"<home>\" (found \"<home>/.git\"); commands that create a Run Worktree refuse it\n",
			code:   0,
		},
		{
			name:   "Surface Transcript 3",
			args:   []string{"reconcile"},
			stderr: "roundfix: reconcile failed: worktree.location must not be inside the repository tree: \"<home>/.roundfix/worktrees\" is inside the Git root \"<home>\" (found \"<home>/.git\")\nRun 'roundfix reconcile --help' for usage.\n",
			code:   2,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			homeDir := t.TempDir()
			drive := filepath.Join(homeDir, "drive")
			mustMkdir(t, drive)
			mustMkdir(t, filepath.Join(homeDir, ".git"))
			setCommandEnvironmentForTest(t, homeDir, drive)
			var stdout, stderr bytes.Buffer
			code := runCLI(t, test.args, &stdout, &stderr)
			wantStderr := strings.ReplaceAll(test.stderr, "<home>", homeDir)
			if code != test.code || stdout.String() != test.stdout || stderr.String() != wantStderr {
				t.Fatalf("transcript: exit=%d stdout=%q stderr=%q; want exit=%d stdout=%q stderr=%q",
					code, stdout.String(), stderr.String(), test.code, test.stdout, wantStderr)
			}
		})
	}
}
