package delivery

import (
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

	return PendingQuestion{
		SpecSlug: selected.SpecSlug,
		Blocker:  selected.Blocker,
		Answer:   ClassifyPark(queue, selected).Next,
		Waiting:  parked - 1,
	}, true
}
