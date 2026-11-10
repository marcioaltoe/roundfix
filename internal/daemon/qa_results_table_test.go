// Suite: QA Report Results table repair
// Invariant: one repair continues the seeded table or settles with the shared refusal.
// Boundary IN: TaskCycle with a scripted QA Agent and a real temporary Git repository
// Boundary OUT: ACP runtime, network and repository Verification (task-cycle fakes)
package daemon

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/agent"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
)

type qaResultsRunner struct {
	*taskFakeRunner

	repairs  []agent.ExecuteRequest
	onRepair func()
}

func (runner *qaResultsRunner) Run(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	if !strings.HasPrefix(req.Prompt, "QA Report Results rows lie outside the Results table.\n") {
		return runner.taskFakeRunner.Run(ctx, req, sink)
	}
	runner.repairs = append(runner.repairs, req)
	if runner.onRepair != nil {
		runner.onRepair()
	}
	return agent.ExecuteResult{LogPath: req.LogPath}, nil
}

func qaResultsFixture(t *testing.T, repair bool) (*taskCycleFixture, TaskCycleResult, *qaResultsRunner, spec.QAResultsBreak) {
	t.Helper()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	plan := fixture.qaPlan()
	report := qaReportForTest(spec.VerdictPass) + "\n## Results\n\n| # | Status | Provenance |\n| --- | --- | --- |\n| 1 | pass | Requirement 1 |\n\n| 2 | pass | Requirement 2 |\n"
	runner := &qaResultsRunner{taskFakeRunner: &taskFakeRunner{
		calls: fixture.calls, gitRoot: fixture.gitRoot, qaReport: report,
	}}
	var brk spec.QAResultsBreak
	runner.afterQA = func() {
		path, err := spec.NewestQAReport(plan.Spec.Dir)
		if err != nil {
			t.Fatal(err)
		}
		parsed, err := spec.ReadQAReportFile(path)
		if err != nil || parsed.ResultsBreak == nil {
			t.Fatalf("scripted report must have a Results break: %+v, %v", parsed, err)
		}
		brk = *parsed.ResultsBreak
	}
	if repair {
		runner.onRepair = func() {
			fixed := strings.Replace(report, "| Requirement 1 |\n\n| 2 |", "| Requirement 1 |\n| 2 |", 1)
			mustWriteForTest(t, filepath.Join(fixture.gitRoot, runner.qaReportPath), fixed)
		}
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, GitCommitter{}, GitWorktreeSnapshotter{})
	result, err := engine.TaskCycle(t.Context(), plan)
	if err != nil {
		t.Fatalf("TaskCycle: %v", err)
	}
	return fixture, result, runner, brk
}

func TestQAResultsRowAfterABlankLineGetsOneRepairTurn(t *testing.T) {
	t.Parallel()
	fixture, result, runner, brk := qaResultsFixture(t, true)
	assertQAResultsRepair(t, fixture, runner, brk)
	if !result.QAAccepted || taskStatusOnDisk(t, fixture.gitRoot, fixture.graph.QATaskID) != string(spec.StatusCompleted) {
		t.Fatalf("repaired QA was not accepted: %+v", result)
	}
	if message := strings.TrimSpace(runGitForTest(t, fixture.gitRoot, "log", "-1", "--format=%B")); message != QACommitMessage(taskCycleSlug, spec.VerdictPass) {
		t.Fatalf("QA commit message = %q", message)
	}
	content := runGitForTest(t, fixture.gitRoot, "show", "HEAD:"+result.QAReportPath)
	if !strings.Contains(content, "| 1 | pass | Requirement 1 |\n| 2 | pass | Requirement 2 |\n") {
		t.Fatalf("QA Report commit lost repaired rows:\n%s", content)
	}
}

func TestQAResultsRowAfterABlankLineAfterRepairSettlesFailed(t *testing.T) {
	t.Parallel()
	fixture, result, runner, brk := qaResultsFixture(t, false)
	assertQAResultsRepair(t, fixture, runner, brk)
	path, err := spec.NewestQAReport(fixture.graph.Spec.Dir)
	if err != nil {
		t.Fatal(err)
	}
	brk.Path = path
	wantReason := "QA verdict pass not accepted: " + brk.Error()
	var reason string
	for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonTask) {
		payload := eventPayloadMap(t, event)
		if event.ReviewIssue == fixture.graph.QATaskID && payload["phase"] == "settled" {
			reason, _ = payload["reason"].(string)
		}
	}
	assertQASettlement(t, result, fixture, spec.VerdictPass, spec.StatusFailed, wantReason, reason)
}

func TestQAPromptTellsTheGateToContinueTheSeededTable(t *testing.T) {
	t.Parallel()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	plan := fixture.qaPlan()
	runner := &taskFakeRunner{
		calls: fixture.calls, gitRoot: fixture.gitRoot,
		qaReport: qaReportForTest(spec.VerdictPass) + "\n## Results\n\n| # | Status | Provenance |\n| --- | --- | --- |\n| 1 | pass | Requirement 1 |\n",
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, GitCommitter{}, GitWorktreeSnapshotter{})
	if _, err := engine.TaskCycle(t.Context(), plan); err != nil {
		t.Fatalf("TaskCycle: %v", err)
	}
	if len(runner.qaPrompts) != 1 {
		t.Fatalf("initial QA prompts = %d, want 1", len(runner.qaPrompts))
	}
	want := "never copy, replay or cite evidence or an executable an earlier Run left.\n" +
		"Add each row of the matrix to the seeded Results table directly below its last row, with no blank line before it and the same three columns, #, Status and Provenance; a row after a blank line is outside the table, and the report is refused naming that line.\n"
	if !strings.Contains(runner.qaPrompts[0], want) {
		t.Fatalf("QA prompt lacks ordered table instructions:\n%s", runner.qaPrompts[0])
	}
}

