package spec

// Suite: Recorded paths section ownership
// Invariant: Only an exact unfenced Recorded paths heading that starts the final H2 section is Daemon-owned.
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
