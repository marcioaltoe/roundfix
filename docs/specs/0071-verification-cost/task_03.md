---
task: task_03
spec: 0071-verification-cost
status: pending
type: test
complexity: high
---

# Task 03: Run the CLI tests in parallel

## Overview

With the process-global dependency removed, the CLI package's tests can declare
parallelism — and this is the only slice that can move the headline number,
because packages already overlap and this package is the floor. Parallel
execution also surfaces shared state that sequential execution hid; those are
defects to fix, not reasons to revert.

## Requirements

1. MUST declare parallelism on every test in the package that no longer mutates
   process state and owns whatever filesystem it touches.
2. MUST leave a one-line reason on every test that stays sequential, so
   sequential execution is a recorded decision rather than an omission.
3. MUST fix, not silence, any test that fails only under parallel execution: a
   test that passes alone and fails alongside others has found shared state.
4. MUST prove the absence of races and cross-test leakage by running the
   package with race detection and repeated execution.
5. MUST leave the coverage record unchanged.
6. MUST make the package measurably faster than its recorded 113.2s baseline.

## Subtasks

- [ ] Declare parallelism where process state is no longer touched.
- [ ] State a reason on every test left sequential.
- [ ] Fix shared-state failures surfaced by overlapping execution.
- [ ] Prove no races and no cross-test leakage.

## Acceptance Criteria

- [ ] The package's parallel-declaring test count rises from one to a
      substantial share of its tests.
- [ ] Every test left sequential carries a one-line reason.
- [ ] The package passes with race detection enabled.
- [ ] The package passes when its tests run twice in the same invocation,
      proving no state leaks between overlapping tests.
- [ ] The coverage record from task 01 is unchanged.
- [ ] The package completes measurably faster than 113.2s on the same machine.
- [ ] `git status --porcelain` shows no path outside `internal/cli/` and this
      task file.

## Context

- interface: `internal/cli/cli_test.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/cli -count=1` — expected: exit 0.
- `go test ./internal/cli -count=2` — expected: exit 0; no state leaks between
  overlapping repeated tests.
- `go test ./internal/cli -race -count=1` — expected: exit 0; no data race.
- `go test ./internal/spec -run '^TestCoverageEquivalence$' -count=1 -v | grep -q -- "--- PASS: TestCoverageEquivalence"`
  — expected: exit 0; no test disappeared.
- `grep -rc 't.Parallel()' internal/cli | grep -qv ':1$'` — expected: exit 0;
  the package no longer has a single parallel declaration.

## References

- `_prd.md` → Core Features 1 and 2; Success Metrics (`internal/cli` faster
  than 113.2s).
- `_techspec.md` → Build Order 3; Risks (parallelising surfaces real defects).
- ADR-0089.
