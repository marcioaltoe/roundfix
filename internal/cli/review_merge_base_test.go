// Suite: merge-base-bound pre-PR review.
// Invariant: review examines only what the current candidate introduced.
// Boundary IN: public review command, real Git repositories, records, and reviewer prompts.
// Boundary OUT: real ACP reviewer runtimes and delivery queue transitions.
package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/agent"
	"roundfix/internal/gittest"
)

func TestReviewDiffsTheCandidateFromItsMergeBase(t *testing.T) {
	t.Parallel()
	runner := &reviewCommandRunner{
		results: []reviewCommandRunResult{{
			result: agent.ExecuteResult{Message: "No findings", StopReason: "end_turn"},
		}},
	}
	fixture := newReviewCommandFixture(t, "codex", runner)
	baseTip := advanceReviewBase(t, fixture, func() {
		mustWrite(t, filepath.Join(fixture.repository, "base-only.txt"), "the base branch moved\n")
		gittest.Run(t, fixture.repository, "add", "base-only.txt")
	})

	code, record, stderr := runReviewCommandAgainstBase(t, fixture, "main")

	if code != exitOK || record.Outcome != reviewOutcomeReviewed {
		t.Fatalf("review exit=%d record=%+v stderr=%q", code, record, stderr)
	}
	if record.BaseCommit != fixture.baseCommit || record.BaseTipCommit != baseTip {
		t.Fatalf("review base commits = merge base %q tip %q, want %q and %q", record.BaseCommit, record.BaseTipCommit, fixture.baseCommit, baseTip)
	}
	if !strings.Contains(runner.request.Prompt, "+the reviewer receives this changed line") {
		t.Fatalf("review prompt omits candidate change:\n%s", runner.request.Prompt)
	}
	if strings.Contains(runner.request.Prompt, "base-only.txt") || strings.Contains(runner.request.Prompt, "the base branch moved") {
		t.Fatalf("review prompt includes a base-only change:\n%s", runner.request.Prompt)
	}
}

func TestReviewListsOnlyTheArchivedSpecsTheCandidateChanged(t *testing.T) {
	t.Parallel()
	const (
		baseSlug      = "0180-base-only-archive"
		candidateSlug = "0181-candidate-archive"
	)
	runner := &reviewCommandRunner{
		results: []reviewCommandRunResult{{
			result: agent.ExecuteResult{Message: "No findings", StopReason: "end_turn"},
		}},
	}
	fixture := newReviewCommandFixture(t, "codex", runner)
	advanceReviewBase(t, fixture, func() {
		writeArchivedReviewSpec(t, fixture.repository, baseSlug, true)
		gittest.Run(t, fixture.repository, "add", filepath.ToSlash(filepath.Join("docs", "history", "specs", baseSlug)))
	})
	writeArchivedReviewSpec(t, fixture.repository, candidateSlug, true)
	gittest.Run(t, fixture.repository, "add", filepath.ToSlash(filepath.Join("docs", "history", "specs", candidateSlug)))
	gittest.Run(t, fixture.repository, "commit", "-m", "archive candidate Spec")
	fixture.headCommit = strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))

	code, record, stderr := runReviewCommandAgainstBase(t, fixture, "main")

	if code != exitOK || record.Outcome != reviewOutcomeReviewed {
		t.Fatalf("review exit=%d record=%+v stderr=%q", code, record, stderr)
	}
	if !reflect.DeepEqual(record.ArchivedSpecs, []string{candidateSlug}) {
		t.Fatalf("review archived Specs = %v, want only %q", record.ArchivedSpecs, candidateSlug)
	}
}

func TestReviewReusesAFindingsVerdictAfterTheBaseBranchMoves(t *testing.T) {
	t.Parallel()
	runner := &reviewCommandRunner{
		results: []reviewCommandRunResult{{
			result: agent.ExecuteResult{Message: "Findings:\n- review.txt:2: candidate finding", StopReason: "end_turn"},
		}},
	}
	fixture := newReviewCommandFixture(t, "codex", runner)

	firstCode, first, firstStderr := runReviewCommandAgainstBase(t, fixture, "main")
	if firstCode != exitRunFailed || first.Outcome != reviewOutcomeFindings {
		t.Fatalf("first review exit=%d record=%+v stderr=%q", firstCode, first, firstStderr)
	}
	advanceReviewBase(t, fixture, func() {
		mustWrite(t, filepath.Join(fixture.repository, "later-base-change.txt"), "later base change\n")
		gittest.Run(t, fixture.repository, "add", "later-base-change.txt")
	})

	secondCode, second, secondStderr := runReviewCommandAgainstBase(t, fixture, "main")

	if secondCode != exitRunFailed || second.Outcome != reviewOutcomeFindings || !second.Reused {
		t.Fatalf("reused review exit=%d record=%+v stderr=%q", secondCode, second, secondStderr)
	}
	if runner.preparedCalls != 1 {
		t.Fatalf("reviewer prompt calls = %d, want 1 across both reviews", runner.preparedCalls)
	}
}

func TestReviewRefusesABaseThatSharesNoHistory(t *testing.T) {
	t.Parallel()
	runner := &reviewCommandRunner{}
	fixture := newReviewCommandFixture(t, "codex", runner)
	tree := strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD^{tree}"))
	unrelated := strings.TrimSpace(gittest.Run(t, fixture.repository, "commit-tree", tree, "-m", "unrelated root"))
	gittest.Run(t, fixture.repository, "branch", "unrelated", unrelated)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"review", "--base", "unrelated"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("review exit = %d, want %d; stdout=%q stderr=%q", code, exitPreflight, stdout.String(), stderr.String())
	}
	for _, want := range []string{fixture.headCommit, unrelated, "head shares no history with the base"} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("review stderr = %q, want %q", stderr.String(), want)
		}
	}
	if !strings.HasSuffix(stderr.String(), "pass --base <ref>\n") {
		t.Fatalf("review stderr = %q, want pass --base <ref> suffix", stderr.String())
	}
	if runner.probeCalls != 0 || runner.prepareCalls != 0 || runner.preparedCalls != 0 {
		t.Fatalf("reviewer activity preceded refusal: probes=%d prepares=%d prompts=%d", runner.probeCalls, runner.prepareCalls, runner.preparedCalls)
	}
	if _, err := os.Stat(filepath.Join(fixture.artifactDir, reviewRecordFileName)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("review record exists after refusal: %v", err)
	}
}

func advanceReviewBase(t *testing.T, fixture reviewCommandFixture, stage func()) string {
	t.Helper()
	gittest.Run(t, fixture.repository, "checkout", "main")
	stage()
	gittest.Run(t, fixture.repository, "commit", "-m", "advance review base")
	tip := strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))
	gittest.Run(t, fixture.repository, "checkout", "feature/review")
	return tip
}

func runReviewCommandAgainstBase(t *testing.T, fixture reviewCommandFixture, base string) (int, reviewRecord, string) {
	t.Helper()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runCLI(t, []string{"review", "--base", base}, &stdout, &stderr)

	var record reviewRecord
	if err := json.Unmarshal(stdout.Bytes(), &record); err != nil {
		t.Fatalf("decode review record %q: %v", stdout.String(), err)
	}
	return code, record, stderr.String()
}
