// Suite: two real QA gate passes
// Invariant: unmoved passing evidence survives into the next committed report.
// Boundary IN: TaskCycle, real Git, snapshot recorder, carry resolver and Run events
// Boundary OUT: ACP Agent and repository Verification (task-cycle fakes)
package daemon

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

	"roundfix/internal/agent"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
	"roundfix/internal/store"
)

const twoPassProvenance = "PRD criterion 1; unchanged repository behavior"

// The second Agent takes its work list from the seed, never manufacturing a
// carried row. Missing carry material therefore causes a fresh observation,
// which the committed-report assertions reject in the unmoved case.
type twoPassQAAgent struct {
	taskFakeRunner
	first bool
	rerun []string
}

func (runner *twoPassQAAgent) Run(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	if !strings.Contains(req.Prompt, "Spec QA gate") {
		return runner.taskFakeRunner.Run(ctx, req, sink)
	}
	path, _, err := seededQAReportPathFromPromptForTest(req.Prompt, runner.gitRoot)
	if err != nil {
		return agent.ExecuteResult{}, err
	}
	seed, err := os.ReadFile(path)
	if err != nil {
		return agent.ExecuteResult{}, err
	}
	if runner.first {
		runner.qaReport = "---\nverdict: fail\nrows_blocked_environment: 0\nrows_blocked_finding: 0\nrows_blocked_declared: 0\n---\n\n## Results\n\n| # | Status | Provenance |\n| --- | --- | --- |\n| 1 | pass | " + twoPassProvenance + " |\n| 2 | fail | PRD criterion 2 |\n" + twoPassRowInputs("1", "carry-input.txt") + twoPassRowInputs("2", "repair-input.txt")
	} else {
		report := strings.Replace(string(seed), "verdict: pending", "verdict: pass", 1)
		var rows, inputs strings.Builder
		for _, row := range []struct{ id, input string }{{"1", "carry-input.txt"}, {"2", "repair-input.txt"}} {
			if strings.Contains(report, "| "+row.id+" | carried (established by:") {
				continue
			}
			runner.rerun = append(runner.rerun, row.id)
			fmt.Fprintf(&rows, "| %s | pass | fresh observation of %s |\n", row.id, row.input)
			inputs.WriteString(twoPassRowInputs(row.id, row.input))
		}
		// Keep the seeded Results rows and carry dispositions byte-identical.
		runner.qaReport = strings.Replace(report, "\n\n## Mechanical skips", "\n"+rows.String()+"\n## Mechanical skips", 1) + inputs.String()
	}
	return runner.taskFakeRunner.Run(ctx, req, sink)
}

func twoPassRowInputs(id, input string) string {
	return fmt.Sprintf("\n### %s evidence\n\n```yaml\ninputs:\n  - kind: repository_path\n    ref: %s\n```\n", id, input)
}

func TestTwoGatePassesCarryAnUnmovedRowIntoTheCommittedReport(t *testing.T) {
	testTwoGatePasses(t, false)
}

func TestTwoGatePassesReRunARowWhoseInputMoved(t *testing.T) {
	testTwoGatePasses(t, true)
}

