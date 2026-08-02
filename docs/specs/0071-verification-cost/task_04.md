---
task: task_04
spec: 0071-verification-cost
status: pending
type: backend
complexity: high
---

# Task 04: Free the Baseline package from process state

## Overview

The Baseline package is the second heaviest at 83.9s, with 16 environment
mutations and one working-directory change blocking parallelism. This Task
applies the same two moves as the CLI package — remove the process-global
dependency, then declare parallelism — on a smaller surface. It reduces the sum
of package times and speeds a single-package run; it cannot move the suite
floor, which the CLI package sets.

## Requirements

1. MUST make functions reading process environment or working directory, for
   values their callers know, receive them as parameters, with the default
   resolved once at the command boundary.
2. MUST declare parallelism on every test that no longer mutates process state
   and owns the filesystem it touches.
3. MUST leave a one-line reason on every test that stays sequential.
4. MUST fix, not silence, tests that fail only under parallel execution.
5. MUST prove no races and no cross-test leakage.
6. MUST leave the coverage record and the package's observable behavior
   unchanged.

## Subtasks

- [ ] Give the process-reading functions their parameters.
- [ ] Declare parallelism where process state is no longer touched.
- [ ] State a reason on every test left sequential.
- [ ] Fix shared-state failures and prove no races.

## Acceptance Criteria

- [ ] Production functions no longer read the process for values callers
      supply.
- [ ] The package's parallel-declaring test count rises above its recorded 28.
- [ ] Every test left sequential carries a one-line reason.
- [ ] The package passes with race detection and with its tests run twice in
      one invocation.
- [ ] The coverage record from task 01 is unchanged.
- [ ] The package completes measurably faster than 83.9s on the same machine.
- [ ] `git status --porcelain` shows no path outside `internal/baseline/`,
      `internal/`, and this task file.

## Context

- interface: `internal/baseline/plan.go`
- interface: `internal/baseline/catalog_test.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/baseline -count=1` — expected: exit 0.
- `go test ./internal/baseline -count=2` — expected: exit 0.
- `go test ./internal/baseline -race -count=1` — expected: exit 0.
- `go test ./internal/spec -run '^TestCoverageEquivalence$' -count=1 -v | grep -q -- "--- PASS: TestCoverageEquivalence"`
  — expected: exit 0.
- `go vet ./internal/baseline` — expected: exit 0.

## References

- `_prd.md` → Core Features 1 and 2.
- `_techspec.md` → Build Order 4; Risks (task 04 helps the sum, not the floor).
- ADR-0089.
