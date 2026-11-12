// Suite: reviewer lineage after merging the default branch.
// Invariant: a clean default-branch merge keeps the delta; every restart names its cause.
// Boundary IN: temporary Git repositories, review/dispose commands and checkout records.
// Boundary OUT: fake reviewer prompts and records; no network or user artifacts.
package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
)

func TestReviewLineageContinuesAcrossAConflictFreeMergeOfTheDefaultBranch(t *testing.T) {
	t.Parallel()
	runner := newLineageRunner("Findings:\n- review.txt:2 Failure: candidate defect", "No findings.")
	fixture := newReviewCommandFixture(t, "codex", runner)
	code, first, stderr := runDefaultBranchReview(t)
	if code != exitRunFailed || stderr != "" || first.Lineage == nil || first.Lineage.Round != 1 {
		t.Fatalf("round one exit=%d record=%+v stderr=%q", code, first, stderr)
	}
	fix := commitLineageFile(t, fixture, "review.txt", "before review\nfixed after round one\n")
	disposeLineageFinding(t, "--fixed-by", fix)
	base := mergeAdvancedDefaultBranch(t, fixture, "upstream.txt", "main moved on\n")
	code, second, stderr := runDefaultBranchReview(t)
	if code != exitOK || stderr != "" || second.Outcome != reviewOutcomeReviewed || second.Lineage == nil || second.Lineage.Round != 2 || second.Lineage.Correction {
		t.Fatalf("round two exit=%d record=%+v stderr=%q", code, second, stderr)
	}
	if second.Lineage.PreviousHead != first.HeadCommit || !reflect.DeepEqual(second.Lineage.PreviousFindings, first.FindingItems) || len(second.Lineage.PreviousDispositions) != 1 || second.Lineage.PreviousDispositions[0].FixedBy != fix {
		t.Fatalf("round two lineage=%+v", second.Lineage)
	}
	if len(runner.prompts) != 2 {
		t.Fatalf("reviewer calls=%d, want 2", len(runner.prompts))
	}
	assertMainMergeReviewDelta(t, fixture, second, first.HeadCommit, base, runner.prompts[1], "ROUND 2", "fixed after round one")
	if !strings.Contains(runner.prompts[1], "F1 | stands | fixed by "+fix) {
		t.Fatalf("round two prompt lost the disposition: %s", runner.prompts[1])
	}
	t.Run("oversized carried delta", func(t *testing.T) {
		t.Parallel()
		runner := newLineageRunner("Findings:\n- review.txt:2 Failure: candidate defect")
		fixture := newReviewCommandFixture(t, "codex", runner)
		code, first, stderr := runDefaultBranchReview(t)
		if code != exitRunFailed || stderr != "" {
			t.Fatalf("round one exit=%d record=%+v stderr=%q", code, first, stderr)
		}
		fix := commitLineageFile(t, fixture, "review.txt", "before review\n"+strings.Repeat("x", reviewDiffBound)+"\n")
		disposeLineageFinding(t, "--fixed-by", fix)
		mergeAdvancedDefaultBranch(t, fixture, "upstream.txt", "main moved on\n")
		probes, prepares := runner.probeCalls, runner.prepareCalls
		code, blocked, stderr := runDefaultBranchReview(t)
		if code != exitPreflight || blocked.Outcome != reviewOutcomeBlocked || blocked.Lineage == nil || blocked.Lineage.Round != 2 || blocked.Lineage.DeltaBase == "" || blocked.DiffBytes <= reviewDiffBound || len(blocked.Segments) != 0 {
			t.Fatalf("bound exit=%d record=%+v stderr=%q", code, blocked, stderr)
		}
		want := fmt.Sprintf("review diff too large: %d bytes after omitting 0 path(s) exceeds the review bound of %d bytes", blocked.DiffBytes, reviewDiffBound)
		want += reviewSizeRefusalTestSuffix(t, fixture, blocked.Lineage.DeltaBase, blocked.HeadCommit, blocked.OmittedPaths)
		if blocked.Reason != want || stderr != "roundfix: review blocked: "+want+"\n" || strings.Contains(blocked.Reason, "upstream.txt") {
			t.Fatalf("refusal reason=%q stderr=%q, want %q", blocked.Reason, stderr, want)
		}
		if runner.probeCalls != probes || runner.prepareCalls != prepares || len(runner.prompts) != 1 {
			t.Fatal("oversized carried delta reached the reviewer")
		}
	})
}

