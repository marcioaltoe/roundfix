// Suite: prior QA pass import
// Invariant: only recorded Run settlements supply unchanged QA files to a later pass.
// Boundary IN: real Git, Run Database and journal, import, prompt and TaskCycle
// Boundary OUT: ACP runtime (task-cycle fixture), carry rules (speccheck/qa_row_carry_test.go)
package daemon

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
	"roundfix/internal/store"
)

func priorFixture(t *testing.T) (*taskCycleFixture, TaskPlan, *Engine) {
	t.Helper()
	f := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	p := f.qaPlan()
	commitTaskFixtureSource(t, f.gitRoot, "test: establish QA task")
	p.HeadSHA = strings.TrimSpace(runGitForTest(t, f.gitRoot, "rev-parse", "HEAD"))
	return f, p, f.engine(t, &taskFakeRunner{}, &taskFakeVerifier{}, GitCommitter{}, GitWorktreeSnapshotter{})
}

func priorReportName(sequence int) string {
	return fmt.Sprintf("qa-report-%s-%02d.md", taskCycleNowForTest().Format("2006-01-02"), sequence)
}

func priorReportContent(status string) string {
	return "---\nverdict: fail\nrows_blocked_environment: 0\nrows_blocked_finding: 0\nrows_blocked_declared: 0\n---\n\n## Results\n\n| # | Status | Provenance | Evidence |\n| --- | --- | --- | --- |\n| 1 | " + status + " | criterion | observed |\n"
}

func recordPriorCommit(t *testing.T, f *taskCycleFixture, p TaskPlan, name, report, message string, recorded bool) priorQAPass {
	t.Helper()
	base := strings.TrimSpace(runGitForTest(t, f.gitRoot, "rev-parse", "HEAD"))
	run := f.run
	branch := store.RunBranchPrefix + run.ID
	refs := strings.TrimSpace(runGitForTest(t, f.gitRoot, "for-each-ref", "--format=%(refname)", "refs/heads/"+branch))
	if refs == "" {
		runGitForTest(t, f.gitRoot, "checkout", "-b", branch)
	} else {
		runGitForTest(t, f.gitRoot, "checkout", branch)
	}
	audited := strings.TrimSpace(runGitForTest(t, f.gitRoot, "rev-parse", "HEAD"))
	rel, err := filepath.Rel(f.gitRoot, p.Spec.Dir)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.ToSlash(filepath.Join(rel, "qa", name))
	evidence := filepath.ToSlash(filepath.Join(rel, "qa", "evidence", name+".bin"))
	if err := os.MkdirAll(filepath.Dir(filepath.Join(f.gitRoot, evidence)), 0o755); err != nil {
		t.Fatal(err)
	}
	mustWriteForTest(t, filepath.Join(f.gitRoot, path), report)
	mustWriteForTest(t, filepath.Join(f.gitRoot, evidence), "\x00evidence\r\n\xff")
	runGitForTest(t, f.gitRoot, "add", path, evidence)
	runGitForTest(t, f.gitRoot, "commit", "-m", message)
	commit := strings.TrimSpace(runGitForTest(t, f.gitRoot, "rev-parse", "HEAD"))
	if recorded {
		payload, err := json.Marshal(map[string]any{"decision": "created", "task": "task_02", "commit": commit, "report": path})
		if err != nil {
			t.Fatal(err)
		}
		_, err = f.store.AppendRunEvent(context.Background(), runevent.RunEvent{RunID: run.ID, Batch: 2, Source: runevent.SourceDaemon, Kind: runevent.KindDaemonCommit, Time: taskCycleNowForTest(), Payload: payload})
		if err != nil {
			t.Fatal(err)
		}
	}
	runGitForTest(t, f.gitRoot, "checkout", "--detach", base)
	return priorQAPass{Commit: commit, Head: audited, Report: path, Files: []string{evidence, path}}
}

func assertPriorOutcome(t *testing.T, f *taskCycleFixture, outcome, reason string) {
	t.Helper()
	events := taskEventsOfKind(f.sink, runevent.KindDaemonQA)
	if len(events) == 0 {
		t.Fatal("missing prior_report event")
	}
	payload := eventPayloadMap(t, events[len(events)-1])
	if payload["phase"] != "prior_report" || payload["outcome"] != outcome || payload["reason"] != reason {
		t.Fatalf("prior event = %+v", payload)
	}
}

