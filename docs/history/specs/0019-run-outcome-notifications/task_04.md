---
task: task_04
spec: 0019-run-outcome-notifications
status: completed
type: docs
complexity: low
---

# Task 04: Docs and skill sync for Run Outcome Notifications

## Overview

Close the SKILL-matches-CLI gate for notifications: document the `notify`
config section and the outcome-notification behavior everywhere the config
and Run lifecycle are described, and re-sync the embedded skill bundle.
Verifiable by the drift check and by reading the docs against shipped
behavior.

## Requirements

1. MUST document `notify.enabled` and `notify.command` in the README Config
   section — defaults, precedence, the `ROUNDFIX_*` environment contract,
   the 30s bound, and the native fallback per platform — truthfully against
   implemented behavior.
2. MUST document the notification behavior in Command Boundaries: which
   commands notify, the failure warning shape, and that failures never change
   outcomes or exit codes.
3. MUST update the operational usage guide's Detached Run flow to mention the
   outcome notification as the unattended-Run signal.
4. MUST update the canonical roundfix skill (SKILL.md) and re-sync the
   embedded bundle so the drift check passes.

## Subtasks

- [x] README Config section: notify keys, env contract, native fallback
- [x] README Command Boundaries: notification behavior and failure shape
- [x] Usage guide Detached Run flow update
- [x] roundfix SKILL.md update and `make skills-sync`

## Acceptance Criteria

- [x] README documents both keys, the environment contract, the bound, and
      the per-platform native behavior.
- [x] The usage guide names the notification in the detached monitoring flow.
- [x] `make skills-sync-check` reports no drift and `roundfix skills check`
      passes.

## Verification

- `rtk make verify` — expected: fmt-check, tests, skills-sync-check, skills
  check, and build all pass.

## References

`_prd.md` → Goals; User Experience. `_techspec.md` → API Contracts; Build
Order 4. CLAUDE.md SKILL.md-matches-CLI HARD RULE.

## Result

Implemented task_04.

- Acceptance 1: README Config now documents `notify.enabled` and
  `notify.command`, built-in defaults, User Config and Project Config
  precedence, the four `ROUNDFIX_*` variables, the 30s command timeout, command
  output handling, and native behavior for macOS (`osascript`), Linux
  (`notify-send`), other platforms, and missing tools.
- Acceptance 2: `docs/usage.md` now names the configured outcome notification
  as the unattended-Run signal in the Detached Run monitoring flow.
- Acceptance 3: updated `.agents/skills/roundfix/SKILL.md` and regenerated
  `skills/roundfix/SKILL.md` with `rtk make skills-sync`. `rtk make
  skills-sync-check` passed with no drift output. `rtk go run -buildvcs=false
  ./cmd/roundfix skills check` passed after rerun with normal Go build-cache
  access.
- Verification: `rtk make verify` exited 0; it ran `rtk go test ./...` with
  878 tests in 19 packages, `roundfix skills check`, and `go build`. The
  verify target also includes `fmt-check` and `skills-sync-check`.
