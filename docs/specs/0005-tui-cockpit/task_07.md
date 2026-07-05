---
task: task_07
spec: 0005-tui-cockpit
status: pending
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

- [ ] Attach parity tests for both Run kinds
- [ ] Skill Live Run View wording + `make skills-sync`
- [ ] Design-contract deviation check
- [ ] Glossary pass

## Acceptance Criteria

- [ ] Attach tests replay both Run kinds through the new cockpit with the
      modal working and Run state untouched.
- [ ] Skill text matches shipped keys and surfaces; drift check passes
      inside the full gate.
- [ ] Design docs match shipped behavior (or carry the dated correction).
- [ ] No new un-glossaried term.

## Verification

- `rtk go run ./cmd/roundfix skills check` — expected: skill artifacts
  validate.
- `make verify` — expected: full gate passes, including the drift check.

## References

`_prd.md` → User Story 4; Core Feature 6; Success Metrics. `_techspec.md` →
Build Order 7. ADR-0009. Repo hard rule (canonical skill ships with behavior
changes).
