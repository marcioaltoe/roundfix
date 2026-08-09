---
task: task_01
spec: 0092-a-run-that-can-hand-back-its-work
status: pending
type: test
complexity: high
---

# Task 01: Record the four dispositions a Run gets wrong today

## Overview

The corpus every later Task is measured against, and unusually load-bearing here:
one of the changes this Spec makes was already attempted on 2026-08-09 and
reverted within the hour, because six existing tests encode the contract it
breaks. Those six are named here, with the contract they assert, so the Task that
rewrites them does so deliberately rather than by discovery.

## Requirements

1. MUST record that the work-started signal is published before the first prompt,
   so a first turn that fails without Agent output arrives with it already set.
2. MUST record that a failed Batch overwrites Review Issues its Agent had already
   settled Terminal.
3. MUST record that a stopped Run leaves its settled Tasks reading `pending` in
   the checkout while their commits exist in the Run Worktree.
4. MUST record that Branch Integrity Preflight refuses a new Run while an
   unintegrated Run Branch exists for the same head branch.
5. MUST name, in one place, the six tests in `internal/cli` and `internal/daemon`
   that assert a failed Batch produces an Unresolved Run, each with the assertion
   it makes and the Task that will rewrite it.
6. MUST NOT change any production behaviour.

## Subtasks

- [ ] Capture the four behaviours.
- [ ] Enumerate the six outcome-contract tests with their assertions.
- [ ] Declare each break against its Task.

## Acceptance Criteria

- [ ] Each of the four behaviours has a test asserting today's outcome.
- [ ] The six outcome-contract tests are enumerated with their current
      assertions.
- [ ] Each declared break names the Task that changes it.

## Bounded scope

This Task may create or modify only:

- `internal/daemon/run_disposition_characterization_test.go`
- `docs/specs/0092-a-run-that-can-hand-back-its-work/task_01.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^TestRunDispositionCharacterization' -count=1 -v 2>&1 | grep -q '^--- PASS: TestRunDispositionCharacterizationWorkStartedPrecedesTheFirstPrompt'` — expected: exits 0. A `-run` pattern selecting no cases exits 0, so this asserts the named case ran.
- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^TestRunDispositionCharacterization' -count=1 -v 2>&1 | grep -q '^--- PASS: TestRunDispositionCharacterizationFailedBatchOverwritesSettledIssues'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^TestRunDispositionCharacterization' -count=1 -v 2>&1 | grep -q '^--- PASS: TestRunDispositionCharacterizationStoppedRunLeavesTasksPending'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^TestRunDispositionCharacterization' -count=1 -v 2>&1 | grep -q '^--- PASS: TestRunDispositionCharacterizationPreflightRefusesOnAnUnintegratedBranch'` — expected: exits 0.
- `test "$(grep -c 'Outcome contract test:' internal/daemon/run_disposition_characterization_test.go)" -eq 6` — expected: exits 0, proving all six are enumerated rather than some.

## References

- `_prd.md` → every Goal.
- `_techspec.md` → Build Order 1; Testing Approach.
