// Suite: bounded commit ranges through the public review command.
// Boundary IN: disposable Git commits; OUT: one prepared session and record.
package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/agent"
	"roundfix/internal/preflight"
	"roundfix/internal/runevent"
)

type segmentReviewRunner struct {
	reviewCommandRunner
	prepared []agent.ExecuteRequest
	requests []agent.ExecuteRequest
}

func (runner *segmentReviewRunner) PrepareSession(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) error {
	runner.prepared = append(runner.prepared, req)
	return runner.reviewCommandRunner.PrepareSession(ctx, req, sink)
}

func (runner *segmentReviewRunner) RunPrepared(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	runner.requests = append(runner.requests, req)
	return runner.reviewCommandRunner.RunPrepared(ctx, req, sink)
}

func commitReviewSegments(t *testing.T, fixture reviewCommandFixture, count int) []string {
	t.Helper()
	heads := []string{fixture.baseCommit}
	for i := 1; i <= count; i++ {
		heads = append(heads, commitScopeFile(t, fixture, fmt.Sprintf("segment%d.txt", i), strings.Repeat("x", reviewDiffBound*2/3)+"\n"))
	}
	return heads
}

func TestReviewSegmentsACandidateAboveTheBound(t *testing.T) {
	t.Parallel()
	runner := &segmentReviewRunner{reviewCommandRunner: reviewCommandRunner{results: []reviewCommandRunResult{
		{result: agent.ExecuteResult{Message: "Findings:\n- segment1.txt:1 Failure: the first segment breaks.", ACPSessionID: "one-session"}},
		{result: agent.ExecuteResult{Message: "Findings:\n- segment1.txt:1 Failure: this anchor belongs to another segment.", ACPSessionID: "one-session"}},
		{result: agent.ExecuteResult{Message: "Findings:\n- segment3.txt:1 Failure: the third segment breaks.", ACPSessionID: "one-session"}},
	}}}
	fixture := newReviewCommandFixture(t, "codex", runner)
	// Candidate context is included only in the first prompt.
	commitScopeFile(t, fixture, "docs/specs/example/_prd.md", "# Example\n\n## Decisions\n\nKeep this decision.\n")
	commitScopeFile(t, fixture, "docs/specs/example/_techspec.md", "# Technical design\n")
	config, err := os.ReadFile(filepath.Join(fixture.repository, ".roundfixrc.yml"))
	if err != nil {
		t.Fatal(err)
	}
	commitScopeFile(t, fixture, ".roundfixrc.yml", string(config)+"delivery:\n  derived_paths:\n    - paths: ['segment*.txt']\n      regenerate: make generate\n")
	heads := commitReviewSegments(t, fixture, 3)
	code, record, stderr := fixture.run(t)
	if code != exitRunFailed || record.Outcome != reviewOutcomeFindings {
		t.Fatalf("exit=%d record=%+v stderr=%q", code, record, stderr)
	}
	if runner.prepareCalls != 1 || len(runner.requests) != 3 || runner.endCalls != 0 {
		t.Fatalf("runner=%+v", runner)
	}
	var fields struct {
		Segments []struct {
			Base, Head string
			DiffBytes  int
		}
		FindingItems []struct {
			ID      string
			Segment int
		}
	}
	data, err := os.ReadFile(filepath.Join(reviewCheckoutDir(fixture.artifactDir, fixture.repository), reviewRecordFileName))
	if err != nil {
		t.Fatal(err)
	}
	// Decode persisted fields as well as checking the in-memory record below.
	if err := json.Unmarshal(data, &fields); err != nil {
		t.Fatal(err)
	}
	if len(fields.Segments) != 3 || len(fields.FindingItems) != 3 {
		t.Fatalf("record=%s", data)
	}
	total := 0
	for i, request := range runner.requests {
		segment := fields.Segments[i]
		want := fmt.Sprintf("Review Segment %d of 3: commits %s..%s of the candidate %s..%s; answer for this segment only.\n", i+1, heads[i], heads[i+1], heads[0], heads[3])
		if !strings.HasPrefix(request.Prompt, want) || request.Session != runner.prepared[0].Session || request.Access != agent.SessionAccessReadOnly {
			t.Fatalf("request %d=%+v", i+1, request)
		}
		if strings.Contains(request.Prompt, "--- BEGIN SPEC CONTEXT:") != (i == 0) {
			t.Fatalf("context on segment %d", i+1)
		}
		diff, err := (preflight.ExecGitRunner{}).RunGit(t.Context(), fixture.repository, "diff", "--no-ext-diff", "--no-textconv", "--no-color", heads[i], heads[i+1], "--")
		if err != nil {
			t.Fatal(err)
		}
		if segment.Base != heads[i] || segment.Head != heads[i+1] || segment.DiffBytes != len(diff) || segment.DiffBytes > reviewDiffBound {
			t.Fatalf("segment=%+v", segment)
		}
		total += len(diff)
		finding := record.FindingItems[i]
		if finding.ID != fmt.Sprintf("F%d", i+1) || fields.FindingItems[i].ID != finding.ID || fields.FindingItems[i].Segment != i+1 {
			t.Fatalf("finding=%+v", finding)
		}
		if reviewFindingDismissedByValidation(finding) != (i == 1) {
			t.Fatalf("wrong segment anchor validation: %+v", finding)
		}
	}
	if len(record.OmittedPaths) != 0 || record.DiffBytes != total || record.Lineage.Round != 1 || !record.Lineage.SessionOpen || !reflect.DeepEqual(record.Lineage.ACPSessionIDs, []string{"one-session"}) {
		t.Fatalf("record=%+v", record)
	}
	answer, err := os.ReadFile(record.AnswerPath)
	if err != nil {
		t.Fatal(err)
	}
	for i := 1; i <= 3; i++ {
		if !strings.Contains(string(answer), fmt.Sprintf("--- Review Segment %d of 3 ---\n", i)) {
			t.Fatalf("answer=%s", answer)
		}
	}
	// A repeat at this candidate head reuses all findings as one verdict.
	_, reused, _ := fixture.run(t)
	if !reused.Reused || len(runner.requests) != 3 {
		t.Fatalf("reuse=%+v prompts=%d", reused, len(runner.requests))
	}
}

