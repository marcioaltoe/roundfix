---
task: task_01
spec: 0019-run-outcome-notifications
status: pending
type: backend
complexity: low
---

# Task 01: notify config section with enabled and command

## Overview

Add the `notify` configuration section — `enabled` (bool, default true) and
`command` (string, default empty) — with the standard precedence and strict
validation. Verifiable on its own through config package tests.

## Requirements

1. MUST add `notify.enabled` (bool, built-in default `true`) and
   `notify.command` (string, built-in default empty) with Project Config over
   User Config over built-in default precedence.
2. MUST keep strict unknown-key validation behavior unchanged for the new
   section.
3. MUST accept an empty command (meaning the native mechanism) and any
   non-empty string as the command; no command validation beyond string
   shape.

## Subtasks

- [ ] Config struct, defaults, and precedence wiring for the notify section
- [ ] Strict-key validation coverage for the section
- [ ] Config tests: defaults, per-scope override, unknown key inside notify

## Acceptance Criteria

- [ ] With no configuration, notifications are enabled with an empty command.
- [ ] Project Config overrides User Config for both keys.
- [ ] An unknown key under `notify` fails strict validation.

## Verification

- `rtk go test ./internal/config/` — expected: all tests pass, including the
  new notify section tests.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all
  pass.

## References

`_prd.md` → Core Feature 4. `_techspec.md` → Data Models; Build Order 1.
