// Suite: QA Report evidence binding
// Invariant: one repair regenerates this Run's evidence or refuses the report without a commit.
// Boundary IN: TaskCycle, real Git, prior-pass import, QA prompts and Run Events
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
	"roundfix/internal/store"
)

type qaBorrowedEvidenceRunner struct {
	*taskFakeRunner
	repairs  []agent.ExecuteRequest
	onRepair func()
}

func (runner *qaBorrowedEvidenceRunner) Run(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	if !strings.HasPrefix(req.Prompt, "QA evidence does not belong to this Run.\n") {
		return runner.taskFakeRunner.Run(ctx, req, sink)
	}
	runner.repairs = append(runner.repairs, req)
	if runner.onRepair != nil {
		runner.onRepair()
	}
	return agent.ExecuteResult{LogPath: req.LogPath}, nil
}

func borrowedQAReportForTest(link string) (string, int) {
	row := "| 1 | pass | criterion | [execution](" + link + ") |\n"
	report := qaReportForTest(spec.VerdictPass) + "\n## Results\n\n| # | Status | Provenance | Evidence |\n| --- | --- | --- | --- |\n" + row
	return report, strings.Count(strings.SplitN(report, row, 2)[0], "\n") + 1
}

func borrowedQAFixture(t *testing.T, repair bool) (*taskCycleFixture, TaskCycleResult, *qaBorrowedEvidenceRunner, string, int) {
	t.Helper()
	f, p, _ := priorFixture(t)
	const link = "evidence/2026-11-06-gate/head.txt"
	oldPath := filepath.Join(p.Spec.Dir, "qa", filepath.FromSlash(link))
	if err := os.MkdirAll(filepath.Dir(oldPath), 0o755); err != nil {
		t.Fatal(err)
	}
	mustWriteForTest(t, oldPath, "earlier Run execution\n")
	commitTaskFixtureSource(t, f.gitRoot, "test: earlier Run evidence")
	head := strings.TrimSpace(runGitForTest(t, f.gitRoot, "rev-parse", "HEAD"))
	report, line := borrowedQAReportForTest(link)
	runner := &qaBorrowedEvidenceRunner{taskFakeRunner: &taskFakeRunner{
		calls: f.calls, gitRoot: f.gitRoot, qaReport: report,
	}}
	if repair {
		runner.onRepair = func() {
			reportPath := filepath.Join(f.gitRoot, runner.qaReportPath)
			newLink := "evidence/" + strings.TrimSuffix(filepath.Base(reportPath), ".md") + "/head.txt"
			newPath := filepath.Join(filepath.Dir(reportPath), filepath.FromSlash(newLink))
			if err := os.MkdirAll(filepath.Dir(newPath), 0o755); err != nil {
				t.Fatal(err)
			}
			mustWriteForTest(t, newPath, "this Run execution at "+head+"\n")
			content, err := os.ReadFile(reportPath)
			if err != nil {
				t.Fatal(err)
			}
			mustWriteForTest(t, reportPath, strings.ReplaceAll(string(content), link, newLink))
		}
	}
	engine := f.engine(t, runner, &taskFakeVerifier{calls: f.calls}, GitCommitter{}, GitWorktreeSnapshotter{})
	result, err := engine.TaskCycle(t.Context(), f.plan())
	if err != nil {
		t.Fatalf("TaskCycle: %v", err)
	}
	return f, result, runner, head, line
}

func TestQABorrowedEvidenceRepairTurnRegeneratesTheEvidence(t *testing.T) {
	t.Parallel()
	f, result, runner, head, line := borrowedQAFixture(t, true)
	if !result.QAAccepted || taskStatusOnDisk(t, f.gitRoot, f.graph.QATaskID) != string(spec.StatusCompleted) {
		t.Fatalf("regenerated evidence was not accepted: %+v", result)
	}
	assertBorrowedQARepair(t, f, runner, head, line, "evidence/2026-11-06-gate/head.txt", spec.BorrowedEvidenceReused)
	assertBorrowedQAEvents(t, f, line, "evidence/2026-11-06-gate/head.txt", spec.BorrowedEvidenceReused, []string{"repair"})
	if message := strings.TrimSpace(runGitForTest(t, f.gitRoot, "log", "-1", "--format=%B")); message != QACommitMessage(taskCycleSlug, spec.VerdictPass) {
		t.Fatalf("QA commit message = %q", message)
	}
	newLink := "evidence/" + strings.TrimSuffix(filepath.Base(result.QAReportPath), ".md") + "/head.txt"
	newPath := filepath.ToSlash(filepath.Join(filepath.Dir(result.QAReportPath), newLink))
	if got := runGitForTest(t, f.gitRoot, "show", "HEAD:"+newPath); got != "this Run execution at "+head+"\n" {
		t.Fatalf("committed evidence = %q", got)
	}
	if got := runGitForTest(t, f.gitRoot, "show", "HEAD:"+result.QAReportPath); !strings.Contains(got, "[execution]("+newLink+")") {
		t.Fatal("committed report lost regenerated evidence link")
	}
}

