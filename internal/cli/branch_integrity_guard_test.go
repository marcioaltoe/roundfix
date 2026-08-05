package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/daemon"
	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

// A failed-QA branch set carries review evidence, not implementation work.
// Branch Integrity must let watch consume that evidence without rewriting any
// ref; reclamation remains reconcile's responsibility.
func TestBranchIntegrityPreflightWatchDisregardsFourFailedQACycles(t *testing.T) {
	t.Parallel()
	const slug = "0066-run-teardown-reclaims-what-it-created"
	homeDir, repoDir := withReviewGitWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withRealReviewPreflight(t, repoDir, true)
	branches, reports := createBranchIntegrityFailedQACycles(t, homeDir, repoDir, slug, 4)
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.listPendingRunWork = runworktree.ListPendingRunWork
		dependencies.supersedingQAReport = runworktree.SupersedingQAReport
		dependencies.integratePendingRunWork = func(_ context.Context, _ string, _ string, branch string) error {
			t.Fatalf("preflight must not integrate disregarded QA branch %q", branch)
			return nil
		}
	})
	refsBefore := gitImplementOutput(t, repoDir, "for-each-ref", "--format=%(refname)=%(objectname)")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{
		"watch", "--source", "coderabbit", "--pr", "123", "--until-clean", "--max-rounds", "1", "--no-input",
	}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected watch to proceed past four failed QA cycles, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	refsAfter := gitImplementOutput(t, repoDir, "for-each-ref", "--format=%(refname)=%(objectname)")
	if refsAfter != refsBefore {
		t.Fatalf("Git refs changed during Branch Integrity Preflight\nbefore:\n%s\nafter:\n%s", refsBefore, refsAfter)
	}
	for index, branch := range branches {
		line := branchIntegrityDiagnosticLine(stderr.String(), branch)
		if line == "" {
			t.Fatalf("expected disregard diagnostic for %q, got %q", branch, stderr.String())
		}
		proof := "proof=superseded by current QA Report"
		if index == len(branches)-1 {
			proof = "proof=current QA Report"
		}
		for _, want := range []string{proof, reports[len(reports)-1], "Git ref left unchanged"} {
			if !strings.Contains(line, want) {
				t.Fatalf("expected diagnostic for %q to contain %q, got %q", branch, want, line)
			}
		}
	}
}

func branchIntegrityDiagnosticLine(diagnostic string, branch string) string {
	for _, line := range strings.Split(diagnostic, "\n") {
		if strings.Contains(line, "disregarded branch="+branch+" ") {
			return line
		}
	}
	return ""
}

func createBranchIntegrityFailedQACycles(
	t *testing.T,
	homeDir string,
	repoDir string,
	slug string,
	count int,
) ([]string, []string) {
	t.Helper()
	ctx := context.Background()
	headSHA := strings.TrimSpace(gitImplementOutput(t, repoDir, "rev-parse", "HEAD"))
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open Branch Integrity Run Database: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Fatalf("close Branch Integrity Run Database: %v", err)
		}
	}()
	branches := make([]string, 0, count)
	reports := make([]string, 0, count)
	for index := 0; index < count; index++ {
		run, err := runStore.CreateRun(ctx, store.CreateRunRequest{
			Kind:        store.KindImplement,
			GitRoot:     repoDir,
			LocalBranch: "feature/review",
			HeadSHA:     headSHA,
			SpecSlug:    slug,
		})
		if err != nil {
			t.Fatalf("create failed QA Run %d: %v", index+1, err)
		}
		ref, err := runworktree.Create(ctx, runworktree.CreateOptions{
			UserRoot: repoDir,
			Location: filepath.Join(homeDir, ".roundfix", "worktrees"),
			RunID:    run.ID,
			HeadSHA:  headSHA,
		})
		if err != nil {
			t.Fatalf("create failed QA Run Worktree %d: %v", index+1, err)
		}
		report := filepath.ToSlash(filepath.Join(
			"docs", "specs", slug, "qa", fmt.Sprintf("qa-report-2026-08-%02d.md", index+1),
		))
		mustMkdir(t, filepath.Dir(filepath.Join(ref.Path, filepath.FromSlash(report))))
		mustWrite(t, filepath.Join(ref.Path, filepath.FromSlash(report)), "---\nverdict: fail\n---\n\n# QA Report\n")
		gitImplement(t, ref.Path, "add", report)
		gitImplement(t, ref.Path, "commit", "-m", daemon.QACommitMessage(slug, "fail"))
		if _, err := runStore.SetRunWorkDir(ctx, run.ID, ref.Path); err != nil {
			t.Fatalf("record failed QA Run Worktree %d: %v", index+1, err)
		}
		if _, err := runStore.CompleteRun(ctx, run.ID, store.StateUnresolved); err != nil {
			t.Fatalf("complete failed QA Run %d: %v", index+1, err)
		}
		branches = append(branches, ref.Branch)
		reports = append(reports, report)
	}
	return branches, reports
}

// A Run Branch recorded for a different branch must not block review
// preflight on the target PR Head Branch: git topology alone cannot
// attribute Run Branches, so the guard consults the Run row.
func TestBranchIntegrityIgnoresRunBranchesOwnedByOtherBranches(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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

func TestBranchIntegrityPreflightRefusesActiveImplementRunBranch(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	active, err := runStore.CreateRun(context.Background(), store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     repoDir,
		LocalBranch: "feature/review",
		HeadSHA:     "abc123",
		SpecSlug:    "0066-run-teardown-reclaims-what-it-created",
	})
	if err != nil {
		t.Fatalf("create Active spec Run: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}
	branch := runworktree.BranchName(active.ID)
	recorder := withBranchIntegrity(t, []runworktree.PendingRunWork{{
		Branch:       branch,
		WorktreePath: filepath.Join(repoDir, "..", active.ID),
		AheadCommits: 1,
		FastForward:  true,
	}}, nil)
	withClassifyRunBranchSet(t, func(_ context.Context, _ string, _ string, _ string, _ []store.Run) (runworktree.BranchSetClassification, error) {
		return runworktree.BranchSetClassification{
			Current:       branch,
			CurrentReport: "docs/specs/0066-run-teardown-reclaims-what-it-created/qa/qa-report-2026-08-04.md",
		}, nil
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"fetch", "--source", "coderabbit", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("expected Active spec Run Branch refusal exit %d, got %d stderr=%q", exitPreflight, code, stderr.String())
	}
	if len(recorder.integrations) != 0 {
		t.Fatalf("Active spec Run Branch must not be integrated, got %#v", recorder.integrations)
	}
	if !strings.Contains(stderr.String(), branch) {
		t.Fatalf("expected refusal to name Active Run Branch %q, got %q", branch, stderr.String())
	}
}

// The clean tracked checkout validation must refuse before any
// fast-forward auto-integration mutates the user's branch.
func TestReviewCleanTreeRefusalPrecedesAutoIntegration(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