func TestCorrectionReviewAfterAMergeOfTheDefaultBranchReadsOnlyTheFix(t *testing.T) {
	t.Parallel()
	runner := newLineageRunner("No findings.", "Findings:\n- delta.txt:1 Failure: round two defect", "No findings.")
	fixture := newReviewCommandFixture(t, "codex", runner)
	code, first, stderr := runDefaultBranchReview(t)
	if code != exitOK || stderr != "" || first.Outcome != reviewOutcomeReviewed {
		t.Fatalf("round one exit=%d record=%+v stderr=%q", code, first, stderr)
	}
	commitLineageFile(t, fixture, "delta.txt", "round two defect\n")
	code, second, stderr := runDefaultBranchReview(t)
	if code != exitRunFailed || stderr != "" || second.Lineage == nil || second.Lineage.Round != 2 || len(second.FindingItems) != 1 || second.FindingItems[0].Validation.Status != "stands" {
		t.Fatalf("round two exit=%d record=%+v stderr=%q", code, second, stderr)
	}
	fix := commitLineageFile(t, fixture, "delta.txt", "fixed after round two\n")
	disposeLineageFinding(t, "--fixed-by", fix)
	base := mergeAdvancedDefaultBranch(t, fixture, "upstream.txt", "main moved on\n")
	code, correction, stderr := runDefaultBranchReview(t)
	if code != exitOK || stderr != "" || correction.Outcome != reviewOutcomeReviewed || correction.Lineage == nil || correction.Lineage.Round != 2 || !correction.Lineage.Correction {
		t.Fatalf("correction exit=%d record=%+v stderr=%q", code, correction, stderr)
	}
	if correction.Lineage.PreviousHead != second.HeadCommit || !reflect.DeepEqual(correction.Lineage.PreviousFindings, second.FindingItems) || len(correction.Lineage.PreviousDispositions) != 1 || correction.Lineage.PreviousDispositions[0].FixedBy != fix {
		t.Fatalf("correction lineage=%+v", correction.Lineage)
	}
	if len(runner.prompts) != 3 {
		t.Fatalf("reviewer calls=%d, want 3", len(runner.prompts))
	}
	assertMainMergeReviewDelta(t, fixture, correction, second.HeadCommit, base, runner.prompts[2], "CORRECTION", "fixed after round two")
	if !strings.Contains(runner.prompts[2], "- F1 | fixed by "+fix) || strings.Contains(runner.prompts[2], "BEGIN ROUND 2 DELTA DIFF") || strings.Contains(runner.prompts[2], "the reviewer receives this changed line") {
		t.Fatalf("correction prompt lost its scope or disposition: %s", runner.prompts[2])
	}
	code, reused, stderr := runDefaultBranchReview(t)
	if code != exitOK || stderr != "" || !reused.Reused || len(runner.prompts) != 3 || reused.Lineage.DeltaBase != correction.Lineage.DeltaBase {
		t.Fatalf("correction reuse exit=%d record=%+v calls=%d stderr=%q", code, reused, len(runner.prompts), stderr)
	}
}

