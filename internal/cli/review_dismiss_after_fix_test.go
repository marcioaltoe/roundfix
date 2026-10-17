// Suite: dismissing a reviewed finding after a committed fix.
// Invariant: a later dismissal preserves the reviewed head and proves the recorded anchor in Git.
// Boundary IN: public dispose command, real disposable Git history, and the disposition ledger.
// Boundary OUT: reviewer execution and delivery publication.
package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
)

const dismissalReviewedBody = "heading\nanchored line \r\nstable second\nfooter\n"

func newDismissalAnchorFixture(t *testing.T, anchor *reviewFindingAnchor) (reviewCommandFixture, reviewRecord) {
	t.Helper()
	f := newReviewDispositionFixture(t)
	f.headCommit = commitLineageFile(t, f, "review.txt", dismissalReviewedBody)
	record := dispositionReviewRecord(f.repository, f.headCommit, "review.txt:1: Failure: claimed defect")
	// The recorded anchor deliberately differs from the finding text's first line.
	record.FindingItems[0].Anchor = anchor
	record.FindingItems[0].Validation = &reviewFindingValidation{Status: "stands"}
	writeDispositionRecord(t, f, record)
	return f, record
}

func dismissalAnchorReason(f reviewCommandFixture, head string, anchor reviewFindingAnchor) string {
	return fmt.Sprintf("dismissal at HEAD %q after reviewed head %q requires the anchor of finding %q (%s:%d-%d) to be unchanged", head, f.headCommit, "F1", anchor.Path, anchor.StartLine, anchor.EndLine)
}

func assertDismissalRefusal(t *testing.T, f reviewCommandFixture, reason string) {
	t.Helper()
	path := filepath.Join(f.artifactDir, reviewDispositionLedgerFileName)
	before, beforeErr := os.ReadFile(path)
	if beforeErr != nil && !os.IsNotExist(beforeErr) {
		t.Fatal(beforeErr)
	}
	code, stdout, stderr := runReviewDispose(t, "F1", "--dismiss", "--evidence", "not a defect")
	want := "roundfix: review dispose refused: " + reason + "\n"
	if code != exitPreflight || stdout != "" || stderr != want {
		t.Fatalf("transcript: exit=%d stdout=%q stderr=%q; want exit=2 stdout=empty stderr=%q", code, stdout, stderr, want)
	}
	after, afterErr := os.ReadFile(path)
	if !sameDispositionLedgerState(before, beforeErr, after, afterErr) {
		t.Fatalf("refusal changed ledger: before=%q (%v), after=%q (%v)", before, beforeErr, after, afterErr)
	}
}

