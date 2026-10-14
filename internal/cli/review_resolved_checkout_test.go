// Boundary: public review commands, real Git worktrees under a symlinked root,
// and persisted correction evidence. Reviewer execution uses a local fake.
package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
)

func TestReviewCheckoutDirIsOneForEverySpellingOfACheckout(t *testing.T) {
	t.Parallel()
	fixture, resolved, sibling := newSymlinkedReviewCheckoutFixture(t, true)
	digest := sha256.Sum256([]byte(resolved))
	want := filepath.Join(fixture.artifactDir, "pre-pr-review", hex.EncodeToString(digest[:])[:16])
	for _, checkout := range []string{fixture.repository, resolved} {
		if got := reviewCheckoutDir(fixture.artifactDir, checkout); got != want {
			t.Errorf("directory for %q = %q, want %q", checkout, got, want)
		}
	}
	if got := reviewCheckoutDir(fixture.artifactDir, sibling); got == want {
		t.Fatal("sibling checkout shares the item's review directory")
	}
}

func TestTheCorrectionProofFindsTheRecordUnderASymlinkedWorktreeRoot(t *testing.T) {
	t.Parallel()
	for _, mode := range []struct {
		name     string
		explicit bool
	}{
		{"explicit-artifact-directory", true},
		{"default-artifact-directory", false},
	} {
		t.Run(mode.name, func(t *testing.T) {
			t.Parallel()
			fixture, _, _ := newSymlinkedReviewCheckoutFixture(t, mode.explicit)
			record := reviewSymlinkedCheckout(t, fixture)
			head := commitReviewCorrectionFile(t, fixture, correctionArchivePath, "# Corrected archived Spec\n")
			disposeSymlinkedReviewFinding(t, fixture, head)
			result := proveSymlinkedReviewCorrection(t, fixture, record.HeadCommit, head, mode.explicit)
			if !result.Accepted || result.Reason != "" {
				t.Fatalf("proof=%+v, want accepted docs-only correction", result)
			}
		})
	}
}

func TestTheCorrectionProofRefusesACodeFixUnderASymlinkedWorktreeRoot(t *testing.T) {
	t.Parallel()
	fixture, _, _ := newSymlinkedReviewCheckoutFixture(t, true)
	record := reviewSymlinkedCheckout(t, fixture)
	head := commitReviewCorrectionFile(t, fixture, "review.txt", "corrected code\n")
	disposeSymlinkedReviewFinding(t, fixture, head)
	result := proveSymlinkedReviewCorrection(t, fixture, record.HeadCommit, head, true)
	if result.Accepted || !strings.Contains(result.Reason, "changed path review.txt is outside the named archived Specs") {
		t.Fatalf("proof=%+v, want refusal of an undeclared code fix", result)
	}
}

func TestTheCorrectionProofRefusesARecordOfAnotherCheckout(t *testing.T) {
	t.Parallel()
	fixture, _, sibling := newSymlinkedReviewCheckoutFixture(t, true)
	record := reviewSymlinkedCheckout(t, fixture)
	head := commitReviewCorrectionFile(t, fixture, correctionArchivePath, "# Corrected archived Spec\n")
	disposeSymlinkedReviewFinding(t, fixture, head)
	record.Repository = sibling
	writeDispositionRecordFile(t, filepath.Join(reviewCheckoutDir(fixture.artifactDir, fixture.repository), reviewRecordFileName), record)
	headPath, ok := reviewHeadRecordPath(fixture.artifactDir, fixture.repository, record.HeadCommit)
	if !ok {
		t.Fatal("candidate is not a full commit name")
	}
	writeDispositionRecordFile(t, headPath, record)
	result := proveSymlinkedReviewCorrection(t, fixture, record.HeadCommit, head, true)
	if result.Accepted || !strings.Contains(result.Reason, "review record must name the item worktree") {
		t.Fatalf("proof=%+v, want refusal of another checkout's record", result)
	}
}

