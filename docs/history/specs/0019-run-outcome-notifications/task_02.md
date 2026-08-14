---
task: task_02
spec: 0019-run-outcome-notifications
status: completed
type: backend
complexity: medium
---

# Task 02: notify package: payload, command and native notifiers

## Overview

Build the notifier: the Outcome payload, the configured-command notifier with
its environment contract and bounded duration, the platform-native desktop
notifiers behind build tags, and the selection logic (command when configured,
native otherwise, no-op when disabled or unavailable). Verifiable in isolation
through package tests with fake commands and a lookup fake.

## Requirements

1. MUST define the Outcome payload carrying the Run id, terminal state,
   command kind, and target, and a Notifier interface with one send
   operation.
2. MUST implement the command notifier: runs the configured command through
   the shell with `ROUNDFIX_RUN_ID`, `ROUNDFIX_OUTCOME`, `ROUNDFIX_KIND`, and
   `ROUNDFIX_TARGET` in the environment, bounded by a 30s timeout, output
   discarded on success and captured into the returned error on failure.
3. MUST implement native desktop notifiers behind platform build tags —
   macOS via `osascript`, Linux via `notify-send` when present — with a
   silent no-op elsewhere or when the tool is missing. The notification names
   the outcome and target under a Roundfix title.
4. MUST implement selection: disabled → no-op; command configured → command
   notifier; otherwise → native notifier.
5. MUST be context-first and never panic; every failure returns an error the
   caller reports best-effort.

## Subtasks

- [x] Outcome payload and Notifier interface
- [x] Command notifier with environment contract, timeout, and error capture
- [x] Native notifiers behind build tags with missing-tool no-op
- [x] Selection logic from config
- [x] Package tests: fake command env/timeout, selection table, no-op paths

## Acceptance Criteria

- [x] A fake configured command receives all four `ROUNDFIX_*` variables and
      its non-zero exit surfaces as an error naming the command.
- [x] A command exceeding the bound is killed and reported as a timeout
      error.
- [x] Selection returns no-op when disabled, the command notifier when a
      command is configured, and the native path otherwise.
- [x] With no native tool available, sending is a silent success no-op.

## Verification

- `rtk go test ./internal/notify/` — expected: all tests pass.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all
  pass.

## References

`_prd.md` → User Stories 2-3; Core Features 2-3. `_techspec.md` →
Interfaces; API Contracts; Build Order 2; Decisions (build-tag pattern, 30s
bound).

## Result

Added `internal/notify` with the `Outcome` payload, `Notifier` interface,
config-based selection, a shell-backed configured-command notifier, native
desktop notifiers behind platform build tags, and a no-op fallback. Command
notifications run with the `ROUNDFIX_*` environment contract and a 30s bound.
Native notifications use `osascript` on macOS, `notify-send` on Linux, and
skip silently when the native tool is unavailable.

Pre-change signal:

- `rtk proxy go test ./internal/notify/` failed because `internal/notify` did
  not exist.
- After adding package tests, `rtk go test ./internal/notify/` failed on the
  missing notify types and constructors.

Verification:

- `rtk go test ./internal/notify/`: passed; 7 tests passed in 1 package.
- `rtk make verify`: passed; `go test ./...` reported 871 tests passed in 19
  packages, `roundfix skills check` passed, and the binary build completed.

Acceptance evidence:

- `TestCommandNotifierPassesEnvironmentAndReportsFailure` verifies a fake
  command receives `ROUNDFIX_RUN_ID`, `ROUNDFIX_OUTCOME`, `ROUNDFIX_KIND`, and
  `ROUNDFIX_TARGET`; the non-zero path returns an error naming the configured
  command and including captured output.
- `TestCommandNotifierTimesOut` verifies a context-bound fake command exceeds
  the configured bound and returns a timeout error wrapping
  `context.DeadlineExceeded`.
- `TestNewSelectsNotifier` verifies disabled config selects no-op, configured
  command selects the command notifier, and enabled config with no command
  selects the native path.
- `TestDesktopNotifierMissingToolIsNoop` verifies a missing native executable
  does not run a command and returns nil.

Follow-ups: CLI terminal-outcome wiring and warning/reporting belong to the
next task.