func assertDismissalLine(t *testing.T, f reviewCommandFixture, record reviewRecord, currentHead string) {
	t.Helper()
	code, stdout, stderr := runReviewDispose(t, "F1", "--dismiss", "--evidence", "not a defect")
	if code != exitOK || stderr != "" {
		t.Fatalf("transcript: exit=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	line := readSingleDispositionLine(t, f.artifactDir)
	entry := decodeDispositionLine(t, line)
	assertDispositionEntry(t, entry, f, "F1", record.FindingItems[0].Text, "dismissed")
	quote := func(value string) string {
		encoded, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		return string(encoded)
	}
	// Copy Surface Transcript 4's full line, varying only repository, heads, text and timestamp.
	want := fmt.Sprintf("{\"repository\":%s,\"headCommit\":%s,\"finding\":\"F1\",\"text\":%s,\"disposition\":\"dismissed\",\"evidence\":\"not a defect\",\"recordedAt\":%s", quote(f.repository), quote(f.headCommit), quote(record.FindingItems[0].Text), quote(entry.RecordedAt))
	if currentHead != "" {
		want += ",\"currentHead\":" + quote(currentHead)
	}
	want += "}\n"
	if stdout != want || line != want {
		t.Fatalf("stdout=%q ledger=%q, want %q", stdout, line, want)
	}
	// Repeating the accepted command cannot append another line.
	assertDisposeRefusedWithoutLedgerChange(t, f.artifactDir, "F1", "--dismiss", "--evidence", "not a defect")
}

func TestReviewDisposeReproducesTheDismissalAfterAFixTranscript(t *testing.T) {
	t.Parallel()
	f, record := newDismissalAnchorFixture(t, &reviewFindingAnchor{Path: "review.txt", StartLine: 2, EndLine: 3})
	head := commitLineageFile(t, f, "review.txt", strings.ReplaceAll(dismissalReviewedBody, "footer", "fixed footer"))
	// Dirty anchored lines must not affect a proof about the two commits.
	mustWrite(t, filepath.Join(f.repository, "review.txt"), "uncommitted replacement\n")
	assertDismissalLine(t, f, record, head)
}

func TestReviewDisposeReproducesTheChangedAnchorTranscript(t *testing.T) {
	t.Parallel()
	anchor := reviewFindingAnchor{Path: "review.txt", StartLine: 2, EndLine: 3}
	f, _ := newDismissalAnchorFixture(t, &anchor)
	head := commitLineageFile(t, f, "review.txt", strings.ReplaceAll(dismissalReviewedBody, "stable second", "changed second"))
	// Restoring the worktree cannot hide a change in the committed anchor.
	mustWrite(t, filepath.Join(f.repository, "review.txt"), dismissalReviewedBody)
	if _, err := appendReviewFindingDisposition(filepath.Join(f.artifactDir, reviewDispositionLedgerFileName), reviewFindingDisposition{Finding: "unrelated"}); err != nil {
		t.Fatal(err)
	}
	assertDismissalRefusal(t, f, dismissalAnchorReason(f, head, anchor))
}

func TestReviewDisposeRefusesADismissalWhenTheAnchorMoved(t *testing.T) {
	t.Parallel()
	anchor := reviewFindingAnchor{Path: "review.txt", StartLine: 2, EndLine: 3}
	f, _ := newDismissalAnchorFixture(t, &anchor)
	head := commitLineageFile(t, f, "review.txt", "inserted above\n"+dismissalReviewedBody)
	assertDismissalRefusal(t, f, dismissalAnchorReason(f, head, anchor))
}

func TestReviewDisposeRefusesADismissalWhenTheAnchoredFileIsGone(t *testing.T) {
	t.Parallel()
	for _, missing := range []string{"current head", "reviewed head"} {
		t.Run(missing, func(t *testing.T) {
			t.Parallel()
			anchor := reviewFindingAnchor{Path: "review.txt", StartLine: 2, EndLine: 3}
			if missing == "reviewed head" {
				anchor.Path = "later.txt"
			}
			f, _ := newDismissalAnchorFixture(t, &anchor)
			var head string
			if missing == "current head" {
				gittest.Run(t, f.repository, "rm", "review.txt")
				gittest.Run(t, f.repository, "commit", "-m", "test: delete anchored file")
				head = strings.TrimSpace(gittest.Run(t, f.repository, "rev-parse", "HEAD"))
			} else {
				head = commitLineageFile(t, f, "later.txt", dismissalReviewedBody)
			}
			assertDismissalRefusal(t, f, dismissalAnchorReason(f, head, anchor))
		})
	}
}

func TestReviewDisposeRefusesADismissalWithoutARecordedAnchor(t *testing.T) {
	t.Parallel()
	f, _ := newDismissalAnchorFixture(t, nil)
	head := commitLineageFile(t, f, "other-fix.txt", "fixed\n")
	reason := fmt.Sprintf("dismissal at HEAD %q after reviewed head %q requires finding %q to record an anchor", head, f.headCommit, "F1")
	assertDismissalRefusal(t, f, reason)
}

func TestReviewDisposeRefusesADismissalAtAHeadThatDoesNotDescend(t *testing.T) {
	t.Parallel()
	// Missing anchor must not mask the earlier ancestry refusal.
	f, _ := newDismissalAnchorFixture(t, nil)
	gittest.Run(t, f.repository, "checkout", "-b", "fix/divergent", f.baseCommit)
	head := commitLineageFile(t, f, "other-fix.txt", "different history\n")
	assertDismissalRefusal(t, f, fmt.Sprintf("dismissal requires current HEAD %q to equal or descend from reviewed head %q", head, f.headCommit))
}

func TestReviewDisposeComparesTheWholeFileForALineZeroAnchor(t *testing.T) {
	t.Parallel()
	for _, changed := range []bool{false, true} {
		t.Run(fmt.Sprintf("changed=%t", changed), func(t *testing.T) {
			t.Parallel()
			anchor := reviewFindingAnchor{Path: "review.txt"}
			f, record := newDismissalAnchorFixture(t, &anchor)
			var head string
			if changed {
				head = commitLineageFile(t, f, "review.txt", strings.ReplaceAll(dismissalReviewedBody, "footer", "fixed footer"))
			} else {
				head = commitLineageFile(t, f, "other-fix.txt", "fixed\n")
			}
			if changed {
				assertDismissalRefusal(t, f, dismissalAnchorReason(f, head, anchor))
			} else {
				assertDismissalLine(t, f, record, head)
			}
		})
	}
}

func TestReviewDisposeAtTheReviewedHeadWritesNoCurrentHead(t *testing.T) {
	t.Parallel()
	// Equality still accepts a legacy finding without an anchor and preserves the old JSON bytes.
	f, record := newDismissalAnchorFixture(t, nil)
	assertDismissalLine(t, f, record, "")
}

func TestReviewCorrectionAcceptsADismissalRecordedAfterTheFix(t *testing.T) {
	t.Parallel()
	f := newReviewDispositionFixture(t)
	f.headCommit = commitReviewCorrectionFile(t, f, correctionArchivePath, "needs correction\nstable evidence\n")
	record := dispositionReviewRecord(f.repository, f.headCommit, "- "+correctionArchivePath+":1: Failure: fix record\n- "+correctionArchivePath+":2: Failure: claimed evidence defect")
	record.ArchivedSpecs = []string{"0225-example"}
	record.Lineage = &reviewLineage{Round: 1}
	for i := range record.FindingItems {
		record.FindingItems[i].Anchor = &reviewFindingAnchor{Path: correctionArchivePath, StartLine: i + 1, EndLine: i + 1}
		record.FindingItems[i].Validation = &reviewFindingValidation{Status: "stands"}
	}
	writeDispositionRecord(t, f, record)
	if err := persistReviewRecord(f.artifactDir, record); err != nil {
		t.Fatal(err)
	}
	head := commitReviewCorrectionFile(t, f, correctionArchivePath, "corrected record\nstable evidence\n")
	for _, args := range [][]string{{"F1", "--fixed-by", head}, {"F2", "--dismiss", "--evidence", "archived QA already records the evidence"}} {
		code, stdout, stderr := runReviewDispose(t, args...)
		if code != exitOK || stderr != "" {
			t.Fatalf("dispose %v: exit=%d stdout=%q stderr=%q", args, code, stdout, stderr)
		}
	}
	ledger, err := readReviewFindingDispositions(filepath.Join(f.artifactDir, reviewDispositionLedgerFileName))
	if err != nil || len(ledger) != 2 {
		t.Fatalf("ledger=%+v err=%v", ledger, err)
	}
	if ledger[0].HeadCommit != f.headCommit || ledger[0].FixedBy != head || ledger[0].CurrentHead != "" || ledger[1].HeadCommit != f.headCommit || ledger[1].CurrentHead != head {
		t.Fatalf("ledger=%+v", ledger)
	}
	config := roundconfig.Builtin()
	config.Defaults.ArtifactDir = f.artifactDir
	workflow := &commandDeliveryWorkflow{git: preflight.ExecGitRunner{}, loaded: roundconfig.Loaded{Config: config, HomeDir: f.homeDir, GitRoot: f.repository}}
	proof, err := workflow.ProveReviewCorrection(t.Context(), f.repository, record.ArchivedSpecs, f.headCommit, head)
	if err != nil || !proof.Accepted || proof.Reason != "" {
		t.Fatalf("correction proof=%+v err=%v", proof, err)
	}
}

func TestReviewDisposeRequiresEveryAnchoredByteAtBothCommits(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		anchor reviewFindingAnchor
		body   string
	}{
		{"final newline", reviewFindingAnchor{"review.txt", 4, 4}, strings.TrimSuffix(dismissalReviewedBody, "\n")},
		{"empty final line in whole blob", reviewFindingAnchor{Path: "review.txt"}, dismissalReviewedBody + "\n"},
		{"carriage return", reviewFindingAnchor{"review.txt", 2, 2}, strings.ReplaceAll(dismissalReviewedBody, "\r", "")},
		{"trailing space", reviewFindingAnchor{"review.txt", 2, 2}, strings.ReplaceAll(dismissalReviewedBody, "line \r", "line\r")},
		{"missing range end at current head", reviewFindingAnchor{"review.txt", 2, 3}, "heading\nanchored line \r\n"},
		{"missing range end at reviewed head", reviewFindingAnchor{"review.txt", 4, 5}, dismissalReviewedBody + "new last line\n"},
		{"empty file at current head", reviewFindingAnchor{"review.txt", 1, 1}, ""},
		{"trailing newline is not another line", reviewFindingAnchor{"review.txt", 5, 5}, "changed heading\nanchored line \r\nstable second\nfooter\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f, _ := newDismissalAnchorFixture(t, &tc.anchor)
			head := commitLineageFile(t, f, "review.txt", tc.body)
			assertDismissalRefusal(t, f, dismissalAnchorReason(f, head, tc.anchor))
		})
	}
}

type dismissalAncestryFailureGit struct{ preflight.ExecGitRunner }

func (git dismissalAncestryFailureGit) RunGit(ctx context.Context, root string, args ...string) (string, error) {
	if len(args) > 0 && args[0] == "merge-base" {
		return "", errors.New("ancestry unavailable")
	}
	return git.ExecGitRunner.RunGit(ctx, root, args...)
}

func TestReviewDisposeReportsAncestryFailureWithoutAppending(t *testing.T) {
	t.Parallel()
	f, _ := newDismissalAnchorFixture(t, nil)
	commitLineageFile(t, f, "other-fix.txt", "fixed\n")
	withReviewSpecGitRunner(t, dismissalAncestryFailureGit{})
	assertDismissalRefusal(t, f, "check dismissal ancestry: ancestry unavailable")
}