func TestAPriorPassIsImportedOnlyFromARecordedRunCommit(t *testing.T) {
	f, p, e := priorFixture(t)
	recordPriorCommit(t, f, p, priorReportName(1), priorReportContent("pass"), QACommitMessage(p.Spec.Slug, "fail"), false)
	runGitForTest(t, f.gitRoot, "branch", "-m", store.RunBranchPrefix+f.run.ID, "test/fabricated-qa")
	_, imported, err := e.importPriorQAPass(context.Background(), p, 1)
	if err != nil || imported {
		t.Fatalf("fabricated import = %v, %v", imported, err)
	}
	assertPriorOutcome(t, f, "none", "")
	if _, err := os.Stat(filepath.Join(p.Spec.Dir, "qa")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("fabricated files created: %v", err)
	}
}

func TestPriorQAPassImportsTheNewestUnintegratedReportByteForByte(t *testing.T) {
	f, p, e := priorFixture(t)
	recordPriorCommit(t, f, p, priorReportName(1), priorReportContent("fail"), QACommitMessage(p.Spec.Slug, "fail"), true)
	// Commit dates, rather than report names or Run creation order, select the pass.
	t.Setenv("GIT_COMMITTER_DATE", "2030-01-01T00:00:00Z")
	newest := recordPriorCommit(t, f, p, priorReportName(2), priorReportContent("pass"), QACommitMessage(p.Spec.Slug, "fail"), true)
	pass, imported, err := e.importPriorQAPass(context.Background(), p, 1)
	if err != nil || !imported || pass.Commit != newest.Commit {
		t.Fatalf("import = %+v %v %v", pass, imported, err)
	}
	for _, path := range pass.Files {
		got, err := os.ReadFile(filepath.Join(f.gitRoot, path))
		if err != nil {
			t.Fatal(err)
		}
		want, err := priorQAGit(context.Background(), f.gitRoot, "cat-file", "blob", pass.Commit+":"+path)
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != string(want) {
			t.Fatalf("bytes differ for %s", path)
		}
		info, err := os.Stat(filepath.Join(f.gitRoot, path))
		if err != nil || info.Mode().Perm() != 0o644 {
			t.Fatalf("mode = %v %v", info, err)
		}
	}
	assertPriorOutcome(t, f, "imported", "")
}

func TestPriorQAPassIgnoresReportsAlreadyInHistory(t *testing.T) {
	f, p, e := priorFixture(t)
	pass := recordPriorCommit(t, f, p, priorReportName(1), priorReportContent("pass"), QACommitMessage(p.Spec.Slug, "pass"), true)
	runGitForTest(t, f.gitRoot, "checkout", "--detach", pass.Commit)
	_, imported, err := e.importPriorQAPass(context.Background(), p, 1)
	if err != nil || imported {
		t.Fatalf("ancestor import = %v %v", imported, err)
	}
	assertPriorOutcome(t, f, "none", "")
}

func TestPriorQAPassIgnoresTaskCommitsAndOtherSpecs(t *testing.T) {
	for _, message := range []string{"feat: task\n\nRoundfix-Spec: " + taskCycleSlug + "\nRoundfix-Task: task_01", QACommitMessage("another-spec", "fail")} {
		t.Run(strings.Split(message, "\n")[0], func(t *testing.T) {
			f, p, e := priorFixture(t)
			recordPriorCommit(t, f, p, priorReportName(1), priorReportContent("pass"), message, true)
			_, imported, err := e.importPriorQAPass(context.Background(), p, 1)
			if err != nil || imported {
				t.Fatalf("unrelated import = %v %v", imported, err)
			}
			assertPriorOutcome(t, f, "none", "")
		})
	}
}

