---
task: task_03
spec: 0172-a-qa-gate-that-tells-the-truth
status: pending
type: backend
complexity: medium
---

# Task 03: A temporary retry keeps the deterministic failure

## Overview

In independent Verification, when command A fails deterministically on the first run, command B fails temporarily, and the exclusive retry then ends temporarily on A, `retainCollectedVerificationFailures` in `internal/daemon/task_engine.go` drops A's first-run failure because the retry reached A, and replaces it with the retry's temporary failure. The Task still settles failed, but its reason names only a temporary failure, so the operator reading the Task file and the Run Event re-runs instead of fixing a defect the first run observed.

## Requirements

1. MUST make `retainCollectedVerificationFailures` treat a command the retry ended temporarily like one it could not observe: its first-run deterministic failure is kept, and the retry's temporary failure of that command does not replace it.
2. MUST keep the retry's `TemporaryFailure` on the merged outcome, so `executeTask` still settles the Task failed without a repair turn, as ADR-0056 requires.
3. MUST make the Task's failure reason, built by `taskVerificationFailureReason`, name A's first-run failure and its first-run diagnostic path.
4. MUST keep a retry verdict replacing a first-run verdict when the retry passed or failed deterministically, and keep `TestUnknownOnRetryKeepsTheFirstRunFailure`, `TestRetryKeepsFirstRunFailuresForCommandsItDidNotReach` and `TestIndependentVerificationReplacesCollectedFailuresForRetriedCommands` green unchanged.
5. MUST put the new tests in `internal/daemon/verification_retry_test.go`, keeping each Task's tests in a file of its own.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] First run A deterministic plus B temporary, retry A temporary: the merged outcome's repair target is A's first-run failure and the Task's reason names A's first-run diagnostic path.
- [ ] Through a Task cycle the same sequence settles the Task failed with that reason and gives the Agent no repair turn.
- [ ] A command that passes on the retry drops its first-run failure.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/daemon/engine.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestTemporaryOnRetryKeepsTheFirstRunDeterministicFailure|TestPassOnRetryReplacesTheFirstRunFailure|TestIndependentVerificationTemporaryRetryNamesTheDeterministicFailure|TestUnknownOnRetryKeepsTheFirstRunFailure|TestRetryKeepsFirstRunFailuresForCommandsItDidNotReach|TestIndependentVerificationReplacesCollectedFailuresForRetriedCommands)$" ./internal/daemon 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestTemporaryOnRetryKeepsTheFirstRunDeterministicFailure TestPassOnRetryReplacesTheFirstRunFailure TestIndependentVerificationTemporaryRetryNamesTheDeterministicFailure TestUnknownOnRetryKeepsTheFirstRunFailure TestRetryKeepsFirstRunFailuresForCommandsItDidNotReach TestIndependentVerificationReplacesCollectedFailuresForRetriedCommands; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task the three new named cases do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — The temporary retry
