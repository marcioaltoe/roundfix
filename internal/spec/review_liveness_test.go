package spec

// Suite: orphan Review Artifact liveness
// Invariant: only the recorded outcome may retire an orphan Review Artifact.
// Boundary IN: outcome.md at the Review Artifact root and isolated Git repositories.
// Boundary OUT: hosting-provider state, credentials, network access, and history relocation.

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/rounds"
)

func TestClassifyReviewIgnoresObjectStoreAvailability(t *testing.T) {
	t.Parallel()

	presentRepo, _ := reviewLivenessRepo(t, "main")
	head := reviewLivenessBranchCommit(t, presentRepo, "feature/reviewed")
	reviewDir := filepath.Join(t.TempDir(), "pr-123")
	persistReviewRound(t, reviewDir, 1, "feature/reviewed", head)
	writeReviewOutcome(t, reviewDir, "open", "")

	absentRepo, _ := reviewLivenessRepo(t, "main")
	fetchedRepo, _ := reviewLivenessRepo(t, "main")
	gittest.Run(t, fetchedRepo, "fetch", presentRepo, head+":refs/remotes/origin/reviewed")

	repositories := []struct {
		name string
		root string
	}{
		{name: "recorded head present on a local branch", root: presentRepo},
		{name: "recorded head absent", root: absentRepo},
		{name: "recorded head reachable only from a fetched ref", root: fetchedRepo},
	}

	var firstReason string
	for _, repository := range repositories {
		t.Run(repository.name, func(t *testing.T) {
			got, reason, err := ClassifyReview(context.Background(), repository.root, reviewDir)
			if err != nil {
				t.Fatalf("ClassifyReview() error = %v", err)
			}
			if got != ReviewLive {
				t.Fatalf("ClassifyReview() = %q, want %q; reason: %s", got, ReviewLive, reason)
			}
			if firstReason == "" {
				firstReason = reason
			}
			if reason != firstReason {
				t.Fatalf("ClassifyReview() reason = %q, want stable recorded-outcome reason %q", reason, firstReason)
			}
		})
	}
}

func TestClassifyReviewAcceptsARecordedSquashReceipt(t *testing.T) {
	t.Parallel()

	repo, head := reviewLivenessRepo(t, "main")
	reviewDir := filepath.Join(t.TempDir(), "pr-123")
	persistReviewRound(t, reviewDir, 1, "feature/squashed", head)
	const absentMergeCommit = "0123456789abcdef0123456789abcdef01234567"
	writeReviewOutcome(t, reviewDir, "merged", absentMergeCommit)

	got, reason, err := ClassifyReview(context.Background(), repo, reviewDir)
	if err != nil {
		t.Fatalf("ClassifyReview() error = %v", err)
	}
	if got != ReviewFinished {
		t.Fatalf("ClassifyReview() = %q, want %q; reason: %s", got, ReviewFinished, reason)
	}
	if !strings.Contains(reason, absentMergeCommit) {
		t.Fatalf("ClassifyReview() reason = %q, want recorded merge commit", reason)
	}
}

func TestClassifyReviewWithoutARecordedOutcomeIsUnknown(t *testing.T) {
	t.Parallel()

	repo, head := reviewLivenessRepo(t, "main")
	reviewDir := filepath.Join(t.TempDir(), "pr-123")
	persistReviewRound(t, reviewDir, 1, "feature/no-outcome", head)

	got, reason, err := ClassifyReview(context.Background(), repo, reviewDir)
	if err != nil {
		t.Fatalf("ClassifyReview() error = %v", err)
	}
	if got != ReviewUndecidable {
		t.Fatalf("ClassifyReview() = %q, want %q; reason: %s", got, ReviewUndecidable, reason)
	}
	if !strings.Contains(reason, "outcome.md") || !strings.Contains(reason, "missing") {
		t.Fatalf("ClassifyReview() reason = %q, want missing outcome.md", reason)
	}
}

func TestClassifyReviewReadsRecordedStates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		state       string
		mergeCommit string
		want        ReviewLiveness
	}{
		{name: "open remains live", state: "open", want: ReviewLive},
		{name: "closed is finished", state: "closed", want: ReviewFinished},
		{name: "merged with a hexadecimal receipt is finished", state: "merged", mergeCommit: "abcdef12", want: ReviewFinished},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reviewDir := filepath.Join(t.TempDir(), "pr-123")
			writeReviewOutcome(t, reviewDir, test.state, test.mergeCommit)

			got, reason, err := ClassifyReview(context.Background(), filepath.Join(t.TempDir(), "not-a-repository"), reviewDir)
			if err != nil {
				t.Fatalf("ClassifyReview() error = %v", err)
			}
			if got != test.want {
				t.Fatalf("ClassifyReview() = %q, want %q; reason: %s", got, test.want, reason)
			}
		})
	}
}

