---
task: task_05
spec: 0074-git-spawn-economy
status: pending
type: backend
complexity: high
---

# Task 05: Give the agent runner its environment explicitly

## Overview

The ACPX child's environment is composed from the process at spawn time,
which forces the package's tests to mutate process env with `t.Setenv` —
126 calls that keep 121 of 136 tests sequential, because Go refuses
`t.Parallel()` beside process mutation. This Task applies ADR-0089 exactly
as Spec 0071 applied it to `internal/cli`: the environment composition takes
an explicit base, the process default resolves once at the boundary, tests
inject per-call, and the sequential wall falls.

This is the prefactor-then-parallelize shape, and its rule carries over: a
test that fails under parallel execution found shared state — fixing it is
the work; reverting to sequential without a stated one-line reason is how
the defect survives.

## Requirements

1. MUST make the ACPX environment composition take an explicit base
   environment, with the process environment resolved once at the command
   boundary as the default.
2. MUST let tests inject the fake-binary path and behavior variables
   per-call, replacing `t.Setenv` in every test where the variable's only
   consumer is the spawned child.
3. MUST declare `t.Parallel()` in every test the seam unblocks; a test that
   legitimately keeps process state states its reason in one line and stays
   sequential.
4. MUST validate the new parallelism with `-race -count=2` on the package;
   a race or order-dependent failure is a defect to fix at the root, not a
   reason to revert the declaration.
5. MUST keep the runner's public behavior byte-identical: selection,
   probing, session setup, and error surfaces change in signature plumbing
   only.

## Subtasks

- [ ] Add the explicit-base composition with the boundary default.
- [ ] Convert the fake-ACPX injection from process env to per-call values.
- [ ] Declare parallelism; state reasons for any test left sequential.
- [ ] Run `-race -count=2` and fix what it surfaces.

## Acceptance Criteria

- [ ] `t.Setenv` count in `internal/agent` drops to the tests that verify
      process-default reading itself, each with a stated reason.
- [ ] Parallel declarations cover the unblocked majority, and
      `go test ./internal/agent -race -count=2` passes.
- [ ] The package's wall clock drops measurably under `go test ./...`,
      recorded in the Task result for task_06's report.
- [ ] `git status --porcelain` shows no path outside `internal/agent/` and
      this task file.

## Verification

- `go test ./internal/agent -count=1 -parallel 16` — expected: exit 0.
- `go test ./internal/agent -race -count=2 -parallel 16` — expected: exit
  0; no data race across the new parallelism.
- `go build -buildvcs=false ./...` — expected: exit 0.

## References

- `_prd.md` → Goal 2.
- `_techspec.md` → Implementation Design (the acpxCommandEnv sketch);
  Risks (parallelising surfaces real defects — as it did in cli and
  daemon).
- ADR-0089.
