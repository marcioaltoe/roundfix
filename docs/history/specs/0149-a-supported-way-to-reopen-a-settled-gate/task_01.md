---
task: task_01
spec: 0149-a-supported-way-to-reopen-a-settled-gate
status: completed
type: backend
complexity: medium
---

# Task 01: Expose the QA Task behind a stale-gate refusal

## Overview

`spec.Load` already decides whether a settled QA gate still describes its graph,
and returns `StaleGateError` when it does not. A caller that wants to act on that
answer cannot: `Load` returns a nil graph, and the error carries Task ids but no
path to the QA Task file.

This Task adds the recovery seam. It changes no existing caller's behavior.

## Requirements

1. MUST provide a way for one caller to obtain the QA Task's file path and the
   stale dependency ids when `Load` refuses with `StaleGateError`.
2. MUST leave `Load` refusing exactly as it does today, with the same error and
   the same message, for every existing caller.
3. MUST NOT re-derive staleness. The loader's own predicate is the authority;
   a second implementation is a second thing that can disagree with the gate.
4. MUST NOT change when the refusal is raised.

## Subtasks

- [x] Add the seam and its unit tests.
- [x] Assert that `Load` still refuses on a stale gate.

## Acceptance Criteria

- [x] A fixture Spec with a completed QA Task above a pending dependency yields
      the QA Task's path and the stale dependency ids through the new seam.
- [x] The same fixture still makes `Load` return `StaleGateError`.
- [x] A healthy fixture yields no stale-gate answer.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/spec/spec.go`
- interface: `internal/spec/errors.go`

## Verification

- `out="$(go test -count=1 -run "^TestLoadForRecovery" ./internal/spec 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.
- `out="$(go test -count=1 -run "^TestLoadStillRefusesAStaleGate" ./internal/spec 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; the regression that keeps every existing caller refusing. Before this Task the run reports no tests to run, so the command fails.

## References

- [_techspec.md](_techspec.md) — The condition

## Result

### Implementation

- Added `LoadForRecovery`, which returns the loader's parsed graph alongside
  its existing wrapped `StaleGateError`. Other validation failures still return
  a nil graph.
- Kept `Load` as the existing refusing surface: it delegates to the recovery
  seam and discards the graph whenever any error is returned.

### Focused checks

- `rtk go test -count=1 -run 'Test(LoadForRecovery|LoadStillRefusesAStaleGate|GateStalenessCharacterizesEachVerdict)$' ./internal/spec`
  passed: 10 tests in 1 package.
- `rtk go test -count=1 ./internal/spec` passed: 457 tests in 1 package.
- The first sandboxed focused-test attempt could not write the normal Go build
  cache; the same command passed after the cache was made available.
- The Daemon-owned commands under `## Verification` were not run.

### Acceptance evidence

1. `TestLoadForRecovery/returns_the_graph_with_the_loader's_stale_gate_error`
   reads the QA Task path `demo/task_02.md` from the returned graph and the
   stale dependency id `task_01` from `StaleGateError`.
2. `TestLoadStillRefusesAStaleGate` proves `Load` returns a nil graph, preserves
   the typed `StaleGateError`, and preserves the exact wrapped error message.
3. `TestLoadForRecovery/returns_no_stale_gate_error_for_a_healthy_graph` proves
   a completed dependency yields a graph with no error.