func newSymlinkedReviewCheckoutFixture(t *testing.T, explicitArtifacts bool) (reviewCommandFixture, string, string) {
	t.Helper()
	home, primary, base, _ := newReviewCommandRepository(t)
	primaryFixture := reviewCommandFixture{repository: primary}
	commitReviewCorrectionFile(t, primaryFixture, "docs/specs/.keep", "Spec root\n")
	candidate := commitReviewCorrectionFile(t, primaryFixture, correctionArchivePath, "# Archived Spec\n")
	root := t.TempDir()
	realRoot := filepath.Join(root, "real")
	if err := os.MkdirAll(realRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	linkRoot := filepath.Join(root, "link")
	if err := os.Symlink(realRoot, linkRoot); err != nil {
		t.Fatal(err)
	}
	item := filepath.Join(linkRoot, "repo-slug", "item")
	gittest.Run(t, primary, "worktree", "add", "-b", "test/review-item", item, candidate)
	sibling := filepath.Join(linkRoot, "repo-slug", "sibling")
	gittest.Run(t, primary, "worktree", "add", "-b", "test/review-sibling", sibling, candidate)
	resolved, err := filepath.EvalSymlinks(item)
	if err != nil {
		t.Fatal(err)
	}
	if resolved == item {
		t.Fatal("fixture item does not cross a symbolic link")
	}
	if got := strings.TrimSpace(gittest.Run(t, item, "rev-parse", "--show-toplevel")); got != resolved {
		t.Fatalf("Git checkout root = %q, want resolved path %q", got, resolved)
	}
	configuredArtifacts := ""
	if explicitArtifacts {
		configuredArtifacts = filepath.Join(root, "artifacts")
	}
	writeReviewCommandConfig(t, item, "codex", configuredArtifacts)
	artifacts, err := roundconfig.ResolveArtifactDirectory(configuredArtifacts, item, home)
	if err != nil {
		t.Fatal(err)
	}
	setCommandEnvironmentForTest(t, home, item)
	withAgentRunner(t, &reviewCommandRunner{results: []reviewCommandRunResult{{
		result: agent.ExecuteResult{Message: "Findings:\n- review.txt:2: failure: the changed line loses the review evidence", StopReason: "end_turn"},
	}}})
	return reviewCommandFixture{homeDir: home, repository: item, artifactDir: artifacts, baseCommit: base, headCommit: candidate, provider: "codex"}, resolved, sibling
}

func reviewSymlinkedCheckout(t *testing.T, fixture reviewCommandFixture) reviewRecord {
	t.Helper()
	code, record, stderr := runReviewAtCheckout(t, reviewCheckoutFixture{baseCommit: fixture.baseCommit}, fixture.repository)
	if code != exitRunFailed || record.Outcome != reviewOutcomeFindings || len(record.FindingItems) != 1 || reviewFindingDismissedByValidation(record.FindingItems[0]) {
		t.Fatalf("review exit=%d record=%+v stderr=%q, want one standing finding", code, record, stderr)
	}
	if record.HeadCommit != fixture.headCommit {
		t.Fatalf("reviewed head=%q, want %q", record.HeadCommit, fixture.headCommit)
	}
	return record
}

func disposeSymlinkedReviewFinding(t *testing.T, fixture reviewCommandFixture, head string) {
	t.Helper()
	code, stdout, stderr := runReviewDisposeAtCheckout(t, reviewCheckoutFixture{}, fixture.repository, "F1", "--fixed-by", head)
	if code != exitOK || stdout == "" || stderr != "" {
		t.Fatalf("dispose exit=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
}

func proveSymlinkedReviewCorrection(t *testing.T, fixture reviewCommandFixture, candidate, head string, explicitArtifacts bool) delivery.ReviewCorrection {
	t.Helper()
	config := roundconfig.Builtin()
	if explicitArtifacts {
		config.Defaults.ArtifactDir = fixture.artifactDir
	}
	workflow := &commandDeliveryWorkflow{
		git:    preflight.ExecGitRunner{},
		loaded: roundconfig.Loaded{Config: config, HomeDir: fixture.homeDir, GitRoot: fixture.repository},
	}
	result, err := workflow.ProveReviewCorrection(t.Context(), fixture.repository, []string{"0225-example"}, candidate, head)
	if err != nil {
		t.Fatalf("prove correction with link spelling: %v", err)
	}
	return result
}