func TestReviewBlocksACandidateNeedingMoreThanFourSegments(t *testing.T) {
	t.Parallel()
	runner := &segmentReviewRunner{}
	fixture := newReviewCommandFixture(t, "codex", runner)
	commitReviewSegments(t, fixture, 5)
	code, record, stderr := fixture.run(t)
	want := fmt.Sprintf("review diff too large: %d bytes after omitting 0 path(s) needs 5 review segments; a review allows at most 4", record.DiffBytes)
	assertBlockedReviewCommand(t, code, record, stderr, want)
	if record.Reason != want || runner.probeCalls != 0 || runner.prepareCalls != 0 || runner.preparedCalls != 0 {
		t.Fatalf("record=%+v runner=%+v", record, runner)
	}
}

func TestReviewSegmentsStopAtTheFirstBlockedAnswer(t *testing.T) {
	t.Parallel()
	for _, answer := range []reviewCommandRunResult{
		{result: agent.ExecuteResult{Message: "unclassifiable answer"}},
		{err: errors.New("review transport failed")},
		{result: agent.ExecuteResult{Message: "Findings:\n- missing anchor"}},
	} {
		t.Run(fmt.Sprint(answer), func(t *testing.T) {
			runner := &segmentReviewRunner{reviewCommandRunner: reviewCommandRunner{results: []reviewCommandRunResult{{result: agent.ExecuteResult{Message: "No findings."}}, answer}}}
			fixture := newReviewCommandFixture(t, "codex", runner)
			commitReviewSegments(t, fixture, 3)
			code, record, stderr := fixture.run(t)
			if code != exitPreflight || record.Outcome != reviewOutcomeBlocked || runner.prepareCalls != 1 || len(runner.requests) != 2 || runner.endCalls != 1 {
				t.Fatalf("exit=%d record=%+v stderr=%q runner=%+v", code, record, stderr, runner)
			}
			persisted, err := os.ReadFile(record.AnswerPath)
			if err != nil || strings.Contains(string(persisted), "--- Review Segment 3 of") {
				t.Fatalf("answer=%s err=%v", persisted, err)
			}
		})
	}
}

func TestReviewGreedilyPacksFourSegments(t *testing.T) {
	t.Parallel()
	runner := &segmentReviewRunner{}
	for range 4 {
		runner.results = append(runner.results, reviewCommandRunResult{result: agent.ExecuteResult{Message: "No findings."}})
	}
	runner.results[0].result = agent.ExecuteResult{Message: "intermediate reviewer message\n\nNo findings.", Messages: []string{"intermediate reviewer message", "No findings."}}
	fixture := newReviewCommandFixture(t, "codex", runner)
	for i := 1; i <= 8; i++ {
		commitScopeFile(t, fixture, fmt.Sprintf("part%d.txt", i), strings.Repeat("x", reviewDiffBound/3)+"\n")
	}
	code, record, stderr := fixture.run(t)
	if code != exitOK || len(record.Segments) != 4 || runner.prepareCalls != 1 || len(runner.requests) != 4 || runner.endCalls != 1 {
		t.Fatalf("exit=%d record=%+v stderr=%q runner=%+v", code, record, stderr, runner)
	}
	answer, err := os.ReadFile(record.AnswerPath)
	if err != nil || !strings.Contains(string(answer), "intermediate reviewer message\n\nNo findings.") {
		t.Fatalf("answer=%s err=%v", answer, err)
	}
}

func TestReviewBlocksAnOversizedCommitHiddenByItsRange(t *testing.T) {
	t.Parallel()
	runner := &segmentReviewRunner{}
	fixture := newReviewCommandFixture(t, "codex", runner)
	commitScopeFile(t, fixture, "rewrite.txt", strings.Repeat("x", reviewDiffBound*2/3)+"\n")
	commitScopeFile(t, fixture, "rewrite.txt", strings.Repeat("y", reviewDiffBound*2/3)+"\n")
	commitScopeFile(t, fixture, "next.txt", strings.Repeat("z", reviewDiffBound*2/3)+"\n")
	code, record, stderr := fixture.run(t)
	want := fmt.Sprintf("review diff too large: %d bytes after omitting 0 path(s) exceeds the review bound of %d bytes", record.DiffBytes, reviewDiffBound)
	assertBlockedReviewCommand(t, code, record, stderr, want)
	if runner.probeCalls != 0 || runner.prepareCalls != 0 || runner.preparedCalls != 0 {
		t.Fatal("oversized commit reached provider")
	}
}
