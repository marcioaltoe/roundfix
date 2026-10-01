package delivery

import (
	"strings"

	"roundfix/internal/store"
)

type ParkClass string

const (
	ParkClassDependency    ParkClass = "dependency"
	ParkClassConflict      ParkClass = "conflict"
	ParkClassEnvironment   ParkClass = "environment"
	ParkClassFlakyCheck    ParkClass = "flaky-check"
	ParkClassFinding       ParkClass = "finding"
	ParkClassBudget        ParkClass = "budget"
	ParkClassReview        ParkClass = "review"
	ParkClassAuthorization ParkClass = "authorization"
	ParkClassUnclassified  ParkClass = "unclassified"
)

type ParkClassification struct {
	Class ParkClass
	Next  string
}

// ClassifyPark supplies the same operator action to status and the Pending Question.
func ClassifyPark(_ store.DeliveryQueue, item store.DeliveryQueueItem) ParkClassification {
	classification := ParkClassification{
		Class: ParkClassUnclassified,
		Next:  "resolve the blocker, then run roundfix deliver retry " + item.SpecSlug,
	}
	blocker, detail, _ := strings.Cut(item.Blocker, ":")
	switch blocker {
	case BlockerChecksTimeout, BlockerItemWorktreeMissing, BlockerDeliveryError:
		classification.Class = ParkClassEnvironment
	case BlockerFlakyCheck:
		classification.Class = ParkClassFlakyCheck
		classification.Next = "the check failed twice in " + strings.TrimSpace(detail) + ", which " + item.SpecSlug +
			" did not change; fix or re-run it, then run roundfix deliver retry " + item.SpecSlug
	case BlockerRunUnresolved, BlockerReviewFindings, BlockerCorrectiveSpecRequired, BlockerGateFailed, BlockerChecksFailed, BlockerRevalidationFailed:
		classification.Class = ParkClassFinding
		if blocker == BlockerRevalidationFailed {
			classification.Next = "amend the Spec on its item branch in " + item.Worktree +
				", then run roundfix deliver retry " + item.SpecSlug
		}
	case BlockerRunBudgetExceeded, BlockerQueueDeadline:
		classification.Class = ParkClassBudget
		if blocker == BlockerQueueDeadline {
			classification.Next = "record a new queue for the remaining Specs with roundfix deliver start"
		}
	case BlockerReviewBlocked, BlockerReviewStale:
		classification.Class = ParkClassReview
	case BlockerUnauthorized:
		classification.Class = ParkClassAuthorization
	}
	return classification
}
