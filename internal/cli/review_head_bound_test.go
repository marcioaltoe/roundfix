// Suite: head-bound pre-PR review verdicts.
// Invariant: findings stand for one repository, base, head, and provider until every exact finding is dismissed.
// Boundary IN: public review and dispose commands, Artifact Directory records and ledgers, and delivery result mapping.
// Boundary OUT: real ACP adapters and Delivery Queue stage transitions.
package cli

import (
	"bytes"
	"strings"
	"testing"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
	"roundfix/internal/gittest"
)

func TestReviewReportsFindingsDismissedWithoutAskingTheReviewer(t *testing.T) {
	t.Parallel()
	const findings = "- review.txt:2: first\n- review.txt:2: second"
	fixture, runner := recordHeadBoundFindings(t, "codex", findings)
	for _, findingID := range []string{"F1", "F2"} {
		code, _, stderr := runReviewDispose(t, findingID, "--dismiss", "--evidence", "verified by the candidate contract")
		if code != exitOK || stderr != "" {
			t.Fatalf("dismiss %s exit=%d stderr=%q, want exit=0", findingID, code, stderr)
		}
	}
	resetReviewCommandRunner(runner, agent.ExecuteResult{})

	code, record, stderr := fixture.run(t)

	if code != exitOK || stderr != "" {
		t.Fatalf("reused review exit=%d stderr=%q, want exit=0", code, stderr)
	}
	if record.Outcome != reviewOutcomeFindingsDismissed || !record.Reused || len(record.Dispositions) != 2 {
		t.Fatalf("reused review record = %+v, want findings-dismissed with two dispositions", record)
	}
	for _, disposition := range record.Dispositions {
		if disposition.Disposition != "dismissed" {
			t.Fatalf("reused disposition = %+v, want dismissed", disposition)
		}
	}
	assertNoReviewAgentActivity(t, runner)

	code, _, stderr = runReviewDispose(t, "F1", "--dismiss", "--evidence", "duplicate")
	if code != exitPreflight || !strings.Contains(stderr, "already has a disposition") {
		t.Fatalf("dispose against findings-dismissed exit=%d stderr=%q, want duplicate-disposition refusal", code, stderr)
	}
}

func TestReviewKeepsStandingFindingsWithoutAskingTheReviewer(t *testing.T) {
	t.Parallel()
	const findings = "- review.txt:2: dismissed\n- review.txt:2: standing"
	fixture, runner := recordHeadBoundFindings(t, "codex", findings)
	code, _, stderr := runReviewDispose(t, "F1", "--dismiss", "--evidence", "not reachable")
	if code != exitOK || stderr != "" {
		t.Fatalf("dismiss F1 exit=%d stderr=%q, want exit=0", code, stderr)
	}
	resetReviewCommandRunner(runner, agent.ExecuteResult{})

	code, record, stderr := fixture.run(t)

	if code != exitRunFailed || record.Outcome != reviewOutcomeFindings || !record.Reused {
		t.Fatalf("standing review exit=%d record=%+v, want reused findings", code, record)
	}
	if len(record.Dispositions) != 1 || record.Dispositions[0].Finding != "F1" {
		t.Fatalf("standing review dispositions = %+v, want only F1", record.Dispositions)
	}
	if !strings.Contains(stderr, "F2") || strings.Contains(stderr, "F1") {
		t.Fatalf("standing review stderr=%q, want only missing F2", stderr)
	}
	assertNoReviewAgentActivity(t, runner)
}

func TestReviewIgnoresADismissalOfDifferentText(t *testing.T) {
	t.Parallel()
	fixture, runner := recordHeadBoundFindings(t, "codex", "- review.txt:2: original text")
	code, _, stderr := runReviewDispose(t, "F1", "--dismiss", "--evidence", "applies only to the original text")
	if code != exitOK || stderr != "" {
		t.Fatalf("dismiss original finding exit=%d stderr=%q, want exit=0", code, stderr)
	}
	record := dispositionReviewRecord(fixture.repository, fixture.headCommit, "- review.txt:2: changed text")
	record.BaseCommit = fixture.baseCommit
	writeDispositionRecord(t, fixture, record)
	resetReviewCommandRunner(runner, agent.ExecuteResult{})

	code, record, stderr = fixture.run(t)

	if code != exitRunFailed || record.Outcome != reviewOutcomeFindings || !record.Reused {
		t.Fatalf("different-text review exit=%d record=%+v, want reused standing findings", code, record)
	}
	if len(record.Dispositions) != 0 || !strings.Contains(stderr, "F1") {
		t.Fatalf("different-text dispositions=%+v stderr=%q, want no match and F1 diagnostic", record.Dispositions, stderr)
	}
	assertNoReviewAgentActivity(t, runner)
}

