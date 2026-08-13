---
task: task_14
spec: 0094-one-history-root-under-docs
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: high
---

# Task 14: Refuse a collision and stop calling the tree current

## Overview

The QA gate found the collision path silent and, worse, self-certifying: an
occupied destination let its sibling move, left both colliding files in place,
exited zero, and the next run reported the repository `current` with no file
changes. A repository with a pending migration declares itself finished. Core
Feature 7 requires a refusal that names both paths; this slice makes the refusal
real and stops the follow-up run from asserting something false.

## Requirements

1. MUST refuse the relocation whose destination already exists, naming the source,
   the destination, and the reason.
2. MUST still perform the relocations of the other files in the same plan, so one
   collision does not refuse its siblings.
3. MUST report, on a run where at least one relocation was refused, that not every
   relocation was performed, and MUST NOT exit as though the run succeeded.
4. MUST NOT report a repository as current while a discoverable relocation remains
   unperformed; an unresolved dual layout is a reported state, never a finished
   one.
5. MUST leave a repository with no collision applying and reporting exactly as it
   does today.

## Subtasks

- [ ] Surface the per-file collision refusal through the public apply path.
- [ ] Make the run's exit reflect that not every relocation was performed.
- [ ] Stop a follow-up run from reporting current while relocations remain.
- [ ] Cover collision, sibling completion, exit status, and the follow-up run.

## Acceptance Criteria

- [ ] An occupied destination refuses that relocation and names both paths and the
      reason in the output.
- [ ] The siblings of a refused relocation still reach their destinations.
- [ ] A run with a refused relocation does not exit as a success.
- [ ] A run following an unresolved collision reports the repository as not
      current and names what remains.
- [ ] A repository with no collision is unaffected, proven by a test that fixes
      the tree and varies only the collision.

## Rehearsal Cases

- Case: a destination that already holds a different file; Observation: that
  relocation is refused naming both paths and the reason, and its siblings
  completed.
- Case: the exit status of a run carrying a refused relocation; Observation: it
  does not report success.
- Case: a second run against the same unresolved tree; Observation: it reports
  the repository as not current and names the outstanding relocation.
- Case: a tree with no collision; Observation: apply and the follow-up run behave
  exactly as before this slice.

## Verification

- `go test -count=1 ./internal/baseline -run 'TestHistoryMoveCollisionRefusesPublicly|TestUnresolvedLayoutIsNotCurrent' -v > /tmp/0094-task-14.log 2>&1; s=$?; grep -q '^--- PASS: TestHistoryMoveCollisionRefusesPublicly' /tmp/0094-task-14.log && grep -q '^--- PASS: TestUnresolvedLayoutIsNotCurrent' /tmp/0094-task-14.log || { cat /tmp/0094-task-14.log; exit 1; }; exit $s` — expected: exits 0 and the log names both passing tests; fails today, where neither exists.
- `! grep -qi 'no tests to run' /tmp/0094-task-14.log` — expected: exits 0, refusing a vacuous run.
- `go test -count=1 ./internal/cli -run 'TestBaselineApplyCollisionExitStatus' -v > /tmp/0094-task-14b.log 2>&1; s=$?; grep -q '^--- PASS: TestBaselineApplyCollisionExitStatus' /tmp/0094-task-14b.log || { cat /tmp/0094-task-14b.log; exit 1; }; exit $s` — expected: exits 0, proving the refusal reaches the command's exit status rather than staying inside the package.

## Context

- interface: `internal/baseline/apply.go`
- interface: `internal/baseline/transaction.go`
- interface: `internal/baseline/history_layout.go`

## References

`_prd.md` → Core Feature 7; User Experience, the refused-relocation sentence.
`_techspec.md` → API Contracts, the collision paragraph; Risks: a repository
holding both layouts. QA report `qa/qa-report-2026-08-13.md` → F-002, the
self-certifying half.
