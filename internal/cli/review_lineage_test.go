// Suite: bounded pre-PR reviewer lineage.
// Invariant: round 2 sees the delta; the ceiling consumes dispositions without a reviewer.
// Boundary IN: temporary Git commits, checkout records, and the public review CLI.
// Boundary OUT: reviewer execution is captured by a fake; no network or user artifacts.
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
	"roundfix/internal/runevent"
)

func TestReviewLineageDecidesEachRound(t *testing.T) {
	fixture := newReviewCommandFixture(t, "codex", &reviewCommandRunner{})
	prior := newReviewRecord(fixture.repository, fixture.baseCommit, fixture.headCommit, reviewPolicyForLineage(), reviewOutcomeReviewed)
	candidate := prior
	candidate.HeadCommit = commitLineageFile(t, fixture, "delta.txt", "round two\n")
	gittest.Run(t, fixture.repository, "checkout", "main")
	rebasedBase := commitLineageFile(t, fixture, "upstream.txt", "upstream advance\n")
	gittest.Run(t, fixture.repository, "checkout", "feature/review")
	gittest.Run(t, fixture.repository, "rebase", "main")
	rebasedHead := strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))
	tests := []struct {
		name           string
		change         func(*reviewRecord, *reviewRecord)
		absent         bool
		round          int
		reuse, ceiling bool
	}{
		{name: "absent", absent: true, round: 1},
		{name: "another provider", change: func(p, c *reviewRecord) { p.Provider = "claude" }, round: 1},
		{name: "moved merge base after rebase", change: func(p, c *reviewRecord) { c.BaseCommit = rebasedBase; c.HeadCommit = rebasedHead }, round: 1},
		{name: "head does not descend", change: func(p, c *reviewRecord) { c.HeadCommit = rebasedHead }, round: 1},
		{name: "same head reuse", change: func(p, c *reviewRecord) { c.HeadCommit = p.HeadCommit }, round: 1, reuse: true},
		{name: "same head blocked repeats", change: func(p, c *reviewRecord) {
			c.HeadCommit = p.HeadCommit
			p.Outcome = reviewOutcomeBlocked
			p.Lineage = &reviewLineage{Round: 2, PreviousHead: fixture.baseCommit}
		}, round: 2},
		{name: "ancestor blocked legacy repeats", change: func(p, c *reviewRecord) { p.Outcome = reviewOutcomeBlocked }, round: 1},
		{name: "ancestor blocked round two repeats", change: func(p, c *reviewRecord) {
			p.Outcome = reviewOutcomeBlocked
			p.Lineage = &reviewLineage{Round: 2, PreviousHead: fixture.baseCommit}
		}, round: 2},
		{name: "legacy reviewed", round: 2},
		{name: "round one reviewed", change: func(p, c *reviewRecord) { p.Lineage = &reviewLineage{Round: 1} }, round: 2},
		{name: "another checkout", change: func(p, c *reviewRecord) { p.Repository = t.TempDir() }, round: 1},
		{name: "same head closed reuse", change: func(p, c *reviewRecord) {
			c.HeadCommit = p.HeadCommit
			p.Outcome = reviewOutcomeCeilingClosed
			p.Lineage = &reviewLineage{Round: 2, ReviewedHead: fixture.baseCommit}
		}, round: 2, reuse: true},
		{name: "round one findings", change: func(p, c *reviewRecord) { p.Outcome = reviewOutcomeFindings; p.Lineage = &reviewLineage{Round: 1} }, round: 2},
		{name: "round one findings dismissed", change: func(p, c *reviewRecord) {
			p.Outcome = reviewOutcomeFindingsDismissed
			p.Lineage = &reviewLineage{Round: 1}
		}, round: 2},
		{name: "round two ceiling", change: func(p, c *reviewRecord) { p.Lineage = &reviewLineage{Round: 2} }, round: 2, ceiling: true},
		{name: "closed ceiling", change: func(p, c *reviewRecord) {
			p.Outcome = reviewOutcomeCeilingClosed
			p.Lineage = &reviewLineage{Round: 2, ReviewedHead: fixture.baseCommit}
		}, round: 2, ceiling: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, c := prior, candidate
			if tt.change != nil {
				tt.change(&p, &c)
			}
			ptr := &p
			if tt.absent {
				ptr = nil
			}
			plan, err := decideReviewLineage(context.Background(), ptr, c, preflight.ExecGitRunner{})
			if err != nil {
				t.Fatal(err)
			}
			if plan.Lineage.Round != tt.round || plan.Reuse != tt.reuse || plan.Ceiling != tt.ceiling {
				t.Fatalf("plan = %+v", plan)
			}
			if tt.name == "ancestor blocked round two repeats" && plan.Lineage.PreviousHead != fixture.baseCommit {
				t.Fatal("blocked round lost previous data")
			}
		})
	}
}

