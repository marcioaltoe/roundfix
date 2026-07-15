---
task: task_01
spec: 0029-launch-and-recovery-fixes
status: completed
type: backend
complexity: high
---

# Task 01: Split the Detached Run handshake into liveness and Run-creation phases

## Overview

Fix the root cause of the silent `--detach` failure: the single 10-second handshake deadline is shorter than a real Preflight Validation (the model probe measured 11.4s on the dogfood machine), so the parent kills a healthy child mid-preflight and relays an empty console temp with no message. The handshake becomes two-phase — a liveness marker immediately on entering child mode, then the existing run-id line once the Run exists — and every failure branch prints an explicit diagnostic.

## Requirements

1. MUST make the detached child write a liveness marker byte on the handshake descriptor immediately on entering child mode, before configuration load and Preflight Validation, followed later by the existing run-id/console-log line unchanged.
2. MUST bound the parent's wait per phase: the existing short deadline (10s) applies only to the liveness marker; Run creation gets its own generous ceiling (5 minutes), and a child exit still ends the wait immediately in either phase.
3. MUST print an explicit stderr diagnostic on every failure branch, following the TechSpec's Interfaces shapes: liveness timeout (killed, exit/signal named), Run-creation timeout (killed, exit/signal named), child exit before handshake with output (named exit code, console relayed), and child exit with no output (stated explicitly).
4. MUST keep the four-line success report, all exit codes, and the console-log rename flow byte-identical to today.
5. MUST make both phase durations injectable for tests, following the existing timeout-constant pattern.

## Subtasks

- [x] Child-side liveness marker before config load, run-id line unchanged after Run creation
- [x] Parent-side two-phase wait with per-phase deadlines and child-exit short-circuit
- [x] Failure-branch diagnostics for all four shapes
- [x] Tests with injectable timings: slow-before-liveness, liveness-then-slow-preflight (passes under the ceiling, fails past a shrunk ceiling), child exits 2 with output, child exits silently
- [x] Existing detach success test still passes unchanged

## Acceptance Criteria

- [x] A child stub that sleeps past the liveness deadline is killed and the parent prints the liveness-timeout diagnostic naming the exit or signal
- [x] A child that signals liveness and then takes longer than the old 10s (but under the ceiling) produces the normal four-line success report
- [x] A child that exits before the handshake has its exit code named and its output relayed; a silent child produces the "produced no output" line
- [x] The success path output is byte-identical to the current contract
- [x] The full test suite passes

## Context

- interface: `internal/cli/detach.go`
- interface: `internal/cli/detach_unix.go`

## Verification

- `grep -q "liveness" internal/cli/detach.go` — expected: exit 0 (two-phase handshake exists)
- `go test ./internal/cli/ -run 'Detach'` — expected: all detach tests pass, including the new phase-matrix tests
- `go test ./internal/cli/...` — expected: all tests pass
- `go build -buildvcs=false ./...` — expected: clean build

## References

`_prd.md` → Goals 1–2, Core Feature 1, Problem 1; `_techspec.md` → Build Order 1, Interfaces (two-phase handshake, failure diagnostics), Risks (phase-2 ceiling); ADR-0028.

## Result

Implemented the two-phase Detached Run handshake. The child writes a liveness marker as soon as child mode is detected, before configuration load and Preflight Validation. The parent now waits for liveness with the existing 10s deadline, then waits separately for Run creation with a 5m ceiling. The existing run-id/console-log line, console-log rename flow, success report, and exit-code behavior are preserved.

Added explicit failure diagnostics for liveness timeout, Run-creation timeout, child exit with console output, and silent child exit. Added phase-matrix tests with injectable liveness and Run-creation durations, plus updated the detached preflight-failure test to expect the new diagnostic before relayed console output.

Evidence:

- Red signal before implementation: `rtk proxy grep -q "liveness" internal/cli/detach.go` exited 1.
- Regression matrix: `rtk go test ./internal/cli/ -run 'TestRunDetachedCommand'` failed before the two-phase implementation, then passed with 5 tests.
- Task verification: `rtk proxy grep -q "liveness" internal/cli/detach.go` passed.
- Task verification: `rtk go test ./internal/cli/ -run 'Detach'` passed with 15 tests.
- Task verification: `rtk go test ./internal/cli/...` passed with 429 tests.
- Task verification: `rtk go build -buildvcs=false ./...` passed.
- Full gate: `rtk make verify` passed: `rtk go test ./...` reported 1240 tests passed in 19 packages, `roundfix skills check` passed, and `go build -buildvcs=false -o bin/roundfix ./cmd/roundfix` passed.

Follow-ups:

- None for this task slice.
