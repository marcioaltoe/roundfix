---
task: task_02
spec: 0299-a-new-delivery-queue-resumes-an-item-whose-spec-is-archived-on-its-branch
status: completed
type: backend
complexity: medium
---

# Task 02: The Delivery Engine resumes a proven archived item at review and parks an unproven one

## Overview

A new queue's `queued` item always runs Delivery Revalidation after its item
branch is created or continued, so an item branch that already archived its
Spec parks `delivery-error` on the missing Spec folder. This Task gives the
Delivery Engine an optional prover of an archived item branch: a proven
archive records the item head as the candidate and enters `reviewing`
without revalidation or an Implement Run, an unproven one parks
`archived-item-unproven: <cause>` in the `finding` Park Class, and a Delivery
Retry of that park returns the item to `queued`.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-22,
   [`references/2026-10-22-a-new-delivery-queue-parks-an-item-whose-spec-is-archived-on-its-branch.md`](references/2026-10-22-a-new-delivery-queue-parks-an-item-whose-spec-is-archived-on-its-branch.md)
   (`_prd.md` → Sources 1), by implementing API Contracts 1-4 of
   `_techspec.md` in `internal/delivery/engine.go` and
   `internal/delivery/park_class.go`.
2. MUST add `ArchivedItemProof`, `ArchivedItemProver` and
   `EngineDependencies.ArchivedItems` exactly as `_techspec.md` → Interfaces
   sketches them; a nil prover keeps the `queued` path's behavior unchanged.
3. MUST call the prover in the `queued` branch of `Engine.advanceItem` after
   `CreateItemBranch` records the branch and worktree and before
   `Revalidate`, with the four outcomes of API Contract 2, including the exact
   log line `roundfix: Delivery Queue item <slug>: item branch <branch>
   already archives the Spec at <commit>; resuming at reviewing`, and continue
   the stage loop from `reviewing` for a proven archive.
4. MUST add `BlockerArchivedItemUnproven = "archived-item-unproven"`, classify
   it `ParkClassFinding` in `ClassifyPark` with the next action `in
   <worktree>, make the item branch hold the Spec folder or one exact archive
   commit of the Spec, then run roundfix deliver retry <slug>` (API Contract
   3), and make `Engine.Retry` return such a park to `queued` as API Contract
   4 states, after the checks that precede the `prerequisite-unmerged` rule.
5. MUST keep every other function of `internal/delivery/engine.go` byte for
   byte, in particular `checkCandidate`, `createPullRequest` and `NewEngine`'s
   check timeout handling, which Spec 0298 changes; only `advanceItem`,
   `Retry`, the blocker constants, the new types and the dependency wiring in
   `EngineDependencies`, `Engine` and `NewEngine` change.
6. MUST add, in the new file `internal/delivery/archived_item_resume_test.go`,
   the five tests `_techspec.md` → Testing Approach lists for that file, each
   top-level test calling `t.Parallel()` first, using the package's fake
   workflow and boundary with a fake prover and the real SQLite store under a
   temporary home, and waiting on no real time:
   `TestANewQueueResumesAProvenArchivedItemAtReview`,
   `TestANewQueueParksAnUnprovenArchivedItem`,
   `TestANewQueueRevalidatesAnItemWithoutAnArchive`,
   `TestANewQueueParksAnArchivedItemProofError` and
   `TestARetryReturnsAnUnprovenArchivedItemToQueued`. The resume test asserts
   that the first review call and the archive call receive the proven head,
   which a resume at `gating` would fail; the retry test asserts the stored
   stage `queued` and the retry count raised by one.
7. MUST add to the table of `TestParkClassesAndNextActions` in
   `internal/delivery/park_class_test.go` a row for
   `BlockerArchivedItemUnproven` with class `finding` and its exact next
   action for a recorded worktree, and keep every existing row.

## Subtasks

- [ ] Add the prover types and the optional engine dependency.
- [ ] Resume a proven archive at review and park an unproven one in the queued path.
- [ ] Add the blocker, its Park Class and its retry rule.
- [ ] Add the engine tests and the Park Class row.

## Acceptance Criteria

- [ ] A proven archived item reaches `merged` through the fakes with no Implement or revalidation call, and its first review sees the proven head.
- [ ] A refusal parks `archived-item-unproven: <refusal>` in the `finding` class with the stated next action; a prover error parks `delivery-error: prove archived item: <error>`.
- [ ] The zero proof and a nil prover revalidate and run as before.
- [ ] A retry of the new park returns the item to `queued` and counts one retry.

## Consumers

- checkpoint: task_03

## Context

