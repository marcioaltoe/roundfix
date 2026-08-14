---
task: task_17
spec: 0094-one-history-root-under-docs
status: completed # pending | in_progress | completed | failed — only implement-task changes this
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

## Result

Implementation:

- A performed History Relocation now removes empty source directories from the
  moved file's parent through the source tree where discovery began. It stops at
  the first directory that contains an entry and never removes the source tree's
  parent.
- Rollback recreates source parents removed by relocation without recording them
  as apply-created destination directories, so the later rollback cleanup leaves
  the restored source and its directories in place.
- The transaction regression suite exercises a finished orphan Review Artifact,
  a source directory with an unmoved file, and rollback after source pruning.

Focused checks:

- Before the implementation, `GOCACHE=/tmp/roundfix-task-17-go-cache rtk go
  test ./internal/baseline -run '^TestHistoryMoveRemovesEmptiedSource$'
  -count=1` failed because the finished Review Artifact directory still existed.
- After the implementation, `GOCACHE=/tmp/roundfix-task-17-go-cache rtk go
  test ./internal/baseline -run
  '^(TestHistoryMoveRemovesEmptiedSource|TestHistoryMoveRollback|TestHistoryMoveCollision)$'
  -count=1` exited 0 with 9 passing tests.
- `rtk git diff --check` exited 0 with no diagnostics.
- The daemon-owned commands under `## Verification` were not run.

Acceptance evidence:

- `TestHistoryMoveRemovesEmptiedSource/finished_Review_Artifact_leaves_no_source_shell_or_rerun_finding`
  verifies that the Review Artifact and Round directories are absent while the
  live Review Artifact root survives.
- `TestHistoryMoveRemovesEmptiedSource/unmoved_file_keeps_its_source_directory`
  verifies that a directory containing an unmoved file survives with the file's
  original bytes.
- The finished-Review subtest calls `planHistoryMoves` after apply and verifies
  that the rerun returns neither moves nor retained-review findings.
- `TestHistoryMoveRollback/failure_after_source_pruning_recreates_source_directories`
  injects a post-pruning failure and verifies the source bytes and directory are
  restored.
- The focused run includes the existing collision refusal, destination identity
  mismatch, and interrupted-recovery cases; all passed without weakening their
  path, byte-identity, or refusal assertions.