func TestQABorrowedEvidenceAfterRepairSettlesFailed(t *testing.T) {
	t.Parallel()
	f, result, runner, head, line := borrowedQAFixture(t, false)
	const link = "evidence/2026-11-06-gate/head.txt"
	assertBorrowedQARepair(t, f, runner, head, line, link, spec.BorrowedEvidenceReused)
	assertBorrowedQAEvents(t, f, line, link, spec.BorrowedEvidenceReused, []string{"repair", "refused"})
	assertBorrowedQARefused(t, f, result, line, link, spec.BorrowedEvidenceReused)
	evidence, err := os.ReadFile(filepath.Join(f.gitRoot, filepath.Dir(result.QAReportPath), filepath.FromSlash(link)))
	if err != nil || string(evidence) != "earlier Run execution\n" {
		t.Fatalf("refusal lost evidence: %q, %v", evidence, err)
	}
}

func TestQABorrowedEvidenceNamesARerunRowCitingImportedEvidence(t *testing.T) {
	t.Parallel()
	f, p, _ := priorFixture(t)
	prior := recordPriorCommit(t, f, p, priorReportName(1), priorReportContent("fail"), QACommitMessage(p.Spec.Slug, "fail"), true)
	if _, err := f.store.CompleteRun(t.Context(), f.run.ID, store.StateFailed); err != nil {
		t.Fatal(err)
	}
	next, err := f.store.CreateRun(t.Context(), store.CreateRunRequest{Kind: store.KindImplement, GitRoot: f.gitRoot, LocalBranch: f.run.LocalBranch, SpecSlug: p.Spec.Slug})
	if err != nil {
		t.Fatal(err)
	}
	f.run = next
	link := "evidence/" + priorReportName(1) + ".bin"
	report, line := borrowedQAReportForTest(link)
	runner := &qaBorrowedEvidenceRunner{taskFakeRunner: &taskFakeRunner{calls: f.calls, gitRoot: f.gitRoot, qaReport: report}}
	engine := f.engine(t, runner, &taskFakeVerifier{calls: f.calls}, GitCommitter{}, GitWorktreeSnapshotter{})
	result, err := engine.TaskCycle(t.Context(), f.plan())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(runner.qaSeed, "| 1 | re-run:") {
		t.Fatalf("imported failed row was not marked for re-run:\n%s", runner.qaSeed)
	}
	assertBorrowedQARepair(t, f, runner, prior.Head, line, link, spec.BorrowedEvidenceImported)
	assertBorrowedQAEvents(t, f, line, link, spec.BorrowedEvidenceImported, []string{"repair", "refused"})
	assertBorrowedQARefused(t, f, result, line, link, spec.BorrowedEvidenceImported)
	for _, path := range prior.Files {
		got, err := os.ReadFile(filepath.Join(f.gitRoot, filepath.FromSlash(path)))
		if err != nil {
			t.Fatal(err)
		}
		want, err := priorQAGit(t.Context(), f.gitRoot, "cat-file", "blob", prior.Commit+":"+path)
		if err != nil || string(got) != string(want) {
			t.Fatalf("import changed %s: %v", path, err)
		}
	}
}