func testTwoGatePasses(t *testing.T, moved bool) {
	t.Helper()
	ctx := context.Background()
	f, _, _ := priorFixture(t)
	mustWriteForTest(t, filepath.Join(f.gitRoot, "carry-input.txt"), "stable input\n")
	mustWriteForTest(t, filepath.Join(f.gitRoot, "repair-input.txt"), "broken input\n")
	runGitForTest(t, f.gitRoot, "add", "carry-input.txt", "repair-input.txt")
	runGitForTest(t, f.gitRoot, "commit", "-m", "test: establish both row inputs")
	base := strings.TrimSpace(runGitForTest(t, f.gitRoot, "rev-parse", "HEAD"))
	runGitForTest(t, f.gitRoot, "checkout", "-b", store.RunBranchPrefix+f.run.ID)
	firstAgent := &twoPassQAAgent{taskFakeRunner: taskFakeRunner{calls: f.calls, gitRoot: f.gitRoot}, first: true}
	firstEngine := f.engine(t, firstAgent, &taskFakeVerifier{calls: f.calls}, GitCommitter{}, GitWorktreeSnapshotter{})
	firstEngine.deps.Sink = runevent.NewFanout([]runevent.Sink{f.store.JournalSink(), f.sink}, nil)
	first, err := firstEngine.TaskCycle(ctx, f.plan())
	if err != nil || first.QAVerdict != spec.VerdictFail {
		t.Fatalf("first gate = %+v, %v", first, err)
	}
	firstCommit := strings.TrimSpace(runGitForTest(t, f.gitRoot, "rev-parse", "HEAD"))
	firstReport := runGitForTest(t, f.gitRoot, "show", firstCommit+":"+first.QAReportPath)
	assertTwoPassSnapshot(t, firstReport, base, map[string]string{"1": "carry-input.txt"}, map[string]string{"carry-input.txt": "stable input\n"})
	assertTwoPassEvents(t, f, first, base, nil, 0, 0, 1)
	if _, err := f.store.CompleteRun(ctx, f.run.ID, store.StateFailed); err != nil {
		t.Fatal(err)
	}
	// The failed QA commit stays only on the first Run Branch. The correction
	// starts at the original audited tree, where the QA Task is still pending.
	runGitForTest(t, f.gitRoot, "checkout", "--detach", base)
	mustWriteForTest(t, filepath.Join(f.gitRoot, "repair-input.txt"), "fixed input\n")
	changed := []string{"repair-input.txt"}
	if moved {
		mustWriteForTest(t, filepath.Join(f.gitRoot, "carry-input.txt"), "moved input\n")
		changed = append(changed, "carry-input.txt")
	}
	runGitForTest(t, f.gitRoot, append([]string{"add"}, changed...)...)
	runGitForTest(t, f.gitRoot, "commit", "-m", "test: correct the failing row")
	secondHead := strings.TrimSpace(runGitForTest(t, f.gitRoot, "rev-parse", "HEAD"))
	if ancestor, err := priorQAAncestor(ctx, f.gitRoot, firstCommit, secondHead); err != nil || ancestor {
		t.Fatalf("first QA commit must be unintegrated: %v, %v", ancestor, err)
	}
	f.run, err = f.store.CreateRun(ctx, store.CreateRunRequest{Kind: store.KindImplement, GitRoot: f.gitRoot, LocalBranch: f.run.LocalBranch, SpecSlug: taskCycleSlug})
	if err != nil {
		t.Fatal(err)
	}
	runGitForTest(t, f.gitRoot, "checkout", "-b", store.RunBranchPrefix+f.run.ID)
	f.reloadGraph()
	f.sink = &captureEventSink{}
	secondAgent := &twoPassQAAgent{taskFakeRunner: taskFakeRunner{calls: f.calls, gitRoot: f.gitRoot}}
	secondEngine := f.engine(t, secondAgent, &taskFakeVerifier{calls: f.calls}, GitCommitter{}, GitWorktreeSnapshotter{})
	secondEngine.deps.Sink = runevent.NewFanout([]runevent.Sink{f.store.JournalSink(), f.sink}, nil)
	second, err := secondEngine.TaskCycle(ctx, f.plan())
	if err != nil || second.QAVerdict != spec.VerdictPass || second.QAReportPath == first.QAReportPath {
		t.Fatalf("second gate = %+v, %v\nprogress:\n%s\nseed:\n%s\nevents: %+v", second, err, f.progress.String(), secondAgent.qaSeed, taskEventsOfKind(f.sink, runevent.KindDaemonQA))
	}
	secondCommit := strings.TrimSpace(runGitForTest(t, f.gitRoot, "rev-parse", "HEAD"))
	committed := runGitForTest(t, f.gitRoot, "show", secondCommit+":"+second.QAReportPath)
	carried := "| 1 | carried (established by: " + first.QAReportPath + "; head: " + base + ") | " + twoPassProvenance + " |"
	wantRerun := []string{"2"}
	wantSnapshots := map[string]string{"2": "repair-input.txt"}
	contents := map[string]string{"repair-input.txt": "fixed input\n"}
	carryCount, rerunCount := 1, 1
	if moved {
		wantRerun = []string{"1", "2"}
		wantSnapshots["1"] = "carry-input.txt"
		contents["carry-input.txt"] = "moved input\n"
		carryCount, rerunCount = 0, 2
		if strings.Contains(committed, "| 1 | carried") || strings.Contains(committed, twoPassProvenance) {
			t.Fatalf("moved row kept old evidence:\n%s", committed)
		}
		if !strings.Contains(committed, "| 1 | pass | fresh observation of carry-input.txt |") || !strings.Contains(committed, "| 1 | "+speccheck.CarryDispositionRerun+speccheck.CarryReasonInputMoved+"carry-input.txt |") {
			t.Fatalf("moved row lacks fresh result or disposition:\n%s", committed)
		}
	} else if strings.Count(committed, carried) != 1 || !strings.Contains(committed, "## Row carry-forward\n\n| Prior row | Disposition |\n| --- | --- |\n| 1 | "+speccheck.CarryDispositionCarried+" |") {
		t.Fatalf("committed report lost carried row or provenance:\n%s", committed)
	}
	if !reflect.DeepEqual(secondAgent.rerun, wantRerun) {
		t.Fatalf("Agent re-ran %v, want %v", secondAgent.rerun, wantRerun)
	}
	if !strings.Contains(committed, "| 2 | pass | fresh observation of repair-input.txt |") || !strings.Contains(committed, "| 2 | "+speccheck.CarryDispositionRerun+speccheck.CarryReasonNotPass+" |") {
		t.Fatalf("corrected row lacks result or disposition:\n%s", committed)
	}
	assertTwoPassSnapshot(t, committed, secondHead, wantSnapshots, contents)
	prior := &priorQAPass{Commit: firstCommit, Report: first.QAReportPath, Files: []string{first.QAReportPath}}
	assertTwoPassEvents(t, f, second, secondHead, prior, carryCount, rerunCount, len(wantSnapshots))
	if imported := runGitForTest(t, f.gitRoot, "show", secondCommit+":"+first.QAReportPath); imported != firstReport {
		t.Fatal("first report was not imported byte-identically")
	}
}

