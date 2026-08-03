---
task: task_08
spec: 0071-verification-cost
status: pending
type: test
complexity: high
---

# Task 08: Finish declaring parallelism in the CLI package

## Overview

The prefactor removed the process-global blockers, but only 207 of the
package's 488 tests declare parallelism. Of the 281 still sequential, just nine
retain a real blocker and none changes the working directory — so roughly 272
tests run one at a time for no remaining reason. Measurement shows this is the
cost: the eight heaviest end-to-end journeys complete in 31.7s together, while
the whole package takes 143.4s.

## Requirements

1. MUST declare parallelism on every test in the package that no longer mutates
   process state and owns whatever filesystem it touches.
2. MUST leave a one-line reason on each test that stays sequential, and the
   only acceptable reason is a process-global dependency the test's own subject
   requires.
3. MUST fix, not silence, any test that fails only under parallel execution.
4. MUST prove the absence of races and cross-test leakage with race detection
   and repeated execution.
5. MUST leave the coverage record unchanged.
6. MUST reduce the package's wall clock materially below the recorded 143.4s.

## Subtasks

- [ ] Declare parallelism on every test with no remaining blocker.
- [ ] State a reason on each test left sequential.
- [ ] Fix shared-state failures the wider overlap surfaces.
- [ ] Confirm no races and no leakage.

## Acceptance Criteria

- [ ] The package's parallel-declaring test count rises substantially above 207.
- [ ] Every test left sequential carries a one-line reason naming its
      process-global dependency.
- [ ] The package passes with race detection enabled.
- [ ] The package passes with its tests run twice in one invocation.
- [ ] The coverage record from task 01 is unchanged.
- [ ] The package completes measurably faster than 143.4s on the same machine.
- [ ] `git status --porcelain` shows no path outside `internal/cli/` and this
      task file.

## Context

- interface: `internal/cli/cli_test.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/cli -count=1` — expected: exit 0.
- `go test ./internal/cli -count=2` — expected: exit 0; no state leaks between
  overlapping repeated tests.
- `go test ./internal/cli -race -count=1` — expected: exit 0.
- `go test ./internal/spec -run '^TestCoverageEquivalence$' -count=1 -v | grep -q -- "--- PASS: TestCoverageEquivalence"`
  — expected: exit 0.
- `go test ./internal/cli -count=1 -v | grep -c -- "--- PASS:"` — expected: a
  count consistent with the recorded coverage set; no test silently dropped.

## References

- `_prd.md` → Core Features 1 and 2; Success Metrics.
- `_techspec.md` → Build Order 3; Risks.
- `baseline/2026-08-03-after.md` → the measurement that identified this cost.
- ADR-0089.
