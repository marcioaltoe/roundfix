---
task: task_02
spec: 0089-an-effort-the-runtime-actually-receives
status: pending
type: backend
complexity: high
---

# Task 02: Plan a deferred effort and fail closed on an unadvertised one

## Overview

Teach the planner that a non-empty effort on a runtime which cannot accept one
before its session's first prompt is its own encoding. Preflight can see the
value advertised; it cannot apply it. Naming that case `runtime_deferred` is
what keeps a readiness surface from reporting an assignment it never made.

## Requirements

1. MUST add the `runtime_deferred` selection encoding and select it when the
   requested effort is non-empty, the runtime defers effort application, and the
   requested value appears among the values that model advertises.
2. MUST produce `SelectionUnsupportedError` when the requested effort is
   non-empty and absent from the advertised values, naming those values.
3. MUST leave `independent`, `model_variant`, `model_managed`, and
   `runtime_managed` selecting exactly as they do today.
4. MUST make the effective-state check accept a `runtime_deferred` assignment
   whose effort has not been applied yet, and MUST NOT let it accept one whose
   model is wrong.
5. MUST derive "this runtime defers effort" from one predicate so no second copy
   of the runtime list can drift from the first.
6. MUST re-record the coverage record in this Task's own commit if any test is
   renamed or removed.

## Subtasks

- [ ] Add the encoding constant with the comment that distinguishes it.
- [ ] Add the deferring-runtime predicate and use it in planning.
- [ ] Select the encoding and fail closed on an unadvertised value.
- [ ] Extend the effective-state check for the new encoding.
- [ ] Edit the break-half characterization tests and declare the breaks.

## Acceptance Criteria

- [ ] An `opencode` selection with an advertised non-empty effort plans
      `runtime_deferred` and reports the requested effort.
- [ ] The same selection with an unadvertised effort produces
      `SelectionUnsupportedError` listing the advertised values.
- [ ] A Codex selection with a non-empty effort still plans `independent`.
- [ ] An `opencode` selection with an empty effort still plans `runtime_managed`.
- [ ] The effective-state check accepts a `runtime_deferred` assignment before
      application and rejects one whose current model differs.

## Context

- interface: `internal/agent/selection_assignment.go`
- interface: `internal/agent/selection_capabilities.go`

## Bounded scope

This Task may create or modify only:

- `internal/agent/selection_assignment.go`
- `internal/agent/selection_assignment_test.go`
- `internal/agent/selection_effort_characterization_test.go`
- `docs/references/coverage-record.json`
- `docs/specs/0089-an-effort-the-runtime-actually-receives/task_02.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exits 0.
- `go test ./internal/agent -count=1` — expected: exits 0.
- `go test ./internal/agent -run 'RuntimeDeferred' -count=1 -v` — expected: exits 0 and names at least one test; `no tests to run` fails this Task.
- `go test ./internal/spec -run '^TestCoverageEquivalence$' -count=1` — expected: exits 0.
- `grep -q 'SelectionEncodingRuntimeDeferred' internal/agent/selection_assignment.go` — expected: exits 0.

## References

- `_prd.md` → Goals 2 and 4; Core Features: a proof split across two moments.
- `_techspec.md` → Implementation Design: Interfaces; Build Order 2.
- ADR-0108.