func TestPriorQAPassRefusesADifferingPath(t *testing.T) {
	f, p, e := priorFixture(t)
	pass := recordPriorCommit(t, f, p, priorReportName(1), priorReportContent("pass"), QACommitMessage(p.Spec.Slug, "fail"), true)
	if err := os.MkdirAll(filepath.Dir(filepath.Join(f.gitRoot, pass.Files[0])), 0o755); err != nil {
		t.Fatal(err)
	}
	mustWriteForTest(t, filepath.Join(f.gitRoot, pass.Files[0]), "local bytes")
	_, imported, err := e.importPriorQAPass(context.Background(), p, 1)
	if err != nil || imported {
		t.Fatalf("conflict import = %v %v", imported, err)
	}
	assertPriorOutcome(t, f, "refused", "path differs: "+pass.Files[0])
	got, err := os.ReadFile(filepath.Join(f.gitRoot, pass.Files[0]))
	if err != nil || string(got) != "local bytes" {
		t.Fatalf("local bytes = %q %v", got, err)
	}
	if _, err := os.Stat(filepath.Join(f.gitRoot, pass.Report)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("report created: %v", err)
	}
}

func TestPriorQAPassRefusesAReportOlderThanTheTreesNewest(t *testing.T) {
	f, p, e := priorFixture(t)
	pass := recordPriorCommit(t, f, p, priorReportName(1), priorReportContent("pass"), QACommitMessage(p.Spec.Slug, "fail"), true)
	newest := filepath.Join(p.Spec.Dir, "qa", priorReportName(2))
	if err := os.MkdirAll(filepath.Dir(newest), 0o755); err != nil {
		t.Fatal(err)
	}
	mustWriteForTest(t, newest, priorReportContent("pass"))
	_, imported, err := e.importPriorQAPass(context.Background(), p, 1)
	if err != nil || imported {
		t.Fatalf("old import = %v %v", imported, err)
	}
	assertPriorOutcome(t, f, "refused", "not newer than "+taskPromptPath(p, newest))
	if _, err := os.Stat(filepath.Join(f.gitRoot, pass.Report)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("old report created: %v", err)
	}
}

func TestPriorQAPassRefusesAReportTheShapeDetectorRefuses(t *testing.T) {
	f, p, e := priorFixture(t)
	pass := recordPriorCommit(t, f, p, priorReportName(1), priorReportContent("pending"), QACommitMessage(p.Spec.Slug, "fail"), true)
	_, imported, err := e.importPriorQAPass(context.Background(), p, 1)
	if err != nil || imported {
		t.Fatalf("shape import = %v %v", imported, err)
	}
	assertPriorOutcome(t, f, "refused", "report shape: "+speccheck.CodeMechanicalReportShape)
	for _, path := range pass.Files {
		if _, err := os.Stat(filepath.Join(f.gitRoot, path)); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("refused file remains: %s %v", path, err)
		}
	}
	if _, err := os.Stat(filepath.Join(p.Spec.Dir, "qa")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("refused QA directory remains: %v", err)
	}
}

func TestPriorQAPassRefusesAReportDatedAfterToday(t *testing.T) {
	f, p, e := priorFixture(t)
	pass := recordPriorCommit(t, f, p, "qa-report-2099-01-01.md", priorReportContent("pass"), QACommitMessage(p.Spec.Slug, "fail"), true)
	_, imported, err := e.importPriorQAPass(context.Background(), p, 1)
	if err != nil || imported {
		t.Fatalf("future import = %v %v", imported, err)
	}
	assertPriorOutcome(t, f, "refused", "dated after today")
	if _, err := os.Stat(filepath.Join(f.gitRoot, pass.Report)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("future report created: %v", err)
	}
}

func TestQAPromptNamesTheImportedPassHead(t *testing.T) {
	f, p, e := priorFixture(t)
	prior := recordPriorCommit(t, f, p, priorReportName(1), priorReportContent("pass"), QACommitMessage(p.Spec.Slug, "fail"), true)
	runGitForTest(t, f.gitRoot, "commit", "--allow-empty", "-m", "test: later task")
	p.HeadSHA = strings.TrimSpace(runGitForTest(t, f.gitRoot, "rev-parse", "HEAD"))
	pass, _, err := e.importPriorQAPass(context.Background(), p, 1)
	if err != nil {
		t.Fatal(err)
	}
	prompt, err := e.buildQAPromptContext(context.Background(), p, p.Tasks[len(p.Tasks)-1], pass)
	if err != nil || prompt.PreviousReportPath != prior.Report || prompt.PreviousReportHead != prior.Head {
		t.Fatalf("prompt = %+v %v", prompt, err)
	}
}

