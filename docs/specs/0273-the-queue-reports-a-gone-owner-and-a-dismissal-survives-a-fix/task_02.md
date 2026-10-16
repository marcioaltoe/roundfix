---
task: task_02
spec: 0273-the-queue-reports-a-gone-owner-and-a-dismissal-survives-a-fix
status: completed
type: backend
complexity: medium
---

# Task 02: A deliver command in an item worktree names the checkout that owns the queue

## Overview

Run inside an item worktree, `roundfix deliver retry` and `roundfix deliver
status` say that no Delivery Queue exists for the worktree and stop there
(`_techspec.md` → Current behavior). This Task adds the read-only store query
and the refusal of API Contract 2, used by `status`, `resume`, `retry` and
`stop`, and documents it in the `deliver` user guide. It is verifiable on its
own through those commands in a linked worktree recorded as an item worktree.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-14,
   [a dismissal after a committed fix needs a reset](references/2026-10-14-a-dismissal-after-a-committed-fix-needs-a-reset.md),
   for its second symptom: the retry refusal inside the item worktree names
   the checkout that owns the queue.
2. MUST add `DeliveryQueueItemWorktrees` and `DeliveryQueueItemWorktree` to
   `internal/store/delivery.go` as API Contract 2 defines them: a read of the
   recorded item rows with a non-empty `worktree`, ordered by `git_root` and
   `position`, that writes nothing and needs no migration, and that works on a
   reader store.
3. MUST add `missingDeliveryQueue` to `internal/cli/deliver.go` and use it
   wherever `status`, `resume` and `stop` refuse a missing queue today; in
   `runDeliverRetry`, read the queue once after opening the store and, when it
   is missing, refuse through `printDeliverRetryRefusal` with
   `retry Delivery Queue item "<slug>": ` and the helper's reason before the
   Delivery Engine runs. Paths compare through
   `config.CanonicalizeResolvedPath` on both sides (ADR-0268); the first match
   in the query's order names the item and its owning checkout. A failed
   lookup or no match keeps today's reason byte for byte. Roundfix MUST NOT
   act on the owning checkout's queue from the item worktree: no item changes,
   no owner record changes and no owner starts.
4. MUST NOT change the `deliver` help text, the `Retry refused` and
   `Preflight failed` layouts, their exit code `2`, or any refusal of a
   repository that owns its queue.
5. MUST add the tests named in Verification, each calling `t.Parallel()`
   first. In the new `internal/store/delivery_item_worktrees_test.go`:
   `TestDeliveryQueueItemWorktreesListsEveryRecordedWorktree` (two queues, an
   item without a worktree left out, the order of API Contract 2). In the new
   `internal/cli/deliver_item_worktree_refusal_test.go`:
   `TestDeliverRetryInAnItemWorktreeReproducesSurfaceTranscript` and
   `TestDeliverStatusInAnItemWorktreeReproducesSurfaceTranscript`, each
   asserting Surface Transcript 2 or 3 whole, with a worktree made by
   `git worktree add` and recorded through `RecordDeliveryQueueItemWorktree`;
   `TestDeliverResumeAndStopInAnItemWorktreeNameTheOwningCheckout` (both
   refuse with the reason and the stubbed owner start is never called);
   `TestDeliverItemWorktreeRefusalComparesResolvedPaths` (the worktree
   recorded through a symbolic link to its parent and the command run in the
   resolved path); and `TestDeliverMissingQueueRefusalKeepsItsTextElsewhere`
   (an unrelated repository keeps today's reason for all four commands).
   Every existing `deliver` and store test MUST pass unchanged.
6. MUST describe in `docs/user-guide/commands/deliver.md`, beside the retry
   refusal it already documents, that the queue commands run in an item
   worktree name the owning checkout, with the phrase
   `this checkout is the item worktree of`, and that they do not act on that
   queue; then re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.

## Subtasks

- [ ] Add the store query and its test.
- [ ] Add `missingDeliveryQueue` and use it in the four commands.
- [ ] Add the five command tests.
- [ ] Describe the refusal in the `deliver` guide and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] `deliver retry`, `status`, `resume` and `stop` in a recorded item worktree name the item and the owning checkout and change nothing.
- [ ] The match holds when the recorded path and the checkout path spell the same directory differently.
- [ ] Outside any item worktree the refusal keeps today's text.
- [ ] The `deliver` guide describes the refusal.

## Context

- instruction: `docs/adr/0279-status-reports-a-queue-owner-proven-gone-and-a-dismissal-survives-a-fix-that-leaves-its-anchor.md`
- instruction: `docs/adr/0268-a-recorded-worktree-path-is-resolved-within-its-worktree-root.md`
- interface: `internal/store/delivery.go`
- interface: `internal/cli/deliver.go`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `docs/references/behavior-surfaces.json`
- creates: `internal/store/delivery_item_worktrees_test.go`
- creates: `internal/cli/deliver_item_worktree_refusal_test.go`

## Consumers

