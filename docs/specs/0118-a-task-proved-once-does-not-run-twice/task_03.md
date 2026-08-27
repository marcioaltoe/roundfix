---
status: pending
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
