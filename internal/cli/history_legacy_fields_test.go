package cli

// Suite: Legacy creation dates through both conversion commands.
// Invariant: records keep PRD, TechSpec, then creation-commit date precedence across promotion re-plans.
// Boundary IN: disposable Git repositories, temporary homes, and public command arguments.
// Boundary OUT: written Archive Records and promoted references.

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/spec"
)

func historyCreationFixture(t *testing.T, h historyFixture, prdDate, techDate string) string {
	t.Helper()
	active := "docs/specs/aaa"
	archived := "docs/history/specs/aaa"
	metadata := "spec: aaa\n"
	if prdDate != "" {
		metadata += "created: " + prdDate + "\n"
	}
	historyWrite(t, filepath.Join(h.repo, active, "_prd.md"), "---\n"+metadata+"---\n\n# Demo\n\nPreserved outcome.\n")
	if techDate != "" {
		historyWrite(t, filepath.Join(h.repo, active, "_techspec.md"), "---\ncreated: "+techDate+"\n---\n\n# Design\n")
	}
	historyWrite(t, filepath.Join(h.repo, active, "references/lesson.md"), "---\ntype: fix\nstatus: promoted\ncreated: 2026-01-01\n---\n\n# Lesson\n\nReusable knowledge.\n")
	gittest.Run(t, h.repo, "add", "--all")
	gittest.Run(t, h.repo, "commit", "--date=2026-01-02T12:00:00+00:00", "-m", "Create Spec")
	created := strings.TrimSpace(gittest.Run(t, h.repo, "show", "-s", "--format=%as", "HEAD"))
	if err := os.MkdirAll(filepath.Join(h.repo, "docs/history/specs"), 0o755); err != nil {
		t.Fatal(err)
	}
	gittest.Run(t, h.repo, "mv", active, archived)
	gittest.Run(t, h.repo, "commit", "--date=2026-02-03T12:00:00+00:00", "-m", "Archive Spec (#42)")
	return created
}

func assertHistoryCreationRecord(t *testing.T, h historyFixture, want string, promoted bool) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(h.repo, "docs/history/specs/aaa.md"))
	if err != nil {
		t.Fatal(err)
	}
	record, err := spec.ParseArchiveRecord(data)
	if err != nil || record.Created != want || !reflect.DeepEqual(record.Sources, []string{"lesson.md"}) {
		t.Fatalf("record = %+v, err = %v, want created %q and lesson source", record, err, want)
	}
	if promoted {
		if !reflect.DeepEqual(record.Promoted, []string{"docs/references/lesson.md"}) {
			t.Fatalf("promoted = %v", record.Promoted)
		}
		if _, err := os.Stat(filepath.Join(h.repo, "docs/references/lesson.md")); err != nil {
			t.Fatal(err)
		}
	}
	assertPathMissing(t, filepath.Join(h.repo, "docs/history/specs/aaa"))
}

func TestHistorySanitizeRecordsTheCreationCommitDate(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, prd, tech string
		promote         bool
	}{
		{name: "creation commit"},
		{name: "creation commit with promotion", promote: true},
		{name: "TechSpec precedes Git", tech: "2026-01-03T12:00:00Z", promote: true},
		{name: "PRD precedes TechSpec", prd: "2026-01-04", tech: "2026-01-03", promote: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			h := newHistoryFixture(t, false)
			want := historyCreationFixture(t, h, tc.prd, tc.tech)
			if tc.tech != "" {
				want = tc.tech[:10]
			}
			if tc.prd != "" {
				want = tc.prd
			}
			historyTag(t, h)
			args := []string{"sanitize", "--apply", "--batch", "1"}
			if tc.promote {
				args = append(args, "--promote", "docs/history/specs/aaa/references/lesson.md")
			}
			code, out, stderr := historyRun(t, h, args...)
			if code != 0 || stderr != "" || !strings.Contains(out, "wrote 1 Archive Record(s)") {
				t.Fatalf("exit = %d, stdout = %s, stderr = %s", code, out, stderr)
			}
			assertHistoryCreationRecord(t, h, want, tc.promote)
		})
	}
}

func TestBaselineUpdateHistoryRecordsTheCreationCommitDate(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, prd, tech string }{
		{name: "creation commit"},
		{name: "TechSpec precedes Git", tech: "2026-01-03T12:00:00Z"},
		{name: "PRD precedes TechSpec", prd: "2026-01-04", tech: "2026-01-03"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			h := baselineCitationFixture(t, false)
			want := historyCreationFixture(t, h, tc.prd, tc.tech)
			if tc.tech != "" {
				want = tc.tech[:10]
			}
			if tc.prd != "" {
				want = tc.prd
			}
			result, out, stderr, code := baselineHistoryRun(t, h, "json", "--yes", "--promote", "docs/history/specs/aaa/references/lesson.md")
			if code != 0 || stderr != "" || result.History.Status != "applied" {
				t.Fatalf("exit = %d, stdout = %s, stderr = %s", code, out, stderr)
			}
			assertHistoryCreationRecord(t, h, want, true)
		})
	}
}
