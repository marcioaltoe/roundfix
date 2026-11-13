---
task: task_03
spec: 0389-a-parked-item-always-has-a-way-back-through-a-delivery-retry
status: pending
type: backend
complexity: medium
---

# Task 03: A Pull Request that holds an earlier candidate returns the item to the gate

## Overview

At `checking`, a Pull Request whose head is an earlier candidate of the item
never received the current candidate: the Queue Owner logs it and returns the
item to `gating`, where the repository gate runs again and publication pushes
the candidate. Any other Pull Request head still parks `review-stale`, and that
park now names its next action in `deliver status` and the Pending Question.
This answers item 2 of the Backlog Entry reported by quaestor on 2026-11-12
([`references/2026-11-12-a-retry-after-resolving-a-pull-request-conflict-finds-no-specs-root-and-no-pushed-head.md`](references/2026-11-12-a-retry-after-resolving-a-pull-request-conflict-finds-no-specs-root-and-no-pushed-head.md)).

## Requirements

1. MUST characterize first, in the Task's `## Result`: on the starting tree an
   item at `checking` with candidates `pushed-head` and `resolved-head` whose
   Pull Request holds `pushed-head` parks `review-stale` with no push.
2. MUST change `checkCandidate` in `internal/delivery/engine.go` and add
   `earlierCandidate` as `_techspec.md` → API Contract 4 says, with the log
   line word for word:
   `roundfix: Delivery Queue item <slug>: pull request #<n> holds earlier candidate <pull request head>, not candidate <candidate>; returning to gating to push it`.
   The item goes to `gating`, never straight to `publishing`.
3. MUST change `ClassifyPark` in `internal/delivery/park_class.go` as API
   Contract 5 says, with the next action word for word; `review-blocked` keeps
   the generic answer.
4. MUST create `internal/delivery/unpushed_candidate_test.go` with
   `TestCheckingReturnsACandidateThePullRequestLacksToTheGate` (subtests
   `earlier candidate` and `foreign head`) and add
   `TestAReviewStaleParkNamesItsNextAction` to
   `internal/delivery/park_class_test.go`, each calling `t.Parallel()` first,
   with the fixtures and assertions `_techspec.md` → Testing Approach lists.
5. MUST update, as `_techspec.md` → Declared breaks says, every test that pins
   the generic `review-stale` answer on the default branch when the Task runs:
   `TestExistingGenericParkAnswersAreUnchanged`,
   `TestPendingQuestionAnswersEachBlockerClass`,
   `TestDeliverStatusPrintsTheItemWorktree` and
   `TestDeliverStatusPrintsAParkLinePerParkedItem`, and any other it finds,
   which it names in the `## Result`; MUST change nothing else in them.
6. MUST change `docs/user-guide/commands/deliver.md` with Texts C and D of API
   Contract 6 (c) and (d), word for word, at the anchors they name.
7. MUST, as its last step and after the last guide edit, re-record the
   Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
8. MUST NOT change the other `review-stale` parks (review, archive, Pull
   Request record, merge), the push, `createPullRequest`, `deliver retry` or
   any exit code; no test reaches the network or waits on elapsed time.

## Subtasks

- [ ] Record the characterization.
- [ ] Return an item whose Pull Request holds an earlier candidate to `gating`.
- [ ] Give `review-stale` its next action and update the tests that pin the generic answer.
- [ ] Add the two tests.
- [ ] Add Texts C and D to the guide, then re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] An item whose Pull Request holds its earlier candidate is gated again, pushed once and merged, and the owner's log names both heads.
- [ ] A Pull Request head the item never had still parks `review-stale` with no push.
- [ ] Status and the Pending Question answer a `review-stale` park with its next action, and the guide says so.

## Context

- instruction: `docs/adr/0398-a-parked-item-always-has-a-way-back-through-a-delivery-retry.md`
- instruction: `docs/adr/0223-a-delivery-keeps-its-work-across-archive-requeue-and-review.md`
- instruction: `internal/delivery/engine_test.go`
- creates: `internal/delivery/unpushed_candidate_test.go`
- interface: `internal/delivery/engine.go`
- interface: `internal/delivery/park_class.go`
- interface: `internal/delivery/park_class_test.go`
- interface: `internal/delivery/question_test.go`
- interface: `internal/cli/deliver_test.go`
- interface: `internal/cli/deliver_park_status_test.go`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The `checking` stage and the `review-stale` Park Class answer — consumers: `./internal/delivery`, `./internal/cli`
- The `deliver` guide and the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 30m -v ./internal/delivery 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestCheckingReturnsACandidateThePullRequestLacksToTheGate/earlier_candidate TestCheckingReturnsACandidateThePullRequestLacksToTheGate/foreign_head TestAReviewStaleParkNamesItsNextAction; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the tests do not exist; afterwards they pass with the whole `./internal/delivery` package, the updated Park Class and Pending Question tests included.

- `grep -qF -- 'merge into the item branch any commit its Pull Request' internal/cli/deliver_test.go || { echo 'deliver_test.go still pins the generic answer' >&2; exit 1; }; grep -qF -- 'merge into the item branch any commit its Pull Request' internal/cli/deliver_park_status_test.go || { echo 'deliver_park_status_test.go still pins the generic answer' >&2; exit 1; }; go test -count=1 -timeout 40m -run 'DeliverStatus|DeliverRetry|PendingQuestion|FollowUp|Park' ./internal/cli` — expected: exit 0. Before this Task the two status tests pin the generic answer; afterwards they pin the new action, and the status, retry, Pending Question and follow-up tests of `./internal/cli` pass.

- `tr -s '[:space:]' ' ' < docs/user-guide/commands/deliver.md | grep -qF -- 'returning to gating to push it' || { echo 'deliver guide: Text C missing' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/deliver.md | grep -qF -- 'whose repository gate runs again before the push' || { echo 'deliver guide: Text C missing' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/deliver.md | grep -qF -- 'merge into the item branch any commit its Pull Request' || { echo 'deliver guide: Text D missing' >&2; exit 1; }; go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the texts are absent; afterwards the guide carries them and the documentation contracts, the Behavior Surface Record included, pass.

## References

- `_prd.md` → Goals; User Story 3; User Story 4; Core Features 3, 4 and 5; Success Metric 3; Success Metric 4; Success Metric 5; User Experience
- `_techspec.md` → Current behavior 5 and 6; API Contracts 4, 5 and 6; Invariants 4 and 5; Declared breaks; Testing Approach; Build Order 3
- ADR-0398
- [The adopted quaestor Backlog Entry](references/2026-11-12-a-retry-after-resolving-a-pull-request-conflict-finds-no-specs-root-and-no-pushed-head.md)
