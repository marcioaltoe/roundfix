package spec

// Suite: Recorded paths section ownership
// Invariant: Only an exact unfenced Recorded paths heading followed by nothing but Daemon-written settlement sections is Daemon-owned.
// Boundary IN: Recorded path section reading and replacement.
// Boundary OUT: Task commit staging and Markdown parsing outside this section.

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestRecordedPathsKeepAFencedHeadingInAuthoredText(t *testing.T) {
	t.Parallel()
	original := []byte("# Task\n\n## Result\n\n```markdown\n## Recorded paths\n\n- `authored.go`\n```\n")
	if got := RecordedTaskPaths(original); len(got) != 0 {
		t.Fatalf("RecordedTaskPaths() = %v, want no paths from a fenced heading", got)
	}
	taskPath := writeRecordedPathsSectionFixture(t, original)

	if err := RecordTaskPaths(taskPath, []string{"internal/new.go"}); err != nil {
		t.Fatalf("RecordTaskPaths: %v", err)
	}

	updated := readRecordedPathsSectionFixture(t, taskPath)
	if !bytes.HasPrefix(updated, original) {
		t.Fatalf("RecordTaskPaths changed fenced authored text:\n%s", updated)
	}
	if got, want := RecordedTaskPaths(updated), []string{"internal/new.go"}; !slices.Equal(got, want) {
		t.Fatalf("RecordedTaskPaths() = %v, want %v", got, want)
	}
}

func TestRecordedPathsKeepAnAuthoredHeadingFollowedBySections(t *testing.T) {
	t.Parallel()
	original := []byte("# Task\n\n## Recorded paths\n\nAuthored explanation.\n\n## Notes\n\nKeep these notes byte-for-byte.\n")
	if got := RecordedTaskPaths(original); len(got) != 0 {
		t.Fatalf("RecordedTaskPaths() = %v, want no paths from a non-final heading", got)
	}
	taskPath := writeRecordedPathsSectionFixture(t, original)

	if err := RecordTaskPaths(taskPath, []string{"internal/new.go"}); err != nil {
		t.Fatalf("RecordTaskPaths: %v", err)
	}

	updated := readRecordedPathsSectionFixture(t, taskPath)
	if !bytes.HasPrefix(updated, original) {
		t.Fatalf("RecordTaskPaths changed authored sections:\n%s", updated)
	}
	if got, want := RecordedTaskPaths(updated), []string{"internal/new.go"}; !slices.Equal(got, want) {
		t.Fatalf("RecordedTaskPaths() = %v, want %v", got, want)
	}
}

func TestRecordedPathsReplaceTheTrailingDaemonSection(t *testing.T) {
	t.Parallel()
	base := []byte("# Task\n\n## Result\n\nKeep this result without a final newline")
	original := append(append([]byte(nil), base...), []byte("\n## Recorded paths\n\nold prose\n\n- `old.go`\n")...)
	taskPath := writeRecordedPathsSectionFixture(t, original)

	if err := RecordTaskPaths(taskPath, []string{"internal/new.go"}); err != nil {
		t.Fatalf("RecordTaskPaths: %v", err)
	}

	updated := readRecordedPathsSectionFixture(t, taskPath)
	if !bytes.HasPrefix(updated, base) {
		t.Fatalf("RecordTaskPaths changed bytes before the trailing section:\n%s", updated)
	}
	if bytes.Contains(updated, []byte("old prose")) || bytes.Contains(updated, []byte("old.go")) {
		t.Fatalf("RecordTaskPaths kept bytes from the replaced section:\n%s", updated)
	}
	if got := bytes.Count(updated, []byte("\n"+RecordedPathsHeading+"\n")); got != 1 {
		t.Fatalf("recorded heading count = %d, want 1:\n%s", got, updated)
	}
	if got, want := RecordedTaskPaths(updated), []string{"internal/new.go"}; !slices.Equal(got, want) {
		t.Fatalf("RecordedTaskPaths() = %v, want %v", got, want)
	}
}

func writeRecordedPathsSectionFixture(t *testing.T, content []byte) string {
	t.Helper()
	taskPath := filepath.Join(t.TempDir(), "task_01.md")
	if err := os.WriteFile(taskPath, content, 0o644); err != nil {
		t.Fatalf("write Task fixture: %v", err)
	}
	return taskPath
}

