// Boundary: command output and temporary Run Database; retry owner uses the existing fake.
package cli

import (
	"bytes"
	"roundfix/internal/gittest"
	"roundfix/internal/store"
	"strings"
	"testing"
)

func TestStatusPrintsWhatTheLaneRecorded(t *testing.T) {
	t.Parallel()
	home, repo := withCLIWorkspace(t)
	gittest.InitRepo(t, repo, "-b", "main")
	db, err := store.Open(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	q, err := db.CreateDeliveryQueue(t.Context(), repo, []string{"parked-spec", "merged-spec"})
	if err != nil {
		t.Fatal(err)
	}
	parked, merged := q.Items[0], q.Items[1]
	parked.Stage = store.DeliveryStageParked
	parked.Blocker = "linux-lane-container: Docker unavailable"
	merged.Stage = store.DeliveryStageMerged
	merged.Warning = "linux-lane-skipped: image missing; with detail"
	for _, item := range []store.DeliveryQueueItem{parked, merged} {
		if err := db.UpdateDeliveryQueueItem(t.Context(), repo, item); err != nil {
			t.Fatal(err)
		}
	}
	var out, diagnostic bytes.Buffer
	code := runCLI(t, []string{"deliver", "status"}, &out, &diagnostic)
	if code != 0 || diagnostic.Len() != 0 || !strings.Contains(out.String(), merged.Warning) || !strings.Contains(out.String(), "Park: parked-spec environment: run roundfix deliver retry parked-spec; the Linux Lane was retired, so the gate no longer runs it") {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, out.String(), diagnostic.String())
	}
}

func TestADeliveryRetryReturnsAnItemTheLaneParked(t *testing.T) {
	t.Parallel()
	f := newRetryOwnerProofFixture(t, false, "live")
	f.before.Blocker = "linux-lane-container: Docker unavailable"
	db, err := store.Open(t.Context(), f.home)
	if err != nil {
		t.Fatal(err)
	}
	queue, err := db.CreateDeliveryQueue(t.Context(), f.repo, []string{implementTestSlug, "merged-spec"})
	if err != nil {
		t.Fatal(err)
	}
	merged := queue.Items[1]
	merged.Stage = store.DeliveryStageMerged
	merged.Warning = "linux-lane-skipped: image missing; with detail"
	if err := db.UpdateDeliveryQueueItem(t.Context(), f.repo, merged); err != nil {
		t.Fatal(err)
	}
	if err := db.UpdateDeliveryQueueItem(t.Context(), f.repo, f.before); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	code, out, diagnostic := f.retry(t)
	queue = openDeliveryQueueForCLI(t, f.home, f.repo)
	item := queue.Items[0]
	if queue.Items[1].Stage != store.DeliveryStageMerged || queue.Items[1].Warning != merged.Warning {
		t.Fatalf("stored merged item changed: %+v", queue.Items[1])
	}
	if code != 0 || diagnostic != "" || item.Stage == store.DeliveryStageParked || item.Blocker != "" || !strings.Contains(out, "Retried "+implementTestSlug+": linux-lane-container: Docker unavailable -> running\n") {
		t.Fatalf("exit=%d stdout=%q stderr=%q item=%+v", code, out, diagnostic, item)
	}
}
