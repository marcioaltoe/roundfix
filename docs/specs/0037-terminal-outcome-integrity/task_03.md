---
task: task_03
spec: 0037-terminal-outcome-integrity
status: pending
type: backend
complexity: high
---

# Task 03: Prove owner exit before Force Stop completion

## Overview

Reorder Force Stop so Roundfix proves the recorded owner process is absent
before storing Stopped or releasing the Active Run lock. Platform-specific
controllers retain the existing graceful-then-force policy while failures
leave the Run Active with actionable diagnostics.

## Requirements

1. MUST cancel registered Agent Sessions before terminating the recorded owner.
2. MUST target only the recorded owner PID and follow existing platform process
   boundaries.
3. MUST attempt graceful termination, wait within the bounded stop window, and
   use the existing force-kill path when required.
4. MUST report success only after owner-process absence is positively proven.
5. MUST call terminal completion only after that proof.
6. MUST leave the Run Active and its lock retained on permission, unsupported,
   reused-PID, or deadline failures.
7. MUST make an already stored Stopped outcome idempotent without repeating
   process or Agent Session actions.

## Subtasks

- [ ] Introduce the owner-process controller seam.
- [ ] Implement Unix exit proof and preserve Windows build behavior.
- [ ] Reorder Force Stop coordination around owner absence.
- [ ] Add actionable failed-step diagnostics.
- [ ] Cover graceful, forced, permission, deadline, and idempotent paths.
- [ ] Add a real helper-process integration case on Unix.

## Acceptance Criteria

- [ ] Successful Force Stop proves owner exit before Stopped is persisted.
- [ ] The Active Run lock remains present until the winning completion.
- [ ] Permission or deadline failure prints no success report, stores no Stopped
      outcome, and leaves the Run Active.
- [ ] The diagnostic names the Run, owner PID, failed step, and retained state.
- [ ] A reused or unprovable PID fails closed.
- [ ] Repeating Force Stop for an already Stopped Run performs no process or
      session action.
- [ ] Unix integration proves the helper owner is absent before store
      completion; non-Unix packages still compile.

## Context

- instruction: `.agents/skills/agentic-cli-design/SKILL.md`
- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/golang-context/SKILL.md`
- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/store/process_unix.go`
- interface: `internal/store/process_windows.go`
- interface: `internal/store/process_other.go`
- interface: `internal/store/process_unix_test.go`

## Verification

- `rtk go test ./internal/cli -run 'Test.*ForceStop.*(Owner|Exit|Permission|Deadline|Idempotent|Lock)' -count=1`
  — expected: Force Stop proves absence before completion and every failure
  retains Active state and lock.
- `rtk go test ./internal/store -run 'Test.*Process|TestStoppedRunReleasesActiveLock' -count=1`
  — expected: platform process proof and winner-only lock release pass.
- `rtk go test -race ./internal/cli ./internal/store -run 'Test.*(ForceStop|OwnerProcess)' -count=1`
  — expected: stop coordination and exit proof are race-free.

## References

- `_prd.md` → Goals 1–2; User Story 1; Core Features 1 and 3; User Experience;
  Success Metrics.
- `_techspec.md` → Interfaces: OwnerProcessController; API Contracts:
  Force Stop; Build Order 3.
- `../../adr/0044-orphaned-run-locks-are-reclaimed-on-proven-owner-death.md` →
  owner-death proof.
- `../../adr/0052-run-completion-is-compare-and-set.md` → completion ordering.
