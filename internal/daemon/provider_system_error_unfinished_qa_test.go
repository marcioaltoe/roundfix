// Suite: Provider System Error recovery after an unfinished QA verdict.
// Invariant: the fallback inherits a pending report and the discarded verdict is recorded.
// Boundary OUT: runtime processes and Verification are package fakes.
package daemon

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
)

type verdictBeforeFailureRunner struct {
	*providerSystemErrorRunner
	verdict          string
	unfinishedReport string
}

func (runner *verdictBeforeFailureRunner) Run(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	if len(runner.seen) == 0 {
		path, _, err := seededQAReportPathFromPromptForTest(req.Prompt, runner.gitRoot)
		if err != nil {
			return agent.ExecuteResult{}, err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return agent.ExecuteResult{}, err
		}
		runner.unfinishedReport = strings.Replace(string(content), "verdict: pending", "verdict: "+runner.verdict, 1)
		if err := os.WriteFile(path, []byte(runner.unfinishedReport), 0o644); err != nil {
			return agent.ExecuteResult{}, err
		}
	}
	return runner.providerSystemErrorRunner.Run(ctx, req, sink)
}

func TestQAProviderSystemErrorAfterAVerdictFallsBackAndReopensIt(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		verdict string
	}{
		{"partial", spec.VerdictPartial},
		{"pass", spec.VerdictPass},
		{"fail", spec.VerdictFail},
		{"a pending report records no discarded verdict", spec.VerdictPending},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
			runner := &verdictBeforeFailureRunner{
				providerSystemErrorRunner: &providerSystemErrorRunner{
					taskFakeRunner: &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot, qaReport: qaReportForTest(spec.VerdictPass)},
					turns:          []error{providerSystemErrorForTest()},
				},
				verdict: tc.verdict,
			}
			engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)
			engine.deps.MechanicalStage = &fakeQAMechanicalStage{}
			plan := fixture.qaPlan()
			plan.AgentSelections = selectionProfilesForTest(map[roundconfig.WorkCategory]roundconfig.AgentSelectionProfile{
				roundconfig.CategoryQA: selectionProfileForTest(selectionForTest("codex", "preferred-model", "high"), selectionForTest("claude", "fallback-model", "high")),
			})
			plan.RuntimeFactory = runtimeFactoryForLifecycleTest(nil)
			result, err := engine.TaskCycle(t.Context(), plan)
			if err != nil || !result.QAAccepted || result.QAVerdict != spec.VerdictPass {
				t.Fatalf("result=%+v err=%v", result, err)
			}
			if len(runner.seen) != 2 || runner.seen[1].Runtime.Model != "fallback-model" || !strings.HasSuffix(runner.seen[1].Session.Name, "-fallback-01") || len(runner.closed) != 2 {
				t.Fatalf("requests=%v closed=%v", runner.seen, runner.closed)
			}
			if runner.seen[1].Prompt != runner.seen[0].Prompt+"\n\n"+providerSystemErrorNoticeForTest {
				t.Fatalf("fallback prompt=%s", runner.seen[1].Prompt)
			}
			wantSection := "\n## Agent runtime fallback\n\n- provider_system_error: codex/preferred-model/high: " + providerSystemErrorDetailForTest + "\n"
			if tc.verdict != spec.VerdictPending {
				wantSection += "- discarded_verdict: " + tc.verdict + "\n"
			}
			wantSection += "- fallback: claude/fallback-model/high\n- retry_spent: false\n"
			wantPending := strings.Replace(runner.unfinishedReport, "verdict: "+tc.verdict, "verdict: pending", 1)
			if runner.qaFallback != wantPending+wantSection {
				t.Fatalf("fallback report=%s\nwant=%s", runner.qaFallback, wantPending+wantSection)
			}
			path, _, err := seededQAReportPathFromPromptForTest(runner.seen[1].Prompt, runner.gitRoot)
			if err != nil {
				t.Fatal(err)
			}
			finalReport, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if strings.Count(string(finalReport), wantSection) != 1 {
				t.Fatalf("final report lost or repeated fallback section: %s", finalReport)
			}
			wantDiscarded := 1
			if tc.verdict == spec.VerdictPending {
				wantDiscarded = 0
			}
			for _, report := range []string{runner.qaFallback, string(finalReport)} {
				if strings.Count(report, "- discarded_verdict:") != wantDiscarded {
					t.Fatalf("discarded verdict count: %s", report)
				}
			}
			assertProviderSystemErrorAttempts(t, fixture, "qa", "qa", 1)
			assertProviderSystemErrorEvents(t, fixture.sink, "qa", "fallback")
			singleEventOfKind(t, fixture.sink, runevent.KindDaemonAgentSelectionFallback)
		})
	}
}

