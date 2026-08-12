---
task: task_06
spec: 0094-one-history-root-under-docs
status: completed # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: high
---

# Task 06: Apply and roll back a relocation inside the transaction

## Overview

Applying performs the planned relocations inside the existing mutation boundary,
verifies each destination against the identity recorded before the move, refuses a
collision without touching its siblings, and returns every moved file to its
source on rollback. This is the slice that makes the fleet migration real, and it
is the one whose failure modes a content write does not already exercise.

## Requirements

1. MUST perform each planned relocation and verify the destination's content
   identity against the one the plan recorded.
2. MUST fail the relocation when the destination's identity does not match what
   was recorded, rather than reporting a move that did not preserve the bytes.
3. MUST refuse a relocation whose destination already exists, name both paths and
   the reason, and MUST still perform the relocations of the other files.
4. MUST report, on a run where at least one relocation was refused, that not
   every relocation was performed.
5. MUST return every relocated file to its source on rollback, leaving no file at
   a destination.
6. MUST record the relocations in the recovery journal so an interruption between
   a source's removal and its destination's creation is recoverable by the
   existing phase machinery.
7. MUST leave a plan carrying no relocation applying exactly as it does today.
8. MUST NOT invoke Git; a relocation is a rename plus its verification.

## Subtasks

- [ ] Perform planned relocations within the existing mutation boundary.
- [ ] Verify each destination against its recorded content identity.
- [ ] Refuse a collision per file and report it without refusing siblings.
- [ ] Record relocations in the recovery journal and reverse them on rollback.
- [ ] Cover collision, identity mismatch, rollback, and the empty-ledger case.

## Acceptance Criteria

- [ ] Applying a plan with relocations leaves every file at its destination with
      the recorded content identity, and nothing at its source.
- [ ] A destination that already exists refuses that one relocation, names both
      paths, and the other files in the same plan still reach their destinations.
- [ ] A run with a refused relocation reports that not every relocation was
      performed.
- [ ] Rolling back after a partial relocation returns every moved file to its
      source, with no file left at a destination.
- [ ] Applying a plan with no relocation produces the same result as before this
      Task.
- [ ] No Git invocation occurs on the relocation path.

## Rehearsal Cases

- Case: a plan whose relocations all succeed; Observation: every destination
  holds the recorded content identity and every source is absent.
- Case: a plan whose destination already holds a file; Observation: that
  relocation is refused naming both paths, and its sibling relocations completed.
- Case: a destination whose content identity does not match the recorded one;
  Observation: the relocation fails rather than settling as performed.
- Case: a rollback after some relocations were performed; Observation: every
  moved file is back at its source and no destination holds it.
- Case: a plan with an empty relocation ledger; Observation: the applied result
  is identical to the pre-Task behavior.

## Verification

- `go test -count=1 ./internal/baseline -run 'HistoryMoveApply|HistoryMoveRollback|HistoryMoveCollision' -v > /tmp/0094-task-06.log 2>&1; s=$?; grep -q '^--- PASS: .*HistoryMoveApply' /tmp/0094-task-06.log && grep -q '^--- PASS: .*HistoryMoveRollback' /tmp/0094-task-06.log && grep -q '^--- PASS: .*HistoryMoveCollision' /tmp/0094-task-06.log || { cat /tmp/0094-task-06.log; exit 1; }; exit $s` — expected: exits 0 and the log names all three passing groups; fails when any one is missing.
- `! grep -qi 'no tests to run' /tmp/0094-task-06.log` — expected: exits 0, refusing a vacuous run.
- `go test -count=1 ./internal/baseline > /tmp/0094-task-06b.log 2>&1 && grep -q 'HistoryMove' internal/baseline/transaction.go || { cat /tmp/0094-task-06b.log; exit 1; }` — expected: exits 0, proving the existing transaction contract still holds once the ledger reaches the transaction. The suite alone passed before any work, so it is anchored to the change it is guarding.

## Context

- interface: `internal/baseline/transaction.go`
- interface: `internal/baseline/apply.go`

## References

`_techspec.md` → Build Order 5; System Architecture: the Baseline Plan and its
transaction; Data Models: the transaction journal; Risks: moving files the tool
does not own. `_prd.md` → Core Features 3, 4 and 5; Goal 2; User Story 2.
ADR-0121.

