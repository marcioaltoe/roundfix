// Boundary: the real queue store and engine, with a fake repository gate.
package delivery

import (
	"testing"

	"roundfix/internal/store"
)

func TestGateParksWithTheBlockerItNames(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, blocker, want string
		class               ParkClass
	}{
		{"container", BlockerLinuxLaneContainer + ": missing tool git; see lane.log", BlockerLinuxLaneContainer + ": missing tool git; see lane.log", ParkClassEnvironment},
		{"change", BlockerGateFailed + ": Linux Lane change failure: go test exited 1; see lane.log", BlockerGateFailed + ": Linux Lane change failure: go test exited 1; see lane.log", ParkClassFinding},
		{"repository", "", BlockerGateFailed, ParkClassFinding},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			runStore := openDeliveryEngineStore(t, ctx)
			const gitRoot = "/repo"
			if _, err := runStore.CreateDeliveryQueue(ctx, gitRoot, []string{"spec-one"}); err != nil {
				t.Fatal(err)
			}
			workflow := newFakeDeliveryWorkflow()
			workflow.gates["spec-one"] = GateResult{Reason: "gate reason", Blocker: tc.blocker}
			engine := newTestDeliveryEngine(runStore, workflow, newFakeDeliveryBoundary())
			if _, err := engine.Run(ctx, gitRoot); err != nil {
				t.Fatal(err)
			}
			queue := readDeliveryQueue(t, ctx, runStore, gitRoot)
			item := queue.Items[0]
			if item.Stage != store.DeliveryStageParked || item.Blocker != tc.want {
				t.Fatalf("item=%+v want blocker=%q", item, tc.want)
			}
			if park := ClassifyPark(queue, item); park.Class != tc.class {
				t.Fatalf("classification=%+v want class=%s", park, tc.class)
			}
		})
	}
}