- instruction: `docs/adr/0312-a-new-delivery-queue-resumes-an-archived-item-branch-at-review.md`
- instruction: `docs/adr/0229-an-operator-archive-resumes-any-park-from-the-run-start.md`
- instruction: `internal/delivery/engine_test.go`
- instruction: `internal/delivery/retry_test.go`
- instruction: `internal/store/delivery.go`
- interface: `internal/delivery/engine.go`
- interface: `internal/delivery/park_class.go`
- interface: `internal/delivery/park_class_test.go`
- creates: `internal/delivery/archived_item_resume_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestANewQueueResumesAProvenArchivedItemAtReview|TestANewQueueParksAnUnprovenArchivedItem|TestANewQueueRevalidatesAnItemWithoutAnArchive|TestANewQueueParksAnArchivedItemProofError|TestARetryReturnsAnUnprovenArchivedItemToQueued|TestParkClassesAndNextActions|TestEveryBlockerHasAParkClass)$' ./internal/delivery 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestANewQueueResumesAProvenArchivedItemAtReview TestANewQueueParksAnUnprovenArchivedItem TestANewQueueRevalidatesAnItemWithoutAnArchive TestANewQueueParksAnArchivedItemProofError TestARetryReturnsAnUnprovenArchivedItemToQueued TestParkClassesAndNextActions TestEveryBlockerHasAParkClass; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; grep -q 'archived-item-unproven' internal/delivery/park_class_test.go || { echo 'missing Park Class row' >&2; exit 1; }; GOOS=windows go vet ./internal/delivery` — expected: exit 0. Before this Task the five new tests do not exist, so the command fails at its first pass check. Afterwards the resume, park, revalidation, proof-error and retry tests pass with the Park Class tests, the Park Class table holds the new row, and the package's tests compile for Windows.

## References

- `_prd.md` → Sources; Goals; Core Features 1-2; Success Metric 3; Success Metric 4; Success Metric 5
- `_techspec.md` → Current behavior; API Contracts 1-4; Invariants 1-4 and 6; Testing Approach; Build Order 2
- ADR-0223
- ADR-0229
- ADR-0312


## Result

Implemented API Contracts 1–4 within this Task's engine slice. The optional
`ArchivedItems` dependency proves the item after branch/worktree creation and
before revalidation. A proven archive records exactly its head as the candidate
and enters `reviewing`; refusal parks `archived-item-unproven`, while proof
errors and an empty archived head take the existing `delivery-error` path.
The new blocker has the specified finding class and worktree-specific next
action. Its retry shares the counted return-to-queued path after the existing
merge observation, token, deadline and retry-limit checks, without requiring
workspace, recovery or revalidation dependencies.

Acceptance evidence from focused implementation checks:

- Proven archive: `TestANewQueueResumesAProvenArchivedItemAtReview` reaches
  `merged` through the existing workflow/forge fakes and real temporary SQLite
  store. It asserts no Implement or revalidation call, the proven head in both
  review and archive calls, the exact candidate list and resume log, and the
  ordered review-through-merge actions. The prover also checks that its
  arguments match the already recorded workspace.
- Refusal and errors: `TestANewQueueParksAnUnprovenArchivedItem` asserts the
  refusal blocker, finding class and exact next action, with no candidate work
  or push. `TestANewQueueParksAnArchivedItemProofError` covers a prover error
  and the empty-head error, with exact delivery-error text and no candidate
  work. The Park Class table retains every old row and adds the new blocker
  and exact next-action assertion for bare and cause-suffixed blockers.
- Compatibility: `TestANewQueueRevalidatesAnItemWithoutAnArchive` covers the
  zero proof and nil prover separately; each revalidates once, invokes
  Implement once and reaches `merged`, with the expected start ordering.
- Retry: `TestARetryReturnsAnUnprovenArchivedItemToQueued` covers both the
  exact blocker and its colon/cause form with workspace, recovery and
  revalidator absent; the persisted stage is `queued`, the blocker clears,
  and the retry count increases by one.

Focused command:

```sh
go test ./internal/delivery -run '^(TestANewQueue|TestARetryReturnsAnUnproven|TestParkClass|TestEveryBlocker|TestDeliveryEngineMerges|TestDeliveryActions|TestRetry|TestArchivedRetry|TestExistingGenericPark)' -count=1
```

Outcome: exit 0, `ok roundfix/internal/delivery`. All new top-level tests call
`t.Parallel()` first and introduce no timed waits. The initial focused run
before implementation failed to compile because the proof types and engine
field did not exist. During test development, the workspace fixture was
corrected to record provisioning as the real workflow does, and retry
subtests received explicit names to avoid a generated `#00` in a SQLite URI
path. The reported successful run follows those corrections.

A byte comparison against the starting engine confirms that every existing
function except `advanceItem`, `Retry` and `NewEngine` is unchanged, and that
`NewEngine` differs only by the new dependency assignment. Check timeout and
Pull Request behavior are untouched. `git diff --check` reports no whitespace
errors.

Execution environment: the sandbox rejects the supplied symlinked writable
root; task-scoped commands ran with escalation in its resolved checkout.
The pre-existing task-file status change was preserved. No other Task or Task
Graph was edited, and no commit, push or Pull Request was made. Command-workflow
proof and CLI wiring remain task_03's slice. Declared Task Verification,
including its Windows vet step, and repository Verification were not run;
the Daemon owns them and Task settlement.

## Carry-forward provenance

- Source Run: `run_20261022T213854Z_c379322ae483e8bf`
- Source commit: `5c7e429a697273bd6d214f5ae4cb894847bf1c8d`