func TestReviewIgnoresAFixWhenClearingAHead(t *testing.T) {
	t.Parallel()
	fixture, runner := recordHeadBoundFindings(t, "codex", "- review.txt:2: fixed later")
	mustWrite(t, fixture.repository+"/fix.txt", "fix\n")
	gittest.Run(t, fixture.repository, "add", "fix.txt")
	gittest.Run(t, fixture.repository, "commit", "-m", "fix review finding")
	fixedBy := strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))
	code, _, stderr := runReviewDispose(t, "F1", "--fixed-by", fixedBy)
	if code != exitOK || stderr != "" {
		t.Fatalf("record fix exit=%d stderr=%q, want exit=0", code, stderr)
	}
	gittest.Run(t, fixture.repository, "checkout", "-B", "feature/review", fixture.headCommit)
	resetReviewCommandRunner(runner, agent.ExecuteResult{})

	code, record, stderr := fixture.run(t)

	if code != exitRunFailed || record.Outcome != reviewOutcomeFindings || !record.Reused {
		t.Fatalf("fixed review exit=%d record=%+v stderr=%q, want reused standing findings", code, record, stderr)
	}
	if len(record.Dispositions) != 1 || record.Dispositions[0].Disposition != "fixed" {
		t.Fatalf("fixed review dispositions = %+v, want recorded fix", record.Dispositions)
	}
	assertNoReviewAgentActivity(t, runner)
}

func TestReviewAsksAgainAfterTheHeadMoves(t *testing.T) {
	t.Parallel()
	fixture, runner := recordHeadBoundFindings(t, "codex", "- review.txt:2: original head")
	mustWrite(t, fixture.repository+"/moved.txt", "moved\n")
	gittest.Run(t, fixture.repository, "add", "moved.txt")
	gittest.Run(t, fixture.repository, "commit", "-m", "move review head")
	fixture.headCommit = strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))
	resetReviewCommandRunner(runner, agent.ExecuteResult{Message: "No findings", StopReason: "end_turn"})

	code, record, stderr := fixture.run(t)

	assertFreshReview(t, code, record, stderr, fixture, runner)
}

func TestReviewAsksAgainForADifferentBase(t *testing.T) {
	t.Parallel()
	fixture, runner := recordHeadBoundFindings(t, "codex", "- review.txt:2: original base")
	fixture.baseCommit = fixture.headCommit
	resetReviewCommandRunner(runner, agent.ExecuteResult{Message: "No findings", StopReason: "end_turn"})

	code, record, stderr := fixture.run(t)

	assertFreshReview(t, code, record, stderr, fixture, runner)
}

func TestReviewAsksAgainForADifferentProvider(t *testing.T) {
	t.Parallel()
	fixture, runner := recordHeadBoundFindings(t, "codex", "- review.txt:2: original provider")
	fixture.provider = "claude"
	writeReviewCommandProfileConfig(t, fixture.repository, "claude", fixture.artifactDir, "claude", "claude")
	resetReviewCommandRunner(runner, agent.ExecuteResult{Message: "No findings", StopReason: "end_turn"})

	code, record, stderr := fixture.run(t)

	assertFreshReview(t, code, record, stderr, fixture, runner)
}

func TestReviewNeverReusesACleanVerdict(t *testing.T) {
	t.Parallel()
	runner := &reviewCommandRunner{results: []reviewCommandRunResult{{
		result: agent.ExecuteResult{Message: "No findings", StopReason: "end_turn"},
	}}}
	fixture := newReviewCommandFixture(t, "codex", runner)
	code, record, stderr := fixture.run(t)
	if code != exitOK || record.Outcome != reviewOutcomeReviewed || stderr != "" {
		t.Fatalf("initial clean review exit=%d record=%+v stderr=%q", code, record, stderr)
	}
	resetReviewCommandRunner(runner, agent.ExecuteResult{Message: "No findings", StopReason: "end_turn"})

	code, record, stderr = fixture.run(t)

	assertFreshReview(t, code, record, stderr, fixture, runner)
}

func TestReviewNeverReusesABlockedVerdict(t *testing.T) {
	t.Parallel()
	runner := &reviewCommandRunner{results: []reviewCommandRunResult{{
		result: agent.ExecuteResult{Message: "ambiguous answer", StopReason: "end_turn"},
	}}}
	fixture := newReviewCommandFixture(t, "codex", runner)
	code, record, stderr := fixture.run(t)
	if code != exitPreflight || record.Outcome != reviewOutcomeBlocked || !strings.Contains(stderr, "neither") {
		t.Fatalf("initial blocked review exit=%d record=%+v stderr=%q", code, record, stderr)
	}
	resetReviewCommandRunner(runner, agent.ExecuteResult{Message: "No findings", StopReason: "end_turn"})

	code, record, stderr = fixture.run(t)

	assertFreshReview(t, code, record, stderr, fixture, runner)
}

func TestDeliveryReviewResultAdvancesDismissedFindings(t *testing.T) {
	t.Parallel()
	record := validDismissedReviewRecord()

	result, err := deliveryReviewResult(record, record.HeadCommit)

	if err != nil {
		t.Fatalf("map dismissed review result: %v", err)
	}
	if result.Outcome != delivery.ReviewOutcomeReviewed || result.Head != record.HeadCommit || result.Reason != "" {
		t.Fatalf("delivery review result = %+v, want reviewed at %s", result, record.HeadCommit)
	}
}

