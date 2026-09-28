package daemon

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/runevent"
	"roundfix/internal/spec"
)

// Suite: QA front matter settlement
// Invariant: the Daemon fails an unreadable seeded QA Report and completes a filled, readable seed.
// Boundary IN: mechanical seeding, QA Agent file edits, shared report parsing, settlement, and commit verdict.
// Boundary OUT: CLI refusal presentation, owned by internal/cli tests.

func TestQASettlementRefusesAnEmptyFrontMatter(t *testing.T) {
	result, fixture, runner, reason := runSeededQAFrontmatterSettlementTest(t, func(seed string) string {
		return strings.Replace(seed, "---\n", "---\n---\n", 1)
	})

	if !strings.HasPrefix(runner.qaSeed, "---\nverdict: pending\n") {
		t.Fatalf("QA Agent did not receive the seeded front matter:\n%s", runner.qaSeed)
	}
	if result.QAVerdict != qaVerdictUnreadable || result.QAAccepted {
		t.Fatalf("QA result = verdict %q accepted %t, want unreadable and refused", result.QAVerdict, result.QAAccepted)
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, fixture.graph.QATaskID); got != string(spec.StatusFailed) {
		t.Fatalf("QA Task status = %q, want %q", got, spec.StatusFailed)
	}
	if !strings.Contains(reason, "QA verdict unreadable:") || !strings.Contains(reason, "QA Report front matter is empty") {
		t.Fatalf("QA Task reason = %q, want unreadable verdict and empty-front-matter cause", reason)
	}
}

func TestQASettlementAcceptsAFilledSeededFrontMatter(t *testing.T) {
	result, fixture, _, reason := runSeededQAFrontmatterSettlementTest(t, func(seed string) string {
		seed = strings.Replace(seed, "verdict: pending", "verdict: pass", 1)
		return strings.Replace(seed, "\n## Mechanical skips", "| R01 | pass | observed CLI output |\n\n## Mechanical skips", 1)
	})

	assertQASettlement(t, result, fixture, spec.VerdictPass, spec.StatusCompleted, "", reason)
}

func runSeededQAFrontmatterSettlementTest(t *testing.T, mutate func(string) string) (TaskCycleResult, *taskCycleFixture, *taskFakeRunner, string) {
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
		preserveQASeed: true,
	}
	runner.afterQA = func() {
		path := runner.qaReportPath
		if !filepath.IsAbs(path) {
			path = filepath.Join(fixture.gitRoot, path)
		}
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read seeded QA Report: %v", err)
		}
		if err := os.WriteFile(path, []byte(mutate(string(content))), 0o644); err != nil {
			t.Fatalf("mutate seeded QA Report: %v", err)
		}
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
