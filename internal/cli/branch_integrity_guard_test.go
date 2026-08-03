package cli

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

// A Run Branch recorded for a different branch must not block review
// preflight on the target PR Head Branch: git topology alone cannot
// attribute Run Branches, so the guard consults the Run row.
func TestBranchIntegrityIgnoresRunBranchesOwnedByOtherBranches(t *testing.T) {
	// Sequential: overrides package-level test seams.
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	other, err := runStore.CreateRun(context.Background(), store.CreateRunRequest{
		Kind:           store.KindWatch,
		HeadRepository: "owner/project",
		HeadBranch:     "other-feature",
		BaseRepository: "owner/project",
		PRNumber:       "77",
		GitRoot:        repoDir,
		LocalBranch:    "other-feature",
		HeadSHA:        "def456",
		ArtifactDir:    filepath.Join(repoDir, ".roundfix"),
		Agent:          "codex",
	})
	if err != nil {
		t.Fatalf("create other-branch Run: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}
	pending := []runworktree.PendingRunWork{{
		Branch:       "roundfix/run-" + other.ID,
		WorktreePath: filepath.Join(repoDir, "..", "run-other"),
		AheadCommits: 4,
		FastForward:  false,
	}}
	recorder := withBranchIntegrity(t, pending, nil)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"fetch", "--source", "coderabbit", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected fetch to ignore the other branch's Run Branch, got %d stderr=%q", code, stderr.String())
	}
	if strings.Contains(stderr.String(), "Branch Integrity Preflight refused") {
		t.Fatalf("expected no Branch Integrity refusal, got %q", stderr.String())
	}
	if len(recorder.integrations) != 0 {
		t.Fatalf("expected no integration of the other branch's work, got %#v", recorder.integrations)
	}
	if !strings.Contains(stdout.String(), "Fetch complete") {
		t.Fatalf("expected fetch success, got %q", stdout.String())
	}
}

// A Run created through the Branch Integrity bypass holds no Active Run
// lock; the guard must still discover it through the runs table so a later
// normal Run cannot start concurrently in the same checkout.
func TestBranchIntegrityRejectsLocklessBypassRun(t *testing.T) {
	// Sequential: overrides package-level test seams.
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withBranchIntegrity(t, nil, nil)
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	lockless, err := runStore.CreateRunSkippingActiveLock(context.Background(), store.CreateRunRequest{
		Kind:           store.KindResolve,
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
		BaseRepository: "owner/project",
		PRNumber:       "123",
		GitRoot:        repoDir,
		LocalBranch:    "feature/review",
		HeadSHA:        "abc123",
		ArtifactDir:    filepath.Join(repoDir, ".roundfix"),
		Agent:          "codex",
		OwnerPID:       os.Getpid(),
	})
	if err != nil {
		t.Fatalf("create lockless bypass Run: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"fetch", "--source", "coderabbit", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("expected lockless Active Run refusal exit 2, got %d stderr=%q", code, stderr.String())
	}
	for _, want := range []string{
		lockless.ID,
		"roundfix stop " + lockless.ID,
		"did not create a Run",
	} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("expected stderr to contain %q, got %q", want, stderr.String())
		}
	}
	assertRunCount(t, store.DatabasePath(homeDir), 1)
}

// The clean tracked checkout validation must refuse before any
// fast-forward auto-integration mutates the user's branch.
func TestReviewCleanTreeRefusalPrecedesAutoIntegration(t *testing.T) {
	// Sequential: overrides package-level test seams.
	_, repoDir := withReviewGitWorkspace(t)
	withRealReviewPreflight(t, repoDir, true)
	pending := []runworktree.PendingRunWork{{
		Branch:       "roundfix/run-ready",
		WorktreePath: filepath.Join(repoDir, "..", "run-ready"),
		AheadCommits: 1,
		FastForward:  true,
	}}
	recorder := withBranchIntegrity(t, pending, nil)
	mustWrite(t, filepath.Join(repoDir, "README.md"), "dirty tracked change\n")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), branchIntegrityCommandArgs("resolve"), &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("expected dirty tracked refusal exit 2, got %d stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "README.md") {
		t.Fatalf("expected dirty tracked path in stderr, got %q", stderr.String())
	}
	if len(recorder.integrations) != 0 {
		t.Fatalf("expected no auto-integration before the clean-tree refusal, got %#v", recorder.integrations)
	}
}

// Profile proof runs before the bypass audit and every Run mutation boundary:
// a proof failure publishes no audit and creates no Run, while an audit publish
// failure happens only after the configured review profile has been proven.
func TestBranchIntegrityBypassAuditFollowsProfileProof(t *testing.T) {
	// Sequential: overrides package-level test seams.
	t.Run("proof failure prevents audit and Run creation", func(t *testing.T) {
		homeDir, repoDir := withReviewGitWorkspace(t)
		persistCLIReviewIssue(t, repoDir, 1, "feature/review")
		withRealReviewPreflight(t, repoDir, true)
		withBranchIntegrity(t, nil, nil)
		runner := &fakeAgentRunner{probeErr: errors.New("model rejected")}
		withAgentRunner(t, runner)
		comments := withPullRequestComments(t, nil)
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := runCLIContext(t, context.Background(), []string{"resolve", "--pr", "123", "--skip-branch-integrity", "--no-input"}, &stdout, &stderr)

		if code != exitPreflight {
			t.Fatalf("expected probe failure exit 2, got %d stderr=%q", code, stderr.String())
		}
		if len(comments.calls) != 0 {
			t.Fatalf("expected profile proof failure before bypass audit, got %#v", comments.calls)
		}
		if len(runner.probeRequests) == 0 {
			t.Fatalf("expected profile proof to run before the audit")
		}
		assertNoRunDatabase(t, homeDir)
	})

	t.Run("audit publish failure follows successful proof", func(t *testing.T) {
		_, repoDir := withReviewGitWorkspace(t)
		persistCLIReviewIssue(t, repoDir, 1, "feature/review")
		withRealReviewPreflight(t, repoDir, true)
		withBranchIntegrity(t, nil, nil)
		runner := &fakeAgentRunner{}
		withAgentRunner(t, runner)
		withPullRequestComments(t, errors.New("comment denied"))
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := runCLIContext(t, context.Background(), []string{"resolve", "--pr", "123", "--skip-branch-integrity", "--no-input"}, &stdout, &stderr)

		if code != exitPreflight {
			t.Fatalf("expected audit publish failure exit 2, got %d stderr=%q", code, stderr.String())
		}
		if len(runner.probeRequests) != 2 {
			t.Fatalf("expected successful review profile preferred and fallback proof before audit failure, got %#v", runner.probeRequests)
		}
		if runner.calls != 0 {
			t.Fatalf("audit failure must prevent durable Agent work, got %d call(s)", runner.calls)
		}
	})
}
