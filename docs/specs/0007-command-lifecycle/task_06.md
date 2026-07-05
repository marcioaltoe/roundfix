---
task: task_06
spec: 0007-command-lifecycle
status: pending
type: docs
complexity: low
---

# Task 06: Sync docs and the Roundfix skill with the lifecycle commands

## Overview

Document the shipped surface: setup, upgrade, the freshness note, graceful
versus force stop, and the push-at-Clean config — in the canonical Roundfix
skill, README, and help text, with the embedded copy regenerated and every
shape cross-checked against the built binary. Verifiable through the skills
drift check inside the full gate.

## Requirements

1. MUST document in the canonical Roundfix skill: `setup` (checks, offers,
   `--yes`), `upgrade`/`--check` and the daily freshness note, the stop
   semantics split (graceful default records a Stop Request; `--force`
   cancels the Agent Session and completes immediately; when to use each),
   and `implement.auto_push` with its Clean-only rule; regenerate the
   embedded copy through the sync target.
2. MUST update the README's install/usage guidance: setup as the documented
   bootstrap path (subsuming the manual acpx install line), upgrade as the
   update path.
3. MUST cross-check every documented flag, outcome line, and report shape
   against the built binary's output.
4. MUST verify every term against the glossary; add missing lifecycle terms
   sparingly (candidates: Setup Command, Upgrade Command — mirror the
   existing command-entry style) rather than leaving gaps.

## Subtasks

- [ ] Skill updates + `make skills-sync`
- [ ] README bootstrap/update guidance
- [ ] Binary output cross-check
- [ ] Glossary pass with sparing additions

## Acceptance Criteria

- [ ] Skill text matches shipped behavior exactly; drift check passes
      inside the full gate.
- [ ] Documented outcome lines appear verbatim in CLI test fixtures.
- [ ] README names setup and upgrade as the paths; the manual pin command
      remains only as the fallback setup itself prints.
- [ ] Glossary covers every new command term used.

## Verification

- `rtk go run ./cmd/roundfix skills check` — expected: skill artifacts
  validate.
- `make verify` — expected: full gate passes, including the drift check.

## References

`_prd.md` → Core Features 1–5; User Experience. `_techspec.md` → Build
Order 6. ADR-0021, ADR-0022. Repo hard rule (canonical skill ships with CLI
behavior changes).
