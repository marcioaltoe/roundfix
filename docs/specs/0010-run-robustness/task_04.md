---
task: task_04
spec: 0010-run-robustness
status: pending
type: docs
complexity: low
---

# Task 04: Sync docs and the Roundfix skill with the robustness changes

## Overview

Document the shipped surface: the config deprecation policy and warning
shape, the session-close reaping reports, and Detached Runs with their
four-line report and follow/stop surfaces — in the canonical Roundfix
skill and README, cross-checked against the binary. Verifiable through the
skills drift check inside the full gate.

## Requirements

1. MUST document in the canonical Roundfix skill: `--detach` on the three
   operational commands (report shape, attach/stop follow-up, console-log
   location, Preflight relay), the deprecation warning behavior, and the
   force-stop/sweep session-close reports; regenerate the embedded copy
   through the sync target.
2. MUST update the README: the config-migration promise (removed keys warn,
   never break) and the Detached Run usage pattern for scripts/CI.
3. MUST cross-check every documented line shape against CLI test fixtures
   and the built binary.
4. MUST verify glossary coverage — Detached Run exists; call out any
   further gap instead of inventing language.

## Subtasks

- [ ] Skill updates + `make skills-sync`
- [ ] README migration promise and detach usage
- [ ] Fixture and binary cross-check
- [ ] Glossary pass

## Acceptance Criteria

- [ ] Skill text matches shipped behavior exactly; drift check passes
      inside the full gate.
- [ ] The four detach lines and the deprecation warning appear verbatim in
      CLI test fixtures.
- [ ] README carries both notes.
- [ ] No new un-glossaried term.

## Verification

- `rtk go run ./cmd/roundfix skills check` — expected: skill artifacts
  validate.
- `make verify` — expected: full gate passes, including the drift check.

## References

`_prd.md` → User Experience; Core Features 1–3. `_techspec.md` → Build
Order 4. ADR-0027, ADR-0028. Repo hard rule (canonical skill ships with CLI
behavior changes).
