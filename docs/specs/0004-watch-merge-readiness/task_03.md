---
task: task_03
spec: 0004-watch-merge-readiness
status: pending
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

- [ ] Source-filter Sink decorator with unit tests
- [ ] Flag wiring on the three operational commands
- [ ] TTY rejection and help text
- [ ] Buffer-captured stderr asserts over a fake Agent run

## Acceptance Criteria

- [ ] A full fake run with `--no-agent-console` yields stderr containing
      every Daemon-source event line and zero Agent-source lines; the
      Journal contains both.
- [ ] Without the flag, stderr is byte-identical to today.
- [ ] `--interactive --no-agent-console` fails with the existing
      conflicting-flags error shape.

## Verification

- `rtk go test ./internal/runevent/ ./internal/cli/` — expected: all tests
  pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 4; Core Feature 4; Decisions (display-only).
`_techspec.md` → Interfaces (NewSourceFilterSink), Build Order 3. Dogfood
finding 4. ADR-0008, ADR-0009.