func assertTwoPassEvents(t *testing.T, f *taskCycleFixture, result TaskCycleResult, head string, prior *priorQAPass, carried, rerun, snapshots int) {
	t.Helper()
	events := taskEventsOfKind(f.sink, runevent.KindDaemonQA)
	if len(events) != 4 {
		t.Fatalf("QA events = %+v, want prior, mechanical, verdict, snapshots", events)
	}
	for i, phase := range []string{"prior_report", "mechanical", "verdict", "evidence_snapshots"} {
		payload := eventPayloadMap(t, events[i])
		if events[i].RunID != f.run.ID || events[i].Batch != 1 || payload["phase"] != phase {
			t.Fatalf("QA event %d = %+v", i, events[i])
		}
	}
	previous := eventPayloadMap(t, events[0])
	if prior == nil {
		if previous["outcome"] != "none" || previous["commit"] != "" || previous["report"] != "" || previous["reason"] != "" || previous["files"] != nil {
			t.Fatalf("first prior_report = %+v", previous)
		}
	} else if previous["outcome"] != "imported" || previous["commit"] != prior.Commit || previous["report"] != prior.Report || previous["reason"] != "" || !reflect.DeepEqual(previous["files"], []any{prior.Report}) {
		t.Fatalf("second prior_report = %+v", previous)
	}
	mechanical := eventPayloadMap(t, events[1])
	if mechanical["report"] != result.QAReportPath || mechanical["blocking"] != false || mechanical["findings"] != float64(0) || mechanical["blocked_rows"] != float64(0) || mechanical["carried_rows"] != float64(carried) || mechanical["rerun_rows"] != float64(rerun) {
		t.Fatalf("mechanical = %+v", mechanical)
	}
	recorded := eventPayloadMap(t, events[3])
	if recorded["outcome"] != "recorded" || recorded["head"] != head || recorded["report"] != result.QAReportPath || recorded["rows"] != float64(snapshots) || recorded["error"] != "" {
		t.Fatalf("evidence_snapshots = %+v", recorded)
	}
}

func assertTwoPassSnapshot(t *testing.T, report, head string, rows, contents map[string]string) {
	t.Helper()
	var frontmatter struct {
		Snapshots map[string]struct {
			Head   string `yaml:"head"`
			Inputs []struct {
				Ref    string `yaml:"ref"`
				Count  int    `yaml:"count"`
				SHA256 string `yaml:"sha256"`
			} `yaml:"inputs"`
		} `yaml:"evidence_snapshots"`
	}
	parts := strings.SplitN(report, "---", 3)
	if len(parts) != 3 {
		t.Fatal("committed report lacks frontmatter")
	}
	if err := yaml.Unmarshal([]byte(parts[1]), &frontmatter); err != nil {
		t.Fatal(err)
	}
	if len(frontmatter.Snapshots) != len(rows) {
		t.Fatalf("snapshot rows = %+v, want %v", frontmatter.Snapshots, rows)
	}
	for id, path := range rows {
		snapshot, ok := frontmatter.Snapshots[id]
		if !ok || snapshot.Head != head || len(snapshot.Inputs) != 1 || snapshot.Inputs[0].Ref != path || snapshot.Inputs[0].Count != 1 {
			t.Fatalf("snapshot %s = %+v", id, snapshot)
		}
		input := snapshot.Inputs[0]
		if input.SHA256 != fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%x  %s\n", sha256.Sum256([]byte(contents[path])), path)))) {
			t.Fatalf("snapshot file %s = %+v", id, input)
		}
	}
}
