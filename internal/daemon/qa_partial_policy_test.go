// Suite: Daemon partial policy at the TaskCycle settlement boundary.
// Boundary IN: reports from taskFakeRunner and temporary Spec fixtures.
// Boundary OUT: command streams and external Agent execution.
package daemon

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"roundfix/internal/runevent"
	"roundfix/internal/spec"
)

func TestTheDaemonSettlesEveryPartialShapeAsTheCommandsDo(t *testing.T) {
	t.Parallel()
	for _, shape := range qaPartialShapes() {
		t.Run(shape.name, func(t *testing.T) {
			fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01"}})
			plan := fixture.qaPlan()
			prd := filepath.Join(plan.Spec.Dir, "_prd.md")
			content, err := os.ReadFile(prd)
			if err != nil {
				t.Fatal(err)
			}
			mustWriteForTest(t, prd, string(content)+shape.unreachable())
			fixture.worktree.snapshots = [][]string{nil, {"src/one.go"}, {"src/one.go"}, {"src/one.go", qaReportRelPathForTest()}}
			runner := &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot, store: fixture.store, statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted}, qaReport: shape.report()}
			engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)
			result, err := engine.TaskCycle(t.Context(), plan)
			if err != nil {
				t.Fatal(err)
			}
			accepted := shape.reason == ""
			wantStatus := spec.StatusFailed
			wantReason := "QA verdict partial not accepted: " + shape.reason
			if accepted {
				wantStatus = spec.StatusCompleted
				wantReason = ""
			}
			if result.QAVerdict != spec.VerdictPartial || result.QAAccepted != accepted || taskStatusOnDisk(t, fixture.gitRoot, fixture.graph.QATaskID) != string(wantStatus) {
				t.Fatalf("result=%+v want QA accepted=%v status=%s", result, accepted, wantStatus)
			}
			var settlement map[string]any
			for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonTask) {
				payload := eventPayloadMap(t, event)
				if event.ReviewIssue == fixture.graph.QATaskID && payload["phase"] == "settled" {
					settlement = payload
					break
				}
			}
			gotReason, _ := settlement["reason"].(string)
			if settlement == nil || gotReason != wantReason {
				t.Fatalf("settlement=%v, want reason %q", settlement, wantReason)
			}
		})
	}
}

// The same measured shapes exercise command and Daemon settlement independently.
type qaPartialShape struct {
	name, rows                                   string
	environment, finding, declared, declarations int
	reason                                       string
}

func qaPartialShapes() []qaPartialShape {
	pr := "| " + spec.QANoOpenPullRequestStatus + " | " + spec.QAPullRequestRowSource + " |\n"
	network := "| " + spec.QANetworkDeniedStatusPrefix + "docs.example.com) | Requirement 7; " + spec.QAOutsideEvidenceRowSource + " (published guide) |\n"
	declared := "| blocked (declared) | Requirement 1 |\n| blocked (declared) | Requirement 2 |\n"
	return []qaPartialShape{
		{"Spec 0213 PR only", pr, 1, 0, 0, 0, ""},
		{"Spec 0220 PR only", pr + pr, 2, 0, 0, 0, ""},
		{"Oraculum PR only", pr, 1, 0, 0, 0, ""},
		{"Spec 0227 annotated PR", "| " + spec.QANoOpenPullRequestStatus + " | Pull Request row (Requirement 10 constrains execution) |\n", 1, 0, 0, 0, ""},
		{"Fluxus covered declarations and PR", pr + declared, 1, 0, 2, 2, ""},
		{"network only", network, 1, 0, 0, 0, ""},
		{"PR and network", pr + network, 2, 0, 0, 0, ""},
		{"all three kinds", pr + network + declared, 2, 0, 2, 2, ""},
		{"declared only", declared, 0, 0, 2, 2, ""},
		{"network without outside evidence", "| " + spec.QANetworkDeniedStatusPrefix + "docs.example.com) | Requirement 8 |\n", 1, 0, 0, 0, "rows_blocked_environment is 1; expected 0"},
		{"empty network host", "| " + spec.QANetworkDeniedStatusPrefix + ") | outside-evidence row |\n", 1, 0, 0, 0, "rows_blocked_environment is 1; expected 0"},
		{"another environment beside network", network + "| blocked (environment: credentials unavailable) | Requirement 8 |\n", 2, 0, 0, 0, "rows_blocked_environment is 2, 1 outside the pre-PR Pull Request row and network-denied outside-evidence rows; expected 0 outside them"},
		{"PR and skipped", pr + "| skipped | Requirement 9 |\n", 1, 0, 0, 0, "partial records 1 skipped row(s); a qualifying partial records none"},
		{"finding", pr + "| blocked (finding) | Requirement 10 |\n", 1, 1, 0, 0, "rows_blocked_finding is 1; expected 0"},
		{"uncovered declaration", pr + declared, 1, 0, 2, 1, "rows_blocked_declared is 2, but Spec declares 1 unreachable acceptance; shortfall is 1"},
		{"no unmet row", "| pass | Requirement 1 |\n", 0, 0, 0, 0, `newest QA Report verdict is "partial"; expected "pass"`},
	}
}

func (shape qaPartialShape) report() string {
	return fmt.Sprintf("---\nverdict: partial\nrows_blocked_environment: %d\nrows_blocked_finding: %d\nrows_blocked_declared: %d\n---\n\n## Results\n\n| Status | Provenance |\n| --- | --- |\n%s", shape.environment, shape.finding, shape.declared, shape.rows)
}

func (shape qaPartialShape) unreachable() string {
	if shape.declarations == 0 {
		return ""
	}
	text := "\n## Unreachable Acceptance\n\n"
	for i := 0; i < shape.declarations; i++ {
		text += fmt.Sprintf("- criterion: acceptance criterion %d\n  reason: unavailable environment\n  satisfied-by: task_01\n", i+1)
	}
	return text
}
