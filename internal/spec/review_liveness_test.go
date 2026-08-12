package spec

// Suite: orphan Review Artifact liveness
// Invariant: only local Git proof from the newest Round may retire an orphan Review Artifact.
// Boundary IN: persisted Round metadata and an isolated repository's Git objects and refs.
// Boundary OUT: hosting-provider state, credentials, network access, and history relocation.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/rounds"
)

func TestClassifyReviewLocalGit(t *testing.T) {
	t.Parallel()

	t.Run("finished when newest head is an ancestor of the default branch", func(t *testing.T) {
		t.Parallel()

		repo, _ := reviewLivenessRepo(t, "main")
		head := reviewLivenessBranchCommit(t, repo, "feature/merged")
		gittest.Run(t, repo, "switch", "main")
		gittest.Run(t, repo, "merge", "--no-ff", "-m", "merge review head", "feature/merged")
		reviewDir := filepath.Join(t.TempDir(), "pr-123")
		persistReviewRound(t, reviewDir, 1, "feature/merged", head)

		assertReviewLiveness(t, repo, reviewDir, ReviewFinished)
	})

	t.Run("finished against a nonstandard default named by origin HEAD", func(t *testing.T) {
		t.Parallel()

		repo, _ := reviewLivenessRepo(t, "trunk")
		gittest.Run(t, repo, "update-ref", "refs/remotes/origin/trunk", "HEAD")
		gittest.Run(t, repo, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/trunk")
		head := reviewLivenessBranchCommit(t, repo, "feature/merged-to-trunk")
		gittest.Run(t, repo, "switch", "trunk")
		gittest.Run(t, repo, "merge", "--no-ff", "-m", "merge review head", "feature/merged-to-trunk")
		reviewDir := filepath.Join(t.TempDir(), "pr-123")
		persistReviewRound(t, reviewDir, 1, "feature/merged-to-trunk", head)

		assertReviewLiveness(t, repo, reviewDir, ReviewFinished)
	})

	t.Run("finished when newest head is unreachable and its branch is gone", func(t *testing.T) {
		t.Parallel()

		repo, _ := reviewLivenessRepo(t, "main")
		head := reviewLivenessBranchCommit(t, repo, "feature/abandoned")
		gittest.Run(t, repo, "switch", "main")
		gittest.Run(t, repo, "branch", "-D", "feature/abandoned")
		reviewDir := filepath.Join(t.TempDir(), "pr-123")
		persistReviewRound(t, reviewDir, 1, "feature/abandoned", head)

		assertReviewLiveness(t, repo, reviewDir, ReviewFinished)
	})

	t.Run("live when newest head remains on a branch", func(t *testing.T) {
		t.Parallel()

		repo, _ := reviewLivenessRepo(t, "main")
		head := reviewLivenessBranchCommit(t, repo, "feature/live")
		gittest.Run(t, repo, "switch", "main")
		reviewDir := filepath.Join(t.TempDir(), "pr-123")
		persistReviewRound(t, reviewDir, 1, "feature/live", head)

		assertReviewLiveness(t, repo, reviewDir, ReviewLive)
	})

	t.Run("live when a tag keeps a deleted branch head reachable", func(t *testing.T) {
		t.Parallel()

		repo, _ := reviewLivenessRepo(t, "main")
		head := reviewLivenessBranchCommit(t, repo, "feature/tagged")
		gittest.Run(t, repo, "tag", "retained-review-head", head)
		gittest.Run(t, repo, "switch", "main")
		gittest.Run(t, repo, "branch", "-D", "feature/tagged")
		reviewDir := filepath.Join(t.TempDir(), "pr-123")
		persistReviewRound(t, reviewDir, 1, "feature/tagged", head)

		assertReviewLiveness(t, repo, reviewDir, ReviewLive)
	})

	t.Run("undecidable when an unreachable head branch still exists", func(t *testing.T) {
		t.Parallel()

		repo, _ := reviewLivenessRepo(t, "main")
		head := reviewLivenessBranchCommit(t, repo, "feature/rewritten")
		gittest.Run(t, repo, "switch", "main")
		gittest.Run(t, repo, "branch", "-f", "feature/rewritten", "main")
		reviewDir := filepath.Join(t.TempDir(), "pr-123")
		persistReviewRound(t, reviewDir, 1, "feature/rewritten", head)

		assertReviewLiveness(t, repo, reviewDir, ReviewUndecidable)
	})

	t.Run("newest Round decides when heads differ", func(t *testing.T) {
		t.Parallel()

		repo, mergedHead := reviewLivenessRepo(t, "main")
		liveHead := reviewLivenessBranchCommit(t, repo, "feature/newest")
		gittest.Run(t, repo, "switch", "main")
		reviewDir := filepath.Join(t.TempDir(), "pr-123")
		persistReviewRound(t, reviewDir, 1, "feature/older", mergedHead)
		persistReviewRound(t, reviewDir, 2, "feature/newest", liveHead)

		assertReviewLiveness(t, repo, reviewDir, ReviewLive)
	})

	t.Run("undecidable when newest Round records no head", func(t *testing.T) {
		t.Parallel()

		repo, head := reviewLivenessRepo(t, "main")
		reviewDir := filepath.Join(t.TempDir(), "pr-123")
		persistReviewRound(t, reviewDir, 1, "feature/missing-head", head)
		metadataPath := filepath.Join(reviewDir, "round-001", "round.md")
		metadata, err := os.ReadFile(metadataPath)
		if err != nil {
			t.Fatalf("read Round metadata fixture: %v", err)
		}
		metadata = []byte(strings.Replace(string(metadata), "head_sha: "+head, "head_sha: \"\"", 1))
		if err := os.WriteFile(metadataPath, metadata, 0o644); err != nil {
			t.Fatalf("remove recorded head from Round metadata fixture: %v", err)
		}

		got, reason, err := ClassifyReview(repo, reviewDir)
		if err != nil {
			t.Fatalf("ClassifyReview() error = %v", err)
		}
		if got == ReviewFinished {
			t.Fatalf("ClassifyReview() = %q, want an answer that remains live", got)
		}
		if got != ReviewUndecidable {
			t.Fatalf("ClassifyReview() = %q, want %q", got, ReviewUndecidable)
		}
		if strings.TrimSpace(reason) == "" {
			t.Fatal("ClassifyReview() returned an empty reason")
		}
	})

	t.Run("undecidable when newest Round metadata is malformed", func(t *testing.T) {
		t.Parallel()

		repo, _ := reviewLivenessRepo(t, "main")
		reviewDir := filepath.Join(t.TempDir(), "pr-123")
		roundDir := filepath.Join(reviewDir, "round-001")
		if err := os.MkdirAll(roundDir, 0o755); err != nil {
			t.Fatalf("create malformed Round fixture: %v", err)
		}
		if err := os.WriteFile(filepath.Join(roundDir, "round.md"), []byte("---\nhead_sha: [\n---\n"), 0o644); err != nil {
			t.Fatalf("write malformed Round fixture: %v", err)
		}

		assertReviewLiveness(t, repo, reviewDir, ReviewUndecidable)
	})

	t.Run("undecidable when local Git cannot identify a default branch", func(t *testing.T) {
		t.Parallel()

		repo, head := reviewLivenessRepo(t, "trunk")
		reviewDir := filepath.Join(t.TempDir(), "pr-123")
		persistReviewRound(t, reviewDir, 1, "feature/unknown", head)

		assertReviewLiveness(t, repo, reviewDir, ReviewUndecidable)
	})
}

func assertReviewLiveness(t *testing.T, repo string, reviewDir string, want ReviewLiveness) {
	t.Helper()

	got, reason, err := ClassifyReview(repo, reviewDir)
	if err != nil {
		t.Fatalf("ClassifyReview() error = %v", err)
	}
	if got != want {
		t.Fatalf("ClassifyReview() = %q, want %q; reason: %s", got, want, reason)
	}
	if strings.TrimSpace(reason) == "" {
		t.Fatal("ClassifyReview() returned an empty reason")
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
