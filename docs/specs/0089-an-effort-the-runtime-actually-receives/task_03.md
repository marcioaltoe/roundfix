---
task: task_03
spec: 0089-an-effort-the-runtime-actually-receives
status: pending
type: backend
complexity: medium
---

# Task 03: Stop refusing an OpenCode reasoning effort

## Overview

Spec 0088 refused a non-empty `reasoning_effort` on `opencode` at two gates:
configuration normalization and runtime validation. Both refusals exist because
the effort could not be applied. Task 02 made it plannable; this Task removes
the refusals so a maintainer can configure one at all.

## Requirements

1. MUST remove the configuration refusal so a non-empty `reasoning_effort` on
   `runtime: opencode` loads.
2. MUST restore `opencode` to the generic reasoning-effort config key so runtime
   validation accepts the selection.
3. MUST keep an empty `reasoning_effort` on `opencode` valid, still planning
   `runtime_managed`.
4. MUST NOT change how Codex or Claude map to their reasoning keys.
5. MUST leave no unreachable remnant of the refusal — the error type and the
   runtime list it consulted go with it if nothing else uses them.
6. MUST re-record the coverage record in this Task's own commit if any test is
   renamed or removed.

## Subtasks

- [ ] Remove the normalization refusal and its runtime list.
- [ ] Restore the `opencode` reasoning-effort key mapping.
- [ ] Remove the now-unused error type if nothing else references it.
- [ ] Edit the break-half characterization tests and declare the breaks.

## Acceptance Criteria

- [ ] Configuration with `runtime: opencode` and a non-empty `reasoning_effort`
      loads and resolves with that effort intact.
- [ ] Configuration with `reasoning_effort: ""` on `opencode` still loads.
- [ ] Runtime validation accepts an `opencode` runtime carrying a non-empty
      effort.
- [ ] Codex and Claude still map to their existing reasoning keys.
- [ ] `grep -rn "must be empty for runtime" internal/config` finds nothing.

## Context

- interface: `internal/config/profiles.go`
- interface: `internal/agent/acpx_runner.go`

## Bounded scope

This Task may create or modify only:

- `internal/config/profiles.go`
- `internal/config/config_test.go`
- `internal/config/opencode_effort_characterization_test.go`
- `internal/agent/acpx_runner.go`
- `internal/agent/acpx_runner_test.go`
- `docs/references/coverage-record.json`
- `docs/specs/0089-an-effort-the-runtime-actually-receives/task_03.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exits 0.
- `go test ./internal/config ./internal/agent -count=1` — expected: exits 0.
- `go test ./internal/config -run 'OpenCodeEffortAccepted' -count=1 -v` — expected: exits 0 and names at least one test; `no tests to run` fails this Task.
- `grep -rn 'must be empty for runtime' internal/config` — expected: exits non-zero with no output, proving the refusal text is gone rather than reworded.
- `go test ./internal/spec -run '^TestCoverageEquivalence$' -count=1` — expected: exits 0.

## References

- `_prd.md` → Goal 1; Core Features: a configuration that stops lying.
- `_techspec.md` → Implementation Design: API Contracts; Build Order 3.
- ADR-0108.
