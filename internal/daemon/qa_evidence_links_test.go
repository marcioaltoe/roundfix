// Suite: QA Report evidence link settlement
// Invariant: one repair turn keeps evidence or refuses its report without a commit.
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

type qaEvidenceLinkRunner struct {
	*taskFakeRunner
	repairs  []agent.ExecuteRequest
	onRepair func()
}

func (runner *qaEvidenceLinkRunner) Run(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	if !strings.HasPrefix(req.Prompt, "QA evidence links do not resolve in the committed tree.\n") {
		return runner.taskFakeRunner.Run(ctx, req, sink)
	}
	runner.repairs = append(runner.repairs, req)
	if runner.onRepair != nil {
		runner.onRepair()
	}
	return agent.ExecuteResult{LogPath: req.LogPath}, nil
}

func qaEvidenceLinkFixture(t *testing.T, extension string, repair bool) (*taskCycleFixture, TaskCycleResult, *qaEvidenceLinkRunner, int) {
	t.Helper()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	plan := fixture.qaPlan()
	mustWriteForTest(t, filepath.Join(fixture.gitRoot, ".gitignore"), "*.log\n")
	commitTaskFixtureSource(t, fixture.gitRoot, "test: establish ignored QA logs")
	link := "evidence/consumer-tests." + extension
	row := "| 1 | pass | [consumer tests](" + link + ") |\n"
	report := qaReportForTest(spec.VerdictPass) + "\n## Results\n\n| # | Status | Evidence |\n| --- | --- | --- |\n" + row
	line := strings.Count(strings.SplitN(report, row, 2)[0], "\n") + 1
	runner := &qaEvidenceLinkRunner{taskFakeRunner: &taskFakeRunner{
		calls: fixture.calls, gitRoot: fixture.gitRoot, qaReport: report,
	}}
	runner.afterQA = func() {
		path := filepath.Join(fixture.gitRoot, filepath.Dir(runner.qaReportPath), link)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		mustWriteForTest(t, path, "consumer tests passed\n")
	}
	if repair {
		runner.onRepair = func() {
			root := filepath.Join(fixture.gitRoot, filepath.Dir(runner.qaReportPath))
			if err := os.Rename(filepath.Join(root, link), filepath.Join(root, "evidence/consumer-tests.txt")); err != nil {
				t.Fatal(err)
			}
			content, err := os.ReadFile(filepath.Join(fixture.gitRoot, runner.qaReportPath))
			if err != nil {
				t.Fatal(err)
			}
			mustWriteForTest(t, filepath.Join(fixture.gitRoot, runner.qaReportPath), strings.ReplaceAll(string(content), link, "evidence/consumer-tests.txt"))
		}
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, GitCommitter{}, GitWorktreeSnapshotter{})
	result, err := engine.TaskCycle(t.Context(), plan)
	if err != nil {
		t.Fatalf("TaskCycle: %v", err)
	}
	return fixture, result, runner, line
}

func TestQAEvidenceLinkRepairTurnKeepsTheEvidence(t *testing.T) {
	t.Parallel()
	fixture, result, runner, line := qaEvidenceLinkFixture(t, "log", true)
	if !result.QAAccepted || taskStatusOnDisk(t, fixture.gitRoot, fixture.graph.QATaskID) != string(spec.StatusCompleted) {
		t.Fatalf("repaired QA was not accepted: %+v", result)
	}
	assertQAEvidenceRepair(t, runner, line)
	assertQAEvidenceLinkEvents(t, fixture, line, []string{"repair"})
	assertQAEvidenceCommitted(t, fixture, result)
	old := filepath.Join(fixture.gitRoot, filepath.Dir(result.QAReportPath), "evidence/consumer-tests.log")
	if _, err := os.Stat(old); !os.IsNotExist(err) {
		t.Fatalf("repair left ignored log: %v", err)
	}
}

func TestQAEvidenceLinkUnresolvedAfterRepairSettlesFailed(t *testing.T) {
	t.Parallel()
	fixture, result, runner, line := qaEvidenceLinkFixture(t, "log", false)
	if result.QAAccepted || taskStatusOnDisk(t, fixture.gitRoot, fixture.graph.QATaskID) != string(spec.StatusFailed) {
		t.Fatalf("unresolved QA must settle failed: %+v; committed files: %s", result, runGitForTest(t, fixture.gitRoot, "ls-tree", "-r", "--name-only", "HEAD"))
	}
	assertQAEvidenceRepair(t, runner, line)
	assertQAEvidenceLinkEvents(t, fixture, line, []string{"repair", "refused"})
	wantReason := fmt.Sprintf("QA Report evidence links do not resolve in the committed tree: evidence/consumer-tests.log (line %d, ignored)", line)
	var reason string
	for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonTask) {
		payload := eventPayloadMap(t, event)
		if event.ReviewIssue == fixture.graph.QATaskID && payload["phase"] == "settled" {
			reason, _ = payload["reason"].(string)
		}
	}
	if reason != wantReason {
		t.Fatalf("reason = %q, want %q", reason, wantReason)
	}
	if history := runGitForTest(t, fixture.gitRoot, "log", "--format=%s"); strings.Contains(history, "docs: qa report for "+taskCycleSlug) {
		t.Fatal("refused report was committed")
	}
	for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonCommit) {
		if strings.HasPrefix(event.Summary, "QA Report commit created:") {
			t.Fatal("refusal published a commit event")
		}
	}
	reportPath := filepath.Join(fixture.gitRoot, result.QAReportPath)
	report, err := os.ReadFile(reportPath)
	if err != nil || !strings.Contains(string(report), "[consumer tests](evidence/consumer-tests.log)") {
		t.Fatalf("refusal lost report or link: %v", err)
	}
	evidence, err := os.ReadFile(filepath.Join(filepath.Dir(reportPath), "evidence/consumer-tests.log"))
	if err != nil || string(evidence) != "consumer tests passed\n" {
		t.Fatalf("refusal lost evidence: %q, %v", evidence, err)
	}
}

