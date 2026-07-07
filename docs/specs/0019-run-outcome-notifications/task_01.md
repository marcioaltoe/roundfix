---
task: task_01
spec: 0019-run-outcome-notifications
status: completed
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

- [x] Config struct, defaults, and precedence wiring for the notify section
- [x] Strict-key validation coverage for the section
- [x] Config tests: defaults, per-scope override, unknown key inside notify

## Acceptance Criteria

- [x] With no configuration, notifications are enabled with an empty command.
- [x] Project Config overrides User Config for both keys.
- [x] An unknown key under `notify` fails strict validation.

## Verification

- `rtk go test ./internal/config/` — expected: all tests pass, including the
  new notify section tests.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all
  pass.

## References

`_prd.md` → Core Feature 4. `_techspec.md` → Data Models; Build Order 1.

## Result

Added the `notify` config section to `internal/config`: `notify.enabled`
defaults to `true`, `notify.command` defaults to empty, User Config overlays
built-ins, and Project Config overlays User Config. The generated default
config now includes the `notify` section. Unknown keys inside `notify` fail
with the same strict section-key style as the existing config sections.

Pre-change signal: after adding the notify config tests, `rtk proxy go test
./internal/config/` failed because `Config.Notify` did not exist.

Verification:

- `rtk go test ./internal/config/`: passed; 70 tests passed in 1 package.
- `rtk make verify`: passed; `go test ./...` reported 864 tests passed in 18
  packages, `roundfix skills check` passed, and the binary build completed.

Acceptance evidence:

- No configuration defaults are covered by
  `TestLoadAppliesNotifyConfigHierarchy/builtin only`, which asserts
  `notify.enabled == true` and `notify.command == ""`.
- Project Config precedence is covered by
  `TestLoadAppliesNotifyConfigHierarchy/project override`, which starts with
  User Config `enabled: false` and `command: user-notify`, then verifies
  Project Config overrides both keys with `enabled: true` and `command: ""`.
- Strict validation is covered by `TestLoadRejectsUnknownNotifyConfigKey`,
  which verifies `notify.channel` fails with
  `notify.channel is not a supported config key`.

Follow-ups: none.