## Result

### Implementation

- The existing Baseline transaction now captures each planned history move's
  exact source preimage and ordered identity ledger in its recovery journal.
  It records the move in mutation order before calling `os.Root.Rename`, then
  verifies that the source is absent and the destination has the recorded
  content identity.
- An occupied destination produces a per-file `HistoryMoveRefusal` naming the
  source, destination, and `already exists` reason. The transaction continues
  later moves, commits the successful siblings, and returns both verified and
  refused move evidence.
- `ApplyPlan` includes move paths when inspecting approved versus already-applied
  state. After a collision run commits its unaffected work, apply returns a
  stale-plan refusal stating that not every history relocation was performed and
  names every refused source and destination.
- Rollback processes history moves before ordinary postimages. A normally moved
  file is renamed back; if destination verification detects changed bytes, the
  captured preimage restores the recorded source bytes and removes the invalid
  destination. Interrupted transactions use the same journal mutation order.
- A plan with no history moves omits the new journal fields and follows the
  pre-existing postimage-only path. The relocation implementation uses anchored
  filesystem operations and adds no Git runner or Git command.

### Focused-check evidence

- Acceptance 1: `TestHistoryMoveApply/moves_every_recorded_identity` passed with
  two real filesystem relocations. It compared verified move evidence to the
  approved ledger, read each destination's recorded SHA-256 identity, and proved
  both sources absent.
- Acceptance 2: `TestHistoryMoveCollision` passed. An occupied destination kept
  its original bytes and its source, the refusal named both paths and the
  `already exists` reason, and sibling moves reached their destinations. The
  public `ApplyPlan` case also proved a sibling completed despite the returned
  refusal.
- Acceptance 3: the public `ApplyPlan` collision case passed with an error
  containing `not every history relocation was performed` plus the refused
  source and destination.
- Acceptance 4: `TestHistoryMoveRollback` passed three cases. A failure before a
  sibling move restored the already-moved file; a destination identity mismatch
  restored the recorded source bytes and removed the destination; and abandoning
  a transaction after a rename let the next transaction recover from the
  journal. Every case found all sources restored and all move destinations
  absent.
- Acceptance 5: `TestHistoryMoveApply/empty_ledger_preserves_the_existing_transaction_result`
  passed with nil move evidence and the same verified postimages. Existing
  `TestTransaction*`, `TestApplyExactDigest`,
  `TestApplyPostimageFailureRollsBack`, and `TestEmptyReapply` regressions also
  passed.
- Acceptance 6: inspection of `relocateHistoryMove` and its helpers found only
  `os.Root` state reads, `Rename`, directory sync, journal writes, and identity
  verification. No Git call or runner was added to the relocation path.

Commands run:

- `GOCACHE=/tmp/roundfix-task-06-go-cache rtk go test ./internal/baseline -run '^TestHistoryMoveApply$' -count=1`
  — before implementation, failed to compile because move evidence and
  relocation transaction behavior did not exist; after the first implementation
  pass it exposed and then verified the empty-ledger compatibility correction.
- `GOCACHE=/tmp/roundfix-task-06-go-cache rtk proxy go test ./internal/baseline -run '^TestHistoryMove(Apply|Collision|Rollback)$' -count=1 -v`
  — passed all three groups and their five subtests after the final code and test
  edits.
- `GOCACHE=/tmp/roundfix-task-06-go-cache rtk proxy go test ./internal/baseline -run '^Test(Transaction|ApplyExactDigest|ApplyPostimageFailureRollsBack|EmptyReapply)' -count=1`
  — passed the existing transaction, apply, rollback, and empty-reapply
  regressions.
- `GOCACHE=/tmp/roundfix-task-06-go-cache rtk go vet ./internal/baseline` —
  exited 0 with no diagnostics.
- `rtk git diff --check` — exited 0.
- `GOCACHE=/tmp/roundfix-task-06-go-cache rtk make verify-incremental` — all Go
  packages passed, including `internal/baseline`; the gate then exited 2 in the
  pre-existing `roundfix/skills` contract because
  `.agents/skills/write-prd/SKILL.md` is missing `Exclude
  \`_archived/specs/\` from automatic link rewrites`. That prior-task file is
  outside Task 06 and was not changed.

The Daemon-owned commands under `## Verification` were not run during this turn.
