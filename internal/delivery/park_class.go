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
func ClassifyPark(queue store.DeliveryQueue, item store.DeliveryQueueItem) ParkClassification {
	classification := ParkClassification{
		Class: ParkClassUnclassified,
		Next:  "resolve the blocker, then run roundfix deliver retry " + item.SpecSlug,
	}
	blocker, detail, _ := strings.Cut(item.Blocker, ":")
	switch blocker {
	case BlockerPrerequisiteUnmerged:
		classification.Class = ParkClassDependency
		prerequisites := strings.Split(strings.TrimSpace(detail), ", ")
		for _, slug := range prerequisites {
			for _, prerequisite := range queue.Items {
				if prerequisite.SpecSlug == slug && prerequisite.Stage == store.DeliveryStageParked {
					classification.Next = "run roundfix deliver retry " + slug + "; " + item.SpecSlug + " returns to the queue once " + slug + " is merged"
					return classification
				}
			}
		}
		classification.Next = "deliver or merge " + strings.TrimSpace(detail) + ", then run roundfix deliver retry " + item.SpecSlug
	case BlockerPullRequestConflict:
		classification.Class = ParkClassConflict
		classification.Next = "merge the default branch into the item branch in " + item.Worktree + ", resolve " + strings.TrimSpace(detail) + ", commit, then run roundfix deliver retry " + item.SpecSlug
	case BlockerRuntimeInfrastructure:
		classification.Class = ParkClassEnvironment
		classification.Next = "run roundfix deliver retry " + item.SpecSlug + "; a retry from runtime-infrastructure is not counted"
	case BlockerLinuxLaneContainer:
		classification.Class = ParkClassEnvironment
		classification.Next = "repair the Linux Lane container cause in the blocker (Docker, the image, or delivery.linux_lane_image on the default branch), then run roundfix deliver retry " + item.SpecSlug
	case BlockerQAEnvironmentPartial:
		classification.Class = ParkClassEnvironment
		classification.Next = "(cd " + item.Worktree + " && roundfix reconcile " + item.RunID +
			" --carry-forward), satisfy the environment-blocked QA rows, run roundfix archive " + item.SpecSlug +
			" --qa-override --approval <source> --reason <text>, then run roundfix deliver retry " + item.SpecSlug
	case BlockerChecksTimeout, BlockerItemWorktreeMissing, BlockerDeliveryError:
		classification.Class = ParkClassEnvironment
	case BlockerFlakyCheck:
		classification.Class = ParkClassFlakyCheck
		classification.Next = "the check failed twice in " + strings.TrimSpace(detail) + ", which " + item.SpecSlug +
			" did not change; fix or re-run it, then run roundfix deliver retry " + item.SpecSlug
	case BlockerArchivedItemUnproven:
		classification.Class = ParkClassFinding
		classification.Next = "in " + item.Worktree + ", make the item branch hold the Spec folder or one exact archive commit of the Spec, then run roundfix deliver retry " + item.SpecSlug
	case BlockerSecretFound:
		classification.Class = ParkClassFinding
		classification.Next = "in " + item.Worktree + ", rewrite the item branch from the earliest commit the blocker names so that each named line holds a value no Secret Pattern matches and nothing else in the final tree changes, then run roundfix deliver retry " + item.SpecSlug
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
	case BlockerQueueTokenCeiling:
		classification.Class = ParkClassBudget
		classification.Next = "record a new queue for the remaining Specs with roundfix deliver start and a higher --max-tokens"
	case BlockerReviewBlocked, BlockerReviewStale:
		classification.Class = ParkClassReview
	case BlockerUnauthorized:
		classification.Class = ParkClassAuthorization
	}
	return classification
}
