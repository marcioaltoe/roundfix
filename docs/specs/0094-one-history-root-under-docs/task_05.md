---
task: task_05
spec: 0094-one-history-root-under-docs
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: medium
---

# Task 05: Plan a history relocation as a ledger of identities

## Overview

Planning learns the relocation ledger: an ordered list of source, destination,
and content identity, carrying no bytes, so a plan's size tracks the change it
makes rather than the history it moves. The plan reports relocations alongside
the file changes it already reports, and a repository needing none serializes
exactly as it does today.

## Requirements

1. MUST carry each relocation as source, destination, and content identity, and
   MUST NOT carry file content in the ledger.
2. MUST order the ledger deterministically, so planning one tree twice produces
   an identical document.
3. MUST omit the ledger entirely from a plan with no relocation, so an unrelated
   plan serializes byte-identically to before this Task.
4. MUST report relocations in the plan's human-readable output at the same level
   of detail the plan already gives the changes it makes.
5. MUST expose relocations in the machine-readable output under their own key.
6. MUST keep the existing change ledger and the preimage and postimage sets
   describing only rendered carriers; a relocation never appears there.
7. MUST treat a repository needing migration as a plan with changes, not as an
   error, leaving the command's exit status unchanged.

## Subtasks

- [ ] Add the relocation ledger to the plan document and its serialization.
- [ ] Populate it from discovery during planning.
- [ ] Report relocations in the human-readable and machine-readable outputs.
- [ ] Cover the empty case, the populated case, and determinism with tests.

## Acceptance Criteria

- [ ] A plan for a repository on an older layout lists every relocation with its
      content identity, and no file content appears in the document.
- [ ] A plan for a repository on the current layout omits the ledger key
      entirely.
- [ ] Planning the same tree twice yields an identical document.
- [ ] The human-readable output names the relocations.
- [ ] The existing change ledger and preimage and postimage sets contain no
      relocation entry.
- [ ] Planning a repository that needs migration exits with the same status as
      planning one that does not.

## Verification

- `go test -count=1 ./internal/baseline -run 'HistoryMove|HistoryRelocationPlan' -v > /tmp/0094-task-05.log 2>&1; s=$?; grep -q '^--- PASS: .*HistoryMove\|^--- PASS: .*HistoryRelocationPlan' /tmp/0094-task-05.log || { cat /tmp/0094-task-05.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing ledger tests.
- `grep -q 'historyMoves' internal/baseline/plan.go` — expected: exits 0, proving the ledger reached the plan document rather than only its tests.
- `grep -q 'historyMoves' internal/baseline/plan.go && ! grep -A6 'HistoryMove struct' internal/baseline/plan.go | grep -qi 'content'` — expected: exits 0, proving the ledger exists and its entry carries no content field. The negative clause alone passed before any work, because a token absent from the file cannot match.

## References

`_techspec.md` → Build Order 4; Interfaces: `HistoryMove`; Data Models; API
Contracts: `roundfix baseline plan`. `_prd.md` → Core Features 6 and 8.
ADR-0121.
