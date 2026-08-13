---
task: task_17
spec: 0094-one-history-root-under-docs
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: low
---

# Task 17: Leave no empty shell where a relocation started

## Overview

The QA gate found that relocating a finished orphan Review Artifact moves its
files and leaves its directories behind. Discovery then reads the empty shell as
a Review Artifact whose metadata is gone, and reports it as undecidable on every
later run — so a review that was migrated correctly keeps being reconsidered
forever.

## Requirements

1. MUST remove a source directory the relocation emptied, up to but never past
   the tree the relocation started from.
2. MUST NOT remove a source directory that still holds anything, whether a file
   the relocation did not move or a directory it did not empty.
3. MUST leave a rerun after a completed relocation reporting nothing about the
   relocated Review Artifact, since there is no longer anything there to classify.
4. MUST keep the relocation's existing guarantees unchanged: the bytes move, the
   destination identity is verified, and a collision still refuses by name.
5. MUST leave rollback able to restore a relocated file to its source, including
   when the source's directory was removed.

## Subtasks

- [ ] Remove source directories the relocation emptied.
- [ ] Stop short of any directory that still holds something.
- [ ] Cover the rerun-reports-nothing case and the rollback-after-removal case.

## Acceptance Criteria

- [ ] After relocating every file under a review directory, neither that
      directory nor its rounds remain.
- [ ] A source directory holding an unmoved file survives with that file intact.
- [ ] A rerun after a completed relocation reports nothing about the relocated
      Review Artifact.
- [ ] Rollback restores a relocated file to its source path even though the
      source directory was removed.
- [ ] Collision refusal and destination identity verification are unchanged.

## Rehearsal Cases

- Case: a review directory whose every file relocates; Observation: the review
  and round directories are gone, and a rerun reports nothing about it.
- Case: a source directory holding one file the relocation did not move;
  Observation: the directory and that file survive.
- Case: a rollback after a relocation whose source directory was removed;
  Observation: the file is back at its source path.
- Case: a plan with a collision; Observation: the refusal names both paths as
  before, and the refused source directory is untouched.

## Verification

- `go test -count=1 ./internal/baseline -run 'TestHistoryMoveRemovesEmptiedSource' -v > /tmp/0094-task-17.log 2>&1; s=$?; grep -q '^--- PASS: TestHistoryMoveRemovesEmptiedSource' /tmp/0094-task-17.log || { cat /tmp/0094-task-17.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing test; fails today, where no such test exists.
- `! grep -qi 'no tests to run' /tmp/0094-task-17.log` — expected: exits 0, refusing a vacuous run.
- `go test -count=1 ./internal/baseline -run 'TestHistoryMoveRollback|TestHistoryMoveCollision' -v > /tmp/0094-task-17b.log 2>&1; s=$?; grep -q '^--- PASS: .*HistoryMoveRollback' /tmp/0094-task-17b.log && grep -q '^--- PASS: .*HistoryMoveCollision' /tmp/0094-task-17b.log && grep -rq 'RemovesEmptiedSource' internal/baseline || { cat /tmp/0094-task-17b.log; exit 1; }; exit $s` — expected: exits 0, proving rollback and collision still hold once the removal exists. The suite alone passes before any work, so it is anchored to the change it guards; an unanchored guard is refused as vacuous, which is what happened to an earlier draft of this command.

## Context

- interface: `internal/baseline/history_layout.go`
- interface: `internal/baseline/transaction.go`

## References

`_prd.md` → Core Features 3, 6 and 7; Goal 4, a finished review stops
accumulating beside live work. QA report `qa/qa-report-2026-08-13.md` → F-002,
which names `relocateHistoryMove` as the root cause. ADR-0123.
