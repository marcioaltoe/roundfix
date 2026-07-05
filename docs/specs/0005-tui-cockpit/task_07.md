---
task: task_07
spec: 0005-tui-cockpit
status: completed
type: docs
complexity: low
---

# Task 07: Attach parity pass and docs sync

## Overview

Prove the redesigned cockpit end to end where users will actually meet it —
Attach — and sync the docs: attach renders the same surfaces for both Run
kinds from Journal replay, and the canonical Roundfix skill's Live Run View
wording matches the shipped cockpit. Verifiable through attach tests and the
skills drift check.

## Requirements

1. MUST add attach tests over a finished review Run and a finished spec Run:
   the replayed cockpit renders the two surfaces, phase row, grouped
   timeline, and working modal, read-only, with attach-mode key differences
   preserved (ADR-0009 discipline).
2. MUST update the canonical Roundfix skill wherever it describes the Live
   Run View (surfaces, detail modal, keys) and regenerate the embedded copy
   through the sync target.
3. MUST leave the `design/` plan and mockups as-is unless the shipped
   behavior deviated — any deviation is corrected in the design doc in this
   same task (PRD decision: the design folder is the binding contract).
4. MUST verify every term against the glossary; call out gaps instead of
   inventing language.

## Subtasks

- [x] Attach parity tests for both Run kinds
- [x] Skill Live Run View wording + `make skills-sync`
- [x] Design-contract deviation check
- [x] Glossary pass

## Acceptance Criteria

- [x] Attach tests replay both Run kinds through the new cockpit with the
      modal working and Run state untouched.
- [x] Skill text matches shipped keys and surfaces; drift check passes
      inside the full gate.
- [x] Design docs match shipped behavior (or carry the dated correction).
- [x] No new un-glossaried term.

## Verification

- `rtk go run ./cmd/roundfix skills check` — expected: skill artifacts
  validate.
- `make verify` — expected: full gate passes, including the drift check.

## References

`_prd.md` → User Story 4; Core Feature 6; Success Metrics. `_techspec.md` →
Build Order 7. ADR-0009. Repo hard rule (canonical skill ships with behavior
changes).

## Result

- Attach parity evidence: `TestCockpitAttachReplaysFinishedReviewRunThroughRedesignedCockpit`
  and `TestCockpitAttachReplaysFinishedSpecRunThroughRedesignedCockpit` replay
  finished review and spec Runs through the cockpit in Attach mode. The tests
  assert the `WORK QUEUE` and `SESSION.TIMELINE` surfaces, phase rows, grouped
  Batch timeline, working Detail Modal, `q detach` Attach footer, and unchanged
  Clean Run state.
- Skill sync evidence: `.agents/skills/roundfix/SKILL.md` now documents the
  shipped Live Run View surfaces, Phase Row, Detail Modal, footer keys, Attach
  detach behavior, owning terminal close behavior, and small-terminal collapse.
  `rtk make skills-sync` regenerated `skills/roundfix`, and
  `rtk diff -r .agents/skills/roundfix skills/roundfix` returned no drift.
- Design-contract evidence: `design/ui-redesign-plan.md` carries a dated
  2026-07-05 correction for the shipped `WORK QUEUE` label, normal/modal
  footer wording, Attach/terminal mode keys, and small-terminal fallback.
- Glossary evidence: `CONTEXT.md` now defines Cockpit, Work Queue, Phase Row,
  Session Timeline, and Detail Modal, so the new skill and design wording use
  glossary-backed terms.
- Verification: `rtk go run ./cmd/roundfix skills check` passed.
- Verification: `rtk make verify` passed with the full Go suite, skill drift
  check, skill validation, and build.