func reviewPolicyForLineage() roundconfig.PrePRReview {
	return roundconfig.PrePRReview{Provider: "codex", Source: "project"}
}
func commitLineageFile(t *testing.T, fixture reviewCommandFixture, name, content string) string {
	t.Helper()
	mustWrite(t, filepath.Join(fixture.repository, name), content)
	gittest.Run(t, fixture.repository, "add", "--", name)
	gittest.Run(t, fixture.repository, "commit", "-m", "test: advance review candidate")
	return strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))
}

type lineageReviewRunner struct {
	reviewCommandRunner
	prompts []string
}

func (runner *lineageReviewRunner) RunPrepared(ctx context.Context, request agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	runner.prompts = append(runner.prompts, request.Prompt)
	return runner.reviewCommandRunner.RunPrepared(ctx, request, sink)
}
func newLineageRunner(answers ...string) *lineageReviewRunner {
	runner := &lineageReviewRunner{}
	for _, answer := range answers {
		runner.results = append(runner.results, reviewCommandRunResult{result: agent.ExecuteResult{Message: answer, StopReason: "end_turn"}})
	}
	return runner
}
func runLineageReview(t *testing.T, fixture reviewCommandFixture) (int, reviewRecord, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"review", "--base", fixture.baseCommit}, &stdout, &stderr)
	var record reviewRecord
	if err := json.Unmarshal(stdout.Bytes(), &record); err != nil {
		t.Fatalf("decode %q: %v (stderr %q)", stdout.String(), err, stderr.String())
	}
	return code, record, stderr.String()
}
func disposeLineageFinding(t *testing.T, args ...string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := runCLI(t, append([]string{"review", "dispose", "F1"}, args...), &stdout, &stderr)
	if code != exitOK {
		t.Fatalf("dispose exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
}
func TestReviewRoundTwoPromptCarriesTheDeltaAndRoundOneFindings(t *testing.T) {
	runner := newLineageRunner("Findings:\n- review.txt:2 Failure: the first defect remains\n- review.txt:1 Failure: the second defect remains\n- review.txt:2 Failure: unresolved defect\n- outside.txt:10 Failure: unanchored defect", "No findings.")
	fixture := newReviewCommandFixture(t, "codex", runner)
	code, first, _ := runLineageReview(t, fixture)
	if code != exitRunFailed {
		t.Fatalf("first exit=%d record=%+v", code, first)
	}
	disposeLineageFinding(t, "--dismiss", "--evidence", "operator evidence")
	fix := commitLineageFile(t, fixture, "fixed.txt", "fixed after round one\n")
	var disposeOut, disposeErr bytes.Buffer
	if code := runCLI(t, []string{"review", "dispose", "F2", "--fixed-by", fix}, &disposeOut, &disposeErr); code != exitOK {
		t.Fatalf("dispose F2=%d: %s", code, disposeErr.String())
	}
	commitLineageFile(t, fixture, "delta.txt", "only added after round one\n")
	code, second, stderr := runLineageReview(t, fixture)
	if code != exitOK || stderr != "" {
		t.Fatalf("second exit=%d record=%+v stderr=%q", code, second, stderr)
	}
	prompt := runner.prompts[1]
	for _, want := range []string{"only added after round one", "Previous head: " + first.HeadCommit, "This is round 2 of 2 of this review.", "F1 | stands | dismissed: operator evidence", "F2 | stands | fixed by " + fix, "F3 | stands | no disposition", "F4 | dismissed-by-validation (unanchored): anchor names no line of the candidate diff | no disposition", "--- BEGIN ROUND 2 DELTA DIFF ---", "--- END ROUND 2 DELTA DIFF ---"} {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing %q:\n%s", want, prompt)
		}
	}
	if strings.Contains(prompt, "the reviewer receives this changed line") || strings.Contains(prompt, "BEGIN CANDIDATE DIFF") {
		t.Fatalf("round two contains full diff: %s", prompt)
	}
	if second.Lineage == nil || second.Lineage.Round != 2 || second.Lineage.PreviousHead != first.HeadCommit || !reflect.DeepEqual(second.Lineage.PreviousFindings, first.FindingItems) || len(second.Lineage.PreviousDispositions) != 2 {
		t.Fatalf("second lineage=%+v", second.Lineage)
	}
	if runner.preparedCalls != 2 || runner.endCalls != 1 {
		t.Fatalf("session counts: %+v", runner)
	}
}
func TestReviewRoundTwoAnchorsAgainstTheFullCandidateDiff(t *testing.T) {
	runner := newLineageRunner("No findings.", "Findings:\n- review.txt:2 Failure: original defect remains")
	fixture := newReviewCommandFixture(t, "codex", runner)
	runLineageReview(t, fixture)
	commitLineageFile(t, fixture, "delta.txt", "new delta\n")
	code, record, stderr := runLineageReview(t, fixture)
	if code != exitRunFailed || len(record.FindingItems) != 1 || record.FindingItems[0].Validation.Status != "stands" {
		t.Fatalf("exit=%d record=%+v stderr=%q", code, record, stderr)
	}
	if strings.Contains(runner.prompts[1], "the reviewer receives this changed line") {
		t.Fatal("original diff leaked into delta prompt")
	}
}
func roundTwoLineageFixture(t *testing.T) (reviewCommandFixture, *lineageReviewRunner, reviewRecord) {
	t.Helper()
	runner := newLineageRunner("No findings.", "Findings:\n- review.txt:2 Failure: still broken")
	fixture := newReviewCommandFixture(t, "codex", runner)
	runLineageReview(t, fixture)
	commitLineageFile(t, fixture, "delta.txt", "round two\n")
	code, record, stderr := runLineageReview(t, fixture)
	if code != exitRunFailed || record.Lineage.Round != 2 {
		t.Fatalf("round two exit=%d record=%+v stderr=%q", code, record, stderr)
	}
	return fixture, runner, record
}
func TestReviewCeilingBlocksWithoutCallingTheReviewer(t *testing.T) {
	fixture, runner, second := roundTwoLineageFixture(t)
	head := commitLineageFile(t, fixture, "fix.txt", "fix candidate\n")
	dir := reviewCheckoutDir(fixture.artifactDir, fixture.repository)
	beforeRecord := mustReadReviewCheckoutFile(t, dir, reviewRecordFileName)
	beforeAnswer := mustReadReviewCheckoutFile(t, dir, reviewAnswerFileName)
	runner.preparedCalls = 0
	runner.prepareCalls = 0
	runner.probeCalls = 0
	runner.prompts = nil
	code, blocked, stderr := runLineageReview(t, fixture)
	wantReason := "review round ceiling reached: round 2 reviewed " + second.HeadCommit + "; dispose finding F1 with roundfix review dispose"
	if code != exitPreflight || blocked.Outcome != reviewOutcomeBlocked || blocked.HeadCommit != head || blocked.Reason != wantReason || stderr != "roundfix: review blocked: "+wantReason+"\n" {
		t.Fatalf("exit=%d record=%+v stderr=%q", code, blocked, stderr)
	}
	if blocked.Lineage.Round != 2 || blocked.Lineage.ReviewedHead != second.HeadCommit || blocked.AnswerPath != "" || blocked.Validation != nil {
		t.Fatalf("blocked transcript=%+v", blocked)
	}
	if runner.preparedCalls != 0 || runner.prepareCalls != 0 || runner.probeCalls != 0 {
		t.Fatalf("reviewer called: %+v", runner)
	}
	if !bytes.Equal(beforeRecord, mustReadReviewCheckoutFile(t, dir, reviewRecordFileName)) || !bytes.Equal(beforeAnswer, mustReadReviewCheckoutFile(t, dir, reviewAnswerFileName)) {
		t.Fatal("blocked ceiling changed record or answer")
	}
}
func TestReviewCeilingClosesOnDispositions(t *testing.T) {
	t.Run("fixed by contained commit", func(t *testing.T) {
		fixture, runner, second := roundTwoLineageFixture(t)
		fix := commitLineageFile(t, fixture, "fix.txt", "fix\n")
		disposeLineageFinding(t, "--fixed-by", fix)
		assertLineageCeilingClosed(t, fixture, runner, second)
	})
	t.Run("dismissed with evidence", func(t *testing.T) {
		fixture, runner, second := roundTwoLineageFixture(t)
		disposeLineageFinding(t, "--dismiss", "--evidence", "checked against contract")
		commitLineageFile(t, fixture, "fix.txt", "final correction\n")
		assertLineageCeilingClosed(t, fixture, runner, second)
	})
	t.Run("fixed by commit outside current head", func(t *testing.T) {
		fixture, runner, second := roundTwoLineageFixture(t)
		gittest.Run(t, fixture.repository, "checkout", "-b", "fix/other-candidate")
		fix := commitLineageFile(t, fixture, "fix.txt", "other fix\n")
		disposeLineageFinding(t, "--fixed-by", fix)
		gittest.Run(t, fixture.repository, "checkout", "feature/review")
		commitLineageFile(t, fixture, "different.txt", "different candidate\n")
		dir := reviewCheckoutDir(fixture.artifactDir, fixture.repository)
		before := mustReadReviewCheckoutFile(t, dir, reviewRecordFileName)
		runner.preparedCalls = 0
		code, record, _ := runLineageReview(t, fixture)
		if code != exitPreflight || record.Outcome != reviewOutcomeBlocked || record.Lineage.ReviewedHead != second.HeadCommit || runner.preparedCalls != 0 {
			t.Fatalf("exit=%d record=%+v calls=%d", code, record, runner.preparedCalls)
		}
		if !bytes.Equal(before, mustReadReviewCheckoutFile(t, dir, reviewRecordFileName)) {
			t.Fatal("blocked fixed disposition changed prior record")
		}
	})
}
func assertLineageCeilingClosed(t *testing.T, fixture reviewCommandFixture, runner *lineageReviewRunner, second reviewRecord) {
	t.Helper()
	dir := reviewCheckoutDir(fixture.artifactDir, fixture.repository)
	answer := mustReadReviewCheckoutFile(t, dir, reviewAnswerFileName)
	ledger := mustReadReviewCheckoutFile(t, fixture.artifactDir, reviewDispositionLedgerFileName)
	runner.preparedCalls = 0
	runner.prepareCalls = 0
	runner.probeCalls = 0
	code, record, stderr := runLineageReview(t, fixture)
	if code != exitOK || stderr != "" || record.Outcome != reviewOutcomeCeilingClosed || record.Lineage.Round != 2 || record.Lineage.ReviewedHead != second.HeadCommit || record.Findings != second.Findings || !reflect.DeepEqual(record.FindingItems, second.FindingItems) || len(record.Dispositions) != 1 {
		t.Fatalf("exit=%d record=%+v stderr=%q", code, record, stderr)
	}
	persisted, err := readReviewRecord(filepath.Join(dir, reviewRecordFileName))
	if err != nil || !reflect.DeepEqual(persisted, record) {
		t.Fatalf("persisted=%+v err=%v", persisted, err)
	}
	if runner.preparedCalls != 0 || runner.prepareCalls != 0 || runner.probeCalls != 0 {
		t.Fatal("ceiling called reviewer")
	}
	if !bytes.Equal(answer, mustReadReviewCheckoutFile(t, dir, reviewAnswerFileName)) || !bytes.Equal(ledger, mustReadReviewCheckoutFile(t, fixture.artifactDir, reviewDispositionLedgerFileName)) {
		t.Fatal("ceiling changed answer or ledger")
	}
	// Another descendant closes on the same round-2 head rather than opening a reviewer.
	commitLineageFile(t, fixture, "later.txt", "later final correction\n")
	code, later, _ := runLineageReview(t, fixture)
	if code != exitOK || later.Lineage.ReviewedHead != second.HeadCommit || runner.preparedCalls != 0 {
		t.Fatalf("later ceiling=%+v code=%d", later, code)
	}
}
func TestDeliveryReviewResultAdvancesACeilingClosedRecord(t *testing.T) {
	for _, tt := range []struct {
		outcome reviewOutcome
		want    delivery.ReviewOutcome
	}{{reviewOutcomeCeilingClosed, delivery.ReviewOutcomeReviewed}, {reviewOutcomeBlocked, delivery.ReviewOutcomeBlocked}} {
		t.Run(string(tt.outcome), func(t *testing.T) {
			record := reviewRecord{HeadCommit: "head", Outcome: tt.outcome, Reason: "original reason"}
			result, err := deliveryReviewResult(record, "head")
			if err != nil || result.Outcome != tt.want {
				t.Fatalf("result=%+v err=%v", result, err)
			}
			if tt.outcome == reviewOutcomeBlocked && result.Reason != record.Reason {
				t.Fatal("blocked reason changed")
			}
		})
	}
}

func TestReviewBlockedRoundTwoRepeatsItsOriginalDelta(t *testing.T) {
	runner := newLineageRunner("Findings:\n- review.txt:2 Failure: unresolved", "invalid reviewer verdict", "No findings.")
	fixture := newReviewCommandFixture(t, "codex", runner)
	_, first, _ := runLineageReview(t, fixture)
	commitLineageFile(t, fixture, "delta.txt", "first fix attempt\n")
	code, blocked, _ := runLineageReview(t, fixture)
	if code != exitPreflight || blocked.Outcome != reviewOutcomeBlocked || blocked.Lineage.Round != 2 {
		t.Fatalf("blocked=%+v code=%d", blocked, code)
	}
	commitLineageFile(t, fixture, "retry.txt", "second fix attempt\n")
	code, retry, stderr := runLineageReview(t, fixture)
	if code != exitOK || stderr != "" || retry.Lineage.Round != 2 || retry.Lineage.PreviousHead != first.HeadCommit || !reflect.DeepEqual(retry.Lineage.PreviousFindings, first.FindingItems) {
		t.Fatalf("retry=%+v code=%d stderr=%q", retry, code, stderr)
	}
	prompt := runner.prompts[2]
	for _, want := range []string{"first fix attempt", "second fix attempt", "Previous head: " + first.HeadCommit, "F1 | stands | no disposition"} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("retry prompt missing %q: %s", want, prompt)
		}
	}
}

