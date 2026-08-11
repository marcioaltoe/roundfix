---
task: task_01
spec: 0092-a-run-that-can-hand-back-its-work
status: completed
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

- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^TestRunDispositionCharacterization' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestRunDispositionCharacterizationWorkStartedPrecedesTheFirstPrompt'` — expected: exits 0. A `-run` pattern selecting no cases exits 0, so this asserts the named case ran.
- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^TestRunDispositionCharacterization' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestRunDispositionCharacterizationFailedBatchOverwritesSettledIssues'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^TestRunDispositionCharacterization' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestRunDispositionCharacterizationStoppedRunLeavesTasksPending'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^TestRunDispositionCharacterization' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestRunDispositionCharacterizationPreflightRefusesOnAnUnintegratedBranch'` — expected: exits 0.
- `test "$(grep -c 'Outcome contract test:' internal/daemon/run_disposition_characterization_test.go)" -eq 6` — expected: exits 0, proving all six are enumerated rather than some.

## References

- `_prd.md` → every Goal.
- `_techspec.md` → Build Order 1; Testing Approach.

## Result

Implemented the current-disposition corpus without changing production code:

- `TestRunDispositionCharacterizationWorkStartedPrecedesTheFirstPrompt` records
  that `agent_work_started` is published before a first prompt that fails
  without Agent output. Task 02 owns the declared break.
- `TestRunDispositionCharacterizationFailedBatchOverwritesSettledIssues`
  records that `MarkBatchFailed` overwrites an already-resolved Review Issue.
  Task 03 owns the declared break.
- `TestRunDispositionCharacterizationStoppedRunLeavesTasksPending` drives a
  Task through Verification and a real settlement commit, then records that a
  Stopped Run leaves the checkout Task `pending` while the Run Worktree Task is
  `completed` in that commit. Task 06 owns the declared break.
- `TestRunDispositionCharacterizationPreflightRefusesOnAnUnintegratedBranch`
  drives the public CLI against a diverged, unintegrated Run Branch and records
  Branch Integrity exit 2, its integration command, and the absence of a new
  Run. Task 05 owns the declared break.
- Six `Outcome contract test:` declarations name the four `internal/cli`
  contracts and two `internal/daemon` contracts that currently derive an
  Unresolved Run from a failed Batch. Each declaration states its assertion and
  names Task 04 as its rewrite owner.

Focused checks:

- `rtk go test ./internal/daemon -run '^TestRunDispositionCharacterization(WorkStartedPrecedesTheFirstPrompt|FailedBatchOverwritesSettledIssues|StoppedRunLeavesTasksPending|PreflightRefusesOnAnUnintegratedBranch)$' -count=1`
  passed all four named cases.
- `rtk rg -c 'Outcome contract test:' internal/daemon/run_disposition_characterization_test.go`
  reported `6`.

Acceptance evidence:

- Each of the four behaviours has a named test asserting today's observable
  outcome; the focused four-case command passed after the last test edit.
- The declaration corpus contains exactly six current tests, with the assertion
  each makes written beside its fully qualified test name.
- Every declared break names its changing Task: Tasks 02, 03, 05, and 06 own
  the four characterization breaks, and Task 04 owns all six outcome-contract
  rewrites.

The authored `## Verification` commands were not run; the Daemon owns them.
