---
status: completed
type: backend
---

# Task: Preflight refuses to re-execute proved work

This is the Task the Spec exists for. An Implement Run that would re-execute
Tasks a prior Run already completed and verified must stop before it spends an
Agent turn, and must say which command recovers them.

## Work

- Add the check to the implement Preflight block, after the Run Database is
  opened because it needs the store, and beside the Run Window check because
  that is the established precedent for refusing with no Run and no side
  effects.
- Refuse only when at least one Task is carriable. No Run created, no Agent
  Session opened, nothing written to Git or the Run Database.
- The refusal names the Run, each Task it would recover, and the exact
  carry-forward invocation. It implements the same `NextAction` extension point
  the Run Window refusal uses, so it renders through the existing printer.
- Where several prior Runs hold carriable work, name the Run with the largest
  carriable set and break a tie by preferring the most recently created. The
  caller effectively gets one attempt, because a carried Task's own file
  becomes a moved input afterwards and makes an overlapping set refuse.
- When prior Runs hold completed Tasks but nothing is carriable, report what
  was found and why, then proceed to create the Run.
- Fail open. Any failure of the inspection itself is reported on stderr and the
  Run proceeds. This check saves Agent turns; it must never become a new way
  for the loop to stop.
- Cover: a refusal with zero Run rows observed in the Run Database, a proceed
  on non-carriable work, a proceed on inspection failure, and the selection
  rule with both its largest-set and tie-break cases.

## References

- `_prd.md` → Goals 1, 3 and 4; User Stories 2, 3 and 4; Core Features 3, 4
  and 5
- `_techspec.md` → Build Order 3; Interfaces:
  `implementCarryForwardAvailableError`
- ADR-0115 rejects a command that mutates to unblock itself, which is why this
  refuses rather than carrying work forward on its own

## Verification
- `grep -q "implementCarryForwardAvailableError" internal/cli/implement.go && grep -q "TestImplementRefusesWhenCarryForwardIsAvailable" internal/cli/implement_test.go && go test -count=1 ./internal/cli -run 'TestImplementRefusesWhenCarryForwardIsAvailable|TestImplementProceedsWhenNothingIsCarriable' 2>&1 | grep -q "^ok"`

## Result

Implementation:

- Added the Task Carry-Forward inspection immediately after the Run Database
  and Run Window checks. It runs before Agent Selection proof, retention,
  terminal-worktree pruning, Run creation, or task Agent work.
- Added `implementCarryForwardAvailableError`. Its reason names the selected
  Run and every carriable Task, and its `NextAction` supplies the exact
  `roundfix reconcile <run-id> --carry-forward` invocation to the existing
  Preflight printer.
- Added deterministic selection by carriable Task count, then by Run creation
  time, with Task identifiers sorted for stable output.
- Added diagnostic-only paths for prior Runs with no carriable Tasks and for
  inspection errors. Both paths continue into normal Run creation; inspection
  errors therefore fail open.
- Routed the inspection through the existing command dependency seam so the
  failure path and selection rules can be exercised without a production-only
  test hook.

Focused checks:

- Red signal: `GOCACHE=/tmp/roundfix-task-03-go-cache rtk go test
  ./internal/cli -run '^TestImplementRefusesWhenCarryForwardIsAvailable$'
  -count=1` reached all three refusal cases and failed because terminal
  worktree pruning ran, proving the carry-forward refusal did not yet exist.
- `GOCACHE=/tmp/roundfix-task-03-go-cache rtk go test ./internal/cli -run
  '^TestImplementRefusesWhenCarryForwardIsAvailable$' -count=1` exited 0 with
  five tests, including a real Run Database and Git carry-forward inspection.
- `GOCACHE=/tmp/roundfix-task-03-go-cache rtk go test ./internal/cli -run
  '^TestImplementProceedsWhenNothingIsCarriable$' -count=1` exited 0 with three
  tests covering a proof refusal and an inspection error.
- `GOCACHE=/tmp/roundfix-task-03-go-cache rtk go test ./internal/cli -run
  '^(TestInspectSpecCarryForwards|TestImplementRunWindow|TestImplementRefusesWhenCarryForwardIsAvailable|TestImplementProceedsWhenNothingIsCarriable)$'
  -count=1` exited 0 with 17 tests.
- `rtk git diff --check` exited 0.

Acceptance evidence:

1. The command calls `inspectSpecCarryForwards` after `store.Open` and the Run
   Window check, and before Agent Selection proof or any mutating preflight
   work. The real-prior-Run subtest observes no Agent calls or proof Sessions.
2. The refusal subtests observe Preflight exit, empty stdout, unchanged HEAD
   and Git status, zero task Agent turns, zero Agent Selection proof Sessions,
   and no additional Run row. The injected refusal cases observe zero total
   Run rows.
3. Every refusal assertion observes the selected Run, all selected Task ids,
   the existing no-side-effects block, and the exact backticked carry-forward
   invocation rendered through `NextAction`.
4. `largest_carriable_set_wins_over_recency` selects a two-Task older Run over
   a one-Task newer Run; `newest_Run_wins_a_carriable-count_tie` selects the
   later `CreatedAt` value.
5. `completed_Tasks_fail_their_carry-forward_proofs` observes the Run and its
   moved-input refusal on stderr, then observes one task Agent turn and one new
   Run row.
6. `inspection_failure_fails_open` observes the injected inspection error and
   the proceed notice on stderr, then observes one task Agent turn and one new
   Run row.
7. The focused matrix covers the real refusal wiring, a zero-row refusal,
   non-carriable proceed, inspection-error proceed, both selection rules, and
   the adjacent Run Window and Spec inspection behavior.

The Daemon-owned `## Verification` command was not run in this Agent turn.
