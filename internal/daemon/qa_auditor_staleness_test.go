// Suite: QA auditor staleness at the Delivery Base.
// Invariant: stale auditors warn exactly once while current and unknown auditors leave the gate unchanged.
// Boundary IN: a real Git repository, the QA step, report materialization, and the Run Event Journal sink.
// Boundary OUT: Git default-branch detection, owned by qa_every_run_audit_test.go.
package daemon

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/agent"
	"roundfix/internal/app"
	"roundfix/internal/gittest"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

type qaAuditorStalenessFixture struct {
	taskFixture  *taskCycleFixture
	plan         TaskPlan
	deliveryBase string
	baseParent   string
}

func newQAAuditorStalenessFixture(t *testing.T, candidateCommits int) *qaAuditorStalenessFixture {
	t.Helper()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	fixture.qaPlan()
	commitTaskFixtureSource(t, fixture.gitRoot, "test: seed QA auditor fixture")
	deliveryBase := strings.TrimSpace(gittest.Run(t, fixture.gitRoot, "rev-parse", "HEAD"))
	baseParent := strings.TrimSpace(gittest.Run(t, fixture.gitRoot, "rev-parse", "HEAD^"))
	gittest.Run(t, fixture.gitRoot, "update-ref", "refs/remotes/origin/main", deliveryBase)
	gittest.Run(t, fixture.gitRoot, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/main")
	gittest.Run(t, fixture.gitRoot, "switch", "-c", "feature/auditor-staleness")
	for index := 1; index <= candidateCommits; index++ {
		path := filepath.Join(fixture.gitRoot, fmt.Sprintf("candidate-%02d.txt", index))
		mustWriteForTest(t, path, fmt.Sprintf("candidate %d\n", index))
		gittest.Run(t, fixture.gitRoot, "add", filepath.Base(path))
		gittest.Run(t, fixture.gitRoot, "commit", "-m", fmt.Sprintf("feat: candidate %d", index), "-m",
			"Roundfix-Spec: "+taskCycleSlug+"\nRoundfix-Task: task_01")
	}

	return &qaAuditorStalenessFixture{
		taskFixture:  fixture,
		plan:         fixture.plan(),
		deliveryBase: deliveryBase,
		baseParent:   baseParent,
	}
}

func (fixture *qaAuditorStalenessFixture) engine(t *testing.T, runner agent.Runner, verifier Verifier, binary app.AuditingBinary, result speccheck.MechanicalResult) *Engine {
	t.Helper()
	engine, err := NewEngine(Dependencies{
		Runner:          runner,
		Verifier:        verifier,
		Committer:       &engineFakeCommitter{calls: fixture.taskFixture.calls},
		Pusher:          fixture.taskFixture.pusher,
		Source:          fixture.taskFixture.source,
		Runs:            fixture.taskFixture.store,
		Worktree:        fixture.taskFixture.worktree,
		MechanicalStage: &fakeQAMechanicalStage{result: result},
		Auditor:         func() app.AuditingBinary { return binary },
		GH:              fixture.taskFixture.github,
		Sink:            fixture.taskFixture.sink,
		Now:             taskCycleNowForTest,
		Progress:        fixture.taskFixture.progress,
	})
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}
	return engine
}

func TestSelfAuditSeedRecordsACurrentAuditor(t *testing.T) {
	t.Parallel()
	fixture := newQAAuditorStalenessFixture(t, 2)
	runner := &taskFakeRunner{calls: fixture.taskFixture.calls, gitRoot: fixture.taskFixture.gitRoot, preserveQASeed: true}
	engine := fixture.engine(t, runner, &qaGateRecordingVerifier{}, app.AuditingBinary{Version: "0.17.0", Commit: fixture.deliveryBase}, speccheck.MechanicalResult{})

	if _, err := engine.TaskCycle(context.Background(), fixture.plan); err != nil {
		t.Fatalf("TaskCycle returned error: %v", err)
	}

	want := `auditor_staleness: "current: commit ancestry: build commit does not predate the delivery base"`
	if !strings.Contains(runner.qaSeed, want) {
		t.Fatalf("seeded QA Report omitted %q:\n%s", want, runner.qaSeed)
	}
}

func TestStaleAuditorWarnsAndTheGateProceeds(t *testing.T) {
	t.Parallel()
	fixture := newQAAuditorStalenessFixture(t, 1)
	runner := &taskFakeRunner{calls: fixture.taskFixture.calls, gitRoot: fixture.taskFixture.gitRoot, preserveQASeed: true}
	verifier := &qaGateRecordingVerifier{}
	engine := fixture.engine(t, runner, verifier, app.AuditingBinary{Version: "0.17.0", Commit: fixture.baseParent}, speccheck.MechanicalResult{})
	fixture.plan.RepositoryVerification = "true"

	if _, err := engine.TaskCycle(context.Background(), fixture.plan); err != nil {
		t.Fatalf("TaskCycle returned error: %v", err)
	}

	if len(verifier.requests) != 1 {
		t.Fatalf("repository Verification ran %d times, want once", len(verifier.requests))
	}
	if len(runner.qaPrompts) != 1 {
		t.Fatalf("QA Agent calls = %d, want one", len(runner.qaPrompts))
	}
	if !strings.Contains(runner.qaSeed, "## Auditor staleness warning") {
		t.Fatalf("seeded QA Report omitted the staleness warning:\n%s", runner.qaSeed)
	}
	if strings.Contains(runner.qaSeed, "precondition_check:") {
		t.Fatalf("stale auditor became a precondition refusal:\n%s", runner.qaSeed)
	}

	events := auditorStalenessEvents(t, fixture.taskFixture.sink)
	if len(events) != 1 {
		t.Fatalf("auditor staleness events = %d, want one: %+v", len(events), events)
	}
	payload := eventPayloadMap(t, events[0])
	wantAction := auditorStalenessAction(fixture.deliveryBase)
	if payload["delivery_base"] != fixture.deliveryBase || payload["action"] != wantAction || payload["report"] != runner.qaReportPath {
		t.Fatalf("auditor staleness payload = %#v, want base %q action %q report %q", payload, fixture.deliveryBase, wantAction, runner.qaReportPath)
	}
	if payload["auditor_staleness"] != "stale: commit ancestry: build commit predates the delivery base" {
		t.Fatalf("auditor staleness payload = %#v", payload)
	}
}