func TestQAPromptNamesTheRunWorktreeAuditedHeadAndEvidenceDirectory(t *testing.T) {
	t.Parallel()
	f, _, runner, head, _ := borrowedQAFixture(t, true)
	if len(runner.qaPrompts) != 1 {
		t.Fatalf("initial QA prompts = %d, want 1", len(runner.qaPrompts))
	}
	evidenceDir := filepath.ToSlash(filepath.Join(filepath.Dir(runner.qaReportPath), "evidence", strings.TrimSuffix(filepath.Base(runner.qaReportPath), ".md")))
	want := fmt.Sprintf("Seeded QA Report: %s\nComplete this report in place, preserving its materialized mechanical rows and skips; do not create another QA Report.\nRun Worktree: %s\nAudited head: %s\nEvidence directory: %s/\nBuild and run every executable a row exercises inside this Run Worktree, list their paths in user_flow_executables, write this pass's evidence under the evidence directory, and never copy, replay or cite evidence or an executable an earlier Run left.\n", runner.qaReportPath, f.gitRoot, head, evidenceDir)
	if !strings.Contains(runner.qaPrompts[0], want) {
		t.Fatalf("QA prompt lacks ordered Run facts:\n%s", runner.qaPrompts[0])
	}
}

func TestQAEvidenceLinksAndBorrowedEvidenceShareOneRepairTurn(t *testing.T) {
	t.Parallel()
	f, p, _ := priorFixture(t)
	const borrowedLink = "evidence/earlier.txt"
	const missingLink = "evidence/missing.txt"
	oldPath := filepath.Join(p.Spec.Dir, "qa", filepath.FromSlash(borrowedLink))
	if err := os.MkdirAll(filepath.Dir(oldPath), 0o755); err != nil {
		t.Fatal(err)
	}
	mustWriteForTest(t, oldPath, "earlier execution\n")
	commitTaskFixtureSource(t, f.gitRoot, "test: earlier evidence for shared repair")
	report, line := borrowedQAReportForTest(borrowedLink)
	report += "| 2 | pass | criterion | [missing](" + missingLink + ") |\n"
	runner := &qaEvidenceLinkRunner{taskFakeRunner: &taskFakeRunner{calls: f.calls, gitRoot: f.gitRoot, qaReport: report}}
	runner.onRepair = func() {
		root := filepath.Join(f.gitRoot, filepath.Dir(runner.qaReportPath))
		mustWriteForTest(t, filepath.Join(root, filepath.FromSlash(borrowedLink)), "this Run execution\n")
		mustWriteForTest(t, filepath.Join(root, filepath.FromSlash(missingLink)), "this Run missing execution\n")
	}
	engine := f.engine(t, runner, &taskFakeVerifier{calls: f.calls}, GitCommitter{}, GitWorktreeSnapshotter{})
	result, err := engine.TaskCycle(t.Context(), f.plan())
	if err != nil {
		t.Fatal(err)
	}
	if !result.QAAccepted || len(runner.repairs) != 1 || len(runner.qaPrompts) != 1 {
		t.Fatalf("shared repair: result=%+v, repair turns=%d, QA prompts=%d", result, len(runner.repairs), len(runner.qaPrompts))
	}
	prompt := runner.repairs[0].Prompt
	wantLink := fmt.Sprintf("- %s (line %d, missing)\n", missingLink, line+1)
	wantBorrowed := fmt.Sprintf("QA evidence does not belong to this Run.\nReport: %s", filepath.Join(f.gitRoot, runner.qaReportPath))
	if !strings.HasPrefix(prompt, "QA evidence links do not resolve in the committed tree.\n") || !strings.Contains(prompt, wantLink) || !strings.Contains(prompt, wantBorrowed) || strings.Index(prompt, wantLink) > strings.Index(prompt, wantBorrowed) {
		t.Fatalf("combined prompt lacks ordered link and binding repairs:\n%s", prompt)
	}
	assertBorrowedQAEvents(t, f, line, borrowedLink, spec.BorrowedEvidenceReused, []string{"repair"})
	for _, link := range []string{borrowedLink, missingLink} {
		path := filepath.ToSlash(filepath.Join(filepath.Dir(result.QAReportPath), link))
		if got := runGitForTest(t, f.gitRoot, "show", "HEAD:"+path); !strings.HasPrefix(got, "this Run ") {
			t.Fatalf("shared repair evidence was not committed: %s = %q", path, got)
		}
	}
}

