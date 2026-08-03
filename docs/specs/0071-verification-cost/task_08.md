---
task: task_08
spec: 0071-verification-cost
status: completed
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

## Result

### Implementation

- Replaced mutable package-level CLI test seams with invocation-scoped
  `commandDependencies` carried by `commandEnvironment` and `context.Context`.
  Parallel tests now receive isolated collaborators, clocks, inputs, runners,
  worktree operations, and command-specific dependencies without mutating
  process state.
- Declared `t.Parallel()` on every runnable top-level test without a remaining
  process-global dependency. The focused AST inventory reports 479 parallel
  tests and eight sequential tests out of 487 runnable tests, up from 207
  parallel declarations.
- Kept only the eight tests that require process-wide `PATH`,
  `ROUNDFIX_TUI`, `ROUNDFIX_COLOR`, or `SIGTERM` handling sequential, with the
  exact dependency named in a one-line comment on each test.
- Fixed failures surfaced during overlap testing: workspace registration now
  merges with an existing per-test dependency override instead of replacing
  it, and the live Attach test now gates its writer on entry to follow mode
  instead of racing command startup.

### Acceptance evidence

- Parallel count: the AST inventory reported
  `parallel=479 sequential=8 runnable=487 parallel_global_assignments=0`.
- Sequential reasons: the same inventory enumerated exactly eight tests; every
  entry names one of `PATH`, `ROUNDFIX_TUI`, `ROUNDFIX_COLOR`, or process-wide
  `SIGTERM` handling.
- Focused race and repeat checks:
  - the non-Implement command group passed with `-race -count=2 -parallel=12`
    in 14.175s;
  - the Implement/Settle group passed with `-race -count=2 -parallel=12` in
    16.250s;
  - the Baseline/profile group passed with `-race -count=2 -parallel=12` in
    64.141s;
  - the Stop/Attach group passed with `-race -count=2 -parallel=12` in 4.563s;
  - the owner-identity and live Attach regressions each passed ten repetitions
    with race detection.
- Compile-only package check passed with `-run '^$' -count=1`.
- Coverage record: `git diff --exit-code --
  docs/specs/0071-verification-cost/coverage-record.json` exited 0.
- Diff hygiene: `git diff --check` exited 0. This task authored changes only
  under `internal/cli/` and this task file. The raw worktree status also contains
  pre-existing Task 06 changes to `Makefile` and `task_06.md`; they were not
  modified by this task.
- The full-package race, repeated-run, coverage-equivalence, pass-count, and
  143.4s wall-clock comparisons remain for the Daemon-owned `## Verification`
  commands, which were intentionally not run during this child-agent turn.
