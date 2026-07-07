---
task: task_04
spec: 0019-run-outcome-notifications
status: pending
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

- [ ] README Config section: notify keys, env contract, native fallback
- [ ] README Command Boundaries: notification behavior and failure shape
- [ ] Usage guide Detached Run flow update
- [ ] roundfix SKILL.md update and `make skills-sync`

## Acceptance Criteria

- [ ] README documents both keys, the environment contract, the bound, and
      the per-platform native behavior.
- [ ] The usage guide names the notification in the detached monitoring flow.
- [ ] `make skills-sync-check` reports no drift and `roundfix skills check`
      passes.

## Verification

- `rtk make verify` — expected: fmt-check, tests, skills-sync-check, skills
  check, and build all pass.

## References

`_prd.md` → Goals; User Experience. `_techspec.md` → API Contracts; Build
Order 4. CLAUDE.md SKILL.md-matches-CLI HARD RULE.
