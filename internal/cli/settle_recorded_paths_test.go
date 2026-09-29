package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"roundfix/internal/spec"
)

func TestSettleRecordsTheUndeclaredPathsOfItsTaskCommit(t *testing.T) {
	t.Parallel()
	_, repoDir := newImplementWorkspace(t, []implementSeed{{
		id:           "task_01",
		title:        "Recover delivery storage",
		status:       string(spec.StatusFailed),
		verification: []string{"test -f internal/store/delivery_test.go"},
	}})
	mustMkdir(t, filepath.Join(repoDir, "internal", "store"))
	mustWrite(t, filepath.Join(repoDir, "internal", "store", "delivery_test.go"), "package store\n")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("settle exit = %d, want %d; stdout=%q stderr=%q", code, exitOK, stdout.String(), stderr.String())
	}
	content, err := os.ReadFile(implementTaskPath(repoDir, "task_01"))
	if err != nil {
		t.Fatalf("read settled Task: %v", err)
	}
	want := []string{"internal/store/delivery_test.go"}
	if got := spec.RecordedTaskPaths(content); !slices.Equal(got, want) {
		t.Fatalf("settled Task recorded paths = %v, want %v", got, want)
	}
}

func TestSettleRecordsNothingWhenEveryPathIsDeclared(t *testing.T) {
	t.Parallel()
	_, repoDir := newImplementWorkspace(t, []implementSeed{{
		id:           "task_01",
		title:        "Recover delivery storage",
		status:       string(spec.StatusFailed),
		verification: []string{"test -f internal/store/delivery.go"},
	}})
	taskPath := implementTaskPath(repoDir, "task_01")
	content := mustRead(t, taskPath)
	contextSection := "## Context\n\n- interface: `internal/store/delivery.go`\n\n"
	content = strings.Replace(content, "## Verification\n", contextSection+"## Verification\n", 1)
	if err := os.WriteFile(taskPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write declared Task Context: %v", err)
	}
	if err := spec.RecordTaskPaths(taskPath, []string{"old/attempt.go"}); err != nil {
		t.Fatalf("seed old recorded paths: %v", err)
	}
	gitImplement(t, repoDir, "add", "-A")
	gitImplement(t, repoDir, "commit", "-m", "seed declared Task Context")
	mustMkdir(t, filepath.Join(repoDir, "internal", "store"))
	mustWrite(t, filepath.Join(repoDir, "internal", "store", "delivery.go"), "package store\n")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_01"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("settle exit = %d, want %d; stdout=%q stderr=%q", code, exitOK, stdout.String(), stderr.String())
	}
	settled, err := os.ReadFile(taskPath)
	if err != nil {
		t.Fatalf("read settled Task: %v", err)
	}
	if got := spec.RecordedTaskPaths(settled); len(got) != 0 {
		t.Fatalf("settled Task recorded paths = %v, want none", got)
	}
	if strings.Contains(string(settled), spec.RecordedPathsHeading) {
		t.Fatalf("settled Task kept an empty recorded-path section:\n%s", settled)
	}
}
