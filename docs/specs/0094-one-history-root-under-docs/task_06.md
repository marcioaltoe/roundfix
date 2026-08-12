---
task: task_06
spec: 0094-one-history-root-under-docs
status: pending # pending | in_progress | completed | failed — only implement-task changes this
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
