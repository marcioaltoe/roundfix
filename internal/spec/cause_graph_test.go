package spec

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func causeGraphFixture(t *testing.T, file string) CauseGraph {
	t.Helper()
	dir := t.TempDir()
	return CauseGraph{Tasks: map[string]string{"task_01": file}, Dir: dir}
}

func TestCauseTaskTextReadsOnlyAPlainFileInTheSpecDirectory(t *testing.T) {
	outside := filepath.Join(t.TempDir(), "secret.md")
	if err := os.WriteFile(outside, []byte("# Outside\n\n## Overview\n\nsecret\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ name, file string }{
		{"parent path", "../secret.md"},
		{"nested path", "sub/task_01.md"},
		{"absolute path", outside},
		{"directory", "sub"},
		{"dot", "."},
	} {
		t.Run(tc.name, func(t *testing.T) {
			graph := causeGraphFixture(t, tc.file)
			// Place real evidence where each unsafe name would resolve, so a
			// reader that follows the name finds it.
			for _, target := range []string{filepath.Join(graph.Dir, "..", "secret.md"), filepath.Join(graph.Dir, "sub", "task_01.md")} {
				if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(target, []byte("# Outside\n\n## Overview\n\nsecret\n"), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			text, err := graph.CauseTaskText("task_01", 4096)
			if err != nil || text != "" {
				t.Fatalf("CauseTaskText(%q) = %q, %v; want no evidence", tc.file, text, err)
			}
		})
	}
	t.Run("symbolic link", func(t *testing.T) {
		graph := causeGraphFixture(t, "task_01.md")
		if err := os.Symlink(outside, filepath.Join(graph.Dir, "task_01.md")); err != nil {
			t.Fatal(err)
		}
		text, err := graph.CauseTaskText("task_01", 4096)
		if err != nil || text != "" {
			t.Fatalf("symlinked Task = %q, %v; want no evidence", text, err)
		}
	})
	t.Run("plain file", func(t *testing.T) {
		graph := causeGraphFixture(t, "task_01.md")
		if err := os.WriteFile(filepath.Join(graph.Dir, "task_01.md"), []byte("# Task 01: Fix it\n\n## Overview\n\nThe overview.\n\n## Requirements\n\nlater\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		text, err := graph.CauseTaskText("task_01", 4096)
		if err != nil || text != "Task 01: Fix it\n\nThe overview." {
			t.Fatalf("plain Task = %q, %v", text, err)
		}
	})
}

func TestCauseTaskTextReadsABoundedPrefixOfALargeFile(t *testing.T) {
	graph := causeGraphFixture(t, "task_01.md")
	content := "# Big\n\n## Overview\n\n" + strings.Repeat("x", causeTaskReadLimit*4) + "\n"
	if err := os.WriteFile(filepath.Join(graph.Dir, "task_01.md"), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	text, err := graph.CauseTaskText("task_01", 1<<30)
	if err != nil {
		t.Fatal(err)
	}
	if len(text) > causeTaskReadLimit {
		t.Fatalf("read %d bytes of evidence, want at most the %d-byte read bound", len(text), causeTaskReadLimit)
	}
}
