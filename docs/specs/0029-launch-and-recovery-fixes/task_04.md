---
task: task_04
spec: 0029-launch-and-recovery-fixes
status: pending
type: backend
complexity: low
---

# Task 04: Report the effective Agent Model probe in the Doctor Command

## Overview

Give the advertised-list drift a self-service diagnostic: `roundfix doctor` gains one `model:` check line after the existing `agent:` line, reporting whether the configured runtime accepts the effective Agent Model — and, on failure, the advertised models and a next action. Reuses the existing disposable-session probe; no new acpx invocation shapes.

## Requirements

1. MUST add a `model:` check line to the Doctor Command rendering the effective Agent Model probe outcome for the configured agent: `model: ok (<model>)` on success, `model: failed (...)` with the advertised list and a `next:` action on failure, following the existing doctor line format.
2. MUST reuse the existing probe path and the typed model-rejection error for the failure content.
3. MUST keep doctor's exit-code behavior consistent with the existing checks (any failed check exits nonzero, skipped checks do not).
4. MUST leave all existing doctor lines byte-identical.

## Subtasks

- [ ] Model check wired into the doctor check sequence after the agent probe
- [ ] Failure rendering with advertised list and next action from the typed error
- [ ] Doctor tests with fake probes: success line, failure line with list and next action, exit-code behavior

## Acceptance Criteria

- [ ] Doctor with a passing probe prints `model: ok (<model>)` for the configured runtime's effective model
- [ ] Doctor with a rejected model prints the failure line including the advertised list and a `next:` action, and exits nonzero
- [ ] All pre-existing doctor lines are unchanged
- [ ] The full test suite passes

## Context

- interface: `internal/cli/doctor.go`
- interface: `internal/agent/acpx_runner.go`

## Verification

- `grep -q "model:" internal/cli/doctor.go` — expected: exit 0 (doctor model line exists)
- `go test ./internal/cli/ -run 'Doctor'` — expected: doctor tests pass, including the model line coverage
- `go test ./internal/cli/...` — expected: all tests pass
- `go build -buildvcs=false ./...` — expected: clean build

## References

`_prd.md` → Goal 3, Core Feature 3; `_techspec.md` → Build Order 4, API Contracts (doctor line), Decisions (batch-time classification over probe hardening).
