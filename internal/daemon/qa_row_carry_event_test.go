// Suite: QA mechanical carry event
// Invariant: the mechanical phase reports the carried and re-run work counts.
// Boundary IN: task-cycle fixture and daemon.qa Run event payload
// Boundary OUT: Git carry proof (internal/speccheck/qa_row_carry_test.go)
package daemon

import (
	"context"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
	"testing"
)

func TestQAMechanicalEventCountsCarriedAndRerunRows(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	engine := fixture.engine(t, &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot}, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)
	engine.deps.MechanicalStage = &fakeQAMechanicalStage{result: speccheck.MechanicalResult{
		Blocking:     true,
		Carried:      []speccheck.CarriedRow{{ID: "R01", EstablishedBy: "report.md", EstablishedHead: "abc"}},
		Dispositions: []speccheck.CarryDisposition{{ID: "R01", Carried: true}, {ID: "R02", Reason: speccheck.CarryReasonNotPass}, {ID: "R03", Reason: speccheck.CarryReasonNoInputs}},
	}}
	if _, err := engine.TaskCycle(context.Background(), fixture.qaPlan()); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonQA) {
		payload := eventPayloadMap(t, event)
		if payload["phase"] != "mechanical" {
			continue
		}
		found = true
		if payload["carried_rows"] != float64(1) || payload["rerun_rows"] != float64(2) {
			t.Fatalf("mechanical counts: %+v", payload)
		}
	}
	if !found {
		t.Fatal("missing mechanical event")
	}
}
