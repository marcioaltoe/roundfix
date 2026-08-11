---
task: task_04
spec: 0092-a-run-that-can-hand-back-its-work
status: pending
type: backend
complexity: high
---

# Task 04: Derive the Run outcome from unresolved work

## Overview

The Task that makes Task 03 safe. "Did this Batch finish?" and "is there
unresolved work left?" are different questions sharing one answer today, which is
what forces a failed Batch to overwrite the outcomes it recorded. This Task
computes the Run outcome from the Review Issues that remain unresolved, and
rewrites the six tests Task 01 enumerated to the new contract.

## Requirements

1. MUST compute a Run's outcome from the count of Review Issues still unresolved,
   not from whether a Batch reported failure.
2. MUST keep a Run reporting `Unresolved` with exit 1 while any issue remains
   unresolved, including after a Batch failure.
3. MUST let a Run reach `Clean` only when no unresolved issue remains, regardless
   of whether a Batch failed on the way.
4. MUST rewrite each of the six outcome-contract tests Task 01 enumerated,
   stating in each what the new contract asserts and why it is better rather
   than merely different.
5. MUST NOT report `Clean` for a Run whose Agent failed while unresolved work
   remains; that state is what a reverted 2026-08-09 attempt produced.

## Subtasks

- [ ] Compute the outcome from unresolved issues.
- [ ] Rewrite the six outcome-contract tests.
- [ ] Prove a crashed Agent with remaining work still ends Unresolved.

## Acceptance Criteria

- [ ] A Batch failure with issues still unresolved ends the Run `Unresolved`.
- [ ] A Batch failure with every issue resolved ends the Run `Clean`.
- [ ] All six rewritten tests pass and each states its new contract.
- [ ] No path reports `Clean` while an unresolved issue exists.

## Rehearsal Cases

- Case: a Batch of two issues where the Agent resolves one and then crashes;
  Observation: the resolved issue keeps `resolved`, the other is `failed`, and
  the Run ends `Unresolved` with exit 1.
- Case: a Batch of one issue the Agent resolves before the runtime crashes on
  teardown; Observation: the issue keeps `resolved` and the Run ends `Clean`.
- Case: a Batch that fails before resolving anything; Observation: unchanged from
  today — every issue `failed`, Run `Unresolved`.

## Bounded scope

This Task may create or modify only:

- `internal/daemon/engine.go`
- `internal/daemon/engine_test.go`
- `internal/rounds/rounds.go`
- `internal/rounds/rounds_test.go`
- `internal/cli/cli_test.go`
- `docs/specs/0092-a-run-that-can-hand-back-its-work/task_04.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^TestRunOutcomeDerived' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestRunOutcomeDerivedFromUnresolvedIssues'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^TestRunOutcomeDerived' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestRunOutcomeDerivedStaysUnresolvedAfterAFailedBatch'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/cli ./internal/daemon -count=1 2>&1 | tee /dev/stderr | grep -c '^ok' | grep -q '^2$'` — expected: exits 0, proving both packages pass rather than one being skipped. The `tee /dev/stderr` is load-bearing: without it the counting pipeline swallows every failing test name, and the Daemon records an empty diagnostic that names no failure to repair.
- `test "$(grep -c 'New contract:' internal/cli/cli_test.go internal/daemon/engine_test.go | awk -F: '{s+=$2} END {print s}')" -ge 6` — expected: exits 0, proving all six rewritten tests state their new contract.

## References

- `_prd.md` → Goals 2 and 5.
- `_techspec.md` → Build Order 4; Risks.
- ADR-0010, ADR-0113.
