package delivery

import (
	"context"
	"fmt"
	"os"
	"slices"
	"sync"
	"time"

	"roundfix/internal/store"
)

type reviewSlotState string

const (
	reviewSlotNone     reviewSlotState = "none"
	reviewSlotClaimed  reviewSlotState = "claimed"
	reviewSlotWaiting  reviewSlotState = "waiting"
	reviewSlotInFlight reviewSlotState = "in-flight"
)

type reviewSlot struct {
	item  store.DeliveryQueueItem
	state reviewSlotState
}

// One Run owns this pool, its poller and every review. The mutex protects
// admission and returned items; only the pass touches the workspace.
type reviewSlots struct {
	mu         sync.Mutex
	engine     *Engine
	gitRoot    string
	ctx        context.Context
	cancel     context.CancelFunc
	workers    sync.WaitGroup
	stopped    bool
	held       int
	generation uint64
	items      map[string]reviewSlot
	ready      map[string]store.DeliveryQueueItem
	err        error
	changed    chan struct{}
}

func newReviewSlots(engine *Engine, gitRoot string) *reviewSlots {
	return &reviewSlots{
		engine: engine, gitRoot: gitRoot,
		items: make(map[string]reviewSlot), ready: make(map[string]store.DeliveryQueueItem),
		changed: make(chan struct{}, 1),
	}
}

func (slots *reviewSlots) start(ctx context.Context) {
	slots.ctx, slots.cancel = context.WithCancel(ctx)
	slots.workers.Add(1)
	go func() {
		defer slots.workers.Done()
		ticker := time.NewTicker(slots.engine.turnPollInterval)
		defer ticker.Stop()
		for {
			select {
			case <-slots.ctx.Done():
				return
			case <-ticker.C:
				slots.poll()
			}
		}
	}()
}

func (slots *reviewSlots) stop() {
	slots.mu.Lock()
	slots.stopped = true
	slots.cancel()
	slots.mu.Unlock()
	slots.workers.Wait()
}

func (slots *reviewSlots) state(slug string) reviewSlotState {
	slots.mu.Lock()
	defer slots.mu.Unlock()
	if entry, ok := slots.items[slug]; ok {
		return entry.state
	}
	return reviewSlotNone
}

// claimForPass checks returns, holds and admission atomically. Otherwise a
// review could return between two checks and the pass would prepare its
// archiving worktree instead of consuming the returned item.
func (slots *reviewSlots) claimForPass(item store.DeliveryQueueItem) (_ store.DeliveryQueueItem, returned, held bool) {
	slots.mu.Lock()
	defer slots.mu.Unlock()
	if ready, ok := slots.ready[item.SpecSlug]; ok {
		delete(slots.ready, item.SpecSlug)
		slots.generation++
		return ready, true, false
	}
	if _, ok := slots.items[item.SpecSlug]; ok {
		return item, false, true
	}
	if item.Stage == store.DeliveryStageReviewing {
		slots.items[item.SpecSlug] = reviewSlot{state: reviewSlotClaimed}
	}
	return item, false, false
}

func (slots *reviewSlots) release(slug string) {
	slots.mu.Lock()
	defer slots.mu.Unlock()
	if slots.items[slug].state == reviewSlotClaimed {
		delete(slots.items, slug)
	}
}

func (slots *reviewSlots) take(item store.DeliveryQueueItem) {
	slots.mu.Lock()
	taken := slots.takeLocked(item)
	slots.mu.Unlock()
	if taken {
		slots.engine.reviewSlotsChanged()
	}
}

// takeLocked also accepts the pass's claim, which prevents the poller from
// reviewing an item while its resume prologue prepares the worktree.
func (slots *reviewSlots) takeLocked(item store.DeliveryQueueItem) bool {
	if slots.stopped || slots.ctx.Err() != nil {
		return false
	}
	if entry, ok := slots.items[item.SpecSlug]; ok && entry.state != reviewSlotClaimed {
		return false
	}
	if _, ok := slots.ready[item.SpecSlug]; ok {
		return false
	}
	if slots.held >= slots.engine.reviewCapacity {
		fmt.Fprintf(slots.engine.log, "roundfix: Delivery Queue item %s: waiting for a review slot; %d of %d review slot(s) held\n", item.SpecSlug, slots.held, slots.engine.reviewCapacity)
	}
	slots.items[item.SpecSlug] = reviewSlot{item: item, state: reviewSlotWaiting}
	slots.generation++
	slots.startWaitingLocked()
	return true
}

func (slots *reviewSlots) startWaitingLocked() {
	if slots.stopped || slots.ctx.Err() != nil {
		return
	}
	waiting := make([]store.DeliveryQueueItem, 0)
	for _, entry := range slots.items {
		if entry.state == reviewSlotWaiting {
			waiting = append(waiting, entry.item)
		}
	}
	slices.SortFunc(waiting, func(a, b store.DeliveryQueueItem) int { return a.Position - b.Position })
	for _, item := range waiting {
		if slots.held >= slots.engine.reviewCapacity {
			break
		}
		slots.items[item.SpecSlug] = reviewSlot{item: item, state: reviewSlotInFlight}
		slots.held++
		slots.workers.Add(1)
		go slots.review(item)
	}
}

func (slots *reviewSlots) review(item store.DeliveryQueueItem) {
	defer slots.workers.Done()
	slots.engine.reviewSlotsChanged()
	err := slots.engine.reviewCandidate(slots.ctx, slots.gitRoot, &item)
	if err != nil && slots.ctx.Err() == nil {
		if parkErr := slots.engine.park(slots.ctx, slots.gitRoot, &item, BlockerDeliveryError+": "+err.Error()); parkErr != nil {
			slots.mu.Lock()
			slots.err = fmt.Errorf("deliver Spec %q: park review: %w", item.SpecSlug, parkErr)
			slots.mu.Unlock()
		}
	}
	slots.mu.Lock()
	delete(slots.items, item.SpecSlug)
	slots.held--
	slots.generation++
	if item.Stage == store.DeliveryStageArchiving {
		slots.ready[item.SpecSlug] = item
	}
	slots.startWaitingLocked()
	slots.mu.Unlock()
	slots.engine.reviewSlotsChanged()
}

func (slots *reviewSlots) snapshot() (busy, returned bool, generation uint64, err error) {
	slots.mu.Lock()
	defer slots.mu.Unlock()
	return len(slots.items) > 0, len(slots.ready) > 0, slots.generation, slots.err
}

func (slots *reviewSlots) poll() {
	// Read under the admission lock: an ending review cannot release its hold
	// between this snapshot and the take and admit a stale reviewing item twice.
	slots.mu.Lock()
	taken := 0
	queue, _, err := slots.engine.store.DeliveryQueue(slots.ctx, slots.gitRoot)
	if err != nil && slots.ctx.Err() == nil {
		slots.err = fmt.Errorf("poll Delivery Queue reviews: %w", err)
	}
	if err == nil {
		for _, item := range queue.Items {
			if item.Stage != store.DeliveryStageReviewing || !item.WorktreeProvisioned || item.Worktree == "" {
				continue
			}
			if _, held := slots.items[item.SpecSlug]; held {
				continue
			}
			info, err := os.Stat(item.Worktree)
			if err == nil && info.IsDir() {
				if slots.takeLocked(item) {
					taken++
				}
			}
		}
	}
	slots.mu.Unlock()
	for range taken {
		slots.engine.reviewSlotsChanged()
	}
	slots.engine.reviewSlotsChanged()
}