func TestQAResultsAndEvidenceLinksShareOneRepairTurn(t *testing.T) {
	t.Parallel()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	plan := fixture.qaPlan()
	const link = "evidence/missing.txt"
	report := qaReportForTest(spec.VerdictPass) + "\n## Results\n\n| # | Status | Provenance |\n| --- | --- | --- |\n| 1 | pass | Requirement 1 |\n\n| 2 | pass | [Requirement 2](" + link + ") |\n"
	runner := &qaEvidenceLinkRunner{taskFakeRunner: &taskFakeRunner{
		calls: fixture.calls, gitRoot: fixture.gitRoot, qaReport: report,
	}}
	var brk spec.QAResultsBreak
	runner.afterQA = func() {
		path := filepath.Join(fixture.gitRoot, runner.qaReportPath)
		results := spec.ReadQAResults(path, []byte(report))
		if results.Break == nil {
			t.Fatal("scripted report has no Results break")
		}
		brk = *results.Break
	}
	runner.onRepair = func() {
		path := filepath.Join(fixture.gitRoot, runner.qaReportPath)
		mustWriteForTest(t, path, strings.Replace(report, "| Requirement 1 |\n\n| 2 |", "| Requirement 1 |\n| 2 |", 1))
		evidence := filepath.Join(filepath.Dir(path), link)
		if err := os.MkdirAll(filepath.Dir(evidence), 0o755); err != nil {
			t.Fatal(err)
		}
		mustWriteForTest(t, evidence, "this Run execution\n")
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, GitCommitter{}, GitWorktreeSnapshotter{})
	result, err := engine.TaskCycle(t.Context(), plan)
	if err != nil {
		t.Fatalf("TaskCycle: %v", err)
	}
	if !result.QAAccepted || len(runner.repairs) != 1 || len(runner.qaPrompts) != 1 {
		t.Fatalf("shared repair: result=%+v, repair turns=%d, QA prompts=%d", result, len(runner.repairs), len(runner.qaPrompts))
	}
	want := qaEvidenceLinkRepairPrompt(brk.Path, []spec.UnresolvedEvidenceLink{{Destination: link, Line: brk.Line, Reason: "missing"}}) + qaResultsRepairPrompt(brk.Path, brk)
	if got := runner.repairs[0].Prompt; got != want {
		t.Fatalf("combined repair prompt = %q, want %q", got, want)
	}
	path := filepath.ToSlash(filepath.Join(filepath.Dir(result.QAReportPath), link))
	if got := runGitForTest(t, fixture.gitRoot, "show", "HEAD:"+path); got != "this Run execution\n" {
		t.Fatalf("shared repair evidence was not committed: %q", got)
	}
}

func assertQAResultsRepair(t *testing.T, fixture *taskCycleFixture, runner *qaResultsRunner, brk spec.QAResultsBreak) {
	t.Helper()
	if len(runner.repairs) != 1 || len(runner.qaPrompts) != 1 {
		t.Fatalf("repair turns = %d, QA prompts = %d, want 1 each", len(runner.repairs), len(runner.qaPrompts))
	}
	repair := runner.repairs[0]
	wantPrompt := fmt.Sprintf("QA Report Results rows lie outside the Results table.\nReport: %s\n%s; %s.\nKeep every row in the seeded Results table, directly below its last row, with no blank line before it and the same three columns; never change a verdict or a row status.\n", brk.Path, brk.Detail(), brk.Repair())
	if got := qaResultsRepairPrompt(brk.Path, brk); got != wantPrompt {
		t.Fatalf("qaResultsRepairPrompt = %q, want %q", got, wantPrompt)
	}
	if repair.Prompt != wantPrompt {
		t.Fatalf("repair prompt = %q, want %q", repair.Prompt, wantPrompt)
	}
	initial := runner.requests[len(runner.requests)-1]
	if !reflect.DeepEqual(initial.Session, repair.Session) || initial.RunID != repair.RunID || !reflect.DeepEqual(initial.Batch, repair.Batch) || initial.LogPath != repair.LogPath {
		t.Fatal("repair did not reuse the QA session and Batch")
	}
	count := 0
	for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonQA) {
		payload := eventPayloadMap(t, event)
		if payload["phase"] != "results_table" {
			continue
		}
		count++
		want := map[string]any{"phase": "results_table", "outcome": "repair", "report": brk.Path, "line": float64(brk.Line)}
		if !reflect.DeepEqual(payload, want) {
			t.Errorf("Results table event payload = %v, want %v", payload, want)
		}
		message := fmt.Sprintf("QA Report Results row outside the Results table at line %d; returning the QA turn once.", brk.Line)
		if event.Summary != message || !strings.Contains(fixture.progress.String(), message+"\n") {
			t.Errorf("event/progress lacks %q: %+v", message, event)
		}
	}
	if count != 1 {
		t.Fatalf("Results table repair events = %d, want 1", count)
	}
}