func TestQAProviderSystemErrorAfterAVerdictWithoutFallbackSettlesRuntimeInfrastructure(t *testing.T) {
	t.Parallel()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	runner := &verdictBeforeFailureRunner{
		providerSystemErrorRunner: &providerSystemErrorRunner{
			taskFakeRunner: &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot},
			turns:          []error{providerSystemErrorForTest()},
		},
		verdict: spec.VerdictPartial,
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)
	engine.deps.MechanicalStage = &fakeQAMechanicalStage{}
	plan := fixture.qaPlan()
	plan.AgentSelections = selectionProfilesForTest(map[roundconfig.WorkCategory]roundconfig.AgentSelectionProfile{
		roundconfig.CategoryQA: selectionProfileForTest(selectionForTest("codex", "preferred-model", "high")),
	})
	plan.RuntimeFactory = runtimeFactoryForLifecycleTest(nil)
	result, err := engine.TaskCycle(t.Context(), plan)
	if err != nil || result.QAAccepted || result.QAVerdict != spec.VerdictPending || len(runner.seen) != 1 || len(runner.closed) != 1 {
		t.Fatalf("result=%+v err=%v requests=%v closed=%v", result, err, runner.seen, runner.closed)
	}
	assertProviderSystemErrorAttempts(t, fixture, "qa", "qa", 1)
	assertProviderSystemErrorEvents(t, fixture.sink, "qa", "exhausted")
	wantReason := "runtime infrastructure: qa qa hit a provider system error: " + providerSystemErrorDetailForTest
	settled := 0
	for _, event := range eventsOfKind(fixture.sink, runevent.KindDaemonTask) {
		var payload map[string]any
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			t.Fatal(err)
		}
		if payload["phase"] == "settled" && payload["status"] == "failed" && payload["reason"] == wantReason {
			settled++
		}
	}
	if settled != 1 || len(eventsOfKind(fixture.sink, runevent.KindDaemonAgentSelectionFallback)) != 0 {
		t.Fatalf("failed QA settlements=%d", settled)
	}
	path, _, err := seededQAReportPathFromPromptForTest(runner.seen[0].Prompt, runner.gitRoot)
	if err != nil {
		t.Fatal(err)
	}
	report, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(report) != runner.unfinishedReport {
		t.Fatalf("exhaustion changed the unfinished report: %s", report)
	}
}

func TestResetUnfinishedQAVerdictPreservesTheRestOfTheReport(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name      string
		content   string
		want      string
		discarded string
	}{
		{
			name:      "only the frontmatter line changes",
			content:   "---\nstatus: closed\nverdict: partial\nrows_blocked_environment: 2\n---\n\n# QA Report\nverdict: fail\ntrailing spaces  ",
			want:      "---\nstatus: closed\nverdict: pending\nrows_blocked_environment: 2\n---\n\n# QA Report\nverdict: fail\ntrailing spaces  ",
			discarded: spec.VerdictPartial,
		},
		{
			name:      "CRLF and quoted verdict",
			content:   "---\r\nverdict: 'pass' # unfinished\r\nsummary: kept\r\n---\r\nbody\r\n",
			want:      "---\r\nverdict: pending\r\nsummary: kept\r\n---\r\nbody\r\n",
			discarded: spec.VerdictPass,
		},
		{name: "pending keeps its bytes", content: "---\nverdict: pending # keep\n---\nbody"},
		{name: "missing verdict keeps its bytes", content: "---\nsummary: no verdict\n---\nverdict: pass\n"},
		{name: "no frontmatter keeps its bytes", content: "# QA Report\nverdict: pass\n"},
		{name: "unclosed frontmatter keeps its bytes", content: "---\nverdict: fail\nbody"},
		{name: "blank verdict keeps its bytes", content: "---\nverdict: \n---\nbody"},
		{name: "unreadable frontmatter keeps its bytes", content: "---\nverdict: pass\nsummary: [\n---\nbody"},
		{name: "duplicate verdict keeps its bytes", content: "---\nverdict: pass\nverdict: fail\n---\nbody"},
		{name: "block scalar keeps its bytes", content: "---\nverdict: |-\n  pass\n---\nbody"},
		{name: "empty report keeps its bytes"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			path := filepath.Join(t.TempDir(), "report.md")
			if err := os.WriteFile(path, []byte(tc.content), 0o600); err != nil {
				t.Fatal(err)
			}
			discarded, err := resetUnfinishedQAVerdict(path)
			if err != nil || discarded != tc.discarded {
				t.Fatalf("discarded=%q want=%q err=%v", discarded, tc.discarded, err)
			}
			got, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			want := tc.want
			if tc.discarded == "" {
				want = tc.content
			}
			if string(got) != want {
				t.Fatalf("report=%q want=%q", got, want)
			}
			info, err := os.Stat(path)
			if err != nil || info.Mode().Perm() != 0o600 {
				t.Fatalf("report permissions changed: info=%v err=%v", info, err)
			}
		})
	}
	t.Run("read failure is returned", func(t *testing.T) {
		t.Parallel()
		discarded, err := resetUnfinishedQAVerdict(filepath.Join(t.TempDir(), "missing.md"))
		if !errors.Is(err, os.ErrNotExist) || discarded != "" {
			t.Fatalf("discarded=%q err=%v", discarded, err)
		}
	})
}