func TestQAReportWithResolvedLinksRunsNoRepairTurn(t *testing.T) {
	t.Parallel()
	fixture, result, runner, line := qaEvidenceLinkFixture(t, "txt", false)
	if !result.QAAccepted || taskStatusOnDisk(t, fixture.gitRoot, fixture.graph.QATaskID) != string(spec.StatusCompleted) {
		t.Fatalf("resolved QA was not accepted: %+v", result)
	}
	if len(runner.repairs) != 0 {
		t.Fatalf("resolved QA ran %d repair turns", len(runner.repairs))
	}
	assertQAEvidenceLinkEvents(t, fixture, line, nil)
	assertQAEvidenceCommitted(t, fixture, result)
}

func assertQAEvidenceRepair(t *testing.T, runner *qaEvidenceLinkRunner, line int) {
	t.Helper()
	if len(runner.repairs) != 1 {
		t.Fatalf("repair turns = %d, want 1", len(runner.repairs))
	}
	repair := runner.repairs[0]
	for _, want := range []string{runner.qaReportPath, fmt.Sprintf("evidence/consumer-tests.log (line %d, ignored)", line), ".txt", "restore", "never change a verdict or a row"} {
		if !strings.Contains(repair.Prompt, want) {
			t.Errorf("repair prompt lacks %q: %s", want, repair.Prompt)
		}
	}
	initial := runner.requests[len(runner.requests)-1]
	if !reflect.DeepEqual(initial.Session, repair.Session) || initial.RunID != repair.RunID || !reflect.DeepEqual(initial.Batch, repair.Batch) || initial.LogPath != repair.LogPath {
		t.Fatal("repair did not reuse the QA session and batch")
	}
}

func assertQAEvidenceLinkEvents(t *testing.T, fixture *taskCycleFixture, line int, outcomes []string) {
	t.Helper()
	var got []string
	for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonQA) {
		payload := eventPayloadMap(t, event)
		if payload["phase"] != "evidence_links" {
			continue
		}
		outcome, _ := payload["outcome"].(string)
		got = append(got, outcome)
		want := map[string]any{"phase": "evidence_links", "outcome": outcome, "links": []any{map[string]any{"link": "evidence/consumer-tests.log", "line": float64(line), "reason": "ignored"}}}
		if !reflect.DeepEqual(payload, want) {
			t.Errorf("event payload = %v, want %v", payload, want)
		}
		message := "QA evidence links unresolved: 1; returning the QA turn once."
		if outcome == "refused" {
			message = "QA Report refused: 1 unresolved evidence link(s)"
		}
		if event.Summary != message || !strings.Contains(fixture.progress.String(), message+"\n") {
			t.Errorf("event/progress lacks %q: %+v", message, event)
		}
	}
	if !reflect.DeepEqual(got, outcomes) {
		t.Fatalf("evidence link outcomes = %v, want %v", got, outcomes)
	}
}

func assertQAEvidenceCommitted(t *testing.T, fixture *taskCycleFixture, result TaskCycleResult) {
	t.Helper()
	if message := strings.TrimSpace(runGitForTest(t, fixture.gitRoot, "log", "-1", "--format=%B")); message != QACommitMessage(taskCycleSlug, spec.VerdictPass) {
		t.Fatalf("QA commit message = %q", message)
	}
	path := filepath.ToSlash(filepath.Join(filepath.Dir(result.QAReportPath), "evidence/consumer-tests.txt"))
	if got := runGitForTest(t, fixture.gitRoot, "show", "HEAD:"+path); got != "consumer tests passed\n" {
		t.Fatalf("committed evidence = %q", got)
	}
	if got := runGitForTest(t, fixture.gitRoot, "show", "HEAD:"+result.QAReportPath); !strings.Contains(got, "[consumer tests](evidence/consumer-tests.txt)") {
		t.Fatal("committed report lost repaired link")
	}
}
