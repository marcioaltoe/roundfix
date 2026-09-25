package daemon

import (
	"context"
	"strings"
	"testing"

	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

// Suite: QA report settlement truthfulness
// Invariant: a QA Task settles completed only when its report records at least one QA row and is otherwise eligible.
// Boundary IN: mechanical seed rendering, report parsing, eligibility, and Daemon Task settlement.
// Boundary OUT: maintainer command presentation, owned by internal/cli tests.

func TestMechanicalQAReportSeedStartsPending(t *testing.T) {
	t.Parallel()
	content, err := mechanicalQAReportContent(speccheck.MechanicalResult{}, spec.AuditorEvidence{})
	if err != nil {
		t.Fatalf("mechanicalQAReportContent: %v", err)
	}
	if !strings.Contains(string(content), "\nverdict: "+spec.VerdictPending+"\n") {
		t.Fatalf("mechanical seed does not carry verdict %q:\n%s", spec.VerdictPending, content)
	}
}

func TestQASettlementRefusesAnUntouchedSeed(t *testing.T) {
	t.Parallel()
	result, fixture, runner, reason := runQAReportSettlementTest(t, nil)

	if !strings.Contains(runner.qaSeed, "\nverdict: "+spec.VerdictPending+"\n") {
		t.Fatalf("untouched mechanical seed does not carry pending verdict:\n%s", runner.qaSeed)
	}
	assertQASettlement(t, result, fixture, spec.VerdictPending, spec.StatusFailed, "QA verdict pending", reason)
}

func TestQASettlementRefusesAHollowPass(t *testing.T) {
	t.Parallel()
	result, fixture, _, reason := runQAReportSettlementTest(t, func(*taskCycleFixture) string {
		return "---\nverdict: pass\n---\n\n# QA Report\n\n## Results\n\n| # | Status | Evidence |\n| - | --- | --- |\n"
	})

	wantReason := `QA verdict pass not accepted: newest QA Report verdict is "pass" but records no QA row`
	assertQASettlement(t, result, fixture, spec.VerdictPass, spec.StatusFailed, wantReason, reason)
}

func TestQASettlementAcceptsAPassWithAResultsRow(t *testing.T) {
	t.Parallel()
	result, fixture, _, reason := runQAReportSettlementTest(t, func(*taskCycleFixture) string {
		return "---\nverdict: pass\n---\n\n# QA Report\n\n## Results\n\n| # | Status | Evidence |\n| - | --- | --- |\n| R01 | pass | observed CLI output |\n"
	})

	assertQASettlement(t, result, fixture, spec.VerdictPass, spec.StatusCompleted, "", reason)
}

func runQAReportSettlementTest(t *testing.T, report func(*taskCycleFixture) string) (TaskCycleResult, *taskCycleFixture, *taskFakeRunner, string) {
	t.Helper()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", title: "Build the feature"}})
	reportRel := qaReportRelPathForTest()
	fixture.worktree.snapshots = [][]string{nil, {"src/one.go"}, {"src/one.go"}, {"src/one.go", reportRel}}
	runner := &taskFakeRunner{
		calls:          fixture.calls,
		gitRoot:        fixture.gitRoot,
		store:          fixture.store,
		statusByTask:   map[string]spec.Status{"task_01": spec.StatusCompleted},
		writeLogs:      true,
		preserveQASeed: report == nil,
	}
	if report != nil {
		runner.qaReport = report(fixture)
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.qaPlan())
	if err != nil {
		t.Fatalf("TaskCycle: %v", err)
	}
	var reason string
	for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonTask) {
		payload := eventPayloadMap(t, event)
		if event.ReviewIssue == fixture.graph.QATaskID && payload["phase"] == "settled" {
			reason, _ = payload["reason"].(string)
			break
		}
	}
	return result, fixture, runner, reason
}

func assertQASettlement(t *testing.T, result TaskCycleResult, fixture *taskCycleFixture, verdict string, status spec.Status, wantReason string, gotReason string) {
	t.Helper()
	if result.QAVerdict != verdict || result.QAAccepted != (status == spec.StatusCompleted) {
		t.Fatalf("QA result = verdict %q accepted %t, want verdict %q accepted %t", result.QAVerdict, result.QAAccepted, verdict, status == spec.StatusCompleted)
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, fixture.graph.QATaskID); got != string(status) {
		t.Fatalf("QA Task status = %q, want %q", got, status)
	}
	if gotReason != wantReason {
		t.Fatalf("QA Task reason = %q, want %q", gotReason, wantReason)
	}
}
