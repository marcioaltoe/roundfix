---
task: task_03
spec: 0037-terminal-outcome-integrity
status: completed
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

- [x] Introduce the owner-process controller seam.
- [x] Implement Unix exit proof and preserve Windows build behavior.
- [x] Reorder Force Stop coordination around owner absence.
- [x] Add actionable failed-step diagnostics.
- [x] Cover graceful, forced, permission, deadline, and idempotent paths.
- [x] Add a real helper-process integration case on Unix.

## Acceptance Criteria

- [x] Successful Force Stop proves owner exit before Stopped is persisted.
- [x] The Active Run lock remains present until the winning completion.
- [x] Permission or deadline failure prints no success report, stores no Stopped
      outcome, and leaves the Run Active.
- [x] The diagnostic names the Run, owner PID, failed step, and retained state.
- [x] A reused or unprovable PID fails closed.
- [x] Repeating Force Stop for an already Stopped Run performs no process or
      session action.
- [x] Unix integration proves the helper owner is absent before store
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

## Result

Force Stop now treats owner-process absence as the completion gate. It cancels
registered Agent Sessions, sends graceful termination to the recorded PID,
escalates through the platform force-kill path when needed, and calls
`CompleteRun` only after the controller proves the PID absent. Owner-control
failures leave the Run Active and the Active Run lock retained. An already
Stopped Run returns its stored result without process, session, or Worktree
actions.

Verification:

- Pre-change reproduction: the focused CLI/store test build failed because the
  owner-process seam, controller errors, and platform implementation did not
  exist.
- `rtk go test ./internal/cli -run 'Test.*ForceStop.*(Owner|Exit|Permission|Deadline|Idempotent|Lock)' -count=1`
  passed with 8 tests.
- `rtk go test ./internal/store -run 'Test.*Process|TestStoppedRunReleasesActiveLock' -count=1`
  passed with 8 tests.
- `rtk go test -race ./internal/cli ./internal/store -run 'Test.*(ForceStop|OwnerProcess)' -count=1`
  passed with 13 tests.
- `rtk go test ./internal/cli ./internal/store -count=1` passed.
- `GOOS=windows GOARCH=amd64 go build ./internal/store ./internal/cli`
  passed through `rtk proxy env`.
- `rtk make verify` passed.
- `rtk git -c core.fsmonitor=false diff --check` passed.

Acceptance evidence:

- `TestRunForceStopOwnerExitPrecedesCompletionAndLockRelease` observed the Run
  as Active and a competing Run blocked by its lock while the owner controller
  was proving exit, then observed Stopped only after proof.
- `TestRunForceStopOwnerPermissionAndDeadlineFailuresRetainActiveLock` proved
  permission and deadline failures produce empty stdout, actionable Run/PID/step
  diagnostics, Active state, and a retained lock.
- `TestRunForceStopOwnerPIDReuseFailsClosed` and
  `TestOwnerProcessControllerRejectsUnprovenCurrentProcess` proved conservative
  refusal when process identity cannot be trusted.
- `TestRunForceStopStoppedRunIsIdempotentWithoutOwnerOrSessionActions` proved a
  stored Stopped replay invokes neither process control nor Agent Session
  cleanup.
- `TestOwnerProcessControllerGracefulExitProof` and
  `TestOwnerProcessControllerForceKillExitProof` exercised real Unix helper
  processes through graceful and forced exit.
- `TestRunForceStopOwnerProcessIntegrationProvesExitBeforeStoreCompletion`
  observed the helper PID absent when Stopped first became visible in the Run
  Database.

Follow-up:

- Windows production packages cross-build. Cross-compiling the full
  `internal/cli` test binary remains blocked by pre-existing Unix-only
  `syscall.SysProcAttr.Setpgid` and `syscall.Kill` references in
  `internal/cli/implement_test.go`; this Task did not change that test surface.
