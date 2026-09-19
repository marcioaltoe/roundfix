package cli

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"roundfix/internal/spec"
)

// Suite: reopen command seam
// Invariant: reopen mutates only a completed terminal QA gate whose dependency closure is stale.
// Boundary IN: public CLI dispatch, configured Spec Root, task status and invalidation record.
// Boundary OUT: Daemon execution and publication, which reopen must never start.
func TestReopenStaleGatePreservesEvidence(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Correct the finding", status: string(spec.StatusPending)},
		implementQAGateSeed(string(spec.StatusCompleted), "task_01"),
	})
	taskPath := implementTaskPath(repoDir, "task_qa")
	const priorResult = "## Result\n\nThe prior gate verified the original graph.\n"
	mustWrite(t, taskPath, mustRead(t, taskPath)+"\n"+priorResult)
	reportPath := filepath.Join(repoDir, "docs", "specs", implementTestSlug, "qa", "qa-report-2026-09-18.md")
	const report = "---\nverdict: pass\n---\n\n# QA Report\n\nPrior evidence.\n"
	mustMkdir(t, filepath.Dir(reportPath))
	mustWrite(t, reportPath, report)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runCLIContext(t, context.Background(), []string{"reopen", "--spec", implementTestSlug}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("reopen exit = %d, want %d; stderr=%q stdout=%q", code, exitOK, stderr.String(), stdout.String())
	}
	if stderr.String() != "" {
		t.Fatalf("reopen stderr = %q, want empty", stderr.String())
	}
	if !strings.Contains(stdout.String(), "reopened task_qa pending") {
		t.Fatalf("reopen stdout = %q, want reopened QA Task", stdout.String())
	}
	if got := mustRead(t, reportPath); got != report {
		t.Fatalf("QA Report changed:\n%s", got)
	}
	gotTask := mustRead(t, taskPath)
	for _, want := range []string{
		"status: pending",
		priorResult,
		"## Invalidation",
		"`qa/qa-report-2026-09-18.md`",
		"`task_01`",
	} {
		if !strings.Contains(gotTask, want) {
			t.Fatalf("QA Task does not contain %q after reopen:\n%s", want, gotTask)
		}
	}
	if _, err := spec.Load(filepath.Join(repoDir, "docs", "specs"), implementTestSlug); err != nil {
		t.Fatalf("spec.Load after reopen: %v", err)
	}
	assertNoRunDatabase(t, homeDir)
}

func TestReopenRefusesPendingQATaskWithoutMutation(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", status: string(spec.StatusCompleted)},
		implementQAGateSeed(string(spec.StatusPending), "task_01"),
	})
	specDir := filepath.Join(repoDir, "docs", "specs", implementTestSlug)
	before := reopenDirectoryDigest(t, specDir)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runCLIContext(t, context.Background(), []string{"reopen", "--spec", implementTestSlug}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("reopen exit = %d, want %d; stderr=%q", code, exitPreflight, stderr.String())
	}
	if stdout.String() != "" {
		t.Fatalf("refusal stdout = %q, want empty", stdout.String())
	}
	for _, want := range []string{"QA Task", "not completed", "pending"} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("refusal stderr = %q, want %q", stderr.String(), want)
		}
	}
	if after := reopenDirectoryDigest(t, specDir); !slices.Equal(before, after) {
		t.Fatalf("Spec directory changed on pending-gate refusal\nbefore=%x\nafter=%x", before, after)
	}
	assertNoRunDatabase(t, homeDir)
}

func TestReopenRefusesHealthyCompletedGateWithoutMutation(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", status: string(spec.StatusCompleted)},
		implementQAGateSeed(string(spec.StatusCompleted), "task_01"),
	})
	specDir := filepath.Join(repoDir, "docs", "specs", implementTestSlug)
	before := reopenDirectoryDigest(t, specDir)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runCLIContext(t, context.Background(), []string{"reopen", "--spec", implementTestSlug}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("reopen exit = %d, want %d; stderr=%q", code, exitPreflight, stderr.String())
	}
	if stdout.String() != "" {
		t.Fatalf("refusal stdout = %q, want empty", stdout.String())
	}
	for _, want := range []string{"not stale", "every dependency is completed"} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("refusal stderr = %q, want %q", stderr.String(), want)
		}
	}
	if after := reopenDirectoryDigest(t, specDir); !slices.Equal(before, after) {
		t.Fatalf("Spec directory changed on healthy-gate refusal\nbefore=%x\nafter=%x", before, after)
	}
	assertNoRunDatabase(t, homeDir)
}

func TestReopenRefusesSpecWithoutTerminalQATask(t *testing.T) {
	t.Parallel()
	homeDir, _ := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", status: string(spec.StatusCompleted)},
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"reopen", "--spec", implementTestSlug}, &stdout, &stderr)

	if code != exitPreflight || !strings.Contains(stderr.String(), "no terminal QA Task") {
		t.Fatalf("reopen exit = %d stderr=%q, want no-terminal-QA refusal", code, stderr.String())
	}
	assertNoRunDatabase(t, homeDir)
}

func TestReopenRefusesUnknownFlag(t *testing.T) {
	t.Parallel()
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"reopen", "--unknown"}, &stdout, &stderr)

	if code != exitPreflight || !strings.Contains(stderr.String(), "flag provided but not defined") {
		t.Fatalf("reopen exit = %d stderr=%q, want unknown-flag refusal", code, stderr.String())
	}
}

func reopenDirectoryDigest(t *testing.T, root string) []byte {
	t.Helper()
	hash := sha256.New()
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintf(hash, "%s\x00%d\x00", filepath.ToSlash(rel), len(content)); err != nil {
			return err
		}
		_, err = hash.Write(content)
		return err
	})
	if err != nil {
		t.Fatalf("digest Spec directory: %v", err)
	}
	return hash.Sum(nil)
}
