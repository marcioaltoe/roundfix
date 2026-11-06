// Boundary: retry transactions in a disposable SQLite queue store.
package store

import (
	"errors"
	"testing"
)

func TestTheStoreDoesNotCountADiskFullRetry(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, blocker string
		atLimit       bool
		want          int
		refused       bool
	}{
		{"disk full", "disk-full: /tmp/x/_pkg_.a: 0 MiB free", false, 0, false},
		{"bare disk full", "disk-full", false, 0, false},
		{"disk full at limit", "disk-full: /tmp/x/_pkg_.a: 0 MiB free", true, 1, false},
		{"other failure", "delivery-error: failed", false, 1, false},
		{"other failure at limit", "delivery-error: failed", true, 1, true},
		{"similar prefix is counted", "disk-full-extra", false, 1, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			s := openTestStore(t, ctx, t.TempDir())
			defer closeStore(t, s)
			const root = "/repo-disk-retry"
			queue, err := s.CreateDeliveryQueueWithLimits(ctx, root, []string{"widget"}, DeliveryQueueLimits{MaxRetries: 1})
			if err != nil {
				t.Fatal(err)
			}
			item := queue.Items[0]
			if tc.atLimit {
				item.Stage = DeliveryStageParked
				item.Blocker = "delivery-error"
				if err := s.UpdateDeliveryQueueItem(ctx, root, item); err != nil {
					t.Fatal(err)
				}
				target := item
				target.Stage = DeliveryStageRunning
				if _, _, err := s.RetryDeliveryQueueItem(ctx, root, target, item.Blocker); err != nil {
					t.Fatal(err)
				}
				item = readDeliveryQueueItem(t, ctx, s, root)
			}
			item.Stage = DeliveryStageParked
			item.Blocker = tc.blocker
			if err := s.UpdateDeliveryQueueItem(ctx, root, item); err != nil {
				t.Fatal(err)
			}
			target := item
			target.Stage = DeliveryStageGating
			_, _, err = s.RetryDeliveryQueueItem(ctx, root, target, item.Blocker)
			if tc.refused {
				if !errors.Is(err, ErrDeliveryRetryLimit) {
					t.Fatalf("error=%v", err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
			got := readDeliveryQueueItem(t, ctx, s, root)
			if got.RetryCount != tc.want {
				t.Fatalf("item=%+v want count=%d", got, tc.want)
			}
			if !tc.refused && (got.Blocker != "" || got.Stage != DeliveryStageGating) {
				t.Fatalf("item=%+v", got)
			}
			if tc.refused && (got.Blocker != item.Blocker || got.Stage != item.Stage) {
				t.Fatalf("refused retry changed item=%+v", got)
			}
		})
	}
}