func TestQAGateCarriesARowFromAnUnintegratedFailedPass(t *testing.T) {
	f, p, _ := priorFixture(t)
	base := strings.TrimSpace(runGitForTest(t, f.gitRoot, "rev-parse", "HEAD"))
	// The failed pass audits a Task that the next head independently re-commits.
	runGitForTest(t, f.gitRoot, "checkout", "-b", "test/first-task")
	const input = "unchanged input\n"
	mustWriteForTest(t, filepath.Join(f.gitRoot, "carry-input.txt"), input)
	runGitForTest(t, f.gitRoot, "add", "carry-input.txt")
	runGitForTest(t, f.gitRoot, "commit", "-m", "feat: first task\n\nRoundfix-Spec: "+p.Spec.Slug+"\nRoundfix-Task: task_01")
	audited := strings.TrimSpace(runGitForTest(t, f.gitRoot, "rev-parse", "HEAD"))
	report := fmt.Sprintf("---\nverdict: fail\nrows_blocked_environment: 0\nrows_blocked_finding: 0\nrows_blocked_declared: 0\nevidence_snapshots:\n  \"1\":\n    head: %s\n    inputs:\n      - ref: carry-input.txt\n        files:\n          - path: carry-input.txt\n            sha256: %x\n---\n\n## Results\n\n| # | Status | Provenance | Evidence |\n| --- | --- | --- | --- |\n| 1 | pass | criterion | observed |\n| 2 | fail | criterion | observed |\n\n### 1 evidence\n\n```yaml\ninputs:\n  - kind: repository_path\n    ref: carry-input.txt\n```\n", audited, sha256.Sum256([]byte(input)))
	prior := recordPriorCommit(t, f, p, priorReportName(1), report, QACommitMessage(p.Spec.Slug, "fail"), true)
	if _, err := f.store.CompleteRun(context.Background(), f.run.ID, store.StateFailed); err != nil {
		t.Fatal(err)
	}
	next, err := f.store.CreateRun(context.Background(), store.CreateRunRequest{Kind: store.KindImplement, GitRoot: f.gitRoot, LocalBranch: f.run.LocalBranch, SpecSlug: p.Spec.Slug})
	if err != nil {
		t.Fatal(err)
	}
	f.run = next
	p.RunID = next.ID
	runGitForTest(t, f.gitRoot, "checkout", "--detach", base)
	mustWriteForTest(t, filepath.Join(f.gitRoot, "carry-input.txt"), input)
	runGitForTest(t, f.gitRoot, "add", "carry-input.txt")
	runGitForTest(t, f.gitRoot, "commit", "-m", "feat: recommitted task\n\nRoundfix-Spec: "+p.Spec.Slug+"\nRoundfix-Task: task_01")
	p = f.plan()
	ancestor, err := priorQAAncestor(context.Background(), f.gitRoot, audited, "HEAD")
	if err != nil || ancestor {
		t.Fatalf("audited head must be outside next history: %v %v", ancestor, err)
	}
	runner := &taskFakeRunner{calls: f.calls, gitRoot: f.gitRoot, qaReport: qaReportForTest(spec.VerdictPass)}
	engine := f.engine(t, runner, &taskFakeVerifier{calls: f.calls}, GitCommitter{}, GitWorktreeSnapshotter{})
	result, err := engine.TaskCycle(context.Background(), p)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(runner.qaSeed, "| 1 | carried (established by: "+prior.Report+"; head: "+audited+") | criterion |") || !strings.Contains(runner.qaSeed, "| 2 | re-run: "+speccheck.CarryReasonNotPass+" |") {
		t.Fatalf("second seed lacks carry dispositions:\n%s", runner.qaSeed)
	}
	if len(runner.qaPrompts) != 1 || !strings.Contains(runner.qaPrompts[0], prior.Report) || !strings.Contains(runner.qaPrompts[0], audited) {
		t.Fatalf("prompt lacks previous pass: %v", runner.qaPrompts)
	}
	for _, path := range prior.Files {
		got, err := priorQAGit(context.Background(), f.gitRoot, "cat-file", "blob", "HEAD:"+path)
		if err != nil {
			t.Fatal(err)
		}
		want, err := priorQAGit(context.Background(), f.gitRoot, "cat-file", "blob", prior.Commit+":"+path)
		if err != nil || string(got) != string(want) {
			t.Fatalf("imported commit bytes for %s differ: %v", path, err)
		}
	}
	if result.QAReportPath == prior.Report {
		t.Fatal("new pass overwrote imported report")
	}
	// The settlement journal payload has the SHA needed by another Run.
	commits := taskEventsOfKind(f.sink, runevent.KindDaemonCommit)
	payload := eventPayloadMap(t, commits[len(commits)-1])
	head := strings.TrimSpace(runGitForTest(t, f.gitRoot, "rev-parse", "HEAD"))
	if payload["commit"] != head || payload["task"] == "" {
		t.Fatalf("settlement lacks recorded commit: %+v", payload)
	}
}

