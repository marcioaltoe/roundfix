---
task: task_02
spec: 0019-run-outcome-notifications
status: pending
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

- [ ] Outcome payload and Notifier interface
- [ ] Command notifier with environment contract, timeout, and error capture
- [ ] Native notifiers behind build tags with missing-tool no-op
- [ ] Selection logic from config
- [ ] Package tests: fake command env/timeout, selection table, no-op paths

## Acceptance Criteria

- [ ] A fake configured command receives all four `ROUNDFIX_*` variables and
      its non-zero exit surfaces as an error naming the command.
- [ ] A command exceeding the bound is killed and reported as a timeout
      error.
- [ ] Selection returns no-op when disabled, the command notifier when a
      command is configured, and the native path otherwise.
- [ ] With no native tool available, sending is a silent success no-op.

## Verification

- `rtk go test ./internal/notify/` — expected: all tests pass.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all
  pass.

## References

`_prd.md` → User Stories 2-3; Core Features 2-3. `_techspec.md` →
Interfaces; API Contracts; Build Order 2; Decisions (build-tag pattern, 30s
bound).