func TestReviewCeilingClosedRecordRequiresDispositionEvidence(t *testing.T) {
	fixture := newReviewCommandFixture(t, "codex", &reviewCommandRunner{})
	head := commitLineageFile(t, fixture, "fix.txt", "fix\n")
	base := newReviewRecord(fixture.repository, fixture.baseCommit, head, reviewPolicyForLineage(), reviewOutcomeCeilingClosed)
	base.Findings = "review.txt:2 Failure: broken"
	base.FindingItems = []reviewFinding{{ID: "F1", Text: base.Findings}}
	base.Lineage = &reviewLineage{Round: 2, ReviewedHead: fixture.headCommit}
	base.Dispositions = []reviewFindingDisposition{{Repository: fixture.repository, HeadCommit: fixture.headCommit, Finding: "F1", Text: base.Findings, Disposition: "fixed", FixedBy: head}}
	if err := validateReviewRecord(base); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name   string
		change func(*reviewRecord)
	}{
		{"missing lineage", func(r *reviewRecord) { r.Lineage = nil }},
		{"wrong round", func(r *reviewRecord) { r.Lineage.Round = 1 }},
		{"same reviewed head", func(r *reviewRecord) { r.Lineage.ReviewedHead = r.HeadCommit }},
		{"missing disposition", func(r *reviewRecord) { r.Dispositions = nil }},
		{"missing fixed commit", func(r *reviewRecord) { r.Dispositions[0].FixedBy = "" }},
		{"dismissal without evidence", func(r *reviewRecord) { r.Dispositions[0].Disposition = "dismissed"; r.Dispositions[0].FixedBy = "" }},
		{"disposition wrong head", func(r *reviewRecord) { r.Dispositions[0].HeadCommit = head }},
		{"disposition wrong text", func(r *reviewRecord) { r.Dispositions[0].Text = "other" }},
		{"duplicate disposition", func(r *reviewRecord) { r.Dispositions = append(r.Dispositions, r.Dispositions[0]) }},
		{"findings without items", func(r *reviewRecord) { r.FindingItems = nil; r.Dispositions = nil }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record := base
			lineage := *base.Lineage
			record.Lineage = &lineage
			record.Dispositions = append([]reviewFindingDisposition{}, base.Dispositions...)
			tt.change(&record)
			if err := validateReviewRecord(record); err == nil {
				t.Fatalf("invalid record accepted: %+v", record)
			}
		})
	}
}
