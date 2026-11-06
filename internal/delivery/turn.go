package delivery

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"roundfix/internal/store"
)

// TurnState describes whether an item's stored step runs or waits.
type TurnState string

const (
	TurnAdvancing        TurnState = "advancing"
	TurnReviewing        TurnState = "reviewing"
	TurnReviewSlotWait   TurnState = "review-slot-wait"
	TurnPrerequisiteWait TurnState = "prerequisite-wait"
	TurnWait             TurnState = "turn-wait"
)

type ItemTurn struct {
	SpecSlug      string              `json:"spec"`
	Stage         store.DeliveryStage `json:"stage"`
	State         TurnState           `json:"turn"`
	Prerequisites []string            `json:"prerequisites,omitempty"`
	Checks        []string            `json:"checks,omitempty"`
}

type Turn struct {
	ReviewCapacity int        `json:"review_capacity"`
	Items          []ItemTurn `json:"items"`
}

type TurnRecorder interface {
	RecordTurn(ctx context.Context, turn Turn) error
}

// TurnLines reconciles a recorded turn with the queue status reads now.
func TurnLines(queue store.DeliveryQueue, turn Turn) []string {
	entries := make(map[string]ItemTurn, len(turn.Items))
	for _, entry := range turn.Items {
		entries[entry.SpecSlug] = entry
	}
	items := slices.Clone(queue.Items)
	slices.SortFunc(items, func(a, b store.DeliveryQueueItem) int { return a.Position - b.Position })
	var activePosition int
	advancing := false
	held := 0
	var ownerWaits, reviewWaits []string
	for _, item := range items {
		if item.Stage == store.DeliveryStageMerged || item.Stage == store.DeliveryStageParked {
			continue
		}
		entry, found := entries[item.SpecSlug]
		if !found {
			continue
		}
		if entry.State == TurnAdvancing {
			activePosition, advancing = item.Position, true
		}
		if entry.Stage != item.Stage {
			continue
		}
		switch entry.State {
		case TurnReviewing:
			held++
		case TurnReviewSlotWait:
			reviewWaits = append(reviewWaits, item.SpecSlug)
		case TurnWait:
			ownerWaits = append(ownerWaits, item.SpecSlug)
		}
	}
	if advancing {
		// The next pass visits positions after the active item before wrapping.
		var after, before []string
		for _, item := range items {
			if slices.Contains(ownerWaits, item.SpecSlug) {
				if item.Position > activePosition {
					after = append(after, item.SpecSlug)
				} else {
					before = append(before, item.SpecSlug)
				}
			}
		}
		ownerWaits = append(after, before...)
	}
	var lines []string
	for _, item := range items {
		if item.Stage == store.DeliveryStageMerged || item.Stage == store.DeliveryStageParked {
			continue
		}
		entry, found := entries[item.SpecSlug]
		text := "waiting for the owner to read it"
		if found && (entry.State == TurnAdvancing || entry.Stage == item.Stage) {
			switch entry.State {
			case TurnAdvancing:
				text = fmt.Sprintf("advancing at %s", item.Stage)
				if item.Stage == store.DeliveryStageChecking && len(entry.Checks) > 0 {
					text += fmt.Sprintf("; waiting for %d check(s): %s", len(entry.Checks), strings.Join(entry.Checks, "; "))
				}
			case TurnReviewing:
				text = fmt.Sprintf("reviewing; %d of %d review slot(s) held", held, turn.ReviewCapacity)
			case TurnReviewSlotWait:
				text = fmt.Sprintf("waiting for a review slot (%d of %d waiting); %d of %d review slot(s) held", slices.Index(reviewWaits, item.SpecSlug)+1, len(reviewWaits), held, turn.ReviewCapacity)
			case TurnPrerequisiteWait:
				text = "waiting for prerequisites: " + strings.Join(entry.Prerequisites, ", ")
			case TurnWait:
				text = fmt.Sprintf("waiting for the owner's turn (%d of %d waiting)", slices.Index(ownerWaits, item.SpecSlug)+1, len(ownerWaits))
			}
		}
		lines = append(lines, "Turn: "+item.SpecSlug+" "+text)
	}
	return lines
}

func (engine *Engine) advancingTurn(ctx context.Context, gitRoot, slug string) {
	if engine.turnRecorder == nil {
		return
	}
	engine.turnMu.Lock()
	defer engine.turnMu.Unlock()
	engine.advancing = slug
	delete(engine.prerequisiteWaits, slug)
	engine.recordTurnLocked(ctx, gitRoot)
}

func (engine *Engine) checkTurn(ctx context.Context, gitRoot, slug string, pending []string) {
	if engine.turnRecorder == nil {
		return
	}
	engine.turnMu.Lock()
	defer engine.turnMu.Unlock()
	if len(pending) == 0 {
		delete(engine.checkWaits, slug)
	} else {
		if engine.checkWaits == nil {
			engine.checkWaits = make(map[string][]string)
		}
		engine.checkWaits[slug] = slices.Clone(pending)
	}
	engine.recordTurnLocked(ctx, gitRoot)
}

func (engine *Engine) prerequisiteTurn(ctx context.Context, gitRoot, slug string, unmet []string) {
	if engine.turnRecorder == nil {
		return
	}
	engine.turnMu.Lock()
	defer engine.turnMu.Unlock()
	engine.prerequisiteWaits[slug] = slices.Clone(unmet)
	engine.advancing = ""
	engine.recordTurnLocked(ctx, gitRoot)
}

func (engine *Engine) recordTurn(ctx context.Context, gitRoot string) {
	if engine.turnRecorder == nil {
		return
	}
	engine.turnMu.Lock()
	defer engine.turnMu.Unlock()
	engine.recordTurnLocked(ctx, gitRoot)
}

// turnMu serializes both the fresh read and its write, so an older snapshot
// cannot replace a newer one. Slot callbacks arrive after releasing their lock.
func (engine *Engine) recordTurnLocked(ctx context.Context, gitRoot string) {
	queue, _, err := engine.store.DeliveryQueue(ctx, gitRoot)
	if err == nil {
		turn := Turn{ReviewCapacity: engine.reviewCapacity, Items: make([]ItemTurn, 0, len(queue.Items))}
		for _, item := range queue.Items {
			if item.Stage == store.DeliveryStageMerged || item.Stage == store.DeliveryStageParked {
				continue
			}
			entry := ItemTurn{SpecSlug: item.SpecSlug, Stage: item.Stage, State: TurnWait}
			state := reviewSlotNone
			if engine.reviewSlots != nil {
				state = engine.reviewSlots.state(item.SpecSlug)
			}
			switch {
			case item.SpecSlug == engine.advancing:
				entry.State = TurnAdvancing
				entry.Checks = slices.Clone(engine.checkWaits[item.SpecSlug])
			case state == reviewSlotInFlight:
				entry.State = TurnReviewing
			case state == reviewSlotWaiting:
				entry.State = TurnReviewSlotWait
			case item.Stage == store.DeliveryStageQueued && len(engine.prerequisiteWaits[item.SpecSlug]) > 0:
				entry.State = TurnPrerequisiteWait
				entry.Prerequisites = slices.Clone(engine.prerequisiteWaits[item.SpecSlug])
			}
			turn.Items = append(turn.Items, entry)
		}
		err = engine.turnRecorder.RecordTurn(ctx, turn)
	}
	if err != nil {
		fmt.Fprintf(engine.log, "roundfix: warning: Delivery Queue turn record: %v\n", err)
	}
}
