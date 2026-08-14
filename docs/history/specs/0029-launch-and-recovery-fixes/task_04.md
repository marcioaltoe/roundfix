---
task: task_04
spec: 0029-launch-and-recovery-fixes
status: completed
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

- [x] Model check wired into the doctor check sequence after the agent probe
- [x] Failure rendering with advertised list and next action from the typed error
- [x] Doctor tests with fake probes: success line, failure line with list and next action, exit-code behavior

## Acceptance Criteria

- [x] Doctor with a passing probe prints `model: ok (<model>)` for the configured runtime's effective model
- [x] Doctor with a rejected model prints the failure line including the advertised list and a `next:` action, and exits nonzero
- [x] All pre-existing doctor lines are unchanged
- [x] The full test suite passes

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

## Result

Implemented the Doctor Command model check line. `roundfix doctor` now prints `model: ok (<model>)` immediately after the existing `agent:` line when the configured Agent probe succeeds, and prints a failed `model:` line with the rejected Agent Model, advertised list, and `next:` action when the existing probe returns `agent.ModelNotAdvertisedError`. Other Agent probe failures skip the model line, so unrelated Doctor failure behavior remains owned by the existing checks.

Evidence:
- Pre-change signal: `rtk proxy grep -q "model:" internal/cli/doctor.go` exited 1 before the implementation.
- Regression signal: `rtk go test ./internal/cli/ -run 'Doctor'` failed before implementation because the new Doctor model-line test contract was absent.
- Task verification: `rtk proxy grep -q "model:" internal/cli/doctor.go && rtk go test ./internal/cli/ -run 'Doctor' && rtk go test ./internal/cli/... && rtk go build -buildvcs=false ./...` passed with 9 Doctor tests and 431 CLI tests.
- Full gate: `rtk make verify` passed with `go test ./...` reporting 1247 tests, `roundfix skills check` passing, and `go build -buildvcs=false -o bin/roundfix ./cmd/roundfix` passing.

Acceptance evidence:
- Passing probe: `TestRunDoctorReportsReadinessChecks` and `TestRunDoctorAcceptsConfiguredEmptyReasoningEffort` assert `model: ok (gpt-5.5)` and `model: ok (gpt-5.6-sol)`.
- Rejected model: `TestRunDoctorReportsModelRejectionWithNextAction` asserts the failed `model:` line includes `advertised: gpt-5.5, gpt-5.1`, a `next:` action from the typed error, exit code `1`, and exactly one configured Agent probe.
- Existing lines unchanged: the exact stdout fixtures in `TestRunDoctorReportsReadinessChecks` keep the existing `node:`, `acpx:`, `adapter:`, `agent:`, and `codex:` lines byte-identical with only the inserted `model:` line.
- Full suite: `rtk make verify` passed after all edits.
