// Suite: short Run ids at read-only command boundaries.
// Invariant: unique prefixes replay the full Run; collisions list every match.
// Boundary IN: CLI dispatch and seeded disposable Run Databases.
// Boundary OUT: terminal programs and real Agent sessions.
package cli

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"roundfix/internal/store"
	roundtui "roundfix/internal/tui"
)

func TestAttachAndRunsShowAcceptTheShortID(t *testing.T) {
	t.Parallel()
	t.Run("attach", func(t *testing.T) {
		_, repo := withCLIWorkspace(t)
		id, _ := runResolveForAttachTest(t, repo)
		withAgentRunner(t, &fakeAgentRunner{probeErr: errors.New("attach must not probe Agents")})
		var full, fullDiag, short, shortDiag bytes.Buffer
		if code := runCLI(t, []string{"attach", id}, &full, &fullDiag); code != 0 || fullDiag.Len() != 0 {
			t.Fatalf("full id: exit=%d stderr=%q", code, fullDiag.String())
		}
		if code := runCLI(t, []string{"attach", roundtui.ShortRunID(id)}, &short, &shortDiag); code != 0 || shortDiag.Len() != 0 || short.String() != full.String() {
			t.Fatalf("short id: exit=%d stderr=%q stdout=%q; want %q", code, shortDiag.String(), short.String(), full.String())
		}
	})
	t.Run("runs show", func(t *testing.T) {
		_, _, run := seedUsageRun(t, false)
		for _, tc := range []struct {
			name  string
			flags []string
		}{
			{"text", nil}, {"json", []string{"--json"}},
		} {
			t.Run(tc.name, func(t *testing.T) {
				var full, fullDiag, short, shortDiag bytes.Buffer
				fullArgs := append([]string{"runs", "show", run.ID}, tc.flags...)
				shortArgs := append([]string{"runs", "show", roundtui.ShortRunID(run.ID)}, tc.flags...)
				if code := runCLI(t, fullArgs, &full, &fullDiag); code != 0 || fullDiag.Len() != 0 {
					t.Fatalf("full id: exit=%d stderr=%q", code, fullDiag.String())
				}
				if code := runCLI(t, shortArgs, &short, &shortDiag); code != 0 || shortDiag.Len() != 0 || short.String() != full.String() {
					t.Fatalf("short id: exit=%d stderr=%q stdout=%q; want %q", code, shortDiag.String(), short.String(), full.String())
				}
			})
		}
	})
}

func TestAnAmbiguousShortIDListsEveryMatch(t *testing.T) {
	t.Parallel()
	home, repo := withCLIWorkspace(t)
	runs := seedRunsForList(t, home, []runListSeed{
		{kind: store.KindResolve, state: store.StateClean, gitRoot: repo, branch: "feature/first", prNumber: "123"},
		{kind: store.KindResolve, state: store.StateClean, gitRoot: repo, branch: "feature/second", prNumber: "124"},
	})
	ids := []string{"run_20260706T120000Z_93cee2a3cccccccc", "run_20260706T120100Z_93cee2a3bbbbbbbb"}
	raw, err := sql.Open("sqlite", "file:"+store.DatabasePath(home)+"?_pragma=busy_timeout(5000)&_pragma=foreign_keys(0)")
	if err != nil {
		t.Fatal(err)
	}
	for i, run := range runs {
		result, err := raw.ExecContext(t.Context(), `UPDATE runs SET id = ? WHERE id = ?`, ids[i], run.ID)
		if err != nil {
			t.Fatal(err)
		}
		if n, err := result.RowsAffected(); err != nil || n != 1 {
			t.Fatalf("rename Run: affected=%d err=%v", n, err)
		}
	}
	if err := raw.Close(); err != nil {
		t.Fatal(err)
	}
	message := `Run id prefix "93cee2a3" matches 2 Runs: run_20260706T120100Z_93cee2a3bbbbbbbb, run_20260706T120000Z_93cee2a3cccccccc; pass more digits or the full Run id`
	for _, tc := range []struct {
		name string
		args []string
		want string
	}{
		{"attach", []string{"attach", "93cee2a3"}, "roundfix attach failed: " + message + "\n"},
		{"runs show", []string{"runs", "show", "93cee2a3"}, "roundfix: runs show failed: " + message + "\nRun 'roundfix runs show --help' for usage.\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var out, diag bytes.Buffer
			if code := runCLI(t, tc.args, &out, &diag); code != 2 || out.Len() != 0 || diag.String() != tc.want {
				t.Fatalf("exit=%d stdout=%q stderr=%q; want %q", code, out.String(), diag.String(), tc.want)
			}
		})
	}
	t.Run("invalid references stay unknown", func(t *testing.T) {
		for _, reference := range []string{"93cee2a", "93CEE2A3", "93cee2az", "ffffffff"} {
			for _, command := range [][]string{{"attach"}, {"runs", "show"}} {
				var out, diag bytes.Buffer
				args := append(append([]string{}, command...), reference)
				want := fmt.Sprintf("roundfix attach failed: Run %q does not exist; picker numbers are not stable Run ids — pass a run id or run 'roundfix attach' to pick interactively\n", reference)
				if len(command) == 2 {
					want = fmt.Sprintf("roundfix: runs show failed: Run %q does not exist\nRun 'roundfix runs show --help' for usage.\n", reference)
				}
				if code := runCLI(t, args, &out, &diag); code != 2 || out.Len() != 0 || diag.String() != want {
					t.Fatalf("args=%v exit=%d stdout=%q stderr=%q; want %q", args, code, out.String(), diag.String(), want)
				}
			}
		}
	})
}
