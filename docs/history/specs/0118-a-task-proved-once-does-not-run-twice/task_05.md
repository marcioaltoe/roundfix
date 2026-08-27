---
status: completed
type: docs
---

# Task: Document both command contracts

A reader meets these two behaviors at the command line. The reconcile page
still describes carry-forward as a stopped-Run act, and the implement page
knows nothing about a Preflight that can refuse a Run over work a prior Run
proved.

## Work

- The reconcile contract records which Run outcomes carry-forward accepts and
  that every other terminal outcome is refused by name.
- The implement contract describes the new Preflight refusal: what triggers it,
  that no Run is created, the command that clears it, and that a Run whose
  stranded work is not carriable proceeds with a report instead.
- Say plainly that the caller effectively gets one carry-forward, because a
  carried Task's own file becomes a moved input afterwards. A reader who
  carries from the wrong Run first should learn why before doing it, not after.
- Recut after QA finding F-02. The first cut of this Task documented the
  decision rule as originally specified, and task_08 then changed it. Document
  the delivered rule: Preflight refuses only when the **complete** candidate
  set would carry, and a candidate's declared inputs are compared against the
  accumulating staged carries rather than against the raw checkout.
- Remove the superseded claims: that one passing Task triggers a refusal, and
  that inputs are compared byte-for-byte with the checkout.
- Claims are read from the delivered code, not from the TechSpec draft. Where
  the two disagree, the code is the fact and the TechSpec is corrected.

## References

- `_prd.md` → User Stories 1, 2 and 4; Core Features 2, 3, 4 and 5
- `_techspec.md` → Build Order 5; API Contracts

## Verification
- `grep -q "staged carries" docs/user-guide/commands.md && ! grep -q "byte-for-byte with the checkout" docs/user-guide/commands.md && ! grep -q "passes Task Carry-Forward.s proofs, Preflight Validation refuses" docs/user-guide/commands.md && grep -q "Stopped or Unresolved" docs/user-guide/commands.md && go test -count=1 -tags docscontract ./internal/docscontract`

## Result

Implementation:

- Updated the Implement Command contract to refuse only when a prior Run's
  complete candidate set would carry, while preserving the no-Run,
  no-Agent-Session, no-Git, and no-Run-Database side-effect boundary and the
  exact recovery command.
- Documented report-and-proceed behavior for a stranded Run whose complete
  candidate set is not carriable, inspection failure, released Run Worktrees,
  and largest-set/newest-tie selection.
- Updated the Reconcile Command contract to name Stopped and Unresolved as
  accepted outcomes, refuse every other terminal outcome by its actual name,
  compare declared inputs against accumulating staged carries, and refuse a
  whole set when any candidate refuses.
- Kept the one-effective-carry warning: the carried Task file becomes a moved
  input, so carrying from the wrong Run first can make a later overlapping set
  refuse.
- Corrected the TechSpec API contract to state the delivered complete-set
  Preflight rule.

Focused checks:

- Pre-change `rtk rg -n -F 'passes Task Carry-Forward' docs/user-guide/commands.md`
  exited `0` at the superseded single-candidate rule.
- Pre-change `rtk rg -n -F 'byte-for-byte with the checkout' docs/user-guide/commands.md`
  exited `0`, and `rtk rg -n -F 'staged carries' docs/user-guide/commands.md`
  exited `1`.
- Pre-change `rtk rg -n -F 'at least one carriable' docs/specs/0118-a-task-proved-once-does-not-run-twice/_techspec.md`
  exited `0` at the stale API contract.
- Post-edit checks were run after the final wording edit: the accepted-outcome,
  complete-set, staged-carries, and delivered-code searches exited `0`; the
  superseded checkout-byte, single-candidate, and stale API searches exited
  `1`; `rtk git -c core.fsmonitor=false diff --check` exited `0`.
- `rtk git -c core.fsmonitor=false diff --name-only` reported only this Task
  file, `docs/user-guide/commands.md`, and the Spec TechSpec.
- The Daemon-owned Verification command was not run during this Agent turn.

Acceptance evidence:

1. `docs/user-guide/commands.md` now names Stopped and Unresolved as the
   accepted carry-forward outcomes and states that every other terminal
   outcome is refused by name.
2. The Implement section now states the complete-candidate-set trigger,
   exit `2`, empty stdout, no Run or Agent Session, no Git or Run Database
   writes, and `roundfix reconcile <run-id> --carry-forward` as the next
   action.
3. The Implement section reports non-carriable stranded work and inspection
   failures before proceeding, and documents largest-set/newest-tie selection.
4. The Reconcile section explains accumulating staged carries and the
   one-effective-carry consequence of the carried Task file becoming a moved
   input.
5. The TechSpec API contract now matches the delivered `wouldCarry()` rule;
   the superseded single-candidate and checkout-byte comparison claims were
   removed from the command guide.
