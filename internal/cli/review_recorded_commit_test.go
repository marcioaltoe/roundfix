// Suite: recorded review commits that a complete repository no longer holds.
// Boundary: disposable Git histories and review artifacts; the reviewer is fake.
package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
)

func TestAReviewWhoseReviewedHeadIsGoneOpensANewLineage(t *testing.T) {
	t.Parallel()
	f := newReviewCommandFixture(t, "codex", newLineageRunner("No findings.", "No findings."))
	code, prior, stderr := runLineageReview(t, f)
	if code != exitOK || stderr != "" || prior.Lineage == nil || prior.Lineage.Round != 1 {
		t.Fatalf("first review exit=%d record=%+v stderr=%q", code, prior, stderr)
	}
	gittest.Run(t, f.repository, "checkout", "-b", "feature/next", "main")
	head := commitLineageFile(t, f, "next.txt", "new candidate\n")
	gittest.Run(t, f.repository, "branch", "-D", "feature/review")
	gittest.Run(t, f.repository, "reflog", "expire", "--expire=now", "--all")
	gittest.Run(t, f.repository, "gc", "--prune=now")
	var out, diagnostic bytes.Buffer
	cmd := exec.CommandContext(t.Context(), "git", "-C", f.repository, "-c", "core.fsmonitor=false", "cat-file", "-e", prior.HeadCommit)
	cmd.Env = gittest.IsolatedEnv()
	cmd.Stdout, cmd.Stderr = &out, &diagnostic
	err := cmd.Run()
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != 1 || out.Len() != 0 || diagnostic.Len() != 0 {
		t.Fatalf("old head probe err=%v stdout=%q stderr=%q", err, out.String(), diagnostic.String())
	}
	candidate := newReviewRecord(f.repository, f.baseCommit, head, reviewPolicyForLineage(), reviewOutcomeBlocked)
	plan, err := decideReviewLineage(t.Context(), &prior, candidate, preflight.ExecGitRunner{})
	wantReason := "the reviewed head " + prior.HeadCommit + " is no longer in the repository"
	if err != nil || plan.Restart != wantReason || plan.Lineage.Round != 1 || plan.prior != nil || plan.Reuse || plan.Ceiling {
		t.Fatalf("plan=%+v err=%v", plan, err)
	}
	loaded, err := loadCommandConfig(commandEnvironmentForTest(t), io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	wantNotice := fmt.Sprintf("Review: 0372-example head %s opens a new Reviewer Lineage at round 1, a full review of the candidate: %s", head, wantReason)
	if notice := reviewLineageNotice(t.Context(), preflight.ExecGitRunner{}, loaded, "0372-example", f.repository); notice != wantNotice {
		t.Fatalf("notice=%q, want %q", notice, wantNotice)
	}
	code, second, stderr := runLineageReview(t, f)
	if code != exitOK || stderr != "" || second.Lineage == nil || second.Lineage.Round != 1 || second.Lineage.PreviousHead != "" || second.HeadCommit != head {
		t.Fatalf("second review exit=%d record=%+v stderr=%q", code, second, stderr)
	}
}

func TestARecordedMergeBaseThatIsGoneOpensANewLineage(t *testing.T) {
	t.Parallel()
	f := newReviewCommandFixture(t, "codex", &reviewCommandRunner{})
	prior := newReviewRecord(f.repository, strings.Repeat("f", 40), f.headCommit, reviewPolicyForLineage(), reviewOutcomeReviewed)
	prior.Lineage = &reviewLineage{Round: 2}
	head := commitLineageFile(t, f, "next.txt", "next\n")
	candidate := newReviewRecord(f.repository, f.baseCommit, head, reviewPolicyForLineage(), reviewOutcomeBlocked)
	plan, err := decideReviewLineage(t.Context(), &prior, candidate, preflight.ExecGitRunner{})
	want := "the recorded merge base " + prior.BaseCommit + " is no longer in the repository"
	if err != nil || plan.Restart != want || plan.Lineage.Round != 1 || plan.prior != nil || plan.Reuse || plan.Ceiling {
		t.Fatalf("plan=%+v err=%v, want restart %q", plan, err, want)
	}
}

func TestARecordedCommitThatIsNotProvenGoneStillBlocks(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"shallow clone", "partial clone", "object store error", "abbreviated commit", "Git failure on a present commit"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			f := newReviewCommandFixture(t, "codex", &reviewCommandRunner{})
			head := commitLineageFile(t, f, "third.txt", "third commit\n")
			root, recorded := f.repository, strings.Repeat("f", 40)
			var git preflight.GitRunner = preflight.ExecGitRunner{}
			if name == "shallow clone" || name == "partial clone" || name == "object store error" {
				root = filepath.Join(t.TempDir(), "clone")
				sourcePath := filepath.ToSlash(f.repository)
				if !strings.HasPrefix(sourcePath, "/") {
					sourcePath = "/" + sourcePath
				}
				source := (&url.URL{Scheme: "file", Path: sourcePath}).String()
				args := []string{"clone"}
				switch name {
				case "shallow clone":
					args = append(args, "--depth", "1")
					recorded = f.baseCommit
				case "partial clone":
					gittest.Run(t, f.repository, "config", "uploadpack.allowFilter", "true")
					args = append(args, "--filter=blob:none")
				case "object store error":
					args = append(args, "--no-local")
				}
				gittest.Run(t, "", append(args, source, root)...)
				gittest.Harden(t, root)
				switch name {
				case "partial clone":
					gittest.Run(t, root, "remote", "set-url", "origin", filepath.Join(t.TempDir(), "missing"))
				case "object store error":
					mustWrite(t, filepath.Join(root, ".git", "objects", "info", "alternates"), filepath.ToSlash(filepath.Join(t.TempDir(), "missing"))+"\n")
				}
			} else if name == "abbreviated commit" {
				recorded = recorded[:12]
			} else {
				recorded = f.headCommit
				git = dismissalAncestryFailureGit{}
			}
			prior := newReviewRecord(root, f.baseCommit, recorded, reviewPolicyForLineage(), reviewOutcomeReviewed)
			prior.Lineage = &reviewLineage{Round: 2}
			candidate := newReviewRecord(root, f.baseCommit, head, reviewPolicyForLineage(), reviewOutcomeBlocked)
			plan, err := decideReviewLineage(t.Context(), &prior, candidate, git)
			if err == nil || plan.Restart != "" {
				t.Fatalf("plan=%+v err=%v, want blocking error", plan, err)
			}
			if name == "shallow clone" || name == "partial clone" {
				want := fmt.Sprintf("check review lineage ancestry: a shallow or partial clone cannot prove that the recorded commit %s is gone; fetch the full history and try again", recorded)
				if err.Error() != want {
					t.Fatalf("err=%q, want %q", err, want)
				}
			} else {
				want := "check review lineage ancestry: git merge-base --is-ancestor " + recorded
				if name == "Git failure on a present commit" {
					want = "check review lineage ancestry: ancestry unavailable"
				}
				if !strings.HasPrefix(err.Error(), want) {
					t.Fatalf("err=%q, want prefix %q", err, want)
				}
			}
		})
	}
}

