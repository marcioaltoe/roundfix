---
task: task_09
spec: 0003-dogfood-polish
status: completed
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

- [x] Canonical skill updates + `make skills-sync`
- [x] Agent-instructions commit-style line
- [x] Glossary pass

## Acceptance Criteria

- [x] Skill text matches shipped behavior exactly (subjects, selectors,
      fields); the drift check passes inside the full gate.
- [x] The agent instructions carry the unscoped-subjects line.
- [x] No new un-glossaried term.

## Verification

- `rtk go run ./cmd/roundfix skills check` — expected: skill artifacts
  validate.
- `make verify` — expected: full gate passes, including the drift check.

## References

`_prd.md` → Core Feature 9; Decisions. `_techspec.md` → Build Order 9.
Dogfood finding 14. Repo hard rule (canonical skill ships with CLI behavior
changes).

## Result

- Updated the canonical Roundfix Skill to match shipped behavior: Task commits
  use the existing type mapping with the first rune of the derived subject
  lowercased, QA Report commits use `docs: qa report for <slug> (<verdict>)`,
  `roundfix stop --spec <slug>` is documented, and implement Interactive Input
  documents the `QA gate [y/N]` / `QA gate [Y/n]` field. Ran
  `rtk make skills-sync`; `diff -r .agents/skills/roundfix skills/roundfix`
  had no output afterward.
- Added the repo agent-instructions line: commit and PR titles are unscoped
  Conventional Commits subjects because `cog.toml` sets `scopes = []`.
- Added a dated 2026-07-05 amendment note in this Spec's techspec instead of
  editing the shipped 0001 techspec retroactively.
- Glossary pass: reused existing product terms from `CONTEXT.md` (`Roundfix
  Skill`, `Implement Command`, `Interactive Input`, `QA Report`, `Task`,
  `Spec`, `Active Run`, `Stop Request`, and `Daemon`). No new Roundfix product
  term was introduced; `Conventional Commits` and `scopes = []` are the
  external commit-format/config terms required by this task.
- Pre-change signal: `rg` found stale skill text in both skill copies:
  `<type>: <title>` and `docs(qa): qa report for <slug> (<verdict>)`.
  Post-change, the same `rg` pattern over `.agents/skills/roundfix`,
  `skills/roundfix`, `AGENTS.md`, and this Spec's techspec returned no
  matches.
- Verification:
  - `rtk go run ./cmd/roundfix skills check`: passed (`Roundfix skill check
    passed: roundfix`).
  - `rtk make verify`: passed (`go test`, `skills check`, and `go build`
    exited 0).