func TestClassifyReviewRejectsMalformedRecordedOutcome(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		outcome    string
		wantReason string
	}{
		{name: "malformed front matter", outcome: "---\npull_request_state: [\n---\n", wantReason: "malformed"},
		{name: "missing pull request state", outcome: "---\nrecorded_at: 2026-09-24T00:00:00Z\n---\n", wantReason: "pull_request_state"},
		{name: "missing recorded at", outcome: "---\npull_request_state: open\n---\n", wantReason: "recorded_at"},
		{name: "unknown pull request state", outcome: "---\npull_request_state: draft\nrecorded_at: 2026-09-24T00:00:00Z\n---\n", wantReason: "pull_request_state"},
		{name: "merged without merge commit", outcome: "---\npull_request_state: merged\nrecorded_at: 2026-09-24T00:00:00Z\n---\n", wantReason: "merge_commit"},
		{name: "merged with non hexadecimal merge commit", outcome: "---\npull_request_state: merged\nmerge_commit: not-a-commit\nrecorded_at: 2026-09-24T00:00:00Z\n---\n", wantReason: "merge_commit"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reviewDir := filepath.Join(t.TempDir(), "pr-123")
			if err := os.MkdirAll(reviewDir, 0o755); err != nil {
				t.Fatalf("create Review Artifact fixture: %v", err)
			}
			if err := os.WriteFile(filepath.Join(reviewDir, "outcome.md"), []byte(test.outcome), 0o644); err != nil {
				t.Fatalf("write outcome fixture: %v", err)
			}

			got, reason, err := ClassifyReview(context.Background(), t.TempDir(), reviewDir)
			if err != nil {
				t.Fatalf("ClassifyReview() error = %v", err)
			}
			if got != ReviewUndecidable {
				t.Fatalf("ClassifyReview() = %q, want %q; reason: %s", got, ReviewUndecidable, reason)
			}
			if !strings.Contains(reason, test.wantReason) {
				t.Fatalf("ClassifyReview() reason = %q, want %q", reason, test.wantReason)
			}
		})
	}
}

func reviewLivenessRepo(t *testing.T, defaultBranch string) (string, string) {
	t.Helper()

	repo := filepath.Join(t.TempDir(), "repo")
	gittest.InitRepo(t, repo, "-b", defaultBranch)
	if err := os.WriteFile(filepath.Join(repo, "tracked.txt"), []byte("base\n"), 0o644); err != nil {
		t.Fatalf("write base fixture: %v", err)
	}
	gittest.Run(t, repo, "add", "tracked.txt")
	gittest.Run(t, repo, "commit", "-m", "base")
	return repo, strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
}

func reviewLivenessBranchCommit(t *testing.T, repo string, branch string) string {
	t.Helper()

	gittest.Run(t, repo, "switch", "-c", branch)
	if err := os.WriteFile(filepath.Join(repo, "tracked.txt"), []byte(branch+"\n"), 0o644); err != nil {
		t.Fatalf("write branch fixture: %v", err)
	}
	gittest.Run(t, repo, "add", "tracked.txt")
	gittest.Run(t, repo, "commit", "-m", "review head")
	return strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
}

func persistReviewRound(t *testing.T, reviewDir string, round int, branch string, head string) {
	t.Helper()

	_, err := rounds.PersistRound(t.Context(), rounds.PersistRequest{
		ArtifactDir:    filepath.Dir(reviewDir),
		ReviewRoot:     reviewDir,
		Source:         "coderabbit",
		PRNumber:       "123",
		HeadRepository: "example/repository",
		HeadBranch:     branch,
		HeadSHA:        head,
		Round:          round,
	})
	if err != nil {
		t.Fatalf("persist Round %03d fixture: %v", round, err)
	}
}

func writeReviewOutcome(t *testing.T, reviewDir string, state string, mergeCommit string) {
	t.Helper()

	if err := os.MkdirAll(reviewDir, 0o755); err != nil {
		t.Fatalf("create Review Artifact fixture: %v", err)
	}
	content := fmt.Sprintf("---\npull_request_state: %s\nmerge_commit: %q\nrecorded_at: 2026-09-24T00:00:00Z\n---\n", state, mergeCommit)
	if err := os.WriteFile(filepath.Join(reviewDir, "outcome.md"), []byte(content), 0o644); err != nil {
		t.Fatalf("write outcome fixture: %v", err)
	}
}
