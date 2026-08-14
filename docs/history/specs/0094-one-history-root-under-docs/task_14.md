---
task: task_14
spec: 0094-one-history-root-under-docs
status: completed # pending | in_progress | completed | failed — only implement-task changes this
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

## Result

### Implementation

- Occupied-destination collisions now enter the ordered `historyMoves` ledger
  with their source identity. The existing transaction therefore refuses each
  occupied move after continuing through its siblings, and the public apply
  path returns the existing actionable `stale` error with both paths and the
  `already exists` reason.
- `baseline update` now treats either managed file changes or planned History
  Relocations as pending work. Its text and JSON results carry the outstanding
  `historyMoves`, so a follow-up run names the unresolved source and destination
  and cannot report `current` while the dual layout remains.
- The no-collision path still applies every History Relocation and reports the
  same verified result. The production change reuses the existing transaction
  refusal and verification behavior rather than adding another mutation path.

### Focused-check evidence

- Red signal before the production edit:
  `rtk go test -count=1 ./internal/baseline -run '^TestHistoryMoveCollisionRefusesPublicly$'`
  exited 1 because the generated Plan omitted the occupied relocation and held
  only its sibling.
- Red signal before the production edit:
  `rtk go test -count=1 ./internal/baseline -run '^TestUnresolvedLayoutIsNotCurrent$'`
  exited 1 because `ApplyPlan` returned nil for the collision-only tree.
- Red signal before the production edit:
  `rtk go test -count=1 ./internal/cli -run '^TestBaselineApplyCollisionExitStatus$'`
  exited 1 because public apply exited 0 and reported `verified`.
- After implementation,
  `rtk go test -count=1 ./internal/baseline -run '^TestHistoryMoveCollisionRefusesPublicly$|^TestUnresolvedLayoutIsNotCurrent$'`
  exited 0 with four passing tests/subtests.
- After implementation,
  `rtk go test -count=1 ./internal/cli -run '^TestBaselineApplyCollisionExitStatus$'`
  exited 0.
- The adjacent baseline check
  `rtk go test -count=1 ./internal/baseline -run 'TestHistoryMove|TestUnresolvedLayout|TestDiscoverHistoryLayoutReportsCollision'`
  exited 0 with 15 passing tests/subtests.
- The adjacent CLI check
  `rtk go test -count=1 ./internal/cli -run 'TestBaselineApply|TestBaselineUpdate(IdempotenceReportsZeroFileChanges|TextReportsHistoryMoves)'`
  exited 0 with 13 passing tests/subtests.
- `rtk make verify-incremental` exited 0. All Go packages passed, including
  `internal/baseline` and `internal/cli`; the focused skill contract,
  `roundfix skills check`, and the production build also passed.
- `rtk git diff --check` exited 0 before this Result was recorded.

### Acceptance evidence

1. `TestHistoryMoveCollisionRefusesPublicly` and
   `TestBaselineApplyCollisionExitStatus` require the refusal output to contain
   the source, destination, `already exists`, and the statement that not every
   History Relocation was performed.
2. Both tests inspect the real filesystem after refusal and prove the sibling
   source is absent and its destination contains the original bytes, while both
   colliding files remain unchanged.
3. `TestBaselineApplyCollisionExitStatus` proves the public command exits 3
   (`exitUnverified`) rather than 0.
4. The same CLI test performs a second public `baseline update` and requires
   `plan ready`, `History moves: 1`, and both outstanding paths while forbidding
   `Baseline update: current`. `TestUnresolvedLayoutIsNotCurrent` independently
   proves the follow-up Plan retains the outstanding identity ledger entry when
   managed file changes are zero.
5. The table-driven no-collision companion fixes the same source tree while
   varying only destination occupancy. It proves both files relocate, both moves
   are verified, and apply returns no error. The existing current-layout
   idempotence test also passed in the adjacent CLI check.

### Not run

- The commands under this Task's `## Verification` section were not run; the
  Roundfix Daemon owns declared Verification and Task settlement.