func TestReviewDisposeRefusesAReviewedHeadThatIsGone(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"dismiss", "fixed-by"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			f := newReviewDispositionFixture(t)
			fix := commitLineageFile(t, f, "fix.txt", "fixed\n")
			reviewed := strings.Repeat("f", 40)
			record := dispositionReviewRecord(f.repository, reviewed, "review.txt:1 Failure: claimed defect")
			writeDispositionRecord(t, f, record)
			ledger := filepath.Join(f.artifactDir, reviewDispositionLedgerFileName)
			before, beforeErr := os.ReadFile(ledger)
			if beforeErr != nil && !errors.Is(beforeErr, os.ErrNotExist) {
				t.Fatal(beforeErr)
			}
			args := []string{"F1", "--dismiss", "--evidence", "not a defect"}
			if name == "fixed-by" {
				args = []string{"F1", "--fixed-by", fix}
			}
			code, stdout, stderr := runReviewDispose(t, args...)
			want := fmt.Sprintf("roundfix: review dispose refused: reviewed head %q is no longer in the repository; run roundfix review to open a new Reviewer Lineage\n", reviewed)
			if code != exitPreflight || stdout != "" || stderr != want {
				t.Fatalf("exit=%d stdout=%q stderr=%q, want %q", code, stdout, stderr, want)
			}
			after, afterErr := os.ReadFile(ledger)
			if !sameDispositionLedgerState(before, beforeErr, after, afterErr) {
				t.Fatalf("refusal changed ledger: before=%q (%v) after=%q (%v)", before, beforeErr, after, afterErr)
			}
		})
	}
}

func TestAReviewCorrectionNamesACommitThatIsGone(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"parked candidate", "fixing commit"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			f := newReviewCorrectionFixture(t)
			gone := strings.Repeat("f", 40)
			want := "parked candidate " + gone + " is no longer in the repository"
			if name == "parked candidate" {
				f.candidate = gone
			} else {
				f.dispositions[0].FixedBy = gone
				f.persist(t)
				want = "finding " + f.record.FindingItems[0].ID + " requires exactly one valid disposition with evidence or a fixing commit between the candidate and item head"
			}
			result, err := f.workflow.ProveReviewCorrection(t.Context(), f.review.repository, []string{"0225-example"}, f.candidate, f.head)
			if err != nil || result.Accepted || result.Reason != want {
				t.Fatalf("proof=%+v err=%v, want refusal %q", result, err, want)
			}
		})
	}
}
