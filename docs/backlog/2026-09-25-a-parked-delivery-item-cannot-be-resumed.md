---
type: fix
status: open
created: 2026-09-25
spec: null
reason: null
---

# A parked Delivery Queue item can never re-enter the queue

## Symptom

After the operator fixes what parked an item (here `run-unresolved` for Spec 0170), no command returns it to the queue: the owner loop skips every `parked` item, `deliver resume` only restarts the owner, and `deliver start` records a new queue. The operator has to leave `deliver` and finish the Spec by hand (implement, review, archive, PR, merge), which is the orchestration cost `deliver` exists to remove.

## Where

`internal/delivery/engine.go` (`Run` skips `DeliveryStageParked`); `internal/cli/deliver.go` (no unpark command).

## Expected

A command such as `roundfix deliver retry <slug>` (or `resume` re-evaluating parked items whose blocker is resolved) returns the item to the stage it parked at, re-running the Run only for unfinished Tasks.

## Evidence

Delivery queue of 2026-09-25, item 0170 parked `run-unresolved` after task_02 failed on a corpus golden count.