func assertMainMergeReviewDelta(t *testing.T, fixture reviewCommandFixture, record reviewRecord, previousHead, base, prompt, marker, fix string) {
	t.Helper()
	if record.BaseCommit != base || record.Lineage.DeltaBase == "" || record.Lineage.DeltaBase == previousHead || len(record.Segments) != 0 {
		t.Fatalf("delta base=%q base=%s segments=%+v, want tree from %s onto %s", record.Lineage.DeltaBase, record.BaseCommit, record.Segments, previousHead, base)
	}
	expectedTree := strings.TrimSpace(gittest.Run(t, fixture.repository, "merge-tree", "--write-tree", previousHead, base))
	if record.Lineage.DeltaBase != expectedTree {
		t.Fatalf("deltaBase=%s, want %s", record.Lineage.DeltaBase, expectedTree)
	}
	if kind := strings.TrimSpace(gittest.Run(t, fixture.repository, "cat-file", "-t", record.Lineage.DeltaBase)); kind != "tree" {
		t.Fatalf("deltaBase kind=%s, want tree", kind)
	}
	line := fmt.Sprintf("The candidate merged the default branch since then; the diff starts from tree %s, the reviewed head with the default branch's changes applied, so it holds only the candidate's own change.\n", expectedTree)
	for _, want := range []string{fix, "--- BEGIN " + marker + " DELTA DIFF ---", "--- END " + marker + " DELTA DIFF ---", line} {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing %q: %s", want, prompt)
		}
	}
	if marker == "ROUND 2" && !strings.Contains(prompt, "The diff below is the change since round 1 reviewed "+previousHead+".\n"+line+"Round 1 findings:\n") {
		t.Errorf("round two measurement line is misplaced: %s", prompt)
	}
	if marker == "CORRECTION" && !strings.Contains(prompt, "it is reviewed once and no further review of this candidate follows.\n"+line+"- F1") {
		t.Errorf("correction measurement line is misplaced: %s", prompt)
	}
	for _, excluded := range []string{"main moved on", "BEGIN CANDIDATE DIFF"} {
		if strings.Contains(prompt, excluded) {
			t.Errorf("prompt includes %q: %s", excluded, prompt)
		}
	}
	persisted, err := readReviewRecord(filepath.Join(reviewCheckoutDir(fixture.artifactDir, fixture.repository), reviewRecordFileName))
	if err != nil || !reflect.DeepEqual(persisted, record) {
		t.Fatalf("persisted=%+v err=%v", persisted, err)
	}
	encoded, err := json.Marshal(record)
	if err != nil || !strings.Contains(string(encoded), `"deltaBase":"`+expectedTree+`"`) {
		t.Fatalf("deltaBase JSON=%s err=%v", encoded, err)
	}
}