func TestStaleAuditorWarningRidesARefusedGate(t *testing.T) {
	t.Parallel()
	fixture := newQAAuditorStalenessFixture(t, 1)
	runner := &taskFakeRunner{calls: fixture.taskFixture.calls, gitRoot: fixture.taskFixture.gitRoot}
	prdPath := filepath.Join(fixture.plan.Spec.Dir, "_prd.md")
	prd := readFileForTest(t, prdPath)
	head, _, found := strings.Cut(prd, "## Project Constraints")
	if !found {
		t.Fatalf("fixture PRD has no Project Constraints to remove:\n%s", prd)
	}
	mustWriteForTest(t, prdPath, head)
	engine := fixture.engine(t, runner, &qaGateRecordingVerifier{}, app.AuditingBinary{Version: "0.17.0", Commit: fixture.baseParent}, speccheck.MechanicalResult{})
	engine.deps.MechanicalStage = SpecCheckQAMechanicalStage{}

	result, err := engine.TaskCycle(context.Background(), fixture.plan)
	if err != nil {
		t.Fatalf("TaskCycle returned error: %v", err)
	}
	report := readQAAuditorReport(t, fixture.plan, result.QAReportPath)
	if !strings.Contains(report, `precondition_check: "`+speccheck.GatePreconditionCheck+`"`) {
		t.Fatalf("refused report lost its precondition check:\n%s", report)
	}
	if !strings.Contains(report, "## Auditor staleness warning") {
		t.Fatalf("refused report omitted the staleness warning:\n%s", report)
	}
	if events := auditorStalenessEvents(t, fixture.taskFixture.sink); len(events) != 1 {
		t.Fatalf("auditor staleness events = %d, want one: %+v", len(events), events)
	}
}

func TestCurrentAuditorEmitsNoStalenessWarning(t *testing.T) {
	t.Parallel()
	fixture := newQAAuditorStalenessFixture(t, 1)
	runner := &taskFakeRunner{calls: fixture.taskFixture.calls, gitRoot: fixture.taskFixture.gitRoot, preserveQASeed: true}
	engine := fixture.engine(t, runner, &qaGateRecordingVerifier{}, app.AuditingBinary{Version: "0.17.0", Commit: fixture.deliveryBase}, speccheck.MechanicalResult{})

	if _, err := engine.TaskCycle(context.Background(), fixture.plan); err != nil {
		t.Fatalf("TaskCycle returned error: %v", err)
	}
	assertNoAuditorStalenessWarning(t, fixture, runner)
}

func TestUnknownAuditorEmitsNoStalenessWarning(t *testing.T) {
	t.Parallel()
	fixture := newQAAuditorStalenessFixture(t, 1)
	runner := &taskFakeRunner{calls: fixture.taskFixture.calls, gitRoot: fixture.taskFixture.gitRoot, preserveQASeed: true}
	engine := fixture.engine(t, runner, &qaGateRecordingVerifier{}, app.AuditingBinary{Version: "0.17.0", Commit: "0123456789abcdef0123456789abcdef01234567"}, speccheck.MechanicalResult{})

	if _, err := engine.TaskCycle(context.Background(), fixture.plan); err != nil {
		t.Fatalf("TaskCycle returned error: %v", err)
	}
	assertNoAuditorStalenessWarning(t, fixture, runner)
}

func auditorStalenessEvents(t *testing.T, sink *captureEventSink) []runevent.RunEvent {
	t.Helper()
	events := make([]runevent.RunEvent, 0, 1)
	for _, event := range taskEventsOfKind(sink, runevent.KindDaemonQA) {
		if eventPayloadMap(t, event)["phase"] == "auditor_staleness" {
			events = append(events, event)
		}
	}
	return events
}

func assertNoAuditorStalenessWarning(t *testing.T, fixture *qaAuditorStalenessFixture, runner *taskFakeRunner) {
	t.Helper()
	if events := auditorStalenessEvents(t, fixture.taskFixture.sink); len(events) != 0 {
		t.Fatalf("auditor staleness events = %d, want none: %+v", len(events), events)
	}
	if strings.Contains(runner.qaSeed, "## Auditor staleness warning") {
		t.Fatalf("seeded QA Report contains a staleness warning:\n%s", runner.qaSeed)
	}
	if len(runner.qaPrompts) != 1 {
		t.Fatalf("QA Agent calls = %d, want one", len(runner.qaPrompts))
	}
}

func readQAAuditorReport(t *testing.T, plan TaskPlan, reportPath string) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(plan.WorkDir, filepath.FromSlash(reportPath)))
	if err != nil {
		t.Fatalf("read QA Report %q: %v", reportPath, err)
	}
	return string(content)
}