func readRecordedPathsSectionFixture(t *testing.T, taskPath string) []byte {
	t.Helper()
	content, err := os.ReadFile(taskPath)
	if err != nil {
		t.Fatalf("read Task fixture: %v", err)
	}
	return content
}

func TestRecordedPathsAreReadBeforeTrailingSettlementSections(t *testing.T) {
	t.Parallel()
	const base = "---\ntask: task_01\nstatus: completed\ntype: backend\n---\n\n# Task\n\n## Context\n\n- interface: declared.go\n"
	const recorded = "\n## Recorded paths\n\n- `z.go`\n- `z.go`\n"
	const provenance = "\n## Carry-forward provenance\n\n- Source Run: synthetic-run\n- Source commit: abc123\n"
	for _, tc := range []struct {
		name, sections string
		want           []string
	}{
		{"carry-forward-provenance", recorded + provenance, []string{"z.go"}},
		{"result-then-provenance", recorded + "\n## Result\n\n- `not-recorded.go`\n" + provenance, []string{"z.go"}},
		{"duplicate-sections", recorded + provenance + "\n## Recorded paths\n\n- `a.go`\n- `z.go`\n", []string{"z.go", "a.go"}},
		{"authored-section-between", recorded + "\n## Notes\n\nAuthored notes.\n" + provenance, nil},
		{"fenced-heading-before-provenance", "\n## Result\n\n```markdown\n## Recorded paths\n\n- `z.go`\n```\n" + provenance, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			content := []byte(base + tc.sections)
			if got := RecordedTaskPaths(content); !slices.Equal(got, tc.want) {
				t.Fatalf("RecordedTaskPaths() = %v, want %v", got, tc.want)
			}
			declared, err := DeclaredTaskPaths(content, "task_01.md")
			if err != nil {
				t.Fatal(err)
			}
			wantDeclared := append([]string{"declared.go"}, tc.want...)
			slices.Sort(wantDeclared)
			if !slices.Equal(declared, wantDeclared) {
				t.Fatalf("DeclaredTaskPaths() = %v, want %v", declared, wantDeclared)
			}
		})
	}
}

func TestRecordTaskPathsReplacesTheSectionBeforeCarryForwardProvenance(t *testing.T) {
	t.Parallel()
	// Hazard: settlement rewrites committed Task files. Dropping provenance or
	// appending a duplicate record must fail the exact-byte rows below.
	const base = "# Task\n\n## Result\n\nKeep result bytes without a final newline"
	const explanation = "The Daemon recorded these paths, which this Task changed without declaring them in `## Context`."
	const oldSection = "\n## Recorded paths\n\n" + explanation + "\n\n- `old.go`\n"
	const newSection = "\n## Recorded paths\n\n" + explanation + "\n\n- `new.go`\n"
	const provenance = "\n## Carry-forward provenance\n\n- Source Run: synthetic-run\n- Source commit: abc123\n\nKeep trailing bytes."
	for _, tc := range []struct {
		name, original, want string
		paths                []string
	}{
		{"same-paths", base + oldSection + provenance, base + oldSection + provenance, []string{"old.go"}},
		{"new-paths", base + oldSection + provenance, base + newSection + provenance, []string{"new.go"}},
		{"no-paths", base + oldSection + provenance, base + provenance, nil},
		{"duplicate-sections", base + oldSection + provenance + oldSection, base + newSection + provenance, []string{"new.go"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			taskPath := writeRecordedPathsSectionFixture(t, []byte(tc.original))
			before, err := os.Stat(taskPath)
			if err != nil {
				t.Fatal(err)
			}
			if err := RecordTaskPaths(taskPath, tc.paths); err != nil {
				t.Fatal(err)
			}
			updated := readRecordedPathsSectionFixture(t, taskPath)
			if !bytes.Equal(updated, []byte(tc.want)) {
				t.Fatalf("RecordTaskPaths() wrote:\n%s\nwant exact bytes:\n%s", updated, tc.want)
			}
			if !bytes.Contains(updated, []byte(provenance)) {
				t.Fatalf("settlement lost provenance bytes:\n%s", updated)
			}
			if tc.name == "same-paths" {
				after, err := os.Stat(taskPath)
				if err != nil {
					t.Fatal(err)
				}
				if !os.SameFile(before, after) {
					t.Fatal("same paths replaced the Task file")
				}
			}
		})
	}
}
