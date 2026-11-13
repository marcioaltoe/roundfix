---
task: task_03
spec: 0389-a-parked-item-always-has-a-way-back-through-a-delivery-retry
status: completed
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

## Result

Skills loaded: implement-task, coding-guidelines, golang-cli, golang-testing, testing-boss, systematic-debugging, no-workarounds, agentic-cli-design, tech-writer, evidence-gate.

Starting-tree characterization (before implementation): a temporary Go overlay
added `TestCharacterizeCheckingAnUnpushedCandidate` without changing repository
sources. An item at `checking` with candidates `pushed-head`, `resolved-head`
and Pull Request `318` holding `pushed-head` ended at `parked`, blocker
`review-stale`, with zero pushes and remote head `pushed-head`.

- Focused command: `go test -count=1 -v -overlay /private/tmp/roundfix-task03-characterization-ex6bxa04/overlay.json ./internal/delivery -run '^TestCharacterizeCheckingAnUnpushedCandidate$'` — exit 0; the asserted starting behavior was observed with the real temporary SQLite store and existing offline workflow/Pull Request fakes.

Implementation and acceptance evidence:

- Criterion 1: `checkCandidate` recognizes only a non-empty head among candidates before the last, writes API Contract 4's exact log line, and persists `gating`. The earlier-candidate subtest observes the same PR `318` merge as `merge-318` after exactly one push, remote head `resolved-head`, and ordered events `checks`, `gate`, `authorization`, `publication`, `push`, `pull-request`, `checks`, `authorization`, `merge`. It asserts the exact owner log line naming both heads and no unmatched action intents.
- Criterion 2: the foreign-head subtest observes `review-stale`, zero pushes and merges, unchanged remote head, only a `checks` event, and no return-to-gating log. No other stale park, publication/push implementation, `createPullRequest`, retry implementation or exit code changed.
- Criterion 3: `ClassifyPark` retains the `review` class and supplies API Contract 5's exact action for `review-stale`, including colon-suffixed details. The new Park Class test checks both the action and Pending Question answer at `/worktrees/example`; the two existing CLI status tests assert their whole stdout, empty stderr and exit 0 with the new action (including the empty-worktree `in ,` case). `review-blocked` retains the generic answer.
- The guide carries Text C immediately after Text A and Text D immediately after the `secret-found` paragraph. A focused Python check confirmed both complete texts word for word after whitespace normalization and their exact anchors.

Declared expectation changes only:

- `TestExistingGenericParkAnswersAreUnchanged`: removed only `BlockerReviewStale` from the generic-answer list.
- `TestPendingQuestionAnswersEachBlockerClass`: changed only the `other blocker` fixture to `BlockerReviewBlocked`.
- `TestDeliverStatusPrintsTheItemWorktree`: changed only the `Park:` and `Answer:` expected actions.
- `TestDeliverStatusPrintsAParkLinePerParkedItem`: changed only the last item's expected `Park:` action.
- Searched repository Go tests for `review-stale`, `BlockerReviewStale` and `resolve the blocker`; no additional test pinned the generic `review-stale` answer.

Focused checks:

- Before the production change, `go test -count=1 -v ./internal/delivery -run '^(TestCheckingReturnsACandidateThePullRequestLacksToTheGate|TestAReviewStaleParkNamesItsNextAction)$'` — exit 1 as expected: earlier candidate still parked with zero pushes, both stale-action cases still returned the generic answer; the foreign-head subtest passed.
- After the production change, `go test -count=1 -v ./internal/delivery -run '^(TestCheckingReturnsACandidateThePullRequestLacksToTheGate|TestAReviewStaleParkNamesItsNextAction|TestExistingGenericParkAnswersAreUnchanged|TestPendingQuestionAnswersEachBlockerClass)$'` — exit 0, all selected tests and subtests passed.
- `go test -count=1 -v ./internal/cli -run '^(TestDeliverStatusPrintsTheItemWorktree|TestDeliverStatusPrintsAParkLinePerParkedItem)$'` — exit 0, both whole-transcript tests passed.
- `go test -count=1 -v ./internal/delivery -run '^(TestDeliveryEngineParksAStaleReview|TestDeliveryParksAnInexactPullRequestRecordCommit|TestAnAlreadyMergedPullRequestMustCarryTheReviewedHead|TestEveryBlockerHasAParkClass|TestParkClassesAndNextActions|TestParkClassRequiresAnExactBlockerOrColonPrefix)$'` — exit 0, all selected tests and subtests passed.
- `git diff --check` — exit 0. Scope inspection found only declared task_03 paths; the authored Task content and Daemon-owned `in_progress` status were preserved. New tests call `t.Parallel()` first and use offline fakes with a temporary SQLite store, without network access or elapsed-time waits.

The required Behavior Surface Record regeneration is the final implementation
step after the last guide edit; its outcome is recorded below. The declared
Verification commands, repository Verification and incremental Verification
remain for the Daemon. No Task Graph, other Task, status, commit, push or Pull
Request mutation was performed. No follow-up work was found in this slice.

- Final implementation step: `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` — exit 0; regenerated `docs/references/behavior-surfaces.json` after the final guide edit.