func TestReviewLineageRestartNamesItsReason(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"clean merge", "provider", "rewritten head", "same head", "base not descending", "merge of another branch", "conflict"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			fixture := newReviewCommandFixture(t, "codex", &reviewCommandRunner{})
			prior := newReviewRecord(fixture.repository, fixture.baseCommit, fixture.headCommit, reviewPolicyForLineage(), reviewOutcomeFindings)
			prior.Lineage = &reviewLineage{Round: 1}
			candidate := prior
			var want string
			if name == "conflict" {
				gittest.Run(t, fixture.repository, "checkout", "main")
				candidate.BaseCommit = commitLineageFile(t, fixture, "review.txt", "before review\nmain rewrote the second line\n")
				gittest.Run(t, fixture.repository, "checkout", "feature/review")
				gittest.Run(t, fixture.repository, "merge", "--no-edit", "-X", "ours", "main")
			} else if name == "merge of another branch" {
				var oldestMerge, oldestParent string
				for _, branch := range []string{"test/other", "test/another"} {
					gittest.Run(t, fixture.repository, "checkout", "-b", branch, fixture.baseCommit)
					parent := commitLineageFile(t, fixture, strings.TrimPrefix(branch, "test/")+".txt", "unreviewed branch\n")
					gittest.Run(t, fixture.repository, "checkout", "feature/review")
					gittest.Run(t, fixture.repository, "merge", "--no-edit", branch)
					if oldestMerge == "" {
						oldestMerge = strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))
						oldestParent = parent
					}
				}
				candidate.BaseCommit = mergeAdvancedDefaultBranch(t, fixture, "upstream.txt", "main moved on\n")
				want = fmt.Sprintf("the merge base moved from %s to %s through merge %s of %s, which the default branch does not contain", prior.BaseCommit, candidate.BaseCommit, oldestMerge, oldestParent)
			} else if name == "base not descending" {
				prior.BaseCommit = mergeAdvancedDefaultBranch(t, fixture, "old-upstream.txt", "old main\n")
				prior.HeadCommit = strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))
				gittest.Run(t, fixture.repository, "branch", "-f", "main", fixture.baseCommit)
				candidate.BaseCommit = mergeAdvancedDefaultBranch(t, fixture, "upstream.txt", "main moved on\n")
				want = fmt.Sprintf("the merge base moved from %s to %s, which does not descend from it", prior.BaseCommit, candidate.BaseCommit)
			} else {
				gittest.Run(t, fixture.repository, "checkout", "main")
				candidate.BaseCommit = commitLineageFile(t, fixture, "upstream.txt", "main moved on\n")
				gittest.Run(t, fixture.repository, "checkout", "feature/review")
				if name == "rewritten head" || name == "provider" {
					gittest.Run(t, fixture.repository, "rebase", "main")
					want = "the head does not descend from the reviewed head " + prior.HeadCommit
					if name == "provider" {
						candidate.Provider = "claude"
						want = "the review provider changed from codex to claude"
					}
				} else if name == "same head" {
					want = fmt.Sprintf("the merge base moved from %s to %s at the reviewed head", prior.BaseCommit, candidate.BaseCommit)
				} else {
					gittest.Run(t, fixture.repository, "merge", "--no-edit", "main")
				}
			}
			candidate.HeadCommit = strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))
			plan, err := decideReviewLineage(t.Context(), &prior, candidate, preflight.ExecGitRunner{})
			if err != nil {
				t.Fatal(err)
			}
			round := 1
			if name == "clean merge" || name == "conflict" {
				round = 2
			}
			if plan.Restart != want || plan.Lineage.Round != round || plan.Reuse || plan.Ceiling {
				t.Fatalf("round=%d restart=%q reuse=%t ceiling=%t, want round=%d restart=%q", plan.Lineage.Round, plan.Restart, plan.Reuse, plan.Ceiling, round, want)
			}
		})
	}
	t.Run("git failures stay errors", func(t *testing.T) {
		t.Parallel()
		fixture := newReviewCommandFixture(t, "codex", &reviewCommandRunner{})
		prior := newReviewRecord(fixture.repository, fixture.baseCommit, fixture.headCommit, reviewPolicyForLineage(), reviewOutcomeFindings)
		candidate := prior
		candidate.BaseCommit = mergeAdvancedDefaultBranch(t, fixture, "upstream.txt", "main moved on\n")
		candidate.HeadCommit = strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))
		for _, tc := range []struct {
			name string
			args []string
		}{
			{"head ancestry", []string{"merge-base", "--is-ancestor", prior.HeadCommit, candidate.HeadCommit}},
			{"base ancestry", []string{"merge-base", "--is-ancestor", prior.BaseCommit, candidate.BaseCommit}},
			{"merge list", []string{"rev-list", "--first-parent", "--merges", "--parents", "--reverse", prior.HeadCommit + ".." + candidate.HeadCommit}},
			{"parent ancestry", []string{"merge-base", "--is-ancestor", candidate.BaseCommit, candidate.BaseCommit}},
			{"carry", []string{"merge-tree", "--write-tree", prior.HeadCommit, candidate.BaseCommit}},
		} {
			t.Run(tc.name, func(t *testing.T) {
				sentinel := errors.New("git object could not be read")
				plan, err := decideReviewLineage(t.Context(), &prior, candidate, mainMergeLineageFailureGit{args: tc.args, err: sentinel})
				if !errors.Is(err, sentinel) || plan.Restart != "" {
					t.Fatalf("restart=%q err=%v, want propagated Git failure", plan.Restart, err)
				}
			})
		}
		prior.HeadCommit = strings.Repeat("f", 40)
		plan, err := decideReviewLineage(t.Context(), &prior, candidate, preflight.ExecGitRunner{})
		// Spec 0372: a recorded head proven gone opens a new lineage.
		want := "the reviewed head " + prior.HeadCommit + " is no longer in the repository"
		if err != nil || plan.Lineage.Round != 1 || plan.Restart != want {
			t.Fatalf("plan=%+v err=%v, want restart %q at round 1", plan, err, want)
		}
	})
}

type mainMergeLineageFailureGit struct {
	preflight.ExecGitRunner
	args []string
	err  error
}

func (git mainMergeLineageFailureGit) RunGit(ctx context.Context, root string, args ...string) (string, error) {
	if reflect.DeepEqual(args, git.args) {
		return "", git.err
	}
	return git.ExecGitRunner.RunGit(ctx, root, args...)
}

func mergeAdvancedDefaultBranch(t *testing.T, fixture reviewCommandFixture, name, content string) string {
	t.Helper()
	gittest.Run(t, fixture.repository, "checkout", "main")
	tip := commitLineageFile(t, fixture, name, content)
	gittest.Run(t, fixture.repository, "checkout", "feature/review")
	gittest.Run(t, fixture.repository, "merge", "--no-edit", "main")
	return tip
}
