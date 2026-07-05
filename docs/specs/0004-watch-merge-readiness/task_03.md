---
task: task_03
spec: 0004-watch-merge-readiness
status: completed
type: backend
complexity: low
---

# Task 03: Agent console suppression for supervisors

## Overview

Non-TTY stderr interleaves Daemon milestones with the full Agent console,
forcing supervisors into fragile pattern filtering. Add a display-only filter:
a Sink decorator that drops Agent-source events, exposed as
`--no-agent-console` on the operational commands. Verifiable through sink
unit tests and buffer-captured stderr asserts.

## Requirements

1. MUST add a filtering Sink decorator in the run-event package that drops
   events whose Source is the Agent and forwards everything else, preserving
   order and the critical-sink error contract of the wrapped sink.
2. MUST add `--no-agent-console` to resolve, watch, and implement: in non-TTY
   mode the stderr writer sink is wrapped; the Journal sink is NEVER wrapped
   (ADR-0008), and the interactive cockpit path rejects the flag with the
   existing validation error shape (the cockpit already separates surfaces).
3. MUST leave header/progress stderr lines (non-event diagnostics) untouched.
4. MUST document the flag in each command's help text.

## Subtasks

- [x] Source-filter Sink decorator with unit tests
- [x] Flag wiring on the three operational commands
- [x] TTY rejection and help text
- [x] Buffer-captured stderr asserts over a fake Agent run

## Acceptance Criteria

- [x] A full fake run with `--no-agent-console` yields stderr containing
      every Daemon-source event line and zero Agent-source lines; the
      Journal contains both.
- [x] Without the flag, stderr is byte-identical to today.
- [x] `--interactive --no-agent-console` fails with the existing
      conflicting-flags error shape.

## Verification

- `rtk go test ./internal/runevent/ ./internal/cli/` — expected: all tests
  pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 4; Core Feature 4; Decisions (display-only).
`_techspec.md` → Interfaces (NewSourceFilterSink), Build Order 3. Dogfood
finding 4. ADR-0008, ADR-0009.

## Result

- Added `runevent.NewSourceFilterSink`, which drops events from one configured
  source and forwards every other Run Event in order while preserving wrapped
  sink errors for forwarded events.
- Added `--no-agent-console` to `resolve`, `watch`, and `implement`. In
  non-TTY mode, only the stderr writer sink is wrapped with the Agent-source
  filter; the Run Event Journal sink is unchanged. The flag is rejected for
  explicit `--interactive` input and for the interactive cockpit path.
- Help text for `resolve`, `watch`, and `implement` now documents
  `--no-agent-console`.
- `TestRunResolveNoAgentConsoleSuppressesAgentDisplayOnly`,
  `TestRunWatchNoAgentConsoleSuppressesAgentDisplayOnly`, and
  `TestRunImplementNoAgentConsoleSuppressesAgentDisplayOnly` prove fake full
  Runs keep daemon/progress stderr lines, hide Agent console lines, and keep
  Agent and Daemon events in the Run Event Journal.
- `TestAgentConsoleDisplaySinkKeepsWriterBytesByDefault` proves the default
  no-flag display path emits the same bytes as `agent.WriterSink`.
- `TestRunOperationalCommandRejectsInvalidInput`,
  `TestRunImplementValidationFailures`, and
  `TestRunNoAgentConsoleRejectsInteractiveCockpit` cover the conflicting flag
  and cockpit rejection paths.

Verification:

- `rtk go test ./internal/runevent/ ./internal/cli/` — passed, 196 tests.
- `rtk go test ./...` — passed, 507 tests across 16 packages.
- `rtk make verify` — passed (`go test ./...`, `roundfix skills check`, and
  `go build`).