func assertBorrowedQARepair(t *testing.T, f *taskCycleFixture, runner *qaBorrowedEvidenceRunner, head string, line int, subject, reason string) {
	t.Helper()
	if len(runner.repairs) != 1 {
		t.Fatalf("repair turns = %d, want 1", len(runner.repairs))
	}
	repair := runner.repairs[0]
	evidenceDir := filepath.ToSlash(filepath.Join(filepath.Dir(runner.qaReportPath), "evidence", strings.TrimSuffix(filepath.Base(runner.qaReportPath), ".md")))
	for _, want := range []string{"Report: ", runner.qaReportPath, "Run Worktree: " + f.gitRoot, "Audited head: " + head, "Evidence directory: " + evidenceDir + "/", fmt.Sprintf("- %s (line %d, %s)", subject, line, reason), "never change a verdict or a row the new evidence does not support."} {
		if !strings.Contains(repair.Prompt, want) {
			t.Errorf("repair prompt lacks %q:\n%s", want, repair.Prompt)
		}
	}
	initial := runner.requests[len(runner.requests)-1]
	if !reflect.DeepEqual(initial.Session, repair.Session) || initial.RunID != repair.RunID || !reflect.DeepEqual(initial.Batch, repair.Batch) || initial.LogPath != repair.LogPath {
		t.Fatal("repair did not reuse the QA session and Batch")
	}
}

func assertBorrowedQAEvents(t *testing.T, f *taskCycleFixture, line int, subject, reason string, outcomes []string) {
	t.Helper()
	var got []string
	for _, event := range taskEventsOfKind(f.sink, runevent.KindDaemonQA) {
		payload := eventPayloadMap(t, event)
		if payload["phase"] != "evidence_binding" {
			continue
		}
		outcome, _ := payload["outcome"].(string)
		got = append(got, outcome)
		want := map[string]any{"phase": "evidence_binding", "outcome": outcome, "items": []any{map[string]any{"subject": subject, "line": float64(line), "reason": reason}}}
		if !reflect.DeepEqual(payload, want) {
			t.Errorf("binding event = %v, want %v", payload, want)
		}
		message := "QA evidence does not belong to this Run: 1 item(s); returning the QA turn once."
		if outcome == "refused" {
			message = "QA Report refused: 1 borrowed evidence item(s)"
		}
		if event.Summary != message || !strings.Contains(f.progress.String(), message+"\n") {
			t.Errorf("event/progress lacks %q: %+v", message, event)
		}
	}
	if !reflect.DeepEqual(got, outcomes) {
		t.Fatalf("binding outcomes = %v, want %v", got, outcomes)
	}
}

func assertBorrowedQARefused(t *testing.T, f *taskCycleFixture, result TaskCycleResult, line int, subject, reason string) {
	t.Helper()
	if result.QAAccepted || taskStatusOnDisk(t, f.gitRoot, f.graph.QATaskID) != string(spec.StatusFailed) {
		t.Fatalf("borrowed evidence was accepted: %+v", result)
	}
	wantReason := fmt.Sprintf("QA Report evidence does not belong to this Run: %s (line %d, %s)", subject, line, reason)
	var gotReason string
	for _, event := range taskEventsOfKind(f.sink, runevent.KindDaemonTask) {
		payload := eventPayloadMap(t, event)
		if event.ReviewIssue == f.graph.QATaskID && payload["phase"] == "settled" {
			gotReason, _ = payload["reason"].(string)
		}
	}
	if gotReason != wantReason {
		t.Fatalf("settled reason = %q, want %q", gotReason, wantReason)
	}
	if history := runGitForTest(t, f.gitRoot, "log", "--format=%s"); strings.Contains(history, "docs: qa report for "+taskCycleSlug) {
		t.Fatal("refused QA Report was committed")
	}
	for _, event := range taskEventsOfKind(f.sink, runevent.KindDaemonCommit) {
		if strings.HasPrefix(event.Summary, "QA Report commit created:") {
			t.Fatal("refusal published a QA Report commit event")
		}
	}
	report, err := os.ReadFile(filepath.Join(f.gitRoot, result.QAReportPath))
	if err != nil || !strings.Contains(string(report), "[execution]("+subject+")") {
		t.Fatalf("refusal lost report or evidence link: %v", err)
	}
}