func TestPriorQAPassKeepsIdenticalExistingEvidence(t *testing.T) {
	f, p, e := priorFixture(t)
	pass := recordPriorCommit(t, f, p, priorReportName(1), priorReportContent("pass"), QACommitMessage(p.Spec.Slug, "fail"), true)
	path := filepath.Join(f.gitRoot, pass.Files[0])
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("\x00evidence\r\n\xff"), 0o600); err != nil {
		t.Fatal(err)
	}
	before, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	_, imported, err := e.importPriorQAPass(context.Background(), p, 1)
	if err != nil || !imported {
		t.Fatalf("identical import = %v %v", imported, err)
	}
	after, err := os.Stat(path)
	if err != nil || !os.SameFile(before, after) || !before.ModTime().Equal(after.ModTime()) || after.Mode().Perm() != 0o600 {
		t.Fatalf("existing evidence replaced: %v %v", after, err)
	}
}

func TestPriorQAPassSkipsAPrunedRunBranch(t *testing.T) {
	f, p, e := priorFixture(t)
	pass := recordPriorCommit(t, f, p, priorReportName(1), priorReportContent("pass"), QACommitMessage(p.Spec.Slug, "fail"), true)
	branch := store.RunBranchPrefix + f.run.ID
	// A derived Task Branch is still present, but cannot stand in for the Run Branch.
	runGitForTest(t, f.gitRoot, "branch", "-D", branch)
	runGitForTest(t, f.gitRoot, "branch", branch+"/task_01", pass.Commit)
	_, imported, err := e.importPriorQAPass(context.Background(), p, 1)
	if err != nil || imported {
		t.Fatalf("pruned import = %v %v", imported, err)
	}
	assertPriorOutcome(t, f, "none", "")
}

func TestPriorQAPassCancellationPublishesStop(t *testing.T) {
	f, p, e := priorFixture(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, imported, err := e.importPriorQAPass(ctx, p, 1)
	if imported || !errors.Is(err, context.Canceled) {
		t.Fatalf("cancel import = %v %v", imported, err)
	}
	if len(taskEventsOfKind(f.sink, runevent.KindDaemonStatus)) != 1 {
		t.Fatal("missing stop event")
	}
}

func TestPriorQAPassGitErrorIsInfrastructure(t *testing.T) {
	f, p, e := priorFixture(t)
	recordPriorCommit(t, f, p, priorReportName(1), priorReportContent("pass"), QACommitMessage(p.Spec.Slug, "fail"), true)
	if err := os.Rename(filepath.Join(f.gitRoot, ".git"), filepath.Join(f.gitRoot, "unavailable-git")); err != nil {
		t.Fatal(err)
	}
	_, imported, err := e.importPriorQAPass(context.Background(), p, 1)
	if imported || err == nil {
		t.Fatalf("Git error was not returned: %v %v", imported, err)
	}
	if len(taskEventsOfKind(f.sink, runevent.KindDaemonQA)) != 0 {
		t.Fatal("Git error reported as a refusal or none")
	}
}