func TestDeliveryReviewResultParksStandingFindings(t *testing.T) {
	t.Parallel()
	record := newReviewRecord(
		"/tmp/repository",
		"base",
		"head",
		roundconfig.PrePRReview{Provider: "codex", Source: "project"},
		reviewOutcomeFindings,
	)
	record.Findings = "internal/cli/review.go:10: standing"

	result, err := deliveryReviewResult(record, record.HeadCommit)

	if err != nil {
		t.Fatalf("map standing review result: %v", err)
	}
	if result.Outcome != delivery.ReviewOutcomeFindings || result.Head != record.HeadCommit || result.Reason != record.Findings {
		t.Fatalf("delivery review result = %+v, want standing findings", result)
	}
}

func TestReviewRecordRefusesInvalidFindingsDismissed(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		mutate func(*reviewRecord)
	}{
		{name: "blank findings", mutate: func(record *reviewRecord) { record.Findings = "" }},
		{name: "reason", mutate: func(record *reviewRecord) { record.Reason = "not allowed" }},
		{name: "no items", mutate: func(record *reviewRecord) { record.FindingItems = nil }},
		{name: "missing disposition", mutate: func(record *reviewRecord) { record.Dispositions = nil }},
		{name: "fixed disposition", mutate: func(record *reviewRecord) { record.Dispositions[0].Disposition = "fixed" }},
		{name: "dismissal with fixed commit", mutate: func(record *reviewRecord) { record.Dispositions[0].FixedBy = "commit" }},
		{name: "duplicate disposition", mutate: func(record *reviewRecord) { record.Dispositions = append(record.Dispositions, record.Dispositions[0]) }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			record := validDismissedReviewRecord()
			test.mutate(&record)
			var output bytes.Buffer
			if err := writeReviewRecord(&output, record); err == nil {
				t.Fatalf("invalid findings-dismissed record wrote %q", output.String())
			}
			if output.Len() != 0 {
				t.Fatalf("invalid findings-dismissed record wrote %q", output.String())
			}
		})
	}
}

func recordHeadBoundFindings(t *testing.T, provider string, findings string) (reviewCommandFixture, *reviewCommandRunner) {
	t.Helper()
	runner := &reviewCommandRunner{results: []reviewCommandRunResult{{
		result: agent.ExecuteResult{Message: "Findings:\n" + findings, StopReason: "end_turn"},
	}}}
	fixture := newReviewCommandFixture(t, provider, runner)
	code, record, stderr := fixture.run(t)
	if code != exitRunFailed || record.Outcome != reviewOutcomeFindings || stderr != "" {
		t.Fatalf("initial findings review exit=%d record=%+v stderr=%q", code, record, stderr)
	}
	if runner.prepareCalls != 1 || runner.preparedCalls != 1 {
		t.Fatalf("initial findings review prepares=%d prompts=%d, want one each", runner.prepareCalls, runner.preparedCalls)
	}
	return fixture, runner
}

func resetReviewCommandRunner(runner *reviewCommandRunner, result agent.ExecuteResult) {
	runner.probeCalls = 0
	runner.prepareCalls = 0
	runner.preparedCalls = 0
	runner.endCalls = 0
	runner.request = agent.ExecuteRequest{}
	runner.prepareErrors = nil
	runner.results = []reviewCommandRunResult{{result: result}}
}

func assertNoReviewAgentActivity(t *testing.T, runner *reviewCommandRunner) {
	t.Helper()
	if runner.probeCalls != 0 || runner.prepareCalls != 0 || runner.preparedCalls != 0 {
		t.Fatalf("reused review Agent activity: probes=%d prepares=%d prompts=%d, want zero", runner.probeCalls, runner.prepareCalls, runner.preparedCalls)
	}
}

func assertFreshReview(
	t *testing.T,
	code int,
	record reviewRecord,
	stderr string,
	fixture reviewCommandFixture,
	runner *reviewCommandRunner,
) {
	t.Helper()
	if code != exitOK || record.Outcome != reviewOutcomeReviewed || record.Reused || stderr != "" {
		t.Fatalf("fresh review exit=%d record=%+v stderr=%q, want new reviewed verdict", code, record, stderr)
	}
	fixture.assertCandidate(t, record)
	if runner.prepareCalls != 1 || runner.preparedCalls != 1 {
		t.Fatalf("fresh review prepares=%d prompts=%d, want one each", runner.prepareCalls, runner.preparedCalls)
	}
}

func validDismissedReviewRecord() reviewRecord {
	policy := roundconfig.PrePRReview{Provider: "codex", Source: "project"}
	record := newReviewRecord("/tmp/repository", "base", "head", policy, reviewOutcomeFindingsDismissed)
	record.Findings = "internal/cli/review.go:10: dismissed"
	record.FindingItems = []reviewFinding{{ID: "F1", Text: record.Findings}}
	record.Dispositions = []reviewFindingDisposition{{
		Repository:  record.Repository,
		HeadCommit:  record.HeadCommit,
		Finding:     "F1",
		Text:        record.Findings,
		Disposition: "dismissed",
		Evidence:    "verified",
		RecordedAt:  "2026-09-28T00:00:00Z",
	}}
	return record
}
