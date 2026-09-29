package delivery

import (
	"strings"

	"roundfix/internal/store"
)

type PendingQuestion struct {
	SpecSlug string
	Blocker  string
	Answer   string
	Waiting  int
}

func PendingQuestionFor(queue store.DeliveryQueue) (PendingQuestion, bool) {
	var selected store.DeliveryQueueItem
	parked := 0
	found := false
	for _, item := range queue.Items {
		if item.Stage != store.DeliveryStageParked {
			continue
		}
		parked++
		if !found || item.Position < selected.Position {
			selected = item
			found = true
		}
	}
	if !found {
		return PendingQuestion{}, false
	}

	answer := "resolve the blocker, then run roundfix deliver retry " + selected.SpecSlug
	switch {
	case strings.HasPrefix(selected.Blocker, BlockerRevalidationFailed):
		answer = "amend the Spec on its item branch in " + selected.Worktree +
			", then run roundfix deliver retry " + selected.SpecSlug
	case selected.Blocker == BlockerQueueDeadline:
		answer = "record a new queue for the remaining Specs with roundfix deliver start"
	}
	return PendingQuestion{
		SpecSlug: selected.SpecSlug,
		Blocker:  selected.Blocker,
		Answer:   answer,
		Waiting:  parked - 1,
	}, true
}
