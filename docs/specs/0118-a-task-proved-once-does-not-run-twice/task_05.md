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
- Claims are read from the delivered code, not from the TechSpec draft. Where
  the two disagree, the code is the fact and the TechSpec is corrected.

## References

- `_prd.md` → User Stories 1, 2 and 4; Core Features 2, 3, 4 and 5
- `_techspec.md` → Build Order 5; API Contracts

## Verification
- `grep -q "Stopped or Unresolved" docs/user-guide/commands.md && grep -q -- "--carry-forward" docs/user-guide/commands.md && go test -count=1 -tags docscontract ./internal/docscontract 2>&1 | grep -q "^ok"`

## Result

Implementation:

- Updated the Reconcile Command contract to state that carry-forward accepts
  Stopped or Unresolved Runs, refuses every other terminal outcome by name,
  and preserves the existing proof and whole-set refusal behavior.
- Updated the Implement Command contract to describe the pre-Run refusal,
  its no-Run/no-side-effects boundary, the exact carry-forward recovery
  command, largest-set selection, and the report-and-proceed paths for
  non-carriable work and inspection failure.
- Added the warning that a carried Task's own file becomes a moved input, so a
  caller effectively gets one carry-forward for overlapping work and must
  choose the source Run carefully.
- Corrected the TechSpec query contract to match the delivered implementation:
  it inspects only Stopped or Unresolved Runs for the Spec whose recorded Run
  Worktree is present and skips released Runs.

Focused checks:

- Pre-change signal: `rtk rg -n -F 'Stopped or Unresolved' docs/user-guide/commands.md`
  exited `1`; the guide had no accepted-outcome wording. The equivalent
  inspection-failure search also exited `1`.
- `rtk rg -n -F 'Stopped or Unresolved' docs/user-guide/commands.md` exited
  `0` and found the reconcile switch and contract at lines 771 and 776.
- `rtk rg -n -F 'effectively gets one carry-forward' docs/user-guide/commands.md`
  exited `0`; the guide explains the moved-input consequence and wrong-Run
  ordering risk.
- `rtk rg -n -F 'per prior Stopped or' docs/specs/0118-a-task-proved-once-does-not-run-twice/_techspec.md`
  exited `0`; the corrected query scope is present in the TechSpec.
- `rtk rg -n 'carryForwardAcceptedStates|Task Carry-Forward inspection failed|reportImplementNonCarriableCarryForwards' internal/cli/carryforward.go internal/cli/implement.go`
  exited `0` and matched the delivered accepted-state, fail-open, and
  non-carriable reporting code.
- `rtk git diff --check` exited `0`.

Acceptance evidence:

1. The reconcile section names Stopped and Unresolved as the accepted
   outcomes, says every other terminal outcome is refused by name, and keeps
   `--carry-forward` alongside the other named mutation switches.
2. The implement section describes the refusal trigger, exit `2`, empty
   stdout, no Run or Agent Session, no Git or Run Database state, and the
   copyable `roundfix reconcile <run-id> --carry-forward` recovery command.
3. The implement section documents report-and-proceed behavior for inspected
   non-carriable work and inspection failure, plus largest-set/newest tie
   selection.
4. The guide explains that carry-forward records provenance in the Task file,
   making that file a moved input afterwards and limiting overlapping work to
   one effective carry-forward.
5. The TechSpec wording now agrees with the delivered query's accepted-state
   filter and released-Run skip behavior.

The Daemon-owned Verification command was not run during this Agent turn.

Final post-edit checks:

- `rtk rg -n -U 'Every other terminal\s+outcome is refused by name' docs/user-guide/commands.md`
  exited `0` and found the outcome-refusal contract.
- `rtk rg -n -F 'roundfix reconcile <run-id> --carry-forward' docs/user-guide/commands.md`
  exited `0` and found both recovery-command examples.
- `rtk rg -n -F 'released Run whose Run Worktree is absent is skipped' docs/specs/0118-a-task-proved-once-does-not-run-twice/_techspec.md`
  exited `0` and found the corrected query behavior.
- `rtk git diff --check` exited `0`.
