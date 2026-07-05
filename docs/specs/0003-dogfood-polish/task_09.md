---
task: task_09
spec: 0003-dogfood-polish
status: pending
type: docs
complexity: low
---

# Task 09: Sync docs and the Roundfix skill with the polish

## Overview

Close the loop on the shipped behavior changes: the canonical Roundfix skill
documents the new commit subjects, the stop selector, and the Interactive
Input QA field; the repo agent instructions gain the commit-style line the
dogfood kept tripping on. Verifiable through the skills drift check inside
the full gate.

## Requirements

1. MUST update the canonical Roundfix skill wherever it documents the Task
   commit subject (now lowercase), the QA Report commit subject (now
   unscoped), stop selectors (now including `--spec`), and implement
   Interactive Input (now including the QA field); regenerate the embedded
   copy through the sync target.
2. MUST add one line to the repo agent instructions stating that commit and
   PR titles are unscoped Conventional Commits subjects in this repository
   (`cog.toml` `scopes = []`).
3. MUST amend the 0001 spec's commit-contract line via a dated note in this
   Spec's docs (never editing 0001's shipped techspec retroactively — the
   change is recorded here and in the skill).
4. MUST verify every term against the glossary; call out gaps instead of
   inventing language.

## Subtasks

- [ ] Canonical skill updates + `make skills-sync`
- [ ] Agent-instructions commit-style line
- [ ] Glossary pass

## Acceptance Criteria

- [ ] Skill text matches shipped behavior exactly (subjects, selectors,
      fields); the drift check passes inside the full gate.
- [ ] The agent instructions carry the unscoped-subjects line.
- [ ] No new un-glossaried term.

## Verification

- `rtk go run ./cmd/roundfix skills check` — expected: skill artifacts
  validate.
- `make verify` — expected: full gate passes, including the drift check.

## References

`_prd.md` → Core Feature 9; Decisions. `_techspec.md` → Build Order 9.
Dogfood finding 14. Repo hard rule (canonical skill ships with CLI behavior
changes).