- `DeliveryQueueItemWorktrees` — consumers: `./internal/store`, `./internal/cli`
- Missing-queue refusal of the four `deliver` commands — consumers: `./internal/cli`
- Behavior Surface Record of the `deliver` guide — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -v -run '^TestDeliveryQueueItemWorktreesListsEveryRecordedWorktree$' ./internal/store 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestDeliveryQueueItemWorktreesListsEveryRecordedWorktree ' || { printf 'missing passing store test\n' >&2; exit 1; }; go test -count=1 -run 'Delivery' ./internal/store || exit 1; out="$(go test -count=1 -v -run '^(TestDeliverRetryInAnItemWorktreeReproducesSurfaceTranscript|TestDeliverStatusInAnItemWorktreeReproducesSurfaceTranscript|TestDeliverResumeAndStopInAnItemWorktreeNameTheOwningCheckout|TestDeliverItemWorktreeRefusalComparesResolvedPaths|TestDeliverMissingQueueRefusalKeepsItsTextElsewhere)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDeliverRetryInAnItemWorktreeReproducesSurfaceTranscript TestDeliverStatusInAnItemWorktreeReproducesSurfaceTranscript TestDeliverResumeAndStopInAnItemWorktreeNameTheOwningCheckout TestDeliverItemWorktreeRefusalComparesResolvedPaths TestDeliverMissingQueueRefusalKeepsItsTextElsewhere; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -run 'Deliver' ./internal/cli || exit 1; tr -s '[:space:]' ' ' < docs/user-guide/commands/deliver.md | grep -qF -- 'this checkout is the item worktree of' || { printf 'deliver guide lacks the item worktree refusal\n' >&2; exit 1; }; go test -count=1 -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract || exit 1` — expected: exit 0. Before this Task the six tests do not exist and the guide lacks the phrase, so the command fails; afterwards they, every existing `Delivery` store test and every `Deliver` command test pass and the Behavior Surface Record is current.

## References

- `_prd.md` → Goals; User Story 3; Core Feature 3; Core Feature 4; Success Metric 3; User Experience
- `_techspec.md` → API Contract 2; API Contract 4; Surface Transcript 2; Surface Transcript 3; Invariant 4; Current behavior; Testing Approach; Build Order 2
- ADR-0279; ADR-0268

## Result

Implemented the Task 02 slice for Daemon Verification. Status remains
Daemon-owned; the authored Verification and repository Verification were not
run. No commit, push or pull request was made.

### Implementation and acceptance evidence

- **Recorded item worktree names its owning checkout without acting on it:**
  added the read-only `DeliveryQueueItemWorktrees` query and its result type,
  without a migration. It excludes empty worktrees and orders by `git_root`,
  then `position`. All four commands use `missingDeliveryQueue`; retry checks
  for the queue before constructing the Delivery Engine. The two transcript
  tests compare all stdout, stderr and exit-code output against Surface
  Transcripts 2 and 3. The resume/stop test checks both complete refusals.
  These tests snapshot the owning queue, including item and owner fields,
  and assert zero engine constructions and owner starts.
  `TestDeliverItemWorktreeRefusalUsesFirstRecordedMatch` also proves that the
  first queue and item in query order supplies the owning-checkout reason.
- **Different spellings of the same directory:**
  `TestDeliverItemWorktreeRefusalComparesResolvedPaths` creates a real linked
  worktree, records it through a symbolic link to its parent, and runs each
  of the four commands from the resolved path. Both sides of the production
  comparison use `config.CanonicalizeResolvedPath`.
- **Existing refusal text elsewhere:**
  `TestDeliverMissingQueueRefusalKeepsItsTextElsewhere` checks complete
  refusals for all four commands in an unrelated Git repository. Its
  unavailable-lookup subtest checks that a closed reader store preserves the
  original reason even for a matching worktree. Retry read errors retain the
  existing `retry Delivery Queue item` prefix and `Retry refused` layout;
  help text and refusal printers were left intact.
- **Guide describes the refusal:** added the owning-checkout reason and
  instructions beside the documented retry refusal, including the phrase
  `this checkout is the item worktree of` and the guarantee that the command
  does not act on the owning queue. Ran the required record generator; its
  only generated change is the deliver guide fingerprint in
  `docs/references/behavior-surfaces.json`.
- **Store contract:**
  `TestDeliveryQueueItemWorktreesListsEveryRecordedWorktree` uses two queues,
  reverse insertion order, item names whose lexical order differs from
  position, and an omitted item with no worktree in each queue. It queries
  through `OpenReader` and confirms both queues, including owner records,
  stay unchanged. Every new top-level test calls `t.Parallel()` first.

### Focused checks

- Before the production edit,
  `rtk proxy go test ./internal/cli -run 'TestDeliver.*(ItemWorktree|MissingQueue)' -count=1`
  exited 1: status/resume/stop printed only the old missing-queue reason, and
  retry reached the stubbed engine and owner instead of refusing first.
- After the production edit,
  `rtk proxy go test ./internal/store -run 'WorktreesLists' -count=1`
  exited 0.
- After adding fallback and ordering coverage,
  `rtk proxy go test ./internal/cli -run 'TestDeliver.*(ItemWorktree|MissingQueue|RetryPrintsARefusal|RetryRefusalStartsNoOwner|RetryHandsTheItem|RetryStartsAnOwner|Status.*Owner)' -count=1`
  exited 0.
- Final focused command check,
  `rtk proxy go test ./internal/cli -run 'TestDeliver(Retry|Resume|Stop|Status|ItemWorktree|MissingQueue)' -count=1`,
  exited 0, covering the new regressions and existing tests for the four
  command flows.
- Required generation:
  `rtk proxy go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
  exited 0.
- `rtk proxy git diff --check` exited 0 before this Result was appended.

The sandbox could not start processes because the configured writable root
contains a symbolic link. Approved elevated execution was used for these
commands. The initial worktree change was only this Task's `status:
in_progress`; it was preserved. No other Task or Task Graph was edited.
The Roundfix skill reference update belongs to the later documentation Task
and was left for that slice. Full declared Verification and settlement remain
with the Daemon.

## Carry-forward provenance

- Source Run: `run_20261016T175249Z_782526ccf3bf55ac`
- Source commit: `2e6da25dfdc9d417a42c0078ed38ca2124422249`
