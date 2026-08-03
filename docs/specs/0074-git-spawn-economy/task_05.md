---
task: task_05
spec: 0074-git-spawn-economy
status: completed
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

## Result

Implementation-ready handoff for Daemon Verification:

- Explicit environment: `ACPXRunner.Environment` now supplies the base for
  ACPX version probes, adapter inspection/config resolution, Codex hygiene
  resolution, session setup, prompts, cancellation, listing, and sealed
  prompts. A nil base resolves `os.Environ()` at the child-command boundary;
  `acpxCommandEnv(base, overrides)` preserves the existing guard filtering and
  per-session override behavior.
- Per-call fake injection: fake ACPX behavior, prompt paths, milestone paths,
  adapter configuration, and Codex-path values now live on each test runner.
  `rtk proxy rg -n 't\.Setenv\(' internal/agent --glob '*_test.go'` reports
  three calls: the zero-value process-environment default test and the two
  process-`PATH` adapter-resolution tests. Each stays sequential with a
  one-line reason.
- Parallel coverage: the top-level inventory is 137 tests, with 129 declaring
  `t.Parallel()` and 8 stated sequential exceptions. The exceptions are the
  three process-environment tests above and five `ProbeFallback` tests that
  exercise its process working-directory default through `t.Chdir`.
- Focused concurrency check:
  `GOCACHE=<worktree>/.gocache go test ./internal/agent -race -count=2 -parallel 12`
  exited 0 with no race report after the final production edit. This differs
  from the Daemon-owned `-parallel 16` Verification command.
- Timing evidence for task_06: before this Task,
  `go test ./internal/agent -count=1 -parallel 12` reported package time
  10.250s and real time 11.01s. After the change, the same command reported
  package time 6.305s and real time 6.72s, a 39% real-time reduction. A fresh
  post-change `go test ./... -count=1 -parallel 12` measurement reported
  `internal/agent` at 20.221s while packages competed for the same machine;
  task_06 owns the full-suite before/after attribution.
- Scope: `git -c core.fsmonitor=false status --porcelain` reports only
  `internal/agent/` and this Task file. `rtk git diff --check` exits 0.

The three commands under `## Verification` were not run; the Daemon owns
them and Task settlement.
