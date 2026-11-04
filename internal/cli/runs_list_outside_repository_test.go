// Suite: Runs List outside a repository
// Invariant: no Git root lists every repository, preserving filters, bounds, and streams.
// Boundary IN: CLI entry point, temporary homes and repositories, and Run Database.
// Boundary OUT: Run execution and terminal rendering, covered by their owning suites.
package cli

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunRunsListOutsideRepositoryListsEveryRepository(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := withCLIWorkspace(t)
	otherRepo := filepath.Join(t.TempDir(), "other-repo")
	mustMkdir(t, filepath.Join(otherRepo, ".git"))
	runs := seedRunsForListColumns(t, homeDir, repoDir, otherRepo)
	outsideDir := t.TempDir()

	for _, test := range []struct {
		name  string
		flags []string
		rows  []int
		note  string
	}{
		{name: "default", rows: []int{2, 1}, note: "(1 terminal Run(s) hidden; use --state all)\n"},
		{name: "all states", flags: []string{"--state", "all"}, rows: []int{2, 1, 0}},
		{name: "limit one", flags: []string{"--limit", "1"}, rows: []int{2}, note: "(1 older Run(s) hidden; use --limit 0)\n"},
		{name: "terminal", flags: []string{"--state", "terminal"}, rows: []int{0}, note: "(2 active Run(s) hidden; use --state all)\n"},
		{name: "unbounded", flags: []string{"--limit", "0"}, rows: []int{2, 1}, note: "(1 terminal Run(s) hidden; use --state all)\n"},
		{name: "all states limited", flags: []string{"--state", "all", "--limit", "1"}, rows: []int{2}, note: "(2 older Run(s) hidden; use --limit 0)\n"},
	} {
		t.Run(test.name, func(t *testing.T) {
			args := append([]string{"runs", "list"}, test.flags...)
			setCommandEnvironmentForTest(t, homeDir, repoDir)
			var all, allErr bytes.Buffer
			allCode := runCLI(t, append(append([]string{}, args...), "--all"), &all, &allErr)
			if allCode != exitOK || allErr.String() != test.note {
				t.Fatalf("explicit --all: exit=%d stderr=%q; want exit=0 stderr=%q", allCode, allErr.String(), test.note)
			}
			rows := strings.Split(strings.TrimSuffix(all.String(), "\n"), "\n")
			if len(rows) != len(test.rows) {
				t.Fatalf("explicit --all: got %d rows, want %d: %q", len(rows), len(test.rows), all.String())
			}
			for i, index := range test.rows {
				if !strings.HasPrefix(rows[i], runs[index].ID+"  ") || !strings.HasSuffix(rows[i], "  "+runs[index].GitRoot) {
					t.Fatalf("row %d must name Run %s and repository %s: %q", i, runs[index].ID, runs[index].GitRoot, rows[i])
				}
			}
			setCommandEnvironmentForTest(t, homeDir, outsideDir)
			for _, outsideArgs := range [][]string{args, append(append([]string{}, args...), "--all")} {
				var outside, outsideErr bytes.Buffer
				code := runCLI(t, outsideArgs, &outside, &outsideErr)
				if code != allCode || outside.String() != all.String() || outsideErr.String() != allErr.String() {
					t.Fatalf("outside %v: exit=%d stdout=%q stderr=%q; want exit=%d stdout=%q stderr=%q",
						outsideArgs, code, outside.String(), outsideErr.String(), allCode, all.String(), allErr.String())
				}
			}
		})
	}

	t.Run("inside keeps repository scope", func(t *testing.T) {
		setCommandEnvironmentForTest(t, homeDir, repoDir)
		var stdout, stderr bytes.Buffer
		code := runCLI(t, []string{"runs", "list", "--state", "all"}, &stdout, &stderr)
		if code != exitOK || stderr.Len() != 0 {
			t.Fatalf("inside: exit=%d stderr=%q", code, stderr.String())
		}
		rows := strings.Split(strings.TrimSuffix(stdout.String(), "\n"), "\n")
		if len(rows) != 2 || !strings.HasPrefix(rows[0], runs[1].ID+"  ") || !strings.HasPrefix(rows[1], runs[0].ID+"  ") || strings.Contains(stdout.String(), runs[2].ID) || strings.Contains(stdout.String(), "  "+repoDir+"\n") {
			t.Fatalf("expected only this repository's two Runs without a repository column, got %q", stdout.String())
		}
	})
}

func TestRunsListOutsideARepositoryTranscript(t *testing.T) {
	t.Parallel()
	setCommandEnvironmentForTest(t, t.TempDir(), t.TempDir())
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"runs", "list"}, &stdout, &stderr)
	if code != 0 || stdout.String() != "No Runs found.\n" || stderr.String() != "" {
		t.Fatalf("Surface Transcript 1: exit=%d stdout=%q stderr=%q; want exit=0 stdout=%q stderr=%q",
			code, stdout.String(), stderr.String(), "No Runs found.\n", "")
	}
}
